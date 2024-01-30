package liquidity_provider

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type RegistrationUseCase struct {
	lbc      blockchain.LiquidityBridgeContract
	provider entities.LiquidityProvider
}

func NewRegistrationUseCase(lbc blockchain.LiquidityBridgeContract, provider entities.LiquidityProvider) *RegistrationUseCase {
	return &RegistrationUseCase{lbc: lbc, provider: provider}
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
	var collateral collateralInfo
	var operational operationalInfo
	var addedCollateral bool
	var id int64
	var err error

	if err = useCase.validateParams(params); err != nil {
		return 0, err
	}

	if collateral, err = useCase.getCollateralInfo(); err != nil {
		return 0, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	}
	if operational, err = useCase.getOperationalInfo(); err != nil {
		return 0, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	}

	if useCase.isProviderOperational(params.Type, operational) {
		return 0, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, usecases.AlreadyRegisteredError)
	}

	if addedCollateral, err = useCase.addPeginCollateral(params, operational, collateral); addedCollateral == true || err != nil {
		return 0, err
	}

	if addedCollateral, err = useCase.addPegoutCollateral(params, operational, collateral); addedCollateral == true || err != nil {
		return 0, err
	}

	log.Debug("Registering new provider...")
	if id, err = useCase.registerProvider(params, collateral); err != nil {
		return 0, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	}
	return id, nil
}

func (useCase *RegistrationUseCase) getCollateralInfo() (collateralInfo, error) {
	var err error
	var peginCollateral, pegoutCollateral, minimumCollateral *entities.Wei

	if minimumCollateral, err = useCase.lbc.GetMinimumCollateral(); err != nil {
		return collateralInfo{}, err
	}
	if peginCollateral, err = useCase.lbc.GetCollateral(useCase.provider.RskAddress()); err != nil {
		return collateralInfo{}, err
	}
	if pegoutCollateral, err = useCase.lbc.GetPegoutCollateral(useCase.provider.RskAddress()); err != nil {
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
	if operationalForPegin, err = useCase.lbc.IsOperationalPegin(useCase.provider.RskAddress()); err != nil {
		return operationalInfo{}, err
	}

	if operationalForPegout, err = useCase.lbc.IsOperationalPegout(useCase.provider.RskAddress()); err != nil {
		return operationalInfo{}, err
	}

	return operationalInfo{
		operationalForPegin:  operationalForPegin,
		operationalForPegout: operationalForPegout,
	}, nil
}

func (useCase *RegistrationUseCase) isProviderOperational(providerType entities.ProviderType, operational operationalInfo) bool {
	return (providerType == entities.FullProvider && operational.operationalForPegin && operational.operationalForPegout) ||
		(providerType == entities.PeginProvider && operational.operationalForPegin) ||
		(providerType == entities.PegoutProvider && operational.operationalForPegout)
}

func (useCase *RegistrationUseCase) registerProvider(params blockchain.ProviderRegistrationParams, collateral collateralInfo) (int64, error) {
	value := new(entities.Wei)
	txConfig := blockchain.NewTransactionConfig(value.Mul(collateral.minimumCollateral, entities.NewUWei(2)), 0, nil)
	return useCase.lbc.RegisterProvider(txConfig, params)
}

func (useCase *RegistrationUseCase) validateParams(params blockchain.ProviderRegistrationParams) error {
	var err error
	if err = entities.ValidateStruct(params); err != nil {
		return usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	} else if !params.Type.IsValid() {
		return usecases.WrapUseCaseError(usecases.ProviderRegistrationId, entities.InvalidProviderTypeError)
	}
	return nil
}

func (useCase *RegistrationUseCase) addPeginCollateral(
	params blockchain.ProviderRegistrationParams,
	operational operationalInfo,
	collateral collateralInfo,
) (bool, error) {
	var err error
	var collateralToAdd *entities.Wei
	if !(params.Type.AcceptsPegin() && !operational.operationalForPegin && collateral.peginCollateral.Cmp(entities.NewWei(0)) != 0) {
		return false, nil
	}
	if err = useCase.lbc.AddCollateral(collateralToAdd.Sub(collateral.minimumCollateral, collateral.peginCollateral)); err != nil {
		return true, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	} else {
		return true, nil
	}
}

func (useCase *RegistrationUseCase) addPegoutCollateral(
	params blockchain.ProviderRegistrationParams,
	operational operationalInfo,
	collateral collateralInfo,
) (bool, error) {
	var err error
	var collateralToAdd *entities.Wei
	if !(params.Type.AcceptsPegout() && !operational.operationalForPegout && collateral.pegoutCollateral.Cmp(entities.NewWei(0)) != 0) {
		return false, nil
	}
	if err = useCase.lbc.AddPegoutCollateral(collateralToAdd.Sub(collateral.minimumCollateral, collateral.pegoutCollateral)); err != nil {
		return true, usecases.WrapUseCaseError(usecases.ProviderRegistrationId, err)
	} else {
		return true, nil
	}
}
