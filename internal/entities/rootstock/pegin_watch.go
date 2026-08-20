package rootstock

import (
	"context"
	"fmt"
	"time"
)

type PegInWatchState string

const (
	PegInWatchDiscovered          PegInWatchState = "discovered"           // stored from AddressRegistered, not yet imported
	PegInWatchImported            PegInWatchState = "imported"             // Bitcoin address imported into the wallet
	PegInWatchUnsupportedEncoding PegInWatchState = "unsupported_encoding" // encoding cannot be converted to a Bitcoin address
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

type PegInWatch struct {
	TxHash           string          `json:"txHash" bson:"tx_hash"`
	LogIndex         uint            `json:"logIndex" bson:"log_index"`
	BlockNumber      uint64          `json:"blockNumber" bson:"block_number"`
	RskAddress       string          `json:"rskAddress" bson:"rsk_address"`
	Registrant       string          `json:"registrant" bson:"registrant"`
	RegistrationRoot [32]byte        `json:"registrationRoot" bson:"registration_root"`
	Encoding         uint8           `json:"encoding" bson:"encoding"`
	BtcAddress       string          `json:"btcAddress" bson:"btc_address"`
	State            PegInWatchState `json:"state" bson:"state"`
	LastSeenAt       *time.Time      `json:"lastSeenAt" bson:"last_seen_at"` // post-import wallet observation
	LastError        string          `json:"lastError" bson:"last_error"`
	CreatedAt        time.Time       `json:"createdAt" bson:"created_at"`
	UpdatedAt        time.Time       `json:"updatedAt" bson:"updated_at"`
}

func NewPegInWatch(
	txHash string,
	logIndex uint,
	blockNumber uint64,
	rskAddress string,
	registrant string,
	registrationRoot [32]byte,
) PegInWatch {
	now := time.Now().UTC()
	return PegInWatch{
		TxHash:           txHash,
		LogIndex:         logIndex,
		BlockNumber:      blockNumber,
		RskAddress:       rskAddress,
		Registrant:       registrant,
		RegistrationRoot: registrationRoot,
		State:            PegInWatchDiscovered,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func (watch *PegInWatch) SetEncoding(encoding PegInAddressRegistryEncoding) {
	watch.Encoding = uint8(encoding)
	watch.UpdatedAt = time.Now().UTC()
	if IsSupportedPegInEncoding(encoding) {
		return
	}
	watch.State = PegInWatchUnsupportedEncoding
	watch.LastError = fmt.Sprintf("unsupported encoding %d", encoding)
}

func (watch *PegInWatch) SameLog(txHash string, logIndex uint) bool {
	return watch.TxHash == txHash && watch.LogIndex == logIndex
}

func (watch *PegInWatch) MarkImported() {
	watch.State = PegInWatchImported
	watch.LastError = ""
	watch.UpdatedAt = time.Now().UTC()
}

func (watch *PegInWatch) RecordError(err error) bool {
	if watch.LastError == err.Error() {
		return false
	}
	watch.LastError = err.Error()
	watch.UpdatedAt = time.Now().UTC()
	return true
}

type PegInWatchCheckpoint struct {
	LocalRoot          [32]byte `json:"localRoot" bson:"local_root"`
	LastProcessedBlock uint64   `json:"lastProcessedBlock" bson:"last_processed_block"`
}

type PegInWatchRepository interface {
	Upsert(ctx context.Context, watch PegInWatch) error
	Get(ctx context.Context, rskAddress string) (*PegInWatch, error)
	List(ctx context.Context) ([]PegInWatch, error)
	Update(ctx context.Context, watch PegInWatch) error
}

type PegInWatchCheckpointRepository interface {
	GetCheckpoint(ctx context.Context) (checkpoint PegInWatchCheckpoint, found bool, err error)
	SetCheckpoint(ctx context.Context, checkpoint PegInWatchCheckpoint) error
	DeleteCheckpoint(ctx context.Context) error
}

// PegInWatchRepositorySet is the Replay port.
// Discover takes only PegInWatchRepository because it never
// reads or writes the checkpoint. Replay needs both the entry methods and the
// checkpoint methods. Keep the split so Discover cannot call DeleteCheckpoint.
type PegInWatchRepositorySet interface {
	PegInWatchRepository
	PegInWatchCheckpointRepository
}

type PegInWatches []*PegInWatch

func (entries PegInWatches) Contains(rskAddress string) bool {
	for _, entry := range entries {
		if entry != nil && entry.RskAddress == rskAddress {
			return true
		}
	}
	return false
}
