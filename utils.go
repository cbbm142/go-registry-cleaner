package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func decodeBody(resp *http.Response) interface{} {
	var body interface{}
	jsonBody, err := ioutil.ReadAll(resp.Body)
	errCheck(err)
	err = json.Unmarshal([]byte(jsonBody), &body)
	errCheck(err)
	return body
}

func buildUrl(host interface{}) string {
	builder := strings.Builder{}
	builder.WriteString("https://")
	builder.WriteString(registryUser)
	builder.WriteString(":")
	builder.WriteString(registryPassword)
	builder.WriteString("@")
	builder.WriteString(host.(string))
	builder.WriteString("/v2/")
	return builder.String()
}

func checkTag(tag string) bool {
	var checkTag bool = true
	for _, skipTag := range ignoreValues.tags {
		if tag == skipTag.(string) {
			checkTag = false
		}
	}
	return checkTag
}

func combineTags(repo interface{}, ignoreValues *ignores) {
	for _, tag := range repo.(map[string]interface{})["tags"].([]interface{}) {
		var addTag bool = true
		for _, v := range ignoreValues.tags {
			if tag == v {
				addTag = false
			}
		}
		if addTag {
			ignoreValues.tags = append(ignoreValues.tags, tag)
		}
	}
}

func errCheck(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func readConfig(configFile string) map[interface{}]interface{} {
	yamlFile, err := ioutil.ReadFile(configFile)
	errCheck(err)
	yamlData := make(map[interface{}]interface{})
	errCheck(yaml.Unmarshal(yamlFile, &yamlData))
	return yamlData
}

func checkStale(date string, ignoreValues ignores) bool {
	var splitDate []int
	for _, v := range strings.Split(date, "-") {
		int, _ := strconv.Atoi(v)
		splitDate = append(splitDate, int)
	}
	timeDate := time.Date(splitDate[0], time.Month(splitDate[1]), splitDate[2], 0, 0, 0, 0, time.UTC)
	age := time.Since(timeDate).Hours() / 24
	return age > float64(ignoreValues.days)
}

func setDays(defaultDays interface{}, repo interface{}, ignoreValues *ignores) {
	if defaultDays.(int) != repo.(map[string]interface{})["days"].(int) {
		ignoreValues.days = repo.(map[string]interface{})["days"].(int)
	}
}
