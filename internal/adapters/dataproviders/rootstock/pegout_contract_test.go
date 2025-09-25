package rootstock_test

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
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
	"strings"
	"testing"
	"time"
)

const (
	penalizedIteratorString = "*bindings.ICollateralManagementPenalizedIterator"
	depositIteratorString   = "*bindings.IPegOutPegOutDepositIterator"
	invalidAddressTest      = "Invalid address"
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

var deposits = []*bindings.IPegOutPegOutDeposit{
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

func TestNewPegoutContractImpl(t *testing.T) {
	contract := rootstock.NewPegoutContractImpl(
		rootstock.NewRskClient(&mocks.RpcClientBindingMock{}),
		test.AnyAddress,
		&mocks.PegoutContractAdapterMock{},
		&mocks.TransactionSignerMock{},
		rootstock.RetryParams{Retries: 1, Sleep: 1},
		time.Duration(1),
		Abis,
	)
	test.AssertNonZeroValues(t, contract)
}

func TestPegoutContractImpl_GetAddress(t *testing.T) {
	pegoutContract := rootstock.NewPegoutContractImpl(dummyClient, test.AnyAddress, nil, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
	assert.Equal(t, test.AnyAddress, pegoutContract.GetAddress())
}

func TestPegoutContractImpl_HashPegoutQuote(t *testing.T) {
	contract := &mocks.PegoutContractAdapterMock{}
	hash := [32]byte{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	pegoutContract := rootstock.NewPegoutContractImpl(
		dummyClient,
		test.AnyAddress,
		contract,
		nil, rootstock.RetryParams{},
		time.Duration(1),
		Abis,
	)
	t.Run("Success", func(t *testing.T) {
		contract.EXPECT().HashPegOutQuote(mock.Anything, parsedPegoutQuote).Return(hash, nil).Once()
		result, err := pegoutContract.HashPegoutQuote(pegoutQuote)
		require.NoError(t, err)
		assert.Equal(t, hex.EncodeToString(hash[:]), result)
		contract.AssertExpectations(t)
	})
	t.Run("Error handling on HashPegoutQuote fail", func(t *testing.T) {
		contract.EXPECT().HashPegOutQuote(mock.Anything, parsedPegoutQuote).Return([32]byte{}, assert.AnError).Once()
		result, err := pegoutContract.HashPegoutQuote(pegoutQuote)
		require.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestPegoutContractImpl_HashPegoutQuote_ParsingErrors(t *testing.T) {
	contract := &mocks.PegoutContractAdapterMock{}
	pegoutContract := rootstock.NewPegoutContractImpl(dummyClient, test.AnyAddress, contract, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
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
	contractBinding := &mocks.PegoutContractAdapterMock{}
	pegoutContract := rootstock.NewPegoutContractImpl(dummyClient, test.AnyAddress, contractBinding, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
	t.Run("Success", func(t *testing.T) {
		contractBinding.EXPECT().IsQuoteCompleted(mock.Anything, parsedQuoteHash).Return(true, nil).Once()
		result, err := pegoutContract.IsPegOutQuoteCompleted(quoteHash)
		require.NoError(t, err)
		assert.True(t, result)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Should handle call error", func(t *testing.T) {
		contractBinding.EXPECT().IsQuoteCompleted(mock.Anything, parsedQuoteHash).Return(false, assert.AnError).Once()
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
	contractBinding := &mocks.PegoutContractAdapterMock{}
	t.Run("Success", func(t *testing.T) {
		contractBinding.EXPECT().GetFeePercentage(mock.Anything).Return(big.NewInt(1), nil).Once()
		pegoutContract := rootstock.NewPegoutContractImpl(dummyClient, test.AnyAddress, contractBinding, nil, rootstock.RetryParams{Retries: 0, Sleep: 0}, time.Duration(1), Abis)
		percentage, err := pegoutContract.DaoFeePercentage()
		require.NoError(t, err)
		require.Equal(t, uint64(1), percentage)
	})
	t.Run("Error handling on ProductFeePercentage call fail", func(t *testing.T) {
		contractBinding.EXPECT().GetFeePercentage(mock.Anything).Return(nil, assert.AnError).Once()
		pegoutContract := rootstock.NewPegoutContractImpl(dummyClient, test.AnyAddress, contractBinding, nil, rootstock.RetryParams{Retries: 0, Sleep: 0}, time.Duration(1), Abis)
		percentage, err := pegoutContract.DaoFeePercentage()
		require.Error(t, err)
		require.Zero(t, percentage)
	})
}

func TestPegoutContractImpl_RefundUserPegOut(t *testing.T) {
	contractBinding := &mocks.PegoutContractAdapterMock{}
	signer := &mocks.TransactionSignerMock{}
	client := &mocks.RpcClientBindingMock{}
	pegoutContract := rootstock.NewPegoutContractImpl(
		rootstock.NewRskClient(client),
		test.AnyAddress,
		contractBinding,
		signer,
		rootstock.RetryParams{},
		time.Duration(1),
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
		tx := prepareTxMocks(client, signer, false)
		contractBinding.On("RefundUserPegOut", mock.Anything, mock.Anything).Return(tx, assert.AnError).Once()

		result, err := pegoutContract.RefundUserPegOut(validHash)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "refund user peg out error")
		assert.Empty(t, result)
	})

	t.Run("should succeed", func(t *testing.T) {
		validHash := strings.Repeat("aa", 32)
		tx := prepareTxMocks(client, signer, true)
		contractBinding.On("RefundUserPegOut", mock.Anything, mock.Anything).Return(tx, nil).Once()

		result, err := pegoutContract.RefundUserPegOut(validHash)
		require.NoError(t, err)
		assert.Equal(t, tx.Hash().String(), result)

		contractBinding.AssertExpectations(t)
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
	txConfig := blockchain.TransactionConfig{Value: entities.NewWei(1234), GasLimit: &gasLimit}
	matchOptsFunc := func(opts *bind.TransactOpts) bool {
		return opts.From.String() == parsedAddress.String() && opts.GasLimit == gasLimit
	}
	t.Run("Success", func(t *testing.T) {
		contractBinding := &mocks.PegoutContractAdapterMock{}
		callerMock := &mocks.ContractCallerBindingMock{}
		pegoutContract := rootstock.NewPegoutContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractBinding, signerMock, rootstock.RetryParams{}, time.Duration(1), Abis)
		contractBinding.On("Caller").Return(callerMock).Once()
		callerMock.On("Call", mock.Anything, mock.Anything, "refundPegOut",
			refundParams.QuoteHash, refundParams.BtcRawTx, refundParams.BtcBlockHeaderHash,
			refundParams.MerkleBranchPath, refundParams.MerkleBranchHashes,
		).Return(nil).Once()
		tx := prepareTxMocks(mockClient, signerMock, true)
		contractBinding.On("RefundPegOut", mock.MatchedBy(matchOptsFunc),
			refundParams.QuoteHash, refundParams.BtcRawTx, refundParams.BtcBlockHeaderHash,
			refundParams.MerkleBranchPath, refundParams.MerkleBranchHashes,
		).Return(tx, nil).Once()
		result, err := pegoutContract.RefundPegout(txConfig, refundParams)
		require.NoError(t, err)
		assert.Equal(t, tx.Hash().String(), result)
		contractBinding.AssertExpectations(t)
		callerMock.AssertExpectations(t)
	})
	t.Run("Error handling (waiting for bridge, not enough confirmations)", func(t *testing.T) {
		contractBinding := &mocks.PegoutContractAdapterMock{}
		callerMock := &mocks.ContractCallerBindingMock{}
		pegoutContract := rootstock.NewPegoutContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractBinding, signerMock, rootstock.RetryParams{}, time.Duration(1), Abis)
		e := NewRskRpcError("transaction reverted", "0xd2506f8c00000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000000000000000002")
		contractBinding.EXPECT().Caller().Return(callerMock).Once()
		callerMock.EXPECT().Call(mock.Anything, mock.Anything, "refundPegOut",
			refundParams.QuoteHash, refundParams.BtcRawTx, refundParams.BtcBlockHeaderHash,
			refundParams.MerkleBranchPath, refundParams.MerkleBranchHashes,
		).Return(e).Once()
		result, err := pegoutContract.RefundPegout(txConfig, refundParams)
		require.ErrorIs(t, err, blockchain.WaitingForBridgeError)
		assert.Empty(t, result)
		contractBinding.AssertExpectations(t)
		contractBinding.AssertNotCalled(t, "RefundPegOut")
		callerMock.AssertExpectations(t)
	})
	t.Run("Error handling (waiting for bridge, tx not seen yet)", func(t *testing.T) {
		contractBinding := &mocks.PegoutContractAdapterMock{}
		callerMock := &mocks.ContractCallerBindingMock{}
		pegoutContract := rootstock.NewPegoutContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractBinding, signerMock, rootstock.RetryParams{}, time.Duration(1), Abis)
		e := NewRskRpcError("transaction reverted", "0xd06e366affffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
		contractBinding.EXPECT().Caller().Return(callerMock).Once()
		callerMock.EXPECT().Call(mock.Anything, mock.Anything, "refundPegOut",
			refundParams.QuoteHash, refundParams.BtcRawTx, refundParams.BtcBlockHeaderHash,
			refundParams.MerkleBranchPath, refundParams.MerkleBranchHashes,
		).Return(e).Once()
		result, err := pegoutContract.RefundPegout(txConfig, refundParams)
		require.ErrorIs(t, err, blockchain.WaitingForBridgeError)
		assert.Empty(t, result)
		contractBinding.AssertExpectations(t)
		contractBinding.AssertNotCalled(t, "RefundPegOut")
		callerMock.AssertExpectations(t)
	})
	t.Run("Error handling (Call error)", func(t *testing.T) {
		contractBinding := &mocks.PegoutContractAdapterMock{}
		callerMock := &mocks.ContractCallerBindingMock{}
		pegoutContract := rootstock.NewPegoutContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractBinding, signerMock, rootstock.RetryParams{}, time.Duration(1), Abis)
		contractBinding.On("Caller").Return(callerMock).Once()
		callerMock.On("Call", mock.Anything, mock.Anything, "refundPegOut",
			refundParams.QuoteHash, refundParams.BtcRawTx, refundParams.BtcBlockHeaderHash,
			refundParams.MerkleBranchPath, refundParams.MerkleBranchHashes,
		).Return(assert.AnError).Once()
		result, err := pegoutContract.RefundPegout(txConfig, refundParams)
		require.Error(t, err)
		assert.Empty(t, result)
		contractBinding.AssertExpectations(t)
		contractBinding.AssertNotCalled(t, "RefundPegOut")
		callerMock.AssertExpectations(t)
	})
	t.Run("Error handling (Transaction send error)", func(t *testing.T) {
		contractBinding := &mocks.PegoutContractAdapterMock{}
		callerMock := &mocks.ContractCallerBindingMock{}
		pegoutContract := rootstock.NewPegoutContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractBinding, signerMock, rootstock.RetryParams{}, time.Duration(1), Abis)
		contractBinding.On("Caller").Return(callerMock).Once()
		callerMock.On("Call", mock.Anything, mock.Anything, "refundPegOut",
			refundParams.QuoteHash, refundParams.BtcRawTx, refundParams.BtcBlockHeaderHash,
			refundParams.MerkleBranchPath, refundParams.MerkleBranchHashes,
		).Return(nil).Once()
		_ = prepareTxMocks(mockClient, signerMock, true)
		contractBinding.On("RefundPegOut", mock.MatchedBy(matchOptsFunc),
			refundParams.QuoteHash, refundParams.BtcRawTx, refundParams.BtcBlockHeaderHash,
			refundParams.MerkleBranchPath, refundParams.MerkleBranchHashes,
		).Return(nil, assert.AnError).Once()
		result, err := pegoutContract.RefundPegout(txConfig, refundParams)
		require.ErrorContains(t, err, "refund pegout error")
		assert.Empty(t, result)
		contractBinding.AssertExpectations(t)
		callerMock.AssertExpectations(t)
	})
	t.Run("Error handling (Transaction reverted)", func(t *testing.T) {
		contractBinding := &mocks.PegoutContractAdapterMock{}
		callerMock := &mocks.ContractCallerBindingMock{}
		pegoutContract := rootstock.NewPegoutContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractBinding, signerMock, rootstock.RetryParams{}, time.Duration(1), Abis)
		contractBinding.On("Caller").Return(callerMock).Once()
		callerMock.On("Call", mock.Anything, mock.Anything, "refundPegOut",
			refundParams.QuoteHash, refundParams.BtcRawTx, refundParams.BtcBlockHeaderHash,
			refundParams.MerkleBranchPath, refundParams.MerkleBranchHashes,
		).Return(nil).Once()
		tx := prepareTxMocks(mockClient, signerMock, false)
		contractBinding.On("RefundPegOut", mock.MatchedBy(matchOptsFunc),
			refundParams.QuoteHash, refundParams.BtcRawTx, refundParams.BtcBlockHeaderHash,
			refundParams.MerkleBranchPath, refundParams.MerkleBranchHashes,
		).Return(tx, nil).Once()
		result, err := pegoutContract.RefundPegout(txConfig, refundParams)
		require.ErrorContains(t, err, "refund pegout error: transaction reverted")
		assert.Equal(t, tx.Hash().String(), result)
		contractBinding.AssertExpectations(t)
		callerMock.AssertExpectations(t)
	})
}

func TestPegoutContractImpl_GetDepositEvents(t *testing.T) {
	contractBinding := &mocks.PegoutContractAdapterMock{}
	iteratorMock := &mocks.EventIteratorAdapterMock[bindings.IPegOutPegOutDeposit]{}
	filterMatchFunc := func(from uint64, to uint64) func(opts *bind.FilterOpts) bool {
		return func(opts *bind.FilterOpts) bool {
			return from == opts.Start && to == *opts.End && opts.Context != nil
		}
	}
	pegoutContract := rootstock.NewPegoutContractImpl(dummyClient, test.AnyAddress, contractBinding, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
	t.Run("Success", func(t *testing.T) {
		var from uint64 = 500
		var to uint64 = 1000
		contractBinding.EXPECT().FilterPegOutDeposit(mock.MatchedBy(filterMatchFunc(from, to)), [][32]uint8(nil), []common.Address(nil), []*big.Int(nil)).
			Return(&bindings.IPegOutPegOutDepositIterator{}, nil).Once()
		contractBinding.On("DepositEventIteratorAdapter", mock.AnythingOfType(depositIteratorString)).
			Return(iteratorMock)
		iteratorMock.On("Next").Return(true).Times(len(deposits))
		iteratorMock.On("Next").Return(false).Once()
		for _, deposit := range deposits {
			iteratorMock.On("Event").Return(deposit).Once()
		}
		iteratorMock.On("Error").Return(nil).Once()
		iteratorMock.On("Close").Return(nil).Once()
		result, err := pegoutContract.GetDepositEvents(context.Background(), from, &to)
		require.NoError(t, err)
		assert.Equal(t, parsedDeposits, result)
		contractBinding.AssertExpectations(t)
		iteratorMock.AssertExpectations(t)
	})
	t.Run("Error handling when failed to get iterator", func(t *testing.T) {
		var from uint64 = 600
		var to uint64 = 1100
		contractBinding.EXPECT().FilterPegOutDeposit(mock.MatchedBy(filterMatchFunc(from, to)), [][32]uint8(nil), []common.Address(nil), []*big.Int(nil)).
			Return(nil, assert.AnError).Once()
		contractBinding.On("DepositEventIteratorAdapter", mock.AnythingOfType(depositIteratorString)).
			Return(nil)
		result, err := pegoutContract.GetDepositEvents(context.Background(), from, &to)
		require.Error(t, err)
		assert.Nil(t, result)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling on iterator error", func(t *testing.T) {
		var from uint64 = 700
		var to uint64 = 1200
		contractBinding.EXPECT().FilterPegOutDeposit(mock.MatchedBy(filterMatchFunc(from, to)), [][32]uint8(nil), []common.Address(nil), []*big.Int(nil)).
			Return(&bindings.IPegOutPegOutDepositIterator{}, nil).Once()
		contractBinding.On("DepositEventIteratorAdapter", mock.AnythingOfType(depositIteratorString)).
			Return(iteratorMock)
		iteratorMock.On("Next").Return(false).Once()
		iteratorMock.On("Error").Return(assert.AnError).Once()
		iteratorMock.On("Close").Return(nil).Once()
		result, err := pegoutContract.GetDepositEvents(context.Background(), from, &to)
		require.Error(t, err)
		assert.Nil(t, result)
		contractBinding.AssertExpectations(t)
		iteratorMock.AssertExpectations(t)
	})
}
