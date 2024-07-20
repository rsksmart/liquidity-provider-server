package pegout_test

import (
	"context"
	"errors"
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

var (
	now                = uint32(time.Now().Unix())
	refundPegoutTxHash = "0x1f1d1b"
)

var retainedQuote = quote.RetainedPegoutQuote{
	QuoteHash:          "1c2d3f",
	DepositAddress:     "0x654321",
	Signature:          "0x112a3b",
	RequiredLiquidity:  entities.NewWei(1000),
	State:              quote.PegoutStateSendPegoutSucceeded,
	UserRskTxHash:      "0x3c2b1a",
	LpBtcTxHash:        "0x3c2b1b",
	RefundPegoutTxHash: "",
	BridgeRefundTxHash: "",
}

var pegoutQuote = quote.PegoutQuote{
	LbcAddress:            retainedQuote.QuoteHash,
	LpRskAddress:          "0x1234",
	BtcRefundAddress:      test.AnyAddress,
	RskRefundAddress:      "0x1234",
	LpBtcAddress:          "0x1234",
	CallFee:               entities.NewWei(3000),
	PenaltyFee:            2,
	Nonce:                 3,
	DepositAddress:        test.AnyAddress,
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

var btcBlockInfoMock = blockchain.BitcoinBlockInformation{
	Hash:   [32]byte{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
	Height: big.NewInt(200),
}

var merkleBranchMock = blockchain.MerkleBranch{
	Hashes: [][32]byte{
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
		{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
	},
	Path: big.NewInt(1),
}

var btcTxInfoMock = blockchain.BitcoinTransactionInformation{
	Hash:          "0x3c2b1b",
	Confirmations: 11,
	Outputs:       map[string][]*entities.Wei{test.AnyAddress: {entities.NewWei(1000)}},
}

var btcRawTxMock = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

func TestRefundPegoutUseCase_Run(t *testing.T) {
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	expectedRetained := quote.RetainedPegoutQuote{
		QuoteHash:          retainedQuote.QuoteHash,
		DepositAddress:     retainedQuote.DepositAddress,
		Signature:          retainedQuote.Signature,
		RequiredLiquidity:  retainedQuote.RequiredLiquidity,
		State:              quote.PegoutStateRefundPegOutSucceeded,
		UserRskTxHash:      retainedQuote.UserRskTxHash,
		LpBtcTxHash:        retainedQuote.LpBtcTxHash,
		RefundPegoutTxHash: refundPegoutTxHash,
	}
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, expectedRetained).Return(nil).Once()
	quoteRepository.On("GetQuote", test.AnyCtx, retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
	lbc := new(mocks.LbcMock)
	lbc.On("RefundPegout", mock.Anything, mock.Anything).Return(refundPegoutTxHash, nil).Once()
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutQuoteCompletedEvent) bool {
		expected := retainedQuote
		expected.RefundPegoutTxHash = refundPegoutTxHash
		expected.State = quote.PegoutStateRefundPegOutSucceeded
		require.NoError(t, event.Error)
		return assert.Equal(t, expected, event.RetainedQuote) && assert.Equal(t, quote.PegoutQuoteCompletedEventId, event.Event.Id())
	})).Return().Once()
	btc := new(mocks.BtcRpcMock)
	btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(btcTxInfoMock, nil).Once()
	btc.On("BuildMerkleBranch", mock.Anything).Return(merkleBranchMock, nil)
	btc.On("GetRawTransaction", mock.Anything).Return(btcRawTxMock, nil)
	btc.On("GetTransactionBlockInfo", mock.Anything).Return(btcBlockInfoMock, nil)
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()

	contracts := blockchain.RskContracts{Lbc: lbc}
	rpc := blockchain.Rpc{Btc: btc}
	useCase := pegout.NewRefundPegoutUseCase(quoteRepository, contracts, eventBus, rpc, mutex)
	err := useCase.Run(context.Background(), retainedQuote)
	quoteRepository.AssertExpectations(t)
	lbc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	btc.AssertExpectations(t)
	require.NoError(t, err)
}

func TestRefundPegoutUseCase_Run_UpdateError(t *testing.T) {
	updateError := errors.New("an update error")
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, mock.Anything).Return(updateError).Once()
	quoteRepository.On("GetQuote", test.AnyCtx, retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
	lbc := new(mocks.LbcMock)
	lbc.On("RefundPegout", mock.Anything, mock.Anything).Return(refundPegoutTxHash, nil).Once()
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutQuoteCompletedEvent) bool {
		expected := retainedQuote
		expected.RefundPegoutTxHash = refundPegoutTxHash
		expected.State = quote.PegoutStateRefundPegOutSucceeded
		require.NoError(t, event.Error)
		return assert.Equal(t, expected, event.RetainedQuote) && assert.Equal(t, quote.PegoutQuoteCompletedEventId, event.Event.Id())
	})).Return().Once()
	btc := new(mocks.BtcRpcMock)
	btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(btcTxInfoMock, nil).Once()
	btc.On("BuildMerkleBranch", mock.Anything).Return(merkleBranchMock, nil)
	btc.On("GetRawTransaction", mock.Anything).Return(btcRawTxMock, nil)
	btc.On("GetTransactionBlockInfo", mock.Anything).Return(btcBlockInfoMock, nil)
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()

	contracts := blockchain.RskContracts{Lbc: lbc}
	rpc := blockchain.Rpc{Btc: btc}
	useCase := pegout.NewRefundPegoutUseCase(quoteRepository, contracts, eventBus, rpc, mutex)
	err := useCase.Run(context.Background(), retainedQuote)

	quoteRepository.AssertExpectations(t)
	lbc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	btc.AssertExpectations(t)
	require.ErrorIs(t, err, updateError)
}

func TestRefundPegoutUseCase_Run_NotPublishRecoverableError(t *testing.T) {
	recoverableSetups := []func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock){
		func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, retainedQuote.QuoteHash).Return(nil, assert.AnError).Once()
		},
		func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(blockchain.BitcoinTransactionInformation{}, assert.AnError).Once()
		},
		func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(btcTxInfoMock, nil).Once()
			btc.On("BuildMerkleBranch", mock.Anything).Return(blockchain.MerkleBranch{}, assert.AnError).Once()
		},
		func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(btcTxInfoMock, nil).Once()
			btc.On("BuildMerkleBranch", mock.Anything).Return(merkleBranchMock, nil).Once()
			btc.On("GetTransactionBlockInfo", mock.Anything).Return(blockchain.BitcoinBlockInformation{}, assert.AnError).Once()
		},
		func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(btcTxInfoMock, nil).Once()
			btc.On("BuildMerkleBranch", mock.Anything).Return(merkleBranchMock, nil).Once()
			btc.On("GetTransactionBlockInfo", mock.Anything).Return(btcBlockInfoMock, nil).Once()
			btc.On("GetRawTransaction", mock.Anything).Return(nil, assert.AnError).Once()

		},
		func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(btcTxInfoMock, nil).Once()
			btc.On("BuildMerkleBranch", mock.Anything).Return(merkleBranchMock, nil).Once()
			btc.On("GetTransactionBlockInfo", mock.Anything).Return(btcBlockInfoMock, nil).Once()
			btc.On("GetRawTransaction", mock.Anything).Return(btcRawTxMock, nil).Once()
			lbc.On("RefundPegout", mock.Anything, mock.Anything).Return("", blockchain.WaitingForBridgeError).Once()
		},
	}
	for _, setup := range recoverableSetups {
		eventBus := new(mocks.EventBusMock)
		mutex := new(mocks.MutexMock)
		mutex.On("Lock").Return()
		mutex.On("Unlock").Return()
		quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
		lbc := new(mocks.LbcMock)
		btc := new(mocks.BtcRpcMock)
		setup(quoteRepository, lbc, btc)
		contracts := blockchain.RskContracts{Lbc: lbc}
		rpc := blockchain.Rpc{Btc: btc}
		useCase := pegout.NewRefundPegoutUseCase(quoteRepository, contracts, eventBus, rpc, mutex)
		err := useCase.Run(context.Background(), retainedQuote)
		lbc.AssertExpectations(t)
		btc.AssertExpectations(t)
		quoteRepository.AssertExpectations(t)
		eventBus.AssertNotCalled(t, "Publish", mock.Anything)
		require.Error(t, err)
	}
}

func TestRefundPegoutUseCase_Run_PublishUnrecoverableError(t *testing.T) {
	unrecoverableSetups := []func(retainedQuote *quote.RetainedPegoutQuote, quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock){
		func(retainedQuote *quote.RetainedPegoutQuote, quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, retainedQuote.QuoteHash).Return(nil, nil).Once()
		},
		func(retainedQuote *quote.RetainedPegoutQuote, quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock) {
			retainedQuote.QuoteHash = "no hex"
			quoteRepository.On("GetQuote", test.AnyCtx, retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(btcTxInfoMock, nil).Once()
			btc.On("BuildMerkleBranch", mock.Anything).Return(merkleBranchMock, nil).Once()
			btc.On("GetTransactionBlockInfo", mock.Anything).Return(btcBlockInfoMock, nil).Once()
			btc.On("GetRawTransaction", mock.Anything).Return(btcRawTxMock, nil).Once()
		},
		func(retainedQuote *quote.RetainedPegoutQuote, quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(btcTxInfoMock, nil).Once()
			btc.On("BuildMerkleBranch", mock.Anything).Return(merkleBranchMock, nil).Once()
			btc.On("GetTransactionBlockInfo", mock.Anything).Return(btcBlockInfoMock, nil).Once()
			btc.On("GetRawTransaction", mock.Anything).Return(btcRawTxMock, nil).Once()
			lbc.On("RefundPegout", mock.Anything, mock.Anything).Return("", assert.AnError).Once()
		},
	}

	for _, setup := range unrecoverableSetups {
		quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
		lbc := new(mocks.LbcMock)
		btc := new(mocks.BtcRpcMock)
		caseQuote := retainedQuote
		setup(&caseQuote, quoteRepository, lbc, btc)
		mutex := new(mocks.MutexMock)
		mutex.On("Lock").Return()
		mutex.On("Unlock").Return()
		eventBus := new(mocks.EventBusMock)
		eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutQuoteCompletedEvent) bool {
			require.Error(t, event.Error)
			return assert.Equal(t, caseQuote.LpBtcTxHash, event.RetainedQuote.LpBtcTxHash) && assert.Equal(t, caseQuote.Signature, event.RetainedQuote.Signature) &&
				assert.Equal(t, caseQuote.QuoteHash, event.RetainedQuote.QuoteHash) && assert.Equal(t, caseQuote.DepositAddress, event.RetainedQuote.DepositAddress) &&
				assert.Equal(t, caseQuote.RequiredLiquidity, event.RetainedQuote.RequiredLiquidity) && assert.Equal(t, quote.PegoutStateRefundPegOutFailed, event.RetainedQuote.State) &&
				assert.Equal(t, caseQuote.UserRskTxHash, event.RetainedQuote.UserRskTxHash) && assert.Equal(t, quote.PegoutQuoteCompletedEventId, event.Event.Id())
		})).Return().Once()
		quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, mock.MatchedBy(
			func(q quote.RetainedPegoutQuote) bool {
				expected := caseQuote
				expected.State = quote.PegoutStateRefundPegOutFailed
				return assert.Equal(t, expected, q)
			})).Return(nil).Once()
		contracts := blockchain.RskContracts{Lbc: lbc}
		rpc := blockchain.Rpc{Btc: btc}
		useCase := pegout.NewRefundPegoutUseCase(quoteRepository, contracts, eventBus, rpc, mutex)
		err := useCase.Run(context.Background(), caseQuote)
		lbc.AssertExpectations(t)
		btc.AssertExpectations(t)
		quoteRepository.AssertExpectations(t)
		eventBus.AssertExpectations(t)
		require.Error(t, err)
	}
}

func TestRefundPegoutUseCase_Run_NoConfirmations(t *testing.T) {
	unconfirmedBlockInfo := btcTxInfoMock
	unconfirmedBlockInfo.Confirmations = 1
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
	lbc := new(mocks.LbcMock)
	eventBus := new(mocks.EventBusMock)
	btc := new(mocks.BtcRpcMock)
	btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(unconfirmedBlockInfo, nil).Once()
	mutex := new(mocks.MutexMock)

	contracts := blockchain.RskContracts{Lbc: lbc}
	rpc := blockchain.Rpc{Btc: btc}
	useCase := pegout.NewRefundPegoutUseCase(quoteRepository, contracts, eventBus, rpc, mutex)
	err := useCase.Run(context.Background(), retainedQuote)

	quoteRepository.AssertExpectations(t)
	btc.AssertExpectations(t)
	eventBus.AssertNotCalled(t, "Publish")
	lbc.AssertNotCalled(t, "RefundPegout")
	lbc.AssertNotCalled(t, "GetAddress")
	mutex.AssertNotCalled(t, "Lock")
	mutex.AssertNotCalled(t, "Unlock")
	require.ErrorIs(t, err, usecases.NoEnoughConfirmationsError)
}

func TestRefundPegoutUseCase_Run_WrongState(t *testing.T) {
	wrongStateQuote := retainedQuote
	wrongStateQuote.State = quote.PegoutStateSendPegoutFailed
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	lbc := new(mocks.LbcMock)
	eventBus := new(mocks.EventBusMock)
	btc := new(mocks.BtcRpcMock)
	mutex := new(mocks.MutexMock)

	contracts := blockchain.RskContracts{Lbc: lbc}
	rpc := blockchain.Rpc{Btc: btc}
	useCase := pegout.NewRefundPegoutUseCase(quoteRepository, contracts, eventBus, rpc, mutex)

	err := useCase.Run(context.Background(), wrongStateQuote)

	quoteRepository.AssertNotCalled(t, "UpdateRetainedQuote")
	quoteRepository.AssertNotCalled(t, "GetQuote")
	quoteRepository.AssertNotCalled(t, "GetRetainedQuote")
	btc.AssertNotCalled(t, "GetTransactionInfo")
	eventBus.AssertNotCalled(t, "Publish")
	lbc.AssertNotCalled(t, "RefundPegout")
	lbc.AssertNotCalled(t, "GetAddress")
	mutex.AssertNotCalled(t, "Lock")
	mutex.AssertNotCalled(t, "Unlock")
	require.ErrorIs(t, err, usecases.WrongStateError)
}

func TestRefundPegoutUseCase_Run_RegisterCoinbase(t *testing.T) {
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	lbc := new(mocks.LbcMock)
	bridge := new(mocks.BridgeMock)
	eventBus := new(mocks.EventBusMock)
	btc := new(mocks.BtcRpcMock)
	mutex := new(mocks.MutexMock)
	coinbaseInfo := blockchain.BtcCoinbaseTransactionInformation{BlockHash: utils.To32Bytes(utils.MustGetRandomBytes(32))}
	// Mocks that don't change per test
	mutex.On("Lock").Return().Times(3)
	mutex.On("Unlock").Return().Times(3)
	quoteRepository.EXPECT().UpdateRetainedQuote(test.AnyCtx, mock.Anything).Return(nil).Twice()
	quoteRepository.EXPECT().GetQuote(test.AnyCtx, retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Times(3)
	tx := btcTxInfoMock
	tx.HasWitness = true
	btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(tx, nil).Times(3)
	btc.On("GetCoinbaseInformation", retainedQuote.LpBtcTxHash).Return(coinbaseInfo, nil).Times(3)
	btc.On("BuildMerkleBranch", mock.Anything).Return(merkleBranchMock, nil).Times(3)
	btc.On("GetRawTransaction", mock.Anything).Return(btcRawTxMock, nil).Times(3)
	btc.On("GetTransactionBlockInfo", mock.Anything).Return(btcBlockInfoMock, nil).Times(3)
	// once as it'll be called only on 1st test
	lbc.On("RefundPegout", mock.Anything, mock.Anything).Return(refundPegoutTxHash, nil).Once()

	contracts := blockchain.RskContracts{Lbc: lbc, Bridge: bridge}
	rpc := blockchain.Rpc{Btc: btc}
	useCase := pegout.NewRefundPegoutUseCase(quoteRepository, contracts, eventBus, rpc, mutex)
	t.Run("Should call RegisterCoinbaseTransaction", func(t *testing.T) {
		bridge.On("RegisterBtcCoinbaseTransaction", coinbaseInfo).Return(test.AnyHash, nil).Once()
		eventBus.On("Publish", mock.MatchedBy(func(e quote.PegoutQuoteCompletedEvent) bool {
			return e.Error == nil
		})).Return().Once()
		err := useCase.Run(context.Background(), retainedQuote)
		require.NoError(t, err)
	})
	t.Run("Should return recoverable error if tx wasn't registered due to waiting for the bridge", func(t *testing.T) {
		bridge.On("RegisterBtcCoinbaseTransaction", coinbaseInfo).Return("", blockchain.WaitingForBridgeError).Once()
		err := useCase.Run(context.Background(), retainedQuote)
		require.Error(t, err)
		require.NotErrorIs(t, err, usecases.NonRecoverableError)
	})
	t.Run("Should return non recoverable error if tx wasn't registered due to any other error", func(t *testing.T) {
		bridge.On("RegisterBtcCoinbaseTransaction", coinbaseInfo).Return("", assert.AnError).Once()
		eventBus.On("Publish", mock.MatchedBy(func(e quote.PegoutQuoteCompletedEvent) bool {
			return errors.Is(e.Error, usecases.NonRecoverableError)
		})).Return().Once()
		err := useCase.Run(context.Background(), retainedQuote)
		require.ErrorIs(t, err, usecases.NonRecoverableError)
	})
	mutex.AssertExpectations(t)
	lbc.AssertExpectations(t)
	bridge.AssertExpectations(t)
	btc.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}
