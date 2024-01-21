package email

import "fmt"

type Email struct {
	ToAddress   []string
	FromAddress string
	Subject     string
	Body        string
}

func NewEmail(to []string, from, subject, body string) *Email {
	return &Email{
		ToAddress:   to,
		FromAddress: from,
		Subject:     subject,
		Body:        body,
	}
}

func (e *Email) String() string {
	bodyLength := len(e.Body)
	if bodyLength > 10 {
		bodyLength = 10
	}
	return fmt.Sprintf("To: %s, From: %s, Subject: %s, Body: %s", e.ToAddress, e.FromAddress, e.Subject, e.Body[:bodyLength])
}
