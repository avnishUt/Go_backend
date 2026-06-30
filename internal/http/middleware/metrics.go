package middleware

import (
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type MetricsCollector struct {
	mu             sync.Mutex
	startedAt      time.Time
	totalRequests  int64
	statusCounts   map[int]int64
	methodCounts   map[string]int64
	totalLatencyMS int64
	inFlight       int64
}

type MetricsSnapshot struct {
	StartedAt     time.Time        `json:"started_at"`
	UptimeSeconds int64            `json:"uptime_seconds"`
	TotalRequests int64            `json:"total_requests"`
	StatusCounts  map[string]int64 `json:"status_counts"`
	MethodCounts  map[string]int64 `json:"method_counts"`
	AvgLatencyMS  float64          `json:"avg_latency_ms"`
	InFlight      int64            `json:"in_flight"`
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		startedAt:    time.Now(),
		statusCounts: map[int]int64{},
		methodCounts: map[string]int64{},
	}
}

func Metrics(collector *MetricsCollector) gin.HandlerFunc {
	return func(c *gin.Context) {
		if collector == nil {
			c.Next()
			return
		}

		startedAt := time.Now()
		collector.begin()
		c.Next()
		collector.finish(c.Request.Method, c.Writer.Status(), time.Since(startedAt))
	}
}

func (m *MetricsCollector) Snapshot() MetricsSnapshot {
	m.mu.Lock()
	defer m.mu.Unlock()

	statusCounts := map[string]int64{}
	for status, count := range m.statusCounts {
		statusCounts[strconv.Itoa(status)] = count
	}

	methodCounts := map[string]int64{}
	for method, count := range m.methodCounts {
		methodCounts[method] = count
	}

	var avgLatency float64
	if m.totalRequests > 0 {
		avgLatency = float64(m.totalLatencyMS) / float64(m.totalRequests)
	}

	return MetricsSnapshot{
		StartedAt:     m.startedAt.UTC(),
		UptimeSeconds: int64(time.Since(m.startedAt).Seconds()),
		TotalRequests: m.totalRequests,
		StatusCounts:  statusCounts,
		MethodCounts:  methodCounts,
		AvgLatencyMS:  avgLatency,
		InFlight:      m.inFlight,
	}
}

func (m *MetricsCollector) begin() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.inFlight++
}

func (m *MetricsCollector) finish(method string, status int, latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.inFlight--
	m.totalRequests++
	m.statusCounts[status]++
	m.methodCounts[method]++
	m.totalLatencyMS += latency.Milliseconds()
}
