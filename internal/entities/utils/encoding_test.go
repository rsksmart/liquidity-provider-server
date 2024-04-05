package utils_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDecodeKey(t *testing.T) {
	type params struct {
		key   string
		bytes int
	}
	cases := test.Table[params, []byte]{
		{
			Value:  params{key: "1234567890abcdef", bytes: 8},
			Result: []byte{18, 52, 86, 120, 144, 171, 205, 239},
		},
		{
			Value:  params{key: "a2fbac02d66202e8468d2a4f1deba4fa5c2491f592e0e22e32fe1e6acac25923", bytes: 32},
			Result: []byte{0xa2, 0xfb, 0xac, 0x2, 0xd6, 0x62, 0x2, 0xe8, 0x46, 0x8d, 0x2a, 0x4f, 0x1d, 0xeb, 0xa4, 0xfa, 0x5c, 0x24, 0x91, 0xf5, 0x92, 0xe0, 0xe2, 0x2e, 0x32, 0xfe, 0x1e, 0x6a, 0xca, 0xc2, 0x59, 0x23},
		},
		{
			Value:  params{key: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", bytes: 32},
			Result: []byte{0x9f, 0x86, 0xd0, 0x81, 0x88, 0x4c, 0x7d, 0x65, 0x9a, 0x2f, 0xea, 0xa0, 0xc5, 0x5a, 0xd0, 0x15, 0xa3, 0xbf, 0x4f, 0x1b, 0x2b, 0xb, 0x82, 0x2c, 0xd1, 0x5d, 0x6c, 0x15, 0xb0, 0xf0, 0xa, 0x8},
		},
		{
			Value:  params{key: "c5ff177a86e82441f93e3772da700d5f6838157fa1bfdc0bb689d7f7e55e7aba", bytes: 32},
			Result: []byte{0xc5, 0xff, 0x17, 0x7a, 0x86, 0xe8, 0x24, 0x41, 0xf9, 0x3e, 0x37, 0x72, 0xda, 0x70, 0xd, 0x5f, 0x68, 0x38, 0x15, 0x7f, 0xa1, 0xbf, 0xdc, 0xb, 0xb6, 0x89, 0xd7, 0xf7, 0xe5, 0x5e, 0x7a, 0xba},
		},
		{
			Value:  params{key: "ab5c2d1f", bytes: 4},
			Result: []byte{0xab, 0x5c, 0x2d, 0x1f},
		},
	}
	test.RunTable(t, cases, func(p params) []byte {
		result, err := utils.DecodeKey(p.key, p.bytes)
		require.NoError(t, err)
		return result
	})
}

func TestDecodeKey_SizeError(t *testing.T) {
	sizes := []int{4, 8, 12, 24, 32}
	key := "1122abcdff1122abcdff"
	for _, size := range sizes {
		result, err := utils.DecodeKey(key, size)
		require.Error(t, err)
		require.Nil(t, result)
		require.Errorf(t, err, "key length is not %d bytes, %s is %d bytes long", size, key, 10)
	}
}

func TestDecodeKey_DecodingError(t *testing.T) {
	cases := []string{
		"no hex",
		"abcde",
		"17",
		"g8ab11",
	}
	for _, key := range cases {
		result, err := utils.DecodeKey(key, 8)
		require.Error(t, err)
		require.Nil(t, result)
	}
}
