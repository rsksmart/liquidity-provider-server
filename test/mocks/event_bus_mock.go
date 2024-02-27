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
