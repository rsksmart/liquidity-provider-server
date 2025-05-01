package mongo

import (
	"context"
	"errors"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const TrustedAccountCollection = "trustedAccounts"

type trustedAccountMongoRepository struct {
	conn *Connection
}

func NewTrustedAccountRepository(conn *Connection) liquidity_provider.TrustedAccountRepository {
	return &trustedAccountMongoRepository{conn: conn}
}

func (repo *trustedAccountMongoRepository) GetTrustedAccount(ctx context.Context, address string) (*liquidity_provider.TrustedAccountDetails, error) {
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

func (repo *trustedAccountMongoRepository) GetAllTrustedAccounts(ctx context.Context) ([]liquidity_provider.TrustedAccountDetails, error) {
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

func (repo *trustedAccountMongoRepository) AddTrustedAccount(ctx context.Context, account liquidity_provider.TrustedAccountDetails) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	collection := repo.conn.Collection(TrustedAccountCollection)
	existingAccount, err := repo.GetTrustedAccount(ctx, account.Address)
	if err == nil && existingAccount != nil {
		return liquidity_provider.ErrDuplicateTrustedAccount
	} else if err != nil && !errors.Is(err, liquidity_provider.ErrTrustedAccountNotFound) {
		return err
	}
	_, err = collection.InsertOne(dbCtx, account)
	if err != nil {
		return err
	}
	logDbInteraction(Insert, account)
	return nil
}

func (repo *trustedAccountMongoRepository) UpdateTrustedAccount(ctx context.Context, account liquidity_provider.TrustedAccountDetails) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	collection := repo.conn.Collection(TrustedAccountCollection)
	_, err := repo.GetTrustedAccount(ctx, account.Address)
	if err != nil && !errors.Is(err, liquidity_provider.ErrTrustedAccountNotFound) {
		return err
	}
	filter := bson.M{"address": account.Address}
	opts := options.Update().SetUpsert(true)
	update := bson.M{"$set": account}
	_, err = collection.UpdateOne(dbCtx, filter, update, opts)
	if err != nil {
		return err
	}
	logDbInteraction(Update, account)
	return nil
}

func (repo *trustedAccountMongoRepository) DeleteTrustedAccount(ctx context.Context, address string) error {
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
