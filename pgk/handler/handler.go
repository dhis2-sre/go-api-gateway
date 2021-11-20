package handler

import (
	"errors"
	"github.com/dhis2-sre/go-rate-limiter/pgk/config"
	"github.com/dhis2-sre/go-rate-limiter/pgk/rule"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
	"strings"
)

func ProvideHandler(c *config.Config, router *rule.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if match, r := router.Match(req); match {
			if r.Authentication == "jwt" {
				valid, err := validateRequest(c.Authentication.Jwt.PublicKey, req)
				if !valid || err != nil {
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}
			// TODO: This shouldn't be necessary if we're running in cluster only accessing services
			fixHost(req, r.Backend)
			r.Handler.ServeHTTP(w, req)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}
}

func fixHost(req *http.Request, ruleBackend string) {
	if !strings.HasSuffix(req.Host, ruleBackend) {
		var backend string
		if strings.HasPrefix(ruleBackend, "https://") {
			backend = strings.TrimPrefix(ruleBackend, "https://")
		}
		if strings.HasPrefix(ruleBackend, "http://") {
			backend = strings.TrimPrefix(ruleBackend, "http://")
		}
		req.Host = backend
	}
}

func validateRequest(publicKey string, req *http.Request) (bool, error) {
	authorizationHeader := req.Header.Get("Authorization")
	if !strings.HasPrefix(authorizationHeader, "Bearer") {
		return false, errors.New("no bearer token prefix found in header")
	}
	tokenString := strings.TrimPrefix(authorizationHeader, "Bearer")

	_, err := jwt.Parse(
		[]byte(tokenString),
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, publicKey),
	)
	if err != nil {
		return true, nil
	}
	return false, err
}
