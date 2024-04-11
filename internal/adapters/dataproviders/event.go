package dataproviders

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	singletonLock     = &sync.Mutex{}
	eventBusSingleton *localEventBus
)

const subscriptionBufferSize = 20

type localEventBus struct {
	topics         map[entities.EventId][]chan<- entities.Event
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
			eventBusSingleton = &localEventBus{topics: topics}
		}
	}
	return eventBusSingleton
}

func (bus *localEventBus) Subscribe(id entities.EventId) <-chan entities.Event {
	var topics []chan<- entities.Event
	var ok bool
	bus.subscribeMutex.Lock()
	defer bus.subscribeMutex.Unlock()
	if topics, ok = bus.topics[id]; !ok {
		topics = make([]chan<- entities.Event, 0)
		bus.topics[id] = topics
	}
	subscription := make(chan entities.Event, subscriptionBufferSize)
	bus.topics[id] = append(topics, subscription)
	return subscription
}

func (bus *localEventBus) Shutdown(closeChannel chan<- bool) {
	for _, topic := range bus.topics {
		for _, subscription := range topic {
			close(subscription)
		}
	}
	closeChannel <- true
	log.Debug("Event bus shut down")
}

func (bus *localEventBus) Publish(event entities.Event) {
	bus.publishMutex.Lock()
	defer bus.publishMutex.Unlock()
	topic, ok := bus.topics[event.Id()]
	if !ok {
		return
	}
	for _, subscription := range topic {
		subscription <- event
	}
}
