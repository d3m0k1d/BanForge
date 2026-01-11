package blocker

import (
	"fmt"
	"net"
)

func validateIP(ip string) error {
	if ip == "" {
		return fmt.Errorf("empty IP")
	}

	if net.ParseIP(ip) == nil {
		return fmt.Errorf("invalid IP: %s", ip)
	}

	return nil
}
