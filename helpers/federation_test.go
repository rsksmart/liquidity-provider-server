package federation

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/assert"
	"sort"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

var testQuotes = []struct {
	BTCRefundAddr string
	LBCAddr       string
	LPBTCAddr     string
	QuoteHash     string
	Expected      string
}{
	{
		LPBTCAddr:     "746231713634706b667a306368773065753936776c6a677663683779357633397270643465786a337276",
		LBCAddr:       "0xD2244D24FDE5353e4b3ba3b6e05821b456e04d95",
		BTCRefundAddr: "74623171336b3463726832367936367335747575617333386733347775756a7074793466647579726e76",
		QuoteHash:     "4a3eca107f22707e5dbc79964f3e6c21ec5e354e0903391245d9fdbe6bd2b2f0",
		Expected:      "bb0afe83ab039e7ac51c581f0c9909a2c3c02e956c9c0f975baf5baefb39a715",
	},
}

func testGetDerivationValueHash(t *testing.T) {
	for _, tt := range testQuotes {
		btcRefAddr, err := hex.DecodeString(tt.BTCRefundAddr)
		if err != nil || len(btcRefAddr) == 0 {
			t.Errorf("Cannot parse BTCRefundAddr correctly. value: %v, error: %v", tt.BTCRefundAddr, err)
		}

		if !common.IsHexAddress(tt.LBCAddr) {
			t.Errorf("invalid address: %v", tt.LBCAddr)
		}

		lbcAddr := common.FromHex(tt.LBCAddr)
		if err != nil || len(lbcAddr) == 0 {
			t.Errorf("Cannot parse LBCAddr correctly. value: %v, error: %v", tt.LBCAddr, err)
		}

		lpBTCAdrr, err := hex.DecodeString(tt.LPBTCAddr)
		if err != nil || len(lpBTCAdrr) == 0 {
			t.Errorf("Cannot parse LPBTCAddr correctly. value: %v, error: %v", tt.LPBTCAddr, err)
		}
		hashBytes, err := hex.DecodeString(tt.QuoteHash)
		if err != nil || len(hashBytes) == 0 {
			t.Errorf("Cannot parse QuoteHash correctly. value: %v, error: %v", tt.QuoteHash, err)
		}

		value, err := GetDerivationValueHash(btcRefAddr, lbcAddr, lpBTCAdrr, hashBytes)

		if err != nil {
			t.Errorf("Unexpected error for quotehash %v: %v", tt.QuoteHash, err)
		}

		result := hex.EncodeToString(value)
		if result != tt.Expected {
			t.Errorf("Unexpected derivation value. value: %v, expected: %v, error: %v", result, tt.Expected, err)
		}
	}
}

func testBuildPowPegRedeemScript(t *testing.T) {
	fedInfo := getFakeFedInfo()

	buf, err := getPowPegRedeemScriptBuf(fedInfo, true)
	if err != nil {
		return
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

	buf2, err := getPowPegRedeemScriptBuf(fedInfo, true)
	if err != nil {
		return
	}
	str2 := hex.EncodeToString(buf2.Bytes())

	assert.EqualValues(t, str2, str)
}

func testBuildErpRedeemScript(t *testing.T) {
	fedInfo := getFakeFedInfo()

	buf, err := getErpRedeemScriptBuf(fedInfo, &chaincfg.MainNetParams)
	if err != nil {
		return
	}

	str := hex.EncodeToString(buf.Bytes())
	assert.True(t, checkSubstrings(str, fedInfo.ErpKeys...))
	assert.EqualValues(t, str, getErpScriptString())
}

func testBuildFlyoverRedeemScript(t *testing.T) {
	fedInfo := getFakeFedInfo()

	buf, err := getFlyoverRedeemScriptBuf(fedInfo, getFlyoverDerivationHash())
	if err != nil {
		return
	}

	str := hex.EncodeToString(buf.Bytes())
	assert.True(t, checkSubstrings(str, fedInfo.PubKeys...))
	assert.EqualValues(t, getFlyoverScriptString(), str)
}

func testBuildFlyoverErpRedeemScript(t *testing.T) {
	fedInfo := getFakeFedInfo()

	buf, err := getFlyoverErpRedeemScriptBuf(fedInfo, getFlyoverDerivationHash(), &chaincfg.MainNetParams)
	if err != nil {
		return
	}

	str := hex.EncodeToString(buf.Bytes())
	assert.True(t, checkSubstrings(str, fedInfo.ErpKeys...))
	assert.EqualValues(t, getFlyoverErpScriptString(), str)
}

func testBuildPowPegAddressHash(t *testing.T) {
	fedInfo := getFakeFedInfo()

	buf, err := getPowPegRedeemScriptBuf(fedInfo, true)
	if err != nil {
		return
	}

	str := hex.EncodeToString(buf.Bytes())
	assert.True(t, checkSubstrings(str, fedInfo.PubKeys...))
	assert.EqualValues(t, str, getPowPegScriptString())

	address, err := btcutil.NewAddressScriptHash(buf.Bytes(), &chaincfg.MainNetParams)
	powPegAddr := "3EDhHutH7XnsotnZaTfRr9CwnnGsNNrhCL"

	assert.EqualValues(t, powPegAddr, address.EncodeAddress())
}

func getPowPegScriptString() interface{} {
	return "522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db53ae"
}

func getErpScriptString() string {
	return "64522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db536702cd50b27553210257c293086c4d4fe8943deda5f890a37d11bebd140e220faa76258a41d077b4d42103c2660a46aa73078ee6016dee953488566426cf55fc8011edd0085634d75395f92103cd3e383ec6e12719a6c69515e5559bcbe037d0aa24c187e1e26ce932e22ad7b32102370a9838e4d15708ad14a104ee5606b36caaaaf739d833e67770ce9fd9b3ec805468ae"
}

func getFlyoverScriptString() interface{} {
	return "20ffe4766f7b5f2fdf374f8ae02270d713c4dcb4b1c5d42bffda61b7f4c1c4c6c975522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db53ae"
}
func getFlyoverErpScriptString() interface{} {
	return "20ffe4766f7b5f2fdf374f8ae02270d713c4dcb4b1c5d42bffda61b7f4c1c4c6c97564522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db536702cd50b27553210257c293086c4d4fe8943deda5f890a37d11bebd140e220faa76258a41d077b4d42103c2660a46aa73078ee6016dee953488566426cf55fc8011edd0085634d75395f92103cd3e383ec6e12719a6c69515e5559bcbe037d0aa24c187e1e26ce932e22ad7b32102370a9838e4d15708ad14a104ee5606b36caaaaf739d833e67770ce9fd9b3ec805468ae"
}

func getFlyoverDerivationHash() string {
	return "ffe4766f7b5f2fdf374f8ae02270d713c4dcb4b1c5d42bffda61b7f4c1c4c6c9"
}

func checkSubstrings(str string, subs ...string) bool {

	isCompleteMatch := true

	fmt.Printf("String: \"%s\", Substrings: %s\n", str, subs)

	for _, sub := range subs {
		if !strings.Contains(str, sub) {
			isCompleteMatch = false
		}
	}

	return isCompleteMatch
}

func getFakeFedInfo() *FedInfo {
	var keys []string
	keys = append(keys, "02cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1")
	keys = append(keys, "0362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a124")
	keys = append(keys, "03c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db")

	var erpPubKeys []string
	erpPubKeys = append(erpPubKeys, "0257c293086c4d4fe8943deda5f890a37d11bebd140e220faa76258a41d077b4d4")
	erpPubKeys = append(erpPubKeys, "03c2660a46aa73078ee6016dee953488566426cf55fc8011edd0085634d75395f9")
	erpPubKeys = append(erpPubKeys, "03cd3e383ec6e12719a6c69515e5559bcbe037d0aa24c187e1e26ce932e22ad7b3")
	erpPubKeys = append(erpPubKeys, "02370a9838e4d15708ad14a104ee5606b36caaaaf739d833e67770ce9fd9b3ec80")

	return &FedInfo{
		ActiveFedBlockHeight: 0,
		ErpKeys:              erpPubKeys,
		FedSize:              len(keys),
		FedThreshold:         len(keys)/2 + 1,
		PubKeys:              keys,
		IrisActivationHeight: -1,
	}
}

func TestFederationHelper(t *testing.T) {
	t.Run("get derivation value hash", testGetDerivationValueHash)
	t.Run("test get powpeg redeem script", testBuildPowPegRedeemScript)
	t.Run("test get erp redeem script", testBuildErpRedeemScript)
	t.Run("test get flyover redeem script", testBuildFlyoverRedeemScript)
	t.Run("test get flyover erp redeem script", testBuildFlyoverErpRedeemScript)
	t.Run("test get powpeg address hash", testBuildPowPegAddressHash)

}
