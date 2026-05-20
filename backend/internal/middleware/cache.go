package middleware

import (
	"fmt"
	"net/http"
)

// CacheControl returns a middleware that sets Cache-Control headers on responses.
// Performance: Enables browser and CDN caching for stable endpoints, reducing
// server load by allowing clients to serve responses from local cache.
func CacheControl(maxAgeSeconds int) func(http.Handler) http.Handler {
	value := fmt.Sprintf("public, max-age=%d", maxAgeSeconds)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", value)
			next.ServeHTTP(w, r)
		})
	}
}

// NoCache returns a middleware that disables caching for dynamic/authenticated endpoints.
func NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		next.ServeHTTP(w, r)
	})
}
