package email

import (
	"fmt"
	"os"

	"github.com/emersion/go-imap/client"
	"github.com/joho/godotenv"
)

const (
	IMAP_SERVICE EmailService = 993
	SMTP_SERVICE EmailService = 587
)

type EmailService int

type Config struct {
	Username   string
	Password   string
	ImapServer string
	SmtpServer string
}

func NewConfig(username, password, imapServer, smtpServer string) *Config {
	return &Config{
		Username:   username,
		Password:   password,
		ImapServer: imapServer,
		SmtpServer: smtpServer,
	}
}

func NewConfigFromEnv() (*Config, error) {
	var username string
	var password string
	var fetchServer string
	var sendServer string

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %s", err)
	}

	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")
	fetchServer = os.Getenv("IMAP_SERVER")
	sendServer = os.Getenv("SMTP_SERVER")

	if len(username) == 0 || len(password) == 0 || len(fetchServer) == 0 || len(sendServer) == 0 {
		return nil, fmt.Errorf("missing configuration variables")
	}

	return NewConfig(username, password, fetchServer, sendServer), nil
}

func (c *Config) ConnectToServer(service EmailService) (*client.Client, error) {
	port := int(service)
	var host string
	if service == SMTP_SERVICE {
		host = c.SmtpServer
	} else {
		host = c.ImapServer
	}
	client, err := client.DialTLS(fmt.Sprintf("%s:%d", host, port), nil)
	if err != nil {
		return nil, err
	}

	if err := client.Login(c.Username, c.Password); err != nil {
		return nil, err
	}

	return client, nil
}
