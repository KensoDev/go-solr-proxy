package main

import (
	"fmt"
	"github.com/KensoDev/go-solr-proxy"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"strings"
)

var (
	listenPort   = kingpin.Flag("listen-port", "Which port should the proxy listen on").Int()
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
		Master: *master,
		Slaves: slaveServers,
		AwsConfig: &proxy.AWSConfig{
			BucketName: *s3BucketName,
			S3Endpoint: *s3EndPoint,
			RegionName: *awsRegion,
		},
	}

	fmt.Printf("You have %d slaves\n", len(slaveServers))

	p := proxy.NewProxy(proxyConfig)
	http.Handle("/", p)

	fmt.Printf("Starting proxy on port %d\n", *listenPort)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *listenPort), nil); err != nil {
		panic(err)
	}
}
