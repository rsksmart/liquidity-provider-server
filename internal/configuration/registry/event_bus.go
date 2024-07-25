package registry

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

func NewEventBus() entities.EventBus {
	return dataproviders.NewLocalEventBus()
}
