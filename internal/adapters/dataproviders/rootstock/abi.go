package rootstock

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	collateral "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/collateral_management"
	discovery "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/discovery"
	flyover "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/flyover"
	pegin "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin"
	pegout "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegout"
)

type FlyoverABIs struct {
	PegIn                *abi.ABI
	PegOut               *abi.ABI
	Discovery            *abi.ABI
	CollateralManagement *abi.ABI
	Flyover              *abi.ABI
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

	return &FlyoverABIs{
		PegIn:                pegInAbi,
		PegOut:               pegOutAbi,
		Discovery:            discoveryAbi,
		CollateralManagement: collateralManagementAbi,
		Flyover:              flyoverAbi,
	}
}
