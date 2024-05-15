package environment_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	// this map is to define the value for the vars that intentionally have a zero value in the sample-config.env file
	var sampleZeroVars = map[string]string{
		"MANAGEMENT_USE_HTTPS":             "true",
		"LBC_ADDR":                         "0x1234",
		"ACCOUNT_NUM":                      "1",
		"CAPTCHA_SECRET_KEY":               "secret",
		"CAPTCHA_SITE_KEY":                 "site",
		"PEGOUT_DEPOSIT_CACHE_START_BLOCK": "1",
	}
	const envFilePath = "../../../sample-config.env"
	envFile, err := os.ReadFile(envFilePath)
	envLines := strings.Split(string(envFile), "\n")
	for _, line := range envLines {
		if line != "" && !strings.HasPrefix(line, "#") {
			parts := strings.Split(line, "=")
			if zeroVar, ok := sampleZeroVars[parts[0]]; ok {
				t.Setenv(parts[0], zeroVar)
			} else {
				require.Lenf(t, parts, 2, "Var %s doesn't have a value", parts[0])
				t.Setenv(parts[0], parts[1])
			}
		}
	}
	require.NoError(t, err)
	env := &environment.Environment{}
	err = environment.Load(env)
	require.NoError(t, err)
	assertNonZeroFieldsRecursive(t, env)
}

func assertNonZeroFieldsRecursive(t *testing.T, aStruct any) {
	value := reflect.ValueOf(aStruct)
	if value.IsZero() {
		t.Errorf("The struct is zero")
	} else if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	for i := 0; i < value.NumField(); i++ {
		if value.Field(i).IsZero() {
			t.Errorf("Field %s is unset", value.Type().Field(i).Name)
		} else if value.Field(i).Kind() == reflect.Struct {
			assertNonZeroFieldsRecursive(t, value.Field(i).Interface())
		}
	}
}
