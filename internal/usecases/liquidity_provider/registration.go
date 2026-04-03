package liquidity_provider

import (
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type RegistrationUseCase struct {
	contracts    blockchain.RskContracts
	provider     liquidity_provider.LiquidityProvider
	pollInterval time.Duration
}

func NewRegistrationUseCase(
	contracts blockchain.RskContracts,
	provider liquidity_provider.LiquidityProvider,
	pollInterval time.Duration,
) *RegistrationUseCase {
	return &RegistrationUseCase{
		contracts:    contracts,
		provider:     provider,
		pollInterval: pollInterval,
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

func (useCase *RegistrationUseCase) Run(params blockchain.ProviderRegistrationParams) (int64, error) {
	state, err := useCase.contracts.Discovery.GetRegistrationState(useCase.provider.RskAddress())
	if err != nil {
		return 0, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	}
	switch state {
	case blockchain.RegistrationStateApproved:
		return useCase.handleApproved()
	case blockchain.RegistrationStatePending:
		return useCase.handlePending()
	case blockchain.RegistrationStateRejected:
		return useCase.handleRejected()
	default:
		return useCase.handleNoneOrWithdrawn(params)
	}
}

func (useCase *RegistrationUseCase) handleApproved() (int64, error) {
	provider, err := useCase.contracts.Discovery.GetProvider(useCase.provider.RskAddress())
	if err != nil {
		return 0, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	}
	return int64(provider.Id), nil
}

func (useCase *RegistrationUseCase) handlePending() (int64, error) {
	log.Info("Registration pending admin approval, waiting...")
	return useCase.waitForApproval()
}

func (useCase *RegistrationUseCase) handleRejected() (int64, error) {
	log.Error("Registration rejected by admin. Contact an admin to approve your registration before restarting.")
	return 0, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, usecases.RegistrationRejectedError)
}

func (useCase *RegistrationUseCase) handleNoneOrWithdrawn(params blockchain.ProviderRegistrationParams) (int64, error) {
	if err := usecases.CheckPauseState(useCase.contracts.Discovery, useCase.contracts.CollateralManagement); err != nil {
		return 0, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	}
	if err := useCase.validateParams(params); err != nil {
		return 0, err
	}
	collateral, err := useCase.getCollateralInfo()
	if err != nil {
		return 0, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	}
	operational, err := useCase.getOperationalInfo()
	if err != nil {
		return 0, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	}
	if _, err = useCase.addPeginCollateral(params, operational, collateral); err != nil {
		return 0, err
	}
	if _, err = useCase.addPegoutCollateral(params, operational, collateral); err != nil {
		return 0, err
	}
	log.Debug("Registering new provider...")
	if _, err = useCase.registerProvider(params, collateral); err != nil {
		return 0, err
	}
	log.Info("Registration submitted, waiting for admin approval...")
	return useCase.waitForApproval()
}

func (useCase *RegistrationUseCase) waitForApproval() (int64, error) {
	for {
		state, err := useCase.contracts.Discovery.GetRegistrationState(useCase.provider.RskAddress())
		if err != nil {
			return 0, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
		}
		if state == blockchain.RegistrationStateApproved {
			return useCase.handleApproved()
		} else if state == blockchain.RegistrationStateRejected {
			return useCase.handleRejected()
		}
		time.Sleep(useCase.pollInterval)
	}
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
