package blocker

import (
	"os/exec"

	"github.com/d3m0k1d/BanForge/internal/logger"
)

type Ufw struct {
	logger *logger.Logger
}

func NewUfw(logger *logger.Logger) *Ufw {
	return &Ufw{
		logger: logger,
	}
}

func (ufw *Ufw) Ban(ip string) error {
	err := validateIP(ip)
	if err != nil {
		return err
	}
	cmd := exec.Command("sudo", "ufw", "--force", "deny", "from", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		ufw.logger.Error(err.Error())
		return err
	}
	ufw.logger.Info("Banning " + ip + " " + string(output))
	return nil
}

func (ufw *Ufw) Unban(ip string) error {
	err := validateIP(ip)
	if err != nil {
		return err
	}
	cmd := exec.Command("sudo", "ufw", "--force", "delete", "deny", "from", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		ufw.logger.Error(err.Error())
		return err
	}
	ufw.logger.Info("Unbanning " + ip + " " + string(output))
	return nil
}
