package blockchain

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
)

const (
	GetPegoutDepositsErrorTemplate = "error executing getting deposits in range [%d, %d]"
)

type RegisterPeginParams struct {
	QuoteSignature        []byte
	BitcoinRawTransaction []byte
	PartialMerkleTree     []byte
	BlockHeight           *big.Int
	Quote                 quote.PeginQuote
}

func (params RegisterPeginParams) String() string {
	return fmt.Sprintf(
		"RegisterPeginParams { QuoteSignature: %s, BitcoinRawTransaction: %s, "+
			"PartialMerkleTree: %s, BlockHeight: %v, Quote: %+v }",
		hex.EncodeToString(params.QuoteSignature),
		hex.EncodeToString(params.BitcoinRawTransaction),
		hex.EncodeToString(params.PartialMerkleTree),
		params.BlockHeight,
		params.Quote,
	)
}

type RefundPegoutParams struct {
	QuoteHash          [32]byte
	BtcRawTx           []byte
	BtcBlockHeaderHash [32]byte
	MerkleBranchPath   *big.Int
	MerkleBranchHashes [][32]byte
}

func (params RefundPegoutParams) String() string {
	return fmt.Sprintf(
		"RefundPegoutParams { QuoteHash: %s, BtcRawTx: %s, "+
			"BtcBlockHeaderHash: %s, MerkleBranchPath: %v, MerkleBranchHashes: %v }",
		hex.EncodeToString(params.QuoteHash[:]),
		hex.EncodeToString(params.BtcRawTx),
		hex.EncodeToString(params.BtcBlockHeaderHash[:]),
		params.MerkleBranchPath,
		params.MerkleBranchHashes,
	)
}

type ProviderRegistrationParams struct {
	Name       string                          `validate:"required"`
	ApiBaseUrl string                          `validate:"required"`
	Status     bool                            `validate:"required"`
	Type       liquidity_provider.ProviderType `validate:"oneof=0 1 2"`
}

func NewProviderRegistrationParams(
	name string,
	apiBaseUrl string,
	status bool,
	providerType liquidity_provider.ProviderType,
) ProviderRegistrationParams {
	return ProviderRegistrationParams{
		Name:       name,
		ApiBaseUrl: apiBaseUrl,
		Status:     status,
		Type:       providerType,
	}
}

type PauseStatus struct {
	IsPaused bool
	Reason   string
	Since    uint64
}

type RegistrationState uint8

const (
	RegistrationStateNone      RegistrationState = 0
	RegistrationStatePending   RegistrationState = 1
	RegistrationStateApproved  RegistrationState = 2
	RegistrationStateRejected  RegistrationState = 3
	RegistrationStateWithdrawn RegistrationState = 4
)

func (s RegistrationState) AllowsRegistration() bool {
	switch s {
	case RegistrationStateNone, RegistrationStateRejected, RegistrationStateWithdrawn:
		return true
	default:
		return false
	}
}

type Pausable interface {
	GetAddress() string
	PausedStatus() (PauseStatus, error)
}

type PeginContract interface {
	Pausable
	GetAddress() string
	GetBalance(address string) (*entities.Wei, error)
	HashPeginQuote(peginQuote quote.PeginQuote) (string, error)
	HashPeginQuoteEIP712(peginQuote quote.PeginQuote) ([32]byte, error)
	CallForUser(txConfig TransactionConfig, peginQuote quote.PeginQuote) (TransactionReceipt, error)
	RegisterPegin(params RegisterPeginParams) (TransactionReceipt, error)
	Withdraw(amount *entities.Wei) error
	// RequestPegIn is the commit-first claim entry point (DoS-removal redesign, EPIC E4/E5).
	// The LP fronts RBTC (msg.value = amount - peg-in fee) to serve a confirmed deposit to a
	// registered RSK address. It is first-mined-wins: if another LP already claimed the same
	// (rskAddr, btcTxHash) the transaction reverts; callers must tolerate AlreadyClaimedError.
	RequestPegIn(params RequestPegInParams) (TransactionReceipt, error)
	// ResolvePegIn settles a previously requested peg-in against the Bridge once the BTC tx has
	// the required confirmations, releasing the RBTC and paying the claiming LP its fee.
	ResolvePegIn(params ResolvePegInParams) (TransactionReceipt, error)
}

// PegInAddressRegistryContract is the port for the new PegInAddressRegistry (EPIC E2).
// It exposes the deterministic BTC deposit address derived from an RSK address and the
// AddressRegistered event stream that drives LPS discovery (EPIC E5).
type PegInAddressRegistryContract interface {
	GetAddress() string
	// IsRegistered reports whether the given RSK address has a registered deposit address.
	IsRegistered(rskAddress string) (bool, error)
	// GetPegInAddress returns the derived BTC deposit address (and its encoding) for a single RSK address.
	GetPegInAddress(rskAddress string) (PegInDepositAddress, error)
	// GetPegInAddresses batch-reads the derived BTC deposit addresses for several RSK addresses.
	GetPegInAddresses(rskAddresses []string) ([]PegInDepositAddress, error)
	// GetRegistrationBlock returns the RSK block at which the address was registered.
	GetRegistrationBlock(rskAddress string) (uint64, error)
	// GetRegistrationRoot returns the running registration-root hash.
	GetRegistrationRoot() ([32]byte, error)
	// GetAddressRegisteredEvents polls the AddressRegistered logs in [fromBlock, toBlock].
	GetAddressRegisteredEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]AddressRegisteredEvent, error)
}

// FlyoverConfigurationsContract is the port for the new FlyoverConfigurations contract (EPIC E1).
// Fees and required confirmations are read from on-chain configuration, not negotiated quotes.
type FlyoverConfigurationsContract interface {
	GetAddress() string
	// CalculatePegInFee returns the protocol peg-in fee for a given peg-in amount.
	CalculatePegInFee(amount *entities.Wei) (*entities.Wei, error)
	// GetRequiredPegInConfirmations returns the BTC confirmations required for the given amount.
	GetRequiredPegInConfirmations(amount *entities.Wei) (uint64, error)
}

// RequestPegInParams are the arguments to PeginContract.RequestPegIn.
type RequestPegInParams struct {
	RskAddress string
	Amount     *entities.Wei
	BtcTxHash  [32]byte
	OpReturn   []byte
	// Value is the RBTC the LP fronts (amount - peg-in fee), set as msg.value.
	Value *entities.Wei
}

// ResolvePegInParams are the arguments to PeginContract.ResolvePegIn.
type ResolvePegInParams struct {
	RskAddress        string
	BtcTxHash         [32]byte
	BtcRawTransaction []byte
	PartialMerkleTree []byte
	Height            *big.Int
	Registrant        string
}

// PegInAddressEncoding mirrors IPegInAddressRegistry.Encoding.
type PegInAddressEncoding uint8

// PegInDepositAddress is a derived BTC deposit address read from the registry.
type PegInDepositAddress struct {
	RskAddress string
	// ScriptOrAddress is the raw bytes the registry returns for the derivation address.
	ScriptOrAddress []byte
	Encoding        PegInAddressEncoding
}

// AddressRegisteredEvent is a decoded AddressRegistered log.
type AddressRegisteredEvent struct {
	RskAddress       string
	RegistrationRoot [32]byte
	BlockNumber      uint64
	TxHash           string
}

type PegoutContract interface {
	Pausable
	GetAddress() string
	HashPegoutQuote(pegoutQuote quote.PegoutQuote) (string, error)
	HashPegoutQuoteEIP712(pegoutQuote quote.PegoutQuote) ([32]byte, error)
	RefundUserPegOut(quoteHash string) (string, error)
	IsPegOutQuoteCompleted(quoteHash string) (bool, error)
	GetDepositEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]quote.PegoutDeposit, error)
	RefundPegout(txConfig TransactionConfig, params RefundPegoutParams) (TransactionReceipt, error)
	ValidatePegout(quoteHash string, btcTx []byte) error
}

type DiscoveryContract interface {
	Pausable
	GetAddress() string
	SetProviderStatus(id uint64, newStatus bool) error
	UpdateProvider(name, url string) (string, error)
	RegisterProvider(txConfig TransactionConfig, params ProviderRegistrationParams) (int64, error)
	GetProviders() ([]liquidity_provider.RegisteredLiquidityProvider, error)
	GetProvider(address string) (liquidity_provider.RegisteredLiquidityProvider, error)
	IsOperational(providerType liquidity_provider.ProviderType, address string) (bool, error)
	GetRegistrationState(address string) (RegistrationState, error)
	WatchRegistrationApproval(ctx context.Context, address string) (RegistrationState, error)
}

type CollateralManagementContract interface {
	Pausable
	GetAddress() string
	GetPenalizedEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]penalization.PenalizedEvent, error)
	ProviderResign() error
	WithdrawCollateral() error
	AddCollateral(amount *entities.Wei) error
	AddPegoutCollateral(amount *entities.Wei) error
	GetCollateral(address string) (*entities.Wei, error)
	GetPegoutCollateral(address string) (*entities.Wei, error)
	GetMinimumCollateral() (*entities.Wei, error)
}
