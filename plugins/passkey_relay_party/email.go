package plugin_passkey_relay_party

import (
	"another_node/plugins/passkey_relay_party/conf"

	"gopkg.in/gomail.v2"
)

func ReplyEmail(replyTo, subject, htmlBody string) error {

	err := sendMail(
		replyTo,
		subject,
		htmlBody,
		"text/html")

	return err
}

func sendMail(to, subject, body, mailtype string) error {
	mailCfg := conf.Get().Mail

	m := gomail.NewMessage()

	m.SetHeader("From", mailCfg.User)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody(mailtype, body)
	d := gomail.NewDialer(mailCfg.Host, mailCfg.Port, mailCfg.User, mailCfg.Password)
	return d.DialAndSend(m)
}
