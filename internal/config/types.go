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
	Action      []Action
}

type Metrics struct {
	Enabled bool `toml:"enabled"`
	Port    int  `toml:"port"`
}

// Actions
type Action struct {
	Type          string            `toml:"type"`
	Enabled       bool              `toml:"enabled"`
	URL           string            `toml:"url"`
	Method        string            `toml:"method"`
	Headers       map[string]string `toml:"headers"`
	Body          string            `toml:"body"`
	Email         string            `toml:"email"`
	EmailSender   string            `toml:"email_sender"`
	EmailSubject  string            `toml:"email_subject"`
	SMTPHost      string            `toml:"smtp_host"`
	SMTPPort      int               `toml:"smtp_port"`
	SMTPUser      string            `toml:"smtp_user"`
	SMTPPassword  string            `toml:"smtp_password"`
	SMTPTLS       bool              `toml:"smtp_tls"`
	Interpretator string            `toml:"interpretator"`
	Script        string            `toml:"script"`
}
