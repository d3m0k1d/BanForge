package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
)

// ============================================
// Tests for SanitizeRuleFilename
// ============================================

func TestSanitizeRuleFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple alphanumeric",
			input:    "nginx404",
			expected: "nginx404",
		},
		{
			name:     "with spaces",
			input:    "nginx 404 error",
			expected: "nginx_404_error",
		},
		{
			name:     "with special chars",
			input:    "nginx/404:error",
			expected: "nginx_404_error",
		},
		{
			name:     "with dashes and underscores",
			input:    "nginx-404_error",
			expected: "nginx-404_error",
		},
		{
			name:     "uppercase to lowercase",
			input:    "NGINX-404",
			expected: "nginx-404",
		},
		{
			name:     "mixed case",
			input:    "Nginx-Admin-Access",
			expected: "nginx-admin-access",
		},
		{
			name:     "with dots",
			input:    "nginx.error.page",
			expected: "nginx_error_page",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only special chars",
			input:    "!@#$%^&*()",
			expected: "__________",
		},
		{
			name:     "russian chars",
			input:    "nginxошибка",
			expected: "nginx______",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeRuleFilename(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeRuleFilename(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// ============================================
// Tests for ParseDurationWithYears
// ============================================

func TestParseDurationWithYears(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    time.Duration
		expectError bool
	}{
		{
			name:        "1 year",
			input:       "1y",
			expected:    365 * 24 * time.Hour,
			expectError: false,
		},
		{
			name:        "2 years",
			input:       "2y",
			expected:    2 * 365 * 24 * time.Hour,
			expectError: false,
		},
		{
			name:        "1 month",
			input:       "1M",
			expected:    30 * 24 * time.Hour,
			expectError: false,
		},
		{
			name:        "6 months",
			input:       "6M",
			expected:    6 * 30 * 24 * time.Hour,
			expectError: false,
		},
		{
			name:        "30 days",
			input:       "30d",
			expected:    30 * 24 * time.Hour,
			expectError: false,
		},
		{
			name:        "1 day",
			input:       "1d",
			expected:    24 * time.Hour,
			expectError: false,
		},
		{
			name:        "1 hour",
			input:       "1h",
			expected:    time.Hour,
			expectError: false,
		},
		{
			name:        "30 minutes",
			input:       "30m",
			expected:    30 * time.Minute,
			expectError: false,
		},
		{
			name:        "30 seconds",
			input:       "30s",
			expected:    30 * time.Second,
			expectError: false,
		},
		{
			name:        "complex duration",
			input:       "1h30m",
			expected:    1*time.Hour + 30*time.Minute,
			expectError: false,
		},
		{
			name:        "invalid year format",
			input:       "abc",
			expected:    0,
			expectError: true,
		},
		{
			name:        "invalid month format",
			input:       "xM",
			expected:    0,
			expectError: true,
		},
		{
			name:        "invalid day format",
			input:       "xd",
			expected:    0,
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    0,
			expectError: true,
		},
		{
			name:        "negative duration",
			input:       "-1h",
			expected:    -1 * time.Hour,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDurationWithYears(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("ParseDurationWithYears(%q) error = %v, expectError %v", tt.input, err, tt.expectError)
				return
			}
			if !tt.expectError && result != tt.expected {
				t.Errorf("ParseDurationWithYears(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// ============================================
// Tests for Rule validation (NewRule)
// ============================================

func TestNewRule_EmptyName(t *testing.T) {
	err := NewRule("", "nginx", "", "", "", "1h", 0)
	if err == nil {
		t.Error("NewRule with empty name should return error")
	}
	if err.Error() != "rule name can't be empty" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestNewRule_DuplicateRule(t *testing.T) {
	// Create temp directory for rules
	tmpDir := t.TempDir()

	// Create first rule
	firstRulePath := filepath.Join(tmpDir, "test-rule.toml")
	cfg := Rules{Rules: []Rule{{Name: "test-rule"}}}
	file, _ := os.Create(firstRulePath)
	toml.NewEncoder(file).Encode(cfg)
	file.Close()

	// Try to create duplicate
	err := NewRule("test-rule", "nginx", "", "", "", "1h", 0)
	if err == nil {
		t.Error("NewRule with duplicate name should return error")
	}
}

// ============================================
// Tests for EditRule validation
// ============================================

func TestEditRule_EmptyName(t *testing.T) {
	err := EditRule("", "nginx", "", "", "")
	if err == nil {
		t.Error("EditRule with empty name should return error")
	}
	if err.Error() != "rule name can't be empty" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

// ============================================
// Tests for Rule struct validation
// ============================================

func TestRule_StructTags(t *testing.T) {
	// Test that Rule struct has correct TOML tags
	rule := Rule{
		Name:        "test-rule",
		ServiceName: "nginx",
		Path:        "/admin/*",
		Status:      "403",
		Method:      "POST",
		MaxRetry:    5,
		BanTime:     "1h",
	}

	// Encode to TOML and verify
	cfg := Rules{Rules: []Rule{rule}}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.toml")

	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		t.Fatalf("Failed to encode rule: %v", err)
	}

	// Read back and verify
	var decoded Rules
	if _, err := toml.DecodeFile(tmpFile, &decoded); err != nil {
		t.Fatalf("Failed to decode rule: %v", err)
	}

	if len(decoded.Rules) != 1 {
		t.Fatalf("Expected 1 rule, got %d", len(decoded.Rules))
	}

	decodedRule := decoded.Rules[0]
	if decodedRule.Name != rule.Name {
		t.Errorf("Name mismatch: %q != %q", decodedRule.Name, rule.Name)
	}
	if decodedRule.ServiceName != rule.ServiceName {
		t.Errorf("ServiceName mismatch: %q != %q", decodedRule.ServiceName, rule.ServiceName)
	}
	if decodedRule.Path != rule.Path {
		t.Errorf("Path mismatch: %q != %q", decodedRule.Path, rule.Path)
	}
	if decodedRule.MaxRetry != rule.MaxRetry {
		t.Errorf("MaxRetry mismatch: %d != %d", decodedRule.MaxRetry, rule.MaxRetry)
	}
	if decodedRule.BanTime != rule.BanTime {
		t.Errorf("BanTime mismatch: %q != %q", decodedRule.BanTime, rule.BanTime)
	}
}

// ============================================
// Tests for Action struct validation
// ============================================

func TestAction_EmailAction(t *testing.T) {
	action := Action{
		Type:         "email",
		Enabled:      true,
		Email:        "admin@example.com",
		EmailSender:  "banforge@example.com",
		EmailSubject: "Alert",
		SMTPHost:     "smtp.example.com",
		SMTPPort:     587,
		SMTPUser:     "user",
		SMTPTLS:      true,
		Body:         "IP {ip} banned",
	}

	cfg := Rules{Rules: []Rule{{
		Name:   "test",
		Action: []Action{action},
	}}}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "action.toml")

	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		t.Fatalf("Failed to encode action: %v", err)
	}

	// Read back and verify
	var decoded Rules
	if _, err := toml.DecodeFile(tmpFile, &decoded); err != nil {
		t.Fatalf("Failed to decode action: %v", err)
	}

	decodedAction := decoded.Rules[0].Action[0]
	if decodedAction.Type != "email" {
		t.Errorf("Expected action type 'email', got %q", decodedAction.Type)
	}
	if decodedAction.Email != "admin@example.com" {
		t.Errorf("Expected email 'admin@example.com', got %q", decodedAction.Email)
	}
	if decodedAction.SMTPPort != 587 {
		t.Errorf("Expected SMTP port 587, got %d", decodedAction.SMTPPort)
	}
}

func TestAction_WebhookAction(t *testing.T) {
	action := Action{
		Type:    "webhook",
		Enabled: true,
		URL:     "https://hooks.example.com/alert",
		Method:  "POST",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer token123",
		},
		Body: `{"ip": "{ip}", "rule": "{rule}"}`,
	}

	cfg := Rules{Rules: []Rule{{
		Name:   "test",
		Action: []Action{action},
	}}}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "webhook.toml")

	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		t.Fatalf("Failed to encode webhook action: %v", err)
	}

	// Read back and verify
	var decoded Rules
	if _, err := toml.DecodeFile(tmpFile, &decoded); err != nil {
		t.Fatalf("Failed to decode webhook action: %v", err)
	}

	decodedAction := decoded.Rules[0].Action[0]
	if decodedAction.Type != "webhook" {
		t.Errorf("Expected action type 'webhook', got %q", decodedAction.Type)
	}
	if decodedAction.URL != "https://hooks.example.com/alert" {
		t.Errorf("Expected URL 'https://hooks.example.com/alert', got %q", decodedAction.URL)
	}
	if decodedAction.Headers["Content-Type"] != "application/json" {
		t.Errorf("Expected Content-Type header, got %q", decodedAction.Headers["Content-Type"])
	}
}

func TestAction_ScriptAction(t *testing.T) {
	action := Action{
		Type:          "script",
		Enabled:       true,
		Script:        "/usr/local/bin/notify.sh",
		Interpretator: "bash",
	}

	cfg := Rules{Rules: []Rule{{
		Name:   "test",
		Action: []Action{action},
	}}}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "script.toml")

	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		t.Fatalf("Failed to encode script action: %v", err)
	}

	// Read back and verify
	var decoded Rules
	if _, err := toml.DecodeFile(tmpFile, &decoded); err != nil {
		t.Fatalf("Failed to decode script action: %v", err)
	}

	decodedAction := decoded.Rules[0].Action[0]
	if decodedAction.Type != "script" {
		t.Errorf("Expected action type 'script', got %q", decodedAction.Type)
	}
	if decodedAction.Script != "/usr/local/bin/notify.sh" {
		t.Errorf("Expected script path '/usr/local/bin/notify.sh', got %q", decodedAction.Script)
	}
	if decodedAction.Interpretator != "bash" {
		t.Errorf("Expected interpretator 'bash', got %q", decodedAction.Interpretator)
	}
}

// ============================================
// Tests for Service struct validation
// ============================================

func TestService_FileLogging(t *testing.T) {
	service := Service{
		Name:    "nginx",
		Logging: "file",
		LogPath: "/var/log/nginx/access.log",
		Enabled: true,
	}

	cfg := Config{Service: []Service{service}}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "service.toml")

	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		t.Fatalf("Failed to encode service: %v", err)
	}

	// Read back and verify
	var decoded Config
	if _, err := toml.DecodeFile(tmpFile, &decoded); err != nil {
		t.Fatalf("Failed to decode service: %v", err)
	}

	if len(decoded.Service) != 1 {
		t.Fatalf("Expected 1 service, got %d", len(decoded.Service))
	}

	decodedService := decoded.Service[0]
	if decodedService.Name != "nginx" {
		t.Errorf("Expected service name 'nginx', got %q", decodedService.Name)
	}
	if decodedService.Logging != "file" {
		t.Errorf("Expected logging type 'file', got %q", decodedService.Logging)
	}
	if decodedService.LogPath != "/var/log/nginx/access.log" {
		t.Errorf("Expected log path '/var/log/nginx/access.log', got %q", decodedService.LogPath)
	}
}

func TestService_JournaldLogging(t *testing.T) {
	service := Service{
		Name:    "sshd",
		Logging: "journald",
		LogPath: "sshd",
		Enabled: true,
	}

	cfg := Config{Service: []Service{service}}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "journald.toml")

	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		t.Fatalf("Failed to encode journald service: %v", err)
	}

	// Read back and verify
	var decoded Config
	if _, err := toml.DecodeFile(tmpFile, &decoded); err != nil {
		t.Fatalf("Failed to decode journald service: %v", err)
	}

	decodedService := decoded.Service[0]
	if decodedService.Logging != "journald" {
		t.Errorf("Expected logging type 'journald', got %q", decodedService.Logging)
	}
}

// ============================================
// Tests for Metrics struct validation
// ============================================

func TestMetrics_Enabled(t *testing.T) {
	metrics := Metrics{
		Enabled: true,
		Port:    9090,
	}

	cfg := Config{Metrics: metrics}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "metrics.toml")

	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		t.Fatalf("Failed to encode metrics: %v", err)
	}

	// Read back and verify
	var decoded Config
	if _, err := toml.DecodeFile(tmpFile, &decoded); err != nil {
		t.Fatalf("Failed to decode metrics: %v", err)
	}

	if !decoded.Metrics.Enabled {
		t.Error("Expected metrics to be enabled")
	}
	if decoded.Metrics.Port != 9090 {
		t.Errorf("Expected metrics port 9090, got %d", decoded.Metrics.Port)
	}
}

func TestMetrics_Disabled(t *testing.T) {
	metrics := Metrics{
		Enabled: false,
		Port:    0,
	}

	cfg := Config{Metrics: metrics}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "metrics-disabled.toml")

	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		t.Fatalf("Failed to encode disabled metrics: %v", err)
	}

	// Read back and verify
	var decoded Config
	if _, err := toml.DecodeFile(tmpFile, &decoded); err != nil {
		t.Fatalf("Failed to decode disabled metrics: %v", err)
	}

	if decoded.Metrics.Enabled {
		t.Error("Expected metrics to be disabled")
	}
}

// ============================================
// Tests for Firewall struct validation
// ============================================

func TestFirewall_Nftables(t *testing.T) {
	firewall := Firewall{
		Name:   "nftables",
		Config: "/etc/nftables.conf",
	}

	cfg := Config{Firewall: firewall}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "firewall.toml")

	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		t.Fatalf("Failed to encode firewall: %v", err)
	}

	// Read back and verify
	var decoded Config
	if _, err := toml.DecodeFile(tmpFile, &decoded); err != nil {
		t.Fatalf("Failed to decode firewall: %v", err)
	}

	if decoded.Firewall.Name != "nftables" {
		t.Errorf("Expected firewall name 'nftables', got %q", decoded.Firewall.Name)
	}
	if decoded.Firewall.Config != "/etc/nftables.conf" {
		t.Errorf("Expected firewall config '/etc/nftables.conf', got %q", decoded.Firewall.Config)
	}
}

func TestFirewall_Iptables(t *testing.T) {
	firewall := Firewall{
		Name:   "iptables",
		Config: "",
	}

	cfg := Config{Firewall: firewall}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "iptables.toml")

	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		t.Fatalf("Failed to encode iptables: %v", err)
	}

	// Read back and verify
	var decoded Config
	if _, err := toml.DecodeFile(tmpFile, &decoded); err != nil {
		t.Fatalf("Failed to decode iptables: %v", err)
	}

	if decoded.Firewall.Name != "iptables" {
		t.Errorf("Expected firewall name 'iptables', got %q", decoded.Firewall.Name)
	}
}

func TestFirewall_Ufw(t *testing.T) {
	firewall := Firewall{
		Name:   "ufw",
		Config: "",
	}

	cfg := Config{Firewall: firewall}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "ufw.toml")

	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		t.Fatalf("Failed to encode ufw: %v", err)
	}

	// Read back and verify
	var decoded Config
	if _, err := toml.DecodeFile(tmpFile, &decoded); err != nil {
		t.Fatalf("Failed to decode ufw: %v", err)
	}

	if decoded.Firewall.Name != "ufw" {
		t.Errorf("Expected firewall name 'ufw', got %q", decoded.Firewall.Name)
	}
}

func TestFirewall_Firewalld(t *testing.T) {
	firewall := Firewall{
		Name:   "firewalld",
		Config: "",
	}

	cfg := Config{Firewall: firewall}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "firewalld.toml")

	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		t.Fatalf("Failed to encode firewalld: %v", err)
	}

	// Read back and verify
	var decoded Config
	if _, err := toml.DecodeFile(tmpFile, &decoded); err != nil {
		t.Fatalf("Failed to decode firewalld: %v", err)
	}

	if decoded.Firewall.Name != "firewalld" {
		t.Errorf("Expected firewall name 'firewalld', got %q", decoded.Firewall.Name)
	}
}

// ============================================
// Integration test: Full config round-trip
// ============================================

func TestConfig_FullRoundTrip(t *testing.T) {
	fullConfig := Config{
		Firewall: Firewall{
			Name:   "nftables",
			Config: "/etc/nftables.conf",
		},
		Metrics: Metrics{
			Enabled: true,
			Port:    9090,
		},
		Service: []Service{
			{
				Name:    "nginx",
				Logging: "file",
				LogPath: "/var/log/nginx/access.log",
				Enabled: true,
			},
			{
				Name:    "sshd",
				Logging: "journald",
				LogPath: "sshd",
				Enabled: true,
			},
		},
	}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "full-config.toml")

	// Encode
	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(fullConfig); err != nil {
		t.Fatalf("Failed to encode full config: %v", err)
	}

	// Decode
	var decoded Config
	if _, err := toml.DecodeFile(tmpFile, &decoded); err != nil {
		t.Fatalf("Failed to decode full config: %v", err)
	}

	// Verify firewall
	if decoded.Firewall.Name != fullConfig.Firewall.Name {
		t.Errorf("Firewall.Name mismatch: %q != %q", decoded.Firewall.Name, fullConfig.Firewall.Name)
	}

	// Verify metrics
	if decoded.Metrics.Enabled != fullConfig.Metrics.Enabled {
		t.Errorf("Metrics.Enabled mismatch: %v != %v", decoded.Metrics.Enabled, fullConfig.Metrics.Enabled)
	}
	if decoded.Metrics.Port != fullConfig.Metrics.Port {
		t.Errorf("Metrics.Port mismatch: %d != %d", decoded.Metrics.Port, fullConfig.Metrics.Port)
	}

	// Verify services
	if len(decoded.Service) != len(fullConfig.Service) {
		t.Fatalf("Services count mismatch: %d != %d", len(decoded.Service), len(fullConfig.Service))
	}

	for i, expected := range fullConfig.Service {
		actual := decoded.Service[i]
		if actual.Name != expected.Name {
			t.Errorf("Service[%d].Name mismatch: %q != %q", i, actual.Name, expected.Name)
		}
		if actual.Logging != expected.Logging {
			t.Errorf("Service[%d].Logging mismatch: %q != %q", i, actual.Logging, expected.Logging)
		}
		if actual.Enabled != expected.Enabled {
			t.Errorf("Service[%d].Enabled mismatch: %v != %v", i, actual.Enabled, expected.Enabled)
		}
	}
}

// ============================================
// Integration test: Full rule with actions round-trip
// ============================================

func TestRule_FullRoundTrip(t *testing.T) {
	fullRule := Rules{
		Rules: []Rule{
			{
				Name:        "nginx-bruteforce",
				ServiceName: "nginx",
				Path:        "/admin/*",
				Status:      "403",
				Method:      "POST",
				MaxRetry:    5,
				BanTime:     "2h",
				Action: []Action{
					{
						Type:         "email",
						Enabled:      true,
						Email:        "admin@example.com",
						EmailSender:  "banforge@example.com",
						EmailSubject: "Ban Alert",
						SMTPHost:     "smtp.example.com",
						SMTPPort:     587,
						SMTPUser:     "user",
						SMTPPassword: "pass",
						SMTPTLS:      true,
						Body:         "IP {ip} banned for rule {rule}",
					},
					{
						Type:    "webhook",
						Enabled: true,
						URL:     "https://hooks.slack.com/services/xxx",
						Method:  "POST",
						Headers: map[string]string{
							"Content-Type": "application/json",
						},
						Body: `{"text": "IP {ip} banned"}`,
					},
					{
						Type:          "script",
						Enabled:       true,
						Script:        "/usr/local/bin/notify.sh",
						Interpretator: "bash",
					},
				},
			},
		},
	}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "full-rule.toml")

	// Encode
	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(fullRule); err != nil {
		t.Fatalf("Failed to encode full rule: %v", err)
	}

	// Decode
	var decoded Rules
	if _, err := toml.DecodeFile(tmpFile, &decoded); err != nil {
		t.Fatalf("Failed to decode full rule: %v", err)
	}

	// Verify rule
	if len(decoded.Rules) != 1 {
		t.Fatalf("Expected 1 rule, got %d", len(decoded.Rules))
	}

	rule := decoded.Rules[0]
	if rule.Name != fullRule.Rules[0].Name {
		t.Errorf("Rule.Name mismatch: %q != %q", rule.Name, fullRule.Rules[0].Name)
	}
	if rule.ServiceName != fullRule.Rules[0].ServiceName {
		t.Errorf("Rule.ServiceName mismatch: %q != %q", rule.ServiceName, fullRule.Rules[0].ServiceName)
	}
	if rule.MaxRetry != fullRule.Rules[0].MaxRetry {
		t.Errorf("Rule.MaxRetry mismatch: %d != %d", rule.MaxRetry, fullRule.Rules[0].MaxRetry)
	}

	// Verify actions
	if len(rule.Action) != 3 {
		t.Fatalf("Expected 3 actions, got %d", len(rule.Action))
	}

	// Email action
	emailAction := rule.Action[0]
	if emailAction.Type != "email" {
		t.Errorf("Action[0].Type mismatch: %q != 'email'", emailAction.Type)
	}
	if emailAction.Email != "admin@example.com" {
		t.Errorf("Email action email mismatch: %q != 'admin@example.com'", emailAction.Email)
	}

	// Webhook action
	webhookAction := rule.Action[1]
	if webhookAction.Type != "webhook" {
		t.Errorf("Action[1].Type mismatch: %q != 'webhook'", webhookAction.Type)
	}
	if webhookAction.Headers["Content-Type"] != "application/json" {
		t.Errorf("Webhook Content-Type header mismatch")
	}

	// Script action
	scriptAction := rule.Action[2]
	if scriptAction.Type != "script" {
		t.Errorf("Action[2].Type mismatch: %q != 'script'", scriptAction.Type)
	}
	if scriptAction.Script != "/usr/local/bin/notify.sh" {
		t.Errorf("Script action script mismatch: %q != '/usr/local/bin/notify.sh'", scriptAction.Script)
	}
}
