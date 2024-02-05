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

func parseEnv(arg reflect.Value) error {
	if arg.Kind() == reflect.Ptr && arg.IsNil() {
		arg.Set(reflect.New(arg.Elem().Type()))
	}

	if arg.Kind() == reflect.Ptr {
		arg = arg.Elem()
	}

	argType := arg.Type()
	var envVarName string
	var fieldType reflect.Kind

	for i := 0; i < arg.NumField(); i++ {
		envVarName = argType.Field(i).Tag.Get("env")
		fieldType = argType.Field(i).Type.Kind()
		if envVarName == "" && fieldType != reflect.Struct {
			return fmt.Errorf("field %s doesn't have envVarName tag", argType.Field(i).Name)
		}

		if argType.Field(i).Type.Kind() == reflect.Struct {
			if err := parseEnv(arg.Field(i)); err != nil {
				return err
			}
		} else if envVarName != "" {
			if err := setEnvValue(arg.Field(i), os.Getenv(envVarName)); err != nil {
				return fmt.Errorf("error reading %s environment variable: %w", envVarName, err)
			}
		}
	}
	return nil
}

func setEnvValue(field reflect.Value, envVar string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(envVar)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		if envVar == "" {
			field.SetUint(0)
		} else if uint64Value, err := strconv.ParseUint(envVar, 10, 64); err != nil {
			return err
		} else {
			field.SetUint(uint64Value)
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		if envVar == "" {
			field.SetInt(0)
		} else if int64Value, err := strconv.ParseInt(envVar, 10, 64); err != nil {
			return err
		} else {
			field.SetInt(int64Value)
		}
	case reflect.Float32, reflect.Float64:
		if envVar == "" {
			field.SetFloat(0)
		} else if float64Value, err := strconv.ParseFloat(envVar, 64); err != nil {
			return err
		} else {
			field.SetFloat(float64Value)
		}
	case reflect.Bool:
		if boolValue, err := strconv.ParseBool(envVar); err != nil {
			return err
		} else {
			field.SetBool(boolValue)
		}
	case reflect.Slice:
		if field.Type().Elem().Kind() != reflect.String {
			return errors.New("unsupported env array")
		}
		for _, value := range strings.Split(envVar, ",") {
			element := reflect.New(field.Type().Elem()).Elem()
			element.SetString(value)
			field.Set(reflect.Append(field, element))
		}
	case reflect.Map:
		mapValue := reflect.New(field.Type())
		mapContent := mapValue.Interface()
		if err := json.Unmarshal([]byte(envVar), &mapContent); err != nil {
			return err
		} else {
			field.Set(mapValue.Elem())
		}
	default:
		return jsonUnmarshalEnvValue(field, envVar)
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
