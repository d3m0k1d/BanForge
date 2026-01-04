package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var Firewalls = map[string]string{
	"iptables":  "iptables",
	"nftables":  "nft",
	"ufw":       "ufw",
	"firewalld": "firewall-cmd",
}

var DetectedFirewall string

const (
	ConfigDir  = "/etc/banforge"
	ConfigFile = "config.toml"
)

func CreateConf() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("you must be root to run this command, use sudo/doas")
	}

	if err := os.MkdirAll(ConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(ConfigDir, ConfigFile)

	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Config file already exists: %s\n", configPath)
		return nil
	}

	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	if err := os.Chmod(configPath, 0644); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	fmt.Printf(" Config file created: %s\n", configPath)
	return nil
}

func GetSysConf() error {
	for name, binary := range Firewalls {
		if _, err := exec.LookPath(binary); err == nil {
			DetectedFirewall = name
			fmt.Printf("found firewall: %s\n", name)
			confstr := "firewall = \"" + name + "\""
			os.WriteFile(ConfigDir+"/"+ConfigFile, []byte(confstr), 0644)
			return nil
		}
	}
	return fmt.Errorf("no firewall found (checked iptables, nftables, ufw, firewalld) please install once of them")
}
