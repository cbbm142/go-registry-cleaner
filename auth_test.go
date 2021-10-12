package main

import (
	"fmt"
	"strings"
	"testing"

	"gotest.tools/assert"
)

func TestAuthDetect(t *testing.T) {
	testCases := []struct {
		response        string
		headerDirective string
		httpsProtocol   bool
		expectedErr     bool
		expectedToken   string
		expectedAuth    bool
	}{
		{"ok", "", false, false, "", false},
		{"basicAuthRequired", "Basic", false, false, "", true},
		{"badRequest", "", false, true, "", false},
		// SSL will fail, but we will check it attempted https later
		{"ok", "", true, true, "", false},
	}
	for _, test := range testCases {
		// Ensure these are at default values before every iteration
		bearerToken = ""
		basicAuth = false

		var hostMap = make(map[interface{}]interface{})
		srv := serverMock("/v2", test.response)
		hostMap["host"] = strings.Split(srv.URL, "//")[1] + "/v2"
		var mockConfig config = config{
			configMap:     hostMap,
			httpsProtocol: test.httpsProtocol,
		}

		err := authDetect(mockConfig)
		assert.Assert(t, bearerToken == test.expectedToken)
		assert.Assert(t, basicAuth == test.expectedAuth)
		if test.expectedErr {
			assert.Assert(t, err != nil)
			if test.httpsProtocol {
				assert.Assert(t, strings.Contains(err.Error(), "http: server gave HTTP response to HTTPS client"))
			}
		} else {
			assert.Assert(t, err == nil)
		}
		srv.Close()
	}
}

func TestGetToken(t *testing.T) {
	srv := serverMock("/test", "body")
	defer srv.Close()
	var mockAuthResponse authResponse = authResponse{
		realm:   srv.URL + "/test",
		service: "registry",
		scope:   "mockscope",
	}
	mockToken := getToken(mockAuthResponse, "testuser", "fdsfsdf-sdfsdfsd-sdfsdfdsf-gwr4faf")
	// mockToken value set in api_test
	assert.Assert(t, mockToken == "adsa-kyuj-426tgv-sdhgb6t5rgf-dq3d2-sdfdsf")

}

func TestParseAuth(t *testing.T) {
	var mockRealm string = `Bearer realm="https://auth.docker.io/token"`
	var mockService string = `service="registry.docker.io"`
	testScopes := []struct {
		action string
	}{
		{`scope="repository:samalba/my-app:push`},
		{`scope="repository:samalba/my-app:push,pull`},
	}
	for _, mockScope := range testScopes {
		var mockDirectives authResponse = parseAuth(fmt.Sprintf("%s,%s,%s", mockRealm, mockService, mockScope.action))
		assert.Assert(t, mockDirectives.realm == strings.Split(strings.ReplaceAll(mockRealm, `"`, ""), "=")[1])
		assert.Assert(t, mockDirectives.service == strings.Split(strings.ReplaceAll(mockService, `"`, ""), "=")[1])
		assert.Assert(t, mockDirectives.service == strings.Split(strings.ReplaceAll(mockService, `"`, ""), "=")[1])
	}
}
