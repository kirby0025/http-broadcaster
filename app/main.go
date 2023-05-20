// Receive http PURGE request and broadcast it to several Varnish servers.
package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	// MethodPurge declaration for Varnish.
	MethodPurge = "PURGE"
)

var (
	varnishList = MakeVarnishList()
	status      = "200 Purged"
)

// MakeVarnishList reads the list of varnish servers from a file on disk.
func MakeVarnishList() []string {
	Data, err := os.ReadFile("./varnish")
	if err != nil {
		log.Fatal(err)
	}
	sliceData := strings.Split(string(Data), ",")
	return sliceData
}

// PurgeHandler handles PURGE request to broadcast it to all varnish instances.
func PurgeHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	remoteAddr := r.RemoteAddr
	status := SendToVarnish(url)
	if status != "200 Purged" {
		w.WriteHeader(405)
	}
	log.Println(remoteAddr + " Requested purge on " + url + " : " + status)
	io.WriteString(w, status)
}

// SendToVarnish send to all varnish servers define in varnishList the PURGE request.
func SendToVarnish(url string) string {
	status = "200 Purged"
	// Take url to ban as argument.
	// Loop over the list of Varnish servers and send PURGE request to each.
	// Update status variable to check if servers have successfully purge url.
	for i := 0; i < len(varnishList); i++ {
		client := &http.Client{}
		domain := strings.Trim(varnishList[i], "\r\n")
		req, err := http.NewRequest(MethodPurge, domain+url, nil)
		if err != nil {
			log.Fatal("Create new request : %s", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Send new request : %s", err)
		}
		if resp.StatusCode != 200 {
			status = "405 Not Allowed"
		}
	}
	return status
}

// HealthHandler handles healthcheck requests and return 200.
func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "OK")
}

func main() {
	logFile, err := os.OpenFile("./log/purge.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
	http.HandleFunc("/", PurgeHandler)
	http.HandleFunc("/healthcheck", HealthHandler)
	log.Fatal(http.ListenAndServe(":6081", nil))
}
