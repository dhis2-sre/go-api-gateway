package gateway

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"log"
	"net/http"
	"strings"
)

func ProvideHandler(c *Config, router *Router) http.HandlerFunc {
	publicKey, err := providePublicKey(c.Authentication.Jwt.PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s%s", req.Method, req.URL.Host, req.URL.Path)
		if match, rule := router.match(req); match {
			if rule.Block {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if rule.Authentication == "jwt" {
				valid, err := validateRequest(publicKey, req)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusForbidden)
					return
				}
				if !valid {
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}
			// TODO: This is only necessary if the service behind this gateway is using the http host header... Maybe this should be done differently?
			fixHost(req, rule.Backend)
			rule.Handler.ServeHTTP(w, req)
			return
		}
		w.WriteHeader(http.StatusMisdirectedRequest)
	}
}

func providePublicKey(publicKeyString string) (*rsa.PublicKey, error) {
	if publicKeyString != "" {
		decode, _ := pem.Decode([]byte(publicKeyString))
		if decode == nil {
			return nil, errors.New("decoding of public key failed")
		}
		publicKey, err := x509.ParsePKIXPublicKey(decode.Bytes)
		if err != nil {
			return nil, err
		}
		return publicKey.(*rsa.PublicKey), nil
	}
	return nil, nil
}

func validateRequest(publicKey *rsa.PublicKey, req *http.Request) (bool, error) {
	_, err := jwt.ParseRequest(req,
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, publicKey),
	)
	if err != nil {
		return false, err
	}
	return true, nil
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
