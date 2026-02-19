package storage

import (
	"database/sql"

	"github.com/d3m0k1d/BanForge/internal/logger"
	_ "modernc.org/sqlite"
)

type RequestWriter struct {
	logger *logger.Logger
	db     *sql.DB
}

func NewRequestsWr() (*RequestWriter, error) {
	db, err := sql.Open(
		"sqlite",
		buildSqliteDsn(ReqDBPath, pragmas),
	)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)
	return &RequestWriter{
		logger: logger.New(false),
		db:     db,
	}, nil
}

type RequestReader struct {
	logger *logger.Logger
	db     *sql.DB
}

func NewRequestsRd() (*RequestReader, error) {
	db, err := sql.Open(
		"sqlite",
		buildSqliteDsn(ReqDBPath, pragmas),
	)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)
	return &RequestReader{
		logger: logger.New(false),
		db:     db,
	}, nil
}

func (r *RequestReader) IsMaxRetryExceeded(ip string, max_retry int) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM requests WHERE ip = ?", ip).Scan(&count)
	if err != nil {
		r.logger.Error("error query count: " + err.Error())
		return false, err
	}
	return count >= max_retry, nil
}
