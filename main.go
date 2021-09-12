package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ignores struct {
	tags []interface{}
	days int
}

var ignoreValues ignores = ignores{}
var bearerToken, registryToken, registryUrl string = "", "", ""

func main() {
	const configFile string = "config.yml"
	config := readConfig(configFile)
	ignoreValues = ignores{
		tags: config["defaultTags"].([]interface{}),
		days: config["defaultDays"].(int),
	}
	errCheck(godotenv.Load())
	registryToken = os.Getenv("token")
	registryUrl = buildUrl(config["host"], registryToken)
	fmt.Println(registryUrl)
	for _, repo := range config["repos"].([]interface{}) {
		log.Println(repo)
		fmt.Println(ignoreValues.tags)
		combineTags(repo, &ignoreValues)
		setDays(config["defaultDays"], repo, &ignoreValues)
		if bearerToken == "" {
			resp := getTags(repo)
			parseAuth(resp.Header["Www-Authenticate"][0])
		}
		fmt.Println(bearerToken)
		foundTags := getTags(repo)
		fmt.Println(foundTags)
	}
}
