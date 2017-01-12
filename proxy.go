package proxy

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"os"
	"regexp"
)

type ProxyConfig struct {
	Master      string
	Slaves      []string
	AwsConfig   *AWSConfig
	LogLocation string
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

type Proxy struct {
	updater *Updater
	reader  *Reader
	config  *ProxyConfig
}

func NewProxy(proxyConfig *ProxyConfig) *Proxy {
	updater := NewUpdater(proxyConfig.Master)
	reader := NewReader(proxyConfig.Slaves)

	if proxyConfig.LogLocation == "stdout" {
		log.SetOutput(os.Stdout)
	} else {
		f, err := os.OpenFile(proxyConfig.LogLocation, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}
		log.SetOutput(f)
	}

	return &Proxy{
		updater: updater,
		reader:  reader,
		config:  proxyConfig,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	re := regexp.MustCompile("\\/solr\\/(.*)\\/update$")
	matches := re.FindStringSubmatch(req.URL.Path)

	req.Close = true

	if len(matches) > 0 {
		p.updater.ServeHTTP(w, req, p.config.AwsConfig, matches[1])
	} else {
		p.reader.ServeHTTP(w, req)
	}
}

func writeLog(message string, params ...string) {
	log.Printf(message, params)
}
