package main

import (
	"fmt"
	"github.com/KensoDev/go-solr-proxy"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
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
		Version:     "1.0.1",
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

	pinger := &Pinger{}
	http.Handle("/ping", pinger)

	p := proxy.NewProxy(proxyConfig)
	http.Handle("/", p)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *listenPort), nil); err != nil {
		panic(err)
	}
}

type Pinger struct{}

func (c *Pinger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	resp, err := http.Get("http://localhost:8982/solr/zh/select?q=*%3A*&wt=json&indent=true")
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	w.Write(body)
}
