package config

type Firewall struct {
	Name    string `toml:"name"`
	Config  string `toml:"config"`
	BanTime int    `toml:"ban_time"`
}

type Service struct {
	Name     string `toml:"name"`
	Log_path string `toml:"log_path"`
	Enabled  bool   `toml:"enabled"`
}
