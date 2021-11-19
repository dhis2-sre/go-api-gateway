package proxy

import (
	"github.com/dhis2-sre/go-rate-limiter/pgk/config"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// TODO: Rename ProvideTransparentProxy
func ProvideProxy(c *config.Config) *Proxy {
	backendUrl, err := url.Parse("http://backend:8080/")
	if err != nil {
		log.Fatal(err)
	}
	return &Proxy{c, backendUrl}
}

type Proxy struct {
	c          *config.Config
	BackendUrl *url.URL
}

func (p *Proxy) TransparentProxyHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s -> %s", req.Method, req.URL.Path, p.BackendUrl)
	proxy := httputil.NewSingleHostReverseProxy(p.BackendUrl)
	proxy.ServeHTTP(w, req)
}

func (p *Proxy) SetBackendUrl(backendUrl string) {
	backendUrlParsed, err := url.Parse(backendUrl)
	if err != nil {
		log.Fatal(err)
	}
	p.BackendUrl = backendUrlParsed
}
