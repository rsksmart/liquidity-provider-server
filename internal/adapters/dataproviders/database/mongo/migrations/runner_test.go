package migrations

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mockDatabase struct {
	mock.Mock
}

func (m *mockDatabase) RunCommand(ctx context.Context, runCommand any, opts ...*options.RunCmdOptions) *mongo.SingleResult {
	args := m.Called(ctx, runCommand)
	result, ok := args.Get(0).(*mongo.SingleResult)
	if !ok {
		return mongo.NewSingleResultFromDocument(bson.D{}, assert.AnError, nil)
	}
	return result
}

func (m *mockDatabase) Collection(name string, opts ...*options.CollectionOptions) CollectionBinding {
	args := m.Called(name)
	collection, ok := args.Get(0).(CollectionBinding)
	if !ok {
		return &mockCollection{}
	}
	return collection
}

type mockCollection struct {
	mock.Mock
}

func (m *mockCollection) FindOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter, opts)
	result, ok := args.Get(0).(*mongo.SingleResult)
	if !ok {
		return mongo.NewSingleResultFromDocument(bson.D{}, assert.AnError, nil)
	}
	return result
}

func (m *mockCollection) UpdateOne(
	ctx context.Context, filter any, update any, opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update, opts)
	result, ok := args.Get(0).(*mongo.UpdateResult)
	if !ok {
		return nil, args.Error(1)
	}
	return result, args.Error(1)
}

func newSuccessResult() *mongo.SingleResult {
	return mongo.NewSingleResultFromDocument(bson.D{{Key: "ok", Value: 1}}, nil, nil)
}

func newVersionResult(version int, dirty bool) *mongo.SingleResult {
	return mongo.NewSingleResultFromDocument(
		migrationRecord{Version: version, Dirty: dirty}, nil, nil,
	)
}

func newNoDocumentsResult() *mongo.SingleResult {
	return mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
}

func migrationStateFilter() bson.D {
	return bson.D{{Key: "_id", Value: migrationStateDocumentID}}
}

func setupFreshDbMocks(db *mockDatabase) *mockCollection {
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("FindOne", mock.Anything, migrationStateFilter(), mock.Anything).
		Return(newNoDocumentsResult()).Once()
	return col
}

func setupVersionedDbMocks(db *mockDatabase, version int) *mockCollection {
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("FindOne", mock.Anything, migrationStateFilter(), mock.Anything).
		Return(newVersionResult(version, false)).Once()
	return col
}

func expectSetVersion(col *mockCollection, version int) {
	filter := migrationStateFilter()
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "version", Value: version},
		{Key: "dirty", Value: false},
	}}}
	col.On("UpdateOne", mock.Anything, filter, update, mock.Anything).
		Return(&mongo.UpdateResult{UpsertedCount: 1}, nil).Once()
}

func TestLoadMigrations(t *testing.T) {
	pairs, err := NewRunner(nil).loadMigrations()
	require.NoError(t, err)
	require.NotEmpty(t, pairs)

	t.Run("parses version and name from embedded scripts", func(t *testing.T) {
		assert.Equal(t, 1, pairs[0].version)
		assert.Equal(t, "bridge_rebalances", pairs[0].name)
	})

	t.Run("loads up commands", func(t *testing.T) {
		assert.NotEmpty(t, pairs[0].up, "up commands should be loaded")
	})

	t.Run("migrations are sorted by version", func(t *testing.T) {
		for i := 1; i < len(pairs); i++ {
			assert.Greater(t, pairs[i].version, pairs[i-1].version)
		}
	})
}

func TestParseVersionAndName(t *testing.T) {
	t.Run("parses version and name", func(t *testing.T) {
		v, n, err := NewRunner(nil).parseVersionAndName("000001_bridge_rebalances.up.json")
		require.NoError(t, err)
		assert.Equal(t, 1, v)
		assert.Equal(t, "bridge_rebalances", n)
	})

	t.Run("parses version without name", func(t *testing.T) {
		v, n, err := NewRunner(nil).parseVersionAndName("000002.up.json")
		require.NoError(t, err)
		assert.Equal(t, 2, v)
		assert.Equal(t, "", n)
	})

	t.Run("returns error when version prefix is not six digits", func(t *testing.T) {
		_, _, err := NewRunner(nil).parseVersionAndName("abc_migration.up.json")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "6 digits")
	})

	t.Run("returns error when version prefix is too short", func(t *testing.T) {
		_, _, err := NewRunner(nil).parseVersionAndName("00001_x.up.json")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "6 digits")
	})
}

func TestRunAll_FreshDatabase(t *testing.T) {
	db := &mockDatabase{}
	col := setupFreshDbMocks(db)

	db.On("RunCommand", mock.Anything, mock.Anything).Return(newSuccessResult())
	expectSetVersion(col, 1)

	err := NewRunner(db).RunAll(context.Background())
	require.NoError(t, err)
	db.AssertExpectations(t)
	col.AssertExpectations(t)
}

func TestRunAll_AlreadyAtLatestVersion(t *testing.T) {
	db := &mockDatabase{}
	col := setupVersionedDbMocks(db, 1)

	err := NewRunner(db).RunAll(context.Background())
	require.NoError(t, err)
	db.AssertExpectations(t)
	col.AssertExpectations(t)
}

func TestRunAll_DirtyVersion(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("FindOne", mock.Anything, migrationStateFilter(), mock.Anything).
		Return(newVersionResult(1, true)).Once()

	err := NewRunner(db).RunAll(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "dirty")
	col.AssertExpectations(t)
}

func TestRunAll_CommandFails_MarksDirty(t *testing.T) {
	db := &mockDatabase{}
	col := setupFreshDbMocks(db)

	cmdErr := mongo.NewSingleResultFromDocument(bson.D{}, assert.AnError, nil)
	db.On("RunCommand", mock.Anything, mock.Anything).Return(cmdErr)

	dirtyFilter := migrationStateFilter()
	dirtyUpdate := bson.D{{Key: "$set", Value: bson.D{
		{Key: "version", Value: 1},
		{Key: "dirty", Value: true},
	}}}
	col.On("UpdateOne", mock.Anything, dirtyFilter, dirtyUpdate, mock.Anything).
		Return(&mongo.UpdateResult{UpsertedCount: 1}, nil).Once()

	err := NewRunner(db).RunAll(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "command 1 failed")
	col.AssertExpectations(t)
}

func TestGetCurrentVersion_FreshDatabase(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("FindOne", mock.Anything, migrationStateFilter(), mock.Anything).
		Return(newNoDocumentsResult()).Once()

	version, err := NewRunner(db).getCurrentVersion(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 0, version)
}

func TestGetCurrentVersion_ExistingVersion(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("FindOne", mock.Anything, migrationStateFilter(), mock.Anything).
		Return(newVersionResult(3, false)).Once()

	version, err := NewRunner(db).getCurrentVersion(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 3, version)
}

func TestGetCurrentVersion_DirtyState(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("FindOne", mock.Anything, migrationStateFilter(), mock.Anything).
		Return(newVersionResult(2, true)).Once()

	_, err := NewRunner(db).getCurrentVersion(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "dirty")
	assert.Contains(t, err.Error(), "manual intervention")
}

func TestSetVersion(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)

	filter := migrationStateFilter()
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "version", Value: 5},
		{Key: "dirty", Value: false},
	}}}
	col.On("UpdateOne", mock.Anything, filter, update, mock.Anything).
		Return(&mongo.UpdateResult{UpsertedCount: 1}, nil).Once()

	err := NewRunner(db).setVersion(context.Background(), 5)
	require.NoError(t, err)
	col.AssertExpectations(t)
}

func TestSetVersion_Error(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&mongo.UpdateResult{}, assert.AnError).Once()

	err := NewRunner(db).setVersion(context.Background(), 5)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error recording migration version")
}

func TestSetVersion_NoDocumentsAffected(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&mongo.UpdateResult{ModifiedCount: 0, UpsertedCount: 0}, nil).Once()

	err := NewRunner(db).setVersion(context.Background(), 5)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "was not recorded")
}

func TestMarkDirty(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)

	filter := migrationStateFilter()
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "version", Value: 2},
		{Key: "dirty", Value: true},
	}}}
	col.On("UpdateOne", mock.Anything, filter, update, mock.Anything).
		Return(&mongo.UpdateResult{UpsertedCount: 1}, nil).Once()

	originalErr := assert.AnError
	err := NewRunner(db).markDirty(context.Background(), 2, originalErr)
	assert.Equal(t, originalErr, err, "should return the original migration error")
	col.AssertExpectations(t)
}

func TestMarkDirty_UpdateErrorJoinsMigrationError(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return((*mongo.UpdateResult)(nil), assert.AnError).Once()

	originalErr := errors.New("migration failed")
	err := NewRunner(db).markDirty(context.Background(), 2, originalErr)
	require.Error(t, err)
	require.ErrorIs(t, err, originalErr)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestAssignMigrationFile_DuplicateVersion(t *testing.T) {
	r := NewRunner(nil)
	m := make(map[int]migration)
	require.NoError(t, r.assignMigrationFile("000001_bridge_rebalances.up.json", m))
	err := r.assignMigrationFile("000001_other_name.up.json", m)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate migration version 1")
}

func TestValidateMigrationSequence(t *testing.T) {
	t.Run("accepts contiguous from 1", func(t *testing.T) {
		err := validateMigrationSequence([]migration{
			{version: 1, fileName: "a"},
			{version: 2, fileName: "b"},
		})
		require.NoError(t, err)
	})

	t.Run("rejects gap", func(t *testing.T) {
		err := validateMigrationSequence([]migration{
			{version: 1, fileName: "a"},
			{version: 3, fileName: "c"},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "contiguous")
		assert.Contains(t, err.Error(), "expected version 2, found 3")
		assert.Contains(t, err.Error(), `"c"`)
	})
}
