package quote_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
)

func TestValidateQuoteHash(t *testing.T) {
	tests := []struct {
		name    string
		hash    string
		wantErr bool
	}{
		{
			name:    "Valid 32-byte hash",
			hash:    "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			wantErr: false,
		},
		{
			name:    "Invalid length - too short",
			hash:    "1234567890abcdef",
			wantErr: true,
		},
		{
			name:    "Invalid length - too long",
			hash:    "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef00",
			wantErr: true,
		},
		{
			name:    "Invalid characters",
			hash:    "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdeg",
			wantErr: true,
		},
		{
			name:    "Empty string",
			hash:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := quote.ValidateQuoteHash(tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateQuoteHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// nolint:funlen
func TestCalculateCallFee(t *testing.T) {
	type testArgs struct {
		amount        *entities.Wei
		feePercentage *utils.BigFloat
		fixedFee      *entities.Wei
		result        *entities.Wei
	}
	testCases := []testArgs{
		{
			amount:        entities.NewWei(5000000000000000),
			feePercentage: utils.NewBigFloat64(0),
			fixedFee:      entities.NewWei(100000000000000),
			result:        entities.NewWei(100000000000000),
		},
		{
			amount:        entities.NewWei(5000000000000000),
			feePercentage: utils.NewBigFloat64(0.33),
			fixedFee:      entities.NewWei(0),
			result:        entities.NewWei(16500000000000),
		},
		{
			amount:        entities.NewWei(5000000000000000),
			feePercentage: utils.NewBigFloat64(0),
			fixedFee:      entities.NewWei(0),
			result:        entities.NewWei(0),
		},
		{
			amount:        entities.NewWei(5000000000000000),
			feePercentage: utils.NewBigFloat64(5.12),
			fixedFee:      entities.NewWei(123456789),
			result:        entities.NewWei(256000123456789),
		},
		{
			amount:        entities.NewWei(7777777777777777789),
			feePercentage: utils.NewBigFloat64(77.33),
			fixedFee:      entities.NewWei(0),
			result:        entities.NewWei(6014555555555555564),
		},
		{
			amount:        entities.NewWei(7777777777777777789),
			feePercentage: utils.NewBigFloat64(77.44),
			fixedFee:      entities.NewWei(0),
			result:        entities.NewWei(6023111111111111120),
		},
		{
			amount:        entities.NewWei(7777777777777777789),
			feePercentage: utils.NewBigFloat64(77.41),
			fixedFee:      entities.NewWei(0),
			result:        entities.NewWei(6020777777777777786),
		},
		{
			amount:        entities.NewWei(7777777777777777789),
			feePercentage: utils.NewBigFloat64(77.86),
			fixedFee:      entities.NewWei(0),
			result:        entities.NewWei(6055777777777777787),
		},
		{
			amount:        entities.NewWei(7777777777777777789),
			feePercentage: utils.NewBigFloat64(77.7),
			fixedFee:      entities.NewWei(1110000031224),
			result:        entities.NewWei(6043334443333364566),
		},
		{
			amount:        entities.NewWei(7777777777777777789),
			feePercentage: nil,
			fixedFee:      entities.NewWei(1110000031224),
			result:        entities.NewWei(1110000031224),
		},
		{
			amount:        entities.NewWei(7777777777777777789),
			feePercentage: utils.NewBigFloat64(77.7),
			fixedFee:      nil,
			result:        entities.NewWei(6043333333333333342),
		},
		{
			amount:        nil,
			feePercentage: utils.NewBigFloat64(77.7),
			fixedFee:      entities.NewWei(1110000031224),
			result:        entities.NewWei(1110000031224),
		},
	}
	log.SetLevel(log.DebugLevel)
	for _, tt := range testCases {
		config := &mocks.PegConfigurationMock{}
		config.EXPECT().GetFeePercentage().Return(tt.feePercentage)
		config.EXPECT().GetFixedFee().Return(tt.fixedFee)

		result := quote.CalculateCallFee(tt.amount, config)
		assert.Equal(t, tt.result, result, "Expected %v, got %v", tt.result, result)
		config.AssertExpectations(t)
	}
}
