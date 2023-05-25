// Package varnish provides functions to build the list of varnish servers that will be used
package varnish

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	vault "http-broadcaster/Vault"
)

var (
	varnishList = InitializeVarnishList()
	status      = "200 Purged"
)

// InitializeVarnishList sets varnishList variable according to the LIST_METHOD env var
func InitializeVarnishList() []string {
	switch method := os.Getenv("LIST_METHOD"); method {
	case "vault":
		return GetVarnishListFromVault()
	case "file":
		return GetVarnishListFromFile()
	default:
		panic("LIST_METHOD empty, no provided method to retrieve varnish list")
	}
}

// GetVarnishListFromVault builds a list of varnish servers from Vault.
func GetVarnishListFromVault() []string {
	var value []string
	client := vault.InitVaultConnection()
	secret, err := client.KVv2("app").Get(context.Background(), "http-broadcaster/stg/envVars")
	if err != nil {
		log.Fatal("unable to read secret: %w", err)
		return value
	}

	// selecting list key from retrieved secret
	list, ok := secret.Data["varnish_list"].(string)
	if !ok {
		log.Fatal("value type assertion failed: %T %#v", secret.Data["varnish_list"], secret.Data["varnish_list"])
		return value
	}
	value = strings.Split(string(list), ",")
	return value
}

// GetVarnishListFromFile reads the list of varnish servers from a file on disk.
func GetVarnishListFromFile() []string {
	Data, err := os.ReadFile("./varnish")
	if err != nil {
		log.Fatal(err)
	}
	sliceData := strings.Split(string(Data), ",")
	return sliceData
}

// SendToVarnish send to all varnish servers define in varnishList the request with the PURGE or BAN method
// and the X-Cache-Tags header if necessary.
func SendToVarnish(method string, url string, tag string) string {
	status = "200 Purged"

	// Take url to ban as argument.
	// Loop over the list of Varnish servers and send PURGE request to each.
	// Update status variable to check if servers have successfully purge url.
	for i := 0; i < len(varnishList); i++ {
		client := &http.Client{}
		domain := strings.Trim(varnishList[i], "\r\n")
		req, err := http.NewRequest(method, domain+url, nil)
		if err != nil {
			log.Fatal("Create new request : %s", err)
		}
		if tag != "" {
			req.Header.Add("X-Cache-Tags", tag)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Send new request : ", err)
		}
		if resp.StatusCode != 200 {
			status = "405 Not Allowed"
		}
	}
	return status
}
