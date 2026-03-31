//go:build integration

package mongodb_test

import (
	"encoding/json"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

func emptyFilter() bson.M {
	return bson.M{}
}

func extJSONToDocument(raw json.RawMessage) (bson.D, error) {
	var doc bson.D
	if err := bson.UnmarshalExtJSON(raw, true, &doc); err != nil {
		return bson.D{}, err
	}
	return doc, nil
}

func indexKeysContainField(keys any, field string) bool {
	switch typed := keys.(type) {
	case map[string]any:
		_, ok := typed[field]
		return ok
	case bson.M:
		_, ok := typed[field]
		return ok
	case bson.D:
		for _, e := range typed {
			if e.Key == field {
				return true
			}
		}
		return false
	}

	// Fallback for driver-specific bson.D-like containers without hard dependency.
	val := reflect.ValueOf(keys)
	if val.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		if item.Kind() == reflect.Struct {
			keyField := item.FieldByName("Key")
			if keyField.IsValid() && keyField.Kind() == reflect.String && keyField.String() == field {
				return true
			}
		}
	}
	return false
}
