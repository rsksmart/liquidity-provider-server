package pegout

import (
	"context"
	"errors"
	"sync"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
)

const (
	// BridgeConversionGasLimit see https://dev.rootstock.io/rsk/rbtc/conversion/networks/
	BridgeConversionGasLimit = 100000
	// BridgeConversionGasPrice see https://dev.rootstock.io/rsk/rbtc/conversion/networks/
	BridgeConversionGasPrice = 60000000
)

// ErrBridgeReleaseRejected wraps errors returned when the Bridge emits release_request_rejected
// for a bridge conversion transaction. Callers can identify the case via errors.Is.
var ErrBridgeReleaseRejected = errors.New("bridge release request rejected")

// ReleaseRejectionParser inspects a receipt for the Bridge release_request_rejected event
// emitted by bridgeAddress. It returns rejected=true together with the decoded reason when found.
type ReleaseRejectionParser = func(receipt blockchain.TransactionReceipt, bridgeAddress string) (bool, blockchain.RejectedPegoutReason, error)

type RebalanceStrategy string

const (
	AllAtOnce RebalanceStrategy = "ALL_AT_ONCE"
	UtxoSplit RebalanceStrategy = "UTXO_SPLIT"
)

type RebalanceHandler interface {
	Execute(
		ctx context.Context,
		pegoutConfig liquidity_provider.PegoutConfiguration,
		watchedQuotes []quote.WatchedPegoutQuote,
	) error
}

func NewRebalanceHandler(
	strategy RebalanceStrategy,
	quoteRepository quote.PegoutQuoteRepository,
	rskWallet blockchain.RootstockWallet,
	contracts blockchain.RskContracts,
	mutex sync.Locker,
	releaseRejectionParser ReleaseRejectionParser,
) RebalanceHandler {
	switch strategy {
	case UtxoSplit:
		return NewUtxoSplitHandler(quoteRepository, rskWallet, contracts, mutex, releaseRejectionParser)
	default:
		return NewAllAtOnceHandler(quoteRepository, rskWallet, contracts, mutex, releaseRejectionParser)
	}
}

type BridgePegoutUseCase struct {
	pegoutProvider liquidity_provider.PegoutLiquidityProvider
	handler        RebalanceHandler
}

func NewBridgePegoutUseCase(
	pegoutProvider liquidity_provider.PegoutLiquidityProvider,
	handler RebalanceHandler,
) *BridgePegoutUseCase {
	return &BridgePegoutUseCase{
		pegoutProvider: pegoutProvider,
		handler:        handler,
	}
}

func (useCase *BridgePegoutUseCase) Run(ctx context.Context, watchedQuotes ...quote.WatchedPegoutQuote) error {
	pegoutConfig := useCase.pegoutProvider.PegoutConfiguration(ctx)
	return useCase.handler.Execute(ctx, pegoutConfig, watchedQuotes)
}
