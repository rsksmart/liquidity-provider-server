package main

import (
	"github.com/rsksmart/liquidity-provider-server/http"
	"github.com/rsksmart/liquidity-provider-server/pegin"
	"github.com/rsksmart/liquidity-provider-server/pegout"
)

type config struct {
	LogFile              string   `env:"LOG_FILE"`
	Debug                bool     `env:"DEBUG"`
	IrisActivationHeight int      `env:"IRIS_ACTIVATION_HEIGHT"`
	ErpKeys              []string `env:"ERP_KEYS"`
	MaxQuoteValue        uint64   `env:"MAX_QUOTE_VALUE"`

	Server struct {
		Port uint `env:"SERVER_PORT"`
	}
	DB struct {
		Regtest struct {
			Host     string `env:"DB_REGTEST_HOST"`
			Database string `env:"DB_REGTEST_DATABASE"`
			User     string `env:"DB_REGTEST_USER"`
			Password string `env:"DB_REGTEST_PASSWORD"`
			Port     uint   `env:"DB_REGTEST_PORT"`
		}
		Path string `env:"DB_PATH"`
	}
	RSK http.LiquidityProviderList
	BTC struct {
		Endpoint string `env:"BTC_ENDPOINT"`
		Username string `env:"BTC_USERNAME"`
		Password string `env:"BTC_PASSWORD"`
		Network  string `env:"BTC_NETWORK"`
	}
	Provider      *pegin.ProviderConfig  `env:",prefix=PEGIN_PROVIDER_"`
	PegoutProvier *pegout.ProviderConfig `env:",prefix=PEGOUT_PROVIDER_"`
}
