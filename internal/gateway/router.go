package gateway

import (
	"github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth/v6/limiter"
	"net/http"
	"net/url"
)

func ProvideRouter(c *Config) (*Router, error) {
	var rules []*Rule
	for _, rule := range c.Rules {
		lmt := newLimiter(rule)

		backend, err := url.Parse(rule.Backend)
		if err != nil {
			return nil, err
		}

		rules = append(rules, &Rule{
			ConfigRule: rule,
			Handler:    tollbooth.LimitFuncHandler(lmt, ProvideTransparentProxy(backend)),
		})
	}

	return &Router{
		Rules: rules,
	}, nil
}

func newLimiter(rule ConfigRule) *limiter.Limiter {
	lmt := tollbooth.NewLimiter(rule.RequestPerSecond, nil)
	if rule.Method != "" {
		lmt.SetMethods([]string{rule.Method})
	}
	lmt.SetBurst(rule.Burst)
	return lmt
}

type Router struct {
	Rules []*Rule
}

func (r Router) Match(req *http.Request) (bool, *Rule) {
	for _, rule := range r.Rules {
		if rule.match(req) {
			return true, rule
		}
	}
	return false, nil
}
