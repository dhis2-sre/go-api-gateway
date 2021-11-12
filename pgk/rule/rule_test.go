package rule

import (
	"github.com/dhis2-sre/go-rate-limite/pgk/config"
	"testing"
)

func TestPathMatch(t *testing.T) {
	expectation := true

	rule := createRuleWithPathPattern("^\\/health$")

	actual := rule.pathMatch("/health")

	if actual != expectation {
		t.Errorf("Expected %v but got %v", expectation, actual)
	}
}

func TestPathNoMatch(t *testing.T) {
	expectation := false

	rule := createRuleWithPathPattern("^\\/health$")

	actual := rule.pathMatch("/health-no-match")

	if actual != expectation {
		t.Errorf("Expected %v but got %v", expectation, actual)
	}
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
