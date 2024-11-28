package scripts_test

import (
	"flag"
	"github.com/rsksmart/liquidity-provider-server/cmd/utils/scripts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	awsMockEndpoint    = "http://localhost:1122"
	rskMockEndpoint    = "http://localhost:3344"
	mockFilePath       = "/file/path"
	mockKeystoreSecret = "UnitTest/Keystore-Secret"
	mockPwdSecret      = "UnitTest/Password-Secret"
)

func TestReadBaseInput(t *testing.T) {
	t.Run("should read base input", func(t *testing.T) {
		input := new(scripts.BaseInput)
		scripts.ReadBaseInput(input)
		require.NoError(t, flag.Set("network", "regtest"))
		require.NoError(t, flag.Set("secret-src", "aws"))
		require.NoError(t, flag.Set("aws-endpoint", awsMockEndpoint))
		require.NoError(t, flag.Set("rsk-endpoint", rskMockEndpoint))
		require.NoError(t, flag.Set("lbc-address", "0xBEd51d83cc4676660E3Fc3819dfAD8238549B975"))
		require.NoError(t, flag.Set("keystore-secret", mockKeystoreSecret))
		require.NoError(t, flag.Set("password-secret", mockPwdSecret))
		require.NoError(t, flag.Set("keystore-file", mockFilePath))
		flag.Parse()
		assert.Equal(t, "regtest", input.Network)
		assert.Equal(t, "aws", input.SecretSource)
		assert.Equal(t, awsMockEndpoint, input.AwsLocalEndpoint)
		assert.Equal(t, rskMockEndpoint, input.RskEndpoint)
		assert.Equal(t, "0xBEd51d83cc4676660E3Fc3819dfAD8238549B975", input.CustomLbcAddress)
		assert.Equal(t, mockKeystoreSecret, input.EncryptedJsonSecret)
		assert.Equal(t, mockPwdSecret, input.EncryptedJsonPasswordSecret)
		assert.Equal(t, mockFilePath, input.KeystoreFile)
		assert.Empty(t, input.KeystorePassword)
	})
}

func TestBaseInput_ToEnv(t *testing.T) {
	t.Run("should convert base input to environment with aws source", func(t *testing.T) {
		input := scripts.BaseInput{
			Network:                     "regtest",
			RskEndpoint:                 rskMockEndpoint,
			CustomLbcAddress:            "0xBEd51d83cc4676660E3Fc3819dfAD8238549B975",
			AwsLocalEndpoint:            awsMockEndpoint,
			SecretSource:                "aws",
			EncryptedJsonSecret:         mockKeystoreSecret,
			EncryptedJsonPasswordSecret: mockPwdSecret,
			KeystoreFile:                mockFilePath,
		}
		env, err := input.ToEnv(func(fd int) ([]byte, error) {
			return []byte(""), nil
		})
		require.NoError(t, err)
		assert.Equal(t, "regtest", env.LpsStage)
		assert.Equal(t, awsMockEndpoint, env.AwsLocalEndpoint)
		assert.Equal(t, "aws", env.SecretSource)
		assert.Equal(t, "native", env.WalletManagement)
		assert.Equal(t, rskMockEndpoint, env.Rsk.Endpoint)
		assert.Equal(t, uint64(33), env.Rsk.ChainId)
		assert.Equal(t, "0xBEd51d83cc4676660E3Fc3819dfAD8238549B975", env.Rsk.LbcAddress)
		assert.Equal(t, "0x0000000000000000000000000000000001000006", env.Rsk.BridgeAddress)
		assert.Equal(t, 0, env.Rsk.AccountNumber)
		assert.Equal(t, mockKeystoreSecret, env.Rsk.EncryptedJsonSecret)
		assert.Equal(t, mockPwdSecret, env.Rsk.EncryptedJsonPasswordSecret)
		assert.Equal(t, "regtest", env.Btc.Network)
	})

	t.Run("should convert base input to environment with env source", func(t *testing.T) {
		input := scripts.BaseInput{
			Network:                     "regtest",
			RskEndpoint:                 rskMockEndpoint,
			AwsLocalEndpoint:            awsMockEndpoint,
			SecretSource:                "env",
			EncryptedJsonSecret:         mockKeystoreSecret,
			EncryptedJsonPasswordSecret: mockPwdSecret,
			KeystoreFile:                mockFilePath,
		}
		env, err := input.ToEnv(func(fd int) ([]byte, error) {
			return []byte("test-pwd"), nil
		})
		require.NoError(t, err)
		assert.Equal(t, "regtest", env.LpsStage)
		assert.Equal(t, awsMockEndpoint, env.AwsLocalEndpoint)
		assert.Equal(t, "env", env.SecretSource)
		assert.Equal(t, "native", env.WalletManagement)
		assert.Equal(t, rskMockEndpoint, env.Rsk.Endpoint)
		assert.Equal(t, uint64(33), env.Rsk.ChainId)
		assert.Equal(t, "0x8901a2bbf639bfd21a97004ba4d7ae2bd00b8da8", env.Rsk.LbcAddress)
		assert.Equal(t, "0x0000000000000000000000000000000001000006", env.Rsk.BridgeAddress)
		assert.Equal(t, 0, env.Rsk.AccountNumber)
		assert.Equal(t, mockFilePath, env.Rsk.KeystoreFile)
		assert.Equal(t, "test-pwd", env.Rsk.KeystorePassword)
		assert.Equal(t, "regtest", env.Btc.Network)
	})
}
