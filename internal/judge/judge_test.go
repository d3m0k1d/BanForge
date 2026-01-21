package judge

import (
	"testing"

	"github.com/d3m0k1d/BanForge/internal/config"
	"github.com/d3m0k1d/BanForge/internal/storage"
)

func TestJudgeLogic(t *testing.T) {
	tests := []struct {
		name      string
		inputRule config.Rule
		inputLog  storage.LogEntry
		wantErr   bool
		wantMatch bool
	}{
		{
			name:      "Empty rule",
			inputRule: config.Rule{Name: "", ServiceName: "", Path: "", Status: "", Method: ""},
			inputLog:  storage.LogEntry{ID: 0, Service: "nginx", IP: "127.0.0.1", Path: "/api", Status: "200", Method: "GET", CreatedAt: ""},
			wantErr:   true,
			wantMatch: false,
		},
		{
			name:      "Matching rule",
			inputRule: config.Rule{Name: "test", ServiceName: "nginx", Path: "/api", Status: "200", Method: "GET"},
			inputLog:  storage.LogEntry{ID: 1, Service: "nginx", IP: "127.0.0.1", Path: "/api", Status: "200", Method: "GET", CreatedAt: ""},
			wantErr:   false,
			wantMatch: true,
		},
		{
			name:      "Non-matching status",
			inputRule: config.Rule{Name: "test", ServiceName: "nginx", Path: "/api", Status: "404", Method: "GET"},
			inputLog:  storage.LogEntry{ID: 2, Service: "nginx", IP: "127.0.0.1", Path: "/api", Status: "200", Method: "GET", CreatedAt: ""},
			wantErr:   false,
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.inputRule.Name == "" {
				if !tt.wantErr {
					t.Errorf("Expected error for empty rule name, but got none")
				}
				return
			}

			result := (tt.inputRule.Method == "" || tt.inputLog.Method == tt.inputRule.Method) &&
				(tt.inputRule.Status == "" || tt.inputLog.Status == tt.inputRule.Status) &&
				(tt.inputRule.Path == "" || tt.inputLog.Path == tt.inputRule.Path) &&
				(tt.inputRule.ServiceName == "" || tt.inputLog.Service == tt.inputRule.ServiceName)

			if result != tt.wantMatch {
				t.Errorf("Expected error: %v, but got: %v", tt.wantErr, result)
			}
		})
	}
}
