package main

import (
	"fmt"
	"github.com/kensodev/go-solr-proxy"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"strings"
)

var (
	listenPort = kingpin.Flag("listen-port", "Which port should the proxy listen on").String()
	master     = kingpin.Flag("master", "Location to your master server").String()
	slaves     = kingpin.Flag("slaves", "Comma separated list of servers that act as slaves").String()
)

func main() {
	kingpin.Parse()

	slaveServers := strings.Split(*slaves, ",")
	fmt.Printf("You have %s slaves", len(slaveServers))

	p := proxy.NewProxy(*master, slaveServers)
	http.Handle("/", p)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", *listenPort), nil); err != nil {
		panic(err)
	}
}
