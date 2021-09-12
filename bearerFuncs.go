// Unused bearer functions
package main

// type authResponse struct {
// 	realm, service, scope string
// }

// func getToken(auth authResponse) {

// 	var params = make(map[string]string)
// 	params["realm"] = auth.realm
// 	params["service"] = auth.service
// 	params["scope"] = auth.scope
// 	jsonParams, err := json.Marshal(&params)
// 	errCheck(err)
// 	req := apiReqData{
// 		url:    auth.realm,
// 		body:   string(jsonParams),
// 		header: "Authentication",
// 		token:  "cbbm142, 4dfd0fdf-d773-4549-bea4-29eb1bd45a49",
// 	}
// 	resp := apiRequest(req)
// 	body := decodeBody(resp)
// 	bearerToken = body["token"]
// }

// func parseAuth(authHeader string) {
// 	var authDirectives authResponse = authResponse{}
// 	for _, directive := range strings.Split(authHeader, ",") {
// 		fmt.Println(directive)
// 		switch {
// 		case strings.Contains(directive, "Bearer realm"):
// 			authDirectives.realm = strings.Split(directive, "\"")[1]
// 		case strings.Contains(directive, "service"):
// 			authDirectives.service = strings.Split(directive, "\"")[1]
// 		case strings.Contains(directive, "scope"):
// 			authDirectives.scope = strings.Split(directive, "\"")[1]
// 		}
// 	}
// 	getToken(authDirectives)
// }
