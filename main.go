package main

import (
	"io/ioutil"
	"log"
	"moonspace/api"

	"gopkg.in/yaml.v3"
)

const configFileLocation = "config.yaml"

func main() {
	log.Default().Println("Starting the service...")

	data, err := ioutil.ReadFile(configFileLocation)
	if err != nil {
		log.Fatalf("Error creating API: Error opening config file: %s", err)
	}

	cfg := api.ServerConfig{}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		panic(err)
	}

	api := api.New(cfg.Server)

	log.Default().Println("API Created with config: ", cfg.String())
	api.Start()
}
