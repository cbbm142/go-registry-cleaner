package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gotest.tools/assert"
)

func TestApiDefaults(t *testing.T) {
	req := apiReqData{
		url: "test/rl",
	}
	req.apiDefaults()
	assert.Assert(t, req.method == "GET")
}

func TestApiRequest(t *testing.T) {
	requestTypes := []struct {
		method, response, body  string
		header, headerDirective []string
	}{
		{"GET", "ok", "", nil, nil},
		{"GET", "ok", "", []string{"FOO", "auth"}, []string{"BAR", "bearer"}},
		{"DELETE", "accepted", "", []string{"Accept"}, []string{"application/vnd.docker.distribution.manifest.v2+json"}},
	}

	for _, request := range requestTypes {
		srv := serverMock("/", request.response)
		req := apiReqData{
			url:             srv.URL,
			header:          request.header,
			headerDirective: request.headerDirective,
			body:            request.body,
			method:          request.method,
		}
		resp := req.apiRequest()
		assert.Assert(t, len(resp.Request.Header) == len(req.header))
		assert.Assert(t, resp.Request.Method == req.method)
		assert.Assert(t, strings.Contains(strings.ToLower(resp.Status), request.response))
		srv.Close()
	}

}

func TestDeleteDigest(t *testing.T) {

	digest := "sha:4252435"
	tag := "1"
	dryRun := false
	name := "/v2/test/image"
	srv := serverMock(fmt.Sprintf("%s/manifests/%s", name, digest), "accepted")
	defer srv.Close()

	err := deleteDigest(srv.URL+name, digest, tag, dryRun)
	assert.Assert(t, err == nil)
}

func serverMock(url string, header string) *httptest.Server {
	handler := http.NewServeMux()
	switch header {
	case "accepted":
		handler.HandleFunc(url, mockAccepted)
	case "ok":
		handler.HandleFunc(url, mockOk)
	}

	srv := httptest.NewServer(handler)

	return srv
}

// Available Status Codes https://golang.org/src/net/http/status.go
func mockAccepted(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)

}

func mockOk(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

}
