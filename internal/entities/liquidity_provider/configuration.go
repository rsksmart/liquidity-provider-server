package liquidity_provider

import (
	"errors"
	"fmt"
	"slices"

	"math/big"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
)

var (
	AmountOutOfRangeError = errors.New("amount out of range")
)

// ConfirmationsPerAmount the key represents the amount in wei serialized as a string, and the value represents the number of confirmations required for that amount.
type ConfirmationsPerAmount map[string]uint16

type PeginConfiguration struct {
	TimeForDeposit uint32          `json:"timeForDeposit" bson:"time_for_deposit" validate:"required"`
	CallTime       uint32          `json:"callTime" bson:"call_time" validate:"required"`
	PenaltyFee     *entities.Wei   `json:"penaltyFee" bson:"penalty_fee" validate:"required"`
	FixedFee       *entities.Wei   `json:"fixedFee" bson:"fixed_fee" validate:"required"`
	FeePercentage  *utils.BigFloat `json:"feePercentage" bson:"fee_percentage" validate:"required"`
	MaxValue       *entities.Wei   `json:"maxValue" bson:"max_value" validate:"required"`
	MinValue       *entities.Wei   `json:"minValue" bson:"min_value" validate:"required"`
}

func (config PeginConfiguration) ValidateAmount(amount *entities.Wei) error {
	return validateRange(config.MinValue, config.MaxValue, amount)
}

func (config PeginConfiguration) GetFixedFee() *entities.Wei {
	return config.FixedFee
}

func (config PeginConfiguration) GetFeePercentage() *utils.BigFloat {
	return config.FeePercentage
}

type PegoutConfiguration struct {
	TimeForDeposit       uint32          `json:"timeForDeposit" bson:"time_for_deposit" validate:"required"`
	ExpireTime           uint32          `json:"expireTime" bson:"expire_time" validate:"required"`
	PenaltyFee           *entities.Wei   `json:"penaltyFee" bson:"penalty_fee" validate:"required"`
	FixedFee             *entities.Wei   `json:"fixedFee" bson:"fixed_fee" validate:"required"`
	FeePercentage        *utils.BigFloat `json:"feePercentage" bson:"fee_percentage" validate:"required"`
	MaxValue             *entities.Wei   `json:"maxValue" bson:"max_value" validate:"required"`
	MinValue             *entities.Wei   `json:"minValue" bson:"min_value" validate:"required"`
	ExpireBlocks         uint64          `json:"expireBlocks" bson:"expire_blocks" validate:"required"`
	BridgeTransactionMin *entities.Wei   `json:"bridgeTransactionMin" bson:"bridge_transaction_min" validate:"required"`
}

func (config PegoutConfiguration) ValidateAmount(amount *entities.Wei) error {
	return validateRange(config.MinValue, config.MaxValue, amount)
}

func (config PegoutConfiguration) GetFixedFee() *entities.Wei {
	return config.FixedFee
}

func (config PegoutConfiguration) GetFeePercentage() *utils.BigFloat {
	return config.FeePercentage
}

type GeneralConfiguration struct {
	RskConfirmations          ConfirmationsPerAmount `json:"rskConfirmations" bson:"rsk_confirmations" validate:"required"`
	BtcConfirmations          ConfirmationsPerAmount `json:"btcConfirmations" bson:"btc_confirmations" validate:"required"`
	PublicLiquidityCheck      bool                   `json:"publicLiquidityCheck" bson:"public_liquidity_check" validate:""`
	MaxLiquidity              *entities.Wei          `json:"maxLiquidity" bson:"max_liquidity" validate:""`
	ExcessToleranceFixed      *entities.Wei          `json:"excessToleranceFixed" bson:"excess_tolerance_fixed" validate:""`
	ExcessTolerancePercentage *utils.BigFloat        `json:"excessTolerancePercentage" bson:"excess_tolerance_percentage" validate:""`
}

type HashedCredentials struct {
	HashedUsername string `bson:"hashed_username"`
	HashedPassword string `bson:"hashed_password"`
	UsernameSalt   string `bson:"username_salt"`
	PasswordSalt   string `bson:"password_salt"`
}

type ConfigurationType interface {
	PeginConfiguration | PegoutConfiguration | GeneralConfiguration | HashedCredentials | TrustedAccountDetails
}

func validateRange(min, max, amount *entities.Wei) error {
	if amount.Cmp(max) <= 0 && amount.Cmp(min) >= 0 {
		return nil
	} else {
		return fmt.Errorf("%w [%v, %v]", AmountOutOfRangeError, min, max)
	}
}

func (confirmations ConfirmationsPerAmount) Max() uint16 {
	// replace with slices.Max(maps.Values(lp.env.BtcConfig.Confirmations)) when its on stable go version
	var maxValue uint16
	for _, value := range confirmations {
		if maxValue < value {
			maxValue = value
		}
	}
	return maxValue
}

func (confirmations ConfirmationsPerAmount) ForValue(value *entities.Wei) uint16 {
	values := make([]*big.Int, 0)
	for key := range confirmations {
		bigIntKey := new(big.Int)
		_, ok := bigIntKey.SetString(key, 10)
		if !ok {
			bigIntKey.SetInt64(0)
		}
		values = append(values, bigIntKey)
	}
	slices.SortFunc(values, func(a, b *big.Int) int {
		return a.Cmp(b)
	})
	index := slices.IndexFunc(values, func(item *big.Int) bool {
		return value.AsBigInt().Cmp(item) <= 0
	})
	if index == -1 {
		return confirmations[values[len(values)-1].String()]
	} else {
		return confirmations[values[index].String()]
	}
}

type ServerInfo struct {
	Version  string
	Revision string
}
