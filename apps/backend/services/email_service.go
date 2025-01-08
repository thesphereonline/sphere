package services

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

// SendEmail sends an email using SMTP
func SendEmail(to, subject, body string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", os.Getenv("SMTP_USERNAME"))
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		587,
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}
	return nil
}
