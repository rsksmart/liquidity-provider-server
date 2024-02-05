package pegout

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type InitPegoutDepositCacheUseCase struct {
	pegoutRepository quote.PegoutQuoteRepository
	lbc              blockchain.LiquidityBridgeContract
	rskRpc           blockchain.RootstockRpcServer
}

func NewInitPegoutDepositCacheUseCase(pegoutRepository quote.PegoutQuoteRepository, lbc blockchain.LiquidityBridgeContract, rskRpc blockchain.RootstockRpcServer) *InitPegoutDepositCacheUseCase {
	return &InitPegoutDepositCacheUseCase{pegoutRepository: pegoutRepository, lbc: lbc, rskRpc: rskRpc}
}

func (useCase *InitPegoutDepositCacheUseCase) Run(ctx context.Context, cacheStartBlock uint64) error {
	var deposits []quote.PegoutDeposit
	var err error
	var height uint64
	if height, err = useCase.rskRpc.GetHeight(ctx); err != nil {
		return usecases.WrapUseCaseError(usecases.InitPegoutDepositCacheId, err)
	}

	if deposits, err = useCase.lbc.GetDepositEvents(ctx, cacheStartBlock, &height); err != nil {
		return usecases.WrapUseCaseError(usecases.InitPegoutDepositCacheId, err)
	}

	if err = useCase.pegoutRepository.UpsertPegoutDeposits(ctx, deposits); err != nil {
		return usecases.WrapUseCaseError(usecases.InitPegoutDepositCacheId, err)
	}
	return nil
}
