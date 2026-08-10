package environment_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"strings"
	"testing"
)

func setUpEnv(t *testing.T) {
	// this map is to define the value for the vars that intentionally have a zero value in the sample-config.env file
	var sampleZeroVars = map[string]string{
		"ENABLE_SECURITY_HEADERS":              "true",
		"MANAGEMENT_USE_HTTPS":                 "true",
		"ENABLE_MANAGEMENT_API":                "true",
		"LBC_ADDR":                             "0x1234",
		"ACCOUNT_NUM":                          "1",
		"CAPTCHA_SECRET_KEY":                   "secret",
		"CAPTCHA_SITE_KEY":                     "site",
		"PEGOUT_DEPOSIT_CACHE_START_BLOCK":     "1",
		"RSK_EXTRA_SOURCES":                    "test1,test2",
		"BTC_EXTRA_SOURCES":                    `[{"format": "rpc", "url": "test3.com"}, {"format": "mempool", "url": "test4.com"}]`,
		"ECLIPSE_RSK_TOLERANCE_THRESHOLD":      "5",
		"ECLIPSE_RSK_MAX_MS_WAIT_FOR_BLOCK":    "1000",
		"ECLIPSE_RSK_WAIT_POLLING_MS_INTERVAL": "500",
		"ECLIPSE_BTC_TOLERANCE_THRESHOLD":      "5",
		"ECLIPSE_BTC_MAX_MS_WAIT_FOR_BLOCK":    "1000",
		"ECLIPSE_BTC_WAIT_POLLING_MS_INTERVAL": "500",
		"ECLIPSE_ALERT_COOLDOWN_SECONDS":       "60",
		"ECLIPSE_CHECK_ENABLED":                "true",
		"BTC_RELEASE_WATCHER_START_BLOCK":      "1",
		"USE_SEGWIT_FEDERATION":                "true",
		"ALLOWED_ORIGINS":                      "http://example.com,http://example2.com",
		"REBALANCE_STRATEGY":                   "ALL_AT_ONCE",
		"RUN_DB_MIGRATIONS":                    "true",
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
}

func TestLoad(t *testing.T) {
	t.Run("env vars are loaded correctly", func(t *testing.T) {
		setUpEnv(t)
		env := &environment.Environment{}
		err := environment.Load(env)
		require.NoError(t, err)
		assertNonZeroFieldsRecursive(t, env)
	})
	t.Run("parses empty string as false for bool values", func(t *testing.T) {
		setUpEnv(t)
		t.Setenv("ENABLE_SECURITY_HEADERS", "")
		env := &environment.Environment{}
		err := environment.Load(env)
		require.NoError(t, err)
		assert.False(t, env.Management.EnableSecurityHeaders)
	})
	t.Run("does not fail when a slice is an empty string", func(t *testing.T) {
		setUpEnv(t)
		t.Setenv("RSK_EXTRA_SOURCES", "")
		env := &environment.Environment{}
		err := environment.Load(env)
		require.NoError(t, err)
	})
	t.Run("preserves an explicitly configured zero start block", func(t *testing.T) {
		t.Setenv("PEGIN_ADDRESS_REGISTRY_WATCHER_START_BLOCK", "0")
		t.Setenv("PEGIN_ADDRESS_REGISTRY_WATCHER_PAGE_SIZE", "10")
		env := &environment.Environment{}
		require.NoError(t, environment.Load(env))

		config, err := env.Pegin.AddressRegistryWatcherConfig()
		require.NoError(t, err)
		assert.True(t, config.Enabled)
		assert.Zero(t, config.StartBlock)
		assert.Equal(t, uint64(10), config.PageSize)
	})
}

func TestPeginEnv_AddressRegistryWatcherConfig(t *testing.T) {
	present := func(value uint64) *environment.OptionalUint64 {
		return &environment.OptionalUint64{Value: value, Present: true}
	}

	tests := []struct {
		name       string
		env        environment.PeginEnv
		enabled    bool
		missingVar string
	}{
		{name: "both present", env: environment.PeginEnv{
			AddressRegistryWatcherStartBlock: present(0),
			AddressRegistryWatcherPageSize:   present(10),
		}, enabled: true},
		{name: "start only", env: environment.PeginEnv{
			AddressRegistryWatcherStartBlock: present(1),
		}, missingVar: "PEGIN_ADDRESS_REGISTRY_WATCHER_PAGE_SIZE"},
		{name: "page only", env: environment.PeginEnv{
			AddressRegistryWatcherPageSize: present(10),
		}, missingVar: "PEGIN_ADDRESS_REGISTRY_WATCHER_START_BLOCK"},
		{name: "neither present", env: environment.PeginEnv{}},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.missingVar != "" {
				defer test.AssertLogContains(t, testCase.missingVar)()
			}
			config, err := testCase.env.AddressRegistryWatcherConfig()
			require.NoError(t, err)
			assert.Equal(t, testCase.enabled, config.Enabled)
		})
	}

	t.Run("rejects zero page size when enabled", func(t *testing.T) {
		env := environment.PeginEnv{
			AddressRegistryWatcherStartBlock: present(1),
			AddressRegistryWatcherPageSize:   present(0),
		}
		_, err := env.AddressRegistryWatcherConfig()
		require.ErrorContains(t, err, "must be greater than zero")
	})
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
