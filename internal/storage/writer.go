package storage

import (
	"time"
)

func Write(db *DB, resultCh <-chan *LogEntry) {
	for result := range resultCh {
		_, err := db.db.Exec(
			"INSERT INTO requests (service, ip, path, method, status, created_at) VALUES (?, ?, ?, ?, ?, ?)",
			result.Service,
			result.IP,
			result.Path,
			result.Method,
			result.Status,
			time.Now().Format(time.RFC3339),
		)
		if err != nil {
			db.logger.Error("Failed to write to database", "error", err)
		}
	}
}
