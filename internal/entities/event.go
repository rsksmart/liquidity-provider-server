package entities

import (
	"sync"
	"time"
)

type EventId string

type Event interface {
	Id() EventId
	CreationTimestamp() time.Time
}

type BaseEvent struct {
	EventId   EventId
	Timestamp time.Time
}

func NewBaseEvent(id EventId) BaseEvent {
	return BaseEvent{EventId: id, Timestamp: time.Now()}
}

func (e BaseEvent) Id() EventId {
	return e.EventId
}

func (e BaseEvent) CreationTimestamp() time.Time {
	return e.Timestamp
}

type EventBus interface {
	Closeable
	Publish(events Event)
	Subscribe(id EventId) <-chan Event
	// Shutdown since subscriptions return event channel, shutdown should close all the subscription channels
	Shutdown(chan<- bool)
}

type ApplicationMutexes interface {
	RskWalletMutex() *sync.Mutex
	BtcWalletMutex() *sync.Mutex
	PeginLiquidityMutex() *sync.Mutex
	PegoutLiquidityMutex() *sync.Mutex
}
