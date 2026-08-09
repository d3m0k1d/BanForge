package actions

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/d3m0k1d/BanForge/internal/config"
)

func TestRunScript_Validation(t *testing.T) {
	tests := []struct {
		name    string
		action  config.Action
		wantErr bool
		errMsg  string
	}{
		{
			name: "disabled action",
			action: config.Action{
				Type:    "script",
				Enabled: false,
				Script:  "/bin/true",
			},
			wantErr: false,
		},
		{
			name: "empty script",
			action: config.Action{
				Type:    "script",
				Enabled: true,
				Script:  "",
			},
			wantErr: true,
			errMsg:  "script on config is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RunScript(tt.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunScript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("RunScript() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestRunScript_DirectExec(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("direct exec is POSIX only")
	}

	script := filepath.Join(t.TempDir(), "test.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	action := config.Action{
		Type:    "script",
		Enabled: true,
		Script:  script,
	}

	if err := RunScript(action); err != nil {
		t.Fatalf("RunScript() unexpected error: %v", err)
	}
}

func TestRunScript_WithInterpretator(t *testing.T) {
	script := filepath.Join(t.TempDir(), "test.sh")
	if err := os.WriteFile(script, []byte("exit 0\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	action := config.Action{
		Type:          "script",
		Enabled:       true,
		Interpretator: "sh",
		Script:        script,
	}

	if err := RunScript(action); err != nil {
		t.Fatalf("RunScript() unexpected error: %v", err)
	}
}

func TestRunScript_DirectExecFails(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("direct exec is POSIX only")
	}

	action := config.Action{
		Type:    "script",
		Enabled: true,
		Script:  "/nonexistent/script.sh",
	}

	if err := RunScript(action); err == nil {
		t.Error("RunScript() expected error for missing script")
	}
}
