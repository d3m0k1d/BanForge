package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/d3m0k1d/BanForge/internal/logger"
	"github.com/d3m0k1d/BanForge/internal/metrics"
)

func LoadMetricsConfig() (*Metrics, error) {
	cfg := &Metrics{}
	_, err := toml.DecodeFile("/etc/banforge/config.toml", cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	if cfg.Enabled && cfg.Port > 0 && cfg.Port < 65535 {
		go metrics.StartMetricsServer(cfg.Port)
	} else if cfg.Enabled {
		fmt.Println("Metrics enabled but port invalid, not starting server")
	}

	return cfg, nil
}

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

func NewRule(
	Name string,
	ServiceName string,
	Path string,
	Status string,
	Method string,
	ttl string,
	max_retry int,
) error {
	r, err := LoadRuleConfig()
	if err != nil {
		r = []Rule{}
	}
	if Name == "" {
		fmt.Printf("Rule name can't be empty\n")
		return nil
	}
	r = append(
		r,
		Rule{
			Name:        Name,
			ServiceName: ServiceName,
			Path:        Path,
			Status:      Status,
			Method:      Method,
			BanTime:     ttl,
			MaxRetry:    max_retry,
		},
	)
	file, err := os.Create("/etc/banforge/rules.toml")
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, file.Close())
	}()
	cfg := Rules{Rules: r}
	err = toml.NewEncoder(file).Encode(cfg)
	if err != nil {
		return err
	}
	return nil
}

func EditRule(Name string, ServiceName string, Path string, Status string, Method string) error {
	if Name == "" {
		return fmt.Errorf("Rule name can't be empty")
	}

	r, err := LoadRuleConfig()
	if err != nil {
		return fmt.Errorf("rules is empty, please use 'banforge add rule' or create rules.toml")
	}

	found := false
	for i, rule := range r {
		if rule.Name == Name {
			found = true

			if ServiceName != "" {
				r[i].ServiceName = ServiceName
			}
			if Path != "" {
				r[i].Path = Path
			}
			if Status != "" {
				r[i].Status = Status
			}
			if Method != "" {
				r[i].Method = Method
			}
			break
		}
	}

	if !found {
		return fmt.Errorf("rule '%s' not found", Name)
	}

	file, err := os.Create("/etc/banforge/rules.toml")
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	cfg := Rules{Rules: r}
	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

func ParseDurationWithYears(s string) (time.Duration, error) {
	if ss, ok := strings.CutSuffix(s, "y"); ok {
		years, err := strconv.Atoi(ss)
		if err != nil {
			return 0, err
		}
		return time.Duration(years) * 365 * 24 * time.Hour, nil
	}

	if ss, ok := strings.CutSuffix(s, "M"); ok {
		months, err := strconv.Atoi(ss)
		if err != nil {
			return 0, err
		}
		return time.Duration(months) * 30 * 24 * time.Hour, nil
	}

	if ss, ok := strings.CutSuffix(s, "d"); ok {
		days, err := strconv.Atoi(ss)
		if err != nil {
			return 0, err
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}

	return time.ParseDuration(s)
}
