package main

import (
	"net/http"
	"sync/atomic"
)

// metricsMiddleware increments request counters for every incoming HTTP request.
func (cfg *apiConfig) metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&cfg.requestsTotal, 1)
		next.ServeHTTP(w, r)
	})
}
