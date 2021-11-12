package rule

import (
	"github.com/dhis2-sre/go-rate-limite/pgk/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestMatch(t *testing.T) {
	rule := &config.Rule{
		PathPattern: "^\\/health$",
	}

	configRules := []config.Rule{*rule}
	c := &config.Config{Rules: configRules}

	rules := ProvideRules(c)

	u, err := url.Parse("http://backend/health")
	assert.NoError(t, err)

	req := &http.Request{URL: u}
	actual, _ := rules.Match(req)

	expected := true

	assert.Equal(t, expected, actual)
}

func TestNoMatch(t *testing.T) {
	rule := &config.Rule{
		PathPattern: "^\\/health$",
	}

	configRules := []config.Rule{*rule}
	c := &config.Config{Rules: configRules}

	rules := ProvideRules(c)

	u, err := url.Parse("http://backend/health-no-match")
	assert.NoError(t, err)

	req := &http.Request{URL: u}
	actual, _ := rules.Match(req)

	expected := false

	assert.Equal(t, expected, actual)
}
