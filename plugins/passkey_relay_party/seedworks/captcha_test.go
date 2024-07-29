package seedworks

import "testing"

func TestCaptchaGenerator(t *testing.T) {
	for i := 4; i <= 8; i++ {
		captcha := GenCaptcha(i)
		if len(captcha) != i {
			t.Errorf("Captcha length is not %d", i)
		}
	}
	for i := 0; i <= 4; i++ {
		captcha := GenCaptcha(i)
		if len(captcha) != 4 {
			t.Errorf("Captcha length is not %d", 4)
		}
	}
}
