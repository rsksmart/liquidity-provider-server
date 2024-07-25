package entities_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"testing"
)

func TestNewBaseEvent(t *testing.T) {
	var id entities.EventId = "any id"
	var event entities.Event = entities.NewBaseEvent(id)
	if event.Id() != id || event.CreationTimestamp().IsZero() {
		t.Error("Base event not initialized properly")
	}
}
