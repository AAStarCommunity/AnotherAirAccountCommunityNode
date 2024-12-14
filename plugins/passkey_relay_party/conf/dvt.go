package conf

// DVT represents the configuration for DVT (Distributed Validation Technology)
type DVT struct {
	Disable   bool     `yaml:"disable"`
	Threshold int      `yaml:"threshold"`
	Nodes     []string `yaml:"nodes"`
}

// GetDVT returns the DVT configuration section from the global config
func GetDVT() *DVT {
	return &Get().DVT
}
