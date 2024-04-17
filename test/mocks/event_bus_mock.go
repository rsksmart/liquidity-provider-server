package mocks

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/stretchr/testify/mock"
	"sync"
)

type MutexMock struct {
	mock.Mock
	sync.Mutex
}

func (m *MutexMock) Lock() {
	m.Called()
}

func (m *MutexMock) Unlock() {
	m.Called()
}

type EventBusMock struct {
	entities.EventBus
	mock.Mock
}

func (m *EventBusMock) Publish(event entities.Event) {
	m.Called(event)
}

func (m *EventBusMock) Subscribe(eventId entities.EventId) <-chan entities.Event {
	args := m.Called(eventId)
	return args.Get(0).(<-chan entities.Event)
}

func (m *EventBusMock) Shutdown(shutdownChannel chan<- bool) {
	m.Called(shutdownChannel)
}
