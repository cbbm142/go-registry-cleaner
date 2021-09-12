package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type authResponse struct {
	realm, service, scope string
}

type apiReqData struct {
	method, url, body, token, header string
}

func (call *apiReqData) apiDefaults() {
	if call.method == "" {
		call.method = "GET"
	}
	if call.header == "" {
		call.header = "Authentication"
	}
	if call.token == "" {
		call.token = "Bearer " + bearerToken
	}
}

func apiRequest(reqInfo apiReqData) *http.Response {

	// Don' load defaults if no token present
	if bearerToken != "" {
		reqInfo.apiDefaults()
	}
	client := &http.Client{}
	req, err := http.NewRequest(reqInfo.method, reqInfo.url, strings.NewReader(reqInfo.body))
	errCheck(err)
	// tokenBuild := "Bearer " + base64.StdEncoding.EncodeToString([]byte(registryToken))
	if reqInfo.header != "" {
		req.Header.Add(reqInfo.header, reqInfo.token)
	}
	resp, err := client.Do(req)
	errCheck(err)
	return resp
}

// func basicAuth() string {
// 	tokenBuilder := strings.Builder{}
// 	tokenBuilder.WriteString("Basic ")
// 	tokenBuilder.WriteString(base64.StdEncoding.EncodeToString([]byte(registryToken)))
// 	return tokenBuilder.String()
// }

// func makeAuth(auth authResponse) string {
// 	authRequestBuilder := strings.Builder{}
// 	authRequestBuilder.WriteString(auth.realm)
// 	authRequestBuilder.WriteString("?service=")
// 	authRequestBuilder.WriteString(auth.service)
// 	authRequestBuilder.WriteString("&scope=")
// 	authRequestBuilder.WriteString(auth.scope)
// 	return authRequestBuilder.String()
// }

func parseAuth(authHeader string) {
	var authDirectives authResponse = authResponse{}
	for _, directive := range strings.Split(authHeader, ",") {
		fmt.Println(directive)
		switch {
		case strings.Contains(directive, "Bearer realm"):
			authDirectives.realm = strings.Split(directive, "\"")[1]
		case strings.Contains(directive, "service"):
			authDirectives.service = strings.Split(directive, "\"")[1]
		case strings.Contains(directive, "scope"):
			authDirectives.scope = strings.Split(directive, "\"")[1]
		}
	}
	getToken(authDirectives)
}

func getTags(repo interface{}) *http.Response {
	// resp, err := http.Get(registryUrl + repo.(map[string]interface{})["name"].(string) + "/tags/list/")
	req := apiReqData{
		// url:    registryUrl + "/namespaces/" + repo.(map[string]interface{})["name"].(string) + "/tags/list/",
		url:    registryUrl + "namespaces/cbbm142/repositories/canvas/images-summary",
		method: "GET",
	}
	resp := apiRequest(req)
	// authresp := parseAuth(resp.Header["Www-Authenticate"])
	fmt.Println(registryUrl + "namespaces/cbbm142/repositories/images-summary")
	fmt.Println(registryUrl + repo.(map[string]interface{})["name"].(string) + "/tags/lists/")
	// errCheck(err)
	return resp
}

func getToken(auth authResponse) {

	var params = make(map[string]string)
	params["realm"] = auth.realm
	params["service"] = auth.service
	params["scope"] = auth.scope
	jsonParams, err := json.Marshal(&params)
	errCheck(err)
	req := apiReqData{
		url:    auth.realm,
		body:   string(jsonParams),
		header: "Authentication",
		token:  "cbbm142, 4dfd0fdf-d773-4549-bea4-29eb1bd45a49",
	}
	resp := apiRequest(req)
	token := make(map[string]string)
	body, err := ioutil.ReadAll(resp.Body)
	errCheck(err)
	json.Unmarshal([]byte(body), &token)
	bearerToken = token["token"]

}
