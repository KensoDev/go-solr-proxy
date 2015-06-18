package proxy

import (
	"bytes"
	"github.com/mailgun/oxy/forward"
	"github.com/mailgun/oxy/roundrobin"
	"net/http"
	"net/url"
)

func main() {
	p := new(Proxy)
	http.Handle("/solr/collection1/update", p)

	if err := http.ListenAndServe(":8080", nil); err != nil {
	}
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

func NewProxy(master string, slaves ...string) *Proxy {
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
		lb.UpsertServer(slaveUrl)
	}

	return &Proxy{master: master, lb: lb}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p.lb.ServeHTTP(w, req)
}