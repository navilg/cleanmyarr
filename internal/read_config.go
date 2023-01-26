package internal

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

func ReadConfig(configFile string) (*Config, error) {
	log.Println("Reading configurations")
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println("Failed to read configuration file.", err.Error())
		return nil, err
	}

	var config Config

	// yaml.Unmarshal(data, &config)
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Println("Failed to read configuration file.", err.Error())
		return nil, err
	}

	return &config, nil
}
