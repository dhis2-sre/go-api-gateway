package gateway

import (
	"net/http"
	"regexp"
)

type Rule struct {
	ConfigRule
	Handler http.Handler
}

func (r *Rule) match(req *http.Request) bool {
	match, err := regexp.MatchString(r.PathPattern, req.URL.Path)
	if err != nil {
		return false
	}
	return match && (req.Method == r.Method || r.Method == "")
}