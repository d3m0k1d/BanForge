package storage

import (
	"database/sql"
	"fmt"
	"github.com/d3m0k1d/BanForge/internal/config"
	"github.com/d3m0k1d/BanForge/internal/logger"
	"github.com/jedib0t/go-pretty/v6/table"
	_ "modernc.org/sqlite"
	"os"
	"time"
)

// Writer block
type BanWriter struct {
	logger *logger.Logger
	db     *sql.DB
}

func NewBanWriter() (*BanWriter, error) {
	db, err := sql.Open("sqlite", "/var/lib/banforge/bans.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(30000)&_pragma=synchronous(NORMAL)")
	if err != nil {
		return nil, err
	}
	return &BanWriter{
		logger: logger.New(false),
		db:     db,
	}, nil
}

func (d *BanWriter) CreateTable() error {
	_, err := d.db.Exec(CreateBansTable)
	if err != nil {
		return err
	}
	d.logger.Info("Created tables")
	return nil
}

func (d *BanWriter) AddBan(ip string, ttl string) error {
	duration, err := config.ParseDurationWithYears(ttl)
	if err != nil {
		d.logger.Error("Invalid duration format", "ttl", ttl, "error", err)
		return fmt.Errorf("invalid duration: %w", err)
	}

	now := time.Now()
	expiredAt := now.Add(duration)

	_, err = d.db.Exec(
		"INSERT INTO bans (ip, reason, banned_at, expired_at) VALUES (?, ?, ?, ?)",
		ip,
		"1",
		now.Format(time.RFC3339),
		expiredAt.Format(time.RFC3339),
	)
	if err != nil {
		d.logger.Error("Failed to add ban", "error", err)
		return err
	}

	return nil
}

func (d *BanWriter) RemoveBan(ip string) error {
	_, err := d.db.Exec("DELETE FROM bans WHERE ip = ?", ip)
	if err != nil {
		d.logger.Error("Failed to remove ban", "error", err)
		return err
	}
	return nil
}

func (w *BanWriter) RemoveExpiredBans() ([]string, error) {
	var ips []string
	now := time.Now().Format(time.RFC3339)

	rows, err := w.db.Query(
		"SELECT ip FROM bans WHERE expired_at < ?",
		now,
	)
	if err != nil {
		w.logger.Error("Failed to get expired bans", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ip string
		err := rows.Scan(&ip)
		if err != nil {
			w.logger.Error("Failed to scan ban", "error", err)
			continue
		}
		ips = append(ips, ip)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	result, err := w.db.Exec(
		"DELETE FROM bans WHERE expired_at < ?",
		now,
	)
	if err != nil {
		w.logger.Error("Failed to remove expired bans", "error", err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected > 0 {
		w.logger.Info("Removed expired bans", "count", rowsAffected, "ips", len(ips))
	}

	return ips, nil
}

func (d *BanWriter) Close() error {
	d.logger.Info("Closing database connection")
	err := d.db.Close()
	if err != nil {
		return err
	}
	return nil
}

// Reader block

type BanReader struct {
	logger *logger.Logger
	db     *sql.DB
}

func NewBanReader() (*BanReader, error) {
	db, err := sql.Open("sqlite",
		"/var/lib/banforge/bans.db?"+
			"mode=ro&"+
			"_pragma=journal_mode(WAL)&"+
			"_pragma=mmap_size(268435456)&"+
			"_pragma=cache_size(-2000)&"+
			"_pragma=query_only(1)")
	if err != nil {
		return nil, err
	}

	return &BanReader{
		logger: logger.New(false),
		db:     db,
	}, nil
}

func (d *BanReader) IsBanned(ip string) (bool, error) {
	var bannedIP string
	err := d.db.QueryRow("SELECT ip FROM bans WHERE ip = ? ", ip).Scan(&bannedIP)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check ban status: %w", err)
	}
	return true, nil
}

func (d *BanReader) BanList() error {

	var count int
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleBold)
	t.AppendHeader(table.Row{"№", "IP", "Banned At"})
	rows, err := d.db.Query("SELECT ip, banned_at  FROM bans")
	if err != nil {
		d.logger.Error("Failed to get ban list", "error", err)
		return err
	}
	for rows.Next() {
		count++
		var ip string
		var bannedAt string
		err := rows.Scan(&ip, &bannedAt)
		if err != nil {
			d.logger.Error("Failed to get ban list", "error", err)
			return err
		}
		t.AppendRow(table.Row{count, ip, bannedAt})

	}
	t.Render()
	return nil
}

func (d *BanReader) Close() error {
	d.logger.Info("Closing database connection")
	err := d.db.Close()
	if err != nil {
		return err
	}
	return nil
}
