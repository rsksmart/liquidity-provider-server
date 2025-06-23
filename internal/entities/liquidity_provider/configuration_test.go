package liquidity_provider_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfirmationsPerAmount_ForValue(t *testing.T) {
	table := test.Table[*entities.Wei, uint16]{
		{Value: entities.NewWei(1), Result: uint16(40)},
		{Value: entities.NewWei(10000000), Result: uint16(40)},
		{Value: entities.NewWei(100000000000000000), Result: uint16(40)},
		{Value: entities.NewWei(100000000000000001), Result: uint16(120)},
		{Value: entities.NewWei(400000000000000000), Result: uint16(120)},
		{Value: entities.NewWei(400000000000000001), Result: uint16(200)},
		{Value: entities.NewWei(400000000050000000), Result: uint16(200)},
		{Value: entities.NewWei(2000000000000000000), Result: uint16(200)},
		{Value: entities.NewWei(2000000000000000001), Result: uint16(400)},
		{Value: entities.NewWei(4000000000000000000), Result: uint16(400)},
		{Value: entities.NewWei(4000000000000000001), Result: uint16(800)},
		{Value: entities.NewWei(4000000005000000000), Result: uint16(800)},
		{Value: entities.NewWei(8000000000000000000), Result: uint16(800)},
		{Value: entities.NewWei(8000000000000000005), Result: uint16(800)},
		{Value: entities.NewWei(9000000000000000000), Result: uint16(800)},
	}
	confirmations := liquidity_provider.DefaultRskConfirmationsPerAmount()
	test.RunTable(t, table, confirmations.ForValue)
}

func TestConfirmationsPerAmount_Max(t *testing.T) {
	table := test.Table[liquidity_provider.ConfirmationsPerAmount, uint16]{
		{
			Value:  liquidity_provider.DefaultRskConfirmationsPerAmount(),
			Result: uint16(800),
		},
		{
			Value:  liquidity_provider.DefaultBtcConfirmationsPerAmount(),
			Result: uint16(40),
		},
		{
			Value:  liquidity_provider.ConfirmationsPerAmount{},
			Result: uint16(0),
		},
		{
			Value:  liquidity_provider.ConfirmationsPerAmount{"100": 10, "200": 20, "300": 30},
			Result: uint16(30),
		},
	}
	test.RunTable(t, table, func(confirmations liquidity_provider.ConfirmationsPerAmount) uint16 {
		return confirmations.Max()
	})
}

func TestPeginConfiguration_ValidateAmount(t *testing.T) {
	config := liquidity_provider.DefaultPeginConfiguration()
	table := test.Table[*entities.Wei, error]{
		{
			Value:  entities.NewWei(1),
			Result: liquidity_provider.AmountOutOfRangeError,
		},
		{
			Value:  entities.NewWei(4999999999999999),
			Result: liquidity_provider.AmountOutOfRangeError,
		},
		{
			Value:  entities.NewWei(5000000000000000),
			Result: nil,
		},
		{
			Value:  entities.NewWei(100000000000000000),
			Result: nil,
		},
		{
			Value:  entities.NewWei(100000000000000001),
			Result: liquidity_provider.AmountOutOfRangeError,
		},
		{
			Value:  entities.NewWei(1000000000000000000),
			Result: liquidity_provider.AmountOutOfRangeError,
		},
	}
	for _, item := range table {
		err := config.ValidateAmount(item.Value)
		require.ErrorIs(t, err, item.Result)
	}
}

func TestPegoutConfiguration_ValidateAmount(t *testing.T) {
	config := liquidity_provider.DefaultPegoutConfiguration()
	table := test.Table[*entities.Wei, error]{
		{
			Value:  entities.NewWei(1),
			Result: liquidity_provider.AmountOutOfRangeError,
		},
		{
			Value:  entities.NewWei(4999999999999999),
			Result: liquidity_provider.AmountOutOfRangeError,
		},
		{
			Value:  entities.NewWei(5000000000000000),
			Result: nil,
		},
		{
			Value:  entities.NewWei(100000000000000000),
			Result: nil,
		},
		{
			Value:  entities.NewWei(100000000000000001),
			Result: liquidity_provider.AmountOutOfRangeError,
		},
		{
			Value:  entities.NewWei(1000000000000000000),
			Result: liquidity_provider.AmountOutOfRangeError,
		},
	}
	for _, item := range table {
		err := config.ValidateAmount(item.Value)
		require.ErrorIs(t, err, item.Result)
	}
}
