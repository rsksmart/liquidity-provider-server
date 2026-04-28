package main

import (
	"flag"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/cmd/utils/scripts"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadResignUtilsInput(t *testing.T) {
	t.Run("should set flag values", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
		input := new(ResignUtilsInput)
		ReadResignUtilsInput(input)

		err := flag.CommandLine.Parse([]string{
			"-network", "regtest",
			"-rsk-endpoint", "http://localhost:4444",
			"-secret-src", "env",
			"-keystore-file", "./keystore.json",
			"-resign",
		})
		require.NoError(t, err)

		assert.Equal(t, "regtest", input.Network)
		assert.Equal(t, "http://localhost:4444", input.RskEndpoint)
		assert.Equal(t, "env", input.SecretSource)
		assert.Equal(t, "./keystore.json", input.KeystoreFile)
		assert.True(t, input.Resign)
		assert.False(t, input.WithdrawCollateral)
	})
}

func TestParseResignUtilsInput(t *testing.T) {
	parse := func() {}

	t.Run("should validate required fields", func(t *testing.T) {
		input := &ResignUtilsInput{
			BaseInput: scripts.BaseInput{
				Network:      "",
				RskEndpoint:  "",
				SecretSource: "",
			},
			Resign: true,
		}
		err := ParseResignUtilsInput(parse, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid input")
	})

	t.Run("should reject missing action", func(t *testing.T) {
		input := &ResignUtilsInput{
			BaseInput: scripts.BaseInput{
				Network:          "regtest",
				RskEndpoint:      "http://localhost:4444",
				SecretSource:     "env",
				AwsLocalEndpoint: "http://localhost:4566",
			},
		}
		err := ParseResignUtilsInput(parse, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "select exactly one action")
	})

	t.Run("should reject multiple actions", func(t *testing.T) {
		input := &ResignUtilsInput{
			BaseInput: scripts.BaseInput{
				Network:          "regtest",
				RskEndpoint:      "http://localhost:4444",
				SecretSource:     "env",
				AwsLocalEndpoint: "http://localhost:4566",
			},
			Resign:             true,
			WithdrawCollateral: true,
		}
		err := ParseResignUtilsInput(parse, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "select exactly one action")
	})

	t.Run("should parse valid input", func(t *testing.T) {
		input := &ResignUtilsInput{
			BaseInput: scripts.BaseInput{
				Network:          "regtest",
				RskEndpoint:      "http://localhost:4444",
				SecretSource:     "env",
				AwsLocalEndpoint: "http://localhost:4566",
			},
			Resign: true,
		}
		err := ParseResignUtilsInput(parse, input)
		require.NoError(t, err)
	})
}

func TestExecuteResign(t *testing.T) {
	collateral := new(mocks.CollateralManagementContractMock)
	collateral.On("ProviderResign").Return(nil).Once()

	err := ExecuteResign(collateral)
	require.NoError(t, err)
	collateral.AssertExpectations(t)
}

func TestExecuteWithdrawCollateral(t *testing.T) {
	t.Run("should execute withdraw collateral", func(t *testing.T) {
		collateral := new(mocks.CollateralManagementContractMock)
		collateral.On("WithdrawCollateral").Return(nil).Once()

		err := ExecuteWithdrawCollateral(collateral)
		require.NoError(t, err)
		collateral.AssertExpectations(t)
	})

	t.Run("should return provider not resigned error", func(t *testing.T) {
		collateral := new(mocks.CollateralManagementContractMock)
		collateral.On("WithdrawCollateral").Return(liquidity_provider.ProviderNotResignedError).Once()

		err := ExecuteWithdrawCollateral(collateral)
		require.Error(t, err)
		require.ErrorIs(t, err, liquidity_provider.ProviderNotResignedError)
		collateral.AssertExpectations(t)
	})
}
