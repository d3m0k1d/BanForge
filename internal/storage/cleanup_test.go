package storage

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestRequestWriter_CleanupRequests(t *testing.T) {
	writer := newCleanupTestWriter(t)
	now := time.Now()

	insertCleanupTestRequest(t, writer, "192.0.2.10", now.Add(-48*time.Hour))
	insertCleanupTestRequest(t, writer, "192.0.2.11", now.Add(-time.Hour))

	if err := writer.CleanupRequests(24 * time.Hour); err != nil {
		t.Fatalf("CleanupRequests() error = %v", err)
	}

	rows, err := writer.db.Query("SELECT ip FROM requests ORDER BY ip")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var ips []string
	for rows.Next() {
		var ip string
		if err := rows.Scan(&ip); err != nil {
			t.Fatal(err)
		}
		ips = append(ips, ip)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}

	if len(ips) != 1 || ips[0] != "192.0.2.11" {
		t.Fatalf("remaining IPs = %v, want [192.0.2.11]", ips)
	}
}

func TestRequestWriter_CleanupRequestsRejectsInvalidRetention(t *testing.T) {
	writer := newCleanupTestWriter(t)

	tests := []struct {
		name      string
		retention time.Duration
	}{
		{name: "zero", retention: 0},
		{name: "negative", retention: -time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writer.CleanupRequests(tt.retention); err == nil {
				t.Fatalf("CleanupRequests(%v) error = nil", tt.retention)
			}
		})
	}
}

func TestRequestWriter_CleanupRequestsReturnsDatabaseError(t *testing.T) {
	writer := newCleanupTestWriter(t)
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	err := writer.CleanupRequests(time.Hour)
	if err == nil {
		t.Fatal("CleanupRequests() error = nil for closed database")
	}
	if !strings.Contains(err.Error(), "failed to cleanup requests") {
		t.Errorf("CleanupRequests() error = %q", err)
	}
}

func TestRequestCleanerRunsImmediatelyAndStopsOnCancellation(t *testing.T) {
	writer := newCleanupTestWriter(t)
	insertCleanupTestRequest(t, writer, "192.0.2.10", time.Now().Add(-48*time.Hour))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := RequestCleaner(ctx, 24*time.Hour, time.Hour, writer); err != nil {
		t.Fatalf("RequestCleaner() error = %v", err)
	}

	count, err := writer.GetRequestCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("request count = %d, want 0", count)
	}
}

func TestRequestCleanerRejectsInvalidArguments(t *testing.T) {
	writer := newCleanupTestWriter(t)

	tests := []struct {
		name      string
		retention time.Duration
		interval  time.Duration
		writer    *RequestWriter
	}{
		{name: "zero retention", retention: 0, interval: time.Hour, writer: writer},
		{name: "zero interval", retention: time.Hour, interval: 0, writer: writer},
		{name: "nil writer", retention: time.Hour, interval: time.Hour, writer: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RequestCleaner(context.Background(), tt.retention, tt.interval, tt.writer)
			if err == nil {
				t.Fatal("RequestCleaner() error = nil")
			}
		})
	}
}

func TestRequestCleanerContinuesAfterInitialFailure(t *testing.T) {
	writer := newCleanupTestWriter(t)
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// A failed initial cleanup must be logged, not fatal: the cleaner keeps
	// running and exits normally on context cancellation.
	if err := RequestCleaner(ctx, 24*time.Hour, time.Hour, writer); err != nil {
		t.Fatalf("RequestCleaner() error = %v, want nil", err)
	}
}

func newCleanupTestWriter(t *testing.T) *RequestWriter {
	t.Helper()

	dbPath := t.TempDir() + "/requests_test.db"
	writer, err := NewRequestWriterWithDBPath(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := writer.CreateTable(); err != nil {
		_ = writer.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = writer.Close()
	})

	return writer
}

func insertCleanupTestRequest(
	t *testing.T,
	writer *RequestWriter,
	ip string,
	createdAt time.Time,
) {
	t.Helper()

	_, err := writer.db.Exec(
		`INSERT INTO requests (service, ip, path, method, status, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		"test", ip, "/", "GET", "401", createdAt.Format(time.RFC3339),
	)
	if err != nil {
		t.Fatal(err)
	}
}
