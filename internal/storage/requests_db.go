package storage

import (
	"database/sql"
	"github.com/d3m0k1d/BanForge/internal/logger"
	_ "modernc.org/sqlite"
)

type Request_Writer struct {
	logger *logger.Logger
	db     *sql.DB
}

func NewRequestsWr() (*Request_Writer, error) {
	db, err := sql.Open("sqlite", "/var/lib/banforge/requests.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(30000)&_pragma=synchronous(NORMAL)")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)
	return &Request_Writer{
		logger: logger.New(false),
		db:     db,
	}, nil
}
