package rootstock

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

const (
	BatchPegOutUpdatedEventId entities.EventId = "BatchPegOutUpdated"
)

type BatchPegOut struct {
	TransactionHash    string   `json:"transactionHash" bson:"transaction_hash" validate:"required"`
	BlockHash          string   `json:"blockHash" bson:"block_hash" validate:"required"`
	BlockNumber        uint64   `json:"blockNumber" bson:"block_number" validate:"required"`
	BtcTxHash          string   `json:"btcTxHash" bson:"btc_tx_hash" validate:"required"`
	ReleaseRskTxHashes []string `json:"releaseRskTxHashes" bson:"release_rsk_tx_hashes" validate:"required"`
}

type BatchPegOutRepository interface {
	UpsertBatch(context context.Context, batch BatchPegOut) error
}

type BatchPegOutUpdatedEvent struct {
	entities.Event
	QuoteHashes []string
	BatchPegOut BatchPegOut
}
