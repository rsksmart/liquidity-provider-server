package main

import (
	"flag"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/cmd/utils/scripts"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadWithdrawScriptInput(t *testing.T) {
	t.Run("should set flag values", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
		input := new(WithdrawScriptInput)
		ReadWithdrawScriptInput(input)

		err := flag.CommandLine.Parse([]string{
			"-network", "regtest",
			"-rsk-endpoint", "http://localhost:4444",
			"-secret-src", "env",
			"-keystore-file", "./keystore.json",
		})
		require.NoError(t, err)

		assert.Equal(t, "regtest", input.Network)
		assert.Equal(t, "http://localhost:4444", input.RskEndpoint)
		assert.Equal(t, "env", input.SecretSource)
		assert.Equal(t, "./keystore.json", input.KeystoreFile)
	})
}

func TestParseWithdrawScriptInput(t *testing.T) {
	parse := func() {}

	t.Run("should validate required fields", func(t *testing.T) {
		input := &WithdrawScriptInput{
			BaseInput: scripts.BaseInput{
				Network:      "",
				RskEndpoint:  "",
				SecretSource: "",
			},
		}
		err := ParseWithdrawScriptInput(parse, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid input")
	})

	t.Run("should parse valid input", func(t *testing.T) {
		input := &WithdrawScriptInput{
			BaseInput: scripts.BaseInput{
				Network:          "regtest",
				RskEndpoint:      "http://localhost:4444",
				SecretSource:     "env",
				AwsLocalEndpoint: "http://localhost:4566",
			},
		}
		err := ParseWithdrawScriptInput(parse, input)
		require.NoError(t, err)
	})
}

func TestExecuteWithdrawCollateral(t *testing.T) {
	lbc := new(mocks.LiquidityBridgeContractMock)
	lbc.On("WithdrawCollateral").Return(nil).Once()

	err := ExecuteWithdrawCollateral(lbc)
	require.NoError(t, err)
	lbc.AssertExpectations(t)
}
