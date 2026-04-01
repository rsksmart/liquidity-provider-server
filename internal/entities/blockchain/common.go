package blockchain

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"time"
)

const (
	NodeEclipseEventId        entities.EventId = "NodeEclipse"
	NodePeerCheckEventId      entities.EventId = "NodePeerCheck"
	NodePeerCheckErrorEventId entities.EventId = "NodePeerCheckError"
	NodePeerAlertSentEventId  entities.EventId = "NodePeerAlertSent"
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

type NodePeerCheckEvent struct {
	entities.BaseEvent
	NodeType       entities.NodeType
	CurrentPeers   int64
	MinPeers       uint64
	BelowThreshold bool
}

type NodePeerCheckErrorEvent struct {
	entities.BaseEvent
	NodeType entities.NodeType
}

type NodePeerAlertSentEvent struct {
	entities.BaseEvent
	NodeType entities.NodeType
}
