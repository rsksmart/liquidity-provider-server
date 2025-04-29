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
const TrustedAccountCollection = "trustedAccounts"

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

func NewTrustedAccountRepository(conn *Connection) liquidity_provider.TrustedAccountRepository {
	return &lpMongoRepository{conn: conn}
}

func (repo *lpMongoRepository) GetPeginConfiguration(ctx context.Context) (*entities.Signed[liquidity_provider.PeginConfiguration], error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	return getConfigurationVerbose[liquidity_provider.PeginConfiguration](dbCtx, repo, peginConfigId)
}

func (repo *lpMongoRepository) GetPegoutConfiguration(ctx context.Context) (*entities.Signed[liquidity_provider.PegoutConfiguration], error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	return getConfigurationVerbose[liquidity_provider.PegoutConfiguration](dbCtx, repo, pegoutConfigId)
}

func (repo *lpMongoRepository) GetGeneralConfiguration(ctx context.Context) (*entities.Signed[liquidity_provider.GeneralConfiguration], error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	return getConfigurationVerbose[liquidity_provider.GeneralConfiguration](dbCtx, repo, generalConfigId)
}

func (repo *lpMongoRepository) UpsertPeginConfiguration(ctx context.Context, signedConfig entities.Signed[liquidity_provider.PeginConfiguration]) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	configToStore := StoredConfiguration[liquidity_provider.PeginConfiguration]{
		Signed: signedConfig,
		Name:   peginConfigId,
	}
	return upsertConfigurationVerbose(dbCtx, repo, configToStore)
}

func (repo *lpMongoRepository) UpsertPegoutConfiguration(ctx context.Context, signedConfig entities.Signed[liquidity_provider.PegoutConfiguration]) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	configToStore := StoredConfiguration[liquidity_provider.PegoutConfiguration]{
		Signed: signedConfig,
		Name:   pegoutConfigId,
	}
	return upsertConfigurationVerbose(dbCtx, repo, configToStore)
}

func (repo *lpMongoRepository) UpsertGeneralConfiguration(ctx context.Context, signedConfig entities.Signed[liquidity_provider.GeneralConfiguration]) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	configToStore := StoredConfiguration[liquidity_provider.GeneralConfiguration]{
		Signed: signedConfig,
		Name:   generalConfigId,
	}
	return upsertConfigurationVerbose(dbCtx, repo, configToStore)
}

func (repo *lpMongoRepository) GetCredentials(ctx context.Context) (*entities.Signed[liquidity_provider.HashedCredentials], error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	return getConfiguration[liquidity_provider.HashedCredentials](dbCtx, repo, credentialsId, false)
}

func (repo *lpMongoRepository) UpsertCredentials(ctx context.Context, credentials entities.Signed[liquidity_provider.HashedCredentials]) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	configToStore := StoredConfiguration[liquidity_provider.HashedCredentials]{
		Signed: credentials,
		Name:   credentialsId,
	}
	return upsertConfiguration(dbCtx, repo, configToStore, false)
}

func (repo *lpMongoRepository) GetTrustedAccount(ctx context.Context, address string) (*liquidity_provider.TrustedAccountDetails, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	account := &liquidity_provider.TrustedAccountDetails{}
	collection := repo.conn.Collection(TrustedAccountCollection)
	filter := bson.M{"address": address}
	err := collection.FindOne(dbCtx, filter).Decode(account)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, liquidity_provider.ErrTrustedAccountNotFound
	} else if err != nil {
		return nil, err
	}
	logDbInteraction(Read, account)
	return account, nil
}

func (repo *lpMongoRepository) GetAllTrustedAccounts(ctx context.Context) ([]liquidity_provider.TrustedAccountDetails, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	collection := repo.conn.Collection(TrustedAccountCollection)
	cursor, err := collection.Find(dbCtx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(dbCtx)
	var accounts []liquidity_provider.TrustedAccountDetails
	if err = cursor.All(dbCtx, &accounts); err != nil {
		return nil, err
	}
	logDbInteraction(Read, accounts)
	return accounts, nil
}

func (repo *lpMongoRepository) AddTrustedAccount(ctx context.Context, account liquidity_provider.TrustedAccountDetails) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	collection := repo.conn.Collection(TrustedAccountCollection)
	existingAccount, err := repo.GetTrustedAccount(ctx, account.Address)
	if err == nil && existingAccount != nil {
		return liquidity_provider.ErrDuplicateAddress
	} else if err != nil && err != liquidity_provider.ErrTrustedAccountNotFound {
		return err
	}
	_, err = collection.InsertOne(dbCtx, account)
	if err != nil {
		return err
	}
	logDbInteraction(Insert, account)
	return nil
}

func (repo *lpMongoRepository) UpdateTrustedAccount(ctx context.Context, account liquidity_provider.TrustedAccountDetails) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	collection := repo.conn.Collection(TrustedAccountCollection)
	_, err := repo.GetTrustedAccount(ctx, account.Address)
	if err != nil && err != liquidity_provider.ErrTrustedAccountNotFound {
		return err
	}
	filter := bson.M{"address": account.Address}
	opts := options.Update().SetUpsert(true)
	update := bson.M{"$set": account}
	_, err = collection.UpdateOne(dbCtx, filter, update, opts)
	if err != nil {
		return err
	}
	logDbInteraction(Insert, account)
	return nil
}

func (repo *lpMongoRepository) DeleteTrustedAccount(ctx context.Context, address string) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	collection := repo.conn.Collection(TrustedAccountCollection)
	filter := bson.M{"address": address}
	result, err := collection.DeleteOne(dbCtx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return liquidity_provider.ErrTrustedAccountNotFound
	}
	logDbInteraction(Delete, filter)
	return nil
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
