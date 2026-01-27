package parser

import (
	"github.com/d3m0k1d/BanForge/internal/logger"
	"github.com/d3m0k1d/BanForge/internal/storage"
	"regexp"
)

type ApacheParser struct {
	pattern *regexp.Regexp
	logger  *logger.Logger
}

func NewApacheParser() *ApacheParser {
	pattern := regexp.MustCompile(
		`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\s+-\s+-\s+\[(.*?)\]\s+"(\w+)\s+(.*?)\s+HTTP/[\d.]+"\s+(\d+)\s+(\d+|-)\s+"(.*?)"\s+"(.*?)"`,
	)
	// Groups:
	// 1: IP
	// 2: Timestamp
	// 3: Method (GET, POST, etc.)
	// 4: Path
	// 5: Status Code (200, 404, 403...)
	// 6: Response Size
	// 7: Referer
	// 8: User-Agent

	return &ApacheParser{
		pattern: pattern,
		logger:  logger.New(false),
	}
}

func (p *ApacheParser) Parse(eventCh <-chan Event, resultCh chan<- *storage.LogEntry) {
	// Group 1: IP, Group 2: Timestamp, Group 3: Method, Group 4: Path, Group 5: Status
	for event := range eventCh {
		matches := p.pattern.FindStringSubmatch(event.Data)
		if matches == nil {
			continue
		}
		path := matches[4]
		status := matches[5]
		method := matches[3]

		resultCh <- &storage.LogEntry{
			Service: "apache",
			IP:      matches[1],
			Path:    path,
			Status:  status,
			Method:  method,
		}
		p.logger.Info(
			"Parsed apache log entry",
			"ip", matches[1],
			"path", path,
			"status", status,
			"method", method,
		)
	}
}
