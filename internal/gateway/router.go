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

		handler, err := newHandler(c.DefaultBackend, rule)
		if err != nil {
			return nil, err
		}

		if c.BasePath != "" {
			rule.PathPrefix = c.BasePath + rule.PathPrefix
		}

		rules = append(rules, &Rule{
			ConfigRule: rule,
			Handler:    handler,
		})
	}

	return &Router{
		Rules: rules,
	}, nil
}

func newHandler(defaultBackend string, rule ConfigRule) (http.Handler, error) {
	backend := defaultBackend

	if rule.Backend != "" {
		backend = rule.Backend
	}

	backendUrl, err := url.Parse(backend)
	if err != nil {
		return nil, err
	}

	transparentProxy := provideTransparentProxy(backendUrl)
	handler := http.Handler(transparentProxy)
	if rule.RequestPerSecond != 0 {
		lmt := newLimiter(rule)
		handler = tollbooth.LimitFuncHandler(lmt, transparentProxy)
	}
	return handler, nil
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
