package gateway

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMatch(t *testing.T) {
	rule := &ConfigRule{
		PathPrefix: "/health",
		Backend:    "backend0",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{Backends: getBackends(), Rules: configRules}

	router, err := ProvideRouter(c)
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	actual, _ := router.match(req)

	assert.Equal(t, true, actual)
}

func TestRuleDefinesBackendOrDefaultBackend(t *testing.T) {
	rule := &ConfigRule{
		PathPrefix: "/health",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{Rules: configRules}

	_, err := ProvideRouter(c)
	assert.NotNil(t, err)
	assert.Equal(t, "either a rule needs to define a backend or a default backend needs to be defined", err.Error())
}

func TestMatchUnmappedBackend(t *testing.T) {
	rule := &ConfigRule{
		PathPrefix: "/health",
		Backend:    "some-undefined-backend",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{Rules: configRules}

	_, err := ProvideRouter(c)
	assert.NotNil(t, err)
	assert.Equal(t, "backend map contains not entry for: some-undefined-backend", err.Error())
}

func TestNoMatch(t *testing.T) {
	rule := &ConfigRule{
		PathPrefix: "/health",
		Backend:    "backend0",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{Backends: getBackends(), Rules: configRules}

	router, err := ProvideRouter(c)
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/no-match", nil)
	assert.NoError(t, err)

	actual, _ := router.match(req)

	assert.Equal(t, false, actual)
}

func TestMatchWithBasePath(t *testing.T) {
	basePath := "/base-path"

	rule := &ConfigRule{
		PathPrefix: "/health",
		Backend:    "backend0",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{BasePath: basePath, Backends: getBackends(), Rules: configRules}

	router, err := ProvideRouter(c)
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/base-path/health", nil)
	assert.NoError(t, err)

	actual, _ := router.match(req)

	assert.Equal(t, true, actual)
}

func TestMatchSamePathAndMethodButDifferentHeaders(t *testing.T) {
	userAgentKey := "User-Agent"

	headers0 := map[string][]string{
		userAgentKey: {"Go tests"},
	}
	rule0 := &ConfigRule{
		PathPrefix: "/health",
		Method:     "GET",
		Headers:    headers0,
		Backend:    "backend0",
	}

	headers1 := map[string][]string{
		userAgentKey: {"Some other client"},
	}
	rule1 := &ConfigRule{
		PathPrefix: "/health",
		Method:     "GET",
		Headers:    headers1,
		Backend:    "backend1",
	}

	configRules := []ConfigRule{*rule0, *rule1}
	c := &Config{Backends: getBackends(), Rules: configRules}

	router, err := ProvideRouter(c)
	assert.NoError(t, err)

	req0, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)
	req0.Header.Set(userAgentKey, headers0[userAgentKey][0])
	actual, actualRule0 := router.match(req0)
	assert.Equal(t, true, actual)
	assert.Equal(t, "backend0", actualRule0.Backend)

	req1, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)
	req1.Header.Set(userAgentKey, headers1[userAgentKey][0])
	actual, actualRule1 := router.match(req1)
	assert.Equal(t, true, actual)
	assert.Equal(t, "backend1", actualRule1.Backend)
}

func TestMatchWithHostname(t *testing.T) {
	rule := &ConfigRule{
		PathPrefix: "/health",
		Hostname:   "url",
		Backend:    "backend0",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{Backends: getBackends(), Rules: configRules}

	router, err := ProvideRouter(c)
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	match, actualRule := router.match(req)

	assert.Equal(t, true, match)
	assert.Equal(t, "url", actualRule.Hostname)
	assert.Equal(t, "backend0", actualRule.Backend)
}

func TestNoMatchWithHostname(t *testing.T) {
	rule := &ConfigRule{
		PathPrefix: "/health",
		Hostname:   "no-match",
		Backend:    "backend0",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{Backends: getBackends(), Rules: configRules}

	router, err := ProvideRouter(c)
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	match, _ := router.match(req)

	assert.Equal(t, false, match)
}
