package main

// Removes stale images from a self hosted registry

import (
	"os"

	"github.com/joho/godotenv"
)

type ignores struct {
	tags []interface{}
	days int
}
type config struct {
	configMap        map[interface{}]interface{}
	ignoreValues     ignores
	registryUser     string
	registryPassword string
	registryToken    string
	dryRun           bool
}

var ignoreValues ignores = ignores{}
var bearerToken string
var basicAuth bool = false

func loadConfig(configFile string) config {
	var appConfig config
	appConfig.configMap = readConfig(configFile)
	appConfig.ignoreValues = ignores{
		tags: appConfig.configMap["defaultTags"].([]interface{}),
		days: appConfig.configMap["defaultDays"].(int),
	}
	_ = godotenv.Load()
	appConfig.registryUser = os.Getenv("username")
	appConfig.registryPassword = os.Getenv("password")
	appConfig.registryToken = os.Getenv("token")
	appConfig.dryRun = appConfig.configMap["dryRun"].(bool)
	os.Setenv("registryUrl", buildUrl(appConfig.configMap["host"], appConfig.registryUser, appConfig.registryPassword))
	return appConfig
}

func main() {
	var configFile = "config.yml"
	var appConfig config = loadConfig(configFile)
	err := authDetect(appConfig)
	errCheck(err)
	for _, repo := range appConfig.configMap["repos"].([]interface{}) {
		var name string = repo.(map[string]interface{})["name"].(string)
		if _, tagsExist := repo.(map[string]interface{})["tags"]; tagsExist {
			combineTags(repo, &appConfig.ignoreValues)
		}
		if _, daysExist := repo.(map[string]interface{})["days"]; daysExist {
			setDays(appConfig.configMap["defaultDays"], repo, &appConfig.ignoreValues)
		}
		registryTags := getTags(name)
		for _, tag := range registryTags {
			if checkTag(tag) {
				// Creation time is associated with manifests
				_, creationTime := getManifest(name, tag)
				if checkStale(creationTime, ignoreValues) {
					// Deletion is done via digest not tag/manifest
					digest := getImageDigest(name, tag)
					err := deleteDigest(name, digest, tag, appConfig.dryRun)
					errCheck(err)
				}
			}
		}
	}
}
