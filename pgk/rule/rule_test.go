package rule

import (
	"github.com/dhis2-sre/go-rate-limite/pgk/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestRuleMatch(t *testing.T) {
	expected := true

	rule := createRuleWithPathPattern("POST", "^\\/health$")

	req := &http.Request{
		Method: "POST",
		URL:    &url.URL{Path: "/health"},
	}

	actual := rule.match(req)

	assert.Equal(t, expected, actual)
}

func TestRuleNoMatchPath(t *testing.T) {
	expected := false

	rule := createRuleWithPathPattern("POST", "^\\/health$")

	req := &http.Request{
		Method: "POST",
		URL:    &url.URL{Path: "/health-no-match"},
	}

	actual := rule.match(req)

	assert.Equal(t, expected, actual)
}

func TestRuleNoMatchMethod(t *testing.T) {
	expected := false

	rule := createRuleWithPathPattern("POST", "^\\/health$")

	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/health"},
	}

	actual := rule.match(req)

	assert.Equal(t, expected, actual)
}

func TestRuleWithoutMethod(t *testing.T) {
	expected := true

	rule := createRuleWithPathPattern("", "^\\/health$")

	req := &http.Request{
		Method: "WHATEVER",
		URL:    &url.URL{Path: "/health"},
	}

	actual := rule.match(req)

	assert.Equal(t, expected, actual)
}

func createRuleWithPathPattern(method, pathPattern string) *Rule {
	rule := &Rule{
		Rule: config.Rule{
			Method:           method,
			PathPattern:      pathPattern,
			RequestPerSecond: 0,
			Burst:            0,
		},
		Handler: nil,
	}
	return rule
}
