package rule

import (
	"github.com/dhis2-sre/go-rate-limiter/pgk/config"
	"github.com/dhis2-sre/go-rate-limiter/pgk/proxy"
	"github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth/v6/limiter"
	"log"
	"net/http"
	"net/url"
)

func ProvideRules(c *config.Config) *Rules {
	var rules []*Rule
	for _, rule := range c.Rules {
		lmt := newLimiter(rule)

		backend, err := url.Parse(rule.Backend)
		if err != nil {
			log.Fatal(err)
		}

		rules = append(rules, &Rule{
			Rule:    rule,
			Handler: tollbooth.LimitFuncHandler(lmt, proxy.ProvideTransparentProxy(backend)),
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

func (r Rules) Match(req *http.Request) (bool, *Rule) {
	for _, rule := range r.Rules {
		if rule.match(req) {
			return true, rule
		}
	}
	return false, nil
}
