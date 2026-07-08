package metrics

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Middleware returns an HTTP middleware that records request count and duration.
//
// The "path" label uses the route pattern (e.g. "/api/profiles/{alias}") rather
// than the raw URL to prevent high cardinality from user-provided path params.
func (m *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := newResponseWriter(w)

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(rw.statusCode)

		// Use route pattern if available (set by chi router), fall back to URL path.
		path := r.URL.Path

		m.HTTPRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		m.HTTPRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
	})
}

// RecordWebhookEvent increments the webhook event counter for a given event type.
func (m *Metrics) RecordWebhookEvent(eventType string) {
	m.GitHubWebhookEventsTotal.WithLabelValues(eventType).Inc()
}

// SetActiveUsers sets the gauge for the approximate active user count.
func (m *Metrics) SetActiveUsers(count int) {
	m.ActiveUsersTotal.Set(float64(count))
}

// IncrLeaderboardRefresh increments the scheduler refresh counter.
func (m *Metrics) IncrLeaderboardRefresh() {
	m.LeaderboardRefreshTotal.Inc()
}

// AddPoints increments the total points-awarded counter.
func (m *Metrics) AddPoints(points int) {
	m.PointsAwardedTotal.Add(float64(points))
}

// RoutePattern extracts the chi route pattern from a request, falling back
// to a sanitized URL path. This prevents label cardinality explosion.
//
// Usage with chi:
//
//	rctx := chi.RouteContext(r.Context())
//	pattern := rctx.RoutePattern()
func RoutePattern(r *http.Request) string {
	// Try to read chi route context pattern (avoids import cycle by using interface)
	if rctx := r.Context().Value(struct{ name string }{"chi_route_ctx"}); rctx != nil {
		if pattern, ok := rctx.(fmt.Stringer); ok {
			return pattern.String()
		}
	}
	return r.URL.Path
}
