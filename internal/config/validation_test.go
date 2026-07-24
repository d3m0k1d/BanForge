package config

import (
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestStorageTOMLFields(t *testing.T) {
	var cfg Config
	if _, err := toml.Decode(`
[storage]
retention_time = "2d"
cleanup_interval = "1h"
`, &cfg); err != nil {
		t.Fatal(err)
	}

	if cfg.Storage.RetentionTime != "2d" {
		t.Errorf("RetentionTime = %q, want %q", cfg.Storage.RetentionTime, "2d")
	}
	if cfg.Storage.CleanupInterval != "1h" {
		t.Errorf("CleanupInterval = %q, want %q", cfg.Storage.CleanupInterval, "1h")
	}
}

func TestStorageDefaults(t *testing.T) {
	cfg := newConfigWithDefaults()
	if _, err := toml.Decode("[metrics]\nenabled = false\n", cfg); err != nil {
		t.Fatal(err)
	}

	if cfg.Storage.RetentionTime != "2d" {
		t.Errorf("RetentionTime = %q, want default %q", cfg.Storage.RetentionTime, "2d")
	}
	if cfg.Storage.CleanupInterval != "1h" {
		t.Errorf("CleanupInterval = %q, want default %q", cfg.Storage.CleanupInterval, "1h")
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config validation failed: %v", err)
	}
}

func TestStorageValidate(t *testing.T) {
	tests := []struct {
		name    string
		storage Storage
		wantErr string
	}{
		{
			name: "valid",
			storage: Storage{
				RetentionTime:   "2d",
				CleanupInterval: "1h",
			},
		},
		{
			name: "years supported",
			storage: Storage{
				RetentionTime:   "1y",
				CleanupInterval: "1d",
			},
		},
		{
			name: "invalid retention",
			storage: Storage{
				RetentionTime:   "never",
				CleanupInterval: "1h",
			},
			wantErr: "retention_time",
		},
		{
			name: "zero retention",
			storage: Storage{
				RetentionTime:   "0h",
				CleanupInterval: "1h",
			},
			wantErr: "greater than zero",
		},
		{
			name: "invalid cleanup interval",
			storage: Storage{
				RetentionTime:   "2d",
				CleanupInterval: "invalid",
			},
			wantErr: "cleanup_interval",
		},
		{
			name: "negative cleanup interval",
			storage: Storage{
				RetentionTime:   "2d",
				CleanupInterval: "-1h",
			},
			wantErr: "greater than zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.storage.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("Validate() error = nil, want error containing %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Validate() error = %q, want substring %q", err, tt.wantErr)
			}
		})
	}
}
