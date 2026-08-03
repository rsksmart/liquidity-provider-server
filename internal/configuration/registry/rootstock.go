package registry

import (
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	bridgeBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/bridge"
	collateralBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/collateral_management"
	discoveryBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/discovery"
	flyoverConfigurationsBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/flyover_configurations"
	peginBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin"
	peginAddressRegistryBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_address_registry"
	pegoutBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegout"
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
	peginAddressRegistry  *peginAddressRegistryBinding.PegInAddressRegistryContract
	flyoverConfigurations *flyoverConfigurationsBinding.FlyoverConfigurationsContract
}

type rskBoundContracts struct {
	bridge                *bind.BoundContract
	peginContract         *bind.BoundContract
	pegoutContract        *bind.BoundContract
	collateralManagement  *bind.BoundContract
	discovery             *bind.BoundContract
	peginAddressRegistry  *bind.BoundContract
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

	// Optional adapters: only constructed when their bound contract exists (see
	// createBoundContracts' boot-safety branch). Left nil otherwise so the server can boot
	// without these contract addresses configured; callers must nil-check before use.
	var peginAddressRegistry blockchain.PegInAddressRegistryContract
	if boundContracts.peginAddressRegistry != nil {
		peginAddressRegistry = rootstock.NewPegInAddressRegistryContractImpl(
			client,
			env.Rsk.PegInAddressRegistryAddress,
			boundContracts.peginAddressRegistry,
			rootstock.DefaultRetryParams,
			contractBindings.peginAddressRegistry,
			abis,
		)
	}
	var flyoverConfigurations blockchain.FlyoverConfigurationsContract
	if boundContracts.flyoverConfigurations != nil {
		flyoverConfigurations = rootstock.NewFlyoverConfigurationsContractImpl(
			client,
			env.Rsk.FlyoverConfigurationsAddress,
			boundContracts.flyoverConfigurations,
			rootstock.DefaultRetryParams,
			contractBindings.flyoverConfigurations,
			abis,
		)
	}

	return &Rootstock{
		Contracts: blockchain.RskContracts{
			PegInAddressRegistry:  peginAddressRegistry,
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

	// PegInAddressRegistry and FlyoverConfigurations are optional wiring slots: both env vars
	// may be absent, and ParseAddress rejects "" as an invalid address rather than treating it
	// as "skip" — so the bound contract (and downstream adapter) is only constructed when the
	// address is actually configured. Leaving it nil here is safe as long as callers nil-check
	// RskContracts.PegInAddressRegistry/.FlyoverConfigurations before use.
	var peginAddressRegistry, flyoverConfigurations *bind.BoundContract
	if env.Rsk.PegInAddressRegistryAddress != "" {
		var peginAddressRegistryAddress common.Address
		if err = rootstock.ParseAddress(&peginAddressRegistryAddress, env.Rsk.PegInAddressRegistryAddress); err != nil {
			return rskBoundContracts{}, err
		}
		peginAddressRegistry = bindings.peginAddressRegistry.Instance(client.Rpc(), peginAddressRegistryAddress)
	}
	if env.Rsk.FlyoverConfigurationsAddress != "" {
		var flyoverConfigurationsAddress common.Address
		if err = rootstock.ParseAddress(&flyoverConfigurationsAddress, env.Rsk.FlyoverConfigurationsAddress); err != nil {
			return rskBoundContracts{}, err
		}
		flyoverConfigurations = bindings.flyoverConfigurations.Instance(client.Rpc(), flyoverConfigurationsAddress)
	}

	return rskBoundContracts{
		bridge:                bridge,
		peginContract:         peginContract,
		pegoutContract:        pegoutContract,
		collateralManagement:  collateralManagement,
		discovery:             discovery,
		peginAddressRegistry:  peginAddressRegistry,
		flyoverConfigurations: flyoverConfigurations,
	}, nil
}

func createContractBindings() rskContractBindings {
	return rskContractBindings{
		bridge:                bridgeBinding.NewRskBridge(),
		peginContract:         peginBinding.NewPeginContract(),
		pegoutContract:        pegoutBinding.NewPegoutContract(),
		collateralManagement:  collateralBinding.NewCollateralManagementContract(),
		discovery:             discoveryBinding.NewFlyoverDiscovery(),
		peginAddressRegistry:  peginAddressRegistryBinding.NewPegInAddressRegistryContract(),
		flyoverConfigurations: flyoverConfigurationsBinding.NewFlyoverConfigurationsContract(),
	}
}
