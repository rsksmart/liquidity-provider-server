package pegout_test

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
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
	bridgeTxHash       = "0x2b2d2f"
)

var retainedQuote = quote.RetainedPegoutQuote{
	QuoteHash:          "1c2d3f",
	DepositAddress:     "0x654321",
	Signature:          "0x112a3b",
	RequiredLiquidity:  entities.NewWei(1000),
	State:              quote.PegoutStateSendPegoutSucceeded,
	UserRskTxHash:      "0x3c2b1a",
	LpBtcTxHash:        "0x3c2b1a",
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
	Hash:          "0x1c2b3a",
	Confirmations: 11,
	Outputs:       map[string][]*entities.Wei{test.AnyAddress: {entities.NewWei(1000)}},
}

var btcRawTxMock = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

func TestRefundPegoutUseCase_Run(t *testing.T) {
	bridgeAddress := "0x1234"
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), quote.RetainedPegoutQuote{
		QuoteHash:          retainedQuote.QuoteHash,
		DepositAddress:     retainedQuote.DepositAddress,
		Signature:          retainedQuote.Signature,
		RequiredLiquidity:  retainedQuote.RequiredLiquidity,
		State:              quote.PegoutStateRefundPegOutSucceeded,
		UserRskTxHash:      retainedQuote.UserRskTxHash,
		LpBtcTxHash:        retainedQuote.LpBtcTxHash,
		RefundPegoutTxHash: refundPegoutTxHash,
	}).Return(nil).Once()
	quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), quote.RetainedPegoutQuote{
		QuoteHash:          retainedQuote.QuoteHash,
		DepositAddress:     retainedQuote.DepositAddress,
		Signature:          retainedQuote.Signature,
		RequiredLiquidity:  retainedQuote.RequiredLiquidity,
		State:              quote.PegoutStateRefundPegOutSucceeded,
		UserRskTxHash:      retainedQuote.UserRskTxHash,
		LpBtcTxHash:        retainedQuote.LpBtcTxHash,
		RefundPegoutTxHash: refundPegoutTxHash,
		BridgeRefundTxHash: bridgeTxHash,
	}).Return(nil).Once()
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
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
	rskWallet := new(mocks.RskWalletMock)
	rskWallet.On("SendRbtc", mock.AnythingOfType("context.backgroundCtx"),
		blockchain.NewTransactionConfig(entities.NewWei(8000), 100000, entities.NewWei(60000000)),
		bridgeAddress).Return(bridgeTxHash, nil).Once()
	bridge := new(mocks.BridgeMock)
	bridge.On("GetAddress").Return(bridgeAddress).Once()
	mutex := new(mocks.MutexMock)
	mutex.On("Unlock").Return().Once()
	mutex.On("Lock").Return().Once()

	useCase := pegout.NewRefundPegoutUseCase(quoteRepository, lbc, eventBus, btc, rskWallet, bridge, mutex)
	err := useCase.Run(context.Background(), retainedQuote)

	quoteRepository.AssertExpectations(t)
	lbc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	btc.AssertExpectations(t)
	rskWallet.AssertExpectations(t)
	bridge.AssertExpectations(t)
	mutex.AssertExpectations(t)
	require.NoError(t, err)
}

func TestRefundPegoutUseCase_Run_UpdateError(t *testing.T) {
	updateError := errors.New("an update error")
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil).Once()
	quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(updateError).Once()
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
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
	rskWallet := new(mocks.RskWalletMock)
	rskWallet.On("SendRbtc", mock.AnythingOfType("context.backgroundCtx"), mock.Anything, mock.Anything).Return(bridgeTxHash, nil).Once()
	bridge := new(mocks.BridgeMock)
	bridge.On("GetAddress").Return("an address").Once()
	mutex := new(mocks.MutexMock)
	mutex.On("Unlock").Return().Once()
	mutex.On("Lock").Return().Once()

	useCase := pegout.NewRefundPegoutUseCase(quoteRepository, lbc, eventBus, btc, rskWallet, bridge, mutex)
	err := useCase.Run(context.Background(), retainedQuote)

	quoteRepository.AssertExpectations(t)
	lbc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	btc.AssertExpectations(t)
	rskWallet.AssertExpectations(t)
	bridge.AssertExpectations(t)
	mutex.AssertExpectations(t)
	require.ErrorIs(t, err, updateError)
}

func TestRefundPegoutUseCase_Run_NotPublishRecoverableError(t *testing.T) {
	bridge := new(mocks.BridgeMock)
	bridge.On("GetAddress").Return("0x1234").Once()
	mutex := new(mocks.MutexMock)
	mutex.On("Unlock").Return()
	mutex.On("Lock").Return()

	recoverableSetups := []func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock, rskWallet *mocks.RskWalletMock){
		func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock, rskWallet *mocks.RskWalletMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(nil, assert.AnError).Once()
		},
		func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock, rskWallet *mocks.RskWalletMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(blockchain.BitcoinTransactionInformation{}, assert.AnError).Once()
		},
		func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock, rskWallet *mocks.RskWalletMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(btcTxInfoMock, nil).Once()
			btc.On("BuildMerkleBranch", mock.Anything).Return(blockchain.MerkleBranch{}, assert.AnError).Once()
		},
		func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock, rskWallet *mocks.RskWalletMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(btcTxInfoMock, nil).Once()
			btc.On("BuildMerkleBranch", mock.Anything).Return(merkleBranchMock, nil).Once()
			btc.On("GetTransactionBlockInfo", mock.Anything).Return(blockchain.BitcoinBlockInformation{}, assert.AnError).Once()
		},
		func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock, rskWallet *mocks.RskWalletMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(btcTxInfoMock, nil).Once()
			btc.On("BuildMerkleBranch", mock.Anything).Return(merkleBranchMock, nil).Once()
			btc.On("GetTransactionBlockInfo", mock.Anything).Return(btcBlockInfoMock, nil).Once()
			btc.On("GetRawTransaction", mock.Anything).Return(nil, assert.AnError).Once()

		},
		func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock, rskWallet *mocks.RskWalletMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(btcTxInfoMock, nil).Once()
			btc.On("BuildMerkleBranch", mock.Anything).Return(merkleBranchMock, nil).Once()
			btc.On("GetTransactionBlockInfo", mock.Anything).Return(btcBlockInfoMock, nil).Once()
			btc.On("GetRawTransaction", mock.Anything).Return(btcRawTxMock, nil).Once()
			lbc.On("RefundPegout", mock.Anything, mock.Anything).Return("", blockchain.WaitingForBridgeError).Once()
		},
	}

	for _, setup := range recoverableSetups {
		eventBus := new(mocks.EventBusMock)
		quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
		lbc := new(mocks.LbcMock)
		btc := new(mocks.BtcRpcMock)
		rskWallet := new(mocks.RskWalletMock)
		setup(quoteRepository, lbc, btc, rskWallet)
		useCase := pegout.NewRefundPegoutUseCase(quoteRepository, lbc, eventBus, btc, rskWallet, bridge, mutex)
		err := useCase.Run(context.Background(), retainedQuote)
		lbc.AssertExpectations(t)
		btc.AssertExpectations(t)
		rskWallet.AssertExpectations(t)
		quoteRepository.AssertExpectations(t)
		eventBus.AssertNotCalled(t, "Publish", mock.Anything)
		require.Error(t, err)
	}
}

func TestRefundPegoutUseCase_Run_PublishUnrecoverableError(t *testing.T) {
	bridge := new(mocks.BridgeMock)
	bridge.On("GetAddress").Return("0x1234").Once()
	mutex := new(mocks.MutexMock)
	mutex.On("Unlock").Return()
	mutex.On("Lock").Return()

	unrecoverableSetups := []func(retainedQuote *quote.RetainedPegoutQuote, quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock, rskWallet *mocks.RskWalletMock){
		func(retainedQuote *quote.RetainedPegoutQuote, quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock, rskWallet *mocks.RskWalletMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(nil, nil).Once()
		},
		func(retainedQuote *quote.RetainedPegoutQuote, quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock, rskWallet *mocks.RskWalletMock) {
			retainedQuote.QuoteHash = "no hex"
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
			btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(btcTxInfoMock, nil).Once()
			btc.On("BuildMerkleBranch", mock.Anything).Return(merkleBranchMock, nil).Once()
			btc.On("GetTransactionBlockInfo", mock.Anything).Return(btcBlockInfoMock, nil).Once()
			btc.On("GetRawTransaction", mock.Anything).Return(btcRawTxMock, nil).Once()
		},
		func(retainedQuote *quote.RetainedPegoutQuote, quoteRepository *mocks.PegoutQuoteRepositoryMock, lbc *mocks.LbcMock, btc *mocks.BtcRpcMock, rskWallet *mocks.RskWalletMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
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
		rskWallet := new(mocks.RskWalletMock)
		caseQuote := retainedQuote
		setup(&caseQuote, quoteRepository, lbc, btc, rskWallet)
		eventBus := new(mocks.EventBusMock)
		eventBus.On("Publish", mock.MatchedBy(func(event quote.PegoutQuoteCompletedEvent) bool {
			require.Error(t, event.Error)
			return assert.Equal(t, caseQuote.LpBtcTxHash, event.RetainedQuote.LpBtcTxHash) && assert.Equal(t, caseQuote.Signature, event.RetainedQuote.Signature) &&
				assert.Equal(t, caseQuote.QuoteHash, event.RetainedQuote.QuoteHash) && assert.Equal(t, caseQuote.DepositAddress, event.RetainedQuote.DepositAddress) &&
				assert.Equal(t, caseQuote.RequiredLiquidity, event.RetainedQuote.RequiredLiquidity) && assert.Equal(t, quote.PegoutStateRefundPegOutFailed, event.RetainedQuote.State) &&
				assert.Equal(t, caseQuote.UserRskTxHash, event.RetainedQuote.UserRskTxHash) && assert.Equal(t, quote.PegoutQuoteCompletedEventId, event.Event.Id())
		})).Return().Once()
		quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.MatchedBy(
			func(q quote.RetainedPegoutQuote) bool {
				expected := caseQuote
				expected.State = quote.PegoutStateRefundPegOutFailed
				return assert.Equal(t, expected, q)
			})).Return(nil).Once()
		useCase := pegout.NewRefundPegoutUseCase(quoteRepository, lbc, eventBus, btc, rskWallet, bridge, mutex)
		err := useCase.Run(context.Background(), caseQuote)
		lbc.AssertExpectations(t)
		btc.AssertExpectations(t)
		rskWallet.AssertExpectations(t)
		quoteRepository.AssertExpectations(t)
		eventBus.AssertExpectations(t)
		require.Error(t, err)
	}
}

func TestRefundPegoutUseCase_Run_NoConfirmations(t *testing.T) {
	unconfirmedBlockInfo := btcTxInfoMock
	unconfirmedBlockInfo.Confirmations = 1
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&pegoutQuote, nil).Once()
	lbc := new(mocks.LbcMock)
	eventBus := new(mocks.EventBusMock)
	btc := new(mocks.BtcRpcMock)
	btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(unconfirmedBlockInfo, nil).Once()
	rskWallet := new(mocks.RskWalletMock)
	bridge := new(mocks.BridgeMock)
	mutex := new(mocks.MutexMock)

	useCase := pegout.NewRefundPegoutUseCase(quoteRepository, lbc, eventBus, btc, rskWallet, bridge, mutex)

	err := useCase.Run(context.Background(), retainedQuote)

	quoteRepository.AssertExpectations(t)
	btc.AssertExpectations(t)
	eventBus.AssertNotCalled(t, "Publish")
	rskWallet.AssertNotCalled(t, "SendRbtc")
	bridge.AssertNotCalled(t, "GetAddress")
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
	rskWallet := new(mocks.RskWalletMock)
	bridge := new(mocks.BridgeMock)
	mutex := new(mocks.MutexMock)

	useCase := pegout.NewRefundPegoutUseCase(quoteRepository, lbc, eventBus, btc, rskWallet, bridge, mutex)

	err := useCase.Run(context.Background(), wrongStateQuote)

	quoteRepository.AssertNotCalled(t, "UpdateRetainedQuote")
	quoteRepository.AssertNotCalled(t, "GetQuote")
	quoteRepository.AssertNotCalled(t, "GetRetainedQuote")
	btc.AssertNotCalled(t, "GetTransactionInfo")
	eventBus.AssertNotCalled(t, "Publish")
	rskWallet.AssertNotCalled(t, "SendRbtc")
	bridge.AssertNotCalled(t, "GetAddress")
	lbc.AssertNotCalled(t, "RefundPegout")
	lbc.AssertNotCalled(t, "GetAddress")
	mutex.AssertNotCalled(t, "Lock")
	mutex.AssertNotCalled(t, "Unlock")
	require.ErrorIs(t, err, usecases.WrongStateError)
}

func TestRefundPegoutUseCase_Run_CorrectBridgeAmount(t *testing.T) {
	bridgeAddress := "0x1234"
	lbc := new(mocks.LbcMock)
	lbc.On("RefundPegout", mock.Anything, mock.Anything).Return(refundPegoutTxHash, nil)
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.Anything)
	btc := new(mocks.BtcRpcMock)
	btc.On("GetTransactionInfo", retainedQuote.LpBtcTxHash).Return(btcTxInfoMock, nil)
	btc.On("BuildMerkleBranch", mock.Anything).Return(merkleBranchMock, nil)
	btc.On("GetRawTransaction", mock.Anything).Return(btcRawTxMock, nil)
	btc.On("GetTransactionBlockInfo", mock.Anything).Return(btcBlockInfoMock, nil)
	bridge := new(mocks.BridgeMock)
	bridge.On("GetAddress").Return(bridgeAddress)
	mutex := new(mocks.MutexMock)
	mutex.On("Unlock")
	mutex.On("Lock")

	cases := getQuotesWithExpectedTotalTable()

	for _, c := range cases {
		quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
		q := c.Value()
		quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil)
		quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil)
		quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&q, nil)
		quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote.QuoteHash).Return(&q, nil)

		rskWallet := new(mocks.RskWalletMock)
		rskWallet.On("SendRbtc", mock.AnythingOfType("context.backgroundCtx"),
			blockchain.NewTransactionConfig(c.Result, 100000, entities.NewWei(60000000)),
			bridgeAddress).
			Return(bridgeTxHash, nil).Once()

		useCase := pegout.NewRefundPegoutUseCase(quoteRepository, lbc, eventBus, btc, rskWallet, bridge, mutex)
		err := useCase.Run(context.Background(), retainedQuote)
		quoteRepository.AssertExpectations(t)
		rskWallet.AssertExpectations(t)
		require.NoError(t, err)
	}
}

func getQuotesWithExpectedTotalTable() test.Table[func() quote.PegoutQuote, *entities.Wei] {
	return test.Table[func() quote.PegoutQuote, *entities.Wei]{
		{
			Value: func() quote.PegoutQuote {
				testQuote := pegoutQuote
				testQuote.Value = entities.NewWei(3000)
				testQuote.GasFee = entities.NewWei(3000)
				testQuote.ProductFeeAmount = 1000
				testQuote.CallFee = entities.NewWei(2000)
				return testQuote
			},
			Result: entities.NewWei(8000),
		},
		{
			Value: func() quote.PegoutQuote {
				testQuote := pegoutQuote
				testQuote.Value = entities.NewWei(3000)
				testQuote.GasFee = entities.NewWei(3000)
				testQuote.ProductFeeAmount = 1000
				testQuote.CallFee = entities.NewWei(2000)
				return testQuote
			},
			Result: entities.NewWei(8000),
		},
		{
			Value: func() quote.PegoutQuote {
				testQuote := pegoutQuote
				testQuote.Value = entities.NewWei(0)
				testQuote.GasFee = entities.NewWei(0)
				testQuote.ProductFeeAmount = 1
				testQuote.CallFee = entities.NewWei(0)
				return testQuote
			},
			Result: entities.NewWei(0),
		},
		{
			Value: func() quote.PegoutQuote {
				testQuote := pegoutQuote
				testQuote.Value = entities.NewWei(15000)
				testQuote.GasFee = entities.NewWei(1)
				testQuote.ProductFeeAmount = 1
				testQuote.CallFee = entities.NewWei(500)
				return testQuote
			},
			Result: entities.NewWei(15501),
		},
	}
}
