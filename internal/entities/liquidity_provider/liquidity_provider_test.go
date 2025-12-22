package liquidity_provider_test

import (
	"encoding/hex"
	"errors"
	"github.com/stretchr/testify/assert"
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
		{Value: liquidity_provider.PeginProvider, Result: true},
		{Value: liquidity_provider.PegoutProvider, Result: true},
		{Value: liquidity_provider.FullProvider, Result: true},
		{Value: -1, Result: false},
		{Value: 5, Result: false},
	}
	test.RunTable(t, cases, func(value liquidity_provider.ProviderType) bool {
		return value.IsValid()
	})
}

func TestProviderType_AcceptsPegin(t *testing.T) {
	cases := test.Table[liquidity_provider.ProviderType, bool]{
		{Value: liquidity_provider.PeginProvider, Result: true},
		{Value: liquidity_provider.PegoutProvider, Result: false},
		{Value: liquidity_provider.FullProvider, Result: true},
		{Value: -1, Result: false},
		{Value: 5, Result: false},
	}
	test.RunTable(t, cases, func(value liquidity_provider.ProviderType) bool {
		return value.AcceptsPegin()
	})
}

func TestProviderType_AcceptsPegout(t *testing.T) {
	cases := test.Table[liquidity_provider.ProviderType, bool]{
		{Value: liquidity_provider.PeginProvider, Result: false},
		{Value: liquidity_provider.PegoutProvider, Result: true},
		{Value: liquidity_provider.FullProvider, Result: true},
		{Value: -1, Result: false},
		{Value: 5, Result: false},
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

	cases := test.Table[int, testResult]{
		{Value: 0, Result: testResult{Result: liquidity_provider.PeginProvider, Error: nil}},
		{Value: 1, Result: testResult{Result: liquidity_provider.PegoutProvider, Error: nil}},
		{Value: 2, Result: testResult{Result: liquidity_provider.FullProvider, Error: nil}},
		{Value: -1, Result: testResult{Result: -1, Error: liquidity_provider.InvalidProviderTypeError}},
		{Value: 5, Result: testResult{Result: -1, Error: liquidity_provider.InvalidProviderTypeError}},
	}

	test.RunTable(t, cases, func(value int) testResult {
		result, err := liquidity_provider.ToProviderType(value)
		return testResult{
			Result: result,
			Error:  err,
		}
	})
}

//nolint:funlen
func TestValidateConfiguration(t *testing.T) {
	t.Run("successful validation", func(t *testing.T) {
		mockSigner := new(mocks.SignerMock)

		mockConfig := liquidity_provider.GeneralConfiguration{
			RskConfirmations: liquidity_provider.ConfirmationsPerAmount{
				"10": 100,
				"20": 200,
			},
			BtcConfirmations: liquidity_provider.ConfirmationsPerAmount{
				"10": 5,
				"20": 10,
			},
			PublicLiquidityCheck:      true,
			MaxLiquidity:              entities.NewWei(1000000),
			ReimbursementWindowBlocks: 100,
		}
		mockConfigBytes := []byte(`{"rskConfirmations":{"10":100,"20":200},"btcConfirmations":{"10":5,"20":10},"publicLiquidityCheck":true,"maxLiquidity":1000000,"reimbursementWindowBlocks":100}`)

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
			mockSigner,
			ethcrypto.Keccak256,
			readFunction,
		)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, mockConfig.RskConfirmations, result.Value.RskConfirmations)
		require.Equal(t, mockConfig.BtcConfirmations, result.Value.BtcConfirmations)
		require.Equal(t, mockConfig.PublicLiquidityCheck, result.Value.PublicLiquidityCheck)
		require.Equal(t, mockConfig.MaxLiquidity, result.Value.MaxLiquidity)
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

		result, err := liquidity_provider.ValidateConfiguration(
			mockSigner,
			ethcrypto.Keccak256,
			readFunction,
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

		result, err := liquidity_provider.ValidateConfiguration(
			mockSigner,
			ethcrypto.Keccak256,
			readFunction,
		)

		require.ErrorIs(t, err, liquidity_provider.ConfigurationNotFoundError)
		require.Nil(t, result)
		mockSigner.AssertExpectations(t)
	})

	t.Run("integrity check failure", func(t *testing.T) {
		mockSigner := new(mocks.SignerMock)

		mockConfig := liquidity_provider.GeneralConfiguration{
			RskConfirmations: liquidity_provider.ConfirmationsPerAmount{
				"10": 100,
				"20": 200,
			},
			BtcConfirmations: liquidity_provider.ConfirmationsPerAmount{
				"10": 5,
				"20": 10,
			},
			PublicLiquidityCheck:      true,
			MaxLiquidity:              entities.NewWei(1000000),
			ReimbursementWindowBlocks: 100,
		}
		mockConfigBytes := []byte(`{"rskConfirmations":{"10":100,"20":200},"btcConfirmations":{"10":5,"20":10},"publicLiquidityCheck":true,"maxLiquidity":1000000,"reimbursementWindowBlocks":100}`)

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

		result, err := liquidity_provider.ValidateConfiguration(
			mockSigner,
			ethcrypto.Keccak256,
			readFunction,
		)
		require.ErrorIs(t, err, liquidity_provider.InvalidSignatureError)
		require.Nil(t, result)
		mockSigner.AssertExpectations(t)
	})
}

func TestProviderType_Name(t *testing.T) {
	tests := []struct {
		p    liquidity_provider.ProviderType
		want string
	}{
		{p: liquidity_provider.PeginProvider, want: "pegin"},
		{p: liquidity_provider.PegoutProvider, want: "pegout"},
		{p: liquidity_provider.FullProvider, want: "both"},
		{p: -5, want: "unknown"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.p.Name())
	}
}
