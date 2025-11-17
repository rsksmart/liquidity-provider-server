package rootstock_test

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
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
	ProductFeeAmount:   entities.NewWei(500),
}

var parsedPeginQuote = bindings.QuotesPegInQuote{
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

func TestNewPeginContractImpl(t *testing.T) {
	contract := rootstock.NewPeginContractImpl(
		rootstock.NewRskClient(&mocks.RpcClientBindingMock{}),
		test.AnyAddress,
		&mocks.PeginContractAdapterMock{},
		&mocks.TransactionSignerMock{},
		rootstock.RetryParams{Retries: 1, Sleep: 1},
		time.Duration(1),
		Abis,
	)
	test.AssertNonZeroValues(t, contract)
}

func TestPeginContractImpl_GetBalance(t *testing.T) {
	contract := &mocks.PeginContractAdapterMock{}
	peginContract := rootstock.NewPeginContractImpl(dummyClient, test.AnyAddress, contract, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
	t.Run("Success", func(t *testing.T) {
		contract.On("GetBalance", mock.Anything, parsedAddress).Return(big.NewInt(600), nil).Once()
		result, err := peginContract.GetBalance(parsedAddress.String())
		require.NoError(t, err)
		assert.Equal(t, entities.NewWei(600), result)
		contract.AssertExpectations(t)
	})
	t.Run("Error handling on GetBalance call error", func(t *testing.T) {
		contract.On("GetBalance", mock.Anything, parsedAddress).Return(nil, assert.AnError).Once()
		result, err := peginContract.GetBalance(parsedAddress.String())
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Error handling on invalid address for getting balance", func(t *testing.T) {
		result, err := peginContract.GetBalance(test.AnyString)
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

// nolint:funlen
func TestPeginContractImpl_CallForUser(t *testing.T) {
	contractBinding := &mocks.PeginContractAdapterMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginContract := rootstock.NewPeginContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		contractBinding,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
		Abis,
	)
	modifiers := []txModifier{valueModifier(big.NewInt(1234)), gasLimitModifier(8000)}
	var gasLimit uint64 = 8000
	txConfig := blockchain.TransactionConfig{Value: entities.NewWei(1234), GasLimit: &gasLimit}
	optsMatchFunction := mock.MatchedBy(func(opts *bind.TransactOpts) bool {
		return opts.Value.Cmp(txConfig.Value.AsBigInt()) == 0 && opts.GasLimit == *txConfig.GasLimit
	})
	t.Run("Success", func(t *testing.T) {
		tx := prepareTxMocks(mockClient, signerMock, true, modifiers...)
		contractBinding.On("CallForUser", optsMatchFunction, parsedPeginQuote).Return(tx, nil).Once()
		expectedReceipt := blockchain.TransactionReceipt{
			TransactionHash:   tx.Hash().String(),
			BlockHash:         "0x0000000000000000000000000000000000000000000000000000000000000456",
			BlockNumber:       123,
			From:              parsedAddress.String(),
			To:                parsedAddress.String(),
			CumulativeGasUsed: big.NewInt(50000),
			GasUsed:           big.NewInt(21000),
			Value:             entities.NewBigWei(big.NewInt(1234)),
			GasPrice:          entities.NewWei(20000000000),
		}
		result, err := peginContract.CallForUser(txConfig, peginQuote)
		require.NoError(t, err)
		assert.Equal(t, expectedReceipt, result)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling when sending callForUser tx", func(t *testing.T) {
		_ = prepareTxMocks(mockClient, signerMock, true)
		contractBinding.On("CallForUser", optsMatchFunction, parsedPeginQuote).Return(nil, assert.AnError).Once()
		result, err := peginContract.CallForUser(txConfig, peginQuote)
		require.Error(t, err)
		assert.Empty(t, result.TransactionHash)
	})
	t.Run("Error handling (callForUser tx reverted)", func(t *testing.T) {
		tx := prepareTxMocks(mockClient, signerMock, false, modifiers...)
		contractBinding.On("CallForUser", mock.Anything, parsedPeginQuote).Return(tx, nil).Once()
		result, err := peginContract.CallForUser(txConfig, peginQuote)
		require.ErrorContains(t, err, "call for user error: transaction reverted")
		expectedReceipt := blockchain.TransactionReceipt{
			TransactionHash:   tx.Hash().String(),
			BlockHash:         "0x0000000000000000000000000000000000000000000000000000000000000456",
			BlockNumber:       123,
			From:              parsedAddress.String(),
			To:                parsedAddress.String(),
			CumulativeGasUsed: big.NewInt(50000),
			GasUsed:           big.NewInt(21000),
			Value:             entities.NewWei(1234),
			GasPrice:          entities.NewWei(20000000000),
		}
		assert.Equal(t, expectedReceipt, result)
	})
	t.Run("Error handling (invalid quote)", func(t *testing.T) {
		invalid := peginQuote
		invalid.LbcAddress = ""
		result, err := peginContract.CallForUser(txConfig, invalid)
		require.Error(t, err, "call for user error: transaction reverted")
		assert.Empty(t, result)
	})
}

func TestPeginContractImpl_GetAddress(t *testing.T) {
	peginContract := rootstock.NewPeginContractImpl(dummyClient, test.AnyAddress, nil, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
	assert.Equal(t, test.AnyAddress, peginContract.GetAddress())
}

func TestPeginContractImpl_DaoFeePercentage(t *testing.T) {
	contractBinding := &mocks.PeginContractAdapterMock{}
	t.Run("Success", func(t *testing.T) {
		contractBinding.EXPECT().GetFeePercentage(mock.Anything).Return(big.NewInt(1), nil).Once()
		peginContract := rootstock.NewPeginContractImpl(dummyClient, test.AnyAddress, contractBinding, nil, rootstock.RetryParams{Retries: 0, Sleep: 0}, time.Duration(1), Abis)
		percentage, err := peginContract.DaoFeePercentage()
		require.NoError(t, err)
		require.Equal(t, uint64(1), percentage)
	})
	t.Run("Error handling on ProductFeePercentage call fail", func(t *testing.T) {
		contractBinding.EXPECT().GetFeePercentage(mock.Anything).Return(nil, assert.AnError).Once()
		peginContract := rootstock.NewPeginContractImpl(dummyClient, test.AnyAddress, contractBinding, nil, rootstock.RetryParams{Retries: 0, Sleep: 0}, time.Duration(1), Abis)
		percentage, err := peginContract.DaoFeePercentage()
		require.Error(t, err)
		require.Zero(t, percentage)
	})
}

func TestPeginContractImpl_HashPeginQuote(t *testing.T) {
	contractBinding := &mocks.PeginContractAdapterMock{}
	hash := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	peginContract := rootstock.NewPeginContractImpl(
		dummyClient,
		test.AnyAddress,
		contractBinding,
		nil, rootstock.RetryParams{},
		time.Duration(1),
		Abis,
	)
	t.Run("Success", func(t *testing.T) {
		contractBinding.EXPECT().HashPegInQuote(mock.Anything, parsedPeginQuote).Return(hash, nil).Once()
		result, err := peginContract.HashPeginQuote(peginQuote)
		require.NoError(t, err)
		assert.Equal(t, hex.EncodeToString(hash[:]), result)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling on HashQuote call fail", func(t *testing.T) {
		contractBinding.EXPECT().HashPegInQuote(mock.Anything, parsedPeginQuote).Return([32]byte{}, assert.AnError).Once()
		result, err := peginContract.HashPeginQuote(peginQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestPeginContractImpl_HashPeginQuote_ParsingErrors(t *testing.T) {
	contractBinding := &mocks.PeginContractAdapterMock{}
	peginContract := rootstock.NewPeginContractImpl(dummyClient, test.AnyAddress, contractBinding, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
	validationFunction := func(peginQuote quote.PeginQuote) {
		result, err := peginContract.HashPeginQuote(peginQuote)
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

func TestPeginContractImpl_RegisterPegin(t *testing.T) {
	contractBinding := &mocks.PeginContractAdapterMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginContract := rootstock.NewPeginContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		contractBinding,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
		Abis,
	)
	registerParams := blockchain.RegisterPeginParams{
		QuoteSignature:        []byte{7, 8, 9},
		BitcoinRawTransaction: []byte{4, 5, 6},
		PartialMerkleTree:     []byte{1, 2, 3},
		BlockHeight:           big.NewInt(5),
		Quote:                 peginQuote,
	}
	callerMock := &mocks.ContractCallerBindingMock{}
	matchOptsFunc := func(opts *bind.TransactOpts) bool {
		return opts.From.String() == parsedAddress.String() && opts.GasLimit == 2500000
	}
	t.Run("Success", func(t *testing.T) {
		contractBinding.On("Caller").Return(callerMock).Once()
		callerMock.On("Call", mock.Anything, mock.Anything, "registerPegIn", parsedPeginQuote,
			registerParams.QuoteSignature, registerParams.BitcoinRawTransaction,
			registerParams.PartialMerkleTree, registerParams.BlockHeight,
		).Return(nil).Once()
		tx := prepareTxMocks(mockClient, signerMock, true)
		contractBinding.On("RegisterPegIn", mock.MatchedBy(matchOptsFunc), parsedPeginQuote,
			registerParams.QuoteSignature, registerParams.BitcoinRawTransaction,
			registerParams.PartialMerkleTree, registerParams.BlockHeight,
		).Return(tx, nil).Once()
		expectedReceipt := blockchain.TransactionReceipt{
			TransactionHash:   tx.Hash().String(),
			BlockHash:         "0x0000000000000000000000000000000000000000000000000000000000000456",
			BlockNumber:       123,
			From:              parsedAddress.String(),
			To:                parsedAddress.String(),
			CumulativeGasUsed: big.NewInt(50000),
			GasUsed:           big.NewInt(21000),
			Value:             entities.NewWei(0), // default value, no modifier applied
			GasPrice:          entities.NewWei(20000000000),
		}
		result, err := peginContract.RegisterPegin(registerParams)
		require.NoError(t, err)
		assert.Equal(t, expectedReceipt, result)
		contractBinding.AssertExpectations(t)
		callerMock.AssertExpectations(t)
	})
}

// nolint:funlen
func TestPeginContractImpl_RegisterPegin_ErrorHandling(t *testing.T) {
	contractBinding := &mocks.PeginContractAdapterMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	callerMock := &mocks.ContractCallerBindingMock{}
	peginContract := rootstock.NewPeginContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractBinding, signerMock, rootstock.RetryParams{}, time.Duration(1), Abis)
	registerParams := blockchain.RegisterPeginParams{QuoteSignature: []byte{7, 8, 9}, BitcoinRawTransaction: []byte{4, 5, 6}, PartialMerkleTree: []byte{1, 2, 3}, BlockHeight: big.NewInt(5), Quote: peginQuote}
	matchOptsFunc := func(opts *bind.TransactOpts) bool {
		return opts.From.String() == parsedAddress.String() && opts.GasLimit == 2500000
	}
	t.Run("Error handling (waiting for bridge)", func(t *testing.T) {
		e := NewRskRpcError("transaction reverted", "0xb9310b56")
		contractBinding.EXPECT().Caller().Return(callerMock).Once()
		callerMock.EXPECT().Call(mock.Anything, mock.Anything, "registerPegIn", parsedPeginQuote,
			registerParams.QuoteSignature, registerParams.BitcoinRawTransaction,
			registerParams.PartialMerkleTree, registerParams.BlockHeight,
		).Return(e).Once()
		result, err := peginContract.RegisterPegin(registerParams)
		require.ErrorIs(t, err, blockchain.WaitingForBridgeError)
		assert.Empty(t, result)
		contractBinding.AssertExpectations(t)
		contractBinding.AssertNotCalled(t, "RegisterPegIn")
		callerMock.AssertExpectations(t)
	})
	t.Run("Error handling (Call error)", func(t *testing.T) {
		contractBinding.On("Caller").Return(callerMock).Once()
		callerMock.On("Call", mock.Anything, mock.Anything, "registerPegIn", parsedPeginQuote,
			registerParams.QuoteSignature, registerParams.BitcoinRawTransaction,
			registerParams.PartialMerkleTree, registerParams.BlockHeight,
		).Return(assert.AnError).Once()
		result, err := peginContract.RegisterPegin(registerParams)
		require.Error(t, err)
		require.NotErrorIs(t, err, blockchain.WaitingForBridgeError)
		assert.Empty(t, result)
		contractBinding.AssertExpectations(t)
		contractBinding.AssertNotCalled(t, "RegisterPegIn")
		callerMock.AssertExpectations(t)
	})
	t.Run("Error handling (Transaction send error)", func(t *testing.T) {
		contractBinding.On("Caller").Return(callerMock).Once()
		callerMock.On("Call", mock.Anything, mock.Anything, "registerPegIn", parsedPeginQuote,
			registerParams.QuoteSignature, registerParams.BitcoinRawTransaction,
			registerParams.PartialMerkleTree, registerParams.BlockHeight,
		).Return(nil).Once()
		_ = prepareTxMocks(mockClient, signerMock, true)
		contractBinding.On("RegisterPegIn", mock.MatchedBy(matchOptsFunc), parsedPeginQuote,
			registerParams.QuoteSignature, registerParams.BitcoinRawTransaction,
			registerParams.PartialMerkleTree, registerParams.BlockHeight,
		).Return(nil, assert.AnError).Once()
		result, err := peginContract.RegisterPegin(registerParams)
		require.ErrorContains(t, err, "register pegin error")
		assert.Empty(t, result)
		contractBinding.AssertExpectations(t)
		callerMock.AssertExpectations(t)
	})
	t.Run("Error handling (Transaction reverted)", func(t *testing.T) {
		contractBinding.On("Caller").Return(callerMock).Once()
		callerMock.On("Call", mock.Anything, mock.Anything, "registerPegIn", parsedPeginQuote,
			registerParams.QuoteSignature, registerParams.BitcoinRawTransaction,
			registerParams.PartialMerkleTree, registerParams.BlockHeight,
		).Return(nil).Once()
		tx := prepareTxMocks(mockClient, signerMock, false)
		contractBinding.On("RegisterPegIn", mock.MatchedBy(matchOptsFunc), parsedPeginQuote,
			registerParams.QuoteSignature, registerParams.BitcoinRawTransaction,
			registerParams.PartialMerkleTree, registerParams.BlockHeight,
		).Return(tx, nil).Once()
		result, err := peginContract.RegisterPegin(registerParams)
		require.ErrorContains(t, err, "register pegin error: transaction reverted")
		// Should return populated receipt even on revert (for gas tracking)
		expectedReceipt := blockchain.TransactionReceipt{
			TransactionHash:   tx.Hash().String(),
			BlockHash:         "0x0000000000000000000000000000000000000000000000000000000000000456",
			BlockNumber:       123,
			From:              parsedAddress.String(),
			To:                parsedAddress.String(),
			CumulativeGasUsed: big.NewInt(50000),
			GasUsed:           big.NewInt(21000),
			Value:             entities.NewWei(0),
			GasPrice:          entities.NewWei(20000000000),
		}
		assert.Equal(t, expectedReceipt, result)
		contractBinding.AssertExpectations(t)
		callerMock.AssertExpectations(t)
	})
	t.Run("Error handling (invalid quote)", func(t *testing.T) {
		invalid := registerParams
		invalid.Quote.LbcAddress = ""
		result, err := peginContract.RegisterPegin(invalid)
		require.Error(t, err)
		assert.Empty(t, result)
	})
}
