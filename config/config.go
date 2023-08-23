package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/connectors"
	"github.com/rsksmart/liquidity-provider-server/http"
	"github.com/rsksmart/liquidity-provider-server/pegin"
	"github.com/rsksmart/liquidity-provider-server/pegout"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	LogFile              string   `env:"LOG_FILE"`
	Debug                bool     `env:"DEBUG"`
	IrisActivationHeight int      `env:"IRIS_ACTIVATION_HEIGHT"`
	ErpKeys              []string `env:"ERP_KEYS"`
	PeginProviderName    string   `env:"PEGIN_PROVIDER_NAME"`
	PegoutProviderName   string   `env:"PEGOUT_PROVIDER_NAME"`
	BaseURL              string   `env:"BASE_URL"`
	QuoteCacheStartBlock uint64   `env:"QUOTE_CACHE_START_BLOCK"`

	Server struct {
		Port uint `env:"SERVER_PORT"`
	}
	DB struct {
		Regtest struct {
			Host     string `env:"DB_REGTEST_HOST"`
			Database string `env:"DB_REGTEST_DATABASE"`
			User     string `env:"DB_REGTEST_USER"`
			Password string `env:"DB_REGTEST_PASSWORD"`
			Port     uint   `env:"DB_REGTEST_PORT"`
		}
		Path string `env:"DB_PATH"`
	}
	RSK           http.LiquidityProviderList
	BTC           connectors.BtcConfig
	Provider      pegin.ProviderConfig  `env:",prefix=PEGIN_PROVIDER_"`
	PegoutProvier pegout.ProviderConfig `env:",prefix=PEGOUT_PROVIDER_"`
}

func LoadEnv(arg any) error {
	return loadEnvWithPrefix(reflect.ValueOf(arg), "")
}

func loadEnvWithPrefix(arg reflect.Value, prefix string) error {
	if arg.Kind() == reflect.Ptr && arg.IsNil() {
		arg.Set(reflect.New(arg.Elem().Type()))
	}

	if arg.Kind() == reflect.Ptr {
		arg = arg.Elem()
	}

	argType := arg.Type()
	var env envTag
	var fieldType reflect.Kind

	for i := 0; i < arg.NumField(); i++ {
		env = parseEnvTag(argType.Field(i).Tag)
		fieldType = argType.Field(i).Type.Kind()
		if env.value == "" && fieldType != reflect.Struct {
			return fmt.Errorf("Field %s doesn't have env tag", argType.Field(i).Name)
		}

		if argType.Field(i).Type.Kind() == reflect.Struct {
			if env.squash {
				env.prefix = prefix
			}
			if err := loadEnvWithPrefix(arg.Field(i), env.prefix); err != nil {
				return err
			}
		} else if env.value != "" {
			if err := setEnvValue(arg.Field(i), os.Getenv(prefix+env.value)); err != nil {
				return err
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

type envTag struct {
	value  string
	prefix string
	squash bool
}

func parseEnvTag(tag reflect.StructTag) envTag {
	envTagParts := strings.Split(tag.Get("env"), ",")
	var result envTag
	result.value = envTagParts[0]
	for _, element := range envTagParts {
		if element == "squash" {
			result.squash = true
		} else if strings.Contains(element, "prefix=") {
			result.prefix = strings.TrimPrefix(element, "prefix=")
		}
	}
	return result
}
