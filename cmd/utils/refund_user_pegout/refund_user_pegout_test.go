package main

import (
	"flag"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/term"
)

func TestReadRefundUserPegOutScriptInput(t *testing.T) {
	t.Run("should set flag values", func(t *testing.T) {
		// Reset flags before test
		flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)

		scriptInput := new(RefundUserPegOutScriptInput)
		ReadRefundUserPegOutScriptInput(scriptInput)

		// Set test values
		err := flag.CommandLine.Parse([]string{
			"-network", "regtest",
			"-quote-hash", "d93f58c82100a6cee4f19ac505c11d51b52cafe220f7f1944b70496f33d277fc",
			"-rsk-endpoint", "http://localhost:4444",
			"-secret-src", "env",
			"-keystore-file", "./keystore.json",
		})
		require.NoError(t, err)

		assert.Equal(t, "regtest", scriptInput.Network)
		assert.Equal(t, "d93f58c82100a6cee4f19ac505c11d51b52cafe220f7f1944b70496f33d277fc", scriptInput.QuoteHashBytes)
		assert.Equal(t, "http://localhost:4444", scriptInput.RskEndpoint)
		assert.Equal(t, "env", scriptInput.SecretSource)
		assert.Equal(t, "./keystore.json", scriptInput.KeystoreFile)
	})
}

func TestParseRefundUserPegOutScriptInput(t *testing.T) {
	t.Run("should validate required fields", func(t *testing.T) {
		scriptInput := &RefundUserPegOutScriptInput{
			Network:        "",
			QuoteHashBytes: "",
			RskEndpoint:    "",
			SecretSource:   "",
		}

		_, err := ParseRefundUserPegOutScriptInput(scriptInput, term.ReadPassword)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid input")
	})

	t.Run("should parse valid input", func(t *testing.T) {
		scriptInput := &RefundUserPegOutScriptInput{
			Network:          "regtest",
			QuoteHashBytes:   "d93f58c82100a6cee4f19ac505c11d51b52cafe220f7f1944b70496f33d277fc",
			RskEndpoint:      "http://localhost:4444",
			SecretSource:     "aws",
			AwsLocalEndpoint: "http://localhost:4566",
		}

		env, err := ParseRefundUserPegOutScriptInput(scriptInput, func(fd int) ([]byte, error) {
			return []byte("password"), nil
		})
		require.NoError(t, err)
		assert.Equal(t, "regtest", env.LpsStage)
		assert.Equal(t, "http://localhost:4444", env.Rsk.Endpoint)
		assert.Equal(t, "aws", env.SecretSource)
		assert.Equal(t, "http://localhost:4566", env.AwsLocalEndpoint)
	})
}

func TestRefundUserPegOut(t *testing.T) {
	t.Run("should execute refund user peg out successfully", func(t *testing.T) {
		lbc := &mocks.LbcMock{}
		quoteHash := "d93f58c82100a6cee4f19ac505c11d51b52cafe220f7f1944b70496f33d277fc"
		expectedTxHash := test.AnyHash

		// Setup mock expectations
		lbc.On("RefundUserPegOut", quoteHash).Return(expectedTxHash, nil)

		txHash, err := ExecuteRefundUserPegOut(lbc, quoteHash)
		require.NoError(t, err)
		assert.Equal(t, expectedTxHash, txHash)

		// Verify all expectations were met
		lbc.AssertExpectations(t)
	})
}
