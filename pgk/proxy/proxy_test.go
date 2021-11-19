package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransparentProxyHandler(t *testing.T) {
	expected := http.StatusOK

	req, err := http.NewRequest("GET", "/status/200", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Transparently("https://httpbin.org").ServeHTTP)

	handler.ServeHTTP(recorder, req)

	actual := recorder.Code
	assert.Equal(t, expected, actual)
}
