// Package prometheus provides useful functions to initialize and populate metrics
package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics structure for storing counter
type Metrics struct {
	ClientHTTPReqs  *prometheus.CounterVec
	BackendHTTPReqs *prometheus.CounterVec
}

var (
	// HTTPCounter Metrics
	HTTPCounter = initializeHTTPReqCounter()
	// Reg is a custom registry to have more control on metrics exported.
	Reg = prometheus.NewRegistry()
	// MetricsEnabled is the flag used to enable the prometheus metrics backend.
	MetricsEnabled = false
)

// InitMetrics enable the metrics functionality if the flags is passed as an argument
func InitMetrics(m bool) {
	if m {
		MetricsEnabled = true
		initPrometheusRegistry()
		// Define custom promhttp handler that expose just our custom registry.
		http.Handle("/metrics", promhttp.HandlerFor(Reg, promhttp.HandlerOpts{
			EnableOpenMetrics: true,
			Registry:          Reg,
		}))
	}
}

// InitPrometheusRegistry initialize registry and counters if metrics flag pass as argument.
func initPrometheusRegistry() {
	// We use a custom registry to better now what metrics are exposed.
	Reg = prometheus.NewRegistry()
	Reg.MustRegister(HTTPCounter.ClientHTTPReqs)
	Reg.MustRegister(HTTPCounter.BackendHTTPReqs)
	Reg.MustRegister(collectors.NewBuildInfoCollector())
}

// IncrementClientCounterVec increments the counter with method label provided.
func IncrementClientCounterVec(m string) {
	HTTPCounter.ClientHTTPReqs.WithLabelValues(m).Inc()
}

// IncrementBackendCounterVec increments the counter with method label provided.
func IncrementBackendCounterVec(m string) {
	HTTPCounter.BackendHTTPReqs.WithLabelValues(m).Inc()
}

// InitializeHTTPReqCounter inits the httpReqs counter that will be exported.
func initializeHTTPReqCounter() *Metrics {
	HTTPCounters := &Metrics{
		ClientHTTPReqs: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "client_http_requests_total",
			Help: "How many HTTP requests processed, partitioned by HTTP method.",
		},
			[]string{"method"},
		),
		BackendHTTPReqs: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "backend_http_requests_total",
			Help: "How many HTTP requests sent to backend, partitioned by HTTP method.",
		},
			[]string{"method"},
		),
	}
	return HTTPCounters
}
