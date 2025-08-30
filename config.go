package rrkurrentdb

type Config struct {
	Address string `mapstructure:"address"`
}

func (cfg *Config) InitDefaults() {
	if cfg.Address == "" {
		cfg.Address = "kurrentdb+discover://admin:changeit@localhost:2113"
	}
}
