package watcher_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

var now = uint32(time.Now().Unix())

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
	PenaltyFee:            100,
	Nonce:                 123456,
	DepositAddress:        "any address",
	Value:                 entities.NewWei(5000),
	AgreementTimestamp:    now - 60,
	DepositDateLimit:      now + 60,
	DepositConfirmations:  10,
	TransferConfirmations: 10,
	TransferTime:          600,
	ExpireDate:            now + 600,
	ExpireBlock:           500,
	GasFee:                entities.NewWei(500),
	ProductFeeAmount:      300,
}

func TestUpdatePegoutQuoteDepositUseCase_Run(t *testing.T) {
	deposit := quote.PegoutDeposit{
		TxHash:      "user rsk tx hash",
		QuoteHash:   depositedRetainedQuote.QuoteHash,
		Amount:      entities.NewWei(6800),
		Timestamp:   time.Now(),
		BlockNumber: 480,
		From:        "0x1a1b1c",
	}
	quoteReporitory := new(test.PegoutQuoteRepositoryMock)
	quoteReporitory.On(
		"UpdateRetainedQuote",
		mock.AnythingOfType("context.backgroundCtx"),
		mock.MatchedBy(func(q quote.RetainedPegoutQuote) bool {
			return q.UserRskTxHash == deposit.TxHash &&
				q.State == quote.PegoutStateWaitingForDepositConfirmations
		}),
	).Return(nil)
	quoteReporitory.On("UpsertPegoutDeposit", mock.AnythingOfType("context.backgroundCtx"), deposit).Return(nil)
	useCase := watcher.NewUpdatePegoutQuoteDepositUseCase(quoteReporitory)
	watchedPegoutQuote, err := useCase.Run(context.Background(), watcher.NewWatchedPegoutQuote(depositedPegoutQuote, depositedRetainedQuote), deposit)
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
				TxHash:      "user rsk tx hash",
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
				TxHash:      "user rsk tx hash",
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
				TxHash:      "user rsk tx hash",
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
			quoteReporitory := new(test.PegoutQuoteRepositoryMock)
			useCase := watcher.NewUpdatePegoutQuoteDepositUseCase(quoteReporitory)
			watchedPegoutQuote, err := useCase.Run(context.Background(), watcher.NewWatchedPegoutQuote(depositedPegoutQuote, depositedRetainedQuote), testCase.deposit)
			quoteReporitory.AssertNotCalled(t, "UpdateRetainedQuote")
			quoteReporitory.AssertNotCalled(t, "UpsertPegoutDeposit")
			assert.Equal(t, watcher.WatchedPegoutQuote{}, watchedPegoutQuote)
			require.Error(t, err)
			assert.True(t, strings.Contains(err.Error(), "deposit not valid for quote"))
		})
	}
}

func TestUpdatePegoutQuoteDepositUseCase_Run_IllegalState(t *testing.T) {
	deposit := quote.PegoutDeposit{
		TxHash:      "user rsk tx hash",
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
		quoteReporitory := new(test.PegoutQuoteRepositoryMock)
		useCase := watcher.NewUpdatePegoutQuoteDepositUseCase(quoteReporitory)
		watchedPegoutQuote, err := useCase.Run(context.Background(), watcher.NewWatchedPegoutQuote(depositedPegoutQuote, retainedQuote), deposit)
		quoteReporitory.AssertNotCalled(t, "UpdateRetainedQuote")
		quoteReporitory.AssertNotCalled(t, "UpsertPegoutDeposit")
		assert.Equal(t, watcher.WatchedPegoutQuote{}, watchedPegoutQuote)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "illegal quote state"))
	}
}

func TestUpdatePegoutQuoteDepositUseCase_Run_ErrorHandling(t *testing.T) {
	deposit := quote.PegoutDeposit{
		TxHash:      "user rsk tx hash",
		QuoteHash:   depositedRetainedQuote.QuoteHash,
		Amount:      entities.NewWei(6800),
		Timestamp:   time.Now(),
		BlockNumber: 480,
		From:        "0x1a1b1c",
	}

	setups := []func(quoteRepository *test.PegoutQuoteRepositoryMock){
		func(quoteRepository *test.PegoutQuoteRepositoryMock) {
			quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(assert.AnError)
		},
		func(quoteRepository *test.PegoutQuoteRepositoryMock) {
			quoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil)
			quoteRepository.On("UpsertPegoutDeposit", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(assert.AnError)
		},
	}

	for _, setup := range setups {
		quoteReporitory := new(test.PegoutQuoteRepositoryMock)
		setup(quoteReporitory)
		useCase := watcher.NewUpdatePegoutQuoteDepositUseCase(quoteReporitory)
		watchedPegoutQuote, err := useCase.Run(context.Background(), watcher.NewWatchedPegoutQuote(depositedPegoutQuote, depositedRetainedQuote), deposit)
		quoteReporitory.AssertExpectations(t)
		assert.Equal(t, watcher.WatchedPegoutQuote{}, watchedPegoutQuote)
		require.Error(t, err)
	}
}
