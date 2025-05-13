package registry

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
)

type Database struct {
	PeginRepository             quote.PeginQuoteRepository
	PegoutRepository            quote.PegoutQuoteRepository
	LiquidityProviderRepository liquidity_provider.LiquidityProviderRepository
	TrustedAccountRepository    liquidity_provider.TrustedAccountRepository
	Connection                  *mongo.Connection
	TrustedAccountRepository    liquidity_provider.TrustedAccountRepository
}

func NewDatabaseRegistry(connection *mongo.Connection) *Database {
	return &Database{
		PeginRepository:             mongo.NewPeginMongoRepository(connection),
		PegoutRepository:            mongo.NewPegoutMongoRepository(connection),
		LiquidityProviderRepository: mongo.NewLiquidityProviderRepository(connection),
		TrustedAccountRepository:    mongo.NewTrustedAccountRepository(connection),
		Connection:                  connection,
	}
}
