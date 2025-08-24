package rrkurrentdb

type Config struct {
	Address string `mapstructure:"address"`
}

func (cfg *Config) InitDefaults() {
	if cfg.Address == "" {
		cfg.Address = "esdb://localhost:2113"
	}
}
