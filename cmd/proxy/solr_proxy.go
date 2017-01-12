package main

import (
	"fmt"
	"github.com/KensoDev/go-solr-proxy"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"strings"
)

var (
	listenPort      = kingpin.Flag("listen-port", "Which port should the proxy listen on").Int()
	master          = kingpin.Flag("master", "Location to your master server").String()
	slaves          = kingpin.Flag("slaves", "Comma separated list of servers that act as slaves").String()
	awsRegion       = kingpin.Flag("aws-region", "Which AWS region should it use for the cache").Default("us-west-2").String()
	s3EndPoint      = kingpin.Flag("aws-endpoint", "AWS Endpoint for s3").Default("https://s3-us-west-2.amazonaws.com").String()
	s3BucketName    = kingpin.Flag("bucket-name", "What's the bucket name you want to save the documents in").String()
	logLocation     = kingpin.Flag("log-location", "Where do you want to keep logs").Default("stdout").String()
	buckerPrefix    = kingpin.Flag("bucket-prefix", "Prefix after the bucket name before the filename").String()
	disableS3Upload = kingpin.Flag("disable-upload", "Disable the S3 upload of documents").Bool()
)

func main() {
	kingpin.Parse()
	slaveServers := strings.Split(*slaves, ",")

	proxyConfig := &proxy.ProxyConfig{
		Master:      *master,
		Slaves:      slaveServers,
		LogLocation: *logLocation,
		AwsConfig: &proxy.AWSConfig{
			BucketName:    *s3BucketName,
			S3Endpoint:    *s3EndPoint,
			RegionName:    *awsRegion,
			BucketPrefix:  *buckerPrefix,
			DisableUpload: *disableS3Upload,
		},
	}

	config := proxy.NewConfigurationRender(proxyConfig)
	http.Handle("/proxy/configuration", config)

	p := proxy.NewProxy(proxyConfig)
	http.Handle("/", p)

	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *listenPort), nil); err != nil {
		panic(err)
	}
}
