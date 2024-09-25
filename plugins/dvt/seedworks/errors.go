package seedworks

type ErrSignatureVerifyFailed struct{}

var _ error = ErrSignatureVerifyFailed{}

func (e ErrSignatureVerifyFailed) Error() string {
	return string("signature verify failed")
}
