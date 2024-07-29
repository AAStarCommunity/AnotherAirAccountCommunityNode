package seedworks

type ErrUserNotFound struct{}
type ErrEmailEmpty struct{}
type ErrInvalidCaptcha struct{}
type ErrUserAlreadyExists struct{}

var _ error = ErrUserNotFound{}
var _ error = ErrEmailEmpty{}
var _ error = ErrInvalidCaptcha{}
var _ error = ErrUserAlreadyExists{}

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
