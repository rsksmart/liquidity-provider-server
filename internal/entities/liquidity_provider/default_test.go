package liquidity_provider_test

import (
	"math/big"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/stretchr/testify/assert"
)

func TestDefaultBtcConfirmationsPerAmount(t *testing.T) {
	confirmations := liquidity_provider.DefaultBtcConfirmationsPerAmount()
	assert.Equal(t, liquidity_provider.ConfirmationsPerAmount{
		100000000000000000:  2,
		400000000000000000:  6,
		2000000000000000000: 10,
		4000000000000000000: 20,
		8000000000000000000: 40,
	}, confirmations)
}

func TestDefaultRskConfirmationsPerAmount(t *testing.T) {
	confirmations := liquidity_provider.DefaultRskConfirmationsPerAmount()
	assert.Equal(t, liquidity_provider.ConfirmationsPerAmount{
		100000000000000000:  40,
		400000000000000000:  120,
		2000000000000000000: 200,
		4000000000000000000: 400,
		8000000000000000000: 800,
	}, confirmations)
}

func TestDefaultPegoutConfiguration(t *testing.T) {
	config := liquidity_provider.DefaultPegoutConfiguration()
	assert.Equal(t, liquidity_provider.PegoutConfiguration{
		TimeForDeposit:       3600,
		ExpireTime:           10800,
		PenaltyFee:           entities.NewWei(10000000000000),
		FixedFee:             entities.NewWei(100000000000000),
		PercentageFee:        big.NewFloat(1.25),
		MaxValue:             entities.NewWei(100000000000000000),
		MinValue:             entities.NewWei(5000000000000000),
		ExpireBlocks:         500,
		BridgeTransactionMin: entities.NewWei(15000000000000000),
	}, config)
}

func TestDefaultPeginConfiguration(t *testing.T) {
	config := liquidity_provider.DefaultPeginConfiguration()
	assert.Equal(t, liquidity_provider.PeginConfiguration{
		TimeForDeposit: 3600,
		CallTime:       7200,
		PenaltyFee:     entities.NewWei(10000000000000),
		FixedFee:       entities.NewWei(500000000000000),
		PercentageFee:  big.NewFloat(1.25),
		MaxValue:       entities.NewWei(100000000000000000),
		MinValue:       entities.NewWei(5000000000000000),
	}, config)
}

func TestDefaultRskConfirmationsPerAmount_Max(t *testing.T) {
	confirmations := liquidity_provider.DefaultRskConfirmationsPerAmount()
	assert.Equal(t, uint16(800), confirmations.Max())
}

func TestDefaultBtcConfirmationsPerAmount_Max(t *testing.T) {
	confirmations := liquidity_provider.DefaultBtcConfirmationsPerAmount()
	assert.Equal(t, uint16(40), confirmations.Max())
}

func TestDefaultGeneralConfiguration(t *testing.T) {
	config := liquidity_provider.DefaultGeneralConfiguration()
	assert.Equal(t, liquidity_provider.GeneralConfiguration{
		RskConfirmations:     liquidity_provider.DefaultRskConfirmationsPerAmount(),
		BtcConfirmations:     liquidity_provider.DefaultBtcConfirmationsPerAmount(),
		PublicLiquidityCheck: false,
	}, config)
}
