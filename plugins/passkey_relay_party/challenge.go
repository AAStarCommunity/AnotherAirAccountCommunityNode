package plugin_passkey_relay_party

import (
	"another_node/plugins/passkey_relay_party/seedworks"
	"strings"
)

func (rp *RelayParty) emailStartChallenge(mail, acceptLanguage string) error {

	captcha := seedworks.GenCaptcha(6)

	var body string
	if strings.EqualFold(acceptLanguage, "zh") {
		body = `
验证您的邮箱
<br />
<h2>` + captcha + `</h2>
此验证码将在 <b>10</b> 分钟后失效，非本人操作请忽略。<br />
`
	} else {
		body = `
Verify
<br />
<h2>` + captcha + `</h2>
Invalidate in <b>10</b> minutes, ignore if you were confused about this mail<br />
`
	}

	if err := sendMail(
		mail,
		"Verify Your Email",
		body,
		"text/html",
	); err != nil {
		return err
	}

	return rp.db.SaveChallenge(mail, captcha)
}

func (rp *RelayParty) emailFinishChallenge(mail, code string) error {
	if !rp.db.Challenge(mail, code) {
		return seedworks.ErrInvalidCaptcha{}
	}

	return nil
}
