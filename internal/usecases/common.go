package usecases

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
)

// used for error logging

type UseCaseId string

const EthereumSignedMessagePrefix = "\x19Ethereum Signed Message:\n32"

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
	UpdateTrustedAccountId     UseCaseId = "UpdateTrustedAccountUseCase"
	AddTrustedAccountId        UseCaseId = "AddTrustedAccountUseCase"
	DeleteTrustedAccountId     UseCaseId = "DeleteTrustedAccountUseCase"
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
	EclipseCheckId             UseCaseId = "EclipseCheck"
	UpdateBtcReleaseId         UseCaseId = "UpdateBtcRelease"
)

var (
	NonRecoverableError             = errors.New("non recoverable")
	TxBelowMinimumError             = errors.New("requested amount below bridge's min transaction value")
	RskAddressNotSupportedError     = errors.New("rsk address not supported")
	QuoteNotFoundError              = errors.New("quote not found")
	QuoteNotAcceptedError           = errors.New("quote not accepted")
	ExpiredQuoteError               = errors.New("expired quote")
	NoLiquidityError                = errors.New("not enough liquidity")
	ProviderConfigurationError      = errors.New("pegin and pegout providers are not using the same account")
	WrongStateError                 = errors.New("quote with wrong state")
	NoEnoughConfirmationsError      = errors.New("not enough confirmations for transaction")
	InsufficientAmountError         = errors.New("insufficient amount")
	AlreadyRegisteredError          = errors.New("liquidity provider already registered")
	ProviderNotResignedError        = errors.New("provided hasn't completed resignation process")
	IllegalQuoteStateError          = errors.New("illegal quote state")
	LockingCapExceededError         = errors.New("locking cap exceeded")
	NonPositiveWeiError             = errors.New("wei value must be positive")
	EmptyConfirmationsMapError      = errors.New("confirmations map cannot be empty")
	NonPositiveConfirmationKeyError = errors.New("confirmation amount key must be positive")
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

func ValidateMinLockValue(useCase UseCaseId, bridge rootstock.Bridge, value *entities.Wei) error {
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
func RegisterCoinbaseTransaction(btcRpc blockchain.BitcoinNetwork, bridgeContract rootstock.Bridge, tx blockchain.BitcoinTransactionInformation) error {
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
func ValidateBridgeUtxoMin(bridge rootstock.Bridge, transaction blockchain.BitcoinTransactionInformation, address string) error {
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

// RecoverSignerAddress recovers the address from a signature. Important function for the management
// of trusted accounts.
func RecoverSignerAddress(quoteHash, signature string) (string, error) {
	if quoteHash == "" {
		return "", errors.New("empty hash provided")
	}

	if signature == "" {
		return "", errors.New("empty signature provided")
	}

	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return "", fmt.Errorf("error decoding signature: %w", err)
	}

	// Ethereum signatures should be 65 bytes (r,s,v) where v is the recovery ID
	if len(signatureBytes) != 65 {
		return "", fmt.Errorf("invalid signature length, expected 65 bytes, got %d", len(signatureBytes))
	}

	hashBytes, err := hex.DecodeString(quoteHash)
	if err != nil {
		return "", fmt.Errorf("error decoding hash: %w", err)
	}

	// Hash should be 32 bytes
	if len(hashBytes) != 32 {
		return "", fmt.Errorf("invalid hash length, expected 32 bytes, got %d", len(hashBytes))
	}

	// The signature's recovery ID (v) needs to be adjusted from Ethereum's convention
	// Ethereum uses 27 or 28 as the v value, but Ecrecover expects 0 or 1
	v := signatureBytes[64]
	if v >= 27 {
		signatureBytes[64] = v - 27
	}

	// Create the Ethereum prefixed message
	var buf bytes.Buffer
	buf.WriteString(EthereumSignedMessagePrefix)
	buf.Write(hashBytes)
	prefixedHash := crypto.Keccak256(buf.Bytes())

	pubKey, err := crypto.Ecrecover(prefixedHash, signatureBytes)
	if err != nil {
		return "", errors.New("error recovering public key: " + err.Error())
	}

	// Convert the public key to an Ethereum address
	pubKeyECDSA, err := crypto.UnmarshalPubkey(pubKey)
	if err != nil {
		return "", errors.New("error unmarshalling public key: " + err.Error())
	}

	address := crypto.PubkeyToAddress(*pubKeyECDSA).Hex()
	return address, nil
}

func ValidatePositiveWeiValues(useCase UseCaseId, weiValues ...*entities.Wei) error {
	if err := entities.ValidatePositiveWei(weiValues...); err != nil {
		return WrapUseCaseError(useCase, NonPositiveWeiError)
	}
	return nil
}

func ValidateConfirmations(useCase UseCaseId, confirmations liquidity_provider.ConfirmationsPerAmount) error {
	if len(confirmations) == 0 {
		return WrapUseCaseError(useCase, EmptyConfirmationsMapError)
	}
	for keyStr := range confirmations {
		intKey, err := strconv.Atoi(keyStr)
		if err != nil || intKey <= 0 {
			args := ErrorArg("key", keyStr)
			return WrapUseCaseErrorArgs(useCase, NonPositiveConfirmationKeyError, args)
		}
	}
	return nil
}
