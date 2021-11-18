package proxy

import (
	"github.com/dhis2-sre/go-rate-limiter/pgk/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTransparentProxyHandler(t *testing.T) {
	expected := http.StatusOK

	c := &config.Config{
		Backend: "https://httpbin.org",
	}

	proxy := ProvideProxy(c)

	req, err := http.NewRequest("GET", "/status/200", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(proxy.TransparentProxyHandler)

	handler.ServeHTTP(recorder, req)

	actual := recorder.Code
	assert.Equal(t, expected, actual)
}
