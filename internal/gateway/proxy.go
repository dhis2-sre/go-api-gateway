package gateway

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func provideTransparentProxy(backendUrl *url.URL) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		proxy := httputil.NewSingleHostReverseProxy(backendUrl)
		proxy.ServeHTTP(w, req)
	}
}
