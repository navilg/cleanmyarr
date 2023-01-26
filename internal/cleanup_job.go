package internal

import (
	"fmt"
	"log"
)

func CleanupJob(config *Config) error {
	log.Println("Starting process")
	if config.Radarr.Enabled {
		fmt.Println(config.Radarr.B64APIKey)
		moviesdata, _ := GetMoviesData(config.Radarr.URL, config.Radarr.B64APIKey)
		fmt.Println(string(moviesdata))
	}

	// fmt.Println(*config)

	return nil
}
