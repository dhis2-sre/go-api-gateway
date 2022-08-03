package gateway

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_jwtAuth_validatePublicKey(t *testing.T) {
	rule := ConfigRule{
		PathPrefix:     "/health",
		Authentication: "jwt",
	}

	configRules := []ConfigRule{rule}
	c := &Config{
		Backends:       getBackends(),
		DefaultBackend: defaultBackend,
		Authentication: Authentication{Jwt: Jwt{publicKey}},
		Rules:          configRules,
	}

	jwtAuth := NewJwtAuth(c)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", validAccessToken)

	assert.NoError(t, jwtAuth.ValidateRequest(req))
}

func Test_jwtAuth_validatePublicKey_bearer(t *testing.T) {
	rule := ConfigRule{
		PathPrefix:     "/health",
		Authentication: "jwt",
	}

	configRules := []ConfigRule{rule}
	c := &Config{
		Backends:       getBackends(),
		DefaultBackend: defaultBackend,
		Authentication: Authentication{Jwt: Jwt{publicKey}},
		Rules:          configRules,
	}

	jwtAuth := NewJwtAuth(c)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+validAccessToken)

	assert.NoError(t, jwtAuth.ValidateRequest(req))
}

func Test_jwtAuth_validatePublicKey_invalidToken(t *testing.T) {
	rule := ConfigRule{
		PathPrefix:     "/health",
		Authentication: "jwt",
	}

	configRules := []ConfigRule{rule}
	c := &Config{
		Backends:       getBackends(),
		DefaultBackend: defaultBackend,
		Authentication: Authentication{Jwt: Jwt{publicKey}},
		Rules:          configRules,
	}

	jwtAuth := NewJwtAuth(c)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "bla bla")

	assert.Error(t, jwtAuth.ValidateRequest(req))
}

func Test_jwtAuth_validateJwks(t *testing.T) {
	rule := ConfigRule{
		PathPrefix:     "/health",
		Authentication: "jwt",
	}

	configRules := []ConfigRule{rule}
	c := &Config{
		Backends:       getBackends(),
		DefaultBackend: defaultBackend,
		Authentication: Authentication{Jwks: Jwks{
			Host:                   "http://jwks/jwks.json",
			Index:                  0,
			MinimumRefreshInterval: 60,
		}},
		Rules: configRules,
	}

	jwtAuth := NewJwtAuth(c)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", validAccessToken)

	assert.NoError(t, jwtAuth.ValidateRequest(req))
}

func Test_jwtAuth_validateJwks_bearer(t *testing.T) {
	rule := ConfigRule{
		PathPrefix:     "/health",
		Authentication: "jwt",
	}

	configRules := []ConfigRule{rule}
	c := &Config{
		Backends:       getBackends(),
		DefaultBackend: defaultBackend,
		Authentication: Authentication{Jwks: Jwks{
			Host:                   "http://jwks/jwks.json",
			Index:                  0,
			MinimumRefreshInterval: 60,
		}},
		Rules: configRules,
	}

	jwtAuth := NewJwtAuth(c)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+validAccessToken)

	assert.NoError(t, jwtAuth.ValidateRequest(req))
}

func Test_jwtAuth_validateJwks_invalidToken(t *testing.T) {
	rule := ConfigRule{
		PathPrefix:     "/health",
		Authentication: "jwt",
	}

	configRules := []ConfigRule{rule}
	c := &Config{
		Backends:       getBackends(),
		DefaultBackend: defaultBackend,
		Authentication: Authentication{Jwks: Jwks{
			Host:                   "http://jwks/jwks.json",
			Index:                  0,
			MinimumRefreshInterval: 60,
		}},
		Rules: configRules,
	}

	jwtAuth := NewJwtAuth(c)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "bla bla")

	assert.Error(t, jwtAuth.ValidateRequest(req))
}
