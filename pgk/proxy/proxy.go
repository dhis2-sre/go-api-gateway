package proxy

import (
	"github.com/dhis2-sre/go-rate-limite/pgk/config"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func ProvideProxy(c config.Config) *Proxy {
	return &Proxy{c}
}

type Proxy struct {
	c config.Config
}

func (p *Proxy) TransparentProxyHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s -> %s", req.Method, req.URL.Path, p.c.Backend)

	backendUrl, err := url.Parse(p.c.Backend)
	if err != nil {
		http.Error(w, "error parsing backend url: "+p.c.Backend, http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(backendUrl)
	proxy.ServeHTTP(w, req)
}
