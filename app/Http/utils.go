// Package http provides functions to handle incoming HTTP requests
package http

import (
	varnish "http-broadcaster/Varnish"
	"io"
	"log"
	"net/http"
)

// RequestHandler handles requests to broadcast to all varnish instances.
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	method := r.Method
	tag := r.Header.Get("X-Cache-Tags")
	remoteAddr := r.RemoteAddr
	status := varnish.SendToVarnish(method, url, tag)
	if status != "200 Purged" {
		w.WriteHeader(405)
	}
	if tag != "" {
		log.Println(remoteAddr + " Requested ban on X-Cache-Tags : " + tag + " , status: " + status)
	} else {
		log.Println(remoteAddr + " Requested purge on URI :" + url + " , status: " + status)
	}
	io.WriteString(w, status)
}

// HealthHandler handles healthcheck requests and return 200.
func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "OK")
}
