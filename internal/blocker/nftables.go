package blocker

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/d3m0k1d/BanForge/internal/logger"
)

type Nftables struct {
	logger *logger.Logger
	config string
}

func NewNftables(logger *logger.Logger, config string) *Nftables {
	return &Nftables{
		logger: logger,
		config: config,
	}
}

// Name returns the blocker engine name
func (n *Nftables) Name() string {
	return "nftables"
}

// IsAvailable checks if nftables is available in the system
func (n *Nftables) IsAvailable() bool {
	cmd := exec.Command("which", "nft")
	return cmd.Run() == nil
}

// Setup initializes nftables with required tables and chains
func (n *Nftables) Setup() error {
	return SetupNftables(n.config)
}

// Ban adds an IP to the banned list
func (n *Nftables) Ban(ip string) error {
	err := validateIP(ip)
	if err != nil {
		return err
	}

	cmd := exec.Command("sudo", "nft", "add", "rule", "inet", "banforge", "banned",
		"ip", "saddr", ip, "drop")
	output, err := cmd.CombinedOutput()
	if err != nil {
		n.logger.Error("failed to ban IP",
			"ip", ip,
			"error", err.Error(),
			"output", string(output))
		return err
	}

	n.logger.Info("IP banned", "ip", ip)

	err = saveNftablesConfig(n.config)
	if err != nil {
		n.logger.Error("failed to save config",
			"config_path", n.config,
			"error", err.Error())
		return err
	}

	n.logger.Info("config saved", "config_path", n.config)
	return nil
}

// Unban removes an IP from the banned list
func (n *Nftables) Unban(ip string) error {
	err := validateIP(ip)
	if err != nil {
		return err
	}

	handle, err := n.findRuleHandle(ip)
	if err != nil {
		n.logger.Error("failed to find rule handle",
			"ip", ip,
			"error", err.Error())
		return err
	}

	if handle == "" {
		n.logger.Warn("no rule found for IP", "ip", ip)
		return fmt.Errorf("no rule found for IP %s", ip)
	}
	// #nosec G204 - handle is extracted from nftables output and validated
	cmd := exec.Command("sudo", "nft", "delete", "rule", "inet", "banforge", "banned",
		"handle", handle)
	output, err := cmd.CombinedOutput()
	if err != nil {
		n.logger.Error("failed to unban IP",
			"ip", ip,
			"handle", handle,
			"error", err.Error(),
			"output", string(output))
		return err
	}

	n.logger.Info("IP unbanned", "ip", ip, "handle", handle)

	err = saveNftablesConfig(n.config)
	if err != nil {
		n.logger.Error("failed to save config",
			"config_path", n.config,
			"error", err.Error())
		return err
	}

	n.logger.Info("config saved", "config_path", n.config)
	return nil
}

// List returns all currently banned IPs
func (n *Nftables) List() ([]string, error) {
	cmd := exec.Command("sudo", "nft", "-a", "list", "chain", "inet", "banforge", "banned")
	output, err := cmd.CombinedOutput()
	if err != nil {
		n.logger.Error("failed to list banned IPs",
			"error", err.Error(),
			"output", string(output))
		return nil, err
	}

	var bannedIPs []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "drop") && strings.Contains(line, "saddr") {
			// Extract IP from line like: ip saddr 10.0.0.1 drop # handle 2
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "saddr" && i+1 < len(parts) {
					ip := parts[i+1]
					if validateIP(ip) == nil {
						bannedIPs = append(bannedIPs, ip)
					}
					break
				}
			}
		}
	}

	return bannedIPs, nil
}

// Close performs cleanup operations (placeholder for future use)
func (n *Nftables) Close() error {
	// No cleanup needed for nftables
	n.logger.Info("nftables blocker closed")
	return nil
}

func SetupNftables(config string) error {
	err := validateConfigPath(config)
	if err != nil {
		return err
	}

	cmd := exec.Command("sudo", "nft", "list", "table", "inet", "banforge")
	if err := cmd.Run(); err != nil {
		cmd = exec.Command("sudo", "nft", "add", "table", "inet", "banforge")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create table: %s", string(output))
		}
	}

	cmd = exec.Command("sudo", "nft", "list", "chain", "inet", "banforge", "input")
	if err := cmd.Run(); err != nil {
		script := "sudo nft 'add chain inet banforge input { type filter hook input priority 0; policy accept; }'"
		cmd = exec.Command("bash", "-c", script)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create input chain: %s", string(output))
		}
	}

	cmd = exec.Command("sudo", "nft", "list", "chain", "inet", "banforge", "banned")
	if err := cmd.Run(); err != nil {
		cmd = exec.Command("sudo", "nft", "add", "chain", "inet", "banforge", "banned")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create banned chain: %s", string(output))
		}
	}

	cmd = exec.Command("sudo", "nft", "-a", "list", "chain", "inet", "banforge", "input")
	output, err := cmd.CombinedOutput()
	if err == nil && !strings.Contains(string(output), "jump banned") {
		cmd = exec.Command("sudo", "nft", "add", "rule", "inet", "banforge", "input", "jump", "banned")
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to add jump rule: %s", string(output))
		}
	}

	err = saveNftablesConfig(config)
	if err != nil {
		return fmt.Errorf("failed to save nftables config: %w", err)
	}

	return nil
}

func (n *Nftables) findRuleHandle(ip string) (string, error) {
	cmd := exec.Command("sudo", "nft", "-a", "list", "chain", "inet", "banforge", "banned")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to list chain rules: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, ip) && strings.Contains(line, "drop") {
			if idx := strings.Index(line, "# handle"); idx != -1 {
				parts := strings.Fields(line[idx:])
				if len(parts) >= 3 && parts[1] == "handle" {
					return parts[2], nil
				}
			}
		}
	}

	return "", nil
}

func saveNftablesConfig(configPath string) error {
	err := validateConfigPath(configPath)
	if err != nil {
		return err
	}

	cmd := exec.Command("sudo", "nft", "list", "ruleset")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to get nftables ruleset: %w", err)
	}

	cmd = exec.Command("sudo", "tee", configPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start tee command: %w", err)
	}

	_, err = stdin.Write(output)
	if err != nil {
		return fmt.Errorf("failed to write to config file: %w", err)
	}
	err = stdin.Close()
	if err != nil {
		return fmt.Errorf("failed to close stdin pipe: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}
