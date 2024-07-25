package registry_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/alerting"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAlertSender(t *testing.T) {
	env := environment.Environment{LpsStage: "testnet", Provider: environment.ProviderEnv{AlertSenderEmail: "fake@email.com"}}
	sender := registry.NewAlertSender(context.Background(), env)
	implementationPointer, ok := sender.(*alerting.SesAlertSender)
	assert.NotNil(t, sender)
	assert.True(t, ok)
	assert.Equal(t, 2, test.CountNonZeroValues(*implementationPointer))
}
