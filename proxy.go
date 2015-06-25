package proxy

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/formatters/logstash"
	"net/http"
	"os"
	"regexp"
)

func init() {
	log.SetFormatter(&logstash.LogstashFormatter{Type: "solr-proxy"})
	log.SetOutput(os.Stdout)
}

type Proxy struct {
	updater *Updater
	reader  *Reader
}

func NewProxy(master string, slaves []string) (p *Proxy) {
	updater := NewUpdater(master)
	reader := NewReader(slaves)

	return &Proxy{
		updater: updater,
		reader:  reader,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	isUpdate, _ := regexp.MatchString("\\/solr\\/(\\S+)\\/update$", req.URL.Path)
	log.Printf("url: %v %b", req.URL.Path, isUpdate)

	if isUpdate {
		p.updater.ServeHTTP(w, req)
	} else {
		p.reader.ServeHTTP(w, req)
	}
}

func writeLog(message string, params ...string) {
	log.Printf(message, params)
}
