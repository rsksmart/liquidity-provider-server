package rootstock

import (
	"context"
	"time"
)

type PegInAddressRegistryWatchState string

const (
	PegInAddressRegistryWatchDiscovered          PegInAddressRegistryWatchState = "discovered"
	PegInAddressRegistryWatchImported            PegInAddressRegistryWatchState = "imported"
	PegInAddressRegistryWatchUnsupportedEncoding PegInAddressRegistryWatchState = "unsupported_encoding"
)

type PegInAddressRegistryWatchEntry struct {
	TxHash           string                         `json:"txHash" bson:"tx_hash"`
	LogIndex         uint                           `json:"logIndex" bson:"log_index"`
	BlockNumber      uint64                         `json:"blockNumber" bson:"block_number"`
	RskAddress       string                         `json:"rskAddress" bson:"rsk_address"`
	Registrant       string                         `json:"registrant" bson:"registrant"`
	RegistrationRoot [32]byte                       `json:"registrationRoot" bson:"registration_root"`
	Encoding         uint8                          `json:"encoding" bson:"encoding"`
	BtcAddress       string                         `json:"btcAddress" bson:"btc_address"`
	State            PegInAddressRegistryWatchState `json:"state" bson:"state"`
	LastError        string                         `json:"lastError" bson:"last_error"`
	CreatedAt        time.Time                      `json:"createdAt" bson:"created_at"`
	UpdatedAt        time.Time                      `json:"updatedAt" bson:"updated_at"`
}

type PegInAddressRegistryWatchCheckpoint struct {
	LocalRoot          [32]byte `json:"localRoot" bson:"local_root"`
	LastProcessedBlock uint64   `json:"lastProcessedBlock" bson:"last_processed_block"`
}

type PegInAddressRegistryWatchRepository interface {
	Upsert(context.Context, PegInAddressRegistryWatchEntry) error
	Get(context.Context, string) (*PegInAddressRegistryWatchEntry, error)
	List(context.Context) ([]PegInAddressRegistryWatchEntry, error)
	Update(context.Context, PegInAddressRegistryWatchEntry) error
}

type PegInAddressRegistryWatchCheckpointRepository interface {
	GetCheckpoint(context.Context) (checkpoint PegInAddressRegistryWatchCheckpoint, found bool, err error)
	SetCheckpoint(context.Context, PegInAddressRegistryWatchCheckpoint) error
	DeleteCheckpoint(context.Context) error
}

type PegInAddressRegistryWatchRepositorySet interface {
	PegInAddressRegistryWatchRepository
	PegInAddressRegistryWatchCheckpointRepository
}
