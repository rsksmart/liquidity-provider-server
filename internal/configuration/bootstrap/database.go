package bootstrap

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	log "github.com/sirupsen/logrus"
)

func Mongo(ctx context.Context, env environment.MongoEnv) (*mongo.Connection, error) {
	client, err := mongo.Connect(ctx, env.Username, env.Password, env.Host, env.Port)
	if err == nil {
		return mongo.NewConnection(mongo.NewClientWrapper(client)), nil
	} else {
		log.Error("Error connecting to MongoDB: ", err)
		return nil, err
	}
}
