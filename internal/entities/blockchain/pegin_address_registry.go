package blockchain

import (
	"context"
	"errors"
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
)

var InvalidRootHashFunctionError = errors.New("invalid PegIn address registry root hash function")

// PegInAddressRegistryEncoding is the on-chain IPegInAddressRegistry.Encoding enum.
type PegInAddressRegistryEncoding = rootstock.PegInAddressRegistryEncoding

const (
	PegInAddressRegistryEncodingBase58  = rootstock.PegInAddressRegistryEncodingBase58
	PegInAddressRegistryEncodingBech32  = rootstock.PegInAddressRegistryEncodingBech32
	PegInAddressRegistryEncodingBech32M = rootstock.PegInAddressRegistryEncodingBech32M
)

func IsSupportedPegInEncoding(encoding PegInAddressRegistryEncoding) bool {
	return rootstock.IsSupportedPegInEncoding(encoding)
}

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

func NewAddressRegisteredFromWatchEntry(watch rootstock.PegInWatch) AddressRegistered {
	return AddressRegistered{
		RskAddress:       watch.RskAddress,
		Registrant:       watch.Registrant,
		RegistrationRoot: watch.RegistrationRoot,
		TxHash:           watch.TxHash,
		BlockNumber:      watch.BlockNumber,
		LogIndex:         watch.LogIndex,
	}
}

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

func FoldPegInAddressRegistryRoot(
	hashFunction entities.HashFunction,
	previousRoot [32]byte,
	rskAddress string,
) ([32]byte, error) {
	if hashFunction == nil {
		return [32]byte{}, fmt.Errorf("%w: hash function is nil", InvalidRootHashFunctionError)
	}
	normalizedAddress, err := NormalizeRskAddress(rskAddress)
	if err != nil {
		return [32]byte{}, err
	}
	addressBytes, err := DecodeStringTrimPrefix(normalizedAddress)
	if err != nil {
		return [32]byte{}, fmt.Errorf("%w: %w", InvalidAddressError, err)
	}

	var preimage [52]byte
	copy(preimage[:32], previousRoot[:])
	copy(preimage[32:], addressBytes)

	hash := hashFunction(preimage[:])
	if len(hash) != 32 {
		return [32]byte{}, fmt.Errorf(
			"%w: hash function returned %d bytes, want 32",
			InvalidRootHashFunctionError,
			len(hash),
		)
	}
	return [32]byte(hash), nil
}
