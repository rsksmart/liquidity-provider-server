package rootstock

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	collateral "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/collateral_management"
	configurations "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/configurations"
	discovery "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/discovery"
	flyover "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/flyover"
	pegin "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin"
	pegout "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegout"
	registry "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/registry"
)

type FlyoverABIs struct {
	PegIn                *abi.ABI
	PegOut               *abi.ABI
	Discovery            *abi.ABI
	CollateralManagement *abi.ABI
	Flyover              *abi.ABI
	// Commit-first peg-in contracts (DoS-removal redesign, EPICs E1/E2).
	PegInAddressRegistry  *abi.ABI
	FlyoverConfigurations *abi.ABI
}

func MustLoadFlyoverABIs() *FlyoverABIs {
	pegInAbi, err := pegin.PeginContractMetaData.ParseABI()
	if err != nil {
		panic("could not load PegIn ABI: " + err.Error())
	}
	pegOutAbi, err := pegout.PegoutContractMetaData.ParseABI()
	if err != nil {
		panic("could not load PegOut ABI: " + err.Error())
	}
	discoveryAbi, err := discovery.FlyoverDiscoveryMetaData.ParseABI()
	if err != nil {
		panic("could not load Discovery ABI: " + err.Error())
	}
	collateralManagementAbi, err := collateral.CollateralManagementContractMetaData.ParseABI()
	if err != nil {
		panic("could not load Collateral Management ABI: " + err.Error())
	}
	flyoverAbi, err := flyover.FlyoverMetaData.ParseABI()
	if err != nil {
		panic("could not load Flyover ABI: " + err.Error())
	}
	registryAbi, err := registry.PegInAddressRegistryMetaData.ParseABI()
	if err != nil {
		panic("could not load PegInAddressRegistry ABI: " + err.Error())
	}
	configurationsAbi, err := configurations.FlyoverConfigurationsMetaData.ParseABI()
	if err != nil {
		panic("could not load FlyoverConfigurations ABI: " + err.Error())
	}

	return &FlyoverABIs{
		PegIn:                 pegInAbi,
		PegOut:                pegOutAbi,
		Discovery:             discoveryAbi,
		CollateralManagement:  collateralManagementAbi,
		Flyover:               flyoverAbi,
		PegInAddressRegistry:  registryAbi,
		FlyoverConfigurations: configurationsAbi,
	}
}
