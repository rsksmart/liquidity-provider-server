package pegin_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test/integration/regtest"
	"github.com/stretchr/testify/require"
)

const (
	// Shorter than the LPS watcher intervals of the local stack, so no pass is missed for long.
	pollInterval  = 2 * time.Second
	phaseTimeout  = 5 * time.Minute
	weiPerSatoshi = 10_000_000_000
	// The largest and smallest v2.6 peg-in quotes the local stack allows (maxValue and minValue in
	// docker-compose/local/lps/configurations/requests.js).
	maxQuoteValueWei = 10_000_000_000_000_000_000
	minQuoteValueWei = 600_000_000_000_000_000
)

// pegInStack holds what every peg-in test drives and observes. The steps on it wait for the LPS;
// the tests assert on what the steps return.
type pegInStack struct {
	btc     *regtest.Bitcoin
	rsk     *regtest.Rootstock
	flyover *regtest.Flyover
	bridge  *regtest.Bridge
	db      *regtest.LpsDB
	lps     *regtest.Lps
	driver  *regtest.Driver
	lp      common.Address
}

func newPegInStack(t *testing.T) *pegInStack {
	t.Helper()
	stack := regtest.LoadStack(t)
	stack.RequireLpsHealthy(t)
	btc := regtest.NewBitcoin(t, stack)
	rsk := regtest.NewRootstock(t, stack)
	flyover := regtest.NewFlyover(t, stack, rsk)
	require.True(t, flyover.IsOperationalPegInProvider(t, stack.LiquidityProvider), "the LP must be registered and approved")
	s := &pegInStack{
		btc:     btc,
		rsk:     rsk,
		flyover: flyover,
		bridge:  regtest.NewBridge(stack, rsk),
		db:      regtest.NewLpsDB(t, stack),
		lps:     regtest.NewLps(stack),
		driver:  regtest.NewDriver(btc, rsk, pollInterval),
		lp:      stack.LiquidityProvider,
	}
	s.waitNoPegInsInFlight(t)
	return s
}

// waitNoPegInsInFlight drives until the LPS has finished every peg-in it started, such as one
// left behind by an interrupted run. Its settlement would otherwise move the LP's balances while
// a test measures them.
func (s *pegInStack) waitNoPegInsInFlight(t *testing.T) {
	t.Helper()
	s.driver.AdvanceUntil(t, "earlier peg-ins settled", phaseTimeout, regtest.Step{BtcBlocks: 1, RskBlocks: 3},
		func(t *testing.T) (bool, string) {
			unsettled := s.db.UnsettledClaims(t)
			return len(unsettled) == 0, fmt.Sprintf("%d peg-ins still in flight", len(unsettled))
		})
}

type deposit struct {
	address   string
	txID      string
	proof     regtest.Proof
	amountWei *big.Int
}

// payDeposit pays amount to the institution's deposit address and drives the chains until the
// bridge sees the payment, which is what registration needs.
func (s *pegInStack) payDeposit(t *testing.T, institution common.Address, amount btcutil.Amount) deposit {
	t.Helper()
	address := s.flyover.PegInAddress(t, institution)
	txID := s.btc.Pay(t, address, amount)
	s.btc.Mine(t, 1)
	proof := s.btc.Proof(t, txID)
	// Powpeg relays BTC headers to the bridge only while both chains move.
	s.driver.AdvanceUntil(t, "bridge sees the deposit", phaseTimeout, regtest.Step{BtcBlocks: 1, RskBlocks: 3},
		func(t *testing.T) (bool, string) {
			confirmations := s.bridge.TxConfirmations(t, txID, proof)
			return confirmations > 0, fmt.Sprintf("bridge confirmations %d (bridge BTC height %d, bitcoind %d)",
				confirmations, s.bridge.BtcHeight(t), s.btc.Height(t))
		})
	return deposit{address: address, txID: txID, proof: proof, amountWei: satoshisToWei(amount)}
}

func (s *pegInStack) waitWatchImported(t *testing.T, institution common.Address) *rootstock.PegInWatch {
	t.Helper()
	var watch *rootstock.PegInWatch
	s.driver.AdvanceUntil(t, "watch imported", phaseTimeout, regtest.Step{RskBlocks: 1},
		func(t *testing.T) (bool, string) {
			watch = s.db.Watch(t, institution)
			if watch == nil {
				return false, "no watch row"
			}
			require.NotEqual(t, rootstock.PegInWatchUnsupportedEncoding, watch.State)
			return watch.State == rootstock.PegInWatchImported, "watch " + string(watch.State)
		})
	return watch
}

// waitClaimed drives until the LPS has fronted the deposit and settled its claim. It returns the
// claim and the PegInRequested event of this deposit.
func (s *pegInStack) waitClaimed(t *testing.T, institution common.Address, paid deposit) (*rootstock.PegInClaim, *regtest.PegInRequested) {
	t.Helper()
	var claim *rootstock.PegInClaim
	s.driver.AdvanceUntil(t, "claim claimed", phaseTimeout, regtest.Step{BtcBlocks: 1, RskBlocks: 1},
		func(t *testing.T) (bool, string) {
			claim = s.db.Claim(t, institution, paid.txID)
			if claim == nil {
				return false, fmt.Sprintf("no claim row (deposit confirmations %d)", s.btc.Confirmations(t, paid.txID))
			}
			require.NotEqual(t, rootstock.PegInClaimRaceLost, claim.State, "the only LP lost its own race")
			return claim.State == rootstock.PegInClaimClaimed, "claim " + string(claim.State)
		})
	require.NotEmpty(t, claim.PegInID, "a claimed row carries the peg-in id")
	for _, event := range s.flyover.PegInRequestedEvents(t, institution) {
		if hex.EncodeToString(event.PegInId[:]) == claim.PegInID {
			return claim, event
		}
	}
	require.FailNowf(t, "no PegInRequested event", "for peg-in %s", claim.PegInID)
	return nil, nil
}

// waitResolved drives until the LPS has resolved the peg-in. It returns the claim and the single
// PegInResolved event of this peg-in.
func (s *pegInStack) waitResolved(
	t *testing.T,
	institution common.Address,
	paid deposit,
	requested *regtest.PegInRequested,
) (*rootstock.PegInClaim, *regtest.PegInResolved) {
	t.Helper()
	var claim *rootstock.PegInClaim
	s.driver.AdvanceUntil(t, "claim resolved", phaseTimeout, regtest.Step{BtcBlocks: 1, RskBlocks: 3},
		func(t *testing.T) (bool, string) {
			claim = s.db.Claim(t, institution, paid.txID)
			require.NotEqual(t, rootstock.PegInClaimResolveFailed, claim.State)
			return claim.State == rootstock.PegInClaimResolved, fmt.Sprintf("claim %s (deposit confirmations %d)",
				claim.State, s.btc.Confirmations(t, paid.txID))
		})
	events := s.flyover.PegInResolvedEvents(t, requested.PegInId)
	require.Len(t, events, 1, "a peg-in is resolved once")
	return claim, events[0]
}

// runPhase runs one phase as a subtest and stops the test when it fails, so later phases do not
// fail on the missing results of an earlier one.
func runPhase(t *testing.T, name string, run func(t *testing.T)) {
	t.Helper()
	if !t.Run(name, run) {
		t.FailNow()
	}
}

// depositAmount is twice the highest of the two minimums a deposit must meet: the registry's,
// and the peg-in configuration's that requestPegIn enforces.
func depositAmount(t *testing.T, flyover *regtest.Flyover) btcutil.Amount {
	t.Helper()
	configuration := flyover.PegInConfiguration(t)
	minimumSats := new(big.Int).Div(configuration.MinAmount, big.NewInt(weiPerSatoshi))
	minimumSats = maxBigInt(minimumSats, flyover.MinDepositSats(t))
	amount := btcutil.Amount(2 * minimumSats.Int64())
	maxAmountSats := new(big.Int).Div(configuration.MaxAmount, big.NewInt(weiPerSatoshi))
	require.LessOrEqual(t, int64(amount), maxAmountSats.Int64(), "the deposit must stay under the peg-in maximum")
	return amount
}

func satoshisToWei(amount btcutil.Amount) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(amount)), big.NewInt(weiPerSatoshi))
}

func minBigInt(a, b *big.Int) *big.Int {
	if a.Cmp(b) < 0 {
		return a
	}
	return b
}

func maxBigInt(a, b *big.Int) *big.Int {
	if a.Cmp(b) > 0 {
		return a
	}
	return b
}
