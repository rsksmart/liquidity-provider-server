package rootstock_test

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	bindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegout"
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

var pegoutQuote = quote.PegoutQuote{
	LbcAddress:            "0x7C4890A0f1D4bBf2C669Ac2d1efFa185c505359b",
	LpRskAddress:          "0x9D93929A9099be4355fC2389FbF253982F9dF47c",
	BtcRefundAddress:      "mjaGtyj74LYn7gApr17prZxDPDnfuUnRa5",
	RskRefundAddress:      "0x0d8A0F1ef26B4b9650d98E1c22c560327cF387FE",
	LpBtcAddress:          "mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6",
	CallFee:               entities.NewWei(1),
	PenaltyFee:            entities.NewWei(2),
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
	ProductFeeAmount:      entities.NewWei(789),
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
	DepositAddress:        []byte{111, 211, 208, 46, 217, 53, 0, 241, 222, 67, 83, 58, 58, 11, 192, 246, 63, 2, 142, 48, 246},
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

var deposits = []geth.Log{
	{
		TxHash:      common.Hash{7},
		BlockNumber: 10,
		Topics: []common.Hash{
			common.HexToHash("0xb1bc7bfc0dab19777eb03aa0a5643378fc9f186c8fc5a36620d21136fbea570f"),
			common.HexToHash("0x0102030000000000000000000000000000000000000000000000000000000000"),
			common.HexToHash("0x0000000000000000000000000100000000000000000000000000000000000000"),
			common.HexToHash("0x00000000000000000000000000000000000000000000000000000000075bcd15"),
		},
		Data: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000022b"),
	},
	{
		TxHash:      common.Hash{8},
		BlockNumber: 11,
		Topics: []common.Hash{
			common.HexToHash("0xb1bc7bfc0dab19777eb03aa0a5643378fc9f186c8fc5a36620d21136fbea570f"),
			common.HexToHash("0x0405060000000000000000000000000000000000000000000000000000000000"),
			common.HexToHash("0x0000000000000000000000000200000000000000000000000000000000000000"),
			common.HexToHash("0x000000000000000000000000000000000000000000000000000000003ade68b1"),
		},
		Data: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000029a"),
	},
	{
		TxHash:      common.Hash{9},
		BlockNumber: 12,
		Topics: []common.Hash{
			common.HexToHash("0xb1bc7bfc0dab19777eb03aa0a5643378fc9f186c8fc5a36620d21136fbea570f"),
			common.HexToHash("0x0708090000000000000000000000000000000000000000000000000000000000"),
			common.HexToHash("0x0000000000000000000000000300000000000000000000000000000000000000"),
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000006a11e3d"),
		},
		Data: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000309"),
	},
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

func TestNewPegoutContractImpl(t *testing.T) {
	contract := rootstock.NewPegoutContractImpl(
		rootstock.NewRskClient(&mocks.RpcClientBindingMock{}),
		test.AnyAddress,
		createBoundContractMock().contract,
		&mocks.TransactionSignerMock{},
		rootstock.RetryParams{Retries: 1, Sleep: 1},
		time.Duration(1),
		bindings.NewPegoutContract(),
		Abis,
	)
	test.AssertNonZeroValues(t, contract)
}

func TestPegoutContractImpl_GetAddress(t *testing.T) {
	pegoutContract := rootstock.NewPegoutContractImpl(dummyClient, test.AnyAddress, nil, nil, rootstock.RetryParams{}, time.Duration(1), nil, Abis)
	assert.Equal(t, test.AnyAddress, pegoutContract.GetAddress())
}

func TestPegoutContractImpl_HashPegoutQuote(t *testing.T) {
	contractMock := createBoundContractMock()
	pegoutBinding := bindings.NewPegoutContract()
	hash := [32]byte{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	pegoutContract := rootstock.NewPegoutContractImpl(
		dummyClient,
		test.AnyAddress,
		contractMock.contract,
		nil, rootstock.RetryParams{},
		time.Duration(1),
		pegoutBinding,
		Abis,
	)
	t.Run("Success", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(pegoutBinding.PackHashPegOutQuote(parsedPegoutQuote)),
			mock.Anything,
		).Return(mustPackBytes32(t, hash), nil).Once()
		result, err := pegoutContract.HashPegoutQuote(pegoutQuote)
		require.NoError(t, err)
		assert.Equal(t, hex.EncodeToString(hash[:]), result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on HashPegoutQuote fail", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(pegoutBinding.PackHashPegOutQuote(parsedPegoutQuote)),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := pegoutContract.HashPegoutQuote(pegoutQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestPegoutContractImpl_HashPegoutQuote_ParsingErrors(t *testing.T) {
	contractMock := createBoundContractMock()
	pegoutBinding := bindings.NewPegoutContract()
	pegoutContract := rootstock.NewPegoutContractImpl(dummyClient, test.AnyAddress, contractMock.contract, nil, rootstock.RetryParams{}, time.Duration(1), pegoutBinding, Abis)
	t.Run("Incomplete quote", func(t *testing.T) {
		testQuote := pegoutQuote
		testQuote.LbcAddress = ""
		result, err := pegoutContract.HashPegoutQuote(testQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
	t.Run("Invalid lbc address", func(t *testing.T) {
		testQuote := pegoutQuote
		testQuote.LbcAddress = test.AnyString
		result, err := pegoutContract.HashPegoutQuote(testQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
	t.Run("Invalid lp rsk address", func(t *testing.T) {
		testQuote := pegoutQuote
		testQuote.LpRskAddress = test.AnyString
		result, err := pegoutContract.HashPegoutQuote(testQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
	t.Run("Invalid rsk refund address", func(t *testing.T) {
		testQuote := pegoutQuote
		testQuote.RskRefundAddress = test.AnyString
		result, err := pegoutContract.HashPegoutQuote(testQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
	t.Run("Invalid btc refund address", func(t *testing.T) {
		testQuote := pegoutQuote
		testQuote.BtcRefundAddress = test.AnyString
		result, err := pegoutContract.HashPegoutQuote(testQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
	t.Run("Invalid lp btc address", func(t *testing.T) {
		testQuote := pegoutQuote
		testQuote.LpBtcAddress = test.AnyString
		result, err := pegoutContract.HashPegoutQuote(testQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
	t.Run("Invalid destination address", func(t *testing.T) {
		testQuote := pegoutQuote
		testQuote.DepositAddress = test.AnyString
		result, err := pegoutContract.HashPegoutQuote(testQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestPegoutContractImpl_IsPegOutQuoteCompleted(t *testing.T) {
	const quoteHash = "762d73db7e80d845dae50d6ddda4d64d59f99352ead28afd51610e5674b08c0a"
	parsedQuoteHash := [32]byte{0x76, 0x2d, 0x73, 0xdb, 0x7e, 0x80, 0xd8, 0x45, 0xda, 0xe5, 0xd, 0x6d, 0xdd, 0xa4, 0xd6, 0x4d, 0x59, 0xf9, 0x93, 0x52, 0xea, 0xd2, 0x8a, 0xfd, 0x51, 0x61, 0xe, 0x56, 0x74, 0xb0, 0x8c, 0xa}
	contractMock := createBoundContractMock()
	pegoutBinding := bindings.NewPegoutContract()
	pegoutContract := rootstock.NewPegoutContractImpl(dummyClient, test.AnyAddress, contractMock.contract, nil, rootstock.RetryParams{}, time.Duration(1), pegoutBinding, Abis)
	t.Run("Success", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(pegoutBinding.PackIsQuoteCompleted(parsedQuoteHash)),
			mock.Anything,
		).Return(mustPackBool(t, true), nil).Once()
		result, err := pegoutContract.IsPegOutQuoteCompleted(quoteHash)
		require.NoError(t, err)
		assert.True(t, result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Should handle call error", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(pegoutBinding.PackIsQuoteCompleted(parsedQuoteHash)),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := pegoutContract.IsPegOutQuoteCompleted(quoteHash)
		require.Error(t, err)
		assert.False(t, result)
	})
	t.Run("Should handle parse error", func(t *testing.T) {
		result, err := pegoutContract.IsPegOutQuoteCompleted(test.AnyString)
		require.Error(t, err)
		assert.False(t, result)
	})
	t.Run("Should return error when quote hash is not long enough", func(t *testing.T) {
		result, err := pegoutContract.IsPegOutQuoteCompleted("0104050302")
		require.ErrorContains(t, err, "quote hash must be 32 bytes long")
		assert.False(t, result)
	})
}

func TestPegoutContractImpl_DaoFeePercentage(t *testing.T) {
	contractMock := createBoundContractMock()
	pegoutBinding := bindings.NewPegoutContract()
	t.Run("Success", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(pegoutBinding.PackGetFeePercentage()),
			mock.Anything,
		).Return(mustPackUint256(t, big.NewInt(1)), nil).Once()
		pegoutContract := rootstock.NewPegoutContractImpl(dummyClient, test.AnyAddress, contractMock.contract, nil, rootstock.RetryParams{}, time.Duration(1), pegoutBinding, Abis)
		percentage, err := pegoutContract.DaoFeePercentage()
		require.NoError(t, err)
		require.Equal(t, uint64(1), percentage)
	})
	t.Run("Error handling on ProductFeePercentage call fail", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(pegoutBinding.PackGetFeePercentage()),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		pegoutContract := rootstock.NewPegoutContractImpl(dummyClient, test.AnyAddress, contractMock.contract, nil, rootstock.RetryParams{}, time.Duration(1), pegoutBinding, Abis)
		percentage, err := pegoutContract.DaoFeePercentage()
		require.Error(t, err)
		require.Zero(t, percentage)
	})
}

func TestPegoutContractImpl_RefundUserPegOut(t *testing.T) {
	contractMock := createBoundContractMock()
	pegoutBinding := bindings.NewPegoutContract()
	signer := &mocks.TransactionSignerMock{}
	client := &mocks.RpcClientBindingMock{}
	pegoutContract := rootstock.NewPegoutContractImpl(
		rootstock.NewRskClient(client),
		test.AnyAddress,
		contractMock.contract,
		signer,
		rootstock.RetryParams{},
		time.Duration(1),
		pegoutBinding,
		Abis,
	)

	t.Run("should fail with invalid hash format", func(t *testing.T) {
		result, err := pegoutContract.RefundUserPegOut("invalid hash")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid quote hash format")
		assert.Empty(t, result)
	})

	t.Run("should fail with invalid hash length", func(t *testing.T) {
		result, err := pegoutContract.RefundUserPegOut("ab")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "quote hash must be 32 bytes long")
		assert.Empty(t, result)
	})

	t.Run("should fail if transaction fails", func(t *testing.T) {
		validHash := strings.Repeat("aa", 32)
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 0, big.NewInt(0), pegoutBinding.PackRefundUserPegOut(test.MustEncode32Bytes(validHash))),
		).Return(assert.AnError).Once()
		signer.EXPECT().Sign(mock.Anything, mock.Anything).RunAndReturn(func(addr common.Address, tx *geth.Transaction) (*geth.Transaction, error) {
			return tx, nil
		})
		prepareTxMocks(&contractMock, client, signer, true)
		result, err := pegoutContract.RefundUserPegOut(validHash)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "refund user peg out error")
		assert.Empty(t, result)
	})

	t.Run("should succeed", func(t *testing.T) {
		validHash := strings.Repeat("aa", 32)
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 0, big.NewInt(0), pegoutBinding.PackRefundUserPegOut(test.MustEncode32Bytes(validHash))),
		).Return(nil).Once()
		prepareTxMocks(&contractMock, client, signer, true)
		result, err := pegoutContract.RefundUserPegOut(validHash)
		require.NoError(t, err)
		assert.Equal(t, "0x"+test.AnyHash, result)
		contractMock.transactor.AssertExpectations(t)
	})
}

// nolint:funlen
func TestPegoutContractImpl_RefundPegout(t *testing.T) {
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
	txConfig := blockchain.TransactionConfig{Value: entities.NewWei(0), GasLimit: &gasLimit}
	pegoutBinding := bindings.NewPegoutContract()
	txData := pegoutBinding.PackRefundPegOut(
		refundParams.QuoteHash, refundParams.BtcRawTx, refundParams.BtcBlockHeaderHash,
		refundParams.MerkleBranchPath, refundParams.MerkleBranchHashes,
	)
	t.Run("Success", func(t *testing.T) {
		contractMock := createBoundContractMock()
		pegoutContract := rootstock.NewPegoutContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractMock.contract, signerMock, rootstock.RetryParams{}, time.Duration(1), pegoutBinding, Abis)
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(txData),
			mock.Anything,
		).Return(nil, nil).Once()
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 500, big.NewInt(0), txData),
		).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, true)
		result, err := pegoutContract.RefundPegout(txConfig, refundParams)
		require.NoError(t, err)
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
			Logs:              []blockchain.TransactionLog{},
		}
		assert.Equal(t, expectedReceipt, result)
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertExpectations(t)
	})
	t.Run("Error handling (waiting for bridge, not enough confirmations)", func(t *testing.T) {
		contractMock := createBoundContractMock()
		pegoutContract := rootstock.NewPegoutContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractMock.contract, signerMock, rootstock.RetryParams{}, time.Duration(1), pegoutBinding, Abis)
		e := NewRskRpcError("transaction reverted", "0xd2506f8c00000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000000000000000002")
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(txData),
			mock.Anything,
		).Return(nil, e).Once()
		result, err := pegoutContract.RefundPegout(txConfig, refundParams)
		require.ErrorIs(t, err, blockchain.WaitingForBridgeError)
		assert.Empty(t, result)
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertNotCalled(t, "SendTransaction")
	})
	t.Run("Error handling (waiting for bridge, tx not seen yet)", func(t *testing.T) {
		contractMock := createBoundContractMock()
		pegoutContract := rootstock.NewPegoutContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractMock.contract, signerMock, rootstock.RetryParams{}, time.Duration(1), pegoutBinding, Abis)
		e := NewRskRpcError("transaction reverted", "0xd06e366affffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(txData),
			mock.Anything,
		).Return(nil, e).Once()
		result, err := pegoutContract.RefundPegout(txConfig, refundParams)
		require.ErrorIs(t, err, blockchain.WaitingForBridgeError)
		assert.Empty(t, result)
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertNotCalled(t, "SendTransaction")
	})
	t.Run("Error handling (Call error)", func(t *testing.T) {
		contractMock := createBoundContractMock()
		pegoutContract := rootstock.NewPegoutContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractMock.contract, signerMock, rootstock.RetryParams{}, time.Duration(1), pegoutBinding, Abis)
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(txData),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := pegoutContract.RefundPegout(txConfig, refundParams)
		require.Error(t, err)
		assert.Empty(t, result)
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertNotCalled(t, "SendTransaction")
	})
	t.Run("Error handling (Transaction send error)", func(t *testing.T) {
		contractMock := createBoundContractMock()
		pegoutContract := rootstock.NewPegoutContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractMock.contract, signerMock, rootstock.RetryParams{}, time.Duration(1), pegoutBinding, Abis)
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(txData),
			mock.Anything,
		).Return(nil, nil).Once()
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 500, big.NewInt(0), txData),
		).Return(assert.AnError).Once()
		signerMock.EXPECT().Sign(mock.Anything, mock.Anything).RunAndReturn(func(addr common.Address, tx *geth.Transaction) (*geth.Transaction, error) {
			return tx, nil
		})
		prepareTxMocks(&contractMock, mockClient, signerMock, true)
		result, err := pegoutContract.RefundPegout(txConfig, refundParams)
		require.ErrorContains(t, err, "refund pegout error")
		assert.Empty(t, result)
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertExpectations(t)
	})
	t.Run("Error handling (Transaction reverted)", func(t *testing.T) {
		contractMock := createBoundContractMock()
		pegoutContract := rootstock.NewPegoutContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractMock.contract, signerMock, rootstock.RetryParams{}, time.Duration(1), pegoutBinding, Abis)
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(txData),
			mock.Anything,
		).Return(nil, nil).Once()
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 500, big.NewInt(0), txData),
		).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, false)
		result, err := pegoutContract.RefundPegout(txConfig, refundParams)
		require.ErrorContains(t, err, "refund pegout error: transaction reverted")
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
			Logs:              []blockchain.TransactionLog{},
		}
		assert.Equal(t, expectedReceipt, result)
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertExpectations(t)
	})
}

func TestPegoutContractImpl_GetDepositEvents(t *testing.T) {
	contractMock := createBoundContractMock()
	pegoutBinding := bindings.NewPegoutContract()
	pegoutContract := rootstock.NewPegoutContractImpl(dummyClient, test.AnyAddress, contractMock.contract, nil, rootstock.RetryParams{}, time.Duration(1), pegoutBinding, Abis)
	t.Run("Success", func(t *testing.T) {
		var from uint64 = 500
		var to uint64 = 1000
		contractMock.filterer.EXPECT().FilterLogs(mock.Anything, mock.MatchedBy(filterMatchFunc(from, to))).Return(deposits, nil).Once()
		result, err := pegoutContract.GetDepositEvents(context.Background(), from, &to)
		require.NoError(t, err)
		assert.Equal(t, parsedDeposits, result)
		contractMock.filterer.AssertExpectations(t)
	})
	t.Run("Error handling when failed to get iterator", func(t *testing.T) {
		var from uint64 = 600
		var to uint64 = 1100
		contractMock.filterer.EXPECT().FilterLogs(mock.Anything, mock.MatchedBy(filterMatchFunc(from, to))).Return(nil, assert.AnError).Once()
		result, err := pegoutContract.GetDepositEvents(context.Background(), from, &to)
		require.Error(t, err)
		assert.Nil(t, result)
		contractMock.filterer.AssertExpectations(t)
	})
}

// nolint:funlen
func TestPegoutContractImpl_ValidatePegout(t *testing.T) {
	const quoteHash = "762d73db7e80d845dae50d6ddda4d64d59f99352ead28afd51610e5674b08c0a"
	parsedQuoteHash := [32]byte{0x76, 0x2d, 0x73, 0xdb, 0x7e, 0x80, 0xd8, 0x45, 0xda, 0xe5, 0xd, 0x6d, 0xdd, 0xa4, 0xd6, 0x4d, 0x59, 0xf9, 0x93, 0x52, 0xea, 0xd2, 0x8a, 0xfd, 0x51, 0x61, 0xe, 0x56, 0x74, 0xb0, 0x8c, 0xa}
	btcTx := []byte{0x01, 0x00, 0x00, 0x00}
	contractMock := createBoundContractMock()
	pegoutBinding := bindings.NewPegoutContract()
	signerMock := &mocks.TransactionSignerMock{}
	pegoutContract := rootstock.NewPegoutContractImpl(dummyClient, test.AnyAddress, contractMock.contract, signerMock, rootstock.RetryParams{}, time.Duration(1), pegoutBinding, Abis)

	t.Run("Success", func(t *testing.T) {
		signerMock.On("Address").Return(parsedAddress).Once()
		expectedQuote := bindings.QuotesPegOutQuote{
			CallFee:               big.NewInt(1),
			PenaltyFee:            big.NewInt(1),
			Value:                 big.NewInt(1),
			ProductFeeAmount:      big.NewInt(1),
			GasFee:                big.NewInt(1),
			LbcAddress:            common.HexToAddress("10"),
			LpRskAddress:          common.HexToAddress("10"),
			RskRefundAddress:      common.HexToAddress("10"),
			Nonce:                 1,
			AgreementTimestamp:    2,
			DepositDateLimit:      3,
			TransferTime:          4,
			ExpireDate:            5,
			ExpireBlock:           6,
			DepositConfirmations:  7,
			TransferConfirmations: 8,
			DepositAddress:        []byte{1},
			BtcRefundAddress:      []byte{1},
			LpBtcAddress:          []byte{1},
		}
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(pegoutBinding.PackValidatePegout(parsedQuoteHash, btcTx)),
			mock.Anything,
		).Return(mustPackPegoutQuote(t, expectedQuote), nil).Once()
		err := pegoutContract.ValidatePegout(quoteHash, btcTx)
		require.NoError(t, err)
		contractMock.caller.AssertExpectations(t)
		signerMock.AssertExpectations(t)
	})

	t.Run("Should handle parseable contract revert (non-recoverable)", func(t *testing.T) {
		signerMock.On("Address").Return(parsedAddress).Once()
		// Contract revert with parseable error data (MalformedTransaction error selector: 0x7201f86d)
		contractRevertError := NewRskRpcError("execution reverted", "0x7201f86d")
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(pegoutBinding.PackValidatePegout(parsedQuoteHash, btcTx)),
			mock.Anything,
		).Return(nil, contractRevertError).Once()
		err := pegoutContract.ValidatePegout(quoteHash, btcTx)
		require.Error(t, err)
		require.ErrorContains(t, err, "reverted with:")
		require.ErrorContains(t, err, "MalformedTransaction")
		contractMock.caller.AssertExpectations(t)
		signerMock.AssertExpectations(t)
	})

	t.Run("Should handle unparseable error (potentially recoverable)", func(t *testing.T) {
		signerMock.On("Address").Return(parsedAddress).Once()
		// Network/RPC error that can't be parsed (no error data) - potentially recoverable
		networkError := assert.AnError
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(pegoutBinding.PackValidatePegout(parsedQuoteHash, btcTx)),
			mock.Anything,
		).Return(nil, networkError).Once()
		err := pegoutContract.ValidatePegout(quoteHash, btcTx)
		require.Error(t, err)
		require.ErrorContains(t, err, "error validating pegout:")
		require.NotContains(t, err.Error(), "reverted with:")
		contractMock.caller.AssertExpectations(t)
		signerMock.AssertExpectations(t)
	})

	t.Run("Should handle invalid quote hash format", func(t *testing.T) {
		signerMock.On("Address").Return(parsedAddress).Once()
		err := pegoutContract.ValidatePegout("not-a-hex-string", btcTx)
		require.Error(t, err)
		signerMock.AssertExpectations(t)
	})

	t.Run("Should reject quote hash with 0x prefix", func(t *testing.T) {
		signerMock.On("Address").Return(parsedAddress).Once()
		err := pegoutContract.ValidatePegout("0x762d73db7e80d845dae50d6ddda4d64d59f99352ead28afd51610e5674b08c0a", btcTx)
		require.Error(t, err)
		signerMock.AssertExpectations(t)
	})

	t.Run("Should return error when quote hash is not 32 bytes long", func(t *testing.T) {
		signerMock.On("Address").Return(parsedAddress).Once()
		err := pegoutContract.ValidatePegout("0104050302", btcTx)
		require.ErrorContains(t, err, "quote hash must be 32 bytes long")
		signerMock.AssertExpectations(t)
	})
}
