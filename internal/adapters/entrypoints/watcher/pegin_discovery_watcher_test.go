package watcher_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	peginEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/pegin"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	discoveryRskAddress = "0x892813507Bf3aBF2890759d2135Ec34f4909Fea5"
	discoveryBtcAddress = "mzMsGv8XYJ3Sj2k7AQ8s8Yj8gqpZ1Z1Z1Z" // printable; treated as the deposit address
)

func newDiscoveryClaimUseCase(
	peginContract *mocks.PeginContractMock,
	registry *mocks.PegInAddressRegistryContractMock,
	configurations *mocks.FlyoverConfigurationsContractMock,
	rsk *mocks.RootstockRpcServerMock,
	btc *mocks.BtcRpcMock,
	lp *mocks.ProviderMock,
	mutex *mocks.MutexMock,
) *pegin.ClaimPegInUseCase {
	contracts := blockchain.RskContracts{
		PegIn:                 peginContract,
		PegInAddressRegistry:  registry,
		FlyoverConfigurations: configurations,
	}
	return pegin.NewClaimPegInUseCase(contracts, blockchain.Rpc{Btc: btc, Rsk: rsk}, lp, lp, mutex)
}

func TestPeginDiscoveryWatcher_Shutdown(t *testing.T) {
	createWatcherShutdownTest(t, func(ticker utils.Ticker) watcher.Watcher {
		return watcher.NewPeginDiscoveryWatcher(nil, new(mocks.PegInAddressRegistryContractMock), new(mocks.FlyoverConfigurationsContractMock), new(mocks.BitcoinWalletMock), blockchain.Rpc{}, new(mocks.EventBusMock), ticker)
	})
}

func TestPeginDiscoveryWatcher_Prepare_SetsCheckpoint(t *testing.T) {
	rsk := new(mocks.RootstockRpcServerMock)
	rsk.EXPECT().GetHeight(mock.Anything).Return(uint64(1000), nil).Once()
	w := watcher.NewPeginDiscoveryWatcher(nil, new(mocks.PegInAddressRegistryContractMock), new(mocks.FlyoverConfigurationsContractMock), new(mocks.BitcoinWalletMock), blockchain.Rpc{Rsk: rsk}, new(mocks.EventBusMock), nil)
	require.NoError(t, w.Prepare(context.Background()))
	rsk.AssertExpectations(t)
}

// TestPeginDiscoveryWatcher_DiscoversAndClaims drives a full tick: a new AddressRegistered event is
// discovered, its deposit address imported and watched, a confirmed BTC deposit is seen, and the LP
// fronts RBTC via requestPegIn.
func TestPeginDiscoveryWatcher_DiscoversAndClaims(t *testing.T) {
	registry := new(mocks.PegInAddressRegistryContractMock)
	configurations := new(mocks.FlyoverConfigurationsContractMock)
	peginContract := new(mocks.PeginContractMock)
	peginContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil).Maybe()
	peginContract.EXPECT().GetAddress().Return("pegin").Maybe()
	rsk := new(mocks.RootstockRpcServerMock)
	btcRpc := new(mocks.BtcRpcMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	lp := new(mocks.ProviderMock)
	lp.On("RskAddress").Return(discoveryRskAddress).Maybe()
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return().Maybe()
	mutex.On("Unlock").Return().Maybe()

	depositAmount := entities.NewWei(1_000_000)

	// Prepare checkpoint at block 100, then discovery sees block 110 with one event.
	rsk.EXPECT().GetHeight(mock.Anything).Return(uint64(100), nil).Once()  // Prepare
	rsk.EXPECT().GetHeight(mock.Anything).Return(uint64(110), nil).Maybe() // discoverNewAddresses
	registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(101), mock.Anything).Return([]blockchain.AddressRegisteredEvent{
		{RskAddress: discoveryRskAddress, BlockNumber: 105},
	}, nil).Once()
	registry.EXPECT().GetPegInAddress(discoveryRskAddress).Return(blockchain.PegInDepositAddress{
		RskAddress:      discoveryRskAddress,
		ScriptOrAddress: []byte(discoveryBtcAddress),
		Encoding:        0,
	}, nil).Once()
	btcWallet.EXPECT().ImportAddress(discoveryBtcAddress).Return(nil).Once()

	// A confirmed deposit appears at the watched address.
	btcWallet.EXPECT().GetTransactions(discoveryBtcAddress).Return([]blockchain.BitcoinTransactionInformation{
		{Hash: "0a1b", Confirmations: 10, Outputs: map[string][]*entities.Wei{discoveryBtcAddress: {depositAmount}}},
	}, nil).Maybe()
	btcRpc.On("GetTransactionInfo", "0a1b").Return(blockchain.BitcoinTransactionInformation{Hash: "0a1b", Confirmations: 10}, nil).Maybe()

	// Claim: registered, fee read, own liquidity sufficient, requestPegIn wins.
	registry.EXPECT().IsRegistered(discoveryRskAddress).Return(true, nil).Maybe()
	configurations.EXPECT().CalculatePegInFee(mock.Anything).Return(entities.NewWei(10_000), nil).Maybe()
	rsk.EXPECT().GetBalance(mock.Anything, mock.Anything).Return(entities.NewWei(5_000_000), nil).Maybe()
	requested := make(chan struct{}, 1)
	peginContract.EXPECT().RequestPegIn(mock.Anything).Run(func(_ blockchain.RequestPegInParams) {
		select {
		case requested <- struct{}{}:
		default:
		}
	}).Return(blockchain.TransactionReceipt{TransactionHash: "0xreq"}, nil).Maybe()

	claimUseCase := newDiscoveryClaimUseCase(peginContract, registry, configurations, rsk, btcRpc, lp, mutex)

	tickerChannel := make(chan time.Time, 1)
	ticker := new(mocks.TickerMock)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop().Return()

	w := watcher.NewPeginDiscoveryWatcher(claimUseCase, registry, configurations, btcWallet, blockchain.Rpc{Btc: btcRpc, Rsk: rsk}, new(mocks.EventBusMock), ticker)
	require.NoError(t, w.Prepare(context.Background()))

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() { w.Start(); wg.Done() }()

	tickerChannel <- time.Now() // drive one tick

	select {
	case <-requested:
		// requestPegIn was invoked — discovery → claim wired end to end.
	case <-time.After(2 * time.Second):
		t.Fatal("requestPegIn was not called after discovery tick")
	}

	// The watched address should now be in Requested state.
	require.Eventually(t, func() bool {
		w, ok := w.GetWatchedAddress(discoveryRskAddress)
		return ok && w.State == peginEntity.WatchedAddressStateRequested
	}, 2*time.Second, 10*time.Millisecond)

	closeChannel := make(chan bool, 1)
	w.Shutdown(closeChannel)
	<-closeChannel
	wg.Wait()

	peginContract.AssertCalled(t, "RequestPegIn", mock.Anything)
	btcWallet.AssertExpectations(t)
}

func TestPeginDiscoveryWatcher_DiscoverSkipsUnprintableAddress(t *testing.T) {
	registry := new(mocks.PegInAddressRegistryContractMock)
	rsk := new(mocks.RootstockRpcServerMock)
	btcWallet := new(mocks.BitcoinWalletMock)

	rsk.EXPECT().GetHeight(mock.Anything).Return(uint64(10), nil).Once()  // Prepare
	rsk.EXPECT().GetHeight(mock.Anything).Return(uint64(20), nil).Maybe() // discover
	registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(11), mock.Anything).Return([]blockchain.AddressRegisteredEvent{
		{RskAddress: discoveryRskAddress, BlockNumber: 15},
	}, nil).Maybe()
	// Non-printable script bytes: the watcher logs and skips, never importing.
	registry.EXPECT().GetPegInAddress(discoveryRskAddress).Return(blockchain.PegInDepositAddress{
		RskAddress:      discoveryRskAddress,
		ScriptOrAddress: []byte{0x00, 0x01, 0x02},
		Encoding:        1,
	}, nil).Maybe()

	ticker := new(mocks.TickerMock)
	tickerChannel := make(chan time.Time, 1)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop().Return()

	w := watcher.NewPeginDiscoveryWatcher(nil, registry, new(mocks.FlyoverConfigurationsContractMock), btcWallet, blockchain.Rpc{Rsk: rsk}, new(mocks.EventBusMock), ticker)
	require.NoError(t, w.Prepare(context.Background()))

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() { w.Start(); wg.Done() }()
	tickerChannel <- time.Now()

	// Give the tick a moment, then confirm the address was NOT added to the watch-set.
	require.Never(t, func() bool {
		_, ok := w.GetWatchedAddress(discoveryRskAddress)
		return ok
	}, 300*time.Millisecond, 30*time.Millisecond)

	closeChannel := make(chan bool, 1)
	w.Shutdown(closeChannel)
	<-closeChannel
	wg.Wait()

	btcWallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
}
