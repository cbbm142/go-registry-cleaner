package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gotest.tools/assert"
)

type blob struct {
	BlobSum string `json:"blobSum"`
}
type testBody struct {
	Name     string        `json:"name"`
	Tags     []string      `json:"tags"`
	Fslayers []interface{} `json:"fsLayers"`
	History  []interface{} `json:"history"`
}
type mockV1Compat struct {
	V1Compatibility string `json:"v1Compatibility"`
}

var mockTags []string = []string{"test", "prod", "miracle"}
var name string = "/v2/test/repo"
var mockSha string = "sha256:6c3c624b58dbbcd3c0dd82b4c53f04194d1247c6eebdaab7c610cf7d66709b3b"

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
	digest := mockSha
	tag := "1"
	name := "/v2/test/image"
	var err error
	var responseCode string
	testsToRun := []struct {
		dryRun      bool
		expectedErr bool
	}{
		{dryRun: true, expectedErr: false},
		{dryRun: false, expectedErr: false},
		{dryRun: false, expectedErr: true},
	}
	for _, testCase := range testsToRun {
		if testCase.dryRun || testCase.expectedErr {
			responseCode = "ok"
		} else {
			responseCode = "accepted"
		}
		srv := serverMock(fmt.Sprintf("%s/manifests/%s", name, digest), responseCode)
		err = deleteDigest(srv.URL+name, digest, tag, testCase.dryRun)
		srv.Close()
		if testCase.expectedErr {
			assert.Assert(t, err != nil)
		} else {
			assert.Assert(t, err == nil)
		}
	}
}

func TestGetImageDigest(t *testing.T) {
	var tag string = "testTag:12.2"
	srv := serverMock(fmt.Sprintf("%s/manifests/%s", name, tag), "digest")
	defer srv.Close()
	digest := getImageDigest(srv.URL+name, tag)
	assert.Assert(t, digest == mockSha)
}

func TestGetTags(t *testing.T) {

	srv := serverMock(fmt.Sprintf("%s/tags/list/", name), "body")
	defer srv.Close()
	returnedTags := getTags(srv.URL + name)

	assert.Assert(t, len(returnedTags) == 3)
	tagPresent := make(map[string]bool)
	for _, tag := range mockTags {
		tagPresent[tag] = false
	}
	for _, tag := range returnedTags {
		tagPresent[tag] = true
	}
	for _, tag := range mockTags {
		assert.Assert(t, tagPresent[tag] == true)
	}

}

func TestGetManifest(t *testing.T) {
	var mockTag string = "testTag:1.2.3"
	srv := serverMock(fmt.Sprintf("%s/manifests/%s", name, mockTag), "body")
	defer srv.Close()
	mockManifest, mockCreation := getManifest(srv.URL+name, mockTag)
	assert.Assert(t, mockManifest == mockSha)
	assert.Assert(t, mockCreation == "2014-12-31")
	fmt.Println(mockManifest + mockCreation)
}

func serverMock(url string, header string) *httptest.Server {
	handler := http.NewServeMux()
	switch header {
	case "accepted":
		handler.HandleFunc(url, mockAccepted)
	case "ok":
		handler.HandleFunc(url, mockOk)
	case "digest":
		handler.HandleFunc(url, returnDigest)
	case "body":
		handler.HandleFunc(url, returnBody)
	}

	srv := httptest.NewServer(handler)

	return srv
}

func createMockBody() testBody {

	var mockBlobSum blob = blob{
		BlobSum: mockSha,
	}
	rawV1 := `"id":"e45a5af57b00862e5ef5782a9925979a02ba2b12dff832fd0991335f4a11e5c5","parent":"31cbccb51277105ba3ae35ce33c22b69c9e3f1002e76e4c736a2e8ebff9d7b5d","created":"2014-12-31T22:57:59.178729048Z","container":"27b45f8fb11795b52e9605b686159729b0d9ca92f76d40fb4f05a62e19c46b4f","container_config":{"Hostname":"8ce6509d66e2","Domainname":"","User":"","Memory":0,"MemorySwap":0,"CpuShares":0,"Cpuset":"","AttachStdin":false,"AttachStdout":false,"AttachStderr":false,"PortSpecs":null,"ExposedPorts":null,"Tty":false,"OpenStdin":false,"StdinOnce":false,"Env":["PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"],"Cmd":["/bin/sh","-c","#(nop) CMD [/hello]"],"Image":"31cbccb51277105ba3ae35ce33c22b69c9e3f1002e76e4c736a2e8ebff9d7b5d","Volumes":null,"WorkingDir":"","Entrypoint":null,"NetworkDisabled":false,"MacAddress":"","OnBuild":[],"SecurityOpt":null,"Labels":null},"docker_version":"1.4.1","config":{"Hostname":"8ce6509d66e2","Domainname":"","User":"","Memory":0,"MemorySwap":0,"CpuShares":0,"Cpuset":"","AttachStdin":false,"AttachStdout":false,"AttachStderr":false,"PortSpecs":null,"ExposedPorts":null,"Tty":false,"OpenStdin":false,"StdinOnce":false,"Env":["PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"],"Cmd":["/hello"],"Image":"31cbccb51277105ba3ae35ce33c22b69c9e3f1002e76e4c736a2e8ebff9d7b5d","Volumes":null,"WorkingDir":"","Entrypoint":null,"NetworkDisabled":false,"MacAddress":"","OnBuild":[],"SecurityOpt":null,"Labels":null},"architecture":"amd64","os":"linux","Size":0`
	jsonV1, _ := json.Marshal(rawV1)
	var mockV1 mockV1Compat = mockV1Compat{
		V1Compatibility: string(jsonV1),
	}
	var mockHistroy []interface{} = []interface{}{mockV1}
	// mockHistory is just to ensure the correct data
	var mockFslayer []interface{} = []interface{}{mockBlobSum, mockHistroy}
	var mockBody testBody = testBody{
		Name:     "test/repo",
		Tags:     mockTags,
		Fslayers: mockFslayer,
		History:  mockHistroy,
	}
	return mockBody
}

// Available Status Codes https://golang.org/src/net/http/status.go
func mockAccepted(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)

}

func mockOk(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

}

func returnDigest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Docker-Content-Digest", mockSha)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func returnBody(w http.ResponseWriter, r *http.Request) {
	var mockBody testBody = createMockBody()
	jsonMessage, _ := json.Marshal(mockBody)
	_, err := w.Write([]byte(jsonMessage))
	errCheck(err)

}
