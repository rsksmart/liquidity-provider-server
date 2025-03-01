package watcher_test

import (
	"context"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// nolint:funlen
func TestPegoutBridgeWatcher_Start(t *testing.T) {
	ticker := &mocks.TickerMock{}
	tickerChannel := make(chan time.Time)
	ticker.EXPECT().C().Return(tickerChannel)
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	providerMock := &mocks.ProviderMock{}
	rskWallet := &mocks.RskWalletMock{}
	bridge := &mocks.BridgeMock{}
	bridge.On("GetAddress").Return(test.AnyAddress)
	mutexes := environment.NewApplicationMutexes()
	bridgeUseCase := pegout.NewBridgePegoutUseCase(pegoutRepository, providerMock, rskWallet, blockchain.RskContracts{Bridge: bridge}, mutexes.RskWalletMutex())
	getUseCase := w.NewGetWatchedPegoutQuoteUseCase(pegoutRepository)
	bridgeWatcher := watcher.NewPegoutBridgeWatcher(getUseCase, bridgeUseCase, ticker)
	resetMocks := func() {
		pegoutRepository.Calls = []mock.Call{}
		pegoutRepository.ExpectedCalls = []*mock.Call{}
		providerMock.Calls = []mock.Call{}
		providerMock.ExpectedCalls = []*mock.Call{}
		rskWallet.Calls = []mock.Call{}
		rskWallet.ExpectedCalls = []*mock.Call{}
	}
	go bridgeWatcher.Start()
	t.Run("should handle error getting quotes", func(t *testing.T) {
		resetMocks()
		checkFunc := test.AssertLogContains(t, "error getting pegout quotes")
		pegoutRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool { return checkFunc() && pegoutRepository.AssertExpectations(t) }, time.Second, 10*time.Millisecond)
	})
	const quoteHash = "0102"
	t.Run("should log error sending tx to the bridge", func(t *testing.T) {
		resetMocks()
		checkFunc := test.AssertLogContains(t, "error sending pegout to bridge")
		pegoutRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, quote.PegoutStateRefundPegOutSucceeded).Return([]quote.RetainedPegoutQuote{
			{QuoteHash: quoteHash, State: quote.PegoutStateRefundPegOutSucceeded},
		}, nil).Once()
		pegoutRepository.EXPECT().GetQuote(mock.Anything, quoteHash).Return(&quote.PegoutQuote{Value: entities.NewBigWei(math.BigPow(10, 19))}, nil).Once()
		pegoutRepository.EXPECT().GetPegoutCreationData(mock.Anything, mock.Anything).Return(quote.PegoutCreationData{GasPrice: entities.NewWei(1)}).Once()
		providerMock.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.DefaultPegoutConfiguration()).Once()
		rskWallet.On("GetBalance", mock.Anything).Return((*entities.Wei)(nil), assert.AnError).Once()
		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			return checkFunc() && rskWallet.AssertExpectations(t) && providerMock.AssertExpectations(t) && pegoutRepository.AssertExpectations(t)
		}, time.Second, 10*time.Millisecond)
	})
	t.Run("should send tx to the bridge successfully", func(t *testing.T) {
		resetMocks()
		log.SetLevel(log.DebugLevel)
		checkFunc := test.AssertLogContains(t, "transaction sent to the bridge successfully")
		pegoutRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, quote.PegoutStateRefundPegOutSucceeded).Return([]quote.RetainedPegoutQuote{
			{QuoteHash: quoteHash, State: quote.PegoutStateRefundPegOutSucceeded},
		}, nil).Once()
		pegoutRepository.EXPECT().GetQuote(mock.Anything, quoteHash).Return(&quote.PegoutQuote{Value: entities.NewBigWei(math.BigPow(10, 19))}, nil).Once()
		providerMock.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.DefaultPegoutConfiguration()).Once()
		rskWallet.On("GetBalance", mock.Anything).Return(entities.NewBigWei(math.BigPow(10, 20)), nil).Once()
		rskWallet.On("SendRbtc", mock.Anything, mock.Anything, mock.Anything).Return(test.AnyHash, nil).Once()
		pegoutRepository.EXPECT().UpdateRetainedQuotes(mock.Anything, mock.Anything).Return(nil).Once()
		pegoutRepository.EXPECT().GetPegoutCreationData(mock.Anything, mock.Anything).Return(quote.PegoutCreationData{GasPrice: entities.NewWei(1)}).Once()
		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			return checkFunc() && rskWallet.AssertExpectations(t) && providerMock.AssertExpectations(t) && pegoutRepository.AssertExpectations(t)
		}, time.Second, 10*time.Millisecond)
	})
}

func TestPegoutBridgeWatcher_Prepare(t *testing.T) {
	bridgeWatcher := watcher.NewPegoutBridgeWatcher(nil, nil, nil)
	err := bridgeWatcher.Prepare(context.Background())
	require.NoError(t, err)
}

func TestPegoutBridgeWatcher_Shutdown(t *testing.T) {
	createWatcherShutdownTest(t, func(ticker watcher.Ticker) watcher.Watcher {
		return watcher.NewPegoutBridgeWatcher(nil, nil, ticker)
	})
}
