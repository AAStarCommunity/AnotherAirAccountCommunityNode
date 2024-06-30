package account

type Provider interface {
	CreateAccount(*HdWallet) (string, error)
	GetRpc() string
}
