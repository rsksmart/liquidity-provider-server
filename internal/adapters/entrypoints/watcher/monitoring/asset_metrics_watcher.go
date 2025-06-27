package monitoring

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	log "github.com/sirupsen/logrus"
)

// GetAssetReportUseCase defines the interface for getting asset reports
type GetAssetReportUseCase interface {
	Run(ctx context.Context) (reports.GetAssetReportResult, error)
}

type AssetReportWatcher struct {
	appMetrics            *Metrics
	getAssetReportUseCase GetAssetReportUseCase
	ticker                watcher.Ticker
	watcherStopChannel    chan bool
}

func NewAssetReportWatcher(
	appMetrics *Metrics,
	getAssetReportUseCase GetAssetReportUseCase,
	ticker watcher.Ticker,
) *AssetReportWatcher {
	watcherStopChannel := make(chan bool, 1)
	return &AssetReportWatcher{
		appMetrics:            appMetrics,
		getAssetReportUseCase: getAssetReportUseCase,
		watcherStopChannel:    watcherStopChannel,
		ticker:                ticker,
	}
}

func (watcher *AssetReportWatcher) Prepare(ctx context.Context) error {
	return nil
}

func (watcher *AssetReportWatcher) Start() {
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			watcher.updateMetrics()
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *AssetReportWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- true
	closeChannel <- true
	log.Debug("AssetReportWatcher shut down")
}

func (watcher *AssetReportWatcher) updateMetrics() {
	report, err := watcher.getAssetReportUseCase.Run(context.Background())
	if err != nil {
		log.WithError(err).Error("Failed to collect asset report")
		return
	}

	watcher.appMetrics.UpdateAssetsFromReport(report)
}
