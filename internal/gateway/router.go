package gateway

import (
	"log"
	"net"
	"net/http"
	"strings"
)

func NewRouter(rules Rules) *router {
	return &router{
		Rules: rules,
	}
}

type Rules interface {
	Lookup(key []byte) (interface{}, bool)
	Len() int
	Walk(fn walkFn)
}

type router struct {
	Rules Rules
}

func (r router) match(req *http.Request) (bool, *Rule) {
	key := req.Method + req.URL.Path
	i, match := r.Rules.Lookup([]byte(key))

	if match {
		rules := i.([]*Rule)
		hostname := r.getHostname(req)
		for _, rule := range rules {
			if r.matchHeaders(rule, req) && r.matchHostname(rule, hostname) {
				return true, rule
			}
		}
	}

	return false, nil
}

func (r router) matchHostname(rule *Rule, hostname string) bool {
	return rule.Hostname == "" ||
		(rule.Hostname[0:2] == "*." && strings.HasSuffix(hostname, rule.Hostname[1:])) ||
		hostname == rule.Hostname
}

func (r router) getHostname(req *http.Request) string {
	hostname, _, err := net.SplitHostPort(req.Host)
	if err != nil {
		if strings.HasSuffix(err.Error(), ": missing port in address") {
			return req.Host
		}
		log.Println("Error:", err)
		log.Println("Request:", req)
		return ""
	}
	return hostname
}

func (r router) matchMethod(rule *Rule, req *http.Request) bool {
	return rule.Method == "" || req.Method == rule.Method
}

func (r router) matchHeaders(rule *Rule, req *http.Request) bool {
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
