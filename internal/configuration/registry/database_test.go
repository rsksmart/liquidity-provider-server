package registry_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDatabaseRegistry(t *testing.T) {
	t.Run("should return a new database registry", func(t *testing.T) {
		client := &mocks.DbClientBindingMock{}
		client.On("Database", mongo.DbName).Return(&mocks.DbBindingMock{})
		conn := mongo.NewConnection(client)
		dbRegistry := registry.NewDatabaseRegistry(conn)
		assert.NotNil(t, dbRegistry.PeginRepository)
		assert.NotNil(t, dbRegistry.PegoutRepository)
		assert.NotNil(t, dbRegistry.LiquidityProviderRepository)
		assert.Equal(t, conn, dbRegistry.Connection)
	})
}
