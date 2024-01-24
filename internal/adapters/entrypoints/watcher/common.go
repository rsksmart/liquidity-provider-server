package watcher

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"time"
)

const (
	quoteCleanInterval               = 10 * time.Minute
	peginDepositWatcherInterval      = 1 * time.Minute
	peginBridgeWatcherInterval       = 3 * time.Minute
	pegoutDepositWatcherInterval     = 1 * time.Minute
	pegoutBtcTransferWatcherInterval = 3 * time.Minute
	liquidityCheckInterval           = 10 * time.Minute
	penalizationCheckInterval        = 10 * time.Minute
)

type Watcher interface {
	entities.Closeable
	Prepare(ctx context.Context) error
	Start()
}
