package internal

import (
	"fmt"
	"log"
)

func Job(statusFile string, isDryRun bool) error {

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
		moviesIgnored, err := GetMoviesIgnored(*ignoreTagId, moviesdata)
		if err != nil {
			return err
		}

		moviesMarkedForDeletion, err := MarkMoviesForDeletion(moviesdata, moviesIgnored, isDryRun)
		if err != nil {
			return err
		}

		moviesDeleted, err := DeleteExpiredMovies(moviesdata, moviesIgnored, isDryRun)
		if err != nil {
			return err
		}

		if !isDryRun {
			UpdateStatusFile(moviesDeleted, moviesIgnored, moviesMarkedForDeletion, statusFile)
		}

		if Config.Radarr.Notification && !isDryRun {

			log.Println("Sending notification")

			subject := "ALERT: [Cleanmyarr] [RADARR] Movies deleted"
			body := `

Movies deleted --> ` + fmt.Sprint(moviesDeleted) + `

Movies Marked for deletion --> ` + fmt.Sprint(moviesMarkedForDeletion) + `

Movies marked for deletion will be deleted in next maintenance schedule.
To protect them from automatically deleting on next maintenance, Add tag "` + Config.IgnoreTag + `" to them in Radarr.
			
Next Maintenance schedule --> ` + State.NextMaintenanceDate

			SendEmailNotification(subject, body)
			SendGotifyNotification(subject, body)

		}
	}

	// fmt.Println(*config)

	return nil
}
