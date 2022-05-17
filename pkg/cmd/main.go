package main

import (
	"log"

	"github.com/e-wrobel/microservice/configuration"
	"github.com/e-wrobel/microservice/streamer"
)

const (
	configName = "config"
	configType = "yaml"
	configPath = "/etc"
)

func main() {
	config, err := configuration.New(configName, configType, configPath)
	if err != nil {
		log.Fatal(err)
	}

	client := streamer.New(config.JSONFile, config.CreateJSONEndpoint)
	err = client.SendJSONFile()
	if err != nil {
		log.Fatalf("Unable to process json file: %v", err)
	}
}
