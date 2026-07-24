package config

import "fmt"

const (
	defaultRetentionTime   = "2d"
	defaultCleanupInterval = "1h"
)

func newConfigWithDefaults() *Config {
	return &Config{
		Storage: Storage{
			RetentionTime:   defaultRetentionTime,
			CleanupInterval: defaultCleanupInterval,
		},
	}
}

func (c Config) Validate() error {
	if err := c.Storage.Validate(); err != nil {
		return fmt.Errorf("storage: %w", err)
	}

	return nil
}

func (s Storage) Validate() error {
	retention, err := ParseDurationWithYears(s.RetentionTime)
	if err != nil {
		return fmt.Errorf("invalid retention_time %q: %w", s.RetentionTime, err)
	}
	if retention <= 0 {
		return fmt.Errorf("retention_time must be greater than zero")
	}

	cleanupInterval, err := ParseDurationWithYears(s.CleanupInterval)
	if err != nil {
		return fmt.Errorf("invalid cleanup_interval %q: %w", s.CleanupInterval, err)
	}
	if cleanupInterval <= 0 {
		return fmt.Errorf("cleanup_interval must be greater than zero")
	}

	return nil
}
