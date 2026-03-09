package pegin_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCallForUserUseCase_Run_Paused(t *testing.T) {
	retainedPeginQuote := quote.RetainedPeginQuote{
		QuoteHash:         "101b1c",
		DepositAddress:    test.AnyAddress,
		Signature:         "signature",
		RequiredLiquidity: entities.NewWei(1500),
		State:             quote.PeginStateWaitingForDepositConfirmations,
		UserBtcTxHash:     "0x121a1b",
	}
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	bridge := new(mocks.BridgeMock)
	btc := new(mocks.BtcRpcMock)
	lp := new(mocks.ProviderMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	rsk := new(mocks.RootstockRpcServerMock)
	peginContract := new(mocks.PeginContractMock)
	peginContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: true, Since: 5, Reason: "test"}, nil)
	peginContract.EXPECT().GetAddress().Return("test-contract")
	contracts := blockchain.RskContracts{Bridge: bridge, PegIn: peginContract}
	useCase := pegin.NewCallForUserUseCase(contracts, quoteRepository, blockchain.Rpc{Rsk: rsk, Btc: btc}, lp, eventBus, mutex)
	err := useCase.Run(context.Background(), retainedPeginQuote)
	require.ErrorIs(t, err, blockchain.ContractPausedError)
}

// nolint:funlen
func TestCallForUserUseCase_Run(t *testing.T) {
	callForUserTxHash := "0x1a1b1c"
	lpRskAddress := testPeginQuote.LpRskAddress
	retainedPeginQuote := quote.RetainedPeginQuote{
		QuoteHash:         "101b1c",
		DepositAddress:    test.AnyAddress,
		Signature:         "signature",
		RequiredLiquidity: entities.NewWei(1500),
		State:             quote.PeginStateWaitingForDepositConfirmations,
		UserBtcTxHash:     "0x121a1b",
	}
	creationData := quote.PeginCreationData{GasPrice: entities.NewWei(5), FeePercentage: utils.NewBigFloat64(1.24), FixedFee: entities.NewWei(100)}
	expectedRetainedQuote := retainedPeginQuote
	expectedRetainedQuote.State = quote.PeginStateCallForUserSucceeded
	expectedRetainedQuote.CallForUserTxHash = callForUserTxHash
	expectedRetainedQuote.CallForUserGasUsed = uint64(21000)
	expectedRetainedQuote.CallForUserGasPrice = entities.NewWei(1000000000)
	bridge := new(mocks.BridgeMock)
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()

	lp := new(mocks.ProviderMock)
	lp.On("RskAddress").Return(lpRskAddress).Once()
	peginContract := new(mocks.PeginContractMock)
	peginContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil).Once()
	peginContract.On("GetBalance", testPeginQuote.LpRskAddress).Return(entities.NewWei(50000), nil).Once()
	txConfig := blockchain.NewTransactionConfig(entities.NewWei(0), uint64(testPeginQuote.GasLimit+pegin.CallForUserExtraGas), nil)
	callForUserReceipt := blockchain.TransactionReceipt{
		TransactionHash:   callForUserTxHash,
		BlockHash:         "0xblock123",
		BlockNumber:       uint64(1000),
		From:              testPeginQuote.LpRskAddress,
		To:                testPeginQuote.ContractAddress,
		CumulativeGasUsed: big.NewInt(21000),
		GasUsed:           big.NewInt(21000),
		Value:             entities.NewWei(0),
		GasPrice:          entities.NewWei(1000000000),
	}
	peginContract.On("CallForUser", txConfig, testPeginQuote).Return(callForUserReceipt, nil).Once()
	btc := new(mocks.BtcRpcMock)
	btc.On("GetTransactionInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
		Hash: retainedPeginQuote.UserBtcTxHash, Confirmations: 10,
		Outputs: map[string][]*entities.Wei{retainedPeginQuote.DepositAddress: {entities.NewWei(30012)}},
	}, nil).Once()
	btc.On("GetTransactionBlockInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinBlockInformation{Hash: [32]byte{1, 2, 3}, Height: big.NewInt(5), Time: time.Now()}, nil).Once()

	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.CallForUserCompletedEvent) bool {
		require.NoError(t, event.Error)
		return assert.Equal(t, testPeginQuote, event.PeginQuote) && assert.Equal(t, expectedRetainedQuote, event.RetainedQuote) &&
			assert.Equal(t, quote.CallForUserCompletedEventId, event.Event.Id())
	})).Return().Once()
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()

	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, retainedPeginQuote.QuoteHash).Return(&testPeginQuote, nil).Once()
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, mock.MatchedBy(func(q quote.RetainedPeginQuote) bool {
		return assert.Equal(t, expectedRetainedQuote, q)
	})).Return(nil).Once()
	quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, retainedPeginQuote.QuoteHash).Return(creationData).Once()
	rsk := new(mocks.RootstockRpcServerMock)

	contracts := blockchain.RskContracts{PegIn: peginContract, Bridge: bridge}
	rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
	useCase := pegin.NewCallForUserUseCase(contracts, quoteRepository, rpc, lp, eventBus, mutex)
	err := useCase.Run(context.Background(), retainedPeginQuote)
	require.NoError(t, err)
	peginContract.AssertExpectations(t)
	btc.AssertExpectations(t)
	lp.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	mutex.AssertExpectations(t)
	bridge.AssertExpectations(t)
}

//nolint:funlen
func TestCallForUserUseCase_Run_AddExtraAmountDuringCall(t *testing.T) {
	callForUserTxHash := "0x1a1b1c"
	lpRskAddress := testPeginQuote.LpRskAddress
	retainedPeginQuote := quote.RetainedPeginQuote{
		QuoteHash:         "101b1c",
		DepositAddress:    test.AnyAddress,
		Signature:         "signature",
		RequiredLiquidity: entities.NewWei(1500),
		State:             quote.PeginStateWaitingForDepositConfirmations,
		UserBtcTxHash:     "0x121a1b",
	}
	creationData := quote.PeginCreationData{GasPrice: entities.NewWei(5), FeePercentage: utils.NewBigFloat64(1.24), FixedFee: entities.NewWei(100)}
	expectedRetainedQuote := retainedPeginQuote
	expectedRetainedQuote.State = quote.PeginStateCallForUserSucceeded
	expectedRetainedQuote.CallForUserTxHash = callForUserTxHash
	expectedRetainedQuote.CallForUserGasUsed = uint64(21000)
	expectedRetainedQuote.CallForUserGasPrice = entities.NewWei(1000000000)

	bridge := new(mocks.BridgeMock)
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
	lp := new(mocks.ProviderMock)
	lp.On("RskAddress").Return(lpRskAddress).Twice()
	peginContract := new(mocks.PeginContractMock)
	peginContract.On("GetBalance", testPeginQuote.LpRskAddress).Return(entities.NewWei(600), nil).Once()
	peginContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil).Once()
	txConfig := blockchain.NewTransactionConfig(entities.NewWei(29400), uint64(testPeginQuote.GasLimit+pegin.CallForUserExtraGas), nil)
	callForUserReceipt := blockchain.TransactionReceipt{
		TransactionHash:   callForUserTxHash,
		BlockHash:         "0xblock123",
		BlockNumber:       uint64(1000),
		From:              testPeginQuote.LpRskAddress,
		To:                testPeginQuote.ContractAddress,
		CumulativeGasUsed: big.NewInt(21000),
		GasUsed:           big.NewInt(21000),
		Value:             entities.NewWei(29400),
		GasPrice:          entities.NewWei(1000000000),
	}
	peginContract.On("CallForUser", txConfig, testPeginQuote).Return(callForUserReceipt, nil).Once()
	btc := new(mocks.BtcRpcMock)
	btc.On("GetTransactionInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
		Hash:          retainedPeginQuote.UserBtcTxHash,
		Confirmations: 10,
		Outputs:       map[string][]*entities.Wei{retainedPeginQuote.DepositAddress: {entities.NewWei(30012)}},
	}, nil).Once()
	btc.On("GetTransactionBlockInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinBlockInformation{Hash: [32]byte{1, 2, 3}, Height: big.NewInt(5), Time: time.Now()}, nil).Once()
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.CallForUserCompletedEvent) bool {
		require.NoError(t, event.Error)
		return assert.Equal(t, testPeginQuote, event.PeginQuote) && assert.Equal(t, expectedRetainedQuote, event.RetainedQuote) && assert.Equal(t, quote.CallForUserCompletedEventId, event.Event.Id())
	})).Return().Once()

	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, retainedPeginQuote.QuoteHash).Return(&testPeginQuote, nil).Once()
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, mock.MatchedBy(func(q quote.RetainedPeginQuote) bool {
		return assert.Equal(t, expectedRetainedQuote, q)
	})).Return(nil).Once()
	quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, retainedPeginQuote.QuoteHash).Return(creationData).Once()

	rsk := new(mocks.RootstockRpcServerMock)
	rsk.On("GetBalance", test.AnyCtx, lpRskAddress).Return(entities.NewWei(80000), nil).Once()
	contracts := blockchain.RskContracts{PegIn: peginContract, Bridge: bridge}
	useCase := pegin.NewCallForUserUseCase(contracts, quoteRepository, blockchain.Rpc{Rsk: rsk, Btc: btc}, lp, eventBus, mutex)
	err := useCase.Run(context.Background(), retainedPeginQuote)

	require.NoError(t, err)
	peginContract.AssertExpectations(t)
	btc.AssertExpectations(t)
	lp.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	mutex.AssertExpectations(t)
	bridge.AssertExpectations(t)
}

func TestCallForUserUseCase_Run_DontPublishRecoverableErrors(t *testing.T) {
	retainedPeginQuote := quote.RetainedPeginQuote{
		QuoteHash:         "101b1c",
		DepositAddress:    test.AnyAddress,
		Signature:         "signature",
		RequiredLiquidity: entities.NewWei(1500),
		State:             quote.PeginStateWaitingForDepositConfirmations,
	}

	setups := callForUserRecoverableErrorSetups()

	for _, setup := range setups {
		lp := new(mocks.ProviderMock)
		lp.On("RskAddress").Return("lp rsk address")
		peginContract := new(mocks.PeginContractMock)
		peginContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
		btc := new(mocks.BtcRpcMock)
		eventBus := new(mocks.EventBusMock)
		rsk := new(mocks.RootstockRpcServerMock)
		mutex := new(mocks.MutexMock)
		mutex.On("Lock").Return().Once()
		mutex.On("Unlock").Return().Once()
		quoteRepository := new(mocks.PeginQuoteRepositoryMock)
		quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, retainedPeginQuote.QuoteHash).Return(quote.PeginCreationDataZeroValue()).Once()
		bridge := new(mocks.BridgeMock)
		caseRetainedQuote := retainedPeginQuote
		setup(&caseRetainedQuote, rsk, peginContract, btc, quoteRepository, bridge)

		contracts := blockchain.RskContracts{PegIn: peginContract, Bridge: bridge}
		rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
		useCase := pegin.NewCallForUserUseCase(contracts, quoteRepository, rpc, lp, eventBus, mutex)
		err := useCase.Run(context.Background(), caseRetainedQuote)
		require.Error(t, err)

	}
}

// nolint:funlen
func callForUserRecoverableErrorSetups() []func(caseRetainedQuote *quote.RetainedPeginQuote, rsk *mocks.RootstockRpcServerMock, peginContract *mocks.PeginContractMock, btc *mocks.BtcRpcMock, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock) {
	now := uint32(time.Now().Unix())
	peginQuote := quote.PeginQuote{
		FedBtcAddress:      "fed address",
		LbcAddress:         "lbc address",
		LpRskAddress:       "lp rsk address",
		BtcRefundAddress:   "btc refund address",
		RskRefundAddress:   "rsk refund address",
		LpBtcAddress:       "lp btc address",
		CallFee:            entities.NewWei(100),
		PenaltyFee:         entities.NewWei(100),
		ContractAddress:    "contract address",
		Data:               "0x1a1b",
		GasLimit:           500,
		Nonce:              123456,
		Value:              entities.NewWei(1000),
		AgreementTimestamp: now,
		TimeForDeposit:     600,
		LpCallTime:         600,
		Confirmations:      10,
		CallOnRegister:     false,
		GasFee:             entities.NewWei(500),
		ChainId:            31,
	}
	return []func(caseRetainedQuote *quote.RetainedPeginQuote, rsk *mocks.RootstockRpcServerMock, peginContract *mocks.PeginContractMock, btc *mocks.BtcRpcMock, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock){
		func(caseRetainedQuote *quote.RetainedPeginQuote, rsk *mocks.RootstockRpcServerMock, peginContract *mocks.PeginContractMock, btc *mocks.BtcRpcMock, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock) {
			caseRetainedQuote.State = quote.PeginStateCallForUserSucceeded
		},
		func(caseRetainedQuote *quote.RetainedPeginQuote, rsk *mocks.RootstockRpcServerMock, peginContract *mocks.PeginContractMock, btc *mocks.BtcRpcMock, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(caseRetainedQuote *quote.RetainedPeginQuote, rsk *mocks.RootstockRpcServerMock, peginContract *mocks.PeginContractMock, btc *mocks.BtcRpcMock, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).
				Return(&peginQuote, nil).Once()
			btc.On("GetTransactionInfo", mock.Anything).Return(blockchain.BitcoinTransactionInformation{}, assert.AnError).Once()
		},
		func(caseRetainedQuote *quote.RetainedPeginQuote, rsk *mocks.RootstockRpcServerMock, peginContract *mocks.PeginContractMock, btc *mocks.BtcRpcMock, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).
				Return(&peginQuote, nil).Once()
			btc.On("GetTransactionInfo", mock.Anything).Return(blockchain.BitcoinTransactionInformation{
				Hash:          "0x1d1e",
				Confirmations: 10,
				Outputs:       map[string][]*entities.Wei{test.AnyAddress: {entities.NewWei(1700)}},
			}, nil).Once()
			btc.On("GetTransactionBlockInfo", mock.Anything).Return(blockchain.BitcoinBlockInformation{}, assert.AnError).Once()
		},
		func(caseRetainedQuote *quote.RetainedPeginQuote, rsk *mocks.RootstockRpcServerMock, peginContract *mocks.PeginContractMock, btc *mocks.BtcRpcMock, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).
				Return(&peginQuote, nil).Once()
			btc.On("GetTransactionInfo", mock.Anything).Return(blockchain.BitcoinTransactionInformation{
				Hash:          "0x1d1e",
				Confirmations: 10,
				Outputs:       map[string][]*entities.Wei{test.AnyAddress: {entities.NewWei(1700)}},
			}, nil).Once()
			btc.On("GetTransactionBlockInfo", mock.Anything).Return(blockchain.BitcoinBlockInformation{
				Hash:   [32]byte{1, 2, 3},
				Height: big.NewInt(5),
				Time:   time.Now(),
			}, nil).Once()
			bridge.On("GetMinimumLockTxValue").Return(nil, assert.AnError).Once()
		},
		func(caseRetainedQuote *quote.RetainedPeginQuote, rsk *mocks.RootstockRpcServerMock, peginContract *mocks.PeginContractMock, btc *mocks.BtcRpcMock, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).
				Return(&peginQuote, nil).Once()
			btc.On("GetTransactionInfo", mock.Anything).Return(blockchain.BitcoinTransactionInformation{
				Hash:          "0x1d1e",
				Confirmations: 10,
				Outputs:       map[string][]*entities.Wei{test.AnyAddress: {entities.NewWei(1700)}},
			}, nil).Once()
			btc.On("GetTransactionBlockInfo", mock.Anything).Return(blockchain.BitcoinBlockInformation{
				Hash:   [32]byte{1, 2, 3},
				Height: big.NewInt(5),
				Time:   time.Now(),
			}, nil).Once()
			bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
			peginContract.On("GetBalance", mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(caseRetainedQuote *quote.RetainedPeginQuote, rsk *mocks.RootstockRpcServerMock, peginContract *mocks.PeginContractMock, btc *mocks.BtcRpcMock, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).
				Return(&peginQuote, nil).Once()
			btc.On("GetTransactionInfo", mock.Anything).Return(blockchain.BitcoinTransactionInformation{
				Hash:          "0x1d1e",
				Confirmations: 10,
				Outputs:       map[string][]*entities.Wei{test.AnyAddress: {entities.NewWei(1700)}},
			}, nil).Once()
			btc.On("GetTransactionBlockInfo", mock.Anything).Return(blockchain.BitcoinBlockInformation{
				Hash:   [32]byte{1, 2, 3},
				Height: big.NewInt(5),
				Time:   time.Now(),
			}, nil).Once()
			bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
			peginContract.On("GetBalance", mock.Anything).Return(entities.NewWei(500), nil).Once()
			rsk.On("GetBalance", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		},
	}
}

func TestCallForUserUseCase_Run_NoConfirmations(t *testing.T) {
	retainedPeginQuote := quote.RetainedPeginQuote{
		QuoteHash:         "101b1c",
		DepositAddress:    test.AnyAddress,
		Signature:         "signature",
		RequiredLiquidity: entities.NewWei(1500),
		State:             quote.PeginStateWaitingForDepositConfirmations,
		UserBtcTxHash:     "0x121a1b",
	}

	lp := new(mocks.ProviderMock)
	peginContract := new(mocks.PeginContractMock)
	peginContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil).Once()

	btc := new(mocks.BtcRpcMock)
	btc.On("GetTransactionInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
		Hash:          retainedPeginQuote.UserBtcTxHash,
		Confirmations: 5,
		Outputs:       map[string][]*entities.Wei{retainedPeginQuote.DepositAddress: {entities.NewWei(2000)}},
	}, nil).Once()

	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, retainedPeginQuote.QuoteHash).
		Return(&testPeginQuote, nil).Once()
	quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, retainedPeginQuote.QuoteHash).Return(quote.PeginCreationDataZeroValue()).Once()

	rsk := new(mocks.RootstockRpcServerMock)

	contracts := blockchain.RskContracts{PegIn: peginContract}
	rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
	useCase := pegin.NewCallForUserUseCase(contracts, quoteRepository, rpc, lp, eventBus, mutex)
	err := useCase.Run(context.Background(), retainedPeginQuote)

	require.ErrorIs(t, err, usecases.NoEnoughConfirmationsError)
	peginContract.AssertExpectations(t)
	btc.AssertExpectations(t)
	lp.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertNotCalled(t, "Publish")
	mutex.AssertExpectations(t)
	peginContract.AssertNotCalled(t, "CallForUser")
	mutex.AssertNotCalled(t, "Lock")
	mutex.AssertNotCalled(t, "Unlock")
	peginContract.AssertNotCalled(t, "GetBalance")
	lp.AssertNotCalled(t, "RskAddress")
}

func TestCallForUserUseCase_Run_ExpiredQuote(t *testing.T) {
	peginContract := new(mocks.PeginContractMock)
	peginContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil).Once()
	btc := new(mocks.BtcRpcMock)
	lp := new(mocks.ProviderMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)

	retainedPeginQuote := quote.RetainedPeginQuote{
		QuoteHash:         "101b1c",
		DepositAddress:    test.AnyAddress,
		Signature:         "signature",
		RequiredLiquidity: entities.NewWei(1500),
		State:             quote.PeginStateWaitingForDepositConfirmations,
		UserBtcTxHash:     test.AnyString,
	}

	btc.On("GetTransactionInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
		Hash:          retainedPeginQuote.UserBtcTxHash,
		Confirmations: 10,
		Outputs:       map[string][]*entities.Wei{retainedPeginQuote.DepositAddress: {entities.NewWei(30012)}},
	}, nil).Once()
	btc.On("GetTransactionBlockInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinBlockInformation{
		Hash:   [32]byte{1, 2, 3},
		Height: big.NewInt(5),
		Time:   time.Now().Add(time.Hour * 10),
	}, nil).Once()

	expiredQuote := testPeginQuote
	expiredQuote.AgreementTimestamp -= 1000
	quoteRepository.On("GetQuote", test.AnyCtx, retainedPeginQuote.QuoteHash).Return(&expiredQuote, nil).Once()

	updatedQuote := retainedPeginQuote
	updatedQuote.State = quote.PeginStateCallForUserFailed

	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, updatedQuote).Return(nil).Once()
	quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, retainedPeginQuote.QuoteHash).Return(quote.PeginCreationDataZeroValue()).Once()
	eventBus.On("Publish", mock.MatchedBy(func(event quote.CallForUserCompletedEvent) bool {
		require.ErrorIs(t, event.Error, usecases.ExpiredQuoteError)
		return assert.Equal(t, expiredQuote, event.PeginQuote) &&
			assert.Equal(t, updatedQuote, event.RetainedQuote) &&
			assert.Equal(t, quote.CallForUserCompletedEventId, event.Event.Id())
	})).Return().Once()

	rsk := new(mocks.RootstockRpcServerMock)

	contracts := blockchain.RskContracts{PegIn: peginContract}
	rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
	useCase := pegin.NewCallForUserUseCase(contracts, quoteRepository, rpc, lp, eventBus, mutex)
	err := useCase.Run(context.Background(), retainedPeginQuote)
	require.ErrorIs(t, err, usecases.ExpiredQuoteError)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	peginContract.AssertNotCalled(t, "GetBalance")
	peginContract.AssertNotCalled(t, "CallForUser")
	mutex.AssertNotCalled(t, "Lock")
	mutex.AssertNotCalled(t, "Unlock")
	btc.AssertNotCalled(t, "GetTransactionInfo")
}

func TestCallForUserUseCase_Run_QuoteNotFound(t *testing.T) {
	peginContract := new(mocks.PeginContractMock)
	peginContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil).Once()
	btc := new(mocks.BtcRpcMock)
	lp := new(mocks.ProviderMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)

	retainedPeginQuote := quote.RetainedPeginQuote{
		QuoteHash:         "101b1c",
		DepositAddress:    test.AnyAddress,
		Signature:         "signature",
		RequiredLiquidity: entities.NewWei(1500),
		State:             quote.PeginStateWaitingForDepositConfirmations,
	}

	updatedQuote := retainedPeginQuote
	updatedQuote.State = quote.PeginStateCallForUserFailed

	quoteRepository.On("GetQuote", test.AnyCtx, retainedPeginQuote.QuoteHash).Return(nil, nil).Once()

	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, updatedQuote).Return(nil).Once()
	eventBus.On("Publish", mock.MatchedBy(func(event quote.CallForUserCompletedEvent) bool {
		require.ErrorIs(t, event.Error, usecases.QuoteNotFoundError)
		return assert.Empty(t, event.PeginQuote) &&
			assert.Equal(t, updatedQuote, event.RetainedQuote) &&
			assert.Equal(t, quote.CallForUserCompletedEventId, event.Event.Id())
	})).Return().Once()

	rsk := new(mocks.RootstockRpcServerMock)

	contracts := blockchain.RskContracts{PegIn: peginContract}
	rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
	useCase := pegin.NewCallForUserUseCase(contracts, quoteRepository, rpc, lp, eventBus, mutex)
	err := useCase.Run(context.Background(), retainedPeginQuote)
	require.ErrorIs(t, err, usecases.QuoteNotFoundError)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	peginContract.AssertNotCalled(t, "GetBalance")
	peginContract.AssertNotCalled(t, "CallForUser")
	mutex.AssertNotCalled(t, "Lock")
	mutex.AssertNotCalled(t, "Unlock")
	btc.AssertNotCalled(t, "GetTransactionInfo")
}

func TestCallForUserUseCase_Run_InsufficientAmount(t *testing.T) {
	retainedPeginQuote := quote.RetainedPeginQuote{
		QuoteHash:         "101b1c",
		DepositAddress:    test.AnyAddress,
		Signature:         "signature",
		RequiredLiquidity: entities.NewWei(1500),
		State:             quote.PeginStateWaitingForDepositConfirmations,
		UserBtcTxHash:     "0x121a1b",
	}

	lp := new(mocks.ProviderMock)
	peginContract := new(mocks.PeginContractMock)
	peginContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil).Once()
	btc := new(mocks.BtcRpcMock)
	btc.On("GetTransactionInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
		Hash:          retainedPeginQuote.UserBtcTxHash,
		Confirmations: 10,
		Outputs:       map[string][]*entities.Wei{retainedPeginQuote.DepositAddress: {entities.NewWei(900)}},
	}, nil).Once()
	btc.On("GetTransactionBlockInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinBlockInformation{Hash: [32]byte{1, 2, 3}, Height: big.NewInt(5), Time: time.Now()}, nil).Once()

	updatedQuote := retainedPeginQuote
	updatedQuote.State = quote.PeginStateCallForUserFailed

	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.CallForUserCompletedEvent) bool {
		require.ErrorIs(t, event.Error, usecases.InsufficientAmountError)
		return assert.Equal(t, testPeginQuote, event.PeginQuote) &&
			assert.Equal(t, updatedQuote, event.RetainedQuote) &&
			assert.Equal(t, quote.CallForUserCompletedEventId, event.Event.Id())
	})).Return().Once()

	mutex := new(mocks.MutexMock)
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, retainedPeginQuote.QuoteHash).
		Return(&testPeginQuote, nil).Once()
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, updatedQuote).
		Return(nil).Once()
	quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, retainedPeginQuote.QuoteHash).Return(quote.PeginCreationDataZeroValue()).Once()

	rsk := new(mocks.RootstockRpcServerMock)
	contracts := blockchain.RskContracts{PegIn: peginContract}
	rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
	useCase := pegin.NewCallForUserUseCase(contracts, quoteRepository, rpc, lp, eventBus, mutex)
	err := useCase.Run(context.Background(), retainedPeginQuote)

	require.ErrorIs(t, err, usecases.InsufficientAmountError)
	peginContract.AssertExpectations(t)
	btc.AssertExpectations(t)
	lp.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	mutex.AssertExpectations(t)
	peginContract.AssertNotCalled(t, "CallForUser")
	mutex.AssertNotCalled(t, "Lock")
	mutex.AssertNotCalled(t, "Unlock")
	peginContract.AssertNotCalled(t, "GetBalance")
	lp.AssertNotCalled(t, "RskAddress")
}

func TestCallForUserUseCase_Run_NoLiquidity(t *testing.T) {
	lpRskAddress := testPeginQuote.LpRskAddress

	retainedPeginQuote := quote.RetainedPeginQuote{
		QuoteHash:         "101b1c",
		DepositAddress:    test.AnyAddress,
		Signature:         "signature",
		RequiredLiquidity: entities.NewWei(1500),
		State:             quote.PeginStateWaitingForDepositConfirmations,
		UserBtcTxHash:     "0x121a1b",
	}

	lp := new(mocks.ProviderMock)
	lp.On("RskAddress").Return(lpRskAddress).Twice()

	peginContract := new(mocks.PeginContractMock)
	peginContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil).Once()
	peginContract.On("GetBalance", testPeginQuote.LpRskAddress).Return(entities.NewWei(500), nil).Once()

	btc := new(mocks.BtcRpcMock)
	btc.On("GetTransactionInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
		Hash:          retainedPeginQuote.UserBtcTxHash,
		Confirmations: 10,
		Outputs:       map[string][]*entities.Wei{retainedPeginQuote.DepositAddress: {entities.NewWei(30012)}},
	}, nil).Once()
	btc.On("GetTransactionBlockInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinBlockInformation{
		Hash:   [32]byte{1, 2, 3},
		Height: big.NewInt(5),
		Time:   time.Now(),
	}, nil).Once()

	eventBus := new(mocks.EventBusMock)

	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()

	bridge := new(mocks.BridgeMock)
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()

	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, retainedPeginQuote.QuoteHash).Return(&testPeginQuote, nil).Once()
	quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, retainedPeginQuote.QuoteHash).Return(quote.PeginCreationDataZeroValue()).Once()

	rsk := new(mocks.RootstockRpcServerMock)
	rsk.On("GetBalance", test.AnyCtx, lpRskAddress).Return(entities.NewWei(20000), nil).Once()

	contracts := blockchain.RskContracts{PegIn: peginContract, Bridge: bridge}
	rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
	useCase := pegin.NewCallForUserUseCase(contracts, quoteRepository, rpc, lp, eventBus, mutex)
	err := useCase.Run(context.Background(), retainedPeginQuote)

	require.ErrorIs(t, err, usecases.NoLiquidityError)
	peginContract.AssertExpectations(t)
	btc.AssertExpectations(t)
	lp.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertNotCalled(t, "Publish")
	mutex.AssertExpectations(t)
	bridge.AssertExpectations(t)
}

//nolint:funlen
func TestCallForUserUseCase_Run_CallForUserFail(t *testing.T) {
	callForUserTxHash := "0x1a1b1c"
	lpRskAddress := testPeginQuote.LpRskAddress
	retainedPeginQuote := quote.RetainedPeginQuote{
		QuoteHash:         "101b1c",
		DepositAddress:    test.AnyAddress,
		Signature:         "signature",
		RequiredLiquidity: entities.NewWei(1500),
		State:             quote.PeginStateWaitingForDepositConfirmations,
		UserBtcTxHash:     "0x121a1b",
	}
	expectedRetainedQuote := retainedPeginQuote
	expectedRetainedQuote.State = quote.PeginStateCallForUserFailed
	expectedRetainedQuote.CallForUserTxHash = callForUserTxHash
	expectedRetainedQuote.CallForUserGasUsed = uint64(21000)
	expectedRetainedQuote.CallForUserGasPrice = entities.NewWei(1000000000)

	lp := new(mocks.ProviderMock)
	lp.On("RskAddress").Return(lpRskAddress).Twice()
	bridge := new(mocks.BridgeMock)
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
	peginContract := new(mocks.PeginContractMock)
	peginContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil).Once()
	peginContract.On("GetBalance", testPeginQuote.LpRskAddress).Return(entities.NewWei(600), nil).Once()
	txConfig := blockchain.NewTransactionConfig(entities.NewWei(29400), uint64(testPeginQuote.GasLimit+pegin.CallForUserExtraGas), nil)
	callForUserReceipt := blockchain.TransactionReceipt{
		TransactionHash:   callForUserTxHash,
		BlockHash:         "0xblock123",
		BlockNumber:       uint64(1000),
		From:              testPeginQuote.LpRskAddress,
		To:                testPeginQuote.ContractAddress,
		CumulativeGasUsed: big.NewInt(21000),
		GasUsed:           big.NewInt(21000),
		Value:             entities.NewWei(1234),
		GasPrice:          entities.NewWei(1000000000),
	}
	peginContract.On("CallForUser", txConfig, testPeginQuote).Return(callForUserReceipt, assert.AnError).Once()
	btc := new(mocks.BtcRpcMock)
	btc.On("GetTransactionInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
		Hash:          retainedPeginQuote.UserBtcTxHash,
		Confirmations: 10,
		Outputs:       map[string][]*entities.Wei{retainedPeginQuote.DepositAddress: {entities.NewWei(30012)}},
	}, nil).Once()
	btc.On("GetTransactionBlockInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinBlockInformation{Hash: [32]byte{1, 2, 3}, Height: big.NewInt(5), Time: time.Now()}, nil).Once()
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.CallForUserCompletedEvent) bool {
		require.Error(t, event.Error)
		return assert.Equal(t, testPeginQuote, event.PeginQuote) && assert.Equal(t, expectedRetainedQuote, event.RetainedQuote) && assert.Equal(t, quote.CallForUserCompletedEventId, event.Event.Id())
	})).Return().Once()

	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, retainedPeginQuote.QuoteHash).Return(&testPeginQuote, nil).Once()
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, mock.MatchedBy(func(q quote.RetainedPeginQuote) bool {
		return assert.Equal(t, expectedRetainedQuote, q)
	})).Return(nil).Once()
	quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, retainedPeginQuote.QuoteHash).Return(quote.PeginCreationDataZeroValue()).Once()

	rsk := new(mocks.RootstockRpcServerMock)
	rsk.On("GetBalance", test.AnyCtx, lpRskAddress).Return(entities.NewWei(80000), nil).Once()
	contracts := blockchain.RskContracts{PegIn: peginContract, Bridge: bridge}
	rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
	useCase := pegin.NewCallForUserUseCase(contracts, quoteRepository, rpc, lp, eventBus, mutex)
	err := useCase.Run(context.Background(), retainedPeginQuote)

	require.Error(t, err)
	peginContract.AssertExpectations(t)
	btc.AssertExpectations(t)
	lp.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	mutex.AssertExpectations(t)
	bridge.AssertExpectations(t)
}

func TestCallForUserUseCase_Run_InvalidUTXOs(t *testing.T) {
	retainedPeginQuote := quote.RetainedPeginQuote{
		QuoteHash:         "abcdef",
		DepositAddress:    test.AnyAddress,
		Signature:         "signature",
		RequiredLiquidity: entities.NewWei(1500),
		State:             quote.PeginStateWaitingForDepositConfirmations,
		UserBtcTxHash:     "0x3c2b1a",
	}
	expectedRetainedQuote := retainedPeginQuote
	expectedRetainedQuote.State = quote.PeginStateCallForUserFailed

	lp := new(mocks.ProviderMock)
	peginContract := new(mocks.PeginContractMock)
	peginContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil).Once()

	bridge := new(mocks.BridgeMock)
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()

	btc := new(mocks.BtcRpcMock)
	btc.On("GetTransactionInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
		Hash:          retainedPeginQuote.UserBtcTxHash,
		Confirmations: 20,
		Outputs:       map[string][]*entities.Wei{retainedPeginQuote.DepositAddress: {entities.NewWei(30000), entities.NewWei(999)}},
	}, nil).Once()
	btc.On("GetTransactionBlockInfo", retainedPeginQuote.UserBtcTxHash).Return(blockchain.BitcoinBlockInformation{Hash: [32]byte{1, 2, 3}, Height: big.NewInt(5), Time: time.Now()}, nil).Once()

	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.CallForUserCompletedEvent) bool {
		require.Error(t, event.Error)
		return assert.Equal(t, testPeginQuote, event.PeginQuote) && assert.Equal(t, expectedRetainedQuote, event.RetainedQuote) &&
			assert.Equal(t, quote.CallForUserCompletedEventId, event.Event.Id())
	})).Return().Once()

	mutex := new(mocks.MutexMock)

	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, retainedPeginQuote.QuoteHash).Return(&testPeginQuote, nil).Once()
	quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, mock.MatchedBy(func(q quote.RetainedPeginQuote) bool {
		return assert.Equal(t, expectedRetainedQuote, q)
	})).Return(nil).Once()
	quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, retainedPeginQuote.QuoteHash).Return(quote.PeginCreationDataZeroValue()).Once()
	rsk := new(mocks.RootstockRpcServerMock)

	contracts := blockchain.RskContracts{PegIn: peginContract, Bridge: bridge}
	rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
	useCase := pegin.NewCallForUserUseCase(contracts, quoteRepository, rpc, lp, eventBus, mutex)
	err := useCase.Run(context.Background(), retainedPeginQuote)
	require.ErrorContains(t, err, "not all the UTXOs are above the min lock value")
	btc.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	bridge.AssertExpectations(t)
	mutex.AssertNotCalled(t, "Lock")
	mutex.AssertNotCalled(t, "Unlock")
	peginContract.AssertNotCalled(t, "CallForUser")
	peginContract.AssertNotCalled(t, "GetBalance")
	lp.AssertNotCalled(t, "RskAddress")
}
