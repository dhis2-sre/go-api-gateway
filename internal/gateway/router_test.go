package gateway

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
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

	u, err := url.Parse("http://backend/health")
	assert.NoError(t, err)

	req := &http.Request{URL: u, Method: "GET"}
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

	u, err := url.Parse("http://backend/no-match")
	assert.NoError(t, err)

	req := &http.Request{URL: u, Method: "GET"}
	actual, _ := router.match(req)

	expected := false

	assert.Equal(t, expected, actual)
}
