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

// IFlyoverConfigurationsConfirmationTier is an auto generated low-level Go binding around an user-defined struct.
type IFlyoverConfigurationsConfirmationTier struct {
	MaxAmount     *big.Int
	Confirmations *big.Int
}

// IFlyoverConfigurationsPegConfiguration is an auto generated low-level Go binding around an user-defined struct.
type IFlyoverConfigurationsPegConfiguration struct {
	FixedFee          *big.Int
	PercentageFee     *big.Int
	PenaltyFee        *big.Int
	ConfirmationTiers []IFlyoverConfigurationsConfirmationTier
	CallTime          *big.Int
	ExpireTime        *big.Int
	ExpireBlocks      *big.Int
	DeliveryGrace     *big.Int
	MinAmount         *big.Int
	MaxAmount         *big.Int
}

// FlyoverConfigurationsMetaData contains all meta data concerning the FlyoverConfigurations contract.
var FlyoverConfigurationsMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"FEE_PERCENTAGE_DENOMINATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SAT_TO_WEI_CONVERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyChange\",\"inputs\":[{\"name\":\"flow\",\"type\":\"uint8\",\"internalType\":\"enumFlyoverConfigurations.Flow\"},{\"name\":\"field\",\"type\":\"uint8\",\"internalType\":\"enumFlyoverConfigurations.Field\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyTiersChange\",\"inputs\":[{\"name\":\"flow\",\"type\":\"uint8\",\"internalType\":\"enumFlyoverConfigurations.Flow\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"beginDefaultAdminTransfer\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"calculatePegInFee\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"calculatePegOutFee\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancelDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"changeDefaultAdminDelay\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"defaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelayIncreaseWait\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPegInConfiguration\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIFlyoverConfigurations.PegConfiguration\",\"components\":[{\"name\":\"fixedFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"percentageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmationTiers\",\"type\":\"tuple[]\",\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"callTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deliveryGrace\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPegInConfigurationBounds\",\"inputs\":[],\"outputs\":[{\"name\":\"min\",\"type\":\"tuple\",\"internalType\":\"structIFlyoverConfigurations.PegConfiguration\",\"components\":[{\"name\":\"fixedFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"percentageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmationTiers\",\"type\":\"tuple[]\",\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"callTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deliveryGrace\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"max\",\"type\":\"tuple\",\"internalType\":\"structIFlyoverConfigurations.PegConfiguration\",\"components\":[{\"name\":\"fixedFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"percentageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmationTiers\",\"type\":\"tuple[]\",\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"callTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deliveryGrace\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPegOutConfiguration\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIFlyoverConfigurations.PegConfiguration\",\"components\":[{\"name\":\"fixedFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"percentageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmationTiers\",\"type\":\"tuple[]\",\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"callTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deliveryGrace\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPegOutConfigurationBounds\",\"inputs\":[],\"outputs\":[{\"name\":\"min\",\"type\":\"tuple\",\"internalType\":\"structIFlyoverConfigurations.PegConfiguration\",\"components\":[{\"name\":\"fixedFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"percentageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmationTiers\",\"type\":\"tuple[]\",\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"callTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deliveryGrace\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"max\",\"type\":\"tuple\",\"internalType\":\"structIFlyoverConfigurations.PegConfiguration\",\"components\":[{\"name\":\"fixedFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"percentageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmationTiers\",\"type\":\"tuple[]\",\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"callTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deliveryGrace\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPendingChange\",\"inputs\":[{\"name\":\"flow\",\"type\":\"uint8\",\"internalType\":\"enumFlyoverConfigurations.Flow\"},{\"name\":\"field\",\"type\":\"uint8\",\"internalType\":\"enumFlyoverConfigurations.Field\"}],\"outputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"eta\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRequiredPegInConfirmations\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRequiredPegOutConfirmations\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTimelockDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"defaultAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"initialDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"timelockDelay\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pegInConfig\",\"type\":\"tuple\",\"internalType\":\"structIFlyoverConfigurations.PegConfiguration\",\"components\":[{\"name\":\"fixedFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"percentageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmationTiers\",\"type\":\"tuple[]\",\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"callTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deliveryGrace\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pegOutConfig\",\"type\":\"tuple\",\"internalType\":\"structIFlyoverConfigurations.PegConfiguration\",\"components\":[{\"name\":\"fixedFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"percentageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmationTiers\",\"type\":\"tuple[]\",\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"callTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deliveryGrace\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pegInMin\",\"type\":\"tuple\",\"internalType\":\"structIFlyoverConfigurations.PegConfiguration\",\"components\":[{\"name\":\"fixedFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"percentageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmationTiers\",\"type\":\"tuple[]\",\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"callTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deliveryGrace\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pegInMax\",\"type\":\"tuple\",\"internalType\":\"structIFlyoverConfigurations.PegConfiguration\",\"components\":[{\"name\":\"fixedFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"percentageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmationTiers\",\"type\":\"tuple[]\",\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"callTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deliveryGrace\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pegOutMin\",\"type\":\"tuple\",\"internalType\":\"structIFlyoverConfigurations.PegConfiguration\",\"components\":[{\"name\":\"fixedFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"percentageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmationTiers\",\"type\":\"tuple[]\",\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"callTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deliveryGrace\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pegOutMax\",\"type\":\"tuple\",\"internalType\":\"structIFlyoverConfigurations.PegConfiguration\",\"components\":[{\"name\":\"fixedFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"percentageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmationTiers\",\"type\":\"tuple[]\",\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"callTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deliveryGrace\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"queueChange\",\"inputs\":[{\"name\":\"flow\",\"type\":\"uint8\",\"internalType\":\"enumFlyoverConfigurations.Flow\"},{\"name\":\"field\",\"type\":\"uint8\",\"internalType\":\"enumFlyoverConfigurations.Field\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"queueTiersChange\",\"inputs\":[{\"name\":\"flow\",\"type\":\"uint8\",\"internalType\":\"enumFlyoverConfigurations.Flow\"},{\"name\":\"tiers\",\"type\":\"tuple[]\",\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rollbackDefaultAdminDelay\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeScheduled\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"effectSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferScheduled\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"acceptSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegInCallTimeChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegInConfirmationTiersChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"newValue\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegInDeliveryGraceChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegInExpireBlocksChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegInExpireTimeChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegInFixedFeeChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegInMaxAmountChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegInMinAmountChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegInPenaltyFeeChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegInPercentageFeeChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutCallTimeChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutConfirmationTiersChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"newValue\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutDeliveryGraceChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutExpireBlocksChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutExpireTimeChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutFixedFeeChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutMaxAmountChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutMinAmountChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutPenaltyFeeChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutPercentageFeeChanged\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminDelay\",\"inputs\":[{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminRules\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlInvalidDefaultAdmin\",\"inputs\":[{\"name\":\"defaultAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"EmptyTiers\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpireTimeNotAfterCallTime\",\"inputs\":[{\"name\":\"callTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoQueuedChange\",\"inputs\":[{\"name\":\"flow\",\"type\":\"uint8\",\"internalType\":\"enumFlyoverConfigurations.Flow\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OutOfBounds\",\"inputs\":[{\"name\":\"flow\",\"type\":\"uint8\",\"internalType\":\"enumFlyoverConfigurations.Flow\"},{\"name\":\"field\",\"type\":\"uint8\",\"internalType\":\"enumFlyoverConfigurations.Field\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"min\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"PaymentNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TiersNotAscending\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TimelockNotElapsed\",\"inputs\":[{\"name\":\"eta\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"nowTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
	ID:  "FlyoverConfigurations",
}

// FlyoverConfigurations is an auto generated Go binding around an Ethereum contract.
type FlyoverConfigurations struct {
	abi abi.ABI
}

// NewFlyoverConfigurations creates a new instance of FlyoverConfigurations.
func NewFlyoverConfigurations() *FlyoverConfigurations {
	parsed, err := FlyoverConfigurationsMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &FlyoverConfigurations{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *FlyoverConfigurations) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (flyoverConfigurations *FlyoverConfigurations) PackDEFAULTADMINROLE() []byte {
	enc, err := flyoverConfigurations.abi.Pack("DEFAULT_ADMIN_ROLE")
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackDEFAULTADMINROLE() ([]byte, error) {
	return flyoverConfigurations.abi.Pack("DEFAULT_ADMIN_ROLE")
}

// UnpackDEFAULTADMINROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (flyoverConfigurations *FlyoverConfigurations) UnpackDEFAULTADMINROLE(data []byte) ([32]byte, error) {
	out, err := flyoverConfigurations.abi.Unpack("DEFAULT_ADMIN_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackFEEPERCENTAGEDENOMINATOR is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x54b01b84.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function FEE_PERCENTAGE_DENOMINATOR() view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) PackFEEPERCENTAGEDENOMINATOR() []byte {
	enc, err := flyoverConfigurations.abi.Pack("FEE_PERCENTAGE_DENOMINATOR")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFEEPERCENTAGEDENOMINATOR is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x54b01b84.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function FEE_PERCENTAGE_DENOMINATOR() view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) TryPackFEEPERCENTAGEDENOMINATOR() ([]byte, error) {
	return flyoverConfigurations.abi.Pack("FEE_PERCENTAGE_DENOMINATOR")
}

// UnpackFEEPERCENTAGEDENOMINATOR is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x54b01b84.
//
// Solidity: function FEE_PERCENTAGE_DENOMINATOR() view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) UnpackFEEPERCENTAGEDENOMINATOR(data []byte) (*big.Int, error) {
	out, err := flyoverConfigurations.abi.Unpack("FEE_PERCENTAGE_DENOMINATOR", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackSATTOWEICONVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb5ecfc06.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function SAT_TO_WEI_CONVERSION() view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) PackSATTOWEICONVERSION() []byte {
	enc, err := flyoverConfigurations.abi.Pack("SAT_TO_WEI_CONVERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSATTOWEICONVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb5ecfc06.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function SAT_TO_WEI_CONVERSION() view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) TryPackSATTOWEICONVERSION() ([]byte, error) {
	return flyoverConfigurations.abi.Pack("SAT_TO_WEI_CONVERSION")
}

// UnpackSATTOWEICONVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb5ecfc06.
//
// Solidity: function SAT_TO_WEI_CONVERSION() view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) UnpackSATTOWEICONVERSION(data []byte) (*big.Int, error) {
	out, err := flyoverConfigurations.abi.Unpack("SAT_TO_WEI_CONVERSION", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xffa1ad74.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function VERSION() view returns(string)
func (flyoverConfigurations *FlyoverConfigurations) PackVERSION() []byte {
	enc, err := flyoverConfigurations.abi.Pack("VERSION")
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackVERSION() ([]byte, error) {
	return flyoverConfigurations.abi.Pack("VERSION")
}

// UnpackVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (flyoverConfigurations *FlyoverConfigurations) UnpackVERSION(data []byte) (string, error) {
	out, err := flyoverConfigurations.abi.Unpack("VERSION", data)
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
func (flyoverConfigurations *FlyoverConfigurations) PackAcceptDefaultAdminTransfer() []byte {
	enc, err := flyoverConfigurations.abi.Pack("acceptDefaultAdminTransfer")
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackAcceptDefaultAdminTransfer() ([]byte, error) {
	return flyoverConfigurations.abi.Pack("acceptDefaultAdminTransfer")
}

// PackApplyChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6c9c5735.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function applyChange(uint8 flow, uint8 field) returns()
func (flyoverConfigurations *FlyoverConfigurations) PackApplyChange(flow uint8, field uint8) []byte {
	enc, err := flyoverConfigurations.abi.Pack("applyChange", flow, field)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApplyChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6c9c5735.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function applyChange(uint8 flow, uint8 field) returns()
func (flyoverConfigurations *FlyoverConfigurations) TryPackApplyChange(flow uint8, field uint8) ([]byte, error) {
	return flyoverConfigurations.abi.Pack("applyChange", flow, field)
}

// PackApplyTiersChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c5d7c04.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function applyTiersChange(uint8 flow) returns()
func (flyoverConfigurations *FlyoverConfigurations) PackApplyTiersChange(flow uint8) []byte {
	enc, err := flyoverConfigurations.abi.Pack("applyTiersChange", flow)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApplyTiersChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c5d7c04.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function applyTiersChange(uint8 flow) returns()
func (flyoverConfigurations *FlyoverConfigurations) TryPackApplyTiersChange(flow uint8) ([]byte, error) {
	return flyoverConfigurations.abi.Pack("applyTiersChange", flow)
}

// PackBeginDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x634e93da.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (flyoverConfigurations *FlyoverConfigurations) PackBeginDefaultAdminTransfer(newAdmin common.Address) []byte {
	enc, err := flyoverConfigurations.abi.Pack("beginDefaultAdminTransfer", newAdmin)
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackBeginDefaultAdminTransfer(newAdmin common.Address) ([]byte, error) {
	return flyoverConfigurations.abi.Pack("beginDefaultAdminTransfer", newAdmin)
}

// PackCalculatePegInFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715a177c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function calculatePegInFee(uint256 amount) view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) PackCalculatePegInFee(amount *big.Int) []byte {
	enc, err := flyoverConfigurations.abi.Pack("calculatePegInFee", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCalculatePegInFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715a177c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function calculatePegInFee(uint256 amount) view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) TryPackCalculatePegInFee(amount *big.Int) ([]byte, error) {
	return flyoverConfigurations.abi.Pack("calculatePegInFee", amount)
}

// UnpackCalculatePegInFee is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x715a177c.
//
// Solidity: function calculatePegInFee(uint256 amount) view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) UnpackCalculatePegInFee(data []byte) (*big.Int, error) {
	out, err := flyoverConfigurations.abi.Unpack("calculatePegInFee", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackCalculatePegOutFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x516ba0b3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function calculatePegOutFee(uint256 amount) view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) PackCalculatePegOutFee(amount *big.Int) []byte {
	enc, err := flyoverConfigurations.abi.Pack("calculatePegOutFee", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCalculatePegOutFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x516ba0b3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function calculatePegOutFee(uint256 amount) view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) TryPackCalculatePegOutFee(amount *big.Int) ([]byte, error) {
	return flyoverConfigurations.abi.Pack("calculatePegOutFee", amount)
}

// UnpackCalculatePegOutFee is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x516ba0b3.
//
// Solidity: function calculatePegOutFee(uint256 amount) view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) UnpackCalculatePegOutFee(data []byte) (*big.Int, error) {
	out, err := flyoverConfigurations.abi.Unpack("calculatePegOutFee", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackCancelDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd602b9fd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (flyoverConfigurations *FlyoverConfigurations) PackCancelDefaultAdminTransfer() []byte {
	enc, err := flyoverConfigurations.abi.Pack("cancelDefaultAdminTransfer")
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackCancelDefaultAdminTransfer() ([]byte, error) {
	return flyoverConfigurations.abi.Pack("cancelDefaultAdminTransfer")
}

// PackChangeDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x649a5ec7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (flyoverConfigurations *FlyoverConfigurations) PackChangeDefaultAdminDelay(newDelay *big.Int) []byte {
	enc, err := flyoverConfigurations.abi.Pack("changeDefaultAdminDelay", newDelay)
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackChangeDefaultAdminDelay(newDelay *big.Int) ([]byte, error) {
	return flyoverConfigurations.abi.Pack("changeDefaultAdminDelay", newDelay)
}

// PackDefaultAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84ef8ffc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function defaultAdmin() view returns(address)
func (flyoverConfigurations *FlyoverConfigurations) PackDefaultAdmin() []byte {
	enc, err := flyoverConfigurations.abi.Pack("defaultAdmin")
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackDefaultAdmin() ([]byte, error) {
	return flyoverConfigurations.abi.Pack("defaultAdmin")
}

// UnpackDefaultAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (flyoverConfigurations *FlyoverConfigurations) UnpackDefaultAdmin(data []byte) (common.Address, error) {
	out, err := flyoverConfigurations.abi.Unpack("defaultAdmin", data)
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
func (flyoverConfigurations *FlyoverConfigurations) PackDefaultAdminDelay() []byte {
	enc, err := flyoverConfigurations.abi.Pack("defaultAdminDelay")
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackDefaultAdminDelay() ([]byte, error) {
	return flyoverConfigurations.abi.Pack("defaultAdminDelay")
}

// UnpackDefaultAdminDelay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (flyoverConfigurations *FlyoverConfigurations) UnpackDefaultAdminDelay(data []byte) (*big.Int, error) {
	out, err := flyoverConfigurations.abi.Unpack("defaultAdminDelay", data)
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
func (flyoverConfigurations *FlyoverConfigurations) PackDefaultAdminDelayIncreaseWait() []byte {
	enc, err := flyoverConfigurations.abi.Pack("defaultAdminDelayIncreaseWait")
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackDefaultAdminDelayIncreaseWait() ([]byte, error) {
	return flyoverConfigurations.abi.Pack("defaultAdminDelayIncreaseWait")
}

// UnpackDefaultAdminDelayIncreaseWait is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (flyoverConfigurations *FlyoverConfigurations) UnpackDefaultAdminDelayIncreaseWait(data []byte) (*big.Int, error) {
	out, err := flyoverConfigurations.abi.Unpack("defaultAdminDelayIncreaseWait", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetPegInConfiguration is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7cd5733c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPegInConfiguration() view returns((uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256))
func (flyoverConfigurations *FlyoverConfigurations) PackGetPegInConfiguration() []byte {
	enc, err := flyoverConfigurations.abi.Pack("getPegInConfiguration")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPegInConfiguration is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7cd5733c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPegInConfiguration() view returns((uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256))
func (flyoverConfigurations *FlyoverConfigurations) TryPackGetPegInConfiguration() ([]byte, error) {
	return flyoverConfigurations.abi.Pack("getPegInConfiguration")
}

// UnpackGetPegInConfiguration is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7cd5733c.
//
// Solidity: function getPegInConfiguration() view returns((uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256))
func (flyoverConfigurations *FlyoverConfigurations) UnpackGetPegInConfiguration(data []byte) (IFlyoverConfigurationsPegConfiguration, error) {
	out, err := flyoverConfigurations.abi.Unpack("getPegInConfiguration", data)
	if err != nil {
		return *new(IFlyoverConfigurationsPegConfiguration), err
	}
	out0 := *abi.ConvertType(out[0], new(IFlyoverConfigurationsPegConfiguration)).(*IFlyoverConfigurationsPegConfiguration)
	return out0, nil
}

// PackGetPegInConfigurationBounds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x79ee8df7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPegInConfigurationBounds() view returns((uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) min, (uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) max)
func (flyoverConfigurations *FlyoverConfigurations) PackGetPegInConfigurationBounds() []byte {
	enc, err := flyoverConfigurations.abi.Pack("getPegInConfigurationBounds")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPegInConfigurationBounds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x79ee8df7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPegInConfigurationBounds() view returns((uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) min, (uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) max)
func (flyoverConfigurations *FlyoverConfigurations) TryPackGetPegInConfigurationBounds() ([]byte, error) {
	return flyoverConfigurations.abi.Pack("getPegInConfigurationBounds")
}

// GetPegInConfigurationBoundsOutput serves as a container for the return parameters of contract
// method GetPegInConfigurationBounds.
type GetPegInConfigurationBoundsOutput struct {
	Min IFlyoverConfigurationsPegConfiguration
	Max IFlyoverConfigurationsPegConfiguration
}

// UnpackGetPegInConfigurationBounds is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x79ee8df7.
//
// Solidity: function getPegInConfigurationBounds() view returns((uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) min, (uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) max)
func (flyoverConfigurations *FlyoverConfigurations) UnpackGetPegInConfigurationBounds(data []byte) (GetPegInConfigurationBoundsOutput, error) {
	out, err := flyoverConfigurations.abi.Unpack("getPegInConfigurationBounds", data)
	outstruct := new(GetPegInConfigurationBoundsOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Min = *abi.ConvertType(out[0], new(IFlyoverConfigurationsPegConfiguration)).(*IFlyoverConfigurationsPegConfiguration)
	outstruct.Max = *abi.ConvertType(out[1], new(IFlyoverConfigurationsPegConfiguration)).(*IFlyoverConfigurationsPegConfiguration)
	return *outstruct, nil
}

// PackGetPegOutConfiguration is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8a8a2124.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPegOutConfiguration() view returns((uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256))
func (flyoverConfigurations *FlyoverConfigurations) PackGetPegOutConfiguration() []byte {
	enc, err := flyoverConfigurations.abi.Pack("getPegOutConfiguration")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPegOutConfiguration is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8a8a2124.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPegOutConfiguration() view returns((uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256))
func (flyoverConfigurations *FlyoverConfigurations) TryPackGetPegOutConfiguration() ([]byte, error) {
	return flyoverConfigurations.abi.Pack("getPegOutConfiguration")
}

// UnpackGetPegOutConfiguration is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8a8a2124.
//
// Solidity: function getPegOutConfiguration() view returns((uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256))
func (flyoverConfigurations *FlyoverConfigurations) UnpackGetPegOutConfiguration(data []byte) (IFlyoverConfigurationsPegConfiguration, error) {
	out, err := flyoverConfigurations.abi.Unpack("getPegOutConfiguration", data)
	if err != nil {
		return *new(IFlyoverConfigurationsPegConfiguration), err
	}
	out0 := *abi.ConvertType(out[0], new(IFlyoverConfigurationsPegConfiguration)).(*IFlyoverConfigurationsPegConfiguration)
	return out0, nil
}

// PackGetPegOutConfigurationBounds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd0f7d23d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPegOutConfigurationBounds() view returns((uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) min, (uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) max)
func (flyoverConfigurations *FlyoverConfigurations) PackGetPegOutConfigurationBounds() []byte {
	enc, err := flyoverConfigurations.abi.Pack("getPegOutConfigurationBounds")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPegOutConfigurationBounds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd0f7d23d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPegOutConfigurationBounds() view returns((uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) min, (uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) max)
func (flyoverConfigurations *FlyoverConfigurations) TryPackGetPegOutConfigurationBounds() ([]byte, error) {
	return flyoverConfigurations.abi.Pack("getPegOutConfigurationBounds")
}

// GetPegOutConfigurationBoundsOutput serves as a container for the return parameters of contract
// method GetPegOutConfigurationBounds.
type GetPegOutConfigurationBoundsOutput struct {
	Min IFlyoverConfigurationsPegConfiguration
	Max IFlyoverConfigurationsPegConfiguration
}

// UnpackGetPegOutConfigurationBounds is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd0f7d23d.
//
// Solidity: function getPegOutConfigurationBounds() view returns((uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) min, (uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) max)
func (flyoverConfigurations *FlyoverConfigurations) UnpackGetPegOutConfigurationBounds(data []byte) (GetPegOutConfigurationBoundsOutput, error) {
	out, err := flyoverConfigurations.abi.Unpack("getPegOutConfigurationBounds", data)
	outstruct := new(GetPegOutConfigurationBoundsOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Min = *abi.ConvertType(out[0], new(IFlyoverConfigurationsPegConfiguration)).(*IFlyoverConfigurationsPegConfiguration)
	outstruct.Max = *abi.ConvertType(out[1], new(IFlyoverConfigurationsPegConfiguration)).(*IFlyoverConfigurationsPegConfiguration)
	return *outstruct, nil
}

// PackGetPendingChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd9c4f6f3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPendingChange(uint8 flow, uint8 field) view returns(uint256 value, uint256 eta)
func (flyoverConfigurations *FlyoverConfigurations) PackGetPendingChange(flow uint8, field uint8) []byte {
	enc, err := flyoverConfigurations.abi.Pack("getPendingChange", flow, field)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPendingChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd9c4f6f3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPendingChange(uint8 flow, uint8 field) view returns(uint256 value, uint256 eta)
func (flyoverConfigurations *FlyoverConfigurations) TryPackGetPendingChange(flow uint8, field uint8) ([]byte, error) {
	return flyoverConfigurations.abi.Pack("getPendingChange", flow, field)
}

// GetPendingChangeOutput serves as a container for the return parameters of contract
// method GetPendingChange.
type GetPendingChangeOutput struct {
	Value *big.Int
	Eta   *big.Int
}

// UnpackGetPendingChange is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd9c4f6f3.
//
// Solidity: function getPendingChange(uint8 flow, uint8 field) view returns(uint256 value, uint256 eta)
func (flyoverConfigurations *FlyoverConfigurations) UnpackGetPendingChange(data []byte) (GetPendingChangeOutput, error) {
	out, err := flyoverConfigurations.abi.Unpack("getPendingChange", data)
	outstruct := new(GetPendingChangeOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Value = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Eta = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackGetRequiredPegInConfirmations is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf5fca2fd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRequiredPegInConfirmations(uint256 amount) view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) PackGetRequiredPegInConfirmations(amount *big.Int) []byte {
	enc, err := flyoverConfigurations.abi.Pack("getRequiredPegInConfirmations", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRequiredPegInConfirmations is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf5fca2fd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRequiredPegInConfirmations(uint256 amount) view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) TryPackGetRequiredPegInConfirmations(amount *big.Int) ([]byte, error) {
	return flyoverConfigurations.abi.Pack("getRequiredPegInConfirmations", amount)
}

// UnpackGetRequiredPegInConfirmations is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf5fca2fd.
//
// Solidity: function getRequiredPegInConfirmations(uint256 amount) view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) UnpackGetRequiredPegInConfirmations(data []byte) (*big.Int, error) {
	out, err := flyoverConfigurations.abi.Unpack("getRequiredPegInConfirmations", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetRequiredPegOutConfirmations is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2ce3bc7e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRequiredPegOutConfirmations(uint256 amount) view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) PackGetRequiredPegOutConfirmations(amount *big.Int) []byte {
	enc, err := flyoverConfigurations.abi.Pack("getRequiredPegOutConfirmations", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRequiredPegOutConfirmations is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2ce3bc7e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRequiredPegOutConfirmations(uint256 amount) view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) TryPackGetRequiredPegOutConfirmations(amount *big.Int) ([]byte, error) {
	return flyoverConfigurations.abi.Pack("getRequiredPegOutConfirmations", amount)
}

// UnpackGetRequiredPegOutConfirmations is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2ce3bc7e.
//
// Solidity: function getRequiredPegOutConfirmations(uint256 amount) view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) UnpackGetRequiredPegOutConfirmations(data []byte) (*big.Int, error) {
	out, err := flyoverConfigurations.abi.Unpack("getRequiredPegOutConfirmations", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (flyoverConfigurations *FlyoverConfigurations) PackGetRoleAdmin(role [32]byte) []byte {
	enc, err := flyoverConfigurations.abi.Pack("getRoleAdmin", role)
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackGetRoleAdmin(role [32]byte) ([]byte, error) {
	return flyoverConfigurations.abi.Pack("getRoleAdmin", role)
}

// UnpackGetRoleAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (flyoverConfigurations *FlyoverConfigurations) UnpackGetRoleAdmin(data []byte) ([32]byte, error) {
	out, err := flyoverConfigurations.abi.Unpack("getRoleAdmin", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGetTimelockDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x481c42a2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getTimelockDelay() view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) PackGetTimelockDelay() []byte {
	enc, err := flyoverConfigurations.abi.Pack("getTimelockDelay")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetTimelockDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x481c42a2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getTimelockDelay() view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) TryPackGetTimelockDelay() ([]byte, error) {
	return flyoverConfigurations.abi.Pack("getTimelockDelay")
}

// UnpackGetTimelockDelay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x481c42a2.
//
// Solidity: function getTimelockDelay() view returns(uint256)
func (flyoverConfigurations *FlyoverConfigurations) UnpackGetTimelockDelay(data []byte) (*big.Int, error) {
	out, err := flyoverConfigurations.abi.Unpack("getTimelockDelay", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (flyoverConfigurations *FlyoverConfigurations) PackGrantRole(role [32]byte, account common.Address) []byte {
	enc, err := flyoverConfigurations.abi.Pack("grantRole", role, account)
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackGrantRole(role [32]byte, account common.Address) ([]byte, error) {
	return flyoverConfigurations.abi.Pack("grantRole", role, account)
}

// PackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (flyoverConfigurations *FlyoverConfigurations) PackHasRole(role [32]byte, account common.Address) []byte {
	enc, err := flyoverConfigurations.abi.Pack("hasRole", role, account)
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackHasRole(role [32]byte, account common.Address) ([]byte, error) {
	return flyoverConfigurations.abi.Pack("hasRole", role, account)
}

// UnpackHasRole is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (flyoverConfigurations *FlyoverConfigurations) UnpackHasRole(data []byte) (bool, error) {
	out, err := flyoverConfigurations.abi.Unpack("hasRole", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x03cce06c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function initialize(address defaultAdmin, uint48 initialDelay, uint256 timelockDelay, (uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) pegInConfig, (uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) pegOutConfig, (uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) pegInMin, (uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) pegInMax, (uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) pegOutMin, (uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) pegOutMax) returns()
func (flyoverConfigurations *FlyoverConfigurations) PackInitialize(defaultAdmin common.Address, initialDelay *big.Int, timelockDelay *big.Int, pegInConfig IFlyoverConfigurationsPegConfiguration, pegOutConfig IFlyoverConfigurationsPegConfiguration, pegInMin IFlyoverConfigurationsPegConfiguration, pegInMax IFlyoverConfigurationsPegConfiguration, pegOutMin IFlyoverConfigurationsPegConfiguration, pegOutMax IFlyoverConfigurationsPegConfiguration) []byte {
	enc, err := flyoverConfigurations.abi.Pack("initialize", defaultAdmin, initialDelay, timelockDelay, pegInConfig, pegOutConfig, pegInMin, pegInMax, pegOutMin, pegOutMax)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x03cce06c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function initialize(address defaultAdmin, uint48 initialDelay, uint256 timelockDelay, (uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) pegInConfig, (uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) pegOutConfig, (uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) pegInMin, (uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) pegInMax, (uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) pegOutMin, (uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256) pegOutMax) returns()
func (flyoverConfigurations *FlyoverConfigurations) TryPackInitialize(defaultAdmin common.Address, initialDelay *big.Int, timelockDelay *big.Int, pegInConfig IFlyoverConfigurationsPegConfiguration, pegOutConfig IFlyoverConfigurationsPegConfiguration, pegInMin IFlyoverConfigurationsPegConfiguration, pegInMax IFlyoverConfigurationsPegConfiguration, pegOutMin IFlyoverConfigurationsPegConfiguration, pegOutMax IFlyoverConfigurationsPegConfiguration) ([]byte, error) {
	return flyoverConfigurations.abi.Pack("initialize", defaultAdmin, initialDelay, timelockDelay, pegInConfig, pegOutConfig, pegInMin, pegInMax, pegOutMin, pegOutMax)
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (flyoverConfigurations *FlyoverConfigurations) PackOwner() []byte {
	enc, err := flyoverConfigurations.abi.Pack("owner")
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackOwner() ([]byte, error) {
	return flyoverConfigurations.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (flyoverConfigurations *FlyoverConfigurations) UnpackOwner(data []byte) (common.Address, error) {
	out, err := flyoverConfigurations.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPendingDefaultAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf6eefb7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (flyoverConfigurations *FlyoverConfigurations) PackPendingDefaultAdmin() []byte {
	enc, err := flyoverConfigurations.abi.Pack("pendingDefaultAdmin")
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackPendingDefaultAdmin() ([]byte, error) {
	return flyoverConfigurations.abi.Pack("pendingDefaultAdmin")
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
func (flyoverConfigurations *FlyoverConfigurations) UnpackPendingDefaultAdmin(data []byte) (PendingDefaultAdminOutput, error) {
	out, err := flyoverConfigurations.abi.Unpack("pendingDefaultAdmin", data)
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
func (flyoverConfigurations *FlyoverConfigurations) PackPendingDefaultAdminDelay() []byte {
	enc, err := flyoverConfigurations.abi.Pack("pendingDefaultAdminDelay")
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackPendingDefaultAdminDelay() ([]byte, error) {
	return flyoverConfigurations.abi.Pack("pendingDefaultAdminDelay")
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
func (flyoverConfigurations *FlyoverConfigurations) UnpackPendingDefaultAdminDelay(data []byte) (PendingDefaultAdminDelayOutput, error) {
	out, err := flyoverConfigurations.abi.Unpack("pendingDefaultAdminDelay", data)
	outstruct := new(PendingDefaultAdminDelayOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.NewDelay = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Schedule = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackQueueChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x656213e3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function queueChange(uint8 flow, uint8 field, uint256 value) returns()
func (flyoverConfigurations *FlyoverConfigurations) PackQueueChange(flow uint8, field uint8, value *big.Int) []byte {
	enc, err := flyoverConfigurations.abi.Pack("queueChange", flow, field, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackQueueChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x656213e3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function queueChange(uint8 flow, uint8 field, uint256 value) returns()
func (flyoverConfigurations *FlyoverConfigurations) TryPackQueueChange(flow uint8, field uint8, value *big.Int) ([]byte, error) {
	return flyoverConfigurations.abi.Pack("queueChange", flow, field, value)
}

// PackQueueTiersChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x80a6c46d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function queueTiersChange(uint8 flow, (uint256,uint256)[] tiers) returns()
func (flyoverConfigurations *FlyoverConfigurations) PackQueueTiersChange(flow uint8, tiers []IFlyoverConfigurationsConfirmationTier) []byte {
	enc, err := flyoverConfigurations.abi.Pack("queueTiersChange", flow, tiers)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackQueueTiersChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x80a6c46d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function queueTiersChange(uint8 flow, (uint256,uint256)[] tiers) returns()
func (flyoverConfigurations *FlyoverConfigurations) TryPackQueueTiersChange(flow uint8, tiers []IFlyoverConfigurationsConfirmationTier) ([]byte, error) {
	return flyoverConfigurations.abi.Pack("queueTiersChange", flow, tiers)
}

// PackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (flyoverConfigurations *FlyoverConfigurations) PackRenounceRole(role [32]byte, account common.Address) []byte {
	enc, err := flyoverConfigurations.abi.Pack("renounceRole", role, account)
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackRenounceRole(role [32]byte, account common.Address) ([]byte, error) {
	return flyoverConfigurations.abi.Pack("renounceRole", role, account)
}

// PackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (flyoverConfigurations *FlyoverConfigurations) PackRevokeRole(role [32]byte, account common.Address) []byte {
	enc, err := flyoverConfigurations.abi.Pack("revokeRole", role, account)
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackRevokeRole(role [32]byte, account common.Address) ([]byte, error) {
	return flyoverConfigurations.abi.Pack("revokeRole", role, account)
}

// PackRollbackDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0aa6220b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (flyoverConfigurations *FlyoverConfigurations) PackRollbackDefaultAdminDelay() []byte {
	enc, err := flyoverConfigurations.abi.Pack("rollbackDefaultAdminDelay")
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackRollbackDefaultAdminDelay() ([]byte, error) {
	return flyoverConfigurations.abi.Pack("rollbackDefaultAdminDelay")
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (flyoverConfigurations *FlyoverConfigurations) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := flyoverConfigurations.abi.Pack("supportsInterface", interfaceId)
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
func (flyoverConfigurations *FlyoverConfigurations) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return flyoverConfigurations.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (flyoverConfigurations *FlyoverConfigurations) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := flyoverConfigurations.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// FlyoverConfigurationsDefaultAdminDelayChangeCanceled represents a DefaultAdminDelayChangeCanceled event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsDefaultAdminDelayChangeCanceled struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsDefaultAdminDelayChangeCanceledEventName = "DefaultAdminDelayChangeCanceled"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsDefaultAdminDelayChangeCanceled) ContractEventName() string {
	return FlyoverConfigurationsDefaultAdminDelayChangeCanceledEventName
}

// UnpackDefaultAdminDelayChangeCanceledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (flyoverConfigurations *FlyoverConfigurations) UnpackDefaultAdminDelayChangeCanceledEvent(log *types.Log) (*FlyoverConfigurationsDefaultAdminDelayChangeCanceled, error) {
	event := "DefaultAdminDelayChangeCanceled"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsDefaultAdminDelayChangeCanceled)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsDefaultAdminDelayChangeScheduled represents a DefaultAdminDelayChangeScheduled event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsDefaultAdminDelayChangeScheduled struct {
	NewDelay       *big.Int
	EffectSchedule *big.Int
	Raw            *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsDefaultAdminDelayChangeScheduledEventName = "DefaultAdminDelayChangeScheduled"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsDefaultAdminDelayChangeScheduled) ContractEventName() string {
	return FlyoverConfigurationsDefaultAdminDelayChangeScheduledEventName
}

// UnpackDefaultAdminDelayChangeScheduledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (flyoverConfigurations *FlyoverConfigurations) UnpackDefaultAdminDelayChangeScheduledEvent(log *types.Log) (*FlyoverConfigurationsDefaultAdminDelayChangeScheduled, error) {
	event := "DefaultAdminDelayChangeScheduled"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsDefaultAdminDelayChangeScheduled)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsDefaultAdminTransferCanceled represents a DefaultAdminTransferCanceled event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsDefaultAdminTransferCanceled struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsDefaultAdminTransferCanceledEventName = "DefaultAdminTransferCanceled"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsDefaultAdminTransferCanceled) ContractEventName() string {
	return FlyoverConfigurationsDefaultAdminTransferCanceledEventName
}

// UnpackDefaultAdminTransferCanceledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminTransferCanceled()
func (flyoverConfigurations *FlyoverConfigurations) UnpackDefaultAdminTransferCanceledEvent(log *types.Log) (*FlyoverConfigurationsDefaultAdminTransferCanceled, error) {
	event := "DefaultAdminTransferCanceled"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsDefaultAdminTransferCanceled)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsDefaultAdminTransferScheduled represents a DefaultAdminTransferScheduled event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsDefaultAdminTransferScheduled struct {
	NewAdmin       common.Address
	AcceptSchedule *big.Int
	Raw            *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsDefaultAdminTransferScheduledEventName = "DefaultAdminTransferScheduled"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsDefaultAdminTransferScheduled) ContractEventName() string {
	return FlyoverConfigurationsDefaultAdminTransferScheduledEventName
}

// UnpackDefaultAdminTransferScheduledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (flyoverConfigurations *FlyoverConfigurations) UnpackDefaultAdminTransferScheduledEvent(log *types.Log) (*FlyoverConfigurationsDefaultAdminTransferScheduled, error) {
	event := "DefaultAdminTransferScheduled"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsDefaultAdminTransferScheduled)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsInitialized represents a Initialized event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsInitialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsInitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsInitialized) ContractEventName() string {
	return FlyoverConfigurationsInitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (flyoverConfigurations *FlyoverConfigurations) UnpackInitializedEvent(log *types.Log) (*FlyoverConfigurationsInitialized, error) {
	event := "Initialized"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsInitialized)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegInCallTimeChanged represents a PegInCallTimeChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegInCallTimeChanged struct {
	OldValue *big.Int
	NewValue *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegInCallTimeChangedEventName = "PegInCallTimeChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegInCallTimeChanged) ContractEventName() string {
	return FlyoverConfigurationsPegInCallTimeChangedEventName
}

// UnpackPegInCallTimeChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegInCallTimeChanged(uint256 oldValue, uint256 newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegInCallTimeChangedEvent(log *types.Log) (*FlyoverConfigurationsPegInCallTimeChanged, error) {
	event := "PegInCallTimeChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegInCallTimeChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegInConfirmationTiersChanged represents a PegInConfirmationTiersChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegInConfirmationTiersChanged struct {
	OldValue []IFlyoverConfigurationsConfirmationTier
	NewValue []IFlyoverConfigurationsConfirmationTier
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegInConfirmationTiersChangedEventName = "PegInConfirmationTiersChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegInConfirmationTiersChanged) ContractEventName() string {
	return FlyoverConfigurationsPegInConfirmationTiersChangedEventName
}

// UnpackPegInConfirmationTiersChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegInConfirmationTiersChanged((uint256,uint256)[] oldValue, (uint256,uint256)[] newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegInConfirmationTiersChangedEvent(log *types.Log) (*FlyoverConfigurationsPegInConfirmationTiersChanged, error) {
	event := "PegInConfirmationTiersChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegInConfirmationTiersChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegInDeliveryGraceChanged represents a PegInDeliveryGraceChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegInDeliveryGraceChanged struct {
	OldValue *big.Int
	NewValue *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegInDeliveryGraceChangedEventName = "PegInDeliveryGraceChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegInDeliveryGraceChanged) ContractEventName() string {
	return FlyoverConfigurationsPegInDeliveryGraceChangedEventName
}

// UnpackPegInDeliveryGraceChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegInDeliveryGraceChanged(uint256 oldValue, uint256 newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegInDeliveryGraceChangedEvent(log *types.Log) (*FlyoverConfigurationsPegInDeliveryGraceChanged, error) {
	event := "PegInDeliveryGraceChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegInDeliveryGraceChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegInExpireBlocksChanged represents a PegInExpireBlocksChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegInExpireBlocksChanged struct {
	OldValue *big.Int
	NewValue *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegInExpireBlocksChangedEventName = "PegInExpireBlocksChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegInExpireBlocksChanged) ContractEventName() string {
	return FlyoverConfigurationsPegInExpireBlocksChangedEventName
}

// UnpackPegInExpireBlocksChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegInExpireBlocksChanged(uint256 oldValue, uint256 newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegInExpireBlocksChangedEvent(log *types.Log) (*FlyoverConfigurationsPegInExpireBlocksChanged, error) {
	event := "PegInExpireBlocksChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegInExpireBlocksChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegInExpireTimeChanged represents a PegInExpireTimeChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegInExpireTimeChanged struct {
	OldValue *big.Int
	NewValue *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegInExpireTimeChangedEventName = "PegInExpireTimeChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegInExpireTimeChanged) ContractEventName() string {
	return FlyoverConfigurationsPegInExpireTimeChangedEventName
}

// UnpackPegInExpireTimeChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegInExpireTimeChanged(uint256 oldValue, uint256 newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegInExpireTimeChangedEvent(log *types.Log) (*FlyoverConfigurationsPegInExpireTimeChanged, error) {
	event := "PegInExpireTimeChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegInExpireTimeChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegInFixedFeeChanged represents a PegInFixedFeeChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegInFixedFeeChanged struct {
	OldValue *big.Int
	NewValue *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegInFixedFeeChangedEventName = "PegInFixedFeeChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegInFixedFeeChanged) ContractEventName() string {
	return FlyoverConfigurationsPegInFixedFeeChangedEventName
}

// UnpackPegInFixedFeeChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegInFixedFeeChanged(uint256 oldValue, uint256 newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegInFixedFeeChangedEvent(log *types.Log) (*FlyoverConfigurationsPegInFixedFeeChanged, error) {
	event := "PegInFixedFeeChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegInFixedFeeChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegInMaxAmountChanged represents a PegInMaxAmountChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegInMaxAmountChanged struct {
	OldValue *big.Int
	NewValue *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegInMaxAmountChangedEventName = "PegInMaxAmountChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegInMaxAmountChanged) ContractEventName() string {
	return FlyoverConfigurationsPegInMaxAmountChangedEventName
}

// UnpackPegInMaxAmountChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegInMaxAmountChanged(uint256 oldValue, uint256 newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegInMaxAmountChangedEvent(log *types.Log) (*FlyoverConfigurationsPegInMaxAmountChanged, error) {
	event := "PegInMaxAmountChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegInMaxAmountChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegInMinAmountChanged represents a PegInMinAmountChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegInMinAmountChanged struct {
	OldValue *big.Int
	NewValue *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegInMinAmountChangedEventName = "PegInMinAmountChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegInMinAmountChanged) ContractEventName() string {
	return FlyoverConfigurationsPegInMinAmountChangedEventName
}

// UnpackPegInMinAmountChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegInMinAmountChanged(uint256 oldValue, uint256 newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegInMinAmountChangedEvent(log *types.Log) (*FlyoverConfigurationsPegInMinAmountChanged, error) {
	event := "PegInMinAmountChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegInMinAmountChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegInPenaltyFeeChanged represents a PegInPenaltyFeeChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegInPenaltyFeeChanged struct {
	OldValue *big.Int
	NewValue *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegInPenaltyFeeChangedEventName = "PegInPenaltyFeeChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegInPenaltyFeeChanged) ContractEventName() string {
	return FlyoverConfigurationsPegInPenaltyFeeChangedEventName
}

// UnpackPegInPenaltyFeeChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegInPenaltyFeeChanged(uint256 oldValue, uint256 newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegInPenaltyFeeChangedEvent(log *types.Log) (*FlyoverConfigurationsPegInPenaltyFeeChanged, error) {
	event := "PegInPenaltyFeeChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegInPenaltyFeeChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegInPercentageFeeChanged represents a PegInPercentageFeeChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegInPercentageFeeChanged struct {
	OldValue *big.Int
	NewValue *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegInPercentageFeeChangedEventName = "PegInPercentageFeeChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegInPercentageFeeChanged) ContractEventName() string {
	return FlyoverConfigurationsPegInPercentageFeeChangedEventName
}

// UnpackPegInPercentageFeeChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegInPercentageFeeChanged(uint256 oldValue, uint256 newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegInPercentageFeeChangedEvent(log *types.Log) (*FlyoverConfigurationsPegInPercentageFeeChanged, error) {
	event := "PegInPercentageFeeChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegInPercentageFeeChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegOutCallTimeChanged represents a PegOutCallTimeChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegOutCallTimeChanged struct {
	OldValue *big.Int
	NewValue *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegOutCallTimeChangedEventName = "PegOutCallTimeChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegOutCallTimeChanged) ContractEventName() string {
	return FlyoverConfigurationsPegOutCallTimeChangedEventName
}

// UnpackPegOutCallTimeChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutCallTimeChanged(uint256 oldValue, uint256 newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegOutCallTimeChangedEvent(log *types.Log) (*FlyoverConfigurationsPegOutCallTimeChanged, error) {
	event := "PegOutCallTimeChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegOutCallTimeChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegOutConfirmationTiersChanged represents a PegOutConfirmationTiersChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegOutConfirmationTiersChanged struct {
	OldValue []IFlyoverConfigurationsConfirmationTier
	NewValue []IFlyoverConfigurationsConfirmationTier
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegOutConfirmationTiersChangedEventName = "PegOutConfirmationTiersChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegOutConfirmationTiersChanged) ContractEventName() string {
	return FlyoverConfigurationsPegOutConfirmationTiersChangedEventName
}

// UnpackPegOutConfirmationTiersChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutConfirmationTiersChanged((uint256,uint256)[] oldValue, (uint256,uint256)[] newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegOutConfirmationTiersChangedEvent(log *types.Log) (*FlyoverConfigurationsPegOutConfirmationTiersChanged, error) {
	event := "PegOutConfirmationTiersChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegOutConfirmationTiersChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegOutDeliveryGraceChanged represents a PegOutDeliveryGraceChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegOutDeliveryGraceChanged struct {
	OldValue *big.Int
	NewValue *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegOutDeliveryGraceChangedEventName = "PegOutDeliveryGraceChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegOutDeliveryGraceChanged) ContractEventName() string {
	return FlyoverConfigurationsPegOutDeliveryGraceChangedEventName
}

// UnpackPegOutDeliveryGraceChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutDeliveryGraceChanged(uint256 oldValue, uint256 newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegOutDeliveryGraceChangedEvent(log *types.Log) (*FlyoverConfigurationsPegOutDeliveryGraceChanged, error) {
	event := "PegOutDeliveryGraceChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegOutDeliveryGraceChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegOutExpireBlocksChanged represents a PegOutExpireBlocksChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegOutExpireBlocksChanged struct {
	OldValue *big.Int
	NewValue *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegOutExpireBlocksChangedEventName = "PegOutExpireBlocksChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegOutExpireBlocksChanged) ContractEventName() string {
	return FlyoverConfigurationsPegOutExpireBlocksChangedEventName
}

// UnpackPegOutExpireBlocksChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutExpireBlocksChanged(uint256 oldValue, uint256 newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegOutExpireBlocksChangedEvent(log *types.Log) (*FlyoverConfigurationsPegOutExpireBlocksChanged, error) {
	event := "PegOutExpireBlocksChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegOutExpireBlocksChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegOutExpireTimeChanged represents a PegOutExpireTimeChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegOutExpireTimeChanged struct {
	OldValue *big.Int
	NewValue *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegOutExpireTimeChangedEventName = "PegOutExpireTimeChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegOutExpireTimeChanged) ContractEventName() string {
	return FlyoverConfigurationsPegOutExpireTimeChangedEventName
}

// UnpackPegOutExpireTimeChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutExpireTimeChanged(uint256 oldValue, uint256 newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegOutExpireTimeChangedEvent(log *types.Log) (*FlyoverConfigurationsPegOutExpireTimeChanged, error) {
	event := "PegOutExpireTimeChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegOutExpireTimeChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegOutFixedFeeChanged represents a PegOutFixedFeeChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegOutFixedFeeChanged struct {
	OldValue *big.Int
	NewValue *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegOutFixedFeeChangedEventName = "PegOutFixedFeeChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegOutFixedFeeChanged) ContractEventName() string {
	return FlyoverConfigurationsPegOutFixedFeeChangedEventName
}

// UnpackPegOutFixedFeeChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutFixedFeeChanged(uint256 oldValue, uint256 newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegOutFixedFeeChangedEvent(log *types.Log) (*FlyoverConfigurationsPegOutFixedFeeChanged, error) {
	event := "PegOutFixedFeeChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegOutFixedFeeChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegOutMaxAmountChanged represents a PegOutMaxAmountChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegOutMaxAmountChanged struct {
	OldValue *big.Int
	NewValue *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegOutMaxAmountChangedEventName = "PegOutMaxAmountChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegOutMaxAmountChanged) ContractEventName() string {
	return FlyoverConfigurationsPegOutMaxAmountChangedEventName
}

// UnpackPegOutMaxAmountChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutMaxAmountChanged(uint256 oldValue, uint256 newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegOutMaxAmountChangedEvent(log *types.Log) (*FlyoverConfigurationsPegOutMaxAmountChanged, error) {
	event := "PegOutMaxAmountChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegOutMaxAmountChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegOutMinAmountChanged represents a PegOutMinAmountChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegOutMinAmountChanged struct {
	OldValue *big.Int
	NewValue *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegOutMinAmountChangedEventName = "PegOutMinAmountChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegOutMinAmountChanged) ContractEventName() string {
	return FlyoverConfigurationsPegOutMinAmountChangedEventName
}

// UnpackPegOutMinAmountChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutMinAmountChanged(uint256 oldValue, uint256 newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegOutMinAmountChangedEvent(log *types.Log) (*FlyoverConfigurationsPegOutMinAmountChanged, error) {
	event := "PegOutMinAmountChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegOutMinAmountChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegOutPenaltyFeeChanged represents a PegOutPenaltyFeeChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegOutPenaltyFeeChanged struct {
	OldValue *big.Int
	NewValue *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegOutPenaltyFeeChangedEventName = "PegOutPenaltyFeeChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegOutPenaltyFeeChanged) ContractEventName() string {
	return FlyoverConfigurationsPegOutPenaltyFeeChangedEventName
}

// UnpackPegOutPenaltyFeeChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutPenaltyFeeChanged(uint256 oldValue, uint256 newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegOutPenaltyFeeChangedEvent(log *types.Log) (*FlyoverConfigurationsPegOutPenaltyFeeChanged, error) {
	event := "PegOutPenaltyFeeChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegOutPenaltyFeeChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsPegOutPercentageFeeChanged represents a PegOutPercentageFeeChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPegOutPercentageFeeChanged struct {
	OldValue *big.Int
	NewValue *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsPegOutPercentageFeeChangedEventName = "PegOutPercentageFeeChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsPegOutPercentageFeeChanged) ContractEventName() string {
	return FlyoverConfigurationsPegOutPercentageFeeChangedEventName
}

// UnpackPegOutPercentageFeeChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutPercentageFeeChanged(uint256 oldValue, uint256 newValue)
func (flyoverConfigurations *FlyoverConfigurations) UnpackPegOutPercentageFeeChangedEvent(log *types.Log) (*FlyoverConfigurationsPegOutPercentageFeeChanged, error) {
	event := "PegOutPercentageFeeChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsPegOutPercentageFeeChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsRoleAdminChanged represents a RoleAdminChanged event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsRoleAdminChangedEventName = "RoleAdminChanged"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsRoleAdminChanged) ContractEventName() string {
	return FlyoverConfigurationsRoleAdminChangedEventName
}

// UnpackRoleAdminChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (flyoverConfigurations *FlyoverConfigurations) UnpackRoleAdminChangedEvent(log *types.Log) (*FlyoverConfigurationsRoleAdminChanged, error) {
	event := "RoleAdminChanged"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsRoleAdminChanged)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsRoleGranted represents a RoleGranted event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsRoleGrantedEventName = "RoleGranted"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsRoleGranted) ContractEventName() string {
	return FlyoverConfigurationsRoleGrantedEventName
}

// UnpackRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (flyoverConfigurations *FlyoverConfigurations) UnpackRoleGrantedEvent(log *types.Log) (*FlyoverConfigurationsRoleGranted, error) {
	event := "RoleGranted"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsRoleGranted)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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

// FlyoverConfigurationsRoleRevoked represents a RoleRevoked event raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const FlyoverConfigurationsRoleRevokedEventName = "RoleRevoked"

// ContractEventName returns the user-defined event name.
func (FlyoverConfigurationsRoleRevoked) ContractEventName() string {
	return FlyoverConfigurationsRoleRevokedEventName
}

// UnpackRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (flyoverConfigurations *FlyoverConfigurations) UnpackRoleRevokedEvent(log *types.Log) (*FlyoverConfigurationsRoleRevoked, error) {
	event := "RoleRevoked"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverConfigurations.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverConfigurationsRoleRevoked)
	if len(log.Data) > 0 {
		if err := flyoverConfigurations.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverConfigurations.abi.Events[event].Inputs {
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
func (flyoverConfigurations *FlyoverConfigurations) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], flyoverConfigurations.abi.Errors["AccessControlBadConfirmation"].ID.Bytes()[:4]) {
		return flyoverConfigurations.UnpackAccessControlBadConfirmationError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverConfigurations.abi.Errors["AccessControlEnforcedDefaultAdminDelay"].ID.Bytes()[:4]) {
		return flyoverConfigurations.UnpackAccessControlEnforcedDefaultAdminDelayError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverConfigurations.abi.Errors["AccessControlEnforcedDefaultAdminRules"].ID.Bytes()[:4]) {
		return flyoverConfigurations.UnpackAccessControlEnforcedDefaultAdminRulesError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverConfigurations.abi.Errors["AccessControlInvalidDefaultAdmin"].ID.Bytes()[:4]) {
		return flyoverConfigurations.UnpackAccessControlInvalidDefaultAdminError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverConfigurations.abi.Errors["AccessControlUnauthorizedAccount"].ID.Bytes()[:4]) {
		return flyoverConfigurations.UnpackAccessControlUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverConfigurations.abi.Errors["EmptyTiers"].ID.Bytes()[:4]) {
		return flyoverConfigurations.UnpackEmptyTiersError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverConfigurations.abi.Errors["ExpireTimeNotAfterCallTime"].ID.Bytes()[:4]) {
		return flyoverConfigurations.UnpackExpireTimeNotAfterCallTimeError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverConfigurations.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return flyoverConfigurations.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverConfigurations.abi.Errors["NoQueuedChange"].ID.Bytes()[:4]) {
		return flyoverConfigurations.UnpackNoQueuedChangeError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverConfigurations.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return flyoverConfigurations.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverConfigurations.abi.Errors["OutOfBounds"].ID.Bytes()[:4]) {
		return flyoverConfigurations.UnpackOutOfBoundsError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverConfigurations.abi.Errors["PaymentNotAllowed"].ID.Bytes()[:4]) {
		return flyoverConfigurations.UnpackPaymentNotAllowedError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverConfigurations.abi.Errors["SafeCastOverflowedUintDowncast"].ID.Bytes()[:4]) {
		return flyoverConfigurations.UnpackSafeCastOverflowedUintDowncastError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverConfigurations.abi.Errors["TiersNotAscending"].ID.Bytes()[:4]) {
		return flyoverConfigurations.UnpackTiersNotAscendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverConfigurations.abi.Errors["TimelockNotElapsed"].ID.Bytes()[:4]) {
		return flyoverConfigurations.UnpackTimelockNotElapsedError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// FlyoverConfigurationsAccessControlBadConfirmation represents a AccessControlBadConfirmation error raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsAccessControlBadConfirmation struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlBadConfirmation()
func FlyoverConfigurationsAccessControlBadConfirmationErrorID() common.Hash {
	return common.HexToHash("0x6697b23232a647058342c0724fe7c415cab25915b54e5dbc03f233173d37b41c")
}

// UnpackAccessControlBadConfirmationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlBadConfirmation()
func (flyoverConfigurations *FlyoverConfigurations) UnpackAccessControlBadConfirmationError(raw []byte) (*FlyoverConfigurationsAccessControlBadConfirmation, error) {
	out := new(FlyoverConfigurationsAccessControlBadConfirmation)
	if err := flyoverConfigurations.abi.UnpackIntoInterface(out, "AccessControlBadConfirmation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverConfigurationsAccessControlEnforcedDefaultAdminDelay represents a AccessControlEnforcedDefaultAdminDelay error raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsAccessControlEnforcedDefaultAdminDelay struct {
	Schedule *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlEnforcedDefaultAdminDelay(uint48 schedule)
func FlyoverConfigurationsAccessControlEnforcedDefaultAdminDelayErrorID() common.Hash {
	return common.HexToHash("0x19ca5ebb8fb33f00e502c9392eddab1501674629178bf69b853cf037aaf4bb5d")
}

// UnpackAccessControlEnforcedDefaultAdminDelayError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlEnforcedDefaultAdminDelay(uint48 schedule)
func (flyoverConfigurations *FlyoverConfigurations) UnpackAccessControlEnforcedDefaultAdminDelayError(raw []byte) (*FlyoverConfigurationsAccessControlEnforcedDefaultAdminDelay, error) {
	out := new(FlyoverConfigurationsAccessControlEnforcedDefaultAdminDelay)
	if err := flyoverConfigurations.abi.UnpackIntoInterface(out, "AccessControlEnforcedDefaultAdminDelay", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverConfigurationsAccessControlEnforcedDefaultAdminRules represents a AccessControlEnforcedDefaultAdminRules error raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsAccessControlEnforcedDefaultAdminRules struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlEnforcedDefaultAdminRules()
func FlyoverConfigurationsAccessControlEnforcedDefaultAdminRulesErrorID() common.Hash {
	return common.HexToHash("0x3fc3c27ae3db78c81b8f6e685172134623efa268ee8cd8d54be38ad2a74fc13b")
}

// UnpackAccessControlEnforcedDefaultAdminRulesError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlEnforcedDefaultAdminRules()
func (flyoverConfigurations *FlyoverConfigurations) UnpackAccessControlEnforcedDefaultAdminRulesError(raw []byte) (*FlyoverConfigurationsAccessControlEnforcedDefaultAdminRules, error) {
	out := new(FlyoverConfigurationsAccessControlEnforcedDefaultAdminRules)
	if err := flyoverConfigurations.abi.UnpackIntoInterface(out, "AccessControlEnforcedDefaultAdminRules", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverConfigurationsAccessControlInvalidDefaultAdmin represents a AccessControlInvalidDefaultAdmin error raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsAccessControlInvalidDefaultAdmin struct {
	DefaultAdmin common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlInvalidDefaultAdmin(address defaultAdmin)
func FlyoverConfigurationsAccessControlInvalidDefaultAdminErrorID() common.Hash {
	return common.HexToHash("0xc22c8022f2a840d6b6a9f113407715f5bbd4e88c1b0dd9434dc00700ba609ed4")
}

// UnpackAccessControlInvalidDefaultAdminError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlInvalidDefaultAdmin(address defaultAdmin)
func (flyoverConfigurations *FlyoverConfigurations) UnpackAccessControlInvalidDefaultAdminError(raw []byte) (*FlyoverConfigurationsAccessControlInvalidDefaultAdmin, error) {
	out := new(FlyoverConfigurationsAccessControlInvalidDefaultAdmin)
	if err := flyoverConfigurations.abi.UnpackIntoInterface(out, "AccessControlInvalidDefaultAdmin", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverConfigurationsAccessControlUnauthorizedAccount represents a AccessControlUnauthorizedAccount error raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsAccessControlUnauthorizedAccount struct {
	Account    common.Address
	NeededRole [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func FlyoverConfigurationsAccessControlUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0xe2517d3fbfae6f8515ef5ff1ccedc3933ab0cbbda0b492c06eb54ad10ef03b3e")
}

// UnpackAccessControlUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func (flyoverConfigurations *FlyoverConfigurations) UnpackAccessControlUnauthorizedAccountError(raw []byte) (*FlyoverConfigurationsAccessControlUnauthorizedAccount, error) {
	out := new(FlyoverConfigurationsAccessControlUnauthorizedAccount)
	if err := flyoverConfigurations.abi.UnpackIntoInterface(out, "AccessControlUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverConfigurationsEmptyTiers represents a EmptyTiers error raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsEmptyTiers struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EmptyTiers()
func FlyoverConfigurationsEmptyTiersErrorID() common.Hash {
	return common.HexToHash("0x68b4e98f560f429565ce5506873230066f695520821666deacb15be9ffb7c3d4")
}

// UnpackEmptyTiersError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EmptyTiers()
func (flyoverConfigurations *FlyoverConfigurations) UnpackEmptyTiersError(raw []byte) (*FlyoverConfigurationsEmptyTiers, error) {
	out := new(FlyoverConfigurationsEmptyTiers)
	if err := flyoverConfigurations.abi.UnpackIntoInterface(out, "EmptyTiers", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverConfigurationsExpireTimeNotAfterCallTime represents a ExpireTimeNotAfterCallTime error raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsExpireTimeNotAfterCallTime struct {
	CallTime   *big.Int
	ExpireTime *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ExpireTimeNotAfterCallTime(uint256 callTime, uint256 expireTime)
func FlyoverConfigurationsExpireTimeNotAfterCallTimeErrorID() common.Hash {
	return common.HexToHash("0x34e8e278dbaf146ab19e93e47bf975c8ba479cb4a0da6a67583f3c23d03dd35a")
}

// UnpackExpireTimeNotAfterCallTimeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ExpireTimeNotAfterCallTime(uint256 callTime, uint256 expireTime)
func (flyoverConfigurations *FlyoverConfigurations) UnpackExpireTimeNotAfterCallTimeError(raw []byte) (*FlyoverConfigurationsExpireTimeNotAfterCallTime, error) {
	out := new(FlyoverConfigurationsExpireTimeNotAfterCallTime)
	if err := flyoverConfigurations.abi.UnpackIntoInterface(out, "ExpireTimeNotAfterCallTime", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverConfigurationsInvalidInitialization represents a InvalidInitialization error raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsInvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func FlyoverConfigurationsInvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (flyoverConfigurations *FlyoverConfigurations) UnpackInvalidInitializationError(raw []byte) (*FlyoverConfigurationsInvalidInitialization, error) {
	out := new(FlyoverConfigurationsInvalidInitialization)
	if err := flyoverConfigurations.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverConfigurationsNoQueuedChange represents a NoQueuedChange error raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsNoQueuedChange struct {
	Flow uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoQueuedChange(uint8 flow)
func FlyoverConfigurationsNoQueuedChangeErrorID() common.Hash {
	return common.HexToHash("0xd3c52278efdced481f64abe27d0cfcc81301842afa78d451e34486288a4b18c4")
}

// UnpackNoQueuedChangeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoQueuedChange(uint8 flow)
func (flyoverConfigurations *FlyoverConfigurations) UnpackNoQueuedChangeError(raw []byte) (*FlyoverConfigurationsNoQueuedChange, error) {
	out := new(FlyoverConfigurationsNoQueuedChange)
	if err := flyoverConfigurations.abi.UnpackIntoInterface(out, "NoQueuedChange", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverConfigurationsNotInitializing represents a NotInitializing error raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsNotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func FlyoverConfigurationsNotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (flyoverConfigurations *FlyoverConfigurations) UnpackNotInitializingError(raw []byte) (*FlyoverConfigurationsNotInitializing, error) {
	out := new(FlyoverConfigurationsNotInitializing)
	if err := flyoverConfigurations.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverConfigurationsOutOfBounds represents a OutOfBounds error raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsOutOfBounds struct {
	Flow  uint8
	Field uint8
	Value *big.Int
	Min   *big.Int
	Max   *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OutOfBounds(uint8 flow, uint8 field, uint256 value, uint256 min, uint256 max)
func FlyoverConfigurationsOutOfBoundsErrorID() common.Hash {
	return common.HexToHash("0x1036875555a8673ef01a271cb2bc1324477aac4914d8a2e48a679a3253d2f3ad")
}

// UnpackOutOfBoundsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OutOfBounds(uint8 flow, uint8 field, uint256 value, uint256 min, uint256 max)
func (flyoverConfigurations *FlyoverConfigurations) UnpackOutOfBoundsError(raw []byte) (*FlyoverConfigurationsOutOfBounds, error) {
	out := new(FlyoverConfigurationsOutOfBounds)
	if err := flyoverConfigurations.abi.UnpackIntoInterface(out, "OutOfBounds", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverConfigurationsPaymentNotAllowed represents a PaymentNotAllowed error raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsPaymentNotAllowed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PaymentNotAllowed()
func FlyoverConfigurationsPaymentNotAllowedErrorID() common.Hash {
	return common.HexToHash("0x8619bd43ab22b4b01742bd29d231dff1e50413ee3a444878bed65970c80c97df")
}

// UnpackPaymentNotAllowedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PaymentNotAllowed()
func (flyoverConfigurations *FlyoverConfigurations) UnpackPaymentNotAllowedError(raw []byte) (*FlyoverConfigurationsPaymentNotAllowed, error) {
	out := new(FlyoverConfigurationsPaymentNotAllowed)
	if err := flyoverConfigurations.abi.UnpackIntoInterface(out, "PaymentNotAllowed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverConfigurationsSafeCastOverflowedUintDowncast represents a SafeCastOverflowedUintDowncast error raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsSafeCastOverflowedUintDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func FlyoverConfigurationsSafeCastOverflowedUintDowncastErrorID() common.Hash {
	return common.HexToHash("0x6dfcc6503a32754ce7a89698e18201fc5294fd4aad43edefee786f88423b1a12")
}

// UnpackSafeCastOverflowedUintDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func (flyoverConfigurations *FlyoverConfigurations) UnpackSafeCastOverflowedUintDowncastError(raw []byte) (*FlyoverConfigurationsSafeCastOverflowedUintDowncast, error) {
	out := new(FlyoverConfigurationsSafeCastOverflowedUintDowncast)
	if err := flyoverConfigurations.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverConfigurationsTiersNotAscending represents a TiersNotAscending error raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsTiersNotAscending struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TiersNotAscending()
func FlyoverConfigurationsTiersNotAscendingErrorID() common.Hash {
	return common.HexToHash("0x5affaa88ccad0e156dff884ac88c0f44cb5bd170a771d2dd8f0ecde2970755d5")
}

// UnpackTiersNotAscendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TiersNotAscending()
func (flyoverConfigurations *FlyoverConfigurations) UnpackTiersNotAscendingError(raw []byte) (*FlyoverConfigurationsTiersNotAscending, error) {
	out := new(FlyoverConfigurationsTiersNotAscending)
	if err := flyoverConfigurations.abi.UnpackIntoInterface(out, "TiersNotAscending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverConfigurationsTimelockNotElapsed represents a TimelockNotElapsed error raised by the FlyoverConfigurations contract.
type FlyoverConfigurationsTimelockNotElapsed struct {
	Eta     *big.Int
	NowTime *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TimelockNotElapsed(uint256 eta, uint256 nowTime)
func FlyoverConfigurationsTimelockNotElapsedErrorID() common.Hash {
	return common.HexToHash("0xaadd29104503eb51c1adbfdf4338cd4aa6a6b2e98f6b115c6745592b51091079")
}

// UnpackTimelockNotElapsedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TimelockNotElapsed(uint256 eta, uint256 nowTime)
func (flyoverConfigurations *FlyoverConfigurations) UnpackTimelockNotElapsedError(raw []byte) (*FlyoverConfigurationsTimelockNotElapsed, error) {
	out := new(FlyoverConfigurationsTimelockNotElapsed)
	if err := flyoverConfigurations.abi.UnpackIntoInterface(out, "TimelockNotElapsed", raw); err != nil {
		return nil, err
	}
	return out, nil
}
