package watcher_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	depositWatcherFirstTxID  = "2ba6da53badd14349c5d6379e88c345e88193598aad714815d4b57c691a9fbdf"
	depositWatcherSecondTxID = "fda421ccdff7324a382067d1746f6a387132435de6af336a0ebbf3f720eaae4d"
)

var (
	depositWatcherFirstWatch = rootstock.PegInWatch{
		RskAddress: "0x79568c2989232dCa1840087D73d403602364c0D4",
		BtcAddress: "2N1GMB8gxHYR5HLPSRgf9CJ9Lunjb9CTnKB",
		State:      rootstock.PegInWatchImported,
	}
	depositWatcherSecondWatch = rootstock.PegInWatch{
		RskAddress: "0x8dccd82443b80ddde3690af86746bfd9d766f8d2",
		BtcAddress: "2MvDyVY4Rbn6BZT5yTLAt6Y7HVt9LfjKbMV",
		State:      rootstock.PegInWatchImported,
	}
	depositWatcherSortedAddresses = []string{depositWatcherSecondWatch.BtcAddress, depositWatcherFirstWatch.BtcAddress}
)

type depositWatcherFixture struct {
	repository *mocks.PegInWatchRepositoryMock
	claim      *mocks.ClaimPegInUseCaseMock
	wallet     *mocks.BitcoinWalletMock
	btcRpc     *mocks.BitcoinNetworkMock
	eventBus   *mocks.EventBusMock
	imported   chan entities.Event
	ticker     *sessionTicker
}

func newDepositWatcherFixture(t *testing.T) *depositWatcherFixture {
	t.Helper()
	fixture := &depositWatcherFixture{
		repository: mocks.NewPegInWatchRepositoryMock(t),
		claim:      mocks.NewClaimPegInUseCaseMock(t),
		wallet:     mocks.NewBitcoinWalletMock(t),
		btcRpc:     mocks.NewBitcoinNetworkMock(t),
		eventBus:   &mocks.EventBusMock{},
		imported:   make(chan entities.Event, 1),
		ticker:     newSessionTicker(),
	}
	fixture.eventBus.Test(t)
	t.Cleanup(func() { fixture.eventBus.AssertExpectations(t) })
	return fixture
}

func (fixture *depositWatcherFixture) newWatcher() *watcher.PeginDepositAddressWatcher {
	useCases := watcher.NewPeginDepositAddressWatcherUseCases(
		fixture.claim,
		w.NewGetWatchedRegisteredAddressesUseCase(fixture.repository),
	)
	return watcher.NewPeginDepositAddressWatcher(
		useCases,
		fixture.wallet,
		blockchain.Rpc{Btc: fixture.btcRpc},
		fixture.eventBus,
		fixture.ticker,
	)
}

func (fixture *depositWatcherFixture) expectSubscription() {
	fixture.eventBus.On("Subscribe", blockchain.PegInAddressesImportedEventId).
		Return((<-chan entities.Event)(fixture.imported)).
		Once()
}

func (fixture *depositWatcherFixture) start(t *testing.T, watches ...rootstock.PegInWatch) *watcherSession {
	t.Helper()
	fixture.expectSubscription()
	fixture.repository.EXPECT().List(mock.Anything).Return(watches, nil).Once()
	return startWatcherSession(t, fixture.newWatcher(), fixture.ticker)
}

func (fixture *depositWatcherFixture) expectHeight(height int64) {
	fixture.btcRpc.EXPECT().GetHeight().Return(big.NewInt(height), nil).Once()
}

func (fixture *depositWatcherFixture) expectOutputs(addresses []string, outputs ...blockchain.BitcoinUnspentOutput) {
	fixture.wallet.EXPECT().GetUnspentOutputs(addresses).Return(outputs, nil).Once()
}

func (fixture *depositWatcherFixture) expectClaim(watch rootstock.PegInWatch, txID string) {
	fixture.claim.EXPECT().Run(mock.Anything, watch, txID).Return(nil).Once()
}

func depositOutput(watch rootstock.PegInWatch, txID string, vout uint32) blockchain.BitcoinUnspentOutput {
	return blockchain.BitcoinUnspentOutput{TxID: txID, Vout: vout, Address: watch.BtcAddress, Confirmations: 1}
}

func TestPeginDepositAddressWatcher_Shutdown(t *testing.T) {
	createWatcherShutdownTest(t, func(ticker utils.Ticker) watcher.Watcher {
		useCases := watcher.NewPeginDepositAddressWatcherUseCases(nil, nil)
		eventBus := &mocks.EventBusMock{}
		return watcher.NewPeginDepositAddressWatcher(useCases, nil, blockchain.Rpc{}, eventBus, ticker)
	})
}

func TestPeginDepositAddressWatcher_Prepare(t *testing.T) {
	t.Run("returns the error when the watched addresses cannot be loaded", func(t *testing.T) {
		fixture := newDepositWatcherFixture(t)
		fixture.expectSubscription()
		fixture.repository.EXPECT().List(mock.Anything).Return(nil, assert.AnError).Once()
		err := fixture.newWatcher().Prepare(context.Background())
		require.ErrorIs(t, err, assert.AnError)
	})
}

func TestPeginDepositAddressWatcher_ClaimsTheDepositsOfTheImportedAddresses(t *testing.T) {
	fixture := newDepositWatcherFixture(t)
	discovered := rootstock.PegInWatch{RskAddress: "0x1", BtcAddress: "2N8hwP1WmJrFF5QWABn38y63uYLhnJYJYTF", State: rootstock.PegInWatchDiscovered}
	unsupported := rootstock.PegInWatch{RskAddress: "0x2", State: rootstock.PegInWatchUnsupportedEncoding}
	session := fixture.start(t, depositWatcherFirstWatch, discovered, depositWatcherSecondWatch, unsupported)
	fixture.expectHeight(100)
	fixture.expectOutputs(
		depositWatcherSortedAddresses,
		depositOutput(depositWatcherFirstWatch, depositWatcherFirstTxID, 0),
		depositOutput(depositWatcherSecondWatch, depositWatcherSecondTxID, 1),
	)
	fixture.expectClaim(depositWatcherFirstWatch, depositWatcherFirstTxID)
	fixture.expectClaim(depositWatcherSecondWatch, depositWatcherSecondTxID)

	session.poll()
	session.stop()
}

func TestPeginDepositAddressWatcher_ChecksOnlyWhenANewBlockArrives(t *testing.T) {
	fixture := newDepositWatcherFixture(t)
	session := fixture.start(t, depositWatcherFirstWatch)
	output := depositOutput(depositWatcherFirstWatch, depositWatcherFirstTxID, 0)

	fixture.expectHeight(100)
	fixture.expectOutputs([]string{depositWatcherFirstWatch.BtcAddress}, output)
	fixture.expectClaim(depositWatcherFirstWatch, depositWatcherFirstTxID)
	session.poll()

	fixture.expectHeight(100)
	session.poll()

	fixture.expectHeight(99)
	session.poll()

	fixture.expectHeight(101)
	fixture.expectOutputs([]string{depositWatcherFirstWatch.BtcAddress}, output)
	fixture.expectClaim(depositWatcherFirstWatch, depositWatcherFirstTxID)
	session.poll()

	session.stop()
}

func TestPeginDepositAddressWatcher_SkipsTheTickWhenTheHeightCannotBeRead(t *testing.T) {
	fixture := newDepositWatcherFixture(t)
	session := fixture.start(t, depositWatcherFirstWatch)
	assertLog := test.AssertLogContains(t, fmt.Sprintf(watcher.LogPeginBtcChainHeight, assert.AnError))
	fixture.btcRpc.EXPECT().GetHeight().Return(nil, assert.AnError).Once()

	session.poll()

	assert.True(t, assertLog())
	session.stop()
}

func TestPeginDepositAddressWatcher_ChecksTheSameBlockAgainWhenTheWalletCallFails(t *testing.T) {
	fixture := newDepositWatcherFixture(t)
	session := fixture.start(t, depositWatcherFirstWatch)
	assertLog := test.AssertLogContains(t, fmt.Sprintf(watcher.LogPeginBtcListDepositsError, assert.AnError))
	fixture.expectHeight(100)
	fixture.wallet.EXPECT().GetUnspentOutputs([]string{depositWatcherFirstWatch.BtcAddress}).Return(nil, assert.AnError).Once()
	session.poll()
	assert.True(t, assertLog())

	fixture.expectHeight(100)
	fixture.expectOutputs([]string{depositWatcherFirstWatch.BtcAddress}, depositOutput(depositWatcherFirstWatch, depositWatcherFirstTxID, 0))
	fixture.expectClaim(depositWatcherFirstWatch, depositWatcherFirstTxID)
	session.poll()

	session.stop()
}

func TestPeginDepositAddressWatcher_ClaimsEachDepositOncePerCheck(t *testing.T) {
	fixture := newDepositWatcherFixture(t)
	session := fixture.start(t, depositWatcherFirstWatch, depositWatcherSecondWatch)
	fixture.expectHeight(100)
	fixture.expectOutputs(
		depositWatcherSortedAddresses,
		depositOutput(depositWatcherFirstWatch, depositWatcherFirstTxID, 0),
		depositOutput(depositWatcherFirstWatch, depositWatcherFirstTxID, 1),
		depositOutput(depositWatcherSecondWatch, depositWatcherFirstTxID, 2),
	)
	fixture.expectClaim(depositWatcherFirstWatch, depositWatcherFirstTxID)
	fixture.expectClaim(depositWatcherSecondWatch, depositWatcherFirstTxID)

	session.poll()
	session.stop()
}

func TestPeginDepositAddressWatcher_IgnoresOutputsToAddressesItDoesNotWatch(t *testing.T) {
	fixture := newDepositWatcherFixture(t)
	session := fixture.start(t, depositWatcherFirstWatch)
	fixture.expectHeight(100)
	fixture.expectOutputs(
		[]string{depositWatcherFirstWatch.BtcAddress},
		depositOutput(depositWatcherSecondWatch, depositWatcherSecondTxID, 0),
	)

	session.poll()
	session.stop()
}

func TestPeginDepositAddressWatcher_KeepsClaimingWhenOneClaimFails(t *testing.T) {
	fixture := newDepositWatcherFixture(t)
	session := fixture.start(t, depositWatcherFirstWatch, depositWatcherSecondWatch)
	assertLog := test.AssertLogContains(t, watcher.LogPeginBtcClaimError(depositWatcherFirstWatch.RskAddress, depositWatcherFirstTxID, assert.AnError))
	fixture.expectHeight(100)
	fixture.expectOutputs(
		depositWatcherSortedAddresses,
		depositOutput(depositWatcherFirstWatch, depositWatcherFirstTxID, 0),
		depositOutput(depositWatcherSecondWatch, depositWatcherSecondTxID, 0),
	)
	fixture.claim.EXPECT().Run(mock.Anything, depositWatcherFirstWatch, depositWatcherFirstTxID).Return(assert.AnError).Once()
	fixture.expectClaim(depositWatcherSecondWatch, depositWatcherSecondTxID)

	session.poll()

	assert.True(t, assertLog())
	session.stop()
}

func TestPeginDepositAddressWatcher_ChecksWithNoAddressesWhenNoneIsImported(t *testing.T) {
	fixture := newDepositWatcherFixture(t)
	session := fixture.start(t)
	fixture.expectHeight(100)
	fixture.expectOutputs([]string(nil))

	session.poll()
	session.stop()
}

func TestPeginDepositAddressWatcher_LoadsTheImportedAddressesAgainAfterARestart(t *testing.T) {
	fixture := newDepositWatcherFixture(t)
	firstSession := fixture.start(t, depositWatcherFirstWatch)
	fixture.expectHeight(100)
	fixture.expectOutputs([]string{depositWatcherFirstWatch.BtcAddress})
	firstSession.poll()
	firstSession.stop()

	fixture.ticker = newSessionTicker()
	fixture.imported = make(chan entities.Event, 1)
	secondSession := fixture.start(t, depositWatcherFirstWatch, depositWatcherSecondWatch)
	fixture.expectHeight(100)
	fixture.expectOutputs(depositWatcherSortedAddresses, depositOutput(depositWatcherSecondWatch, depositWatcherSecondTxID, 0))
	fixture.expectClaim(depositWatcherSecondWatch, depositWatcherSecondTxID)
	secondSession.poll()
	secondSession.stop()
}

func (fixture *depositWatcherFixture) publish(t *testing.T, session *watcherSession, event entities.Event) {
	t.Helper()
	idleBefore := fixture.ticker.selects.Load()
	fixture.imported <- event
	session.waitUntilIdle(idleBefore)
}

func importedEvent(watches ...rootstock.PegInWatch) blockchain.PegInAddressesImportedEvent {
	return blockchain.PegInAddressesImportedEvent{
		BaseEvent: entities.NewBaseEvent(blockchain.PegInAddressesImportedEventId),
		Watches:   watches,
	}
}

func TestPeginDepositAddressWatcher_WatchesAnAddressImportedWhileRunning(t *testing.T) {
	fixture := newDepositWatcherFixture(t)
	session := fixture.start(t, depositWatcherFirstWatch)
	discovered := rootstock.PegInWatch{RskAddress: "0x1", BtcAddress: "2N8hwP1WmJrFF5QWABn38y63uYLhnJYJYTF", State: rootstock.PegInWatchDiscovered}
	fixture.publish(t, session, importedEvent(depositWatcherSecondWatch, discovered))
	fixture.expectHeight(100)
	fixture.expectOutputs(depositWatcherSortedAddresses, depositOutput(depositWatcherSecondWatch, depositWatcherSecondTxID, 0))
	fixture.expectClaim(depositWatcherSecondWatch, depositWatcherSecondTxID)

	session.poll()
	session.stop()
}

func TestPeginDepositAddressWatcher_ChecksAnAddressImportedAfterTheBlockWasChecked(t *testing.T) {
	fixture := newDepositWatcherFixture(t)
	session := fixture.start(t, depositWatcherFirstWatch)
	fixture.expectHeight(100)
	fixture.expectOutputs([]string{depositWatcherFirstWatch.BtcAddress})
	session.poll()

	fixture.publish(t, session, importedEvent(depositWatcherSecondWatch))
	fixture.expectHeight(100)
	fixture.expectOutputs(depositWatcherSortedAddresses, depositOutput(depositWatcherSecondWatch, depositWatcherSecondTxID, 0))
	fixture.expectClaim(depositWatcherSecondWatch, depositWatcherSecondTxID)
	session.poll()

	fixture.expectHeight(100)
	session.poll()

	session.stop()
}

func TestPeginDepositAddressWatcher_ChecksAnImportedAddressAgainWhenTheWalletCallFails(t *testing.T) {
	fixture := newDepositWatcherFixture(t)
	session := fixture.start(t, depositWatcherFirstWatch)
	fixture.expectHeight(100)
	fixture.expectOutputs([]string{depositWatcherFirstWatch.BtcAddress})
	session.poll()

	fixture.publish(t, session, importedEvent(depositWatcherSecondWatch))
	fixture.expectHeight(100)
	fixture.wallet.EXPECT().GetUnspentOutputs(depositWatcherSortedAddresses).Return(nil, assert.AnError).Once()
	session.poll()

	fixture.expectHeight(100)
	fixture.expectOutputs(depositWatcherSortedAddresses, depositOutput(depositWatcherSecondWatch, depositWatcherSecondTxID, 0))
	fixture.expectClaim(depositWatcherSecondWatch, depositWatcherSecondTxID)
	session.poll()

	session.stop()
}

func TestPeginDepositAddressWatcher_IgnoresAnEventOfAnotherType(t *testing.T) {
	fixture := newDepositWatcherFixture(t)
	session := fixture.start(t, depositWatcherFirstWatch)
	fixture.expectHeight(100)
	fixture.expectOutputs([]string{depositWatcherFirstWatch.BtcAddress})
	session.poll()

	assertLog := test.AssertLogContains(t, watcher.LogPeginBtcWrongEvent)
	fixture.publish(t, session, blockchain.PegInAddressRegistryRebuildEvent{
		BaseEvent: entities.NewBaseEvent(blockchain.PegInAddressesImportedEventId),
	})
	assert.True(t, assertLog())
	fixture.expectHeight(100)
	session.poll()

	session.stop()
}

func TestPeginDepositAddressWatcher_StopsReadingTheEventChannelOnceItIsClosed(t *testing.T) {
	fixture := newDepositWatcherFixture(t)
	session := fixture.start(t, depositWatcherFirstWatch)
	assertNoLog := test.AssertNoLog(t)
	idleBefore := fixture.ticker.selects.Load()
	close(fixture.imported)
	session.waitUntilIdle(idleBefore)
	idleAfterClose := fixture.ticker.selects.Load()
	assert.Never(t, func() bool { return fixture.ticker.selects.Load() > idleAfterClose }, 50*time.Millisecond, time.Millisecond)

	fixture.expectHeight(100)
	fixture.expectOutputs([]string{depositWatcherFirstWatch.BtcAddress})
	session.poll()
	session.stop()
	assertNoLog()
}
