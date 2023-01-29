package internal

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
)

func SendEmailNotification(subject, body string) error {
	if !Config.NotificationChannel.SMTP.Enabled {
		log.Println("SMTP not enabled.")
		return errors.New("SMTP not enabled")
	}
	log.Println("Preparing to send email notification")
	password, err := Base64Decode(Config.NotificationChannel.SMTP.B64Password)
	if err != nil {
		log.Println("Failed to send Email notification", err.Error())
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
