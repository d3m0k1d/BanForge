package judge

import (
	"fmt"
	"github.com/d3m0k1d/BanForge/internal/blocker"
	"github.com/d3m0k1d/BanForge/internal/config"
	"github.com/d3m0k1d/BanForge/internal/logger"
	"github.com/d3m0k1d/BanForge/internal/storage"
	"strings"
	"time"
)

type Judge struct {
	db             *storage.DB
	logger         *logger.Logger
	Blocker        blocker.BlockerEngine
	rulesByService map[string][]config.Rule
}

func New(db *storage.DB, b blocker.BlockerEngine) *Judge {
	return &Judge{
		db:             db,
		logger:         logger.New(false),
		rulesByService: make(map[string][]config.Rule),
		Blocker:        b,
	}
}

func (j *Judge) LoadRules(rules []config.Rule) {
	j.rulesByService = make(map[string][]config.Rule)
	for _, rule := range rules {
		j.rulesByService[rule.ServiceName] = append(
			j.rulesByService[rule.ServiceName],
			rule,
		)
	}
	j.logger.Info("Rules loaded and indexed by service")
}

func (j *Judge) ProcessUnviewed() error {
	rows, err := j.db.SearchUnViewed()
	if err != nil {
		j.logger.Error(fmt.Sprintf("Failed to query database: %v", err))
		return err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			j.logger.Error(fmt.Sprintf("Failed to close database connection: %v", err))
		}
	}()
	for rows.Next() {
		var entry storage.LogEntry
		err = rows.Scan(
			&entry.ID,
			&entry.Service,
			&entry.IP,
			&entry.Path,
			&entry.Status,
			&entry.Method,
			&entry.IsViewed,
			&entry.CreatedAt,
		)
		if err != nil {
			j.logger.Error(fmt.Sprintf("Failed to scan database row: %v", err))
			continue
		}

		rules, serviceExists := j.rulesByService[entry.Service]
		if serviceExists {
			for _, rule := range rules {
				if (rule.Method == "" || entry.Method == rule.Method) &&
					(rule.Status == "" || entry.Status == rule.Status) &&
					matchPath(entry.Path, rule.Path) {

					j.logger.Info(
						fmt.Sprintf(
							"Rule matched for IP: %s, Service: %s",
							entry.IP,
							entry.Service,
						),
					)
					ban_status, err := j.db.IsBanned(entry.IP)
					if err != nil {
						j.logger.Error(fmt.Sprintf("Failed to check ban status: %v", err))
						return err
					}
					if !ban_status {
						err = j.Blocker.Ban(entry.IP)
						if err != nil {
							j.logger.Error(fmt.Sprintf("Failed to ban IP: %v", err))
						}
						j.logger.Info(fmt.Sprintf("IP banned: %s", entry.IP))
						err = j.db.AddBan(entry.IP, rule.BanTime)
						if err != nil {
							j.logger.Error(fmt.Sprintf("Failed to add ban: %v", err))
						}
					}
					break
				}
			}
		}

		err = j.db.MarkAsViewed(entry.ID)
		if err != nil {
			j.logger.Error(fmt.Sprintf("Failed to mark entry as viewed: %v", err))
		} else {
			j.logger.Info(fmt.Sprintf("Entry marked as viewed: ID=%d", entry.ID))
		}
	}

	if err = rows.Err(); err != nil {
		j.logger.Error(fmt.Sprintf("Error iterating rows: %v", err))
		return err
	}

	return nil
}

func (j *Judge) UnbanChecker() {
	tick := time.NewTicker(5 * time.Minute)
	defer tick.Stop()

	for range tick.C {
		ips, err := j.db.CheckExpiredBans()
		if err != nil {
			j.logger.Error(fmt.Sprintf("Failed to check expired bans: %v", err))
			continue
		}

		for _, ip := range ips {
			err = j.db.RemoveBan(ip)
			if err != nil {
				j.logger.Error(fmt.Sprintf("Failed to remove ban: %v", err))
			}
			if err := j.Blocker.Unban(ip); err != nil {
				j.logger.Error(fmt.Sprintf("Failed to unban IP %s: %v", ip, err))
				continue
			}
			j.logger.Info(fmt.Sprintf("IP unbanned: %s", ip))
		}
	}
}

func matchPath(path string, rulePath string) bool {
	if rulePath == "" {
		return true
	}
	if strings.HasPrefix(rulePath, "/*") {
		prefix := strings.TrimPrefix(rulePath, "*")
		return strings.HasPrefix(path, prefix)
	}
	if strings.HasSuffix(rulePath, "*") {
		suffix := strings.TrimSuffix(rulePath, "*")
		return strings.HasSuffix(path, suffix)
	}
	return path == rulePath
}
