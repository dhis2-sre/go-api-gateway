package gateway

import (
	"errors"
	"github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth/v6/limiter"
	"github.com/hashicorp/go-immutable-radix"
	"log"
	"net"
	"net/http"
	"net/url"
)

func ProvideRouter(c *Config) (*Router, error) {
	backendMap, err := mapBackends(c)
	if err != nil {
		return nil, err
	}

	ruleMap, err := mapRules(c, backendMap)
	if err != nil {
		return nil, err
	}

	ruleTree := iradix.New()
	for key, rules := range ruleMap {
		ruleTree, _, _ = ruleTree.Insert([]byte(key), rules)
	}

	return &Router{
		Rules: ruleTree,
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

func mapRules(c *Config, backendMap map[string]*url.URL) (map[string][]*Rule, error) {
	httpMethods := []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE"}

	ruleMap := map[string][]*Rule{}
	for _, configRule := range c.Rules {

		// Use default backend if rule doesn't specify one
		if configRule.Backend == "" {
			configRule.Backend = c.DefaultBackend
		}

		if configRule.Backend == "" {
			return nil, errors.New("either a rule needs to define a backend or a default backend needs to be defined")
		}

		// Prefix with basePath if it's defined
		if c.BasePath != "" {
			configRule.PathPrefix = c.BasePath + configRule.PathPrefix
		}

		// Create handler
		backendUrl := backendMap[configRule.Backend]
		if backendUrl == nil {
			return nil, errors.New("backend map contains not entry for: " + configRule.Backend)
		}
		handler, err := newHandler(configRule, backendUrl)
		if err != nil {
			return nil, err
		}

		rule := &Rule{
			ConfigRule: configRule,
			Handler:    handler,
		}

		if configRule.Method != "" {
			key := configRule.Hostname + configRule.Method + configRule.PathPrefix
			ruleMap[key] = append(ruleMap[key], rule)
		} else {
			for _, method := range httpMethods {
				key := configRule.Hostname + method + configRule.PathPrefix
				ruleMap[key] = append(ruleMap[key], rule)
			}
		}
	}
	return ruleMap, nil
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
	Rules *iradix.Tree
}

func (r Router) match(req *http.Request) (bool, *Rule) {
	hostname := r.getHostname(req)
	key := hostname + req.Method + req.URL.Path
	_, i, match := r.Rules.Root().LongestPrefix([]byte(key))

	if !match {
		key := req.Method + req.URL.Path
		_, i, match = r.Rules.Root().LongestPrefix([]byte(key))
	}

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

func (r Router) getHostname(req *http.Request) string {
	hostname, _, err := net.SplitHostPort(req.Host)
	if err != nil {
		// TODO:
		log.Fatalln(err)
	}
	return hostname
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
