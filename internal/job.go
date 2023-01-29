package internal

func Job(isDryRun bool) error {

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
		err = MarkMoviesForDeletion(moviesdata, *ignoreTagId, isDryRun)
		if err != nil {
			return err
		}

		err = DeleteExpiredMovies(moviesdata, *ignoreTagId, isDryRun)
		if err != nil {
			return err
		}
	}

	// fmt.Println(*config)

	return nil
}
