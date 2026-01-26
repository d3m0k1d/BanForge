package storage

import (
	"database/sql"

	"github.com/d3m0k1d/BanForge/internal/logger"
	"path/filepath"
	"testing"
)

func TestBanWriter_AddBan(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "bans_test.db")

	writer, err := NewBanWriterWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create BanWriter: %v", err)
	}
	defer writer.Close()

	err = writer.CreateTable()
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	ip := "192.168.1.1"
	ttl := "1h"

	err = writer.AddBan(ip, ttl, "test")
	if err != nil {
		t.Errorf("AddBan failed: %v", err)
	}

	reader, err := NewBanReaderWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create BanReader: %v", err)
	}
	defer reader.Close()

	isBanned, err := reader.IsBanned(ip)
	if err != nil {
		t.Errorf("IsBanned failed: %v", err)
	}
	if !isBanned {
		t.Error("Expected IP to be banned, but it's not")
	}
}

func TestBanWriter_RemoveBan(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "bans_test.db")

	writer, err := NewBanWriterWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create BanWriter: %v", err)
	}
	defer writer.Close()

	err = writer.CreateTable()
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	ip := "192.168.1.2"
	err = writer.AddBan(ip, "1h", "test")
	if err != nil {
		t.Fatalf("Failed to add ban: %v", err)
	}

	reader, err := NewBanReaderWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create BanReader: %v", err)
	}
	defer reader.Close()

	isBanned, err := reader.IsBanned(ip)
	if err != nil {
		t.Fatalf("IsBanned failed: %v", err)
	}
	if !isBanned {
		t.Fatal("Expected IP to be banned before removal")
	}

	err = writer.RemoveBan(ip)
	if err != nil {
		t.Errorf("RemoveBan failed: %v", err)
	}

	isBanned, err = reader.IsBanned(ip)
	if err != nil {
		t.Errorf("IsBanned failed after removal: %v", err)
	}
	if isBanned {
		t.Error("Expected IP to be unbanned after removal, but it's still banned")
	}
}

func TestBanWriter_RemoveExpiredBans(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "bans_test.db")

	writer, err := NewBanWriterWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create BanWriter: %v", err)
	}
	defer writer.Close()

	err = writer.CreateTable()
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	expiredIP := "192.168.1.3"
	err = writer.AddBan(expiredIP, "-1h", "tes")
	if err != nil {
		t.Fatalf("Failed to add expired ban: %v", err)
	}

	activeIP := "192.168.1.4"
	err = writer.AddBan(activeIP, "1h", "test")
	if err != nil {
		t.Fatalf("Failed to add active ban: %v", err)
	}

	removedIPs, err := writer.RemoveExpiredBans()
	if err != nil {
		t.Errorf("RemoveExpiredBans failed: %v", err)
	}

	found := false
	for _, ip := range removedIPs {
		if ip == expiredIP {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected expired IP to be in removed list")
	}

	if len(removedIPs) != 1 {
		t.Errorf("Expected 1 removed IP, got %d", len(removedIPs))
	}

	reader, err := NewBanReaderWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create BanReader: %v", err)
	}
	defer reader.Close()

	isExpiredBanned, err := reader.IsBanned(expiredIP)
	if err != nil {
		t.Errorf("IsBanned failed for expired IP: %v", err)
	}
	if isExpiredBanned {
		t.Error("Expected expired IP to be unbanned, but it's still banned")
	}

	isActiveBanned, err := reader.IsBanned(activeIP)
	if err != nil {
		t.Errorf("IsBanned failed for active IP: %v", err)
	}
	if !isActiveBanned {
		t.Error("Expected active IP to still be banned, but it's not")
	}
}

func TestBanReader_IsBanned(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "bans_test.db")

	writer, err := NewBanWriterWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create BanWriter: %v", err)
	}
	defer writer.Close()

	err = writer.CreateTable()
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	ip := "192.168.1.5"
	err = writer.AddBan(ip, "1h", "test")
	if err != nil {
		t.Fatalf("Failed to add ban: %v", err)
	}

	reader, err := NewBanReaderWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create BanReader: %v", err)
	}
	defer reader.Close()

	isBanned, err := reader.IsBanned(ip)
	if err != nil {
		t.Errorf("IsBanned failed for banned IP: %v", err)
	}
	if !isBanned {
		t.Error("Expected IP to be banned")
	}

	isBanned, err = reader.IsBanned("192.168.1.6")
	if err != nil {
		t.Errorf("IsBanned failed for non-banned IP: %v", err)
	}
	if isBanned {
		t.Error("Expected IP to not be banned")
	}
}

func TestBanWriter_Close(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "bans_test.db")

	writer, err := NewBanWriterWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create BanWriter: %v", err)
	}

	err = writer.CreateTable()
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	err = writer.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	_, err = writer.db.Exec("SELECT 1")
	if err == nil {
		t.Error("Expected error when using closed connection")
	}
}

func TestBanReader_Close(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "bans_test.db")

	writer, err := NewBanWriterWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create BanWriter: %v", err)
	}
	defer writer.Close()

	err = writer.CreateTable()
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	reader, err := NewBanReaderWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create BanReader: %v", err)
	}

	err = reader.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	_, err = reader.db.Query("SELECT 1")
	if err == nil {
		t.Error("Expected error when using closed connection")
	}
}

func TestBanWriter_AddBan_InvalidDuration(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "bans_test.db")

	writer, err := NewBanWriterWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create BanWriter: %v", err)
	}
	defer writer.Close()

	err = writer.CreateTable()
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	err = writer.AddBan("192.168.1.7", "invalid_duration", "test")
	if err == nil {
		t.Error("Expected error for invalid duration")
	} else if err.Error() == "" || err.Error() == "<nil>" {
		t.Error("Expected meaningful error message for invalid duration")
	}
}

func TestMultipleBans(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "bans_test.db")

	writer, err := NewBanWriterWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create BanWriter: %v", err)
	}
	defer writer.Close()

	err = writer.CreateTable()
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	ips := []string{"192.168.1.8", "192.168.1.9", "192.168.1.10"}

	for _, ip := range ips {
		err := writer.AddBan(ip, "1h", "test")
		if err != nil {
			t.Errorf("Failed to add ban for IP %s: %v", ip, err)
		}
	}

	reader, err := NewBanReaderWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create BanReader: %v", err)
	}
	defer reader.Close()

	for _, ip := range ips {
		isBanned, err := reader.IsBanned(ip)
		if err != nil {
			t.Errorf("IsBanned failed for IP %s: %v", ip, err)
			continue
		}
		if !isBanned {
			t.Errorf("Expected IP %s to be banned", ip)
		}
	}
}

func TestRemoveNonExistentBan(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "bans_test.db")

	writer, err := NewBanWriterWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create BanWriter: %v", err)
	}
	defer writer.Close()

	err = writer.CreateTable()
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	err = writer.RemoveBan("192.168.1.11")
	if err != nil {
		t.Errorf("RemoveBan should not return error for non-existent ban: %v", err)
	}
}
func NewBanWriterWithDBPath(dbPath string) (*BanWriter, error) {
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(30000)&_pragma=synchronous(NORMAL)")
	if err != nil {
		return nil, err
	}
	return &BanWriter{
		logger: logger.New(false),
		db:     db,
	}, nil
}

func NewBanReaderWithDBPath(dbPath string) (*BanReader, error) {
	db, err := sql.Open("sqlite",
		dbPath+"?"+
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
