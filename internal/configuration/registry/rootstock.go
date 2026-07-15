package registry

import (
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	bridgeBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/bridge"
	collateralBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/collateral_management"
	configurationsBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/configurations"
	discoveryBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/discovery"
	peginBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin"
	pegoutBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegout"
	registryBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap/wallet"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

type Rootstock struct {
	Contracts blockchain.RskContracts
	Wallet    rootstock.RskSignerWallet
	Client    *rootstock.RskClient
}

type rskContractBindings struct {
	bridge                *bridgeBinding.RskBridge
	peginContract         *peginBinding.PeginContract
	pegoutContract        *pegoutBinding.PegoutContract
	collateralManagement  *collateralBinding.CollateralManagementContract
	discovery             *discoveryBinding.FlyoverDiscovery
	pegInAddressRegistry  *registryBinding.PegInAddressRegistry
	flyoverConfigurations *configurationsBinding.FlyoverConfigurations
}

type rskBoundContracts struct {
	bridge                *bind.BoundContract
	peginContract         *bind.BoundContract
	pegoutContract        *bind.BoundContract
	collateralManagement  *bind.BoundContract
	discovery             *bind.BoundContract
	pegInAddressRegistry  *bind.BoundContract
	flyoverConfigurations *bind.BoundContract
}

// nolint:funlen
func NewRootstockRegistry(env environment.Environment, client *rootstock.RskClient, walletFactory wallet.AbstractFactory, timeouts environment.ApplicationTimeouts) (*Rootstock, error) {
	contractBindings := createContractBindings()

	boundContracts, err := createBoundContracts(env, contractBindings, client)
	if err != nil {
		return nil, err
	}

	wallet, err := walletFactory.RskWallet()
	if err != nil {
		return nil, err
	}

	btcParams, err := env.Btc.GetNetworkParams()
	if err != nil {
		return nil, err
	}

	abis := rootstock.MustLoadFlyoverABIs()

	// Commit-first peg-in contracts (DoS-removal redesign, EPICs E1/E2/E5). They are optional:
	// only constructed when their addresses are configured, so legacy deployments are unaffected.
	pegInAddressRegistry := buildPegInAddressRegistry(env, contractBindings, boundContracts, client)
	flyoverConfigurations := buildFlyoverConfigurations(env, contractBindings, boundContracts, client)

	return &Rootstock{
		Contracts: blockchain.RskContracts{
			PegInAddressRegistry:  pegInAddressRegistry,
			FlyoverConfigurations: flyoverConfigurations,
			Bridge: rootstock.NewRskBridgeImpl(
				rootstock.RskBridgeConfig{
					Address:               env.Rsk.BridgeAddress,
					RequiredConfirmations: env.Rsk.BridgeRequiredConfirmations,
					ErpKeys:               env.Rsk.ErpKeys,
					UseSegwitFederation:   env.Rsk.UseSegwitFederation,
				},
				boundContracts.bridge,
				client,
				btcParams,
				rootstock.DefaultRetryParams,
				wallet,
				contractBindings.bridge,
				timeouts.MiningWait.Seconds(),
			),
			PegIn: rootstock.NewPeginContractImpl(
				client,
				env.Rsk.PeginContractAddress,
				boundContracts.peginContract,
				wallet,
				rootstock.DefaultRetryParams,
				timeouts.MiningWait.Seconds(),
				contractBindings.peginContract,
				abis,
			),
			PegOut: rootstock.NewPegoutContractImpl(
				client,
				env.Rsk.PegoutContractAddress,
				boundContracts.pegoutContract,
				wallet,
				rootstock.DefaultRetryParams,
				timeouts.MiningWait.Seconds(),
				contractBindings.pegoutContract,
				abis,
			),
			CollateralManagement: rootstock.NewCollateralManagementContractImpl(
				client,
				wallet.Address().String(),
				env.Rsk.CollateralManagementAddress,
				boundContracts.collateralManagement,
				wallet,
				contractBindings.collateralManagement,
				rootstock.DefaultRetryParams,
				timeouts.MiningWait.Seconds(),
				abis,
			),
			Discovery: rootstock.NewDiscoveryContractImpl(
				client,
				env.Rsk.DiscoveryAddress,
				boundContracts.discovery,
				wallet,
				rootstock.DefaultRetryParams,
				timeouts.MiningWait.Seconds(),
				time.Duration(env.Provider.FillWithDefaults().RegistrationPollIntervalSeconds)*time.Second,
				contractBindings.discovery,
				abis,
			),
		},
		Wallet: wallet,
		Client: client,
	}, nil
}

func createBoundContracts(
	env environment.Environment,
	bindings rskContractBindings,
	client *rootstock.RskClient,
) (rskBoundContracts, error) {
	var (
		err                         error
		bridgeAddress               common.Address
		peginContractAddress        common.Address
		pegoutContractAddress       common.Address
		collateralManagementAddress common.Address
		discoveryAddress            common.Address
	)
	if err = rootstock.ParseAddress(&peginContractAddress, env.Rsk.PeginContractAddress); err != nil {
		return rskBoundContracts{}, err
	}
	if err = rootstock.ParseAddress(&pegoutContractAddress, env.Rsk.PegoutContractAddress); err != nil {
		return rskBoundContracts{}, err
	}
	if err = rootstock.ParseAddress(&collateralManagementAddress, env.Rsk.CollateralManagementAddress); err != nil {
		return rskBoundContracts{}, err
	}
	if err = rootstock.ParseAddress(&discoveryAddress, env.Rsk.DiscoveryAddress); err != nil {
		return rskBoundContracts{}, err
	}
	if err = rootstock.ParseAddress(&bridgeAddress, env.Rsk.BridgeAddress); err != nil {
		return rskBoundContracts{}, err
	}

	peginContract := bindings.peginContract.Instance(client.Rpc(), peginContractAddress)
	pegoutContract := bindings.pegoutContract.Instance(client.Rpc(), pegoutContractAddress)
	collateralManagement := bindings.collateralManagement.Instance(client.Rpc(), collateralManagementAddress)
	discovery := bindings.discovery.Instance(client.Rpc(), discoveryAddress)
	bridge := bindings.bridge.Instance(client.Rpc(), bridgeAddress)

	bound := rskBoundContracts{
		bridge:               bridge,
		peginContract:        peginContract,
		pegoutContract:       pegoutContract,
		collateralManagement: collateralManagement,
		discovery:            discovery,
	}

	// Optional commit-first peg-in contracts (EPICs E1/E2/E5).
	if env.Rsk.PegInAddressRegistryAddress != "" {
		var addr common.Address
		if err = rootstock.ParseAddress(&addr, env.Rsk.PegInAddressRegistryAddress); err != nil {
			return rskBoundContracts{}, err
		}
		bound.pegInAddressRegistry = bindings.pegInAddressRegistry.Instance(client.Rpc(), addr)
	}
	if env.Rsk.FlyoverConfigurationsAddress != "" {
		var addr common.Address
		if err = rootstock.ParseAddress(&addr, env.Rsk.FlyoverConfigurationsAddress); err != nil {
			return rskBoundContracts{}, err
		}
		bound.flyoverConfigurations = bindings.flyoverConfigurations.Instance(client.Rpc(), addr)
	}

	return bound, nil
}

func createContractBindings() rskContractBindings {
	return rskContractBindings{
		bridge:                bridgeBinding.NewRskBridge(),
		peginContract:         peginBinding.NewPeginContract(),
		pegoutContract:        pegoutBinding.NewPegoutContract(),
		collateralManagement:  collateralBinding.NewCollateralManagementContract(),
		discovery:             discoveryBinding.NewFlyoverDiscovery(),
		pegInAddressRegistry:  registryBinding.NewPegInAddressRegistry(),
		flyoverConfigurations: configurationsBinding.NewFlyoverConfigurations(),
	}
}

// buildPegInAddressRegistry constructs the registry adapter when its address is configured,
// otherwise returns nil so the rest of the system stays on the legacy path.
func buildPegInAddressRegistry(
	env environment.Environment,
	bindings rskContractBindings,
	bound rskBoundContracts,
	client *rootstock.RskClient,
) blockchain.PegInAddressRegistryContract {
	if bound.pegInAddressRegistry == nil {
		return nil
	}
	return rootstock.NewPegInAddressRegistryImpl(
		client,
		env.Rsk.PegInAddressRegistryAddress,
		bound.pegInAddressRegistry,
		rootstock.DefaultRetryParams,
		bindings.pegInAddressRegistry,
	)
}

// buildFlyoverConfigurations constructs the configurations adapter when its address is configured.
func buildFlyoverConfigurations(
	env environment.Environment,
	bindings rskContractBindings,
	bound rskBoundContracts,
	client *rootstock.RskClient,
) blockchain.FlyoverConfigurationsContract {
	if bound.flyoverConfigurations == nil {
		return nil
	}
	return rootstock.NewFlyoverConfigurationsImpl(
		client,
		env.Rsk.FlyoverConfigurationsAddress,
		bound.flyoverConfigurations,
		rootstock.DefaultRetryParams,
		bindings.flyoverConfigurations,
	)
}
