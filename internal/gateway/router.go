package gateway

import (
	"github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth/v6/limiter"
	"github.com/hashicorp/go-immutable-radix"
	"net/http"
	"net/url"
)

func ProvideRouter(c *Config) (*Router, error) {
	r := iradix.New()

	for _, rule := range c.Rules {

		if rule.Backend == "" {
			rule.Backend = c.DefaultBackend
		}

		if c.BasePath != "" {
			rule.PathPrefix = c.BasePath + rule.PathPrefix
		}

		handler, err := newHandler(rule)
		if err != nil {
			return nil, err
		}

		r, _, _ = r.Insert([]byte(rule.PathPrefix), &Rule{
			ConfigRule: rule,
			Handler:    handler,
		})
	}

	return &Router{
		Rules: r,
	}, nil
}

func newHandler(rule ConfigRule) (http.Handler, error) {
	backendUrl, err := url.Parse(rule.Backend)
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
	Rules *iradix.Tree
}

func (r Router) match(req *http.Request) (bool, *Rule) {
	_, i, match := r.Rules.Root().LongestPrefix([]byte(req.URL.Path))
	if match {
		rule := i.(*Rule)
		if rule.Method == "" || req.Method == rule.Method {
			return true, rule
		}
	}
	return false, nil
}
