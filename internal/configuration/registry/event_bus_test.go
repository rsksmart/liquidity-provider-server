package registry_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEventBus(t *testing.T) {
	bus := registry.NewEventBus()
	assert.IsType(t, &dataproviders.LocalEventBus{}, bus)
}
