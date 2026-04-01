package common

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

func CheckBalance(ctx context.Context, useCaseId usecases.UseCaseId, rskWallet blockchain.RootstockWallet, requiredBalance *entities.Wei) error {
	balance, err := rskWallet.GetBalance(ctx)
	if err != nil {
		return usecases.WrapUseCaseError(useCaseId, err)
	}
	if balance.Cmp(requiredBalance) < 0 {
		return usecases.WrapUseCaseError(useCaseId, usecases.InsufficientAmountError)
	}
	return nil
}

func CalculateTotalToPegout(watchedQuotes []quote.WatchedPegoutQuote) *entities.Wei {
	totalValue := new(entities.Wei)
	for _, watchedQuote := range watchedQuotes {
		totalValue.Add(totalValue, watchedQuote.Remaining())
	}
	return totalValue
}
