// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// LiquidityBridgeContractLiquidityProvider is an auto generated low-level Go binding around an user-defined struct.
type LiquidityBridgeContractLiquidityProvider struct {
	Id                      *big.Int
	Provider                common.Address
	Name                    string
	Fee                     *big.Int
	QuoteExpiration         *big.Int
	MinTransactionValue     *big.Int
	MaxTransactionValue     *big.Int
	ApiBaseUrl              string
	Status                  bool
	ProviderType            string
}

// LiquidityBridgeContractPegOutQuote is an auto generated low-level Go binding around an user-defined struct.
type LiquidityBridgeContractPegOutQuote struct {
	LbcAddress            common.Address
	LpRskAddress          common.Address
	BtcRefundAddress      []byte
	RskRefundAddress      common.Address
	LpBtcAddress          []byte
	CallFee               *big.Int
	PenaltyFee            *big.Int
	Nonce                 int64
	DeposityAddress       []byte
	GasLimit              uint32
	Value                 *big.Int
	AgreementTimestamp    uint32
	DepositDateLimit      uint32
	DepositConfirmations  uint16
	TransferConfirmations uint16
	TransferTime          uint32
	ExpireDate            uint32
	ExpireBlock           uint32
}

// LiquidityBridgeContractPegOutQuoteState is an auto generated low-level Go binding around an user-defined struct.
type LiquidityBridgeContractPegOutQuoteState struct {
	ReceivedAmount *big.Int
	StatusCode     uint8
	Refunded       bool
}

// LiquidityBridgeContractQuote is an auto generated low-level Go binding around an user-defined struct.
type LiquidityBridgeContractQuote struct {
	FedBtcAddress               [20]byte
	LbcAddress                  common.Address
	LiquidityProviderRskAddress common.Address
	BtcRefundAddress            []byte
	RskRefundAddress            common.Address
	LiquidityProviderBtcAddress []byte
	CallFee                     *big.Int
	PenaltyFee                  *big.Int
	ContractAddress             common.Address
	Data                        []byte
	GasLimit                    uint32
	Nonce                       int64
	Value                       *big.Int
	AgreementTimestamp          uint32
	TimeForDeposit              uint32
	CallTime                    uint32
	DepositConfirmations        uint16
	CallOnRegister              bool
}

// LiquidityBridgeContractMetaData contains all meta data concerning the LiquidityBridgeContract contract.
var LiquidityBridgeContractMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BalanceDecrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BalanceIncrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"errorCode\",\"type\":\"int256\"}],\"name\":\"BridgeCapExceeded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"errorCode\",\"type\":\"int256\"}],\"name\":\"BridgeError\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"CallForUser\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"CollateralIncrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quotehash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"processed\",\"type\":\"uint256\"}],\"name\":\"PegOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PegOutBalanceDecrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PegOutBalanceIncrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"accumulatedAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"PegOutDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"PegOutRefunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"userAddress\",\"type\":\"address\"}],\"name\":\"PegOutUserRefunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PegoutCollateralIncrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PegoutWithdrawCollateral\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"liquidityProvider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"penalty\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"Penalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"Refund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Register\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"Resigned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawCollateral\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"addCollateral\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addPegoutCollateral\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"}],\"internalType\":\"structLiquidityBridgeContract.Quote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"callForUser\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lpRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"lpBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"bytes\",\"name\":\"deposityAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"}],\"internalType\":\"structLiquidityBridgeContract.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"depositPegout\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"header\",\"type\":\"bytes\"}],\"name\":\"getBtcBlockTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDustThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxQuoteValue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinPegIn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"getPegOutQuoteState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"receivedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"statusCode\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"refunded\",\"type\":\"bool\"}],\"internalType\":\"structLiquidityBridgeContract.PegOutQuoteState\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getPegoutCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProviderIds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"providerIds\",\"type\":\"uint256[]\"}],\"name\":\"getProviders\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"quoteExpiration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minTransactionValue\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxTransactionValue\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"apiBaseUrl\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"providerType\",\"type\":\"string\"}],\"internalType\":\"structLiquidityBridgeContract.LiquidityProvider[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"getRegisteredPegOutQuote\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lpRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"lpBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"bytes\",\"name\":\"deposityAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"}],\"internalType\":\"structLiquidityBridgeContract.PegOutQuote\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getResignDelayBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRewardPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lpRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"lpBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"bytes\",\"name\":\"deposityAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"}],\"internalType\":\"structLiquidityBridgeContract.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"hashPegoutQuote\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"}],\"internalType\":\"structLiquidityBridgeContract.Quote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"hashQuote\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_bridgeAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minimumCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_minimumPegIn\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"_rewardPercentage\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_resignDelayBlocks\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"_dustThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxQuoteValue\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isOperational\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isOperationalForPegout\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lpRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"lpBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"bytes\",\"name\":\"deposityAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"}],\"internalType\":\"structLiquidityBridgeContract.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"btcTxHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"btcBlockHeaderHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"partialMerkleTree\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"merkleBranchHashes\",\"type\":\"bytes32[]\"}],\"name\":\"refundPegOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lpRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"lpBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"bytes\",\"name\":\"deposityAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"}],\"internalType\":\"structLiquidityBridgeContract.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"refundUserPegOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_quoteExpiration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_minTransactionValue\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxTransactionValue\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_apiBaseUrl\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"_status\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"_providerType\",\"type\":\"string\"}],\"name\":\"register\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"}],\"internalType\":\"structLiquidityBridgeContract.Quote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"btcRawTransaction\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"partialMerkleTree\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"name\":\"registerPegIn\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lpRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"lpBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"bytes\",\"name\":\"deposityAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"}],\"internalType\":\"structLiquidityBridgeContract.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"registerPegOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resign\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_providerId\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"}],\"name\":\"setProviderStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawPegoutCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// LiquidityBridgeContractABI is the input ABI used to generate the binding from.
// Deprecated: Use LiquidityBridgeContractMetaData.ABI instead.
var LiquidityBridgeContractABI = LiquidityBridgeContractMetaData.ABI

// LiquidityBridgeContract is an auto generated Go binding around an Ethereum contract.
type LiquidityBridgeContract struct {
	LiquidityBridgeContractCaller     // Read-only binding to the contract
	LiquidityBridgeContractTransactor // Write-only binding to the contract
	LiquidityBridgeContractFilterer   // Log filterer for contract events
}

// LiquidityBridgeContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type LiquidityBridgeContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidityBridgeContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LiquidityBridgeContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidityBridgeContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LiquidityBridgeContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidityBridgeContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LiquidityBridgeContractSession struct {
	Contract     *LiquidityBridgeContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// LiquidityBridgeContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LiquidityBridgeContractCallerSession struct {
	Contract *LiquidityBridgeContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// LiquidityBridgeContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LiquidityBridgeContractTransactorSession struct {
	Contract     *LiquidityBridgeContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// LiquidityBridgeContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type LiquidityBridgeContractRaw struct {
	Contract *LiquidityBridgeContract // Generic contract binding to access the raw methods on
}

// LiquidityBridgeContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LiquidityBridgeContractCallerRaw struct {
	Contract *LiquidityBridgeContractCaller // Generic read-only contract binding to access the raw methods on
}

// LiquidityBridgeContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LiquidityBridgeContractTransactorRaw struct {
	Contract *LiquidityBridgeContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLiquidityBridgeContract creates a new instance of LiquidityBridgeContract, bound to a specific deployed contract.
func NewLiquidityBridgeContract(address common.Address, backend bind.ContractBackend) (*LiquidityBridgeContract, error) {
	contract, err := bindLiquidityBridgeContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContract{LiquidityBridgeContractCaller: LiquidityBridgeContractCaller{contract: contract}, LiquidityBridgeContractTransactor: LiquidityBridgeContractTransactor{contract: contract}, LiquidityBridgeContractFilterer: LiquidityBridgeContractFilterer{contract: contract}}, nil
}

// NewLiquidityBridgeContractCaller creates a new read-only instance of LiquidityBridgeContract, bound to a specific deployed contract.
func NewLiquidityBridgeContractCaller(address common.Address, caller bind.ContractCaller) (*LiquidityBridgeContractCaller, error) {
	contract, err := bindLiquidityBridgeContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractCaller{contract: contract}, nil
}

// NewLiquidityBridgeContractTransactor creates a new write-only instance of LiquidityBridgeContract, bound to a specific deployed contract.
func NewLiquidityBridgeContractTransactor(address common.Address, transactor bind.ContractTransactor) (*LiquidityBridgeContractTransactor, error) {
	contract, err := bindLiquidityBridgeContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractTransactor{contract: contract}, nil
}

// NewLiquidityBridgeContractFilterer creates a new log filterer instance of LiquidityBridgeContract, bound to a specific deployed contract.
func NewLiquidityBridgeContractFilterer(address common.Address, filterer bind.ContractFilterer) (*LiquidityBridgeContractFilterer, error) {
	contract, err := bindLiquidityBridgeContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractFilterer{contract: contract}, nil
}

// bindLiquidityBridgeContract binds a generic wrapper to an already deployed contract.
func bindLiquidityBridgeContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LiquidityBridgeContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LiquidityBridgeContract *LiquidityBridgeContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LiquidityBridgeContract.Contract.LiquidityBridgeContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LiquidityBridgeContract *LiquidityBridgeContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.LiquidityBridgeContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LiquidityBridgeContract *LiquidityBridgeContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.LiquidityBridgeContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LiquidityBridgeContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.contract.Transact(opts, method, params...)
}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetBalance(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getBalance", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetBalance(addr common.Address) (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetBalance(&_LiquidityBridgeContract.CallOpts, addr)
}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetBalance(addr common.Address) (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetBalance(&_LiquidityBridgeContract.CallOpts, addr)
}

// GetBridgeAddress is a free data retrieval call binding the contract method 0xfb32c508.
//
// Solidity: function getBridgeAddress() view returns(address)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetBridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getBridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBridgeAddress is a free data retrieval call binding the contract method 0xfb32c508.
//
// Solidity: function getBridgeAddress() view returns(address)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetBridgeAddress() (common.Address, error) {
	return _LiquidityBridgeContract.Contract.GetBridgeAddress(&_LiquidityBridgeContract.CallOpts)
}

// GetBridgeAddress is a free data retrieval call binding the contract method 0xfb32c508.
//
// Solidity: function getBridgeAddress() view returns(address)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetBridgeAddress() (common.Address, error) {
	return _LiquidityBridgeContract.Contract.GetBridgeAddress(&_LiquidityBridgeContract.CallOpts)
}

// GetBtcBlockTimestamp is a free data retrieval call binding the contract method 0xa0cd70fc.
//
// Solidity: function getBtcBlockTimestamp(bytes header) pure returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetBtcBlockTimestamp(opts *bind.CallOpts, header []byte) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getBtcBlockTimestamp", header)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBtcBlockTimestamp is a free data retrieval call binding the contract method 0xa0cd70fc.
//
// Solidity: function getBtcBlockTimestamp(bytes header) pure returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetBtcBlockTimestamp(header []byte) (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetBtcBlockTimestamp(&_LiquidityBridgeContract.CallOpts, header)
}

// GetBtcBlockTimestamp is a free data retrieval call binding the contract method 0xa0cd70fc.
//
// Solidity: function getBtcBlockTimestamp(bytes header) pure returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetBtcBlockTimestamp(header []byte) (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetBtcBlockTimestamp(&_LiquidityBridgeContract.CallOpts, header)
}

// GetCollateral is a free data retrieval call binding the contract method 0x9b56d6c9.
//
// Solidity: function getCollateral(address addr) view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetCollateral(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getCollateral", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCollateral is a free data retrieval call binding the contract method 0x9b56d6c9.
//
// Solidity: function getCollateral(address addr) view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetCollateral(addr common.Address) (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetCollateral(&_LiquidityBridgeContract.CallOpts, addr)
}

// GetCollateral is a free data retrieval call binding the contract method 0x9b56d6c9.
//
// Solidity: function getCollateral(address addr) view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetCollateral(addr common.Address) (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetCollateral(&_LiquidityBridgeContract.CallOpts, addr)
}

// GetDustThreshold is a free data retrieval call binding the contract method 0x33f07ad3.
//
// Solidity: function getDustThreshold() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetDustThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getDustThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDustThreshold is a free data retrieval call binding the contract method 0x33f07ad3.
//
// Solidity: function getDustThreshold() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetDustThreshold() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetDustThreshold(&_LiquidityBridgeContract.CallOpts)
}

// GetDustThreshold is a free data retrieval call binding the contract method 0x33f07ad3.
//
// Solidity: function getDustThreshold() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetDustThreshold() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetDustThreshold(&_LiquidityBridgeContract.CallOpts)
}

// GetMaxQuoteValue is a free data retrieval call binding the contract method 0xa8828575.
//
// Solidity: function getMaxQuoteValue() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetMaxQuoteValue(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getMaxQuoteValue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxQuoteValue is a free data retrieval call binding the contract method 0xa8828575.
//
// Solidity: function getMaxQuoteValue() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetMaxQuoteValue() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetMaxQuoteValue(&_LiquidityBridgeContract.CallOpts)
}

// GetMaxQuoteValue is a free data retrieval call binding the contract method 0xa8828575.
//
// Solidity: function getMaxQuoteValue() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetMaxQuoteValue() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetMaxQuoteValue(&_LiquidityBridgeContract.CallOpts)
}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetMinCollateral(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getMinCollateral")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetMinCollateral() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetMinCollateral(&_LiquidityBridgeContract.CallOpts)
}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetMinCollateral() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetMinCollateral(&_LiquidityBridgeContract.CallOpts)
}

// GetMinPegIn is a free data retrieval call binding the contract method 0xfa88dcde.
//
// Solidity: function getMinPegIn() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetMinPegIn(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getMinPegIn")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinPegIn is a free data retrieval call binding the contract method 0xfa88dcde.
//
// Solidity: function getMinPegIn() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetMinPegIn() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetMinPegIn(&_LiquidityBridgeContract.CallOpts)
}

// GetMinPegIn is a free data retrieval call binding the contract method 0xfa88dcde.
//
// Solidity: function getMinPegIn() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetMinPegIn() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetMinPegIn(&_LiquidityBridgeContract.CallOpts)
}

// GetPegOutQuoteState is a free data retrieval call binding the contract method 0xc6d1579b.
//
// Solidity: function getPegOutQuoteState(bytes32 quoteHash) view returns((uint256,uint8,bool))
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetPegOutQuoteState(opts *bind.CallOpts, quoteHash [32]byte) (LiquidityBridgeContractPegOutQuoteState, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getPegOutQuoteState", quoteHash)

	if err != nil {
		return *new(LiquidityBridgeContractPegOutQuoteState), err
	}

	out0 := *abi.ConvertType(out[0], new(LiquidityBridgeContractPegOutQuoteState)).(*LiquidityBridgeContractPegOutQuoteState)

	return out0, err

}

// GetPegOutQuoteState is a free data retrieval call binding the contract method 0xc6d1579b.
//
// Solidity: function getPegOutQuoteState(bytes32 quoteHash) view returns((uint256,uint8,bool))
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetPegOutQuoteState(quoteHash [32]byte) (LiquidityBridgeContractPegOutQuoteState, error) {
	return _LiquidityBridgeContract.Contract.GetPegOutQuoteState(&_LiquidityBridgeContract.CallOpts, quoteHash)
}

// GetPegOutQuoteState is a free data retrieval call binding the contract method 0xc6d1579b.
//
// Solidity: function getPegOutQuoteState(bytes32 quoteHash) view returns((uint256,uint8,bool))
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetPegOutQuoteState(quoteHash [32]byte) (LiquidityBridgeContractPegOutQuoteState, error) {
	return _LiquidityBridgeContract.Contract.GetPegOutQuoteState(&_LiquidityBridgeContract.CallOpts, quoteHash)
}

// GetPegoutCollateral is a free data retrieval call binding the contract method 0xbd519eff.
//
// Solidity: function getPegoutCollateral(address addr) view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetPegoutCollateral(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getPegoutCollateral", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPegoutCollateral is a free data retrieval call binding the contract method 0xbd519eff.
//
// Solidity: function getPegoutCollateral(address addr) view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetPegoutCollateral(addr common.Address) (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetPegoutCollateral(&_LiquidityBridgeContract.CallOpts, addr)
}

// GetPegoutCollateral is a free data retrieval call binding the contract method 0xbd519eff.
//
// Solidity: function getPegoutCollateral(address addr) view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetPegoutCollateral(addr common.Address) (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetPegoutCollateral(&_LiquidityBridgeContract.CallOpts, addr)
}

// GetProviderIds is a free data retrieval call binding the contract method 0x0a9cb4a7.
//
// Solidity: function getProviderIds() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetProviderIds(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getProviderIds")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetProviderIds is a free data retrieval call binding the contract method 0x0a9cb4a7.
//
// Solidity: function getProviderIds() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetProviderIds() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetProviderIds(&_LiquidityBridgeContract.CallOpts)
}

// GetProviderIds is a free data retrieval call binding the contract method 0x0a9cb4a7.
//
// Solidity: function getProviderIds() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetProviderIds() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetProviderIds(&_LiquidityBridgeContract.CallOpts)
}

// GetProviders is a free data retrieval call binding the contract method 0x668dbd83.
//
// Solidity: function getProviders(uint256[] providerIds) view returns((uint256,address,string,uint256,uint256,uint256,uint256,uint256,string,bool,string)[])
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetProviders(opts *bind.CallOpts, providerIds []*big.Int) ([]LiquidityBridgeContractLiquidityProvider, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getProviders", providerIds)

	if err != nil {
		return *new([]LiquidityBridgeContractLiquidityProvider), err
	}

	out0 := *abi.ConvertType(out[0], new([]LiquidityBridgeContractLiquidityProvider)).(*[]LiquidityBridgeContractLiquidityProvider)

	return out0, err

}

// GetProviders is a free data retrieval call binding the contract method 0x668dbd83.
//
// Solidity: function getProviders(uint256[] providerIds) view returns((uint256,address,string,uint256,uint256,uint256,uint256,uint256,string,bool,string)[])
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetProviders(providerIds []*big.Int) ([]LiquidityBridgeContractLiquidityProvider, error) {
	return _LiquidityBridgeContract.Contract.GetProviders(&_LiquidityBridgeContract.CallOpts, providerIds)
}

// GetProviders is a free data retrieval call binding the contract method 0x668dbd83.
//
// Solidity: function getProviders(uint256[] providerIds) view returns((uint256,address,string,uint256,uint256,uint256,uint256,uint256,string,bool,string)[])
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetProviders(providerIds []*big.Int) ([]LiquidityBridgeContractLiquidityProvider, error) {
	return _LiquidityBridgeContract.Contract.GetProviders(&_LiquidityBridgeContract.CallOpts, providerIds)
}

// GetRegisteredPegOutQuote is a free data retrieval call binding the contract method 0xe90d2ddb.
//
// Solidity: function getRegisteredPegOutQuote(bytes32 quoteHash) view returns((address,address,bytes,address,bytes,uint256,uint256,int64,bytes,uint32,uint256,uint32,uint32,uint16,uint16,uint32,uint32,uint32))
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetRegisteredPegOutQuote(opts *bind.CallOpts, quoteHash [32]byte) (LiquidityBridgeContractPegOutQuote, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getRegisteredPegOutQuote", quoteHash)

	if err != nil {
		return *new(LiquidityBridgeContractPegOutQuote), err
	}

	out0 := *abi.ConvertType(out[0], new(LiquidityBridgeContractPegOutQuote)).(*LiquidityBridgeContractPegOutQuote)

	return out0, err

}

// GetRegisteredPegOutQuote is a free data retrieval call binding the contract method 0xe90d2ddb.
//
// Solidity: function getRegisteredPegOutQuote(bytes32 quoteHash) view returns((address,address,bytes,address,bytes,uint256,uint256,int64,bytes,uint32,uint256,uint32,uint32,uint16,uint16,uint32,uint32,uint32))
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetRegisteredPegOutQuote(quoteHash [32]byte) (LiquidityBridgeContractPegOutQuote, error) {
	return _LiquidityBridgeContract.Contract.GetRegisteredPegOutQuote(&_LiquidityBridgeContract.CallOpts, quoteHash)
}

// GetRegisteredPegOutQuote is a free data retrieval call binding the contract method 0xe90d2ddb.
//
// Solidity: function getRegisteredPegOutQuote(bytes32 quoteHash) view returns((address,address,bytes,address,bytes,uint256,uint256,int64,bytes,uint32,uint256,uint32,uint32,uint16,uint16,uint32,uint32,uint32))
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetRegisteredPegOutQuote(quoteHash [32]byte) (LiquidityBridgeContractPegOutQuote, error) {
	return _LiquidityBridgeContract.Contract.GetRegisteredPegOutQuote(&_LiquidityBridgeContract.CallOpts, quoteHash)
}

// GetResignDelayBlocks is a free data retrieval call binding the contract method 0xbd5798c3.
//
// Solidity: function getResignDelayBlocks() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetResignDelayBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getResignDelayBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetResignDelayBlocks is a free data retrieval call binding the contract method 0xbd5798c3.
//
// Solidity: function getResignDelayBlocks() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetResignDelayBlocks() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetResignDelayBlocks(&_LiquidityBridgeContract.CallOpts)
}

// GetResignDelayBlocks is a free data retrieval call binding the contract method 0xbd5798c3.
//
// Solidity: function getResignDelayBlocks() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetResignDelayBlocks() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetResignDelayBlocks(&_LiquidityBridgeContract.CallOpts)
}

// GetRewardPercentage is a free data retrieval call binding the contract method 0xc7213163.
//
// Solidity: function getRewardPercentage() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetRewardPercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getRewardPercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRewardPercentage is a free data retrieval call binding the contract method 0xc7213163.
//
// Solidity: function getRewardPercentage() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetRewardPercentage() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetRewardPercentage(&_LiquidityBridgeContract.CallOpts)
}

// GetRewardPercentage is a free data retrieval call binding the contract method 0xc7213163.
//
// Solidity: function getRewardPercentage() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetRewardPercentage() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetRewardPercentage(&_LiquidityBridgeContract.CallOpts)
}

// HashPegoutQuote is a free data retrieval call binding the contract method 0xf691ceb2.
//
// Solidity: function hashPegoutQuote((address,address,bytes,address,bytes,uint256,uint256,int64,bytes,uint32,uint256,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote) view returns(bytes32)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) HashPegoutQuote(opts *bind.CallOpts, quote LiquidityBridgeContractPegOutQuote) ([32]byte, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "hashPegoutQuote", quote)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashPegoutQuote is a free data retrieval call binding the contract method 0xf691ceb2.
//
// Solidity: function hashPegoutQuote((address,address,bytes,address,bytes,uint256,uint256,int64,bytes,uint32,uint256,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote) view returns(bytes32)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) HashPegoutQuote(quote LiquidityBridgeContractPegOutQuote) ([32]byte, error) {
	return _LiquidityBridgeContract.Contract.HashPegoutQuote(&_LiquidityBridgeContract.CallOpts, quote)
}

// HashPegoutQuote is a free data retrieval call binding the contract method 0xf691ceb2.
//
// Solidity: function hashPegoutQuote((address,address,bytes,address,bytes,uint256,uint256,int64,bytes,uint32,uint256,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote) view returns(bytes32)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) HashPegoutQuote(quote LiquidityBridgeContractPegOutQuote) ([32]byte, error) {
	return _LiquidityBridgeContract.Contract.HashPegoutQuote(&_LiquidityBridgeContract.CallOpts, quote)
}

// HashQuote is a free data retrieval call binding the contract method 0x1b032188.
//
// Solidity: function hashQuote((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool) quote) view returns(bytes32)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) HashQuote(opts *bind.CallOpts, quote LiquidityBridgeContractQuote) ([32]byte, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "hashQuote", quote)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashQuote is a free data retrieval call binding the contract method 0x1b032188.
//
// Solidity: function hashQuote((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool) quote) view returns(bytes32)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) HashQuote(quote LiquidityBridgeContractQuote) ([32]byte, error) {
	return _LiquidityBridgeContract.Contract.HashQuote(&_LiquidityBridgeContract.CallOpts, quote)
}

// HashQuote is a free data retrieval call binding the contract method 0x1b032188.
//
// Solidity: function hashQuote((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool) quote) view returns(bytes32)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) HashQuote(quote LiquidityBridgeContractQuote) ([32]byte, error) {
	return _LiquidityBridgeContract.Contract.HashQuote(&_LiquidityBridgeContract.CallOpts, quote)
}

// IsOperational is a free data retrieval call binding the contract method 0x457385f2.
//
// Solidity: function isOperational(address addr) view returns(bool)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) IsOperational(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "isOperational", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperational is a free data retrieval call binding the contract method 0x457385f2.
//
// Solidity: function isOperational(address addr) view returns(bool)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) IsOperational(addr common.Address) (bool, error) {
	return _LiquidityBridgeContract.Contract.IsOperational(&_LiquidityBridgeContract.CallOpts, addr)
}

// IsOperational is a free data retrieval call binding the contract method 0x457385f2.
//
// Solidity: function isOperational(address addr) view returns(bool)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) IsOperational(addr common.Address) (bool, error) {
	return _LiquidityBridgeContract.Contract.IsOperational(&_LiquidityBridgeContract.CallOpts, addr)
}

// IsOperationalForPegout is a free data retrieval call binding the contract method 0x4d0ec971.
//
// Solidity: function isOperationalForPegout(address addr) view returns(bool)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) IsOperationalForPegout(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "isOperationalForPegout", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperationalForPegout is a free data retrieval call binding the contract method 0x4d0ec971.
//
// Solidity: function isOperationalForPegout(address addr) view returns(bool)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) IsOperationalForPegout(addr common.Address) (bool, error) {
	return _LiquidityBridgeContract.Contract.IsOperationalForPegout(&_LiquidityBridgeContract.CallOpts, addr)
}

// IsOperationalForPegout is a free data retrieval call binding the contract method 0x4d0ec971.
//
// Solidity: function isOperationalForPegout(address addr) view returns(bool)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) IsOperationalForPegout(addr common.Address) (bool, error) {
	return _LiquidityBridgeContract.Contract.IsOperationalForPegout(&_LiquidityBridgeContract.CallOpts, addr)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) Owner() (common.Address, error) {
	return _LiquidityBridgeContract.Contract.Owner(&_LiquidityBridgeContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) Owner() (common.Address, error) {
	return _LiquidityBridgeContract.Contract.Owner(&_LiquidityBridgeContract.CallOpts)
}

// AddCollateral is a paid mutator transaction binding the contract method 0x9e816999.
//
// Solidity: function addCollateral() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) AddCollateral(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "addCollateral")
}

// AddCollateral is a paid mutator transaction binding the contract method 0x9e816999.
//
// Solidity: function addCollateral() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) AddCollateral() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.AddCollateral(&_LiquidityBridgeContract.TransactOpts)
}

// AddCollateral is a paid mutator transaction binding the contract method 0x9e816999.
//
// Solidity: function addCollateral() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) AddCollateral() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.AddCollateral(&_LiquidityBridgeContract.TransactOpts)
}

// AddPegoutCollateral is a paid mutator transaction binding the contract method 0x4198687e.
//
// Solidity: function addPegoutCollateral() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) AddPegoutCollateral(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "addPegoutCollateral")
}

// AddPegoutCollateral is a paid mutator transaction binding the contract method 0x4198687e.
//
// Solidity: function addPegoutCollateral() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) AddPegoutCollateral() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.AddPegoutCollateral(&_LiquidityBridgeContract.TransactOpts)
}

// AddPegoutCollateral is a paid mutator transaction binding the contract method 0x4198687e.
//
// Solidity: function addPegoutCollateral() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) AddPegoutCollateral() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.AddPegoutCollateral(&_LiquidityBridgeContract.TransactOpts)
}

// CallForUser is a paid mutator transaction binding the contract method 0xac29d744.
//
// Solidity: function callForUser((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool) quote) payable returns(bool)
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) CallForUser(opts *bind.TransactOpts, quote LiquidityBridgeContractQuote) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "callForUser", quote)
}

// CallForUser is a paid mutator transaction binding the contract method 0xac29d744.
//
// Solidity: function callForUser((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool) quote) payable returns(bool)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) CallForUser(quote LiquidityBridgeContractQuote) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.CallForUser(&_LiquidityBridgeContract.TransactOpts, quote)
}

// CallForUser is a paid mutator transaction binding the contract method 0xac29d744.
//
// Solidity: function callForUser((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool) quote) payable returns(bool)
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) CallForUser(quote LiquidityBridgeContractQuote) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.CallForUser(&_LiquidityBridgeContract.TransactOpts, quote)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) Deposit() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Deposit(&_LiquidityBridgeContract.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) Deposit() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Deposit(&_LiquidityBridgeContract.TransactOpts)
}

// DepositPegout is a paid mutator transaction binding the contract method 0x09a02525.
//
// Solidity: function depositPegout(bytes32 quoteHash, address lpAddress) payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) DepositPegout(opts *bind.TransactOpts, quoteHash [32]byte, lpAddress common.Address) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "depositPegout", quoteHash, lpAddress)
}

// DepositPegout is a paid mutator transaction binding the contract method 0x09a02525.
//
// Solidity: function depositPegout(bytes32 quoteHash, address lpAddress) payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) DepositPegout(quoteHash [32]byte, lpAddress common.Address) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.DepositPegout(&_LiquidityBridgeContract.TransactOpts, quoteHash, lpAddress)
}

// DepositPegout is a paid mutator transaction binding the contract method 0x09a02525.
//
// Solidity: function depositPegout(bytes32 quoteHash, address lpAddress) payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) DepositPegout(quoteHash [32]byte, lpAddress common.Address) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.DepositPegout(&_LiquidityBridgeContract.TransactOpts, quoteHash, lpAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xf0f4419c.
//
// Solidity: function initialize(address _bridgeAddress, uint256 _minimumCollateral, uint256 _minimumPegIn, uint32 _rewardPercentage, uint32 _resignDelayBlocks, uint256 _dustThreshold, uint256 _maxQuoteValue) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) Initialize(opts *bind.TransactOpts, _bridgeAddress common.Address, _minimumCollateral *big.Int, _minimumPegIn *big.Int, _rewardPercentage uint32, _resignDelayBlocks uint32, _dustThreshold *big.Int, _maxQuoteValue *big.Int) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "initialize", _bridgeAddress, _minimumCollateral, _minimumPegIn, _rewardPercentage, _resignDelayBlocks, _dustThreshold, _maxQuoteValue)
}

// Initialize is a paid mutator transaction binding the contract method 0xf0f4419c.
//
// Solidity: function initialize(address _bridgeAddress, uint256 _minimumCollateral, uint256 _minimumPegIn, uint32 _rewardPercentage, uint32 _resignDelayBlocks, uint256 _dustThreshold, uint256 _maxQuoteValue) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) Initialize(_bridgeAddress common.Address, _minimumCollateral *big.Int, _minimumPegIn *big.Int, _rewardPercentage uint32, _resignDelayBlocks uint32, _dustThreshold *big.Int, _maxQuoteValue *big.Int) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Initialize(&_LiquidityBridgeContract.TransactOpts, _bridgeAddress, _minimumCollateral, _minimumPegIn, _rewardPercentage, _resignDelayBlocks, _dustThreshold, _maxQuoteValue)
}

// Initialize is a paid mutator transaction binding the contract method 0xf0f4419c.
//
// Solidity: function initialize(address _bridgeAddress, uint256 _minimumCollateral, uint256 _minimumPegIn, uint32 _rewardPercentage, uint32 _resignDelayBlocks, uint256 _dustThreshold, uint256 _maxQuoteValue) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) Initialize(_bridgeAddress common.Address, _minimumCollateral *big.Int, _minimumPegIn *big.Int, _rewardPercentage uint32, _resignDelayBlocks uint32, _dustThreshold *big.Int, _maxQuoteValue *big.Int) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Initialize(&_LiquidityBridgeContract.TransactOpts, _bridgeAddress, _minimumCollateral, _minimumPegIn, _rewardPercentage, _resignDelayBlocks, _dustThreshold, _maxQuoteValue)
}

// RefundPegOut is a paid mutator transaction binding the contract method 0xde51a249.
//
// Solidity: function refundPegOut((address,address,bytes,address,bytes,uint256,uint256,int64,bytes,uint32,uint256,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote, bytes32 btcTxHash, bytes32 btcBlockHeaderHash, uint256 partialMerkleTree, bytes32[] merkleBranchHashes) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) RefundPegOut(opts *bind.TransactOpts, quote LiquidityBridgeContractPegOutQuote, btcTxHash [32]byte, btcBlockHeaderHash [32]byte, partialMerkleTree *big.Int, merkleBranchHashes [][32]byte) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "refundPegOut", quote, btcTxHash, btcBlockHeaderHash, partialMerkleTree, merkleBranchHashes)
}

// RefundPegOut is a paid mutator transaction binding the contract method 0xde51a249.
//
// Solidity: function refundPegOut((address,address,bytes,address,bytes,uint256,uint256,int64,bytes,uint32,uint256,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote, bytes32 btcTxHash, bytes32 btcBlockHeaderHash, uint256 partialMerkleTree, bytes32[] merkleBranchHashes) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) RefundPegOut(quote LiquidityBridgeContractPegOutQuote, btcTxHash [32]byte, btcBlockHeaderHash [32]byte, partialMerkleTree *big.Int, merkleBranchHashes [][32]byte) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.RefundPegOut(&_LiquidityBridgeContract.TransactOpts, quote, btcTxHash, btcBlockHeaderHash, partialMerkleTree, merkleBranchHashes)
}

// RefundPegOut is a paid mutator transaction binding the contract method 0xde51a249.
//
// Solidity: function refundPegOut((address,address,bytes,address,bytes,uint256,uint256,int64,bytes,uint32,uint256,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote, bytes32 btcTxHash, bytes32 btcBlockHeaderHash, uint256 partialMerkleTree, bytes32[] merkleBranchHashes) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) RefundPegOut(quote LiquidityBridgeContractPegOutQuote, btcTxHash [32]byte, btcBlockHeaderHash [32]byte, partialMerkleTree *big.Int, merkleBranchHashes [][32]byte) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.RefundPegOut(&_LiquidityBridgeContract.TransactOpts, quote, btcTxHash, btcBlockHeaderHash, partialMerkleTree, merkleBranchHashes)
}

// RefundUserPegOut is a paid mutator transaction binding the contract method 0xabc8d2dc.
//
// Solidity: function refundUserPegOut((address,address,bytes,address,bytes,uint256,uint256,int64,bytes,uint32,uint256,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote, bytes signature) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) RefundUserPegOut(opts *bind.TransactOpts, quote LiquidityBridgeContractPegOutQuote, signature []byte) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "refundUserPegOut", quote, signature)
}

// RefundUserPegOut is a paid mutator transaction binding the contract method 0xabc8d2dc.
//
// Solidity: function refundUserPegOut((address,address,bytes,address,bytes,uint256,uint256,int64,bytes,uint32,uint256,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote, bytes signature) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) RefundUserPegOut(quote LiquidityBridgeContractPegOutQuote, signature []byte) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.RefundUserPegOut(&_LiquidityBridgeContract.TransactOpts, quote, signature)
}

// RefundUserPegOut is a paid mutator transaction binding the contract method 0xabc8d2dc.
//
// Solidity: function refundUserPegOut((address,address,bytes,address,bytes,uint256,uint256,int64,bytes,uint32,uint256,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote, bytes signature) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) RefundUserPegOut(quote LiquidityBridgeContractPegOutQuote, signature []byte) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.RefundUserPegOut(&_LiquidityBridgeContract.TransactOpts, quote, signature)
}

// Register is a paid mutator transaction binding the contract method 0xd2330788.
//
// Solidity: function register(string _name, uint256 _fee, uint256 _quoteExpiration, uint256 _minTransactionValue, uint256 _maxTransactionValue, string _apiBaseUrl, bool _status, string _providerType) payable returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) Register(opts *bind.TransactOpts, _name string, _fee *big.Int, _quoteExpiration *big.Int, _minTransactionValue *big.Int, _maxTransactionValue *big.Int, _apiBaseUrl string, _status bool, _providerType string) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "register", _name, _fee, _quoteExpiration, _minTransactionValue, _maxTransactionValue, _apiBaseUrl, _status, _providerType)
}

// Register is a paid mutator transaction binding the contract method 0xd2330788.
//
// Solidity: function register(string _name, uint256 _fee, uint256 _quoteExpiration, uint256 _minTransactionValue, uint256 _maxTransactionValue, string _apiBaseUrl, bool _status, string _providerType) payable returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) Register(_name string, _fee *big.Int, _quoteExpiration *big.Int, _minTransactionValue *big.Int, _maxTransactionValue *big.Int, _apiBaseUrl string, _status bool, _providerType string) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Register(&_LiquidityBridgeContract.TransactOpts, _name, _fee, _quoteExpiration, _minTransactionValue, _maxTransactionValue, _apiBaseUrl, _status, _providerType)
}

// Register is a paid mutator transaction binding the contract method 0xd2330788.
//
// Solidity: function register(string _name, uint256 _fee, uint256 _quoteExpiration, uint256 _minTransactionValue, uint256 _maxTransactionValue, string _apiBaseUrl, bool _status, string _providerType) payable returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) Register(_name string, _fee *big.Int, _quoteExpiration *big.Int, _minTransactionValue *big.Int, _maxTransactionValue *big.Int, _apiBaseUrl string, _status bool, _providerType string) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Register(&_LiquidityBridgeContract.TransactOpts, _name, _fee, _quoteExpiration, _minTransactionValue, _maxTransactionValue, _apiBaseUrl, _status, _providerType)
}

// RegisterPegIn is a paid mutator transaction binding the contract method 0x6e2e8c70.
//
// Solidity: function registerPegIn((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool) quote, bytes signature, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256)
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) RegisterPegIn(opts *bind.TransactOpts, quote LiquidityBridgeContractQuote, signature []byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "registerPegIn", quote, signature, btcRawTransaction, partialMerkleTree, height)
}

// RegisterPegIn is a paid mutator transaction binding the contract method 0x6e2e8c70.
//
// Solidity: function registerPegIn((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool) quote, bytes signature, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) RegisterPegIn(quote LiquidityBridgeContractQuote, signature []byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.RegisterPegIn(&_LiquidityBridgeContract.TransactOpts, quote, signature, btcRawTransaction, partialMerkleTree, height)
}

// RegisterPegIn is a paid mutator transaction binding the contract method 0x6e2e8c70.
//
// Solidity: function registerPegIn((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool) quote, bytes signature, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256)
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) RegisterPegIn(quote LiquidityBridgeContractQuote, signature []byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.RegisterPegIn(&_LiquidityBridgeContract.TransactOpts, quote, signature, btcRawTransaction, partialMerkleTree, height)
}

// RegisterPegOut is a paid mutator transaction binding the contract method 0xea85347e.
//
// Solidity: function registerPegOut((address,address,bytes,address,bytes,uint256,uint256,int64,bytes,uint32,uint256,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote, bytes signature) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) RegisterPegOut(opts *bind.TransactOpts, quote LiquidityBridgeContractPegOutQuote, signature []byte) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "registerPegOut", quote, signature)
}

// RegisterPegOut is a paid mutator transaction binding the contract method 0xea85347e.
//
// Solidity: function registerPegOut((address,address,bytes,address,bytes,uint256,uint256,int64,bytes,uint32,uint256,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote, bytes signature) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) RegisterPegOut(quote LiquidityBridgeContractPegOutQuote, signature []byte) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.RegisterPegOut(&_LiquidityBridgeContract.TransactOpts, quote, signature)
}

// RegisterPegOut is a paid mutator transaction binding the contract method 0xea85347e.
//
// Solidity: function registerPegOut((address,address,bytes,address,bytes,uint256,uint256,int64,bytes,uint32,uint256,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote, bytes signature) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) RegisterPegOut(quote LiquidityBridgeContractPegOutQuote, signature []byte) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.RegisterPegOut(&_LiquidityBridgeContract.TransactOpts, quote, signature)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.RenounceOwnership(&_LiquidityBridgeContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.RenounceOwnership(&_LiquidityBridgeContract.TransactOpts)
}

// Resign is a paid mutator transaction binding the contract method 0x69652fcf.
//
// Solidity: function resign() returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) Resign(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "resign")
}

// Resign is a paid mutator transaction binding the contract method 0x69652fcf.
//
// Solidity: function resign() returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) Resign() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Resign(&_LiquidityBridgeContract.TransactOpts)
}

// Resign is a paid mutator transaction binding the contract method 0x69652fcf.
//
// Solidity: function resign() returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) Resign() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Resign(&_LiquidityBridgeContract.TransactOpts)
}

// SetProviderStatus is a paid mutator transaction binding the contract method 0x72cbf4e8.
//
// Solidity: function setProviderStatus(uint256 _providerId, bool status) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) SetProviderStatus(opts *bind.TransactOpts, _providerId *big.Int, status bool) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "setProviderStatus", _providerId, status)
}

// SetProviderStatus is a paid mutator transaction binding the contract method 0x72cbf4e8.
//
// Solidity: function setProviderStatus(uint256 _providerId, bool status) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) SetProviderStatus(_providerId *big.Int, status bool) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.SetProviderStatus(&_LiquidityBridgeContract.TransactOpts, _providerId, status)
}

// SetProviderStatus is a paid mutator transaction binding the contract method 0x72cbf4e8.
//
// Solidity: function setProviderStatus(uint256 _providerId, bool status) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) SetProviderStatus(_providerId *big.Int, status bool) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.SetProviderStatus(&_LiquidityBridgeContract.TransactOpts, _providerId, status)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.TransferOwnership(&_LiquidityBridgeContract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.TransferOwnership(&_LiquidityBridgeContract.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "withdraw", amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Withdraw(&_LiquidityBridgeContract.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Withdraw(&_LiquidityBridgeContract.TransactOpts, amount)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x59c153be.
//
// Solidity: function withdrawCollateral() returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) WithdrawCollateral(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "withdrawCollateral")
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x59c153be.
//
// Solidity: function withdrawCollateral() returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) WithdrawCollateral() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.WithdrawCollateral(&_LiquidityBridgeContract.TransactOpts)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x59c153be.
//
// Solidity: function withdrawCollateral() returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) WithdrawCollateral() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.WithdrawCollateral(&_LiquidityBridgeContract.TransactOpts)
}

// WithdrawPegoutCollateral is a paid mutator transaction binding the contract method 0x35510a7d.
//
// Solidity: function withdrawPegoutCollateral() returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) WithdrawPegoutCollateral(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "withdrawPegoutCollateral")
}

// WithdrawPegoutCollateral is a paid mutator transaction binding the contract method 0x35510a7d.
//
// Solidity: function withdrawPegoutCollateral() returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) WithdrawPegoutCollateral() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.WithdrawPegoutCollateral(&_LiquidityBridgeContract.TransactOpts)
}

// WithdrawPegoutCollateral is a paid mutator transaction binding the contract method 0x35510a7d.
//
// Solidity: function withdrawPegoutCollateral() returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) WithdrawPegoutCollateral() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.WithdrawPegoutCollateral(&_LiquidityBridgeContract.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) Receive() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Receive(&_LiquidityBridgeContract.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) Receive() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Receive(&_LiquidityBridgeContract.TransactOpts)
}

// LiquidityBridgeContractBalanceDecreaseIterator is returned from FilterBalanceDecrease and is used to iterate over the raw logs and unpacked data for BalanceDecrease events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractBalanceDecreaseIterator struct {
	Event *LiquidityBridgeContractBalanceDecrease // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractBalanceDecreaseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractBalanceDecrease)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractBalanceDecrease)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractBalanceDecreaseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractBalanceDecreaseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractBalanceDecrease represents a BalanceDecrease event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractBalanceDecrease struct {
	Dest   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBalanceDecrease is a free log retrieval operation binding the contract event 0x8e51a4493a6f66c76e13fd9e3b754eafbfe21343c04508deb61be8ccc0064587.
//
// Solidity: event BalanceDecrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterBalanceDecrease(opts *bind.FilterOpts) (*LiquidityBridgeContractBalanceDecreaseIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "BalanceDecrease")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractBalanceDecreaseIterator{contract: _LiquidityBridgeContract.contract, event: "BalanceDecrease", logs: logs, sub: sub}, nil
}

// WatchBalanceDecrease is a free log subscription operation binding the contract event 0x8e51a4493a6f66c76e13fd9e3b754eafbfe21343c04508deb61be8ccc0064587.
//
// Solidity: event BalanceDecrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchBalanceDecrease(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractBalanceDecrease) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "BalanceDecrease")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractBalanceDecrease)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "BalanceDecrease", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBalanceDecrease is a log parse operation binding the contract event 0x8e51a4493a6f66c76e13fd9e3b754eafbfe21343c04508deb61be8ccc0064587.
//
// Solidity: event BalanceDecrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseBalanceDecrease(log types.Log) (*LiquidityBridgeContractBalanceDecrease, error) {
	event := new(LiquidityBridgeContractBalanceDecrease)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "BalanceDecrease", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractBalanceIncreaseIterator is returned from FilterBalanceIncrease and is used to iterate over the raw logs and unpacked data for BalanceIncrease events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractBalanceIncreaseIterator struct {
	Event *LiquidityBridgeContractBalanceIncrease // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractBalanceIncreaseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractBalanceIncrease)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractBalanceIncrease)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractBalanceIncreaseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractBalanceIncreaseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractBalanceIncrease represents a BalanceIncrease event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractBalanceIncrease struct {
	Dest   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBalanceIncrease is a free log retrieval operation binding the contract event 0x42cfb81a915ac5a674852db250bf722637bee705a267633b68cab3a2dde06f53.
//
// Solidity: event BalanceIncrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterBalanceIncrease(opts *bind.FilterOpts) (*LiquidityBridgeContractBalanceIncreaseIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "BalanceIncrease")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractBalanceIncreaseIterator{contract: _LiquidityBridgeContract.contract, event: "BalanceIncrease", logs: logs, sub: sub}, nil
}

// WatchBalanceIncrease is a free log subscription operation binding the contract event 0x42cfb81a915ac5a674852db250bf722637bee705a267633b68cab3a2dde06f53.
//
// Solidity: event BalanceIncrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchBalanceIncrease(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractBalanceIncrease) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "BalanceIncrease")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractBalanceIncrease)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "BalanceIncrease", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBalanceIncrease is a log parse operation binding the contract event 0x42cfb81a915ac5a674852db250bf722637bee705a267633b68cab3a2dde06f53.
//
// Solidity: event BalanceIncrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseBalanceIncrease(log types.Log) (*LiquidityBridgeContractBalanceIncrease, error) {
	event := new(LiquidityBridgeContractBalanceIncrease)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "BalanceIncrease", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractBridgeCapExceededIterator is returned from FilterBridgeCapExceeded and is used to iterate over the raw logs and unpacked data for BridgeCapExceeded events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractBridgeCapExceededIterator struct {
	Event *LiquidityBridgeContractBridgeCapExceeded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractBridgeCapExceededIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractBridgeCapExceeded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractBridgeCapExceeded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractBridgeCapExceededIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractBridgeCapExceededIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractBridgeCapExceeded represents a BridgeCapExceeded event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractBridgeCapExceeded struct {
	QuoteHash [32]byte
	ErrorCode *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBridgeCapExceeded is a free log retrieval operation binding the contract event 0xfb209329d5ab5b7bcb2e92f45f4534814b6e68fa5ad1f171dabc1d17d26f0ebe.
//
// Solidity: event BridgeCapExceeded(bytes32 quoteHash, int256 errorCode)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterBridgeCapExceeded(opts *bind.FilterOpts) (*LiquidityBridgeContractBridgeCapExceededIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "BridgeCapExceeded")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractBridgeCapExceededIterator{contract: _LiquidityBridgeContract.contract, event: "BridgeCapExceeded", logs: logs, sub: sub}, nil
}

// WatchBridgeCapExceeded is a free log subscription operation binding the contract event 0xfb209329d5ab5b7bcb2e92f45f4534814b6e68fa5ad1f171dabc1d17d26f0ebe.
//
// Solidity: event BridgeCapExceeded(bytes32 quoteHash, int256 errorCode)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchBridgeCapExceeded(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractBridgeCapExceeded) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "BridgeCapExceeded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractBridgeCapExceeded)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "BridgeCapExceeded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBridgeCapExceeded is a log parse operation binding the contract event 0xfb209329d5ab5b7bcb2e92f45f4534814b6e68fa5ad1f171dabc1d17d26f0ebe.
//
// Solidity: event BridgeCapExceeded(bytes32 quoteHash, int256 errorCode)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseBridgeCapExceeded(log types.Log) (*LiquidityBridgeContractBridgeCapExceeded, error) {
	event := new(LiquidityBridgeContractBridgeCapExceeded)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "BridgeCapExceeded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractBridgeErrorIterator is returned from FilterBridgeError and is used to iterate over the raw logs and unpacked data for BridgeError events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractBridgeErrorIterator struct {
	Event *LiquidityBridgeContractBridgeError // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractBridgeErrorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractBridgeError)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractBridgeError)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractBridgeErrorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractBridgeErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractBridgeError represents a BridgeError event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractBridgeError struct {
	QuoteHash [32]byte
	ErrorCode *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBridgeError is a free log retrieval operation binding the contract event 0xa0f8bae2e63548ef07d0f252b12cda04ea27800c1e2605af7b822cdef64e756f.
//
// Solidity: event BridgeError(bytes32 quoteHash, int256 errorCode)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterBridgeError(opts *bind.FilterOpts) (*LiquidityBridgeContractBridgeErrorIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "BridgeError")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractBridgeErrorIterator{contract: _LiquidityBridgeContract.contract, event: "BridgeError", logs: logs, sub: sub}, nil
}

// WatchBridgeError is a free log subscription operation binding the contract event 0xa0f8bae2e63548ef07d0f252b12cda04ea27800c1e2605af7b822cdef64e756f.
//
// Solidity: event BridgeError(bytes32 quoteHash, int256 errorCode)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchBridgeError(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractBridgeError) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "BridgeError")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractBridgeError)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "BridgeError", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBridgeError is a log parse operation binding the contract event 0xa0f8bae2e63548ef07d0f252b12cda04ea27800c1e2605af7b822cdef64e756f.
//
// Solidity: event BridgeError(bytes32 quoteHash, int256 errorCode)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseBridgeError(log types.Log) (*LiquidityBridgeContractBridgeError, error) {
	event := new(LiquidityBridgeContractBridgeError)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "BridgeError", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractCallForUserIterator is returned from FilterCallForUser and is used to iterate over the raw logs and unpacked data for CallForUser events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractCallForUserIterator struct {
	Event *LiquidityBridgeContractCallForUser // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractCallForUserIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractCallForUser)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractCallForUser)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractCallForUserIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractCallForUserIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractCallForUser represents a CallForUser event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractCallForUser struct {
	From      common.Address
	Dest      common.Address
	GasLimit  *big.Int
	Value     *big.Int
	Data      []byte
	Success   bool
	QuoteHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCallForUser is a free log retrieval operation binding the contract event 0xbfc7404e6fe464f0646fe2c6ab942b92d56be722bb39f8c6bc4830d2d32fb80d.
//
// Solidity: event CallForUser(address indexed from, address indexed dest, uint256 gasLimit, uint256 value, bytes data, bool success, bytes32 quoteHash)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterCallForUser(opts *bind.FilterOpts, from []common.Address, dest []common.Address) (*LiquidityBridgeContractCallForUserIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var destRule []interface{}
	for _, destItem := range dest {
		destRule = append(destRule, destItem)
	}

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "CallForUser", fromRule, destRule)
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractCallForUserIterator{contract: _LiquidityBridgeContract.contract, event: "CallForUser", logs: logs, sub: sub}, nil
}

// WatchCallForUser is a free log subscription operation binding the contract event 0xbfc7404e6fe464f0646fe2c6ab942b92d56be722bb39f8c6bc4830d2d32fb80d.
//
// Solidity: event CallForUser(address indexed from, address indexed dest, uint256 gasLimit, uint256 value, bytes data, bool success, bytes32 quoteHash)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchCallForUser(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractCallForUser, from []common.Address, dest []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var destRule []interface{}
	for _, destItem := range dest {
		destRule = append(destRule, destItem)
	}

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "CallForUser", fromRule, destRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractCallForUser)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "CallForUser", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCallForUser is a log parse operation binding the contract event 0xbfc7404e6fe464f0646fe2c6ab942b92d56be722bb39f8c6bc4830d2d32fb80d.
//
// Solidity: event CallForUser(address indexed from, address indexed dest, uint256 gasLimit, uint256 value, bytes data, bool success, bytes32 quoteHash)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseCallForUser(log types.Log) (*LiquidityBridgeContractCallForUser, error) {
	event := new(LiquidityBridgeContractCallForUser)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "CallForUser", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractCollateralIncreaseIterator is returned from FilterCollateralIncrease and is used to iterate over the raw logs and unpacked data for CollateralIncrease events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractCollateralIncreaseIterator struct {
	Event *LiquidityBridgeContractCollateralIncrease // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractCollateralIncreaseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractCollateralIncrease)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractCollateralIncrease)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractCollateralIncreaseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractCollateralIncreaseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractCollateralIncrease represents a CollateralIncrease event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractCollateralIncrease struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCollateralIncrease is a free log retrieval operation binding the contract event 0x456e0f4ea86ac283092c750200e8c877f6ad8901ae575f90e02081acd455af84.
//
// Solidity: event CollateralIncrease(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterCollateralIncrease(opts *bind.FilterOpts) (*LiquidityBridgeContractCollateralIncreaseIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "CollateralIncrease")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractCollateralIncreaseIterator{contract: _LiquidityBridgeContract.contract, event: "CollateralIncrease", logs: logs, sub: sub}, nil
}

// WatchCollateralIncrease is a free log subscription operation binding the contract event 0x456e0f4ea86ac283092c750200e8c877f6ad8901ae575f90e02081acd455af84.
//
// Solidity: event CollateralIncrease(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchCollateralIncrease(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractCollateralIncrease) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "CollateralIncrease")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractCollateralIncrease)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "CollateralIncrease", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCollateralIncrease is a log parse operation binding the contract event 0x456e0f4ea86ac283092c750200e8c877f6ad8901ae575f90e02081acd455af84.
//
// Solidity: event CollateralIncrease(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseCollateralIncrease(log types.Log) (*LiquidityBridgeContractCollateralIncrease, error) {
	event := new(LiquidityBridgeContractCollateralIncrease)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "CollateralIncrease", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractDepositIterator struct {
	Event *LiquidityBridgeContractDeposit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractDeposit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractDeposit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractDeposit represents a Deposit event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractDeposit struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterDeposit(opts *bind.FilterOpts) (*LiquidityBridgeContractDepositIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractDepositIterator{contract: _LiquidityBridgeContract.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractDeposit) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractDeposit)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Deposit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDeposit is a log parse operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseDeposit(log types.Log) (*LiquidityBridgeContractDeposit, error) {
	event := new(LiquidityBridgeContractDeposit)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractInitializedIterator struct {
	Event *LiquidityBridgeContractInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractInitialized represents a Initialized event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterInitialized(opts *bind.FilterOpts) (*LiquidityBridgeContractInitializedIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractInitializedIterator{contract: _LiquidityBridgeContract.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractInitialized) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractInitialized)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseInitialized(log types.Log) (*LiquidityBridgeContractInitialized, error) {
	event := new(LiquidityBridgeContractInitialized)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractOwnershipTransferredIterator struct {
	Event *LiquidityBridgeContractOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractOwnershipTransferred represents a OwnershipTransferred event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*LiquidityBridgeContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractOwnershipTransferredIterator{contract: _LiquidityBridgeContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractOwnershipTransferred)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseOwnershipTransferred(log types.Log) (*LiquidityBridgeContractOwnershipTransferred, error) {
	event := new(LiquidityBridgeContractOwnershipTransferred)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractPegOutIterator is returned from FilterPegOut and is used to iterate over the raw logs and unpacked data for PegOut events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegOutIterator struct {
	Event *LiquidityBridgeContractPegOut // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractPegOutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractPegOut)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractPegOut)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractPegOutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractPegOutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractPegOut represents a PegOut event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegOut struct {
	From      common.Address
	Amount    *big.Int
	Quotehash [32]byte
	Processed *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPegOut is a free log retrieval operation binding the contract event 0xed3e6789842b3369f529c844ab6575be53f29ffeabd4d8b84c04c8431535b1ec.
//
// Solidity: event PegOut(address from, uint256 amount, bytes32 quotehash, uint256 processed)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterPegOut(opts *bind.FilterOpts) (*LiquidityBridgeContractPegOutIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "PegOut")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractPegOutIterator{contract: _LiquidityBridgeContract.contract, event: "PegOut", logs: logs, sub: sub}, nil
}

// WatchPegOut is a free log subscription operation binding the contract event 0xed3e6789842b3369f529c844ab6575be53f29ffeabd4d8b84c04c8431535b1ec.
//
// Solidity: event PegOut(address from, uint256 amount, bytes32 quotehash, uint256 processed)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchPegOut(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractPegOut) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "PegOut")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractPegOut)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegOut", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePegOut is a log parse operation binding the contract event 0xed3e6789842b3369f529c844ab6575be53f29ffeabd4d8b84c04c8431535b1ec.
//
// Solidity: event PegOut(address from, uint256 amount, bytes32 quotehash, uint256 processed)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParsePegOut(log types.Log) (*LiquidityBridgeContractPegOut, error) {
	event := new(LiquidityBridgeContractPegOut)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegOut", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractPegOutBalanceDecreaseIterator is returned from FilterPegOutBalanceDecrease and is used to iterate over the raw logs and unpacked data for PegOutBalanceDecrease events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegOutBalanceDecreaseIterator struct {
	Event *LiquidityBridgeContractPegOutBalanceDecrease // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractPegOutBalanceDecreaseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractPegOutBalanceDecrease)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractPegOutBalanceDecrease)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractPegOutBalanceDecreaseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractPegOutBalanceDecreaseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractPegOutBalanceDecrease represents a PegOutBalanceDecrease event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegOutBalanceDecrease struct {
	Dest   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPegOutBalanceDecrease is a free log retrieval operation binding the contract event 0x7862918831efd3b8f1079c5d7e9344c4a47c2940a11a6b449a03224c4fa38334.
//
// Solidity: event PegOutBalanceDecrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterPegOutBalanceDecrease(opts *bind.FilterOpts) (*LiquidityBridgeContractPegOutBalanceDecreaseIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "PegOutBalanceDecrease")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractPegOutBalanceDecreaseIterator{contract: _LiquidityBridgeContract.contract, event: "PegOutBalanceDecrease", logs: logs, sub: sub}, nil
}

// WatchPegOutBalanceDecrease is a free log subscription operation binding the contract event 0x7862918831efd3b8f1079c5d7e9344c4a47c2940a11a6b449a03224c4fa38334.
//
// Solidity: event PegOutBalanceDecrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchPegOutBalanceDecrease(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractPegOutBalanceDecrease) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "PegOutBalanceDecrease")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractPegOutBalanceDecrease)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegOutBalanceDecrease", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePegOutBalanceDecrease is a log parse operation binding the contract event 0x7862918831efd3b8f1079c5d7e9344c4a47c2940a11a6b449a03224c4fa38334.
//
// Solidity: event PegOutBalanceDecrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParsePegOutBalanceDecrease(log types.Log) (*LiquidityBridgeContractPegOutBalanceDecrease, error) {
	event := new(LiquidityBridgeContractPegOutBalanceDecrease)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegOutBalanceDecrease", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractPegOutBalanceIncreaseIterator is returned from FilterPegOutBalanceIncrease and is used to iterate over the raw logs and unpacked data for PegOutBalanceIncrease events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegOutBalanceIncreaseIterator struct {
	Event *LiquidityBridgeContractPegOutBalanceIncrease // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractPegOutBalanceIncreaseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractPegOutBalanceIncrease)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractPegOutBalanceIncrease)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractPegOutBalanceIncreaseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractPegOutBalanceIncreaseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractPegOutBalanceIncrease represents a PegOutBalanceIncrease event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegOutBalanceIncrease struct {
	Dest   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPegOutBalanceIncrease is a free log retrieval operation binding the contract event 0x7eb93adae7d5cfb024d663c84ccba97b3104e572bbe138ca323c854f666ae531.
//
// Solidity: event PegOutBalanceIncrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterPegOutBalanceIncrease(opts *bind.FilterOpts) (*LiquidityBridgeContractPegOutBalanceIncreaseIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "PegOutBalanceIncrease")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractPegOutBalanceIncreaseIterator{contract: _LiquidityBridgeContract.contract, event: "PegOutBalanceIncrease", logs: logs, sub: sub}, nil
}

// WatchPegOutBalanceIncrease is a free log subscription operation binding the contract event 0x7eb93adae7d5cfb024d663c84ccba97b3104e572bbe138ca323c854f666ae531.
//
// Solidity: event PegOutBalanceIncrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchPegOutBalanceIncrease(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractPegOutBalanceIncrease) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "PegOutBalanceIncrease")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractPegOutBalanceIncrease)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegOutBalanceIncrease", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePegOutBalanceIncrease is a log parse operation binding the contract event 0x7eb93adae7d5cfb024d663c84ccba97b3104e572bbe138ca323c854f666ae531.
//
// Solidity: event PegOutBalanceIncrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParsePegOutBalanceIncrease(log types.Log) (*LiquidityBridgeContractPegOutBalanceIncrease, error) {
	event := new(LiquidityBridgeContractPegOutBalanceIncrease)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegOutBalanceIncrease", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractPegOutDepositIterator is returned from FilterPegOutDeposit and is used to iterate over the raw logs and unpacked data for PegOutDeposit events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegOutDepositIterator struct {
	Event *LiquidityBridgeContractPegOutDeposit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractPegOutDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractPegOutDeposit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractPegOutDeposit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractPegOutDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractPegOutDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractPegOutDeposit represents a PegOutDeposit event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegOutDeposit struct {
	QuoteHash         [32]byte
	AccumulatedAmount *big.Int
	Timestamp         *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterPegOutDeposit is a free log retrieval operation binding the contract event 0x8be9eb0562490aeda2f448fddb7c23045f3a1b431cffdfe5700f12a0cec580b1.
//
// Solidity: event PegOutDeposit(bytes32 quoteHash, uint256 accumulatedAmount, uint256 timestamp)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterPegOutDeposit(opts *bind.FilterOpts) (*LiquidityBridgeContractPegOutDepositIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "PegOutDeposit")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractPegOutDepositIterator{contract: _LiquidityBridgeContract.contract, event: "PegOutDeposit", logs: logs, sub: sub}, nil
}

// WatchPegOutDeposit is a free log subscription operation binding the contract event 0x8be9eb0562490aeda2f448fddb7c23045f3a1b431cffdfe5700f12a0cec580b1.
//
// Solidity: event PegOutDeposit(bytes32 quoteHash, uint256 accumulatedAmount, uint256 timestamp)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchPegOutDeposit(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractPegOutDeposit) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "PegOutDeposit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractPegOutDeposit)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegOutDeposit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePegOutDeposit is a log parse operation binding the contract event 0x8be9eb0562490aeda2f448fddb7c23045f3a1b431cffdfe5700f12a0cec580b1.
//
// Solidity: event PegOutDeposit(bytes32 quoteHash, uint256 accumulatedAmount, uint256 timestamp)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParsePegOutDeposit(log types.Log) (*LiquidityBridgeContractPegOutDeposit, error) {
	event := new(LiquidityBridgeContractPegOutDeposit)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegOutDeposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractPegOutRefundedIterator is returned from FilterPegOutRefunded and is used to iterate over the raw logs and unpacked data for PegOutRefunded events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegOutRefundedIterator struct {
	Event *LiquidityBridgeContractPegOutRefunded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractPegOutRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractPegOutRefunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractPegOutRefunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractPegOutRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractPegOutRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractPegOutRefunded represents a PegOutRefunded event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegOutRefunded struct {
	QuoteHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPegOutRefunded is a free log retrieval operation binding the contract event 0xb781856ec73fd0dc39351043d1634ea22cd3277b0866ab93e7ec1801766bb384.
//
// Solidity: event PegOutRefunded(bytes32 quoteHash)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterPegOutRefunded(opts *bind.FilterOpts) (*LiquidityBridgeContractPegOutRefundedIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "PegOutRefunded")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractPegOutRefundedIterator{contract: _LiquidityBridgeContract.contract, event: "PegOutRefunded", logs: logs, sub: sub}, nil
}

// WatchPegOutRefunded is a free log subscription operation binding the contract event 0xb781856ec73fd0dc39351043d1634ea22cd3277b0866ab93e7ec1801766bb384.
//
// Solidity: event PegOutRefunded(bytes32 quoteHash)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchPegOutRefunded(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractPegOutRefunded) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "PegOutRefunded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractPegOutRefunded)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegOutRefunded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePegOutRefunded is a log parse operation binding the contract event 0xb781856ec73fd0dc39351043d1634ea22cd3277b0866ab93e7ec1801766bb384.
//
// Solidity: event PegOutRefunded(bytes32 quoteHash)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParsePegOutRefunded(log types.Log) (*LiquidityBridgeContractPegOutRefunded, error) {
	event := new(LiquidityBridgeContractPegOutRefunded)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegOutRefunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractPegOutUserRefundedIterator is returned from FilterPegOutUserRefunded and is used to iterate over the raw logs and unpacked data for PegOutUserRefunded events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegOutUserRefundedIterator struct {
	Event *LiquidityBridgeContractPegOutUserRefunded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractPegOutUserRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractPegOutUserRefunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractPegOutUserRefunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractPegOutUserRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractPegOutUserRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractPegOutUserRefunded represents a PegOutUserRefunded event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegOutUserRefunded struct {
	QuoteHash   [32]byte
	Value       *big.Int
	UserAddress common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPegOutUserRefunded is a free log retrieval operation binding the contract event 0x9ccbeffc442024e2a6ade18ff0978af9a4c4d6562ae38adb51ccf8256cf42b41.
//
// Solidity: event PegOutUserRefunded(bytes32 quoteHash, uint256 value, address userAddress)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterPegOutUserRefunded(opts *bind.FilterOpts) (*LiquidityBridgeContractPegOutUserRefundedIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "PegOutUserRefunded")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractPegOutUserRefundedIterator{contract: _LiquidityBridgeContract.contract, event: "PegOutUserRefunded", logs: logs, sub: sub}, nil
}

// WatchPegOutUserRefunded is a free log subscription operation binding the contract event 0x9ccbeffc442024e2a6ade18ff0978af9a4c4d6562ae38adb51ccf8256cf42b41.
//
// Solidity: event PegOutUserRefunded(bytes32 quoteHash, uint256 value, address userAddress)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchPegOutUserRefunded(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractPegOutUserRefunded) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "PegOutUserRefunded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractPegOutUserRefunded)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegOutUserRefunded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePegOutUserRefunded is a log parse operation binding the contract event 0x9ccbeffc442024e2a6ade18ff0978af9a4c4d6562ae38adb51ccf8256cf42b41.
//
// Solidity: event PegOutUserRefunded(bytes32 quoteHash, uint256 value, address userAddress)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParsePegOutUserRefunded(log types.Log) (*LiquidityBridgeContractPegOutUserRefunded, error) {
	event := new(LiquidityBridgeContractPegOutUserRefunded)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegOutUserRefunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractPegoutCollateralIncreaseIterator is returned from FilterPegoutCollateralIncrease and is used to iterate over the raw logs and unpacked data for PegoutCollateralIncrease events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegoutCollateralIncreaseIterator struct {
	Event *LiquidityBridgeContractPegoutCollateralIncrease // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractPegoutCollateralIncreaseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractPegoutCollateralIncrease)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractPegoutCollateralIncrease)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractPegoutCollateralIncreaseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractPegoutCollateralIncreaseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractPegoutCollateralIncrease represents a PegoutCollateralIncrease event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegoutCollateralIncrease struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPegoutCollateralIncrease is a free log retrieval operation binding the contract event 0x873d5a2949567203ad4f0cceef41c2813c87b9a397ee777d87a8acdaec2c6fa9.
//
// Solidity: event PegoutCollateralIncrease(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterPegoutCollateralIncrease(opts *bind.FilterOpts) (*LiquidityBridgeContractPegoutCollateralIncreaseIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "PegoutCollateralIncrease")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractPegoutCollateralIncreaseIterator{contract: _LiquidityBridgeContract.contract, event: "PegoutCollateralIncrease", logs: logs, sub: sub}, nil
}

// WatchPegoutCollateralIncrease is a free log subscription operation binding the contract event 0x873d5a2949567203ad4f0cceef41c2813c87b9a397ee777d87a8acdaec2c6fa9.
//
// Solidity: event PegoutCollateralIncrease(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchPegoutCollateralIncrease(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractPegoutCollateralIncrease) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "PegoutCollateralIncrease")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractPegoutCollateralIncrease)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegoutCollateralIncrease", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePegoutCollateralIncrease is a log parse operation binding the contract event 0x873d5a2949567203ad4f0cceef41c2813c87b9a397ee777d87a8acdaec2c6fa9.
//
// Solidity: event PegoutCollateralIncrease(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParsePegoutCollateralIncrease(log types.Log) (*LiquidityBridgeContractPegoutCollateralIncrease, error) {
	event := new(LiquidityBridgeContractPegoutCollateralIncrease)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegoutCollateralIncrease", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractPegoutWithdrawCollateralIterator is returned from FilterPegoutWithdrawCollateral and is used to iterate over the raw logs and unpacked data for PegoutWithdrawCollateral events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegoutWithdrawCollateralIterator struct {
	Event *LiquidityBridgeContractPegoutWithdrawCollateral // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractPegoutWithdrawCollateralIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractPegoutWithdrawCollateral)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractPegoutWithdrawCollateral)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractPegoutWithdrawCollateralIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractPegoutWithdrawCollateralIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractPegoutWithdrawCollateral represents a PegoutWithdrawCollateral event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegoutWithdrawCollateral struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPegoutWithdrawCollateral is a free log retrieval operation binding the contract event 0xfc72299650b405e7b0480ca8fb0fb3948fb10a77ac02f797cc9de1f4aaa55db7.
//
// Solidity: event PegoutWithdrawCollateral(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterPegoutWithdrawCollateral(opts *bind.FilterOpts) (*LiquidityBridgeContractPegoutWithdrawCollateralIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "PegoutWithdrawCollateral")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractPegoutWithdrawCollateralIterator{contract: _LiquidityBridgeContract.contract, event: "PegoutWithdrawCollateral", logs: logs, sub: sub}, nil
}

// WatchPegoutWithdrawCollateral is a free log subscription operation binding the contract event 0xfc72299650b405e7b0480ca8fb0fb3948fb10a77ac02f797cc9de1f4aaa55db7.
//
// Solidity: event PegoutWithdrawCollateral(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchPegoutWithdrawCollateral(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractPegoutWithdrawCollateral) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "PegoutWithdrawCollateral")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractPegoutWithdrawCollateral)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegoutWithdrawCollateral", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePegoutWithdrawCollateral is a log parse operation binding the contract event 0xfc72299650b405e7b0480ca8fb0fb3948fb10a77ac02f797cc9de1f4aaa55db7.
//
// Solidity: event PegoutWithdrawCollateral(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParsePegoutWithdrawCollateral(log types.Log) (*LiquidityBridgeContractPegoutWithdrawCollateral, error) {
	event := new(LiquidityBridgeContractPegoutWithdrawCollateral)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegoutWithdrawCollateral", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractPenalizedIterator is returned from FilterPenalized and is used to iterate over the raw logs and unpacked data for Penalized events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPenalizedIterator struct {
	Event *LiquidityBridgeContractPenalized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractPenalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractPenalized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractPenalized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractPenalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractPenalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractPenalized represents a Penalized event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPenalized struct {
	LiquidityProvider common.Address
	Penalty           *big.Int
	QuoteHash         [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterPenalized is a free log retrieval operation binding the contract event 0x9685484093cc596fdaeab51abf645b1753dbb7d869bfd2eb21e2c646e47a36f4.
//
// Solidity: event Penalized(address liquidityProvider, uint256 penalty, bytes32 quoteHash)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterPenalized(opts *bind.FilterOpts) (*LiquidityBridgeContractPenalizedIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "Penalized")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractPenalizedIterator{contract: _LiquidityBridgeContract.contract, event: "Penalized", logs: logs, sub: sub}, nil
}

// WatchPenalized is a free log subscription operation binding the contract event 0x9685484093cc596fdaeab51abf645b1753dbb7d869bfd2eb21e2c646e47a36f4.
//
// Solidity: event Penalized(address liquidityProvider, uint256 penalty, bytes32 quoteHash)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchPenalized(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractPenalized) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "Penalized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractPenalized)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Penalized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePenalized is a log parse operation binding the contract event 0x9685484093cc596fdaeab51abf645b1753dbb7d869bfd2eb21e2c646e47a36f4.
//
// Solidity: event Penalized(address liquidityProvider, uint256 penalty, bytes32 quoteHash)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParsePenalized(log types.Log) (*LiquidityBridgeContractPenalized, error) {
	event := new(LiquidityBridgeContractPenalized)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Penalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractRefundIterator is returned from FilterRefund and is used to iterate over the raw logs and unpacked data for Refund events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractRefundIterator struct {
	Event *LiquidityBridgeContractRefund // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractRefundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractRefund)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractRefund)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractRefundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractRefundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractRefund represents a Refund event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractRefund struct {
	Dest      common.Address
	Amount    *big.Int
	Success   bool
	QuoteHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRefund is a free log retrieval operation binding the contract event 0x3052ea2f7e0d74fdc1c1e1f858ff1ae3d91ab1609717c3efedb95db603b255f6.
//
// Solidity: event Refund(address dest, uint256 amount, bool success, bytes32 quoteHash)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterRefund(opts *bind.FilterOpts) (*LiquidityBridgeContractRefundIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "Refund")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractRefundIterator{contract: _LiquidityBridgeContract.contract, event: "Refund", logs: logs, sub: sub}, nil
}

// WatchRefund is a free log subscription operation binding the contract event 0x3052ea2f7e0d74fdc1c1e1f858ff1ae3d91ab1609717c3efedb95db603b255f6.
//
// Solidity: event Refund(address dest, uint256 amount, bool success, bytes32 quoteHash)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchRefund(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractRefund) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "Refund")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractRefund)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Refund", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRefund is a log parse operation binding the contract event 0x3052ea2f7e0d74fdc1c1e1f858ff1ae3d91ab1609717c3efedb95db603b255f6.
//
// Solidity: event Refund(address dest, uint256 amount, bool success, bytes32 quoteHash)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseRefund(log types.Log) (*LiquidityBridgeContractRefund, error) {
	event := new(LiquidityBridgeContractRefund)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Refund", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractRegisterIterator is returned from FilterRegister and is used to iterate over the raw logs and unpacked data for Register events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractRegisterIterator struct {
	Event *LiquidityBridgeContractRegister // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractRegisterIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractRegister)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractRegister)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractRegisterIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractRegisterIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractRegister represents a Register event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractRegister struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRegister is a free log retrieval operation binding the contract event 0xa9d44d6e13bb3fee938c3f66d1103e91f8dc6b12d4405a55eea558e8f275aa6e.
//
// Solidity: event Register(uint256 indexed id, address indexed from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterRegister(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*LiquidityBridgeContractRegisterIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "Register", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractRegisterIterator{contract: _LiquidityBridgeContract.contract, event: "Register", logs: logs, sub: sub}, nil
}

// WatchRegister is a free log subscription operation binding the contract event 0xa9d44d6e13bb3fee938c3f66d1103e91f8dc6b12d4405a55eea558e8f275aa6e.
//
// Solidity: event Register(uint256 indexed id, address indexed from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchRegister(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractRegister) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "Register")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractRegister)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Register", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRegister is a log parse operation binding the contract event 0xa9d44d6e13bb3fee938c3f66d1103e91f8dc6b12d4405a55eea558e8f275aa6e.
//
// Solidity: event Register(uint256 indexed id, address indexed from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseRegister(log types.Log) (*LiquidityBridgeContractRegister, error) {
	event := new(LiquidityBridgeContractRegister)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Register", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractResignedIterator is returned from FilterResigned and is used to iterate over the raw logs and unpacked data for Resigned events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractResignedIterator struct {
	Event *LiquidityBridgeContractResigned // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractResignedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractResigned)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractResigned)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractResignedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractResignedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractResigned represents a Resigned event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractResigned struct {
	From common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterResigned is a free log retrieval operation binding the contract event 0xa6674aa33cd1b7435474751667707bf05fde99e537d67043ec5f907782577d86.
//
// Solidity: event Resigned(address from)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterResigned(opts *bind.FilterOpts) (*LiquidityBridgeContractResignedIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "Resigned")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractResignedIterator{contract: _LiquidityBridgeContract.contract, event: "Resigned", logs: logs, sub: sub}, nil
}

// WatchResigned is a free log subscription operation binding the contract event 0xa6674aa33cd1b7435474751667707bf05fde99e537d67043ec5f907782577d86.
//
// Solidity: event Resigned(address from)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchResigned(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractResigned) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "Resigned")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractResigned)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Resigned", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseResigned is a log parse operation binding the contract event 0xa6674aa33cd1b7435474751667707bf05fde99e537d67043ec5f907782577d86.
//
// Solidity: event Resigned(address from)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseResigned(log types.Log) (*LiquidityBridgeContractResigned, error) {
	event := new(LiquidityBridgeContractResigned)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Resigned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractWithdrawCollateralIterator is returned from FilterWithdrawCollateral and is used to iterate over the raw logs and unpacked data for WithdrawCollateral events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractWithdrawCollateralIterator struct {
	Event *LiquidityBridgeContractWithdrawCollateral // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractWithdrawCollateralIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractWithdrawCollateral)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractWithdrawCollateral)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractWithdrawCollateralIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractWithdrawCollateralIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractWithdrawCollateral represents a WithdrawCollateral event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractWithdrawCollateral struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawCollateral is a free log retrieval operation binding the contract event 0xa8e76b822fc682be77f3b1c822ea81f6bda5aed92ba82e6873bfd889f328d1d2.
//
// Solidity: event WithdrawCollateral(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterWithdrawCollateral(opts *bind.FilterOpts) (*LiquidityBridgeContractWithdrawCollateralIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "WithdrawCollateral")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractWithdrawCollateralIterator{contract: _LiquidityBridgeContract.contract, event: "WithdrawCollateral", logs: logs, sub: sub}, nil
}

// WatchWithdrawCollateral is a free log subscription operation binding the contract event 0xa8e76b822fc682be77f3b1c822ea81f6bda5aed92ba82e6873bfd889f328d1d2.
//
// Solidity: event WithdrawCollateral(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchWithdrawCollateral(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractWithdrawCollateral) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "WithdrawCollateral")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractWithdrawCollateral)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "WithdrawCollateral", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawCollateral is a log parse operation binding the contract event 0xa8e76b822fc682be77f3b1c822ea81f6bda5aed92ba82e6873bfd889f328d1d2.
//
// Solidity: event WithdrawCollateral(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseWithdrawCollateral(log types.Log) (*LiquidityBridgeContractWithdrawCollateral, error) {
	event := new(LiquidityBridgeContractWithdrawCollateral)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "WithdrawCollateral", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractWithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractWithdrawalIterator struct {
	Event *LiquidityBridgeContractWithdrawal // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidityBridgeContractWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractWithdrawal)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidityBridgeContractWithdrawal)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidityBridgeContractWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractWithdrawal represents a Withdrawal event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractWithdrawal struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterWithdrawal(opts *bind.FilterOpts) (*LiquidityBridgeContractWithdrawalIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "Withdrawal")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractWithdrawalIterator{contract: _LiquidityBridgeContract.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractWithdrawal) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "Withdrawal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractWithdrawal)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Withdrawal", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawal is a log parse operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseWithdrawal(log types.Log) (*LiquidityBridgeContractWithdrawal, error) {
	event := new(LiquidityBridgeContractWithdrawal)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
