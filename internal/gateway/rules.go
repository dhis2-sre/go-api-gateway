package gateway

import (
	"errors"
	"net/http"
	"net/url"
	"sort"

	"github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth/v6/limiter"
	iradix "github.com/hashicorp/go-immutable-radix"
)

func NewRules(c *Config) (*rules, error) {
	backendMap, err := mapBackends(c)
	if err != nil {
		return &rules{}, err
	}

	ruleMap, err := mapRules(c, backendMap)
	if err != nil {
		return &rules{}, err
	}

	ruleTree := iradix.New()
	for key, rules := range ruleMap {
		// Sort by hostname length ensuring we're matching against the longest hostname first
		sort.Slice(rules, func(i, j int) bool {
			return len(rules[i].Hostname) > len(rules[j].Hostname)
		})
		ruleTree, _, _ = ruleTree.Insert([]byte(key), rules)
	}

	return &rules{ruleTree}, nil
}

type rules struct {
	rulesTree *iradix.Tree
}

func (r rules) Lookup(key []byte) (interface{}, bool) {
	_, i, match := r.rulesTree.Root().LongestPrefix(key)
	return i, match
}

func (r rules) Len() int {
	return r.rulesTree.Len()
}

type walkFn func(v interface{}) bool

func (r rules) Walk(fn walkFn) {
	r.rulesTree.Root().Walk(func(_ []byte, v interface{}) bool {
		return fn(v)
	})
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

		if c.BasePath != "" {
			configRule.PathPrefix = c.BasePath + configRule.PathPrefix
			// don't add trailing / if it's a catch-all rule
			if configRule.PathPrefix == c.BasePath+"/" {
				configRule.PathPrefix = c.BasePath
			}
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
			key := configRule.Method + configRule.PathPrefix
			ruleMap[key] = append(ruleMap[key], rule)
		} else {
			for _, method := range httpMethods {
				key := method + configRule.PathPrefix
				ruleMap[key] = append(ruleMap[key], rule)
			}
		}
	}
	return ruleMap, nil
}

func newHandler(rule ConfigRule, backendUrl *url.URL) (http.Handler, error) {
	transparentProxy := newTransparentProxy(backendUrl)
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
