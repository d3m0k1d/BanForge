package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var DetectedFirewall string

const (
	ConfigDir  = "/etc/banforge"
	ConfigFile = "config.toml"
)

func createFileWithPermissions(path string, perm os.FileMode) error {
	// #nosec G304 - path is controlled by config package not user
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	if err := os.Chmod(path, perm); err != nil {
		_ = file.Close()
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}
	return nil
}

func CreateConf() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("you must be root to run this command, use sudo/doas")
	}

	configPath := filepath.Join(ConfigDir, ConfigFile)

	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Config file already exists: %s\n", configPath)
		return nil
	}

	if err := os.MkdirAll(ConfigDir, 0750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(configPath, []byte(Base_config), 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	fmt.Printf("Config file created: %s\n", configPath)

	rulesDir := filepath.Join(ConfigDir, "rules.d")
	if err := os.MkdirAll(rulesDir, 0750); err != nil {
		return fmt.Errorf("failed to create rules directory: %w", err)
	}
	fmt.Printf("Rules directory created: %s\n", rulesDir)

	bansDBDir := filepath.Dir("/var/lib/banforge/bans.db")
	if err := os.MkdirAll(bansDBDir, 0750); err != nil {
		return fmt.Errorf("failed to create bans database directory: %w", err)
	}

	reqDBDir := filepath.Dir("/var/lib/banforge/requests.db")
	if err := os.MkdirAll(reqDBDir, 0750); err != nil {
		return fmt.Errorf("failed to create requests database directory: %w", err)
	}

	bansDBPath := "/var/lib/banforge/bans.db"
	if err := createFileWithPermissions(bansDBPath, 0600); err != nil {
		return fmt.Errorf("failed to create bans database file: %w", err)
	}
	fmt.Printf("Bans database file created: %s\n", bansDBPath)

	reqDBPath := "/var/lib/banforge/requests.db"
	if err := createFileWithPermissions(reqDBPath, 0600); err != nil {
		return fmt.Errorf("failed to create requests database file: %w", err)
	}
	fmt.Printf("Requests database file created: %s\n", reqDBPath)

	return nil
}

func FindFirewall() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("firewall settings needs sudo privileges")
	}

	firewalls := []string{"nft", "firewall-cmd", "iptables", "ufw"}
	for _, firewall := range firewalls {
		_, err := exec.LookPath(firewall)
		if err == nil {
			switch firewall {
			case "firewall-cmd":
				DetectedFirewall = "firewalld"
			case "nft":
				DetectedFirewall = "nftables"
			default:
				DetectedFirewall = firewall
			}

			fmt.Printf("Detected firewall: %s\n", DetectedFirewall)

			cfg := &Config{}
			_, err := toml.DecodeFile("/etc/banforge/config.toml", cfg)
			if err != nil {
				return fmt.Errorf("failed to decode config: %w", err)
			}

			cfg.Firewall.Name = DetectedFirewall

			file, err := os.Create("/etc/banforge/config.toml")
			if err != nil {
				return fmt.Errorf("failed to create config file: %w", err)
			}

			encoder := toml.NewEncoder(file)
			if err := encoder.Encode(cfg); err != nil {
				_ = file.Close()
				return fmt.Errorf("failed to encode config: %w", err)
			}

			if err := file.Close(); err != nil {
				return fmt.Errorf("failed to close file: %w", err)
			}

			fmt.Printf("Config updated with firewall: %s\n", DetectedFirewall)
			return nil
		}
	}

	return fmt.Errorf("firewall not found")
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	_, err := toml.DecodeFile("/etc/banforge/config.toml", cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}
	return cfg, nil
}
