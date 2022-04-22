package gateway

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	rule := &ConfigRule{
		PathPrefix: "/health",
		Backend:    "backend0",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{Backends: getBackends(), Rules: configRules}

	rules, err := ProvideRules(c)
	assert.NoError(t, err)

	router := ProvideRouter(rules)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	actual, _ := router.match(req)

	assert.Equal(t, true, actual)
}

func TestNoMatch(t *testing.T) {
	rule := &ConfigRule{
		PathPrefix: "/health",
		Backend:    "backend0",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{Backends: getBackends(), Rules: configRules}

	rules, err := ProvideRules(c)
	assert.NoError(t, err)

	router := ProvideRouter(rules)

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

	rules, err := ProvideRules(c)
	assert.NoError(t, err)

	router := ProvideRouter(rules)

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

	rules, err := ProvideRules(c)
	assert.NoError(t, err)

	router := ProvideRouter(rules)

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

func TestMatchHostname(t *testing.T) {
	rule := &ConfigRule{
		PathPrefix: "/health",
		Hostname:   "url",
		Backend:    "backend0",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{Backends: getBackends(), Rules: configRules}

	rules, err := ProvideRules(c)
	assert.NoError(t, err)

	router := ProvideRouter(rules)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	match, actualRule := router.match(req)

	assert.Equal(t, true, match)
	assert.Equal(t, "url", actualRule.Hostname)
	assert.Equal(t, "backend0", actualRule.Backend)
}

func TestNoMatchHostname(t *testing.T) {
	rule := &ConfigRule{
		PathPrefix: "/health",
		Hostname:   "no-match",
		Backend:    "backend0",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{Backends: getBackends(), Rules: configRules}

	rules, err := ProvideRules(c)
	assert.NoError(t, err)

	router := ProvideRouter(rules)

	req, err := http.NewRequest("GET", defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	match, _ := router.match(req)

	assert.Equal(t, false, match)
}

func TestMatchConfigWithMultipleHostnames(t *testing.T) {
	ruleA := &ConfigRule{
		PathPrefix: "/",
		Hostname:   "a.domain.org",
		Backend:    "backend0",
	}

	ruleB := &ConfigRule{
		PathPrefix: "/",
		Hostname:   "a.b.domain.org",
		Backend:    "backend0",
	}

	ruleC := &ConfigRule{
		PathPrefix: "/",
		Hostname:   "a.b.c.domain.org",
		Backend:    "backend0",
	}

	configRules := []ConfigRule{*ruleA, *ruleB, *ruleC}
	c := &Config{Backends: getBackends(), Rules: configRules}

	rules, err := ProvideRules(c)
	assert.NoError(t, err)

	router := ProvideRouter(rules)

	reqA, err := http.NewRequest("GET", "http://a.domain.org/", nil)
	assert.NoError(t, err)
	assertMatch(t, router, reqA, ruleA.Hostname)

	reqB, err := http.NewRequest("GET", "http://a.b.domain.org/", nil)
	assert.NoError(t, err)
	assertMatch(t, router, reqB, ruleB.Hostname)

	reqC, err := http.NewRequest("GET", "http://a.b.c.domain.org/", nil)
	assert.NoError(t, err)
	assertMatch(t, router, reqC, ruleC.Hostname)
}

func TestMatchSubdomain(t *testing.T) {
	ruleA := &ConfigRule{
		PathPrefix: "/",
		Hostname:   "*.a.domain.org",
		Backend:    "backend0",
	}

	ruleB := &ConfigRule{
		PathPrefix: "/",
		Hostname:   "a.domain.org",
		Backend:    "backend0",
	}

	configRules := []ConfigRule{*ruleA, *ruleB}
	c := &Config{Backends: getBackends(), Rules: configRules}

	rules, err := ProvideRules(c)
	assert.NoError(t, err)

	router := ProvideRouter(rules)

	reqA, err := http.NewRequest("GET", "http://sub.a.domain.org/", nil)
	assert.NoError(t, err)
	assertMatch(t, router, reqA, ruleA.Hostname)

	reqB, err := http.NewRequest("GET", "http://a.domain.org/", nil)
	assert.NoError(t, err)
	assertMatch(t, router, reqB, ruleB.Hostname)
}

func assertMatch(t *testing.T, router *Router, req *http.Request, hostname string) {
	match, actualRule := router.match(req)

	assert.Equal(t, true, match)
	assert.Equal(t, hostname, actualRule.Hostname)
	assert.Equal(t, "backend0", actualRule.Backend)
}
