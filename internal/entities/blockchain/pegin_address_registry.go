package blockchain

import "context"

// PegInAddressRegistryEncoding mirrors the on-chain IPegInAddressRegistry.Encoding enum.
type PegInAddressRegistryEncoding uint8

const (
	PegInAddressRegistryEncodingBase58 PegInAddressRegistryEncoding = iota
	PegInAddressRegistryEncodingBech32
	PegInAddressRegistryEncodingBech32M
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
