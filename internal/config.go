package internal

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Interval string

const (
	Daily      Interval = "daily"
	Every3Days Interval = "every3days"
	Weekly     Interval = "weekly"
	Bimonthly  Interval = "bimonthly"
	Monthly    Interval = "monthly"
)

type Security string

const (
	None Security = "none"
	TLS  Security = "tls"
)

type SMTPConfig struct {
	Enabled     bool     `yaml:"enabled"`
	Server      string   `yaml:"server"`
	Port        int      `yaml:"port"`
	Security    Security `yaml:"security"`
	Username    string   `yaml:"username"`
	B64Password string   `yaml:"b64Password"`
	FromEmail   string   `yaml:"fromEmail"`
	ToEmail     string   `yaml:"toEmail"`
	CcEmail     string   `yaml:"ccEmail"`
	BccEmail    string   `yaml:"bccEmail"`
}

type GotifyConfig struct {
	Enabled     bool   `yaml:"enabled"`
	URL         string `yaml:"url"`
	B64AppToken string `yaml:"b64AppToken"`
}

type TelegramConfig struct {
	Enabled     bool   `yaml:"enabled"`
	B64BotToken string `yaml:"b64BotToken"`
	ChatId      string `yaml:"chatId"`
}

type NotificationChannel struct {
	SMTP     SMTPConfig
	Gotify   GotifyConfig
	Telegram TelegramConfig
}

type RadarrConfig struct {
	Enabled      bool   `yaml:"enabled"`
	URL          string `yaml:"url"`
	B64APIKey    string `yaml:"b64ApiKey"`
	Notification bool   `yaml:"notification"`
}

type SonarrConfig struct {
	Enabled      bool   `yaml:"enabled"`
	URL          string `yaml:"url"`
	B64APIKey    string `yaml:"b64ApiKey"`
	Notification bool   `yaml:"notification"`
}

type Configuration struct {
	MaintenanceCycle    Interval            `yaml:"maintenanceCycle"`
	DeleteAfterDays     int                 `yaml:"deleteAfterDays"`
	NotificationChannel NotificationChannel `yaml:"notificationChannel"`
	Radarr              RadarrConfig        `yaml:"radarr"`
	Sonarr              SonarrConfig        `yaml:"sonarr"`
}

func MaintenanceCycleInInt(period Interval) int {
	if period == Daily {
		return 1
	} else if period == Every3Days {
		return 3
	} else if period == Weekly {
		return 7
	} else if period == Bimonthly {
		return 15
	} else if period == Monthly {
		return 30
	} else {
		return 0
	}
}

var Config Configuration

func ReadConfig(configFile string) (*Configuration, error) {
	log.Println("Reading configurations")
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println("Failed to read configuration file.", err.Error())
		return nil, err
	}

	// yaml.Unmarshal(data, &config)
	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		log.Println("Failed to read configuration file.", err.Error())
		return nil, err
	}

	return &Config, nil
}
