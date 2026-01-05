package parser

import (
	"os"
	"testing"
)

func TestNewScanner(t *testing.T) {
	file, err := os.CreateTemp("", "test.log")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	defer os.Remove(file.Name())
	s, err := NewScanner(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("Scanner is nil")
	}
}
