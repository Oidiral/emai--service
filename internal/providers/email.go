package providers

import (
	"context"

	"github.com/Oidiral/emai--service/internal/models"
	"gopkg.in/gomail.v2"
)

type EmailProvider struct {
	Login    string
	Pass     string
	SmtpHost string
	SmtpPort int
}

func NewEmailProvider(login, pass, smtpHost string, smtpPort int) *EmailProvider {
	return &EmailProvider{
		Login:    login,
		Pass:     pass,
		SmtpHost: smtpHost,
		SmtpPort: smtpPort,
	}
}

func (p *EmailProvider) Slug() string {
	return models.EmailProvider
}

func (p *EmailProvider) Send(ctx context.Context, email *models.ProviderEmail) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", p.Login)
	msg.SetHeader("To", email.Email)
	msg.SetHeader("Subject", email.Subject)
	msg.SetBody("text/html", email.Text)

	n := gomail.NewDialer(p.SmtpHost, p.SmtpPort, p.Login, p.Pass)

	return n.DialAndSend(msg)
}
