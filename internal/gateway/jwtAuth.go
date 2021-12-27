package gateway

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"log"
	"net/http"
	"time"
)

type JwtAuth interface {
	ValidateRequest(req *http.Request) (bool, error)
}

func ProvideJwtAuth(c *Config) JwtAuth {
	publicKey, err := providePublicKey(c.Authentication.Jwt.PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	if publicKey != nil {
		return jwtAuth{c, publicKey, nil}
	}

	jwksHost := c.Authentication.Jwks.Host
	autoRefresh, err := provideJwkAutoRefresh(jwksHost, c.Authentication.Jwks.MinimumRefreshInterval*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	if autoRefresh != nil {
		return jwtAuth{c, nil, autoRefresh}
	}

	return nil
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

// TODO: https://github.com/lestrrat-go/jwx/blob/main/examples/jwk_example_test.go#L188
func provideJwkAutoRefresh(host string, minRefreshInterval time.Duration) (*jwk.AutoRefresh, error) {
	if host != "" {
		ctx := context.TODO()
		ar := jwk.NewAutoRefresh(ctx)
		ar.Configure(host, jwk.WithMinRefreshInterval(minRefreshInterval))

		_, err := ar.Refresh(ctx, host)
		if err != nil {
			return nil, err
		}

		return ar, nil
	}
	return nil, nil
}

type jwtAuth struct {
	config      *Config
	publicKey   *rsa.PublicKey
	autoRefresh *jwk.AutoRefresh
}

func (j jwtAuth) ValidateRequest(req *http.Request) (bool, error) {
	if j.autoRefresh != nil {
		return j.validateJwks(req)
	}

	if j.publicKey != nil {
		return j.validatePublicKey(req)
	}

	return false, errors.New("no validator configured")
}

func (j jwtAuth) validatePublicKey(req *http.Request) (bool, error) {
	return j.validateStaticPublicKey(j.publicKey, req)
}

func (j jwtAuth) validateStaticPublicKey(publicKey *rsa.PublicKey, req *http.Request) (bool, error) {
	_, err := jwt.ParseRequest(req,
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, publicKey),
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (j jwtAuth) validateJwks(req *http.Request) (bool, error) {
	keySet, err := j.autoRefresh.Fetch(context.TODO(), j.config.Authentication.Jwks.Host)
	if err != nil {
		return false, err
	}

	if key, ok := keySet.Get(j.config.Authentication.Jwks.Index); ok {
		publicKey := &rsa.PublicKey{}
		err := key.Raw(publicKey)
		if err != nil {
			return false, err
		}
		return j.validateStaticPublicKey(publicKey, req)
	}

	return false, nil
}
