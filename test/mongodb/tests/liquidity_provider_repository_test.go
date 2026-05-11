//go:build integration

package mongodb_test

import (
	"context"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/test/mongodb/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLP_UpsertAndGetPeginConfig(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	config := utils.NewTestPeginConfig()
	err := lpRepo.UpsertPeginConfiguration(ctx, config)
	require.NoError(t, err)

	got, err := lpRepo.GetPeginConfiguration(ctx)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, config.Signature, got.Signature)
	assert.Equal(t, config.Hash, got.Hash)
	assert.Equal(t, config.Value.TimeForDeposit, got.Value.TimeForDeposit)
	assert.Equal(t, config.Value.CallTime, got.Value.CallTime)
	assertWeiEqual(t, config.Value.PenaltyFee, got.Value.PenaltyFee)
	assertWeiEqual(t, config.Value.FixedFee, got.Value.FixedFee)
	assertWeiEqual(t, config.Value.MaxValue, got.Value.MaxValue)
	assertWeiEqual(t, config.Value.MinValue, got.Value.MinValue)

	// Upsert overwrites
	config.Value.CallTime = 9999
	err = lpRepo.UpsertPeginConfiguration(ctx, config)
	require.NoError(t, err)

	got, err = lpRepo.GetPeginConfiguration(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint32(9999), got.Value.CallTime)
}

func TestLP_Get_NotFound(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	tests := []struct {
		name string
		get  func(context.Context) (any, error)
	}{
		{
			name: "pegin_config",
			get: func(ctx context.Context) (any, error) {
				return lpRepo.GetPeginConfiguration(ctx)
			},
		},
		{
			name: "pegout_config",
			get: func(ctx context.Context) (any, error) {
				return lpRepo.GetPegoutConfiguration(ctx)
			},
		},
		{
			name: "general_config",
			get: func(ctx context.Context) (any, error) {
				return lpRepo.GetGeneralConfiguration(ctx)
			},
		},
		{
			name: "credentials",
			get: func(ctx context.Context) (any, error) {
				return lpRepo.GetCredentials(ctx)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.get(ctx)
			require.NoError(t, err)
			assert.Nil(t, got)
		})
	}
}

func TestLP_UpsertAndGetPegoutConfig(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	config := utils.NewTestPegoutConfig()
	err := lpRepo.UpsertPegoutConfiguration(ctx, config)
	require.NoError(t, err)

	got, err := lpRepo.GetPegoutConfiguration(ctx)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, config.Signature, got.Signature)
	assert.Equal(t, config.Value.TimeForDeposit, got.Value.TimeForDeposit)
	assert.Equal(t, config.Value.ExpireTime, got.Value.ExpireTime)
	assert.Equal(t, config.Value.ExpireBlocks, got.Value.ExpireBlocks)
	assertWeiEqual(t, config.Value.PenaltyFee, got.Value.PenaltyFee)
	assertWeiEqual(t, config.Value.BridgeTransactionMin, got.Value.BridgeTransactionMin)
}

func TestLP_UpsertAndGetGeneralConfig(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	config := utils.NewTestGeneralConfig()
	err := lpRepo.UpsertGeneralConfiguration(ctx, config)
	require.NoError(t, err)

	got, err := lpRepo.GetGeneralConfiguration(ctx)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, config.Signature, got.Signature)
	assert.Equal(t, config.Value.PublicLiquidityCheck, got.Value.PublicLiquidityCheck)
	assert.Equal(t, config.Value.RskConfirmations, got.Value.RskConfirmations)
	assert.Equal(t, config.Value.BtcConfirmations, got.Value.BtcConfirmations)
}

func TestLP_UpsertAndGetCredentials(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	creds := utils.NewTestCredentials()
	err := lpRepo.UpsertCredentials(ctx, creds)
	require.NoError(t, err)

	got, err := lpRepo.GetCredentials(ctx)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, creds.Value.HashedUsername, got.Value.HashedUsername)
	assert.Equal(t, creds.Value.HashedPassword, got.Value.HashedPassword)
	assert.Equal(t, creds.Value.UsernameSalt, got.Value.UsernameSalt)
	assert.Equal(t, creds.Value.PasswordSalt, got.Value.PasswordSalt)
}
