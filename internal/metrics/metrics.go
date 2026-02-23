package metrics

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	metricsMu sync.RWMutex
	metrics   = make(map[string]int64)
)

func IncBan(service string) {
	metricsMu.Lock()
	metrics["ban_count"]++
	metrics[service+"_bans"]++
	metricsMu.Unlock()
}

func IncUnban(service string) {
	metricsMu.Lock()
	metrics["unban_count"]++
	metrics[service+"_unbans"]++
	metricsMu.Unlock()
}

func IncRuleMatched(rule_name string) {
	metricsMu.Lock()
	metrics[rule_name+"_rule_matched"]++
	metricsMu.Unlock()
}

func IncLogParsed() {
	metricsMu.Lock()
	metrics["log_parsed"]++
	metricsMu.Unlock()
}

func IncError() {
	metricsMu.Lock()
	metrics["error_count"]++
	metricsMu.Unlock()
}

func IncBanAttempt(firewall string) {
	metricsMu.Lock()
	metrics["ban_attempt_count"]++
	metrics[firewall+"_ban_attempts"]++
	metricsMu.Unlock()
}

func IncUnbanAttempt(firewall string) {
	metricsMu.Lock()
	metrics["unban_attempt_count"]++
	metrics[firewall+"_unban_attempts"]++
	metricsMu.Unlock()
}

func IncPortOperation(operation string, protocol string) {
	metricsMu.Lock()
	key := "port_" + operation + "_" + protocol
	metrics[key]++
	metricsMu.Unlock()
}

func IncParserEvent(service string) {
	metricsMu.Lock()
	metrics[service+"_parsed_events"]++
	metricsMu.Unlock()
}

func IncScannerEvent(service string) {
	metricsMu.Lock()
	metrics[service+"_scanner_events"]++
	metricsMu.Unlock()
}

func IncDBOperation(operation string, table string) {
	metricsMu.Lock()
	key := "db_" + operation + "_" + table
	metrics[key]++
	metricsMu.Unlock()
}

func IncRequestCount(service string) {
	metricsMu.Lock()
	metrics[service+"_request_count"]++
	metricsMu.Unlock()
}

func MetricsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricsMu.RLock()
		snapshot := make(map[string]int64, len(metrics))
		for k, v := range metrics {
			snapshot[k] = v
		}
		metricsMu.RUnlock()

		w.Header().Set("Content-Type", "text/plain; version=0.0.4")

		for name, value := range snapshot {
			metricName := name + "_total"
			_, _ = fmt.Fprintf(w, "# TYPE %s counter\n", metricName)
			_, _ = fmt.Fprintf(w, "%s %d\n", metricName, value)
		}
	})
}

func StartMetricsServer(port int) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", MetricsHandler())

	server := &http.Server{
		Addr:         "localhost:" + strconv.Itoa(port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Printf("Starting metrics server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("metrics server failed: %w", err)
	}
	return nil
}
