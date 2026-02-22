package metrics

import (
	"net/http"
)

func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		for k, v := range metrics {
			w.Write([]byte(k + " " + string(v) + "\n"))
		}
	})
}
