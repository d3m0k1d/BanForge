package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/d3m0k1d/BanForge/internal/logger"
)

func LoadRuleConfig() ([]Rule, error) {
	log := logger.New(false)
	var cfg Rules

	_, err := toml.DecodeFile("/etc/banforge/rules.toml", &cfg)
	if err != nil {
		log.Error(fmt.Sprintf("failed to decode config: %v", err))
		return nil, err
	}

	log.Info(fmt.Sprintf("loaded %d rules", len(cfg.Rules)))
	return cfg.Rules, nil
}
