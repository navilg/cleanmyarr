package internal

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"net/url"

	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client/message"
	"github.com/gotify/go-api-client/v2/gotify"
	"github.com/gotify/go-api-client/v2/models"
)

func SendEmailNotification(subject, body string) error {
	if !Config.NotificationChannel.SMTP.Enabled {
		log.Println("SMTP not enabled.")
		return errors.New("SMTP not enabled")
	}
	log.Println("Preparing to send email notification")
	password, err := Base64Decode(Config.NotificationChannel.SMTP.B64Password)
	if err != nil {
		log.Println("Failed to send Email notification")
		return err
	}
	auth := smtp.PlainAuth("", Config.NotificationChannel.SMTP.Username, password, Config.NotificationChannel.SMTP.Server)
	var CcListString string = ""
	for _, cc := range Config.NotificationChannel.SMTP.CcEmail {
		CcListString = CcListString + "," + cc
	}

	var BccListString string = ""
	for _, bcc := range Config.NotificationChannel.SMTP.BccEmail {
		BccListString = BccListString + "," + bcc
	}

	var ToListString string = ""
	for _, to := range Config.NotificationChannel.SMTP.ToEmail {
		ToListString = ToListString + "," + to
	}

	message := []byte("To: " + ToListString + "\r\n" + "Cc: " + CcListString + "\r\n" + "Bcc: " + BccListString + "\r\n" + "Subject: " + subject + "\r\n" + body)

	err = smtp.SendMail(Config.NotificationChannel.SMTP.Server+fmt.Sprintf(":%d", Config.NotificationChannel.SMTP.Port), auth, Config.NotificationChannel.SMTP.FromEmail, Config.NotificationChannel.SMTP.ToEmail, message)
	if err != nil {
		log.Println("Failed to send Email notification", err.Error())
		return err
	}

	log.Println("Email notification sent successfully.")
	log.Println("To:", Config.NotificationChannel.SMTP.ToEmail)
	log.Println("CC:", Config.NotificationChannel.SMTP.CcEmail)
	log.Println("BCC:", Config.NotificationChannel.SMTP.BccEmail)

	return nil
}

func SendGotifyNotification(title, body string) error {
	if !Config.NotificationChannel.Gotify.Enabled {
		log.Println("Gotify not enabled.")
		return errors.New("Gotify not enabled")
	}
	log.Println("Preparing to send gotify notification")
	gotifyURL := Config.NotificationChannel.Gotify.URL
	url, _ := url.Parse(gotifyURL)

	appToken, err := Base64Decode(Config.NotificationChannel.Gotify.B64AppToken)
	if err != nil {
		log.Println("Failed to send Gotify notification")
		return err
	}

	client := gotify.NewClient(url, &http.Client{})

	params := message.NewCreateMessageParams()
	params.Body = &models.MessageExternal{
		Title:    title,
		Message:  body,
		Priority: Config.NotificationChannel.Gotify.Priority,
	}

	_, err = client.Message.CreateMessage(params, auth.TokenAuth(appToken))
	if err != nil {
		log.Println("Failed to send Gotify notification", err.Error())
		return err
	}

	log.Println("Notification sent to Gotify URL", gotifyURL)

	return nil
}
