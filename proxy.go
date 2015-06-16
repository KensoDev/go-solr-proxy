package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	p := new(Proxy)
	http.Handle("/solr/collection1/update", p)

	if err := http.ListenAndServe(":8080", nil); err != nil {
	}
}

type Proxy struct {
	proxy http.Handler
	err   error
}

type RequestReader struct {
	*bytes.Buffer
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	u, err := url.Parse("http://192.168.99.98:8983")
	if err != nil {
		p.err = err
		return
	}

	contents, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Printf("%s", err)
	}

	log.Printf("request: %s", contents)
	p.proxy = httputil.NewSingleHostReverseProxy(u)
	p.proxy.ServeHTTP(w, req)
}