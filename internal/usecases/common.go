package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
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
)

var (
	NonRecoverableError         = errors.New("non recoverable")
	TxBelowMinimumError         = errors.New("requested amount below bridge's min pegin transaction value")
	BtcAddressNotSupportedError = errors.New("btc address not supported")
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
	AmountOutOfRangeError       = errors.New("amount out of range")
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
	jsonString, _ := json.Marshal(args)
	return string(jsonString)
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
	daoFeeAmount := new(entities.Wei)
	daoGasAmount := new(entities.Wei)
	var err error
	if daoFeePercentage == 0 {
		return DaoAmounts{}, nil
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
