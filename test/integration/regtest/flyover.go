package regtest

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	discoveryBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/discovery"
	configurationsBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/flyover_configurations"
	peginBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin"
	registryBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_address_registry"
	commitFirstBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_commit_first"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/stretchr/testify/require"
)

// Functions the generated bindings lack or decode with an outdated struct: the registry binding
// has no getMinDepositSats, and the configurations binding predates registrantFee in
// PegConfiguration.
const missingFunctionsABI = `[
	{"type":"function","name":"getMinDepositSats","inputs":[],"stateMutability":"view",
		"outputs":[{"name":"minDepositSats","type":"uint256"}]},
	{"type":"function","name":"getPegInConfiguration","inputs":[],"stateMutability":"view",
		"outputs":[{"name":"configuration","type":"tuple","components":[
			{"name":"fixedFee","type":"uint256"},{"name":"percentageFee","type":"uint256"},
			{"name":"minAmount","type":"uint256"},{"name":"maxAmount","type":"uint256"},
			{"name":"registrantFee","type":"uint256"},
			{"name":"confirmationTiers","type":"tuple[]","components":[
				{"name":"maxAmount","type":"uint256"},{"name":"confirmations","type":"uint256"}]}]}]}
]`

// PegInConfiguration mirrors IFlyoverConfigurations.PegConfiguration as deployed. The json tags
// match the type go-ethereum decodes the tuple into, so abi.ConvertType can convert it.
type PegInConfiguration struct {
	FixedFee          *big.Int `json:"fixedFee"`
	PercentageFee     *big.Int `json:"percentageFee"`
	MinAmount         *big.Int `json:"minAmount"`
	MaxAmount         *big.Int `json:"maxAmount"`
	RegistrantFee     *big.Int `json:"registrantFee"`
	ConfirmationTiers []struct {
		MaxAmount     *big.Int `json:"maxAmount"`
		Confirmations *big.Int `json:"confirmations"`
	} `json:"confirmationTiers"`
}

type (
	AddressRegistered = registryBindings.PegInAddressRegistryContractAddressRegistered
	PegInRequested    = commitFirstBindings.PeginCommitFirstContractPegInRequested
	PegInResolved     = commitFirstBindings.PeginCommitFirstContractPegInResolved
)

// Flyover reads and calls the Flyover contracts the commit-first peg-in goes through.
type Flyover struct {
	rsk            *Rootstock
	contracts      Contracts
	registry       *registryBindings.PegInAddressRegistryContract
	commitFirst    *commitFirstBindings.PeginCommitFirstContract
	pegIn          *peginBindings.PeginContract
	configurations *configurationsBindings.FlyoverConfigurationsContract
	discovery      *discoveryBindings.FlyoverDiscovery
	missing        abi.ABI
}

func NewFlyover(t *testing.T, stack Stack, rsk *Rootstock) *Flyover {
	t.Helper()
	missing, err := abi.JSON(strings.NewReader(missingFunctionsABI))
	require.NoError(t, err)
	return &Flyover{
		rsk:            rsk,
		contracts:      stack.Contracts,
		registry:       registryBindings.NewPegInAddressRegistryContract(),
		commitFirst:    commitFirstBindings.NewPeginCommitFirstContract(),
		pegIn:          peginBindings.NewPeginContract(),
		configurations: configurationsBindings.NewFlyoverConfigurationsContract(),
		discovery:      discoveryBindings.NewFlyoverDiscovery(),
		missing:        missing,
	}
}

// PegInAddress is the BTC deposit address the registry derives for rskAddr, decoded the way the
// LPS decodes it.
func (flyover *Flyover) PegInAddress(t *testing.T, rskAddr common.Address) string {
	t.Helper()
	output, err := flyover.registry.UnpackGetPegInAddress(
		flyover.call(t, flyover.contracts.PegInAddressRegistry, flyover.registry.PackGetPegInAddress(rskAddr)))
	require.NoError(t, err)
	require.Equal(t, uint8(rootstock.PegInAddressRegistryEncodingBase58), output.Encoding, "the registry only produces BASE58")
	address, err := bitcoin.EncodeAddressBase58(output.Payload)
	require.NoError(t, err)
	return address
}

func (flyover *Flyover) MinDepositSats(t *testing.T) *big.Int {
	t.Helper()
	output := flyover.callMissing(t, flyover.contracts.PegInAddressRegistry, "getMinDepositSats")
	minDepositSats, ok := abi.ConvertType(output, new(big.Int)).(*big.Int)
	require.True(t, ok, "getMinDepositSats did not return a uint256")
	return minDepositSats
}

func (flyover *Flyover) IsRegistered(t *testing.T, rskAddr common.Address) bool {
	t.Helper()
	registered, err := flyover.registry.UnpackIsRegistered(
		flyover.call(t, flyover.contracts.PegInAddressRegistry, flyover.registry.PackIsRegistered(rskAddr)))
	require.NoError(t, err)
	return registered
}

// RegisterAddress registers rskAddr with registrant as the msg.sender, proving the deposit.
func (flyover *Flyover) RegisterAddress(t *testing.T, registrant Account, rskAddr common.Address, proof Proof) *types.Receipt {
	t.Helper()
	data := flyover.registry.PackRegisterAddress(rskAddr, proof.RawTx, proof.BlockHash, proof.Branch.Path, proof.Branch.Hashes)
	return flyover.rsk.Transact(t, registrant, flyover.contracts.PegInAddressRegistry, data, nil)
}

func (flyover *Flyover) PegInConfiguration(t *testing.T) PegInConfiguration {
	t.Helper()
	output := flyover.callMissing(t, flyover.contracts.FlyoverConfigurations, "getPegInConfiguration")
	configuration, ok := abi.ConvertType(output, new(PegInConfiguration)).(*PegInConfiguration)
	require.True(t, ok, "getPegInConfiguration did not return a PegConfiguration")
	return *configuration
}

// PegInFee is the LP fee the configurations contract charges for amountWei.
func (flyover *Flyover) PegInFee(t *testing.T, amountWei *big.Int) *big.Int {
	t.Helper()
	fee, err := flyover.configurations.UnpackCalculatePegInFee(flyover.call(t,
		flyover.contracts.FlyoverConfigurations, flyover.configurations.PackCalculatePegInFee(amountWei)))
	require.NoError(t, err)
	return fee
}

// RequiredPegInConfirmations is how many BTC confirmations requestPegIn needs for amountWei.
func (flyover *Flyover) RequiredPegInConfirmations(t *testing.T, amountWei *big.Int) uint64 {
	t.Helper()
	confirmations, err := flyover.configurations.UnpackGetRequiredPegInBtcConfirmations(flyover.call(t,
		flyover.contracts.FlyoverConfigurations, flyover.configurations.PackGetRequiredPegInBtcConfirmations(amountWei)))
	require.NoError(t, err)
	return confirmations.Uint64()
}

// Balance is what PegInContract holds for address, such as the LP's repaid advances.
func (flyover *Flyover) Balance(t *testing.T, address common.Address) *big.Int {
	t.Helper()
	balance, err := flyover.pegIn.UnpackGetBalance(flyover.call(t, flyover.contracts.PegIn, flyover.pegIn.PackGetBalance(address)))
	require.NoError(t, err)
	return balance
}

func (flyover *Flyover) IsOperationalPegInProvider(t *testing.T, provider common.Address) bool {
	t.Helper()
	operational, err := flyover.discovery.UnpackIsOperational(flyover.call(t, flyover.contracts.Discovery,
		flyover.discovery.PackIsOperational(uint8(liquidity_provider.PeginProvider), provider)))
	require.NoError(t, err)
	return operational
}

// AddressRegisteredEvents lists the registrations of rskAddr.
func (flyover *Flyover) AddressRegisteredEvents(t *testing.T, rskAddr common.Address) []*AddressRegistered {
	t.Helper()
	logs := flyover.logs(t, flyover.contracts.PegInAddressRegistry,
		&registryBindings.PegInAddressRegistryContractMetaData, registryBindings.PegInAddressRegistryContractAddressRegisteredEventName, 1, common.BytesToHash(rskAddr.Bytes()))
	events := make([]*AddressRegistered, 0, len(logs))
	for i := range logs {
		event, err := flyover.registry.UnpackAddressRegisteredEvent(&logs[i])
		require.NoError(t, err)
		events = append(events, event)
	}
	return events
}

// PegInRequestedEvents lists the claims sent for peg-ins to rskAddr.
func (flyover *Flyover) PegInRequestedEvents(t *testing.T, rskAddr common.Address) []*PegInRequested {
	t.Helper()
	logs := flyover.logs(t, flyover.contracts.PegIn,
		&commitFirstBindings.PeginCommitFirstContractMetaData, commitFirstBindings.PeginCommitFirstContractPegInRequestedEventName, 3, common.BytesToHash(rskAddr.Bytes()))
	events := make([]*PegInRequested, 0, len(logs))
	for i := range logs {
		event, err := flyover.commitFirst.UnpackPegInRequestedEvent(&logs[i])
		require.NoError(t, err)
		events = append(events, event)
	}
	return events
}

// PegInResolvedEvents lists the settlements of pegInID.
func (flyover *Flyover) PegInResolvedEvents(t *testing.T, pegInID [32]byte) []*PegInResolved {
	t.Helper()
	logs := flyover.logs(t, flyover.contracts.PegIn,
		&commitFirstBindings.PeginCommitFirstContractMetaData, commitFirstBindings.PeginCommitFirstContractPegInResolvedEventName, 1, pegInID)
	events := make([]*PegInResolved, 0, len(logs))
	for i := range logs {
		event, err := flyover.commitFirst.UnpackPegInResolvedEvent(&logs[i])
		require.NoError(t, err)
		events = append(events, event)
	}
	return events
}

// callMissing calls a function from missingFunctionsABI and returns its single decoded output.
func (flyover *Flyover) callMissing(t *testing.T, to common.Address, function string) any {
	t.Helper()
	data, err := flyover.missing.Pack(function)
	require.NoError(t, err)
	values, err := flyover.missing.Unpack(function, flyover.call(t, to, data))
	require.NoError(t, err)
	require.Len(t, values, 1)
	return values[0]
}

func (flyover *Flyover) call(t *testing.T, to common.Address, data []byte) []byte {
	t.Helper()
	output, err := flyover.rsk.Client().CallContract(t.Context(), ethereum.CallMsg{To: &to, Data: data}, nil)
	require.NoErrorf(t, err, "the call to %s failed", to)
	return output
}

// logs filters the whole chain by one indexed topic. Regtest chains are short enough for that.
func (flyover *Flyover) logs(t *testing.T, address common.Address, metadata *bind.MetaData, event string, topicIndex int, topic common.Hash) []types.Log {
	t.Helper()
	parsed, err := metadata.ParseABI()
	require.NoError(t, err)
	definition, found := parsed.Events[event]
	require.Truef(t, found, "the binding has no %s event", event)
	topics := make([][]common.Hash, topicIndex+1)
	topics[0] = []common.Hash{definition.ID}
	topics[topicIndex] = []common.Hash{topic}
	logs, err := flyover.rsk.Client().FilterLogs(t.Context(), ethereum.FilterQuery{
		FromBlock: new(big.Int), Addresses: []common.Address{address}, Topics: topics,
	})
	require.NoError(t, err)
	return logs
}
