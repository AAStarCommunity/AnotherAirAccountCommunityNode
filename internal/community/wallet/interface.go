package wallet

type Provider interface {
	CreateAccount(*HdWallet) (string, error)
	GetRpc() string
}
