package bitcoin_test

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type ScriptEntry struct {
	Name          string
	ScriptHex     string
	SegwitProgram string
	P2SHP2WSHAddr string
}

var ScriptDataset = []ScriptEntry{
	{
		Name:          "2-of-3 Multisig",
		ScriptHex:     "52210279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f817982103828d9817f20b62f5c8c6fba48e33cfad4b9a2ed81ebb57e947d48f5e9f0f5b1d21027d1fe0e10b7a65f2aa5c3b2e0f593c9f8db3c9f5a639146e0d8df0f0bdcb0a3653ae",
		SegwitProgram: "00208bf61295ed4b6005a0452aa64873a49eabbc9204db0f5b51cf82c02805a3e6d9",
		P2SHP2WSHAddr: "3BwNcbMRRbsQ2V3dT5GLZeos2FGLh46NAY",
	},
	{
		Name:          "CSV",
		ScriptHex:     "a9144733f37cf4db86fbc2efed2500b4f4e49f3120238763b175",
		SegwitProgram: "002097dfc3e76a746db584e46036ab39e3473b4d7ce704da1f514e06be7a7056b110",
		P2SHP2WSHAddr: "3KgrN1XJms8tGm9zkLuUWv7R4aP3h2oQSL",
	},
	{
		Name:          "IF/ELSE",
		ScriptHex:     "63210279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f8179867522103828d9817f20b62f5c8c6fba48e33cfad4b9a2ed81ebb57e947d48f5e9f0f5b1d68ac",
		SegwitProgram: "00204f6f9370bd1dd3783180f01c78fc31bd03c5cb62fdbc48d6973362044ff36fe1",
		P2SHP2WSHAddr: "3622fSbuttx3uVUwgHnRRPZeNtSRHQsQB8",
	},
}

func TestScriptToAddressP2shP2wsh(t *testing.T) {
	for _, test := range ScriptDataset {
		t.Run("should parse "+test.Name, func(t *testing.T) {
			script, err := hex.DecodeString(test.ScriptHex)
			require.NoError(t, err)
			result, err := bitcoin.ScriptToAddressP2shP2wsh(script, &chaincfg.MainNetParams)
			require.NoError(t, err)
			assert.Equal(t, test.P2SHP2WSHAddr, result.EncodeAddress())
		})
	}
	t.Run("should return error for empty script", func(t *testing.T) {
		const errorMsg = "script cannot be empty"
		nilResult, nilErr := bitcoin.ScriptToAddressP2shP2wsh(nil, &chaincfg.MainNetParams)
		emptyResult, emptyErr := bitcoin.ScriptToAddressP2shP2wsh([]byte{}, &chaincfg.MainNetParams)
		require.ErrorContains(t, nilErr, errorMsg)
		require.Error(t, emptyErr, errorMsg)
		assert.Nil(t, nilResult)
		assert.Nil(t, emptyResult)
	})

	t.Run("should return error for empty network params", func(t *testing.T) {
		script, err := hex.DecodeString(ScriptDataset[0].ScriptHex)
		require.NoError(t, err)
		result, err := bitcoin.ScriptToAddressP2shP2wsh(script, nil)
		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func TestScriptToP2shP2wsh(t *testing.T) {
	for _, test := range ScriptDataset {
		t.Run("should parse "+test.Name, func(t *testing.T) {
			script, err := hex.DecodeString(test.ScriptHex)
			require.NoError(t, err)
			result, err := bitcoin.ScriptToP2shP2wsh(script)
			require.NoError(t, err)
			assert.Equal(t, test.SegwitProgram, hex.EncodeToString(result))
		})
	}
	t.Run("should return error for empty script", func(t *testing.T) {
		const errorMsg = "script cannot be empty"
		nilResult, nilErr := bitcoin.ScriptToP2shP2wsh(nil)
		emptyResult, emptyErr := bitcoin.ScriptToP2shP2wsh([]byte{})
		require.ErrorContains(t, nilErr, errorMsg)
		require.Error(t, emptyErr, errorMsg)
		assert.Nil(t, nilResult)
		assert.Nil(t, emptyResult)
	})
}
