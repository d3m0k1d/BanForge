package blocker

import (
	"fmt"
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

func (u *Ufw) Ban(ip string) error {
	err := validateIP(ip)
	if err != nil {
		return err
	}

	cmd := exec.Command("sudo", "ufw", "--force", "deny", "from", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		u.logger.Error("failed to ban IP",
			"ip", ip,
			"error", err.Error(),
			"output", string(output))
		return fmt.Errorf("failed to ban IP %s: %w", ip, err)
	}

	u.logger.Info("IP banned", "ip", ip, "output", string(output))
	return nil
}
func (u *Ufw) Unban(ip string) error {
	err := validateIP(ip)
	if err != nil {
		return err
	}

	cmd := exec.Command("sudo", "ufw", "--force", "delete", "deny", "from", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		u.logger.Error("failed to unban IP",
			"ip", ip,
			"error", err.Error(),
			"output", string(output))
		return fmt.Errorf("failed to unban IP %s: %w", ip, err)
	}

	u.logger.Info("IP unbanned", "ip", ip, "output", string(output))
	return nil
}

func (u *Ufw) Setup(config string) error {
	if config != "" {
		fmt.Printf("Ufw dont support config file\n")
		cmd := exec.Command("sudo", "ufw", "enable")
		output, err := cmd.CombinedOutput()
		if err != nil {
			u.logger.Error("failed to enable ufw",
				"error", err.Error(),
				"output", string(output))
			return fmt.Errorf("failed to enable ufw: %w", err)
		}
	}
	if config == "" {
		cmd := exec.Command("sudo", "ufw", "enable")
		output, err := cmd.CombinedOutput()
		if err != nil {
			u.logger.Error("failed to enable ufw",
				"error", err.Error(),
				"output", string(output))
			return fmt.Errorf("failed to enable ufw: %w", err)
		}
	}
	return nil
}
