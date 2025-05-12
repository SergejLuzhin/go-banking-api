package service

import (
	"crypto/tls"
	"fmt"

	mail "github.com/go-mail/mail/v2"
)

type EmailService struct {
	Dialer *mail.Dialer
	From   string
}

func NewEmailService(host string, port int, user, pass string) *EmailService {
	d := mail.NewDialer(host, port, user, pass)
	d.TLSConfig = &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: false,
	}
	return &EmailService{
		Dialer: d,
		From:   user,
	}
}

func (s *EmailService) SendEmail(to string, subject string, body string) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if err := s.Dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("не удалось отправить email: %w", err)
	}
	return nil
}
