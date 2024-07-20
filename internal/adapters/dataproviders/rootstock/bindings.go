package rootstock

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"math/big"
)

type RpcClientBinding interface {
	bind.ContractTransactor
	bind.ContractCaller
	bind.ContractFilterer
	Close()
	ChainID(ctx context.Context) (*big.Int, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
}

type RskBridgeBinding interface {
	GetFederationAddress(opts *bind.CallOpts) (string, error)
	GetMinimumLockTxValue(opts *bind.CallOpts) (*big.Int, error)
	GetActivePowpegRedeemScript(opts *bind.CallOpts) ([]byte, error)
	GetFederationSize(opts *bind.CallOpts) (*big.Int, error)
	GetFederatorPublicKeyOfType(opts *bind.CallOpts, index *big.Int, atype string) ([]byte, error)
	GetFederationThreshold(opts *bind.CallOpts) (*big.Int, error)
	GetActiveFederationCreationBlockHeight(opts *bind.CallOpts) (*big.Int, error)
	IsBtcTxHashAlreadyProcessed(opts *bind.CallOpts, hash string) (bool, error)
	HasBtcBlockCoinbaseTransactionInformation(opts *bind.CallOpts, blockHash [32]byte) (bool, error)
	GetBtcBlockchainBestChainHeight(opts *bind.CallOpts) (*big.Int, error)
	RegisterBtcCoinbaseTransaction(opts *bind.TransactOpts, btcTxSerialized []byte, blockHash [32]byte, pmtSerialized []byte, witnessMerkleRoot [32]byte, witnessReservedValue [32]byte) (*types.Transaction, error)
}

type LbcBinding interface {
	HashQuote(opts *bind.CallOpts, quote bindings.QuotesPeginQuote) ([32]byte, error)
	HashPegoutQuote(opts *bind.CallOpts, quote bindings.QuotesPegOutQuote) ([32]byte, error)
	GetProviderIds(opts *bind.CallOpts) (*big.Int, error)
	GetProviders(opts *bind.CallOpts) ([]bindings.LiquidityBridgeContractLiquidityProvider, error)
	GetProvider(opts *bind.CallOpts, providerAddress common.Address) (bindings.LiquidityBridgeContractLiquidityProvider, error)
	Resign(opts *bind.TransactOpts) (*types.Transaction, error)
	SetProviderStatus(opts *bind.TransactOpts, _providerId *big.Int, status bool) (*types.Transaction, error)
	GetCollateral(opts *bind.CallOpts, addr common.Address) (*big.Int, error)
	GetPegoutCollateral(opts *bind.CallOpts, addr common.Address) (*big.Int, error)
	GetMinCollateral(opts *bind.CallOpts) (*big.Int, error)
	AddCollateral(opts *bind.TransactOpts) (*types.Transaction, error)
	AddPegoutCollateral(opts *bind.TransactOpts) (*types.Transaction, error)
	WithdrawCollateral(opts *bind.TransactOpts) (*types.Transaction, error)
	GetBalance(opts *bind.CallOpts, addr common.Address) (*big.Int, error)
	CallForUser(opts *bind.TransactOpts, quote bindings.QuotesPeginQuote) (*types.Transaction, error)
	RegisterPegIn(opts *bind.TransactOpts, quote bindings.QuotesPeginQuote, signature []byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) (*types.Transaction, error)
	RefundPegOut(opts *bind.TransactOpts, quoteHash [32]byte, btcTx []byte, btcBlockHeaderHash [32]byte, partialMerkleTree *big.Int, merkleBranchHashes [][32]byte) (*types.Transaction, error)
	IsOperational(opts *bind.CallOpts, addr common.Address) (bool, error)
	IsOperationalForPegout(opts *bind.CallOpts, addr common.Address) (bool, error)
	Register(opts *bind.TransactOpts, _name string, _apiBaseUrl string, _status bool, _providerType string) (*types.Transaction, error)
	FilterPegOutDeposit(opts *bind.FilterOpts, quoteHash [][32]byte, sender []common.Address) (*bindings.LiquidityBridgeContractPegOutDepositIterator, error)
	FilterPenalized(opts *bind.FilterOpts) (*bindings.LiquidityBridgeContractPenalizedIterator, error)
	ParseRegister(log types.Log) (*bindings.LiquidityBridgeContractRegister, error)
	ProductFeePercentage(opts *bind.CallOpts) (*big.Int, error)
	IsPegOutQuoteCompleted(opts *bind.CallOpts, quoteHash [32]byte) (bool, error)
}

type LbcAdapter interface {
	LbcBinding
	Caller() ContractCallerBinding
	DepositEventIteratorAdapter(rawIterator *bindings.LiquidityBridgeContractPegOutDepositIterator) EventIteratorAdapter[bindings.LiquidityBridgeContractPegOutDeposit]
	PenalizedEventIteratorAdapter(rawIterator *bindings.LiquidityBridgeContractPenalizedIterator) EventIteratorAdapter[bindings.LiquidityBridgeContractPenalized]
}

type EventIteratorAdapter[T any] interface {
	Next() bool
	Close() error
	Event() *T
	Error() error
}

type ContractCallerBinding interface {
	Call(opts *bind.CallOpts, result *[]any, method string, params ...any) error
}

type depositEventIteratorAdapter struct {
	*bindings.LiquidityBridgeContractPegOutDepositIterator
}

func (i *depositEventIteratorAdapter) Event() *bindings.LiquidityBridgeContractPegOutDeposit {
	return i.LiquidityBridgeContractPegOutDepositIterator.Event
}

type penalizedEventIteratorAdapter struct {
	*bindings.LiquidityBridgeContractPenalizedIterator
}

func (i *penalizedEventIteratorAdapter) Event() *bindings.LiquidityBridgeContractPenalized {
	return i.LiquidityBridgeContractPenalizedIterator.Event
}

type lbcAdapter struct {
	*bindings.LiquidityBridgeContract
}

func NewLbcAdapter(liquidityBridgeContract *bindings.LiquidityBridgeContract) LbcAdapter {
	return &lbcAdapter{LiquidityBridgeContract: liquidityBridgeContract}
}

func (lbc *lbcAdapter) Caller() ContractCallerBinding {
	return &bindings.LiquidityBridgeContractCallerRaw{
		Contract: &lbc.LiquidityBridgeContract.LiquidityBridgeContractCaller,
	}
}

func (lbc *lbcAdapter) DepositEventIteratorAdapter(rawIterator *bindings.LiquidityBridgeContractPegOutDepositIterator) EventIteratorAdapter[bindings.LiquidityBridgeContractPegOutDeposit] {
	return &depositEventIteratorAdapter{LiquidityBridgeContractPegOutDepositIterator: rawIterator}
}

func (lbc *lbcAdapter) PenalizedEventIteratorAdapter(rawIterator *bindings.LiquidityBridgeContractPenalizedIterator) EventIteratorAdapter[bindings.LiquidityBridgeContractPenalized] {
	return &penalizedEventIteratorAdapter{LiquidityBridgeContractPenalizedIterator: rawIterator}
}
