package blocker

import (
	"os/exec"
	"strconv"

	"github.com/d3m0k1d/BanForge/internal/logger"
	"github.com/d3m0k1d/BanForge/internal/metrics"
)

type Iptables struct {
	logger *logger.Logger
	config string
}

func NewIptables(logger *logger.Logger, config string) *Iptables {
	return &Iptables{
		logger: logger,
		config: config,
	}
}

func (f *Iptables) Ban(ip string) error {
	err := validateIP(ip)
	if err != nil {
		return err
	}
	metrics.IncBanAttempt("iptables")
	err = validateConfigPath(f.config)
	if err != nil {
		return err
	}
	// #nosec G204 - f.config is validated above via validateConfigPath()
	cmd := exec.Command("iptables", "-A", "INPUT", "-s", ip, "-j", "DROP")
	output, err := cmd.CombinedOutput()
	if err != nil {
		f.logger.Error("failed to ban IP",
			"ip", ip,
			"error", err.Error(),
			"output", string(output))
		metrics.IncError()
		return err
	}
	f.logger.Info("IP banned",
		"ip", ip,
		"output", string(output))
	metrics.IncBan("iptables")

	err = validateConfigPath(f.config)
	if err != nil {
		return err
	}
	// #nosec G204 - f.config is validated above via validateConfigPath()
	cmd = exec.Command("iptables-save", "-f", f.config)
	output, err = cmd.CombinedOutput()
	if err != nil {
		f.logger.Error("failed to save config",
			"config_path", f.config,
			"error", err.Error(),
			"output", string(output))
		metrics.IncError()
		return err
	}
	f.logger.Info("config saved",
		"config_path", f.config,
		"output", string(output))
	return nil
}

func (f *Iptables) Unban(ip string) error {
	err := validateIP(ip)
	if err != nil {
		return err
	}
	metrics.IncUnbanAttempt("iptables")
	err = validateConfigPath(f.config)
	if err != nil {
		return err
	}
	// #nosec G204 - f.config is validated above via validateConfigPath()
	cmd := exec.Command("iptables", "-D", "INPUT", "-s", ip, "-j", "DROP")
	output, err := cmd.CombinedOutput()
	if err != nil {
		f.logger.Error("failed to unban IP",
			"ip", ip,
			"error", err.Error(),
			"output", string(output))
		metrics.IncError()
		return err
	}
	f.logger.Info("IP unbanned",
		"ip", ip,
		"output", string(output))
	metrics.IncUnban("iptables")

	err = validateConfigPath(f.config)
	if err != nil {
		return err
	}
	// #nosec G204 - f.config is validated above via validateConfigPath()
	cmd = exec.Command("iptables-save", "-f", f.config)
	output, err = cmd.CombinedOutput()
	if err != nil {
		f.logger.Error("failed to save config",
			"config_path", f.config,
			"error", err.Error(),
			"output", string(output))
		metrics.IncError()
		return err
	}
	f.logger.Info("config saved",
		"config_path", f.config,
		"output", string(output))
	return nil
}

func (f *Iptables) PortOpen(port int, protocol string) error {
	if port >= 0 && port <= 65535 {
		if protocol != "tcp" && protocol != "udp" {
			f.logger.Error("invalid protocol")
			return nil
		}
		s := strconv.Itoa(port)
		metrics.IncPortOperation("open", protocol)
		// #nosec G204 - managed by system adminstartor
		cmd := exec.Command("iptables", "-A", "INPUT", "-p", protocol, "--dport", s, "-j", "ACCEPT")
		output, err := cmd.CombinedOutput()
		if err != nil {
			f.logger.Error(err.Error())
			metrics.IncError()
			return err
		}
		f.logger.Info("Add port " + s + " " + string(output))
		// #nosec G204 - f.config is validated above via validateConfigPath()
		cmd = exec.Command("iptables-save", "-f", f.config)
		output, err = cmd.CombinedOutput()
		if err != nil {
			f.logger.Error("failed to save config",
				"config_path", f.config,
				"error", err.Error(),
				"output", string(output))
			metrics.IncError()
			return err
		}
	}
	return nil
}

func (f *Iptables) PortClose(port int, protocol string) error {
	if port >= 0 && port <= 65535 {
		if protocol != "tcp" && protocol != "udp" {
			f.logger.Error("invalid protocol")
			return nil
		}
		s := strconv.Itoa(port)
		metrics.IncPortOperation("close", protocol)
		// #nosec G204 - managed by system adminstartor
		cmd := exec.Command("iptables", "-D", "INPUT", "-p", protocol, "--dport", s, "-j", "ACCEPT")
		output, err := cmd.CombinedOutput()
		if err != nil {
			f.logger.Error(err.Error())
			metrics.IncError()
			return err
		}
		f.logger.Info("Add port " + s + " " + string(output))
		// #nosec G204 - f.config is validated above via validateConfigPath()
		cmd = exec.Command("iptables-save", "-f", f.config)
		output, err = cmd.CombinedOutput()
		if err != nil {
			f.logger.Error("failed to save config",
				"config_path", f.config,
				"error", err.Error(),
				"output", string(output))
			metrics.IncError()
			return err
		}
	}
	return nil
}

func (f *Iptables) Setup(config string) error {
	return nil
}
