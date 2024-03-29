package gateway

import (
	"log"
	"net/http"
	"strings"
)

type auth interface {
	ValidateRequest(req *http.Request) (bool, error)
}

func NewHandler(config *Config, router *router, auth auth) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s%s", req.Method, req.URL.Host, req.URL.Path)

		if match, rule := router.match(req); match {
			if rule.Block {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			target := rule.PathReplace.Target
			if target != "" {
				path := strings.Replace(req.URL.Path, target, rule.PathReplace.Replacement, 1)
				req.URL.Path = path
			}

			if rule.Authentication == "jwt" && (req.Method != "OPTIONS" || config.Authentication.AuthenticateHttpOptionsMethod) {
				valid, err := auth.ValidateRequest(req)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusForbidden)
					return
				}
				if !valid {
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}

			contentType := req.Header.Get("Content-Type")
			if strings.HasPrefix(contentType, "multipart/form-data") {
				req.Body = http.MaxBytesReader(w, req.Body, config.MaxMultipartSize<<20)
			}

			// TODO: This is only necessary if the service behind this gateway is using the http host header... Maybe this should be done differently?
			fixHost(req, rule.Backend)
			rule.Handler.ServeHTTP(w, req)
			return
		}

		log.Printf("No match: %+v", req)
		w.WriteHeader(http.StatusMisdirectedRequest)
	}
}

func fixHost(req *http.Request, ruleBackend string) {
	if strings.HasSuffix(req.Host, ruleBackend) {
		return
	}

	if strings.HasPrefix(ruleBackend, "https://") {
		req.Host = strings.TrimPrefix(ruleBackend, "https://")
		return
	}

	if strings.HasPrefix(ruleBackend, "http://") {
		req.Host = strings.TrimPrefix(ruleBackend, "http://")
		return
	}
}
