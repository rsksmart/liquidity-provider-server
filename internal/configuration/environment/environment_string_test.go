package environment_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/stretchr/testify/require"
)

const expectedSecretMask = "********"

//nolint:funlen 
func TestEnvironment_String(t *testing.T) {
	t.Run("redacts only secret fields", func(t *testing.T) {
		env := environment.Environment{
			LpsStage:         "regtest",
			Port:             8080,
			LogLevel:         "debug",
			AwsLocalEndpoint: "http://localstack:4566",
			SecretSource:     "aws",
			WalletManagement: "native",
			AllowedOrigins:   []string{"http://localhost:3000"},
			Management: environment.ManagementEnv{
				EnableManagementApi:   true,
				SessionAuthKey:        "auth-secret",
				SessionEncryptionKey:  "encryption-secret",
				SessionTokenAuthKey:   "token-secret",
				UseHttps:              true,
				EnableSecurityHeaders: true,
			},
			Mongo: environment.MongoEnv{
				Username:      "mongo-user",
				Password:      "mongo-secret",
				Host:          "mongodb.local",
				Port:          27017,
				RunMigrations: true,
			},
			Rsk: environment.RskEnv{
				Endpoint:                    "http://rsk.local",
				ChainId:                     31,
				PeginContractAddress:        "0x0000000000000000000000000000000000000001",
				PegoutContractAddress:       "0x0000000000000000000000000000000000000002",
				CollateralManagementAddress: "0x0000000000000000000000000000000000000003",
				DiscoveryAddress:            "0x0000000000000000000000000000000000000004",
				BridgeAddress:               "0x0000000000000000000000000000000000000005",
				BridgeRequiredConfirmations: 10,
				ErpKeys:                     []string{"erp-key"},
				EncryptedJsonSecret:         "key-secret-id",
				EncryptedJsonPasswordSecret: "password-secret-id",
				KeystoreFile:                "keystore.json",
				KeystorePassword:            "keystore-secret",
			},
			Btc: environment.BtcEnv{
				Network:  "regtest",
				Username: "btc-user",
				Password: "btc-secret",
				Endpoint: "http://btc.local",
			},
			Provider: environment.ProviderEnv{
				AlertSenderEmail:    "sender@example.com",
				AlertRecipientEmail: "recipient@example.com",
				Name:                "provider-name",
				ApiBaseUrl:          "http://lps.local",
				ProviderTypeName:    "both",
			},
			Captcha: environment.CaptchaEnv{
				SecretKey: "captcha-secret",
				SiteKey:   "captcha-site",
				Threshold: 0.5,
				Disabled:  false,
				Url:       "http://captcha.local",
			},
		}

		output := fmt.Sprintf("%+v", env)

		require.NotContains(t, output, "mongo-secret")
		require.NotContains(t, output, "auth-secret")
		require.NotContains(t, output, "encryption-secret")
		require.NotContains(t, output, "token-secret")
		require.NotContains(t, output, "keystore-secret")
		require.NotContains(t, output, "btc-secret")
		require.NotContains(t, output, "captcha-secret")
		require.Contains(t, output, "regtest")
		require.Contains(t, output, "mongodb.local")
		require.Contains(t, output, "btc-user")
		require.Contains(t, output, "captcha-site")
		require.Contains(t, output, "key-secret-id")
		require.Contains(t, output, "password-secret-id")
		require.Equal(t, 7, strings.Count(output, expectedSecretMask))
	})

	t.Run("keeps empty secret fields empty", func(t *testing.T) {
		env := environment.Environment{
			Management: environment.ManagementEnv{},
			Mongo:      environment.MongoEnv{},
			Rsk:        environment.RskEnv{},
			Btc:        environment.BtcEnv{},
			Captcha:    environment.CaptchaEnv{},
		}

		require.NotContains(t, env.String(), expectedSecretMask)
	})

	t.Run("redacts pointer formatting", func(t *testing.T) {
		env := environment.Environment{
			Mongo: environment.MongoEnv{Password: "mongo-secret"},
		}

		output := fmt.Sprintf("%+v", &env)

		require.NotContains(t, output, "mongo-secret")
		require.Contains(t, output, expectedSecretMask)
	})

	t.Run("redacts every field tagged as secret", func(t *testing.T) {
		env := environment.Environment{}
		secretValues := setSecretTaggedStrings(reflect.ValueOf(&env).Elem(), "Environment")

		output := env.String()

		require.NotEmpty(t, secretValues)
		for _, secretValue := range secretValues {
			require.NotContains(t, output, secretValue)
		}
		require.Equal(t, len(secretValues), strings.Count(output, expectedSecretMask))
	})
}

func setSecretTaggedStrings(value reflect.Value, path string) []string {
	valueType := value.Type()
	secretValues := make([]string, 0)
	for i := 0; i < value.NumField(); i++ {
		field := valueType.Field(i)
		fieldValue := value.Field(i)
		fieldPath := path + "." + field.Name
		if _, isSecret := field.Tag.Lookup("secret"); isSecret && fieldValue.Kind() == reflect.String {
			secretValue := "secret-value-for-" + fieldPath
			fieldValue.SetString(secretValue)
			secretValues = append(secretValues, secretValue)
		} else if fieldValue.Kind() == reflect.Struct {
			secretValues = append(secretValues, setSecretTaggedStrings(fieldValue, fieldPath)...)
		}
	}
	return secretValues
}
