package internal

import (
	"fmt"
	"log"
	"time"
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

		nextMaintenanceCycle, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", State.NextMaintenanceDate)
		if err != nil {
			log.Println("Failed get net maintenanance cycle", err.Error())
		}

		moviesdata, _ := GetMoviesData()
		moviesIgnored, err := GetMoviesIgnored(*ignoreTagId, moviesdata)
		if err != nil {
			return err
		}

		log.Println("Movies ignored", moviesIgnored)
		moviesMarkedForDeletion, err := MarkMoviesForDeletion(moviesdata, moviesIgnored, nextMaintenanceCycle, isDryRun)
		if err != nil {
			return err
		}

		moviesDeleted, err := DeleteExpiredMovies(moviesdata, moviesIgnored, nextMaintenanceCycle, isDryRun)
		if err != nil {
			return err
		}

		if !isDryRun {
			UpdateStatusFile(Now.UTC().String(), moviesDeleted, moviesIgnored, moviesMarkedForDeletion, statusFile)
			_, err = ReadStatus(statusFile)
		}

		log.Println("Next maintenance time:", State.NextMaintenanceDate)

		if Config.Radarr.Notification && !isDryRun {

			log.Println("Sending notification")

			subject := "ALERT: [Cleanmyarr] [RADARR] Movies deleted"
			body := `

[CLEANMYARR]

Movies deleted --> ` + fmt.Sprint(moviesDeleted) + `
Movies Marked for deletion --> ` + fmt.Sprint(moviesMarkedForDeletion) + `
Next maintenance time --> ` + fmt.Sprint(State.NextMaintenanceDate) + `

Movies marked for deletion will be deleted in next maintenance schedule.
To protect them from automatically deleting on next maintenance, Add tag "` + Config.IgnoreTag + `" to them in Radarr.
			
Next Maintenance schedule --> ` + State.NextMaintenanceDate

			SendEmailNotification(subject, body)
			SendGotifyNotification(subject, body)
			SendTelegramNotification(body)

		}
	}

	// fmt.Println(*config)

	return nil
}
