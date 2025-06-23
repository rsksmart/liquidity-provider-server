package environment

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func Load(arg *Environment) error {
	return parseEnv(reflect.ValueOf(arg))
}

func parseEnv(value reflect.Value) error {
	if value.Kind() == reflect.Ptr && value.IsNil() {
		value.Set(reflect.New(value.Elem().Type()))
	}

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	valueType := value.Type()
	var err error

	for i := 0; i < value.NumField(); i++ {
		if err = parseField(i, value, valueType); err != nil {
			return err
		}
	}
	return nil
}

func parseField(fieldNumber int, value reflect.Value, valueType reflect.Type) error {
	envVarName := valueType.Field(fieldNumber).Tag.Get("env")
	fieldType := valueType.Field(fieldNumber).Type.Kind()
	if envVarName == "" && fieldType != reflect.Struct {
		return fmt.Errorf("field %s doesn't have envVarName tag", valueType.Field(fieldNumber).Name)
	}

	if valueType.Field(fieldNumber).Type.Kind() == reflect.Struct {
		if err := parseEnv(value.Field(fieldNumber)); err != nil {
			return err
		}
	} else if envVarName != "" {
		if err := setEnvValue(value.Field(fieldNumber), os.Getenv(envVarName)); err != nil {
			return fmt.Errorf("error reading %s environment variable: %w", envVarName, err)
		}
	}
	return nil
}

func setEnvValue(field reflect.Value, envVar string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(envVar)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return parseUint(envVar, field)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return parseInt(envVar, field)
	case reflect.Float32, reflect.Float64:
		return parseFloat(envVar, field)
	case reflect.Bool:
		return parseBool(envVar, field)
	case reflect.Slice:
		return parseSlice(envVar, field)
	case reflect.Map:
		return parseMap(envVar, field)
	default:
		return jsonUnmarshalEnvValue(field, envVar)
	}
	return nil
}

func parseUint(envVar string, field reflect.Value) error {
	if envVar == "" {
		field.SetUint(0)
	} else if uint64Value, err := strconv.ParseUint(envVar, 10, 64); err != nil {
		return err
	} else {
		field.SetUint(uint64Value)
	}
	return nil
}

func parseInt(envVar string, field reflect.Value) error {
	if envVar == "" {
		field.SetInt(0)
	} else if int64Value, err := strconv.ParseInt(envVar, 10, 64); err != nil {
		return err
	} else {
		field.SetInt(int64Value)
	}
	return nil
}

func parseFloat(envVar string, field reflect.Value) error {
	if envVar == "" {
		field.SetFloat(0)
	} else if float64Value, err := strconv.ParseFloat(envVar, 64); err != nil {
		return err
	} else {
		field.SetFloat(float64Value)
	}
	return nil
}

func parseBool(envVar string, field reflect.Value) error {
	if envVar == "" {
		field.SetBool(false)
		return nil
	}
	if boolValue, err := strconv.ParseBool(envVar); err != nil {
		return err
	} else {
		field.SetBool(boolValue)
	}
	return nil
}

func parseSlice(envVar string, field reflect.Value) error {
	if field.Type().Elem().Kind() != reflect.String {
		return errors.New("unsupported env array")
	}
	if envVar == "" {
		return nil
	}
	for _, value := range strings.Split(envVar, ",") {
		element := reflect.New(field.Type().Elem()).Elem()
		element.SetString(value)
		field.Set(reflect.Append(field, element))
	}
	return nil
}

func parseMap(envVar string, field reflect.Value) error {
	mapValue := reflect.New(field.Type())
	mapContent := mapValue.Interface()
	if err := json.Unmarshal([]byte(envVar), &mapContent); err != nil {
		return err
	} else {
		field.Set(mapValue.Elem())
	}
	return nil
}

func jsonUnmarshalEnvValue(field reflect.Value, envVar string) error {
	var fieldType reflect.Type
	if field.Kind() == reflect.Ptr {
		fieldType = field.Type().Elem()
	} else {
		fieldType = field.Type()
	}

	if field.Kind() == reflect.Ptr && field.IsNil() {
		field.Set(reflect.New(fieldType))
	}

	unmarshaler, ok := field.Interface().(json.Unmarshaler)
	if !ok {
		return fmt.Errorf("cant unmarshal field %s", field.Type().Name())
	}
	return unmarshaler.UnmarshalJSON([]byte(envVar))
}
