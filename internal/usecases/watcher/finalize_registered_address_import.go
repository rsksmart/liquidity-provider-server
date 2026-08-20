package watcher

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type FinalizeRegisteredAddressImportUseCase struct {
	repository rootstock.PegInAddressRegistryWatchRepository
}

func NewFinalizeRegisteredAddressImportUseCase(
	repository rootstock.PegInAddressRegistryWatchRepository,
) *FinalizeRegisteredAddressImportUseCase {
	return &FinalizeRegisteredAddressImportUseCase{repository: repository}
}

func (useCase *FinalizeRegisteredAddressImportUseCase) Run(
	ctx context.Context,
	pending []*rootstock.PegInAddressRegistryWatchEntry,
	rescanErr error,
) error {
	if len(pending) == 0 {
		return nil
	}
	if rescanErr != nil {
		for _, entry := range pending {
			if err := persistWatchError(
				ctx,
				useCase.repository,
				entry,
				rescanErr,
				usecases.FinalizeRegisteredAddressImportId,
			); err != nil {
				return err
			}
		}
		return nil
	}
	for _, entry := range pending {
		if err := persistImported(
			ctx,
			useCase.repository,
			entry,
			usecases.FinalizeRegisteredAddressImportId,
		); err != nil {
			return err
		}
	}
	return nil
}
