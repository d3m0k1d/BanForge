package actions

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/d3m0k1d/BanForge/internal/config"
)

func TestSendWebhook_BodyWithoutHeaders(t *testing.T) {
	var gotContentType string
	var gotBody string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		buf := make([]byte, 1024)
		n, _ := r.Body.Read(buf)
		gotBody = string(buf[:n])
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	action := config.Action{
		Type:    "webhook",
		Enabled: true,
		URL:     server.URL,
		Method:  "POST",
		Body:    `{"test": true}`,
	}

	if err := SendWebhook(action); err != nil {
		t.Fatalf("SendWebhook() unexpected error: %v", err)
	}

	if gotContentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", gotContentType, "application/json")
	}
	if gotBody != `{"test": true}` {
		t.Errorf("body = %q, want %q", gotBody, `{"test": true}`)
	}
}

func TestSendWebhook_NoBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	action := config.Action{
		Type:    "webhook",
		Enabled: true,
		URL:     server.URL,
		Method:  "POST",
	}

	if err := SendWebhook(action); err != nil {
		t.Fatalf("SendWebhook() unexpected error: %v", err)
	}
}

func TestSendWebhook_Validation(t *testing.T) {
	tests := []struct {
		name    string
		action  config.Action
		wantErr bool
	}{
		{
			name: "disabled action",
			action: config.Action{
				Type:    "webhook",
				Enabled: false,
				URL:     "",
			},
			wantErr: false,
		},
		{
			name: "empty URL",
			action: config.Action{
				Type:    "webhook",
				Enabled: true,
				URL:     "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendWebhook(tt.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendWebhook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
