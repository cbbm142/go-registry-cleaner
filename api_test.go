package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
	srv := serverMock("/", "ok")
	defer srv.Close()

	req := apiReqData{
		url: srv.URL,
	}
	resp := req.apiRequest()
	assert.Assert(t, len(resp.Request.Header) == 0)
	assert.Assert(t, resp.Request.Method == "GET")

	req.header = "FOO"
	req.headerDirective = "Bar"
	resp = req.apiRequest()
	assert.Assert(t, len(resp.Request.Header) == 1)
	assert.Assert(t, resp.Request.Method == "GET")

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
