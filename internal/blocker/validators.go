package blocker

import (
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"strings"
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

func validateConfigPath(pathIn string) error {
	if pathIn == "" {
		return errors.New("config path cannot be empty")
	}

	cleanPath := filepath.Clean(pathIn)

	if !filepath.IsAbs(cleanPath) {
		return fmt.Errorf("config path must be absolute, got: %s", cleanPath)
	}

	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("config path contains path traversal: %s", cleanPath)
	}

	return nil
}
