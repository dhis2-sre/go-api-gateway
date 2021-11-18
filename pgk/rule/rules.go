package rule

import (
	"github.com/dhis2-sre/go-rate-limite/pgk/config"
	"github.com/dhis2-sre/go-rate-limite/pgk/proxy"
	"github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth/v6/limiter"
	"net/http"
)

func ProvideRules(c *config.Config) *Rules {
	var rules []*Rule
	for _, rule := range c.Rules {
		lmt := newLimiter(rule)

		p := proxy.ProvideProxy(c)
		rules = append(rules, &Rule{
			Rule:    rule,
			Handler: tollbooth.LimitFuncHandler(lmt, p.TransparentProxyHandler),
		})
	}

	return &Rules{
		Rules: rules,
	}
}

func newLimiter(rule config.Rule) *limiter.Limiter {
	lmt := tollbooth.NewLimiter(rule.RequestPerSecond, nil)
	if rule.Method != "" {
		lmt.SetMethods([]string{rule.Method})
	}
	lmt.SetBurst(rule.Burst)
	return lmt
}

type Rules struct {
	Rules []*Rule
}

func (r Rules) Match(req *http.Request) (bool, http.Handler) {
	for _, rule := range r.Rules {
		if rule.match(req) {
			return true, rule.Handler
		}
	}
	return false, nil
}
