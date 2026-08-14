package watcher

import (
	"context"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

// Watchers intervals
const (
	quoteCleanInterval               = 10 * time.Minute
	peginDepositWatcherInterval      = 1 * time.Minute
	peginBridgeWatcherInterval       = 3 * time.Minute
	pegoutDepositWatcherInterval     = 1 * time.Minute
	pegoutBtcTransferWatcherInterval = 3 * time.Minute
	pegoutBridgeWatcherInterval      = 5 * time.Minute
	liquidityCheckInterval           = 10 * time.Minute
	penalizationCheckInterval        = 10 * time.Minute
	bitcoinEclipseCheckInterval      = 5 * time.Minute
	rskEclipseCheckInterval          = 15 * time.Second
	btcReleaseCheckInterval          = 3 * time.Minute
	assetMetricsUpdateInterval       = 1 * time.Minute
	transferColdWalletInterval       = 5 * time.Minute
	bitcoinReorgCheckInterval        = 5 * time.Minute
	rootstockReorgCheckInterval      = 30 * time.Second
	bitcoinPeerCheckInterval         = 1 * time.Minute
	rootstockPeerCheckInterval       = 1 * time.Minute
	// Accepted miss window, not a protocol invariant. Deposits older than 100 BTC blocks
	// (~16.7h, Powpeg's confirmation window) stay invisible: late registration, LPS downtime,
	// and first enable against an already-populated registry will miss them.
	peginAddressRegistryRescanDepthBlocks int64 = 100
)

type Watcher interface {
	entities.Closeable
	Prepare(ctx context.Context) error
	Start()
}
