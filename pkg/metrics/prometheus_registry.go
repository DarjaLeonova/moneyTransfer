package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HttpRequestCountWithPath = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests.",
		},
		[]string{"url"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Response time of HTTP request.",
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(HttpRequestCountWithPath)
	prometheus.MustRegister(HttpRequestDuration)
}
