package plugin_passkey_relay_party

import (
	"another_node/plugins/passkey_relay_party/seedworks"
	"os"
	"testing"
)

func TestSendEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	os.Setenv("Env", "dev")
	defer os.Unsetenv("Env")

	os.Chdir("../../")

	captcha := seedworks.GenCaptcha(6)

	body := `
验证您的邮箱
<br />
<h2>` + captcha + `</h2>
此验证码将在 <b>10</b> 分钟后失效，非本人操作请忽略。<br />
`

	err := ReplyEmail("邮箱验证", "993921@qq.com", "AAStar验证邮箱", body)

	if err != nil {
		t.Errorf("Send email failed: %v", err)
	}
}
