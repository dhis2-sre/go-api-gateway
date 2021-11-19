package rule

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/dhis2-sre/go-rate-limite/pgk/config"
	"github.com/stretchr/testify/assert"
)

func TestPathMatch(t *testing.T) {
	expected := true

	rule := createRuleWithPathPattern("^\\/health$")

	actual := rule.pathMatch("/health")

	assert.Equal(t, expected, actual)
}

func TestPathNoMatch(t *testing.T) {
	expected := false

	rule := createRuleWithPathPattern("^\\/health$")

	actual := rule.pathMatch("/health-no-match")

	assert.Equal(t, expected, actual)
}

func createRuleWithPathPattern(pathPattern string) *Rule {
	rule := &Rule{
		Rule: config.Rule{
			Method:           "",
			PathPattern:      pathPattern,
			RequestPerSecond: 0,
			Burst:            0,
		},
		Handler: nil,
	}
	return rule
}

func TestMatch(t *testing.T) {
	rule := &config.Rule{
		PathPattern: "^\\/health$",
	}

	configRules := []config.Rule{*rule}
	c := &config.Config{Rules: configRules}

	rules := NewRules(c)

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

	rules := NewRules(c)

	u, err := url.Parse("http://backend/health-no-match")
	assert.NoError(t, err)

	req := &http.Request{URL: u}
	actual, _ := rules.Match(req)

	expected := false

	assert.Equal(t, expected, actual)
}
