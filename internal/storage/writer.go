package storage

import (
	"time"
)

func Write(db *Request_Writer, resultCh <-chan *LogEntry) {
	db.logger.Info("Starting log writer")
	const batchSize = 100
	const flushInterval = 1 * time.Second

	batch := make([]*LogEntry, 0, batchSize)
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}

		tx, err := db.db.Begin()
		if err != nil {
			db.logger.Error("Failed to begin transaction", "error", err)
			return
		}

		stmt, err := tx.Prepare(
			"INSERT INTO requests (service, ip, path, method, status, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		)
		if err != nil {
			db.logger.Error("Failed to prepare statement", "error", err)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				db.logger.Error("Failed to rollback transaction", "error", rollbackErr)
			}
			return
		}
		defer func() {
			if closeErr := stmt.Close(); closeErr != nil {
				db.logger.Error("Failed to close statement", "error", closeErr)
			}
		}()

		for _, entry := range batch {
			_, err := stmt.Exec(
				entry.Service,
				entry.IP,
				entry.Path,
				entry.Method,
				entry.Status,
				time.Now().Format(time.RFC3339),
			)
			if err != nil {
				db.logger.Error("Failed to insert entry", "error", err)
			}
		}

		if err := tx.Commit(); err != nil {
			db.logger.Error("Failed to commit transaction", "error", err)
			return
		}

		db.logger.Debug("Flushed batch", "count", len(batch))
		batch = batch[:0]
	}

	for {
		select {
		case result, ok := <-resultCh:
			if !ok {
				flush()
				return
			}

			batch = append(batch, result)
			if len(batch) >= batchSize {
				flush()
			}

		case <-ticker.C:
			flush()
		}
	}
}
