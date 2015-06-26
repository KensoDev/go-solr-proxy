package proxy

import (
	"bytes"
	"github.com/mailgun/oxy/forward"
	"github.com/mailgun/oxy/roundrobin"
	"io/ioutil"
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

func (updater *Updater) ServeHTTP(w http.ResponseWriter, req *http.Request, awsConfig *AWSConfig) {
	buf, _ := ioutil.ReadAll(req.Body)
	rdr1 := RequestReader{bytes.NewBuffer(buf)}
	rdr2 := RequestReader{bytes.NewBuffer(buf)}
	req.Body = rdr2

	writeLog("Updating: %v", req.URL.Path)
	updater.lb.ServeHTTP(w, req)

	content, _ := ioutil.ReadAll(rdr1)
	doc := ParseXMLDocument(content)
	solrDoc := doc.GetSolrDocument()
	solrDoc.Cache(awsConfig)
}
