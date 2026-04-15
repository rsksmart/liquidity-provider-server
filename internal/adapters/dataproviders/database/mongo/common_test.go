package mongo_test

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func TestNewConnection(t *testing.T) {
	client := &mocks.DbClientBindingMock{}
	client.On("Database", mongo.DbName).Return(&mocks.DbBindingMock{})
	conn := mongo.NewConnection(client, time.Duration(1))
	test.AssertNonZeroValues(t, conn)
}

func TestConnection_GetDb(t *testing.T) {
	client := &mocks.DbClientBindingMock{}
	client.On("Database", mongo.DbName).Return(&mocks.DbBindingMock{})
	conn := mongo.NewConnection(client, time.Duration(1))
	assert.NotNil(t, conn.GetDb())
}

func TestConnection_CheckConnection(t *testing.T) {
	t.Run("Connection ok", func(t *testing.T) {
		client := &mocks.DbClientBindingMock{}
		client.On("Database", mongo.DbName).Return(&mocks.DbBindingMock{})
		client.On("Ping", test.AnyCtx, (*readpref.ReadPref)(nil)).Return(nil)
		conn := mongo.NewConnection(client, time.Duration(1))
		result := conn.CheckConnection(context.Background())
		assert.True(t, result)
		client.AssertExpectations(t)
	})
	t.Run("Connection error", func(t *testing.T) {
		client := &mocks.DbClientBindingMock{}
		client.On("Database", mongo.DbName).Return(&mocks.DbBindingMock{})
		client.On("Ping", test.AnyCtx, (*readpref.ReadPref)(nil)).Return(assert.AnError)
		conn := mongo.NewConnection(client, time.Duration(1))
		result := conn.CheckConnection(context.Background())
		assert.False(t, result)
		client.AssertExpectations(t)
	})
}

func TestConnection_Shutdown(t *testing.T) {
	t.Run("Disconnect ok", func(t *testing.T) {
		client := &mocks.DbClientBindingMock{}
		client.On("Database", mongo.DbName).Return(&mocks.DbBindingMock{})
		client.On("Disconnect", mock.Anything).Return(nil)
		conn := mongo.NewConnection(client, time.Duration(1))
		closeChannel := make(chan bool)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-closeChannel
		}()
		conn.Shutdown(closeChannel)
		wg.Wait()
		client.AssertExpectations(t)
	})
	t.Run("Disconnect error", func(t *testing.T) {
		client := &mocks.DbClientBindingMock{}
		client.On("Database", mongo.DbName).Return(&mocks.DbBindingMock{})
		client.On("Disconnect", mock.Anything).Return(assert.AnError)
		conn := mongo.NewConnection(client, time.Duration(1))
		closeChannel := make(chan bool)
		defer test.AssertLogContains(t, "Error disconnecting from MongoDB")()
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-closeChannel
		}()
		conn.Shutdown(closeChannel)
		wg.Wait()
		client.AssertExpectations(t)
	})
}

func TestConnection_Collection(t *testing.T) {
	collectionName := test.AnyString
	client := &mocks.DbClientBindingMock{}
	db := &mocks.DbBindingMock{}
	client.On("Database", mongo.DbName).Return(db)
	db.On("Collection", collectionName).Return(&mocks.CollectionBindingMock{})
	conn := mongo.NewConnection(client, time.Duration(1))
	assert.NotNil(t, conn.Collection(collectionName))
}

// withUpsert returns a testify matcher that verifies an options builder has Upsert=true
func withUpsert() interface{} {
	return mock.MatchedBy(func(opt options.Lister[options.ReplaceOptions]) bool {
		resolved := &options.ReplaceOptions{}
		for _, fn := range opt.List() {
			if err := fn(resolved); err != nil {
				return false
			}
		}
		return resolved.Upsert != nil && *resolved.Upsert
	})
}

// sortedBy returns a testify matcher that verifies a Find options builder has the given sort document
func sortedBy(sort any) interface{} {
	return mock.MatchedBy(func(opt options.Lister[options.FindOptions]) bool {
		resolved := &options.FindOptions{}
		for _, fn := range opt.List() {
			if err := fn(resolved); err != nil {
				return false
			}
		}
		return reflect.DeepEqual(resolved.Sort, sort)
	})
}

func assertDbInteractionLog(t *testing.T, interaction mongo.DbInteraction) (assertFunc func() bool) {
	return test.AssertLogContains(t, string(interaction))
}

func getClientAndCollectionMocks(collectionName string) (*mocks.DbClientBindingMock, *mocks.CollectionBindingMock) {
	client := &mocks.DbClientBindingMock{}
	db := &mocks.DbBindingMock{}
	client.On("Database", mongo.DbName).Return(db)
	collection := &mocks.CollectionBindingMock{}
	db.On("Collection", collectionName).Return(collection)
	return client, collection
}

func getClientAndDatabaseMocks() (*mocks.DbClientBindingMock, *mocks.DbBindingMock) {
	client := &mocks.DbClientBindingMock{}
	db := &mocks.DbBindingMock{}
	client.On("Database", mongo.DbName).Return(db)
	return client, db
}
