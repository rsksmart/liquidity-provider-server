package blockchain

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

// PegInAddressRegistryEncoding mirrors the on-chain IPegInAddressRegistry.Encoding enum.
type PegInAddressRegistryEncoding uint8

const (
	PegInAddressRegistryEncodingBase58  PegInAddressRegistryEncoding = 0
	PegInAddressRegistryEncodingBech32  PegInAddressRegistryEncoding = 1
	PegInAddressRegistryEncodingBech32M PegInAddressRegistryEncoding = 2
)

// PegInRegistration is the decoded registration record for a single RSK destination address.
type PegInRegistration struct {
	Registrant        string
	RegistrationBlock *entities.Wei
}

// AddressRegistered is the decoded AddressRegistered event emitted by the PegInAddressRegistry contract.
type AddressRegistered struct {
	RskAddress       string
	Registrant       string
	RegistrationRoot [32]byte
	TxHash           string
	BlockNumber      uint64
}

// PegInAddressRegistryContract is a read-only port over the frozen IPegInAddressRegistry ABI.
// Registration (registerAddress) is intentionally not exposed: it is the watchtower's job
// (FLY-2446), never the LPS's.
type PegInAddressRegistryContract interface {
	GetAddress() string
	GetPegInAddress(rskAddr string) (payload []byte, encoding PegInAddressRegistryEncoding, err error)
	GetPegInAddresses(rskAddrs []string) (payloads [][]byte, encoding PegInAddressRegistryEncoding, err error)
	IsRegistered(rskAddr string) (bool, error)
	GetRegistration(rskAddr string) (PegInRegistration, error)
	GetRegistrationRoot() ([32]byte, error)
	GetAddressRegisteredEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]AddressRegistered, error)
}
