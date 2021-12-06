package gateway

import (
	"github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth/v6/limiter"
	"github.com/hashicorp/go-immutable-radix"
	"net/http"
	"net/url"
)

func ProvideRouter(c *Config) (*Router, error) {
	rules := iradix.New()

	backendMap, err := mapBackends(c)
	if err != nil {
		return nil, err
	}

	ruleMap, catchAllRule, err := mapRules(c, backendMap)
	if err != nil {
		return nil, err
	}

	for path, r := range ruleMap {
		rules, _, _ = rules.Insert([]byte(path), r)
	}

	return &Router{
		Rules:        rules,
		CatchAllRule: catchAllRule,
	}, nil
}

func mapBackends(c *Config) (map[string]*url.URL, error) {
	backendMap := map[string]*url.URL{}
	for _, backend := range c.Backends {
		backendUrl, err := url.Parse(backend.Url)
		if err != nil {
			return nil, err
		}
		backendMap[backend.Name] = backendUrl
	}
	return backendMap, nil
}

func mapRules(c *Config, backendMap map[string]*url.URL) (map[string][]*Rule, *Rule, error) {
	httpMethods := []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE"}

	ruleMap := map[string][]*Rule{}
	var catchAll *Rule = nil

	for _, configRule := range c.Rules {

		// Use default backend if rule doesn't specify one
		if configRule.Backend == "" {
			configRule.Backend = c.DefaultBackend
		}

		// Prefix with basePath if it's defined
		if c.BasePath != "" {
			configRule.PathPrefix = c.BasePath + configRule.PathPrefix
		}

		// Create handler
		backendUrl := backendMap[configRule.Backend]
		handler, err := newHandler(configRule, backendUrl)
		if err != nil {
			return nil, nil, err
		}

		rule := &Rule{
			ConfigRule: configRule,
			Handler:    handler,
		}

		// Detect catch all rule
		if configRule.PathPrefix == c.BasePath+"/" {
			catchAll = rule
			continue
		}

		// Only method and path prefix is indexed in the radix tree, so we might have multiple rules with overlapping which only differs based on headers
		if configRule.Method != "" {
			key := configRule.Method + configRule.PathPrefix
			ruleMap[key] = append(ruleMap[key], rule)
		} else {
			for _, method := range httpMethods {
				key := method + configRule.PathPrefix
				ruleMap[key] = append(ruleMap[key], rule)
			}
		}
	}
	return ruleMap, catchAll, nil
}

func newHandler(rule ConfigRule, backendUrl *url.URL) (http.Handler, error) {
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
	Rules        *iradix.Tree
	CatchAllRule *Rule
}

func (r Router) match(req *http.Request) (bool, *Rule) {
	match, rule := r.matchRule(req)
	if match {
		return true, rule
	}

	if r.CatchAllRule != nil && r.matchMethod(r.CatchAllRule, req) && r.matchHeaders(r.CatchAllRule, req) {
		return true, r.CatchAllRule
	}

	return false, nil
}

func (r Router) matchRule(req *http.Request) (bool, *Rule) {
	key := req.Method + req.URL.Path
	_, i, match := r.Rules.Root().LongestPrefix([]byte(key))
	if match {
		rules := i.([]*Rule)
		for _, rule := range rules {
			if r.matchHeaders(rule, req) {
				return true, rule
			}
		}
	}
	return false, nil
}

func (r Router) matchMethod(rule *Rule, req *http.Request) bool {
	return rule.Method == "" || req.Method == rule.Method
}

func (r Router) matchHeaders(rule *Rule, req *http.Request) bool {
	for ruleHeader := range rule.Headers {
		requestHeaderValues, exists := req.Header[ruleHeader]
		if !exists {
			return false
		}

		for _, ruleHeaderValue := range rule.Headers[ruleHeader] {
			if !stringInSlice(ruleHeaderValue, requestHeaderValues) {
				return false
			}
		}
	}
	return true
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
