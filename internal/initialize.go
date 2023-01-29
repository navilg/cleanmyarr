package internal

import (
	"log"
	"os"
)

var sample_config_file_url string = "https://raw.githubusercontent.com/navilg/cleanmyarr/main/sample-config.yaml"

func Initialize(configFile string) error {

	log.Println("Initializing...")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Println("Configuration file not found. Downloading sample configuration.")
		err := DownloadFile(sample_config_file_url, configFile)
		if err != nil {
			log.Println("Failed to initialize.")
			return err
		}
	}

	log.Println("Initialization completed.")
	return nil
}
