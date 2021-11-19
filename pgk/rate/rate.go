package rate

import (
	"net/http"

	"github.com/dhis2-sre/go-rate-limite/pgk/rule"
)

func Limit(rules *rule.Rules) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ok, h := rules.Match(r); ok {
				h.ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
