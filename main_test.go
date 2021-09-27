package main

import (
	"os"
	"reflect"
	"testing"

	"gotest.tools/assert"
)

func TestLoadConfig(t *testing.T) {
	configFile := "config.yml.example"
	os.Setenv("username", "testUser")
	os.Setenv("password", "testPass")
	mockConfig := loadConfig(configFile)
	registryValue := os.Getenv("registryUrl")
	assert.Assert(t, mockConfig.registryUser == "testUser")
	assert.Assert(t, registryValue == "https://testUser:testPass@registry.example.com/v2/")
	assert.Assert(t, mockConfig.ignoreValues.days == 30)
	assert.Assert(t, reflect.ValueOf(mockConfig.configMap).Len() == 5)
	assert.Assert(t, reflect.ValueOf(mockConfig.ignoreValues.tags).Len() == 2)
	assert.Assert(t, mockConfig.dryRun == true)
}
