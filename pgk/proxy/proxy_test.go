package proxy

import (
	"github.com/dhis2-sre/go-rate-limite/pgk/config"
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
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(proxy.TransparentProxyHandler)

	handler.ServeHTTP(recorder, req)

	if actual := recorder.Code; actual != expected {
		t.Errorf("handler returned wrong actual code: got %v want %v", actual, expected)
	}
}
