package connectors

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"sort"
	"strings"
	"testing"
)

func testDerivationComplete(t *testing.T) {
	for _, tt := range TestQuotes {
		tt.FedInfo.IrisActivationHeight = 1
		if !common.IsHexAddress(tt.LBCAddr) {
			t.Errorf("invalid address: %v", tt.LBCAddr)
			continue
		}
		lbcAddr := common.FromHex(tt.LBCAddr)
		hashBytes, err := hex.DecodeString(tt.QuoteHash)
		if err != nil || len(hashBytes) == 0 {
			t.Errorf("Cannot parse QuoteHash correctly. value: %v, error: %v", tt.QuoteHash, err)
			continue
		}
		userBtcRefundAddr, err := DecodeBTCAddressWithVersion(tt.BTCRefundAddr)
		if err != nil {
			t.Errorf("Unexpected error in getBytesFromBtcAddress. error: %v", err)
			continue
		}
		lpBtcAddress, err := DecodeBTCAddressWithVersion(tt.LPBTCAddr)
		if err != nil {
			t.Errorf("Unexpected error in getBytesFromBtcAddress. error: %v", err)
			continue
		}
		value, err := getDerivationValueHash(userBtcRefundAddr, lbcAddr, lpBtcAddress, hashBytes)
		if err != nil {
			t.Errorf("Unexpected error in GetDerivationValueHash. value: %v, expected: %v, error: %v", value, tt.ExpectedDerivationValueHash, err)
			continue
		}
		result := hex.EncodeToString(value)
		assert.EqualValues(t, tt.ExpectedDerivationValueHash, result)
		buf, err := getFlyoverPrefix(value)
		if err != nil {
			t.Errorf("Unexpected error in getFlyoverPrefix. error: %v", err)
			continue
		}
		btc, err := NewBTC(tt.NetworkParams)
		if err != nil {
			t.Errorf("error initializing BTC: %v", err)
			continue
		}
		fedInfo := getFakeFedInfo()
		if btc.params.Name == chaincfg.TestNet3Params.Name {
			fedInfo.FedAddress = "2N5muMepJizJE1gR7FbHJU6CD18V3BpNF9p"
		} else {
			fedInfo.FedAddress = "3EDhHutH7XnsotnZaTfRr9CwnnGsNNrhCL"
		}
		scriptBuf, err := fedInfo.getPowPegRedeemScriptBuf(true)
		if err != nil {
			t.Errorf("Unexpected error in getPowPegRedeemScriptBuf. error: %v", err)
			continue
		}
		buf.Write(scriptBuf.Bytes())
		addr, err := btcutil.NewAddressScriptHash(buf.Bytes(), &btc.params)
		if err != nil {
			t.Errorf("Unexpected error in NewAddressScriptHash. error: %v", err)
			continue
		}
		assert.EqualValues(t, tt.ExpectedAddressHash, addr.EncodeAddress())
		err = fedInfo.validateRedeemScript(btc.params, scriptBuf.Bytes())
		if err != nil {
			t.Errorf("error in validateRedeemScript: %v", err)
		}
	}
}

func testBuildPowPegRedeemScript(t *testing.T) {
	fedInfo := getFakeFedInfo()
	buf, err := fedInfo.getPowPegRedeemScriptBuf(true)
	if err != nil {
		t.Fatalf("error in getPowPegRedeemScriptBuf: %v", err)
	}

	str := hex.EncodeToString(buf.Bytes())
	assert.True(t, checkSubstrings(str, fedInfo.PubKeys...))

	op2 := fmt.Sprintf("%02x", txscript.OP_2)
	assert.EqualValues(t, str[0:2], op2)

	op3 := fmt.Sprintf("%02x", txscript.OP_3)
	assert.EqualValues(t, str[len(str)-4:len(str)-2], op3)

	sort.Slice(fedInfo.PubKeys, func(i, j int) bool {
		return fedInfo.PubKeys[i] < fedInfo.PubKeys[j]
	})

	buf2, err := fedInfo.getPowPegRedeemScriptBuf(true)
	if err != nil {
		t.Errorf("error in getPowPegRedeemScriptBuf: %v", err)
	}
	str2 := hex.EncodeToString(buf2.Bytes())
	assert.EqualValues(t, str2, str)
}

func testBuildErpRedeemScript(t *testing.T) {
	btc, err := NewBTC("mainnet")
	if err != nil {
		t.Fatalf("error initializing BTC: %v", err)
	}

	fedInfo := getFakeFedInfo()
	buf, err := fedInfo.getErpRedeemScriptBuf(btc.params)
	if err != nil {
		t.Fatalf("error in getErpRedeemScriptBuf: %v", err)
	}

	str := hex.EncodeToString(buf.Bytes())
	assert.True(t, checkSubstrings(str, fedInfo.ErpKeys...))
	assert.EqualValues(t, str, ErpScriptString)
}

func testBuildFlyoverRedeemScript(t *testing.T) {
	btc, err := NewBTC("mainnet")
	if err != nil {
		t.Fatalf("error initializing BTC: %v", err)
	}
	fedInfo := getFakeFedInfo()
	fedInfo.FedAddress = "3EDhHutH7XnsotnZaTfRr9CwnnGsNNrhCL"
	fedInfo.IrisActivationHeight = 1

	derivationValue, err := getFlyoverDerivationHash()
	if err != nil {
		t.Fatalf("error in getFlyoverDerivationHash: %v", err)
	}

	fedRedeemScript, err := fedInfo.getFedRedeemScript(btc.params)
	if err != nil {
		t.Fatalf("error in getFedRedeemScript: %v", err)
	}
	bts, err := fedInfo.getFlyoverRedeemScript(derivationValue, fedRedeemScript)
	if err != nil {
		t.Fatalf("error in getFlyoverRedeemScript: %v", err)
	}

	str := hex.EncodeToString(bts)
	assert.True(t, checkSubstrings(str, fedInfo.PubKeys...))
	assert.EqualValues(t, FlyoverScriptString, str)
}

func testBuildFlyoverErpRedeemScript(t *testing.T) {
	btc, err := NewBTC("mainnet")
	if err != nil {
		t.Fatalf("error initializing BTC: %v", err)
	}
	fedInfo := getFakeFedInfo()
	fedInfo.FedAddress = "3C8e41MpbE2MB8XDqaYnQ2FbtRwPYLJtto"
	fedInfo.IrisActivationHeight = -1

	derivationValue, err := getFlyoverDerivationHash()
	if err != nil {
		t.Fatalf("error in getFlyoverDerivationHash: %v", err)
	}

	fedRedeemScript, err := fedInfo.getFedRedeemScript(btc.params)
	if err != nil {
		t.Fatalf("error in getFedRedeemScript: %v", err)
	}
	bts, err := fedInfo.getFlyoverRedeemScript(derivationValue, fedRedeemScript)
	if err != nil {
		t.Fatalf("error in getFlyoverRedeemScript: %v", err)
	}

	str := hex.EncodeToString(bts)
	assert.True(t, checkSubstrings(str, fedInfo.ErpKeys...))
	assert.EqualValues(t, FlyoverErpScriptString, str)
}

func testBuildFlyoverErpRedeemScriptFallback(t *testing.T) {
	btc, err := NewBTC("mainnet")
	if err != nil {
		t.Fatalf("error initializing BTC: %v", err)
	}
	fedInfo := getFakeFedInfo()
	fedInfo.FedAddress = "3EDhHutH7XnsotnZaTfRr9CwnnGsNNrhCL"
	fedInfo.IrisActivationHeight = -1

	derivationValue, err := getFlyoverDerivationHash()
	if err != nil {
		t.Fatalf("error in getFlyoverDerivationHash: %v", err)
	}

	fedRedeemScript, err := fedInfo.getFedRedeemScript(btc.params)
	if err != nil {
		t.Fatalf("error in getFedRedeemScript: %v", err)
	}
	bts, err := fedInfo.getFlyoverRedeemScript(derivationValue, fedRedeemScript)
	if err != nil {
		t.Fatalf("error in getFlyoverRedeemScript: %v", err)
	}

	str := hex.EncodeToString(bts)
	assert.True(t, checkSubstrings(str, fedInfo.PubKeys...))
	assert.EqualValues(t, FlyoverScriptString, str)
}

func testBuildPowPegAddressHash(t *testing.T) {
	fedInfo := getFakeFedInfo()
	fedInfo.IrisActivationHeight = 1

	buf, err := fedInfo.getPowPegRedeemScriptBuf(true)
	if err != nil {
		t.Fatalf("error in getPowPegRedeemScriptBuf: %v", err)
	}

	str := hex.EncodeToString(buf.Bytes())
	assert.True(t, checkSubstrings(str, fedInfo.PubKeys...))
	assert.EqualValues(t, str, PowPegScriptString)

	address, err := btcutil.NewAddressScriptHash(buf.Bytes(), &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("error in NewAddressScriptHash: %v", err)
	}

	expectedAddr := "3EDhHutH7XnsotnZaTfRr9CwnnGsNNrhCL"

	assert.EqualValues(t, expectedAddr, address.EncodeAddress())
}

func testBuildFlyoverPowPegAddressHash(t *testing.T) {
	btc, err := NewBTC("mainnet")
	if err != nil {
		t.Fatalf("error initializing BTC: %v", err)
	}
	fedInfo := getFakeFedInfo()
	fedInfo.FedAddress = "3EDhHutH7XnsotnZaTfRr9CwnnGsNNrhCL"
	fedInfo.IrisActivationHeight = 1

	derivationValue, err := getFlyoverDerivationHash()
	if err != nil {
		t.Fatalf("error in getFlyoverDerivationHash: %v", err)
	}

	fedRedeemScript, err := fedInfo.getFedRedeemScript(btc.params)
	if err != nil {
		t.Fatalf("error in getFedRedeemScript: %v", err)
	}
	bts, err := fedInfo.getFlyoverRedeemScript(derivationValue, fedRedeemScript)
	if err != nil {
		t.Fatalf("error in getFlyoverRedeemScript: %v", err)
	}

	str := hex.EncodeToString(bts)
	assert.True(t, checkSubstrings(str, fedInfo.PubKeys...))
	assert.EqualValues(t, FlyoverScriptString, str)

	address, err := btcutil.NewAddressScriptHash(bts, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("error in NewAddressScriptHash: %v", err)
	}
	expectedAddr := "34TNebhLLHsE6FHQVMmeHAhTFpaAWhfweR"

	assert.EqualValues(t, expectedAddr, address.EncodeAddress())
}

func testBuildFlyoverErpAddressHash(t *testing.T) {
	btc, err := NewBTC("mainnet")
	if err != nil {
		t.Fatalf("error initializing BTC: %v", err)
	}
	fedInfo := getFakeFedInfo()
	fedInfo.FedAddress = "3C8e41MpbE2MB8XDqaYnQ2FbtRwPYLJtto"
	fedInfo.IrisActivationHeight = -1

	derivationValue, err := getFlyoverDerivationHash()
	if err != nil {
		t.Fatalf("error in getFlyoverDerivationHash: %v", err)
	}

	fedRedeemScript, err := fedInfo.getFedRedeemScript(btc.params)
	if err != nil {
		t.Fatalf("error in getFedRedeemScript: %v", err)
	}
	bts, err := fedInfo.getFlyoverRedeemScript(derivationValue, fedRedeemScript)
	if err != nil {
		t.Fatalf("error in getFlyoverRedeemScript: %v", err)
	}

	str := hex.EncodeToString(bts)
	assert.True(t, checkSubstrings(str, fedInfo.ErpKeys...))
	assert.EqualValues(t, FlyoverErpScriptString, str)

	address, err := btcutil.NewAddressScriptHash(bts, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("error in NewAddressScriptHash: %v", err)
	}
	expectedAddr := "3PS2FEphLJMbJURMdYYFNAZR6zLasX51RC"

	assert.EqualValues(t, expectedAddr, address.EncodeAddress())
}

func testBuildFlyoverErpAddressHashFallback(t *testing.T) {
	btc, err := NewBTC("mainnet")
	if err != nil {
		t.Fatalf("error initializing BTC: %v", err)
	}
	fedInfo := getFakeFedInfo()
	fedInfo.FedAddress = "3EDhHutH7XnsotnZaTfRr9CwnnGsNNrhCL"
	fedInfo.IrisActivationHeight = -1

	derivationValue, err := getFlyoverDerivationHash()
	if err != nil {
		t.Fatalf("error in getFlyoverDerivationHash: %v", err)
	}

	fedRedeemScript, err := fedInfo.getFedRedeemScript(btc.params)
	if err != nil {
		t.Fatalf("error in getFedRedeemScript: %v", err)
	}
	bts, err := fedInfo.getFlyoverRedeemScript(derivationValue, fedRedeemScript)
	if err != nil {
		t.Fatalf("error in getFlyoverRedeemScript: %v", err)
	}

	str := hex.EncodeToString(bts)
	assert.True(t, checkSubstrings(str, fedInfo.PubKeys...))
	assert.EqualValues(t, FlyoverScriptString, str)

	address, err := btcutil.NewAddressScriptHash(bts, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("error in NewAddressScriptHash: %v", err)
	}
	expectedAddr := "34TNebhLLHsE6FHQVMmeHAhTFpaAWhfweR"

	assert.EqualValues(t, expectedAddr, address.EncodeAddress())
}

func getFlyoverDerivationHash() ([]byte, error) {
	sHash := "ffe4766f7b5f2fdf374f8ae02270d713c4dcb4b1c5d42bffda61b7f4c1c4c6c9"
	return hex.DecodeString(sHash)
}

func checkSubstrings(str string, subs ...string) bool {
	isCompleteMatch := true
	for _, sub := range subs {
		if !strings.Contains(str, sub) {
			isCompleteMatch = false
		}
	}

	return isCompleteMatch
}

func TestFedInfo(t *testing.T) {
	t.Run("test derivation complete", testDerivationComplete)
	t.Run("test get powpeg redeem script", testBuildPowPegRedeemScript)
	t.Run("test get erp redeem script", testBuildErpRedeemScript)
	t.Run("test get flyover redeem script", testBuildFlyoverRedeemScript)
	t.Run("test get flyover erp redeem script", testBuildFlyoverErpRedeemScript)
	t.Run("test get flyover erp redeem script fallback", testBuildFlyoverErpRedeemScriptFallback)
	t.Run("test get powpeg address hash", testBuildPowPegAddressHash)
	t.Run("test get flyover powpeg address hash", testBuildFlyoverPowPegAddressHash)
	t.Run("test get flyover erp address hash", testBuildFlyoverErpAddressHash)
	t.Run("test get flyover erp address hash fallback", testBuildFlyoverErpAddressHashFallback)
}
