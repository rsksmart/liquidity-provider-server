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
	DepositTxID      string                         `json:"depositTxId" bson:"deposit_txid"`
	Confirmations    uint64                         `json:"confirmations" bson:"confirmations"`
	LastSeenAt       *time.Time                     `json:"lastSeenAt" bson:"last_seen_at"` // post-import wallet observation
	LastError        string                         `json:"lastError" bson:"last_error"`
	CreatedAt        time.Time                      `json:"createdAt" bson:"created_at"`
	UpdatedAt        time.Time                      `json:"updatedAt" bson:"updated_at"`
}

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

type PegInAddressRegistryWatchRepository interface {
	Upsert(ctx context.Context, watch PegInAddressRegistryWatch) error
	Get(ctx context.Context, txHash string, logIndex uint) (*PegInAddressRegistryWatch, error)
	List(ctx context.Context) ([]PegInAddressRegistryWatch, error)
	Update(ctx context.Context, watch PegInAddressRegistryWatch) error
	GetCursor(ctx context.Context) (lastScannedBlock uint64, found bool, err error)
	SetCursor(ctx context.Context, lastScannedBlock uint64) error
}
