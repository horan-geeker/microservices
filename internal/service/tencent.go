package service

import (
	"gopkg.in/gomail.v2"
	"microservices/internal/config"
	"microservices/pkg/consts"
	"net/mail"
)

// SendMailByTencent .
func SendMailByTencent(mailTo string, from string, title string, content string) error {
	m := gomail.NewMessage()
	fromMail := mail.Address{from, config.Env.MailServerAddress}
	m.SetHeader("From", fromMail.String())
	m.SetHeader("To", mailTo)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)

	d := gomail.NewDialer(consts.TencentSmtpServer, consts.TencentSmtpPort, config.Env.MailServerAddress,
		config.Env.MailServerPassword)
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
