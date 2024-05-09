package usecases

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"math/big"
)

// used for error logging

type UseCaseId string

const (
	GetPeginQuoteId            UseCaseId = "GetPeginQuote"
	GetPegoutQuoteId           UseCaseId = "GetPegoutQuote"
	AcceptPeginQuoteId         UseCaseId = "AcceptPeginQuote"
	AcceptPegoutQuoteId        UseCaseId = "AcceptPegoutQuote"
	ProviderDetailId           UseCaseId = "ProviderDetail"
	GetProvidersId             UseCaseId = "GetProviders"
	GetUserQuotesId            UseCaseId = "GetUserQuotes"
	ProviderResignId           UseCaseId = "ProviderResign"
	ChangeProviderStatusId     UseCaseId = "ChangeProviderStatus"
	GetCollateralId            UseCaseId = "GetCollateral"
	GetPegoutCollateralId      UseCaseId = "GetPegoutCollateral"
	AddCollateralId            UseCaseId = "AddCollateral"
	AddPegoutCollateralId      UseCaseId = "AddPegoutCollateral"
	WithdrawCollateralId       UseCaseId = "WithdrawCollateral"
	WithdrawPegoutCollateralId UseCaseId = "WithdrawPegoutCollateral"
	CallForUserId              UseCaseId = "CallForUser"
	RegisterPeginId            UseCaseId = "RegisterPegin"
	SendPegoutId               UseCaseId = "SendPegout"
	RefundPegoutId             UseCaseId = "RefundPegout"
	ProviderRegistrationId     UseCaseId = "ProviderRegistration"
	GetWatchedPeginQuoteId     UseCaseId = "GetWatchedPeginQuote"
	GetWatchedPegoutQuoteId    UseCaseId = "GetWatchedPegoutQuote"
	ExpiredPeginQuoteId        UseCaseId = "ExpiredPeginQuote"
	ExpiredPegoutQuoteId       UseCaseId = "ExpiredPegoutQuote"
	UpdatePegoutDepositId      UseCaseId = "UpdatePegoutDeposit"
	InitPegoutDepositCacheId   UseCaseId = "InitPegoutDepositCache"
	CheckLiquidityId           UseCaseId = "CheckLiquidity"
	PenalizationId             UseCaseId = "Penalization"
	SetPeginConfigId           UseCaseId = "SetPeginConfigUseCase"
	SetPegoutConfigId          UseCaseId = "SetPegoutConfigUseCase"
	SetGeneralConfigId         UseCaseId = "SetGeneralConfigUseCase"
	LoginId                    UseCaseId = "Login"
	ChangeCredentialsId        UseCaseId = "ChangeCredentials"
	DefaultCredentialsId       UseCaseId = "GenerateDefaultCredentials"
	GetManagementUiId          UseCaseId = "GetManagementUi"
	BridgePegoutId             UseCaseId = "BridgePegout"
)

var (
	NonRecoverableError         = errors.New("non recoverable")
	TxBelowMinimumError         = errors.New("requested amount below bridge's min transaction value")
	RskAddressNotSupportedError = errors.New("rsk address not supported")
	QuoteNotFoundError          = errors.New("quote not found")
	ExpiredQuoteError           = errors.New("expired quote")
	NoLiquidityError            = errors.New("not enough liquidity")
	ProviderConfigurationError  = errors.New("pegin and pegout providers are not using the same account")
	ProviderNotFoundError       = errors.New("liquidity provider not found")
	WrongStateError             = errors.New("quote with wrong state")
	NoEnoughConfirmationsError  = errors.New("not enough confirmations for transaction")
	InsufficientAmountError     = errors.New("insufficient amount")
	AlreadyRegisteredError      = errors.New("liquidity provider already registered")
	ProviderNotResignedError    = errors.New("provided hasn't completed resignation process")
)

type ErrorArgs map[string]string

func NewErrorArgs() ErrorArgs {
	return make(ErrorArgs)
}

func ErrorArg(key, value string) ErrorArgs {
	return ErrorArgs{key: value}
}

func (args ErrorArgs) String() string {
	if jsonString, err := json.Marshal(args); err != nil {
		return ""
	} else {
		return string(jsonString)
	}
}

func WrapUseCaseError(useCase UseCaseId, err error) error {
	return WrapUseCaseErrorArgs(useCase, err, make(ErrorArgs, 0))
}

func WrapUseCaseErrorArgs(useCase UseCaseId, err error, args ErrorArgs) error {
	if len(args) == 0 {
		return fmt.Errorf("%s: %w", useCase, err)
	} else {
		return fmt.Errorf("%s: %w. Args: %v", useCase, err, args)
	}
}

type DaoAmounts struct {
	DaoGasAmount *entities.Wei
	DaoFeeAmount *entities.Wei
}

func CalculateDaoAmounts(ctx context.Context, rsk blockchain.RootstockRpcServer, value *entities.Wei, daoFeePercentage uint64, feeCollectorAddress string) (DaoAmounts, error) {
	var daoGasAmount *entities.Wei
	daoFeeAmount := new(entities.Wei)
	var err error
	if daoFeePercentage == 0 {
		return DaoAmounts{
			DaoFeeAmount: entities.NewWei(0),
			DaoGasAmount: entities.NewWei(0),
		}, nil
	}

	daoFeeAmount.Mul(value, entities.NewUWei(daoFeePercentage))
	daoFeeAmount.AsBigInt().Div(daoFeeAmount.AsBigInt(), big.NewInt(100))
	daoGasAmount, err = rsk.EstimateGas(ctx, feeCollectorAddress, daoFeeAmount, make([]byte, 0))
	if err != nil {
		return DaoAmounts{}, err
	}
	return DaoAmounts{
		DaoFeeAmount: daoFeeAmount,
		DaoGasAmount: daoGasAmount,
	}, nil
}

func ValidateMinLockValue(useCase UseCaseId, bridge blockchain.RootstockBridge, value *entities.Wei) error {
	var err error
	var minLockTxValue *entities.Wei

	errorArgs := NewErrorArgs()
	if minLockTxValue, err = bridge.GetMinimumLockTxValue(); err != nil {
		return WrapUseCaseError(useCase, err)
	}
	if value.Cmp(minLockTxValue) < 0 {
		errorArgs["minimum"] = minLockTxValue.String()
		errorArgs["value"] = value.String()
		return WrapUseCaseErrorArgs(useCase, TxBelowMinimumError, errorArgs)
	}
	return nil
}

func SignConfiguration[C liquidity_provider.ConfigurationType](
	useCaseId UseCaseId,
	signer entities.Signer,
	hashFunction entities.HashFunction,
	config C,
) (entities.Signed[C], error) {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return entities.Signed[C]{}, WrapUseCaseError(useCaseId, err)
	}
	hash := hashFunction(configBytes)
	signature, err := signer.SignBytes(hash)
	if err != nil {
		return entities.Signed[C]{}, WrapUseCaseError(useCaseId, err)
	}
	signedConfig := entities.Signed[C]{
		Value:     config,
		Hash:      hex.EncodeToString(hash),
		Signature: hex.EncodeToString(signature),
	}
	return signedConfig, nil
}
