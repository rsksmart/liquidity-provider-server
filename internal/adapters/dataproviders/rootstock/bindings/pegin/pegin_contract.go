// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// PegInContractPegInClaim is an auto generated low-level Go binding around an user-defined struct.
type PegInContractPegInClaim struct {
	Claimer      common.Address
	Amount       *big.Int
	Fee          *big.Int
	RequestBlock *big.Int
	Resolved     bool
}

// QuotesPegInQuote is an auto generated low-level Go binding around an user-defined struct.
type QuotesPegInQuote struct {
	ChainId                     *big.Int
	CallFee                     *big.Int
	PenaltyFee                  *big.Int
	Value                       *big.Int
	GasFee                      *big.Int
	FedBtcAddress               [20]byte
	LbcAddress                  common.Address
	LiquidityProviderRskAddress common.Address
	ContractAddress             common.Address
	RskRefundAddress            common.Address
	Nonce                       int64
	GasLimit                    uint32
	AgreementTimestamp          uint32
	TimeForDeposit              uint32
	CallTime                    uint32
	DepositConfirmations        uint16
	CallOnRegister              bool
	BtcRefundAddress            []byte
	LiquidityProviderBtcAddress []byte
	Data                        []byte
}

// PeginContractMetaData contains all meta data concerning the PeginContract contract.
var PeginContractMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"NAME\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"beginDefaultAdminTransfer\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"callForUser\",\"inputs\":[{\"name\":\"quote\",\"type\":\"tuple\",\"internalType\":\"structQuotes.PegInQuote\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"callFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fedBtcAddress\",\"type\":\"bytes20\",\"internalType\":\"bytes20\"},{\"name\":\"lbcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rskRefundAddress\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"nonce\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"agreementTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"timeForDeposit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"callTime\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"callOnRegister\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"btcRefundAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"cancelDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"changeDefaultAdminDelay\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"defaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelayIncreaseWait\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"dustThreshold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBalance\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinPegIn\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPegInClaim\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"btcTxHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPegInContract.PegInClaim\",\"components\":[{\"name\":\"claimer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requestBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"resolved\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuoteStatus\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIPegIn.PegInStates\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hashPegInQuote\",\"inputs\":[{\"name\":\"quote\",\"type\":\"tuple\",\"internalType\":\"structQuotes.PegInQuote\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"callFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fedBtcAddress\",\"type\":\"bytes20\",\"internalType\":\"bytes20\"},{\"name\":\"lbcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rskRefundAddress\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"nonce\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"agreementTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"timeForDeposit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"callTime\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"callOnRegister\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"btcRefundAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hashPegInQuoteEIP712\",\"inputs\":[{\"name\":\"quote\",\"type\":\"tuple\",\"internalType\":\"structQuotes.PegInQuote\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"callFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fedBtcAddress\",\"type\":\"bytes20\",\"internalType\":\"bytes20\"},{\"name\":\"lbcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rskRefundAddress\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"nonce\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"agreementTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"timeForDeposit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"callTime\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"callOnRegister\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"btcRefundAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"defaultAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bridge\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"dustThreshold_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minPegIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"collateralManagement\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"mainnet\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"pauseRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauseRegistry\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauseRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPauseRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauseStatus\",\"inputs\":[],\"outputs\":[{\"name\":\"isPaused\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"since\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pegInId\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"btcTxHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerPegIn\",\"inputs\":[{\"name\":\"quote\",\"type\":\"tuple\",\"internalType\":\"structQuotes.PegInQuote\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"callFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fedBtcAddress\",\"type\":\"bytes20\",\"internalType\":\"bytes20\"},{\"name\":\"lbcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rskRefundAddress\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"nonce\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"agreementTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"timeForDeposit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"callTime\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"callOnRegister\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"btcRefundAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"btcRawTransaction\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"partialMerkleTree\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"height\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestPegIn\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"btcTxHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"opReturn\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"callSuccess\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"resolvePegIn\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"btcTxHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"btcRawTransaction\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"partialMerkleTree\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"height\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"registrant\",\"type\":\"address\",\"internalType\":\"addresspayable\"}],\"outputs\":[{\"name\":\"bridgeResult\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rollbackDefaultAdminDelay\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setCollateralManagement\",\"inputs\":[{\"name\":\"collateralManagement\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDustThreshold\",\"inputs\":[{\"name\":\"threshold\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinPegIn\",\"inputs\":[{\"name\":\"minPegIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPegInClaimDependencies\",\"inputs\":[{\"name\":\"registry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"configurations\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"claimDeadlineBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"registrantFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slashUnclaimedPegIn\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"btcTxHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validatePegInDepositAddress\",\"inputs\":[{\"name\":\"quote\",\"type\":\"tuple\",\"internalType\":\"structQuotes.PegInQuote\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"callFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fedBtcAddress\",\"type\":\"bytes20\",\"internalType\":\"bytes20\"},{\"name\":\"lbcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rskRefundAddress\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"nonce\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"agreementTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"timeForDeposit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"callTime\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"callOnRegister\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"btcRefundAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"depositAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"BalanceDecrease\",\"inputs\":[{\"name\":\"dest\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BalanceIncrease\",\"inputs\":[{\"name\":\"dest\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BridgeCapExceeded\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"errorCode\",\"type\":\"int256\",\"indexed\":true,\"internalType\":\"int256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CallForUser\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"dest\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"success\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CollateralManagementSet\",\"inputs\":[{\"name\":\"oldAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeScheduled\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"effectSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferScheduled\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"acceptSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DustThresholdSet\",\"inputs\":[{\"name\":\"oldThreshold\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"newThreshold\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinPegInSet\",\"inputs\":[{\"name\":\"oldMinPegIn\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"newMinPegIn\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegInClaimDependenciesSet\",\"inputs\":[{\"name\":\"registry\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"configurations\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegInRegistered\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"transferredAmount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegInRequested\",\"inputs\":[{\"name\":\"pegInId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"claimer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"rskAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"netToUser\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"callSuccess\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegInResolved\",\"inputs\":[{\"name\":\"pegInId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"claimer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fronted\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Refund\",\"inputs\":[{\"name\":\"dest\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"success\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnclaimedPegInSlashed\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawal\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminDelay\",\"inputs\":[{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminRules\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlInvalidDefaultAdmin\",\"inputs\":[{\"name\":\"defaultAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressNotRegistered\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AmountUnderMinimum\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"BridgeSettlementFailed\",\"inputs\":[{\"name\":\"bridgeResult\",\"type\":\"int256\",\"internalType\":\"int256\"}]},{\"type\":\"error\",\"name\":\"ClaimDeadlineNotReached\",\"inputs\":[{\"name\":\"deadlineBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"DependencyNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyBlockHeader\",\"inputs\":[{\"name\":\"heightOrHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectContract\",\"inputs\":[{\"name\":\"expected\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"actual\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"IncorrectFronting\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"provided\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"IncorrectSignature\",\"inputs\":[{\"name\":\"expectedAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"usedHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InsufficientAmount\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"target\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientConfirmations\",\"inputs\":[{\"name\":\"have\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"required\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientGas\",\"inputs\":[{\"name\":\"gasLeft\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasRequired\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidChainId\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRefundAddress\",\"inputs\":[{\"name\":\"refundAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidSender\",\"inputs\":[{\"name\":\"expected\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"actual\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NoBalance\",\"inputs\":[{\"name\":\"wanted\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NoContract\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotEnoughConfirmations\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Overflow\",\"inputs\":[{\"name\":\"passedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"PaymentFailed\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"reason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"PaymentNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PegInAlreadyProcessed\",\"inputs\":[{\"name\":\"pegInId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"PegInNotClaimed\",\"inputs\":[{\"name\":\"pegInId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ProviderNotRegistered\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"QuoteAlreadyProcessed\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnexpectedBridgeError\",\"inputs\":[{\"name\":\"errorCode\",\"type\":\"int256\",\"internalType\":\"int256\"}]}]",
	ID:  "PeginContract",
}

// PeginContract is an auto generated Go binding around an Ethereum contract.
type PeginContract struct {
	abi abi.ABI
}

// NewPeginContract creates a new instance of PeginContract.
func NewPeginContract() *PeginContract {
	parsed, err := PeginContractMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &PeginContract{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *PeginContract) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (peginContract *PeginContract) PackDEFAULTADMINROLE() []byte {
	enc, err := peginContract.abi.Pack("DEFAULT_ADMIN_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (peginContract *PeginContract) TryPackDEFAULTADMINROLE() ([]byte, error) {
	return peginContract.abi.Pack("DEFAULT_ADMIN_ROLE")
}

// UnpackDEFAULTADMINROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (peginContract *PeginContract) UnpackDEFAULTADMINROLE(data []byte) ([32]byte, error) {
	out, err := peginContract.abi.Unpack("DEFAULT_ADMIN_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackNAME is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa3f4df7e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function NAME() view returns(string)
func (peginContract *PeginContract) PackNAME() []byte {
	enc, err := peginContract.abi.Pack("NAME")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackNAME is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa3f4df7e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function NAME() view returns(string)
func (peginContract *PeginContract) TryPackNAME() ([]byte, error) {
	return peginContract.abi.Pack("NAME")
}

// UnpackNAME is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa3f4df7e.
//
// Solidity: function NAME() view returns(string)
func (peginContract *PeginContract) UnpackNAME(data []byte) (string, error) {
	out, err := peginContract.abi.Unpack("NAME", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xffa1ad74.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function VERSION() view returns(string)
func (peginContract *PeginContract) PackVERSION() []byte {
	enc, err := peginContract.abi.Pack("VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xffa1ad74.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function VERSION() view returns(string)
func (peginContract *PeginContract) TryPackVERSION() ([]byte, error) {
	return peginContract.abi.Pack("VERSION")
}

// UnpackVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (peginContract *PeginContract) UnpackVERSION(data []byte) (string, error) {
	out, err := peginContract.abi.Unpack("VERSION", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackAcceptDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcefc1429.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (peginContract *PeginContract) PackAcceptDefaultAdminTransfer() []byte {
	enc, err := peginContract.abi.Pack("acceptDefaultAdminTransfer")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAcceptDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcefc1429.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (peginContract *PeginContract) TryPackAcceptDefaultAdminTransfer() ([]byte, error) {
	return peginContract.abi.Pack("acceptDefaultAdminTransfer")
}

// PackBeginDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x634e93da.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (peginContract *PeginContract) PackBeginDefaultAdminTransfer(newAdmin common.Address) []byte {
	enc, err := peginContract.abi.Pack("beginDefaultAdminTransfer", newAdmin)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBeginDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x634e93da.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (peginContract *PeginContract) TryPackBeginDefaultAdminTransfer(newAdmin common.Address) ([]byte, error) {
	return peginContract.abi.Pack("beginDefaultAdminTransfer", newAdmin)
}

// PackCallForUser is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc7a3dc3c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function callForUser((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) payable returns(bool)
func (peginContract *PeginContract) PackCallForUser(quote QuotesPegInQuote) []byte {
	enc, err := peginContract.abi.Pack("callForUser", quote)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCallForUser is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc7a3dc3c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function callForUser((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) payable returns(bool)
func (peginContract *PeginContract) TryPackCallForUser(quote QuotesPegInQuote) ([]byte, error) {
	return peginContract.abi.Pack("callForUser", quote)
}

// UnpackCallForUser is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc7a3dc3c.
//
// Solidity: function callForUser((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) payable returns(bool)
func (peginContract *PeginContract) UnpackCallForUser(data []byte) (bool, error) {
	out, err := peginContract.abi.Unpack("callForUser", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackCancelDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd602b9fd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (peginContract *PeginContract) PackCancelDefaultAdminTransfer() []byte {
	enc, err := peginContract.abi.Pack("cancelDefaultAdminTransfer")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCancelDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd602b9fd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (peginContract *PeginContract) TryPackCancelDefaultAdminTransfer() ([]byte, error) {
	return peginContract.abi.Pack("cancelDefaultAdminTransfer")
}

// PackChangeDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x649a5ec7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (peginContract *PeginContract) PackChangeDefaultAdminDelay(newDelay *big.Int) []byte {
	enc, err := peginContract.abi.Pack("changeDefaultAdminDelay", newDelay)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackChangeDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x649a5ec7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (peginContract *PeginContract) TryPackChangeDefaultAdminDelay(newDelay *big.Int) ([]byte, error) {
	return peginContract.abi.Pack("changeDefaultAdminDelay", newDelay)
}

// PackDefaultAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84ef8ffc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function defaultAdmin() view returns(address)
func (peginContract *PeginContract) PackDefaultAdmin() []byte {
	enc, err := peginContract.abi.Pack("defaultAdmin")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDefaultAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84ef8ffc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function defaultAdmin() view returns(address)
func (peginContract *PeginContract) TryPackDefaultAdmin() ([]byte, error) {
	return peginContract.abi.Pack("defaultAdmin")
}

// UnpackDefaultAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (peginContract *PeginContract) UnpackDefaultAdmin(data []byte) (common.Address, error) {
	out, err := peginContract.abi.Unpack("defaultAdmin", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcc8463c8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (peginContract *PeginContract) PackDefaultAdminDelay() []byte {
	enc, err := peginContract.abi.Pack("defaultAdminDelay")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcc8463c8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (peginContract *PeginContract) TryPackDefaultAdminDelay() ([]byte, error) {
	return peginContract.abi.Pack("defaultAdminDelay")
}

// UnpackDefaultAdminDelay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (peginContract *PeginContract) UnpackDefaultAdminDelay(data []byte) (*big.Int, error) {
	out, err := peginContract.abi.Unpack("defaultAdminDelay", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDefaultAdminDelayIncreaseWait is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x022d63fb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (peginContract *PeginContract) PackDefaultAdminDelayIncreaseWait() []byte {
	enc, err := peginContract.abi.Pack("defaultAdminDelayIncreaseWait")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDefaultAdminDelayIncreaseWait is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x022d63fb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (peginContract *PeginContract) TryPackDefaultAdminDelayIncreaseWait() ([]byte, error) {
	return peginContract.abi.Pack("defaultAdminDelayIncreaseWait")
}

// UnpackDefaultAdminDelayIncreaseWait is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (peginContract *PeginContract) UnpackDefaultAdminDelayIncreaseWait(data []byte) (*big.Int, error) {
	out, err := peginContract.abi.Unpack("defaultAdminDelayIncreaseWait", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd0e30db0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function deposit() payable returns()
func (peginContract *PeginContract) PackDeposit() []byte {
	enc, err := peginContract.abi.Pack("deposit")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd0e30db0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function deposit() payable returns()
func (peginContract *PeginContract) TryPackDeposit() ([]byte, error) {
	return peginContract.abi.Pack("deposit")
}

// PackDustThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe8462e8f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function dustThreshold() view returns(uint256)
func (peginContract *PeginContract) PackDustThreshold() []byte {
	enc, err := peginContract.abi.Pack("dustThreshold")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDustThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe8462e8f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function dustThreshold() view returns(uint256)
func (peginContract *PeginContract) TryPackDustThreshold() ([]byte, error) {
	return peginContract.abi.Pack("dustThreshold")
}

// UnpackDustThreshold is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe8462e8f.
//
// Solidity: function dustThreshold() view returns(uint256)
func (peginContract *PeginContract) UnpackDustThreshold(data []byte) (*big.Int, error) {
	out, err := peginContract.abi.Unpack("dustThreshold", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackEip712Domain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84b0196e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (peginContract *PeginContract) PackEip712Domain() []byte {
	enc, err := peginContract.abi.Pack("eip712Domain")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEip712Domain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84b0196e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (peginContract *PeginContract) TryPackEip712Domain() ([]byte, error) {
	return peginContract.abi.Pack("eip712Domain")
}

// Eip712DomainOutput serves as a container for the return parameters of contract
// method Eip712Domain.
type Eip712DomainOutput struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}

// UnpackEip712Domain is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (peginContract *PeginContract) UnpackEip712Domain(data []byte) (Eip712DomainOutput, error) {
	out, err := peginContract.abi.Unpack("eip712Domain", data)
	outstruct := new(Eip712DomainOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)
	return *outstruct, nil
}

// PackGetBalance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf8b2cb4f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (peginContract *PeginContract) PackGetBalance(addr common.Address) []byte {
	enc, err := peginContract.abi.Pack("getBalance", addr)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetBalance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf8b2cb4f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (peginContract *PeginContract) TryPackGetBalance(addr common.Address) ([]byte, error) {
	return peginContract.abi.Pack("getBalance", addr)
}

// UnpackGetBalance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf8b2cb4f.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (peginContract *PeginContract) UnpackGetBalance(data []byte) (*big.Int, error) {
	out, err := peginContract.abi.Unpack("getBalance", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetMinPegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfa88dcde.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getMinPegIn() view returns(uint256)
func (peginContract *PeginContract) PackGetMinPegIn() []byte {
	enc, err := peginContract.abi.Pack("getMinPegIn")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetMinPegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfa88dcde.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getMinPegIn() view returns(uint256)
func (peginContract *PeginContract) TryPackGetMinPegIn() ([]byte, error) {
	return peginContract.abi.Pack("getMinPegIn")
}

// UnpackGetMinPegIn is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfa88dcde.
//
// Solidity: function getMinPegIn() view returns(uint256)
func (peginContract *PeginContract) UnpackGetMinPegIn(data []byte) (*big.Int, error) {
	out, err := peginContract.abi.Unpack("getMinPegIn", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetPegInClaim is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa3146dfb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPegInClaim(address rskAddr, bytes32 btcTxHash) view returns((address,uint256,uint256,uint256,bool))
func (peginContract *PeginContract) PackGetPegInClaim(rskAddr common.Address, btcTxHash [32]byte) []byte {
	enc, err := peginContract.abi.Pack("getPegInClaim", rskAddr, btcTxHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPegInClaim is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa3146dfb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPegInClaim(address rskAddr, bytes32 btcTxHash) view returns((address,uint256,uint256,uint256,bool))
func (peginContract *PeginContract) TryPackGetPegInClaim(rskAddr common.Address, btcTxHash [32]byte) ([]byte, error) {
	return peginContract.abi.Pack("getPegInClaim", rskAddr, btcTxHash)
}

// UnpackGetPegInClaim is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa3146dfb.
//
// Solidity: function getPegInClaim(address rskAddr, bytes32 btcTxHash) view returns((address,uint256,uint256,uint256,bool))
func (peginContract *PeginContract) UnpackGetPegInClaim(data []byte) (PegInContractPegInClaim, error) {
	out, err := peginContract.abi.Unpack("getPegInClaim", data)
	if err != nil {
		return *new(PegInContractPegInClaim), err
	}
	out0 := *abi.ConvertType(out[0], new(PegInContractPegInClaim)).(*PegInContractPegInClaim)
	return out0, nil
}

// PackGetQuoteStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf93c8ec2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getQuoteStatus(bytes32 quoteHash) view returns(uint8)
func (peginContract *PeginContract) PackGetQuoteStatus(quoteHash [32]byte) []byte {
	enc, err := peginContract.abi.Pack("getQuoteStatus", quoteHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetQuoteStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf93c8ec2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getQuoteStatus(bytes32 quoteHash) view returns(uint8)
func (peginContract *PeginContract) TryPackGetQuoteStatus(quoteHash [32]byte) ([]byte, error) {
	return peginContract.abi.Pack("getQuoteStatus", quoteHash)
}

// UnpackGetQuoteStatus is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf93c8ec2.
//
// Solidity: function getQuoteStatus(bytes32 quoteHash) view returns(uint8)
func (peginContract *PeginContract) UnpackGetQuoteStatus(data []byte) (uint8, error) {
	out, err := peginContract.abi.Unpack("getQuoteStatus", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (peginContract *PeginContract) PackGetRoleAdmin(role [32]byte) []byte {
	enc, err := peginContract.abi.Pack("getRoleAdmin", role)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (peginContract *PeginContract) TryPackGetRoleAdmin(role [32]byte) ([]byte, error) {
	return peginContract.abi.Pack("getRoleAdmin", role)
}

// UnpackGetRoleAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (peginContract *PeginContract) UnpackGetRoleAdmin(data []byte) ([32]byte, error) {
	out, err := peginContract.abi.Unpack("getRoleAdmin", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (peginContract *PeginContract) PackGrantRole(role [32]byte, account common.Address) []byte {
	enc, err := peginContract.abi.Pack("grantRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (peginContract *PeginContract) TryPackGrantRole(role [32]byte, account common.Address) ([]byte, error) {
	return peginContract.abi.Pack("grantRole", role, account)
}

// PackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (peginContract *PeginContract) PackHasRole(role [32]byte, account common.Address) []byte {
	enc, err := peginContract.abi.Pack("hasRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (peginContract *PeginContract) TryPackHasRole(role [32]byte, account common.Address) ([]byte, error) {
	return peginContract.abi.Pack("hasRole", role, account)
}

// UnpackHasRole is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (peginContract *PeginContract) UnpackHasRole(data []byte) (bool, error) {
	out, err := peginContract.abi.Unpack("hasRole", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackHashPegInQuote is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf218a7d8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hashPegInQuote((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) view returns(bytes32)
func (peginContract *PeginContract) PackHashPegInQuote(quote QuotesPegInQuote) []byte {
	enc, err := peginContract.abi.Pack("hashPegInQuote", quote)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHashPegInQuote is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf218a7d8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hashPegInQuote((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) view returns(bytes32)
func (peginContract *PeginContract) TryPackHashPegInQuote(quote QuotesPegInQuote) ([]byte, error) {
	return peginContract.abi.Pack("hashPegInQuote", quote)
}

// UnpackHashPegInQuote is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf218a7d8.
//
// Solidity: function hashPegInQuote((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) view returns(bytes32)
func (peginContract *PeginContract) UnpackHashPegInQuote(data []byte) ([32]byte, error) {
	out, err := peginContract.abi.Unpack("hashPegInQuote", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackHashPegInQuoteEIP712 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x928f4598.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hashPegInQuoteEIP712((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) view returns(bytes32)
func (peginContract *PeginContract) PackHashPegInQuoteEIP712(quote QuotesPegInQuote) []byte {
	enc, err := peginContract.abi.Pack("hashPegInQuoteEIP712", quote)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHashPegInQuoteEIP712 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x928f4598.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hashPegInQuoteEIP712((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) view returns(bytes32)
func (peginContract *PeginContract) TryPackHashPegInQuoteEIP712(quote QuotesPegInQuote) ([]byte, error) {
	return peginContract.abi.Pack("hashPegInQuoteEIP712", quote)
}

// UnpackHashPegInQuoteEIP712 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x928f4598.
//
// Solidity: function hashPegInQuoteEIP712((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) view returns(bytes32)
func (peginContract *PeginContract) UnpackHashPegInQuoteEIP712(data []byte) ([32]byte, error) {
	out, err := peginContract.abi.Unpack("hashPegInQuoteEIP712", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe19b6563.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function initialize(address defaultAdmin, address bridge, uint256 dustThreshold_, uint256 minPegIn, address collateralManagement, bool mainnet, address pauseRegistry) returns()
func (peginContract *PeginContract) PackInitialize(defaultAdmin common.Address, bridge common.Address, dustThreshold *big.Int, minPegIn *big.Int, collateralManagement common.Address, mainnet bool, pauseRegistry common.Address) []byte {
	enc, err := peginContract.abi.Pack("initialize", defaultAdmin, bridge, dustThreshold, minPegIn, collateralManagement, mainnet, pauseRegistry)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe19b6563.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function initialize(address defaultAdmin, address bridge, uint256 dustThreshold_, uint256 minPegIn, address collateralManagement, bool mainnet, address pauseRegistry) returns()
func (peginContract *PeginContract) TryPackInitialize(defaultAdmin common.Address, bridge common.Address, dustThreshold *big.Int, minPegIn *big.Int, collateralManagement common.Address, mainnet bool, pauseRegistry common.Address) ([]byte, error) {
	return peginContract.abi.Pack("initialize", defaultAdmin, bridge, dustThreshold, minPegIn, collateralManagement, mainnet, pauseRegistry)
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (peginContract *PeginContract) PackOwner() []byte {
	enc, err := peginContract.abi.Pack("owner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function owner() view returns(address)
func (peginContract *PeginContract) TryPackOwner() ([]byte, error) {
	return peginContract.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (peginContract *PeginContract) UnpackOwner(data []byte) (common.Address, error) {
	out, err := peginContract.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPauseRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1c82732d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pauseRegistry() view returns(address)
func (peginContract *PeginContract) PackPauseRegistry() []byte {
	enc, err := peginContract.abi.Pack("pauseRegistry")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPauseRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1c82732d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pauseRegistry() view returns(address)
func (peginContract *PeginContract) TryPackPauseRegistry() ([]byte, error) {
	return peginContract.abi.Pack("pauseRegistry")
}

// UnpackPauseRegistry is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1c82732d.
//
// Solidity: function pauseRegistry() view returns(address)
func (peginContract *PeginContract) UnpackPauseRegistry(data []byte) (common.Address, error) {
	out, err := peginContract.abi.Unpack("pauseRegistry", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPauseStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x466916ca.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (peginContract *PeginContract) PackPauseStatus() []byte {
	enc, err := peginContract.abi.Pack("pauseStatus")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPauseStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x466916ca.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (peginContract *PeginContract) TryPackPauseStatus() ([]byte, error) {
	return peginContract.abi.Pack("pauseStatus")
}

// PauseStatusOutput serves as a container for the return parameters of contract
// method PauseStatus.
type PauseStatusOutput struct {
	IsPaused bool
	Reason   string
	Since    uint64
}

// UnpackPauseStatus is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x466916ca.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (peginContract *PeginContract) UnpackPauseStatus(data []byte) (PauseStatusOutput, error) {
	out, err := peginContract.abi.Unpack("pauseStatus", data)
	outstruct := new(PauseStatusOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.IsPaused = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Reason = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Since = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	return *outstruct, nil
}

// PackPegInId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x441064de.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pegInId(address rskAddr, bytes32 btcTxHash) pure returns(bytes32)
func (peginContract *PeginContract) PackPegInId(rskAddr common.Address, btcTxHash [32]byte) []byte {
	enc, err := peginContract.abi.Pack("pegInId", rskAddr, btcTxHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPegInId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x441064de.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pegInId(address rskAddr, bytes32 btcTxHash) pure returns(bytes32)
func (peginContract *PeginContract) TryPackPegInId(rskAddr common.Address, btcTxHash [32]byte) ([]byte, error) {
	return peginContract.abi.Pack("pegInId", rskAddr, btcTxHash)
}

// UnpackPegInId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x441064de.
//
// Solidity: function pegInId(address rskAddr, bytes32 btcTxHash) pure returns(bytes32)
func (peginContract *PeginContract) UnpackPegInId(data []byte) ([32]byte, error) {
	out, err := peginContract.abi.Unpack("pegInId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackPendingDefaultAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf6eefb7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (peginContract *PeginContract) PackPendingDefaultAdmin() []byte {
	enc, err := peginContract.abi.Pack("pendingDefaultAdmin")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPendingDefaultAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf6eefb7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (peginContract *PeginContract) TryPackPendingDefaultAdmin() ([]byte, error) {
	return peginContract.abi.Pack("pendingDefaultAdmin")
}

// PendingDefaultAdminOutput serves as a container for the return parameters of contract
// method PendingDefaultAdmin.
type PendingDefaultAdminOutput struct {
	NewAdmin common.Address
	Schedule *big.Int
}

// UnpackPendingDefaultAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcf6eefb7.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (peginContract *PeginContract) UnpackPendingDefaultAdmin(data []byte) (PendingDefaultAdminOutput, error) {
	out, err := peginContract.abi.Unpack("pendingDefaultAdmin", data)
	outstruct := new(PendingDefaultAdminOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.NewAdmin = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Schedule = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackPendingDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa1eda53c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (peginContract *PeginContract) PackPendingDefaultAdminDelay() []byte {
	enc, err := peginContract.abi.Pack("pendingDefaultAdminDelay")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPendingDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa1eda53c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (peginContract *PeginContract) TryPackPendingDefaultAdminDelay() ([]byte, error) {
	return peginContract.abi.Pack("pendingDefaultAdminDelay")
}

// PendingDefaultAdminDelayOutput serves as a container for the return parameters of contract
// method PendingDefaultAdminDelay.
type PendingDefaultAdminDelayOutput struct {
	NewDelay *big.Int
	Schedule *big.Int
}

// UnpackPendingDefaultAdminDelay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa1eda53c.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (peginContract *PeginContract) UnpackPendingDefaultAdminDelay(data []byte) (PendingDefaultAdminDelayOutput, error) {
	out, err := peginContract.abi.Unpack("pendingDefaultAdminDelay", data)
	outstruct := new(PendingDefaultAdminDelayOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.NewDelay = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Schedule = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackRegisterPegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3823c753.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function registerPegIn((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes signature, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256)
func (peginContract *PeginContract) PackRegisterPegIn(quote QuotesPegInQuote, signature []byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) []byte {
	enc, err := peginContract.abi.Pack("registerPegIn", quote, signature, btcRawTransaction, partialMerkleTree, height)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRegisterPegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3823c753.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function registerPegIn((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes signature, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256)
func (peginContract *PeginContract) TryPackRegisterPegIn(quote QuotesPegInQuote, signature []byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) ([]byte, error) {
	return peginContract.abi.Pack("registerPegIn", quote, signature, btcRawTransaction, partialMerkleTree, height)
}

// UnpackRegisterPegIn is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3823c753.
//
// Solidity: function registerPegIn((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes signature, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256)
func (peginContract *PeginContract) UnpackRegisterPegIn(data []byte) (*big.Int, error) {
	out, err := peginContract.abi.Unpack("registerPegIn", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (peginContract *PeginContract) PackRenounceRole(role [32]byte, account common.Address) []byte {
	enc, err := peginContract.abi.Pack("renounceRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (peginContract *PeginContract) TryPackRenounceRole(role [32]byte, account common.Address) ([]byte, error) {
	return peginContract.abi.Pack("renounceRole", role, account)
}

// PackRequestPegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x53cf6888.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function requestPegIn(address rskAddr, uint256 amount, bytes32 btcTxHash, bytes opReturn) payable returns(bool callSuccess)
func (peginContract *PeginContract) PackRequestPegIn(rskAddr common.Address, amount *big.Int, btcTxHash [32]byte, opReturn []byte) []byte {
	enc, err := peginContract.abi.Pack("requestPegIn", rskAddr, amount, btcTxHash, opReturn)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRequestPegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x53cf6888.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function requestPegIn(address rskAddr, uint256 amount, bytes32 btcTxHash, bytes opReturn) payable returns(bool callSuccess)
func (peginContract *PeginContract) TryPackRequestPegIn(rskAddr common.Address, amount *big.Int, btcTxHash [32]byte, opReturn []byte) ([]byte, error) {
	return peginContract.abi.Pack("requestPegIn", rskAddr, amount, btcTxHash, opReturn)
}

// UnpackRequestPegIn is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x53cf6888.
//
// Solidity: function requestPegIn(address rskAddr, uint256 amount, bytes32 btcTxHash, bytes opReturn) payable returns(bool callSuccess)
func (peginContract *PeginContract) UnpackRequestPegIn(data []byte) (bool, error) {
	out, err := peginContract.abi.Unpack("requestPegIn", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackResolvePegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x68248b5a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function resolvePegIn(address rskAddr, bytes32 btcTxHash, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height, address registrant) returns(int256 bridgeResult)
func (peginContract *PeginContract) PackResolvePegIn(rskAddr common.Address, btcTxHash [32]byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int, registrant common.Address) []byte {
	enc, err := peginContract.abi.Pack("resolvePegIn", rskAddr, btcTxHash, btcRawTransaction, partialMerkleTree, height, registrant)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackResolvePegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x68248b5a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function resolvePegIn(address rskAddr, bytes32 btcTxHash, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height, address registrant) returns(int256 bridgeResult)
func (peginContract *PeginContract) TryPackResolvePegIn(rskAddr common.Address, btcTxHash [32]byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int, registrant common.Address) ([]byte, error) {
	return peginContract.abi.Pack("resolvePegIn", rskAddr, btcTxHash, btcRawTransaction, partialMerkleTree, height, registrant)
}

// UnpackResolvePegIn is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x68248b5a.
//
// Solidity: function resolvePegIn(address rskAddr, bytes32 btcTxHash, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height, address registrant) returns(int256 bridgeResult)
func (peginContract *PeginContract) UnpackResolvePegIn(data []byte) (*big.Int, error) {
	out, err := peginContract.abi.Unpack("resolvePegIn", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (peginContract *PeginContract) PackRevokeRole(role [32]byte, account common.Address) []byte {
	enc, err := peginContract.abi.Pack("revokeRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (peginContract *PeginContract) TryPackRevokeRole(role [32]byte, account common.Address) ([]byte, error) {
	return peginContract.abi.Pack("revokeRole", role, account)
}

// PackRollbackDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0aa6220b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (peginContract *PeginContract) PackRollbackDefaultAdminDelay() []byte {
	enc, err := peginContract.abi.Pack("rollbackDefaultAdminDelay")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRollbackDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0aa6220b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (peginContract *PeginContract) TryPackRollbackDefaultAdminDelay() ([]byte, error) {
	return peginContract.abi.Pack("rollbackDefaultAdminDelay")
}

// PackSetCollateralManagement is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x05dfbf62.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setCollateralManagement(address collateralManagement) returns()
func (peginContract *PeginContract) PackSetCollateralManagement(collateralManagement common.Address) []byte {
	enc, err := peginContract.abi.Pack("setCollateralManagement", collateralManagement)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetCollateralManagement is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x05dfbf62.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setCollateralManagement(address collateralManagement) returns()
func (peginContract *PeginContract) TryPackSetCollateralManagement(collateralManagement common.Address) ([]byte, error) {
	return peginContract.abi.Pack("setCollateralManagement", collateralManagement)
}

// PackSetDustThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad7e55ba.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setDustThreshold(uint256 threshold) returns()
func (peginContract *PeginContract) PackSetDustThreshold(threshold *big.Int) []byte {
	enc, err := peginContract.abi.Pack("setDustThreshold", threshold)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetDustThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad7e55ba.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setDustThreshold(uint256 threshold) returns()
func (peginContract *PeginContract) TryPackSetDustThreshold(threshold *big.Int) ([]byte, error) {
	return peginContract.abi.Pack("setDustThreshold", threshold)
}

// PackSetMinPegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe4eae927.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setMinPegIn(uint256 minPegIn) returns()
func (peginContract *PeginContract) PackSetMinPegIn(minPegIn *big.Int) []byte {
	enc, err := peginContract.abi.Pack("setMinPegIn", minPegIn)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetMinPegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe4eae927.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setMinPegIn(uint256 minPegIn) returns()
func (peginContract *PeginContract) TryPackSetMinPegIn(minPegIn *big.Int) ([]byte, error) {
	return peginContract.abi.Pack("setMinPegIn", minPegIn)
}

// PackSetPegInClaimDependencies is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc66cd39e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setPegInClaimDependencies(address registry, address configurations, uint256 claimDeadlineBlocks, uint256 registrantFee) returns()
func (peginContract *PeginContract) PackSetPegInClaimDependencies(registry common.Address, configurations common.Address, claimDeadlineBlocks *big.Int, registrantFee *big.Int) []byte {
	enc, err := peginContract.abi.Pack("setPegInClaimDependencies", registry, configurations, claimDeadlineBlocks, registrantFee)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetPegInClaimDependencies is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc66cd39e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setPegInClaimDependencies(address registry, address configurations, uint256 claimDeadlineBlocks, uint256 registrantFee) returns()
func (peginContract *PeginContract) TryPackSetPegInClaimDependencies(registry common.Address, configurations common.Address, claimDeadlineBlocks *big.Int, registrantFee *big.Int) ([]byte, error) {
	return peginContract.abi.Pack("setPegInClaimDependencies", registry, configurations, claimDeadlineBlocks, registrantFee)
}

// PackSlashUnclaimedPegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x988a0156.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function slashUnclaimedPegIn(address rskAddr, uint256 amount, bytes32 btcTxHash) returns()
func (peginContract *PeginContract) PackSlashUnclaimedPegIn(rskAddr common.Address, amount *big.Int, btcTxHash [32]byte) []byte {
	enc, err := peginContract.abi.Pack("slashUnclaimedPegIn", rskAddr, amount, btcTxHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSlashUnclaimedPegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x988a0156.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function slashUnclaimedPegIn(address rskAddr, uint256 amount, bytes32 btcTxHash) returns()
func (peginContract *PeginContract) TryPackSlashUnclaimedPegIn(rskAddr common.Address, amount *big.Int, btcTxHash [32]byte) ([]byte, error) {
	return peginContract.abi.Pack("slashUnclaimedPegIn", rskAddr, amount, btcTxHash)
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (peginContract *PeginContract) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := peginContract.abi.Pack("supportsInterface", interfaceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (peginContract *PeginContract) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return peginContract.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (peginContract *PeginContract) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := peginContract.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackValidatePegInDepositAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe9accea2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function validatePegInDepositAddress((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes depositAddress) view returns(bool)
func (peginContract *PeginContract) PackValidatePegInDepositAddress(quote QuotesPegInQuote, depositAddress []byte) []byte {
	enc, err := peginContract.abi.Pack("validatePegInDepositAddress", quote, depositAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackValidatePegInDepositAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe9accea2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function validatePegInDepositAddress((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes depositAddress) view returns(bool)
func (peginContract *PeginContract) TryPackValidatePegInDepositAddress(quote QuotesPegInQuote, depositAddress []byte) ([]byte, error) {
	return peginContract.abi.Pack("validatePegInDepositAddress", quote, depositAddress)
}

// UnpackValidatePegInDepositAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe9accea2.
//
// Solidity: function validatePegInDepositAddress((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes depositAddress) view returns(bool)
func (peginContract *PeginContract) UnpackValidatePegInDepositAddress(data []byte) (bool, error) {
	out, err := peginContract.abi.Unpack("validatePegInDepositAddress", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2e1a7d4d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdraw(uint256 amount) returns()
func (peginContract *PeginContract) PackWithdraw(amount *big.Int) []byte {
	enc, err := peginContract.abi.Pack("withdraw", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2e1a7d4d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdraw(uint256 amount) returns()
func (peginContract *PeginContract) TryPackWithdraw(amount *big.Int) ([]byte, error) {
	return peginContract.abi.Pack("withdraw", amount)
}

// PeginContractBalanceDecrease represents a BalanceDecrease event raised by the PeginContract contract.
type PeginContractBalanceDecrease struct {
	Dest   common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const PeginContractBalanceDecreaseEventName = "BalanceDecrease"

// ContractEventName returns the user-defined event name.
func (PeginContractBalanceDecrease) ContractEventName() string {
	return PeginContractBalanceDecreaseEventName
}

// UnpackBalanceDecreaseEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BalanceDecrease(address indexed dest, uint256 indexed amount)
func (peginContract *PeginContract) UnpackBalanceDecreaseEvent(log *types.Log) (*PeginContractBalanceDecrease, error) {
	event := "BalanceDecrease"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractBalanceDecrease)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractBalanceIncrease represents a BalanceIncrease event raised by the PeginContract contract.
type PeginContractBalanceIncrease struct {
	Dest   common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const PeginContractBalanceIncreaseEventName = "BalanceIncrease"

// ContractEventName returns the user-defined event name.
func (PeginContractBalanceIncrease) ContractEventName() string {
	return PeginContractBalanceIncreaseEventName
}

// UnpackBalanceIncreaseEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BalanceIncrease(address indexed dest, uint256 indexed amount)
func (peginContract *PeginContract) UnpackBalanceIncreaseEvent(log *types.Log) (*PeginContractBalanceIncrease, error) {
	event := "BalanceIncrease"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractBalanceIncrease)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractBridgeCapExceeded represents a BridgeCapExceeded event raised by the PeginContract contract.
type PeginContractBridgeCapExceeded struct {
	QuoteHash [32]byte
	ErrorCode *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const PeginContractBridgeCapExceededEventName = "BridgeCapExceeded"

// ContractEventName returns the user-defined event name.
func (PeginContractBridgeCapExceeded) ContractEventName() string {
	return PeginContractBridgeCapExceededEventName
}

// UnpackBridgeCapExceededEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BridgeCapExceeded(bytes32 indexed quoteHash, int256 indexed errorCode)
func (peginContract *PeginContract) UnpackBridgeCapExceededEvent(log *types.Log) (*PeginContractBridgeCapExceeded, error) {
	event := "BridgeCapExceeded"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractBridgeCapExceeded)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractCallForUser represents a CallForUser event raised by the PeginContract contract.
type PeginContractCallForUser struct {
	From      common.Address
	Dest      common.Address
	QuoteHash [32]byte
	GasLimit  *big.Int
	Value     *big.Int
	Data      []byte
	Success   bool
	Raw       *types.Log // Blockchain specific contextual infos
}

const PeginContractCallForUserEventName = "CallForUser"

// ContractEventName returns the user-defined event name.
func (PeginContractCallForUser) ContractEventName() string {
	return PeginContractCallForUserEventName
}

// UnpackCallForUserEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CallForUser(address indexed from, address indexed dest, bytes32 indexed quoteHash, uint256 gasLimit, uint256 value, bytes data, bool success)
func (peginContract *PeginContract) UnpackCallForUserEvent(log *types.Log) (*PeginContractCallForUser, error) {
	event := "CallForUser"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractCallForUser)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractCollateralManagementSet represents a CollateralManagementSet event raised by the PeginContract contract.
type PeginContractCollateralManagementSet struct {
	OldAddress common.Address
	NewAddress common.Address
	Raw        *types.Log // Blockchain specific contextual infos
}

const PeginContractCollateralManagementSetEventName = "CollateralManagementSet"

// ContractEventName returns the user-defined event name.
func (PeginContractCollateralManagementSet) ContractEventName() string {
	return PeginContractCollateralManagementSetEventName
}

// UnpackCollateralManagementSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CollateralManagementSet(address indexed oldAddress, address indexed newAddress)
func (peginContract *PeginContract) UnpackCollateralManagementSetEvent(log *types.Log) (*PeginContractCollateralManagementSet, error) {
	event := "CollateralManagementSet"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractCollateralManagementSet)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractDefaultAdminDelayChangeCanceled represents a DefaultAdminDelayChangeCanceled event raised by the PeginContract contract.
type PeginContractDefaultAdminDelayChangeCanceled struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const PeginContractDefaultAdminDelayChangeCanceledEventName = "DefaultAdminDelayChangeCanceled"

// ContractEventName returns the user-defined event name.
func (PeginContractDefaultAdminDelayChangeCanceled) ContractEventName() string {
	return PeginContractDefaultAdminDelayChangeCanceledEventName
}

// UnpackDefaultAdminDelayChangeCanceledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (peginContract *PeginContract) UnpackDefaultAdminDelayChangeCanceledEvent(log *types.Log) (*PeginContractDefaultAdminDelayChangeCanceled, error) {
	event := "DefaultAdminDelayChangeCanceled"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractDefaultAdminDelayChangeCanceled)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractDefaultAdminDelayChangeScheduled represents a DefaultAdminDelayChangeScheduled event raised by the PeginContract contract.
type PeginContractDefaultAdminDelayChangeScheduled struct {
	NewDelay       *big.Int
	EffectSchedule *big.Int
	Raw            *types.Log // Blockchain specific contextual infos
}

const PeginContractDefaultAdminDelayChangeScheduledEventName = "DefaultAdminDelayChangeScheduled"

// ContractEventName returns the user-defined event name.
func (PeginContractDefaultAdminDelayChangeScheduled) ContractEventName() string {
	return PeginContractDefaultAdminDelayChangeScheduledEventName
}

// UnpackDefaultAdminDelayChangeScheduledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (peginContract *PeginContract) UnpackDefaultAdminDelayChangeScheduledEvent(log *types.Log) (*PeginContractDefaultAdminDelayChangeScheduled, error) {
	event := "DefaultAdminDelayChangeScheduled"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractDefaultAdminDelayChangeScheduled)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractDefaultAdminTransferCanceled represents a DefaultAdminTransferCanceled event raised by the PeginContract contract.
type PeginContractDefaultAdminTransferCanceled struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const PeginContractDefaultAdminTransferCanceledEventName = "DefaultAdminTransferCanceled"

// ContractEventName returns the user-defined event name.
func (PeginContractDefaultAdminTransferCanceled) ContractEventName() string {
	return PeginContractDefaultAdminTransferCanceledEventName
}

// UnpackDefaultAdminTransferCanceledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminTransferCanceled()
func (peginContract *PeginContract) UnpackDefaultAdminTransferCanceledEvent(log *types.Log) (*PeginContractDefaultAdminTransferCanceled, error) {
	event := "DefaultAdminTransferCanceled"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractDefaultAdminTransferCanceled)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractDefaultAdminTransferScheduled represents a DefaultAdminTransferScheduled event raised by the PeginContract contract.
type PeginContractDefaultAdminTransferScheduled struct {
	NewAdmin       common.Address
	AcceptSchedule *big.Int
	Raw            *types.Log // Blockchain specific contextual infos
}

const PeginContractDefaultAdminTransferScheduledEventName = "DefaultAdminTransferScheduled"

// ContractEventName returns the user-defined event name.
func (PeginContractDefaultAdminTransferScheduled) ContractEventName() string {
	return PeginContractDefaultAdminTransferScheduledEventName
}

// UnpackDefaultAdminTransferScheduledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (peginContract *PeginContract) UnpackDefaultAdminTransferScheduledEvent(log *types.Log) (*PeginContractDefaultAdminTransferScheduled, error) {
	event := "DefaultAdminTransferScheduled"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractDefaultAdminTransferScheduled)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractDustThresholdSet represents a DustThresholdSet event raised by the PeginContract contract.
type PeginContractDustThresholdSet struct {
	OldThreshold *big.Int
	NewThreshold *big.Int
	Raw          *types.Log // Blockchain specific contextual infos
}

const PeginContractDustThresholdSetEventName = "DustThresholdSet"

// ContractEventName returns the user-defined event name.
func (PeginContractDustThresholdSet) ContractEventName() string {
	return PeginContractDustThresholdSetEventName
}

// UnpackDustThresholdSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DustThresholdSet(uint256 indexed oldThreshold, uint256 indexed newThreshold)
func (peginContract *PeginContract) UnpackDustThresholdSetEvent(log *types.Log) (*PeginContractDustThresholdSet, error) {
	event := "DustThresholdSet"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractDustThresholdSet)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractEIP712DomainChanged represents a EIP712DomainChanged event raised by the PeginContract contract.
type PeginContractEIP712DomainChanged struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const PeginContractEIP712DomainChangedEventName = "EIP712DomainChanged"

// ContractEventName returns the user-defined event name.
func (PeginContractEIP712DomainChanged) ContractEventName() string {
	return PeginContractEIP712DomainChangedEventName
}

// UnpackEIP712DomainChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EIP712DomainChanged()
func (peginContract *PeginContract) UnpackEIP712DomainChangedEvent(log *types.Log) (*PeginContractEIP712DomainChanged, error) {
	event := "EIP712DomainChanged"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractEIP712DomainChanged)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractInitialized represents a Initialized event raised by the PeginContract contract.
type PeginContractInitialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const PeginContractInitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (PeginContractInitialized) ContractEventName() string {
	return PeginContractInitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (peginContract *PeginContract) UnpackInitializedEvent(log *types.Log) (*PeginContractInitialized, error) {
	event := "Initialized"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractInitialized)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractMinPegInSet represents a MinPegInSet event raised by the PeginContract contract.
type PeginContractMinPegInSet struct {
	OldMinPegIn *big.Int
	NewMinPegIn *big.Int
	Raw         *types.Log // Blockchain specific contextual infos
}

const PeginContractMinPegInSetEventName = "MinPegInSet"

// ContractEventName returns the user-defined event name.
func (PeginContractMinPegInSet) ContractEventName() string {
	return PeginContractMinPegInSetEventName
}

// UnpackMinPegInSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event MinPegInSet(uint256 indexed oldMinPegIn, uint256 indexed newMinPegIn)
func (peginContract *PeginContract) UnpackMinPegInSetEvent(log *types.Log) (*PeginContractMinPegInSet, error) {
	event := "MinPegInSet"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractMinPegInSet)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractPegInClaimDependenciesSet represents a PegInClaimDependenciesSet event raised by the PeginContract contract.
type PeginContractPegInClaimDependenciesSet struct {
	Registry       common.Address
	Configurations common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const PeginContractPegInClaimDependenciesSetEventName = "PegInClaimDependenciesSet"

// ContractEventName returns the user-defined event name.
func (PeginContractPegInClaimDependenciesSet) ContractEventName() string {
	return PeginContractPegInClaimDependenciesSetEventName
}

// UnpackPegInClaimDependenciesSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegInClaimDependenciesSet(address indexed registry, address indexed configurations)
func (peginContract *PeginContract) UnpackPegInClaimDependenciesSetEvent(log *types.Log) (*PeginContractPegInClaimDependenciesSet, error) {
	event := "PegInClaimDependenciesSet"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractPegInClaimDependenciesSet)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractPegInRegistered represents a PegInRegistered event raised by the PeginContract contract.
type PeginContractPegInRegistered struct {
	QuoteHash         [32]byte
	TransferredAmount *big.Int
	Raw               *types.Log // Blockchain specific contextual infos
}

const PeginContractPegInRegisteredEventName = "PegInRegistered"

// ContractEventName returns the user-defined event name.
func (PeginContractPegInRegistered) ContractEventName() string {
	return PeginContractPegInRegisteredEventName
}

// UnpackPegInRegisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegInRegistered(bytes32 indexed quoteHash, uint256 indexed transferredAmount)
func (peginContract *PeginContract) UnpackPegInRegisteredEvent(log *types.Log) (*PeginContractPegInRegistered, error) {
	event := "PegInRegistered"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractPegInRegistered)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractPegInRequested represents a PegInRequested event raised by the PeginContract contract.
type PeginContractPegInRequested struct {
	PegInId     [32]byte
	Claimer     common.Address
	RskAddr     common.Address
	Amount      *big.Int
	NetToUser   *big.Int
	CallSuccess bool
	Raw         *types.Log // Blockchain specific contextual infos
}

const PeginContractPegInRequestedEventName = "PegInRequested"

// ContractEventName returns the user-defined event name.
func (PeginContractPegInRequested) ContractEventName() string {
	return PeginContractPegInRequestedEventName
}

// UnpackPegInRequestedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegInRequested(bytes32 indexed pegInId, address indexed claimer, address indexed rskAddr, uint256 amount, uint256 netToUser, bool callSuccess)
func (peginContract *PeginContract) UnpackPegInRequestedEvent(log *types.Log) (*PeginContractPegInRequested, error) {
	event := "PegInRequested"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractPegInRequested)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractPegInResolved represents a PegInResolved event raised by the PeginContract contract.
type PeginContractPegInResolved struct {
	PegInId [32]byte
	Claimer common.Address
	Fronted *big.Int
	Fee     *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const PeginContractPegInResolvedEventName = "PegInResolved"

// ContractEventName returns the user-defined event name.
func (PeginContractPegInResolved) ContractEventName() string {
	return PeginContractPegInResolvedEventName
}

// UnpackPegInResolvedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegInResolved(bytes32 indexed pegInId, address indexed claimer, uint256 fronted, uint256 fee)
func (peginContract *PeginContract) UnpackPegInResolvedEvent(log *types.Log) (*PeginContractPegInResolved, error) {
	event := "PegInResolved"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractPegInResolved)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractRefund represents a Refund event raised by the PeginContract contract.
type PeginContractRefund struct {
	Dest      common.Address
	QuoteHash [32]byte
	Amount    *big.Int
	Success   bool
	Raw       *types.Log // Blockchain specific contextual infos
}

const PeginContractRefundEventName = "Refund"

// ContractEventName returns the user-defined event name.
func (PeginContractRefund) ContractEventName() string {
	return PeginContractRefundEventName
}

// UnpackRefundEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Refund(address indexed dest, bytes32 indexed quoteHash, uint256 indexed amount, bool success)
func (peginContract *PeginContract) UnpackRefundEvent(log *types.Log) (*PeginContractRefund, error) {
	event := "Refund"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractRefund)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractRoleAdminChanged represents a RoleAdminChanged event raised by the PeginContract contract.
type PeginContractRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               *types.Log // Blockchain specific contextual infos
}

const PeginContractRoleAdminChangedEventName = "RoleAdminChanged"

// ContractEventName returns the user-defined event name.
func (PeginContractRoleAdminChanged) ContractEventName() string {
	return PeginContractRoleAdminChangedEventName
}

// UnpackRoleAdminChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (peginContract *PeginContract) UnpackRoleAdminChangedEvent(log *types.Log) (*PeginContractRoleAdminChanged, error) {
	event := "RoleAdminChanged"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractRoleAdminChanged)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractRoleGranted represents a RoleGranted event raised by the PeginContract contract.
type PeginContractRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const PeginContractRoleGrantedEventName = "RoleGranted"

// ContractEventName returns the user-defined event name.
func (PeginContractRoleGranted) ContractEventName() string {
	return PeginContractRoleGrantedEventName
}

// UnpackRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (peginContract *PeginContract) UnpackRoleGrantedEvent(log *types.Log) (*PeginContractRoleGranted, error) {
	event := "RoleGranted"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractRoleGranted)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractRoleRevoked represents a RoleRevoked event raised by the PeginContract contract.
type PeginContractRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const PeginContractRoleRevokedEventName = "RoleRevoked"

// ContractEventName returns the user-defined event name.
func (PeginContractRoleRevoked) ContractEventName() string {
	return PeginContractRoleRevokedEventName
}

// UnpackRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (peginContract *PeginContract) UnpackRoleRevokedEvent(log *types.Log) (*PeginContractRoleRevoked, error) {
	event := "RoleRevoked"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractRoleRevoked)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractUnclaimedPegInSlashed represents a UnclaimedPegInSlashed event raised by the PeginContract contract.
type PeginContractUnclaimedPegInSlashed struct {
	RskAddr common.Address
	Amount  *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const PeginContractUnclaimedPegInSlashedEventName = "UnclaimedPegInSlashed"

// ContractEventName returns the user-defined event name.
func (PeginContractUnclaimedPegInSlashed) ContractEventName() string {
	return PeginContractUnclaimedPegInSlashedEventName
}

// UnpackUnclaimedPegInSlashedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event UnclaimedPegInSlashed(address indexed rskAddr, uint256 indexed amount)
func (peginContract *PeginContract) UnpackUnclaimedPegInSlashedEvent(log *types.Log) (*PeginContractUnclaimedPegInSlashed, error) {
	event := "UnclaimedPegInSlashed"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractUnclaimedPegInSlashed)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractWithdrawal represents a Withdrawal event raised by the PeginContract contract.
type PeginContractWithdrawal struct {
	From   common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const PeginContractWithdrawalEventName = "Withdrawal"

// ContractEventName returns the user-defined event name.
func (PeginContractWithdrawal) ContractEventName() string {
	return PeginContractWithdrawalEventName
}

// UnpackWithdrawalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Withdrawal(address indexed from, uint256 indexed amount)
func (peginContract *PeginContract) UnpackWithdrawalEvent(log *types.Log) (*PeginContractWithdrawal, error) {
	event := "Withdrawal"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractWithdrawal)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (peginContract *PeginContract) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], peginContract.abi.Errors["AccessControlBadConfirmation"].ID.Bytes()[:4]) {
		return peginContract.UnpackAccessControlBadConfirmationError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["AccessControlEnforcedDefaultAdminDelay"].ID.Bytes()[:4]) {
		return peginContract.UnpackAccessControlEnforcedDefaultAdminDelayError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["AccessControlEnforcedDefaultAdminRules"].ID.Bytes()[:4]) {
		return peginContract.UnpackAccessControlEnforcedDefaultAdminRulesError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["AccessControlInvalidDefaultAdmin"].ID.Bytes()[:4]) {
		return peginContract.UnpackAccessControlInvalidDefaultAdminError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["AccessControlUnauthorizedAccount"].ID.Bytes()[:4]) {
		return peginContract.UnpackAccessControlUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["AddressNotRegistered"].ID.Bytes()[:4]) {
		return peginContract.UnpackAddressNotRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["AmountUnderMinimum"].ID.Bytes()[:4]) {
		return peginContract.UnpackAmountUnderMinimumError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["BridgeSettlementFailed"].ID.Bytes()[:4]) {
		return peginContract.UnpackBridgeSettlementFailedError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["ClaimDeadlineNotReached"].ID.Bytes()[:4]) {
		return peginContract.UnpackClaimDeadlineNotReachedError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["DependencyNotSet"].ID.Bytes()[:4]) {
		return peginContract.UnpackDependencyNotSetError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["EmptyBlockHeader"].ID.Bytes()[:4]) {
		return peginContract.UnpackEmptyBlockHeaderError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["EnforcedPause"].ID.Bytes()[:4]) {
		return peginContract.UnpackEnforcedPauseError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["IncorrectContract"].ID.Bytes()[:4]) {
		return peginContract.UnpackIncorrectContractError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["IncorrectFronting"].ID.Bytes()[:4]) {
		return peginContract.UnpackIncorrectFrontingError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["IncorrectSignature"].ID.Bytes()[:4]) {
		return peginContract.UnpackIncorrectSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["InsufficientAmount"].ID.Bytes()[:4]) {
		return peginContract.UnpackInsufficientAmountError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["InsufficientConfirmations"].ID.Bytes()[:4]) {
		return peginContract.UnpackInsufficientConfirmationsError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["InsufficientGas"].ID.Bytes()[:4]) {
		return peginContract.UnpackInsufficientGasError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["InvalidChainId"].ID.Bytes()[:4]) {
		return peginContract.UnpackInvalidChainIdError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return peginContract.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["InvalidRefundAddress"].ID.Bytes()[:4]) {
		return peginContract.UnpackInvalidRefundAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["InvalidSender"].ID.Bytes()[:4]) {
		return peginContract.UnpackInvalidSenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["NoBalance"].ID.Bytes()[:4]) {
		return peginContract.UnpackNoBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["NoContract"].ID.Bytes()[:4]) {
		return peginContract.UnpackNoContractError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["NotEnoughConfirmations"].ID.Bytes()[:4]) {
		return peginContract.UnpackNotEnoughConfirmationsError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return peginContract.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["Overflow"].ID.Bytes()[:4]) {
		return peginContract.UnpackOverflowError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["PaymentFailed"].ID.Bytes()[:4]) {
		return peginContract.UnpackPaymentFailedError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["PaymentNotAllowed"].ID.Bytes()[:4]) {
		return peginContract.UnpackPaymentNotAllowedError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["PegInAlreadyProcessed"].ID.Bytes()[:4]) {
		return peginContract.UnpackPegInAlreadyProcessedError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["PegInNotClaimed"].ID.Bytes()[:4]) {
		return peginContract.UnpackPegInNotClaimedError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["ProviderNotRegistered"].ID.Bytes()[:4]) {
		return peginContract.UnpackProviderNotRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["QuoteAlreadyProcessed"].ID.Bytes()[:4]) {
		return peginContract.UnpackQuoteAlreadyProcessedError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["ReentrancyGuardReentrantCall"].ID.Bytes()[:4]) {
		return peginContract.UnpackReentrancyGuardReentrantCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["SafeCastOverflowedUintDowncast"].ID.Bytes()[:4]) {
		return peginContract.UnpackSafeCastOverflowedUintDowncastError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["UnexpectedBridgeError"].ID.Bytes()[:4]) {
		return peginContract.UnpackUnexpectedBridgeErrorError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// PeginContractAccessControlBadConfirmation represents a AccessControlBadConfirmation error raised by the PeginContract contract.
type PeginContractAccessControlBadConfirmation struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlBadConfirmation()
func PeginContractAccessControlBadConfirmationErrorID() common.Hash {
	return common.HexToHash("0x6697b23232a647058342c0724fe7c415cab25915b54e5dbc03f233173d37b41c")
}

// UnpackAccessControlBadConfirmationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlBadConfirmation()
func (peginContract *PeginContract) UnpackAccessControlBadConfirmationError(raw []byte) (*PeginContractAccessControlBadConfirmation, error) {
	out := new(PeginContractAccessControlBadConfirmation)
	if err := peginContract.abi.UnpackIntoInterface(out, "AccessControlBadConfirmation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractAccessControlEnforcedDefaultAdminDelay represents a AccessControlEnforcedDefaultAdminDelay error raised by the PeginContract contract.
type PeginContractAccessControlEnforcedDefaultAdminDelay struct {
	Schedule *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlEnforcedDefaultAdminDelay(uint48 schedule)
func PeginContractAccessControlEnforcedDefaultAdminDelayErrorID() common.Hash {
	return common.HexToHash("0x19ca5ebb8fb33f00e502c9392eddab1501674629178bf69b853cf037aaf4bb5d")
}

// UnpackAccessControlEnforcedDefaultAdminDelayError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlEnforcedDefaultAdminDelay(uint48 schedule)
func (peginContract *PeginContract) UnpackAccessControlEnforcedDefaultAdminDelayError(raw []byte) (*PeginContractAccessControlEnforcedDefaultAdminDelay, error) {
	out := new(PeginContractAccessControlEnforcedDefaultAdminDelay)
	if err := peginContract.abi.UnpackIntoInterface(out, "AccessControlEnforcedDefaultAdminDelay", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractAccessControlEnforcedDefaultAdminRules represents a AccessControlEnforcedDefaultAdminRules error raised by the PeginContract contract.
type PeginContractAccessControlEnforcedDefaultAdminRules struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlEnforcedDefaultAdminRules()
func PeginContractAccessControlEnforcedDefaultAdminRulesErrorID() common.Hash {
	return common.HexToHash("0x3fc3c27ae3db78c81b8f6e685172134623efa268ee8cd8d54be38ad2a74fc13b")
}

// UnpackAccessControlEnforcedDefaultAdminRulesError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlEnforcedDefaultAdminRules()
func (peginContract *PeginContract) UnpackAccessControlEnforcedDefaultAdminRulesError(raw []byte) (*PeginContractAccessControlEnforcedDefaultAdminRules, error) {
	out := new(PeginContractAccessControlEnforcedDefaultAdminRules)
	if err := peginContract.abi.UnpackIntoInterface(out, "AccessControlEnforcedDefaultAdminRules", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractAccessControlInvalidDefaultAdmin represents a AccessControlInvalidDefaultAdmin error raised by the PeginContract contract.
type PeginContractAccessControlInvalidDefaultAdmin struct {
	DefaultAdmin common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlInvalidDefaultAdmin(address defaultAdmin)
func PeginContractAccessControlInvalidDefaultAdminErrorID() common.Hash {
	return common.HexToHash("0xc22c8022f2a840d6b6a9f113407715f5bbd4e88c1b0dd9434dc00700ba609ed4")
}

// UnpackAccessControlInvalidDefaultAdminError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlInvalidDefaultAdmin(address defaultAdmin)
func (peginContract *PeginContract) UnpackAccessControlInvalidDefaultAdminError(raw []byte) (*PeginContractAccessControlInvalidDefaultAdmin, error) {
	out := new(PeginContractAccessControlInvalidDefaultAdmin)
	if err := peginContract.abi.UnpackIntoInterface(out, "AccessControlInvalidDefaultAdmin", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractAccessControlUnauthorizedAccount represents a AccessControlUnauthorizedAccount error raised by the PeginContract contract.
type PeginContractAccessControlUnauthorizedAccount struct {
	Account    common.Address
	NeededRole [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func PeginContractAccessControlUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0xe2517d3fbfae6f8515ef5ff1ccedc3933ab0cbbda0b492c06eb54ad10ef03b3e")
}

// UnpackAccessControlUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func (peginContract *PeginContract) UnpackAccessControlUnauthorizedAccountError(raw []byte) (*PeginContractAccessControlUnauthorizedAccount, error) {
	out := new(PeginContractAccessControlUnauthorizedAccount)
	if err := peginContract.abi.UnpackIntoInterface(out, "AccessControlUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractAddressNotRegistered represents a AddressNotRegistered error raised by the PeginContract contract.
type PeginContractAddressNotRegistered struct {
	RskAddr common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressNotRegistered(address rskAddr)
func PeginContractAddressNotRegisteredErrorID() common.Hash {
	return common.HexToHash("0xf3a8c4eea98f33ac9ed01788a298f73b7cc71b3532c7507ff436e3d37bd70a2a")
}

// UnpackAddressNotRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressNotRegistered(address rskAddr)
func (peginContract *PeginContract) UnpackAddressNotRegisteredError(raw []byte) (*PeginContractAddressNotRegistered, error) {
	out := new(PeginContractAddressNotRegistered)
	if err := peginContract.abi.UnpackIntoInterface(out, "AddressNotRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractAmountUnderMinimum represents a AmountUnderMinimum error raised by the PeginContract contract.
type PeginContractAmountUnderMinimum struct {
	Amount *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AmountUnderMinimum(uint256 amount)
func PeginContractAmountUnderMinimumErrorID() common.Hash {
	return common.HexToHash("0x12b21ac5391b7c8c532d3e1a87c2e8173cdc084bbd0c28c2cabc2c2c296d2f1a")
}

// UnpackAmountUnderMinimumError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AmountUnderMinimum(uint256 amount)
func (peginContract *PeginContract) UnpackAmountUnderMinimumError(raw []byte) (*PeginContractAmountUnderMinimum, error) {
	out := new(PeginContractAmountUnderMinimum)
	if err := peginContract.abi.UnpackIntoInterface(out, "AmountUnderMinimum", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractBridgeSettlementFailed represents a BridgeSettlementFailed error raised by the PeginContract contract.
type PeginContractBridgeSettlementFailed struct {
	BridgeResult *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BridgeSettlementFailed(int256 bridgeResult)
func PeginContractBridgeSettlementFailedErrorID() common.Hash {
	return common.HexToHash("0xd4fb298cb8cbd68b7728ffaa64e915a5c316b593be120215bd8728f5eb1828f4")
}

// UnpackBridgeSettlementFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BridgeSettlementFailed(int256 bridgeResult)
func (peginContract *PeginContract) UnpackBridgeSettlementFailedError(raw []byte) (*PeginContractBridgeSettlementFailed, error) {
	out := new(PeginContractBridgeSettlementFailed)
	if err := peginContract.abi.UnpackIntoInterface(out, "BridgeSettlementFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractClaimDeadlineNotReached represents a ClaimDeadlineNotReached error raised by the PeginContract contract.
type PeginContractClaimDeadlineNotReached struct {
	DeadlineBlock *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ClaimDeadlineNotReached(uint256 deadlineBlock)
func PeginContractClaimDeadlineNotReachedErrorID() common.Hash {
	return common.HexToHash("0xe44567bb7bcfc4d6e430408d231f8b6a7b66f40ad586d1605d8b4dd776e25216")
}

// UnpackClaimDeadlineNotReachedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ClaimDeadlineNotReached(uint256 deadlineBlock)
func (peginContract *PeginContract) UnpackClaimDeadlineNotReachedError(raw []byte) (*PeginContractClaimDeadlineNotReached, error) {
	out := new(PeginContractClaimDeadlineNotReached)
	if err := peginContract.abi.UnpackIntoInterface(out, "ClaimDeadlineNotReached", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractDependencyNotSet represents a DependencyNotSet error raised by the PeginContract contract.
type PeginContractDependencyNotSet struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DependencyNotSet()
func PeginContractDependencyNotSetErrorID() common.Hash {
	return common.HexToHash("0x215af8cd3fa877313f255ca855809ebe76b7fb0dc776b632f6b0a9dcbb175e0a")
}

// UnpackDependencyNotSetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DependencyNotSet()
func (peginContract *PeginContract) UnpackDependencyNotSetError(raw []byte) (*PeginContractDependencyNotSet, error) {
	out := new(PeginContractDependencyNotSet)
	if err := peginContract.abi.UnpackIntoInterface(out, "DependencyNotSet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractEmptyBlockHeader represents a EmptyBlockHeader error raised by the PeginContract contract.
type PeginContractEmptyBlockHeader struct {
	HeightOrHash [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EmptyBlockHeader(bytes32 heightOrHash)
func PeginContractEmptyBlockHeaderErrorID() common.Hash {
	return common.HexToHash("0xc1a923b4e595599b5ebca706a34bfaa111ec5aad01c417609e91334f899d99e4")
}

// UnpackEmptyBlockHeaderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EmptyBlockHeader(bytes32 heightOrHash)
func (peginContract *PeginContract) UnpackEmptyBlockHeaderError(raw []byte) (*PeginContractEmptyBlockHeader, error) {
	out := new(PeginContractEmptyBlockHeader)
	if err := peginContract.abi.UnpackIntoInterface(out, "EmptyBlockHeader", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractEnforcedPause represents a EnforcedPause error raised by the PeginContract contract.
type PeginContractEnforcedPause struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EnforcedPause()
func PeginContractEnforcedPauseErrorID() common.Hash {
	return common.HexToHash("0xd93c0665d6c96d04a8f174024fc4ddd66c250604aff22bbec808de86dd3637e3")
}

// UnpackEnforcedPauseError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EnforcedPause()
func (peginContract *PeginContract) UnpackEnforcedPauseError(raw []byte) (*PeginContractEnforcedPause, error) {
	out := new(PeginContractEnforcedPause)
	if err := peginContract.abi.UnpackIntoInterface(out, "EnforcedPause", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractIncorrectContract represents a IncorrectContract error raised by the PeginContract contract.
type PeginContractIncorrectContract struct {
	Expected common.Address
	Actual   common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error IncorrectContract(address expected, address actual)
func PeginContractIncorrectContractErrorID() common.Hash {
	return common.HexToHash("0x367b77278f2bd6dd9afab1117babaedb89c7f420646aa9a343c7f6bd654b7740")
}

// UnpackIncorrectContractError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error IncorrectContract(address expected, address actual)
func (peginContract *PeginContract) UnpackIncorrectContractError(raw []byte) (*PeginContractIncorrectContract, error) {
	out := new(PeginContractIncorrectContract)
	if err := peginContract.abi.UnpackIntoInterface(out, "IncorrectContract", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractIncorrectFronting represents a IncorrectFronting error raised by the PeginContract contract.
type PeginContractIncorrectFronting struct {
	Expected *big.Int
	Provided *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error IncorrectFronting(uint256 expected, uint256 provided)
func PeginContractIncorrectFrontingErrorID() common.Hash {
	return common.HexToHash("0x4c0c9bd9730b66d0ffcd547ddbe6f891acd58a39d5bbd9bd3984d4e760e964fe")
}

// UnpackIncorrectFrontingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error IncorrectFronting(uint256 expected, uint256 provided)
func (peginContract *PeginContract) UnpackIncorrectFrontingError(raw []byte) (*PeginContractIncorrectFronting, error) {
	out := new(PeginContractIncorrectFronting)
	if err := peginContract.abi.UnpackIntoInterface(out, "IncorrectFronting", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractIncorrectSignature represents a IncorrectSignature error raised by the PeginContract contract.
type PeginContractIncorrectSignature struct {
	ExpectedAddress common.Address
	UsedHash        [32]byte
	Signature       []byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error IncorrectSignature(address expectedAddress, bytes32 usedHash, bytes signature)
func PeginContractIncorrectSignatureErrorID() common.Hash {
	return common.HexToHash("0xf6c6db712ba972f36fab8afd40173e51d123cf1086d0e678bb4ec5dbae15bfe6")
}

// UnpackIncorrectSignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error IncorrectSignature(address expectedAddress, bytes32 usedHash, bytes signature)
func (peginContract *PeginContract) UnpackIncorrectSignatureError(raw []byte) (*PeginContractIncorrectSignature, error) {
	out := new(PeginContractIncorrectSignature)
	if err := peginContract.abi.UnpackIntoInterface(out, "IncorrectSignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractInsufficientAmount represents a InsufficientAmount error raised by the PeginContract contract.
type PeginContractInsufficientAmount struct {
	Amount *big.Int
	Target *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientAmount(uint256 amount, uint256 target)
func PeginContractInsufficientAmountErrorID() common.Hash {
	return common.HexToHash("0x6d400e382e49fdfa5a03b18b5b2c938638a3fb351ac4810276f70a093eb3f20f")
}

// UnpackInsufficientAmountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientAmount(uint256 amount, uint256 target)
func (peginContract *PeginContract) UnpackInsufficientAmountError(raw []byte) (*PeginContractInsufficientAmount, error) {
	out := new(PeginContractInsufficientAmount)
	if err := peginContract.abi.UnpackIntoInterface(out, "InsufficientAmount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractInsufficientConfirmations represents a InsufficientConfirmations error raised by the PeginContract contract.
type PeginContractInsufficientConfirmations struct {
	Have     *big.Int
	Required *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientConfirmations(uint256 have, uint256 required)
func PeginContractInsufficientConfirmationsErrorID() common.Hash {
	return common.HexToHash("0x22c53f1f9da074136409181cd7a7bfae7388312a843098ca3ab118f5de3d1890")
}

// UnpackInsufficientConfirmationsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientConfirmations(uint256 have, uint256 required)
func (peginContract *PeginContract) UnpackInsufficientConfirmationsError(raw []byte) (*PeginContractInsufficientConfirmations, error) {
	out := new(PeginContractInsufficientConfirmations)
	if err := peginContract.abi.UnpackIntoInterface(out, "InsufficientConfirmations", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractInsufficientGas represents a InsufficientGas error raised by the PeginContract contract.
type PeginContractInsufficientGas struct {
	GasLeft     *big.Int
	GasRequired *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientGas(uint256 gasLeft, uint256 gasRequired)
func PeginContractInsufficientGasErrorID() common.Hash {
	return common.HexToHash("0x23e228cb30ba888d86279f51c27a6690ec0fe016ea074bbd6332de5b2571f9b3")
}

// UnpackInsufficientGasError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientGas(uint256 gasLeft, uint256 gasRequired)
func (peginContract *PeginContract) UnpackInsufficientGasError(raw []byte) (*PeginContractInsufficientGas, error) {
	out := new(PeginContractInsufficientGas)
	if err := peginContract.abi.UnpackIntoInterface(out, "InsufficientGas", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractInvalidChainId represents a InvalidChainId error raised by the PeginContract contract.
type PeginContractInvalidChainId struct {
	Expected *big.Int
	Actual   *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidChainId(uint256 expected, uint256 actual)
func PeginContractInvalidChainIdErrorID() common.Hash {
	return common.HexToHash("0x9fba672f6672f44dc27784b208145b99fe892bc4b6a497de84dfb629e18d300e")
}

// UnpackInvalidChainIdError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidChainId(uint256 expected, uint256 actual)
func (peginContract *PeginContract) UnpackInvalidChainIdError(raw []byte) (*PeginContractInvalidChainId, error) {
	out := new(PeginContractInvalidChainId)
	if err := peginContract.abi.UnpackIntoInterface(out, "InvalidChainId", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractInvalidInitialization represents a InvalidInitialization error raised by the PeginContract contract.
type PeginContractInvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func PeginContractInvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (peginContract *PeginContract) UnpackInvalidInitializationError(raw []byte) (*PeginContractInvalidInitialization, error) {
	out := new(PeginContractInvalidInitialization)
	if err := peginContract.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractInvalidRefundAddress represents a InvalidRefundAddress error raised by the PeginContract contract.
type PeginContractInvalidRefundAddress struct {
	RefundAddress []byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidRefundAddress(bytes refundAddress)
func PeginContractInvalidRefundAddressErrorID() common.Hash {
	return common.HexToHash("0x17b97374cf0657ace771eafcd95943e733cce4db0d5b4e147411067d764c8f99")
}

// UnpackInvalidRefundAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidRefundAddress(bytes refundAddress)
func (peginContract *PeginContract) UnpackInvalidRefundAddressError(raw []byte) (*PeginContractInvalidRefundAddress, error) {
	out := new(PeginContractInvalidRefundAddress)
	if err := peginContract.abi.UnpackIntoInterface(out, "InvalidRefundAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractInvalidSender represents a InvalidSender error raised by the PeginContract contract.
type PeginContractInvalidSender struct {
	Expected common.Address
	Actual   common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSender(address expected, address actual)
func PeginContractInvalidSenderErrorID() common.Hash {
	return common.HexToHash("0xe1130dbad6e77228912cd79cc3b53cd156f090ef6a73d9fdb2720c4f9d40d9d3")
}

// UnpackInvalidSenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSender(address expected, address actual)
func (peginContract *PeginContract) UnpackInvalidSenderError(raw []byte) (*PeginContractInvalidSender, error) {
	out := new(PeginContractInvalidSender)
	if err := peginContract.abi.UnpackIntoInterface(out, "InvalidSender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractNoBalance represents a NoBalance error raised by the PeginContract contract.
type PeginContractNoBalance struct {
	Wanted *big.Int
	Actual *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoBalance(uint256 wanted, uint256 actual)
func PeginContractNoBalanceErrorID() common.Hash {
	return common.HexToHash("0x292266533ee1631c0f0faf752ebfa5783238c0e3e2fcdef002a1685294062289")
}

// UnpackNoBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoBalance(uint256 wanted, uint256 actual)
func (peginContract *PeginContract) UnpackNoBalanceError(raw []byte) (*PeginContractNoBalance, error) {
	out := new(PeginContractNoBalance)
	if err := peginContract.abi.UnpackIntoInterface(out, "NoBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractNoContract represents a NoContract error raised by the PeginContract contract.
type PeginContractNoContract struct {
	Addr common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoContract(address addr)
func PeginContractNoContractErrorID() common.Hash {
	return common.HexToHash("0x5f15d672b6235f8600ffc72925d8d2f9dcea14be067296327891153847185a5c")
}

// UnpackNoContractError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoContract(address addr)
func (peginContract *PeginContract) UnpackNoContractError(raw []byte) (*PeginContractNoContract, error) {
	out := new(PeginContractNoContract)
	if err := peginContract.abi.UnpackIntoInterface(out, "NoContract", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractNotEnoughConfirmations represents a NotEnoughConfirmations error raised by the PeginContract contract.
type PeginContractNotEnoughConfirmations struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotEnoughConfirmations()
func PeginContractNotEnoughConfirmationsErrorID() common.Hash {
	return common.HexToHash("0xb9310b562727f0fb376475537f3a4e5f39f5fed59dbda43f984b828c8ef037d0")
}

// UnpackNotEnoughConfirmationsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotEnoughConfirmations()
func (peginContract *PeginContract) UnpackNotEnoughConfirmationsError(raw []byte) (*PeginContractNotEnoughConfirmations, error) {
	out := new(PeginContractNotEnoughConfirmations)
	if err := peginContract.abi.UnpackIntoInterface(out, "NotEnoughConfirmations", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractNotInitializing represents a NotInitializing error raised by the PeginContract contract.
type PeginContractNotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func PeginContractNotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (peginContract *PeginContract) UnpackNotInitializingError(raw []byte) (*PeginContractNotInitializing, error) {
	out := new(PeginContractNotInitializing)
	if err := peginContract.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractOverflow represents a Overflow error raised by the PeginContract contract.
type PeginContractOverflow struct {
	PassedAmount *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Overflow(uint256 passedAmount)
func PeginContractOverflowErrorID() common.Hash {
	return common.HexToHash("0xe0fb6a7ce291b396fa814871fbb6fcc26c1a1454a6e18a2e7c911a8763b928dc")
}

// UnpackOverflowError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Overflow(uint256 passedAmount)
func (peginContract *PeginContract) UnpackOverflowError(raw []byte) (*PeginContractOverflow, error) {
	out := new(PeginContractOverflow)
	if err := peginContract.abi.UnpackIntoInterface(out, "Overflow", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractPaymentFailed represents a PaymentFailed error raised by the PeginContract contract.
type PeginContractPaymentFailed struct {
	Addr   common.Address
	Amount *big.Int
	Reason []byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PaymentFailed(address addr, uint256 amount, bytes reason)
func PeginContractPaymentFailedErrorID() common.Hash {
	return common.HexToHash("0xadca8d516d2aaa483a86cefb25d722eccb15750e54cc37c21033c80cc79b13e3")
}

// UnpackPaymentFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PaymentFailed(address addr, uint256 amount, bytes reason)
func (peginContract *PeginContract) UnpackPaymentFailedError(raw []byte) (*PeginContractPaymentFailed, error) {
	out := new(PeginContractPaymentFailed)
	if err := peginContract.abi.UnpackIntoInterface(out, "PaymentFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractPaymentNotAllowed represents a PaymentNotAllowed error raised by the PeginContract contract.
type PeginContractPaymentNotAllowed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PaymentNotAllowed()
func PeginContractPaymentNotAllowedErrorID() common.Hash {
	return common.HexToHash("0x8619bd43ab22b4b01742bd29d231dff1e50413ee3a444878bed65970c80c97df")
}

// UnpackPaymentNotAllowedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PaymentNotAllowed()
func (peginContract *PeginContract) UnpackPaymentNotAllowedError(raw []byte) (*PeginContractPaymentNotAllowed, error) {
	out := new(PeginContractPaymentNotAllowed)
	if err := peginContract.abi.UnpackIntoInterface(out, "PaymentNotAllowed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractPegInAlreadyProcessed represents a PegInAlreadyProcessed error raised by the PeginContract contract.
type PeginContractPegInAlreadyProcessed struct {
	PegInId [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PegInAlreadyProcessed(bytes32 pegInId)
func PeginContractPegInAlreadyProcessedErrorID() common.Hash {
	return common.HexToHash("0x61ef28332480fb1a452c181f028122bc7068d7c84387e53bbca7f60f4f3cb6f9")
}

// UnpackPegInAlreadyProcessedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PegInAlreadyProcessed(bytes32 pegInId)
func (peginContract *PeginContract) UnpackPegInAlreadyProcessedError(raw []byte) (*PeginContractPegInAlreadyProcessed, error) {
	out := new(PeginContractPegInAlreadyProcessed)
	if err := peginContract.abi.UnpackIntoInterface(out, "PegInAlreadyProcessed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractPegInNotClaimed represents a PegInNotClaimed error raised by the PeginContract contract.
type PeginContractPegInNotClaimed struct {
	PegInId [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PegInNotClaimed(bytes32 pegInId)
func PeginContractPegInNotClaimedErrorID() common.Hash {
	return common.HexToHash("0x913cecc833ab31c9001ed04ef71cd13c55d40d8de5a2e4b9ec4ce9bd89ae245d")
}

// UnpackPegInNotClaimedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PegInNotClaimed(bytes32 pegInId)
func (peginContract *PeginContract) UnpackPegInNotClaimedError(raw []byte) (*PeginContractPegInNotClaimed, error) {
	out := new(PeginContractPegInNotClaimed)
	if err := peginContract.abi.UnpackIntoInterface(out, "PegInNotClaimed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractProviderNotRegistered represents a ProviderNotRegistered error raised by the PeginContract contract.
type PeginContractProviderNotRegistered struct {
	From common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ProviderNotRegistered(address from)
func PeginContractProviderNotRegisteredErrorID() common.Hash {
	return common.HexToHash("0x232cb27a4b9e96657e43917628f0b0ddd34885ba8495a2108b78da7512210fb9")
}

// UnpackProviderNotRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ProviderNotRegistered(address from)
func (peginContract *PeginContract) UnpackProviderNotRegisteredError(raw []byte) (*PeginContractProviderNotRegistered, error) {
	out := new(PeginContractProviderNotRegistered)
	if err := peginContract.abi.UnpackIntoInterface(out, "ProviderNotRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractQuoteAlreadyProcessed represents a QuoteAlreadyProcessed error raised by the PeginContract contract.
type PeginContractQuoteAlreadyProcessed struct {
	QuoteHash [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error QuoteAlreadyProcessed(bytes32 quoteHash)
func PeginContractQuoteAlreadyProcessedErrorID() common.Hash {
	return common.HexToHash("0xda4bb665b5917c51a08a9d79f0cf72c95ee90f52ca97526abf4b52ae2d737c77")
}

// UnpackQuoteAlreadyProcessedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error QuoteAlreadyProcessed(bytes32 quoteHash)
func (peginContract *PeginContract) UnpackQuoteAlreadyProcessedError(raw []byte) (*PeginContractQuoteAlreadyProcessed, error) {
	out := new(PeginContractQuoteAlreadyProcessed)
	if err := peginContract.abi.UnpackIntoInterface(out, "QuoteAlreadyProcessed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractReentrancyGuardReentrantCall represents a ReentrancyGuardReentrantCall error raised by the PeginContract contract.
type PeginContractReentrancyGuardReentrantCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ReentrancyGuardReentrantCall()
func PeginContractReentrancyGuardReentrantCallErrorID() common.Hash {
	return common.HexToHash("0x3ee5aeb571de7fc460830b4d0017439a1ca56fb0bc39062227ade4fe4a24c1ca")
}

// UnpackReentrancyGuardReentrantCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ReentrancyGuardReentrantCall()
func (peginContract *PeginContract) UnpackReentrancyGuardReentrantCallError(raw []byte) (*PeginContractReentrancyGuardReentrantCall, error) {
	out := new(PeginContractReentrancyGuardReentrantCall)
	if err := peginContract.abi.UnpackIntoInterface(out, "ReentrancyGuardReentrantCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractSafeCastOverflowedUintDowncast represents a SafeCastOverflowedUintDowncast error raised by the PeginContract contract.
type PeginContractSafeCastOverflowedUintDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func PeginContractSafeCastOverflowedUintDowncastErrorID() common.Hash {
	return common.HexToHash("0x6dfcc6503a32754ce7a89698e18201fc5294fd4aad43edefee786f88423b1a12")
}

// UnpackSafeCastOverflowedUintDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func (peginContract *PeginContract) UnpackSafeCastOverflowedUintDowncastError(raw []byte) (*PeginContractSafeCastOverflowedUintDowncast, error) {
	out := new(PeginContractSafeCastOverflowedUintDowncast)
	if err := peginContract.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractUnexpectedBridgeError represents a UnexpectedBridgeError error raised by the PeginContract contract.
type PeginContractUnexpectedBridgeError struct {
	ErrorCode *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnexpectedBridgeError(int256 errorCode)
func PeginContractUnexpectedBridgeErrorErrorID() common.Hash {
	return common.HexToHash("0xab2e19b07f35862bb35f42ec9b8ce9397c135ed8ed749f2117205465c4166a09")
}

// UnpackUnexpectedBridgeErrorError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnexpectedBridgeError(int256 errorCode)
func (peginContract *PeginContract) UnpackUnexpectedBridgeErrorError(raw []byte) (*PeginContractUnexpectedBridgeError, error) {
	out := new(PeginContractUnexpectedBridgeError)
	if err := peginContract.abi.UnpackIntoInterface(out, "UnexpectedBridgeError", raw); err != nil {
		return nil, err
	}
	return out, nil
}
