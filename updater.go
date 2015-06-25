package proxy

import (
	"github.com/mailgun/oxy/forward"
	"github.com/mailgun/oxy/roundrobin"
	"net/http"
	"net/url"
)

type Updater struct {
	lb     *roundrobin.RoundRobin
	master string
}

func NewUpdater(master string) (updater *Updater) {
	fwd, err := forward.New()
	lb, err := roundrobin.New(fwd)
	masterUrl, err := url.Parse(master)
	lb.UpsertServer(masterUrl)

	if err != nil {
		panic(err)
	}

	return &Updater{lb: lb, master: master}
}

func (updater *Updater) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	writeLog("Updating: %v", req.URL.Path)
	updater.lb.ServeHTTP(w, req)
}
