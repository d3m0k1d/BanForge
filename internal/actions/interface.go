package actions

import "github.com/d3m0k1d/BanForge/internal/config"

type Executor struct {
	Action config.Action
}

func (e *Executor) Execute() error {
	switch e.Action.Type {
	case "email":
		return SendEmail(e.Action)
	case "webhook":
		return SendWebhook(e.Action)
	case "script":
		return RunScript(e.Action)
	}
	return nil
}
