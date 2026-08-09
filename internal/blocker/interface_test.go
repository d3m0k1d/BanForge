package blocker

import "testing"

func TestGetBlockerKnown(t *testing.T) {
	known := []string{"ufw", "iptables", "nftables", "firewalld"}
	for _, fw := range known {
		t.Run(fw, func(t *testing.T) {
			b, err := GetBlocker(fw, "/nonexistent/config")
			if err != nil {
				t.Fatalf("GetBlocker(%q) unexpected error: %v", fw, err)
			}
			if b == nil {
				t.Fatalf("GetBlocker(%q) returned nil engine", fw)
			}
		})
	}
}

func TestGetBlockerUnknown(t *testing.T) {
	b, err := GetBlocker("nonexistent", "")
	if err == nil {
		t.Fatal("GetBlocker() expected error for unknown firewall")
	}
	if b != nil {
		t.Errorf("GetBlocker() returned engine %v, want nil", b)
	}
}
