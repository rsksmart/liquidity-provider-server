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

const LiquidityProviderCollection = "liquidityProvider"

type ConfigurationName string

const (
	peginConfigId   ConfigurationName = "pegin"
	pegoutConfigId  ConfigurationName = "pegout"
	generalConfigId ConfigurationName = "general"
	credentialsId   ConfigurationName = "credentials"
)

type lpMongoRepository struct {
	conn *Connection
}

type StoredConfiguration[C liquidity_provider.ConfigurationType] struct {
	entities.Signed[C] `bson:",inline"`
	Name               ConfigurationName `json:"name" bson:"name"`
}

func NewLiquidityProviderRepository(conn *Connection) liquidity_provider.LiquidityProviderRepository {
	return &lpMongoRepository{conn: conn}
}

func (repo *lpMongoRepository) GetPeginConfiguration(ctx context.Context) (*entities.Signed[liquidity_provider.PeginConfiguration], error) {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return getConfigurationVerbose[liquidity_provider.PeginConfiguration](dbCtx, repo, peginConfigId)
}

func (repo *lpMongoRepository) GetPegoutConfiguration(ctx context.Context) (*entities.Signed[liquidity_provider.PegoutConfiguration], error) {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return getConfigurationVerbose[liquidity_provider.PegoutConfiguration](dbCtx, repo, pegoutConfigId)
}

func (repo *lpMongoRepository) GetGeneralConfiguration(ctx context.Context) (*entities.Signed[liquidity_provider.GeneralConfiguration], error) {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return getConfigurationVerbose[liquidity_provider.GeneralConfiguration](dbCtx, repo, generalConfigId)
}

func (repo *lpMongoRepository) UpsertPeginConfiguration(ctx context.Context, signedConfig entities.Signed[liquidity_provider.PeginConfiguration]) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	configToStore := StoredConfiguration[liquidity_provider.PeginConfiguration]{
		Signed: signedConfig,
		Name:   peginConfigId,
	}
	return upsertConfigurationVerbose(dbCtx, repo, configToStore)
}

func (repo *lpMongoRepository) UpsertPegoutConfiguration(ctx context.Context, signedConfig entities.Signed[liquidity_provider.PegoutConfiguration]) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	configToStore := StoredConfiguration[liquidity_provider.PegoutConfiguration]{
		Signed: signedConfig,
		Name:   pegoutConfigId,
	}
	return upsertConfigurationVerbose(dbCtx, repo, configToStore)
}

func (repo *lpMongoRepository) UpsertGeneralConfiguration(ctx context.Context, signedConfig entities.Signed[liquidity_provider.GeneralConfiguration]) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	configToStore := StoredConfiguration[liquidity_provider.GeneralConfiguration]{
		Signed: signedConfig,
		Name:   generalConfigId,
	}
	return upsertConfigurationVerbose(dbCtx, repo, configToStore)
}

func (repo *lpMongoRepository) GetCredentials(ctx context.Context) (*entities.Signed[liquidity_provider.HashedCredentials], error) {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return getConfiguration[liquidity_provider.HashedCredentials](dbCtx, repo, credentialsId, false)
}

func (repo *lpMongoRepository) UpsertCredentials(ctx context.Context, credentials entities.Signed[liquidity_provider.HashedCredentials]) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	configToStore := StoredConfiguration[liquidity_provider.HashedCredentials]{
		Signed: credentials,
		Name:   credentialsId,
	}
	return upsertConfiguration(dbCtx, repo, configToStore, false)
}

func upsertConfigurationVerbose[C liquidity_provider.ConfigurationType](
	ctx context.Context,
	repo *lpMongoRepository,
	config StoredConfiguration[C],
) error {
	return upsertConfiguration(ctx, repo, config, true)
}

func getConfigurationVerbose[C liquidity_provider.ConfigurationType](
	ctx context.Context,
	repo *lpMongoRepository,
	name ConfigurationName,
) (*entities.Signed[C], error) {
	return getConfiguration[C](ctx, repo, name, true)
}

func upsertConfiguration[C liquidity_provider.ConfigurationType](
	ctx context.Context,
	repo *lpMongoRepository,
	config StoredConfiguration[C],
	logInteraction bool,
) error {
	collection := repo.conn.Collection(LiquidityProviderCollection)
	opts := options.Replace().SetUpsert(true)
	filter := bson.D{primitive.E{Key: "name", Value: config.Name}}
	_, err := collection.ReplaceOne(ctx, filter, config, opts)
	if err != nil {
		return err
	}
	if logInteraction {
		logDbInteraction(Insert, config)
	}
	return nil
}

func getConfiguration[C liquidity_provider.ConfigurationType](
	ctx context.Context,
	repo *lpMongoRepository,
	name ConfigurationName,
	logInteraction bool,
) (*entities.Signed[C], error) {
	config := &StoredConfiguration[C]{}
	collection := repo.conn.Collection(LiquidityProviderCollection)
	filter := bson.D{primitive.E{Key: "name", Value: name}}

	err := collection.FindOne(ctx, filter).Decode(config)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if logInteraction {
		logDbInteraction(Read, config.Signed)
	}
	return &config.Signed, nil
}
