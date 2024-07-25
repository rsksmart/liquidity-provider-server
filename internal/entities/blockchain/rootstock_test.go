package blockchain_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewTransactionConfig(t *testing.T) {
	config := blockchain.NewTransactionConfig(entities.NewWei(1), 2, entities.NewWei(3))
	var value uint64 = 2
	assert.Equal(t, entities.NewWei(1), config.Value)
	assert.Equal(t, &value, config.GasLimit)
	assert.Equal(t, entities.NewWei(3), config.GasPrice)
}

func TestIsRskAddress(t *testing.T) {
	goodAddresses := []string{
		"0x7C4890A0f1D4bBf2C669Ac2d1efFa185c505359b",
		"0x79568c2989232dCa1840087D73d403602364c0D4",
		"0xd5f00ABfbEA7A0B193836CAc6833c2Ad9D06cEa8",
		"0x892813507Bf3aBF2890759d2135Ec34f4909Fea5",
		"0x5dE07e2BE63595854C396E2da291e0d1EdE15112",
		"0x0D8Fb5d32704DB2931e05DB91F64BcA6f76Ce573",
		"0x8dccd82443B80DDdE3690af86746bfd9d766F8D2",
		"0xa2011668bd246f9Aa10623f3Cfea704E3b6c0C3b",
		"0xBb519e5dCB3f98ED0c48238b42BFa3fd4d1a5E45",
		"0xe8d8c8f343522fd53c45c71723B93D735b149220",
	}

	badAddresses := []string{
		"mwtKGvtdDno6zzoioQHgWbV9A2i2kbfWcX",
		"0xe753be697499877faabae44049e7305afdfccd24fcf8b10f9e16ad0eec4aee6c",
		"0xe8d8c8f343522fd53c45c71723B93D735b149220c1",
		"0x8dccd82443B80DDdE3690af86746bfd9d766F8",
		"0x892813507Bf3aBF2890759d2135Ec34f4909ea5",
		"TCNtTa1rveKkovHR2ebABu4K66U6ocUCZX",
	}

	for _, address := range goodAddresses {
		assert.Truef(t, blockchain.IsRskAddress(address), "Address %s should be valid", address)
	}

	for _, address := range badAddresses {
		assert.Falsef(t, blockchain.IsRskAddress(address), "Address %s should not be valid", address)
	}
}

func TestDecodeStringTrimPrefix(t *testing.T) {
	cases := test.Table[string, []byte]{
		{Value: "0x7C4890A0f1D4bBf2C669Ac2d1efFa185c505359b", Result: []byte{124, 72, 144, 160, 241, 212, 187, 242, 198, 105, 172, 45, 30, 255, 161, 133, 197, 5, 53, 155}},
		{Value: "0x79568c2989232dCa1840087D73d403602364c0D4", Result: []byte{121, 86, 140, 41, 137, 35, 45, 202, 24, 64, 8, 125, 115, 212, 3, 96, 35, 100, 192, 212}},
		{Value: "0xd5f00ABfbEA7A0B193836CAc6833c2Ad9D06cEa8", Result: []byte{213, 240, 10, 191, 190, 167, 160, 177, 147, 131, 108, 172, 104, 51, 194, 173, 157, 6, 206, 168}},
		{Value: "0x892813507Bf3aBF2890759d2135Ec34f4909Fea5", Result: []byte{137, 40, 19, 80, 123, 243, 171, 242, 137, 7, 89, 210, 19, 94, 195, 79, 73, 9, 254, 165}},
		{Value: "0x5dE07e2BE63595854C396E2da291e0d1EdE15112", Result: []byte{93, 224, 126, 43, 230, 53, 149, 133, 76, 57, 110, 45, 162, 145, 224, 209, 237, 225, 81, 18}},
		{Value: "0x0D8Fb5d32704DB2931e05DB91F64BcA6f76Ce573", Result: []byte{13, 143, 181, 211, 39, 4, 219, 41, 49, 224, 93, 185, 31, 100, 188, 166, 247, 108, 229, 115}},
		{Value: "0x8dccd82443B80DDdE3690af86746bfd9d766F8D2", Result: []byte{141, 204, 216, 36, 67, 184, 13, 221, 227, 105, 10, 248, 103, 70, 191, 217, 215, 102, 248, 210}},
		{Value: "0xa2011668bd246f9Aa10623f3Cfea704E3b6c0C3b", Result: []byte{162, 1, 22, 104, 189, 36, 111, 154, 161, 6, 35, 243, 207, 234, 112, 78, 59, 108, 12, 59}},
		{Value: "0xBb519e5dCB3f98ED0c48238b42BFa3fd4d1a5E45", Result: []byte{187, 81, 158, 93, 203, 63, 152, 237, 12, 72, 35, 139, 66, 191, 163, 253, 77, 26, 94, 69}},
		{Value: "0xe8d8c8f343522fd53c45c71723B93D735b149220", Result: []byte{232, 216, 200, 243, 67, 82, 47, 213, 60, 69, 199, 23, 35, 185, 61, 115, 91, 20, 146, 32}},
	}
	var bytes []byte
	test.RunTable(t, cases, func(address string) []byte {
		bytes, _ = blockchain.DecodeStringTrimPrefix(address)
		return bytes
	})
}

func TestDecodeStringTrim_Fail(t *testing.T) {
	badAddresses := []string{
		"mwtKGvtdDno6zzoioQHgWbV9A2i2kbfWcX",
		"0x892813507Bf3aBF2890759d2135Ec34f4909ea5",
		"TCNtTa1rveKkovHR2ebABu4K66U6ocUCZX",
	}

	var err error
	for _, address := range badAddresses {
		_, err = blockchain.DecodeStringTrimPrefix(address)
		require.Error(t, err)
	}
}
