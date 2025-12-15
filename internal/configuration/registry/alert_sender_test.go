package registry_test

import (
	"context"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/alerting"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/stretchr/testify/assert"
)

func TestNewAlertSender(t *testing.T) {
	env := environment.Environment{LpsStage: "testnet", Provider: environment.ProviderEnv{AlertSenderEmail: "fake@email.com"}}
	sender := registry.NewAlertSender(context.Background(), env)
	implementationPointer, ok := sender.(*alerting.LogAlertSender)
	assert.NotNil(t, sender)
	assert.True(t, ok)
	assert.NotNil(t, implementationPointer)
}
