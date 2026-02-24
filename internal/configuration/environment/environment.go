package environment

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/go-playground/validator/v10"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
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
	Eclipse          EclipseEnv
	ColdWallet       ColdWalletEnv
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
	PeginContractAddress        string   `env:"PEGIN_CONTRACT_ADDRESS" validate:"required"`
	PegoutContractAddress       string   `env:"PEGOUT_CONTRACT_ADDRESS" validate:"required"`
	CollateralManagementAddress string   `env:"COLLATERAL_MANAGEMENT_ADDRESS" validate:"required"`
	DiscoveryAddress            string   `env:"DISCOVERY_ADDRESS" validate:"required"`
	BridgeAddress               string   `env:"RSK_BRIDGE_ADDR" validate:"required"`
	BridgeRequiredConfirmations uint64   `env:"RSK_REQUIRED_BRIDGE_CONFIRMATIONS" validate:"required"`
	ErpKeys                     []string `env:"ERP_KEYS" validate:"required"`
	UseSegwitFederation         bool     `env:"USE_SEGWIT_FEDERATION"`
	AccountNumber               int      `env:"ACCOUNT_NUM"` // no validation because 0 works fine
	FeeCollectorAddress         string   `env:"DAO_FEE_COLLECTOR_ADDRESS" validate:"required"`
	// Only if secret source is aws & wallet is native
	WalletSecret   string `env:"WALLET_SECRET"`
	PasswordSecret string `env:"PASSWORD_SECRET"`
	// Only if secret source is env & wallet is native
	WalletFile       string   `env:"WALLET_FILE"`
	KeystorePassword string   `env:"KEYSTORE_PWD"`
	RskExtraSources  []string `env:"RSK_EXTRA_SOURCES"`
}

type BtcExtraSource struct {
	Format string `json:"format" validate:"required,oneof=rpc,mempool"`
	Url    string `json:"url" validate:"required,url"`
}

type BtcEnv struct {
	Network         string           `env:"BTC_NETWORK" validate:"required"`
	Username        string           `env:"BTC_USERNAME" validate:"required"`
	Password        string           `env:"BTC_PASSWORD" validate:"required"`
	Endpoint        string           `env:"BTC_ENDPOINT" validate:"required"`
	BtcExtraSources []BtcExtraSource `env:"BTC_EXTRA_SOURCES"`
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
	BtcReleaseCheck     uint64 `env:"BTC_RELEASE_CHECK_TIMEOUT"`
}

type EclipseEnv struct {
	Enabled                  bool   `env:"ECLIPSE_CHECK_ENABLED"`
	RskToleranceThreshold    uint8  `env:"ECLIPSE_RSK_TOLERANCE_THRESHOLD"`
	RskMaxMsWaitForBlock     uint64 `env:"ECLIPSE_RSK_MAX_MS_WAIT_FOR_BLOCK"`
	RskWaitPollingMsInterval uint64 `env:"ECLIPSE_RSK_WAIT_POLLING_MS_INTERVAL"`
	BtcToleranceThreshold    uint8  `env:"ECLIPSE_BTC_TOLERANCE_THRESHOLD"`
	BtcMaxMsWaitForBlock     uint64 `env:"ECLIPSE_BTC_MAX_MS_WAIT_FOR_BLOCK"`
	BtcWaitPollingMsInterval uint64 `env:"ECLIPSE_BTC_WAIT_POLLING_MS_INTERVAL"`
	AlertCooldownSeconds     uint64 `env:"ECLIPSE_ALERT_COOLDOWN_SECONDS"`
}

func (env *EclipseEnv) FillWithDefaults() *EclipseEnv {
	defaults := EclipseEnv{
		RskToleranceThreshold:    50,
		RskMaxMsWaitForBlock:     10_000,
		RskWaitPollingMsInterval: 1000,
		BtcToleranceThreshold:    50,
		BtcMaxMsWaitForBlock:     60_000,
		BtcWaitPollingMsInterval: 10_000,
		AlertCooldownSeconds:     30 * 60, // 30 min
	}
	env.RskToleranceThreshold = utils.FirstNonZero(env.RskToleranceThreshold, defaults.RskToleranceThreshold)
	env.RskMaxMsWaitForBlock = utils.FirstNonZero(env.RskMaxMsWaitForBlock, defaults.RskMaxMsWaitForBlock)
	env.RskWaitPollingMsInterval = utils.FirstNonZero(env.RskWaitPollingMsInterval, defaults.RskWaitPollingMsInterval)
	env.BtcToleranceThreshold = utils.FirstNonZero(env.BtcToleranceThreshold, defaults.BtcToleranceThreshold)
	env.BtcMaxMsWaitForBlock = utils.FirstNonZero(env.BtcMaxMsWaitForBlock, defaults.BtcMaxMsWaitForBlock)
	env.BtcWaitPollingMsInterval = utils.FirstNonZero(env.BtcWaitPollingMsInterval, defaults.BtcWaitPollingMsInterval)
	env.AlertCooldownSeconds = utils.FirstNonZero(env.AlertCooldownSeconds, defaults.AlertCooldownSeconds)
	return env
}

func (env *EclipseEnv) ToConfig() watcher.EclipseCheckConfig {
	return watcher.EclipseCheckConfig{
		RskToleranceThreshold:    env.RskToleranceThreshold,
		RskMaxMsWaitForBlock:     env.RskMaxMsWaitForBlock,
		RskWaitPollingMsInterval: env.RskWaitPollingMsInterval,
		BtcToleranceThreshold:    env.BtcToleranceThreshold,
		BtcMaxMsWaitForBlock:     env.BtcMaxMsWaitForBlock,
		BtcWaitPollingMsInterval: env.BtcWaitPollingMsInterval,
	}
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
	AlertSenderEmail    string `env:"ALERT_SENDER_EMAIL"  validate:"required"`
	AlertRecipientEmail string `env:"ALERT_RECIPIENT_EMAIL"  validate:"required"`
	Name                string `env:"PROVIDER_NAME"  validate:"required"`
	ApiBaseUrl          string `env:"BASE_URL"  validate:"required"`
	ProviderTypeName    string `env:"PROVIDER_TYPE"  validate:"required,oneof=pegin pegout both"`
}

func (env *ProviderEnv) ProviderType() liquidity_provider.ProviderType {
	switch env.ProviderTypeName {
	case "pegin":
		return liquidity_provider.PeginProvider
	case "pegout":
		return liquidity_provider.PegoutProvider
	case "both":
		return liquidity_provider.FullProvider
	default:
		return -1
	}
}

// PeginEnv This structure was kept just in case, right now all the parameters are manipulated through management API
type PeginEnv struct{}

type PegoutEnv struct {
	DepositCacheStartBlock      uint64 `env:"PEGOUT_DEPOSIT_CACHE_START_BLOCK"`
	BtcReleaseWatcherStartBlock uint64 `env:"BTC_RELEASE_WATCHER_START_BLOCK"`
	BtcReleaseWatcherPageSize   uint64 `env:"BTC_RELEASE_WATCHER_PAGE_SIZE"`
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

type ColdWalletEnv struct {
	BtcMinTransferFeeMultiplier   uint64 `env:"BTC_MIN_TRANSFER_FEE_MULTIPLIER"`
	RbtcMinTransferFeeMultiplier  uint64 `env:"RBTC_MIN_TRANSFER_FEE_MULTIPLIER"`
	ForceTransferAfterSeconds     uint64 `env:"COLD_WALLET_FORCE_TRANSFER_AFTER_SECONDS"`
	HotWalletLowLiquidityWarning  uint64 `env:"HOT_WALLET_LOW_LIQUIDITY_WARNING"`
	HotWalletLowLiquidityCritical uint64 `env:"HOT_WALLET_LOW_LIQUIDITY_CRITICAL"`
}

func (env *ColdWalletEnv) FillWithDefaults() *ColdWalletEnv {
	defaults := ColdWalletEnv{
		BtcMinTransferFeeMultiplier:   5,
		RbtcMinTransferFeeMultiplier:  100,
		ForceTransferAfterSeconds:     1209600, // 2 weeks (14 days * 24 hours * 60 minutes * 60 seconds)
		HotWalletLowLiquidityWarning:  3,
		HotWalletLowLiquidityCritical: 1,
	}
	env.BtcMinTransferFeeMultiplier = utils.FirstNonZero(env.BtcMinTransferFeeMultiplier, defaults.BtcMinTransferFeeMultiplier)
	env.RbtcMinTransferFeeMultiplier = utils.FirstNonZero(env.RbtcMinTransferFeeMultiplier, defaults.RbtcMinTransferFeeMultiplier)
	env.ForceTransferAfterSeconds = utils.FirstNonZero(env.ForceTransferAfterSeconds, defaults.ForceTransferAfterSeconds)
	env.HotWalletLowLiquidityWarning = utils.FirstNonZero(env.HotWalletLowLiquidityWarning, defaults.HotWalletLowLiquidityWarning)
	env.HotWalletLowLiquidityCritical = utils.FirstNonZero(env.HotWalletLowLiquidityCritical, defaults.HotWalletLowLiquidityCritical)
	if env.HotWalletLowLiquidityCritical >= env.HotWalletLowLiquidityWarning {
		log.Fatal("HOT_WALLET_LOW_LIQUIDITY_CRITICAL must be less than HOT_WALLET_LOW_LIQUIDITY_WARNING")
	}
	return env
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
