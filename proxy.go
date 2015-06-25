package proxy

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/mailgun/oxy/forward"
	"github.com/mailgun/oxy/roundrobin"
	"net/http"
	"net/url"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

type Proxy struct {
	proxy  http.Handler
	lb     *roundrobin.RoundRobin
	master string
	err    error
}

type RequestReader struct {
	*bytes.Buffer
}

func NewProxy(master string, slaves []string) *Proxy {
	// oxy lb from slaves
	fwd, err := forward.New()
	lb, err := roundrobin.New(fwd)
	if err != nil {
		panic(err)
	}

	for _, slave := range slaves {
		slaveUrl, err := url.Parse(slave)
		if err != nil {
			panic(err)
		}
		log.Printf("Adding upsert server %v", slave)
		lb.UpsertServer(slaveUrl)
	}

	return &Proxy{master: master, lb: lb}
}

type SmartUpdater struct {
	fwd *forward.Forwarder
}

func (updater *SmartUpdater) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("Updating: %v", req.URL.Path)
	updater.fwd.ServeHTTP(w, req)
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	isUpdate := false // grep /update/ req.URL.Path
	var next http.Handler

	if isUpdate {
		fwd, _ := forward.New()
		updater := &SmartUpdater{fwd: fwd}
		next = updater
	} else {
		log.Printf("Reading: %v", req.URL.Path)
		next = p.lb
	}
	next.ServeHTTP(w, req)
}
