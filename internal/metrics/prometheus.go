package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "HTTP request latency",
		},
		[]string{"method", "path"},
	)

	ReservationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "reservations_total",
			Help: "Total reservation attempts",
		},
		[]string{"result"},
	)

	CacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of Redis cache hits",
		},
	)

	CacheMiss = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_miss_total",
			Help: "Total number of Redis cache misses",
		},
	)
)

func Register() {

	prometheus.MustRegister(
		HTTPRequestsTotal,
		HTTPRequestDuration,
		ReservationsTotal,
		CacheHits,
		CacheMiss,
	)

}
