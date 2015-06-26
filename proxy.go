package proxy

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/formatters/logstash"
	"net/http"
	"os"
	"regexp"
)

type ProxyConfig struct {
	Master    string
	Slaves    []string
	AwsConfig *AWSConfig
}

func init() {
	log.SetFormatter(&logstash.LogstashFormatter{Type: "solr-proxy"})
	log.SetOutput(os.Stdout)
}

type Proxy struct {
	updater *Updater
	reader  *Reader
	config  *ProxyConfig
}

func NewProxy(proxyConfig *ProxyConfig) *Proxy {
	updater := NewUpdater(proxyConfig.Master)
	reader := NewReader(proxyConfig.Slaves)

	return &Proxy{
		updater: updater,
		reader:  reader,
		config:  proxyConfig,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	isUpdate, _ := regexp.MatchString("\\/solr\\/(\\S+)\\/update$", req.URL.Path)
	log.Printf("url: %v %b", req.URL.Path, isUpdate)

	if isUpdate {
		p.updater.ServeHTTP(w, req, p.config.AwsConfig)
	} else {
		p.reader.ServeHTTP(w, req)
	}
}

func writeLog(message string, params ...string) {
	log.Printf(message, params)
}
