package gateway

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransparentProxyHandler(t *testing.T) {
	expected := http.StatusOK

	backend := "http://backend0:8080"

	c := &Config{
		Authentication: Authentication{
			Jwt: Jwt{publicKey},
		},
	}

	auth := NewJwtAuth(c)
	proxy, err := newTransparentProxy(ConfigRule{Backend: backend}, auth)
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	proxy.ServeHTTP(recorder, req)

	actual := recorder.Code
	assert.Equal(t, expected, actual)
}
