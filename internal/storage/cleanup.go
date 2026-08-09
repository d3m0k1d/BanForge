package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/d3m0k1d/BanForge/internal/metrics"
)

func (w *RequestWriter) CleanupRequests(retention time.Duration) error {
	if retention <= 0 {
		return fmt.Errorf("retention must be greater than 0")
	}
	result, err := w.db.Exec(
		"DELETE FROM requests WHERE created_at < ?",
		time.Now().Add(-retention).Format(time.RFC3339),
	)
	if err != nil {
		w.logger.Error("Failed to cleanup requests", "error", err)
		metrics.IncError()
		return fmt.Errorf("failed to cleanup requests: %w", err)
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		w.logger.Error("Failed to get deleted requests count", "error", err)
		metrics.IncError()
		return fmt.Errorf("failed to get deleted requests count: %w", err)
	}

	if deleted > 0 {
		w.logger.Info("Request cleanup completed", "deleted", deleted)
		metrics.IncDBOperation("delete", "requests")
	}

	return nil
}
func RequestCleaner(
	ctx context.Context,
	retention time.Duration,
	interval time.Duration,
	w *RequestWriter,
) error {
	if retention <= 0 {
		return fmt.Errorf("retention must be greater than zero")
	}
	if interval <= 0 {
		return fmt.Errorf("interval must be greater than zero")
	}
	if w == nil {
		return fmt.Errorf("request writer is nil")
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	if err := w.CleanupRequests(retention); err != nil {
		w.logger.Error("request cleanup failed", "error", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			if err := w.CleanupRequests(retention); err != nil {
				w.logger.Error("request cleanup failed", "error", err)
			}
		}
	}
}
