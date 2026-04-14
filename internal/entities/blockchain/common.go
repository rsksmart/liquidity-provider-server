package blockchain

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"time"
)

const (
	NodeEclipseEventId         entities.EventId = "NodeEclipse"
	NodeReorgCheckEventId      entities.EventId = "NodeReorgCheck"
	NodeReorgCheckErrorEventId entities.EventId = "NodeReorgCheckError"
	NodeReorgAlertSentEventId  entities.EventId = "NodeReorgAlertSent"
	NodePeerCheckEventId       entities.EventId = "NodePeerCheck"
	NodePeerCheckErrorEventId  entities.EventId = "NodePeerCheckError"
	NodePeerAlertSentEventId   entities.EventId = "NodePeerAlertSent"
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

type NodeReorgCheckEvent struct {
	entities.BaseEvent
	NodeType        entities.NodeType
	CurrentDepth    uint64
	MaxAllowedDepth uint64
	AboveThreshold  bool
}

type NodeReorgCheckErrorEvent struct {
	entities.BaseEvent
	NodeType entities.NodeType
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

type NodeReorgAlertSentEvent struct {
	entities.BaseEvent
	NodeType      entities.NodeType
	DetectedDepth uint64
}

type NodePeerAlertSentEvent struct {
	entities.BaseEvent
	NodeType entities.NodeType
}
