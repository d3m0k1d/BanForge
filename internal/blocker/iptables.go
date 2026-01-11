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
	cmd := exec.Command("sudo", "iptables", "-A", "INPUT", "-s", ip, "-j", "DROP")
	output, err := cmd.CombinedOutput()
	if err != nil {
		f.logger.Error(err.Error())
		return err
	}
	f.logger.Info("Banning " + ip + " " + string(output))
	cmd = exec.Command("sudo", "iptables-save", "-f", f.config)
	output, err = cmd.CombinedOutput()
	if err != nil {
		f.logger.Error(err.Error())
		return err
	}
	f.logger.Info("Config saved " + string(output))
	return nil
}

func (f *Iptables) Unban(ip string) error {
	err := validateIP(ip)
	if err != nil {
		return err
	}
	cmd := exec.Command("sudo", "iptables", "-D", "INPUT", "-s", ip, "-j", "DROP")
	output, err := cmd.CombinedOutput()
	if err != nil {
		f.logger.Error(err.Error())
		return err
	}
	f.logger.Info("Unbanning " + ip + " " + string(output))
	cmd = exec.Command("sudo", "iptables-save", "-f", f.config)
	output, err = cmd.CombinedOutput()
	if err != nil {
		f.logger.Error(err.Error())
		return err
	}
	f.logger.Info("Config saved " + string(output))
	return nil
}
