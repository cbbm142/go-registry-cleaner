package main

import (
	"strings"
	"testing"
	"time"
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
			t.Errorf("Tag %s returned true when it should be false", tag)
		}
		trueTag := "return_true_" + tag
		if !checkTag(trueTag) {
			t.Errorf("Tag %s returned false when it should be true", trueTag)
		}
	}
}

func TestCheckStale(t *testing.T) {
	keepDate := strings.ReplaceAll(time.Now().Add(time.Hour*24*-14).Format("2006/01/02"), "/", "-")

	if checkStale(keepDate, ignoreValues) {
		t.Errorf("Date %s returned true, when it should return false", keepDate)
	}
	staleDate := strings.ReplaceAll(time.Now().Add(time.Hour*24*-35).Format("2006/01/02"), "/", "-")
	if !checkStale(staleDate, ignoreValues) {
		t.Errorf("Date %s returned false, when it should return true", staleDate)
	}
}
