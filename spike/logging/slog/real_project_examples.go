// Real LPS examples: imports project packages and rewrites representative
// pegout-rebalance logs into structured slog fields.
package main

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

func demoRealLPSUseCaseLogging(ctx context.Context) {
	adjustedTotal := entities.SatoshiToWei(75_000_000)
	bridgeMin := entities.EtherToWei(1)

	slog.InfoContext(ctx, "bridge pegout threshold not met",
		slog.String("useCase", string(usecases.BridgePegoutId)),
		slog.String("adjustedTotalWei", adjustedTotal.String()),
		slog.String("adjustedTotalRbtc", adjustedTotal.ToRbtc().Text('f', 18)),
		slog.String("bridgeMinWei", bridgeMin.String()),
		slog.String("bridgeMinRbtc", bridgeMin.ToRbtc().Text('f', 18)),
	)

	slog.InfoContext(ctx, "bridge pegout split transaction persisted",
		slog.String("useCase", string(usecases.BridgePegoutId)),
		slog.Int("chunkIndex", 1),
		slog.Int("totalChunks", 3),
		slog.String("transactionHash", "0xfeedbeef"),
	)
}

// demoRealLPSUseCaseError shows the LPS sentinel + WrapUseCaseError pattern
func demoRealLPSUseCaseError(ctx context.Context) {
	adjustedTotal := entities.SatoshiToWei(50_000_000)
	bridgeMin := entities.EtherToWei(1)

	wrapped := usecases.WrapUseCaseError(usecases.BridgePegoutId, usecases.TxBelowMinimumError)

	logError(ctx, "bridge pegout threshold not met", string(usecases.BridgePegoutId), wrapped,
		slog.String("adjustedTotalRbtc", adjustedTotal.ToRbtc().Text('f', 18)),
		slog.String("bridgeMinRbtc", bridgeMin.ToRbtc().Text('f', 18)),
		slog.Bool("matchesSentinel", errors.Is(wrapped, usecases.TxBelowMinimumError)),
	)
}
