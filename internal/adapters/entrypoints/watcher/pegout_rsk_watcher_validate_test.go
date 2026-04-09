package watcher

import (
	"math/big"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
)

// nolint:funlen
func TestValidateDepositedPegoutQuote(t *testing.T) {
	const (
		lbcAddress  = "0x5678"
		quoteHash1  = "0102030000000000000000000000000000000000000000000000000000000000"
		quoteHash2  = "0405060000000000000000000000000000000000000000000000000000000000"
		blockNumber = uint64(10)
		height      = uint64(16) // blockNumber + DepositConfirmations(5) < height: 15 < 16
	)

	baseQuote := quote.PegoutQuote{
		LbcAddress:           lbcAddress,
		Value:                entities.NewWei(3),
		GasFee:               entities.NewWei(0),
		CallFee:              entities.NewWei(0),
		DepositConfirmations: 5,
	}
	baseRetained := quote.RetainedPegoutQuote{
		QuoteHash: quoteHash1,
		State:     quote.PegoutStateWaitingForDepositConfirmations,
	}
	watchedQuote := quote.WatchedPegoutQuote{
		PegoutQuote:   baseQuote,
		RetainedQuote: baseRetained,
	}

	newReceipt := func() *blockchain.TransactionReceipt {
		return &blockchain.TransactionReceipt{
			TransactionHash:   "0xaabbcc",
			BlockHash:         "0xddeeff",
			BlockNumber:       blockNumber,
			From:              "0x1234",
			CumulativeGasUsed: big.NewInt(100),
			GasUsed:           big.NewInt(100),
			Value:             entities.NewWei(0),
		}
	}

	t.Run("single-event receipt with matching hash and sufficient amount returns true", func(t *testing.T) {
		receipt := newReceipt()
		receipt = test.AppendDepositLogFromQuote(t, receipt, baseQuote, baseRetained)
		assert.True(t, validateDepositedPegoutQuote(watchedQuote, *receipt, height))
	})

	t.Run("two-event receipt returns true when second quote matches and amount is sufficient", func(t *testing.T) {
		otherRetained := baseRetained
		otherRetained.QuoteHash = quoteHash2
		receipt := newReceipt()
		receipt = test.AppendDepositLogFromQuote(t, receipt, baseQuote, otherRetained)
		receipt = test.AppendDepositLogFromQuote(t, receipt, baseQuote, baseRetained)
		assert.True(t, validateDepositedPegoutQuote(watchedQuote, *receipt, height))
	})

	t.Run("two-event receipt returns false when second quote has insufficient amount", func(t *testing.T) {
		otherRetained := baseRetained
		otherRetained.QuoteHash = quoteHash2
		insufficientQuote := baseQuote
		insufficientQuote.Value = entities.NewWei(2) // total = 2 < watchedQuote.Total() = 3
		receipt := newReceipt()
		receipt = test.AppendDepositLogFromQuote(t, receipt, baseQuote, otherRetained)
		receipt = test.AppendDepositLogFromQuote(t, receipt, insufficientQuote, baseRetained)
		assert.False(t, validateDepositedPegoutQuote(watchedQuote, *receipt, height))
	})

	t.Run("wrong LBC address in event returns false", func(t *testing.T) {
		wrongAddrQuote := baseQuote
		wrongAddrQuote.LbcAddress = "0x9999"
		receipt := newReceipt()
		receipt = test.AppendDepositLogFromQuote(t, receipt, wrongAddrQuote, baseRetained)
		assert.False(t, validateDepositedPegoutQuote(watchedQuote, *receipt, height))
	})

	t.Run("not enough confirmations returns false", func(t *testing.T) {
		// blockNumber(10) + DepositConfirmations(5) = 15; need < height, but 15 < 15 is false
		insufficientHeight := uint64(15)
		receipt := newReceipt()
		receipt = test.AppendDepositLogFromQuote(t, receipt, baseQuote, baseRetained)
		assert.False(t, validateDepositedPegoutQuote(watchedQuote, *receipt, insufficientHeight))
	})
}
