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

type DB struct {
	logger *logger.Logger
	db     *sql.DB
}

func NewDB() (*DB, error) {
	db, err := sql.Open(
		"sqlite",
		"/var/lib/banforge/storage.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(30000)&_pragma=synchronous(NORMAL)",
	)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DB{
		logger: logger.New(false),
		db:     db,
	}, nil
}

func (d *DB) Close() error {
	d.logger.Info("Closing database connection")
	err := d.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) CreateTable() error {
	_, err := d.db.Exec(CreateTables)
	if err != nil {
		return err
	}
	d.logger.Info("Created tables")
	return nil
}

func (d *DB) IsBanned(ip string) (bool, error) {
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

func (d *DB) AddBan(ip string, ttl string) error {
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

func (d *DB) RemoveBan(ip string) error {
	_, err := d.db.Exec("DELETE FROM bans WHERE ip = ?", ip)
	if err != nil {
		d.logger.Error("Failed to remove ban", "error", err)
		return err
	}
	return nil
}

func (d *DB) BanList() error {

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

func (d *DB) CheckExpiredBans() ([]string, error) {
	var ips []string
	rows, err := d.db.Query(
		"SELECT ip FROM bans WHERE expired_at < ?",
		time.Now().Format(time.RFC3339),
	)
	if err != nil {
		d.logger.Error("Failed to get ban list", "error", err)
		return nil, err
	}
	for rows.Next() {
		var ip string
		r, err := d.db.Exec("DELETE FROM bans WHERE ip = ?", ip)
		if err != nil {
			d.logger.Error("Failed to get ban list", "error", err)
			return nil, err
		}
		d.logger.Info("Ban removed", "ip", ip, "rows", r)
		err = rows.Scan(&ip)
		if err != nil {
			d.logger.Error("Failed to get ban list", "error", err)
			return nil, err
		}
		ips = append(ips, ip)
	}
	return ips, nil
}
