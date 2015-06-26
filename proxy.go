package proxy

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/formatters/logstash"
	"net/http"
	"os"
	"regexp"
)

type ProxyConfig struct {
	master    string
	slaves    []string
	awsConfig *AWSConfig
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

func NewProxy(proxyConfig *ProxyConfig) (p *Proxy) {
	updater := NewUpdater(proxyConfig.master)
	reader := NewReader(proxyConfig.slaves)

	p = &Proxy{
		updater: updater,
		reader:  reader,
		config:  proxyConfig,
	}
	return p
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	isUpdate, _ := regexp.MatchString("\\/solr\\/(\\S+)\\/update$", req.URL.Path)
	log.Printf("url: %v %b", req.URL.Path, isUpdate)

	if isUpdate {
		p.updater.ServeHTTP(w, req, p.config.awsConfig)
	} else {
		p.reader.ServeHTTP(w, req)
	}
}

func writeLog(message string, params ...string) {
	log.Printf(message, params)
}
