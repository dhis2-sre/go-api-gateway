package gateway

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestTransparentProxyHandler(t *testing.T) {
	expected := http.StatusOK

	backend := "http://backend0:8080"

	backendUrl, err := url.Parse(backend)
	assert.NoError(t, err)

	proxy := provideTransparentProxy(backendUrl)

	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	proxy.ServeHTTP(recorder, req)

	actual := recorder.Code
	assert.Equal(t, expected, actual)
}
