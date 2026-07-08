package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds all Prometheus metric descriptors for SWU_OSR.
type Metrics struct {
	// HTTPRequestsTotal counts incoming HTTP requests by method, path pattern, and status code.
	HTTPRequestsTotal *prometheus.CounterVec

	// HTTPRequestDuration tracks the latency distribution of HTTP requests.
	HTTPRequestDuration *prometheus.HistogramVec

	// GitHubWebhookEventsTotal counts processed GitHub webhook events by event type.
	GitHubWebhookEventsTotal *prometheus.CounterVec

	// ActiveUsers tracks how many unique users are currently authenticated (approximate).
	ActiveUsersTotal prometheus.Gauge

	// LeaderboardRefreshTotal counts how many times the leaderboard scheduler ran.
	LeaderboardRefreshTotal prometheus.Counter

	// PointsAwardedTotal counts total points awarded across all users.
	PointsAwardedTotal prometheus.Counter
}

// New registers and returns all Prometheus metrics.
func New() *Metrics {
	return &Metrics{
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "swu_osr",
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests partitioned by method, path, and status code.",
			},
			[]string{"method", "path", "status"},
		),

		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "swu_osr",
				Name:      "http_request_duration_seconds",
				Help:      "HTTP request latency distribution in seconds.",
				Buckets:   prometheus.DefBuckets, // 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2.5s, 5s, 10s
			},
			[]string{"method", "path"},
		),

		GitHubWebhookEventsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "swu_osr",
				Name:      "github_webhook_events_total",
				Help:      "Total GitHub webhook events received and processed by event type.",
			},
			[]string{"event_type"},
		),

		ActiveUsersTotal: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "swu_osr",
				Name:      "active_users_total",
				Help:      "Approximate number of active (registered) users in the platform.",
			},
		),

		LeaderboardRefreshTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: "swu_osr",
				Name:      "leaderboard_refresh_total",
				Help:      "Total number of leaderboard refresh cycles completed by the scheduler.",
			},
		),

		PointsAwardedTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: "swu_osr",
				Name:      "points_awarded_total",
				Help:      "Total gamification points awarded across all users since server start.",
			},
		),
	}
}

// Handler returns the Prometheus HTTP handler for the /metrics endpoint.
func Handler() http.Handler {
	return promhttp.Handler()
}
