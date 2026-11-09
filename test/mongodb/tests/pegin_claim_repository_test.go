//go:build integration

package mongodb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	mongoAdapter "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func newIntegrationPegInClaim(ctx context.Context, t *testing.T, state rootstock.PegInClaimState) rootstock.PegInClaim {
	t.Helper()
	now := time.Now().UTC().Truncate(time.Millisecond)
	claim := rootstock.PegInClaim{
		RskAddress:  fmt.Sprintf("0xclaimrsk-%d", time.Now().UnixNano()),
		DepositTxID: fmt.Sprintf("claimdeposit-%d", time.Now().UnixNano()),
		BtcAddress:  "2N4bCrVMPu6ixHvxRS7qXcpRkdFh6WW4wEz",
		PegInID:     "0x6c0c1c2ba8a9cd4f0a3f0e15cc1d0a1b2c3d4e5f60718293a4b5c6d7e8f90a1b",
		State:       state,
		Amount:      entities.NewWei(5000000000000000),
		Fee:         entities.NewWei(100000000000000),
		GasFee:      entities.NewWei(20000000000000),
		Value:       entities.NewWei(4880000000000000),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	collection := mongoClient.Database(mongoAdapter.DbName).Collection(mongoAdapter.PegInClaimCollection)
	t.Cleanup(func() {
		_, err := collection.DeleteMany(ctx, bson.M{"rsk_address": claim.RskAddress})
		assert.NoError(t, err)
	})
	return claim
}

func TestPegInClaimRepository_AddKeepsAClaimThatWasNotARetryableFailure(t *testing.T) {
	ctx := context.Background()
	repository := mongoAdapter.NewPegInClaimMongoRepository(conn)
	collection := mongoClient.Database(mongoAdapter.DbName).Collection(mongoAdapter.PegInClaimCollection)

	for _, state := range []rootstock.PegInClaimState{
		rootstock.PegInClaimSubmitting,
		rootstock.PegInClaimClaimed,
		rootstock.PegInClaimRaceLost,
	} {
		t.Run(string(state), func(t *testing.T) {
			stored := newIntegrationPegInClaim(ctx, t, state)
			stored.TxHash = "0xsent"
			require.NoError(t, repository.Add(ctx, stored))

			second := stored
			second.State = rootstock.PegInClaimSubmitting
			second.TxHash = ""
			require.ErrorIs(t, repository.Add(ctx, second), rootstock.PegInClaimExistsError)

			result, err := repository.Get(ctx, stored.RskAddress, stored.DepositTxID)
			require.NoError(t, err)
			assert.Equal(t, &stored, result)
			count, err := collection.CountDocuments(ctx, bson.M{"rsk_address": stored.RskAddress})
			require.NoError(t, err)
			assert.Equal(t, int64(1), count)
		})
	}
}

func TestPegInClaimRepository_AddReplacesARetryableFailure(t *testing.T) {
	ctx := context.Background()
	repository := mongoAdapter.NewPegInClaimMongoRepository(conn)
	collection := mongoClient.Database(mongoAdapter.DbName).Collection(mongoAdapter.PegInClaimCollection)

	failed := newIntegrationPegInClaim(ctx, t, rootstock.PegInClaimRetryableFailure)
	require.NoError(t, repository.Add(ctx, failed))

	retry := failed
	retry.State = rootstock.PegInClaimSubmitting
	retry.Value = entities.NewWei(4870000000000000)
	retry.UpdatedAt = failed.UpdatedAt.Add(time.Minute)
	require.NoError(t, repository.Add(ctx, retry))

	result, err := repository.Get(ctx, failed.RskAddress, failed.DepositTxID)
	require.NoError(t, err)
	assert.Equal(t, &retry, result)
	count, err := collection.CountDocuments(ctx, bson.M{"rsk_address": failed.RskAddress})
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestPegInClaimRepository_UpdateReplacesTheStoredClaim(t *testing.T) {
	ctx := context.Background()
	repository := mongoAdapter.NewPegInClaimMongoRepository(conn)

	claim := newIntegrationPegInClaim(ctx, t, rootstock.PegInClaimSubmitting)
	require.NoError(t, repository.Add(ctx, claim))

	claimed := claim
	claimed.State = rootstock.PegInClaimClaimed
	claimed.TxHash = "0xclaimtx"
	claimed.CallSuccess = true
	claimed.UpdatedAt = claim.UpdatedAt.Add(time.Minute)
	require.NoError(t, repository.Update(ctx, claimed))

	result, err := repository.Get(ctx, claim.RskAddress, claim.DepositTxID)
	require.NoError(t, err)
	assert.Equal(t, &claimed, result)
}

func TestPegInClaimRepository_UpdateFailsForAnUnknownClaim(t *testing.T) {
	ctx := context.Background()
	repository := mongoAdapter.NewPegInClaimMongoRepository(conn)
	claim := newIntegrationPegInClaim(ctx, t, rootstock.PegInClaimClaimed)
	require.ErrorIs(t, repository.Update(ctx, claim), rootstock.PegInClaimNotFoundError)
}

func TestPegInClaimRepository_GetReturnsNilForAnUnknownClaim(t *testing.T) {
	repository := mongoAdapter.NewPegInClaimMongoRepository(conn)
	result, err := repository.Get(context.Background(), "0xunknownclaimrsk", "unknown-deposit")
	require.NoError(t, err)
	assert.Nil(t, result)
}
