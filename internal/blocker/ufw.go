package blocker

import (
	"github.com/d3m0k1d/BanForge/internal/logger"
	"os/exec"
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
	validateIP(ip)
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
	validateIP(ip)
	cmd := exec.Command("sudo", "ufw", "--force", "delete", "deny", "from", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		ufw.logger.Error(err.Error())
		return err
	}
	ufw.logger.Info("Unbanning " + ip + " " + string(output))
	return nil
}
