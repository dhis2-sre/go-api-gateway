package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func ProvideTransparentProxy(backendUrl *url.URL) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s -> %s", req.Method, req.URL.Path, backendUrl)
		proxy := httputil.NewSingleHostReverseProxy(backendUrl)
		proxy.ServeHTTP(w, req)
	}
}
