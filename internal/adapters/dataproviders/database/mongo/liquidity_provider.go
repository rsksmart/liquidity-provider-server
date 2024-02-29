package mongo

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const liquidityProviderCollection = "liquidityProvider"

type configurationName string

const (
	peginConfigId   configurationName = "pegin"
	pegoutConfigId  configurationName = "pegout"
	generalConfigId configurationName = "general"
)

type lpMongoRepository struct {
	conn *Connection
}

type storedConfiguration[C liquidity_provider.ConfigurationType] struct {
	entities.Signed[C] `bson:",inline"`
	Name               configurationName `json:"name" bson:"name"`
}

func NewLiquidityProviderRepository(conn *Connection) liquidity_provider.LiquidityProviderRepository {
	return &lpMongoRepository{conn: conn}
}

func (repo *lpMongoRepository) GetPeginConfiguration(ctx context.Context) (*entities.Signed[liquidity_provider.PeginConfiguration], error) {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return getConfiguration[liquidity_provider.PeginConfiguration](dbCtx, repo, peginConfigId)
}

func (repo *lpMongoRepository) GetPegoutConfiguration(ctx context.Context) (*entities.Signed[liquidity_provider.PegoutConfiguration], error) {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return getConfiguration[liquidity_provider.PegoutConfiguration](dbCtx, repo, pegoutConfigId)
}

func (repo *lpMongoRepository) GetGeneralConfiguration(ctx context.Context) (*entities.Signed[liquidity_provider.GeneralConfiguration], error) {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return getConfiguration[liquidity_provider.GeneralConfiguration](dbCtx, repo, generalConfigId)
}

func (repo *lpMongoRepository) UpsertPeginConfiguration(ctx context.Context, signedConfig entities.Signed[liquidity_provider.PeginConfiguration]) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	configToStore := storedConfiguration[liquidity_provider.PeginConfiguration]{
		Signed: signedConfig,
		Name:   peginConfigId,
	}
	return upsertConfiguration(dbCtx, repo, configToStore)
}

func (repo *lpMongoRepository) UpsertPegoutConfiguration(ctx context.Context, signedConfig entities.Signed[liquidity_provider.PegoutConfiguration]) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	configToStore := storedConfiguration[liquidity_provider.PegoutConfiguration]{
		Signed: signedConfig,
		Name:   pegoutConfigId,
	}
	return upsertConfiguration(dbCtx, repo, configToStore)
}

func (repo *lpMongoRepository) UpsertGeneralConfiguration(ctx context.Context, signedConfig entities.Signed[liquidity_provider.GeneralConfiguration]) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	configToStore := storedConfiguration[liquidity_provider.GeneralConfiguration]{
		Signed: signedConfig,
		Name:   generalConfigId,
	}
	return upsertConfiguration(dbCtx, repo, configToStore)
}

func upsertConfiguration[C liquidity_provider.ConfigurationType](
	ctx context.Context,
	repo *lpMongoRepository,
	config storedConfiguration[C],
) error {
	collection := repo.conn.Collection(liquidityProviderCollection)
	options := options.Replace().SetUpsert(true)
	filter := bson.D{primitive.E{Key: "name", Value: config.Name}}
	_, err := collection.ReplaceOne(ctx, filter, config, options)
	if err != nil {
		return err
	} else {
		logDbInteraction(insert, config)
		return nil
	}
}

func getConfiguration[C liquidity_provider.ConfigurationType](
	ctx context.Context,
	repo *lpMongoRepository,
	name configurationName,
) (*entities.Signed[C], error) {
	config := &storedConfiguration[C]{}
	collection := repo.conn.Collection(liquidityProviderCollection)
	filter := bson.D{primitive.E{Key: "name", Value: name}}

	err := collection.FindOne(ctx, filter).Decode(config)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	logDbInteraction(read, config.Signed)
	return &config.Signed, nil
}
