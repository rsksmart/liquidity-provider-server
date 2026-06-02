// Real LPS examples: imports project packages and rewrites representative
// pegout-rebalance logs into structured zap fields.
package main

import (
	"context"
	"errors"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"go.uber.org/zap"
)

func demoRealLPSUseCaseLogging(ctx context.Context) {
	adjustedTotal := entities.SatoshiToWei(75_000_000)
	bridgeMin := entities.EtherToWei(1)

	loggerWithTrace(ctx).Info("bridge pegout threshold not met",
		zap.String("useCase", string(usecases.BridgePegoutId)),
		zap.String("adjustedTotalWei", adjustedTotal.String()),
		zap.String("adjustedTotalRbtc", adjustedTotal.ToRbtc().Text('f', 18)),
		zap.String("bridgeMinWei", bridgeMin.String()),
		zap.String("bridgeMinRbtc", bridgeMin.ToRbtc().Text('f', 18)),
	)

	loggerWithTrace(ctx).Info("bridge pegout split transaction persisted",
		zap.String("useCase", string(usecases.BridgePegoutId)),
		zap.Int("chunkIndex", 1),
		zap.Int("totalChunks", 3),
		zap.String("transactionHash", "0xfeedbeef"),
	)
}

// demoRealLPSUseCaseError shows the LPS sentinel + WrapUseCaseError pattern
func demoRealLPSUseCaseError(ctx context.Context) {
	adjustedTotal := entities.SatoshiToWei(50_000_000)
	bridgeMin := entities.EtherToWei(1)

	wrapped := usecases.WrapUseCaseError(usecases.BridgePegoutId, usecases.TxBelowMinimumError)

	logError(ctx, "bridge pegout threshold not met", string(usecases.BridgePegoutId), wrapped,
		zap.String("adjustedTotalRbtc", adjustedTotal.ToRbtc().Text('f', 18)),
		zap.String("bridgeMinRbtc", bridgeMin.ToRbtc().Text('f', 18)),
		zap.Bool("matchesSentinel", errors.Is(wrapped, usecases.TxBelowMinimumError)),
	)
}
