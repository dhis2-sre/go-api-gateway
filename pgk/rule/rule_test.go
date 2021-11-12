package rule

import (
	"github.com/dhis2-sre/go-rate-limite/pgk/config"
	"github.com/stretchr/testify/assert"
	"testing"
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
