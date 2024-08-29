package plugin_passkey_relay_party

import (
	"another_node/plugins/passkey_relay_party/seedworks"
	"strings"
)

func (rp *RelayParty) emailStartChallenge(mail, acceptLanguage string) error {

	captcha := seedworks.GenCaptcha(6)

	if err := rp.db.SaveChallenge(mail, captcha); err != nil {
		return err
	}

	var title, subject, body string
	if strings.EqualFold(acceptLanguage, "zh") || strings.EqualFold(acceptLanguage, "zh-cn") {
		title = "请验证您的AirAccount注册邮箱"
		subject = "验证您的邮箱"
		body = `
验证您的邮箱
<br />
<h2>` + captcha + `</h2>
此验证码将在 <b>10</b> 分钟后失效，非本人操作请忽略。<br />
`
	} else {
		title = "Verification Code for your AirAccount"
		subject = "Verify Your Email"
		body = `
Verify
<br />
<h2>` + captcha + `</h2>
Invalidate in <b>10</b> minutes, ignore it if you were confused about this mail<br />
`
	}

	if err := sendMail(
		title,
		mail,
		subject,
		body,
		"text/html",
	); err != nil {
		return err
	}

	return nil
}

func (rp *RelayParty) emailChallenge(mail, code string) error {
	if !rp.db.Challenge(mail, code) {
		return seedworks.ErrInvalidCaptcha{}
	}

	return nil
}
