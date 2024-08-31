package seedworks

type ErrUserNotFound struct{}
type ErrEmailEmpty struct{}
type ErrInvalidCaptcha struct{}
type ErrUserAlreadyExists struct{}
type ErrWalletNotFound struct{}
type ErrChainNotFound struct{}

var _ error = ErrUserNotFound{}
var _ error = ErrEmailEmpty{}
var _ error = ErrInvalidCaptcha{}
var _ error = ErrUserAlreadyExists{}
var _ error = ErrWalletNotFound{}
var _ error = ErrChainNotFound{}

func (e ErrUserNotFound) Error() string {
	return string("user not found")
}

func (e ErrEmailEmpty) Error() string {
	return string("email is empty")
}

func (e ErrInvalidCaptcha) Error() string {
	return string("invalid captcha")
}

func (e ErrUserAlreadyExists) Error() string {
	return string("user already exists")
}

func (e ErrWalletNotFound) Error() string {
	return string("wallet not found")
}

func (e ErrChainNotFound) Error() string {
	return string("chain not found")
}
