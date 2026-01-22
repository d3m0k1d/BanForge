package storage

import (
	"database/sql"
	"github.com/d3m0k1d/BanForge/internal/logger"
	_ "modernc.org/sqlite"
	"path/filepath"
	"testing"
	"time"
)

func TestWrite_BatchInsert(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "requests_test.db")

	writer, err := NewRequestWriterWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create RequestWriter: %v", err)
	}
	defer writer.Close()

	err = writer.CreateTable()
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	resultCh := make(chan *LogEntry, 100)

	done := make(chan bool)
	go func() {
		WriteReq(writer, resultCh)
		close(done)
	}()

	entries := []*LogEntry{
		{Service: "service1", IP: "192.168.1.1", Path: "/path1", Method: "GET", Status: "200"},
		{Service: "service2", IP: "192.168.1.2", Path: "/path2", Method: "POST", Status: "404"},
		{Service: "service3", IP: "192.168.1.3", Path: "/path3", Method: "PUT", Status: "500"},
		{Service: "service4", IP: "192.168.1.4", Path: "/path4", Method: "DELETE", Status: "200"},
		{Service: "service5", IP: "192.168.1.5", Path: "/path5", Method: "GET", Status: "301"},
	}

	for _, entry := range entries {
		resultCh <- entry
	}

	close(resultCh)
	<-done

	count, err := writer.GetRequestCount()
	if err != nil {
		t.Fatalf("Failed to get request count: %v", err)
	}

	if count != len(entries) {
		t.Errorf("Expected %d entries, got %d", len(entries), count)
	}
	rows, err := writer.db.Query("SELECT service, ip, path, method, status FROM requests ORDER BY id")
	if err != nil {
		t.Fatalf("Failed to query requests: %v", err)
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var service, ip, path, method, status string
		err := rows.Scan(&service, &ip, &path, &method, &status)
		if err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}

		if i >= len(entries) {
			t.Fatal("More rows returned than expected")
		}

		expected := entries[i]
		if service != expected.Service {
			t.Errorf("Expected service %s, got %s", expected.Service, service)
		}
		if ip != expected.IP {
			t.Errorf("Expected IP %s, got %s", expected.IP, ip)
		}
		if path != expected.Path {
			t.Errorf("Expected path %s, got %s", expected.Path, path)
		}
		if method != expected.Method {
			t.Errorf("Expected method %s, got %s", expected.Method, method)
		}
		if status != expected.Status {
			t.Errorf("Expected status %s, got %s", expected.Status, status)
		}

		i++
	}

	if i != len(entries) {
		t.Errorf("Expected to read %d entries, got %d", len(entries), i)
	}
}

func TestWrite_BatchSizeTrigger(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "requests_test.db")

	writer, err := NewRequestWriterWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create RequestWriter: %v", err)
	}
	defer writer.Close()

	err = writer.CreateTable()
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	resultCh := make(chan *LogEntry, 100)
	done := make(chan bool)
	go func() {
		WriteReq(writer, resultCh)
		close(done)
	}()

	batchSize := 100
	entries := make([]*LogEntry, batchSize)
	for i := 0; i < batchSize; i++ {
		entries[i] = &LogEntry{
			Service: "service" + string(rune(i+'0')),
			IP:      "192.168.1." + string(rune(i+'0')),
			Path:    "/path" + string(rune(i+'0')),
			Method:  "GET",
			Status:  "200",
		}
	}

	for _, entry := range entries {
		resultCh <- entry
	}

	close(resultCh)
	<-done

	count, err := writer.GetRequestCount()
	if err != nil {
		t.Fatalf("Failed to get request count: %v", err)
	}

	if count != batchSize {
		t.Errorf("Expected %d entries, got %d", batchSize, count)
	}
}

func TestWrite_FlushInterval(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "requests_test.db")

	writer, err := NewRequestWriterWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create RequestWriter: %v", err)
	}
	defer writer.Close()

	err = writer.CreateTable()
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	resultCh := make(chan *LogEntry, 100)

	done := make(chan bool)
	go func() {
		WriteReq(writer, resultCh)
		close(done)
	}()

	entries := []*LogEntry{
		{Service: "service1", IP: "192.168.1.1", Path: "/path1", Method: "GET", Status: "200"},
		{Service: "service2", IP: "192.168.1.2", Path: "/path2", Method: "POST", Status: "404"},
		{Service: "service3", IP: "192.168.1.3", Path: "/path3", Method: "PUT", Status: "500"},
		{Service: "service4", IP: "192.168.1.4", Path: "/path4", Method: "DELETE", Status: "200"},
		{Service: "service5", IP: "192.168.1.5", Path: "/path5", Method: "GET", Status: "301"},
	}

	for _, entry := range entries {
		resultCh <- entry
	}
	time.Sleep(1500 * time.Millisecond)

	close(resultCh)
	<-done

	count, err := writer.GetRequestCount()
	if err != nil {
		t.Fatalf("Failed to get request count: %v", err)
	}

	if count != len(entries) {
		t.Errorf("Expected %d entries, got %d", len(entries), count)
	}
}

func TestWrite_EmptyBatch(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "requests_test.db")

	writer, err := NewRequestWriterWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create RequestWriter: %v", err)
	}
	defer writer.Close()

	err = writer.CreateTable()
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	resultCh := make(chan *LogEntry, 100)

	done := make(chan bool)
	go func() {
		WriteReq(writer, resultCh)
		close(done)
	}()

	close(resultCh)
	<-done
	count, err := writer.GetRequestCount()
	if err != nil {
		t.Fatalf("Failed to get request count: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected 0 entries for empty batch, got %d", count)
	}
}

func TestWrite_ChannelClosed(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "requests_test.db")

	writer, err := NewRequestWriterWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create RequestWriter: %v", err)
	}
	defer writer.Close()

	err = writer.CreateTable()
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	resultCh := make(chan *LogEntry, 100)

	done := make(chan bool)
	go func() {
		WriteReq(writer, resultCh)
		close(done)
	}()

	entries := []*LogEntry{
		{Service: "service1", IP: "192.168.1.1", Path: "/path1", Method: "GET", Status: "200"},
		{Service: "service2", IP: "192.168.1.2", Path: "/path2", Method: "POST", Status: "404"},
	}

	for _, entry := range entries {
		resultCh <- entry
	}

	close(resultCh)

	<-done

	count, err := writer.GetRequestCount()
	if err != nil {
		t.Fatalf("Failed to get request count: %v", err)
	}

	if count != len(entries) {
		t.Errorf("Expected %d entries, got %d", len(entries), count)
	}
}

func NewRequestWriterWithDBPath(dbPath string) (*Request_Writer, error) {
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(30000)&_pragma=synchronous(NORMAL)")
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

func (w *Request_Writer) CreateTable() error {
	_, err := w.db.Exec(CreateRequestsTable)
	if err != nil {
		return err
	}
	w.logger.Info("Created requests table")
	return nil
}

func (w *Request_Writer) Close() error {
	w.logger.Info("Closing request database connection")
	err := w.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (w *Request_Writer) GetRequestCount() (int, error) {
	var count int
	err := w.db.QueryRow("SELECT COUNT(*) FROM requests").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
