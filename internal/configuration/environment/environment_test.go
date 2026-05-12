package environment_test

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/test"
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

func TestEclipseEnv_FillWithDefaults(t *testing.T) {
	env := &environment.EclipseEnv{
		RskToleranceThreshold:    0,
		RskMaxMsWaitForBlock:     0,
		RskWaitPollingMsInterval: 0,
		BtcToleranceThreshold:    0,
		BtcMaxMsWaitForBlock:     0,
		BtcWaitPollingMsInterval: 0,
		AlertCooldownSeconds:     0,
	}
	defaults := env.FillWithDefaults()
	require.Equal(t, uint8(50), defaults.RskToleranceThreshold)
	require.Equal(t, uint64(10_000), defaults.RskMaxMsWaitForBlock)
	require.Equal(t, uint64(1000), defaults.RskWaitPollingMsInterval)
	require.Equal(t, uint8(50), defaults.BtcToleranceThreshold)
	require.Equal(t, uint64(60_000), defaults.BtcMaxMsWaitForBlock)
	require.Equal(t, uint64(10_000), defaults.BtcWaitPollingMsInterval)
	require.Equal(t, uint64(30*60), defaults.AlertCooldownSeconds) // 30 min
	test.AssertMaxZeroValues(t, defaults, 1)
}

func TestEclipseEnv_ToConfig(t *testing.T) {
	env := &environment.EclipseEnv{
		RskToleranceThreshold:    50,
		RskMaxMsWaitForBlock:     10000,
		RskWaitPollingMsInterval: 1000,
		BtcToleranceThreshold:    50,
		BtcMaxMsWaitForBlock:     60000,
		BtcWaitPollingMsInterval: 10000,
	}
	config := env.ToConfig()
	require.Equal(t, uint8(50), config.RskToleranceThreshold)
	require.Equal(t, uint64(10000), config.RskMaxMsWaitForBlock)
	require.Equal(t, uint64(1000), config.RskWaitPollingMsInterval)
	require.Equal(t, uint8(50), config.BtcToleranceThreshold)
	require.Equal(t, uint64(60000), config.BtcMaxMsWaitForBlock)
	require.Equal(t, uint64(10000), config.BtcWaitPollingMsInterval)
	test.AssertNonZeroValues(t, config)
}
