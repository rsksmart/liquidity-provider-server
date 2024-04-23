package dataproviders

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	singletonLock     = &sync.Mutex{}
	eventBusSingleton *LocalEventBus
)

const (
	subscriptionBufferSize = 20
	closedBusMessage       = "Trying to interact with closed bus"
)

type LocalEventBus struct {
	Topics         map[entities.EventId][]chan<- entities.Event
	subscribeMutex sync.Mutex
	publishMutex   sync.Mutex
}

func NewLocalEventBus() entities.EventBus {
	if eventBusSingleton == nil {
		singletonLock.Lock()
		defer singletonLock.Unlock()
		// we need to check if it still not created after the getting into the critical section
		if eventBusSingleton == nil {
			topics := make(map[entities.EventId][]chan<- entities.Event)
			eventBusSingleton = &LocalEventBus{Topics: topics}
		}
	}
	return eventBusSingleton
}

func (bus *LocalEventBus) Subscribe(id entities.EventId) <-chan entities.Event {
	if eventBusSingleton == nil {
		log.Error(closedBusMessage)
		return nil
	}
	var topics []chan<- entities.Event
	var ok bool
	bus.subscribeMutex.Lock()
	defer bus.subscribeMutex.Unlock()
	if topics, ok = bus.Topics[id]; !ok {
		topics = make([]chan<- entities.Event, 0)
		bus.Topics[id] = topics
	}
	subscription := make(chan entities.Event, subscriptionBufferSize)
	bus.Topics[id] = append(topics, subscription)
	return subscription
}

func (bus *LocalEventBus) Shutdown(closeChannel chan<- bool) {
	if eventBusSingleton == nil {
		log.Error(closedBusMessage)
		return
	}
	for key, topic := range bus.Topics {
		for _, subscription := range topic {
			close(subscription)
		}
		delete(bus.Topics, key)
	}
	singletonLock.Lock()
	defer singletonLock.Unlock()
	eventBusSingleton = nil
	closeChannel <- true
	log.Debug("Event bus shut down")
}

func (bus *LocalEventBus) Publish(event entities.Event) {
	if eventBusSingleton == nil {
		log.Error(closedBusMessage)
		return
	}
	bus.publishMutex.Lock()
	defer bus.publishMutex.Unlock()
	topic, ok := bus.Topics[event.Id()]
	if !ok {
		return
	}
	for _, subscription := range topic {
		subscription <- event
	}
}
