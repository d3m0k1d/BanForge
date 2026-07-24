package blocker

import (
	"testing"
)

func TestValidateConfigPath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "empty", input: "", wantErr: true},
		{name: "valid path", input: "/path/to/config", wantErr: false},
		{name: "invalid path", input: "path/to/config", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfigPath(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfigPath(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateIP(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "empty", input: "", wantErr: true},
		{name: "invalid IP", input: "1.1.1", wantErr: true},
		{name: "valid IP", input: "1.1.1.1", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateIP(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateIP(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestNftablesSetForIP(t *testing.T) {
	tests := []struct {
		name        string
		ip          string
		wantSet     string
		wantAddress string
		wantErr     bool
	}{
		{
			name:        "IPv4",
			ip:          "192.0.2.10",
			wantSet:     "blocked_ipv4",
			wantAddress: "192.0.2.10",
		},
		{
			name:        "IPv6",
			ip:          "2001:0db8:0:0:0:0:0:1",
			wantSet:     "blocked_ipv6",
			wantAddress: "2001:db8::1",
		},
		{
			name:    "invalid",
			ip:      "not-an-ip",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set, address, err := nftablesSetForIP(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Fatalf("nftablesSetForIP(%q) error = %v, wantErr %v", tt.ip, err, tt.wantErr)
			}
			if set != tt.wantSet {
				t.Errorf("nftablesSetForIP(%q) set = %q, want %q", tt.ip, set, tt.wantSet)
			}
			if address != tt.wantAddress {
				t.Errorf("nftablesSetForIP(%q) address = %q, want %q", tt.ip, address, tt.wantAddress)
			}
		})
	}
}
