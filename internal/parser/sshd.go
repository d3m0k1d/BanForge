package parser

import (
	"regexp"

	"github.com/d3m0k1d/BanForge/internal/logger"
	"github.com/d3m0k1d/BanForge/internal/storage"
)

type SshdParser struct {
	pattern *regexp.Regexp
	logger  *logger.Logger
}

func NewSshdParser() *SshdParser {
	pattern := regexp.MustCompile(
		`^([A-Za-z]{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2})\s+(\S+)\s+sshd(?:-session)?\[(\d+)\]:\s+Failed\s+(\w+)\s+for\s+(?:invalid\s+user\s+)?(\S+)\s+from\s+(\S+)\s+port\s+(\d+)`,
	)
	return &SshdParser{
		pattern: pattern,
		logger:  logger.New(false),
	}
}

func (p *SshdParser) Parse(eventCh <-chan Event, resultCh chan<- *storage.LogEntry) {
	// Group 1: Timestamp, Group 2: hostame, Group 3: pid, Group 4: Method auth, Group 5: User, Group 6: IP, Group 7: port
	for event := range eventCh {
		matches := p.pattern.FindStringSubmatch(event.Data)
		if matches == nil {
			continue
		}
		resultCh <- &storage.LogEntry{
			Service: "ssh",
			IP:      matches[6],
			Path:    matches[5], // user
			Status:  "Failed",
			Method:  matches[4], // method auth
		}
		p.logger.Info(
			"Parsed ssh log entry",
			"ip",
			matches[6],
			"user",
			matches[5],
			"method",
			matches[4],
			"status",
			"Failed",
		)
	}
}
