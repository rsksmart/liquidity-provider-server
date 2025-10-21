package main

import (
	"flag"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/term"

	"github.com/rsksmart/liquidity-provider-server/cmd/utils/scripts"
)

const (
	testRskEndpoint      = "http://localhost:4444"
	testQuoteHash        = "d93f58c82100a6cee4f19ac505c11d51b52cafe220f7f1944b70496f33d277fc"
	testAwsLocalEndpoint = "http://localhost:4566"
	testNetwork          = "regtest"
	testKeystoreFile     = "./keystore.json"
)

func TestReadRefundUserPegOutScriptInput(t *testing.T) {
	t.Run("should set flag values", func(t *testing.T) {
		// Reset flags before test
		flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)

		scriptInput := new(RefundUserPegOutScriptInput)
		ReadRefundUserPegOutScriptInput(scriptInput)

		// Set test values
		err := flag.CommandLine.Parse([]string{
			"-network", testNetwork,
			"-quote-hash", testQuoteHash,
			"-rsk-endpoint", testRskEndpoint,
			"-secret-src", "env",
			"-keystore-file", testKeystoreFile,
		})
		require.NoError(t, err)

		assert.Equal(t, testNetwork, scriptInput.Network)
		assert.Equal(t, testQuoteHash, scriptInput.QuoteHashBytes)
		assert.Equal(t, testRskEndpoint, scriptInput.RskEndpoint)
		assert.Equal(t, "env", scriptInput.SecretSource)
		assert.Equal(t, testKeystoreFile, scriptInput.KeystoreFile)
	})
}

func TestParseRefundUserPegOutScriptInput(t *testing.T) {

	parse := func() { // parse is a no-op function used as a placeholder in tests since the actual parsing
		// functionality is not relevant for these test cases
	}

	t.Run("should validate required fields", func(t *testing.T) {
		scriptInput := &RefundUserPegOutScriptInput{
			BaseInput: scripts.BaseInput{
				Network:      "",
				RskEndpoint:  "",
				SecretSource: "",
			},
			QuoteHashBytes: "",
		}

		_, err := ParseRefundUserPegOutScriptInput(parse, scriptInput, term.ReadPassword)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid input")
	})

	t.Run("should parse valid input", func(t *testing.T) {
		scriptInput := &RefundUserPegOutScriptInput{
			BaseInput: scripts.BaseInput{
				Network:          testNetwork,
				RskEndpoint:      testRskEndpoint,
				SecretSource:     "aws",
				AwsLocalEndpoint: testAwsLocalEndpoint,
			},
			QuoteHashBytes: testQuoteHash,
		}

		env, err := ParseRefundUserPegOutScriptInput(parse, scriptInput, func(fd int) ([]byte, error) {
			return []byte("password"), nil
		})
		require.NoError(t, err)
		assert.Equal(t, testNetwork, env.LpsStage)
		assert.Equal(t, testRskEndpoint, env.Rsk.Endpoint)
		assert.Equal(t, "aws", env.SecretSource)
		assert.Equal(t, testAwsLocalEndpoint, env.AwsLocalEndpoint)
	})
}

func TestRefundUserPegOut(t *testing.T) {
	t.Run("should execute refund user peg out successfully", func(t *testing.T) {
		pegoutContract := &mocks.PegoutContractMock{}
		expectedTxHash := test.AnyHash

		// Setup mock expectations
		pegoutContract.On("RefundUserPegOut", testQuoteHash).Return(expectedTxHash, nil)

		txHash, err := ExecuteRefundUserPegOut(pegoutContract, testQuoteHash)
		require.NoError(t, err)
		assert.Equal(t, expectedTxHash, txHash)

		// Verify all expectations were met
		pegoutContract.AssertExpectations(t)
	})
}
