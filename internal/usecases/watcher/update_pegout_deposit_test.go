package watcher_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	now       = uint32(time.Now().Unix())
	userRskTx = "user rsk tx hash"
)

var depositedRetainedQuote = quote.RetainedPegoutQuote{
	QuoteHash:         "02011d",
	DepositAddress:    "0xabcd",
	Signature:         "signature",
	RequiredLiquidity: entities.NewWei(5500),
	State:             quote.PegoutStateWaitingForDeposit,
}

var depositedPegoutQuote = quote.PegoutQuote{
	LbcAddress:            depositedRetainedQuote.DepositAddress,
	LpRskAddress:          "0x1234",
	BtcRefundAddress:      "0x1234",
	RskRefundAddress:      "0x1234",
	LpBtcAddress:          "0x1234",
	CallFee:               entities.NewWei(1000),
	PenaltyFee:            entities.NewWei(100),
	Nonce:                 quote.NewNonce(123456),
	DepositAddress:        test.AnyAddress,
	Value:                 entities.NewWei(5000),
	AgreementTimestamp:    now - 60,
	DepositDateLimit:      now + 60,
	DepositConfirmations:  10,
	TransferConfirmations: 10,
	TransferTime:          600,
	ExpireDate:            now + 600,
	ExpireBlock:           500,
	GasFee:                entities.NewWei(500),
	ProductFeeAmount:      entities.NewWei(300),
}

var depositedPegoutCreationData = quote.PegoutCreationData{
	GasPrice:      entities.NewWei(5),
	FeePercentage: utils.NewBigFloat64(1.5),
	FeeRate:       utils.NewBigFloat64(111.57),
	FixedFee:      entities.NewWei(7),
}

func TestUpdatePegoutQuoteDepositUseCase_Run(t *testing.T) {
	deposit := quote.PegoutDeposit{
		TxHash:      userRskTx,
		QuoteHash:   depositedRetainedQuote.QuoteHash,
		Amount:      entities.NewWei(6800),
		Timestamp:   time.Now(),
		BlockNumber: 480,
		From:        "0x1a1b1c",
	}
	quoteReporitory := new(mocks.PegoutQuoteRepositoryMock)
	quoteReporitory.On(
		"UpdateRetainedQuote",
		test.AnyCtx,
		mock.MatchedBy(func(q quote.RetainedPegoutQuote) bool {
			return q.UserRskTxHash == deposit.TxHash &&
				q.State == quote.PegoutStateWaitingForDepositConfirmations
		}),
	).Return(nil)
	quoteReporitory.On("UpsertPegoutDeposit", test.AnyCtx, deposit).Return(nil)
	useCase := watcher.NewUpdatePegoutQuoteDepositUseCase(quoteReporitory)
	watchedPegoutQuote, err := useCase.Run(context.Background(), quote.NewWatchedPegoutQuote(depositedPegoutQuote, depositedRetainedQuote, depositedPegoutCreationData), deposit)
	quoteReporitory.AssertExpectations(t)
	require.NoError(t, err)
	assert.Equal(t, quote.PegoutStateWaitingForDepositConfirmations, watchedPegoutQuote.RetainedQuote.State)
	assert.Equal(t, deposit.TxHash, watchedPegoutQuote.RetainedQuote.UserRskTxHash)
}

func TestUpdatePegoutQuoteDepositUseCase_Run_NotValid(t *testing.T) {
	cases := []struct {
		name    string
		deposit quote.PegoutDeposit
	}{
		{
			name: "Should fail by value",
			deposit: quote.PegoutDeposit{
				TxHash:      userRskTx,
				QuoteHash:   depositedRetainedQuote.QuoteHash,
				Amount:      entities.NewWei(6000),
				Timestamp:   time.Now(),
				BlockNumber: 480,
				From:        "0x1a1b1c",
			},
		},
		{
			name: "Should fail by time",
			deposit: quote.PegoutDeposit{
				TxHash:      userRskTx,
				QuoteHash:   depositedRetainedQuote.QuoteHash,
				Amount:      entities.NewWei(6500),
				Timestamp:   time.Unix(time.Now().Unix()+660, 0),
				BlockNumber: 480,
				From:        "0x1a1b1c",
			},
		},
		{
			name: "Should fail by confirmations",
			deposit: quote.PegoutDeposit{
				TxHash:      userRskTx,
				QuoteHash:   depositedRetainedQuote.QuoteHash,
				Amount:      entities.NewWei(6500),
				Timestamp:   time.Now(),
				BlockNumber: 501,
				From:        "0x1a1b1c",
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			quoteReporitory := new(mocks.PegoutQuoteRepositoryMock)
			useCase := watcher.NewUpdatePegoutQuoteDepositUseCase(quoteReporitory)
			creationData := quote.PegoutCreationDataZeroValue()
			watchedPegoutQuote, err := useCase.Run(context.Background(), quote.NewWatchedPegoutQuote(depositedPegoutQuote, depositedRetainedQuote, creationData), testCase.deposit)
			quoteReporitory.AssertNotCalled(t, "UpdateRetainedQuote")
			quoteReporitory.AssertNotCalled(t, "UpsertPegoutDeposit")
			assert.Equal(t, quote.WatchedPegoutQuote{}, watchedPegoutQuote)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "deposit not valid for quote")
		})
	}
}

func TestUpdatePegoutQuoteDepositUseCase_Run_IllegalState(t *testing.T) {
	deposit := quote.PegoutDeposit{
		TxHash:      userRskTx,
		QuoteHash:   "02011d",
		Amount:      entities.NewWei(6800),
		Timestamp:   time.Now(),
		BlockNumber: 480,
		From:        "0x1a1b1c",
	}
	quotes := []quote.RetainedPegoutQuote{
		{State: quote.PegoutStateWaitingForDepositConfirmations},
		{State: quote.PegoutStateSendPegoutSucceeded},
		{State: quote.PegoutStateSendPegoutFailed},
		{State: quote.PegoutStateRefundPegOutSucceeded},
		{State: quote.PegoutStateRefundPegOutFailed},
		{State: quote.PegoutStateTimeForDepositElapsed},
	}
	for _, retainedQuote := range quotes {
		quoteReporitory := new(mocks.PegoutQuoteRepositoryMock)
		useCase := watcher.NewUpdatePegoutQuoteDepositUseCase(quoteReporitory)
		creationData := quote.PegoutCreationDataZeroValue()
		watchedPegoutQuote, err := useCase.Run(context.Background(), quote.NewWatchedPegoutQuote(depositedPegoutQuote, retainedQuote, creationData), deposit)
		quoteReporitory.AssertNotCalled(t, "UpdateRetainedQuote")
		quoteReporitory.AssertNotCalled(t, "UpsertPegoutDeposit")
		assert.Equal(t, quote.WatchedPegoutQuote{}, watchedPegoutQuote)
		require.Error(t, err)
		require.ErrorIs(t, err, usecases.IllegalQuoteStateError)
	}
}

func TestUpdatePegoutQuoteDepositUseCase_Run_ErrorHandling(t *testing.T) {
	deposit := quote.PegoutDeposit{
		TxHash:      userRskTx,
		QuoteHash:   depositedRetainedQuote.QuoteHash,
		Amount:      entities.NewWei(6800),
		Timestamp:   time.Now(),
		BlockNumber: 480,
		From:        "0x1a1b1c",
	}

	setups := []func(quoteRepository *mocks.PegoutQuoteRepositoryMock){
		func(quoteRepository *mocks.PegoutQuoteRepositoryMock) {
			quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, mock.Anything).Return(assert.AnError)
		},
		func(quoteRepository *mocks.PegoutQuoteRepositoryMock) {
			quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, mock.Anything).Return(nil)
			quoteRepository.On("UpsertPegoutDeposit", test.AnyCtx, mock.Anything).Return(assert.AnError)
		},
	}

	for _, setup := range setups {
		quoteReporitory := new(mocks.PegoutQuoteRepositoryMock)
		setup(quoteReporitory)
		useCase := watcher.NewUpdatePegoutQuoteDepositUseCase(quoteReporitory)
		creationData := quote.PegoutCreationDataZeroValue()
		watchedPegoutQuote, err := useCase.Run(context.Background(), quote.NewWatchedPegoutQuote(depositedPegoutQuote, depositedRetainedQuote, creationData), deposit)
		quoteReporitory.AssertExpectations(t)
		assert.Equal(t, quote.WatchedPegoutQuote{}, watchedPegoutQuote)
		require.Error(t, err)
	}
}
