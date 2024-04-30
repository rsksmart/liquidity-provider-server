package pegout

import (
	"context"
	"encoding/hex"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

var now = uint32(time.Now().Unix())

const (
	blockHash          = "0x6e6f6a"
	blockNumber uint64 = 440
)

var retainedQuote = quote.RetainedPegoutQuote{
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

var pegoutQuote = quote.PegoutQuote{
	LbcAddress:            "0x5678",
	LpRskAddress:          "0x1234",
	BtcRefundAddress:      test.AnyAddress,
	RskRefundAddress:      "0x1234",
	LpBtcAddress:          "0x1234",
	CallFee:               entities.NewWei(3000),
	PenaltyFee:            2,
	Nonce:                 3,
	DepositAddress:        retainedQuote.DepositAddress,
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

func TestSendPegoutUseCase_Run(t *testing.T) {
	btcTxHash := "0x5b5c5d"
	btcWallet := new(mocks.BtcWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(10000), nil).Once()
	quoteHash, _ := hex.DecodeString(retainedQuote.QuoteHash)
	btcWallet.On("SendWithOpReturn", retainedQuote.DepositAddress, pegoutQuote.Value, quoteHash).Return(btcTxHash, nil).Once()
	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
		expected := retainedQuote
		expected.LpBtcTxHash = btcTxHash
		expected.State = quote.PegoutStateSendPegoutSucceeded
		require.NoError(t, event.Error)
		return assert.Equal(t, pegoutQuote, event.PegoutQuote) &&
			assert.Equal(t, expected, event.RetainedQuote) &&
			assert.Equal(t, quote.PegoutBtcSentEventId, event.Event.Id())
	})).Return().Once()
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
	rsk.On("GetTransactionReceipt", test.AnyCtx, retainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
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
	quoteRepository.On("GetQuote", test.AnyCtx, retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
	updatedQuote := retainedQuote
	updatedQuote.LpBtcTxHash = btcTxHash
	updatedQuote.State = quote.PegoutStateSendPegoutSucceeded
	quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), updatedQuote).Return(nil).Once()

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, mutex)
	err := useCase.Run(context.Background(), retainedQuote)

	require.NoError(t, err)
	btcWallet.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
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
		caseQuote := retainedQuote
		setup(&caseQuote, btcWallet, rsk, quoteRepository)
		rpc := blockchain.Rpc{Rsk: rsk}
		useCase := NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, mutex)
		err := useCase.Run(context.Background(), caseQuote)
		btcWallet.AssertExpectations(t)
		rsk.AssertExpectations(t)
		quoteRepository.AssertExpectations(t)
		eventBus.AssertNotCalled(t, "Publish", mock.Anything)
		require.Error(t, err)
	}
}

func getRecoverableSetups() []func(retainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock) {
	return []func(retainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock){
		func(retainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock) {
			retainedQuote.State = quote.PegoutStateWaitingForDeposit
		},
		func(retainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock) {
			retainedQuote.UserRskTxHash = ""
		},
		func(retainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(nil, assert.AnError).Once()
		},
		func(retainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(0), assert.AnError).Once()
		},
		func(retainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(450), nil).Once()
			rsk.On("GetTransactionReceipt", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(blockchain.TransactionReceipt{}, assert.AnError).Once()
		},
		func(retainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(450), nil).Once()
			rsk.On("GetTransactionReceipt", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
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
		func(retainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BtcWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(450), nil).Once()
			rsk.On("GetTransactionReceipt", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
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
			btcWallet.On("GetBalance").Return(entities.NewWei(0), assert.AnError).Once()
		},
	}
}

func TestSendPegoutUseCase_Run_InsufficientAmount(t *testing.T) {
	btcWallet := new(mocks.BtcWalletMock)
	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
		expected := retainedQuote
		expected.State = quote.PegoutStateSendPegoutFailed
		require.ErrorIs(t, event.Error, usecases.InsufficientAmountError)
		return assert.Equal(t, pegoutQuote, event.PegoutQuote) &&
			assert.Equal(t, expected, event.RetainedQuote) &&
			assert.Equal(t, quote.PegoutBtcSentEventId, event.Event.Id())
	})).Return().Once()
	mutex := new(mocks.MutexMock)
	rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(450), nil).Once()
	rsk.On("GetTransactionReceipt", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
		TransactionHash:   retainedQuote.UserRskTxHash,
		BlockHash:         "0x6e6f6a",
		BlockNumber:       blockNumber,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8000),
	}, nil).Once()
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
	updatedQuote := retainedQuote
	updatedQuote.State = quote.PegoutStateSendPegoutFailed
	quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), updatedQuote).Return(nil).Once()

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, mutex)
	err := useCase.Run(context.Background(), retainedQuote)

	require.ErrorIs(t, err, usecases.InsufficientAmountError)
	quoteRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	mutex.AssertNotCalled(t, "Lock")
	mutex.AssertNotCalled(t, "Unlock")
	btcWallet.AssertNotCalled(t, "GetBalance")
	btcWallet.AssertNotCalled(t, "SendWithOpReturn")
}

func TestSendPegoutUseCase_Run_NoConfirmations(t *testing.T) {
	btcWallet := new(mocks.BtcWalletMock)
	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(450), nil).Once()
	rsk.On("GetTransactionReceipt", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
		TransactionHash:   retainedQuote.UserRskTxHash,
		BlockHash:         "0x6e6f6a",
		BlockNumber:       445,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8500),
	}, nil).Once()
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, mutex)
	err := useCase.Run(context.Background(), retainedQuote)

	require.ErrorIs(t, err, usecases.NoEnoughConfirmationsError)
	quoteRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	eventBus.AssertNotCalled(t, "Publish")
	quoteRepository.AssertNotCalled(t, "UpdateRetainedQuote")
}

func TestSendPegoutUseCase_Run_ExpiredQuote(t *testing.T) {
	btcWallet := new(mocks.BtcWalletMock)
	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	expiredQuote := pegoutQuote
	expiredQuote.ExpireDate = now - 60
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&expiredQuote, nil).Once()
	updatedQuote := retainedQuote
	updatedQuote.State = quote.PegoutStateSendPegoutFailed
	quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), updatedQuote).Return(nil).Once()
	eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
		expected := retainedQuote
		expected.State = quote.PegoutStateSendPegoutFailed
		require.ErrorIs(t, event.Error, usecases.ExpiredQuoteError)
		return assert.Equal(t, expiredQuote, event.PegoutQuote) &&
			assert.Equal(t, expected, event.RetainedQuote) &&
			assert.Equal(t, quote.PegoutBtcSentEventId, event.Event.Id())
	})).Return().Once()
	rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
	rsk.On("GetTransactionReceipt", test.AnyCtx, retainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
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

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, mutex)
	err := useCase.Run(context.Background(), retainedQuote)

	require.ErrorIs(t, err, usecases.ExpiredQuoteError)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	rsk.AssertNotCalled(t, "GetTransactionReceipt")
	rsk.AssertNotCalled(t, "GetHeight")
}

func TestSendPegoutUseCase_Run_NoLiquidity(t *testing.T) {
	btcWallet := new(mocks.BtcWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(100), nil).Once()
	rsk := new(mocks.RootstockRpcServerMock)
	rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(450), nil).Once()
	rsk.On("GetTransactionReceipt", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
		TransactionHash:   retainedQuote.UserRskTxHash,
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
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, mutex)
	err := useCase.Run(context.Background(), retainedQuote)

	require.ErrorIs(t, err, usecases.NoLiquidityError)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertNotCalled(t, "Publish")
	btcWallet.AssertExpectations(t)
	rsk.AssertExpectations(t)
	mutex.AssertExpectations(t)
}

func TestSendPegoutUseCase_Run_QuoteNotFound(t *testing.T) {
	btcWallet := new(mocks.BtcWalletMock)
	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(nil, nil).Once()
	updatedQuote := retainedQuote
	updatedQuote.State = quote.PegoutStateSendPegoutFailed
	quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), updatedQuote).Return(nil).Once()
	eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
		expected := retainedQuote
		expected.State = quote.PegoutStateSendPegoutFailed
		require.ErrorIs(t, event.Error, usecases.QuoteNotFoundError)
		return assert.Equal(t, quote.PegoutQuote{}, event.PegoutQuote) && assert.Equal(t, expected, event.RetainedQuote) && assert.Equal(t, quote.PegoutBtcSentEventId, event.Event.Id())
	})).Return().Once()

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, mutex)
	err := useCase.Run(context.Background(), retainedQuote)

	require.ErrorIs(t, err, usecases.QuoteNotFoundError)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	mutex.AssertNotCalled(t, "Lock")
	mutex.AssertNotCalled(t, "Unlock")
	rsk.AssertNotCalled(t, "GetTransactionReceipt")
	rsk.AssertNotCalled(t, "GetHeight")
}

func TestSendPegoutUseCase_Run_BtcTxFail(t *testing.T) {
	btcWallet := new(mocks.BtcWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(10000), nil).Once()
	quoteHash, _ := hex.DecodeString(retainedQuote.QuoteHash)
	btcWallet.On("SendWithOpReturn", retainedQuote.DepositAddress, pegoutQuote.Value, quoteHash).Return("", assert.AnError).Once()
	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
		expected := retainedQuote
		expected.State = quote.PegoutStateSendPegoutFailed
		require.Error(t, event.Error)
		return assert.Equal(t, pegoutQuote, event.PegoutQuote) && assert.Equal(t, expected, event.RetainedQuote) && assert.Equal(t, quote.PegoutBtcSentEventId, event.Event.Id())
	})).Return().Once()
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(450), nil).Once()
	rsk.On("GetTransactionReceipt", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
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
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
	updatedQuote := retainedQuote
	updatedQuote.State = quote.PegoutStateSendPegoutFailed
	quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), updatedQuote).Return(nil).Once()

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, mutex)
	err := useCase.Run(context.Background(), retainedQuote)

	require.Error(t, err)
	btcWallet.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	mutex.AssertExpectations(t)
}

func TestSendPegoutUseCase_Run_UpdateError(t *testing.T) {
	btcTxHash := "0x5b5c5d"
	btcWallet := new(mocks.BtcWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(10000), nil)
	quoteHash, _ := hex.DecodeString(retainedQuote.QuoteHash)
	btcWallet.On("SendWithOpReturn", retainedQuote.DepositAddress, pegoutQuote.Value, quoteHash).Return(btcTxHash, nil)
	rsk := new(mocks.RootstockRpcServerMock)
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return()
	mutex.On("Unlock").Return()
	rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(450), nil)
	rsk.On("GetTransactionReceipt", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.UserRskTxHash).Return(blockchain.TransactionReceipt{
		TransactionHash: retainedQuote.UserRskTxHash, Value: entities.NewWei(8500),
		BlockHash: blockHash, BlockNumber: blockNumber,
		From: "0x1234", To: "0x5678",
		CumulativeGasUsed: big.NewInt(500), GasUsed: big.NewInt(500),
	}, nil)
	rsk.On("GetBlockByHash", test.AnyCtx, blockHash).Return(blockchain.BlockInfo{
		Hash: blockHash, Number: blockNumber,
		Timestamp: time.Unix(int64(now+10), 0), Nonce: 1,
	}, nil)

	setups := []func(retainedQuote *quote.RetainedPegoutQuote, quoteRepository *mocks.PegoutQuoteRepositoryMock, eventBus *mocks.EventBusMock){
		func(retainedQuote *quote.RetainedPegoutQuote, quoteRepository *mocks.PegoutQuoteRepositoryMock, eventBus *mocks.EventBusMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(assert.AnError).Once()
			eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
				expected := *retainedQuote
				expected.LpBtcTxHash = btcTxHash
				expected.State = quote.PegoutStateSendPegoutSucceeded
				require.NoError(t, event.Error)
				return assert.Equal(t, pegoutQuote, event.PegoutQuote) && assert.Equal(t, expected, event.RetainedQuote) && assert.Equal(t, quote.PegoutBtcSentEventId, event.Event.Id())
			})).Return().Once()
		},
		func(retainedQuote *quote.RetainedPegoutQuote, quoteRepository *mocks.PegoutQuoteRepositoryMock, eventBus *mocks.EventBusMock) {
			retainedQuote.QuoteHash = "no hex"
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(assert.AnError).Once()
			eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
				expected := *retainedQuote
				expected.State = quote.PegoutStateSendPegoutFailed
				require.Error(t, event.Error)
				return assert.Equal(t, pegoutQuote, event.PegoutQuote) && assert.Equal(t, expected, event.RetainedQuote) && assert.Equal(t, quote.PegoutBtcSentEventId, event.Event.Id())
			})).Return().Once()
		},
	}
	for _, setup := range setups {
		caseQuote := retainedQuote
		quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
		eventBus := new(mocks.EventBusMock)
		setup(&caseQuote, quoteRepository, eventBus)
		rpc := blockchain.Rpc{Rsk: rsk}
		useCase := NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, mutex)
		err := useCase.Run(context.Background(), caseQuote)
		quoteRepository.AssertExpectations(t)
		require.Error(t, err)
		btcWallet.AssertExpectations(t)
		rsk.AssertExpectations(t)
		eventBus.AssertExpectations(t)
		mutex.AssertExpectations(t)
	}
}
