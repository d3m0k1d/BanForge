package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/d3m0k1d/BanForge/internal/metrics"
)

func WriteReq(db *RequestWriter, resultCh <-chan *LogEntry) {
	db.logger.Info("Starting log writer")
	const batchSize = 100
	const flushInterval = 1 * time.Second

	batch := make([]*LogEntry, 0, batchSize)
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	flush := func() {
		defer db.logger.Debug("Flushed batch", "count", len(batch))
		err := func() (err error) {
			if len(batch) == 0 {
				return nil
			}

			tx, err := db.db.Begin()
			if err != nil {
				return fmt.Errorf("failed to begin transaction: %w", err)
			}
			defer func() {
				if rollbackErr := tx.Rollback(); rollbackErr != nil &&
					!errors.Is(rollbackErr, sql.ErrTxDone) {
					err = errors.Join(
						err,
						fmt.Errorf("failed to rollback transaction: %w", rollbackErr),
					)
				}
			}()

			stmt, err := tx.Prepare(
				"INSERT INTO requests (service, ip, path, method, status, created_at) VALUES (?, ?, ?, ?, ?, ?)",
			)
			if err != nil {
				err = fmt.Errorf("failed to prepare statement: %w", err)
				return err
			}
			defer func() {
				if closeErr := stmt.Close(); closeErr != nil {
					err = errors.Join(err, fmt.Errorf("failed to close statement: %w", closeErr))
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
					db.logger.Error(fmt.Errorf("failed to insert entry: %w", err).Error())
					metrics.IncError()
				} else {
					metrics.IncRequestCount(entry.Service)
				}
			}

			if err := tx.Commit(); err != nil {
				return fmt.Errorf("failed to commit transaction: %w", err)
			}

			batch = batch[:0]
			metrics.IncDBOperation("insert", "requests")
			return err
		}()
		if err != nil {
			db.logger.Error(err.Error())
		}
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

func (w *RequestWriter) GetRequestCount() (int, error) {
	var count int
	err := w.db.QueryRow("SELECT COUNT(*) FROM requests").Scan(&count)
	return count, err
}

func (w *RequestWriter) Close() error {
	return w.db.Close()
}
