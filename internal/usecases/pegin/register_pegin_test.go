package pegin_test

import (
	"context"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

var btcRawTxMock = []byte{50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60}
var pmtMock = []byte{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
var btcBlockInfoMock = blockchain.BitcoinBlockInformation{
	Hash:   [32]byte{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
	Height: big.NewInt(200),
}

func TestRegisterPeginUseCase_Run(t *testing.T) {
	registerPeginTx := "register tx hash"
	retainedPeginQuote := quote.RetainedPeginQuote{
		QuoteHash:         "101b1c",
		DepositAddress:    "any address",
		Signature:         "0102031f1b",
		RequiredLiquidity: entities.NewWei(1500),
		State:             quote.PeginStateCallForUserSucceeded,
		UserBtcTxHash:     "btc tx hash",
		CallForUserTxHash: "cfu tx hash",
	}
	expectedRetainedQuote := retainedPeginQuote
	expectedRetainedQuote.State = quote.PeginStateRegisterPegInSucceeded
	expectedRetainedQuote.RegisterPeginTxHash = registerPeginTx

	lbc := new(test.LbcMock)
	lbc.On("RegisterPegin", blockchain.RegisterPeginParams{
		QuoteSignature:        []byte{1, 2, 3, 31, 27},
		BitcoinRawTransaction: btcRawTxMock,
		PartialMerkleTree:     pmtMock,
		BlockHeight:           btcBlockInfoMock.Height,
		Quote:                 testPeginQuote,
	}).Return(registerPeginTx, nil).Once()
	quoteRepository := new(test.PeginQuoteRepositoryMock)
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedPeginQuote.QuoteHash).Return(&testPeginQuote, nil).Once()
	quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.MatchedBy(func(q quote.RetainedPeginQuote) bool {
		return assert.Equal(t, expectedRetainedQuote, q)
	})).Return(nil).Once()
	eventBus := new(test.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.RegisterPeginCompletedEvent) bool {
		require.NoError(t, event.Error)
		return assert.Equal(t, expectedRetainedQuote, event.RetainedQuote) && assert.Equal(t, quote.RegisterPeginCompletedEventId, event.Event.Id())
	})).Return().Once()
	bridge := new(test.BridgeMock)
	bridge.On("GetRequiredTxConfirmations").Return(uint64(10)).Once()
	btc := new(test.BtcRpcMock)
	btc.On("GetTransactionInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
		Hash:          retainedPeginQuote.UserBtcTxHash,
		Confirmations: 11,
		Outputs:       map[string][]*entities.Wei{retainedPeginQuote.DepositAddress: {entities.NewWei(30012)}},
	}, nil).Once()
	btc.On("GetRawTransaction", retainedPeginQuote.UserBtcTxHash).Return(btcRawTxMock, nil).Once()
	btc.On("GetPartialMerkleTree", retainedPeginQuote.UserBtcTxHash).Return(pmtMock, nil).Once()
	btc.On("GetTransactionBlockInfo", retainedPeginQuote.UserBtcTxHash).Return(btcBlockInfoMock, nil)
	mutex := new(test.MutexMock)
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()

	useCase := pegin.NewRegisterPeginUseCase(lbc, quoteRepository, eventBus, bridge, btc, mutex)
	err := useCase.Run(context.Background(), retainedPeginQuote)

	require.NoError(t, err)
	lbc.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	bridge.AssertExpectations(t)
	btc.AssertExpectations(t)
	mutex.AssertExpectations(t)
}

func TestRegisterPeginUseCase_Run_DontPublishRecoverableErrors(t *testing.T) {
	retainedPeginQuote := quote.RetainedPeginQuote{
		QuoteHash:         "101b1c",
		DepositAddress:    "any address",
		Signature:         "0102031f1b",
		RequiredLiquidity: entities.NewWei(1500),
		State:             quote.PeginStateCallForUserSucceeded,
		UserBtcTxHash:     "btc tx hash",
		CallForUserTxHash: "cfu tx hash",
	}

	setups := registerPeginRecoverableErrorSetups()

	for _, setup := range setups {
		lbc := new(test.LbcMock)
		quoteRepository := new(test.PeginQuoteRepositoryMock)
		btc := new(test.BtcRpcMock)
		eventBus := new(test.EventBusMock)

		bridge := new(test.BridgeMock)
		bridge.On("GetRequiredTxConfirmations").Return(uint64(10))

		mutex := new(test.MutexMock)
		mutex.On("Lock").Return()
		mutex.On("Unlock").Return()

		caseQuote := retainedPeginQuote
		setup(&caseQuote, lbc, quoteRepository, btc)
		useCase := pegin.NewRegisterPeginUseCase(lbc, quoteRepository, eventBus, bridge, btc, mutex)
		err := useCase.Run(context.Background(), caseQuote)

		require.Error(t, err)
		eventBus.AssertNotCalled(t, "Publish")
		lbc.AssertExpectations(t)
		quoteRepository.AssertExpectations(t)
		btc.AssertExpectations(t)
	}
}

func registerPeginRecoverableErrorSetups() []func(caseQuote *quote.RetainedPeginQuote, lbc *test.LbcMock, quoteRepository *test.PeginQuoteRepositoryMock, btc *test.BtcRpcMock) {
	return []func(caseQuote *quote.RetainedPeginQuote, lbc *test.LbcMock, quoteRepository *test.PeginQuoteRepositoryMock, btc *test.BtcRpcMock){
		func(caseQuote *quote.RetainedPeginQuote, lbc *test.LbcMock, quoteRepository *test.PeginQuoteRepositoryMock, btc *test.BtcRpcMock) {
			caseQuote.State = quote.PeginStateWaitingForDeposit
		},
		func(caseQuote *quote.RetainedPeginQuote, lbc *test.LbcMock, quoteRepository *test.PeginQuoteRepositoryMock, btc *test.BtcRpcMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), caseQuote.QuoteHash).
				Return(nil, assert.AnError).Once()
		},
		func(caseQuote *quote.RetainedPeginQuote, lbc *test.LbcMock, quoteRepository *test.PeginQuoteRepositoryMock, btc *test.BtcRpcMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), caseQuote.QuoteHash).
				Return(&testPeginQuote, nil).Once()
			btc.On("GetTransactionInfo", caseQuote.UserBtcTxHash).
				Return(blockchain.BitcoinTransactionInformation{}, assert.AnError).Once()
		},
		func(caseQuote *quote.RetainedPeginQuote, lbc *test.LbcMock, quoteRepository *test.PeginQuoteRepositoryMock, btc *test.BtcRpcMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), caseQuote.QuoteHash).
				Return(&testPeginQuote, nil).Once()
			btc.On("GetTransactionInfo", caseQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
				Hash:          caseQuote.UserBtcTxHash,
				Confirmations: 11,
				Outputs:       map[string][]*entities.Wei{caseQuote.DepositAddress: {entities.NewWei(30012)}},
			}, nil).Once()
			caseQuote.Signature = "malformed signature"
		},
		func(caseQuote *quote.RetainedPeginQuote, lbc *test.LbcMock, quoteRepository *test.PeginQuoteRepositoryMock, btc *test.BtcRpcMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), caseQuote.QuoteHash).
				Return(&testPeginQuote, nil).Once()
			btc.On("GetTransactionInfo", caseQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
				Hash:          caseQuote.UserBtcTxHash,
				Confirmations: 11,
				Outputs:       map[string][]*entities.Wei{caseQuote.DepositAddress: {entities.NewWei(30012)}},
			}, nil).Once()
			btc.On("GetRawTransaction", caseQuote.UserBtcTxHash).Return([]byte{}, assert.AnError).Once()
		},
		func(caseQuote *quote.RetainedPeginQuote, lbc *test.LbcMock, quoteRepository *test.PeginQuoteRepositoryMock, btc *test.BtcRpcMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), caseQuote.QuoteHash).
				Return(&testPeginQuote, nil).Once()
			btc.On("GetTransactionInfo", caseQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
				Hash:          caseQuote.UserBtcTxHash,
				Confirmations: 11,
				Outputs:       map[string][]*entities.Wei{caseQuote.DepositAddress: {entities.NewWei(30012)}},
			}, nil).Once()
			btc.On("GetRawTransaction", caseQuote.UserBtcTxHash).Return(btcRawTxMock, nil).Once()
			btc.On("GetPartialMerkleTree", caseQuote.UserBtcTxHash).Return([]byte{}, assert.AnError).Once()
		},
		func(caseQuote *quote.RetainedPeginQuote, lbc *test.LbcMock, quoteRepository *test.PeginQuoteRepositoryMock, btc *test.BtcRpcMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), caseQuote.QuoteHash).
				Return(&testPeginQuote, nil).Once()
			btc.On("GetTransactionInfo", caseQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
				Hash:          caseQuote.UserBtcTxHash,
				Confirmations: 11,
				Outputs:       map[string][]*entities.Wei{caseQuote.DepositAddress: {entities.NewWei(30012)}},
			}, nil).Once()
			btc.On("GetRawTransaction", caseQuote.UserBtcTxHash).Return(btcRawTxMock, nil).Once()
			btc.On("GetPartialMerkleTree", caseQuote.UserBtcTxHash).Return(pmtMock, nil).Once()
			btc.On("GetTransactionBlockInfo", caseQuote.UserBtcTxHash).Return(blockchain.BitcoinBlockInformation{}, assert.AnError)
		},
	}
}

func TestRegisterPeginUseCase_Run_QuoteNotFound(t *testing.T) {
	retainedPeginQuote := quote.RetainedPeginQuote{
		QuoteHash:         "101b1c",
		DepositAddress:    "any address",
		Signature:         "0102031f1b",
		RequiredLiquidity: entities.NewWei(1500),
		State:             quote.PeginStateCallForUserSucceeded,
		UserBtcTxHash:     "btc tx hash",
		CallForUserTxHash: "cfu tx hash",
	}

	quoteRepository := new(test.PeginQuoteRepositoryMock)
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedPeginQuote.QuoteHash).Return(nil, nil).Once()
	quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.MatchedBy(func(q quote.RetainedPeginQuote) bool {
		expected := retainedPeginQuote
		expected.State = quote.PeginStateRegisterPegInFailed
		return assert.Equal(t, expected, q)
	})).Return(nil).Once()

	eventBus := new(test.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.RegisterPeginCompletedEvent) bool {
		expected := retainedPeginQuote
		expected.State = quote.PeginStateRegisterPegInFailed
		require.ErrorIs(t, event.Error, usecases.QuoteNotFoundError)
		return assert.Equal(t, expected, event.RetainedQuote) &&
			assert.Equal(t, quote.RegisterPeginCompletedEventId, event.Event.Id())
	})).Return().Once()

	lbc := new(test.LbcMock)
	bridge := new(test.BridgeMock)
	btc := new(test.BtcRpcMock)
	mutex := new(test.MutexMock)

	useCase := pegin.NewRegisterPeginUseCase(lbc, quoteRepository, eventBus, bridge, btc, mutex)

	err := useCase.Run(context.Background(), retainedPeginQuote)

	require.ErrorIs(t, err, usecases.QuoteNotFoundError)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)

	lbc.AssertNotCalled(t, "RegisterPegin")
	bridge.AssertNotCalled(t, "GetRequiredTxConfirmations")
	btc.AssertNotCalled(t, "GetTransactionInfo")
	btc.AssertNotCalled(t, "GetRawTransaction")
	btc.AssertNotCalled(t, "GetPartialMerkleTree")
	btc.AssertNotCalled(t, "GetTransactionBlockInfo")
	mutex.AssertNotCalled(t, "Lock")
	mutex.AssertNotCalled(t, "Unlock")
}

func TestRegisterPeginUseCase_Run_RegisterPeginFailed(t *testing.T) {
	registerPeginTx := "register tx hash"
	retainedPeginQuote := quote.RetainedPeginQuote{
		QuoteHash:         "101b1c",
		DepositAddress:    "any address",
		Signature:         "0102031f1b",
		RequiredLiquidity: entities.NewWei(1500),
		State:             quote.PeginStateCallForUserSucceeded,
		UserBtcTxHash:     "btc tx hash",
		CallForUserTxHash: "cfu tx hash",
	}
	expectedRetainedQuote := retainedPeginQuote
	expectedRetainedQuote.State = quote.PeginStateRegisterPegInFailed
	expectedRetainedQuote.RegisterPeginTxHash = registerPeginTx

	lbc := new(test.LbcMock)
	lbc.On("RegisterPegin", blockchain.RegisterPeginParams{
		QuoteSignature:        []byte{1, 2, 3, 31, 27},
		BitcoinRawTransaction: btcRawTxMock,
		PartialMerkleTree:     pmtMock,
		BlockHeight:           btcBlockInfoMock.Height,
		Quote:                 testPeginQuote,
	}).Return(registerPeginTx, assert.AnError).Once()

	quoteRepository := new(test.PeginQuoteRepositoryMock)
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedPeginQuote.QuoteHash).Return(&testPeginQuote, nil).Once()
	quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.MatchedBy(func(q quote.RetainedPeginQuote) bool {
		return assert.Equal(t, expectedRetainedQuote, q)
	})).Return(nil).Once()

	eventBus := new(test.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.RegisterPeginCompletedEvent) bool {
		require.Error(t, event.Error)
		return assert.Equal(t, expectedRetainedQuote, event.RetainedQuote) && assert.Equal(t, quote.RegisterPeginCompletedEventId, event.Event.Id())
	})).Return().Once()
	bridge := new(test.BridgeMock)
	bridge.On("GetRequiredTxConfirmations").Return(uint64(10)).Once()

	btc := new(test.BtcRpcMock)
	btc.On("GetTransactionInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
		Hash:          retainedPeginQuote.UserBtcTxHash,
		Confirmations: 11,
		Outputs:       map[string][]*entities.Wei{retainedPeginQuote.DepositAddress: {entities.NewWei(30012)}},
	}, nil).Once()
	btc.On("GetRawTransaction", retainedPeginQuote.UserBtcTxHash).Return(btcRawTxMock, nil).Once()
	btc.On("GetPartialMerkleTree", retainedPeginQuote.UserBtcTxHash).Return(pmtMock, nil).Once()
	btc.On("GetTransactionBlockInfo", retainedPeginQuote.UserBtcTxHash).Return(btcBlockInfoMock, nil)
	mutex := new(test.MutexMock)
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()

	useCase := pegin.NewRegisterPeginUseCase(lbc, quoteRepository, eventBus, bridge, btc, mutex)
	err := useCase.Run(context.Background(), retainedPeginQuote)

	require.Error(t, err)
	lbc.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	bridge.AssertExpectations(t)
	btc.AssertExpectations(t)
	mutex.AssertExpectations(t)
}

func TestRegisterPeginUseCase_Run_NotEnoughConfirmations(t *testing.T) {
	retainedPeginQuote := quote.RetainedPeginQuote{
		QuoteHash:         "101b1c",
		DepositAddress:    "any address",
		Signature:         "0102031f1b",
		RequiredLiquidity: entities.NewWei(1500),
		State:             quote.PeginStateCallForUserSucceeded,
		UserBtcTxHash:     "btc tx hash",
		CallForUserTxHash: "cfu tx hash",
	}

	setups := registerPeginNotEnoughConfirmationsSetups(retainedPeginQuote)

	for _, testCase := range setups {
		t.Run(testCase.description, func(t *testing.T) {
			lbc := new(test.LbcMock)
			quoteRepository := new(test.PeginQuoteRepositoryMock)
			eventBus := new(test.EventBusMock)
			btc := new(test.BtcRpcMock)
			mutex := new(test.MutexMock)
			mutex.On("Lock").Return()
			mutex.On("Unlock").Return()
			bridge := new(test.BridgeMock)
			bridge.On("GetRequiredTxConfirmations").Return(uint64(30))

			testCase.setup(lbc, quoteRepository, btc)
			useCase := pegin.NewRegisterPeginUseCase(lbc, quoteRepository, eventBus, bridge, btc, mutex)
			err := useCase.Run(context.Background(), retainedPeginQuote)

			require.ErrorIs(t, err, testCase.err)
			lbc.AssertExpectations(t)
			quoteRepository.AssertExpectations(t)
			btc.AssertExpectations(t)
			eventBus.AssertNotCalled(t, "Publish")
		})
	}
}

func registerPeginNotEnoughConfirmationsSetups(retainedPeginQuote quote.RetainedPeginQuote) []struct {
	description string
	setup       func(lbc *test.LbcMock, quoteRepository *test.PeginQuoteRepositoryMock, btc *test.BtcRpcMock)
	err         error
} {
	return []struct {
		description string
		setup       func(lbc *test.LbcMock, quoteRepository *test.PeginQuoteRepositoryMock, btc *test.BtcRpcMock)
		err         error
	}{
		{
			description: "Should fail when tx has less confirmations than required from bridge",
			setup: func(lbc *test.LbcMock, quoteRepository *test.PeginQuoteRepositoryMock, btc *test.BtcRpcMock) {
				quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedPeginQuote.QuoteHash).
					Return(&testPeginQuote, nil).Once()
				btc.On("GetTransactionInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
					Hash:          retainedPeginQuote.UserBtcTxHash,
					Confirmations: 10,
					Outputs:       map[string][]*entities.Wei{retainedPeginQuote.DepositAddress: {entities.NewWei(30012)}},
				}, nil).Once()
			},
			err: usecases.NoEnoughConfirmationsError,
		},
		{
			description: "Should fail when confirmations weren't processed from RSK bridge yet",
			setup: func(lbc *test.LbcMock, quoteRepository *test.PeginQuoteRepositoryMock, btc *test.BtcRpcMock) {
				quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedPeginQuote.QuoteHash).
					Return(&testPeginQuote, nil).Once()
				btc.On("GetTransactionInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
					Hash:          retainedPeginQuote.UserBtcTxHash,
					Confirmations: 100,
					Outputs:       map[string][]*entities.Wei{retainedPeginQuote.DepositAddress: {entities.NewWei(30012)}},
				}, nil).Once()
				btc.On("GetRawTransaction", retainedPeginQuote.UserBtcTxHash).Return(btcRawTxMock, nil).Once()
				btc.On("GetPartialMerkleTree", retainedPeginQuote.UserBtcTxHash).Return(pmtMock, nil).Once()
				btc.On("GetTransactionBlockInfo", retainedPeginQuote.UserBtcTxHash).Return(btcBlockInfoMock, nil)
				lbc.On("RegisterPegin", blockchain.RegisterPeginParams{
					QuoteSignature:        []byte{1, 2, 3, 31, 27},
					BitcoinRawTransaction: btcRawTxMock,
					PartialMerkleTree:     pmtMock,
					BlockHeight:           btcBlockInfoMock.Height,
					Quote:                 testPeginQuote,
				}).Return("register tx hash", fmt.Errorf("some wrapper: %w", blockchain.WaitingForBridgeError)).Once()
			},
			err: blockchain.WaitingForBridgeError,
		},
	}
}

func TestRegisterPeginUseCase_Run_UpdateError(t *testing.T) {
	registerPeginTx := "register tx hash"
	retainedPeginQuote := quote.RetainedPeginQuote{
		QuoteHash:         "101b1c",
		DepositAddress:    "any address",
		Signature:         "0102031f1b",
		RequiredLiquidity: entities.NewWei(1500),
		State:             quote.PeginStateCallForUserSucceeded,
		UserBtcTxHash:     "btc tx hash",
		CallForUserTxHash: "cfu tx hash",
	}

	setups := registerPeginUpdateErrorSetups(t, registerPeginTx, retainedPeginQuote)

	bridge := new(test.BridgeMock)
	bridge.On("GetRequiredTxConfirmations").Return(uint64(10))

	btc := new(test.BtcRpcMock)
	btc.On("GetTransactionInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
		Hash:          retainedPeginQuote.UserBtcTxHash,
		Confirmations: 11,
		Outputs:       map[string][]*entities.Wei{retainedPeginQuote.DepositAddress: {entities.NewWei(30012)}},
	}, nil)
	btc.On("GetRawTransaction", retainedPeginQuote.UserBtcTxHash).Return(btcRawTxMock, nil)
	btc.On("GetPartialMerkleTree", retainedPeginQuote.UserBtcTxHash).Return(pmtMock, nil)
	btc.On("GetTransactionBlockInfo", retainedPeginQuote.UserBtcTxHash).Return(btcBlockInfoMock, nil)

	mutex := new(test.MutexMock)
	mutex.On("Lock").Return()
	mutex.On("Unlock").Return()

	lbc := new(test.LbcMock)
	lbc.On("RegisterPegin", blockchain.RegisterPeginParams{
		QuoteSignature:        []byte{1, 2, 3, 31, 27},
		BitcoinRawTransaction: btcRawTxMock,
		PartialMerkleTree:     pmtMock,
		BlockHeight:           btcBlockInfoMock.Height,
		Quote:                 testPeginQuote,
	}).Return(registerPeginTx, nil)

	for _, setup := range setups {
		quoteRepository := new(test.PeginQuoteRepositoryMock)
		eventBus := new(test.EventBusMock)

		setup(quoteRepository, eventBus)
		useCase := pegin.NewRegisterPeginUseCase(lbc, quoteRepository, eventBus, bridge, btc, mutex)
		err := useCase.Run(context.Background(), retainedPeginQuote)

		require.Error(t, err)
		quoteRepository.AssertExpectations(t)
		eventBus.AssertExpectations(t)
	}
}

func registerPeginUpdateErrorSetups(t *testing.T, registerPeginTx string, retainedPeginQuote quote.RetainedPeginQuote) []func(quoteRepository *test.PeginQuoteRepositoryMock, eventBus *test.EventBusMock) {
	return []func(quoteRepository *test.PeginQuoteRepositoryMock, eventBus *test.EventBusMock){
		func(quoteRepository *test.PeginQuoteRepositoryMock, eventBus *test.EventBusMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedPeginQuote.QuoteHash).
				Return(nil, nil).Once()
			quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.MatchedBy(func(q quote.RetainedPeginQuote) bool {
				expected := retainedPeginQuote
				expected.State = quote.PeginStateRegisterPegInFailed
				return assert.Equal(t, expected, q)
			})).Return(assert.AnError).Once()
			eventBus.On("Publish", mock.MatchedBy(func(event quote.RegisterPeginCompletedEvent) bool {
				expected := retainedPeginQuote
				expected.State = quote.PeginStateRegisterPegInFailed
				require.Error(t, event.Error)
				return assert.Equal(t, expected, event.RetainedQuote) &&
					assert.Equal(t, quote.RegisterPeginCompletedEventId, event.Event.Id())
			})).Return().Once()
		},
		func(quoteRepository *test.PeginQuoteRepositoryMock, eventBus *test.EventBusMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedPeginQuote.QuoteHash).
				Return(&testPeginQuote, nil).Once()
			quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.MatchedBy(func(q quote.RetainedPeginQuote) bool {
				expected := retainedPeginQuote
				expected.State = quote.PeginStateRegisterPegInSucceeded
				expected.RegisterPeginTxHash = registerPeginTx
				return assert.Equal(t, expected, q)
			})).Return(assert.AnError).Once()
			eventBus.On("Publish", mock.MatchedBy(func(event quote.RegisterPeginCompletedEvent) bool {
				expected := retainedPeginQuote
				expected.State = quote.PeginStateRegisterPegInSucceeded
				expected.RegisterPeginTxHash = registerPeginTx
				require.NoError(t, event.Error)
				return assert.Equal(t, expected, event.RetainedQuote) &&
					assert.Equal(t, quote.RegisterPeginCompletedEventId, event.Event.Id())
			})).Return().Once()
		},
	}
}
