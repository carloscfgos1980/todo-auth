package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds all Prometheus metrics for the application
type Metrics struct {
	httpRequestsTotal   prometheus.Counter
	httpRequestDuration prometheus.Histogram
	todosCreated        prometheus.Counter
	todosUpdated        prometheus.Counter
	todosDeleted        prometheus.Counter
	todosFetched        prometheus.Counter
}

// NewMetrics creates and registers all metrics
func NewMetrics() *Metrics {
	m := &Metrics{
		httpRequestsTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
		),
		httpRequestDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
		),
		todosCreated: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "todos_created_total",
				Help: "Total number of todos created",
			},
		),
		todosUpdated: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "todos_updated_total",
				Help: "Total number of todos updated",
			},
		),
		todosDeleted: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "todos_deleted_total",
				Help: "Total number of todos deleted",
			},
		),
		todosFetched: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "todos_fetched_total",
				Help: "Total number of todos fetched",
			},
		),
	}

	// Register all metrics
	prometheus.MustRegister(
		m.httpRequestsTotal,
		m.httpRequestDuration,
		m.todosCreated,
		m.todosUpdated,
		m.todosDeleted,
		m.todosFetched,
	)

	return m
}

// IncrementTodosCreated increments the todos created counter
func (m *Metrics) IncrementTodosCreated() {
	m.todosCreated.Inc()
}

// IncrementTodosUpdated increments the todos updated counter
func (m *Metrics) IncrementTodosUpdated() {
	m.todosUpdated.Inc()
}

// IncrementTodosDeleted increments the todos deleted counter
func (m *Metrics) IncrementTodosDeleted() {
	m.todosDeleted.Inc()
}

// IncrementTodosFetched increments the todos fetched counter
func (m *Metrics) IncrementTodosFetched() {
	m.todosFetched.Inc()
}

// MetricsMiddleware wraps HTTP handlers to track metrics
func MetricsMiddleware(m *Metrics) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			m.httpRequestsTotal.Inc()

			// Call the next handler
			next.ServeHTTP(w, r)

			// Record the duration
			duration := time.Since(start).Seconds()
			m.httpRequestDuration.Observe(duration)
		})
	}
}

// Handler returns an HTTP handler for serving Prometheus metrics
func Handler() http.Handler {
	return promhttp.Handler()
}
