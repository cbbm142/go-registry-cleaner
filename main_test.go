package main

import (
	"strings"
	"testing"
)

var config map[interface{}]interface{}
var ignoreTags = []string{"latest", "prod"}

func init() {
	configFile = "config.yml.example"
	config = readConfig(configFile)
	ignoreValues = ignores{
		days: 30,
	}
	for _, val := range ignoreTags {
		ignoreValues.tags = append(ignoreValues.tags, val)
	}
}

func TestBuildUrl(t *testing.T) {
	registryUser = "testUser"
	registryPassword = "testPass"
	testUrl := buildUrl(config["host"])
	var testItems = []string{registryUser, registryPassword, "registry.example.com"}
	for _, item := range testItems {
		if !strings.Contains(testUrl, item) {
			t.Errorf("Missing %s from url %s", item, testUrl)
		}
	}
}

func TestCheckTag(t *testing.T) {
	for _, tag := range ignoreTags {
		if checkTag(tag) {
			t.Errorf("Tag %s returned wrong value", tag)
		}
	}
}
