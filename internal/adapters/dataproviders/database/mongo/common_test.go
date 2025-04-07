package mongo_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

type TestStoredQuote struct {
	Hash      string
	TestQuote TestQuote
}

type TestQuote struct {
	Value int
}

type TestRetainedQuote struct {
	QuoteHash string
	State     string
}

type errorCursor struct {
	*mongodriver.Cursor
	err error
}

func (c *errorCursor) All(ctx context.Context, results interface{}) error {
	return c.err
}

type mockCursor struct {
	err     error
	docs    []interface{}
	current int
}

func newMockCursor(err error, docs []interface{}) *mockCursor {
	return &mockCursor{
		err:     err,
		docs:    docs,
		current: -1,
	}
}

func (m *mockCursor) ID() int64 { return 0 }

func (m *mockCursor) Next(ctx context.Context) bool {
	if m.err != nil {
		return false
	}
	m.current++
	return m.current < len(m.docs)
}

func (m *mockCursor) Decode(val interface{}) error {
	if m.err != nil {
		return m.err
	}
	if m.current < 0 || m.current >= len(m.docs) {
		return fmt.Errorf("no document to decode")
	}
	data, err := bson.Marshal(m.docs[m.current])
	if err != nil {
		return err
	}
	return bson.Unmarshal(data, val)
}

func (m *mockCursor) Err() error { return m.err }

func (m *mockCursor) Close(ctx context.Context) error { return nil }

func (m *mockCursor) All(ctx context.Context, results interface{}) error {
	if m.err != nil {
		return m.err
	}
	data, err := bson.Marshal(m.docs)
	if err != nil {
		return err
	}
	return bson.Unmarshal(data, results)
}

func TestListQuotesByDateRange(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	startTimestamp := startDate.Unix()
	endTimestamp := endDate.Unix()

	t.Run("successful retrieval of quotes", func(t *testing.T) {
		client := &mocks.DbClientBindingMock{}
		db := &mocks.DbBindingMock{}
		client.On("Database", mongo.DbName).Return(db)

		quoteCollection := &mocks.CollectionBindingMock{}
		retainedCollection := &mocks.CollectionBindingMock{}
		db.On("Collection", "quoteCollection").Return(quoteCollection)
		db.On("Collection", "retainedCollection").Return(retainedCollection)

		storedQuotes := []TestStoredQuote{
			{Hash: "hash1", TestQuote: TestQuote{Value: 1}},
			{Hash: "hash2", TestQuote: TestQuote{Value: 2}},
		}
		retainedQuotes := []TestRetainedQuote{
			{QuoteHash: "hash1", State: "state1"},
			{QuoteHash: "hash2", State: "state2"},
		}

		quoteFilter := bson.D{
			{Key: "agreement_timestamp", Value: bson.D{
				{Key: "$gte", Value: startTimestamp},
				{Key: "$lte", Value: endTimestamp},
			}},
		}

		// Convert storedQuotes to []interface{}
		storedQuotesInterface := make([]interface{}, len(storedQuotes))
		for i, v := range storedQuotes {
			storedQuotesInterface[i] = v
		}
		quoteCursor, err := mongodriver.NewCursorFromDocuments(storedQuotesInterface, nil, nil)
		require.NoError(t, err)
		quoteCollection.On("Find", mock.Anything, quoteFilter).Return(quoteCursor, nil)

		retainedFilter := bson.D{
			{Key: "quote_hash", Value: bson.D{
				{Key: "$in", Value: []string{"hash1", "hash2"}},
			}},
		}

		// Convert retainedQuotes to []interface{}
		retainedQuotesInterface := make([]interface{}, len(retainedQuotes))
		for i, v := range retainedQuotes {
			retainedQuotesInterface[i] = v
		}
		retainedCursor, err := mongodriver.NewCursorFromDocuments(retainedQuotesInterface, nil, nil)
		require.NoError(t, err)
		retainedCollection.On("Find", mock.Anything, retainedFilter).Return(retainedCursor, nil)

		conn := mongo.NewConnection(client, time.Duration(1))
		quotes, retained, err := mongo.ListQuotesByDateRange[TestStoredQuote, TestQuote, TestRetainedQuote](
			context.Background(),
			conn,
			startDate,
			endDate,
			"quoteCollection",
			"retainedCollection",
			func(stored TestStoredQuote) (string, TestQuote) {
				return stored.Hash, stored.TestQuote
			},
		)

		require.NoError(t, err)
		assert.Len(t, quotes, 2)
		assert.Len(t, retained, 2)
		assert.Equal(t, 1, quotes[0].Value)
		assert.Equal(t, 2, quotes[1].Value)
		assert.Equal(t, "state1", retained[0].State)
		assert.Equal(t, "state2", retained[1].State)
	})

	t.Run("empty result set", func(t *testing.T) {
		client := &mocks.DbClientBindingMock{}
		db := &mocks.DbBindingMock{}
		client.On("Database", mongo.DbName).Return(db)

		quoteCollection := &mocks.CollectionBindingMock{}
		retainedCollection := &mocks.CollectionBindingMock{}
		db.On("Collection", "quoteCollection").Return(quoteCollection)
		db.On("Collection", "retainedCollection").Return(retainedCollection)

		quoteFilter := bson.D{
			{Key: "agreement_timestamp", Value: bson.D{
				{Key: "$gte", Value: startTimestamp},
				{Key: "$lte", Value: endTimestamp},
			}},
		}
		quoteCursor, err := mongodriver.NewCursorFromDocuments([]interface{}{}, nil, nil)
		require.NoError(t, err)
		quoteCollection.On("Find", mock.Anything, quoteFilter).Return(quoteCursor, nil)

		retainedFilter := bson.D{
			{Key: "quote_hash", Value: bson.D{
				{Key: "$in", Value: []string{}},
			}},
		}
		retainedCursor, err := mongodriver.NewCursorFromDocuments([]interface{}{}, nil, nil)
		require.NoError(t, err)
		retainedCollection.On("Find", mock.Anything, retainedFilter).Return(retainedCursor, nil)

		conn := mongo.NewConnection(client, time.Duration(1))
		quotes, retained, err := mongo.ListQuotesByDateRange[TestStoredQuote, TestQuote, TestRetainedQuote](
			context.Background(),
			conn,
			startDate,
			endDate,
			"quoteCollection",
			"retainedCollection",
			func(stored TestStoredQuote) (string, TestQuote) {
				return stored.Hash, stored.TestQuote
			},
		)

		require.NoError(t, err)
		assert.Empty(t, quotes)
		assert.Empty(t, retained)
	})

	t.Run("database error on quote collection", func(t *testing.T) {
		client := &mocks.DbClientBindingMock{}
		db := &mocks.DbBindingMock{}
		client.On("Database", mongo.DbName).Return(db)

		quoteCollection := &mocks.CollectionBindingMock{}
		db.On("Collection", "quoteCollection").Return(quoteCollection)
		quoteCollection.On("Find", mock.Anything, mock.Anything).Return(nil, assert.AnError)

		conn := mongo.NewConnection(client, time.Duration(1))
		quotes, retained, err := mongo.ListQuotesByDateRange[TestStoredQuote, TestQuote, TestRetainedQuote](
			context.Background(),
			conn,
			startDate,
			endDate,
			"quoteCollection",
			"retainedCollection",
			func(stored TestStoredQuote) (string, TestQuote) {
				return stored.Hash, stored.TestQuote
			},
		)

		require.Error(t, err)
		assert.Nil(t, quotes)
		assert.Nil(t, retained)
	})

	t.Run("database error on retained collection", func(t *testing.T) {
		client := &mocks.DbClientBindingMock{}
		db := &mocks.DbBindingMock{}
		client.On("Database", mongo.DbName).Return(db)

		quoteCollection := &mocks.CollectionBindingMock{}
		retainedCollection := &mocks.CollectionBindingMock{}
		db.On("Collection", "quoteCollection").Return(quoteCollection)
		db.On("Collection", "retainedCollection").Return(retainedCollection)

		storedQuotes := []TestStoredQuote{
			{Hash: "hash1", TestQuote: TestQuote{Value: 1}},
		}
		quoteFilter := bson.D{
			{Key: "agreement_timestamp", Value: bson.D{
				{Key: "$gte", Value: startTimestamp},
				{Key: "$lte", Value: endTimestamp},
			}},
		}

		storedQuotesInterface := make([]interface{}, len(storedQuotes))
		for i, v := range storedQuotes {
			storedQuotesInterface[i] = v
		}
		quoteCursor, err := mongodriver.NewCursorFromDocuments(storedQuotesInterface, nil, nil)
		require.NoError(t, err)
		quoteCollection.On("Find", mock.Anything, quoteFilter).Return(quoteCursor, nil)

		retainedFilter := bson.D{
			{Key: "quote_hash", Value: bson.D{
				{Key: "$in", Value: []string{"hash1"}},
			}},
		}
		retainedCollection.On("Find", mock.Anything, retainedFilter).Return(nil, assert.AnError)

		conn := mongo.NewConnection(client, time.Duration(1))
		quotes, retained, err := mongo.ListQuotesByDateRange[TestStoredQuote, TestQuote, TestRetainedQuote](
			context.Background(),
			conn,
			startDate,
			endDate,
			"quoteCollection",
			"retainedCollection",
			func(stored TestStoredQuote) (string, TestQuote) {
				return stored.Hash, stored.TestQuote
			},
		)

		require.Error(t, err)
		assert.Nil(t, quotes)
		assert.Nil(t, retained)
	})

	t.Run("error_in_quote_cursor_All", func(t *testing.T) {
		client := &mocks.DbClientBindingMock{}
		db := &mocks.DbBindingMock{}
		client.On("Database", mongo.DbName).Return(db)

		quoteCollection := &mocks.CollectionBindingMock{}
		retainedCollection := &mocks.CollectionBindingMock{}
		db.On("Collection", "quoteCollection").Return(quoteCollection)
		db.On("Collection", "retainedCollection").Return(retainedCollection)

		errorCursor := newMockCursor(assert.AnError, []interface{}{})

		quoteFilter := bson.D{
			{Key: "agreement_timestamp", Value: bson.D{
				{Key: "$gte", Value: startTimestamp},
				{Key: "$lte", Value: endTimestamp},
			}},
		}
		quoteCollection.On("Find", mock.Anything, quoteFilter).Return(errorCursor, nil)

		conn := mongo.NewConnection(client, time.Duration(1))
		quotes, retained, err := mongo.ListQuotesByDateRange[TestStoredQuote, TestQuote, TestRetainedQuote](
			context.Background(),
			conn,
			startDate,
			endDate,
			"quoteCollection",
			"retainedCollection",
			func(stored TestStoredQuote) (string, TestQuote) {
				return stored.Hash, stored.TestQuote
			},
		)

		assert.Error(t, err)
		assert.Nil(t, quotes)
		assert.Nil(t, retained)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("error_in_retained_cursor_All", func(t *testing.T) {
		client := &mocks.DbClientBindingMock{}
		db := &mocks.DbBindingMock{}
		client.On("Database", mongo.DbName).Return(db)

		quoteCollection := &mocks.CollectionBindingMock{}
		retainedCollection := &mocks.CollectionBindingMock{}
		db.On("Collection", "quoteCollection").Return(quoteCollection)
		db.On("Collection", "retainedCollection").Return(retainedCollection)

		storedQuotes := []TestStoredQuote{
			{Hash: "hash1", TestQuote: TestQuote{Value: 1}},
		}
		quoteCursor, err := mongodriver.NewCursorFromDocuments([]interface{}{storedQuotes[0]}, nil, nil)
		assert.NoError(t, err)

		quoteFilter := bson.D{
			{Key: "agreement_timestamp", Value: bson.D{
				{Key: "$gte", Value: startTimestamp},
				{Key: "$lte", Value: endTimestamp},
			}},
		}
		quoteCollection.On("Find", mock.Anything, quoteFilter).Return(quoteCursor, nil)

		retainedFilter := bson.D{
			{Key: "quote_hash", Value: bson.D{
				{Key: "$in", Value: []string{"hash1"}},
			}},
		}
		retainedCollection.On("Find", mock.Anything, retainedFilter).Return(nil, assert.AnError)

		conn := mongo.NewConnection(client, time.Duration(1))
		quotes, retained, err := mongo.ListQuotesByDateRange[TestStoredQuote, TestQuote, TestRetainedQuote](
			context.Background(),
			conn,
			startDate,
			endDate,
			"quoteCollection",
			"retainedCollection",
			func(stored TestStoredQuote) (string, TestQuote) {
				return stored.Hash, stored.TestQuote
			},
		)

		assert.Error(t, err)
		assert.Nil(t, quotes)
		assert.Nil(t, retained)
		assert.Equal(t, assert.AnError, err)
	})
}
