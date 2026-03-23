package pegout

import (
	"context"
	"sync"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
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
) RebalanceHandler {
	switch strategy {
	case UtxoSplit:
		return NewUtxoSplitHandler(quoteRepository, rskWallet, contracts)
	default:
		return NewAllAtOnceHandler(quoteRepository, rskWallet, contracts)
	}
}

type BridgePegoutUseCase struct {
	pegoutProvider liquidity_provider.PegoutLiquidityProvider
	rskWalletMutex sync.Locker
	handler        RebalanceHandler
}

func NewBridgePegoutUseCase(
	pegoutProvider liquidity_provider.PegoutLiquidityProvider,
	rskWalletMutex sync.Locker,
	handler RebalanceHandler,
) *BridgePegoutUseCase {
	return &BridgePegoutUseCase{
		pegoutProvider: pegoutProvider,
		rskWalletMutex: rskWalletMutex,
		handler:        handler,
	}
}

func (useCase *BridgePegoutUseCase) Run(ctx context.Context, watchedQuotes ...quote.WatchedPegoutQuote) error {
	pegoutConfig := useCase.pegoutProvider.PegoutConfiguration(ctx)

	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

	return useCase.handler.Execute(ctx, pegoutConfig, watchedQuotes)
}

func checkBalance(ctx context.Context, rskWallet blockchain.RootstockWallet, requiredBalance *entities.Wei) error {
	balance, err := rskWallet.GetBalance(ctx)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, err)
	}
	if balance.Cmp(requiredBalance) < 0 {
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, usecases.InsufficientAmountError)
	}
	return nil
}

func calculateTotalToPegout(watchedQuotes []quote.WatchedPegoutQuote) *entities.Wei {
	totalValue := new(entities.Wei)
	for _, watchedQuote := range watchedQuotes {
		totalValue.Add(totalValue, quoteContribution(watchedQuote))
	}
	return totalValue
}

func quoteContribution(wq quote.WatchedPegoutQuote) *entities.Wei {
	contribution := new(entities.Wei)
	if wq.PegoutQuote.Value != nil {
		contribution.Add(contribution, wq.PegoutQuote.Value)
	}
	if wq.PegoutQuote.CallFee != nil {
		contribution.Add(contribution, wq.PegoutQuote.CallFee)
	}
	if wq.PegoutQuote.GasFee != nil {
		contribution.Add(contribution, wq.PegoutQuote.GasFee)
	}
	return contribution
}
