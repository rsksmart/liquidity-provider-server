package main

import (
	"github.com/rsksmart/liquidity-provider-server/http"
	"github.com/rsksmart/liquidity-provider/providers"
)

type config struct {
	LogFile              string
	Debug                bool
	IrisActivationHeight int
	ErpKeys              []string
	MaxQuoteValue        uint64
	SimultaneouslyQuotes int

	Server struct {
		Port uint
	}
	DB struct {
		Path string
	}
	RSK http.LiquidityProviderList
	BTC struct {
		Endpoint string
		Username string
		Password string
		Network  string
	}
	Provider providers.ProviderConfig
}
