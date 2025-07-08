package blockchain

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"time"
)

const (
	NodeEclipseEventId entities.EventId = "NodeEclipse"
)

type Rpc struct {
	Btc BitcoinNetwork
	Rsk RootstockRpcServer
}

type NodeEclipseEvent struct {
	entities.BaseEvent
	NodeType            entities.NodeType
	EclipsedBlockNumber uint64
	EclipsedBlockHash   string
	DetectionTime       time.Time
}
