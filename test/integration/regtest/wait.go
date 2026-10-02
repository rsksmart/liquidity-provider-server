package regtest

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/stretchr/testify/require"
)

const (
	PollInterval      = 500 * time.Millisecond
	WaitTimeout       = 20 * time.Second
	OneClaimTick      = 5 * time.Second
	FinalityDepth     = 2
	RelayWaitTimeout  = 3 * time.Minute
	RelayPollInterval = 2 * time.Second
	RelayMinePerPoll  = 3
)

func (stack *MongoStack) WaitNoWatch(t *testing.T, user common.Address) {
	t.Helper()
	require.Never(t, func() bool {
		return stack.GetWatch(t, user) != nil
	}, WaitTimeout, PollInterval, "peginWatch appeared for unregistered %s", user.Hex())
}

func (stack *MongoStack) WaitWatchImported(t *testing.T, user common.Address) *rootstock.PegInWatch {
	t.Helper()
	var watch *rootstock.PegInWatch
	require.Eventually(t, func() bool {
		watch = stack.GetWatch(t, user)
		return watch != nil && watch.State == rootstock.PegInWatchImported
	}, WaitTimeout, PollInterval, "peginWatch did not become imported for %s", user.Hex())
	return watch
}

func (stack *MongoStack) WaitClaimStates(
	t *testing.T,
	user common.Address,
	depositTxID string,
	states ...rootstock.PegInClaimState,
) *rootstock.PegInClaim {
	t.Helper()
	allowed := map[rootstock.PegInClaimState]struct{}{}
	for _, state := range states {
		allowed[state] = struct{}{}
	}
	var claim *rootstock.PegInClaim
	require.Eventually(t, func() bool {
		claim = stack.GetClaim(t, user, depositTxID)
		if claim == nil {
			return false
		}
		_, ok := allowed[claim.State]
		return ok
	}, WaitTimeout, PollInterval, "peginClaims row for %s / %s did not reach %v", user.Hex(), depositTxID, states)
	return claim
}

func (stack *MongoStack) WaitClaimed(t *testing.T, user common.Address, depositTxID string) *rootstock.PegInClaim {
	t.Helper()
	claim := stack.WaitClaimStates(t, user, depositTxID, rootstock.PegInClaimClaimed)
	require.NotEmpty(t, claim.TxHash, "claimed peginClaims row for %s / %s has empty tx_hash", user.Hex(), depositTxID)
	return claim
}

func (stack *RootstockStack) WaitBalanceIncreasedBy(t *testing.T, user common.Address, before *big.Int, net *entities.Wei) {
	t.Helper()
	expected := new(big.Int).Add(before, net.AsBigInt())
	require.Eventually(t, func() bool {
		return stack.Balance(t, user).Cmp(expected) == 0
	}, WaitTimeout, PollInterval, "user %s balance did not increase by net %s", user.Hex(), net.AsBigInt())
}

func (stack *MongoStack) WaitNoPaidClaim(t *testing.T, user common.Address, depositTxID string) {
	t.Helper()
	require.Never(t, func() bool {
		claim := stack.GetClaim(t, user, depositTxID)
		if claim == nil {
			return false
		}
		return claim.State == rootstock.PegInClaimSubmitting || claim.State == rootstock.PegInClaimClaimed
	}, OneClaimTick, PollInterval, "LPS broadcast a claim for %s / %s", user.Hex(), depositTxID)
}

func (stack *RootstockStack) WaitBalanceUnchanged(t *testing.T, user common.Address, before *big.Int) {
	t.Helper()
	require.Never(t, func() bool {
		return stack.Balance(t, user).Cmp(before) != 0
	}, OneClaimTick, PollInterval, "user %s balance changed", user.Hex())
}

// WaitBridgeBtcAtLeast follows fed-migrator primeBtcRelay: mine Rootstock
// blocks so federator transactions can be included, then read the Bridge
// Bitcoin height until it reaches minHeight.
func (stack *RootstockStack) WaitBridgeBtcAtLeast(t *testing.T, minHeight int64) {
	t.Helper()
	require.Eventually(t, func() bool {
		stack.Advance(t, RelayMinePerPoll)
		return stack.BridgeBtcHeight(t) >= minHeight
	}, RelayWaitTimeout, RelayPollInterval, "Bridge Bitcoin height did not reach %d", minHeight)
}
