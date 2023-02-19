package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"
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

	type reqBody struct {
		Title    string `json:"title"`
		Message  string `json:"message"`
		Priority int    `json:"priority"`
	}

	gotifyURL := Config.NotificationChannel.Gotify.URL
	apiUrl := gotifyURL + "/message"

	appToken, err := Base64Decode(Config.NotificationChannel.Gotify.B64AppToken)
	if err != nil {
		return err
	}

	// Create request

	var reqBodyValue reqBody

	reqBodyValue.Title = title
	reqBodyValue.Message = body
	reqBodyValue.Priority = Config.NotificationChannel.Gotify.Priority

	reqBodyValueJson, err := json.Marshal(reqBodyValue)
	if err != nil {
		log.Println("Failed to send Gotify notification", err.Error())
		return err
	}

	requestBodyJson := bytes.NewReader(reqBodyValueJson)
	if err != nil {
		log.Println("Failed to send Gotify notification", err.Error())
		return err
	}

	req, err := http.NewRequest(http.MethodPost, apiUrl, requestBodyJson)
	if err != nil {
		log.Println("Failed to send Gotify notification", err.Error())
		return err
	}
	req.Header.Set("Authorization", "Bearer "+appToken)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Create client
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request
	res, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send Gotify notification", err.Error())
		return err
	}

	if res.StatusCode/100 != 2 {
		log.Println("Failed to send Gotify notification", res.Status)
		return errors.New("Failed to send Gotify notification")
	}

	log.Println("Notification sent to Gotify server", gotifyURL)

	return nil
}
