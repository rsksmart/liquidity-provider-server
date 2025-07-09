package pkg_test

import (
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
)

func TestToAvailableLiquidityDTO(t *testing.T) {
	peginLiquidity := new(big.Int)
	peginLiquidity.SetString("1234567890987654321", 10)
	pegoutLiquidity := new(big.Int)
	pegoutLiquidity.SetString("9876543210123456789", 10)

	liquidity := liquidity_provider.AvailableLiquidity{
		PeginLiquidity:  entities.NewBigWei(peginLiquidity),
		PegoutLiquidity: entities.NewBigWei(pegoutLiquidity),
	}
	dto := pkg.ToAvailableLiquidityDTO(liquidity)
	assert.Equal(t, "1234567890987654321", dto.PeginLiquidityAmount.String())
	assert.Equal(t, "9876543210123456789", dto.PegoutLiquidityAmount.String())
}

func TestFromPeginConfigurationDTO(t *testing.T) {
	dto := pkg.PeginConfigurationDTO{
		TimeForDeposit: 10,
		CallTime:       200,
		PenaltyFee:     "3000000000000000000000",
		FixedFee:       "5000000000000000000000",
		FeePercentage:  5.443321101,
		MaxValue:       "7000000000000000000000",
		MinValue:       "6000000000000000000000",
	}
	configuration := pkg.FromPeginConfigurationDTO(dto)
	assert.Equal(t, uint32(10), configuration.TimeForDeposit)
	assert.Equal(t, uint32(200), configuration.CallTime)
	assert.Equal(t, "3000000000000000000000", configuration.PenaltyFee.AsBigInt().String())
	assert.Equal(t, "5000000000000000000000", configuration.FixedFee.AsBigInt().String())
	assert.Equal(t, "5.443321101", configuration.FeePercentage.Native().String())
	assert.Equal(t, "7000000000000000000000", configuration.MaxValue.AsBigInt().String())
	assert.Equal(t, "6000000000000000000000", configuration.MinValue.AsBigInt().String())
	test.AssertNonZeroValues(t, dto)
}

func TestFromPegoutConfigurationDTO(t *testing.T) {
	dto := pkg.PegoutConfigurationDTO{
		TimeForDeposit:       10,
		ExpireTime:           200,
		PenaltyFee:           "3000000000000000000000",
		FixedFee:             "5000000000000000000000",
		FeePercentage:        0.5123333,
		MaxValue:             "7000000000000000000000",
		MinValue:             "6000000000000000000000",
		ExpireBlocks:         20,
		BridgeTransactionMin: "8000000000000000000000",
	}
	configuration := pkg.FromPegoutConfigurationDTO(dto)
	assert.Equal(t, uint32(10), configuration.TimeForDeposit)
	assert.Equal(t, uint32(200), configuration.ExpireTime)
	assert.Equal(t, "3000000000000000000000", configuration.PenaltyFee.AsBigInt().String())
	assert.Equal(t, "5000000000000000000000", configuration.FixedFee.AsBigInt().String())
	assert.Equal(t, "0.5123333", configuration.FeePercentage.Native().String())
	assert.Equal(t, "7000000000000000000000", configuration.MaxValue.AsBigInt().String())
	assert.Equal(t, "6000000000000000000000", configuration.MinValue.AsBigInt().String())
	assert.Equal(t, uint64(20), configuration.ExpireBlocks)
	assert.Equal(t, "8000000000000000000000", configuration.BridgeTransactionMin.AsBigInt().String())
	test.AssertNonZeroValues(t, dto)
}

func TestToServerInfoDTO(t *testing.T) {
	serverInfo := liquidity_provider.ServerInfo{
		Version:  "1.0.0",
		Revision: "1234567890",
	}
	dto := pkg.ToServerInfoDTO(serverInfo)
	assert.Equal(t, "1.0.0", dto.Version)
	assert.Equal(t, "1234567890", dto.Revision)
}

// nolint:funlen
func TestLocalLiquidityProvider_ProviderDTOValidation(t *testing.T) {
	t.Run("Test FromPegoutConfigurationDTO conversion", func(t *testing.T) {
		dto := pkg.PegoutConfigurationDTO{
			TimeForDeposit:       3600,
			ExpireTime:           7200,
			PenaltyFee:           "1000000000000000",
			FixedFee:             "2000000000000000",
			FeePercentage:        1.5,
			MaxValue:             "1000000000000000000",
			MinValue:             "100000000000000000",
			ExpireBlocks:         500,
			BridgeTransactionMin: "50000000000000000",
		}
		penaltyFeeBigInt := new(big.Int)
		penaltyFeeBigInt.SetString(dto.PenaltyFee, 10)
		fixedFeeBigInt := new(big.Int)
		fixedFeeBigInt.SetString(dto.FixedFee, 10)
		maxValueBigInt := new(big.Int)
		maxValueBigInt.SetString(dto.MaxValue, 10)
		minValueBigInt := new(big.Int)
		minValueBigInt.SetString(dto.MinValue, 10)
		bridgeTransactionMinBigInt := new(big.Int)
		bridgeTransactionMinBigInt.SetString(dto.BridgeTransactionMin, 10)
		expectedConfig := liquidity_provider.PegoutConfiguration{
			TimeForDeposit:       dto.TimeForDeposit,
			ExpireTime:           dto.ExpireTime,
			PenaltyFee:           entities.NewBigWei(penaltyFeeBigInt),
			FixedFee:             entities.NewBigWei(fixedFeeBigInt),
			FeePercentage:        utils.NewBigFloat64(dto.FeePercentage),
			MaxValue:             entities.NewBigWei(maxValueBigInt),
			MinValue:             entities.NewBigWei(minValueBigInt),
			ExpireBlocks:         dto.ExpireBlocks,
			BridgeTransactionMin: entities.NewBigWei(bridgeTransactionMinBigInt),
		}
		config := pkg.FromPegoutConfigurationDTO(dto)
		assert.Equal(t, expectedConfig, config)
		test.AssertNonZeroValues(t, dto)
	})
	t.Run("Test ToPeginConfigurationDTO conversion", func(t *testing.T) {
		config := liquidity_provider.PeginConfiguration{
			TimeForDeposit: 3600,
			CallTime:       7200,
			PenaltyFee:     entities.NewWei(1000000000000000),
			FixedFee:       entities.NewWei(2000000000000000),
			FeePercentage:  utils.NewBigFloat64(1.5),
			MaxValue:       entities.NewWei(1000000000000000000),
			MinValue:       entities.NewWei(100000000000000000),
		}
		dto := pkg.ToPeginConfigurationDTO(config)
		feePercentage, _ := config.FeePercentage.Native().Float64()
		expectedDTO := pkg.PeginConfigurationDTO{
			TimeForDeposit: config.TimeForDeposit,
			CallTime:       config.CallTime,
			PenaltyFee:     config.PenaltyFee.AsBigInt().String(),
			FixedFee:       config.FixedFee.AsBigInt().String(),
			FeePercentage:  feePercentage,
			MaxValue:       config.MaxValue.AsBigInt().String(),
			MinValue:       config.MinValue.AsBigInt().String(),
		}
		assert.Equal(t, expectedDTO, dto)
		test.AssertNonZeroValues(t, dto)
	})
}

func TestToTrustedAccountDTO(t *testing.T) {
	btcLockingCap := new(big.Int)
	btcLockingCap.SetString("5000000000000000000", 10)
	rbtcLockingCap := new(big.Int)
	rbtcLockingCap.SetString("7000000000000000000", 10)
	trustedAccount := liquidity_provider.TrustedAccountDetails{
		Address:        "0x1234567890abcdef",
		Name:           "Test Trusted Account",
		BtcLockingCap:  entities.NewBigWei(btcLockingCap),
		RbtcLockingCap: entities.NewBigWei(rbtcLockingCap),
	}
	dto := pkg.ToTrustedAccountDTO(trustedAccount)
	assert.Equal(t, "0x1234567890abcdef", dto.Address)
	assert.Equal(t, "Test Trusted Account", dto.Name)
	assert.Equal(t, "5000000000000000000", dto.BtcLockingCap.String())
	assert.Equal(t, "7000000000000000000", dto.RbtcLockingCap.String())
}

func TestToTrustedAccountsDTO(t *testing.T) {
	btcLockingCap1 := new(big.Int)
	btcLockingCap1.SetString("5000000000000000000", 10)
	rbtcLockingCap1 := new(big.Int)
	rbtcLockingCap1.SetString("7000000000000000000", 10)
	btcLockingCap2 := new(big.Int)
	btcLockingCap2.SetString("9000000000000000000", 10)
	rbtcLockingCap2 := new(big.Int)
	rbtcLockingCap2.SetString("3000000000000000000", 10)
	account1 := liquidity_provider.TrustedAccountDetails{
		Address:        "0x1234567890abcdef",
		Name:           "Test Trusted Account 1",
		BtcLockingCap:  entities.NewBigWei(btcLockingCap1),
		RbtcLockingCap: entities.NewBigWei(rbtcLockingCap1),
	}
	account2 := liquidity_provider.TrustedAccountDetails{
		Address:        "0xabcdef1234567890",
		Name:           "Test Trusted Account 2",
		BtcLockingCap:  entities.NewBigWei(btcLockingCap2),
		RbtcLockingCap: entities.NewBigWei(rbtcLockingCap2),
	}

	signedAccounts := []entities.Signed[liquidity_provider.TrustedAccountDetails]{
		{
			Value:     account1,
			Signature: "signature1",
			Hash:      "hash1",
		},
		{
			Value:     account2,
			Signature: "signature2",
			Hash:      "hash2",
		},
	}

	dtos := pkg.ToTrustedAccountsDTO(signedAccounts)
	assert.Len(t, dtos, 2)
	assert.Equal(t, "0x1234567890abcdef", dtos[0].Address)
	assert.Equal(t, "Test Trusted Account 1", dtos[0].Name)
	assert.Equal(t, "5000000000000000000", dtos[0].BtcLockingCap.String())
	assert.Equal(t, "7000000000000000000", dtos[0].RbtcLockingCap.String())
	assert.Equal(t, "0xabcdef1234567890", dtos[1].Address)
	assert.Equal(t, "Test Trusted Account 2", dtos[1].Name)
	assert.Equal(t, "9000000000000000000", dtos[1].BtcLockingCap.String())
	assert.Equal(t, "3000000000000000000", dtos[1].RbtcLockingCap.String())
}

func TestFromGeneralConfigurationDTO(t *testing.T) {
	t.Run("converts valid configuration", func(t *testing.T) {
		dto := pkg.GeneralConfigurationDTO{
			RskConfirmations: map[string]uint16{
				"1000000000000000000": 5,
				"2000000000000000000": 10,
			},
			BtcConfirmations: map[string]uint16{
				"3000000000000000000": 15,
				"4000000000000000000": 20,
			},
			PublicLiquidityCheck: true,
		}

		config, err := pkg.FromGeneralConfigurationDTO(dto)

		require.NoError(t, err)
		assert.Equal(t, dto.RskConfirmations, map[string]uint16(config.RskConfirmations))
		assert.Equal(t, dto.BtcConfirmations, map[string]uint16(config.BtcConfirmations))
		assert.Equal(t, dto.PublicLiquidityCheck, config.PublicLiquidityCheck)
		test.AssertNonZeroValues(t, dto)
	})

	t.Run("returns error on invalid numeric keys", func(t *testing.T) {
		invalidBtc := pkg.GeneralConfigurationDTO{
			RskConfirmations: map[string]uint16{
				"1000000000000000000": 5,
			},
			BtcConfirmations: map[string]uint16{
				"3000000000000000000": 15,
				"notanumber":          20,
			},
			PublicLiquidityCheck: true,
		}
		invalidRsk := pkg.GeneralConfigurationDTO{
			RskConfirmations: map[string]uint16{
				"1000000000000000000": 5,
				"invalid":             10,
			},
			BtcConfirmations: map[string]uint16{
				"3000000000000000000": 15,
			},
			PublicLiquidityCheck: false,
		}

		config, err := pkg.FromGeneralConfigurationDTO(invalidBtc)
		assert.Empty(t, config)
		require.ErrorContains(t, err, "cannot deserialize BTC confirmations key notanumber")

		config, err = pkg.FromGeneralConfigurationDTO(invalidRsk)
		assert.Empty(t, config)
		require.ErrorContains(t, err, "cannot deserialize RSK confirmations key invalid")
	})
}
