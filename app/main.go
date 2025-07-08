// Receive http PURGE request and broadcast it to several Varnish servers.
package main

import (
	h "http-broadcaster/Http"
	prometheus2 "http-broadcaster/Prometheus"
	tools "http-broadcaster/Tools"
	varnish "http-broadcaster/Varnish"
	"log"
	"net/http"
	"os"

	"github.com/alexflint/go-arg"
)

var (
	args struct {
		Log     string `arg:"-l,--logfile" help:"location of output logfile." default:"/app/http-broadcaster.log"`
		EnvFile string `arg:"-e,--envfile" help:"location of file containing environment variables." default:"/vault/secrets/.env"`
		Metrics bool   `arg:"--metrics" help:"enable prometheus exporter on /metrics." default:"false"`
	}
)

func main() {
	arg.MustParse(&args)
	tools.InitLog(args.Log)
	tools.ReadDotEnvFile(args.EnvFile)
	tools.ClientList = tools.InitAllowedIPList(os.Getenv("CLIENT_LIST"))
	varnish.VarnishList = varnish.InitializeVarnishList(os.Getenv("VARNISH_SERVERS"))
	prometheus2.InitMetrics(args.Metrics)

	http.HandleFunc("/", h.RequestHandler)
	http.HandleFunc("/healthcheck", h.HealthHandler)
	log.Fatal(http.ListenAndServe(":6081", nil))
}
