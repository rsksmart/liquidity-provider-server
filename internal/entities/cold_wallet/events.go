package cold_wallet

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

const (
	RbtcTransferredDueToThresholdEventId   entities.EventId = "RbtcTransferredDueToThreshold"
	BtcTransferredDueToThresholdEventId    entities.EventId = "BtcTransferredDueToThreshold"
	RbtcTransferredDueToTimeForcingEventId entities.EventId = "RbtcTransferredDueToTimeForcing"
	BtcTransferredDueToTimeForcingEventId  entities.EventId = "BtcTransferredDueToTimeForcing"
)

type RbtcTransferredDueToThresholdEvent struct {
	entities.Event
	Amount *entities.Wei
	TxHash string
	Fee    *entities.Wei
}

type BtcTransferredDueToThresholdEvent struct {
	entities.Event
	Amount *entities.Wei
	TxHash string
	Fee    *entities.Wei
}

type RbtcTransferredDueToTimeForcingEvent struct {
	entities.Event
	Amount *entities.Wei
	TxHash string
	Fee    *entities.Wei
}

type BtcTransferredDueToTimeForcingEvent struct {
	entities.Event
	Amount *entities.Wei
	TxHash string
	Fee    *entities.Wei
}
