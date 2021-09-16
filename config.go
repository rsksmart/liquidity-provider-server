package main

import "github.com/rsksmart/liquidity-provider/providers"

type config struct {
	LogFile              string
	Debug                bool
	IrisActivationHeight int
	ErpKeys              []string

	Server struct {
		Port uint
	}
	DB struct {
		Path string
	}
	RSK struct {
		Endpoint                    string
		LBCAddr                     string
		BridgeAddr                  string
		RequiredBridgeConfirmations int64
	}
	BTC struct {
		Endpoint string
		Username string
		Password string
		Network  string
	}
	Provider providers.ProviderConfig
}
