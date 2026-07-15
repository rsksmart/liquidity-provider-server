// Package pegin holds entities for the commit-first peg-in flow introduced by the liquidity-DoS
// removal redesign (EPICs E2/E4/E5). In this flow the user commits BTC first to a deterministic
// deposit address derived from their RSK address; LPs discover that commitment on-chain and compete
// to serve it. No liquidity is reserved before the user commits.
package pegin

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

// WatchedAddressState is the lifecycle of a registered RSK address tracked by LPS discovery.
type WatchedAddressState string

const (
	// WatchedAddressStateWaitingForDeposit means the address is registered and LPS is monitoring
	// its BTC deposit address for an incoming, confirmed deposit.
	WatchedAddressStateWaitingForDeposit WatchedAddressState = "WaitingForDeposit"
	// WatchedAddressStateRequested means LPS has fronted RBTC via requestPegIn and is waiting for
	// enough BTC confirmations to settle.
	WatchedAddressStateRequested WatchedAddressState = "Requested"
	// WatchedAddressStateResolved means resolvePegIn settled the peg-in successfully.
	WatchedAddressStateResolved WatchedAddressState = "Resolved"
	// WatchedAddressStateClaimedByOther means another LP won the first-mined-wins race.
	WatchedAddressStateClaimedByOther WatchedAddressState = "ClaimedByOther"
)

// WatchedRegisteredAddress is a registered RSK address LPS is watching for a commit-first peg-in,
// together with the derived BTC deposit address and any in-progress claim state.
type WatchedRegisteredAddress struct {
	// RskAddress is the user's RSK address that derived the deposit address.
	RskAddress string
	// DerivationAddress is the BTC deposit address derived for RskAddress, encoded as a string
	// suitable for the Bitcoin wallet (e.g. a base58/bech32 address).
	DerivationAddress string
	// RegistrationBlock is the RSK block where the address was registered (discovery checkpoint).
	RegistrationBlock uint64
	// State is the current claim lifecycle state.
	State WatchedAddressState
	// DepositTxHash is the BTC tx hash of the detected deposit, once seen.
	DepositTxHash string
}

// ConfirmedDeposit is a confirmed BTC deposit to a watched derivation address that is ready to be
// claimed via requestPegIn.
type ConfirmedDeposit struct {
	RskAddress    string
	BtcTxHash     string
	BtcTxHashRaw  [32]byte
	Amount        *entities.Wei
	Confirmations uint64
	// OpReturn carries the optional smart-contract-call payload parsed from the deposit's
	// OP_RETURN output (EPIC E4.2); empty for a plain peg-in.
	OpReturn []byte
}
