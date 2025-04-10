package mongo_test

import (
	"context"
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

func createCursorFromList[T any](t *testing.T, documents []T) *mongodriver.Cursor {
	docsInterface := make([]interface{}, len(documents))
	for i, v := range documents {
		docsInterface[i] = v
	}
	cursor, err := mongodriver.NewCursorFromDocuments(docsInterface, nil, nil)
	assert.NoError(t, err)
	return cursor
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

func (r TestRetainedQuote) GetQuoteHash() string {
	return r.QuoteHash
}

func setupQuoteTestData() (startDate, endDate time.Time, startTimestamp, endTimestamp int64) {
	startDate = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate = time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	startTimestamp = startDate.Unix()
	endTimestamp = endDate.Unix()
	return
}

func setupTestCollections() (*mocks.DbClientBindingMock, *mocks.CollectionBindingMock, *mocks.CollectionBindingMock) {
	client, db := getClientAndDatabaseMocks()
	quoteCollection := &mocks.CollectionBindingMock{}
	retainedCollection := &mocks.CollectionBindingMock{}
	db.On("Collection", "quoteCollection").Return(quoteCollection)
	db.On("Collection", "retainedCollection").Return(retainedCollection)
	return client, quoteCollection, retainedCollection
}

func createQuoteFilter(startTimestamp, endTimestamp int64) bson.D {
	return bson.D{
		{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startTimestamp},
			{Key: "$lte", Value: endTimestamp},
		}},
	}
}

func createRetainedFilter(quoteHashes []string, startTimestamp, endTimestamp int64) bson.D {
	return bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "quote_hash", Value: bson.D{
				{Key: "$in", Value: quoteHashes},
			}}},
			bson.D{
				{Key: "created_at", Value: bson.D{
					{Key: "$gte", Value: startTimestamp},
					{Key: "$lte", Value: endTimestamp},
				}},
			},
		}},
	}
}

func createTestQuery(client *mocks.DbClientBindingMock, startDate, endDate time.Time) mongo.QuoteQuery {
	conn := mongo.NewConnection(client, time.Duration(1))
	return mongo.QuoteQuery{
		Ctx:                context.Background(),
		Conn:               conn,
		StartDate:          startDate,
		EndDate:            endDate,
		QuoteCollection:    "quoteCollection",
		RetainedCollection: "retainedCollection",
	}
}

func createQuoteExtractor() func(doc bson.D) TestQuote {
	return func(doc bson.D) TestQuote {
		var stored TestStoredQuote
		bsonBytes, err := bson.Marshal(doc)
		if err != nil {
			return TestQuote{}
		}
		if err := bson.Unmarshal(bsonBytes, &stored); err != nil {
			return TestQuote{}
		}
		return stored.TestQuote
	}
}

func TestListQuotesByDateRange_SuccessfulRetrieval(t *testing.T) {
	startDate, endDate, startTimestamp, endTimestamp := setupQuoteTestData()
	client, quoteCollection, retainedCollection := setupTestCollections()
	storedQuotes := []TestStoredQuote{
		{Hash: "hash1", TestQuote: TestQuote{Value: 1}},
		{Hash: "hash2", TestQuote: TestQuote{Value: 2}},
	}
	retainedQuotes := []TestRetainedQuote{
		{QuoteHash: "hash1", State: "state1"},
		{QuoteHash: "hash2", State: "state2"},
	}
	quoteFilter := createQuoteFilter(startTimestamp, endTimestamp)
	quoteCursor := createCursorFromList(t, storedQuotes)
	quoteCollection.On("Find", mock.Anything, quoteFilter).Return(quoteCursor, nil)
	retainedFilter := createRetainedFilter([]string{"hash1", "hash2"}, startTimestamp, endTimestamp)
	retainedCursor := createCursorFromList(t, retainedQuotes)
	retainedCollection.On("Find", mock.Anything, retainedFilter).Return(retainedCursor, nil)
	query := createTestQuery(client, startDate, endDate)
	result, err := mongo.ListQuotesByDateRange[TestQuote, TestRetainedQuote](query, createQuoteExtractor())
	require.NoError(t, err)
	assert.Len(t, result.Quotes, 2)
	assert.Len(t, result.RetainedQuotes, 2)
	assert.Equal(t, TestQuote{Value: 1}, result.Quotes[0])
	assert.Equal(t, TestQuote{Value: 2}, result.Quotes[1])
	assert.Equal(t, TestRetainedQuote{QuoteHash: "hash1", State: "state1"}, result.RetainedQuotes[0])
	assert.Equal(t, TestRetainedQuote{QuoteHash: "hash2", State: "state2"}, result.RetainedQuotes[1])
}

func TestListQuotesByDateRange_EmptyResultSet(t *testing.T) {
	startDate, endDate, startTimestamp, endTimestamp := setupQuoteTestData()
	client, quoteCollection, retainedCollection := setupTestCollections()
	quoteFilter := createQuoteFilter(startTimestamp, endTimestamp)
	quoteCursor := createCursorFromList(t, []TestStoredQuote{})
	quoteCollection.On("Find", mock.Anything, quoteFilter).Return(quoteCursor, nil)
	retainedFilter := createRetainedFilter([]string{}, startTimestamp, endTimestamp)
	retainedCursor := createCursorFromList(t, []TestRetainedQuote{})
	retainedCollection.On("Find", mock.Anything, retainedFilter).Return(retainedCursor, nil)
	query := createTestQuery(client, startDate, endDate)
	result, err := mongo.ListQuotesByDateRange[TestQuote, TestRetainedQuote](query, createQuoteExtractor())
	require.NoError(t, err)
	assert.Empty(t, result.Quotes)
	assert.Empty(t, result.RetainedQuotes)
}

func TestListQuotesByDateRange_DatabaseErrorOnQuoteCollection(t *testing.T) {
	startDate, endDate, startTimestamp, endTimestamp := setupQuoteTestData()
	client, quoteCollection, _ := setupTestCollections()
	quoteFilter := createQuoteFilter(startTimestamp, endTimestamp)
	quoteCollection.On("Find", mock.Anything, quoteFilter).Return(nil, assert.AnError)
	query := createTestQuery(client, startDate, endDate)
	result, err := mongo.ListQuotesByDateRange[TestQuote, TestRetainedQuote](query, createQuoteExtractor())
	require.Error(t, err)
	assert.Empty(t, result.Quotes)
	assert.Empty(t, result.RetainedQuotes)
}

func TestListQuotesByDateRange_DatabaseErrorOnRetainedCollection(t *testing.T) {
	startDate, endDate, startTimestamp, endTimestamp := setupQuoteTestData()
	client, quoteCollection, retainedCollection := setupTestCollections()
	storedQuotes := []TestStoredQuote{
		{Hash: "hash1", TestQuote: TestQuote{Value: 1}},
	}
	quoteFilter := createQuoteFilter(startTimestamp, endTimestamp)
	quoteCursor := createCursorFromList(t, storedQuotes)
	quoteCollection.On("Find", mock.Anything, quoteFilter).Return(quoteCursor, nil)
	retainedFilter := createRetainedFilter([]string{"hash1"}, startTimestamp, endTimestamp)
	retainedCollection.On("Find", mock.Anything, retainedFilter).Return(nil, assert.AnError)
	query := createTestQuery(client, startDate, endDate)
	result, err := mongo.ListQuotesByDateRange[TestQuote, TestRetainedQuote](query, createQuoteExtractor())
	require.Error(t, err)
	assert.Empty(t, result.Quotes)
	assert.Empty(t, result.RetainedQuotes)
}

func TestListQuotesByDateRange_FetchesAdditionalQuotes(t *testing.T) {
	startDate, endDate, startTimestamp, endTimestamp := setupQuoteTestData()
	client, quoteCollection, retainedCollection := setupTestCollections()
	storedQuotes := []TestStoredQuote{
		{Hash: "hash1", TestQuote: TestQuote{Value: 1}},
	}
	retainedQuotes := []TestRetainedQuote{
		{QuoteHash: "hash1", State: "state1"},
		{QuoteHash: "hash2", State: "state2"},
	}
	additionalQuotes := []TestStoredQuote{
		{Hash: "hash2", TestQuote: TestQuote{Value: 2}},
	}
	quoteFilter := createQuoteFilter(startTimestamp, endTimestamp)
	quoteCursor := createCursorFromList(t, storedQuotes)
	quoteCollection.On("Find", mock.Anything, quoteFilter).Return(quoteCursor, nil)
	retainedFilter := createRetainedFilter([]string{"hash1"}, startTimestamp, endTimestamp)
	retainedCursor := createCursorFromList(t, retainedQuotes)
	retainedCollection.On("Find", mock.Anything, retainedFilter).Return(retainedCursor, nil)
	additionalFilter := bson.D{
		{Key: "hash", Value: bson.D{
			{Key: "$in", Value: []string{"hash2"}},
		}},
	}
	additionalCursor := createCursorFromList(t, additionalQuotes)
	quoteCollection.On("Find", mock.Anything, additionalFilter).Return(additionalCursor, nil)
	query := createTestQuery(client, startDate, endDate)
	result, err := mongo.ListQuotesByDateRange[TestQuote, TestRetainedQuote](query, createQuoteExtractor())
	require.NoError(t, err)
	assert.Len(t, result.Quotes, 2)
	assert.Len(t, result.RetainedQuotes, 2)
	assert.Contains(t, []int{1, 2}, result.Quotes[0].Value)
	assert.Contains(t, []int{1, 2}, result.Quotes[1].Value)
	assert.Contains(t, []string{"state1", "state2"}, result.RetainedQuotes[0].State)
	assert.Contains(t, []string{"state1", "state2"}, result.RetainedQuotes[1].State)
}

func TestListQuotesByDateRange_HandlesErrorWhenFetchingAdditionalQuotes(t *testing.T) {
	startDate, endDate, startTimestamp, endTimestamp := setupQuoteTestData()
	client, quoteCollection, retainedCollection := setupTestCollections()
	storedQuotes := []TestStoredQuote{
		{Hash: "hash1", TestQuote: TestQuote{Value: 1}},
	}
	retainedQuotes := []TestRetainedQuote{
		{QuoteHash: "hash1", State: "state1"},
		{QuoteHash: "hash2", State: "state2"},
	}
	quoteFilter := createQuoteFilter(startTimestamp, endTimestamp)
	quoteCursor := createCursorFromList(t, storedQuotes)
	quoteCollection.On("Find", mock.Anything, quoteFilter).Return(quoteCursor, nil)
	retainedFilter := createRetainedFilter([]string{"hash1"}, startTimestamp, endTimestamp)
	retainedCursor := createCursorFromList(t, retainedQuotes)
	retainedCollection.On("Find", mock.Anything, retainedFilter).Return(retainedCursor, nil)
	additionalFilter := bson.D{
		{Key: "hash", Value: bson.D{
			{Key: "$in", Value: []string{"hash2"}},
		}},
	}
	quoteCollection.On("Find", mock.Anything, additionalFilter).Return(nil, assert.AnError)
	query := createTestQuery(client, startDate, endDate)
	result, err := mongo.ListQuotesByDateRange[TestQuote, TestRetainedQuote](query, createQuoteExtractor())
	require.NoError(t, err)
	assert.Len(t, result.Quotes, 1)
	assert.Len(t, result.RetainedQuotes, 2)
	assert.Equal(t, TestQuote{Value: 1}, result.Quotes[0])
	assert.Equal(t, "state1", result.RetainedQuotes[0].State)
	assert.Equal(t, "state2", result.RetainedQuotes[1].State)
}
