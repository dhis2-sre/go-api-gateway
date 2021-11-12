package handler

import (
	"github.com/dhis2-sre/go-rate-limite/pgk/proxy"
	"github.com/dhis2-sre/go-rate-limite/pgk/rule"
	"net/http"
)

func ProvideHandler(rules *rule.Rules, proxy *proxy.Proxy) Handler {
	return Handler{rules, proxy}
}

type Handler struct {
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
