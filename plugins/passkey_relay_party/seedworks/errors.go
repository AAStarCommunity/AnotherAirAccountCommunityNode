package seedworks

type UserNotFoundError struct{}
type EmailEmptyError struct{}

var _ error = UserNotFoundError{}
var _ error = EmailEmptyError{}

func (e UserNotFoundError) Error() string {
	return string("user not found")
}

func (e EmailEmptyError) Error() string {
	return string("email is empty")
}
