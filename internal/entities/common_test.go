package entities_test

import (
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestSigned_CheckIntegrity(t *testing.T) {
	peginConfig := liquidity_provider.PeginConfiguration{
		TimeForDeposit: 3600,
		CallTime:       7200,
		PenaltyFee:     entities.NewUWei(1000000000000000),
		FixedFee:       entities.NewUWei(10000000000000000),
		FeePercentage:  utils.NewBigFloat64(0.01),
		MaxValue:       entities.NewUWei(10000000000000000000),
		MinValue:       entities.NewUWei(600000000000000000),
	}
	pegoutConfig := liquidity_provider.PegoutConfiguration{
		TimeForDeposit:       3600,
		ExpireTime:           7200,
		PenaltyFee:           entities.NewUWei(1000000000000000),
		FixedFee:             entities.NewUWei(10000000000000000),
		FeePercentage:        utils.NewBigFloat64(0.01),
		MaxValue:             entities.NewUWei(10000000000000000000),
		MinValue:             entities.NewUWei(600000000000000000),
		ExpireBlocks:         500,
		BridgeTransactionMin: entities.NewWei(1500000000000000000),
	}
	generalConfig := liquidity_provider.GeneralConfiguration{
		RskConfirmations: map[int]uint16{
			4000000000000000000: 400,
			8000000000000000000: 800,
			9000000000000000000: 801,
			100000000000000000:  41,
			2000000000000000000: 200,
			400000000000000000:  120,
		},
		BtcConfirmations: map[int]uint16{
			400000000000000000:  6,
			4000000000000000000: 20,
			8000000000000000000: 40,
			9000000000000000001: 45,
			100000000000000000:  3,
			2000000000000000000: 10,
		},
	}

	tests := []struct {
		signed entities.Signed[any]
		err    error
	}{
		{signed: entities.Signed[any]{Value: peginConfig, Hash: "2f8528af602b84e2b8083ed9bbbf21bca3bcfcc78948ed970c65aa58f4271870"}},
		{signed: entities.Signed[any]{Value: pegoutConfig, Hash: "e405a29c1e0469a284f45aa8b26c95a84d8ee993664e0cf2cc5f1a963212e432"}},
		{signed: entities.Signed[any]{Value: generalConfig, Hash: "77a1d9b2426955a2dbeb4e6b561607fbd8bd044de7a60c1ed77126e72ea3cb18"}},
		{signed: entities.Signed[any]{Value: peginConfig, Hash: "f3daab424654d2eeb2b50dc00b3e453e24ca1c690d80015f5f54d5f1fefaf900"}, err: entities.IntegrityError},
		{signed: entities.Signed[any]{Value: pegoutConfig, Hash: "3b3e7b075eb60b8c249f44a117f406c64992bafda1273f540277448abd14077e"}, err: entities.IntegrityError},
		{signed: entities.Signed[any]{Value: generalConfig, Hash: "3fecc42296c21a63dff80885f972ea88caf5038e47f014b1c91bb9b80529b757"}, err: entities.IntegrityError},
		{signed: entities.Signed[any]{Value: generalConfig, Hash: "not a hash"}, err: hex.InvalidByteError('n')},
		{signed: entities.Signed[any]{Value: map[string]int{"test": 5}, Hash: "17bdb7aeb84082e4f0bf751ba78ee1fea05982f93d01e41016d1aeaaa718e18b"}},
	}

	for _, testCase := range tests {
		err := testCase.signed.CheckIntegrity(crypto.Keccak256)
		require.ErrorIs(t, err, testCase.err)
	}
}

func TestSigned_CheckIntegrity_encoding(t *testing.T) {
	var encodingErr *json.UnsupportedValueError
	err := entities.Signed[any]{
		Value: map[string]float64{"test": math.NaN()},
		Hash:  "17bdb7aeb84082e4f0bf751ba78ee1fea05982f93d01e41016d1aeaaa718e18b",
	}.CheckIntegrity(crypto.Keccak256)
	require.ErrorAs(t, err, &encodingErr)
}
