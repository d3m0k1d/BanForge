package actions

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/d3m0k1d/BanForge/internal/config"
)

var defaultClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	},
}

func SendWebhook(action config.Action) error {
	if !action.Enabled {
		return nil
	}
	if action.URL == "" {
		return fmt.Errorf("URL on config is empty")
	}

	method := action.Method
	if method == "" {
		method = "POST"
	}

	var bodyReader io.Reader
	if action.Body != "" {
		bodyReader = strings.NewReader(action.Body)
		if action.Headers["Content-Type"] == "" && action.Headers["content-type"] == "" {
			action.Headers["Content-Type"] = "application/json"
		}
	}

	req, err := http.NewRequest(method, action.URL, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	for key, value := range action.Headers {
		req.Header.Add(key, value)
	}

	// #nosec G704 - HTTP request validation by system administrators
	resp, err := defaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}
