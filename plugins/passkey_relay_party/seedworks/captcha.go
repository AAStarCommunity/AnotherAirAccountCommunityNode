package seedworks

import "math/rand"

const seed = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenCaptcha generate captcha
func GenCaptcha(n int) string {
	if n < 4 {
		n = 4
	}
	captcha := make([]byte, n)
	for i := 0; i < n; i++ {
		rnd := rand.Int() % len(seed)
		captcha[i] = seed[rnd]
	}
	return string(captcha)
}
