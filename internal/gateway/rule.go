package gateway

import (
	"net/http"
	"strings"
)

type Rule struct {
	ConfigRule
	Handler http.Handler
}

func (r *Rule) match(req *http.Request) bool {
	match := strings.HasPrefix(req.URL.Path, r.PathPrefix)
	return match && (req.Method == r.Method || r.Method == "")
}
