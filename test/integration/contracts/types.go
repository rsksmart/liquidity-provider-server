package contracts

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/stretchr/testify/suite"
	"time"
)

type LiquidityBridgeContractExecutor interface {
	DepositPegout(s TestSuite, opts *bind.TransactOpts, pegoutQuote pkg.PegoutQuoteDTO, hexSignature string) (*types.Receipt, *types.Transaction)
	GetRefundPegoutEvent(s TestSuite, timeout time.Duration, quoteHash string) RefundPegoutEvent
}

type RefundPegoutEvent struct {
	QuoteHash string
	RawEvent  types.Log
}

type TestSuite interface {
	Raw() *suite.Suite
	RskClient() *ethclient.Client
}
