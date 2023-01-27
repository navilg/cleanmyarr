package internal

import (
	"log"
)

func Job() error {
	log.Println("Starting process")
	if Config.Radarr.Enabled {
		// fmt.Println(config.Radarr.B64APIKey)
		ignoreTagId, err := GetTagIdFromRadarr(Config.IgnoreTag)
		if err != nil {
			return err
		}
		if ignoreTagId == nil {
			ignoreTagId, err = CreateTagInRadarr(Config.IgnoreTag)
			if err != nil {
				return err
			}
		}
		moviesdata, _ := GetMoviesData()
		err = MarkMoviesForDeletion(moviesdata, *ignoreTagId)
		if err != nil {
			return err
		}

		err = DeleteExpiredMovies(moviesdata, *ignoreTagId)
		if err != nil {
			return err
		}
	}

	// fmt.Println(*config)

	return nil
}
