package internal

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
	Enabled   bool     `yaml:"enabled"`
	Server    string   `yaml:"server"`
	Port      int      `yaml:"port"`
	Security  Security `yaml:"security"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password"`
	FromEmail string   `yaml:"fromEmail"`
	ToEmail   string   `yaml:"toEmail"`
	CcEmail   string   `yaml:"ccEmail"`
	BccEmail  string   `yaml:"bccEmail"`
}

type GotifyConfig struct {
	Enabled  bool   `yaml:"enabled"`
	URL      string `yaml:"url"`
	AppToken string `yaml:"appToken"`
}

type TelegramConfig struct {
	Enabled  bool   `yaml:"enabled"`
	BotToken string `yaml:"botToken"`
	ChatId   string `yaml:"chatId"`
}

type NotificationChannel struct {
	SMTP     SMTPConfig
	Gotify   GotifyConfig
	Telegram TelegramConfig
}

type RadarrConfig struct {
	Enabled      bool   `yaml:"enabled"`
	URL          string `yaml:"url"`
	APIKey       string `yaml:"apiKey"`
	Notification bool   `yaml:"notification"`
}

type SonarrConfig struct {
	Enabled      bool   `yaml:"enabled"`
	URL          string `yaml:"url"`
	APIKey       string `yaml:"apiKey"`
	Notification bool   `yaml:"notification"`
}

type Config struct {
	Period              Interval            `yaml:"period"`
	DefaultCleanupTime  int                 `yaml:"defaultCleanupTime"`
	NotificationChannel NotificationChannel `yaml:"notificationChannel"`
	Radarr              RadarrConfig        `yaml:"radarr"`
	Sonarr              SonarrConfig        `yaml:"sonarr"`
}
