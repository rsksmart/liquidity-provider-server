package pegout_test

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
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
	PenaltyFee:            entities.NewWei(2),
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
	ProductFeeAmount:      entities.NewWei(500),
}

func TestSendPegoutUseCase_Run_Paused(t *testing.T) {
	mutex := new(mocks.MutexMock)
	eventBus := new(mocks.EventBusMock)
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rsk := new(mocks.RootstockRpcServerMock)
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: true, Since: 5, Reason: "test"}, nil)
	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{PegOut: pegoutContract}, mutex, rootstock.ParseDepositEvent)
	err := useCase.Run(context.Background(), sendPegoutRetainedQuote)
	require.ErrorIs(t, err, blockchain.ContractPausedError)
}

// nolint:funlen
func TestSendPegoutUseCase_Run_InternalTransaction(t *testing.T) {
	btcTxHash := "0x5b5c5d"
	btcWallet := new(mocks.BitcoinWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(10000), nil).Once()
	creationData := quote.PegoutCreationData{FeeRate: utils.NewBigFloat64(1.5), FeePercentage: utils.NewBigFloat64(0.5), GasPrice: entities.NewWei(100), FixedFee: entities.NewWei(100)}
	quoteHash, err := hex.DecodeString(sendPegoutRetainedQuote.QuoteHash)
	require.NoError(t, err)
	unfundedTx := []byte{0x01, 0x02, 0x03, 0x04}
	btcWallet.On("CreateUnfundedTransactionWithOpReturn", sendPegoutRetainedQuote.DepositAddress, sendPegoutTestQuote.Value, quoteHash).Return(unfundedTx, nil).Once()
	btcWallet.On("SendWithOpReturn", sendPegoutRetainedQuote.DepositAddress, sendPegoutTestQuote.Value, quoteHash).Return(blockchain.BitcoinTransactionResult{
		Hash: btcTxHash,
		Fee:  entities.NewWei(1000),
	}, nil).Once()
	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
		expected := sendPegoutRetainedQuote
		expected.LpBtcTxHash = btcTxHash
		expected.SendPegoutBtcFee = entities.NewWei(1000)
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
	receipt := &blockchain.TransactionReceipt{
		TransactionHash:   "0x5b5c5d",
		BlockHash:         blockHash,
		BlockNumber:       blockNumber,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8500),
		Logs: []blockchain.TransactionLog{
			{
				Address: "0x1122",
				Topics: [][32]byte{
					{1, 2, 3},
					{4, 5, 6},
				},
				Data:        []byte{11, 22, 33, 44},
				BlockNumber: 123,
				TxHash:      test.AnyHash,
				TxIndex:     0,
				BlockHash:   test.AnyHash,
				Index:       0,
				Removed:     false,
			},
		},
	}
	receipt = test.AddDepositLogFromQuote(t, receipt, sendPegoutTestQuote, sendPegoutRetainedQuote)
	receipt.Value = entities.NewWei(0)
	receipt.To = "0x1122"
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(*receipt, nil).Once()
	rsk.On("GetBlockByHash", test.AnyCtx, blockHash).Return(blockchain.BlockInfo{
		Hash: blockHash, Number: blockNumber,
		Timestamp: time.Unix(int64(now+10), 0), Nonce: 1,
	}, nil).Once()
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
	quoteRepository.EXPECT().GetPegoutCreationData(test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(creationData).Once()
	updatedQuote := sendPegoutRetainedQuote
	updatedQuote.LpBtcTxHash = btcTxHash
	updatedQuote.SendPegoutBtcFee = entities.NewWei(1000)
	updatedQuote.State = quote.PegoutStateSendPegoutSucceeded
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, updatedQuote).Return(nil).Once()
	rpc := blockchain.Rpc{Rsk: rsk}
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	pegoutContract.EXPECT().IsPegOutQuoteCompleted(sendPegoutRetainedQuote.QuoteHash).Return(false, nil).Once()
	pegoutContract.On("ValidatePegout", sendPegoutRetainedQuote.QuoteHash, unfundedTx).Return(nil).Once()
	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{PegOut: pegoutContract}, mutex, rootstock.ParseDepositEvent)
	err = useCase.Run(context.Background(), sendPegoutRetainedQuote)

	require.NoError(t, err)
	btcWallet.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	pegoutContract.AssertExpectations(t)
	mutex.AssertExpectations(t)
}

// nolint:funlen
func TestSendPegoutUseCase_Run(t *testing.T) {
	btcTxHash := "0x5b5c5d"
	btcWallet := new(mocks.BitcoinWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(10000), nil).Once()
	creationData := quote.PegoutCreationData{FeeRate: utils.NewBigFloat64(1.5), FeePercentage: utils.NewBigFloat64(0.5), GasPrice: entities.NewWei(100), FixedFee: entities.NewWei(100)}
	quoteHash, err := hex.DecodeString(sendPegoutRetainedQuote.QuoteHash)
	require.NoError(t, err)
	unfundedTx := []byte{0x01, 0x02, 0x03, 0x04}
	btcFee := entities.NewWei(5000)
	btcResult := blockchain.BitcoinTransactionResult{Hash: btcTxHash, Fee: btcFee}
	btcWallet.On("CreateUnfundedTransactionWithOpReturn", sendPegoutRetainedQuote.DepositAddress, sendPegoutTestQuote.Value, quoteHash).Return(unfundedTx, nil).Once()
	btcWallet.On("SendWithOpReturn", sendPegoutRetainedQuote.DepositAddress, sendPegoutTestQuote.Value, quoteHash).Return(btcResult, nil).Once()
	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
		expected := sendPegoutRetainedQuote
		expected.LpBtcTxHash = btcTxHash
		expected.SendPegoutBtcFee = btcFee
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
	receipt := &blockchain.TransactionReceipt{
		TransactionHash:   "0x5b5c5d",
		BlockHash:         blockHash,
		BlockNumber:       blockNumber,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8500),
	}
	receipt = test.AddDepositLogFromQuote(t, receipt, sendPegoutTestQuote, sendPegoutRetainedQuote)
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(*receipt, nil).Once()
	rsk.On("GetBlockByHash", test.AnyCtx, blockHash).Return(blockchain.BlockInfo{
		Hash: blockHash, Number: blockNumber,
		Timestamp: time.Unix(int64(now+10), 0), Nonce: 1,
	}, nil).Once()
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
	quoteRepository.EXPECT().GetPegoutCreationData(test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(creationData).Once()
	updatedQuote := sendPegoutRetainedQuote
	updatedQuote.LpBtcTxHash = btcTxHash
	updatedQuote.SendPegoutBtcFee = btcFee
	updatedQuote.State = quote.PegoutStateSendPegoutSucceeded
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, updatedQuote).Return(nil).Once()
	rpc := blockchain.Rpc{Rsk: rsk}
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	pegoutContract.On("IsPegOutQuoteCompleted", sendPegoutRetainedQuote.QuoteHash).Return(false, nil).Once()
	pegoutContract.On("ValidatePegout", sendPegoutRetainedQuote.QuoteHash, unfundedTx).Return(nil).Once()
	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{PegOut: pegoutContract}, mutex, rootstock.ParseDepositEvent)
	err = useCase.Run(context.Background(), sendPegoutRetainedQuote)

	require.NoError(t, err)
	btcWallet.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	pegoutContract.AssertExpectations(t)
	pegoutContract.AssertCalled(t, "ValidatePegout", sendPegoutRetainedQuote.QuoteHash, unfundedTx)
	mutex.AssertExpectations(t)
}

func TestSendPegoutUseCase_Run_ShouldNotPublishRecoverableError(t *testing.T) {
	mutex := new(mocks.MutexMock)
	mutex.On("Lock")
	mutex.On("Unlock")

	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish")

	recoverableSetups := getRecoverableSetups(t)

	for _, setup := range recoverableSetups {
		quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
		btcWallet := new(mocks.BitcoinWalletMock)
		rsk := new(mocks.RootstockRpcServerMock)
		caseQuote := sendPegoutRetainedQuote
		pegoutContract := new(mocks.PegoutContractMock)
		pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
		setup(&caseQuote, btcWallet, rsk, quoteRepository, pegoutContract)
		rpc := blockchain.Rpc{Rsk: rsk}

		useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{PegOut: pegoutContract}, mutex, rootstock.ParseDepositEvent)
		err := useCase.Run(context.Background(), caseQuote)
		btcWallet.AssertExpectations(t)
		rsk.AssertExpectations(t)
		quoteRepository.AssertExpectations(t)
		eventBus.AssertNotCalled(t, "Publish", mock.Anything)
		pegoutContract.AssertExpectations(t)
		require.Error(t, err)
	}
}

// nolint:funlen
func getRecoverableSetups(t *testing.T) []func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BitcoinWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, pegoutContract *mocks.PegoutContractMock) {
	return []func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BitcoinWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, pegoutContract *mocks.PegoutContractMock){
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BitcoinWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, pegoutContract *mocks.PegoutContractMock) {
			sendPegoutRetainedQuote.State = quote.PegoutStateWaitingForDeposit
		},
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BitcoinWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, pegoutContract *mocks.PegoutContractMock) {
			sendPegoutRetainedQuote.UserRskTxHash = ""
		},
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BitcoinWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, pegoutContract *mocks.PegoutContractMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(nil, assert.AnError).Once()
		},
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BitcoinWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, pegoutContract *mocks.PegoutContractMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
			rsk.On("GetHeight", test.AnyCtx).Return(uint64(0), assert.AnError).Once()
		},
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BitcoinWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, pegoutContract *mocks.PegoutContractMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
			rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
			rsk.On("GetTransactionReceipt", test.AnyCtx, mock.Anything).Return(blockchain.TransactionReceipt{}, assert.AnError).Once()
		},
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BitcoinWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, pegoutContract *mocks.PegoutContractMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
			rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
			receipt := &blockchain.TransactionReceipt{
				TransactionHash:   "0x5b5c5d",
				BlockHash:         "0x6e6f6a",
				BlockNumber:       blockNumber,
				From:              "0x1234",
				To:                "0x5678",
				CumulativeGasUsed: big.NewInt(500),
				GasUsed:           big.NewInt(500),
				Value:             entities.NewWei(8500),
			}
			receipt = test.AddDepositLogFromQuote(t, receipt, sendPegoutTestQuote, *sendPegoutRetainedQuote)
			rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(*receipt, nil).Once()
			rsk.On("GetBlockByHash", test.AnyCtx, mock.Anything).Return(blockchain.BlockInfo{}, assert.AnError).Once()
		},
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BitcoinWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, pegoutContract *mocks.PegoutContractMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
			rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
			receipt := &blockchain.TransactionReceipt{
				TransactionHash:   "0x5b5c5d",
				BlockHash:         "0x6e6f6a",
				BlockNumber:       blockNumber,
				From:              "0x1234",
				To:                "0x5678",
				CumulativeGasUsed: big.NewInt(500),
				GasUsed:           big.NewInt(500),
				Value:             entities.NewWei(8500),
			}
			receipt = test.AddDepositLogFromQuote(t, receipt, sendPegoutTestQuote, *sendPegoutRetainedQuote)
			rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(*receipt, nil).Once()
			rsk.On("GetBlockByHash", test.AnyCtx, mock.Anything).
				Return(blockchain.BlockInfo{Timestamp: time.Unix(int64(now), 0)}, nil).Once()
			pegoutContract.On("IsPegOutQuoteCompleted", sendPegoutRetainedQuote.QuoteHash).Return(false, assert.AnError).Once()
		},
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BitcoinWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, pegoutContract *mocks.PegoutContractMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
			rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
			receipt := &blockchain.TransactionReceipt{
				TransactionHash:   "0x5b5c5d",
				BlockHash:         "0x6e6f6a",
				BlockNumber:       blockNumber,
				From:              "0x1234",
				To:                "0x5678",
				CumulativeGasUsed: big.NewInt(500),
				GasUsed:           big.NewInt(500),
				Value:             entities.NewWei(8500),
			}
			receipt = test.AddDepositLogFromQuote(t, receipt, sendPegoutTestQuote, *sendPegoutRetainedQuote)
			rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(*receipt, nil).Once()
			rsk.On("GetBlockByHash", test.AnyCtx, mock.Anything).
				Return(blockchain.BlockInfo{Timestamp: time.Unix(int64(now), 0)}, nil).Once()
			pegoutContract.On("IsPegOutQuoteCompleted", sendPegoutRetainedQuote.QuoteHash).Return(false, nil).Once()
			btcWallet.On("GetBalance").Return(entities.NewWei(0), assert.AnError).Once()
		},
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, btcWallet *mocks.BitcoinWalletMock, rsk *mocks.RootstockRpcServerMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, pegoutContract *mocks.PegoutContractMock) {
			quoteHash, err := hex.DecodeString(sendPegoutRetainedQuote.QuoteHash)
			require.NoError(t, err, "test data should have valid quote hash")
			unfundedTx := []byte{0x01, 0x02, 0x03, 0x04}
			quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
			rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
			receipt := &blockchain.TransactionReceipt{
				TransactionHash:   "0x5b5c5d",
				BlockHash:         "0x6e6f6a",
				BlockNumber:       blockNumber,
				From:              "0x1234",
				To:                "0x5678",
				CumulativeGasUsed: big.NewInt(500),
				GasUsed:           big.NewInt(500),
				Value:             entities.NewWei(8500),
			}
			receipt = test.AddDepositLogFromQuote(t, receipt, sendPegoutTestQuote, *sendPegoutRetainedQuote)
			rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(*receipt, nil).Once()
			rsk.On("GetBlockByHash", test.AnyCtx, mock.Anything).
				Return(blockchain.BlockInfo{Timestamp: time.Unix(int64(now), 0)}, nil).Once()
			pegoutContract.On("IsPegOutQuoteCompleted", sendPegoutRetainedQuote.QuoteHash).Return(false, nil).Once()
			btcWallet.On("GetBalance").Return(entities.NewWei(10000), nil).Once()
			btcWallet.On("CreateUnfundedTransactionWithOpReturn", sendPegoutRetainedQuote.DepositAddress, sendPegoutTestQuote.Value, quoteHash).Return(unfundedTx, nil).Once()
			// Simulate a network/RPC error that can't be parsed (no "reverted with:" marker)
			pegoutContract.On("ValidatePegout", sendPegoutRetainedQuote.QuoteHash, unfundedTx).Return(errors.New("error validating pegout: connection timeout")).Once()
		},
	}
}

func TestSendPegoutUseCase_Run_InsufficientAmount(t *testing.T) {
	btcWallet := new(mocks.BitcoinWalletMock)
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
	receipt := &blockchain.TransactionReceipt{
		TransactionHash:   sendPegoutRetainedQuote.UserRskTxHash,
		BlockHash:         "0x6e6f6a",
		BlockNumber:       blockNumber,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8000),
	}
	receipt = test.AddDepositLogFromQuote(t, receipt, sendPegoutTestQuote, sendPegoutRetainedQuote)
	amountHex := fmt.Sprintf("%064x", new(big.Int).Sub(pegoutQuote.Total().AsBigInt(), big.NewInt(1000)))
	timestampHex := fmt.Sprintf("%064x", uint64(pegoutQuote.DepositDateLimit-500))
	parsedData, err := hex.DecodeString(amountHex + timestampHex)
	require.NoError(t, err)
	receipt.Logs[0].Data = parsedData
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(*receipt, nil).Once()
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
	updatedQuote := sendPegoutRetainedQuote
	updatedQuote.State = quote.PegoutStateSendPegoutFailed
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, updatedQuote).Return(nil).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{PegOut: pegoutContract}, mutex, rootstock.ParseDepositEvent)
	err = useCase.Run(context.Background(), sendPegoutRetainedQuote)

	require.ErrorIs(t, err, usecases.InsufficientAmountError)
	quoteRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	mutex.AssertNotCalled(t, "Lock")
	mutex.AssertNotCalled(t, "Unlock")
	btcWallet.AssertNotCalled(t, "GetBalance")
	btcWallet.AssertNotCalled(t, "SendWithOpReturn")
	pegoutContract.AssertNotCalled(t, "IsPegOutQuoteCompleted")
}

func TestSendPegoutUseCase_Run_NoConfirmations(t *testing.T) {
	btcWallet := new(mocks.BitcoinWalletMock)
	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
	receipt := &blockchain.TransactionReceipt{
		TransactionHash:   sendPegoutRetainedQuote.UserRskTxHash,
		BlockHash:         "0x6e6f6a",
		BlockNumber:       445,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8500),
	}
	receipt = test.AddDepositLogFromQuote(t, receipt, sendPegoutTestQuote, sendPegoutRetainedQuote)
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(*receipt, nil).Once()
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{PegOut: pegoutContract}, mutex, rootstock.ParseDepositEvent)
	err := useCase.Run(context.Background(), sendPegoutRetainedQuote)

	require.ErrorIs(t, err, usecases.NoEnoughConfirmationsError)
	quoteRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	eventBus.AssertNotCalled(t, "Publish")
	quoteRepository.AssertNotCalled(t, "UpdateRetainedQuote")
	pegoutContract.AssertNotCalled(t, "IsPegOutQuoteCompleted")
}

func TestSendPegoutUseCase_Run_ExpiredQuote(t *testing.T) {
	btcWallet := new(mocks.BitcoinWalletMock)
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
	receipt := &blockchain.TransactionReceipt{
		TransactionHash:   "0x5b5c5d",
		BlockHash:         blockHash,
		BlockNumber:       blockNumber,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8500),
	}
	receipt = test.AddDepositLogFromQuote(t, receipt, sendPegoutTestQuote, sendPegoutRetainedQuote)
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(*receipt, nil).Once()
	rsk.On("GetBlockByHash", test.AnyCtx, blockHash).Return(blockchain.BlockInfo{
		Hash:      blockHash,
		Number:    blockNumber,
		Timestamp: time.Unix(int64(now+800), 0),
		Nonce:     1,
	}, nil).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{PegOut: pegoutContract}, mutex, rootstock.ParseDepositEvent)
	err := useCase.Run(context.Background(), sendPegoutRetainedQuote)

	require.ErrorIs(t, err, usecases.ExpiredQuoteError)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	rsk.AssertNotCalled(t, "GetTransactionReceipt")
	rsk.AssertNotCalled(t, "GetHeight")
	pegoutContract.AssertNotCalled(t, "IsPegOutQuoteCompleted")
}

func TestSendPegoutUseCase_Run_NoLiquidity(t *testing.T) {
	btcWallet := new(mocks.BitcoinWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(100), nil).Once()
	rsk := new(mocks.RootstockRpcServerMock)
	rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
	receipt := &blockchain.TransactionReceipt{
		TransactionHash:   sendPegoutRetainedQuote.UserRskTxHash,
		BlockHash:         blockHash,
		BlockNumber:       blockNumber,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8500),
	}
	receipt = test.AddDepositLogFromQuote(t, receipt, sendPegoutTestQuote, sendPegoutRetainedQuote)
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(*receipt, nil).Once()
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
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.On("IsPegOutQuoteCompleted", sendPegoutRetainedQuote.QuoteHash).Return(false, nil).Once()
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{PegOut: pegoutContract}, mutex, rootstock.ParseDepositEvent)
	err := useCase.Run(context.Background(), sendPegoutRetainedQuote)

	require.ErrorIs(t, err, usecases.NoLiquidityError)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertNotCalled(t, "Publish")
	btcWallet.AssertExpectations(t)
	rsk.AssertExpectations(t)
	mutex.AssertExpectations(t)
	pegoutContract.AssertExpectations(t)
}

// nolint:funlen
func TestSendPegoutUseCase_Run_QuoteNotFound(t *testing.T) {
	btcWallet := new(mocks.BitcoinWalletMock)
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
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{PegOut: pegoutContract}, mutex, rootstock.ParseDepositEvent)
	err := useCase.Run(context.Background(), sendPegoutRetainedQuote)

	require.ErrorIs(t, err, usecases.QuoteNotFoundError)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	mutex.AssertNotCalled(t, "Lock")
	mutex.AssertNotCalled(t, "Unlock")
	rsk.AssertNotCalled(t, "GetTransactionReceipt")
	rsk.AssertNotCalled(t, "GetHeight")
	pegoutContract.AssertNotCalled(t, "IsPegOutQuoteCompleted")
}

// nolint:funlen
func TestSendPegoutUseCase_Run_BtcTxFail(t *testing.T) {
	btcWallet := new(mocks.BitcoinWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(10000), nil).Once()
	quoteHash, err := hex.DecodeString(sendPegoutRetainedQuote.QuoteHash)
	creationData := quote.PegoutCreationData{
		FeeRate: utils.NewBigFloat64(55.5), FeePercentage: utils.NewBigFloat64(0.5),
		GasPrice: entities.NewWei(100), FixedFee: entities.NewWei(100),
	}
	require.NoError(t, err)
	unfundedTx := []byte{0x01, 0x02, 0x03, 0x04}
	btcWallet.On("CreateUnfundedTransactionWithOpReturn", sendPegoutRetainedQuote.DepositAddress, sendPegoutTestQuote.Value, quoteHash).Return(unfundedTx, nil).Once()
	btcWallet.On("SendWithOpReturn", sendPegoutRetainedQuote.DepositAddress, sendPegoutTestQuote.Value, quoteHash).Return(blockchain.BitcoinTransactionResult{}, assert.AnError).Once()
	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
		expected := sendPegoutRetainedQuote
		expected.State = quote.PegoutStateSendPegoutFailed
		require.Error(t, event.Error)
		return assert.Equal(t, sendPegoutTestQuote, event.PegoutQuote) && assert.Equal(t, expected, event.RetainedQuote) &&
			assert.Equal(t, quote.PegoutBtcSentEventId, event.Event.Id()) && assert.Equal(t, creationData, event.CreationData)
	})).Return().Once()
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
	receipt := &blockchain.TransactionReceipt{
		TransactionHash:   "0x5b5c5d",
		BlockHash:         blockHash,
		BlockNumber:       blockNumber,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8500),
	}
	receipt = test.AddDepositLogFromQuote(t, receipt, sendPegoutTestQuote, sendPegoutRetainedQuote)
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(*receipt, nil).Once()
	rsk.On("GetBlockByHash", test.AnyCtx, blockHash).Return(blockchain.BlockInfo{
		Hash:      blockHash,
		Number:    blockNumber,
		Timestamp: time.Unix(int64(now+10), 0),
		Nonce:     1,
	}, nil).Once()
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
	updatedQuote := sendPegoutRetainedQuote
	updatedQuote.State = quote.PegoutStateSendPegoutFailed
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, updatedQuote).Return(nil).Once()
	quoteRepository.EXPECT().GetPegoutCreationData(test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(creationData).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	pegoutContract.On("IsPegOutQuoteCompleted", sendPegoutRetainedQuote.QuoteHash).Return(false, nil).Once()
	pegoutContract.On("ValidatePegout", sendPegoutRetainedQuote.QuoteHash, unfundedTx).Return(nil).Once()

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{PegOut: pegoutContract}, mutex, rootstock.ParseDepositEvent)
	err = useCase.Run(context.Background(), sendPegoutRetainedQuote)
	require.Error(t, err)
	btcWallet.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	mutex.AssertExpectations(t)
	pegoutContract.AssertExpectations(t)
}

// nolint: funlen
func TestSendPegoutUseCase_Run_UpdateError(t *testing.T) {
	btcTxHash := "0x5b5c5d"
	btcWallet := new(mocks.BitcoinWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(10000), nil)
	quoteHash, err := hex.DecodeString(sendPegoutRetainedQuote.QuoteHash)
	require.NoError(t, err)
	unfundedTx := []byte{0x01, 0x02, 0x03, 0x04}
	btcFee := entities.NewWei(3000)
	btcResult2 := blockchain.BitcoinTransactionResult{Hash: btcTxHash, Fee: btcFee}
	btcWallet.On("CreateUnfundedTransactionWithOpReturn", sendPegoutRetainedQuote.DepositAddress, sendPegoutTestQuote.Value, quoteHash).Return(unfundedTx, nil)
	btcWallet.On("SendWithOpReturn", sendPegoutRetainedQuote.DepositAddress, sendPegoutTestQuote.Value, quoteHash).Return(btcResult2, nil)
	rsk := new(mocks.RootstockRpcServerMock)
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return()
	mutex.On("Unlock").Return()
	rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil)
	receipt := &blockchain.TransactionReceipt{
		TransactionHash: sendPegoutRetainedQuote.UserRskTxHash, Value: entities.NewWei(8500), BlockHash: blockHash,
		BlockNumber: blockNumber, From: "0x1234", To: "0x5678", CumulativeGasUsed: big.NewInt(500), GasUsed: big.NewInt(500),
	}
	receipt = test.AddDepositLogFromQuote(t, receipt, sendPegoutTestQuote, sendPegoutRetainedQuote)
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(*receipt, nil)
	rsk.On("GetBlockByHash", test.AnyCtx, blockHash).Return(blockchain.BlockInfo{
		Hash: blockHash, Number: blockNumber, Timestamp: time.Unix(int64(now+10), 0), Nonce: 1,
	}, nil)
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.On("IsPegOutQuoteCompleted", mock.Anything).Return(false, nil)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	pegoutContract.On("ValidatePegout", mock.Anything, mock.Anything).Return(nil)

	setups := []func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, quoteRepository *mocks.PegoutQuoteRepositoryMock, eventBus *mocks.EventBusMock){
		func(sendPegoutRetainedQuote *quote.RetainedPegoutQuote, quoteRepository *mocks.PegoutQuoteRepositoryMock, eventBus *mocks.EventBusMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
			quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, mock.Anything).Return(assert.AnError).Once()
			quoteRepository.EXPECT().GetPegoutCreationData(test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(quote.PegoutCreationDataZeroValue())
			eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
				expected := *sendPegoutRetainedQuote
				expected.LpBtcTxHash = btcTxHash
				expected.SendPegoutBtcFee = btcFee
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
				expected.SendPegoutBtcFee = nil
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
		useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{PegOut: pegoutContract}, mutex, rootstock.ParseDepositEvent)
		err = useCase.Run(context.Background(), caseQuote)
		quoteRepository.AssertExpectations(t)
		require.Error(t, err)
		btcWallet.AssertExpectations(t)
		rsk.AssertExpectations(t)
		eventBus.AssertExpectations(t)
		mutex.AssertExpectations(t)
	}
}

// nolint:funlen
func TestSendPegoutUseCase_Run_QuoteAlreadyCompleted(t *testing.T) {
	const errorMsg = "quote e64215867af36cad04e8c2e3e8336618b358f68923529f2a1e5dbc6dd4af4df1 was already completed"
	btcWallet := new(mocks.BitcoinWalletMock)
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
	receipt := &blockchain.TransactionReceipt{
		TransactionHash:   "0x5b5c5d",
		BlockHash:         blockHash,
		BlockNumber:       blockNumber,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8500),
	}
	receipt = test.AddDepositLogFromQuote(t, receipt, sendPegoutTestQuote, sendPegoutRetainedQuote)
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(*receipt, nil).Once()
	rsk.On("GetBlockByHash", test.AnyCtx, blockHash).Return(blockchain.BlockInfo{
		Hash:      blockHash,
		Number:    blockNumber,
		Timestamp: time.Unix(int64(now+10), 0),
		Nonce:     1,
	}, nil).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	pegoutContract.On("IsPegOutQuoteCompleted", sendPegoutRetainedQuote.QuoteHash).Return(true, nil).Once()

	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{PegOut: pegoutContract}, mutex, rootstock.ParseDepositEvent)
	err := useCase.Run(context.Background(), sendPegoutRetainedQuote)

	require.ErrorContains(t, err, errorMsg)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	rsk.AssertNotCalled(t, "GetTransactionReceipt")
	rsk.AssertNotCalled(t, "GetHeight")
	pegoutContract.AssertExpectations(t)
}

// nolint:funlen
func TestSendPegoutUseCase_Run_ValidationFailure(t *testing.T) {
	btcWallet := new(mocks.BitcoinWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(10000), nil).Once()
	quoteHash, err := hex.DecodeString(sendPegoutRetainedQuote.QuoteHash)
	require.NoError(t, err)

	unfundedTx := []byte{0x01, 0x02, 0x03, 0x04}
	btcWallet.On("CreateUnfundedTransactionWithOpReturn", sendPegoutRetainedQuote.DepositAddress, sendPegoutTestQuote.Value, quoteHash).Return(unfundedTx, nil).Once()

	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
		expected := sendPegoutRetainedQuote
		expected.State = quote.PegoutStateSendPegoutFailed
		expected.UserRskTxHash = "0x5b5c5d"
		require.Error(t, event.Error)
		require.ErrorContains(t, event.Error, "reverted with:")
		return assert.Equal(t, sendPegoutTestQuote, event.PegoutQuote) &&
			assert.Equal(t, expected, event.RetainedQuote) &&
			assert.Equal(t, quote.PegoutBtcSentEventId, event.Event.Id())
	})).Return().Once()
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
	receipt := &blockchain.TransactionReceipt{
		TransactionHash:   "0x5b5c5d",
		BlockHash:         blockHash,
		BlockNumber:       blockNumber,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8500),
	}
	receipt = test.AddDepositLogFromQuote(t, receipt, sendPegoutTestQuote, sendPegoutRetainedQuote)
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(*receipt, nil).Once()
	rsk.On("GetBlockByHash", test.AnyCtx, blockHash).Return(blockchain.BlockInfo{
		Hash: blockHash, Number: blockNumber,
		Timestamp: time.Unix(int64(now+10), 0), Nonce: 1,
	}, nil).Once()
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
	// GetPegoutCreationData is NOT called on validation failure - publishErrorEvent uses zero value
	updatedQuote := sendPegoutRetainedQuote
	updatedQuote.State = quote.PegoutStateSendPegoutFailed
	updatedQuote.UserRskTxHash = "0x5b5c5d" // Receipt transaction hash
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, updatedQuote).Return(nil).Once()
	rpc := blockchain.Rpc{Rsk: rsk}
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	pegoutContract.EXPECT().IsPegOutQuoteCompleted(sendPegoutRetainedQuote.QuoteHash).Return(false, nil).Once()
	// Mock validation failure - simulate a parsed contract revert
	pegoutContract.On("ValidatePegout", sendPegoutRetainedQuote.QuoteHash, unfundedTx).Return(errors.New("validatePegout reverted with: InvalidDestination")).Once()

	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{PegOut: pegoutContract}, mutex, rootstock.ParseDepositEvent)
	err = useCase.Run(context.Background(), sendPegoutRetainedQuote)

	require.Error(t, err)
	require.ErrorContains(t, err, "reverted with:")
	btcWallet.AssertExpectations(t)
	btcWallet.AssertNotCalled(t, "SendWithOpReturn")
	quoteRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	pegoutContract.AssertExpectations(t)
	mutex.AssertExpectations(t)
}

func TestSendPegoutUseCase_Run_UnfundedTxCreationError(t *testing.T) {
	btcWallet := new(mocks.BitcoinWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(10000), nil).Once()
	quoteHash, err := hex.DecodeString(sendPegoutRetainedQuote.QuoteHash)
	require.NoError(t, err)

	// Mock unfunded transaction creation error
	btcWallet.On("CreateUnfundedTransactionWithOpReturn", sendPegoutRetainedQuote.DepositAddress, sendPegoutTestQuote.Value, quoteHash).Return([]byte(nil), errors.New("failed to create transaction")).Once()
	rsk := new(mocks.RootstockRpcServerMock)
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutBtcSentToUserEvent) bool {
		expected := sendPegoutRetainedQuote
		expected.State = quote.PegoutStateSendPegoutFailed
		expected.UserRskTxHash = "0x5b5c5d" // Receipt transaction hash
		require.Error(t, event.Error)
		require.ErrorContains(t, event.Error, "failed to create unfunded transaction")
		return assert.Equal(t, sendPegoutTestQuote, event.PegoutQuote) &&
			assert.Equal(t, expected, event.RetainedQuote) &&
			assert.Equal(t, quote.PegoutBtcSentEventId, event.Event.Id())
	})).Return().Once()
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	rsk.On("GetHeight", test.AnyCtx).Return(uint64(450), nil).Once()
	receipt := &blockchain.TransactionReceipt{
		TransactionHash:   "0x5b5c5d",
		BlockHash:         blockHash,
		BlockNumber:       blockNumber,
		From:              "0x1234",
		To:                "0x5678",
		CumulativeGasUsed: big.NewInt(500),
		GasUsed:           big.NewInt(500),
		Value:             entities.NewWei(8500),
	}
	receipt = test.AddDepositLogFromQuote(t, receipt, sendPegoutTestQuote, sendPegoutRetainedQuote)
	rsk.On("GetTransactionReceipt", test.AnyCtx, sendPegoutRetainedQuote.UserRskTxHash).Return(*receipt, nil).Once()
	rsk.On("GetBlockByHash", test.AnyCtx, blockHash).Return(blockchain.BlockInfo{
		Hash: blockHash, Number: blockNumber,
		Timestamp: time.Unix(int64(now+10), 0), Nonce: 1,
	}, nil).Once()
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, sendPegoutRetainedQuote.QuoteHash).Return(&sendPegoutTestQuote, nil).Once()
	// GetPegoutCreationData is NOT called on tx creation error - publishErrorEvent uses zero value
	updatedQuote := sendPegoutRetainedQuote
	updatedQuote.State = quote.PegoutStateSendPegoutFailed
	updatedQuote.UserRskTxHash = "0x5b5c5d" // Receipt transaction hash
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, updatedQuote).Return(nil).Once()
	rpc := blockchain.Rpc{Rsk: rsk}
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	pegoutContract.EXPECT().IsPegOutQuoteCompleted(sendPegoutRetainedQuote.QuoteHash).Return(false, nil).Once()

	useCase := pegout.NewSendPegoutUseCase(btcWallet, quoteRepository, rpc, eventBus, blockchain.RskContracts{PegOut: pegoutContract}, mutex, rootstock.ParseDepositEvent)
	err = useCase.Run(context.Background(), sendPegoutRetainedQuote)

	require.Error(t, err)
	require.ErrorContains(t, err, "failed to create unfunded transaction")
	btcWallet.AssertExpectations(t)
	btcWallet.AssertNotCalled(t, "SendWithOpReturn")
	quoteRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	pegoutContract.AssertExpectations(t)
	pegoutContract.AssertNotCalled(t, "ValidatePegout")
	mutex.AssertExpectations(t)
}
