package proxy

import (
	"encoding/json"
	"net/http"
)

type Configuration struct {
	proxyconfig *ProxyConfig
}

func (c *Configuration) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	json, err := json.Marshal(c.proxyconfig)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func NewConfigurationRender(proxyConfig *ProxyConfig) (config *Configuration) {
	config = &Configuration{
		proxyconfig: proxyConfig,
	}
	return
}
