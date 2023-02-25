package internal

import (
	"io/ioutil"
	"log"
	"os"
)

// var sample_config_file_url string = "https://raw.githubusercontent.com/navilg/cleanmyarr/main/sample-config.yaml"
var sampleConfigFile string = "sample-config.yaml"

func Initialize(configFile string) error {

	log.Println("Initializing...")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Println("Configuration file not found.")
		log.Println("Creating default configuration file.")
		// err := DownloadFile(sample_config_file_url, configFile)
		input, err := ioutil.ReadFile(sampleConfigFile)
		if err != nil {
			log.Println("Failed to initialize.")
			return err
		}

		err = ioutil.WriteFile(configFile, input, 0644)
		if err != nil {
			log.Println("Failed to initialize.")
			return err
		}
	}

	log.Println("Initialization completed.")
	return nil
}
