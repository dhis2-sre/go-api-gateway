package rule

import (
	"github.com/dhis2-sre/go-rate-limite/pgk/config"
	"net/http"
	"regexp"
)

type Rule struct {
	config.Rule
	Handler http.Handler
}

func (r *Rule) pathMatch(path string) bool {
	match, err := regexp.MatchString(r.PathPattern, path)
	if err != nil {
		return false
	}
	return match
}
