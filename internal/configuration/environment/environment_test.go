package environment_test

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

const expectedSecretMask = "********"

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

func TestEnvironment_String_RedactsSecrets(t *testing.T) {
	env := environment.Environment{
		LpsStage:         "regtest",
		Port:             8080,
		LogLevel:         "debug",
		AwsLocalEndpoint: "http://localstack:4566",
		SecretSource:     "aws",
		WalletManagement: "native",
		AllowedOrigins:   []string{"http://localhost:3000"},
		Management: environment.ManagementEnv{
			EnableManagementApi:   true,
			SessionAuthKey:        "auth-secret",
			SessionEncryptionKey:  "encryption-secret",
			SessionTokenAuthKey:   "token-secret",
			UseHttps:              true,
			EnableSecurityHeaders: true,
		},
		Mongo: environment.MongoEnv{
			Username: "mongo-user",
			Password: "mongo-secret",
			Host:     "mongodb.local",
			Port:     27017,
		},
		Rsk: environment.RskEnv{
			Endpoint:                    "http://rsk.local",
			ChainId:                     31,
			EncryptedJsonSecret:         "key-secret-id",
			EncryptedJsonPasswordSecret: "password-secret-id",
			KeystoreFile:                "keystore.json",
			KeystorePassword:            "keystore-secret",
		},
		Btc: environment.BtcEnv{
			Network:  "regtest",
			Username: "btc-user",
			Password: "btc-secret",
			Endpoint: "http://btc.local",
		},
		Captcha: environment.CaptchaEnv{
			SecretKey: "captcha-secret",
			SiteKey:   "captcha-site",
			Url:       "http://captcha.local",
		},
	}

	output := env.String()

	require.NotContains(t, output, "mongo-secret")
	require.NotContains(t, output, "auth-secret")
	require.NotContains(t, output, "encryption-secret")
	require.NotContains(t, output, "token-secret")
	require.NotContains(t, output, "keystore-secret")
	require.NotContains(t, output, "btc-secret")
	require.NotContains(t, output, "captcha-secret")
	require.Contains(t, output, "regtest")
	require.Contains(t, output, "mongodb.local")
	require.Contains(t, output, "btc-user")
	require.Contains(t, output, "captcha-site")
	require.Contains(t, output, "key-secret-id")
	require.Contains(t, output, "password-secret-id")
	require.Equal(t, 7, strings.Count(output, expectedSecretMask))
}

func TestEnvironment_String_KeepsEmptySecretFieldsEmpty(t *testing.T) {
	env := environment.Environment{}

	require.NotContains(t, env.String(), expectedSecretMask)
}

func TestEnvironment_String_RedactsThroughDebugLogFormat(t *testing.T) {
	env := &environment.Environment{
		Mongo: environment.MongoEnv{Password: "mongo-secret"},
	}

	output := fmt.Sprintf("Environment loaded: %+v", env)

	require.NotContains(t, output, "mongo-secret")
	require.Contains(t, output, expectedSecretMask)
}
