package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type apiReqData struct {
	method, url, body, headerDirective, header string
}

func (setDefault *apiReqData) apiDefaults() {
	if setDefault.method == "" {
		setDefault.method = "GET"
	}
	setDefault.url = registryUrl + setDefault.url
}

func (apiCall apiReqData) apiRequest() *http.Response {
	apiCall.apiDefaults()
	client := &http.Client{}
	req, err := http.NewRequest(apiCall.method, apiCall.url, strings.NewReader(apiCall.body))
	errCheck(err)
	if apiCall.header != "" {
		req.Header.Add(apiCall.header, apiCall.headerDirective)
	}
	resp, err := client.Do(req)
	errCheck(err)
	return resp
}

func deleteDigest(name string, digest string, tag string, dryRun bool) error {
	fmt.Printf("Deleting digest %s for tag %s for repo %s.", digest, tag, name)
	if !dryRun {
		req := apiReqData{
			url:    name + "/manifests/" + digest,
			method: "DELETE",
		}
		resp := req.apiRequest()
		if resp.StatusCode != 202 {
			return errors.New(fmt.Sprintf("There was an error while attempting to delete %s", digest))
		}
	}
	return nil
}

func getImageDigest(name string, tag string) string {
	// Same endpoint as manifests, but adding header make registry reply with digest
	req := apiReqData{
		url:             name + "/manifests/" + tag,
		header:          "Accept",
		headerDirective: "application/vnd.docker.distribution.manifest.v2+json",
	}
	resp := req.apiRequest()
	return resp.Header["Docker-Content-Digest"][0]
}

func getTags(name string) []string {
	req := apiReqData{
		url: name + "/tags/list/",
	}
	resp := req.apiRequest()
	body := decodeBody(resp)
	intTags := body.(map[string]interface{})["tags"].([]interface{})
	var tags []string = nil
	for _, tag := range intTags {
		tags = append(tags, tag.(string))
	}
	return tags
}

func getManifest(name string, tag string) (string, string) {
	req := apiReqData{
		url: name + "/manifests/" + tag,
	}
	resp := req.apiRequest()
	body := decodeBody(resp)
	// Only care about newest layer
	manifest := body.(map[string]interface{})["fsLayers"].([]interface{})[0].(map[string]interface{})["blobSum"].(string)
	v1compat := body.(map[string]interface{})["history"].([]interface{})[0].(map[string]interface{})["v1Compatibility"].(string)
	var creationTime string = ""
	for _, v := range strings.Split(v1compat, ",") {
		if strings.Contains(v, "created") {
			creationTime = strings.Split(strings.Split(v, "\"")[3], "T")[0]
		}
	}
	return manifest, creationTime

}
