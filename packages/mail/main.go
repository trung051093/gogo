package main

import (
	"os"
	"user_management/components/appctx"
	"user_management/components/mailer"
)

func main() {
	config := appctx.GetConfig()
	mail := os.Getenv("mail")
	mailService := mailer.NewMailer(config.GetMailConfig())
	mailService.SendMail(mailer.Mail{
		Sender:  config.Mail.Sender,
		To:      []string{mail},
		Subject: "Test send mail",
		Body:    "Hello world!!!",
	})
}
