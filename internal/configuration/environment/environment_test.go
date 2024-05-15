package environment_test

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBtcEnv_GetNetworkParams(t *testing.T) {
	t.Run("should return testnet params", func(t *testing.T) {
		env := &environment.BtcEnv{Network: "testnet"}
		params, err := env.GetNetworkParams()
		require.NoError(t, err)
		require.Equal(t, &chaincfg.TestNet3Params, params)
	})
	t.Run("should return mainnet params", func(t *testing.T) {
		env := &environment.BtcEnv{Network: "mainnet"}
		params, err := env.GetNetworkParams()
		require.NoError(t, err)
		require.Equal(t, &chaincfg.MainNetParams, params)
	})
	t.Run("should return regtest params", func(t *testing.T) {
		env := &environment.BtcEnv{Network: "regtest"}
		params, err := env.GetNetworkParams()
		require.NoError(t, err)
		require.Equal(t, &chaincfg.RegressionNetParams, params)
	})
	t.Run("should return error on unknown network", func(t *testing.T) {
		env := &environment.BtcEnv{Network: "simnet"}
		params, err := env.GetNetworkParams()
		require.ErrorContains(t, err, "invalid network name: simnet")
		require.Nil(t, params)
	})
}
