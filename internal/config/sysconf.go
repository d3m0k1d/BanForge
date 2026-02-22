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

func CreateConf() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("you must be root to run this command, use sudo/doas")
	}

	configPath := filepath.Join(ConfigDir, ConfigFile)

	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Config file already exists: %s\n", configPath)
		return nil
	}

	file, err := os.Create("/etc/banforge/config.toml")
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	if err := os.Chmod(configPath, 0600); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}
	err = os.WriteFile(configPath, []byte(Base_config), 0600)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	fmt.Printf(" Config file created: %s\n", configPath)
	file, err = os.Create("/etc/banforge/rules.toml")
	if err != nil {
		return fmt.Errorf("failed to create rules file: %w", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	bansDBPath := "/var/lib/banforge/bans.db"
	reqDBPath := "/var/lib/banforge/requests.db"

	file, err = os.Create(bansDBPath)
	if err != nil {
		return fmt.Errorf("failed to create bans database file: %w", err)
	}
	err = os.Chmod(bansDBPath, 0600)
	if err != nil {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
		return fmt.Errorf("failed to set permissions: %w", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	file, err = os.Create(reqDBPath)
	if err != nil {
		return fmt.Errorf("failed to create requests database file: %w", err)
	}
	err = os.Chmod(reqDBPath, 0600)
	if err != nil {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
		return fmt.Errorf("failed to set permissions: %w", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	fmt.Printf(" Rules file created: %s\n", configPath)
	return nil
}

func FindFirewall() error {
	if os.Getegid() != 0 {
		fmt.Printf("Firewall settings needs sudo privileges\n")
		os.Exit(1)
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
				err = file.Close()
				if err != nil {
					return fmt.Errorf("failed to close file: %w", err)
				}
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
