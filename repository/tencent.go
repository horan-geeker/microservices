package repository

import (
	"gopkg.in/gomail.v2"
	"microservices/entity/options"
	"microservices/pkg/consts"
	"net/mail"
)

type Tencent interface {
	SendMailTo(mailTo string, from string, title string, content string) error
}

type tencent struct {
	mailServerAddress  string
	mailServerPassword string
}

// SendMailTo .
func (t *tencent) SendMailTo(mailTo string, from string, title string, content string) error {
	m := gomail.NewMessage()
	fromMail := mail.Address{from, t.mailServerAddress}
	m.SetHeader("From", fromMail.String())
	m.SetHeader("To", mailTo)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)

	d := gomail.NewDialer(consts.TencentSmtpServer, consts.TencentSmtpPort, t.mailServerAddress,
		t.mailServerPassword)
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func newTencent(tencentOpt *options.TencentOptions) Tencent {
	return &tencent{
		mailServerAddress:  tencentOpt.MailServerAddress,
		mailServerPassword: tencentOpt.MailServerPassword,
	}
}
