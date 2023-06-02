// Package prometheus provides useful functions to initialize and populate metrics
package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics structure for storing counter
type Metrics struct {
	ClientHTTPReqs  *prometheus.CounterVec
	BackendHTTPReqs *prometheus.CounterVec
}

var (
	//HTTPCounter Metrics
	HTTPCounter = InitializeHTTPReqCounter(Reg)
	Reg         = prometheus.NewRegistry()
)

// IncrementClientCounterVec increments the counter with method label provided.
func IncrementClientCounterVec(m string) {
	HTTPCounter.ClientHTTPReqs.WithLabelValues(m).Inc()
	return
}

// IncrementBackendCounterVec increments the counter with method label provided.
func IncrementBackendCounterVec(m string) {
	HTTPCounter.BackendHTTPReqs.WithLabelValues(m).Inc()
	return
}

// InitializeHTTPReqCounter inits the httpReqs counter that will be exported.
func InitializeHTTPReqCounter(reg prometheus.Registerer) *Metrics {
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
