package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func EmptyFilter() bson.M {
	return bson.M{}
}

func ExtJSONToDocument(raw json.RawMessage) (bson.D, error) {
	var doc bson.D
	if err := bson.UnmarshalExtJSON(raw, true, &doc); err != nil {
		return bson.D{}, err
	}
	return doc, nil
}

func IndexKeysContainField(keys any, field string) bool {
	switch typed := keys.(type) {
	case map[string]any:
		return mapHasField(typed, field)
	case bson.M:
		return mapHasField(map[string]any(typed), field)
	case bson.D:
		return bsonDHasField(typed, field)
	default:
		return sliceOfKeyedElementsHasField(keys, field)
	}
}

// IndexKeysIsSingleFieldAscending reports whether keys is exactly {field: 1} —
// catches accidental compound indexes or descending shapes that IndexKeysContainField
// would silently accept.
func IndexKeysIsSingleFieldAscending(keys any, field string) bool {
	switch typed := keys.(type) {
	case map[string]any:
		return mapIsSingleFieldAscending(typed, field)
	case bson.M:
		return mapIsSingleFieldAscending(map[string]any(typed), field)
	case bson.D:
		return len(typed) == 1 && typed[0].Key == field && isAscending(typed[0].Value)
	default:
		return sliceOfKeyedElementsIsSingleFieldAscending(keys, field)
	}
}

func mapIsSingleFieldAscending(m map[string]any, field string) bool {
	if len(m) != 1 {
		return false
	}
	v, ok := m[field]
	return ok && isAscending(v)
}

func sliceOfKeyedElementsIsSingleFieldAscending(keys any, field string) bool {
	val := reflect.ValueOf(keys)
	if val.Kind() != reflect.Slice || val.Len() != 1 {
		return false
	}
	item := val.Index(0)
	if item.Kind() == reflect.Pointer {
		item = item.Elem()
	}
	if item.Kind() != reflect.Struct {
		return false
	}
	keyField := item.FieldByName("Key")
	valueField := item.FieldByName("Value")
	if !keyField.IsValid() || keyField.Kind() != reflect.String || keyField.String() != field {
		return false
	}
	if !valueField.IsValid() {
		return false
	}
	return isAscending(valueField.Interface())
}

func isAscending(v any) bool {
	switch n := v.(type) {
	case int:
		return n == 1
	case int32:
		return n == 1
	case int64:
		return n == 1
	case float64:
		return n == 1
	default:
		return false
	}
}

func mapHasField(m map[string]any, field string) bool {
	_, ok := m[field]
	return ok
}

func bsonDHasField(d bson.D, field string) bool {
	for _, e := range d {
		if e.Key == field {
			return true
		}
	}
	return false
}

func sliceOfKeyedElementsHasField(keys any, field string) bool {
	// Used for driver-specific bson.D-like containers without hard dependency.
	val := reflect.ValueOf(keys)
	if val.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		if item.Kind() == reflect.Pointer {
			item = item.Elem()
		}
		if item.Kind() != reflect.Struct {
			continue
		}
		keyField := item.FieldByName("Key")
		if keyField.IsValid() && keyField.Kind() == reflect.String && keyField.String() == field {
			return true
		}
	}
	return false
}

func EnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func EnvOrUint(key string, fallback uint) (uint, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseUint(v, 10, 0)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
}

func FixturesPath() string {
	if v := os.Getenv("FIXTURES_DIR"); v != "" {
		return v
	}
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return filepath.Join("test", "mongodb", "fixtures")
	}
	return filepath.Join(filepath.Dir(filename), "..", "fixtures")
}
