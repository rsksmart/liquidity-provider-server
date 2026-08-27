package blockchain

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
)

// PegInAddressRegistryEncoding mirrors the on-chain IPegInAddressRegistry.Encoding enum.
type PegInAddressRegistryEncoding uint8

const (
	PegInAddressRegistryEncodingBase58 PegInAddressRegistryEncoding = iota
	PegInAddressRegistryEncodingBech32
	PegInAddressRegistryEncodingBech32M
)

func IsSupportedPegInEncoding(encoding PegInAddressRegistryEncoding) bool {
	return encoding == PegInAddressRegistryEncodingBase58
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
	GetRegistrationRoot() ([32]byte, error)
	// GetAddressRegisteredEvents returns the AddressRegistered events in [fromBlock, toBlock].
	// A nil toBlock reads up to the latest block.
	GetAddressRegisteredEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]AddressRegistered, error)
}
