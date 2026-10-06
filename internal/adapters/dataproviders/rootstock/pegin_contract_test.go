package rootstock_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	bindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin"
	commitfirst "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_commit_first"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"strings"
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
	ChainId:            31,
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
	GasFee:                      big.NewInt(100),
	ChainId:                     big.NewInt(31),
}

func TestNewPeginContractImpl(t *testing.T) {
	contract := rootstock.NewPeginContractImpl(
		rootstock.NewRskClient(&mocks.RpcClientBindingMock{}),
		test.AnyAddress,
		createBoundContractMock().contract,
		&mocks.TransactionSignerMock{},
		rootstock.RetryParams{Retries: 1, Sleep: 1},
		time.Duration(1),
		bindings.NewPeginContract(),
		Abis,
	)
	test.AssertNonZeroValues(t, contract)
}

func TestPeginContractImpl_GetBalance(t *testing.T) {
	contractMock := createBoundContractMock()
	peginBinding := bindings.NewPeginContract()
	peginContract := rootstock.NewPeginContractImpl(dummyClient, test.AnyAddress, contractMock.contract, nil, rootstock.RetryParams{}, time.Duration(1), peginBinding, Abis)
	t.Run("Success", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(peginBinding.PackGetBalance(parsedAddress)),
			mock.Anything,
		).Return(mustPackUint256(t, big.NewInt(600)), nil).Once()
		result, err := peginContract.GetBalance(parsedAddress.String())
		require.NoError(t, err)
		assert.Equal(t, entities.NewWei(600), result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on GetBalance call error", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(peginBinding.PackGetBalance(parsedAddress)),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
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
	contractMock := createBoundContractMock()
	peginBinding := bindings.NewPeginContract()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginContract := rootstock.NewPeginContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		contractMock.contract,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
		peginBinding,
		Abis,
	)
	var gasLimit uint64 = 8000
	txConfig := blockchain.TransactionConfig{Value: entities.NewWei(1234), GasLimit: &gasLimit}
	t.Run("Success", func(t *testing.T) {
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, big.NewInt(1234), peginBinding.PackCallForUser(parsedPeginQuote)),
		).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, true)
		expectedReceipt := blockchain.TransactionReceipt{
			TransactionHash:   "0x" + test.AnyHash,
			BlockHash:         "0x0000000000000000000000000000000000000000000000000000000000000456",
			BlockNumber:       123,
			From:              parsedAddress.String(),
			To:                test.AnyRskAddress,
			CumulativeGasUsed: big.NewInt(50000),
			GasUsed:           big.NewInt(21000),
			Value:             entities.NewBigWei(big.NewInt(1234)),
			GasPrice:          entities.NewWei(20000000000),
		}
		result, err := peginContract.CallForUser(txConfig, peginQuote)
		require.NoError(t, err)
		assert.Equal(t, expectedReceipt, result)
		contractMock.transactor.AssertExpectations(t)
	})
	t.Run("Error handling when sending callForUser tx", func(t *testing.T) {
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, big.NewInt(1234), peginBinding.PackCallForUser(parsedPeginQuote)),
		).Return(assert.AnError).Once()
		signerMock.EXPECT().Sign(mock.Anything, mock.Anything).RunAndReturn(func(addr common.Address, tx *geth.Transaction) (*geth.Transaction, error) {
			return tx, nil
		})
		result, err := peginContract.CallForUser(txConfig, peginQuote)
		require.Error(t, err)
		assert.Empty(t, result.TransactionHash)
	})
	t.Run("Error handling (callForUser tx reverted)", func(t *testing.T) {
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, big.NewInt(1234), peginBinding.PackCallForUser(parsedPeginQuote)),
		).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, false)
		result, err := peginContract.CallForUser(txConfig, peginQuote)
		require.ErrorContains(t, err, "call for user error: transaction reverted")
		expectedReceipt := blockchain.TransactionReceipt{
			TransactionHash:   "0x" + test.AnyHash,
			BlockHash:         "0x0000000000000000000000000000000000000000000000000000000000000456",
			BlockNumber:       123,
			From:              parsedAddress.String(),
			To:                test.AnyRskAddress,
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
	peginContract := rootstock.NewPeginContractImpl(dummyClient, test.AnyAddress, nil, nil, rootstock.RetryParams{}, time.Duration(1), nil, Abis)
	assert.Equal(t, test.AnyAddress, peginContract.GetAddress())
}

func TestPeginContractImpl_HashPeginQuote(t *testing.T) {
	contractMock := createBoundContractMock()
	peginBinding := bindings.NewPeginContract()
	hash := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	peginContract := rootstock.NewPeginContractImpl(
		dummyClient,
		test.AnyAddress,
		contractMock.contract,
		nil, rootstock.RetryParams{},
		time.Duration(1),
		peginBinding,
		Abis,
	)
	t.Run("Success", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(peginBinding.PackHashPegInQuote(parsedPeginQuote)),
			mock.Anything,
		).Return(mustPackBytes32(t, hash), nil).Once()
		result, err := peginContract.HashPeginQuote(peginQuote)
		require.NoError(t, err)
		assert.Equal(t, hex.EncodeToString(hash[:]), result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on HashQuote call fail", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(peginBinding.PackHashPegInQuote(parsedPeginQuote)),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := peginContract.HashPeginQuote(peginQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestPeginContractImpl_HashPeginQuote_ParsingErrors(t *testing.T) {
	contractMock := createBoundContractMock()
	peginBinding := bindings.NewPeginContract()
	peginContract := rootstock.NewPeginContractImpl(dummyClient, test.AnyAddress, contractMock.contract, nil, rootstock.RetryParams{}, time.Duration(1), peginBinding, Abis)
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
	contractMock := createBoundContractMock()
	peginBinding := bindings.NewPeginContract()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginContract := rootstock.NewPeginContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		contractMock.contract,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
		peginBinding,
		Abis,
	)
	registerParams := blockchain.RegisterPeginParams{
		QuoteSignature:        []byte{7, 8, 9},
		BitcoinRawTransaction: []byte{4, 5, 6},
		PartialMerkleTree:     []byte{1, 2, 3},
		BlockHeight:           big.NewInt(5),
		Quote:                 peginQuote,
	}
	t.Run("Success", func(t *testing.T) {
		txData := peginBinding.PackRegisterPegIn(parsedPeginQuote, registerParams.QuoteSignature, registerParams.BitcoinRawTransaction, registerParams.PartialMerkleTree, registerParams.BlockHeight)
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchRequestPegInCall(txData, parsedAddress, common.HexToAddress(test.AnyRskAddress), nil),
			mock.Anything,
		).Return(nil, nil).Once()
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 2500000, big.NewInt(0), txData),
		).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, true)
		expectedReceipt := blockchain.TransactionReceipt{
			TransactionHash:   "0x" + test.AnyHash,
			BlockHash:         "0x0000000000000000000000000000000000000000000000000000000000000456",
			BlockNumber:       123,
			From:              parsedAddress.String(),
			To:                test.AnyRskAddress,
			CumulativeGasUsed: big.NewInt(50000),
			GasUsed:           big.NewInt(21000),
			Value:             entities.NewWei(0), // default value, no modifier applied
			GasPrice:          entities.NewWei(20000000000),
		}
		result, err := peginContract.RegisterPegin(registerParams)
		require.NoError(t, err)
		assert.Equal(t, expectedReceipt, result)
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertExpectations(t)
	})
}

// nolint:funlen
func TestPeginContractImpl_RegisterPegin_ErrorHandling(t *testing.T) {
	peginBinding := bindings.NewPeginContract()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	registerParams := blockchain.RegisterPeginParams{QuoteSignature: []byte{7, 8, 9}, BitcoinRawTransaction: []byte{4, 5, 6}, PartialMerkleTree: []byte{1, 2, 3}, BlockHeight: big.NewInt(5), Quote: peginQuote}
	txData := peginBinding.PackRegisterPegIn(parsedPeginQuote, registerParams.QuoteSignature, registerParams.BitcoinRawTransaction, registerParams.PartialMerkleTree, registerParams.BlockHeight)
	signerMock.On("Address").Return(parsedAddress)
	t.Run("Error handling (waiting for bridge)", func(t *testing.T) {
		contractMock := createBoundContractMock()
		peginContract := rootstock.NewPeginContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractMock.contract, signerMock, rootstock.RetryParams{}, time.Duration(1), peginBinding, Abis)
		e := NewRskRpcError("transaction reverted", "0xb9310b56")
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchRequestPegInCall(txData, parsedAddress, common.HexToAddress(test.AnyRskAddress), nil),
			mock.Anything,
		).Return(nil, e).Once()
		result, err := peginContract.RegisterPegin(registerParams)
		require.ErrorIs(t, err, blockchain.WaitingForBridgeError)
		assert.Empty(t, result)
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertNotCalled(t, "SendTransaction")
	})
	t.Run("Error handling (Call error)", func(t *testing.T) {
		contractMock := createBoundContractMock()
		peginContract := rootstock.NewPeginContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractMock.contract, signerMock, rootstock.RetryParams{}, time.Duration(1), peginBinding, Abis)
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchRequestPegInCall(txData, parsedAddress, common.HexToAddress(test.AnyRskAddress), nil),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := peginContract.RegisterPegin(registerParams)
		require.Error(t, err)
		require.NotErrorIs(t, err, blockchain.WaitingForBridgeError)
		assert.Empty(t, result)
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertNotCalled(t, "SendTransaction")
	})
	t.Run("Error handling (Transaction send error)", func(t *testing.T) {
		contractMock := createBoundContractMock()
		peginContract := rootstock.NewPeginContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractMock.contract, signerMock, rootstock.RetryParams{}, time.Duration(1), peginBinding, Abis)
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchRequestPegInCall(txData, parsedAddress, common.HexToAddress(test.AnyRskAddress), nil),
			mock.Anything,
		).Return(nil, nil).Once()
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 2500000, big.NewInt(0), txData),
		).Return(assert.AnError).Once()
		signerMock.EXPECT().Sign(mock.Anything, mock.Anything).RunAndReturn(func(addr common.Address, tx *geth.Transaction) (*geth.Transaction, error) {
			return tx, nil
		})
		prepareTxMocks(&contractMock, mockClient, signerMock, true)
		result, err := peginContract.RegisterPegin(registerParams)
		require.ErrorContains(t, err, "register pegin error")
		assert.Empty(t, result)
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertExpectations(t)
	})
	t.Run("Error handling (Transaction reverted)", func(t *testing.T) {
		contractMock := createBoundContractMock()
		peginContract := rootstock.NewPeginContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractMock.contract, signerMock, rootstock.RetryParams{}, time.Duration(1), peginBinding, Abis)
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchRequestPegInCall(txData, parsedAddress, common.HexToAddress(test.AnyRskAddress), nil),
			mock.Anything,
		).Return(nil, nil).Once()
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 2500000, big.NewInt(0), txData),
		).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, false)
		result, err := peginContract.RegisterPegin(registerParams)
		require.ErrorContains(t, err, "register pegin error: transaction reverted")
		// Should return populated receipt even on revert (for gas tracking)
		expectedReceipt := blockchain.TransactionReceipt{
			TransactionHash:   "0x" + test.AnyHash,
			BlockHash:         "0x0000000000000000000000000000000000000000000000000000000000000456",
			BlockNumber:       123,
			From:              parsedAddress.String(),
			To:                test.AnyRskAddress,
			CumulativeGasUsed: big.NewInt(50000),
			GasUsed:           big.NewInt(21000),
			Value:             entities.NewWei(0),
			GasPrice:          entities.NewWei(20000000000),
		}
		assert.Equal(t, expectedReceipt, result)
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertExpectations(t)
	})
	t.Run("Error handling (invalid quote)", func(t *testing.T) {
		contractMock := createBoundContractMock()
		peginContract := rootstock.NewPeginContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractMock.contract, signerMock, rootstock.RetryParams{}, time.Duration(1), peginBinding, Abis)
		invalid := registerParams
		invalid.Quote.LbcAddress = ""
		result, err := peginContract.RegisterPegin(invalid)
		require.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestPeginContractImpl_Withdraw(t *testing.T) {
	contractMock := createBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginBinding := bindings.NewPeginContract()
	peginContract := rootstock.NewPeginContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractMock.contract, signerMock, rootstock.RetryParams{}, time.Duration(1), peginBinding, Abis)
	withdrawAmount := entities.NewWei(5000000000000000000)
	t.Run("Success", func(t *testing.T) {
		txData := peginBinding.PackWithdraw(withdrawAmount.AsBigInt())
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchRequestPegInCall(txData, parsedAddress, common.HexToAddress(test.AnyRskAddress), nil),
			mock.Anything,
		).Return(nil, nil).Once()
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 0, big.NewInt(0), txData),
		).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, true)
		err := peginContract.Withdraw(withdrawAmount)
		require.NoError(t, err)
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertExpectations(t)
	})
}

// nolint:funlen
func TestPeginContractImpl_Withdraw_ErrorHandling(t *testing.T) {
	contractMock := createBoundContractMock()
	peginBinding := bindings.NewPeginContract()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginContract := rootstock.NewPeginContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractMock.contract, signerMock, rootstock.RetryParams{}, time.Duration(1), peginBinding, Abis)
	withdrawAmount := entities.NewWei(5000000000000000000)
	signerMock.On("Address").Return(parsedAddress)
	t.Run("Error handling (dry-run revert NoBalance)", func(t *testing.T) {
		e := NewRskRpcError("transaction reverted", "0x29226653")
		contractMock.caller.EXPECT().CallContract(mock.Anything, mock.Anything, mock.Anything).Return(nil, e).Once()
		err := peginContract.Withdraw(withdrawAmount)
		require.ErrorContains(t, err, "withdraw reverted with: NoBalance")
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertNotCalled(t, "Transact")
	})
	t.Run("Error handling (dry-run parse error)", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		err := peginContract.Withdraw(withdrawAmount)
		require.ErrorContains(t, err, "error parsing withdraw result")
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertNotCalled(t, "Transact")
	})
	t.Run("Error handling (transaction send error)", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()
		contractMock.transactor.EXPECT().SendTransaction(mock.Anything, mock.Anything).Return(assert.AnError).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, true)
		signerMock.EXPECT().Sign(mock.Anything, mock.Anything).RunAndReturn(func(addr common.Address, tx *geth.Transaction) (*geth.Transaction, error) {
			return tx, nil
		})
		contractMock.transactor.EXPECT().HeaderByNumber(mock.Anything, mock.Anything).Return(&geth.Header{}, nil).Once()
		contractMock.transactor.EXPECT().SuggestGasPrice(mock.Anything).Return(big.NewInt(1), nil).Once()
		contractMock.transactor.EXPECT().PendingNonceAt(mock.Anything, mock.Anything).Return(1, nil).Once()
		err := peginContract.Withdraw(withdrawAmount)
		require.ErrorContains(t, err, "withdraw error")
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertExpectations(t)
	})
	t.Run("Error handling (transaction reverted)", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()
		contractMock.transactor.EXPECT().SendTransaction(mock.Anything, mock.Anything).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, false)
		signerMock.EXPECT().Sign(mock.Anything, mock.Anything).RunAndReturn(func(addr common.Address, tx *geth.Transaction) (*geth.Transaction, error) {
			return tx, nil
		})
		contractMock.transactor.EXPECT().HeaderByNumber(mock.Anything, mock.Anything).Return(&geth.Header{}, nil).Once()
		contractMock.transactor.EXPECT().SuggestGasPrice(mock.Anything).Return(big.NewInt(1), nil).Once()
		contractMock.transactor.EXPECT().PendingNonceAt(mock.Anything, mock.Anything).Return(1, nil).Once()
		err := peginContract.Withdraw(withdrawAmount)
		require.ErrorContains(t, err, "withdraw error: transaction failed")
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertExpectations(t)
	})
}

func TestPeginContractImpl_PausedStatus(t *testing.T) {
	contractMock := createBoundContractMock()
	peginBinding := bindings.NewPeginContract()
	contract := rootstock.NewPeginContractImpl(dummyClient, test.AnyAddress, contractMock.contract, nil, rootstock.RetryParams{}, time.Duration(1), peginBinding, Abis)
	t.Run("should return pause status result", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(peginBinding.PackPauseStatus()),
			mock.Anything,
		).Return(mustPackPauseStatus(t, generalPauseStatus{IsPaused: true, Reason: "test", Since: 123}), nil).Once()
		result, err := contract.PausedStatus()
		require.NoError(t, err)
		assert.Equal(t, blockchain.PauseStatus{IsPaused: true, Reason: "test", Since: 123}, result)
	})
	t.Run("should handle error checking pause status", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(peginBinding.PackPauseStatus()),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := contract.PausedStatus()
		require.Error(t, err)
		assert.Empty(t, result)
	})
	contractMock.caller.AssertExpectations(t)
}

func TestPeginContractImpl_HashPeginQuoteEIP712(t *testing.T) {
	t.Run("should return hash pegin quote eip712", func(t *testing.T) {
		contractMock := createBoundContractMock()
		peginBinding := bindings.NewPeginContract()
		hash := [32]byte{1, 2, 3}
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(peginBinding.PackHashPegInQuoteEIP712(parsedPeginQuote)),
			mock.Anything,
		).Return(mustPackBytes32(t, hash), nil).Once()
		contract := rootstock.NewPeginContractImpl(dummyClient, test.AnyAddress, contractMock.contract, nil, rootstock.RetryParams{}, time.Duration(1), peginBinding, Abis)
		result, err := contract.HashPeginQuoteEIP712(peginQuote)
		require.NoError(t, err)
		assert.Equal(t, hash, result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("should handle error hashing pegin quote eip712", func(t *testing.T) {
		contractMock := createBoundContractMock()
		peginBinding := bindings.NewPeginContract()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(peginBinding.PackHashPegInQuoteEIP712(parsedPeginQuote)),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		contract := rootstock.NewPeginContractImpl(dummyClient, test.AnyAddress, contractMock.contract, nil, rootstock.RetryParams{}, time.Duration(1), peginBinding, Abis)
		result, err := contract.HashPeginQuoteEIP712(peginQuote)
		require.Error(t, err)
		assert.Empty(t, result)
		contractMock.caller.AssertExpectations(t)
	})
}

// pinnedResolvePegInABI is the deployed resolvePegIn signature, parsed independently of the generated packer.
const pinnedResolvePegInABI = `[{"type":"function","name":"resolvePegIn","stateMutability":"nonpayable","inputs":[{"name":"rskAddr","type":"address"},{"name":"btcRawTransaction","type":"bytes"},{"name":"partialMerkleTree","type":"bytes"},{"name":"height","type":"uint256"}],"outputs":[{"name":"registerResult","type":"int256"}]}]`

const resolvePegInGasLimit = 2500000

var resolvePegInParams = blockchain.ResolvePegInParams{
	RskAddress:            test.AnyRskAddress,
	BitcoinRawTransaction: []byte{4, 5, 6},
	PartialMerkleTree:     []byte{1, 2, 3},
	BlockHeight:           big.NewInt(5),
}

type resolvePegInMocks struct {
	contractMock boundContractMock
	signerMock   *mocks.TransactionSignerMock
	mockClient   *mocks.RpcClientBindingMock
	pegin        blockchain.PeginContract
}

func newResolvePegInMocks() resolvePegInMocks {
	contractMock := createBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	signerMock.On("Address").Return(parsedAddress)
	pegin := rootstock.NewPeginContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyRskAddress,
		contractMock.contract,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
		nil,
		Abis,
	)
	return resolvePegInMocks{contractMock: contractMock, signerMock: signerMock, mockClient: mockClient, pegin: pegin}
}

func pinnedResolvePegIn(t *testing.T) abi.Method {
	t.Helper()
	parsed, err := abi.JSON(strings.NewReader(pinnedResolvePegInABI))
	require.NoError(t, err)
	return parsed.Methods["resolvePegIn"]
}

func resolvePegInCallData(t *testing.T) []byte {
	t.Helper()
	method := pinnedResolvePegIn(t)
	args, err := method.Inputs.Pack(
		common.HexToAddress(resolvePegInParams.RskAddress),
		resolvePegInParams.BitcoinRawTransaction,
		resolvePegInParams.PartialMerkleTree,
		resolvePegInParams.BlockHeight,
	)
	require.NoError(t, err)
	return append(method.ID, args...)
}

func resolvePegInOutput(t *testing.T, registerResult int64) []byte {
	t.Helper()
	output, err := pinnedResolvePegIn(t).Outputs.Pack(big.NewInt(registerResult))
	require.NoError(t, err)
	return output
}

func (m resolvePegInMocks) expectDryRun(t *testing.T, output []byte, revert error) {
	m.contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchRequestPegInCall(resolvePegInCallData(t), parsedAddress, common.HexToAddress(test.AnyRskAddress), nil),
		mock.Anything,
	).Return(output, revert).Once()
}

func (m resolvePegInMocks) expectSend(t *testing.T, success bool, logs ...*geth.Log) {
	m.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(m.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), resolvePegInGasLimit, big.NewInt(0), resolvePegInCallData(t)),
	).Return(nil).Once()
	prepareTxMocks(&m.contractMock, m.mockClient, m.signerMock, success, logs...)
}

func (m resolvePegInMocks) assertNotSent(t *testing.T) {
	t.Helper()
	m.contractMock.transactor.AssertNotCalled(t, "SendTransaction", mock.Anything, mock.Anything)
}

func resolveRevert(t *testing.T, errorName string, args ...any) error {
	t.Helper()
	parsed, err := commitfirst.PeginCommitFirstContractMetaData.ParseABI()
	require.NoError(t, err)
	abiError := parsed.Errors[errorName]
	data, err := abiError.Inputs.Pack(args...)
	require.NoError(t, err)
	return NewRskRpcError("transaction reverted", "0x"+hex.EncodeToString(append(abiError.ID.Bytes()[:4], data...)))
}

type pegInResolvedFixture struct {
	pegInID       [32]byte
	claimer       common.Address
	registrant    common.Address
	released      *entities.Wei
	claimerPayout *entities.Wei
	registrantFee *entities.Wei
	userPayout    *entities.Wei
}

var resolvedFixture = pegInResolvedFixture{
	pegInID:       [32]byte{0xaa},
	claimer:       common.HexToAddress("0x5dE07e2BE63595854C396E2da291e0d1EdE15112"),
	registrant:    common.HexToAddress("0x892813507Bf3aBF2890759d2135Ec34f4909Fea5"),
	released:      entities.NewWei(1000),
	claimerPayout: entities.NewWei(990),
	registrantFee: entities.NewWei(10),
	userPayout:    entities.NewWei(0),
}

func mustPegInResolvedLog(t *testing.T, fixture pegInResolvedFixture, address common.Address) *geth.Log {
	t.Helper()
	parsed, err := commitfirst.PeginCommitFirstContractMetaData.ParseABI()
	require.NoError(t, err)
	event := parsed.Events["PegInResolved"]
	data, err := event.Inputs.NonIndexed().Pack(
		fixture.released.AsBigInt(),
		fixture.claimerPayout.AsBigInt(),
		fixture.registrantFee.AsBigInt(),
		fixture.userPayout.AsBigInt(),
	)
	require.NoError(t, err)
	return &geth.Log{
		Address: address,
		Topics: []common.Hash{
			event.ID,
			common.BytesToHash(fixture.pegInID[:]),
			common.BytesToHash(fixture.claimer.Bytes()),
			common.BytesToHash(fixture.registrant.Bytes()),
		},
		Data: data,
	}
}

func TestPeginContractImpl_ResolvePegIn(t *testing.T) {
	m := newResolvePegInMocks()
	m.expectDryRun(t, resolvePegInOutput(t, 1000), nil)
	m.expectSend(t, true, mustPegInResolvedLog(t, resolvedFixture, common.HexToAddress(test.AnyRskAddress)))

	result, err := m.pegin.ResolvePegIn(resolvePegInParams)

	require.NoError(t, err)
	assert.Equal(t, "0x"+test.AnyHash, result.Receipt.TransactionHash)
	assert.Equal(t, uint64(blockchain.SuccessfulTxStatus), result.Receipt.Status)
	assert.Equal(t, resolvedFixture.claimerPayout, result.ClaimerPayout)
	assert.Equal(t, resolvedFixture.registrantFee, result.RegistrantFee)
	m.contractMock.caller.AssertExpectations(t)
	m.contractMock.transactor.AssertExpectations(t)
}

func TestPeginContractImpl_ResolvePegIn_WaitingForBridge(t *testing.T) {
	m := newResolvePegInMocks()
	m.expectDryRun(t, resolvePegInOutput(t, -303), nil)

	result, err := m.pegin.ResolvePegIn(resolvePegInParams)

	require.ErrorIs(t, err, blockchain.WaitingForBridgeError)
	assert.Empty(t, result)
	m.assertNotSent(t)
}

func TestPeginContractImpl_ResolvePegIn_BridgeRejected(t *testing.T) {
	codes := []int64{0, -200, -302, -900}
	for _, code := range codes {
		t.Run(big.NewInt(code).String(), func(t *testing.T) {
			m := newResolvePegInMocks()
			m.expectDryRun(t, resolvePegInOutput(t, code), nil)

			result, err := m.pegin.ResolvePegIn(resolvePegInParams)

			require.ErrorIs(t, err, blockchain.ErrBridgeRejectedPegIn)
			require.ErrorContains(t, err, fmt.Sprintf("code %d", code))
			require.NotErrorIs(t, err, blockchain.WaitingForBridgeError)
			assert.Empty(t, result)
			m.assertNotSent(t)
		})
	}
}

// nolint:funlen
func TestPeginContractImpl_ResolvePegIn_DryRunErrors(t *testing.T) {
	enforcedPause := NewRskRpcError("transaction reverted", "0x"+hex.EncodeToString(crypto.Keccak256([]byte("EnforcedPause()"))[:4]))
	genericRevert := NewRskRpcError("revert", "0x08c379a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000047465737400000000000000000000000000000000000000000000000000000000")
	tests := []struct {
		name        string
		revert      error
		expectedErr error
		errContains string
	}{
		{name: "already processed", revert: resolveRevert(t, "PegInAlreadyProcessed", [32]byte{0xaa}), expectedErr: blockchain.ErrPegInAlreadyProcessed},
		{name: "not claimed", revert: resolveRevert(t, "PegInNotClaimed", [32]byte{0xaa}), expectedErr: blockchain.ErrPegInNotClaimed},
		{name: "commit-first error not expected from resolvePegIn", revert: resolveRevert(t, "InsufficientConfirmations", big.NewInt(1), big.NewInt(2)), errContains: "resolvePegIn reverted with: InsufficientConfirmations"},
		{name: "error not in the commit-first ABI", revert: enforcedPause, errContains: "error parsing resolvePegIn result"},
		{name: "generic revert", revert: genericRevert, errContains: "error parsing resolvePegIn result: found generic error: test"},
		{name: "call error without revert data", revert: assert.AnError, errContains: "error parsing resolvePegIn result"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := newResolvePegInMocks()
			m.expectDryRun(t, nil, tc.revert)

			result, err := m.pegin.ResolvePegIn(resolvePegInParams)

			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.ErrorContains(t, err, tc.errContains)
				require.NotErrorIs(t, err, blockchain.ErrPegInAlreadyProcessed)
				require.NotErrorIs(t, err, blockchain.ErrPegInNotClaimed)
			}
			assert.Empty(t, result)
			m.assertNotSent(t)
		})
	}
}

func TestPeginContractImpl_ResolvePegIn_SendErrors(t *testing.T) {
	t.Run("send fails", func(t *testing.T) {
		m := newResolvePegInMocks()
		m.expectDryRun(t, resolvePegInOutput(t, 1000), nil)
		m.contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(m.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), resolvePegInGasLimit, big.NewInt(0), resolvePegInCallData(t)),
		).Return(assert.AnError).Once()
		prepareTxMocks(&m.contractMock, m.mockClient, m.signerMock, true)

		result, err := m.pegin.ResolvePegIn(resolvePegInParams)

		require.ErrorIs(t, err, assert.AnError)
		require.ErrorContains(t, err, "resolve pegin error")
		assert.Empty(t, result)
	})
	t.Run("transaction reverted keeps the receipt", func(t *testing.T) {
		m := newResolvePegInMocks()
		m.expectDryRun(t, resolvePegInOutput(t, 1000), nil)
		m.expectSend(t, false)

		result, err := m.pegin.ResolvePegIn(resolvePegInParams)

		require.ErrorContains(t, err, "resolve pegin error: transaction reverted")
		assert.Equal(t, "0x"+test.AnyHash, result.Receipt.TransactionHash)
		assert.Nil(t, result.ClaimerPayout)
		assert.Nil(t, result.RegistrantFee)
	})
	t.Run("successful receipt without PegInResolved keeps the receipt", func(t *testing.T) {
		m := newResolvePegInMocks()
		m.expectDryRun(t, resolvePegInOutput(t, 1000), nil)
		otherContract := common.HexToAddress("0x0D8Fb5d32704DB2931e05DB91F64BcA6f76Ce573")
		m.expectSend(t, true, mustPegInResolvedLog(t, resolvedFixture, otherContract))

		result, err := m.pegin.ResolvePegIn(resolvePegInParams)

		require.ErrorIs(t, err, blockchain.ErrPegInNotResolved)
		assert.Equal(t, "0x"+test.AnyHash, result.Receipt.TransactionHash)
		assert.Nil(t, result.ClaimerPayout)
		assert.Nil(t, result.RegistrantFee)
	})
}

func TestPeginContractImpl_ResolvePegIn_InvalidAddress(t *testing.T) {
	m := newResolvePegInMocks()
	params := resolvePegInParams
	params.RskAddress = "invalid"

	result, err := m.pegin.ResolvePegIn(params)

	require.ErrorIs(t, err, blockchain.InvalidAddressError)
	assert.Empty(t, result)
	m.contractMock.caller.AssertNotCalled(t, "CallContract", mock.Anything, mock.Anything, mock.Anything)
	m.assertNotSent(t)
}

// pinnedRequestPegInABI is parsed independently of the generated packer under test.
const pinnedRequestPegInABI = `[{"type":"function","name":"requestPegIn","stateMutability":"payable","inputs":[{"name":"rskAddr","type":"address"},{"name":"btcTxSerialized","type":"bytes"},{"name":"btcBlockHash","type":"bytes32"},{"name":"merkleBranchPath","type":"uint256"},{"name":"merkleBranchHashes","type":"bytes32[]"}],"outputs":[{"name":"pegInId","type":"bytes32"}]}]`

const requestPegInEstimatedGas = uint64(1000)

var (
	strippedRawTx    = []byte{0x01, 0x00, 0x00, 0x00, 0x01, 0xff, 0xaa, 0xbb}
	requestBlockHash = [32]byte{0x11}
	requestPath      = big.NewInt(1)
	requestHashes    = [][32]byte{{0x22}}
)

func packPinnedRequestPegIn(t *testing.T, rskAddr common.Address, rawTx []byte, blockHash [32]byte, path *big.Int, hashes [][32]byte) []byte {
	t.Helper()
	parsed, err := abi.JSON(strings.NewReader(pinnedRequestPegInABI))
	require.NoError(t, err)
	calldata, err := parsed.Pack("requestPegIn", rskAddr, rawTx, blockHash, path, hashes)
	require.NoError(t, err)
	return calldata
}

func paddedRequestPegInGas() uint64 {
	return requestPegInEstimatedGas * 12 / 10
}

func stubEstimateRequestPegInGas(mockClient *mocks.RpcClientBindingMock) {
	mockClient.On("EstimateGas", mock.Anything, mock.Anything).Return(requestPegInEstimatedGas, nil).Once()
}

// stubRequestPegInSigner prepares the signing and code checks a broadcast needs,
// without stubbing the receipt lookup. Use it when the test drives the receipt
// wait itself.
func stubRequestPegInSigner(h requestPegInHarness) {
	h.signerMock.On("Address").Return(parsedAddress)
	key := testSignerKey()
	h.signerMock.EXPECT().Sign(mock.Anything, mock.Anything).RunAndReturn(func(_ common.Address, transaction *geth.Transaction) (*geth.Transaction, error) {
		return geth.SignTx(transaction, geth.HomesteadSigner{}, key)
	}).Once()
	h.contractMock.transactor.EXPECT().PendingCodeAt(mock.Anything, mock.Anything).Return([]byte{1}, nil).Maybe()
	h.contractMock.caller.EXPECT().CodeAt(mock.Anything, mock.Anything, mock.Anything).Return([]byte{1}, nil).Maybe()
}

func stubRequestPegInDryRun(h requestPegInHarness, expectedData []byte, value *big.Int, revert error) *ethereum.CallMsg {
	captured := new(ethereum.CallMsg)
	h.signerMock.On("Address").Return(parsedAddress).Maybe()
	h.mockClient.EXPECT().CallContract(
		mock.Anything,
		matchRequestPegInCall(expectedData, parsedAddress, common.HexToAddress(test.AnyRskAddress), value),
		mock.Anything,
	).RunAndReturn(func(_ context.Context, msg ethereum.CallMsg, _ *big.Int) ([]byte, error) {
		*captured = msg
		return nil, revert
	}).Once()
	return captured
}

// guardRejectedRequestPegIn stops a rejected request from panicking on a dry run
// that must never happen, so the assertions report the failure instead.
func guardRejectedRequestPegIn(h requestPegInHarness) {
	h.signerMock.On("Address").Return(parsedAddress).Maybe()
	h.contractMock.caller.EXPECT().CodeAt(mock.Anything, mock.Anything, mock.Anything).Return([]byte{1}, nil).Maybe()
	h.contractMock.caller.EXPECT().CallContract(mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	h.mockClient.EXPECT().CallContract(mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Maybe()
}

func assertPayableDryRun(t *testing.T, h requestPegInHarness, expectedData []byte, value *big.Int) {
	t.Helper()
	h.mockClient.AssertCalled(
		t,
		"CallContract",
		mock.Anything,
		matchRequestPegInCall(expectedData, parsedAddress, common.HexToAddress(test.AnyRskAddress), value),
		mock.Anything,
	)
	h.contractMock.caller.AssertNotCalled(t, "CallContract", mock.Anything, mock.Anything, mock.Anything)
}

func assertNoDryRun(t *testing.T, h requestPegInHarness) {
	t.Helper()
	h.mockClient.AssertNotCalled(t, "CallContract", mock.Anything, mock.Anything, mock.Anything)
	h.contractMock.caller.AssertNotCalled(t, "CallContract", mock.Anything, mock.Anything, mock.Anything)
}

func assertRequestPegInNotSent(t *testing.T, h requestPegInHarness) {
	t.Helper()
	h.contractMock.transactor.AssertNotCalled(t, "SendTransaction", mock.Anything, mock.Anything)
	h.mockClient.AssertNotCalled(t, "SendTransaction", mock.Anything, mock.Anything)
}

func assertRequestPegInNotEstimated(t *testing.T, h requestPegInHarness) {
	t.Helper()
	h.contractMock.transactor.AssertNotCalled(t, "EstimateGas", mock.Anything, mock.Anything)
	h.mockClient.AssertNotCalled(t, "EstimateGas", mock.Anything, mock.Anything)
}

type pegInRequestedFixture struct {
	pegInID     [32]byte
	claimer     common.Address
	rskAddress  common.Address
	amount      *entities.Wei
	netToUser   *entities.Wei
	callSuccess bool
}

func assertPegInRequestedEvent(t *testing.T, fixture pegInRequestedFixture, event blockchain.PegInRequestedEvent) {
	t.Helper()
	assert.Equal(t, fixture.pegInID, event.PegInId)
	assert.Equal(t, fixture.claimer.Hex(), event.Claimer)
	assert.Equal(t, fixture.rskAddress.Hex(), event.RskAddress)
	assert.Equal(t, fixture.amount, event.Amount)
	assert.Equal(t, fixture.netToUser, event.NetToUser)
	assert.Equal(t, fixture.callSuccess, event.CallSuccess)
}

type requestPegInHarness struct {
	contractMock boundContractMock
	signerMock   *mocks.TransactionSignerMock
	mockClient   *mocks.RpcClientBindingMock
	pegin        blockchain.PeginContract
}

func newRequestPegInHarness(t *testing.T) requestPegInHarness {
	t.Helper()
	contractMock := createBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	return requestPegInHarness{
		contractMock: contractMock,
		signerMock:   signerMock,
		mockClient:   mockClient,
		pegin:        newRequestPegInContract(t, contractMock, mockClient, signerMock),
	}
}

func newRequestPegInContract(
	t *testing.T,
	contractMock boundContractMock,
	mockClient *mocks.RpcClientBindingMock,
	signerMock *mocks.TransactionSignerMock,
) blockchain.PeginContract {
	t.Helper()
	return rootstock.NewPeginContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyRskAddress,
		contractMock.contract,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
		nil,
		Abis,
	)
}

func sampleRequestPegInParams(amount, fee *entities.Wei) blockchain.RequestPegInParams {
	return blockchain.RequestPegInParams{
		RskAddress:         parsedAddress.Hex(),
		BitcoinRawTx:       strippedRawTx,
		BtcBlockHash:       requestBlockHash,
		MerkleBranchPath:   requestPath,
		MerkleBranchHashes: requestHashes,
		Amount:             amount,
		Fee:                fee,
	}
}

func revertHexFromErrorID(t *testing.T, id common.Hash, tail []byte) string {
	t.Helper()
	return "0x" + hex.EncodeToString(append(id.Bytes()[:4], tail...))
}

func mustPegInRequestedLog(t *testing.T, fixture pegInRequestedFixture) *geth.Log {
	t.Helper()
	parsed, err := commitfirst.PeginCommitFirstContractMetaData.ParseABI()
	require.NoError(t, err)
	event := parsed.Events["PegInRequested"]
	data, err := event.Inputs.NonIndexed().Pack(
		fixture.amount.AsBigInt(),
		fixture.netToUser.AsBigInt(),
		fixture.callSuccess,
	)
	require.NoError(t, err)
	return &geth.Log{
		Address: common.HexToAddress(test.AnyRskAddress),
		Topics: []common.Hash{
			event.ID,
			common.BytesToHash(fixture.pegInID[:]),
			common.BytesToHash(fixture.claimer.Bytes()),
			common.BytesToHash(fixture.rskAddress.Bytes()),
		},
		Data: data,
	}
}

func transactionLogFromGeth(eventLog *geth.Log) blockchain.TransactionLog {
	topics := make([][32]byte, len(eventLog.Topics))
	for i, topic := range eventLog.Topics {
		topics[i] = topic
	}
	return blockchain.TransactionLog{
		Address: eventLog.Address.Hex(),
		Topics:  topics,
		Data:    eventLog.Data,
		Removed: eventLog.Removed,
	}
}

func TestPeginContractImpl_RequestPegIn_PackingMatchesPinnedABI(t *testing.T) {
	h := newRequestPegInHarness(t)

	amount := entities.SatoshiToWei(1000)
	fee := entities.SatoshiToWei(100)
	expectedValue := new(entities.Wei).Sub(amount, fee)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	fixture := pegInRequestedFixture{
		pegInID:     [32]byte{0xab},
		claimer:     parsedAddress,
		rskAddress:  parsedAddress,
		amount:      amount,
		netToUser:   expectedValue,
		callSuccess: true,
	}
	eventLog := mustPegInRequestedLog(t, fixture)
	gasLimit := paddedRequestPegInGas()

	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, expectedValue.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&h.contractMock, h.mockClient, h.signerMock, true, eventLog)
	stubEstimateRequestPegInGas(h.mockClient)
	stubRequestPegInDryRun(h, expectedData, expectedValue.AsBigInt(), nil)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(amount, fee))
	require.NoError(t, err)
	assert.Equal(t, expectedValue, result.Receipt.Value)
	assert.True(t, strings.HasPrefix(hex.EncodeToString(expectedData), "fc73bbd3"))
	assertPegInRequestedEvent(t, fixture, result.Event)
	assertPayableDryRun(t, h, expectedData, expectedValue.AsBigInt())
	h.contractMock.transactor.AssertExpectations(t)
	h.mockClient.AssertCalled(t, "EstimateGas", mock.Anything, mock.Anything)
}

func TestPeginContractImpl_RequestPegIn_FirstOutputValueNotSum(t *testing.T) {
	address := "2N2Sg8C2uX1YtugYSxEQvRqf9V2EivxcWER"
	txInfo := blockchain.BitcoinTransactionInformation{
		Outputs: map[string][]*entities.Wei{
			address: {entities.SatoshiToWei(1000), entities.SatoshiToWei(2000)},
		},
	}
	first := txInfo.FirstOutputToAddress(address)
	sum := txInfo.AmountToAddress(address)
	require.Equal(t, entities.SatoshiToWei(1000), first)
	require.Equal(t, entities.SatoshiToWei(3000), sum)

	h := newRequestPegInHarness(t)
	fee := entities.SatoshiToWei(100)
	expectedValue := new(entities.Wei).Sub(first, fee)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	fixture := pegInRequestedFixture{
		pegInID:     [32]byte{0x01},
		claimer:     parsedAddress,
		rskAddress:  parsedAddress,
		amount:      first,
		netToUser:   expectedValue,
		callSuccess: true,
	}
	eventLog := mustPegInRequestedLog(t, fixture)
	gasLimit := paddedRequestPegInGas()

	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, expectedValue.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&h.contractMock, h.mockClient, h.signerMock, true, eventLog)
	stubEstimateRequestPegInGas(h.mockClient)
	stubRequestPegInDryRun(h, expectedData, expectedValue.AsBigInt(), nil)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(first, fee))
	require.NoError(t, err)
	assert.Equal(t, expectedValue, result.Receipt.Value)
	assert.NotEqual(t, new(entities.Wei).Sub(sum, fee), result.Receipt.Value)
	assertPegInRequestedEvent(t, fixture, result.Event)
	assertPayableDryRun(t, h, expectedData, expectedValue.AsBigInt())
}

func TestPeginContractImpl_RequestPegIn_SatToWeiBoundary(t *testing.T) {
	oneSat := entities.SatoshiToWei(1)
	require.Equal(t, entities.NewWei(10_000_000_000), oneSat)

	h := newRequestPegInHarness(t)
	fee := entities.NewWei(0)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	fixture := pegInRequestedFixture{
		pegInID:     [32]byte{0x02},
		claimer:     parsedAddress,
		rskAddress:  parsedAddress,
		amount:      oneSat,
		netToUser:   oneSat,
		callSuccess: true,
	}
	eventLog := mustPegInRequestedLog(t, fixture)
	gasLimit := paddedRequestPegInGas()

	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, oneSat.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&h.contractMock, h.mockClient, h.signerMock, true, eventLog)
	stubEstimateRequestPegInGas(h.mockClient)
	stubRequestPegInDryRun(h, expectedData, oneSat.AsBigInt(), nil)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(oneSat, fee))
	require.NoError(t, err)
	assert.Equal(t, oneSat, result.Receipt.Value)
	assertPegInRequestedEvent(t, fixture, result.Event)
	assertPayableDryRun(t, h, expectedData, oneSat.AsBigInt())
}

func TestPeginContractImpl_RequestPegIn_StatusZeroDoesNotClassifyRaceLoss(t *testing.T) {
	h := newRequestPegInHarness(t)
	amount := entities.SatoshiToWei(1000)
	fee := entities.NewWei(0)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	gasLimit := paddedRequestPegInGas()

	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, amount.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&h.contractMock, h.mockClient, h.signerMock, false)
	stubEstimateRequestPegInGas(h.mockClient)
	stubRequestPegInDryRun(h, expectedData, amount.AsBigInt(), nil)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(amount, fee))
	require.Error(t, err)
	require.ErrorContains(t, err, "request pegin error: transaction reverted")
	require.NotErrorIs(t, err, blockchain.ErrPegInAlreadyProcessed)
	assert.NotEmpty(t, result.Receipt.TransactionHash)
	assert.Equal(t, blockchain.PegInRequestedEvent{}, result.Event)
	assertPayableDryRun(t, h, expectedData, amount.AsBigInt())
}

func TestPeginContractImpl_RequestPegIn_DoesNotCheckPause(t *testing.T) {
	h := newRequestPegInHarness(t)
	amount := entities.SatoshiToWei(500)
	fee := entities.NewWei(0)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	fixture := pegInRequestedFixture{
		pegInID:     [32]byte{0x03},
		claimer:     parsedAddress,
		rskAddress:  parsedAddress,
		amount:      amount,
		netToUser:   amount,
		callSuccess: true,
	}
	eventLog := mustPegInRequestedLog(t, fixture)
	gasLimit := paddedRequestPegInGas()

	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, amount.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&h.contractMock, h.mockClient, h.signerMock, true, eventLog)
	stubEstimateRequestPegInGas(h.mockClient)
	stubRequestPegInDryRun(h, expectedData, amount.AsBigInt(), nil)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(amount, fee))
	require.NoError(t, err)
	assert.NotEmpty(t, result.Receipt.TransactionHash)
	assertPegInRequestedEvent(t, fixture, result.Event)
	assertPayableDryRun(t, h, expectedData, amount.AsBigInt())
	h.contractMock.transactor.AssertExpectations(t)
}

// ParseReceipt recovers the sender from the signature, not TransactOpts.From.
// This test prevents the signer mock from claiming an address its key does not own.
func TestPeginContractImpl_RequestPegIn_ReceiptSenderMatchesSigner(t *testing.T) {
	h := newRequestPegInHarness(t)
	amount := entities.SatoshiToWei(1000)
	fee := entities.NewWei(0)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	eventLog := mustPegInRequestedLog(t, pegInRequestedFixture{
		pegInID:     [32]byte{0x06},
		claimer:     parsedAddress,
		rskAddress:  parsedAddress,
		amount:      amount,
		netToUser:   amount,
		callSuccess: true,
	})

	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), paddedRequestPegInGas(), amount.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&h.contractMock, h.mockClient, h.signerMock, true, eventLog)
	stubEstimateRequestPegInGas(h.mockClient)
	stubRequestPegInDryRun(h, expectedData, amount.AsBigInt(), nil)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(amount, fee))
	require.NoError(t, err)
	assert.Equal(t, parsedAddress.Hex(), result.Receipt.From)
}

func TestPeginContractImpl_RequestPegIn_PreflightAlreadyProcessed(t *testing.T) {
	h := newRequestPegInHarness(t)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	pegInId := [32]byte{0x44}
	amount := entities.SatoshiToWei(1000)
	selector := commitfirst.PeginCommitFirstContractPegInAlreadyProcessedErrorID().Bytes()[:4]
	revertHex := "0x" + hex.EncodeToString(append(selector, mustPackBytes32(t, pegInId)...))

	stubRequestPegInDryRun(h, expectedData, amount.AsBigInt(), NewRskRpcError("execution reverted", revertHex))

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(amount, entities.NewWei(0)))
	require.ErrorIs(t, err, blockchain.ErrPegInAlreadyProcessed)
	assert.Empty(t, result.Receipt.TransactionHash)
	assertPayableDryRun(t, h, expectedData, amount.AsBigInt())
	assertRequestPegInNotSent(t, h)
	assertRequestPegInNotEstimated(t, h)
}

func TestPeginContractImpl_RequestPegIn_AmountBelowFee(t *testing.T) {
	h := newRequestPegInHarness(t)
	guardRejectedRequestPegIn(h)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(entities.NewWei(1), entities.NewWei(2)))
	require.ErrorIs(t, err, blockchain.ErrIncorrectFronting)
	assert.Empty(t, result.Receipt.TransactionHash)
	assertNoDryRun(t, h)
	assertRequestPegInNotSent(t, h)
	assertRequestPegInNotEstimated(t, h)
}

func TestPeginContractImpl_RequestPegIn_PreflightTypedErrors(t *testing.T) {
	pegInId := [32]byte{0x44}
	btcTxHash := [32]byte{0x55}
	cases := []struct {
		name    string
		errorID common.Hash
		tail    []byte
		want    error
	}{
		{"AlreadyProcessed", commitfirst.PeginCommitFirstContractPegInAlreadyProcessedErrorID(), mustPackBytes32(t, pegInId), blockchain.ErrPegInAlreadyProcessed},
		{"AddressNotRegistered", commitfirst.PeginCommitFirstContractAddressNotRegisteredErrorID(), mustPackAddress(t, parsedAddress), blockchain.ErrAddressNotRegistered},
		{"DepositOutputNotFound", commitfirst.PeginCommitFirstContractDepositOutputNotFoundErrorID(), append(mustPackAddress(t, parsedAddress), mustPackBytes32(t, btcTxHash)...), blockchain.ErrDepositOutputNotFound},
		{"InsufficientConfirmations", commitfirst.PeginCommitFirstContractInsufficientConfirmationsErrorID(), append(mustPackUint256(t, big.NewInt(1)), mustPackUint256(t, big.NewInt(6))...), blockchain.ErrInsufficientConfirmations},
		{"IncorrectFronting", commitfirst.PeginCommitFirstContractIncorrectFrontingErrorID(), append(mustPackUint256(t, big.NewInt(1000)), mustPackUint256(t, big.NewInt(500))...), blockchain.ErrIncorrectFronting},
		{"PegInBelowMinimum", commitfirst.PeginCommitFirstContractPegInBelowMinimumErrorID(), append(mustPackUint256(t, big.NewInt(1000)), mustPackUint256(t, big.NewInt(5000))...), blockchain.ErrPegInBelowMinimum},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertPreflightTypedError(t, tc.errorID, tc.tail, tc.want)
		})
	}
}

func assertPreflightTypedError(t *testing.T, errorID common.Hash, tail []byte, want error) {
	t.Helper()
	h := newRequestPegInHarness(t)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	amount := entities.SatoshiToWei(1000)

	stubRequestPegInDryRun(h, expectedData, amount.AsBigInt(), NewRskRpcError("execution reverted", revertHexFromErrorID(t, errorID, tail)))

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(amount, entities.NewWei(0)))
	require.ErrorIs(t, err, want)
	assert.Empty(t, result.Receipt.TransactionHash)
	assertPayableDryRun(t, h, expectedData, amount.AsBigInt())
	h.contractMock.transactor.AssertNotCalled(t, "SendTransaction", mock.Anything, mock.Anything)
	h.mockClient.AssertNotCalled(t, "SendTransaction", mock.Anything, mock.Anything)
	assertRequestPegInNotEstimated(t, h)
}

func TestPeginContractImpl_RequestPegIn_PreflightDoesNotInventRaceLoss(t *testing.T) {
	errorTestHex := "0x08c379a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000047465737400000000000000000000000000000000000000000000000000000000"

	t.Run("generic Error(string) revert", func(t *testing.T) {
		// The reason string must reach the caller. "Unknown error" would hide why the call reverted.
		assertPreflightJunkRevert(t, NewRskRpcError("execution reverted", errorTestHex), "requestPegIn reverted: test")
	})
	t.Run("Panic(uint256) revert", func(t *testing.T) {
		// Panic payloads must stay on the generic path. UnpackError would hide the mapped panic reason.
		const panicHex = "0x4e487b710000000000000000000000000000000000000000000000000000000000000011"
		wantMessage, unpackErr := abi.UnpackRevert(common.FromHex(panicHex))
		require.NoError(t, unpackErr)
		assertPreflightJunkRevert(t, NewRskRpcError("execution reverted", panicHex), "requestPegIn reverted: "+wantMessage)
	})
	t.Run("non-DataError", func(t *testing.T) {
		assertPreflightJunkRevert(t, assert.AnError, "error parsing requestPegIn result")
	})
	t.Run("short revert data", func(t *testing.T) {
		for _, revertHex := range []string{"0x", "0xaabbcc"} {
			t.Run(revertHex, func(t *testing.T) {
				assertPreflightJunkRevert(t, NewRskRpcError("execution reverted", revertHex), "requestPegIn reverted")
			})
		}
	})
}

func assertPreflightJunkRevert(t *testing.T, revert error, wantContains string) {
	t.Helper()
	h := newRequestPegInHarness(t)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	amount := entities.SatoshiToWei(1000)

	stubRequestPegInDryRun(h, expectedData, amount.AsBigInt(), revert)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(amount, entities.NewWei(0)))
	require.ErrorContains(t, err, wantContains)
	require.NotErrorIs(t, err, blockchain.ErrPegInAlreadyProcessed)
	assert.Empty(t, result.Receipt.TransactionHash)
	assertPayableDryRun(t, h, expectedData, amount.AsBigInt())
	assertRequestPegInNotSent(t, h)
	assertRequestPegInNotEstimated(t, h)
}

func TestPeginContractImpl_RequestPegIn_StatusOneWithoutPegInRequested(t *testing.T) {
	h := newRequestPegInHarness(t)
	amount := entities.SatoshiToWei(1000)
	fee := entities.NewWei(0)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	gasLimit := paddedRequestPegInGas()

	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, amount.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&h.contractMock, h.mockClient, h.signerMock, true)
	stubEstimateRequestPegInGas(h.mockClient)
	stubRequestPegInDryRun(h, expectedData, amount.AsBigInt(), nil)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(amount, fee))
	require.ErrorContains(t, err, "PegInRequested event not found")
	assert.NotEmpty(t, result.Receipt.TransactionHash)
	assert.Equal(t, blockchain.PegInRequestedEvent{}, result.Event)
	assertPayableDryRun(t, h, expectedData, amount.AsBigInt())
}

func TestPeginContractImpl_RequestPegIn_MatchingTopicUnpackError(t *testing.T) {
	h := newRequestPegInHarness(t)
	amount := entities.SatoshiToWei(1000)
	fee := entities.NewWei(0)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	gasLimit := paddedRequestPegInGas()
	parsed, err := commitfirst.PeginCommitFirstContractMetaData.ParseABI()
	require.NoError(t, err)
	badLog := &geth.Log{
		Address: common.HexToAddress(test.AnyRskAddress),
		Topics:  []common.Hash{parsed.Events["PegInRequested"].ID},
		Data:    []byte{0x01},
	}

	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, amount.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&h.contractMock, h.mockClient, h.signerMock, true, badLog)
	stubEstimateRequestPegInGas(h.mockClient)
	stubRequestPegInDryRun(h, expectedData, amount.AsBigInt(), nil)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(amount, fee))
	require.Error(t, err)
	require.NotContains(t, err.Error(), "PegInRequested event not found")
	assert.Equal(t, blockchain.PegInRequestedEvent{}, result.Event)
	assertPayableDryRun(t, h, expectedData, amount.AsBigInt())
}

func TestPeginContractImpl_RequestPegIn_NilAmountOrFee(t *testing.T) {
	t.Run("nil amount", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		guardRejectedRequestPegIn(h)

		result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(nil, entities.NewWei(0)))
		require.ErrorIs(t, err, blockchain.ErrIncorrectFronting)
		assert.Empty(t, result.Receipt.TransactionHash)
		assertNoDryRun(t, h)
		assertRequestPegInNotSent(t, h)
		assertRequestPegInNotEstimated(t, h)
	})
	t.Run("nil fee", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		guardRejectedRequestPegIn(h)

		result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(entities.SatoshiToWei(1000), nil))
		require.ErrorIs(t, err, blockchain.ErrIncorrectFronting)
		assert.Empty(t, result.Receipt.TransactionHash)
		assertNoDryRun(t, h)
		assertRequestPegInNotSent(t, h)
		assertRequestPegInNotEstimated(t, h)
	})
}

func TestPeginContractImpl_RequestPegIn_InvalidAddress(t *testing.T) {
	h := newRequestPegInHarness(t)
	params := sampleRequestPegInParams(entities.SatoshiToWei(1000), entities.NewWei(0))
	params.RskAddress = "not-an-address"
	guardRejectedRequestPegIn(h)

	result, err := h.pegin.RequestPegIn(params)
	require.ErrorIs(t, err, blockchain.InvalidAddressError)
	assert.Empty(t, result.Receipt.TransactionHash)
	assertNoDryRun(t, h)
	assertRequestPegInNotSent(t, h)
	assertRequestPegInNotEstimated(t, h)
}

func TestPeginContractImpl_RequestPegIn_SendError(t *testing.T) {
	h := newRequestPegInHarness(t)
	amount := entities.SatoshiToWei(1000)
	fee := entities.NewWei(0)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	gasLimit := paddedRequestPegInGas()

	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, amount.AsBigInt(), expectedData),
	).Return(assert.AnError).Once()
	prepareTxMocks(&h.contractMock, h.mockClient, h.signerMock, true)
	stubEstimateRequestPegInGas(h.mockClient)
	stubRequestPegInDryRun(h, expectedData, amount.AsBigInt(), nil)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(amount, fee))
	require.ErrorContains(t, err, "request pegin error")
	assert.Empty(t, result.Receipt.TransactionHash)
	assert.Equal(t, blockchain.PegInRequestedEvent{}, result.Event)
	assertPayableDryRun(t, h, expectedData, amount.AsBigInt())
}

// A receipt wait failure does not cancel a broadcast transaction. Returning its
// hash lets the caller recover the result instead of broadcasting a duplicate.
func TestPeginContractImpl_RequestPegIn_KeepsHashWhenReceiptWaitFails(t *testing.T) {
	h := newRequestPegInHarness(t)
	amount := entities.SatoshiToWei(1000)
	fee := entities.NewWei(0)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)

	var sent *geth.Transaction
	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), paddedRequestPegInGas(), amount.AsBigInt(), expectedData),
	).RunAndReturn(func(_ context.Context, tx *geth.Transaction) error {
		sent = tx
		return nil
	}).Once()
	stubRequestPegInSigner(h)
	h.mockClient.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).Return(nil, assert.AnError).Maybe()
	stubEstimateRequestPegInGas(h.mockClient)
	stubRequestPegInDryRun(h, expectedData, amount.AsBigInt(), nil)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(amount, fee))
	require.ErrorContains(t, err, "request pegin error")
	require.NotNil(t, sent)
	assert.Equal(t, sent.Hash().String(), result.Receipt.TransactionHash)
	assert.Equal(t, blockchain.PegInRequestedEvent{}, result.Event)
	h.contractMock.transactor.AssertNumberOfCalls(t, "SendTransaction", 1)
}

func TestPeginContractImpl_RequestPegIn_EstimateGasFailure(t *testing.T) {
	h := newRequestPegInHarness(t)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	amount := entities.SatoshiToWei(1000)

	h.mockClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(0), assert.AnError).Once()
	stubRequestPegInDryRun(h, expectedData, amount.AsBigInt(), nil)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(amount, entities.NewWei(0)))
	require.Error(t, err)
	assert.Empty(t, result.Receipt.TransactionHash)
	assertPayableDryRun(t, h, expectedData, amount.AsBigInt())
	assertRequestPegInNotSent(t, h)
}

func TestPeginContractImpl_EstimateRequestPegInGas(t *testing.T) {
	h := newRequestPegInHarness(t)
	h.signerMock.On("Address").Return(parsedAddress)
	h.mockClient.On("EstimateGas", mock.Anything, mock.Anything).Return(requestPegInEstimatedGas, nil).Once()

	gas, err := h.pegin.EstimateRequestPegInGas(sampleRequestPegInParams(entities.SatoshiToWei(1000), entities.NewWei(0)))
	require.NoError(t, err)
	assert.Equal(t, paddedRequestPegInGas(), gas)
	assertRequestPegInNotSent(t, h)
	assertNoDryRun(t, h)
}

func TestPeginContractImpl_SimulateRequestPegIn_RejectsWithoutSending(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*blockchain.RequestPegInParams)
		want   error
	}{
		{
			name:   "invalid address",
			mutate: func(params *blockchain.RequestPegInParams) { params.RskAddress = "not-an-address" },
			want:   blockchain.InvalidAddressError,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newRequestPegInHarness(t)
			params := sampleRequestPegInParams(entities.SatoshiToWei(1000), entities.NewWei(0))
			tc.mutate(&params)
			guardRejectedRequestPegIn(h)

			err := h.pegin.SimulateRequestPegIn(params)
			require.ErrorIs(t, err, tc.want)
			assertNoDryRun(t, h)
			assertRequestPegInNotSent(t, h)
			assertRequestPegInNotEstimated(t, h)
		})
	}
}

// TestPeginContractImpl_SimulateRequestPegIn_RejectsIncorrectFronting states that
// the simulation must apply the same fronting rule as RequestPegIn. A dry run
// with a value the provider cannot front tells the caller nothing useful, so the
// call must be rejected before the dry run.
func TestPeginContractImpl_SimulateRequestPegIn_RejectsIncorrectFronting(t *testing.T) {
	cases := []struct {
		name   string
		amount *entities.Wei
		fee    *entities.Wei
	}{
		{name: "fee above amount", amount: entities.NewWei(1), fee: entities.NewWei(2)},
		{name: "nil amount", amount: nil, fee: entities.NewWei(0)},
		{name: "nil fee", amount: entities.SatoshiToWei(1000), fee: nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newRequestPegInHarness(t)
			guardRejectedRequestPegIn(h)

			err := h.pegin.SimulateRequestPegIn(sampleRequestPegInParams(tc.amount, tc.fee))
			require.ErrorIs(t, err, blockchain.ErrIncorrectFronting)
			assertNoDryRun(t, h)
			assertRequestPegInNotSent(t, h)
			assertRequestPegInNotEstimated(t, h)
		})
	}
}

func TestPeginContractImpl_SimulateRequestPegIn_DoesNotSend(t *testing.T) {
	h := newRequestPegInHarness(t)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	amount := entities.SatoshiToWei(1000)

	stubRequestPegInDryRun(h, expectedData, amount.AsBigInt(), nil)

	err := h.pegin.SimulateRequestPegIn(sampleRequestPegInParams(amount, entities.NewWei(0)))
	require.NoError(t, err)
	assertPayableDryRun(t, h, expectedData, amount.AsBigInt())
	assertRequestPegInNotSent(t, h)
	h.mockClient.AssertNotCalled(t, "TransactionReceipt", mock.Anything, mock.Anything)
}

// TestPeginContractImpl_SimulateRequestPegIn_DryRunMatchesRequestPegIn states that
// the simulation must dry-run the exact message RequestPegIn would send. If the
// two differ in From, To, or Value, the simulation can accept a request that RequestPegIn
// later reverts on, or reject one that would succeed.
func TestPeginContractImpl_SimulateRequestPegIn_DryRunMatchesRequestPegIn(t *testing.T) {
	amount := entities.SatoshiToWei(1000)
	fee := entities.SatoshiToWei(100)
	expectedValue := new(entities.Wei).Sub(amount, fee)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	params := sampleRequestPegInParams(amount, fee)

	simulate := newRequestPegInHarness(t)
	simulateCall := stubRequestPegInDryRun(simulate, expectedData, expectedValue.AsBigInt(), nil)

	require.NoError(t, simulate.pegin.SimulateRequestPegIn(params))
	assertPayableDryRun(t, simulate, expectedData, expectedValue.AsBigInt())
	assertRequestPegInNotSent(t, simulate)

	request := newRequestPegInHarness(t)
	eventLog := mustPegInRequestedLog(t, pegInRequestedFixture{
		pegInID:     [32]byte{0x05},
		claimer:     parsedAddress,
		rskAddress:  parsedAddress,
		amount:      amount,
		netToUser:   expectedValue,
		callSuccess: true,
	})
	request.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(request.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), paddedRequestPegInGas(), expectedValue.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&request.contractMock, request.mockClient, request.signerMock, true, eventLog)
	stubEstimateRequestPegInGas(request.mockClient)
	requestCall := stubRequestPegInDryRun(request, expectedData, expectedValue.AsBigInt(), nil)

	_, err := request.pegin.RequestPegIn(params)
	require.NoError(t, err)
	assertPayableDryRun(t, request, expectedData, expectedValue.AsBigInt())

	assert.Equal(t, requestCall.From, simulateCall.From)
	require.NotNil(t, simulateCall.To)
	require.NotNil(t, requestCall.To)
	assert.Equal(t, common.HexToAddress(test.AnyRskAddress), *simulateCall.To)
	assert.Equal(t, *requestCall.To, *simulateCall.To)
	require.NotNil(t, simulateCall.Value)
	require.NotNil(t, requestCall.Value)
	assert.Equal(t, 0, requestCall.Value.Cmp(simulateCall.Value))
	assert.Equal(t, expectedValue.AsBigInt(), simulateCall.Value)
	assert.Equal(t, parsedAddress, simulateCall.From)
}

func TestPeginContractImpl_UnpackPegInRequested(t *testing.T) {
	claimer := common.HexToAddress("0x00000000000000000000000000000000000000aa")
	rskAddress := common.HexToAddress("0x00000000000000000000000000000000000000bb")
	amount := entities.SatoshiToWei(1000)
	netToUser := entities.SatoshiToWei(900)
	cases := []struct {
		name        string
		pegInID     [32]byte
		callSuccess bool
	}{
		{name: "call succeeded", pegInID: [32]byte{0x11, 0x22}, callSuccess: true},
		{name: "call failed", pegInID: [32]byte{0x33, 0x44}, callSuccess: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newRequestPegInHarness(t)
			fixture := pegInRequestedFixture{
				pegInID:     tc.pegInID,
				claimer:     claimer,
				rskAddress:  rskAddress,
				amount:      amount,
				netToUser:   netToUser,
				callSuccess: tc.callSuccess,
			}
			eventLog := mustPegInRequestedLog(t, fixture)

			event, err := h.pegin.UnpackPegInRequested(blockchain.TransactionReceipt{
				Status: blockchain.SuccessfulTxStatus,
				Logs:   []blockchain.TransactionLog{transactionLogFromGeth(eventLog)},
			})
			require.NoError(t, err)
			assertPegInRequestedEvent(t, fixture, event)
			assertRequestPegInNotSent(t, h)
			assertNoDryRun(t, h)
		})
	}
}

func TestPeginContractImpl_UnpackPegInRequested_SkipsForeignLog(t *testing.T) {
	h := newRequestPegInHarness(t)
	claimer := common.HexToAddress("0x00000000000000000000000000000000000000aa")
	rskAddress := common.HexToAddress("0x00000000000000000000000000000000000000bb")
	amount := entities.SatoshiToWei(1000)
	netToUser := entities.SatoshiToWei(900)
	foreign := mustPegInRequestedLog(t, pegInRequestedFixture{
		pegInID:     [32]byte{0x11},
		claimer:     claimer,
		rskAddress:  rskAddress,
		amount:      amount,
		netToUser:   netToUser,
		callSuccess: true,
	})
	foreign.Address = common.HexToAddress("0x00000000000000000000000000000000000000aa")
	real := pegInRequestedFixture{
		pegInID:     [32]byte{0x22},
		claimer:     claimer,
		rskAddress:  rskAddress,
		amount:      amount,
		netToUser:   netToUser,
		callSuccess: false,
	}

	event, err := h.pegin.UnpackPegInRequested(blockchain.TransactionReceipt{
		Logs: []blockchain.TransactionLog{
			transactionLogFromGeth(foreign),
			transactionLogFromGeth(mustPegInRequestedLog(t, real)),
		},
	})

	require.NoError(t, err)
	assertPegInRequestedEvent(t, real, event)
}

func TestPeginContractImpl_UnpackPegInRequested_RejectsForeignLog(t *testing.T) {
	h := newRequestPegInHarness(t)
	foreign := mustPegInRequestedLog(t, pegInRequestedFixture{
		pegInID:     [32]byte{0x11},
		claimer:     common.HexToAddress("0x00000000000000000000000000000000000000aa"),
		rskAddress:  common.HexToAddress("0x00000000000000000000000000000000000000bb"),
		amount:      entities.NewWei(1000),
		netToUser:   entities.NewWei(900),
		callSuccess: true,
	})
	foreign.Address = common.HexToAddress("0x00000000000000000000000000000000000000aa")

	_, err := h.pegin.UnpackPegInRequested(blockchain.TransactionReceipt{
		Logs: []blockchain.TransactionLog{transactionLogFromGeth(foreign)},
	})

	require.ErrorContains(t, err, "request pegin error: PegInRequested event not found")
}

func TestPeginContractImpl_UnpackPegInRequested_RejectsEmptyData(t *testing.T) {
	h := newRequestPegInHarness(t)
	empty := mustPegInRequestedLog(t, pegInRequestedFixture{
		pegInID:     [32]byte{0x11},
		claimer:     common.HexToAddress("0x00000000000000000000000000000000000000aa"),
		rskAddress:  common.HexToAddress("0x00000000000000000000000000000000000000bb"),
		amount:      entities.NewWei(1000),
		netToUser:   entities.NewWei(900),
		callSuccess: true,
	})
	empty.Data = nil

	_, err := h.pegin.UnpackPegInRequested(blockchain.TransactionReceipt{
		Logs: []blockchain.TransactionLog{transactionLogFromGeth(empty)},
	})

	require.ErrorContains(t, err, "request pegin error: PegInRequested event not found")
}

func TestPeginContractImpl_UnpackPegInRequested_SkipsEmptyDataBeforeValidEvent(t *testing.T) {
	h := newRequestPegInHarness(t)
	claimer := common.HexToAddress("0x00000000000000000000000000000000000000aa")
	rskAddress := common.HexToAddress("0x00000000000000000000000000000000000000bb")
	amount := entities.SatoshiToWei(1000)
	netToUser := entities.SatoshiToWei(900)
	empty := mustPegInRequestedLog(t, pegInRequestedFixture{
		pegInID:     [32]byte{0x11},
		claimer:     claimer,
		rskAddress:  rskAddress,
		amount:      amount,
		netToUser:   netToUser,
		callSuccess: true,
	})
	empty.Data = nil
	valid := pegInRequestedFixture{
		pegInID:     [32]byte{0x22},
		claimer:     claimer,
		rskAddress:  rskAddress,
		amount:      amount,
		netToUser:   netToUser,
		callSuccess: false,
	}

	event, err := h.pegin.UnpackPegInRequested(blockchain.TransactionReceipt{
		Logs: []blockchain.TransactionLog{
			transactionLogFromGeth(empty),
			transactionLogFromGeth(mustPegInRequestedLog(t, valid)),
		},
	})

	require.NoError(t, err)
	assertPegInRequestedEvent(t, valid, event)
}

func TestPeginContractImpl_UnpackPegInRequested_SkipsRemovedLog(t *testing.T) {
	h := newRequestPegInHarness(t)
	claimer := common.HexToAddress("0x00000000000000000000000000000000000000aa")
	rskAddress := common.HexToAddress("0x00000000000000000000000000000000000000bb")
	amount := entities.SatoshiToWei(1000)
	netToUser := entities.SatoshiToWei(900)
	removed := mustPegInRequestedLog(t, pegInRequestedFixture{
		pegInID:     [32]byte{0x11},
		claimer:     claimer,
		rskAddress:  rskAddress,
		amount:      amount,
		netToUser:   netToUser,
		callSuccess: true,
	})
	removed.Removed = true
	live := pegInRequestedFixture{
		pegInID:     [32]byte{0x22},
		claimer:     claimer,
		rskAddress:  rskAddress,
		amount:      amount,
		netToUser:   netToUser,
		callSuccess: true,
	}

	event, err := h.pegin.UnpackPegInRequested(blockchain.TransactionReceipt{
		Logs: []blockchain.TransactionLog{
			transactionLogFromGeth(removed),
			transactionLogFromGeth(mustPegInRequestedLog(t, live)),
		},
	})

	require.NoError(t, err)
	assertPegInRequestedEvent(t, live, event)
}

func TestPeginContractImpl_UnpackPegInRequested_AcceptsMixedCaseAddress(t *testing.T) {
	h := newRequestPegInHarness(t)
	fixture := pegInRequestedFixture{
		pegInID:     [32]byte{0x11},
		claimer:     common.HexToAddress("0x00000000000000000000000000000000000000aa"),
		rskAddress:  common.HexToAddress("0x00000000000000000000000000000000000000bb"),
		amount:      entities.SatoshiToWei(1000),
		netToUser:   entities.SatoshiToWei(900),
		callSuccess: true,
	}
	eventLog := transactionLogFromGeth(mustPegInRequestedLog(t, fixture))
	eventLog.Address = strings.ToLower(test.AnyRskAddress)

	event, err := h.pegin.UnpackPegInRequested(blockchain.TransactionReceipt{
		Logs: []blockchain.TransactionLog{eventLog},
	})

	require.NoError(t, err)
	assertPegInRequestedEvent(t, fixture, event)
}
