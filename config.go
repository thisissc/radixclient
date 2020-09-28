package radixclient

var (
	_RadixConfigMap = make(map[string]RadixConfig, 0)
)

type RadixConfig struct {
	Name          string
	Addr          string
	Password      string
	MinPool       int
	MaxPool       int
	DrainInterval int64
	IsCluster     bool
}

type Config []RadixConfig

func (c *Config) Init() error {
	for _, rcfg := range []RadixConfig(*c) {
		_RadixConfigMap[rcfg.Name] = rcfg
	}

	return nil
}
