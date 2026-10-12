package regtest_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/rsksmart/liquidity-provider-server/test/integration/regtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStackIsUp checks the helpers against the running stack, so a broken helper shows up here
// rather than in the middle of a peg-in.
//
//nolint:funlen // one function checks every helper against the same running stack
func TestStackIsUp(t *testing.T) {
	stack := regtest.LoadStack(t)
	stack.RequireLpsHealthy(t)
	btc := regtest.NewBitcoin(t, stack)
	rsk := regtest.NewRootstock(t, stack)

	t.Run("mines bitcoin blocks", func(t *testing.T) {
		before := btc.Height(t)
		btc.Mine(t, 2)
		assert.Equal(t, before+2, btc.Height(t))
	})

	t.Run("mines rootstock blocks", func(t *testing.T) {
		before := rsk.BlockNumber(t)
		rsk.Mine(t, 2)
		assert.Equal(t, before+2, rsk.BlockNumber(t))
	})

	t.Run("pays and proves a bitcoin transaction", func(t *testing.T) {
		txID := btc.Pay(t, btc.NewAddress(t).EncodeAddress(), btcutil.Amount(100_000))
		assert.Zero(t, btc.Confirmations(t, txID))
		btc.Mine(t, 1)
		assert.Equal(t, uint64(1), btc.Confirmations(t, txID))

		proof := btc.Proof(t, txID)
		assert.Equal(t, txID, chainhash.DoubleHashH(proof.RawTx).String(), "the raw transaction hashes to its txid")
		hash, err := chainhash.NewHashFromStr(txID)
		require.NoError(t, err)
		mined, err := btc.Node().GetRawTransactionVerbose(hash)
		require.NoError(t, err)
		assert.Equal(t, mined.BlockHash, hex.EncodeToString(proof.BlockHash[:]), "the block hash is in display order")
		assert.NotNil(t, proof.Branch.Path)
	})

	t.Run("signs and sends from a new account", func(t *testing.T) {
		funding := big.NewInt(1_000_000_000_000_000)
		account := rsk.NewAccount(t, funding)
		assert.Equal(t, funding, rsk.Balance(t, account.Address))

		receipt := rsk.Transact(t, account, account.Address, nil, nil)
		assert.NotNil(t, receipt.BlockNumber)
	})

	t.Run("reads the flyover contracts", func(t *testing.T) {
		flyover := regtest.NewFlyover(t, stack, rsk)
		assert.True(t, flyover.IsOperationalPegInProvider(t, stack.LiquidityProvider), "the LP is registered and approved")
		assert.Positive(t, flyover.MinDepositSats(t).Sign())
		configuration := flyover.PegInConfiguration(t)
		assert.Positive(t, configuration.MinAmount.Sign())
		confirmations := flyover.RequiredPegInConfirmations(t, configuration.MinAmount)
		assert.Positive(t, confirmations)
		t.Logf("peg-in: min %s wei (registry min %s sats), fixed fee %s, percentage fee %s, registrant fee %s, %d confirmations at the minimum",
			configuration.MinAmount, flyover.MinDepositSats(t), configuration.FixedFee, configuration.PercentageFee,
			configuration.RegistrantFee, confirmations)
		assert.NotNil(t, flyover.Balance(t, stack.LiquidityProvider))

		user := rsk.NewAccount(t, new(big.Int)).Address
		depositAddress := flyover.PegInAddress(t, user)
		decoded, err := btcutil.DecodeAddress(depositAddress, btc.Params(t))
		require.NoError(t, err)
		assert.True(t, decoded.IsForNet(btc.Params(t)))
		assert.Equal(t, depositAddress, flyover.PegInAddress(t, user), "the derivation is deterministic")
		assert.False(t, flyover.IsRegistered(t, user))
		assert.Empty(t, flyover.AddressRegisteredEvents(t, user))
		assert.Empty(t, flyover.PegInRequestedEvents(t, user))
	})

	t.Run("reads the LPS database", func(t *testing.T) {
		db := regtest.NewLpsDB(t, stack)
		user := rsk.NewAccount(t, new(big.Int)).Address
		assert.Nil(t, db.Watch(t, user))
		assert.Nil(t, db.Claim(t, user, "0000000000000000000000000000000000000000000000000000000000000000"))
	})

	t.Run("the bridge learns bitcoin blocks through powpeg", func(t *testing.T) {
		bridge := regtest.NewBridge(stack, rsk)
		driver := regtest.NewDriver(btc, rsk, 2*time.Second)
		txID := btc.Pay(t, btc.NewAddress(t).EncodeAddress(), btcutil.Amount(100_000))
		btc.Mine(t, 1)
		proof := btc.Proof(t, txID)
		// Powpeg relays headers when it sees new BTC blocks and sends them in RSK transactions,
		// so both chains have to move.
		driver.AdvanceUntil(t, "bridge relay (needs powpeg)", 3*time.Minute, regtest.Step{BtcBlocks: 1, RskBlocks: 3},
			func(t *testing.T) (bool, string) {
				confirmations := bridge.TxConfirmations(t, txID, proof)
				return confirmations > 0, fmt.Sprintf("bridge BTC height %d, bitcoind %d, tx confirmations seen by the bridge %d",
					bridge.BtcHeight(t), btc.Height(t), confirmations)
			})
	})

	t.Run("driver advances until the condition holds", func(t *testing.T) {
		driver := regtest.NewDriver(btc, rsk, 100*time.Millisecond)
		target := btc.Height(t) + 3
		driver.AdvanceUntil(t, "three bitcoin blocks", 30*time.Second, regtest.Step{BtcBlocks: 1, RskBlocks: 1},
			func(t *testing.T) (bool, string) {
				height := btc.Height(t)
				return height >= target, "bitcoin height " + big.NewInt(height).String()
			})
	})
}
