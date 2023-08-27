package internal

import (
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func InitializeConfig(configFile string) error {

	var err error = nil

	log.Println("Initializing...")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Println("Configuration file not found.")
		log.Println("Creating default configuration file.")
		f, err := os.Create(configFile)
		f.Close()

		if err != nil {
			log.Println("Failed to initialize.")
			return err
		}
	}

	maintenanceCycle := os.Getenv("CMA_MAINTENANCE_CYCLE")
	if maintenanceCycle != "" {
		Config.MaintenanceCycle = Interval(maintenanceCycle)
	}

	deleteAfterDays := os.Getenv("CMA_DELETE_AFTER_DAYS")
	if deleteAfterDays != "" {
		Config.DeleteAfterDays, err = strconv.Atoi(deleteAfterDays)
		if err != nil {
			return err
		}
	}

	ignoreTag := os.Getenv("CMA_IGNORE_TAG")
	if ignoreTag != "" {
		Config.IgnoreTag = ignoreTag
	}

	smtpEnabled := os.Getenv("CMA_ENABLE_EMAIL_NOTIFICATION")
	if smtpEnabled != "" {
		Config.NotificationChannel.SMTP.Enabled, err = strconv.ParseBool(smtpEnabled)
		if err != nil {
			return err
		}
	}

	smtpServer := os.Getenv("CMA_SMTP_SERVER")
	if smtpServer != "" {
		Config.NotificationChannel.SMTP.Server = smtpServer
	}

	smtpPort := os.Getenv("CMA_SMTP_PORT")
	if smtpPort != "" {
		Config.NotificationChannel.SMTP.Port, err = strconv.Atoi(smtpPort)
		if err != nil {
			return err
		}
	}

	smtpSecurity := os.Getenv("CMA_SMTP_SECURITY")
	if smtpSecurity != "" {
		Config.NotificationChannel.SMTP.Security = Security(smtpSecurity)
	}

	smtpUsername := os.Getenv("CMA_SMTP_USERNAME")
	if smtpUsername != "" {
		Config.NotificationChannel.SMTP.Username = smtpUsername
	}

	smtpb64Pass := os.Getenv("CMA_SMTP_ENCODED_PASSWORD")
	if smtpb64Pass != "" {
		Config.NotificationChannel.SMTP.B64Password = smtpb64Pass
	}

	smtpFromEmail := os.Getenv("CMA_SMTP_FROM_EMAIL")
	if smtpFromEmail != "" {
		Config.NotificationChannel.SMTP.FromEmail = smtpFromEmail
	}

	smtpToEmails := os.Getenv("CMA_SMTP_TO_EMAILS")
	if smtpToEmails != "" {
		Config.NotificationChannel.SMTP.ToEmail = strings.Split(smtpToEmails, ",")
	}

	ccEmails := os.Getenv("CMA_SMTP_CC_EMAILS")
	if ccEmails != "" {
		Config.NotificationChannel.SMTP.CcEmail = strings.Split(ccEmails, ",")
	}

	bccEmails := os.Getenv("CMA_SMTP_BCC_EMAILS")
	if bccEmails != "" {
		Config.NotificationChannel.SMTP.BccEmail = strings.Split(bccEmails, ",")
	}

	argConfigData, err := yaml.Marshal(Config)
	if err != nil {
		log.Println("Failed to initialize configuration file", err.Error())
		return err
	}

	err = os.WriteFile(configFile, argConfigData, 0644)
	if err != nil {
		log.Println("Failed to initialize configuration file", err.Error())
		return err
	}

	log.Println("Initialization completed.")
	return nil
}
