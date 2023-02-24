package gateway

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getBackends() []Backend {
	return []Backend{
		{Name: "backend0", Url: "http://backend0:8080"},
		{Name: "backend1", Url: "http://backend1:8080"},
	}
}

const defaultBackend = "backend0"
const defaultRequestUrl = "http://url:80"

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtYrBsSkVGXZKQL13lbmd
xFCQcvi6KIssjz3KOHIko/Da6sxE2w67OL84t98wCYbmIuq6xTK6qpEqEs1LaqQS
DnCs2VNDTLk4D1J42R63OpJQfOfebzhTJLx6KldyK2FRGXWILY7AzcoqyuLk433s
lHk6/yFDYgBA4COofeXZvXtUazuzpBWTZCxpEh341ob6XQ5juLYrqr/80XLYzXiu
N1iz24ulxSnD0GV4cRfHEnnzN3oYFzoYTcTQB6dffNAs/ADHNA9IemyLbT0ugvbf
L5MOEBOftYLRwmGFWrXf5s9jccku0FPid2wtZEwsv5Sa+Yvr36KHtrr+PSFksOB1
0QIDAQAB
-----END PUBLIC KEY-----`

const validAccessToken = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQ3OTExNTQ0MTUsImlhdCI6MTYzNzU1NDQxNX0.PtQp6_k5bQ9KE9uk520i4emVnUmxFD8DxyeZsfzgT6CY2oMyXEm7zlIA-4_xz2Q7CrSeqnWxpy0coK9MN0EPE2vhFomTrP6D3l7_lX6Dyn1gH6zWpjC_dRqOSRv3AqS3buZiC-vNwCatLhu6WE74cykBAE2veIr8Gp_ebiITXJKiHBNaTlPk2WEfcJ1NL3g7nafy6l-V4h2-Vj3tapJQiLfpgReIXYIswFYH7En7qy94fL0eOUbZzQI9fOuiXvAN-owR3GYcbwz9Hll23VACWsekMJdDBEgUSdek9JOmRHGxko6FE79-_ClYvF1dGUgZB2mDwY_xF2TOG2q3XDi9Aw`

func TestHandler(t *testing.T) {
	expected := http.StatusOK

	rule := ConfigRule{
		PathPrefix: "/health",
	}

	configRules := []ConfigRule{rule}
	c := &Config{DefaultBackend: defaultBackend, Backends: getBackends(), Rules: configRules}

	rules, err := NewRules(c)
	assert.NoError(t, err)

	router := NewRouter(rules)

	jwtAuth := NewJwtAuth(c)

	handler := NewHandler(c, router, jwtAuth)

	req, err := http.NewRequest(http.MethodGet, defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	actual := recorder.Code
	assert.Equal(t, expected, actual)
}

func TestMaxMultipart(t *testing.T) {
	// TODO:
	/*
		expected := http.StatusOK

		rule := ConfigRule{
			PathPrefix: "/health",
		}

		configRules := []ConfigRule{rule}
		c := &Config{DefaultBackend: defaultBackend, MaxMultipartSize: 2, Backends: getBackends(), Rules: configRules}

		rules, err := NewRules(c)
		assert.NoError(t, err)

		router := NewRouter(rules)

		jwtAuth := NewJwtAuth(c)

		handler := NewHandler(c, router, jwtAuth)

		var b bytes.Buffer
		w := multipart.NewWriter(&b)

		file, err := ioutil.TempFile("/tmp", "go-api-gateway-test")
		assert.NoError(t, err)

		defer func(name string) {
			err := os.Remove(name)
			assert.NoError(t, err)
		}(file.Name())

		data := make([]byte, 20<<20)
		_, err = file.Write(data)
		assert.NoError(t, err)

		var fw io.Writer
		if fw, err = w.CreateFormFile("whatever", file.Name()); err != nil {
			assert.NoError(t, err)
		}

		if _, err = io.Copy(fw, file); err != nil {
			assert.NoError(t, err)
		}

		req, err := http.NewRequest(http.MethodGet, defaultRequestUrl+"/health", &b)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", w.FormDataContentType())

		err = w.Close()
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		actual := recorder.Code
		assert.Equal(t, expected, actual)
	*/
}

func TestHandlerBlock(t *testing.T) {
	expected := http.StatusForbidden

	rule := ConfigRule{
		PathPrefix: "/health",
		Block:      true,
	}

	configRules := []ConfigRule{rule}
	c := &Config{DefaultBackend: defaultBackend, Backends: getBackends(), Rules: configRules}

	rules, err := NewRules(c)
	assert.NoError(t, err)

	router := NewRouter(rules)

	jwtAuth := NewJwtAuth(c)

	handler := NewHandler(c, router, jwtAuth)

	req, err := http.NewRequest(http.MethodGet, defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	actual := recorder.Code
	assert.Equal(t, expected, actual)
}

func TestHandlerBlockFalse(t *testing.T) {
	expected := http.StatusOK

	rule := ConfigRule{
		PathPrefix: "/health",
		Block:      false,
	}

	configRules := []ConfigRule{rule}
	c := &Config{DefaultBackend: defaultBackend, Backends: getBackends(), Rules: configRules}

	rules, err := NewRules(c)
	assert.NoError(t, err)

	router := NewRouter(rules)

	jwtAuth := NewJwtAuth(c)

	handler := NewHandler(c, router, jwtAuth)

	req, err := http.NewRequest(http.MethodGet, defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	actual := recorder.Code
	assert.Equal(t, expected, actual)
}

func TestHandlerRateLimited(t *testing.T) {
	rule := ConfigRule{
		PathPrefix:       "/health",
		RequestPerSecond: 1,
		Burst:            1,
	}

	configRules := []ConfigRule{rule}
	c := &Config{DefaultBackend: defaultBackend, Backends: getBackends(), Rules: configRules}

	rules, err := NewRules(c)
	assert.NoError(t, err)

	router := NewRouter(rules)

	jwtAuth := NewJwtAuth(c)

	handler := NewHandler(c, router, jwtAuth)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/health", nil)
	assert.NoError(t, err)

	client := &http.Client{}

	response0, err := client.Do(req)
	assert.NoError(t, err)
	defer response0.Body.Close()

	response1, err := client.Do(req)
	assert.NoError(t, err)
	defer response1.Body.Close()

	actual0 := response0.StatusCode
	expected0 := http.StatusOK
	assert.Equal(t, expected0, actual0)

	actual1 := response1.StatusCode
	expected1 := http.StatusTooManyRequests
	assert.Equal(t, expected1, actual1)
}

func TestHandlerUserAgentHeader(t *testing.T) {
	expected := http.StatusOK

	rule := ConfigRule{
		PathPrefix: "/health",
		Headers: map[string][]string{
			"User-Agent": {"Go tests"},
		},
	}

	configRules := []ConfigRule{rule}
	c := &Config{DefaultBackend: defaultBackend, Backends: getBackends(), Rules: configRules}

	rules, err := NewRules(c)
	assert.NoError(t, err)

	router := NewRouter(rules)

	jwtAuth := NewJwtAuth(c)

	handler := NewHandler(c, router, jwtAuth)

	req, err := http.NewRequest(http.MethodGet, defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	req.Header.Set("User-Agent", "Go tests")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	actual := recorder.Code
	assert.Equal(t, expected, actual)
}

func TestHandlerUserAgentHeaderNoMatch(t *testing.T) {
	expected := http.StatusMisdirectedRequest

	rule := ConfigRule{
		PathPrefix: "/health",
		Headers: map[string][]string{
			"User-Agent": {"Go tests"},
		},
	}

	configRules := []ConfigRule{rule}
	c := &Config{DefaultBackend: defaultBackend, Backends: getBackends(), Rules: configRules}

	rules, err := NewRules(c)
	assert.NoError(t, err)

	router := NewRouter(rules)

	jwtAuth := NewJwtAuth(c)

	handler := NewHandler(c, router, jwtAuth)

	req, err := http.NewRequest(http.MethodGet, defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0)")

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	actual := recorder.Code
	assert.Equal(t, expected, actual)
}

func TestHandlerNoMatch(t *testing.T) {
	expected := http.StatusMisdirectedRequest

	rule := ConfigRule{
		PathPrefix: "/health",
	}

	configRules := []ConfigRule{rule}
	c := &Config{DefaultBackend: defaultBackend, Backends: getBackends(), Rules: configRules}

	rules, err := NewRules(c)
	assert.NoError(t, err)

	router := NewRouter(rules)

	jwtAuth := NewJwtAuth(c)

	handler := NewHandler(c, router, jwtAuth)

	req, err := http.NewRequest(http.MethodGet, defaultRequestUrl+"/no-match", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	actual := recorder.Code
	assert.Equal(t, expected, actual)
}

func TestHandlerJwtAuthentication(t *testing.T) {
	expected := http.StatusOK

	rule := ConfigRule{
		PathPrefix:     "/health",
		Authentication: "jwt",
	}

	configRules := []ConfigRule{rule}

	c := &Config{
		Backends:       getBackends(),
		DefaultBackend: defaultBackend,
		Authentication: Authentication{Jwt: Jwt{publicKey}},
		Rules:          configRules,
	}

	rules, err := NewRules(c)
	assert.NoError(t, err)

	router := NewRouter(rules)

	jwtAuth := NewJwtAuth(c)

	handler := NewHandler(c, router, jwtAuth)

	req, err := http.NewRequest(http.MethodGet, defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", validAccessToken)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	actual := recorder.Code
	assert.Equal(t, expected, actual)
}

func TestHandlerJwtAuthenticationInvalidToken(t *testing.T) {
	expected := http.StatusForbidden

	rule := ConfigRule{
		PathPrefix:     "/health",
		Authentication: "jwt",
	}

	configRules := []ConfigRule{rule}

	c := &Config{
		Backends:       getBackends(),
		DefaultBackend: defaultBackend,
		Authentication: Authentication{Jwt: Jwt{publicKey}},
		Rules:          configRules,
	}

	rules, err := NewRules(c)
	assert.NoError(t, err)

	router := NewRouter(rules)

	jwtAuth := NewJwtAuth(c)

	handler := NewHandler(c, router, jwtAuth)

	req, err := http.NewRequest(http.MethodGet, defaultRequestUrl+"/health", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "bla bla")

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	actual := recorder.Code
	assert.Equal(t, expected, actual)
}

func TestHandlerPathReplacePostfix(t *testing.T) {
	expected := http.StatusOK

	rule := ConfigRule{
		PathPrefix: "/health/backend0",
		PathReplace: PathReplace{
			Target:      "/backend0",
			Replacement: "",
		},
	}

	configRules := []ConfigRule{rule}
	c := &Config{DefaultBackend: defaultBackend, Backends: getBackends(), Rules: configRules}

	rules, err := NewRules(c)
	assert.NoError(t, err)

	router := NewRouter(rules)

	jwtAuth := NewJwtAuth(c)

	handler := NewHandler(c, router, jwtAuth)

	req, err := http.NewRequest(http.MethodGet, defaultRequestUrl+"/health/backend0", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	actual := recorder.Code
	assert.Equal(t, expected, actual)
}

func TestHandlerPathReplacePrefix(t *testing.T) {
	expected := http.StatusOK

	rule := ConfigRule{
		PathPrefix: "/backend0/health",
		PathReplace: PathReplace{
			Target:      "/backend0",
			Replacement: "",
		},
	}

	configRules := []ConfigRule{rule}
	c := &Config{DefaultBackend: defaultBackend, Backends: getBackends(), Rules: configRules}

	rules, err := NewRules(c)
	assert.NoError(t, err)

	router := NewRouter(rules)

	jwtAuth := NewJwtAuth(c)

	handler := NewHandler(c, router, jwtAuth)

	req, err := http.NewRequest(http.MethodGet, defaultRequestUrl+"/backend0/health", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	actual := recorder.Code
	assert.Equal(t, expected, actual)
}

func TestHandlerPathReplaceWithReplacement(t *testing.T) {
	expected := http.StatusOK

	rule := ConfigRule{
		PathPrefix: "/something",
		PathReplace: PathReplace{
			Target:      "thing",
			Replacement: "body",
		},
	}

	configRules := []ConfigRule{rule}
	c := &Config{DefaultBackend: defaultBackend, Backends: getBackends(), Rules: configRules}

	rules, err := NewRules(c)
	assert.NoError(t, err)

	router := NewRouter(rules)

	jwtAuth := NewJwtAuth(c)

	handler := NewHandler(c, router, jwtAuth)

	req, err := http.NewRequest(http.MethodGet, defaultRequestUrl+"/something", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	actual := recorder.Code
	assert.Equal(t, expected, actual)
}
