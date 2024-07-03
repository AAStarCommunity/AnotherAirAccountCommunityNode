package seedworks

type ErrUserNotFound struct{}
type ErrEmailEmpty struct{}
type ErrInvalidCaptcha struct{}

var _ error = ErrUserNotFound{}
var _ error = ErrEmailEmpty{}
var _ error = ErrInvalidCaptcha{}

func (e ErrUserNotFound) Error() string {
	return string("user not found")
}

func (e ErrEmailEmpty) Error() string {
	return string("email is empty")
}

func (e ErrInvalidCaptcha) Error() string {
	return string("invalid captcha")
}
