package regtest

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	// unlockedFunder is the prefunded regtest account rskj keeps unlocked. The node signs for it,
	// and rskj's autoMine mines a block for every transaction it receives.
	unlockedFunder = "0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826"
	receiptTimeout = 30 * time.Second
	receiptPolling = 100 * time.Millisecond
	// rskj mines a block for each transaction it receives and answers this when that block does
	// not include it, which happens when another sender's transaction races it (the LPS sends raw
	// transactions concurrently). The transaction is still in the pool and goes into a later block.
	minedInAnotherBlock = "transaction wasn't mined"
	// How many empty blocks a funding transfer that raced may take to land.
	racedTransferBlocks = 5
)

// Account is a test-owned Rootstock account that signs its own transactions.
type Account struct {
	Key     *ecdsa.PrivateKey
	Address common.Address
}

// Rootstock acts as the users and the miner on the regtest rskj.
type Rootstock struct {
	client  *ethclient.Client
	chainID *big.Int
	funder  common.Address
}

// NewRootstock connects to rskj. The client closes when the test ends.
func NewRootstock(t *testing.T, stack Stack) *Rootstock {
	t.Helper()
	client, err := ethclient.DialContext(t.Context(), stack.RskURL)
	require.NoErrorf(t, err, "cannot connect to rskj at %s", stack.RskURL)
	t.Cleanup(client.Close)
	chainID, err := client.ChainID(t.Context())
	require.NoErrorf(t, err, "rskj is not reachable at %s. Run lps-local.sh first", stack.RskURL)
	require.Equal(t, stack.ChainID, chainID.String(), "rskj runs a different chain than the env file says")
	return &Rootstock{client: client, chainID: chainID, funder: common.HexToAddress(unlockedFunder)}
}

// Client is the raw client, for contract bindings and calls this helper does not cover.
func (rsk *Rootstock) Client() *ethclient.Client {
	return rsk.client
}

func (rsk *Rootstock) BlockNumber(t *testing.T) uint64 {
	t.Helper()
	number, err := rsk.client.BlockNumber(t.Context())
	require.NoError(t, err)
	return number
}

func (rsk *Rootstock) Balance(t *testing.T, address common.Address) *big.Int {
	t.Helper()
	balance, err := rsk.client.BalanceAt(t.Context(), address, nil)
	require.NoError(t, err)
	return balance
}

// Mine adds blocks with empty transfers from the unlocked funder.
func (rsk *Rootstock) Mine(t *testing.T, blocks int) {
	t.Helper()
	for range blocks {
		rsk.sendFromFunder(t, rsk.funder, new(big.Int))
	}
}

// NewAccount creates a random account funded with balance.
func (rsk *Rootstock) NewAccount(t *testing.T, balance *big.Int) Account {
	t.Helper()
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	account := Account{Key: key, Address: crypto.PubkeyToAddress(key.PublicKey)}
	if balance.Sign() > 0 {
		rsk.sendFromFunder(t, account.Address, balance)
	}
	return account
}

// Transact signs and sends a transaction from account and fails the test if it reverts.
func (rsk *Rootstock) Transact(t *testing.T, from Account, to common.Address, data []byte, value *big.Int) *types.Receipt {
	t.Helper()
	ctx := t.Context()
	if value == nil {
		value = new(big.Int)
	}
	call := ethereum.CallMsg{From: from.Address, To: &to, Data: data, Value: value}
	gas, err := rsk.client.EstimateGas(ctx, call)
	require.NoErrorf(t, err, "the transaction from %s to %s would fail", from.Address, to)
	gasPrice, err := rsk.client.SuggestGasPrice(ctx)
	require.NoError(t, err)
	nonce, err := rsk.client.PendingNonceAt(ctx, from.Address)
	require.NoError(t, err)

	tx, err := types.SignNewTx(from.Key, types.NewEIP155Signer(rsk.chainID), &types.LegacyTx{
		Nonce: nonce, GasPrice: gasPrice, Gas: gas, To: &to, Value: value, Data: data,
	})
	require.NoError(t, err)
	if err = rsk.client.SendTransaction(ctx, tx); err != nil {
		require.Truef(t, raced(err), "cannot send %s from %s to %s: %v", tx.Hash(), from.Address, to, err)
		t.Logf("[rskj] %s raced another transaction; mining a block to include it", tx.Hash())
		rsk.Mine(t, 1)
	}
	receipt := rsk.WaitForReceipt(t, tx.Hash())
	if receipt.Status != types.ReceiptStatusSuccessful {
		_, replayErr := rsk.client.CallContract(ctx, call, receipt.BlockNumber)
		require.FailNowf(t, "transaction reverted", "%s from %s to %s: %v", tx.Hash(), from.Address, to, replayErr)
	}
	return receipt
}

// WaitForReceipt waits until the transaction is mined.
func (rsk *Rootstock) WaitForReceipt(t *testing.T, hash common.Hash) *types.Receipt {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), receiptTimeout)
	defer cancel()
	for {
		receipt, err := rsk.client.TransactionReceipt(ctx, hash)
		if err == nil {
			return receipt
		}
		require.ErrorIsf(t, err, ethereum.NotFound, "cannot read the receipt of %s", hash)
		select {
		case <-ctx.Done():
			require.FailNowf(t, "transaction not mined", "%s has no receipt after %s", hash, receiptTimeout)
		case <-time.After(receiptPolling):
		}
	}
}

func (rsk *Rootstock) sendFromFunder(t *testing.T, to common.Address, value *big.Int) {
	t.Helper()
	var balanceBefore *big.Int
	if value.Sign() > 0 {
		balanceBefore = rsk.Balance(t, to)
	}
	var hash common.Hash
	err := rsk.client.Client().CallContext(t.Context(), &hash, "eth_sendTransaction", map[string]string{
		"from":  rsk.funder.Hex(),
		"to":    to.Hex(),
		"value": hexutil.EncodeBig(value),
	})
	if err == nil {
		receipt := rsk.WaitForReceipt(t, hash)
		require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
		return
	}
	require.Truef(t, raced(err), "the unlocked funder cannot send: %v", err)
	if value.Sign() == 0 {
		// An empty transfer only has to make rskj mine, and the racing transaction did.
		t.Log("[rskj] an empty transfer raced another transaction; a block was mined anyway")
		return
	}
	t.Logf("[rskj] a transfer of %s wei to %s raced another transaction; mining until it lands", value, to)
	expected := new(big.Int).Add(balanceBefore, value)
	for range racedTransferBlocks {
		rsk.Mine(t, 1)
		if rsk.Balance(t, to).Cmp(expected) >= 0 {
			return
		}
	}
	require.FailNowf(t, "transfer not mined", "%s wei to %s did not land after %d blocks", value, to, racedTransferBlocks)
}

// raced tells whether rskj refused to confirm a transaction only because its block did not
// include it.
func raced(err error) bool {
	return strings.Contains(err.Error(), minedInAnotherBlock)
}
