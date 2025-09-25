package rootstock

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
)

type FlyoverABIs struct {
	PegIn                *abi.ABI
	PegOut               *abi.ABI
	Discovery            *abi.ABI
	CollateralManagement *abi.ABI
	DaoContributor       *abi.ABI
	Flyover              *abi.ABI
}

func MustLoadFlyoverABIs() *FlyoverABIs {
	pegInAbi, err := bindings.IPegInMetaData.GetAbi()
	if err != nil {
		panic("could not load PegIn ABI: " + err.Error())
	}
	pegOutAbi, err := bindings.IPegOutMetaData.GetAbi()
	if err != nil {
		panic("could not load PegOut ABI: " + err.Error())
	}
	discoveryAbi, err := bindings.IFlyoverDiscoveryMetaData.GetAbi()
	if err != nil {
		panic("could not load Discovery ABI: " + err.Error())
	}
	collateralManagementAbi, err := bindings.IPegOutMetaData.GetAbi()
	if err != nil {
		panic("could not load Collateral Management ABI: " + err.Error())
	}
	daoContributorAbi, err := bindings.IDaoContributorMetaData.GetAbi()
	if err != nil {
		panic("could not load DAO Contributor ABI: " + err.Error())
	}
	flyoverAbi, err := bindings.FlyoverMetaData.GetAbi()
	if err != nil {
		panic("could not load Flyover ABI: " + err.Error())
	}

	return &FlyoverABIs{
		PegIn:                pegInAbi,
		PegOut:               pegOutAbi,
		Discovery:            discoveryAbi,
		CollateralManagement: collateralManagementAbi,
		DaoContributor:       daoContributorAbi,
		Flyover:              flyoverAbi,
	}
}
