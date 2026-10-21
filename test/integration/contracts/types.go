package contracts

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/stretchr/testify/suite"
	"math/big"
	"time"
)

type LiquidityBridgeContractExecutor interface {
	DepositPegout(s TestSuite, opts *bind.TransactOpts, pegoutQuote pkg.PegoutQuoteDTO, hexSignature string) (*types.Receipt, *types.Transaction)
	GetRefundPegoutEvent(s TestSuite, timeout time.Duration, quoteHash string) RefundPegoutEvent
	GetCallForUserEvent(s TestSuite, timeout time.Duration, userAddress, providerAddress string) CallForUserEvent
	GetPeginRegisteredEvent(s TestSuite, timeout time.Duration, quoteHash string) PegInRegisteredEvent
}

type RefundPegoutEvent struct {
	QuoteHash string
	RawEvent  types.Log
}

type PegInRegisteredEvent struct {
	QuoteHash string
	Amount    *entities.Wei
	RawEvent  types.Log
}

type CallForUserEvent struct {
	From      string
	To        string
	QuoteHash string
	GasLimit  uint64
	Value     *big.Int
	Data      []byte
	Success   bool
}

type TestSuite interface {
	Raw() *suite.Suite
	RskClient() *ethclient.Client
}
