package gateway

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRuleDefinesBackendOrDefaultBackend(t *testing.T) {
	rule := &ConfigRule{
		PathPrefix: "/health",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{Rules: configRules}

	_, err := ProvideRules(c)

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

	_, err := ProvideRules(c)

	assert.NotNil(t, err)
	assert.Equal(t, "backend map contains not entry for: some-undefined-backend", err.Error())
}

func TestLen(t *testing.T) {
	rule := &ConfigRule{
		PathPrefix: "/health",
		Method:     "GET",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{DefaultBackend: defaultBackend, Backends: getBackends(), Rules: configRules}

	rules, err := ProvideRules(c)
	assert.NoError(t, err)

	assert.Equal(t, 1, rules.Len())
}

func TestLenNoMethod(t *testing.T) {
	rule := &ConfigRule{
		PathPrefix: "/health",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{DefaultBackend: defaultBackend, Backends: getBackends(), Rules: configRules}

	rules, err := ProvideRules(c)
	assert.NoError(t, err)

	assert.Equal(t, 9, rules.Len())
}

func TestLookup(t *testing.T) {
	rule := &ConfigRule{
		PathPrefix: "/health",
		Method:     "GET",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{DefaultBackend: defaultBackend, Backends: getBackends(), Rules: configRules}

	r, err := ProvideRules(c)
	assert.NoError(t, err)

	i, match := r.Lookup([]byte("GET" + rule.PathPrefix))
	assert.True(t, match)

	rules := i.([]*Rule)
	assert.Equal(t, rule.PathPrefix, rules[0].PathPrefix)
	assert.Equal(t, rule.Method, rules[0].Method)
}

func TestLookupWithTwoRules(t *testing.T) {
	pathPrefix := "/health"
	method := "GET"

	ruleA := &ConfigRule{
		PathPrefix: pathPrefix,
		Method:     method,
		Hostname:   "domain.org",
	}

	ruleB := &ConfigRule{
		PathPrefix: pathPrefix,
		Method:     method,
		Hostname:   "other.domain.org",
	}

	configRules := []ConfigRule{*ruleA, *ruleB}
	c := &Config{DefaultBackend: defaultBackend, Backends: getBackends(), Rules: configRules}

	r, err := ProvideRules(c)
	assert.NoError(t, err)

	i, match := r.Lookup([]byte(method + pathPrefix))
	assert.True(t, match)

	rules := i.([]*Rule)

	assert.Equal(t, pathPrefix, rules[0].PathPrefix)
	assert.Equal(t, method, rules[0].Method)
	assert.Equal(t, ruleB.Hostname, rules[0].Hostname)

	assert.Equal(t, pathPrefix, rules[1].PathPrefix)
	assert.Equal(t, method, rules[1].Method)
	assert.Equal(t, ruleA.Hostname, rules[1].Hostname)
}

func TestWalk(t *testing.T) {
	rule := &ConfigRule{
		PathPrefix: "/health",
		Method:     "GET",
	}

	configRules := []ConfigRule{*rule}
	c := &Config{DefaultBackend: defaultBackend, Backends: getBackends(), Rules: configRules}

	rules, err := ProvideRules(c)
	assert.NoError(t, err)

	rules.Walk(func(v interface{}) bool {
		r := v.([]*Rule)[0]
		assert.Equal(t, rule.PathPrefix, r.PathPrefix)
		assert.Equal(t, rule.Method, r.Method)
		return false
	})
}
