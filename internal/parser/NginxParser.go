package parser

import (
	"regexp"

	"github.com/d3m0k1d/BanForge/internal/logger"
	"github.com/d3m0k1d/BanForge/internal/storage"
)

type NginxParser struct {
	pattern *regexp.Regexp
	logger  *logger.Logger
}

func NewNginxParser() *NginxParser {
	pattern := regexp.MustCompile(
		`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}).*\[(.*?)\]\s+"(\w+)\s+(.*?)\s+HTTP.*"\s+(\d+)`,
	)
	return &NginxParser{
		pattern: pattern,
		logger:  logger.New(false),
	}
}

func (p *NginxParser) Parse(eventCh <-chan Event, resultCh chan<- *storage.LogEntry) {
	// Group 1: IP, Group 2: Timestamp, Group 3: Method, Group 4: Path, Group 5: Status
	go func() {
		for event := range eventCh {
			matches := p.pattern.FindStringSubmatch(event.Data)
			if matches == nil {
				continue
			}
			path := matches[4]
			status := matches[5]
			method := matches[3]

			resultCh <- &storage.LogEntry{
				Service:  "nginx",
				IP:       matches[1],
				Path:     path,
				Status:   status,
				Method:   method,
				IsViewed: false,
			}
			p.logger.Info(
				"Parsed nginx log entry",
				"ip",
				matches[1],
				"path",
				path,
				"status",
				status,
				"method",
				method,
			)
		}
	}()
}
