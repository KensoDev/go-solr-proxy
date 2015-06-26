package main

import (
	"fmt"
	"github.com/kensodev/go-solr-proxy"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"strings"
)

var (
	listenPort   = kingpin.Flag("listen-port", "Which port should the proxy listen on").String()
	master       = kingpin.Flag("master", "Location to your master server").String()
	slaves       = kingpin.Flag("slaves", "Comma separated list of servers that act as slaves").String()
	awsRegion    = kingpin.Flag("aws-region", "Which AWS region should it use for the cache").Default("us-west-2").String()
	s3EndPoint   = kingpin.Flag("aws-endpoint", "AWS Endpoint for s3").Default("https://s3-us-west-2.amazonaws.com").String()
	s3BucketName = kingpin.Flag("bucket-name", "What's the bucket name you want to save the documents in").String()
)

func main() {
	kingpin.Parse()
	slaveServers := strings.Split(*slaves, ",")

	proxyConfig := &proxy.ProxyConfig{
		master: *master,
		slaves: slaveServers,
		awsConfig: &proxy.AWSConfig{
			BucketName: *s3BucketName,
			S3Endpoint: *s3EndPoint,
			RegionName: *awsRegion,
		},
	}

	fmt.Printf("You have %s slaves", len(slaveServers))

	p := proxy.NewProxy(*master, slaveServers)
	http.Handle("/", p)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", *listenPort), nil); err != nil {
		panic(err)
	}
}
