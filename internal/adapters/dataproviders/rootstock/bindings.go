package rootstock

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
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
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
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
	FilterBatchPegoutCreated(opts *bind.FilterOpts, btcTxHash [][32]byte) (*bindings.RskBridgeBatchPegoutCreatedIterator, error)
}

type PegoutBinding interface {
	HashPegOutQuote(opts *bind.CallOpts, quote bindings.QuotesPegOutQuote) ([32]byte, error)
	RefundPegOut(opts *bind.TransactOpts, quoteHash [32]byte, btcTx []byte, btcBlockHeaderHash [32]byte, partialMerkleTree *big.Int, merkleBranchHashes [][32]byte) (*types.Transaction, error)
	FilterPegOutDeposit(opts *bind.FilterOpts, quoteHash [][32]byte, sender []common.Address, timestamp []*big.Int) (*bindings.IPegOutPegOutDepositIterator, error)
	IsQuoteCompleted(opts *bind.CallOpts, quoteHash [32]byte) (bool, error)
	RefundUserPegOut(opts *bind.TransactOpts, quoteHash [32]byte) (*types.Transaction, error)
	GetFeePercentage(opts *bind.CallOpts) (*big.Int, error)
	ValidatePegout(opts *bind.CallOpts, quoteHash [32]byte, btcTx []byte) (bindings.QuotesPegOutQuote, error)
}

type PegoutContractAdapter interface {
	PegoutBinding
	Caller() ContractCallerBinding
	DepositEventIteratorAdapter(rawIterator *bindings.IPegOutPegOutDepositIterator) EventIteratorAdapter[bindings.IPegOutPegOutDeposit]
}

type PeginBinding interface {
	HashPegInQuote(opts *bind.CallOpts, quote bindings.QuotesPegInQuote) ([32]byte, error)
	RegisterPegIn(opts *bind.TransactOpts, quote bindings.QuotesPegInQuote, signature []byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) (*types.Transaction, error)
	CallForUser(opts *bind.TransactOpts, quote bindings.QuotesPegInQuote) (*types.Transaction, error)
	GetBalance(opts *bind.CallOpts, addr common.Address) (*big.Int, error)
	GetFeePercentage(opts *bind.CallOpts) (*big.Int, error)
}

type PeginContractAdapter interface {
	PeginBinding
	Caller() ContractCallerBinding
}

type DiscoveryBinding interface {
	IsOperational(opts *bind.CallOpts, providerType uint8, addr common.Address) (bool, error)
	GetProviders(opts *bind.CallOpts) ([]bindings.FlyoverLiquidityProvider, error)
	GetProvider(opts *bind.CallOpts, providerAddress common.Address) (bindings.FlyoverLiquidityProvider, error)
	Register(opts *bind.TransactOpts, name string, apiBaseUrl string, status bool, providerType uint8) (*types.Transaction, error)
	UpdateProvider(opts *bind.TransactOpts, _name string, _url string) (*types.Transaction, error)
	SetProviderStatus(opts *bind.TransactOpts, _providerId *big.Int, status bool) (*types.Transaction, error)
	ParseRegister(log types.Log) (*bindings.IFlyoverDiscoveryRegister, error)
}

type CollateralManagementBinding interface {
	FilterPenalized(opts *bind.FilterOpts, liquidityProvider []common.Address, punisher []common.Address, quoteHash [][32]byte) (*bindings.ICollateralManagementPenalizedIterator, error)
	Resign(opts *bind.TransactOpts) (*types.Transaction, error)
	GetPegInCollateral(opts *bind.CallOpts, addr common.Address) (*big.Int, error)
	GetPegOutCollateral(opts *bind.CallOpts, addr common.Address) (*big.Int, error)
	GetMinCollateral(opts *bind.CallOpts) (*big.Int, error)
	AddPegInCollateral(opts *bind.TransactOpts) (*types.Transaction, error)
	AddPegOutCollateral(opts *bind.TransactOpts) (*types.Transaction, error)
	WithdrawCollateral(opts *bind.TransactOpts) (*types.Transaction, error)
}

type CollateralManagementAdapter interface {
	CollateralManagementBinding
	Caller() ContractCallerBinding
	PenalizedEventIteratorAdapter(rawIterator *bindings.ICollateralManagementPenalizedIterator) EventIteratorAdapter[bindings.ICollateralManagementPenalized]
}

type RskBridgeAdapter interface {
	RskBridgeBinding
	BatchPegOutCreatedIteratorAdapter(rawIterator *bindings.RskBridgeBatchPegoutCreatedIterator) EventIteratorAdapter[bindings.RskBridgeBatchPegoutCreated]
}

type ContractCallerBinding interface {
	Call(opts *bind.CallOpts, result *[]any, method string, params ...any) error
}
