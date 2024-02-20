package registry

import (
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
)

type UseCaseRegistry interface {
	GetPeginQuoteUseCase() *pegin.GetQuoteUseCase
	GetAcceptPeginQuoteUseCase() *pegin.AcceptQuoteUseCase
	GetProviderDetailUseCase() *liquidity_provider.GetDetailUseCase
	GetPegoutQuoteUseCase() *pegout.GetQuoteUseCase
	GetAcceptPegoutQuoteUseCase() *pegout.AcceptQuoteUseCase
	GetUserDepositsUseCase() *pegout.GetUserDepositsUseCase
	GetProvidersUseCase() *liquidity_provider.GetProvidersUseCase
	GetPeginCollateralUseCase() *pegin.GetCollateralUseCase
	GetPegoutCollateralUseCase() *pegout.GetCollateralUseCase
	WithdrawPeginCollateralUseCase() *pegin.WithdrawCollateralUseCase
	WithdrawPegoutCollateralUseCase() *pegout.WithdrawCollateralUseCase
	HealthUseCase() *usecases.HealthUseCase
	ResignationUseCase() *liquidity_provider.ResignUseCase
	ChangeStatusUseCase() *liquidity_provider.ChangeStatusUseCase
	AddPeginCollateralUseCase() *pegin.AddCollateralUseCase
	AddPegoutCollateralUseCase() *pegout.AddCollateralUseCase
	SetPeginConfigUseCase() *liquidity_provider.SetPeginConfigUseCase
	SetPegoutConfigUseCase() *liquidity_provider.SetPegoutConfigUseCase
	SetGeneralConfigUseCase() *liquidity_provider.SetGeneralConfigUseCase
	GetConfigurationUseCase() *liquidity_provider.GetConfigUseCase
}
