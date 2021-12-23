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
