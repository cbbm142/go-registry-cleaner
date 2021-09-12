package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type apiReqData struct {
	method, url, body, headerDirective, header string
}

func (call *apiReqData) apiDefaults() {
	if call.method == "" {
		call.method = "GET"
	}
	call.url = registryUrl + call.url
}

func apiRequest(reqInfo apiReqData) *http.Response {
	reqInfo.apiDefaults()
	client := &http.Client{}
	req, err := http.NewRequest(reqInfo.method, reqInfo.url, strings.NewReader(reqInfo.body))
	errCheck(err)
	if reqInfo.header != "" {
		req.Header.Add(reqInfo.header, reqInfo.headerDirective)
	}
	resp, err := client.Do(req)
	errCheck(err)
	return resp
}

func deleteDigest(name string, digest string, tag string, dryRun bool) {
	fmt.Printf("Deleting digest %s for tag %s for repo %s.", digest, tag, name)
	if !dryRun {
		req := apiReqData{
			url:    name + "/manifests/" + digest,
			method: "DELETE",
		}
		resp := apiRequest(req)
		if resp.StatusCode != 202 {
			log.Fatalf("There was an error while attempting to delete %s", digest)
		}
	}
}

func getImageDigest(name string, tag string) string {
	// Same endpoint as manifests, but adding header make registry reply with digest
	req := apiReqData{
		url:             name + "/manifests/" + tag,
		header:          "Accept",
		headerDirective: "application/vnd.docker.distribution.manifest.v2+json",
	}
	resp := apiRequest(req)
	return resp.Header["Docker-Content-Digest"][0]
}

func getTags(name string) []string {
	req := apiReqData{
		url: name + "/tags/list/",
	}
	resp := apiRequest(req)
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
	resp := apiRequest(req)
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
