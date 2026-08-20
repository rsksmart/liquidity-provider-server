package rootstock

import (
	"context"
	"time"
)

type PegInAddressRegistryWatchState string

const (
	PegInAddressRegistryWatchDiscovered          PegInAddressRegistryWatchState = "discovered"           // stored from AddressRegistered, not yet imported
	PegInAddressRegistryWatchImported            PegInAddressRegistryWatchState = "imported"             // Bitcoin address imported into the wallet
	PegInAddressRegistryWatchUnsupportedEncoding PegInAddressRegistryWatchState = "unsupported_encoding" // encoding cannot be converted to a Bitcoin address
)

type PegInAddressRegistryWatch struct {
	TxHash           string                         `json:"txHash" bson:"tx_hash"`
	LogIndex         uint                           `json:"logIndex" bson:"log_index"`
	BlockNumber      uint64                         `json:"blockNumber" bson:"block_number"`
	RskAddress       string                         `json:"rskAddress" bson:"rsk_address"`
	Registrant       string                         `json:"registrant" bson:"registrant"`
	RegistrationRoot [32]byte                       `json:"registrationRoot" bson:"registration_root"`
	Encoding         uint8                          `json:"encoding" bson:"encoding"`
	BtcAddress       string                         `json:"btcAddress" bson:"btc_address"`
	State            PegInAddressRegistryWatchState `json:"state" bson:"state"`
	LastSeenAt       *time.Time                     `json:"lastSeenAt" bson:"last_seen_at"` // post-import wallet observation
	LastError        string                         `json:"lastError" bson:"last_error"`
	CreatedAt        time.Time                      `json:"createdAt" bson:"created_at"`
	UpdatedAt        time.Time                      `json:"updatedAt" bson:"updated_at"`
}

// PegInAddressRegistryWatchEntry is the FLY-2515 name for the watch document.
// Replay/Finalize/helpers use this alias; mongo and 2514 tests keep Watch.
type PegInAddressRegistryWatchEntry = PegInAddressRegistryWatch

func NewPegInAddressRegistryWatch(
	txHash string,
	logIndex uint,
	blockNumber uint64,
	rskAddress string,
	registrant string,
	registrationRoot [32]byte,
) PegInAddressRegistryWatch {
	now := time.Now().UTC()
	return PegInAddressRegistryWatch{
		TxHash:           txHash,
		LogIndex:         logIndex,
		BlockNumber:      blockNumber,
		RskAddress:       rskAddress,
		Registrant:       registrant,
		RegistrationRoot: registrationRoot,
		State:            PegInAddressRegistryWatchDiscovered,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func (watch *PegInAddressRegistryWatch) MarkUnsupportedEncoding(encoding uint8) {
	watch.Encoding = encoding
	watch.State = PegInAddressRegistryWatchUnsupportedEncoding
	watch.LastError = ""
	watch.UpdatedAt = time.Now().UTC()
}

func (watch *PegInAddressRegistryWatch) MarkImported() {
	watch.State = PegInAddressRegistryWatchImported
	watch.LastError = ""
	watch.UpdatedAt = time.Now().UTC()
}

func (watch *PegInAddressRegistryWatch) RecordError(err error) bool {
	if watch.LastError == err.Error() {
		return false
	}
	watch.LastError = err.Error()
	watch.UpdatedAt = time.Now().UTC()
	return true
}

type PegInAddressRegistryWatchCheckpoint struct {
	LocalRoot          [32]byte `json:"localRoot" bson:"local_root"`
	LastProcessedBlock uint64   `json:"lastProcessedBlock" bson:"last_processed_block"`
}

type PegInAddressRegistryWatchRepository interface {
	Upsert(ctx context.Context, watch PegInAddressRegistryWatch) error
	Get(ctx context.Context, rskAddress string) (*PegInAddressRegistryWatch, error)
	List(ctx context.Context) ([]PegInAddressRegistryWatch, error)
	Update(ctx context.Context, watch PegInAddressRegistryWatch) error
}

type PegInAddressRegistryWatchCheckpointRepository interface {
	GetCheckpoint(ctx context.Context) (checkpoint PegInAddressRegistryWatchCheckpoint, found bool, err error)
	SetCheckpoint(ctx context.Context, checkpoint PegInAddressRegistryWatchCheckpoint) error
	DeleteCheckpoint(ctx context.Context) error
}

// PegInAddressRegistryWatchRepositorySet is the Replay port.
// Discover takes only PegInAddressRegistryWatchRepository because it never
// reads or writes the checkpoint. Replay needs both the entry methods and the
// checkpoint methods. Keep the split so Discover cannot call DeleteCheckpoint.
type PegInAddressRegistryWatchRepositorySet interface {
	PegInAddressRegistryWatchRepository
	PegInAddressRegistryWatchCheckpointRepository
}

type PegInAddressRegistryWatchEntries []*PegInAddressRegistryWatchEntry

func (entries PegInAddressRegistryWatchEntries) Contains(rskAddress string) bool {
	for _, entry := range entries {
		if entry != nil && entry.RskAddress == rskAddress {
			return true
		}
	}
	return false
}
