package internal

import (
	"log"
)

func Job() error {
	log.Println("Starting process")
	if Config.Radarr.Enabled {
		// fmt.Println(config.Radarr.B64APIKey)
		moviesdata, _ := GetMoviesData()
		err := MarkMoviesForDeletion(moviesdata)
		if err != nil {
			return err
		}
	}

	// fmt.Println(*config)

	return nil
}
