package blocker

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"

	"github.com/d3m0k1d/BanForge/internal/logger"
	"github.com/d3m0k1d/BanForge/internal/metrics"
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

func (n *Nftables) Ban(ip string) error {
	set, address, err := nftablesSetForIP(ip)
	if err != nil {
		return err
	}
	metrics.IncBanAttempt("nftables")

	exists, err := nftablesElementExists(set, address)
	if err != nil {
		metrics.IncError()
		return err
	}
	if exists {
		n.logger.Info("IP already banned", "ip", address, "set", set)
		return nil
	}

	// #nosec G204 - address is parsed and normalized by net.ParseIP
	cmd := exec.Command(
		"nft", "add", "element", "inet", "banforge", set, "{", address, "}",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		n.logger.Error("failed to ban IP",
			"ip", address,
			"set", set,
			"error", err.Error(),
			"output", string(output))
		metrics.IncError()
		return fmt.Errorf("failed to add IP to nftables set: %w", err)
	}

	n.logger.Info("IP banned", "ip", address, "set", set)
	metrics.IncBan("nftables")
	return nil
}

func (n *Nftables) Unban(ip string) error {
	set, address, err := nftablesSetForIP(ip)
	if err != nil {
		return err
	}
	metrics.IncUnbanAttempt("nftables")

	exists, err := nftablesElementExists(set, address)
	if err != nil {
		metrics.IncError()
		return err
	}
	if !exists {
		n.logger.Info("IP already unbanned", "ip", address, "set", set)
		return nil
	}

	// #nosec G204 - address is parsed and normalized by net.ParseIP
	cmd := exec.Command(
		"nft", "delete", "element", "inet", "banforge", set, "{", address, "}",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		n.logger.Error("failed to unban IP",
			"ip", address,
			"set", set,
			"error", err.Error(),
			"output", string(output))
		metrics.IncError()
		return fmt.Errorf("failed to delete IP from nftables set: %w", err)
	}

	n.logger.Info("IP unbanned", "ip", address, "set", set)
	metrics.IncUnban("nftables")
	return nil
}

func nftablesSetForIP(ip string) (string, string, error) {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return "", "", fmt.Errorf("invalid IP: %s", ip)
	}

	if ipv4 := parsed.To4(); ipv4 != nil {
		return "blocked_ipv4", ipv4.String(), nil
	}

	return "blocked_ipv6", parsed.String(), nil
}

func nftablesElementExists(set string, ip string) (bool, error) {
	// #nosec G204 - set is selected internally and ip is normalized by net.ParseIP
	cmd := exec.Command(
		"nft", "get", "element", "inet", "banforge", set, "{", ip, "}",
	)
	if err := cmd.Run(); err == nil {
		return true, nil
	}

	// #nosec G204 - set is selected internally by nftablesSetForIP
	cmd = exec.Command("nft", "list", "set", "inet", "banforge", set)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to inspect nftables set %s: %w: %s", set, err, output)
	}

	return false, nil
}

func (n *Nftables) Setup(config string) error {
	if err := validateConfigPath(config); err != nil {
		return fmt.Errorf("path error: %w", err)
	}

	const banforgeConfigPath = "/var/lib/banforge/banforge.nft"
	const chains = `# managed by BanForge IPS system
table inet banforge {
	set blocked_ipv4 {
		type ipv4_addr
		flags timeout
	}

	set blocked_ipv6 {
		type ipv6_addr
		flags timeout
	}

	chain input {
		type filter hook input priority -100; policy accept;

		ip saddr @blocked_ipv4 drop
		ip6 saddr @blocked_ipv6 drop
	}
}
`
	const include = `include "/var/lib/banforge/banforge.nft"`
	const nftConfig = "\n# managed by BanForge IPS system\n" + include + "\n"

	if err := os.WriteFile(banforgeConfigPath, []byte(chains), 0600); err != nil {
		return fmt.Errorf("failed to write BanForge nftables config: %w", err)
	}
	if err := os.Chmod(banforgeConfigPath, 0600); err != nil {
		return fmt.Errorf("failed to set BanForge nftables config permissions: %w", err)
	}

	// #nosec G304 - config is an absolute administrator-managed path validated above
	file, err := os.ReadFile(config)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if !bytes.Contains(file, []byte(include)) {
		// #nosec G304 - config is an absolute administrator-managed path validated above
		conf, err := os.OpenFile(config, os.O_APPEND|os.O_WRONLY, 0)
		if err != nil {
			return fmt.Errorf("failed to open config file: %w", err)
		}

		if _, err := conf.WriteString(nftConfig); err != nil {
			_ = conf.Close()
			return fmt.Errorf("failed to add BanForge include: %w", err)
		}

		if err := conf.Close(); err != nil {
			return fmt.Errorf("failed to close config file: %w", err)
		}
	}

	tableExists := exec.Command("nft", "list", "table", "inet", "banforge").Run() == nil
	if tableExists {
		ipv4SetExists := exec.Command(
			"nft", "list", "set", "inet", "banforge", "blocked_ipv4",
		).Run() == nil
		ipv6SetExists := exec.Command(
			"nft", "list", "set", "inet", "banforge", "blocked_ipv6",
		).Run() == nil
		if ipv4SetExists && ipv6SetExists {
			return nil
		}

		output, err := exec.Command("nft", "delete", "table", "inet", "banforge").CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to replace old BanForge table: %s", string(output))
		}
	}

	// #nosec G204 - config is an absolute administrator-managed path validated above
	cmd := exec.Command("nft", "-f", config)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to load nftables config: %s", string(output))
	}

	return nil
}

func (n *Nftables) PortOpen(port int, protocol string) error {
	if port >= 0 && port <= 65535 {
		if protocol != "tcp" && protocol != "udp" {
			n.logger.Error("invalid protocol")
			metrics.IncError()
			return fmt.Errorf("invalid protocol")
		}
		s := strconv.Itoa(port)
		metrics.IncPortOperation("open", protocol)
		// #nosec G204 - managed by system adminstartor
		cmd := exec.Command(
			"nft",
			"add",
			"rule",
			"inet",
			"banforge",
			"input",
			protocol,
			"dport",
			s,
			"accept",
		)
		output, err := cmd.CombinedOutput()
		if err != nil {
			n.logger.Error(err.Error())
			metrics.IncError()
			return err
		}
		n.logger.Info("Add port " + s + " " + string(output))
	}
	return nil
}

func (n *Nftables) PortClose(port int, protocol string) error {
	if port >= 0 && port <= 65535 {
		if protocol != "tcp" && protocol != "udp" {
			n.logger.Error("invalid protocol")
			metrics.IncError()
			return fmt.Errorf("invalid protocol")
		}
		s := strconv.Itoa(port)
		metrics.IncPortOperation("close", protocol)
		// #nosec G204 - managed by system adminstartor
		cmd := exec.Command(
			"nft",
			"add",
			"rule",
			"inet",
			"banforge",
			"input",
			protocol,
			"dport",
			s,
			"drop",
		)
		output, err := cmd.CombinedOutput()
		if err != nil {
			n.logger.Error(err.Error())
			metrics.IncError()
			return err
		}
		n.logger.Info("Add port " + s + " " + string(output))

	}
	return nil
}
