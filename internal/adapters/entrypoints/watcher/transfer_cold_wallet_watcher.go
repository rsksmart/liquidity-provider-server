package watcher

import (
	"context"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	log "github.com/sirupsen/logrus"
)

const transferColdWalletTimeout = 60 * time.Second

type TransferExcessToColdWalletUseCase interface {
	Run(ctx context.Context) (*liquidity_provider.TransferToColdWalletResult, error)
}

type TransferColdWalletWatcher struct {
	transferUseCase    TransferExcessToColdWalletUseCase
	watcherStopChannel chan bool
	ticker             utils.Ticker
}

func NewTransferColdWalletWatcher(
	transferUseCase TransferExcessToColdWalletUseCase,
	ticker utils.Ticker,
) *TransferColdWalletWatcher {
	watcherStopChannel := make(chan bool, 1)
	return &TransferColdWalletWatcher{
		transferUseCase:    transferUseCase,
		watcherStopChannel: watcherStopChannel,
		ticker:             ticker,
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
			ctx, cancel := context.WithTimeout(context.Background(), transferColdWalletTimeout)
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

	// Log BTC transfer result
	switch result.BtcResult.Status {
	case liquidity_provider.TransferStatusSuccess:
		log.Infof("TransferColdWalletWatcher: BTC transfer successful - TxHash: %s, Amount: %s, Fee: %s",
			result.BtcResult.TxHash,
			result.BtcResult.Amount.String(),
			result.BtcResult.Fee.String())
	case liquidity_provider.TransferStatusSkippedNoExcess:
		log.Info("TransferColdWalletWatcher: BTC transfer skipped - no excess liquidity")
	case liquidity_provider.TransferStatusSkippedNotEconomical:
		log.Info("TransferColdWalletWatcher: BTC transfer skipped - not economical: ", result.BtcResult.Message)
	case liquidity_provider.TransferStatusFailed:
		log.Errorf("TransferColdWalletWatcher: BTC transfer failed - %s: %v",
			result.BtcResult.Message,
			result.BtcResult.Error)
	}

	// Log RSK transfer result
	switch result.RskResult.Status {
	case liquidity_provider.TransferStatusSuccess:
		log.Infof("TransferColdWalletWatcher: RSK transfer successful - TxHash: %s, Amount: %s, Fee: %s",
			result.RskResult.TxHash,
			result.RskResult.Amount.String(),
			result.RskResult.Fee.String())
	case liquidity_provider.TransferStatusSkippedNoExcess:
		log.Info("TransferColdWalletWatcher: RSK transfer skipped - no excess liquidity")
	case liquidity_provider.TransferStatusSkippedNotEconomical:
		log.Info("TransferColdWalletWatcher: RSK transfer skipped - not economical: ", result.RskResult.Message)
	case liquidity_provider.TransferStatusFailed:
		log.Errorf("TransferColdWalletWatcher: RSK transfer failed - %s: %v",
			result.RskResult.Message,
			result.RskResult.Error)
	}
}
