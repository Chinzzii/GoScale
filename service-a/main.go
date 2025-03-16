package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define custom metrics
var (
	REQUEST_COUNT = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Total HTTP Requests",
		},
		[]string{"method", "status", "path"},
	)

	REQUEST_LATENCY = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP Request Duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status", "path"},
	)

	REQUEST_IN_PROGRESS = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_progress",
			Help: "HTTP Requests in progress",
		},
		[]string{"method", "path"},
	)

	// These system metrics would need to be updated with actual data, e.g., via runtime or OS calls.
	CPU_USAGE = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "process_cpu_usage",
		Help: "Current CPU usage in percent",
	})

	MEMORY_USAGE = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "process_memory_usage_bytes",
		Help: "Current memory usage in bytes",
	})
)

// Register the metrics during initialization
func init() {
	prometheus.MustRegister(REQUEST_COUNT, REQUEST_LATENCY, REQUEST_IN_PROGRESS, CPU_USAGE, MEMORY_USAGE)
}

// Middleware to instrument HTTP requests
func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use the full path (e.g., "/a") and HTTP method for labeling.
		path := c.FullPath()
		if path == "" {
			// fallback in case route is not named
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		// Increment gauge for requests in progress
		REQUEST_IN_PROGRESS.WithLabelValues(method, path).Inc()

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get HTTP status text (e.g., "OK", "Not Found")
		status := http.StatusText(c.Writer.Status())

		// Update metrics
		REQUEST_COUNT.WithLabelValues(method, status, path).Inc()
		REQUEST_LATENCY.WithLabelValues(method, status, path).Observe(duration)
		REQUEST_IN_PROGRESS.WithLabelValues(method, path).Dec()
	}
}

func main() {
	// Create a new Gin router
	r := gin.Default()

	// Use our Prometheus middleware to instrument all requests.
	r.Use(prometheusMiddleware())

	// Sample endpoint returning plain text.
	r.GET("/a", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from Service A")
	})

	// Expose Prometheus metrics at /metrics.
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Start the server on port 8000.
	r.Run(":8000")
}
