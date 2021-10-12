package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type authResponse struct {
	realm, service, scope string
}

func authDetect(appConfig config) error {
	var protocol string
	if appConfig.httpsProtocol {
		protocol = "https://"
	} else {
		protocol = "http://"
	}
	resp, err := http.Get(fmt.Sprintf("%s%s", protocol, appConfig.configMap["host"].(string)))
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case 200:
		// No auth needed
		bearerToken = ""
		basicAuth = false
	case 401:
		if strings.Contains(resp.Header["Www-Authenticate"][0], "Bearer realm") {
			var parsedResponse authResponse = parseAuth(resp.Header["Www-Authenticate"][0])
			bearerToken = getToken(parsedResponse, appConfig.registryUser, appConfig.registryToken)
		} else if strings.Contains(resp.Header["Www-Authenticate"][0], "Basic") {
			bearerToken = ""
			basicAuth = true
		} else {
			return fmt.Errorf("unable to determine auth type")
		}
	default:
		return fmt.Errorf("it appears the registry either doesn't exist, or doesn't support the v2 api")
	}
	return nil
}

func getToken(auth authResponse, user string, token string) string {

	var params = make(map[string]string)
	params["realm"] = auth.realm
	params["service"] = auth.service
	params["scope"] = auth.scope
	jsonParams, err := json.Marshal(&params)
	errCheck(err)
	req := apiReqData{
		url:             auth.realm,
		body:            string(jsonParams),
		header:          []string{"Authentication"},
		headerDirective: []string{fmt.Sprintf("%s, %s", user, token)},
	}
	resp := req.apiRequest()
	body := decodeBody(resp)
	return body.(map[string]interface{})["token"].(string)
}

func parseAuth(authHeader string) authResponse {
	// Only used for oauth
	var authDirectives authResponse = authResponse{}
	for _, directive := range strings.Split(authHeader, ",") {
		if authDirectives.scope != "" {
			authDirectives.scope = authDirectives.scope + "," + strings.Split(directive, "\"")[0]
		}
		switch {
		case strings.Contains(directive, "Bearer realm"):
			authDirectives.realm = strings.Split(directive, "\"")[1]
		case strings.Contains(directive, "service"):
			authDirectives.service = strings.Split(directive, "\"")[1]
		case strings.Contains(directive, "scope"):
			authDirectives.scope = strings.Split(directive, "\"")[1]
		}
	}
	return authDirectives
}
