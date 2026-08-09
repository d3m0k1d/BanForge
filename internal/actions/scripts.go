package actions

import (
	"fmt"
	"os/exec"

	"github.com/d3m0k1d/BanForge/internal/config"
)

func RunScript(action config.Action) error {
	if !action.Enabled {
		return nil
	}
	if action.Script == "" {
		return fmt.Errorf("script on config is empty")
	}
	// #nosec G204 - managed by system adminstartor
	cmd := exec.Command(action.Script)
	if action.Interpretator != "" {
		// #nosec G204 - managed by system adminstartor
		cmd = exec.Command(action.Interpretator, action.Script)
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run script: %w", err)
	}
	return nil
}
