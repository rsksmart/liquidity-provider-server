package rootstock_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	penalizedIteratorString = "*bindings.LiquidityBridgeContractPenalizedIterator"
	depositIteratorString   = "*bindings.LiquidityBridgeContractPegOutDepositIterator"
	invalidAddressTest      = "Invalid address"
)

var peginQuote = quote.PeginQuote{
	FedBtcAddress:      "2MzQwSSnBHWHqSAqtTVQ6v47XtaisrJa1Vc",
	LbcAddress:         "0xd5f00ABfbEA7A0B193836CAc6833c2Ad9D06cEa8",
	LpRskAddress:       "0x892813507Bf3aBF2890759d2135Ec34f4909Fea5",
	BtcRefundAddress:   "mmExR8uifX9SdfZ7V9st7pSdUhf1rmaPGk",
	RskRefundAddress:   "0x5dE07e2BE63595854C396E2da291e0d1EdE15112",
	LpBtcAddress:       "mipcBbFg9gMiCh81Kj8tqqdgoZub1ZJRfn",
	CallFee:            entities.NewWei(150),
	PenaltyFee:         entities.NewWei(99),
	ContractAddress:    "0x0D8Fb5d32704DB2931e05DB91F64BcA6f76Ce573",
	Data:               "0x12a1",
	GasLimit:           8000,
	Nonce:              11223344,
	Value:              entities.NewWei(1234),
	AgreementTimestamp: 20,
	TimeForDeposit:     30,
	LpCallTime:         40,
	Confirmations:      50,
	CallOnRegister:     true,
	GasFee:             entities.NewWei(100),
	ProductFeeAmount:   500,
}

var parsedPeginQuote = bindings.QuotesPeginQuote{
	FedBtcAddress:               [20]byte{78, 159, 57, 202, 70, 136, 255, 16, 33, 40, 234, 76, 205, 163, 65, 5, 50, 67, 5, 176},
	LbcAddress:                  common.Address{0xd5, 0xf0, 0x0A, 0xBf, 0xbE, 0xA7, 0xA0, 0xB1, 0x93, 0x83, 0x6C, 0xAc, 0x68, 0x33, 0xc2, 0xAd, 0x9D, 0x06, 0xcE, 0xa8},
	LiquidityProviderRskAddress: common.Address{0x89, 0x28, 0x13, 0x50, 0x7B, 0xf3, 0xaB, 0xF2, 0x89, 0x07, 0x59, 0xd2, 0x13, 0x5E, 0xc3, 0x4f, 0x49, 0x09, 0xFe, 0xa5},
	BtcRefundAddress:            []byte{111, 62, 202, 51, 133, 156, 181, 52, 157, 247, 31, 35, 0, 74, 185, 49, 162, 115, 243, 129, 220},
	RskRefundAddress:            common.Address{0x5d, 0xE0, 0x7e, 0x2B, 0xE6, 0x35, 0x95, 0x85, 0x4C, 0x39, 0x6E, 0x2d, 0xa2, 0x91, 0xe0, 0xd1, 0xEd, 0xE1, 0x51, 0x12},
	LiquidityProviderBtcAddress: []byte{111, 36, 63, 19, 148, 244, 69, 84, 244, 206, 63, 214, 134, 73, 193, 154, 220, 72, 60, 233, 36},
	CallFee:                     big.NewInt(150),
	PenaltyFee:                  big.NewInt(99),
	ContractAddress:             common.Address{0x0D, 0x8F, 0xb5, 0xd3, 0x27, 0x04, 0xDB, 0x29, 0x31, 0xe0, 0x5D, 0xB9, 0x1F, 0x64, 0xBc, 0xA6, 0xf7, 0x6C, 0xe5, 0x73},
	Data:                        []byte{18, 161},
	GasLimit:                    8000,
	Nonce:                       11223344,
	Value:                       big.NewInt(1234),
	AgreementTimestamp:          20,
	TimeForDeposit:              30,
	CallTime:                    40,
	DepositConfirmations:        50,
	CallOnRegister:              true,
	ProductFeeAmount:            big.NewInt(500),
	GasFee:                      big.NewInt(100),
}

var pegoutQuote = quote.PegoutQuote{
	LbcAddress:            "0x7C4890A0f1D4bBf2C669Ac2d1efFa185c505359b",
	LpRskAddress:          "0x9D93929A9099be4355fC2389FbF253982F9dF47c",
	BtcRefundAddress:      "mjaGtyj74LYn7gApr17prZxDPDnfuUnRa5",
	RskRefundAddress:      "0x0d8A0F1ef26B4b9650d98E1c22c560327cF387FE",
	LpBtcAddress:          "mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6",
	CallFee:               entities.NewWei(1),
	PenaltyFee:            2,
	Nonce:                 3,
	DepositAddress:        "mzpvGLbwqFK7zvtFanb3r83WpwDxNQUYcz",
	Value:                 entities.NewWei(11),
	AgreementTimestamp:    10,
	DepositDateLimit:      5,
	DepositConfirmations:  50,
	TransferConfirmations: 40,
	TransferTime:          800,
	ExpireDate:            700,
	ExpireBlock:           600,
	GasFee:                entities.NewWei(888),
	ProductFeeAmount:      789,
}

var parsedPegoutQuote = bindings.QuotesPegOutQuote{
	LbcAddress:            common.Address{0x7C, 0x48, 0x90, 0xA0, 0xf1, 0xD4, 0xbB, 0xf2, 0xC6, 0x69, 0xAc, 0x2d, 0x1e, 0xfF, 0xa1, 0x85, 0xc5, 0x05, 0x35, 0x9b},
	LpRskAddress:          common.Address{0x9D, 0x93, 0x92, 0x9A, 0x90, 0x99, 0xbe, 0x43, 0x55, 0xfC, 0x23, 0x89, 0xFb, 0xF2, 0x53, 0x98, 0x2F, 0x9d, 0xF4, 0x7c},
	BtcRefundAddress:      []byte{111, 44, 129, 71, 129, 50, 181, 221, 166, 79, 252, 72, 74, 13, 34, 80, 150, 196, 178, 42, 213},
	RskRefundAddress:      common.Address{0x0d, 0x8A, 0x0F, 0x1e, 0xf2, 0x6B, 0x4b, 0x96, 0x50, 0xd9, 0x8E, 0x1c, 0x22, 0xc5, 0x60, 0x32, 0x7c, 0xF3, 0x87, 0xFE},
	LpBtcAddress:          []byte{111, 77, 25, 29, 213, 89, 98, 105, 251, 153, 254, 72, 87, 69, 181, 236, 217, 126, 72, 185, 152},
	CallFee:               big.NewInt(1),
	PenaltyFee:            big.NewInt(2),
	Nonce:                 3,
	DeposityAddress:       []byte{111, 211, 208, 46, 217, 53, 0, 241, 222, 67, 83, 58, 58, 11, 192, 246, 63, 2, 142, 48, 246},
	Value:                 big.NewInt(11),
	AgreementTimestamp:    10,
	DepositDateLimit:      5,
	DepositConfirmations:  50,
	TransferConfirmations: 40,
	TransferTime:          800,
	ExpireDate:            700,
	ExpireBlock:           600,
	ProductFeeAmount:      big.NewInt(789),
	GasFee:                big.NewInt(888),
}

var contractProviders = []bindings.LiquidityBridgeContractLiquidityProvider{
	{
		Id:           big.NewInt(1),
		Name:         "test",
		ApiBaseUrl:   "http://test.com",
		Status:       true,
		Provider:     common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
		ProviderType: "pegin",
	},
	{
		Id:           big.NewInt(2),
		Name:         "test2",
		ApiBaseUrl:   "http://test2.com",
		Status:       true,
		Provider:     common.HexToAddress("0x1233557799abcdef1234567890abcdef12345678"),
		ProviderType: "pegout",
	},
}
var parsedProviders = []liquidity_provider.RegisteredLiquidityProvider{
	{
		Id:           1,
		Address:      "0x1234567890AbcdEF1234567890aBcdef12345678",
		Name:         "test",
		ApiBaseUrl:   "http://test.com",
		Status:       true,
		ProviderType: "pegin",
	},
	{
		Id:           2,
		Address:      "0x1233557799ABcDEf1234567890abCdef12345678",
		Name:         "test2",
		ApiBaseUrl:   "http://test2.com",
		Status:       true,
		ProviderType: "pegout",
	},
}

var penalizations = []*bindings.LiquidityBridgeContractPenalized{
	{QuoteHash: [32]byte{1, 2, 3}, LiquidityProvider: common.Address{1}, Penalty: big.NewInt(555)},
	{QuoteHash: [32]byte{4, 5, 6}, LiquidityProvider: common.Address{2}, Penalty: big.NewInt(666)},
	{QuoteHash: [32]byte{7, 8, 9}, LiquidityProvider: common.Address{3}, Penalty: big.NewInt(777)},
}

var parsedPenalizations = []penalization.PenalizedEvent{
	{
		QuoteHash:         "0102030000000000000000000000000000000000000000000000000000000000",
		Penalty:           entities.NewWei(555),
		LiquidityProvider: "0x0100000000000000000000000000000000000000",
	},
	{
		QuoteHash:         "0405060000000000000000000000000000000000000000000000000000000000",
		Penalty:           entities.NewWei(666),
		LiquidityProvider: "0x0200000000000000000000000000000000000000",
	},
	{
		QuoteHash:         "0708090000000000000000000000000000000000000000000000000000000000",
		Penalty:           entities.NewWei(777),
		LiquidityProvider: "0x0300000000000000000000000000000000000000",
	},
}

var deposits = []*bindings.LiquidityBridgeContractPegOutDeposit{
	{Raw: geth.Log{TxHash: common.Hash{7}, BlockNumber: 10}, QuoteHash: [32]byte{1, 2, 3}, Sender: common.Address{1}, Amount: big.NewInt(555), Timestamp: big.NewInt(123456789)},
	{Raw: geth.Log{TxHash: common.Hash{8}, BlockNumber: 11}, QuoteHash: [32]byte{4, 5, 6}, Sender: common.Address{2}, Amount: big.NewInt(666), Timestamp: big.NewInt(987654321)},
	{Raw: geth.Log{TxHash: common.Hash{9}, BlockNumber: 12}, QuoteHash: [32]byte{7, 8, 9}, Sender: common.Address{3}, Amount: big.NewInt(777), Timestamp: big.NewInt(111222333)},
}

var parsedDeposits = []quote.PegoutDeposit{
	{
		TxHash:      "0x0700000000000000000000000000000000000000000000000000000000000000",
		QuoteHash:   "0102030000000000000000000000000000000000000000000000000000000000",
		Amount:      entities.NewWei(555),
		Timestamp:   time.Unix(123456789, 0),
		BlockNumber: 10,
		From:        "0x0100000000000000000000000000000000000000",
	},
	{
		TxHash:      "0x0800000000000000000000000000000000000000000000000000000000000000",
		QuoteHash:   "0405060000000000000000000000000000000000000000000000000000000000",
		Amount:      entities.NewWei(666),
		Timestamp:   time.Unix(987654321, 0),
		BlockNumber: 11,
		From:        "0x0200000000000000000000000000000000000000",
	},
	{
		TxHash:      "0x0900000000000000000000000000000000000000000000000000000000000000",
		QuoteHash:   "0708090000000000000000000000000000000000000000000000000000000000",
		Amount:      entities.NewWei(777),
		Timestamp:   time.Unix(111222333, 0),
		BlockNumber: 12,
		From:        "0x0300000000000000000000000000000000000000",
	},
}

var parsedAddress = common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")

func TestNewLiquidityBridgeContractImpl(t *testing.T) {
	lbc := rootstock.NewLiquidityBridgeContractImpl(
		rootstock.NewRskClient(&mocks.RpcClientBindingMock{}),
		test.AnyAddress, &mocks.LbcAdapterMock{},
		&mocks.TransactionSignerMock{},
		rootstock.RetryParams{Retries: 1, Sleep: time.Duration(1)},
		time.Duration(1),
	)
	test.AssertNonZeroValues(t, lbc)
}

func TestLiquidityBridgeContractImpl_GetAddress(t *testing.T) {
	lbc := rootstock.NewLiquidityBridgeContractImpl(dummyClient, test.AnyAddress, nil, nil, rootstock.RetryParams{}, time.Duration(1))
	assert.Equal(t, test.AnyAddress, lbc.GetAddress())
}

func TestLiquidityBridgeContractImpl_HashPeginQuote(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	hash := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	lbc := rootstock.NewLiquidityBridgeContractImpl(
		dummyClient,
		test.AnyAddress,
		lbcMock,
		nil, rootstock.RetryParams{},
		time.Duration(1),
	)
	t.Run("Success", func(t *testing.T) {
		lbcMock.On("HashQuote", mock.Anything, parsedPeginQuote).Return(hash, nil).Once()
		result, err := lbc.HashPeginQuote(peginQuote)
		require.NoError(t, err)
		assert.Equal(t, hex.EncodeToString(hash[:]), result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling on HashQuote call fail", func(t *testing.T) {
		lbcMock.On("HashQuote", mock.Anything, parsedPeginQuote).Return(nil, assert.AnError).Once()
		result, err := lbc.HashPeginQuote(peginQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestLiquidityBridgeContractImpl_HashPeginQuote_ParsingErrors(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(dummyClient, test.AnyAddress, lbcMock, nil, rootstock.RetryParams{}, time.Duration(1))
	validationFunction := func(peginQuote quote.PeginQuote) {
		result, err := lbc.HashPeginQuote(peginQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	}
	t.Run("Incomplete quote", func(t *testing.T) {
		testQuote := peginQuote
		testQuote.LbcAddress = ""
		validationFunction(testQuote)
	})
	t.Run("Invalid federation address", func(t *testing.T) {
		testQuote := peginQuote
		testQuote.FedBtcAddress = test.AnyString
		validationFunction(testQuote)
	})
	t.Run("Invalid lp btc address", func(t *testing.T) {
		testQuote := peginQuote
		testQuote.LpBtcAddress = test.AnyString
		validationFunction(testQuote)
	})
	t.Run("Invalid btc refund address", func(t *testing.T) {
		testQuote := peginQuote
		testQuote.BtcRefundAddress = test.AnyString
		validationFunction(testQuote)
	})
	t.Run("Invalid lbc address", func(t *testing.T) {
		testQuote := peginQuote
		testQuote.LbcAddress = test.AnyString
		validationFunction(testQuote)
	})
	t.Run("Invalid rsk refund address", func(t *testing.T) {
		testQuote := peginQuote
		testQuote.RskRefundAddress = test.AnyString
		validationFunction(testQuote)
	})
	t.Run("Invalid lp rsk address", func(t *testing.T) {
		testQuote := peginQuote
		testQuote.LpRskAddress = test.AnyString
		validationFunction(testQuote)
	})
	t.Run("Invalid destination address", func(t *testing.T) {
		testQuote := peginQuote
		testQuote.ContractAddress = test.AnyString
		validationFunction(testQuote)
	})
	t.Run("Invalid data", func(t *testing.T) {
		testQuote := peginQuote
		testQuote.Data = test.AnyString
		validationFunction(testQuote)
	})
}

func TestLiquidityBridgeContractImpl_HashPegoutQuote(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	hash := [32]byte{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	lbc := rootstock.NewLiquidityBridgeContractImpl(
		dummyClient,
		test.AnyAddress,
		lbcMock,
		nil, rootstock.RetryParams{},
		time.Duration(1),
	)
	t.Run("Success", func(t *testing.T) {
		lbcMock.On("HashPegoutQuote", mock.Anything, parsedPegoutQuote).Return(hash, nil).Once()
		result, err := lbc.HashPegoutQuote(pegoutQuote)
		require.NoError(t, err)
		assert.Equal(t, hex.EncodeToString(hash[:]), result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling on HashPegoutQuote fail", func(t *testing.T) {
		lbcMock.On("HashPegoutQuote", mock.Anything, parsedPegoutQuote).Return(nil, assert.AnError).Once()
		result, err := lbc.HashPegoutQuote(pegoutQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestLiquidityBridgeContractImpl_HashPegoutQuote_ParsingErrors(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(dummyClient, test.AnyAddress, lbcMock, nil, rootstock.RetryParams{}, time.Duration(1))
	t.Run("Incomplete quote", func(t *testing.T) {
		testQuote := pegoutQuote
		testQuote.LbcAddress = ""
		result, err := lbc.HashPegoutQuote(testQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
	t.Run("Invalid lbc address", func(t *testing.T) {
		testQuote := pegoutQuote
		testQuote.LbcAddress = test.AnyString
		result, err := lbc.HashPegoutQuote(testQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
	t.Run("Invalid lp rsk address", func(t *testing.T) {
		testQuote := pegoutQuote
		testQuote.LpRskAddress = test.AnyString
		result, err := lbc.HashPegoutQuote(testQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
	t.Run("Invalid rsk refund address", func(t *testing.T) {
		testQuote := pegoutQuote
		testQuote.RskRefundAddress = test.AnyString
		result, err := lbc.HashPegoutQuote(testQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
	t.Run("Invalid btc refund address", func(t *testing.T) {
		testQuote := pegoutQuote
		testQuote.BtcRefundAddress = test.AnyString
		result, err := lbc.HashPegoutQuote(testQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
	t.Run("Invalid lp btc address", func(t *testing.T) {
		testQuote := pegoutQuote
		testQuote.LpBtcAddress = test.AnyString
		result, err := lbc.HashPegoutQuote(testQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
	t.Run("Invalid destination address", func(t *testing.T) {
		testQuote := pegoutQuote
		testQuote.DepositAddress = test.AnyString
		result, err := lbc.HashPegoutQuote(testQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestLiquidityBridgeContractImpl_GetProviders(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(dummyClient, test.AnyAddress, lbcMock, nil, rootstock.RetryParams{}, time.Duration(1))
	t.Run("Success", func(t *testing.T) {
		lbcMock.On("GetProviders", mock.Anything).Return(contractProviders, nil).Once()
		result, err := lbc.GetProviders()
		require.NoError(t, err)
		assert.Equal(t, parsedProviders, result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling on invalid provider type", func(t *testing.T) {
		invalidProviders := contractProviders
		invalidProviders[0].ProviderType = "invalid type"
		lbcMock.On("GetProviders", mock.Anything).Return(invalidProviders, nil).Once()
		result, err := lbc.GetProviders()
		require.ErrorIs(t, err, liquidity_provider.InvalidProviderTypeError)
		assert.Nil(t, result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling GetProviders", func(t *testing.T) {
		lbcMock.On("GetProviders", mock.Anything).Return(nil, assert.AnError).Once()
		result, err := lbc.GetProviders()
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestLiquidityBridgeContractImpl_ProviderResign(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		lbcMock,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
	)
	t.Run("Success", func(t *testing.T) {
		tx, _ := prepareTxMocks(mockClient, signerMock, true)
		lbcMock.On("Resign", mock.Anything).Return(tx, nil).Once()
		err := lbc.ProviderResign()
		require.NoError(t, err)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling when sending resign tx", func(t *testing.T) {
		_, _ = prepareTxMocks(mockClient, signerMock, true)
		lbcMock.On("Resign", mock.Anything).Return(nil, assert.AnError).Once()
		err := lbc.ProviderResign()
		require.Error(t, err)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling (resign tx reverted)", func(t *testing.T) {
		tx, _ := prepareTxMocks(mockClient, signerMock, false)
		lbcMock.On("Resign", mock.Anything).Return(tx, nil).Once()
		err := lbc.ProviderResign()
		require.ErrorContains(t, err, "resign transaction failed")
		lbcMock.AssertExpectations(t)
	})
}

func TestLiquidityBridgeContractImpl_SetProviderStatus(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		lbcMock,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
	)
	t.Run("Success", func(t *testing.T) {
		tx, _ := prepareTxMocks(mockClient, signerMock, true)
		lbcMock.On("SetProviderStatus", mock.Anything, big.NewInt(2), true).Return(tx, nil).Once()
		err := lbc.SetProviderStatus(2, true)
		require.NoError(t, err)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling when sending setProviderStatus tx", func(t *testing.T) {
		_, _ = prepareTxMocks(mockClient, signerMock, true)
		lbcMock.On("SetProviderStatus", mock.Anything, big.NewInt(1), true).Return(nil, assert.AnError).Once()
		err := lbc.SetProviderStatus(1, true)
		require.Error(t, err)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling (setProviderStatus tx reverted)", func(t *testing.T) {
		tx, _ := prepareTxMocks(mockClient, signerMock, false)
		lbcMock.On("SetProviderStatus", mock.Anything, big.NewInt(1), false).Return(tx, nil).Once()
		err := lbc.SetProviderStatus(1, false)
		require.ErrorContains(t, err, "setProviderStatus transaction failed")
		lbcMock.AssertExpectations(t)
	})
}

func TestLiquidityBridgeContractImpl_UpdateProvider(t *testing.T) {
	const (
		name = "test name"
		url  = "http://test.update.example.com"
	)

	lbcMock := &mocks.LbcAdapterMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		lbcMock,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
	)
	t.Run("Success", func(t *testing.T) {
		tx, _ := prepareTxMocks(mockClient, signerMock, true)
		lbcMock.EXPECT().UpdateProvider(mock.Anything, name, url).Return(tx, nil).Once()
		result, err := lbc.UpdateProvider(name, url)
		require.NoError(t, err)
		assert.Equal(t, tx.Hash().String(), result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling when sending updateProvider tx", func(t *testing.T) {
		_, _ = prepareTxMocks(mockClient, signerMock, true)
		lbcMock.EXPECT().UpdateProvider(mock.Anything, name, url).Return(nil, assert.AnError).Once()
		result, err := lbc.UpdateProvider(name, url)
		require.Error(t, err)
		assert.Empty(t, result)
		lbcMock.AssertExpectations(t)

		_, _ = prepareTxMocks(mockClient, signerMock, true)
		lbcMock.EXPECT().UpdateProvider(mock.Anything, name, url).Return(nil, nil).Once()
		result, err = lbc.UpdateProvider(name, url)
		require.Error(t, err)
		assert.Empty(t, result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling (updateProvider tx reverted)", func(t *testing.T) {
		tx, _ := prepareTxMocks(mockClient, signerMock, false)
		lbcMock.EXPECT().UpdateProvider(mock.Anything, name, url).Return(tx, nil).Once()
		result, err := lbc.UpdateProvider(name, url)
		require.ErrorContains(t, err, "update provider error")
		lbcMock.AssertExpectations(t)
		assert.Equal(t, tx.Hash().String(), result)
	})
}

func TestLiquidityBridgeContractImpl_GetCollateral(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(dummyClient, test.AnyAddress, lbcMock, nil, rootstock.RetryParams{}, time.Duration(1))
	t.Run("Success", func(t *testing.T) {
		lbcMock.On("GetCollateral", mock.Anything, parsedAddress).Return(big.NewInt(500), nil).Once()
		result, err := lbc.GetCollateral(parsedAddress.String())
		require.NoError(t, err)
		assert.Equal(t, entities.NewWei(500), result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling on GetCollateral call error", func(t *testing.T) {
		lbcMock.On("GetCollateral", mock.Anything, parsedAddress).Return(nil, assert.AnError).Once()
		result, err := lbc.GetCollateral(parsedAddress.String())
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Error handling on invalid address for getting collateral", func(t *testing.T) {
		result, err := lbc.GetCollateral(test.AnyString)
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestLiquidityBridgeContractImpl_GetPegoutCollateral(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(dummyClient, test.AnyAddress, lbcMock, nil, rootstock.RetryParams{}, time.Duration(1))
	t.Run("Success", func(t *testing.T) {
		lbcMock.On("GetPegoutCollateral", mock.Anything, parsedAddress).Return(big.NewInt(500), nil).Once()
		result, err := lbc.GetPegoutCollateral(parsedAddress.String())
		require.NoError(t, err)
		assert.Equal(t, entities.NewWei(500), result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling on GetPegoutCollateral call error", func(t *testing.T) {
		lbcMock.On("GetPegoutCollateral", mock.Anything, parsedAddress).Return(nil, assert.AnError).Once()
		result, err := lbc.GetPegoutCollateral(parsedAddress.String())
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Error handling on invalid address for getting pegout collateral", func(t *testing.T) {
		result, err := lbc.GetPegoutCollateral(test.AnyString)
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestLiquidityBridgeContractImpl_GetMinimumCollateral(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(dummyClient, test.AnyAddress, lbcMock, nil, rootstock.RetryParams{}, time.Duration(1))
	t.Run("Success", func(t *testing.T) {
		lbcMock.On("GetMinCollateral", mock.Anything).Return(big.NewInt(500), nil).Once()
		result, err := lbc.GetMinimumCollateral()
		require.NoError(t, err)
		assert.Equal(t, entities.NewWei(500), result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling on GetMinCollateral call fail", func(t *testing.T) {
		lbcMock.On("GetMinCollateral", mock.Anything).Return(nil, assert.AnError).Once()
		result, err := lbc.GetMinimumCollateral()
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestLiquidityBridgeContractImpl_AddCollateral(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		lbcMock,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
	)
	txMatchFunction := mock.MatchedBy(func(opts *bind.TransactOpts) bool {
		return opts.Value.Cmp(big.NewInt(500)) == 0 && bytes.Equal(opts.From.Bytes(), parsedAddress.Bytes())
	})
	t.Run("Success", func(t *testing.T) {
		tx, _ := prepareTxMocks(mockClient, signerMock, true, valueModifier(big.NewInt(500)))
		lbcMock.On("AddCollateral", txMatchFunction).Return(tx, nil).Once()
		err := lbc.AddCollateral(entities.NewWei(500))
		require.NoError(t, err)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling when sending addCollateral tx", func(t *testing.T) {
		_, _ = prepareTxMocks(mockClient, signerMock, true)
		lbcMock.On("AddCollateral", txMatchFunction).Return(nil, assert.AnError).Once()
		err := lbc.AddCollateral(entities.NewWei(500))
		require.Error(t, err)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling (addCollateral tx reverted)", func(t *testing.T) {
		tx, _ := prepareTxMocks(mockClient, signerMock, false, valueModifier(big.NewInt(500)))
		lbcMock.On("AddCollateral", txMatchFunction).Return(tx, nil).Once()
		err := lbc.AddCollateral(entities.NewWei(500))
		require.ErrorContains(t, err, "error adding pegin collateral")
		lbcMock.AssertExpectations(t)
	})
}

func TestLiquidityBridgeContractImpl_AddPegoutCollateral(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		lbcMock,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
	)
	txMatchFunction := mock.MatchedBy(func(opts *bind.TransactOpts) bool {
		return opts.Value.Cmp(big.NewInt(777)) == 0 && bytes.Equal(opts.From.Bytes(), parsedAddress.Bytes())
	})
	t.Run("Success", func(t *testing.T) {
		tx, _ := prepareTxMocks(mockClient, signerMock, true, valueModifier(big.NewInt(777)))
		lbcMock.On("AddPegoutCollateral", txMatchFunction).Return(tx, nil).Once()
		err := lbc.AddPegoutCollateral(entities.NewWei(777))
		require.NoError(t, err)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling when sending addPegoutCollateral tx", func(t *testing.T) {
		_, _ = prepareTxMocks(mockClient, signerMock, true)
		lbcMock.On("AddPegoutCollateral", txMatchFunction).Return(nil, assert.AnError).Once()
		err := lbc.AddPegoutCollateral(entities.NewWei(777))
		require.Error(t, err)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling (addPegoutCollateral tx reverted)", func(t *testing.T) {
		tx, _ := prepareTxMocks(mockClient, signerMock, false, valueModifier(big.NewInt(777)))
		lbcMock.On("AddPegoutCollateral", txMatchFunction).Return(tx, nil).Once()
		err := lbc.AddPegoutCollateral(entities.NewWei(777))
		require.ErrorContains(t, err, "error adding pegout collateral")
		lbcMock.AssertExpectations(t)
	})
}

func TestLiquidityBridgeContractImpl_WithdrawCollateral(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		lbcMock,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
	)
	t.Run("Success", func(t *testing.T) {
		tx, _ := prepareTxMocks(mockClient, signerMock, true)
		lbcMock.On("WithdrawCollateral", mock.Anything).Return(tx, nil).Once()
		err := lbc.WithdrawCollateral()
		require.NoError(t, err)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling when sending withdrawCollateral tx", func(t *testing.T) {
		_, _ = prepareTxMocks(mockClient, signerMock, true)
		lbcMock.On("WithdrawCollateral", mock.Anything).Return(nil, assert.AnError).Once()
		err := lbc.WithdrawCollateral()
		require.Error(t, err)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling (withdrawCollateral tx reverted)", func(t *testing.T) {
		tx, _ := prepareTxMocks(mockClient, signerMock, false)
		lbcMock.On("WithdrawCollateral", mock.Anything).Return(tx, nil).Once()
		err := lbc.WithdrawCollateral()
		require.ErrorContains(t, err, "withdraw pegin collateral error")
		lbcMock.AssertExpectations(t)
	})
}

func TestLiquidityBridgeContractImpl_GetBalance(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(dummyClient, test.AnyAddress, lbcMock, nil, rootstock.RetryParams{}, time.Duration(1))
	t.Run("Success", func(t *testing.T) {
		lbcMock.On("GetBalance", mock.Anything, parsedAddress).Return(big.NewInt(600), nil).Once()
		result, err := lbc.GetBalance(parsedAddress.String())
		require.NoError(t, err)
		assert.Equal(t, entities.NewWei(600), result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling on GetBalance call error", func(t *testing.T) {
		lbcMock.On("GetBalance", mock.Anything, parsedAddress).Return(nil, assert.AnError).Once()
		result, err := lbc.GetBalance(parsedAddress.String())
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Error handling on invalid address for getting balance", func(t *testing.T) {
		result, err := lbc.GetBalance(test.AnyString)
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestLiquidityBridgeContractImpl_CallForUser(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		lbcMock,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
	)
	modifiers := []txModifier{valueModifier(big.NewInt(1234)), gasLimitModifier(8000)}
	var gasLimit uint64 = 8000
	txConfig := blockchain.TransactionConfig{Value: entities.NewWei(1234), GasLimit: &gasLimit}
	optsMatchFunction := mock.MatchedBy(func(opts *bind.TransactOpts) bool {
		return opts.Value.Cmp(txConfig.Value.AsBigInt()) == 0 && opts.GasLimit == *txConfig.GasLimit
	})
	t.Run("Success", func(t *testing.T) {
		tx, receipt := prepareTxMocks(mockClient, signerMock, true, modifiers...)

		lbcMock.On("CallForUser", optsMatchFunction, parsedPeginQuote).Return(tx, nil).Once()
		callForUserReturn, err := lbc.CallForUser(txConfig, peginQuote)
		fmt.Printf("callForUserReturn: %+v\n", callForUserReturn)
		require.NoError(t, err)
		assert.Equal(t, tx.Hash().String(), callForUserReturn.TxHash)
		assert.Equal(t, receipt.GasUsed, callForUserReturn.GasUsed)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling when sending callForUser tx", func(t *testing.T) {
		_, _ = prepareTxMocks(mockClient, signerMock, true)
		lbcMock.On("CallForUser", optsMatchFunction, parsedPeginQuote).Return(nil, assert.AnError).Once()
		callForUserReturn, err := lbc.CallForUser(txConfig, peginQuote)
		require.Error(t, err)
		assert.Empty(t, callForUserReturn.TxHash)
	})
	t.Run("Error handling (callForUser tx reverted)", func(t *testing.T) {
		tx, receipt := prepareTxMocks(mockClient, signerMock, false, modifiers...)
		lbcMock.On("CallForUser", mock.Anything, parsedPeginQuote).Return(tx, nil).Once()
		callForUserReturn, err := lbc.CallForUser(txConfig, peginQuote)
		require.ErrorContains(t, err, "call for user error: transaction reverted")
		assert.Equal(t, tx.Hash().String(), callForUserReturn.TxHash)
		assert.Equal(t, receipt.GasUsed, callForUserReturn.GasUsed)
	})
	t.Run("Error handling (invalid quote)", func(t *testing.T) {
		invalid := peginQuote
		invalid.LbcAddress = ""
		callForUserReturn, err := lbc.CallForUser(txConfig, invalid)
		require.Error(t, err, "call for user error: transaction reverted")
		assert.Empty(t, callForUserReturn.TxHash)
	})
}

func TestLiquidityBridgeContractImpl_IsPegOutQuoteCompleted(t *testing.T) {
	const quoteHash = "762d73db7e80d845dae50d6ddda4d64d59f99352ead28afd51610e5674b08c0a"
	parsedQuoteHash := [32]byte{0x76, 0x2d, 0x73, 0xdb, 0x7e, 0x80, 0xd8, 0x45, 0xda, 0xe5, 0xd, 0x6d, 0xdd, 0xa4, 0xd6, 0x4d, 0x59, 0xf9, 0x93, 0x52, 0xea, 0xd2, 0x8a, 0xfd, 0x51, 0x61, 0xe, 0x56, 0x74, 0xb0, 0x8c, 0xa}
	lbcMock := &mocks.LbcAdapterMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(dummyClient, test.AnyAddress, lbcMock, nil, rootstock.RetryParams{}, time.Duration(1))
	t.Run("Success", func(t *testing.T) {
		lbcMock.On("IsPegOutQuoteCompleted", mock.Anything, parsedQuoteHash).Return(true, nil).Once()
		result, err := lbc.IsPegOutQuoteCompleted(quoteHash)
		require.NoError(t, err)
		assert.True(t, result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Should handle call error", func(t *testing.T) {
		lbcMock.On("IsPegOutQuoteCompleted", mock.Anything, parsedQuoteHash).Return(false, assert.AnError).Once()
		result, err := lbc.IsPegOutQuoteCompleted(quoteHash)
		require.Error(t, err)
		assert.False(t, result)
	})
	t.Run("Should handle parse error", func(t *testing.T) {
		result, err := lbc.IsPegOutQuoteCompleted(test.AnyString)
		require.Error(t, err)
		assert.False(t, result)
	})
	t.Run("Should return error when quote hash is not long enough", func(t *testing.T) {
		result, err := lbc.IsPegOutQuoteCompleted("0104050302")
		require.ErrorContains(t, err, "quote hash must be 32 bytes long")
		assert.False(t, result)
	})

}

func TestLiquidityBridgeContractImpl_RegisterPegin(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		lbcMock,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
	)
	registerParams := blockchain.RegisterPeginParams{
		QuoteSignature:        []byte{7, 8, 9},
		BitcoinRawTransaction: []byte{4, 5, 6},
		PartialMerkleTree:     []byte{1, 2, 3},
		BlockHeight:           big.NewInt(5),
		Quote:                 peginQuote,
	}
	callerMock := &mocks.LbcCallerBindingMock{}
	matchOptsFunc := func(opts *bind.TransactOpts) bool {
		return opts.From.String() == parsedAddress.String() && opts.GasLimit == 2500000
	}
	t.Run("Success", func(t *testing.T) {
		lbcMock.On("Caller").Return(callerMock).Once()
		callerMock.On("Call", mock.Anything, mock.Anything, "registerPegIn", parsedPeginQuote,
			registerParams.QuoteSignature, registerParams.BitcoinRawTransaction,
			registerParams.PartialMerkleTree, registerParams.BlockHeight,
		).Return(nil).Once()
		tx, receipt := prepareTxMocks(mockClient, signerMock, true)
		lbcMock.On("RegisterPegIn", mock.MatchedBy(matchOptsFunc), parsedPeginQuote,
			registerParams.QuoteSignature, registerParams.BitcoinRawTransaction,
			registerParams.PartialMerkleTree, registerParams.BlockHeight,
		).Return(tx, nil).Once()
		registerPeginReturn, err := lbc.RegisterPegin(registerParams)
		require.NoError(t, err)
		assert.Equal(t, tx.Hash().String(), registerPeginReturn.TxHash)
		assert.Equal(t, receipt.GasUsed, registerPeginReturn.GasUsed)
		lbcMock.AssertExpectations(t)
		callerMock.AssertExpectations(t)
	})
}

// nolint:funlen
func TestLiquidityBridgeContractImpl_RegisterPegin_ErrorHandling(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	callerMock := &mocks.LbcCallerBindingMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, lbcMock, signerMock, rootstock.RetryParams{}, time.Duration(1))
	registerParams := blockchain.RegisterPeginParams{QuoteSignature: []byte{7, 8, 9}, BitcoinRawTransaction: []byte{4, 5, 6}, PartialMerkleTree: []byte{1, 2, 3}, BlockHeight: big.NewInt(5), Quote: peginQuote}
	matchOptsFunc := func(opts *bind.TransactOpts) bool {
		return opts.From.String() == parsedAddress.String() && opts.GasLimit == 2500000
	}
	t.Run("Error handling (waiting for bridge)", func(t *testing.T) {
		lbcMock.On("Caller").Return(callerMock).Once()
		callerMock.On("Call", mock.Anything, mock.Anything, "registerPegIn", parsedPeginQuote,
			registerParams.QuoteSignature, registerParams.BitcoinRawTransaction,
			registerParams.PartialMerkleTree, registerParams.BlockHeight,
		).Return(errors.New("LBC031")).Once()
		result, err := lbc.RegisterPegin(registerParams)
		require.ErrorIs(t, err, blockchain.WaitingForBridgeError)
		assert.Empty(t, result)
		lbcMock.AssertExpectations(t)
		lbcMock.AssertNotCalled(t, "RegisterPegIn")
		callerMock.AssertExpectations(t)
	})
	t.Run("Error handling (Call error)", func(t *testing.T) {
		lbcMock.On("Caller").Return(callerMock).Once()
		callerMock.On("Call", mock.Anything, mock.Anything, "registerPegIn", parsedPeginQuote,
			registerParams.QuoteSignature, registerParams.BitcoinRawTransaction,
			registerParams.PartialMerkleTree, registerParams.BlockHeight,
		).Return(assert.AnError).Once()
		result, err := lbc.RegisterPegin(registerParams)
		require.Error(t, err)
		require.NotErrorIs(t, err, blockchain.WaitingForBridgeError)
		assert.Empty(t, result)
		lbcMock.AssertExpectations(t)
		lbcMock.AssertNotCalled(t, "RegisterPegIn")
		callerMock.AssertExpectations(t)
	})
	t.Run("Error handling (Transaction send error)", func(t *testing.T) {
		lbcMock.On("Caller").Return(callerMock).Once()
		callerMock.On("Call", mock.Anything, mock.Anything, "registerPegIn", parsedPeginQuote,
			registerParams.QuoteSignature, registerParams.BitcoinRawTransaction,
			registerParams.PartialMerkleTree, registerParams.BlockHeight,
		).Return(nil).Once()
		_, _ = prepareTxMocks(mockClient, signerMock, true)
		lbcMock.On("RegisterPegIn", mock.MatchedBy(matchOptsFunc), parsedPeginQuote,
			registerParams.QuoteSignature, registerParams.BitcoinRawTransaction,
			registerParams.PartialMerkleTree, registerParams.BlockHeight,
		).Return(nil, assert.AnError).Once()
		result, err := lbc.RegisterPegin(registerParams)
		require.ErrorContains(t, err, "register pegin error")
		assert.Empty(t, result)
		lbcMock.AssertExpectations(t)
		callerMock.AssertExpectations(t)
	})
	t.Run("Error handling (Transaction reverted)", func(t *testing.T) {
		lbcMock.On("Caller").Return(callerMock).Once()
		callerMock.On("Call", mock.Anything, mock.Anything, "registerPegIn", parsedPeginQuote,
			registerParams.QuoteSignature, registerParams.BitcoinRawTransaction,
			registerParams.PartialMerkleTree, registerParams.BlockHeight,
		).Return(nil).Once()
		tx, receipt := prepareTxMocks(mockClient, signerMock, false)
		lbcMock.On("RegisterPegIn", mock.MatchedBy(matchOptsFunc), parsedPeginQuote,
			registerParams.QuoteSignature, registerParams.BitcoinRawTransaction,
			registerParams.PartialMerkleTree, registerParams.BlockHeight,
		).Return(tx, nil).Once()
		registerPeginReturn, err := lbc.RegisterPegin(registerParams)
		require.ErrorContains(t, err, "register pegin error: transaction reverted")
		assert.Equal(t, tx.Hash().String(), registerPeginReturn.TxHash)
		assert.Equal(t, receipt.GasUsed, registerPeginReturn.GasUsed)
		lbcMock.AssertExpectations(t)
		callerMock.AssertExpectations(t)
	})
	t.Run("Error handling (invalid quote)", func(t *testing.T) {
		invalid := registerParams
		invalid.Quote.LbcAddress = ""
		result, err := lbc.RegisterPegin(invalid)
		require.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestLiquidityBridgeContractImpl_RefundUserPegOut(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	signer := &mocks.TransactionSignerMock{}
	client := &mocks.RpcClientBindingMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(
		rootstock.NewRskClient(client),
		test.AnyAddress,
		lbcMock,
		signer,
		rootstock.RetryParams{},
		time.Duration(1),
	)

	t.Run("should fail with invalid hash format", func(t *testing.T) {
		result, err := lbc.RefundUserPegOut("invalid hash")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid quote hash format")
		assert.Empty(t, result)
	})

	t.Run("should fail with invalid hash length", func(t *testing.T) {
		result, err := lbc.RefundUserPegOut("ab")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "quote hash must be 32 bytes long")
		assert.Empty(t, result)
	})

	t.Run("should fail if transaction fails", func(t *testing.T) {
		validHash := strings.Repeat("aa", 32)
		tx, _ := prepareTxMocks(client, signer, false)
		lbcMock.On("RefundUserPegOut", mock.Anything, mock.Anything).Return(tx, assert.AnError).Once()

		result, err := lbc.RefundUserPegOut(validHash)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "refund user peg out error")
		assert.Empty(t, result)
	})

	t.Run("should succeed", func(t *testing.T) {
		validHash := strings.Repeat("aa", 32)
		tx, _ := prepareTxMocks(client, signer, true)
		lbcMock.On("RefundUserPegOut", mock.Anything, mock.Anything).Return(tx, nil).Once()

		result, err := lbc.RefundUserPegOut(validHash)
		require.NoError(t, err)
		assert.Equal(t, tx.Hash().String(), result)

		lbcMock.AssertExpectations(t)
	})
}

// nolint:funlen
func TestLiquidityBridgeContractImpl_RefundPegout(t *testing.T) {
	var gasLimit uint64 = 500
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	refundParams := blockchain.RefundPegoutParams{
		QuoteHash:          [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
		BtcRawTx:           []byte{1, 2, 3},
		BtcBlockHeaderHash: [32]byte{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		MerkleBranchPath:   big.NewInt(5),
		MerkleBranchHashes: [][32]byte{{3, 2, 1}, {6, 5, 4}, {9, 8, 7}},
	}
	txConfig := blockchain.TransactionConfig{Value: entities.NewWei(1234), GasLimit: &gasLimit}
	matchOptsFunc := func(opts *bind.TransactOpts) bool {
		return opts.From.String() == parsedAddress.String() && opts.GasLimit == gasLimit
	}
	t.Run("Success", func(t *testing.T) {
		lbcMock := &mocks.LbcAdapterMock{}
		callerMock := &mocks.LbcCallerBindingMock{}
		lbc := rootstock.NewLiquidityBridgeContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, lbcMock, signerMock, rootstock.RetryParams{}, time.Duration(1))
		lbcMock.On("Caller").Return(callerMock).Once()
		callerMock.On("Call", mock.Anything, mock.Anything, "refundPegOut",
			refundParams.QuoteHash, refundParams.BtcRawTx, refundParams.BtcBlockHeaderHash,
			refundParams.MerkleBranchPath, refundParams.MerkleBranchHashes,
		).Return(nil).Once()
		tx, receipt := prepareTxMocks(mockClient, signerMock, true)
		lbcMock.On("RefundPegOut", mock.MatchedBy(matchOptsFunc),
			refundParams.QuoteHash, refundParams.BtcRawTx, refundParams.BtcBlockHeaderHash,
			refundParams.MerkleBranchPath, refundParams.MerkleBranchHashes,
		).Return(tx, nil).Once()
		result, err := lbc.RefundPegout(txConfig, refundParams)
		require.NoError(t, err)
		assert.Equal(t, tx.Hash().String(), result.TxHash)
		assert.Equal(t, receipt.GasUsed, result.GasUsed)
		lbcMock.AssertExpectations(t)
		callerMock.AssertExpectations(t)
	})
	t.Run("Error handling (waiting for bridge)", func(t *testing.T) {
		lbcMock := &mocks.LbcAdapterMock{}
		callerMock := &mocks.LbcCallerBindingMock{}
		lbc := rootstock.NewLiquidityBridgeContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, lbcMock, signerMock, rootstock.RetryParams{}, time.Duration(1))
		lbcMock.On("Caller").Return(callerMock).Once()
		callerMock.On("Call", mock.Anything, mock.Anything, "refundPegOut",
			refundParams.QuoteHash, refundParams.BtcRawTx, refundParams.BtcBlockHeaderHash,
			refundParams.MerkleBranchPath, refundParams.MerkleBranchHashes,
		).Return(errors.New("LBC049")).Once()
		signerMock.On("Address").Return(common.HexToAddress("0x1234567890123456789012345678901234567890")).Maybe()
		result, err := lbc.RefundPegout(txConfig, refundParams)
		require.ErrorIs(t, err, blockchain.WaitingForBridgeError)
		assert.Empty(t, result)
		lbcMock.AssertExpectations(t)
		lbcMock.AssertNotCalled(t, "RefundPegOut")
		callerMock.AssertExpectations(t)
	})
	t.Run("Error handling (Call error)", func(t *testing.T) {
		lbcMock := &mocks.LbcAdapterMock{}
		callerMock := &mocks.LbcCallerBindingMock{}
		lbc := rootstock.NewLiquidityBridgeContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, lbcMock, signerMock, rootstock.RetryParams{}, time.Duration(1))
		lbcMock.On("Caller").Return(callerMock).Once()
		callerMock.On("Call", mock.Anything, mock.Anything, "refundPegOut",
			refundParams.QuoteHash, refundParams.BtcRawTx, refundParams.BtcBlockHeaderHash,
			refundParams.MerkleBranchPath, refundParams.MerkleBranchHashes,
		).Return(assert.AnError).Once()
		result, err := lbc.RefundPegout(txConfig, refundParams)
		require.Error(t, err)
		assert.Empty(t, result)
		lbcMock.AssertExpectations(t)
		lbcMock.AssertNotCalled(t, "RefundPegOut")
		callerMock.AssertExpectations(t)
	})
	t.Run("Error handling (Transaction send error)", func(t *testing.T) {
		lbcMock := &mocks.LbcAdapterMock{}
		callerMock := &mocks.LbcCallerBindingMock{}
		lbc := rootstock.NewLiquidityBridgeContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, lbcMock, signerMock, rootstock.RetryParams{}, time.Duration(1))
		lbcMock.On("Caller").Return(callerMock).Once()
		callerMock.On("Call", mock.Anything, mock.Anything, "refundPegOut",
			refundParams.QuoteHash, refundParams.BtcRawTx, refundParams.BtcBlockHeaderHash,
			refundParams.MerkleBranchPath, refundParams.MerkleBranchHashes,
		).Return(nil).Once()
		_, _ = prepareTxMocks(mockClient, signerMock, true)
		lbcMock.On("RefundPegOut", mock.MatchedBy(matchOptsFunc),
			refundParams.QuoteHash, refundParams.BtcRawTx, refundParams.BtcBlockHeaderHash,
			refundParams.MerkleBranchPath, refundParams.MerkleBranchHashes,
		).Return(nil, assert.AnError).Once()
		result, err := lbc.RefundPegout(txConfig, refundParams)
		require.ErrorContains(t, err, "refund pegout error")
		assert.Empty(t, result)
		lbcMock.AssertExpectations(t)
		callerMock.AssertExpectations(t)
	})
	t.Run("Error handling (Transaction reverted)", func(t *testing.T) {
		lbcMock := &mocks.LbcAdapterMock{}
		callerMock := &mocks.LbcCallerBindingMock{}
		lbc := rootstock.NewLiquidityBridgeContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, lbcMock, signerMock, rootstock.RetryParams{}, time.Duration(1))
		lbcMock.On("Caller").Return(callerMock).Once()
		callerMock.On("Call", mock.Anything, mock.Anything, "refundPegOut",
			refundParams.QuoteHash, refundParams.BtcRawTx, refundParams.BtcBlockHeaderHash,
			refundParams.MerkleBranchPath, refundParams.MerkleBranchHashes,
		).Return(nil).Once()
		tx, _ := prepareTxMocks(mockClient, signerMock, false)
		lbcMock.On("RefundPegOut", mock.MatchedBy(matchOptsFunc),
			refundParams.QuoteHash, refundParams.BtcRawTx, refundParams.BtcBlockHeaderHash,
			refundParams.MerkleBranchPath, refundParams.MerkleBranchHashes,
		).Return(tx, nil).Once()
		result, err := lbc.RefundPegout(txConfig, refundParams)
		require.ErrorContains(t, err, "refund pegout error: transaction reverted")
		assert.Equal(t, tx.Hash().String(), result.TxHash)
		lbcMock.AssertExpectations(t)
		callerMock.AssertExpectations(t)
	})
}

func TestLiquidityBridgeContractImpl_IsOperationalPegin(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(dummyClient, test.AnyAddress, lbcMock, nil, rootstock.RetryParams{}, time.Duration(1))
	t.Run("Success", func(t *testing.T) {
		lbcMock.On("IsOperational", mock.Anything, parsedAddress).Return(true, nil).Once()
		result, err := lbc.IsOperationalPegin(parsedAddress.String())
		require.NoError(t, err)
		assert.True(t, result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling on IsOperational call fail", func(t *testing.T) {
		lbcMock.On("IsOperational", mock.Anything, parsedAddress).Return(true, assert.AnError).Once()
		_, err := lbc.IsOperationalPegin(parsedAddress.String())
		require.Error(t, err)
	})
	t.Run(invalidAddressTest, func(t *testing.T) {
		result, err := lbc.IsOperationalPegin(test.AnyString)
		require.ErrorIs(t, err, blockchain.InvalidAddressError)
		assert.False(t, result)
	})
}

func TestLiquidityBridgeContractImpl_IsOperationalPegout(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(dummyClient, test.AnyAddress, lbcMock, nil, rootstock.RetryParams{}, time.Duration(1))
	t.Run("Success", func(t *testing.T) {
		lbcMock.On("IsOperationalForPegout", mock.Anything, parsedAddress).Return(true, nil).Once()
		result, err := lbc.IsOperationalPegout(parsedAddress.String())
		require.NoError(t, err)
		assert.True(t, result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling on IsOperationalForPegout call fail", func(t *testing.T) {
		lbcMock.On("IsOperationalForPegout", mock.Anything, parsedAddress).Return(true, assert.AnError).Once()
		_, err := lbc.IsOperationalPegout(parsedAddress.String())
		require.Error(t, err)
	})
	t.Run(invalidAddressTest, func(t *testing.T) {
		result, err := lbc.IsOperationalPegout(test.AnyString)
		require.ErrorIs(t, err, blockchain.InvalidAddressError)
		assert.False(t, result)
	})
}

func TestLiquidityBridgeContractImpl_RegisterProvider(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		lbcMock,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
	)
	txConfig := blockchain.TransactionConfig{Value: entities.NewWei(800)}
	params := blockchain.ProviderRegistrationParams{
		Name:       "mock provider",
		ApiBaseUrl: "url.com",
		Status:     true,
		Type:       "both",
	}
	t.Run("Success", func(t *testing.T) {
		tx, _ := prepareTxMocks(mockClient, signerMock, true)
		receipt, err := mockClient.TransactionReceipt(context.Background(), tx.Hash())
		require.NoError(t, err)
		data, err := hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000d529ae9e860000")
		require.NoError(t, err)
		receipt.Logs = append(receipt.Logs, &geth.Log{
			Address: common.HexToAddress("0xAa9caf1e3967600578727f975F283446a3dA6612"),
			Topics: []common.Hash{
				common.HexToHash("0xa9d44d6e13bb3fee938c3f66d1103e91f8dc6b12d4405a55eea558e8f275aa6e"),
				common.HexToHash("0x0000000000000000000000004202bac9919c3412fc7c8be4e678e26279386603"),
			},
			Data:        data,
			BlockNumber: 5778711,
			TxHash:      common.HexToHash("0x37e52bd50866063727188751052e35510b8bc7d5de72541b84168cb2cb8b9c6c"),
			TxIndex:     0,
			BlockHash:   common.HexToHash("0xdc48007bd41ed3d8027aaac9c67fe1142107453b80b1fead090490fa8cbd751a"),
			Index:       0,
			Removed:     false,
		})
		lbcMock.On("ParseRegister", *receipt.Logs[0]).Return(&bindings.LiquidityBridgeContractRegister{
			Id:     big.NewInt(1),
			From:   parsedAddress,
			Amount: txConfig.Value.AsBigInt(),
			Raw:    *receipt.Logs[0],
		}, nil)
		mockClient.On("TransactionReceipt", mock.Anything, mock.Anything).Return(receipt, nil).Once()
		lbcMock.On("Register", mock.Anything,
			params.Name, params.ApiBaseUrl, params.Status, string(params.Type)).
			Return(tx, nil).Once()
		result, err := lbc.RegisterProvider(txConfig, params)
		require.NoError(t, err)
		assert.Equal(t, int64(1), result)
		lbcMock.AssertExpectations(t)
	})

}

func TestLiquidityBridgeContractImpl_RegisterProvider_ErrorHandling(t *testing.T) {
	const incompleteReceipt = "incomplete receipt"
	lbcMock := &mocks.LbcAdapterMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		lbcMock,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
	)
	txConfig := blockchain.TransactionConfig{Value: entities.NewWei(800)}
	params := blockchain.ProviderRegistrationParams{Name: "mock provider", ApiBaseUrl: "url.com", Status: true, Type: "both"}
	t.Run("Error handling (send transaction error)", func(t *testing.T) {
		_, _ = prepareTxMocks(mockClient, signerMock, true)
		lbcMock.On("Register", mock.Anything,
			params.Name, params.ApiBaseUrl, params.Status, string(params.Type)).
			Return(nil, assert.AnError).Once()
		result, err := lbc.RegisterProvider(txConfig, params)
		require.Error(t, err)
		assert.Equal(t, int64(0), result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling (receipt without event)", func(t *testing.T) {
		tx, _ := prepareTxMocks(mockClient, signerMock, true)
		lbcMock.On("Register", mock.Anything,
			params.Name, params.ApiBaseUrl, params.Status, string(params.Type)).
			Return(tx, nil).Once()
		result, err := lbc.RegisterProvider(txConfig, params)
		require.ErrorContains(t, err, incompleteReceipt)
		assert.Equal(t, int64(0), result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling (transaction revert)", func(t *testing.T) {
		tx, _ := prepareTxMocks(mockClient, signerMock, false)
		lbcMock.On("Register", mock.Anything,
			params.Name, params.ApiBaseUrl, params.Status, string(params.Type)).
			Return(tx, nil).Once()
		result, err := lbc.RegisterProvider(txConfig, params)
		require.ErrorContains(t, err, incompleteReceipt)
		assert.Equal(t, int64(0), result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling (parsing error)", func(t *testing.T) {
		tx, _ := prepareTxMocks(mockClient, signerMock, true)
		receipt, err := mockClient.TransactionReceipt(context.Background(), tx.Hash())
		require.NoError(t, err)
		receipt.Logs = append(receipt.Logs, &geth.Log{})
		lbcMock.On("ParseRegister", *receipt.Logs[0]).Return(nil, assert.AnError)
		mockClient.On("TransactionReceipt", mock.Anything, mock.Anything).Return(receipt, nil).Once()
		lbcMock.On("Register", mock.Anything,
			params.Name, params.ApiBaseUrl, params.Status, string(params.Type)).
			Return(tx, nil).Once()
		result, err := lbc.RegisterProvider(txConfig, params)
		require.ErrorContains(t, err, "error parsing register event")
		assert.Equal(t, int64(0), result)
		lbcMock.AssertExpectations(t)
	})
}

func TestLiquidityBridgeContractImpl_GetDepositEvents(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	iteratorMock := &mocks.EventIteratorAdapterMock[bindings.LiquidityBridgeContractPegOutDeposit]{}
	filterMatchFunc := func(from uint64, to uint64) func(opts *bind.FilterOpts) bool {
		return func(opts *bind.FilterOpts) bool {
			return from == opts.Start && to == *opts.End && opts.Context != nil
		}
	}
	lbc := rootstock.NewLiquidityBridgeContractImpl(dummyClient, test.AnyAddress, lbcMock, nil, rootstock.RetryParams{}, time.Duration(1))
	t.Run("Success", func(t *testing.T) {
		var from uint64 = 500
		var to uint64 = 1000
		lbcMock.On("FilterPegOutDeposit", mock.MatchedBy(filterMatchFunc(from, to)), [][32]uint8(nil), []common.Address(nil)).
			Return(&bindings.LiquidityBridgeContractPegOutDepositIterator{}, nil).Once()
		lbcMock.On("DepositEventIteratorAdapter", mock.AnythingOfType(depositIteratorString)).
			Return(iteratorMock)
		iteratorMock.On("Next").Return(true).Times(len(deposits))
		iteratorMock.On("Next").Return(false).Once()
		for _, deposit := range deposits {
			iteratorMock.On("Event").Return(deposit).Once()
		}
		iteratorMock.On("Error").Return(nil).Once()
		iteratorMock.On("Close").Return(nil).Once()
		result, err := lbc.GetDepositEvents(context.Background(), from, &to)
		require.NoError(t, err)
		assert.Equal(t, parsedDeposits, result)
		lbcMock.AssertExpectations(t)
		iteratorMock.AssertExpectations(t)
	})
	t.Run("Error handling when failed to get iterator", func(t *testing.T) {
		var from uint64 = 600
		var to uint64 = 1100
		lbcMock.On("FilterPegOutDeposit", mock.MatchedBy(filterMatchFunc(from, to)), [][32]uint8(nil), []common.Address(nil)).
			Return(nil, assert.AnError).Once()
		lbcMock.On("DepositEventIteratorAdapter", mock.AnythingOfType(depositIteratorString)).
			Return(nil)
		result, err := lbc.GetDepositEvents(context.Background(), from, &to)
		require.Error(t, err)
		assert.Nil(t, result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling on iterator error", func(t *testing.T) {
		var from uint64 = 700
		var to uint64 = 1200
		lbcMock.On("FilterPegOutDeposit", mock.MatchedBy(filterMatchFunc(from, to)), [][32]uint8(nil), []common.Address(nil)).
			Return(&bindings.LiquidityBridgeContractPegOutDepositIterator{}, nil).Once()
		lbcMock.On("DepositEventIteratorAdapter", mock.AnythingOfType(depositIteratorString)).
			Return(iteratorMock)
		iteratorMock.On("Next").Return(false).Once()
		iteratorMock.On("Error").Return(assert.AnError).Once()
		iteratorMock.On("Close").Return(nil).Once()
		result, err := lbc.GetDepositEvents(context.Background(), from, &to)
		require.Error(t, err)
		assert.Nil(t, result)
		lbcMock.AssertExpectations(t)
		iteratorMock.AssertExpectations(t)
	})
}

func TestLiquidityBridgeContractImpl_GetPunishmentEvents(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	iteratorMock := &mocks.EventIteratorAdapterMock[bindings.LiquidityBridgeContractPenalized]{}
	filterMatchFunc := func(from uint64, to uint64) func(opts *bind.FilterOpts) bool {
		return func(opts *bind.FilterOpts) bool {
			return from == opts.Start && to == *opts.End && opts.Context != nil
		}
	}
	lbc := rootstock.NewLiquidityBridgeContractImpl(dummyClient, test.AnyAddress, lbcMock, nil, rootstock.RetryParams{}, time.Duration(1))
	t.Run("Success", func(t *testing.T) {
		var from uint64 = 500
		var to uint64 = 1000
		lbcMock.On("FilterPenalized", mock.MatchedBy(filterMatchFunc(from, to))).
			Return(&bindings.LiquidityBridgeContractPenalizedIterator{}, nil).Once()
		lbcMock.On("PenalizedEventIteratorAdapter", mock.AnythingOfType(penalizedIteratorString)).
			Return(iteratorMock)
		iteratorMock.On("Next").Return(true).Times(len(penalizations))
		iteratorMock.On("Next").Return(false).Once()
		for _, deposit := range penalizations {
			iteratorMock.On("Event").Return(deposit).Once()
		}
		iteratorMock.On("Error").Return(nil).Once()
		iteratorMock.On("Close").Return(nil).Once()
		result, err := lbc.GetPenalizedEvents(context.Background(), from, &to)
		require.NoError(t, err)
		assert.Equal(t, parsedPenalizations, result)
		lbcMock.AssertExpectations(t)
		iteratorMock.AssertExpectations(t)
	})
	t.Run("Error handling when failed to get iterator", func(t *testing.T) {
		var from uint64 = 600
		var to uint64 = 1100
		lbcMock.On("FilterPenalized", mock.MatchedBy(filterMatchFunc(from, to))).
			Return(nil, assert.AnError).Once()
		lbcMock.On("PenalizedEventIteratorAdapter", mock.AnythingOfType(penalizedIteratorString)).
			Return(nil)
		result, err := lbc.GetPenalizedEvents(context.Background(), from, &to)
		require.Error(t, err)
		assert.Nil(t, result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling on iterator error", func(t *testing.T) {
		var from uint64 = 700
		var to uint64 = 1200
		lbcMock.On("FilterPenalized", mock.MatchedBy(filterMatchFunc(from, to))).
			Return(&bindings.LiquidityBridgeContractPenalizedIterator{}, nil).Once()
		lbcMock.On("PenalizedEventIteratorAdapter", mock.AnythingOfType(penalizedIteratorString)).
			Return(iteratorMock)
		iteratorMock.On("Next").Return(false).Once()
		iteratorMock.On("Error").Return(assert.AnError).Once()
		iteratorMock.On("Close").Return(nil).Once()
		result, err := lbc.GetPenalizedEvents(context.Background(), from, &to)
		require.Error(t, err)
		assert.Nil(t, result)
		lbcMock.AssertExpectations(t)
		iteratorMock.AssertExpectations(t)
	})
}

func TestLiquidityBridgeContractImpl_GetProvider(t *testing.T) {
	lbcMock := &mocks.LbcAdapterMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(dummyClient, test.AnyAddress, lbcMock, nil, rootstock.RetryParams{}, time.Duration(1))
	t.Run("Success", func(t *testing.T) {
		lbcMock.EXPECT().GetProvider(mock.Anything, parsedAddress).Return(bindings.LiquidityBridgeContractLiquidityProvider{
			Id:           big.NewInt(5),
			Provider:     parsedAddress,
			Name:         test.AnyString,
			ApiBaseUrl:   test.AnyUrl,
			Status:       true,
			ProviderType: string(liquidity_provider.FullProvider),
		}, nil).Once()
		result, err := lbc.GetProvider(parsedAddress.String())
		require.NoError(t, err)
		assert.Equal(t, liquidity_provider.RegisteredLiquidityProvider{
			Id:           5,
			Address:      parsedAddress.String(),
			Name:         test.AnyString,
			ApiBaseUrl:   test.AnyUrl,
			Status:       true,
			ProviderType: liquidity_provider.FullProvider,
		}, result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error handling on GetProvider call fail", func(t *testing.T) {
		lbcMock.On("GetProvider", mock.Anything, parsedAddress).Return(bindings.LiquidityBridgeContractLiquidityProvider{}, assert.AnError).Once()
		result, err := lbc.GetProvider(parsedAddress.String())
		require.Error(t, err)
		assert.Empty(t, result)
		lbcMock.AssertExpectations(t)
	})
	t.Run(invalidAddressTest, func(t *testing.T) {
		result, err := lbc.GetProvider(test.AnyString)
		require.ErrorIs(t, err, blockchain.InvalidAddressError)
		assert.Empty(t, result)
	})
	t.Run("Invalid type", func(t *testing.T) {
		lbcMock.EXPECT().GetProvider(mock.Anything, parsedAddress).Return(bindings.LiquidityBridgeContractLiquidityProvider{
			Id:           big.NewInt(5),
			Provider:     parsedAddress,
			Name:         test.AnyString,
			ApiBaseUrl:   test.AnyUrl,
			Status:       true,
			ProviderType: test.AnyString,
		}, nil).Once()
		result, err := lbc.GetProvider(parsedAddress.String())
		require.ErrorIs(t, err, liquidity_provider.InvalidProviderTypeError)
		assert.Empty(t, result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Provider not found", func(t *testing.T) {
		lbcMock.EXPECT().GetProvider(mock.Anything, parsedAddress).Return(bindings.LiquidityBridgeContractLiquidityProvider{}, errors.New("LBC001")).Once()
		result, err := lbc.GetProvider(parsedAddress.String())
		require.ErrorIs(t, err, liquidity_provider.ProviderNotFoundError)
		assert.Empty(t, result)
		lbcMock.AssertExpectations(t)
	})
}

type txModifier func(tx *geth.LegacyTx)

func valueModifier(value *big.Int) txModifier {
	return func(tx *geth.LegacyTx) {
		tx.Value = value
	}
}
func gasLimitModifier(value uint64) txModifier {
	return func(tx *geth.LegacyTx) {
		tx.Gas = value
	}
}

func prepareTxMocks(
	mockClient *mocks.RpcClientBindingMock,
	signerMock *mocks.TransactionSignerMock,
	success bool,
	txModifiers ...txModifier,
) (*geth.Transaction, *geth.Receipt) {
	legacyTx := &geth.LegacyTx{
		Nonce:    1,
		To:       &parsedAddress,
		Gas:      1,
		GasPrice: big.NewInt(1),
		Data:     nil,
	}

	mockClient.Calls = []mock.Call{}
	mockClient.ExpectedCalls = []*mock.Call{}
	signerMock.Calls = []mock.Call{}
	signerMock.ExpectedCalls = []*mock.Call{}

	for _, modifier := range txModifiers {
		modifier(legacyTx)
	}

	tx := geth.NewTx(legacyTx)

	receipt := &geth.Receipt{}
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = uint64(1000)
	if success == true {
		receipt.Status = 1
	}
	mockClient.On("TransactionReceipt", mock.Anything, mock.Anything).Return(receipt, nil).Once()
	signerMock.On("Sign", mock.Anything, mock.Anything).Return(tx, nil).Once()
	signerMock.On("Address").Return(parsedAddress)
	return tx, receipt
}
