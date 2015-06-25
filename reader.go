package proxy

import (
	"github.com/mailgun/oxy/forward"
	"github.com/mailgun/oxy/roundrobin"
	"net/http"
	"net/url"
)

type Reader struct {
	slaves []string
	lb     *roundrobin.RoundRobin
}

func NewReader(slaves []string) (r *Reader) {
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
		lb.UpsertServer(slaveUrl)
	}

	return &Reader{lb: lb, slaves: slaves}
}

func (reader *Reader) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	writeLog("Reading: %v", req.URL.Path)
	reader.lb.ServeHTTP(w, req)
}
