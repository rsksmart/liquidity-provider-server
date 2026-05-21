// Real LPS examples: imports project packages and rewrites representative
// pegout-rebalance logs into structured logrus fields.
package main

import (
	"context"
	"errors"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/sirupsen/logrus"
)

func demoRealLPSUseCaseLogging(ctx context.Context) {
	adjustedTotal := entities.SatoshiToWei(75_000_000)
	bridgeMin := entities.EtherToWei(1)

	baseEntry(ctx).WithFields(logrus.Fields{
		"useCase":           string(usecases.BridgePegoutId),
		"adjustedTotalWei":  adjustedTotal.String(),
		"adjustedTotalRbtc": adjustedTotal.ToRbtc().Text('f', 18),
		"bridgeMinWei":      bridgeMin.String(),
		"bridgeMinRbtc":     bridgeMin.ToRbtc().Text('f', 18),
	}).Info("bridge pegout threshold not met")

	baseEntry(ctx).WithFields(logrus.Fields{
		"useCase":         string(usecases.BridgePegoutId),
		"chunkIndex":      1,
		"totalChunks":     3,
		"transactionHash": "0xfeedbeef",
	}).Info("bridge pegout split transaction persisted")
}

// demoRealLPSUseCaseError shows the LPS sentinel + WrapUseCaseError pattern
func demoRealLPSUseCaseError(ctx context.Context) {
	adjustedTotal := entities.SatoshiToWei(50_000_000)
	bridgeMin := entities.EtherToWei(1)

	wrapped := usecases.WrapUseCaseError(usecases.BridgePegoutId, usecases.TxBelowMinimumError)

	logError(ctx, "bridge pegout threshold not met", string(usecases.BridgePegoutId), wrapped, logrus.Fields{
		"adjustedTotalRbtc": adjustedTotal.ToRbtc().Text('f', 18),
		"bridgeMinRbtc":     bridgeMin.ToRbtc().Text('f', 18),
		"matchesSentinel":   errors.Is(wrapped, usecases.TxBelowMinimumError),
	})
}
