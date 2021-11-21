package gateway

import (
	"net/http"
)

type Rule struct {
	ConfigRule
	Handler http.Handler
}
