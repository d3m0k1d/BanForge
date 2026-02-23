package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

func LoadRuleConfig() ([]Rule, error) {
	const rulesDir = "/etc/banforge/rules.d"

	var cfg Rules

	files, err := os.ReadDir(rulesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read rules directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".toml") {
			continue
		}

		filePath := filepath.Join(rulesDir, file.Name())
		var fileCfg Rules

		if _, err := toml.DecodeFile(filePath, &fileCfg); err != nil {
			return nil, fmt.Errorf("failed to parse rule file %s: %w", filePath, err)
		}

		cfg.Rules = append(cfg.Rules, fileCfg.Rules...)
	}

	return cfg.Rules, nil
}

func NewRule(
	name string,
	serviceName string,
	path string,
	status string,
	method string,
	ttl string,
	maxRetry int,
) error {
	if name == "" {
		return fmt.Errorf("rule name can't be empty")
	}

	rule := Rule{
		Name:        name,
		ServiceName: serviceName,
		Path:        path,
		Status:      status,
		Method:      method,
		BanTime:     ttl,
		MaxRetry:    maxRetry,
	}

	filePath := filepath.Join("/etc/banforge/rules.d", SanitizeRuleFilename(name)+".toml")

	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("rule with name '%s' already exists", name)
	}

	cfg := Rules{Rules: []Rule{rule}}

	// #nosec G304 - validate by sanitizeRuleFilename
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create rule file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("warning: failed to close rule file: %v\n", closeErr)
		}
	}()

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode rule: %w", err)
	}

	return nil
}

func EditRule(name string, serviceName string, path string, status string, method string) error {
	if name == "" {
		return fmt.Errorf("rule name can't be empty")
	}

	rules, err := LoadRuleConfig()
	if err != nil {
		return fmt.Errorf("failed to load rules: %w", err)
	}

	found := false
	var updatedRule *Rule
	for i, rule := range rules {
		if rule.Name == name {
			found = true
			updatedRule = &rules[i]

			if serviceName != "" {
				updatedRule.ServiceName = serviceName
			}
			if path != "" {
				updatedRule.Path = path
			}
			if status != "" {
				updatedRule.Status = status
			}
			if method != "" {
				updatedRule.Method = method
			}
			break
		}
	}

	if !found {
		return fmt.Errorf("rule '%s' not found", name)
	}

	filePath := filepath.Join("/etc/banforge/rules.d", SanitizeRuleFilename(name)+".toml")
	cfg := Rules{Rules: []Rule{*updatedRule}}

	// #nosec G304 - validate by sanitizeRuleFilename
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to update rule file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("warning: failed to close rule file: %v\n", closeErr)
		}
	}()

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode updated rule: %w", err)
	}

	return nil
}

func SanitizeRuleFilename(name string) string {
	result := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, name)
	return strings.ToLower(result)
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
