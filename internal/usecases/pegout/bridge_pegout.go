package pegout

import (
	"context"
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
) RebalanceHandler {
	switch strategy {
	case UtxoSplit:
		return NewUtxoSplitHandler(quoteRepository, rskWallet, contracts, mutex)
	default:
		return NewAllAtOnceHandler(quoteRepository, rskWallet, contracts, mutex)
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
