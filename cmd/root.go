/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func driver() {
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

	var parsedLastMaintenanceRun time.Time

	if internal.State.LastMaintenanceRun != "" {
		parsedLastMaintenanceRun, err = time.Parse("2006-01-02 15:04:05 +0000 UTC", internal.State.LastMaintenanceRun)
		if err != nil {
			log.Println("Failed get last maintenenace run time", err.Error())
		}
	}
	maintenanceCycleDays := internal.MaintenanceCycleInInt(internal.Config.MaintenanceCycle)
	nextMaintenanceCycle := parsedLastMaintenanceRun.Add(time.Duration(maintenanceCycleDays) * time.Hour * 24)

	if time.Now().After(nextMaintenanceCycle) {
		err = internal.Job(isDryRun)
		if err == nil && !isDryRun {
			internal.UpdateStatusFile(statusFile)
		}
	}

	jobSyncInterval := internal.JobSyncInterval * time.Hour // Job syncs with config in every 1 hours
	// jobSyncInterval := 5 * time.Second

	ticker := time.NewTicker(jobSyncInterval)

	for range ticker.C {
		retryCount := 0
		for {
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

		if internal.State.LastMaintenanceRun != "" {
			parsedLastMaintenanceRun, err = time.Parse("2006-01-02 15:04:05 +0000 UTC", internal.State.LastMaintenanceRun)
			if err != nil {
				log.Println("Failed get last maintenenace run time", err.Error())
			}
		}

		maintenanceCycleDays := internal.MaintenanceCycleInInt(internal.Config.MaintenanceCycle)
		nextMaintenanceCycle := parsedLastMaintenanceRun.Add(time.Duration(maintenanceCycleDays) * time.Hour * 24)

		if time.Now().After(nextMaintenanceCycle) {
			err = internal.Job(isDryRun)
			if err == nil && !isDryRun {
				internal.UpdateStatusFile(statusFile)
			}
		}
	}
}
