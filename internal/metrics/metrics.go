package metrics

import (
	"fmt"
	"log"
	"net/http"
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

func StartMetricsServer(port int) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", MetricsHandler())

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Printf("Starting metrics server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("Metrics server error: %v", err)
	}
}
