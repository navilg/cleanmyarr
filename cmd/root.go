/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/navilg/cleanmyarr/internal"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cleanmyarr",
	Short: "A lightweight utility to delete movies and shows from Radarr and Sonarr after specified time.",
	Long:  `A lightweight utility to delete movies and shows from Radarr and Sonarr after specified time`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		driver()
	},
}
var cfgFile string
var isDryRun bool

// var debug bool

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "/config/config.yaml", "config file (default is /config/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&isDryRun, "dry-run", false, "Dry run (default is false")
	// rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Generate more logs for debuging (default is false)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func driver() {
	internal.Now = time.Now()
	err := internal.Initialize(cfgFile)
	if err != nil {
		os.Exit(1)
	}

	_, err = internal.ReadConfig(cfgFile)
	if err != nil {
		os.Exit(1)
	}

	cfgFileDir := filepath.Dir(cfgFile)
	statusFile := cfgFileDir + "/" + internal.StatusFileName

	_, err = internal.ReadStatus(statusFile)

	if !isDryRun {
		log.Println("Process running")
	} else {
		log.Println("Process running [DRY RUN]")
	}

	nextMaintenanceCycle, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", internal.State.NextMaintenanceDate)
	if err != nil {
		log.Println("Failed to get next maintenanance cycle", err.Error())
	}

	if internal.Now.After(nextMaintenanceCycle) {
		err = internal.Job(statusFile, isDryRun)
	} else {
		log.Println("Next maintenance cycle is at", nextMaintenanceCycle)
		log.Println("Napping...")
	}

	jobSyncInterval := internal.JobSyncInterval * time.Hour // Job syncs with config in every 1 hours
	// jobSyncInterval := 5 * time.Second // For test

	ticker := time.NewTicker(jobSyncInterval)

	for range ticker.C {
		retryCount := 0
		for {
			internal.Now = time.Now()
			_, err := internal.ReadConfig(cfgFile)

			if err != nil && retryCount < 10 {
				log.Println("Failed to read config file. Retrying in 1 min.")
				time.Sleep(time.Minute)
				retryCount = retryCount + 1
				continue
			} else if retryCount == 10 {
				log.Println("Failed to read config file.")
			}
			break
		}

		_, err = internal.ReadStatus(statusFile)

		nextMaintenanceCycle, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", internal.State.NextMaintenanceDate)
		if err != nil {
			log.Println("Failed get net maintenanance cycle", err.Error())
		}

		if internal.Now.After(nextMaintenanceCycle) {
			err = internal.Job(statusFile, isDryRun)
		} else {
			retryCount = 0
			for {
				if retryCount >= 3 {
					log.Println("Failed to mark movies for deletion.")
					break
				}

				ignoreTagId, err := internal.GetTagIdFromRadarr(internal.Config.IgnoreTag)
				if err != nil {
					log.Println("Failed to get ignore tag id from radarr. Retrying in 1 min.")
					time.Sleep(time.Minute)
					retryCount = retryCount + 1
					continue
				}
				if ignoreTagId == nil {
					ignoreTagId, err = internal.CreateTagInRadarr(internal.Config.IgnoreTag)
					if err != nil {
						log.Println("Failed to create ignore tag in radarr. Retrying in 1 min.")
						time.Sleep(time.Minute)
						retryCount = retryCount + 1
						continue
					}
				}

				moviesdata, _ := internal.GetMoviesData()
				moviesIgnored, err := internal.GetMoviesIgnored(*ignoreTagId, moviesdata)
				if err != nil {
					log.Println("Failed to get movies data from radarr. Retrying in 1 min.")
					time.Sleep(time.Minute)
					retryCount = retryCount + 1
					continue
				}
				log.Println("Movies ignored", moviesIgnored)
				moviesMarkedForDeletion, err := internal.MarkMoviesForDeletion(moviesdata, moviesIgnored, nextMaintenanceCycle, isDryRun)

				if err != nil {
					log.Println("Failed to mark movies for deletion. Retrying in 1 min.")
					time.Sleep(time.Minute)
					retryCount = retryCount + 1
					continue
				}

				var newMoviesMarkedForDeletion []string

				for _, movie := range moviesMarkedForDeletion {
					var isAlreadyMarked bool = false
					for _, moviesAlreadyMarked := range internal.State.MoviesMarkedForDeletion {
						if movie == moviesAlreadyMarked {
							isAlreadyMarked = true
							break
						}
					}

					if !isAlreadyMarked {
						newMoviesMarkedForDeletion = append(newMoviesMarkedForDeletion, movie)
					}
				}

				if !isDryRun {
					internal.UpdateStatusFile(internal.State.LastMaintenanceDate, internal.State.DeletedMovies, internal.State.IgnoredMovies, moviesMarkedForDeletion, statusFile)
					_, err = internal.ReadStatus(statusFile)
				}

				log.Println("Next maintenance time:", internal.State.NextMaintenanceDate)

				if internal.Config.Radarr.Notification && len(newMoviesMarkedForDeletion) > 0 && !isDryRun {

					log.Println("Sending notification")

					subject := "ALERT: [Cleanmyarr] [RADARR] New movies marked for deletion"
					body := `

[CLEANMYARR]

New movies Marked for deletion --> ` + fmt.Sprint(newMoviesMarkedForDeletion) + `
Movies marked for deletion --> ` + fmt.Sprint(moviesMarkedForDeletion) + `
Next maintenance time --> ` + internal.State.NextMaintenanceDate + `
	
Movies marked for deletion will be deleted in next maintenance schedule.
To protect them from automatically deleting on next maintenance, Add tag "` + internal.Config.IgnoreTag + `" to them in Radarr.
				
Next Maintenance schedule --> ` + internal.State.NextMaintenanceDate

					internal.SendEmailNotification(subject, body)
					internal.SendGotifyNotification(subject, body)
					internal.SendTelegramNotification(body)

				}

				break
			}
		}

		log.Println("Napping...")
	}
}
