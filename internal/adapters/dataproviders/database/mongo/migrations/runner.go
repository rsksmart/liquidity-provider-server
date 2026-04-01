package migrations

import (
	"cmp"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"slices"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:embed scripts/*.json
var migrationScripts embed.FS

const (
	migrationsCollection      = "schema_migrations"
	migrationStateDocumentID  = "migration_state"
	upSuffix                  = ".up.json"
	migrationVersionPrefixLen = 6
)

type (
	// DatabaseBinding is the interface the migration runner needs from the database.
	DatabaseBinding interface {
		RunCommand(ctx context.Context, runCommand any, opts ...*options.RunCmdOptions) *mongo.SingleResult
		Collection(name string, opts ...*options.CollectionOptions) CollectionBinding
	}

	// CollectionBinding is the interface for collection operations needed by version tracking.
	CollectionBinding interface {
		FindOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult
		UpdateOne(ctx context.Context, filter any, update any,
			opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	}

	migrationRecord struct {
		Version int  `bson:"version"`
		Dirty   bool `bson:"dirty"`
	}

	migration struct {
		version  int
		name     string
		fileName string
		up       []bson.D
	}
)

// Runner applies database migrations against a DatabaseBinding.
type Runner struct {
	db DatabaseBinding
}

func NewRunner(db DatabaseBinding) *Runner {
	return &Runner{db: db}
}

func (r *Runner) parseVersionAndName(fileName string) (int, string, error) {
	baseName := strings.TrimSuffix(fileName, upSuffix)
	parts := strings.SplitN(baseName, "_", 2)
	verStr := parts[0]
	if len(verStr) != migrationVersionPrefixLen {
		return 0, "", fmt.Errorf(
			"migration version prefix must be exactly %d digits, got %q in %s",
			migrationVersionPrefixLen, verStr, fileName,
		)
	}
	for _, c := range verStr {
		if c < '0' || c > '9' {
			return 0, "", fmt.Errorf("migration version must be decimal digits, got %q in %s", verStr, fileName)
		}
	}
	version, err := strconv.Atoi(verStr)
	if err != nil {
		return 0, "", fmt.Errorf("invalid migration version in %s: %w", fileName, err)
	}
	name := ""
	if len(parts) > 1 {
		name = parts[1]
	}
	return version, name, nil
}

func (r *Runner) parseCommandsFromFile(fileName string) ([]bson.D, error) {
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
			return nil, fmt.Errorf("error parsing command %d in %s: %w", i+1, fileName, err)
		}
	}
	return commands, nil
}

func (r *Runner) assignMigrationFile(fileName string, byVersion map[int]migration) error {
	if !strings.HasSuffix(fileName, upSuffix) {
		return nil
	}

	version, name, err := r.parseVersionAndName(fileName)
	if err != nil {
		return err
	}

	if prev, dup := byVersion[version]; dup {
		return fmt.Errorf(
			"duplicate migration version %d: %q and %q",
			version, prev.fileName, fileName,
		)
	}

	commands, err := r.parseCommandsFromFile(fileName)
	if err != nil {
		return err
	}

	byVersion[version] = migration{version: version, name: name, fileName: fileName, up: commands}
	return nil
}

func (r *Runner) loadMigrations() ([]migration, error) {
	files, err := migrationScripts.ReadDir("scripts")
	if err != nil {
		return nil, fmt.Errorf("error reading migration scripts: %w", err)
	}

	byVersion := make(map[int]migration)
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		if err = r.assignMigrationFile(f.Name(), byVersion); err != nil {
			return nil, err
		}
	}

	ordered := make([]migration, 0, len(byVersion))
	for _, m := range byVersion {
		ordered = append(ordered, m)
	}
	slices.SortFunc(ordered, func(a, b migration) int { return cmp.Compare(a.version, b.version) })
	if err := validateMigrationSequence(ordered); err != nil {
		return nil, err
	}
	return ordered, nil
}

func validateMigrationSequence(ordered []migration) error {
	for i := range ordered {
		want := i + 1
		if ordered[i].version != want {
			return fmt.Errorf(
				"migration versions must be contiguous from 1: expected version %d, found %d in %q",
				want, ordered[i].version, ordered[i].fileName,
			)
		}
	}
	return nil
}

func (r *Runner) logResult(applied, previousVersion, latestVersion int) {
	if applied == 0 {
		log.Infof("Database already at version %d, no migrations needed", previousVersion)
	} else {
		log.Infof("Database migrations completed, applied %d migration(s), version: %d", applied, latestVersion)
	}
}

// RunAll applies all pending up migrations.
func (r *Runner) RunAll(ctx context.Context) error {
	loaded, err := r.loadMigrations()
	if err != nil {
		return err
	}
	if len(loaded) == 0 {
		log.Info("No migration scripts found")
		return nil
	}

	current, err := r.getCurrentVersion(ctx)
	if err != nil {
		return err
	}

	applied := 0
	for _, m := range loaded {
		if m.version <= current {
			continue
		}
		if err = r.applyUp(ctx, m); err != nil {
			return err
		}
		applied++
	}

	r.logResult(applied, current, loaded[len(loaded)-1].version)
	return nil
}

func (r *Runner) getCurrentVersion(ctx context.Context) (int, error) {
	var record migrationRecord
	filter := bson.D{{Key: "_id", Value: migrationStateDocumentID}}
	err := r.db.Collection(migrationsCollection).FindOne(ctx, filter).Decode(&record)
	if errors.Is(err, mongo.ErrNoDocuments) {
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

func (r *Runner) applyUp(ctx context.Context, m migration) error {
	log.Infof("Applying migration %06d_%s (up)", m.version, m.name)
	if err := r.runCommands(ctx, m); err != nil {
		return r.markDirty(ctx, m.version, err)
	}
	return r.setVersion(ctx, m.version)
}

func (r *Runner) runCommands(ctx context.Context, m migration) error {
	for i, cmd := range m.up {
		if err := r.db.RunCommand(ctx, cmd).Err(); err != nil {
			return fmt.Errorf("migration %06d_%s command %d failed: %w", m.version, m.name, i+1, err)
		}
	}
	return nil
}

func (r *Runner) setVersion(ctx context.Context, version int) error {
	filter := bson.D{{Key: "_id", Value: migrationStateDocumentID}}
	setDoc := bson.D{
		{Key: "version", Value: version},
		{Key: "dirty", Value: false},
	}
	result, err := r.db.Collection(migrationsCollection).UpdateOne(ctx,
		filter,
		bson.D{{Key: "$set", Value: setDoc}},
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

func (r *Runner) markDirty(ctx context.Context, version int, migrationErr error) error {
	filter := bson.D{{Key: "_id", Value: migrationStateDocumentID}}
	setDoc := bson.D{
		{Key: "version", Value: version},
		{Key: "dirty", Value: true},
	}
	result, err := r.db.Collection(migrationsCollection).UpdateOne(ctx,
		filter,
		bson.D{{Key: "$set", Value: setDoc}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return errors.Join(
			migrationErr,
			fmt.Errorf("failed to mark migration %d as dirty: %w", version, err),
		)
	}
	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		return errors.Join(
			migrationErr,
			fmt.Errorf("failed to mark migration %d as dirty: no documents modified or upserted", version),
		)
	}
	return migrationErr
}
