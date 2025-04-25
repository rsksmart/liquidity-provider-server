package bootstrap

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	log "github.com/sirupsen/logrus"
)

func Mongo(ctx context.Context, env environment.MongoEnv, timeouts environment.ApplicationTimeouts) (*mongo.Connection, error) {
	client, err := mongo.Connect(ctx, timeouts.DatabaseConnection.Seconds(), env.Username, env.Password, env.Host, env.Port)
	if err == nil {
		return mongo.NewConnection(mongo.NewClientWrapper(client), timeouts.DatabaseInteraction.Seconds()), nil
	} else {
		log.Error("Error connecting to MongoDB: ", err)
		return nil, err
	}
}
