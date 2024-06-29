package bls_tss

type ErrThresholdGreaterThanTotal struct {
}

func (e *ErrThresholdGreaterThanTotal) Error() string {
	return "threshold must be less than or equal to total"
}

type ErrNotEnoughSigners struct {
}

func (e *ErrNotEnoughSigners) Error() string {
	return "not enough signers"
}

type ErrThresholdGreaterThanZero struct {
}

func (e *ErrThresholdGreaterThanZero) Error() string {
	return "threshold must be greater than zero"
}
