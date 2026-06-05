package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type MetricsCollector struct {
	mu             sync.RWMutex
	totalRequests  uint64
	totalLatencyMs uint64
	byRoute        map[string]uint64
	byStatus       map[int]uint64
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		byRoute:  make(map[string]uint64),
		byStatus: make(map[int]uint64),
	}
}

func (m *MetricsCollector) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		key := c.Request.Method + " " + route
		status := c.Writer.Status()
		latencyMs := uint64(time.Since(start).Milliseconds())

		m.mu.Lock()
		m.totalRequests++
		m.totalLatencyMs += latencyMs
		m.byRoute[key]++
		m.byStatus[status]++
		m.mu.Unlock()
	}
}

func (m *MetricsCollector) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.mu.RLock()
		totalRequests := m.totalRequests
		totalLatencyMs := m.totalLatencyMs
		byRoute := make(map[string]uint64, len(m.byRoute))
		for k, v := range m.byRoute {
			byRoute[k] = v
		}
		byStatus := make(map[int]uint64, len(m.byStatus))
		for k, v := range m.byStatus {
			byStatus[k] = v
		}
		m.mu.RUnlock()

		avgLatencyMs := float64(0)
		if totalRequests > 0 {
			avgLatencyMs = float64(totalLatencyMs) / float64(totalRequests)
		}

		c.JSON(http.StatusOK, gin.H{
			"total_requests":   totalRequests,
			"average_latency_ms": avgLatencyMs,
			"requests_by_route": byRoute,
			"requests_by_status": byStatus,
		})
	}
}