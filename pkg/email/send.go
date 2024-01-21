package email

import (
	"net/smtp"
	"strconv"
	"strings"
)

func SendEmail(email *Email) error {
	config, err := NewConfigFromEnv()
	if err != nil {
		return err
	}

	if err == nil {
		err = sendEmail(config.Username, config.Password, config.SmtpServer, int(SMTP_SERVICE), email.ToAddress, email.Subject, email.Body)
	}

	return err
}

func sendEmail(username, password, server string, port int, to []string, subject, body string) error {
	auth := smtp.PlainAuth("", username, password, server)

	msg := "From: " + username + "\r\n" +
		"To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body

	err := smtp.SendMail(server+":"+strconv.Itoa(port), auth, username, to, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
