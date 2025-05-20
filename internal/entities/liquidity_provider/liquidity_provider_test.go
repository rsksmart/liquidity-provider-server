package liquidity_provider_test

import (
	"encoding/hex"
	"errors"
	"testing"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/require"
)

func TestProviderType_IsValid(t *testing.T) {
	cases := test.Table[liquidity_provider.ProviderType, bool]{
		{Value: "pegin", Result: true},
		{Value: "pegout", Result: true},
		{Value: "both", Result: true},
		{Value: "", Result: false},
		{Value: test.AnyString, Result: false},
	}
	test.RunTable(t, cases, func(value liquidity_provider.ProviderType) bool {
		return value.IsValid()
	})
}

func TestProviderType_AcceptsPegin(t *testing.T) {
	cases := test.Table[liquidity_provider.ProviderType, bool]{
		{Value: "pegin", Result: true},
		{Value: "pegout", Result: false},
		{Value: "both", Result: true},
		{Value: "", Result: false},
		{Value: test.AnyString, Result: false},
	}
	test.RunTable(t, cases, func(value liquidity_provider.ProviderType) bool {
		return value.AcceptsPegin()
	})
}

func TestProviderType_AcceptsPegout(t *testing.T) {
	cases := test.Table[liquidity_provider.ProviderType, bool]{
		{Value: "pegin", Result: false},
		{Value: "pegout", Result: true},
		{Value: "both", Result: true},
		{Value: "", Result: false},
		{Value: test.AnyString, Result: false},
	}
	test.RunTable(t, cases, func(value liquidity_provider.ProviderType) bool {
		return value.AcceptsPegout()
	})
}

func TestToProviderType(t *testing.T) {
	type testResult struct {
		Result liquidity_provider.ProviderType
		Error  error
	}

	cases := test.Table[string, testResult]{
		{Value: "pegin", Result: testResult{Result: liquidity_provider.PeginProvider, Error: nil}},
		{Value: "pegout", Result: testResult{Result: liquidity_provider.PegoutProvider, Error: nil}},
		{Value: "both", Result: testResult{Result: liquidity_provider.FullProvider, Error: nil}},
		{Value: "", Result: testResult{Result: "", Error: liquidity_provider.InvalidProviderTypeError}},
		{Value: test.AnyString, Result: testResult{Result: "", Error: liquidity_provider.InvalidProviderTypeError}},
	}

	test.RunTable(t, cases, func(value string) testResult {
		result, err := liquidity_provider.ToProviderType(value)
		return testResult{
			Result: result,
			Error:  err,
		}
	})
}

//nolint:funlen
func TestValidateConfiguration(t *testing.T) {
	displayName := "test-config"

	t.Run("successful validation", func(t *testing.T) {
		mockSigner := new(mocks.SignerMock)

		mockConfig := liquidity_provider.GeneralConfiguration{
			RskConfirmations: liquidity_provider.ConfirmationsPerAmount{
				10: 100,
				20: 200,
			},
			BtcConfirmations: liquidity_provider.ConfirmationsPerAmount{
				10: 5,
				20: 10,
			},
			PublicLiquidityCheck: true,
		}
		mockConfigBytes := []byte(`{"rskConfirmations":{"10":100,"20":200},"btcConfirmations":{"10":5,"20":10},"publicLiquidityCheck":true}`)

		hash := ethcrypto.Keccak256(mockConfigBytes)
		hashHex := hex.EncodeToString(hash)

		mockSignature := []byte("mock-signature")
		mockSigner.On("SignBytes", hash).Return(mockSignature, nil)

		signature, err := mockSigner.SignBytes(hash)
		require.NoError(t, err)
		signatureHex := hex.EncodeToString(signature)

		signedConfig := &entities.Signed[liquidity_provider.GeneralConfiguration]{
			Value:     mockConfig,
			Hash:      hashHex,
			Signature: signatureHex,
		}

		mockSigner.On("Validate", signatureHex, hashHex).Return(true)

		readFunction := func() (*entities.Signed[liquidity_provider.GeneralConfiguration], error) {
			return signedConfig, nil
		}

		result, err := liquidity_provider.ValidateConfiguration(
			displayName,
			mockSigner,
			readFunction,
			ethcrypto.Keccak256,
		)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, mockConfig.RskConfirmations, result.Value.RskConfirmations)
		require.Equal(t, mockConfig.BtcConfirmations, result.Value.BtcConfirmations)
		require.Equal(t, mockConfig.PublicLiquidityCheck, result.Value.PublicLiquidityCheck)
		require.Equal(t, hashHex, result.Hash)
		require.Equal(t, signatureHex, result.Signature)

		mockSigner.AssertExpectations(t)
	})

	t.Run("error getting configuration", func(t *testing.T) {
		mockSigner := new(mocks.SignerMock)

		expectedErr := errors.New("database error")
		readFunction := func() (*entities.Signed[liquidity_provider.GeneralConfiguration], error) {
			return nil, expectedErr
		}

		defer test.AssertLogContains(t, "Error getting test-config configuration, using default configuration. Error: database error")()

		result, err := liquidity_provider.ValidateConfiguration(
			displayName,
			mockSigner,
			readFunction,
			ethcrypto.Keccak256,
		)

		require.Error(t, err)
		require.Equal(t, expectedErr, err)
		require.Nil(t, result)
		mockSigner.AssertExpectations(t)
	})

	t.Run("nil configuration", func(t *testing.T) {
		mockSigner := new(mocks.SignerMock)

		readFunction := func() (*entities.Signed[liquidity_provider.GeneralConfiguration], error) {
			return nil, nil
		}

		defer test.AssertLogContains(t, "Custom test-config configuration not found. Using default configuration.")()

		result, err := liquidity_provider.ValidateConfiguration(
			displayName,
			mockSigner,
			readFunction,
			ethcrypto.Keccak256,
		)

		require.Error(t, err)
		require.Equal(t, "configuration not found", err.Error())
		require.Nil(t, result)
		mockSigner.AssertExpectations(t)
	})

	t.Run("integrity check failure", func(t *testing.T) {
		mockSigner := new(mocks.SignerMock)

		mockConfig := liquidity_provider.GeneralConfiguration{
			RskConfirmations: liquidity_provider.ConfirmationsPerAmount{
				10: 100,
				20: 200,
			},
			BtcConfirmations: liquidity_provider.ConfirmationsPerAmount{
				10: 5,
				20: 10,
			},
			PublicLiquidityCheck: true,
		}
		mockConfigBytes := []byte(`{"rskConfirmations":{"10":100,"20":200},"btcConfirmations":{"10":5,"20":10},"publicLiquidityCheck":true}`)

		hash := ethcrypto.Keccak256(mockConfigBytes)
		hashHex := hex.EncodeToString(hash)

		mockSignature := []byte("mock-signature")
		mockSigner.On("SignBytes", hash).Return(mockSignature, nil)

		signature, err := mockSigner.SignBytes(hash)
		require.NoError(t, err)
		signatureHex := hex.EncodeToString(signature)

		signedConfig := &entities.Signed[liquidity_provider.GeneralConfiguration]{
			Value:     mockConfig,
			Hash:      hashHex,
			Signature: signatureHex,
		}

		mockSigner.On("Validate", signatureHex, hashHex).Return(false)

		readFunction := func() (*entities.Signed[liquidity_provider.GeneralConfiguration], error) {
			return signedConfig, nil
		}

		defer test.AssertLogContains(t, "Invalid test-config configuration signature. Using default configuration.")()

		result, err := liquidity_provider.ValidateConfiguration(
			displayName,
			mockSigner,
			readFunction,
			ethcrypto.Keccak256,
		)

		require.Error(t, err)
		require.Equal(t, "invalid signature", err.Error())
		require.Nil(t, result)
		mockSigner.AssertExpectations(t)
	})
}
