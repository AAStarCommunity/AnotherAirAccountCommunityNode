package conf

type Provider struct {
	Alchemy string
}

func GetProvider() *Provider {
	return &Provider{
		Alchemy: getConf().Provider.Alchemy,
	}
}
