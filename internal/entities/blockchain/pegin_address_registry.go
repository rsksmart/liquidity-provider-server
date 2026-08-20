package blockchain

import (
	"context"
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"golang.org/x/crypto/sha3"
)

// PegInAddressRegistryEncoding mirrors the on-chain IPegInAddressRegistry.Encoding enum.
type PegInAddressRegistryEncoding uint8

const (
	PegInAddressRegistryEncodingBase58 PegInAddressRegistryEncoding = iota
	PegInAddressRegistryEncodingBech32
	PegInAddressRegistryEncodingBech32M

	PegInAddressRegistryRootMismatchEventId  entities.EventId = "PegInAddressRegistryRootMismatch"
	PegInAddressRegistryResyncStartedEventId entities.EventId = "PegInAddressRegistryResyncStarted"
)

// PegInAddress is a registered BTC address payload together with the encoding needed to read it.
type PegInAddress struct {
	Payload  []byte
	Encoding PegInAddressRegistryEncoding
}

// PegInAddressBatch is the result of a batch lookup. The registry returns a single encoding
// shared by every payload in the batch.
type PegInAddressBatch struct {
	Payloads [][]byte
	Encoding PegInAddressRegistryEncoding
}

// PegInRegistration is the decoded registration record for a single RSK destination address.
type PegInRegistration struct {
	Registrant        string
	RegistrationBlock uint64
}

// AddressRegistered is the decoded AddressRegistered event of the PegInAddressRegistry contract.
type AddressRegistered struct {
	RskAddress       string
	Registrant       string
	RegistrationRoot [32]byte
	TxHash           string
	BlockNumber      uint64
	LogIndex         uint
}

func NewAddressRegisteredFromWatchEntry(entry rootstock.PegInAddressRegistryWatchEntry) AddressRegistered {
	return AddressRegistered{
		RskAddress:       entry.RskAddress,
		Registrant:       entry.Registrant,
		RegistrationRoot: entry.RegistrationRoot,
		TxHash:           entry.TxHash,
		BlockNumber:      entry.BlockNumber,
		LogIndex:         entry.LogIndex,
	}
}

type PegInAddressRegistryRootMismatchEvent struct {
	entities.BaseEvent
	BlockNumber uint64
	LocalRoot   [32]byte
	ChainRoot   [32]byte
}

type PegInAddressRegistryResyncStartedEvent struct {
	entities.BaseEvent
	Reason string
}

// FoldPegInAddressRegistryRoot mirrors
// keccak256(abi.encodePacked(previousRoot, rskAddress)): exactly 32 root bytes
// followed by the unpadded 20-byte RSK address.
func FoldPegInAddressRegistryRoot(previousRoot [32]byte, rskAddress string) ([32]byte, error) {
	normalizedAddress, err := NormalizeRskAddress(rskAddress)
	if err != nil {
		return [32]byte{}, err
	}
	addressBytes, err := DecodeStringTrimPrefix(normalizedAddress)
	if err != nil {
		return [32]byte{}, fmt.Errorf("%w: %w", InvalidAddressError, err)
	}
	if len(addressBytes) != 20 {
		return [32]byte{}, fmt.Errorf("%w: decoded RSK address length %d", InvalidAddressError, len(addressBytes))
	}

	var preimage [52]byte
	copy(preimage[:32], previousRoot[:])
	copy(preimage[32:], addressBytes)

	hasher := sha3.NewLegacyKeccak256()
	if _, err = hasher.Write(preimage[:]); err != nil {
		return [32]byte{}, fmt.Errorf("fold PegIn address registry root: %w", err)
	}
	var root [32]byte
	copy(root[:], hasher.Sum(nil))
	return root, nil
}

// PegInAddressRegistryContract is a read-only port over the frozen IPegInAddressRegistry ABI.
// Registration (registerAddress) is intentionally not exposed: writing registrations is the
// responsibility of a separate on-chain watcher process, not the liquidity provider server.
type PegInAddressRegistryContract interface {
	GetAddress() string
	GetPegInAddress(rskAddr string) (PegInAddress, error)
	GetPegInAddresses(rskAddrs []string) (PegInAddressBatch, error)
	IsRegistered(rskAddr string) (bool, error)
	GetRegistration(rskAddr string) (PegInRegistration, error)
	GetRegistrationRoot(ctx context.Context, blockNumber uint64) ([32]byte, error)
	// GetAddressRegisteredEvents returns the AddressRegistered events in [fromBlock, toBlock].
	// A nil toBlock reads up to the latest block.
	GetAddressRegisteredEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]AddressRegistered, error)
}
