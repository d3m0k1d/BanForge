package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	fmt.Printf(" Config file created: %s\n", configPath)
	return nil
}

func FindFirewall() error {

	if os.Getegid() != 0 {
		fmt.Printf("Firewall settings needs sudo privileges\n")
		os.Exit(1)
	}
	firewalls := []string{"iptables", "nft", "firewall-cmd", "ufw"}
	for _, firewall := range firewalls {
		_, err := exec.LookPath(firewall)
		if err == nil {
			if firewall == "firewall-cmd" {
				DetectedFirewall = "firewalld"
			}
			if firewall == "nft" {
				DetectedFirewall = "nftables"
			}
			DetectedFirewall = firewall
			fmt.Printf("Detected firewall: %s\n", firewall)
			return nil
		}
	}
	return fmt.Errorf("no firewall found (checked ufw, firewall-cmd, iptables, nft) please install one of them")
}
