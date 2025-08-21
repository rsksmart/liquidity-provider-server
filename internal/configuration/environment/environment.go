package environment

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/go-playground/validator/v10"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	log "github.com/sirupsen/logrus"
)

type Environment struct {
	LpsStage         string `env:"LPS_STAGE" validate:"required,oneof=regtest testnet mainnet"`
	Port             uint   `env:"SERVER_PORT" validate:"required"`
	LogLevel         string `env:"LOG_LEVEL" validate:"required"`
	LogFile          string `env:"LOG_FILE"`
	AwsLocalEndpoint string `env:"AWS_LOCAL_ENDPOINT"`
	SecretSource     string `env:"SECRET_SRC" validate:"required,oneof=aws env"`
	WalletManagement string `env:"WALLET" validate:"required,oneof=native fireblocks"`
	Management       ManagementEnv
	Mongo            MongoEnv
	Rsk              RskEnv
	Btc              BtcEnv
	Provider         ProviderEnv
	Pegout           PegoutEnv
	Captcha          CaptchaEnv
	Timeouts         TimeoutEnv
}

type MongoEnv struct {
	Username string `env:"MONGODB_USER" validate:"required"`
	Password string `env:"MONGODB_PASSWORD" validate:"required"`
	Host     string `env:"MONGODB_HOST" validate:"required"`
	Port     uint   `env:"MONGODB_PORT" validate:"required"`
}

type RskEnv struct {
	Endpoint                    string   `env:"RSK_ENDPOINT" validate:"required"`
	ChainId                     uint64   `env:"CHAIN_ID" validate:"required"`
	LbcAddress                  string   `env:"LBC_ADDR" validate:"required"`
	BridgeAddress               string   `env:"RSK_BRIDGE_ADDR" validate:"required"`
	BridgeRequiredConfirmations uint64   `env:"RSK_REQUIRED_BRIDGE_CONFIRMATIONS" validate:"required"`
	ErpKeys                     []string `env:"ERP_KEYS" validate:"required"`
	AccountNumber               int      `env:"ACCOUNT_NUM"` // no validation because 0 works fine
	FeeCollectorAddress         string   `env:"DAO_FEE_COLLECTOR_ADDRESS" validate:"required"`
	// Only if secret source is aws & wallet is native
	EncryptedJsonSecret         string `env:"KEY_SECRET"`
	EncryptedJsonPasswordSecret string `env:"PASSWORD_SECRET"`
	// Only if secret source is env & wallet is native
	KeystoreFile     string `env:"KEYSTORE_FILE"`
	KeystorePassword string `env:"KEYSTORE_PWD"`
}

type BtcEnv struct {
	Network  string `env:"BTC_NETWORK" validate:"required"`
	Username string `env:"BTC_USERNAME" validate:"required"`
	Password string `env:"BTC_PASSWORD" validate:"required"`
	Endpoint string `env:"BTC_ENDPOINT" validate:"required"`
}

type TimeoutEnv struct {
	Bootstrap           uint64 `env:"BOOTSTRAP_TIMEOUT"`
	WatcherPreparation  uint64 `env:"WATCHER_PREPARATION_TIMEOUT"`
	WatcherValidation   uint64 `env:"WATCHER_VALIDATION_TIMEOUT"`
	DatabaseInteraction uint64 `env:"DATABASE_INTERACTION_TIMEOUT"`
	MiningWait          uint64 `env:"MINING_WAIT_TIMEOUT"`
	DatabaseConnection  uint64 `env:"DATABASE_CONNECTION_TIMEOUT"`
	ServerReadHeader    uint64 `env:"SERVER_READ_HEADER_TIMEOUT"`
	ServerWrite         uint64 `env:"SERVER_WRITE_TIMEOUT"`
	ServerIdle          uint64 `env:"SERVER_IDLE_TIMEOUT"`
	PegoutDepositCheck  uint64 `env:"PEGOUT_DEPOSIT_CHECK_TIMEOUT"`
}

func (env BtcEnv) GetNetworkParams() (*chaincfg.Params, error) {
	switch env.Network {
	case "mainnet":
		return &chaincfg.MainNetParams, nil
	case "testnet":
		return &chaincfg.TestNet3Params, nil
	case "regtest":
		return &chaincfg.RegressionNetParams, nil
	default:
		return nil, fmt.Errorf("invalid network name: %v", env.Network)
	}
}

type ProviderEnv struct {
	AlertSenderEmail    string                          `env:"ALERT_SENDER_EMAIL"  validate:"required"`
	AlertRecipientEmail string                          `env:"ALERT_RECIPIENT_EMAIL"  validate:"required"`
	Name                string                          `env:"PROVIDER_NAME"  validate:"required"`
	ApiBaseUrl          string                          `env:"BASE_URL"  validate:"required"`
	ProviderType        liquidity_provider.ProviderType `env:"PROVIDER_TYPE"  validate:"required,oneof=pegin pegout both"`
}

// PeginEnv This structure was kept just in case, right now all the parameters are manipulated through management API
type PeginEnv struct{}

type PegoutEnv struct {
	DepositCacheStartBlock uint64 `env:"PEGOUT_DEPOSIT_CACHE_START_BLOCK"`
}

type CaptchaEnv struct {
	SecretKey string  `env:"CAPTCHA_SECRET_KEY"`
	SiteKey   string  `env:"CAPTCHA_SITE_KEY"`
	Threshold float32 `env:"CAPTCHA_THRESHOLD"`
	Disabled  bool    `env:"DISABLE_CAPTCHA"`
	Url       string  `env:"CAPTCHA_URL"`
}

type ManagementEnv struct {
	EnableManagementApi   bool   `env:"ENABLE_MANAGEMENT_API"`
	SessionAuthKey        string `env:"MANAGEMENT_AUTH_KEY"`
	SessionEncryptionKey  string `env:"MANAGEMENT_ENCRYPTION_KEY"`
	SessionTokenAuthKey   string `env:"MANAGEMENT_TOKEN_AUTH_KEY"`
	UseHttps              bool   `env:"MANAGEMENT_USE_HTTPS"`
	EnableSecurityHeaders bool   `env:"ENABLE_SECURITY_HEADERS"`
}

func LoadEnv() *Environment {
	validate := validator.New(validator.WithRequiredStructEnabled())
	env := &Environment{}
	if err := Load(env); err != nil {
		log.Fatal("Error reading environment: ", err)
	} else if err = validate.Struct(env); err != nil {
		log.Fatal("Environment incomplete: ", err)
	}

	return env
}
