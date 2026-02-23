package blocker

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/d3m0k1d/BanForge/internal/logger"
	"github.com/d3m0k1d/BanForge/internal/metrics"
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
	metrics.IncBanAttempt("firewalld")
	// #nosec G204 - ip is validated
	cmd := exec.Command("firewall-cmd", "--zone=drop", "--add-source", ip, "--permanent")
	output, err := cmd.CombinedOutput()
	if err != nil {
		f.logger.Error(err.Error())
		metrics.IncError()
		return err
	}
	f.logger.Info("Add source " + ip + " " + string(output))
	output, err = exec.Command("firewall-cmd", "--reload").CombinedOutput()
	if err != nil {
		f.logger.Error(err.Error())
		metrics.IncError()
		return err
	}
	f.logger.Info("Reload " + string(output))
	metrics.IncBan("firewalld")
	return nil
}

func (f *Firewalld) Unban(ip string) error {
	err := validateIP(ip)
	if err != nil {
		return err
	}
	metrics.IncUnbanAttempt("firewalld")
	// #nosec G204 - ip is validated
	cmd := exec.Command("firewall-cmd", "--zone=drop", "--remove-source", ip, "--permanent")
	output, err := cmd.CombinedOutput()
	if err != nil {
		f.logger.Error(err.Error())
		metrics.IncError()
		return err
	}
	f.logger.Info("Remove source " + ip + " " + string(output))
	output, err = exec.Command("firewall-cmd", "--reload").CombinedOutput()
	if err != nil {
		f.logger.Error(err.Error())
		metrics.IncError()
		return err
	}
	f.logger.Info("Reload " + string(output))
	metrics.IncUnban("firewalld")
	return nil
}

func (f *Firewalld) PortOpen(port int, protocol string) error {
	// #nosec G204 - handle is extracted from Firewalld output and validated
	if port >= 0 && port <= 65535 {
		if protocol != "tcp" && protocol != "udp" {
			f.logger.Error("invalid protocol")
			return fmt.Errorf("invalid protocol")
		}
		s := strconv.Itoa(port)
		metrics.IncPortOperation("open", protocol)
		cmd := exec.Command(
			"firewall-cmd",
			"--zone=public",
			"--add-port="+s+"/"+protocol,
			"--permanent",
		)
		output, err := cmd.CombinedOutput()
		if err != nil {
			f.logger.Error(err.Error())
			metrics.IncError()
			return err
		}
		f.logger.Info("Add port " + s + " " + string(output))
		output, err = exec.Command("firewall-cmd", "--reload").CombinedOutput()
		if err != nil {
			f.logger.Error(err.Error())
			metrics.IncError()
			return err
		}
		f.logger.Info("Reload " + string(output))
	}
	return nil
}

func (f *Firewalld) PortClose(port int, protocol string) error {
	// #nosec G204 - handle is extracted from nftables output and validated
	if port >= 0 && port <= 65535 {
		if protocol != "tcp" && protocol != "udp" {
			return fmt.Errorf("invalid protocol")
		}
		s := strconv.Itoa(port)
		metrics.IncPortOperation("close", protocol)
		cmd := exec.Command(
			"firewall-cmd",
			"--zone=public",
			"--remove-port="+s+"/"+protocol,
			"--permanent",
		)
		output, err := cmd.CombinedOutput()
		if err != nil {
			metrics.IncError()
			return err
		}
		f.logger.Info("Remove port " + s + " " + string(output))
		output, err = exec.Command("firewall-cmd", "--reload").CombinedOutput()
		if err != nil {
			metrics.IncError()
			return err
		}
		f.logger.Info("Reload " + string(output))
	}
	return nil
}

func (f *Firewalld) Setup(config string) error {
	return nil
}
