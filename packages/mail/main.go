package main

import (
	"gogo/components/appctx"
	"gogo/components/mailer"
	"os"
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
