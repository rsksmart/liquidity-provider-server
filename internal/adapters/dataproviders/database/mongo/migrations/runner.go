package migrations

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:embed scripts/*.json
var migrationScripts embed.FS

const migrationsCollection = "schema_migrations"

// DatabaseBinding is the interface the migration runner needs from the database.
type DatabaseBinding interface {
	RunCommand(ctx context.Context, runCommand any, opts ...*options.RunCmdOptions) *mongo.SingleResult
	Collection(name string, opts ...*options.CollectionOptions) CollectionBinding
}

// CollectionBinding is the interface for collection operations needed by version tracking.
type CollectionBinding interface {
	FindOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult
	UpdateOne(ctx context.Context, filter any, update any,
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

// MongoDatabaseAdapter wraps *mongo.Database to satisfy DatabaseBinding.
type MongoDatabaseAdapter struct {
	DB *mongo.Database
}

func (a *MongoDatabaseAdapter) RunCommand(
	ctx context.Context, runCommand any, opts ...*options.RunCmdOptions,
) *mongo.SingleResult {
	return a.DB.RunCommand(ctx, runCommand, opts...)
}

func (a *MongoDatabaseAdapter) Collection(
	name string, opts ...*options.CollectionOptions,
) CollectionBinding {
	return a.DB.Collection(name, opts...)
}

type migrationRecord struct {
	Version int  `bson:"version"`
	Dirty   bool `bson:"dirty"`
}

type migrationStep struct {
	version  int
	name     string
	commands []bson.D
}

type migrationPair struct {
	version int
	name    string
	up      []bson.D
	down    []bson.D
}

// RunAll applies all pending up migrations.
func RunAll(ctx context.Context, db DatabaseBinding) error {
	pairs, err := loadMigrations()
	if err != nil {
		return err
	}
	if len(pairs) == 0 {
		log.Info("No migration scripts found")
		return nil
	}

	current, err := getCurrentVersion(ctx, db)
	if err != nil {
		return err
	}

	applied := 0
	for _, p := range pairs {
		if p.version <= current {
			continue
		}
		if err = applyUp(ctx, db, p); err != nil {
			return err
		}
		applied++
	}

	logResult(applied, current, pairs[len(pairs)-1].version)
	return nil
}

// Down rolls back the last applied migration.
func Down(ctx context.Context, db DatabaseBinding) error {
	current, err := getCurrentVersion(ctx, db)
	if err != nil {
		return err
	}
	if current == 0 {
		log.Info("No migrations to roll back")
		return nil
	}

	pairs, err := loadMigrations()
	if err != nil {
		return err
	}

	pair, found := findVersion(pairs, current)
	if !found {
		return fmt.Errorf("migration version %d not found in scripts", current)
	}
	if len(pair.down) == 0 {
		return fmt.Errorf("no down migration for version %d", current)
	}

	return applyDown(ctx, db, pair)
}

// MigrateTo migrates the database to the specified target version,
// applying up or down migrations as needed.
func MigrateTo(ctx context.Context, db DatabaseBinding, target int) error {
	pairs, err := loadMigrations()
	if err != nil {
		return err
	}

	current, err := getCurrentVersion(ctx, db)
	if err != nil {
		return err
	}

	if target == current {
		log.Infof("Database already at version %d", current)
		return nil
	}

	if target > current {
		return migrateUp(ctx, db, pairs, current, target)
	}
	return migrateDown(ctx, db, pairs, current, target)
}

func migrateUp(ctx context.Context, db DatabaseBinding, pairs []migrationPair, current, target int) error {
	for _, p := range pairs {
		if p.version <= current || p.version > target {
			continue
		}
		if err := applyUp(ctx, db, p); err != nil {
			return err
		}
	}
	log.Infof("Migrated up to version %d", target)
	return nil
}

func migrateDown(ctx context.Context, db DatabaseBinding, pairs []migrationPair, current, target int) error {
	for i := len(pairs) - 1; i >= 0; i-- {
		p := pairs[i]
		if p.version > current || p.version <= target {
			continue
		}
		if len(p.down) == 0 {
			return fmt.Errorf("no down migration for version %d, cannot roll back", p.version)
		}
		if err := applyDown(ctx, db, p); err != nil {
			return err
		}
	}
	log.Infof("Migrated down to version %d", target)
	return nil
}

func applyUp(ctx context.Context, db DatabaseBinding, p migrationPair) error {
	step := migrationStep{version: p.version, name: p.name, commands: p.up}
	log.Infof("Applying migration %06d_%s (up)", p.version, p.name)
	if err := runCommands(ctx, db, step); err != nil {
		return markDirty(ctx, db, p.version, err)
	}
	return setVersion(ctx, db, p.version)
}

func applyDown(ctx context.Context, db DatabaseBinding, p migrationPair) error {
	step := migrationStep{version: p.version, name: p.name, commands: p.down}
	log.Infof("Rolling back migration %06d_%s (down)", p.version, p.name)
	if err := runCommands(ctx, db, step); err != nil {
		return markDirty(ctx, db, p.version, err)
	}
	return removeVersion(ctx, db, p.version)
}

func runCommands(ctx context.Context, db DatabaseBinding, step migrationStep) error {
	for i, cmd := range step.commands {
		if err := db.RunCommand(ctx, cmd).Err(); err != nil {
			return fmt.Errorf("migration %06d_%s command %d failed: %w", step.version, step.name, i, err)
		}
	}
	return nil
}

func logResult(applied, previousVersion, latestVersion int) {
	if applied == 0 {
		log.Infof("Database already at version %d, no migrations needed", previousVersion)
	} else {
		log.Infof("Database migrations completed, applied %d migration(s), version: %d", applied, latestVersion)
	}
}


func loadMigrations() ([]migrationPair, error) {
	files, err := migrationScripts.ReadDir("scripts")
	if err != nil {
		return nil, fmt.Errorf("error reading migration scripts: %w", err)
	}

	pairMap := make(map[int]*migrationPair)
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		if err = assignMigrationFile(f.Name(), pairMap); err != nil {
			return nil, err
		}
	}

	pairs := make([]migrationPair, 0, len(pairMap))
	for _, p := range pairMap {
		pairs = append(pairs, *p)
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].version < pairs[j].version })
	return pairs, nil
}

var migrationSuffixes = map[string]string{
	".up.json":   "up",
	".down.json": "down",
}

func assignMigrationFile(fileName string, pairMap map[int]*migrationPair) error {
	suffix, direction := matchSuffix(fileName)
	if direction == "" {
		return nil
	}

	version, name, err := parseVersionAndName(strings.TrimSuffix(fileName, suffix), fileName)
	if err != nil {
		return err
	}

	commands, err := parseCommandsFromFile(fileName)
	if err != nil {
		return err
	}

	p, ok := pairMap[version]
	if !ok {
		p = &migrationPair{version: version, name: name}
		pairMap[version] = p
	}
	switch direction {
	case "up":
		p.up = commands
	case "down":
		p.down = commands
	}
	return nil
}

func matchSuffix(fileName string) (string, string) {
	for suffix, direction := range migrationSuffixes {
		if strings.HasSuffix(fileName, suffix) {
			return suffix, direction
		}
	}
	return "", ""
}

func parseVersionAndName(baseName, fileName string) (int, string, error) {
	parts := strings.SplitN(baseName, "_", 2)
	version, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", fmt.Errorf("invalid migration version in %s: %w", fileName, err)
	}
	name := ""
	if len(parts) > 1 {
		name = parts[1]
	}
	return version, name, nil
}

func parseCommandsFromFile(fileName string) ([]bson.D, error) {
	data, err := migrationScripts.ReadFile(path.Join("scripts", fileName))
	if err != nil {
		return nil, fmt.Errorf("error reading migration file %s: %w", fileName, err)
	}

	var rawCommands []json.RawMessage
	if err = json.Unmarshal(data, &rawCommands); err != nil {
		return nil, fmt.Errorf("error parsing migration file %s: %w", fileName, err)
	}

	commands := make([]bson.D, len(rawCommands))
	for i, raw := range rawCommands {
		if err = bson.UnmarshalExtJSON(raw, false, &commands[i]); err != nil {
			return nil, fmt.Errorf("error parsing command %d in %s: %w", i, fileName, err)
		}
	}
	return commands, nil
}

func findVersion(pairs []migrationPair, version int) (migrationPair, bool) {
	for _, p := range pairs {
		if p.version == version {
			return p, true
		}
	}
	return migrationPair{}, false
}


func getCurrentVersion(ctx context.Context, db DatabaseBinding) (int, error) {
	var record migrationRecord
	err := db.Collection(migrationsCollection).FindOne(ctx, bson.D{},
		options.FindOne().SetSort(bson.D{primitive.E{Key: "version", Value: -1}}),
	).Decode(&record)
	if err == mongo.ErrNoDocuments {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("error reading migration version: %w", err)
	}
	if record.Dirty {
		return 0, fmt.Errorf("database migration version %d is dirty, manual intervention required", record.Version)
	}
	return record.Version, nil
}

func setVersion(ctx context.Context, db DatabaseBinding, version int) error {
	result, err := db.Collection(migrationsCollection).UpdateOne(ctx,
		bson.D{primitive.E{Key: "version", Value: version}},
		bson.D{primitive.E{Key: "$set", Value: migrationRecord{Version: version, Dirty: false}}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("error recording migration version %d: %w", version, err)
	}
	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		return fmt.Errorf("migration version %d was not recorded: no documents modified or upserted", version)
	}
	return nil
}

func removeVersion(ctx context.Context, db DatabaseBinding, version int) error {
	result, err := db.Collection(migrationsCollection).DeleteOne(
		ctx, bson.D{primitive.E{Key: "version", Value: version}},
	)
	if err != nil {
		return fmt.Errorf("error removing migration version %d: %w", version, err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("migration version %d was not found for removal", version)
	}
	return nil
}

func markDirty(ctx context.Context, db DatabaseBinding, version int, migrationErr error) error {
	result, err := db.Collection(migrationsCollection).UpdateOne(ctx,
		bson.D{primitive.E{Key: "version", Value: version}},
		bson.D{primitive.E{Key: "$set", Value: migrationRecord{Version: version, Dirty: true}}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Errorf("Failed to mark migration %d as dirty: %v", version, err)
	} else if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		log.Errorf("Failed to mark migration %d as dirty: no documents modified or upserted", version)
	}
	return migrationErr
}
