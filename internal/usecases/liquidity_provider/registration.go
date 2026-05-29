package liquidity_provider

import (
	"context"
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type RegistrationUseCase struct {
	contracts blockchain.RskContracts
	provider  liquidity_provider.LiquidityProvider
}

func NewRegistrationUseCase(
	contracts blockchain.RskContracts,
	provider liquidity_provider.LiquidityProvider,
) *RegistrationUseCase {
	return &RegistrationUseCase{
		contracts: contracts,
		provider:  provider,
	}
}

type collateralInfo struct {
	peginCollateral   *entities.Wei
	pegoutCollateral  *entities.Wei
	minimumCollateral *entities.Wei
}

type operationalInfo struct {
	operationalForPegin  bool
	operationalForPegout bool
}

func (useCase *RegistrationUseCase) Run(ctx context.Context, params blockchain.ProviderRegistrationParams) (int64, error) {
	state, err := useCase.contracts.Discovery.GetRegistrationState(useCase.provider.RskAddress())
	if err != nil {
		return 0, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	}

	if state == blockchain.RegistrationStateNone {
		if err = useCase.registerForApproval(params); err != nil {
			return 0, err
		}
		state = blockchain.RegistrationStatePending
	}

	if state == blockchain.RegistrationStatePending {
		state, err = useCase.contracts.Discovery.WatchRegistrationApproval(ctx, useCase.provider.RskAddress())
		if err != nil {
			return 0, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
		}
	}

	switch state {
	case blockchain.RegistrationStateApproved:
		provider, providerErr := useCase.contracts.Discovery.GetProvider(useCase.provider.RskAddress())
		if providerErr != nil {
			return 0, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, providerErr)
		}
		return int64(provider.Id), nil
	case blockchain.RegistrationStateRejected:
		return 0, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, usecases.RegistrationRejectedError)
	case blockchain.RegistrationStateWithdrawn:
		return 0, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, usecases.RegistrationWithdrawnError)
	default:
		return 0, usecases.WrapUseCaseError(usecases.ProviderRegistrationId,
			fmt.Errorf("unexpected registration state %d", state))
	}
}

func (useCase *RegistrationUseCase) registerForApproval(params blockchain.ProviderRegistrationParams) error {
	if err := usecases.CheckPauseState(useCase.contracts.Discovery, useCase.contracts.CollateralManagement); err != nil {
		return usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	}
	if err := useCase.validateParams(params); err != nil {
		return err
	}
	collateral, err := useCase.getCollateralInfo()
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	}
	operational, err := useCase.getOperationalInfo()
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	}
	if _, err = useCase.addPeginCollateral(params, operational, collateral); err != nil {
		return err
	}
	if _, err = useCase.addPegoutCollateral(params, operational, collateral); err != nil {
		return err
	}
	log.Debug("Registering new provider...")
	if _, err = useCase.registerProvider(params, collateral); err != nil {
		return err
	}
	log.Info("Registration submitted, waiting for admin approval...")
	return nil
}

func (useCase *RegistrationUseCase) getCollateralInfo() (collateralInfo, error) {
	var err error
	var peginCollateral, pegoutCollateral, minimumCollateral *entities.Wei

	if minimumCollateral, err = useCase.contracts.CollateralManagement.GetMinimumCollateral(); err != nil {
		return collateralInfo{}, err
	}
	if peginCollateral, err = useCase.contracts.CollateralManagement.GetCollateral(useCase.provider.RskAddress()); err != nil {
		return collateralInfo{}, err
	}
	if pegoutCollateral, err = useCase.contracts.CollateralManagement.GetPegoutCollateral(useCase.provider.RskAddress()); err != nil {
		return collateralInfo{}, err
	}
	return collateralInfo{
		peginCollateral:   peginCollateral,
		pegoutCollateral:  pegoutCollateral,
		minimumCollateral: minimumCollateral,
	}, nil
}

func (useCase *RegistrationUseCase) getOperationalInfo() (operationalInfo, error) {
	var operationalForPegin, operationalForPegout bool
	var err error
	if operationalForPegin, err = useCase.contracts.Discovery.IsOperational(liquidity_provider.PeginProvider, useCase.provider.RskAddress()); err != nil {
		return operationalInfo{}, err
	}
	if operationalForPegout, err = useCase.contracts.Discovery.IsOperational(liquidity_provider.PegoutProvider, useCase.provider.RskAddress()); err != nil {
		return operationalInfo{}, err
	}
	return operationalInfo{
		operationalForPegin:  operationalForPegin,
		operationalForPegout: operationalForPegout,
	}, nil
}

func (useCase *RegistrationUseCase) registerProvider(params blockchain.ProviderRegistrationParams, collateral collateralInfo) (int64, error) {
	value := new(entities.Wei)
	txConfig := blockchain.NewTransactionConfig(value.Mul(collateral.minimumCollateral, entities.NewUWei(2)), 0, nil)
	if id, err := useCase.contracts.Discovery.RegisterProvider(txConfig, params); err != nil {
		return 0, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	} else {
		return id, nil
	}
}

func (useCase *RegistrationUseCase) validateParams(params blockchain.ProviderRegistrationParams) error {
	var err error
	if err = entities.ValidateStruct(params); err != nil {
		return usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	} else if !params.Type.IsValid() {
		return usecases.WrapUseCaseError(usecases.ProviderRegistrationId, liquidity_provider.InvalidProviderTypeError)
	}
	return nil
}

func (useCase *RegistrationUseCase) addPeginCollateral(
	params blockchain.ProviderRegistrationParams,
	operational operationalInfo,
	collateral collateralInfo,
) (bool, error) {
	if !(params.Type.AcceptsPegin() && !operational.operationalForPegin && collateral.peginCollateral.Cmp(entities.NewWei(0)) != 0) {
		return false, nil
	}
	collateralToAdd := new(entities.Wei)
	log.Debug("Adding pegin collateral...")
	if err := useCase.contracts.CollateralManagement.AddCollateral(collateralToAdd.Sub(collateral.minimumCollateral, collateral.peginCollateral)); err != nil {
		return false, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	}
	return true, nil
}

func (useCase *RegistrationUseCase) addPegoutCollateral(
	params blockchain.ProviderRegistrationParams,
	operational operationalInfo,
	collateral collateralInfo,
) (bool, error) {
	if !(params.Type.AcceptsPegout() && !operational.operationalForPegout && collateral.pegoutCollateral.Cmp(entities.NewWei(0)) != 0) {
		return false, nil
	}
	collateralToAdd := new(entities.Wei)
	log.Debug("Adding pegout collateral...")
	if err := useCase.contracts.CollateralManagement.AddPegoutCollateral(collateralToAdd.Sub(collateral.minimumCollateral, collateral.pegoutCollateral)); err != nil {
		return false, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	}
	return true, nil
}
