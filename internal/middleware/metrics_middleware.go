package middleware

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

// MetricsCollector is a struct that collects metrics for HTTP requests, including total requests, total latency in milliseconds, requests by route, and requests by status code. It uses a mutex to ensure thread-safe access to the metrics data.
type MetricsCollector struct {
	totalRequests  atomic.Uint64
	totalLatencyMs atomic.Uint64
	byRoute        sync.Map // map[string]*atomic.Uint64
	byStatus       sync.Map // map[int]*atomic.Uint64
}

// NewMetricsCollector creates and returns a new instance of MetricsCollector with initialized fields.
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{}
}

// Middleware is a Gin middleware function that collects metrics for each incoming HTTP request. It records the start time of the request, processes the request, and then calculates the latency and updates the metrics for total requests, total latency, requests by route, and requests by status code. It uses the request method and route as the key for tracking requests by route, and the response status code for tracking requests by status code. If the route is not matched, it uses "unmatched" as the key for tracking requests by route.
func (m *MetricsCollector) Middleware() gin.HandlerFunc {
	// Return a Gin handler function that will be called for each incoming HTTP request. This function will record the start time of the request, process the request, and then calculate the latency and update the metrics accordingly.
	return func(c *gin.Context) {
		// Record the start time of the request to calculate latency later.
		start := time.Now()
		// Process the request by calling the next handler in the chain.
		c.Next()
		// Get the full path of the route that was matched for the request. If no route was matched, use "unmatched" as the key for tracking requests by route.
		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		//	Create a key for tracking requests by route using the request method and route. Get the response status code and calculate the latency in milliseconds. Then, update the total requests, total latency, requests by route, and requests by status code metrics using atomic operations and a helper function to increment the counters in the sync.Map.
		key := c.Request.Method + " " + route
		status := c.Writer.Status()
		latencyMs := uint64(time.Since(start).Milliseconds())
		// Increment the total requests and total latency metrics using atomic operations.
		m.totalRequests.Add(1)
		m.totalLatencyMs.Add(latencyMs)
		incrementCounter(&m.byRoute, key)
		incrementCounter(&m.byStatus, status)
	}
}

// Handler is a Gin handler function that returns the collected metrics as a JSON response. It calculates the total requests, average latency, requests by route, and requests by status code, and sends this information in the response.
func (m *MetricsCollector) Handler() gin.HandlerFunc {
	// Return a Gin handler function that will be called when the /metrics endpoint is accessed. This function will gather the collected metrics and return them as a JSON response.
	return func(c *gin.Context) {
		// Load the total requests and total latency metrics using atomic operations. Then, iterate over the byRoute and byStatus sync.Maps to gather the counts for each route and status code. Finally, calculate the average latency and return all the metrics in a JSON response.
		totalRequests := m.totalRequests.Load()
		totalLatencyMs := m.totalLatencyMs.Load()
		// Create maps to hold the counts for requests by route and requests by status code. Use the Range method of sync.Map to iterate over the entries in byRoute and byStatus, and load the counts using atomic operations.
		byRoute := make(map[string]uint64)
		m.byRoute.Range(func(k, v any) bool {
			key, ok := k.(string)
			if !ok {
				return true
			}
			counter, ok := v.(*atomic.Uint64)
			if !ok {
				return true
			}
			byRoute[key] = counter.Load()
			return true
		})
		// Create a map to hold the counts for requests by status code. Use the Range method of sync.Map to iterate over the entries in byStatus, and load the counts using atomic operations.
		byStatus := make(map[int]uint64)
		m.byStatus.Range(func(k, v any) bool {
			statusCode, ok := k.(int)
			if !ok {
				return true
			}
			counter, ok := v.(*atomic.Uint64)
			if !ok {
				return true
			}
			byStatus[statusCode] = counter.Load()
			return true
		})
		// Calculate the average latency in milliseconds by dividing the total latency by the total number of requests. If there are no requests, set the average latency to 0 to avoid division by zero.
		avgLatencyMs := float64(0)
		if totalRequests > 0 {
			avgLatencyMs = float64(totalLatencyMs) / float64(totalRequests)
		}
		// Return the collected metrics in a JSON response with fields for total requests, average latency, requests by route, and requests by status code.
		c.JSON(http.StatusOK, gin.H{
			"total_requests":     totalRequests,
			"average_latency_ms": avgLatencyMs,
			"requests_by_route":  byRoute,
			"requests_by_status": byStatus,
		})
	}
}

// incrementCounter is a helper function that increments the counter for a given key in a sync.Map. It uses the LoadOrStore method to get or create an atomic.Uint64 counter for the key, and then increments the counter using the Add method.
func incrementCounter(counterMap *sync.Map, key any) {
	value, _ := counterMap.LoadOrStore(key, &atomic.Uint64{})
	counter, ok := value.(*atomic.Uint64)
	if !ok {
		return
	}
	counter.Add(1)
}
