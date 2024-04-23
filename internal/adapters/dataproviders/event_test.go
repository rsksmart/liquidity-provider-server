package dataproviders_test

import (
	"bytes"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

const (
	testEventId         entities.EventId = "test_event"
	channelNotClosedMsg string           = "Channel was not closed"
)

type testEvent struct {
	entities.Event
}

func TestNewLocalEventBus(t *testing.T) {
	bus1 := dataproviders.NewLocalEventBus()
	bus2 := dataproviders.NewLocalEventBus()
	t.Run("Should always return the same bus", func(t *testing.T) {
		assert.NotNil(t, bus1)
		assert.Same(t, bus1, bus2)
	})
	t.Run("Should initialize topics", func(t *testing.T) {
		assert.NotNil(t, bus1.(*dataproviders.LocalEventBus).Topics)
	})
}

func TestLocalEventBus_Shutdown(t *testing.T) {
	t.Run("Should send signal to close channel after closing", func(t *testing.T) {
		bus1 := dataproviders.NewLocalEventBus()
		closeChannel := make(chan bool, 1)
		bus1.Shutdown(closeChannel)
		select {
		case result := <-closeChannel:
			assert.True(t, result)
		default:
			assert.Fail(t, channelNotClosedMsg)
		}
	})

	t.Run("Should clear topic map", func(t *testing.T) {
		bus := dataproviders.NewLocalEventBus()
		_ = bus.Subscribe(testEventId)
		assert.Len(t, bus.(*dataproviders.LocalEventBus).Topics, 1)
		closeChannel := make(chan bool, 1)
		bus.Shutdown(closeChannel)
		select {
		case <-closeChannel:
			assert.Empty(t, bus.(*dataproviders.LocalEventBus).Topics)
		default:
			assert.Fail(t, channelNotClosedMsg)
		}
	})

	t.Run("Should not allow to interact after closed", func(t *testing.T) {
		const expectedBuff = "Trying to interact with closed bus"
		message := make([]byte, 100)
		buff := new(bytes.Buffer)
		log.SetOutput(buff)
		bus := dataproviders.NewLocalEventBus()
		closeChannel := make(chan bool, 1)
		bus.Shutdown(closeChannel)
		select {
		case <-closeChannel:
			assert.Empty(t, bus.(*dataproviders.LocalEventBus).Topics)
		default:
			assert.Fail(t, channelNotClosedMsg)
		}
		require.NotPanics(t, func() {
			result := bus.Subscribe(testEventId)
			assert.Nil(t, result)
			_, err := buff.Read(message)
			require.NoError(t, err)
			assert.Contains(t, string(message), expectedBuff)
		})
		require.NotPanics(t, func() {
			bus.Publish(testEvent{})
			_, err := buff.Read(message)
			require.NoError(t, err)
			assert.Contains(t, string(message), expectedBuff)
		})
		require.NotPanics(t, func() {
			bus.Shutdown(make(chan bool))
			_, err := buff.Read(message)
			require.NoError(t, err)
			assert.Contains(t, string(message), expectedBuff)
		})
	})
}

func TestLocalEventBus_Subscribe(t *testing.T) {
	bus := dataproviders.NewLocalEventBus()
	firstSub := bus.Subscribe(testEventId)
	secondSub := bus.Subscribe(testEventId)
	t.Run("Should create one channel per subscription", func(t *testing.T) {
		assert.NotEqual(t, firstSub, secondSub)
		assert.Len(t, bus.(*dataproviders.LocalEventBus).Topics[testEventId], 2)
	})
}

func TestLocalEventBus_Publish(t *testing.T) {
	t.Run("Should send message to every subscription", func(t *testing.T) {
		assert.Eventually(t, func() bool {
			var wg sync.WaitGroup
			wg.Add(2)
			bus := dataproviders.NewLocalEventBus()
			sub1 := bus.Subscribe(testEventId)
			sub2 := bus.Subscribe(testEventId)
			event := testEvent{entities.NewBaseEvent(testEventId)}
			bus.Publish(event)
			go func() {
				assert.Equal(t, event, <-sub1)
				wg.Done()
			}()
			go func() {
				assert.Equal(t, event, <-sub2)
				wg.Done()
			}()
			wg.Wait()
			return true
		}, time.Second*1, time.Millisecond*10)
	})
	t.Run("Should not send message to other subscriptions", func(t *testing.T) {
		bus := dataproviders.NewLocalEventBus()
		sub1 := bus.Subscribe(testEventId)
		sub2 := bus.Subscribe(testEventId + "-copy")
		event := testEvent{entities.NewBaseEvent(testEventId)}
		bus.Publish(event)
		assert.Equal(t, event, <-sub1)
		select {
		case <-sub2:
			assert.Fail(t, "Should not receive message")
		default:
		}
	})
	t.Run("Should not send message on non existing event", func(t *testing.T) {
		const errorMessage = "Should not receive message"
		bus := dataproviders.NewLocalEventBus()
		sub1 := bus.Subscribe(testEventId)
		sub2 := bus.Subscribe(testEventId)
		event := testEvent{entities.NewBaseEvent(testEventId + "-other")}
		bus.Publish(event)
		select {
		case <-sub2:
			assert.Fail(t, errorMessage)
		case <-sub1:
			assert.Fail(t, errorMessage)
		default:
		}
	})
}
