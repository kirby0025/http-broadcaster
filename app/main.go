// Receive http PURGE request and broadcast it to several Varnish servers.
package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	h "http-broadcaster/Http"
	prometheus2 "http-broadcaster/Prometheus"
	"log"
	"net/http"
	"os"
)

func main() {
	logFile, err := os.OpenFile("./log/purge.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
	prometheus2.Reg = prometheus.NewRegistry()
	prometheus2.Reg.MustRegister(prometheus2.HTTPCounter.ClientHTTPReqs)
	prometheus2.Reg.MustRegister(prometheus2.HTTPCounter.BackendHTTPReqs)
	prometheus2.Reg.MustRegister(collectors.NewBuildInfoCollector())
	http.HandleFunc("/", h.RequestHandler)
	http.HandleFunc("/healthcheck", h.HealthHandler)
	http.Handle("/metrics", promhttp.HandlerFor(prometheus2.Reg, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
		Registry:          prometheus2.Reg,
	}))
	log.Fatal(http.ListenAndServe(":6081", nil))
}
