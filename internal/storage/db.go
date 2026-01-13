package storage

import (
	"database/sql"

	"github.com/d3m0k1d/BanForge/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	logger *logger.Logger
	db     *sql.DB
}

func NewDB() (*DB, error) {
	db, err := sql.Open("sqlite3", "/var/lib/banforge/storage.db")
	if err != nil {
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

func (d *DB) SearchUnViewed() (*sql.Rows, error) {
	rows, err := d.db.Query("SELECT id, service, ip, path, status, method, viewed, created_at FROM requests WHERE viewed = 0")
	if err != nil {
		d.logger.Error("Failed to query database")
		return nil, err
	}
	return rows, nil
}

func (d *DB) MarkAsViewed(id int) error {
	_, err := d.db.Exec("UPDATE requests SET viewed = 1 WHERE id = ?", id)
	if err != nil {
		d.logger.Error("Failed to mark as viewed", "error", err)
		return err
	}
	return nil
}
