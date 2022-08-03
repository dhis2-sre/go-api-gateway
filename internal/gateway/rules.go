package gateway

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth/v6/limiter"
)

func NewRuleList(c *Config) ([]ConfigRule, error) {
	backendMap, err := mapBackends(c)
	if err != nil {
		return nil, nil
	}

	ruleMap, err := mapRules(c, backendMap)
	if err != nil {
		return nil, nil
	}

	return ruleMap, nil
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

func mapRules(c *Config, backendMap map[string]*url.URL) ([]ConfigRule, error) {
	ruleList := make([]ConfigRule, len(c.Rules))
	for i, configRule := range c.Rules {
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

		backendUrl := backendMap[configRule.Backend]
		if backendUrl == nil {
			return nil, errors.New("backend map contains not entry for: " + configRule.Backend)
		}
		configRule.Backend = backendUrl.String()
		ruleList[i] = configRule
	}
	return ruleList, nil
}

func NewHandler2(rule ConfigRule, auth auth) (http.Handler, error) {
	return newHandler(rule, auth)
}

func newHandler(rule ConfigRule, auth auth) (http.Handler, error) {
	transparentProxy, err := newTransparentProxy(rule, auth)
	if err != nil {
		return nil, nil
	}

	handler := http.Handler(transparentProxy)
	if rule.RequestPerSecond != 0 {
		lmt := newLimiter(rule.Method, rule.RequestPerSecond, rule.Burst)
		handler = tollbooth.LimitFuncHandler(lmt, transparentProxy)
	}
	return handler, nil
}

func newLimiter(method string, requestPerSecond float64, burst int) *limiter.Limiter {
	lmt := tollbooth.NewLimiter(requestPerSecond, nil)
	if method != "" {
		lmt.SetMethods([]string{method})
	}
	lmt.SetBurst(burst)
	return lmt
}
