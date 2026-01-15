package storage

import (
	"testing"
	"time"
)

func TestWrite(t *testing.T) {
	var ip string
	d := createTestDBStruct(t)

	err := d.CreateTable()
	if err != nil {
		t.Fatal(err)
	}

	resultCh := make(chan *LogEntry)

	go Write(d, resultCh)

	resultCh <- &LogEntry{
		Service: "test",
		IP:      "127.0.0.1",
		Path:    "/test",
		Method:  "GET",
		Status:  "200",
	}

	close(resultCh)

	time.Sleep(100 * time.Millisecond)

	err = d.db.QueryRow("SELECT ip FROM requests LIMIT 1").Scan(&ip)
	if err != nil {
		t.Fatal(err)
	}
	if ip != "127.0.0.1" {
		t.Fatal("ip should be 127.0.0.1")
	}
}
