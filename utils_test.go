package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"gotest.tools/assert"
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

func TestDecodeBody(t *testing.T) {
	var buff bytes.Buffer
	var testMessage string = "Test body message"
	jsonMessage, _ := json.Marshal(testMessage)
	buff.WriteString(string(jsonMessage))
	recorder := httptest.NewRecorder()
	recorder.Body = &buff
	simResp := recorder.Result()
	decodedBody := decodeBody(simResp)
	assert.Assert(t, decodedBody.(string) == testMessage)

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

func TestCombineTags(t *testing.T) {
	var fakeIgnores ignores
	fakeRepo := map[string]interface{}{
		"tags": []interface{}{"latest", "tag1"},
	}
	fakeIgnores.tags = []interface{}{"latest", "prod"}
	combineTags(fakeRepo, &fakeIgnores)

	assert.Assert(t, reflect.ValueOf(fakeIgnores.tags).Len() == 3)
	for _, tag := range fakeIgnores.tags {
		var countOccurances int
		for _, tagsToCheck := range fakeIgnores.tags {
			if tag == tagsToCheck {
				countOccurances++
			}
		}
		assert.Assert(t, countOccurances == 1)
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

func TestSetDays(t *testing.T) {
	var defaultDays interface{}
	var mockIgnores ignores
	defaultDays = 30

	var daysToCheck = []int{100, 35, 365, 1}
	for _, day := range daysToCheck {
		fakeRepo := map[string]interface{}{
			"days": day,
		}
		setDays(defaultDays, fakeRepo, &mockIgnores)
		assert.Assert(t, mockIgnores.days == day)
	}
}
