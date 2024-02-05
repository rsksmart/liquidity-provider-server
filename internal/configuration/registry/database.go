package registry

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
)

type Database struct {
	PeginRepository  quote.PeginQuoteRepository
	PegoutRepository quote.PegoutQuoteRepository
	Connection       *mongo.Connection
}

func NewDatabaseRegistry(connection *mongo.Connection) *Database {
	return &Database{
		PeginRepository:  mongo.NewPeginMongoRepository(connection),
		PegoutRepository: mongo.NewPegoutMongoRepository(connection),
		Connection:       connection,
	}
}
