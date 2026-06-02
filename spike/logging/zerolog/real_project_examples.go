// Real LPS examples: imports project packages and rewrites representative
// pegout-rebalance logs into structured zerolog fields.
package main

import (
	"context"
	"errors"
	"strconv"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

func demoRealLPSUseCaseLogging(ctx context.Context) {
	adjustedTotal := entities.SatoshiToWei(75_000_000)
	bridgeMin := entities.EtherToWei(1)

	loggerWithTrace(ctx).Info().
		Str("useCase", string(usecases.BridgePegoutId)).
		Str("adjustedTotalWei", adjustedTotal.String()).
		Str("adjustedTotalRbtc", adjustedTotal.ToRbtc().Text('f', 18)).
		Str("bridgeMinWei", bridgeMin.String()).
		Str("bridgeMinRbtc", bridgeMin.ToRbtc().Text('f', 18)).
		Msg("bridge pegout threshold not met")

	loggerWithTrace(ctx).Info().
		Str("useCase", string(usecases.BridgePegoutId)).
		Int("chunkIndex", 1).
		Int("totalChunks", 3).
		Str("transactionHash", "0xfeedbeef").
		Msg("bridge pegout split transaction persisted")
}

// demoRealLPSUseCaseError shows the LPS sentinel + WrapUseCaseError pattern
func demoRealLPSUseCaseError(ctx context.Context) {
	adjustedTotal := entities.SatoshiToWei(50_000_000)
	bridgeMin := entities.EtherToWei(1)

	wrapped := usecases.WrapUseCaseError(usecases.BridgePegoutId, usecases.TxBelowMinimumError)

	logError(ctx, "bridge pegout threshold not met", string(usecases.BridgePegoutId), wrapped, map[string]string{
		"adjustedTotalRbtc": adjustedTotal.ToRbtc().Text('f', 18),
		"bridgeMinRbtc":     bridgeMin.ToRbtc().Text('f', 18),
		"matchesSentinel":   strconv.FormatBool(errors.Is(wrapped, usecases.TxBelowMinimumError)),
	})
}
