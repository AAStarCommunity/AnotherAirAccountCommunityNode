package plugin_passkey_relay_party

import (
	"another_node/plugins/passkey_relay_party/conf"

	"gopkg.in/gomail.v2"
)

func ReplyEmail(title, replyTo, subject, htmlBody string) error {

	err := sendMail(
		title,
		replyTo,
		subject,
		htmlBody,
		"text/html")

	return err
}

func sendMail(title, to, subject, body, mailtype string) error {
	mailCfg := conf.Get().Mail

	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress(mailCfg.User, title))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody(mailtype, body)
	d := gomail.NewDialer(mailCfg.Host, mailCfg.Port, mailCfg.User, mailCfg.Password)
	return d.DialAndSend(m)
}
