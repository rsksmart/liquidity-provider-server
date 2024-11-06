package main

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/term"
	"testing"
)

func TestNewUpdateProviderArgs(t *testing.T) {
	t.Run("should fail on invalid url", func(t *testing.T) {
		invalidResult, err := NewUpdateProviderArgs("name", "1http://example.com", "regtest")
		require.Error(t, err)
		assert.Empty(t, invalidResult)

		validResult, err := NewUpdateProviderArgs("name", "https://example.com", "regtest")
		require.NoError(t, err)
		assert.NotEmpty(t, validResult)
	})
}

func TestUpdateProviderArgs_Validate(t *testing.T) {
	t.Run("should fail on empty name", func(t *testing.T) {
		args := UpdateProviderArgs{Name: "", url: nil, network: ""}
		err := args.Validate()
		require.Error(t, err)
	})
	t.Run("should fail on incomplete urls", func(t *testing.T) {
		urls := []string{"", "example.com", "https://"}
		for _, url := range urls {
			args, err := NewUpdateProviderArgs("name", url, "regtest")
			require.NoError(t, err)
			err = args.Validate()
			require.Error(t, err)
		}
	})
	t.Run("should enforce https in testnet or mainnet", func(t *testing.T) {
		testCases := []struct {
			url         string
			network     string
			expectError bool
		}{
			{url: "http://example11.com", network: "testnet", expectError: true},
			{url: "https://example22.com", network: "testnet", expectError: false},
			{url: "http://example33.com", network: "mainnet", expectError: true},
			{url: "https://example44.com", network: "mainnet", expectError: false},
		}
		for _, tc := range testCases {
			args, err := NewUpdateProviderArgs("name", tc.url, tc.network)
			require.NoError(t, err)
			err = args.Validate()
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}
	})
	t.Run("should not enforce https in regtest", func(t *testing.T) {
		args, err := NewUpdateProviderArgs("name", "http://example.com", "regtest")
		require.NoError(t, err)
		err = args.Validate()
		require.NoError(t, err)
	})
}

func TestUpdateProviderArgs_Url(t *testing.T) {
	testCases := []struct {
		url     string
		network string
		result  string
	}{
		{url: "http://example1.com/", network: "regtest", result: "http://example1.com"},
		{url: "https://example2.com/path", network: "testnet", result: "https://example2.com"},
		{url: "https://example3.com:1234", network: "mainnet", result: "https://example3.com:1234"},
		{url: "https://example4.com:1234/", network: "testnet", result: "https://example4.com:1234"},
		{url: "https://example5.com:1234/path", network: "testnet", result: "https://example5.com:1234"},
	}
	for _, tc := range testCases {
		args, err := NewUpdateProviderArgs("name", tc.url, tc.network)
		require.NoError(t, err)
		assert.Equal(t, tc.result, args.Url())
	}
}

func TestParseUpdateProviderScriptInput(t *testing.T) {
	scriptInput := new(UpdateProviderScriptInput)
	ReadUpdateProviderScriptInput(scriptInput)
	t.Run("should parse with aws secret source", func(t *testing.T) {
		require.NoError(t, flag.Set("network", "regtest"))
		require.NoError(t, flag.Set("provider-url", "https://example.com"))
		require.NoError(t, flag.Set("provider-name", "a name"))
		require.NoError(t, flag.Set("secret-src", "aws"))
		require.NoError(t, flag.Set("aws-endpoint", "http://localhost:1122"))
		require.NoError(t, flag.Set("rsk-endpoint", "http://localhost:3344"))
		require.NoError(t, flag.Set("lbc-address", "0xBEd51d83cc4676660E3Fc3819dfAD8238549B975"))
		require.NoError(t, flag.Set("keystore-secret", "UnitTest/Keystore-Secret"))
		require.NoError(t, flag.Set("password-secret", "UnitTest/Password-Secret"))
		env, err := ParseUpdateProviderScriptInput(scriptInput, term.ReadPassword)
		require.NoError(t, err)
		assert.Equal(t, "regtest", env.LpsStage)
		assert.Equal(t, "http://localhost:1122", env.AwsLocalEndpoint)
		assert.Equal(t, "aws", env.SecretSource)
		assert.Equal(t, "native", env.WalletManagement)
		assert.Equal(t, "http://localhost:3344", env.Rsk.Endpoint)
		assert.Equal(t, uint64(33), env.Rsk.ChainId)
		assert.Equal(t, "0xBEd51d83cc4676660E3Fc3819dfAD8238549B975", env.Rsk.LbcAddress)
		assert.Equal(t, "0x0000000000000000000000000000000001000006", env.Rsk.BridgeAddress)
		assert.Equal(t, 0, env.Rsk.AccountNumber)
		assert.Equal(t, "UnitTest/Keystore-Secret", env.Rsk.EncryptedJsonSecret)
		assert.Equal(t, "UnitTest/Password-Secret", env.Rsk.EncryptedJsonPasswordSecret)
		assert.Equal(t, "regtest", env.Btc.Network)

		assert.Equal(t, "a name", scriptInput.ProviderName)
		assert.Equal(t, "https://example.com", scriptInput.ProviderUrl)
	})

	t.Run("should parse with env secret source", func(t *testing.T) {
		require.NoError(t, flag.Set("network", "testnet"))
		require.NoError(t, flag.Set("provider-url", "https://example2.com"))
		require.NoError(t, flag.Set("provider-name", "a name 2"))
		require.NoError(t, flag.Set("secret-src", "env"))
		require.NoError(t, flag.Set("rsk-endpoint", "http://localhost:5566"))
		require.NoError(t, flag.Set("lbc-address", "0x64DCC3BcbEAE8CE586CabDeF79104986bEAFcAD6"))
		require.NoError(t, flag.Set("keystore-file", "path/to/a/file"))
		env, err := ParseUpdateProviderScriptInput(scriptInput, func(fd int) ([]byte, error) {
			return []byte("secret-password-123"), nil
		})
		require.NoError(t, err)
		assert.Equal(t, "testnet", env.LpsStage)
		assert.Equal(t, "env", env.SecretSource)
		assert.Equal(t, "native", env.WalletManagement)
		assert.Equal(t, "http://localhost:5566", env.Rsk.Endpoint)
		assert.Equal(t, uint64(31), env.Rsk.ChainId)
		assert.Equal(t, "0x64DCC3BcbEAE8CE586CabDeF79104986bEAFcAD6", env.Rsk.LbcAddress)
		assert.Equal(t, "0x0000000000000000000000000000000001000006", env.Rsk.BridgeAddress)
		assert.Equal(t, 0, env.Rsk.AccountNumber)
		assert.Equal(t, "path/to/a/file", env.Rsk.KeystoreFile)
		assert.Equal(t, "secret-password-123", env.Rsk.KeystorePassword)
		assert.Equal(t, "testnet", env.Btc.Network)

		assert.Equal(t, "a name 2", scriptInput.ProviderName)
		assert.Equal(t, "https://example2.com", scriptInput.ProviderUrl)
	})
}
