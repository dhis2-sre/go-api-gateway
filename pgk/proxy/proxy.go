package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Transparently(backend string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s -> %s", req.Method, req.URL.Path, backend)

		backendUrl, err := url.Parse(backend)
		if err != nil {
			http.Error(w, "error parsing backend url: "+backend, http.StatusInternalServerError)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(backendUrl)
		proxy.ServeHTTP(w, req)
	})
}
