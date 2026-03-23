package migrations

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type mockDatabase struct {
	mock.Mock
}

func (m *mockDatabase) RunCommand(ctx context.Context, runCommand any, opts ...*options.RunCmdOptions) *mongo.SingleResult {
	args := m.Called(ctx, runCommand)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *mockDatabase) Collection(name string, opts ...*options.CollectionOptions) CollectionBinding {
	args := m.Called(name)
	return args.Get(0).(CollectionBinding)
}

type mockCollection struct {
	mock.Mock
}

func (m *mockCollection) FindOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *mockCollection) UpdateOne(
	ctx context.Context, filter any, update any, opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update, opts)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *mockCollection) DeleteOne(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}


func newSuccessResult() *mongo.SingleResult {
	return mongo.NewSingleResultFromDocument(bson.D{primitive.E{Key: "ok", Value: 1}}, nil, nil)
}

func newVersionResult(version int, dirty bool) *mongo.SingleResult {
	return mongo.NewSingleResultFromDocument(
		migrationRecord{Version: version, Dirty: dirty}, nil, nil,
	)
}

func newNoDocumentsResult() *mongo.SingleResult {
	return mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
}

func setupFreshDbMocks(db *mockDatabase) *mockCollection {
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("FindOne", mock.Anything, bson.D{}, mock.Anything).
		Return(newNoDocumentsResult()).Once()
	return col
}

func setupVersionedDbMocks(db *mockDatabase, version int) *mockCollection {
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("FindOne", mock.Anything, bson.D{}, mock.Anything).
		Return(newVersionResult(version, false)).Once()
	return col
}

func expectSetVersion(col *mockCollection, version int) {
	filter := bson.D{primitive.E{Key: "version", Value: version}}
	update := bson.D{primitive.E{Key: "$set", Value: migrationRecord{Version: version, Dirty: false}}}
	col.On("UpdateOne", mock.Anything, filter, update, mock.Anything).
		Return(&mongo.UpdateResult{UpsertedCount: 1}, nil).Once()
}

func expectRemoveVersion(col *mockCollection, version int) {
	filter := bson.D{primitive.E{Key: "version", Value: version}}
	col.On("DeleteOne", mock.Anything, filter, mock.Anything).
		Return(&mongo.DeleteResult{DeletedCount: 1}, nil).Once()
}


func TestLoadMigrations(t *testing.T) {
	pairs, err := loadMigrations()
	require.NoError(t, err)
	require.NotEmpty(t, pairs)

	t.Run("parses version and name from embedded scripts", func(t *testing.T) {
		assert.Equal(t, 1, pairs[0].version)
		assert.Equal(t, "bridge_rebalances", pairs[0].name)
	})

	t.Run("loads both up and down commands", func(t *testing.T) {
		assert.NotEmpty(t, pairs[0].up, "up commands should be loaded")
		assert.NotEmpty(t, pairs[0].down, "down commands should be loaded")
	})

	t.Run("migrations are sorted by version", func(t *testing.T) {
		for i := 1; i < len(pairs); i++ {
			assert.Greater(t, pairs[i].version, pairs[i-1].version)
		}
	})
}

func TestParseVersionAndName(t *testing.T) {
	t.Run("parses version and name", func(t *testing.T) {
		v, n, err := parseVersionAndName("000001_bridge_rebalances", "000001_bridge_rebalances.up.json")
		require.NoError(t, err)
		assert.Equal(t, 1, v)
		assert.Equal(t, "bridge_rebalances", n)
	})

	t.Run("parses version without name", func(t *testing.T) {
		v, n, err := parseVersionAndName("000002", "000002.up.json")
		require.NoError(t, err)
		assert.Equal(t, 2, v)
		assert.Equal(t, "", n)
	})

	t.Run("returns error for invalid version", func(t *testing.T) {
		_, _, err := parseVersionAndName("abc_migration", "abc_migration.up.json")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid migration version")
	})
}

func TestFindVersion(t *testing.T) {
	pairs := []migrationPair{
		{version: 1, name: "first"},
		{version: 2, name: "second"},
		{version: 3, name: "third"},
	}

	t.Run("finds existing version", func(t *testing.T) {
		p, found := findVersion(pairs, 2)
		require.True(t, found)
		assert.Equal(t, "second", p.name)
	})

	t.Run("returns false for missing version", func(t *testing.T) {
		_, found := findVersion(pairs, 99)
		assert.False(t, found)
	})
}


func TestRunAll_FreshDatabase(t *testing.T) {
	db := &mockDatabase{}
	col := setupFreshDbMocks(db)

	// Expect RunCommand for commands in all pending up migrations
	db.On("RunCommand", mock.Anything, mock.Anything).Return(newSuccessResult())
	expectSetVersion(col, 1)
	expectSetVersion(col, 2)

	err := RunAll(context.Background(), db)
	require.NoError(t, err)
	db.AssertExpectations(t)
	col.AssertExpectations(t)
}

func TestRunAll_AlreadyAtLatestVersion(t *testing.T) {
	db := &mockDatabase{}
	_ = setupVersionedDbMocks(db, 2)

	err := RunAll(context.Background(), db)
	require.NoError(t, err)
	db.AssertExpectations(t)
}

func TestRunAll_DirtyVersion(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("FindOne", mock.Anything, bson.D{}, mock.Anything).
		Return(newVersionResult(1, true)).Once()

	err := RunAll(context.Background(), db)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "dirty")
}

func TestRunAll_CommandFails_MarksDirty(t *testing.T) {
	db := &mockDatabase{}
	col := setupFreshDbMocks(db)

	cmdErr := mongo.NewSingleResultFromDocument(bson.D{}, assert.AnError, nil)
	db.On("RunCommand", mock.Anything, mock.Anything).Return(cmdErr)

	// Expect markDirty call
	dirtyFilter := bson.D{primitive.E{Key: "version", Value: 1}}
	dirtyUpdate := bson.D{primitive.E{Key: "$set", Value: migrationRecord{Version: 1, Dirty: true}}}
	col.On("UpdateOne", mock.Anything, dirtyFilter, dirtyUpdate, mock.Anything).
		Return(&mongo.UpdateResult{UpsertedCount: 1}, nil).Once()

	err := RunAll(context.Background(), db)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "command 0 failed")
	col.AssertExpectations(t)
}


func TestDown_RollsBackLastMigration(t *testing.T) {
	db := &mockDatabase{}
	col := setupVersionedDbMocks(db, 2)

	db.On("RunCommand", mock.Anything, mock.Anything).Return(newSuccessResult())
	expectRemoveVersion(col, 2)

	err := Down(context.Background(), db)
	require.NoError(t, err)
	db.AssertExpectations(t)
	col.AssertExpectations(t)
}

func TestDown_NothingToRollBack(t *testing.T) {
	db := &mockDatabase{}
	_ = setupFreshDbMocks(db)

	err := Down(context.Background(), db)
	require.NoError(t, err)
}

func TestDown_DirtyVersion(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("FindOne", mock.Anything, bson.D{}, mock.Anything).
		Return(newVersionResult(1, true)).Once()

	err := Down(context.Background(), db)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "dirty")
}

func TestDown_VersionNotInScripts(t *testing.T) {
	db := &mockDatabase{}
	_ = setupVersionedDbMocks(db, 999)

	err := Down(context.Background(), db)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found in scripts")
}


func TestMigrateTo_AlreadyAtTarget(t *testing.T) {
	db := &mockDatabase{}
	_ = setupVersionedDbMocks(db, 1)

	err := MigrateTo(context.Background(), db, 1)
	require.NoError(t, err)
}

func TestMigrateTo_UpFromFresh(t *testing.T) {
	db := &mockDatabase{}
	col := setupFreshDbMocks(db)

	db.On("RunCommand", mock.Anything, mock.Anything).Return(newSuccessResult())
	expectSetVersion(col, 1)

	err := MigrateTo(context.Background(), db, 1)
	require.NoError(t, err)
	db.AssertExpectations(t)
	col.AssertExpectations(t)
}

func TestMigrateTo_DownToZero(t *testing.T) {
	db := &mockDatabase{}
	col := setupVersionedDbMocks(db, 1)

	db.On("RunCommand", mock.Anything, mock.Anything).Return(newSuccessResult())
	expectRemoveVersion(col, 1)

	err := MigrateTo(context.Background(), db, 0)
	require.NoError(t, err)
	db.AssertExpectations(t)
	col.AssertExpectations(t)
}


func TestGetCurrentVersion_FreshDatabase(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("FindOne", mock.Anything, bson.D{}, mock.Anything).
		Return(newNoDocumentsResult()).Once()

	version, err := getCurrentVersion(context.Background(), db)
	require.NoError(t, err)
	assert.Equal(t, 0, version)
}

func TestGetCurrentVersion_ExistingVersion(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("FindOne", mock.Anything, bson.D{}, mock.Anything).
		Return(newVersionResult(3, false)).Once()

	version, err := getCurrentVersion(context.Background(), db)
	require.NoError(t, err)
	assert.Equal(t, 3, version)
}

func TestGetCurrentVersion_DirtyState(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("FindOne", mock.Anything, bson.D{}, mock.Anything).
		Return(newVersionResult(2, true)).Once()

	_, err := getCurrentVersion(context.Background(), db)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "dirty")
	assert.Contains(t, err.Error(), "manual intervention")
}

func TestSetVersion(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)

	filter := bson.D{primitive.E{Key: "version", Value: 5}}
	update := bson.D{primitive.E{Key: "$set", Value: migrationRecord{Version: 5, Dirty: false}}}
	col.On("UpdateOne", mock.Anything, filter, update, mock.Anything).
		Return(&mongo.UpdateResult{UpsertedCount: 1}, nil).Once()

	err := setVersion(context.Background(), db, 5)
	require.NoError(t, err)
	col.AssertExpectations(t)
}

func TestSetVersion_Error(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&mongo.UpdateResult{}, assert.AnError).Once()

	err := setVersion(context.Background(), db, 5)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error recording migration version")
}

func TestRemoveVersion(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)

	filter := bson.D{primitive.E{Key: "version", Value: 3}}
	col.On("DeleteOne", mock.Anything, filter, mock.Anything).
		Return(&mongo.DeleteResult{DeletedCount: 1}, nil).Once()

	err := removeVersion(context.Background(), db, 3)
	require.NoError(t, err)
	col.AssertExpectations(t)
}

func TestRemoveVersion_Error(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("DeleteOne", mock.Anything, mock.Anything, mock.Anything).
		Return(&mongo.DeleteResult{}, assert.AnError).Once()

	err := removeVersion(context.Background(), db, 3)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error removing migration version")
}

func TestSetVersion_NoDocumentsAffected(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&mongo.UpdateResult{ModifiedCount: 0, UpsertedCount: 0}, nil).Once()

	err := setVersion(context.Background(), db, 5)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "was not recorded")
}

func TestRemoveVersion_NoDocumentsDeleted(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)
	col.On("DeleteOne", mock.Anything, mock.Anything, mock.Anything).
		Return(&mongo.DeleteResult{DeletedCount: 0}, nil).Once()

	err := removeVersion(context.Background(), db, 3)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "was not found for removal")
}

func TestMarkDirty(t *testing.T) {
	db := &mockDatabase{}
	col := &mockCollection{}
	db.On("Collection", migrationsCollection).Return(col)

	filter := bson.D{primitive.E{Key: "version", Value: 2}}
	update := bson.D{primitive.E{Key: "$set", Value: migrationRecord{Version: 2, Dirty: true}}}
	col.On("UpdateOne", mock.Anything, filter, update, mock.Anything).
		Return(&mongo.UpdateResult{UpsertedCount: 1}, nil).Once()

	originalErr := assert.AnError
	err := markDirty(context.Background(), db, 2, originalErr)
	assert.Equal(t, originalErr, err, "should return the original migration error")
	col.AssertExpectations(t)
}
