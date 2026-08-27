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
	DepositTxID      string          `json:"depositTxId" bson:"deposit_txid"`
	Confirmations    uint64          `json:"confirmations" bson:"confirmations"`
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

// pegInWatchEncodingBase58 matches blockchain.PegInAddressRegistryEncodingBase58 (iota 0).
// rootstock cannot import blockchain without an import cycle.
const pegInWatchEncodingBase58 uint8 = 0

func (watch *PegInWatch) SetEncoding(encoding uint8) {
	watch.Encoding = encoding
	watch.UpdatedAt = time.Now().UTC()
	if encoding == pegInWatchEncodingBase58 {
		return
	}
	watch.State = PegInWatchUnsupportedEncoding
	watch.LastError = fmt.Sprintf("unsupported encoding %d", encoding)
}

func (watch PegInWatch) SameLog(txHash string, logIndex uint) bool {
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

type PegInWatchRepository interface {
	Upsert(ctx context.Context, watch PegInWatch) error
	Get(ctx context.Context, txHash string, logIndex uint) (*PegInWatch, error)
	List(ctx context.Context) ([]PegInWatch, error)
	Update(ctx context.Context, watch PegInWatch) error
	GetCursor(ctx context.Context) (lastScannedBlock uint64, found bool, err error)
	SetCursor(ctx context.Context, lastScannedBlock uint64) error
}
