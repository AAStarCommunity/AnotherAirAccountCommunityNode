package wallet_storage

type ErrInvalidWalletData struct{}

var _ error = ErrInvalidWalletData{}

func (e ErrInvalidWalletData) Error() string {
	return string("invalid wallet data")
}
