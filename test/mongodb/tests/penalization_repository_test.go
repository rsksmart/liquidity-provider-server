//go:build integration

package mongodb_test

import (
	"context"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/test/mongodb/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPenalization_InsertAndGetByHashes(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	hash1 := utils.RandomHash()
	hash2 := utils.RandomHash()
	hash3 := utils.RandomHash()
	event1 := utils.NewTestPenalizedEvent(hash1)
	event2 := utils.NewTestPenalizedEvent(hash2)
	event3 := utils.NewTestPenalizedEvent(hash3)

	require.NoError(t, penaltyRepo.InsertPenalization(ctx, event1))
	require.NoError(t, penaltyRepo.InsertPenalization(ctx, event2))
	require.NoError(t, penaltyRepo.InsertPenalization(ctx, event3))

	results, err := penaltyRepo.GetPenalizationsByQuoteHashes(ctx, []string{hash1, hash3})
	require.NoError(t, err)
	assert.Len(t, results, 2)

	resultHashes := make([]string, 0, len(results))
	for _, r := range results {
		resultHashes = append(resultHashes, r.QuoteHash)
	}
	assert.ElementsMatch(t, []string{hash1, hash3}, resultHashes)
}

func TestPenalization_GetByHashes_NoMatches(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	results, err := penaltyRepo.GetPenalizationsByQuoteHashes(ctx, []string{utils.RandomHash()})
	require.NoError(t, err)
	assert.Empty(t, results)
}
