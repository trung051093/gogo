package mailer

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type MailService interface {
	SendMail(mail Mail) error
}

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

type MailConfig struct {
	Sender   string
	Host     string
	Port     int
	Username string
	Password string
}

type mailService struct {
	config MailConfig
}

func NewMailer(config MailConfig) MailService {
	return &mailService{config: config}
}

func (m *mailService) GetAuth() smtp.Auth {
	return smtp.PlainAuth("", m.config.Sender, m.config.Password, m.config.Host)
}

func (m *mailService) GetAddress() string {
	return fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)
}

func (m *mailService) SendMail(mail Mail) error {
	msg := m.BuildMessage(mail)
	err := smtp.SendMail(m.GetAddress(), m.GetAuth(), mail.Sender, mail.To, []byte(msg))
	if err != nil {
		log.Println("Send mail error: ", err)
	}
	return nil
}

func (m *mailService) BuildMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}
