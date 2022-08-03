package gateway

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type auth interface {
	ValidateRequest(req *http.Request) error
}

func newTransparentProxy(rule ConfigRule, auth auth) (http.HandlerFunc, error) {
	backendUrl, err := url.Parse(rule.Backend)

	return func(w http.ResponseWriter, req *http.Request) {
		if rule.Block {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		target := rule.PathReplace.Target
		if target != "" {
			path := strings.Replace(req.URL.Path, target, rule.PathReplace.Replacement, 1)
			req.URL.Path = path
		}

		if rule.Authentication == "jwt" {
			err := auth.ValidateRequest(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
		}

		proxy := httputil.NewSingleHostReverseProxy(backendUrl)
		proxy.ServeHTTP(w, req)
	}, err
}

func NewBlockingProxy(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	return
}
