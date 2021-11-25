package gateway

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMatch(t *testing.T) {
	rule := &ConfigRule{
		PathPrefix: "/health",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{Rules: configRules}

	router, err := ProvideRouter(c)
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "http://backend/health", nil)
	assert.NoError(t, err)

	actual, _ := router.match(req)

	expected := true

	assert.Equal(t, expected, actual)
}

func TestNoMatch(t *testing.T) {
	rule := &ConfigRule{
		PathPrefix: "/health",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{Rules: configRules}

	router, err := ProvideRouter(c)
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "http://backend/no-match", nil)
	assert.NoError(t, err)

	actual, _ := router.match(req)

	expected := false

	assert.Equal(t, expected, actual)
}

func TestMatchWithBasePath(t *testing.T) {
	basePath := "/base-path"

	rule := &ConfigRule{
		PathPrefix: "/health",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{BasePath: basePath, Rules: configRules}

	router, err := ProvideRouter(c)
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "http://backend/base-path/health", nil)
	assert.NoError(t, err)

	actual, _ := router.match(req)

	expected := true

	assert.Equal(t, expected, actual)
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
	c := &Config{Rules: configRules}

	router, err := ProvideRouter(c)
	assert.NoError(t, err)

	req0, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)
	req0.Header.Set(userAgentKey, headers0[userAgentKey][0])
	actual, actualRule0 := router.match(req0)
	assert.Equal(t, true, actual)
	assert.Equal(t, "backend0", actualRule0.Backend)

	req1, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)
	req1.Header.Set(userAgentKey, headers1[userAgentKey][0])
	actual, actualRule1 := router.match(req1)
	assert.Equal(t, true, actual)
	assert.Equal(t, "backend1", actualRule1.Backend)
}
