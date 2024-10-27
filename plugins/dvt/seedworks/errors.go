package seedworks

type ErrSignatureVerifyFailed struct{}
type ErrThresholdGreaterThanTotal struct{}
type ErrNotEnoughSigners struct{}
type ErrThresholdGreaterThanZero struct{}

var _ error = ErrSignatureVerifyFailed{}
var _ error = ErrThresholdGreaterThanTotal{}
var _ error = ErrNotEnoughSigners{}
var _ error = ErrThresholdGreaterThanZero{}

func (e ErrSignatureVerifyFailed) Error() string {
	return string("signature verify failed")
}

func (e ErrThresholdGreaterThanTotal) Error() string {
	return "threshold must be less than or equal to total"
}

func (e ErrNotEnoughSigners) Error() string {
	return "not enough signers"
}

func (e ErrThresholdGreaterThanZero) Error() string {
	return "threshold must be greater than zero"
}
