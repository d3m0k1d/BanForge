package storage

import (
	"database/sql"
	"github.com/d3m0k1d/BanForge/internal/logger"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func createTestDB(t *testing.T) *sql.DB {
	tmpDir, err := os.MkdirTemp("", "banforge-test-*")
	if err != nil {
		t.Fatal(err)
	}

	filePath := filepath.Join(tmpDir, "test.db")
	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
		os.RemoveAll(tmpDir)
	})

	return db
}

func createTestDBStruct(t *testing.T) *DB {
	tmpDir, err := os.MkdirTemp("", "banforge-test-*")
	if err != nil {
		t.Fatal(err)
	}

	filePath := filepath.Join(tmpDir, "test.db")
	sqlDB, err := sql.Open("sqlite", filePath)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		sqlDB.Close()
		os.RemoveAll(tmpDir)
	})

	return &DB{
		logger: logger.New(false),
		db:     sqlDB,
	}
}

func TestCreateTable(t *testing.T) {
	d := createTestDBStruct(t)

	err := d.CreateTable()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := d.db.Query("SELECT 1 FROM requests LIMIT 1")
	if err != nil {
		t.Fatal("requests table should exist:", err)
	}
	rows.Close()

	rows, err = d.db.Query("SELECT 1 FROM bans LIMIT 1")
	if err != nil {
		t.Fatal("bans table should exist:", err)
	}
	rows.Close()
}

func TestIsBanned(t *testing.T) {
	d := createTestDBStruct(t)

	err := d.CreateTable()
	if err != nil {
		t.Fatal(err)
	}

	_, err = d.db.Exec("INSERT INTO bans (ip, banned_at) VALUES (?, ?)", "127.0.0.1", time.Now().Format(time.RFC3339))
	if err != nil {
		t.Fatal(err)
	}

	isBanned, err := d.IsBanned("127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	if !isBanned {
		t.Fatal("should be banned")
	}
}

func TestAddBan(t *testing.T) {
	d := createTestDBStruct(t)

	err := d.CreateTable()
	if err != nil {
		t.Fatal(err)
	}

	err = d.AddBan("127.0.0.1", "7h")
	if err != nil {
		t.Fatal(err)
	}

	var ip string
	err = d.db.QueryRow("SELECT ip FROM bans WHERE ip = ?", "127.0.0.1").Scan(&ip)
	if err != nil {
		t.Fatal(err)
	}

	if ip != "127.0.0.1" {
		t.Fatal("ip should be 127.0.0.1")
	}
}

func TestBanList(t *testing.T) {
	d := createTestDBStruct(t)

	err := d.CreateTable()
	if err != nil {
		t.Fatal(err)
	}

	_, err = d.db.Exec("INSERT INTO bans (ip, banned_at) VALUES (?, ?)", "127.0.0.1", time.Now().Format(time.RFC3339))
	if err != nil {
		t.Fatal(err)
	}

	err = d.BanList()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClose(t *testing.T) {
	d := createTestDBStruct(t)

	err := d.Close()
	if err != nil {
		t.Fatal(err)
	}
}
