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
	PeginProviderName    string   `env:"PEGIN_PROVIDER_NAME"`
	PeginFee		     uint  	  `env:"PEGIN_FEE"`
	PeginQuoteExp	     uint  	  `env:"PEGIN_QUOTE_EXPIRATION"`
	PeginAcceptedQuoteExp uint    `env:"PEGIN_ACCEPTED_QUOTE_EXPIRATION"`
	PeginMinTransactValue uint64    `env:"PEGIN_MIN_TRANSACTION_VALUE"`
	PeginMaxTransactValue uint64   `env:"PEGIN_MAX_TRANSACTION_VALUE"`
	PegoutProviderName    string  `env:"PEGOUT_PROVIDER_NAME"`
	PegoutFee		      uint    `env:"PEGOUT_FEE"`
	PegoutQuoteExp	      uint    `env:"PEGOUT_QUOTE_EXPIRATION"`
	PegoutAcceptedQuoteExp uint   `env:"PEGOUT_ACCEPTED_QUOTE_EXPIRATION"`
	PegoutMinTransactValue uint64   `env:"PEGOUT_MIN_TRANSACTION_VALUE"`
	PegoutMaxTransactValue uint64   `env:"PEGOUT_MAX_TRANSACTION_VALUE"`
	BaseURL string   `env:"BASE_URL"`

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
