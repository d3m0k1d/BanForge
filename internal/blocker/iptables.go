package blocker

import (
	"os/exec"

	"github.com/d3m0k1d/BanForge/internal/logger"
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
	err = validateConfigPath(f.config)
	if err != nil {
		return err
	}
	cmd := exec.Command("iptables", "-A", "INPUT", "-s", ip, "-j", "DROP")
	output, err := cmd.CombinedOutput()
	if err != nil {
		f.logger.Error("failed to ban IP",
			"ip", ip,
			"error", err.Error(),
			"output", string(output))
		return err
	}
	f.logger.Info("IP banned",
		"ip", ip,
		"output", string(output))

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
	err = validateConfigPath(f.config)
	if err != nil {
		return err
	}
	cmd := exec.Command("iptables", "-D", "INPUT", "-s", ip, "-j", "DROP")
	output, err := cmd.CombinedOutput()
	if err != nil {
		f.logger.Error("failed to unban IP",
			"ip", ip,
			"error", err.Error(),
			"output", string(output))
		return err
	}
	f.logger.Info("IP unbanned",
		"ip", ip,
		"output", string(output))

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
		return err
	}
	f.logger.Info("config saved",
		"config_path", f.config,
		"output", string(output))
	return nil
}

func (f *Iptables) PortOpen(port int) error {
	return nil
}

func (f *Iptables) PortClose(port int) error {
	return nil
}

func (f *Iptables) Setup(config string) error {
	return nil
}
