package pegout

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type InitPegoutDepositCacheUseCase struct {
	pegoutRepository quote.PegoutQuoteRepository
	contracts        blockchain.RskContracts
	rpc              blockchain.Rpc
}

func NewInitPegoutDepositCacheUseCase(pegoutRepository quote.PegoutQuoteRepository, contracts blockchain.RskContracts, rpc blockchain.Rpc) *InitPegoutDepositCacheUseCase {
	return &InitPegoutDepositCacheUseCase{pegoutRepository: pegoutRepository, contracts: contracts, rpc: rpc}
}

func (useCase *InitPegoutDepositCacheUseCase) Run(ctx context.Context, cacheStartBlock uint64) error {
	var deposits []quote.PegoutDeposit
	var err error
	var height uint64
	if height, err = useCase.rpc.Rsk.GetHeight(ctx); err != nil {
		return usecases.WrapUseCaseError(usecases.InitPegoutDepositCacheId, err)
	}

	if deposits, err = useCase.contracts.PegOut.GetDepositEvents(ctx, cacheStartBlock, &height); err != nil {
		return usecases.WrapUseCaseError(usecases.InitPegoutDepositCacheId, err)
	}

	if err = useCase.pegoutRepository.UpsertPegoutDeposits(ctx, deposits); err != nil {
		return usecases.WrapUseCaseError(usecases.InitPegoutDepositCacheId, err)
	}
	return nil
}
