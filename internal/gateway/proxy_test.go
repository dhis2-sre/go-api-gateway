package gateway

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
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
