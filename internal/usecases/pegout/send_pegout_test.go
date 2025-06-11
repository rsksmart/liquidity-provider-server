package pegout_test

import (
	"context"
	"encoding/hex"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

const (
	blockHash          = "0x6e6f6a"
	blockNumber uint64 = 440
)

var sendPegoutRetainedQuote = quote.RetainedPegoutQuote{
	QuoteHash:          "e64215867af36cad04e8c2e3e8336618b358f68923529f2a1e5dbc6dd4af4df1",
	DepositAddress:     "0x654321",
	Signature:          "0x112a3b",
	RequiredLiquidity:  entities.NewWei(1000),
	State:              quote.PegoutStateWaitingForDepositConfirmations,
	UserRskTxHash:      "0x3c2b1a",
	LpBtcTxHash:        "",
	RefundPegoutTxHash: "",
	BridgeRefundTxHash: "",
}

var sendPegoutTestQuote = quote.PegoutQuote{
	LbcAddress:            "0x5678",
	LpRskAddress:          "0x1234",
	BtcRefundAddress:      test.AnyAddress,
	RskRefundAddress:      "0x1234",
	LpBtcAddress:          "0x1234",
	CallFee:               entities.NewWei(3000),
	PenaltyFee:            2,
	Nonce:                 3,
	DepositAddress:        sendPegoutRetainedQuote.DepositAddress,
	Value:                 entities.NewWei(4000),
	AgreementTimestamp:    now,
	DepositDateLimit:      now + 60,
	DepositConfirmations:  10,
	TransferConfirmations: 10,
	TransferTime:          60,
	ExpireDate:            now + 60,
	ExpireBlock:           500,
	GasFee:                entities.NewWei(1000),
	ProductFeeAmount:      500,
}

// nolint:funlen
func TestSendPegoutUseCase_Run(t *testing.T) {
	receiptData := blockchain.BitcoinTransactionResult{
		Hash: "0x5b5c5d",
		Fee:  entities.NewWei(100),
	}
	btcWallet := new(mocks.BtcWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(10000), nil).Once()
	creationData := quote.PegoutCreationData{FeeRate: utils.NewBigFloat64(1.5), FeePercentage: utils.NewBigFloat64(0.5), GasPrice: entities.NewWei(100), FixedFee: entities.NewWei(100)}
	quoteHash, err := hex.DecodeString(sendPegoutRetainedQuote.QuoteHash)
	require.NoError(t, err)
	btcWallet.On("SendWithOpReturn", sendPegoutRetainedQuote.DepositAddress, sendPegoutTestQuote.Value, quoteHash).Return(receiptData, nil).Once()
	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
		expected := sendPegoutRetainedQuote
		expected.LpBtcTxHash = receiptData.Hash
		if receiptData.Fee != nil {
			expected.LpBtcTxFee = receiptData.Fee
		}
		expected.State = quote.PegoutStateSendPegoutSucceeded
		require.NoError(t, event.Error)
		return assert.Equal(t, sendPegoutTestQuote, event.PegoutQuote) &&
			assert.Equal(t, expected, event.RetainedQuote) &&
			assert.Equal(t, creationData, event.CreationData) &&
			assert.Equal(t, quote.PegoutBtcSentEventId, event.Event.Id())
	})).Return().Once()
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
		TransactionHash:   "0x5b5c5d",
		BlockHash:         blockHash,
		BlockNumber:       blockNumber,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8500),
	}, nil).Once()
	rsk.On("GetBlockByHash", test.AnyCtx, blockHash).Return(blockchain.BlockInfo{
		Hash:      blockHash,
		Number:    blockNumber,
		Timestamp: time.Unix(int64(now+10), 0),
		Nonce:     1,
	}, nil).Once()
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
	quoteRepository.EXPECT().GetPegoutCreationData(test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(creationData).Once()
	updatedQuote := sendPegoutRetainedQuote
	updatedQuote.LpBtcTxHash = receiptData.Hash
	if receiptData.Fee != nil {
		updatedQuote.LpBtcTxFee = receiptData.Fee
	}
	updatedQuote.State = quote.PegoutStateSendPegoutSucceeded
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, updatedQuote).Return(nil).Once()
	rpc := blockchain.Rpc{Rsk: rsk}
	lbc := new(mocks.LbcMock)
	lbc.On("IsPegOutQuoteCompleted", sendPegoutRetainedQuote.QuoteHash).Return(false, nil).Once()
	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{Lbc: lbc}, mutex)
	err = useCase.Run(context.Background(), sendPegoutRetainedQuote)

	require.NoError(t, err)
	btcWallet.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	lbc.AssertExpectations(t)
	mutex.AssertExpectations(t)
}

func TestSendPegoutUseCase_Run_ShouldNotPublishRecoverableError(t *testing.T) {
	mutex := new(mocks.MutexMock)
	mutex.On("Lock")
	mutex.On("Unlock")

	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish")

	recoverableSetups := getRecoverableSetups()

	for _, setup := range recoverableSetups {
		quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
		btcWallet := new(mocks.BtcWalletMock)
		rsk := new(mocks.RootstockRpcServerMock)
		caseQuote := sendPegoutRetainedQuote
		lbc := new(mocks.LbcMock)
		setup(&caseQuote, btcWallet, rsk, quoteRepository, lbc)
		rpc := blockchain.Rpc{Rsk: rsk}
		useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{Lbc: lbc}, mutex)
		err := useCase.Run(context.Background(), caseQuote)
		btcWallet.AssertExpectations(t)
		rsk.AssertExpectations(t)
		quoteRepository.AssertExpectations(t)
		eventBus.AssertNotCalled(t, "Publish", mock.Anything)
		lbc.AssertExpectations(t)
		require.Error(t, err)
	}
}

// nolint:funlen
func getRecoverableSetups() []func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock) {
	return []func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock){
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock) {
			sendPegoutRetainedQuote.State = quote.PegoutStateWaitingForDeposit
		},
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock) {
			sendPegoutRetainedQuote.UserRskTxHash = ""
		},
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(nil, assert.AnError).Once()
		},
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
			rsk.On("GetHeight", test.AnyCtx).Return(uint64(0), assert.AnError).Once()
		},
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
			rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
			rsk.On("GetTransactionReceipt", test.AnyCtx, mock.Anything).Return(blockchain.TransactionReceipt{}, assert.AnError).Once()
		},
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
			rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
			rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
				TransactionHash:   "0x5b5c5d",
				BlockHash:         "0x6e6f6a",
				BlockNumber:       blockNumber,
				From:              "0x1234",
				To:                "0x5678",
				CumulativeGasUsed: big.NewInt(500),
				GasUsed:           big.NewInt(500),
				Value:             entities.NewWei(8500),
			}, nil).Once()
			rsk.On("GetBlockByHash", test.AnyCtx, mock.Anything).Return(blockchain.BlockInfo{}, assert.AnError).Once()
		},
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
			rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
			rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
				TransactionHash:   "0x5b5c5d",
				BlockHash:         "0x6e6f6a",
				BlockNumber:       blockNumber,
				From:              "0x1234",
				To:                "0x5678",
				CumulativeGasUsed: big.NewInt(500),
				GasUsed:           big.NewInt(500),
				Value:             entities.NewWei(8500),
			}, nil).Once()
			rsk.On("GetBlockByHash", test.AnyCtx, mock.Anything).
				Return(blockchain.BlockInfo{Timestamp: time.Unix(int64(now), 0)}, nil).Once()
			lbc.On("IsPegOutQuoteCompleted", sendPegoutRetainedQuote.QuoteHash).Return(false, assert.AnError).Once()
		},
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
			rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
			rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
				TransactionHash:   "0x5b5c5d",
				BlockHash:         "0x6e6f6a",
				BlockNumber:       blockNumber,
				From:              "0x1234",
				To:                "0x5678",
				CumulativeGasUsed: big.NewInt(500),
				GasUsed:           big.NewInt(500),
				Value:             entities.NewWei(8500),
			}, nil).Once()
			rsk.On("GetBlockByHash", test.AnyCtx, mock.Anything).
				Return(blockchain.BlockInfo{Timestamp: time.Unix(int64(now), 0)}, nil).Once()
			lbc.On("IsPegOutQuoteCompleted", sendPegoutRetainedQuote.QuoteHash).Return(false, nil).Once()
			btcWallet.On("GetBalance").Return(entities.NewWei(0), assert.AnError).Once()
		},
	}
}

func TestSendPegoutUseCase_Run_InsufficientAmount(t *testing.T) {
	btcWallet := new(mocks.BtcWalletMock)
	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
		expected := sendPegoutRetainedQuote
		expected.State = quote.PegoutStateSendPegoutFailed
		require.ErrorIs(t, event.Error, usecases.InsufficientAmountError)
		return assert.Equal(t, sendPegoutTestQuote, event.PegoutQuote) &&
			assert.Equal(t, expected, event.RetainedQuote) &&
			assert.Equal(t, quote.PegoutBtcSentEventId, event.Event.Id())
	})).Return().Once()
	mutex := new(mocks.MutexMock)
	rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
		TransactionHash:   sendPegoutRetainedQuote.UserRskTxHash,
		BlockHash:         "0x6e6f6a",
		BlockNumber:       blockNumber,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8000),
	}, nil).Once()
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
	updatedQuote := sendPegoutRetainedQuote
	updatedQuote.State = quote.PegoutStateSendPegoutFailed
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, updatedQuote).Return(nil).Once()
	lbc := new(mocks.LbcMock)

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{Lbc: lbc}, mutex)
	err := useCase.Run(context.Background(), sendPegoutRetainedQuote)

	require.ErrorIs(t, err, usecases.InsufficientAmountError)
	quoteRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	mutex.AssertNotCalled(t, "Lock")
	mutex.AssertNotCalled(t, "Unlock")
	btcWallet.AssertNotCalled(t, "GetBalance")
	btcWallet.AssertNotCalled(t, "SendWithOpReturn")
	lbc.AssertNotCalled(t, "IsPegOutQuoteCompleted")
}

func TestSendPegoutUseCase_Run_NoConfirmations(t *testing.T) {
	btcWallet := new(mocks.BtcWalletMock)
	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
		TransactionHash:   sendPegoutRetainedQuote.UserRskTxHash,
		BlockHash:         "0x6e6f6a",
		BlockNumber:       445,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8500),
	}, nil).Once()
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
	lbc := new(mocks.LbcMock)

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{Lbc: lbc}, mutex)
	err := useCase.Run(context.Background(), sendPegoutRetainedQuote)

	require.ErrorIs(t, err, usecases.NoEnoughConfirmationsError)
	quoteRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	eventBus.AssertNotCalled(t, "Publish")
	quoteRepository.AssertNotCalled(t, "UpdateRetainedQuote")
	lbc.AssertNotCalled(t, "IsPegOutQuoteCompleted")
}

func TestSendPegoutUseCase_Run_ExpiredQuote(t *testing.T) {
	btcWallet := new(mocks.BtcWalletMock)
	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	expiredQuote := sendPegoutTestQuote
	expiredQuote.ExpireDate = now - 60
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&expiredQuote, nil).Once()
	updatedQuote := sendPegoutRetainedQuote
	updatedQuote.State = quote.PegoutStateSendPegoutFailed
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, updatedQuote).Return(nil).Once()
	eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
		expected := sendPegoutRetainedQuote
		expected.State = quote.PegoutStateSendPegoutFailed
		require.ErrorIs(t, event.Error, usecases.ExpiredQuoteError)
		return assert.Equal(t, expiredQuote, event.PegoutQuote) &&
			assert.Equal(t, expected, event.RetainedQuote) &&
			assert.Equal(t, quote.PegoutBtcSentEventId, event.Event.Id())
	})).Return().Once()
	rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
		TransactionHash:   "0x5b5c5d",
		BlockHash:         blockHash,
		BlockNumber:       blockNumber,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8500),
	}, nil).Once()
	rsk.On("GetBlockByHash", test.AnyCtx, blockHash).Return(blockchain.BlockInfo{
		Hash:      blockHash,
		Number:    blockNumber,
		Timestamp: time.Unix(int64(now+800), 0),
		Nonce:     1,
	}, nil).Once()
	lbc := new(mocks.LbcMock)

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{Lbc: lbc}, mutex)
	err := useCase.Run(context.Background(), sendPegoutRetainedQuote)

	require.ErrorIs(t, err, usecases.ExpiredQuoteError)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	rsk.AssertNotCalled(t, "GetTransactionReceipt")
	rsk.AssertNotCalled(t, "GetHeight")
	lbc.AssertNotCalled(t, "IsPegOutQuoteCompleted")
}

func TestSendPegoutUseCase_Run_NoLiquidity(t *testing.T) {
	btcWallet := new(mocks.BtcWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(100), nil).Once()
	rsk := new(mocks.RootstockRpcServerMock)
	rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
		TransactionHash:   sendPegoutRetainedQuote.UserRskTxHash,
		BlockHash:         blockHash,
		BlockNumber:       blockNumber,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8500),
	}, nil).Once()
	rsk.On("GetBlockByHash", test.AnyCtx, blockHash).Return(blockchain.BlockInfo{
		Hash:      blockHash,
		Number:    blockNumber,
		Timestamp: time.Unix(int64(now+10), 0),
		Nonce:     1,
	}, nil).Once()
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
	lbc := new(mocks.LbcMock)
	lbc.On("IsPegOutQuoteCompleted", sendPegoutRetainedQuote.QuoteHash).Return(false, nil).Once()

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{Lbc: lbc}, mutex)
	err := useCase.Run(context.Background(), sendPegoutRetainedQuote)

	require.ErrorIs(t, err, usecases.NoLiquidityError)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertNotCalled(t, "Publish")
	btcWallet.AssertExpectations(t)
	rsk.AssertExpectations(t)
	mutex.AssertExpectations(t)
	lbc.AssertExpectations(t)
}

func TestSendPegoutUseCase_Run_QuoteNotFound(t *testing.T) {
	btcWallet := new(mocks.BtcWalletMock)
	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(nil, nil).Once()
	updatedQuote := sendPegoutRetainedQuote
	updatedQuote.State = quote.PegoutStateSendPegoutFailed
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, updatedQuote).Return(nil).Once()
	eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
		expected := sendPegoutRetainedQuote
		expected.State = quote.PegoutStateSendPegoutFailed
		require.ErrorIs(t, event.Error, usecases.QuoteNotFoundError)
		return assert.Equal(t, quote.PegoutQuote{}, event.PegoutQuote) && assert.Equal(t, expected, event.RetainedQuote) && assert.Equal(t, quote.PegoutBtcSentEventId, event.Event.Id())
	})).Return().Once()
	lbc := new(mocks.LbcMock)

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{Lbc: lbc}, mutex)
	err := useCase.Run(context.Background(), sendPegoutRetainedQuote)

	require.ErrorIs(t, err, usecases.QuoteNotFoundError)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	mutex.AssertNotCalled(t, "Lock")
	mutex.AssertNotCalled(t, "Unlock")
	rsk.AssertNotCalled(t, "GetTransactionReceipt")
	rsk.AssertNotCalled(t, "GetHeight")
	lbc.AssertNotCalled(t, "IsPegOutQuoteCompleted")
}

// nolint:funlen
func TestSendPegoutUseCase_Run_BtcTxFail(t *testing.T) {
	dataReceipt := blockchain.BitcoinTransactionResult{
		Hash: "",
		Fee:  entities.NewWei(0),
	}
	btcWallet := new(mocks.BtcWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(10000), nil).Once()
	quoteHash, err := hex.DecodeString(sendPegoutRetainedQuote.QuoteHash)
	creationData := quote.PegoutCreationData{
		FeeRate:       utils.NewBigFloat64(55.5),
		FeePercentage: utils.NewBigFloat64(0.5),
		GasPrice:      entities.NewWei(100),
		FixedFee:      entities.NewWei(100),
	}
	require.NoError(t, err)
	btcWallet.On("SendWithOpReturn", sendPegoutRetainedQuote.DepositAddress, sendPegoutTestQuote.Value, quoteHash).Return(dataReceipt, assert.AnError).Once()
	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
		expected := sendPegoutRetainedQuote
		if dataReceipt.Fee != nil {
			expected.LpBtcTxFee = dataReceipt.Fee
		}
		expected.State = quote.PegoutStateSendPegoutFailed
		require.Error(t, event.Error)
		return assert.Equal(t, sendPegoutTestQuote, event.PegoutQuote) && assert.Equal(t, expected, event.RetainedQuote) &&
			assert.Equal(t, quote.PegoutBtcSentEventId, event.Event.Id()) && assert.Equal(t, creationData, event.CreationData)
	})).Return().Once()
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
		TransactionHash:   "0x5b5c5d",
		BlockHash:         blockHash,
		BlockNumber:       blockNumber,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8500),
	}, nil).Once()
	rsk.On("GetBlockByHash", test.AnyCtx, blockHash).Return(blockchain.BlockInfo{
		Hash:      blockHash,
		Number:    blockNumber,
		Timestamp: time.Unix(int64(now+10), 0),
		Nonce:     1,
	}, nil).Once()
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
	updatedQuote := sendPegoutRetainedQuote
	if dataReceipt.Fee != nil {
		updatedQuote.LpBtcTxFee = dataReceipt.Fee
	}
	updatedQuote.State = quote.PegoutStateSendPegoutFailed
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, updatedQuote).Return(nil).Once()
	quoteRepository.EXPECT().GetPegoutCreationData(test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(creationData).Once()
	lbc := new(mocks.LbcMock)
	lbc.On("IsPegOutQuoteCompleted", sendPegoutRetainedQuote.QuoteHash).Return(false, nil).Once()

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{Lbc: lbc}, mutex)
	if dataReceipt.Fee != nil {
		sendPegoutRetainedQuote.LpBtcTxFee = dataReceipt.Fee
	}
	err = useCase.Run(context.Background(), sendPegoutRetainedQuote)

	require.Error(t, err)
	btcWallet.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	mutex.AssertExpectations(t)
	lbc.AssertExpectations(t)
}

// nolint: funlen
func TestSendPegoutUseCase_Run_UpdateError(t *testing.T) {
	dataReceipt := blockchain.BitcoinTransactionResult{
		Hash: "0x5b5c5d",
		Fee:  entities.NewWei(100),
	}
	btcWallet := new(mocks.BtcWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(10000), nil)
	quoteHash, err := hex.DecodeString(sendPegoutRetainedQuote.QuoteHash)
	require.NoError(t, err)
	btcWallet.On("SendWithOpReturn", sendPegoutRetainedQuote.DepositAddress, sendPegoutTestQuote.Value, quoteHash).Return(dataReceipt, nil)
	rsk := new(mocks.RootstockRpcServerMock)
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return()
	mutex.On("Unlock").Return()
	rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil)
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
		TransactionHash: sendPegoutRetainedQuote.UserRskTxHash, Value: entities.NewWei(8500), BlockHash: blockHash,
		BlockNumber: blockNumber, From: "0x1234", To: "0x5678", CumulativeGasUsed: big.NewInt(500), GasUsed: big.NewInt(500),
	}, nil)
	rsk.On("GetBlockByHash", test.AnyCtx, blockHash).Return(blockchain.BlockInfo{
		Hash: blockHash, Number: blockNumber, Timestamp: time.Unix(int64(now+10), 0), Nonce: 1,
	}, nil)
	lbc := new(mocks.LbcMock)
	lbc.On("IsPegOutQuoteCompleted", mock.Anything).Return(false, nil)

	setups := []func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, quoteRepository *mocks.PegoutQuoteRepositoryMock, eventBus *mocks.EventBusMock){
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, quoteRepository *mocks.PegoutQuoteRepositoryMock, eventBus *mocks.EventBusMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
			quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, mock.Anything).Return(assert.AnError).Once()
			quoteRepository.EXPECT().GetPegoutCreationData(test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(quote.PegoutCreationDataZeroValue())
			eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
				expected := *sendPegoutRetainedQuote
				expected.LpBtcTxHash = dataReceipt.Hash
				if dataReceipt.Fee != nil {
					expected.LpBtcTxFee = dataReceipt.Fee
				}
				expected.State = quote.PegoutStateSendPegoutSucceeded
				require.NoError(t, event.Error)
				return assert.Equal(t, sendPegoutTestQuote, event.PegoutQuote) && assert.Equal(t, expected, event.RetainedQuote) && assert.Equal(t, quote.PegoutBtcSentEventId, event.Event.Id())
			})).Return().Once()
		},
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, quoteRepository *mocks.PegoutQuoteRepositoryMock, eventBus *mocks.EventBusMock) {
			sendPegoutRetainedQuote.QuoteHash = "no hex"
			quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
			quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, mock.Anything).Return(assert.AnError).Once()
			eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
				expected := *sendPegoutRetainedQuote
				expected.State = quote.PegoutStateSendPegoutFailed
				require.Error(t, event.Error)
				return assert.Equal(t, sendPegoutTestQuote, event.PegoutQuote) && assert.Equal(t, expected, event.RetainedQuote) && assert.Equal(t, quote.PegoutBtcSentEventId, event.Event.Id())
			})).Return().Once()
		},
	}
	for _, setup := range setups {
		caseQuote := sendPegoutRetainedQuote
		quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
		eventBus := new(mocks.EventBusMock)
		setup(&caseQuote, quoteRepository, eventBus)
		rpc := blockchain.Rpc{Rsk: rsk}
		useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{Lbc: lbc}, mutex)
		err = useCase.Run(context.Background(), caseQuote)
		quoteRepository.AssertExpectations(t)
		require.Error(t, err)
		btcWallet.AssertExpectations(t)
		rsk.AssertExpectations(t)
		eventBus.AssertExpectations(t)
		mutex.AssertExpectations(t)
	}
}

func TestSendPegoutUseCase_Run_QuoteAlreadyCompleted(t *testing.T) {
	const errorMsg = "quote e64215867af36cad04e8c2e3e8336618b358f68923529f2a1e5dbc6dd4af4df1 was already completed"
	btcWallet := new(mocks.BtcWalletMock)
	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
	updatedQuote := sendPegoutRetainedQuote
	updatedQuote.State = quote.PegoutStateSendPegoutFailed
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, updatedQuote).Return(nil).Once()
	eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
		expected := sendPegoutRetainedQuote
		expected.State = quote.PegoutStateSendPegoutFailed
		require.ErrorContains(t, event.Error, errorMsg)
		return assert.Equal(t, sendPegoutTestQuote, event.PegoutQuote) &&
			assert.Equal(t, expected, event.RetainedQuote) &&
			assert.Equal(t, quote.PegoutBtcSentEventId, event.Event.Id())
	})).Return().Once()
	rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
		TransactionHash:   "0x5b5c5d",
		BlockHash:         blockHash,
		BlockNumber:       blockNumber,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8500),
	}, nil).Once()
	rsk.On("GetBlockByHash", test.AnyCtx, blockHash).Return(blockchain.BlockInfo{
		Hash:      blockHash,
		Number:    blockNumber,
		Timestamp: time.Unix(int64(now+10), 0),
		Nonce:     1,
	}, nil).Once()
	lbc := new(mocks.LbcMock)
	lbc.On("IsPegOutQuoteCompleted", sendPegoutRetainedQuote.QuoteHash).Return(true, nil).Once()

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{Lbc: lbc}, mutex)
	err := useCase.Run(context.Background(), sendPegoutRetainedQuote)

	require.ErrorContains(t, err, errorMsg)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	rsk.AssertNotCalled(t, "GetTransactionReceipt")
	rsk.AssertNotCalled(t, "GetHeight")
	lbc.AssertExpectations(t)
}
