package metrics

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func NewPrometheusMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)

			duration := time.Since(start).Seconds()
			path := r.URL.Path

			HttpRequestCountWithPath.WithLabelValues(path).Inc()
			HttpRequestDuration.WithLabelValues(path).Observe(duration)
		})
	}
}
