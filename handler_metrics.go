package main

import (
	"net/http"
	"sync/atomic"
	"time"
)

// handlerMetrics returns basic runtime metrics for the API process.
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	uptimeSeconds := int64(time.Since(cfg.startedAt).Seconds())
	requestsTotal := atomic.LoadUint64(&cfg.requestsTotal)

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"uptime_seconds": uptimeSeconds,
		"requests_total": requestsTotal,
	})
}
