package integration

import (
	"encoding/json"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

type InterfacesConfig struct {
	Lps struct {
		Url string `json:"url"`
	} `json:"lps"`
	Btc struct {
		RpcEndpoint         string `json:"rpcEndpoint"`
		User                string `json:"user"`
		Password            string `json:"password"`
		WalletPassword      string `json:"walletPassword"`
		WalletUnlockSeconds int64  `json:"walletUnlockSeconds"`
	} `json:"btc"`
	Rsk struct {
		RpcUrl         string `json:"rpcUrl"`
		LbcAddress     string `json:"lbcAddress"`
		UserPrivateKey string `json:"userPrivateKey"`
	} `json:"rsk"`
}

type TestConfig struct {
	Mempool struct {
	} `json:"mempool"`
	Pegin struct {
		Value              *big.Int `json:"value"`
		DestinationAddress string   `json:"destinationAddress"`
		RefundAddress      string   `json:"refundAddress"`
		Data               string   `json:"data"`
		CallForUserWait    int64    `json:"callForUserWait"`
		RegisterPeginWait  int64    `json:"registerPeginWait"`
	} `json:"pegin"`
	Pegout struct {
		Value              *big.Int `json:"value"`
		DestinationAddress string   `json:"destinationAddress"`
		RefundAddress      string   `json:"refundAddress"`
		PegoutPaymentWait  int64    `json:"pegoutPaymentWait"`
		PegoutRefundWait   int64    `json:"pegoutRefundWait"`
	} `json:"pegout"`
}

type Config struct {
	InterfacesConfig
	Network string     `json:"network"`
	Tests   TestConfig `json:"tests"`
}

func ReadTestConfig(t *testing.T) *Config {
	const configPath = "test/integration/integration-test.config.json"
	config := new(Config)
	configBytes := test.ReadFile(t, configPath)
	require.NoError(t, json.Unmarshal(configBytes, config))
	return config
}
