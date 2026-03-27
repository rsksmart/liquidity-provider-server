package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func exportFixturesAsExtJSON(ctx context.Context, client *mongodriver.Client) (string, error) {
	fixturesDir := fixturesPath()
	if err := ensureFixturesDir(fixturesDir); err != nil {
		return "", err
	}

	rawDB := client.Database(mongo.DbName)
	for _, c := range fixtureCollections {
		if err := exportOneCollectionAsExtJSON(ctx, rawDB, fixturesDir, c.collection, c.fileName); err != nil {
			return "", err
		}
	}

	return fixturesDir, nil
}

func ensureFixturesDir(fixturesDir string) error {
	if err := os.MkdirAll(fixturesDir, 0o755); err != nil {
		return fmt.Errorf("create fixtures dir %s: %w", fixturesDir, err)
	}
	return nil
}

func exportOneCollectionAsExtJSON(ctx context.Context, rawDB *mongodriver.Database, fixturesDir, collName, fileName string) error {
	cursor, err := rawDB.Collection(collName).Find(
		ctx,
		bson.M{},
		options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}),
	)
	if err != nil {
		return fmt.Errorf("read %s: %w", collName, err)
	}
	defer cursor.Close(ctx)

	var docs []bson.M
	if err = cursor.All(ctx, &docs); err != nil {
		return fmt.Errorf("decode %s: %w", collName, err)
	}

	for i := range docs {
		delete(docs[i], "_id")
	}

	extDocs := make([]json.RawMessage, 0, len(docs))
	for _, doc := range docs {
		ext, marshalErr := bson.MarshalExtJSON(doc, true, true)
		if marshalErr != nil {
			return fmt.Errorf("marshal extjson for %s: %w", collName, marshalErr)
		}
		extDocs = append(extDocs, ext)
	}

	data, err := json.MarshalIndent(extDocs, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", collName, err)
	}

	filePath := filepath.Join(fixturesDir, fileName)
	if err := os.WriteFile(filePath, data, 0o600); err != nil {
		return fmt.Errorf("write %s: %w", filePath, err)
	}

	return nil
}

func fixturesPath() string {
	if v := os.Getenv("FIXTURES_DIR"); v != "" {
		return v
	}
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "fixtures"
	}
	return filepath.Join(filepath.Dir(filename), "..", "..", "fixtures")
}
