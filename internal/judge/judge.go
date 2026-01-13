package Judge

import (
	"fmt"

	"github.com/d3m0k1d/BanForge/internal/blocker"
	"github.com/d3m0k1d/BanForge/internal/config"
	"github.com/d3m0k1d/BanForge/internal/logger"
	"github.com/d3m0k1d/BanForge/internal/storage"
)

type Judge struct {
	db             *storage.DB
	logger         *logger.Logger
	Blocker        blocker.BlockerEngine
	rulesByService map[string][]config.Rule
}

func New(db *storage.DB) *Judge {
	return &Judge{
		db:             db,
		logger:         logger.New(false),
		rulesByService: make(map[string][]config.Rule),
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
		err = rows.Scan(&entry.ID, &entry.Service, &entry.IP, &entry.Path, &entry.Status, &entry.Method, &entry.IsViewed, &entry.CreatedAt)
		if err != nil {
			j.logger.Error(fmt.Sprintf("Failed to scan database row: %v", err))
			continue
		}

		rules, serviceExists := j.rulesByService[entry.Service]
		if serviceExists {
			for _, rule := range rules {
				if (rule.Method == "" || entry.Method == rule.Method) &&
					(rule.Status == "" || entry.Status == rule.Status) &&
					(rule.Path == "" || entry.Path == rule.Path) {

					j.logger.Info(fmt.Sprintf("Rule matched for IP: %s, Service: %s", entry.IP, entry.Service))
					j.Blocker.Ban(entry.IP)
					j.logger.Info(fmt.Sprintf("IP banned: %s", entry.IP))
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
