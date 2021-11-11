package handler

import (
	"github.com/dhis2-sre/go-rate-limite/pgk/config"
	"github.com/dhis2-sre/go-rate-limite/pgk/proxy"
	"github.com/dhis2-sre/go-rate-limite/pgk/rule"
	"net/http"
)

func ProvideHandler(c config.Config, rules *rule.Rules, proxy *proxy.Proxy) Handler {
	return Handler{c, rules, proxy}
}

type Handler struct {
	c     config.Config
	rules *rule.Rules
	proxy *proxy.Proxy
}

func (h *Handler) RateLimitingProxyHandler(w http.ResponseWriter, r *http.Request) {
	if match, rateLimitingProxyHandler := h.rules.Match(r); match {
		rateLimitingProxyHandler.ServeHTTP(w, r)
	} else {
		h.proxy.TransparentProxyHandler(w, r)
	}
}
