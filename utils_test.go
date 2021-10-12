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

var ignoreTags = []string{"latest", "prod"}

func init() {
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
	testCase := []struct {
		useBasic, https                              bool
		registryUser, registryPassword, mockRegistry string
	}{
		{true, false, "exampleUser", "examplePass", "registry.corp.example.com"},
		{false, false, "exampleUser2", "anotherPass", "registry.localhost"},
		{true, true, "exampleUser", "examplePass", "registry.corp.example.com"},
		{false, true, "exampleUser2", "anotherPass", "registry.localhost"},
	}
	for _, scenario := range testCase {
		basicAuth = false
		basicAuth = scenario.useBasic
		testUrl := buildUrl(scenario.mockRegistry, scenario.https, scenario.registryUser, scenario.registryPassword)
		assert.Assert(t, strings.Contains(testUrl, scenario.mockRegistry))
		if scenario.useBasic {
			assert.Assert(t, strings.Contains(testUrl, scenario.registryUser))
			assert.Assert(t, strings.Contains(testUrl, "@"))
			assert.Assert(t, strings.Contains(testUrl, scenario.registryPassword))
		} else {
			assert.Assert(t, !strings.Contains(testUrl, scenario.registryUser))
			assert.Assert(t, !strings.Contains(testUrl, "@"))
			assert.Assert(t, !strings.Contains(testUrl, scenario.registryPassword))
		}
		if scenario.https {
			assert.Assert(t, strings.Contains(testUrl, "https://"))
		} else {
			assert.Assert(t, strings.Contains(testUrl, "http://"))
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
