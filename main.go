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

var ignoreValues ignores = ignores{}
var configFile, registryPassword, registryUrl, registryUser string = "", "", "", ""

func main() {
	configFile = "config.yml"
	config := readConfig(configFile)
	ignoreValues = ignores{
		tags: config["defaultTags"].([]interface{}),
		days: config["defaultDays"].(int),
	}
	errCheck(godotenv.Load())
	registryUser = os.Getenv("username")
	registryPassword = os.Getenv("password")
	registryUrl = buildUrl(config["host"])
	dryRun := config["dryRun"].(bool)
	for _, repo := range config["repos"].([]interface{}) {
		var name string = repo.(map[string]interface{})["name"].(string)
		combineTags(repo, &ignoreValues)
		setDays(config["defaultDays"], repo, &ignoreValues)
		registryTags := getTags(name)
		for _, tag := range registryTags {
			if checkTag(tag) {
				// Creation time is associated with manifests
				_, creationTime := getManifest(name, tag)
				if checkStale(creationTime, ignoreValues) {
					// Deletion is done via digest not tag/manifest
					digest := getImageDigest(name, tag)
					err := deleteDigest(name, digest, tag, dryRun)
					errCheck(err)
				}
			}
		}
	}
}
