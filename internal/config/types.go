package config

type Firewall struct {
	Name   string `toml:"name"`
	Config string `toml:"config"`
}

type Service struct {
	Name    string `toml:"name"`
	Logging string `toml:"logging"`
	LogPath string `toml:"log_path"`
	Enabled bool   `toml:"enabled"`
}

type Config struct {
	Firewall Firewall  `toml:"firewall"`
	Metrics  Metrics   `toml:"metrics"`
	Service  []Service `toml:"service"`
}

// Rules
type Rules struct {
	Rules []Rule `toml:"rule"`
}

type Rule struct {
	Name        string `toml:"name"`
	ServiceName string `toml:"service"`
	Path        string `toml:"path"`
	Status      string `toml:"status"`
	Method      string `toml:"method"`
	MaxRetry    int    `toml:"max_retry"`
	BanTime     string `toml:"ban_time"`
}

type Metrics struct {
	Enabled bool `toml:"enabled"`
	Port    int  `toml:"port"`
}
