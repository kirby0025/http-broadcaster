// Package http provides functions to handle incoming HTTP requests
package http

import (
	prometheus "http-broadcaster/Prometheus"
	tools "http-broadcaster/Tools"
	varnish "http-broadcaster/Varnish"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// logRequest print the requests and wanted informations in log file
func logRequest(t time.Time, r *http.Request, s int, h map[string]string) {
	// Test if X-Cache-Tags header is empty
	if len(h) == 0 {
		log.Printf("%s %s - - %s \"%s %s %s\" %d 0 \"-\" \"%s\" %d\n",
			r.Host,
			r.Header["X-Forwarded-For"][0],
			t.Format("[02/Jan/2006:15:04:05 -0700]"),
			r.Method,
			r.URL.Path,
			r.Proto,
			s,
			r.UserAgent(),
			time.Since(t).Milliseconds(),
		)
	} else {
		var header string
		if h["X-Cache-Tags"] != "" {
			header = h["X-Cache-Tags"]
		} else {
			header = h["ApiPlatform-Ban-Regex"]
		}
		log.Printf("%s %s - - %s \"%s %s %s\" %d 0 \"-\" \"%s\" %d %s\n",
			r.Host,
			r.Header["X-Forwarded-For"][0],
			t.Format("[02/Jan/2006:15:04:05 -0700]"),
			r.Method,
			r.URL.Path,
			r.Proto,
			s,
			r.UserAgent(),
			time.Since(t).Milliseconds(),
			header,
		)
	}
}

// checkAllowedIP verify if the IPs is authorized to do BAN/PURGE request.
func checkAllowedIP(ip string) bool {
	return tools.IPAllowed(ip)
}

// RequestHandler handles requests to broadcast to all varnish instances.
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	var tag = make(map[string]string)
	ipAddress := r.RemoteAddr
	// check x-forwarded-for instead of RemoteAddr header because kube
	//ip, err := netip.ParseAddr(r.Header["X-Forwarded-For"][0])
	fwdAddress := r.Header.Get("X-Forwarded-For")
	if fwdAddress != "" {
		// Case there is a single IP in the header
		ipAddress = fwdAddress

		ips := strings.Split(fwdAddress, ",")
		if len(ips) > 1 {
			ipAddress = ips[0]
		}
	}

	// If IP is not authorized to do purge/ban requests, respond with 401.
	if !checkAllowedIP(ipAddress) {
		log.Printf("Client ip not authorized : %v", ipAddress)
		w.WriteHeader(401)
		_, _ = io.WriteString(w, strconv.Itoa(401))
		return
	}
	// If metrics are not enabled, return 404 on /metrics path.
	if r.URL.Path == "/metrics" && !prometheus.MetricsEnabled {
		w.WriteHeader(404)
		_, _ = io.WriteString(w, strconv.Itoa(404))
		return
	}
	t := time.Now()
	url := r.URL.String()
	method := r.Method
	h := r.Header.Get("X-Cache-Tags")
	if h != "" {
		tag["X-Cache-Tags"] = h
	}
	h = r.Header.Get("ApiPlatform-Ban-Regex")
	if h != "" {
		tag["ApiPlatform-Ban-Regex"] = h
	}
	status := varnish.SendToVarnish(method, url, tag)
	if prometheus.MetricsEnabled {
		prometheus.IncrementClientCounterVec(method)
	}
	// Return HTTP code 405 if not all varnish servers returned 200.
	if status != 200 {
		w.WriteHeader(405)
	}
	logRequest(t, r, status, tag)
	_, _ = io.WriteString(w, strconv.Itoa(status))
}

// HealthHandler handles healthcheck requests and return 200.
func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = io.WriteString(w, "OK")
}
