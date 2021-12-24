package gateway

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
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

	jwtAuth := ProvideJwtAuth(c)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", validAccessToken)

	valid, err := jwtAuth.ValidateRequest(req)
	assert.NoError(t, err)

	assert.True(t, valid)
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

	jwtAuth := ProvideJwtAuth(c)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+validAccessToken)

	valid, err := jwtAuth.ValidateRequest(req)
	assert.NoError(t, err)

	assert.True(t, valid)
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

	jwtAuth := ProvideJwtAuth(c)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "bla bla")

	valid, err := jwtAuth.ValidateRequest(req)
	assert.Error(t, err)

	assert.False(t, valid)
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

	jwtAuth := ProvideJwtAuth(c)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", validAccessToken)

	valid, err := jwtAuth.ValidateRequest(req)
	assert.NoError(t, err)

	assert.True(t, valid)
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

	jwtAuth := ProvideJwtAuth(c)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+validAccessToken)

	valid, err := jwtAuth.ValidateRequest(req)
	assert.NoError(t, err)

	assert.True(t, valid)
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

	jwtAuth := ProvideJwtAuth(c)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "bla bla")

	valid, err := jwtAuth.ValidateRequest(req)
	assert.Error(t, err)

	assert.False(t, valid)
}
