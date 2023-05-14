package saver

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type SaverConfig struct {
	SubscribedInstruments []string
	SaverPort             uint
}

func ParseConfig(configPath string) SaverConfig {
	jsonFile, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}

	byteConfig, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()

	var config SaverConfig

	var rawConfig map[string]interface{}
	json.Unmarshal(byteConfig, &rawConfig)

	for _, elem := range rawConfig["subscribed_instruments"].([]interface{}) {
		instrument := elem.(string)
		config.SubscribedInstruments = append(config.SubscribedInstruments, instrument)
	}
	config.SaverPort = uint(rawConfig["port"].(float64))

	return config
}
