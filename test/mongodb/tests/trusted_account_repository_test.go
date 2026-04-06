//go:build integration

package mongodb_test

import (
	"context"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mongodb/support"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrusted_AddAndGet(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	account := support.NewTestTrustedAccount("0xaaaa1111bbbb2222cccc3333dddd4444eeee5555")
	err := trustedRepo.AddTrustedAccount(ctx, account)
	require.NoError(t, err)

	got, err := trustedRepo.GetTrustedAccount(ctx, account.Value.Address)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, account.Value.Address, got.Value.Address)
	assert.Equal(t, account.Value.Name, got.Value.Name)
	assert.Equal(t, account.Signature, got.Signature)
	assertWeiEqual(t, account.Value.BtcLockingCap, got.Value.BtcLockingCap)
	assertWeiEqual(t, account.Value.RbtcLockingCap, got.Value.RbtcLockingCap)
}

func TestTrusted_Get_NotFound(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	_, err := trustedRepo.GetTrustedAccount(ctx, "0xnonexistent")
	require.ErrorIs(t, err, liquidity_provider.TrustedAccountNotFoundError)
}

func TestTrusted_AddDuplicate(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	account := support.NewTestTrustedAccount("0xdddd1111eeee2222ffff3333aaaa4444bbbb5555")
	err := trustedRepo.AddTrustedAccount(ctx, account)
	require.NoError(t, err)

	err = trustedRepo.AddTrustedAccount(ctx, account)
	require.ErrorIs(t, err, liquidity_provider.DuplicateTrustedAccountError)
}

func TestTrusted_GetAll(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	addr1 := "0x1111222233334444555566667777888899990001"
	addr2 := "0x1111222233334444555566667777888899990002"
	require.NoError(t, trustedRepo.AddTrustedAccount(ctx, support.NewTestTrustedAccount(addr1)))
	require.NoError(t, trustedRepo.AddTrustedAccount(ctx, support.NewTestTrustedAccount(addr2)))

	all, err := trustedRepo.GetAllTrustedAccounts(ctx)
	require.NoError(t, err)
	assert.Len(t, all, 2)
	gotAddresses := make([]string, 0, len(all))
	for _, account := range all {
		gotAddresses = append(gotAddresses, account.Value.Address)
	}
	assert.ElementsMatch(t, []string{addr1, addr2}, gotAddresses)
}

func TestTrusted_Update(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	account := support.NewTestTrustedAccount("0xeeee1111ffff2222aaaa3333bbbb4444cccc5555")
	require.NoError(t, trustedRepo.AddTrustedAccount(ctx, account))

	account.Value.Name = "UpdatedName"
	account.Signature = "updated-sig"
	err := trustedRepo.UpdateTrustedAccount(ctx, account)
	require.NoError(t, err)

	got, err := trustedRepo.GetTrustedAccount(ctx, account.Value.Address)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, account.Value.Address, got.Value.Address)
	assert.Equal(t, "UpdatedName", got.Value.Name)
	assertWeiEqual(t, account.Value.BtcLockingCap, got.Value.BtcLockingCap)
	assertWeiEqual(t, account.Value.RbtcLockingCap, got.Value.RbtcLockingCap)
	assert.Equal(t, "updated-sig", got.Signature)
}

func TestTrusted_Update_NotFound(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	account := support.NewTestTrustedAccount("0xnonexistent000000000000000000000000000000")
	err := trustedRepo.UpdateTrustedAccount(ctx, account)
	require.ErrorIs(t, err, liquidity_provider.TrustedAccountNotFoundError)
}

func TestTrusted_Delete(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	account := support.NewTestTrustedAccount("0xffff1111aaaa2222bbbb3333cccc4444dddd5555")
	require.NoError(t, trustedRepo.AddTrustedAccount(ctx, account))

	err := trustedRepo.DeleteTrustedAccount(ctx, account.Value.Address)
	require.NoError(t, err)

	_, err = trustedRepo.GetTrustedAccount(ctx, account.Value.Address)
	require.ErrorIs(t, err, liquidity_provider.TrustedAccountNotFoundError)
}

func TestTrusted_Delete_NotFound(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	err := trustedRepo.DeleteTrustedAccount(ctx, "0xnonexistent000000000000000000000000000000")
	require.ErrorIs(t, err, liquidity_provider.TrustedAccountNotFoundError)
}
