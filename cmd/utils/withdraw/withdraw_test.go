package main

import (
	"flag"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/cmd/utils/scripts"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
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

func TestParseWithdrawScriptInput_ValidateRequiredFields(t *testing.T) {
	parse := func() {}
	input := &WithdrawScriptInput{
		BaseInput: scripts.BaseInput{
			Network:      "",
			RskEndpoint:  "",
			SecretSource: "",
		},
		All: true, // Add flag to pass flag validation and reach struct validation
	}
	err := ParseWithdrawScriptInput(parse, input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid input")
}

func TestParseWithdrawScriptInput_AllFlag(t *testing.T) {
	parse := func() {}
	input := &WithdrawScriptInput{
		BaseInput: scripts.BaseInput{
			Network:          "regtest",
			RskEndpoint:      "http://localhost:4444",
			SecretSource:     "env",
			AwsLocalEndpoint: "http://localhost:4566",
		},
		All: true,
	}
	err := ParseWithdrawScriptInput(parse, input)
	require.NoError(t, err)
}

func TestParseWithdrawScriptInput_AmountFlag(t *testing.T) {
	parse := func() {}
	input := &WithdrawScriptInput{
		BaseInput: scripts.BaseInput{
			Network:          "regtest",
			RskEndpoint:      "http://localhost:4444",
			SecretSource:     "env",
			AwsLocalEndpoint: "http://localhost:4566",
		},
		Amount: "1000000000000000000",
	}
	err := ParseWithdrawScriptInput(parse, input)
	require.NoError(t, err)
}

func TestParseWithdrawScriptInput_BothFlags(t *testing.T) {
	parse := func() {}
	input := &WithdrawScriptInput{
		BaseInput: scripts.BaseInput{
			Network:          "regtest",
			RskEndpoint:      "http://localhost:4444",
			SecretSource:     "env",
			AwsLocalEndpoint: "http://localhost:4566",
		},
		All:    true,
		Amount: "1000000000000000000",
	}
	err := ParseWithdrawScriptInput(parse, input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot use both --all and --amount flags")
}

func TestParseWithdrawScriptInput_NoFlags(t *testing.T) {
	parse := func() {}
	input := &WithdrawScriptInput{
		BaseInput: scripts.BaseInput{
			Network:          "regtest",
			RskEndpoint:      "http://localhost:4444",
			SecretSource:     "env",
			AwsLocalEndpoint: "http://localhost:4566",
		},
	}
	err := ParseWithdrawScriptInput(parse, input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must provide either --all or --amount flag")
}

func TestExecuteWithdraw(t *testing.T) {
	address := "0x1234567890123456789012345678901234567890"

	t.Run("should withdraw all funds when --all flag is used", func(t *testing.T) {
		lbc := new(mocks.LiquidityBridgeContractMock)
		balance := entities.NewWei(1000000000000000000)
		lbc.On("GetBalance", address).Return(balance, nil).Once()
		lbc.On("Withdraw", balance).Return(nil).Once()

		err := ExecuteWithdraw(lbc, address, true, "")
		require.NoError(t, err)
		lbc.AssertExpectations(t)
	})

	t.Run("should withdraw specific amount when --amount flag is used", func(t *testing.T) {
		lbc := new(mocks.LiquidityBridgeContractMock)
		amount := entities.NewWei(500000000000000000)
		lbc.On("Withdraw", amount).Return(nil).Once()

		err := ExecuteWithdraw(lbc, address, false, "500000000000000000")
		require.NoError(t, err)
		lbc.AssertExpectations(t)
	})

	t.Run("should fail when GetBalance fails", func(t *testing.T) {
		lbc := new(mocks.LiquidityBridgeContractMock)
		lbc.On("GetBalance", address).Return(nil, assert.AnError).Once()

		err := ExecuteWithdraw(lbc, address, true, "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get balance")
		lbc.AssertExpectations(t)
	})

	t.Run("should fail when amount is invalid", func(t *testing.T) {
		lbc := new(mocks.LiquidityBridgeContractMock)

		err := ExecuteWithdraw(lbc, address, false, "invalid")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid amount format")
	})

	t.Run("should fail when Withdraw fails", func(t *testing.T) {
		lbc := new(mocks.LiquidityBridgeContractMock)
		amount := entities.NewWei(500000000000000000)
		lbc.On("Withdraw", amount).Return(assert.AnError).Once()

		err := ExecuteWithdraw(lbc, address, false, "500000000000000000")
		require.Error(t, err)
		lbc.AssertExpectations(t)
	})
}
