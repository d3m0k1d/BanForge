package judge

import (
	"fmt"
	"strings"
	"time"

	"github.com/d3m0k1d/BanForge/internal/blocker"
	"github.com/d3m0k1d/BanForge/internal/config"
	"github.com/d3m0k1d/BanForge/internal/logger"
	"github.com/d3m0k1d/BanForge/internal/storage"
)

type Judge struct {
	db_r           *storage.BanReader
	db_w           *storage.BanWriter
	logger         *logger.Logger
	Blocker        blocker.BlockerEngine
	rulesByService map[string][]config.Rule
	entryCh        chan *storage.LogEntry
	resultCh       chan *storage.LogEntry
}

func New(
	db_r *storage.BanReader,
	db_w *storage.BanWriter,
	b blocker.BlockerEngine,
	resultCh chan *storage.LogEntry,
	entryCh chan *storage.LogEntry,
) *Judge {
	return &Judge{
		db_w:           db_w,
		db_r:           db_r,
		logger:         logger.New(false),
		rulesByService: make(map[string][]config.Rule),
		Blocker:        b,
		entryCh:        entryCh,
		resultCh:       resultCh,
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

func (j *Judge) Tribunal() {
	j.logger.Info("Tribunal started")

	for entry := range j.entryCh {
		j.logger.Debug(
			"Processing entry",
			"ip",
			entry.IP,
			"service",
			entry.Service,
			"status",
			entry.Status,
		)

		rules, serviceExists := j.rulesByService[entry.Service]
		if !serviceExists {
			j.logger.Debug("No rules for service", "service", entry.Service)
			continue
		}

		ruleMatched := false
		for _, rule := range rules {
			methodMatch := rule.Method == "" || entry.Method == rule.Method
			statusMatch := rule.Status == "" || entry.Status == rule.Status
			pathMatch := matchPath(entry.Path, rule.Path)

			j.logger.Debug(
				"Testing rule",
				"rule", rule.Name,
				"method_match", methodMatch,
				"status_match", statusMatch,
				"path_match", pathMatch,
			)

			if methodMatch && statusMatch && pathMatch {
				ruleMatched = true
				j.logger.Info("Rule matched", "rule", rule.Name, "ip", entry.IP)

				banned, err := j.db_r.IsBanned(entry.IP)
				if err != nil {
					j.logger.Error("Failed to check ban status", "ip", entry.IP, "error", err)
					break
				}

				if banned {
					j.logger.Info("IP already banned", "ip", entry.IP)
					j.resultCh <- entry
					break
				}

				err = j.db_w.AddBan(entry.IP, rule.BanTime, rule.Name)
				if err != nil {
					j.logger.Error(
						"Failed to add ban to database",
						"ip",
						entry.IP,
						"ban_time",
						rule.BanTime,
						"error",
						err,
					)
					break
				}

				if err := j.Blocker.Ban(entry.IP); err != nil {
					j.logger.Error("Failed to ban IP at firewall", "ip", entry.IP, "error", err)
					break
				}
				j.logger.Info(
					"IP banned successfully",
					"ip",
					entry.IP,
					"rule",
					rule.Name,
					"ban_time",
					rule.BanTime,
				)
				j.resultCh <- entry
				break
			}
		}

		if !ruleMatched {
			j.logger.Debug("No rules matched", "ip", entry.IP, "service", entry.Service)
		}
	}

	j.logger.Info("Tribunal stopped - entryCh closed")
}

func (j *Judge) UnbanChecker() {
	tick := time.NewTicker(5 * time.Minute)
	defer tick.Stop()

	for range tick.C {
		ips, err := j.db_w.RemoveExpiredBans()
		if err != nil {
			j.logger.Error(fmt.Sprintf("Failed to check expired bans: %v", err))
			continue
		}

		for _, ip := range ips {
			if err := j.Blocker.Unban(ip); err != nil {
				j.logger.Error(fmt.Sprintf("Failed to unban IP at firewall: %v", err))
			}
		}
	}
}

func matchPath(path string, rulePath string) bool {
	if rulePath == "" {
		return true
	}

	if strings.HasPrefix(rulePath, "*") {
		suffix := strings.TrimPrefix(rulePath, "*")
		return strings.HasSuffix(path, suffix)
	}

	if strings.HasPrefix(rulePath, "/*") {
		suffix := strings.TrimPrefix(rulePath, "/*")
		return strings.HasSuffix(path, suffix)
	}

	if strings.HasSuffix(rulePath, "*") {
		prefix := strings.TrimSuffix(rulePath, "*")
		return strings.HasPrefix(path, prefix)
	}
	return path == rulePath
}
