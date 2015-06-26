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
	rdr := RequestReader{bytes.NewBuffer(buf)}
	content, err := ioutil.ReadAll(rdr)
	if err != nil {
		http.Error(w, "Could not read request body\n"+err.Error(), http.StatusBadRequest)
		return
	}

	doc, err := ParseXMLDocument(content)
	if err != nil {
		http.Error(w, "Could not deserialize request body\n"+err.Error(), http.StatusBadRequest)
		return
	}

	solrDoc := doc.GetSolrDocument()
	err = solrDoc.Cache(awsConfig)
	if err != nil {
		http.Error(w, "Could not cache solr document\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	rdr = RequestReader{bytes.NewBuffer(buf)}
	req.Body = rdr

	writeLog("Updating: %v", req.URL.Path)
	updater.lb.ServeHTTP(w, req)
}
