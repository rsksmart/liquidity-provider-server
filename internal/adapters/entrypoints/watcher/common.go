package watcher

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"time"
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
)

type Watcher interface {
	entities.Closeable
	Prepare(ctx context.Context) error
	Start()
}
