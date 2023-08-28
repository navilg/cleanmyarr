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

	// Get all configurations from environment variable

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

	gotifyEnabled := os.Getenv("CMA_ENABLE_GOTIFY_NOTIFICATION")
	if gotifyEnabled != "" {
		Config.NotificationChannel.Gotify.Enabled, err = strconv.ParseBool(gotifyEnabled)
		if err != nil {
			return err
		}
	}

	gotifyURL := os.Getenv("CMA_GOTIFY_URL")
	if gotifyURL != "" {
		Config.NotificationChannel.Gotify.URL = gotifyURL
	}

	gotifyb64AppToken := os.Getenv("CMA_GOTIFY_ENCODED_APP_TOKEN")
	if gotifyb64AppToken != "" {
		Config.NotificationChannel.Gotify.B64AppToken = gotifyb64AppToken
	}

	gotifyPriority := os.Getenv("CMA_GOTIFY_PRIORITY")
	if gotifyPriority != "" {
		Config.NotificationChannel.Gotify.Priority, err = strconv.Atoi(gotifyPriority)
		if err != nil {
			return err
		}
	}

	telegramEnabled := os.Getenv("CMA_ENABLE_TELEGRAM_NOTIFICATION")
	if telegramEnabled != "" {
		Config.NotificationChannel.Telegram.Enabled, err = strconv.ParseBool(telegramEnabled)
		if err != nil {
			return err
		}
	}

	telegramb64BotToken := os.Getenv("CMA_TELEGRAM_ENCODED_BOT_TOKEN")
	if telegramb64BotToken != "" {
		Config.NotificationChannel.Telegram.B64BotToken = telegramb64BotToken
	}

	telegramChatId := os.Getenv("CMA_TELEGRAM_CHAT_ID")
	if telegramChatId != "" {
		Config.NotificationChannel.Telegram.ChatId = telegramChatId
	}

	enableRadarr := os.Getenv("CMA_MONITOR_RADARR")
	if enableRadarr != "" {
		Config.Radarr.Enabled, err = strconv.ParseBool(enableRadarr)
		if err != nil {
			return err
		}
	}

	radarrUrl := os.Getenv("CMA_RADARR_URL")
	if radarrUrl != "" {
		Config.Radarr.URL = radarrUrl
	}

	radarrb64ApiKey := os.Getenv("CMA_RADARR_ENCODED_API_KEY")
	if radarrb64ApiKey != "" {
		Config.Radarr.B64APIKey = radarrb64ApiKey
	}

	radarrNotification := os.Getenv("CMA_RADARR_ENABLE_NOTIFICATION")
	if radarrNotification != "" {
		Config.Radarr.Notification, err = strconv.ParseBool(radarrNotification)
		if err != nil {
			return err
		}
	}

	enableSonarr := os.Getenv("CMA_MONITOR_SONARR")
	if enableSonarr != "" {
		Config.Sonarr.Enabled, err = strconv.ParseBool(enableSonarr)
		if err != nil {
			return err
		}
	}

	sonarrUrl := os.Getenv("CMA_SONARR_URL")
	if sonarrUrl != "" {
		Config.Sonarr.URL = sonarrUrl
	}

	sonarrb64ApiKey := os.Getenv("CMA_SONARR_ENCODED_API_KEY")
	if sonarrb64ApiKey != "" {
		Config.Sonarr.B64APIKey = sonarrb64ApiKey
	}

	sonarrNotification := os.Getenv("CMA_SONARR_ENABLE_NOTIFICATION")
	if sonarrNotification != "" {
		Config.Sonarr.Notification, err = strconv.ParseBool(sonarrNotification)
		if err != nil {
			return err
		}
	}

	// Write configuration

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
