// Package varnish provides functions to build the list of varnish servers that will be used
package varnish

import (
	prometheus2 "http-broadcaster/Prometheus"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	// VarnishList contains the list of varnish servers.
	VarnishList []string
	status      = 200
)

// InitializeVarnishList sets varnishList variable according to the LIST_METHOD env var
func InitializeVarnishList(l string) []string {
	data := os.Getenv("VARNISH_SERVERS")
	sliceData := strings.Split(string(data), ",")
	return sliceData
}

// SendToVarnish send to all varnish servers define in varnishList the request with the PURGE or BAN method
// and the X-Cache-Tags header if necessary.
func SendToVarnish(method string, url string, tag map[string]string) int {
	status = 200

	// Take url to ban as argument.
	// Loop over the list of Varnish servers and send PURGE request to each.
	// Update status variable to check if servers have successfully purge url.
	for i := 0; i < len(VarnishList); i++ {
		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		// sanitize varnish server host.
		domain := strings.Trim(VarnishList[i], "\r\n")
		req, err := http.NewRequest(method, domain+url, nil)
		if err != nil {
			log.Println("Create new request : ", err)
		}
		// If X-Cache-Tags header is not empty with pass it to varnish.
		if tag["X-Cache-Tags"] != "" {
			req.Header.Add("X-Cache-Tags", tag["X-Cache-Tags"])
		}
		if tag["ApiPlatform-Ban-Regex"] != "" {
			req.Header.Add("ApiPlatform-Ban-Regex", tag["ApiPlatform-Ban-Regex"])
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Send new request : ", err)
		}
		if prometheus2.MetricsEnabled {
			prometheus2.IncrementBackendCounterVec(method)
		}
		if resp.StatusCode != 200 {
			status = 405
		}
		defer resp.Body.Close() //nolint:all https://github.com/securego/gosec/pull/935
	}
	return status
}
