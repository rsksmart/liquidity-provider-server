package watcher

import (
	"context"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	log "github.com/sirupsen/logrus"
)

type TransferExcessToColdWalletUseCase interface {
	Run(ctx context.Context) (*liquidity_provider.TransferToColdWalletResult, error)
}

type TransferColdWalletWatcher struct {
	transferUseCase    TransferExcessToColdWalletUseCase
	watcherStopChannel chan bool
	ticker             utils.Ticker
	timeout            time.Duration
}

func NewTransferColdWalletWatcher(
	transferUseCase TransferExcessToColdWalletUseCase,
	ticker utils.Ticker,
	timeout time.Duration,
) *TransferColdWalletWatcher {
	watcherStopChannel := make(chan bool, 1)
	return &TransferColdWalletWatcher{
		transferUseCase:    transferUseCase,
		watcherStopChannel: watcherStopChannel,
		ticker:             ticker,
		timeout:            timeout,
	}
}

func (watcher *TransferColdWalletWatcher) Prepare(ctx context.Context) error {
	return nil
}

func (watcher *TransferColdWalletWatcher) Start() {
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			ctx, cancel := context.WithTimeout(context.Background(), watcher.timeout)
			result, err := watcher.transferUseCase.Run(ctx)
			if err != nil {
				log.Error("TransferColdWalletWatcher: Error executing transfer to cold wallet: ", err)
			} else {
				watcher.logTransferResult(result)
			}
			cancel()
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *TransferColdWalletWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- true
	closeChannel <- true
	log.Debug("TransferColdWalletWatcher shut down")
}

func (watcher *TransferColdWalletWatcher) logTransferResult(result *liquidity_provider.TransferToColdWalletResult) {
	if result == nil {
		return
	}
	watcher.logNetworkTransferResult("BTC", result.BtcResult)
	watcher.logNetworkTransferResult("RSK", result.RskResult)
}

func (watcher *TransferColdWalletWatcher) logNetworkTransferResult(network string, result liquidity_provider.NetworkTransferResult) {
	switch result.Status {
	case liquidity_provider.TransferStatusSuccess:
		log.Infof("TransferColdWalletWatcher: %s transfer successful - TxHash: %s, Amount: %s, Fee: %s",
			network, result.TxHash, result.Amount.String(), result.Fee.String())
	case liquidity_provider.TransferStatusSkippedNoExcess:
		log.Infof("TransferColdWalletWatcher: %s transfer skipped - no excess liquidity", network)
	case liquidity_provider.TransferStatusSkippedNotEconomical:
		log.Infof("TransferColdWalletWatcher: %s transfer skipped - not economical: %s", network, result.Message)
	case liquidity_provider.TransferStatusSkippedCooldown:
		log.Infof("TransferColdWalletWatcher: %s transfer skipped - liquidity target cooldown active", network)
	case liquidity_provider.TransferStatusFailed:
		log.Errorf("TransferColdWalletWatcher: %s transfer failed - %s: %v", network, result.Message, result.Error)
	}
}
