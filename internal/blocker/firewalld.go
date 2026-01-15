package blocker

import (
	"os/exec"

	"github.com/d3m0k1d/BanForge/internal/logger"
)

type Firewalld struct {
	logger *logger.Logger
}

func NewFirewalld(logger *logger.Logger) *Firewalld {
	return &Firewalld{
		logger: logger,
	}
}

func (f *Firewalld) Ban(ip string) error {
	err := validateIP(ip)
	if err != nil {
		return err
	}
	cmd := exec.Command("sudo", "firewall-cmd", "--zone=drop", "--add-source", ip, "--permanent")
	output, err := cmd.CombinedOutput()
	if err != nil {
		f.logger.Error(err.Error())
		return err
	}
	f.logger.Info("Add source " + ip + " " + string(output))
	output, err = exec.Command("sudo", "firewall-cmd", "--reload").CombinedOutput()
	if err != nil {
		f.logger.Error(err.Error())
		return err
	}
	f.logger.Info("Reload " + string(output))
	return nil
}

func (f *Firewalld) Unban(ip string) error {
	err := validateIP(ip)
	if err != nil {
		return err
	}
	cmd := exec.Command("sudo", "firewall-cmd", "--zone=drop", "--remove-source", ip, "--permanent")
	output, err := cmd.CombinedOutput()
	if err != nil {
		f.logger.Error(err.Error())
		return err
	}
	f.logger.Info("Remove source " + ip + " " + string(output))
	output, err = exec.Command("sudo", "firewall-cmd", "--reload").CombinedOutput()
	if err != nil {
		f.logger.Error(err.Error())
		return err
	}
	f.logger.Info("Reload " + string(output))
	return nil
}

func (f *Firewalld) Setup(config string) error {
	return nil
}
