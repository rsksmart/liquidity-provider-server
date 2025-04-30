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
	PeginQuoteStatusId         UseCaseId = "PeginQuoteStatus"
	PegoutQuoteStatusId        UseCaseId = "PegoutQuoteStatus"
	GetAvailableLiquidityId    UseCaseId = "GetAvailableLiquidity"
	UpdatePeginDepositId       UseCaseId = "UpdatePeginDeposit"
	ServerInfoId               UseCaseId = "ServerInfo"
	GetPeginReportId           UseCaseId = "GetPeginReport"
	GetPegoutReportId          UseCaseId = "GetPegoutReport"
)

var (
	NonRecoverableError         = errors.New("non recoverable")
	TxBelowMinimumError         = errors.New("requested amount below bridge's min transaction value")
	RskAddressNotSupportedError = errors.New("rsk address not supported")
	QuoteNotFoundError          = errors.New("quote not found")
	QuoteNotAcceptedError       = errors.New("quote not accepted")
	ExpiredQuoteError           = errors.New("expired quote")
	NoLiquidityError            = errors.New("not enough liquidity")
	ProviderConfigurationError  = errors.New("pegin and pegout providers are not using the same account")
	WrongStateError             = errors.New("quote with wrong state")
	NoEnoughConfirmationsError  = errors.New("not enough confirmations for transaction")
	InsufficientAmountError     = errors.New("insufficient amount")
	AlreadyRegisteredError      = errors.New("liquidity provider already registered")
	ProviderNotResignedError    = errors.New("provided hasn't completed resignation process")
	IllegalQuoteStateError      = errors.New("illegal quote state")
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

// RegisterCoinbaseTransaction registers the information of the coinbase transaction of the block of a specific transaction in the Rootstock Bridge.
// IMPORTANT: this function should not be called right now for security reasons. It is in the codebase for future compatibility but should not be used for now.
func RegisterCoinbaseTransaction(btcRpc blockchain.BitcoinNetwork, bridgeContract blockchain.RootstockBridge, tx blockchain.BitcoinTransactionInformation) error {
	if !tx.HasWitness {
		return nil
	}

	coinbaseInfo, err := btcRpc.GetCoinbaseInformation(tx.Hash)
	if err != nil {
		return err
	}
	_, err = bridgeContract.RegisterBtcCoinbaseTransaction(coinbaseInfo)
	return err
}

// ValidateBridgeUtxoMin checks that all the UTXOs to an address of a Bitcoin transaction are above the Rootstock Bridge minimum
func ValidateBridgeUtxoMin(bridge blockchain.RootstockBridge, transaction blockchain.BitcoinTransactionInformation, address string) error {
	minLockTxValueInWei, err := bridge.GetMinimumLockTxValue()
	if err != nil {
		return err
	}
	utxos := transaction.UTXOsToAddress(address)
	if len(utxos) == 0 {
		err = fmt.Errorf("no UTXO directed to address %s present in transaction", address)
		return errors.Join(TxBelowMinimumError, err)
	}
	for _, utxo := range utxos {
		if minLockTxValueInWei.Cmp(utxo) > 0 {
			err = errors.New("not all the UTXOs are above the min lock value")
			return errors.Join(TxBelowMinimumError, err)
		}
	}
	return nil
}
