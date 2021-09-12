package main

import (
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v3"
)

func buildUrl(host interface{}, token string) string {
	builder := strings.Builder{}
	builder.WriteString("https://")
	// builder.WriteString("cbbm142:")
	// builder.WriteString(registryToken)
	// builder.WriteString("@")
	builder.WriteString(host.(string))
	builder.WriteString("/v2/")
	return builder.String()
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

func setDays(defaultDays interface{}, repo interface{}, ignoreValues *ignores) {
	if defaultDays.(int) != repo.(map[string]interface{})["days"].(int) {
		ignoreValues.days = repo.(map[string]interface{})["days"].(int)
	}
}
