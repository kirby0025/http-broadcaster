// Receive http PURGE request and broadcast it to several Varnish servers.
package main

import (
	h "http-broadcaster/Http"
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
	http.HandleFunc("/", h.RequestHandler)
	http.HandleFunc("/healthcheck", h.HealthHandler)
	log.Fatal(http.ListenAndServe(":6081", nil))
}
