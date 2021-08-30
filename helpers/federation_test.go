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
	BTCRefundAddr               string
	LBCAddr                     string
	LPBTCAddr                   string
	QuoteHash                   string
	ExpectedDerivationValueHash string
	ExpectedAddressHash         string
	NetworkParams               *chaincfg.Params
	FedInfo                     *FedInfo
}{
	{
		LPBTCAddr:                   "mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk",
		BTCRefundAddr:               "mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk",
		LBCAddr:                     "2ff74F841b95E000625b3A77fed03714874C4fEa",
		QuoteHash:                   "4a3eca107f22707e5dbc79964f3e6c21ec5e354e0903391245d9fdbe6bd2b2f0",
		ExpectedAddressHash:         "2Mx7jaPHtsgJTbqGnjU5UqBpkekHgfigXay",
		ExpectedDerivationValueHash: "ff883edd54f8cb22464a8181ed62652fcdb0028e0ada18f9828afd76e0df2c72",
		NetworkParams:               &chaincfg.TestNet3Params,
		FedInfo:                     getFakeFedInfo(),
	},
	{
		LPBTCAddr:                   "2NDjJznHgtH1rzq63eeFG3SiDi5wxE25FSz",
		BTCRefundAddr:               "2NDjJznHgtH1rzq63eeFG3SiDi5wxE25FSz",
		LBCAddr:                     "2ff74F841b95E000625b3A77fed03714874C4fEa",
		QuoteHash:                   "4a3eca107f22707e5dbc79964f3e6c21ec5e354e0903391245d9fdbe6bd2b2f0",
		ExpectedAddressHash:         "2N6LxcNDYkKzeyXh7xjZUNZnS9G4Sq3mysi",
		ExpectedDerivationValueHash: "4cd8a9037f5342217092a9ccc027ab0af1be60bf015e4228afc87214f86f2e51",
		NetworkParams:               &chaincfg.TestNet3Params,
		FedInfo:                     getFakeFedInfo(),
	},
	{
		LPBTCAddr:                   "17VZNX1SN5NtKa8UQFxwQbFeFc3iqRYhem",
		BTCRefundAddr:               "17VZNX1SN5NtKa8UQFxwQbFeFc3iqRYhem",
		LBCAddr:                     "2ff74F841b95E000625b3A77fed03714874C4fEa",
		QuoteHash:                   "4a3eca107f22707e5dbc79964f3e6c21ec5e354e0903391245d9fdbe6bd2b2f0",
		ExpectedAddressHash:         "38r8PQdgw5vdebE9h12Eum6saVnWEXxbve",
		ExpectedDerivationValueHash: "f07f644aa9123cd339f232be7f02ec536d40247f6f0c89a93d625ee57918c544",
		NetworkParams:               &chaincfg.MainNetParams,
		FedInfo:                     getFakeFedInfo(),
	},
	{
		LPBTCAddr:                   "3EktnHQD7RiAE6uzMj2ZifT9YgRrkSgzQX",
		BTCRefundAddr:               "3EktnHQD7RiAE6uzMj2ZifT9YgRrkSgzQX",
		LBCAddr:                     "2ff74F841b95E000625b3A77fed03714874C4fEa",
		QuoteHash:                   "4a3eca107f22707e5dbc79964f3e6c21ec5e354e0903391245d9fdbe6bd2b2f0",
		ExpectedAddressHash:         "33P85aACtqezxcGjhrferYkpg6djBtvstk",
		ExpectedDerivationValueHash: "edb9cfe28705fa1619fe1c1bc70e55d5eee4965aea0de631bcf56434a7c454cc",
		NetworkParams:               &chaincfg.MainNetParams,
		FedInfo:                     getFakeFedInfo(),
	},
}

const (
	PowPegScriptString     = "522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db53ae"
	ErpScriptString        = "64522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db536702cd50b27553210257c293086c4d4fe8943deda5f890a37d11bebd140e220faa76258a41d077b4d42103c2660a46aa73078ee6016dee953488566426cf55fc8011edd0085634d75395f92103cd3e383ec6e12719a6c69515e5559bcbe037d0aa24c187e1e26ce932e22ad7b32102370a9838e4d15708ad14a104ee5606b36caaaaf739d833e67770ce9fd9b3ec805468ae"
	FlyoverScriptString    = "20ffe4766f7b5f2fdf374f8ae02270d713c4dcb4b1c5d42bffda61b7f4c1c4c6c975522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db53ae"
	FlyoverErpScriptString = "20ffe4766f7b5f2fdf374f8ae02270d713c4dcb4b1c5d42bffda61b7f4c1c4c6c97564522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db536702cd50b27553210257c293086c4d4fe8943deda5f890a37d11bebd140e220faa76258a41d077b4d42103c2660a46aa73078ee6016dee953488566426cf55fc8011edd0085634d75395f92103cd3e383ec6e12719a6c69515e5559bcbe037d0aa24c187e1e26ce932e22ad7b32102370a9838e4d15708ad14a104ee5606b36caaaaf739d833e67770ce9fd9b3ec805468ae"
)

func testGetDerivationValueHash(t *testing.T) {
	for _, tt := range testQuotes {
		if !common.IsHexAddress(tt.LBCAddr) {
			t.Errorf("invalid address: %v", tt.LBCAddr)
		}
		lbcAddr := common.FromHex(tt.LBCAddr)
		hashBytes, err := hex.DecodeString(tt.QuoteHash)
		if err != nil || len(hashBytes) == 0 {
			t.Errorf("Cannot parse QuoteHash correctly. value: %v, error: %v", tt.QuoteHash, err)
		}
		value, _ := GetDerivationValueHash(tt.BTCRefundAddr, lbcAddr, tt.LPBTCAddr, hashBytes)
		result := hex.EncodeToString(value)
		if result != tt.ExpectedDerivationValueHash {
			t.Errorf("Unexpected derivation value. value: %v, expected: %v, error: %v", result, tt.ExpectedDerivationValueHash, err)
		}
	}
}

func testDerivationComplete(t *testing.T) {
	for _, tt := range testQuotes {
		tt.FedInfo.IrisActivationHeight = 1

		lbcAddr := common.FromHex(tt.LBCAddr)
		hashBytes, _ := hex.DecodeString(tt.QuoteHash)
		value, err := GetDerivationValueHash(tt.BTCRefundAddr, lbcAddr, tt.LPBTCAddr, hashBytes)
		if err != nil {
			t.Errorf("Unexpected error in GetDerivationValueHash. value: %v, expected: %v, error: %v", value, tt.ExpectedDerivationValueHash, err)
		}
		result := hex.EncodeToString(value)
		assert.EqualValues(t, tt.ExpectedDerivationValueHash, result)
		buf, err := getFlyoverPrefix(value)
		if err != nil {
			t.Errorf("Unexpected error in getFlyoverPrefix. error: %v", err)
		}
		scriptBuf, err := getPowPegRedeemScriptBuf(tt.FedInfo, true)
		buf.Write(scriptBuf.Bytes())
		if err != nil {
			t.Errorf("Unexpected error in getPowPegRedeemScriptBuf. error: %v", err)
		}
		addr, err := btcutil.NewAddressScriptHash(buf.Bytes(), tt.NetworkParams)
		if err != nil {
			t.Errorf("Unexpected error in NewAddressScriptHash. error: %v", err)
		}
		assert.EqualValues(t, tt.ExpectedAddressHash, addr.EncodeAddress())
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
	assert.EqualValues(t, str, ErpScriptString)
}

func testBuildFlyoverRedeemScript(t *testing.T) {
	fedInfo := getFakeFedInfo()
	hash, err := getFlyoverDerivationHash()
	if err != nil {
		return
	}
	fedInfo.IrisActivationHeight = 1
	bts, err := GetRedeemScript(fedInfo, hash, &chaincfg.MainNetParams)
	if err != nil {
		return
	}

	str := hex.EncodeToString(bts)
	assert.True(t, checkSubstrings(str, fedInfo.PubKeys...))
	assert.EqualValues(t, FlyoverScriptString, str)
}

func testBuildFlyoverErpRedeemScript(t *testing.T) {
	fedInfo := getFakeFedInfo()
	fedInfo.IrisActivationHeight = -1

	hash, err := getFlyoverDerivationHash()
	if err != nil {
		return
	}

	bts, err := GetRedeemScript(fedInfo, hash, &chaincfg.MainNetParams)
	if err != nil {
		return
	}

	str := hex.EncodeToString(bts)
	assert.True(t, checkSubstrings(str, fedInfo.ErpKeys...))
	assert.EqualValues(t, FlyoverErpScriptString, str)
}

func testBuildPowPegAddressHash(t *testing.T) {
	fedInfo := getFakeFedInfo()
	fedInfo.IrisActivationHeight = 1

	buf, err := getPowPegRedeemScriptBuf(fedInfo, true)
	if err != nil {
		return
	}

	str := hex.EncodeToString(buf.Bytes())
	assert.True(t, checkSubstrings(str, fedInfo.PubKeys...))
	assert.EqualValues(t, str, PowPegScriptString)

	address, err := btcutil.NewAddressScriptHash(buf.Bytes(), &chaincfg.MainNetParams)
	expectedAddr := "3EDhHutH7XnsotnZaTfRr9CwnnGsNNrhCL"

	assert.EqualValues(t, expectedAddr, address.EncodeAddress())
}

func testBuildFlyoverPowPegAddressHash(t *testing.T) {
	fedInfo := getFakeFedInfo()
	fedInfo.IrisActivationHeight = 1
	hash, err := getFlyoverDerivationHash()
	if err != nil {
		return
	}
	bts, err := GetRedeemScript(fedInfo, hash, &chaincfg.MainNetParams)
	if err != nil {
		return
	}

	str := hex.EncodeToString(bts)
	assert.True(t, checkSubstrings(str, fedInfo.PubKeys...))
	assert.EqualValues(t, FlyoverScriptString, str)

	address, err := btcutil.NewAddressScriptHash(bts, &chaincfg.MainNetParams)
	expectedAddr := "34TNebhLLHsE6FHQVMmeHAhTFpaAWhfweR"

	assert.EqualValues(t, expectedAddr, address.EncodeAddress())
}

func testBuildFlyoverErpAddressHash(t *testing.T) {
	fedInfo := getFakeFedInfo()
	hash, err := getFlyoverDerivationHash()
	fedInfo.IrisActivationHeight = -1

	if err != nil {
		return
	}
	bts, err := GetRedeemScript(fedInfo, hash, &chaincfg.MainNetParams)
	if err != nil {
		return
	}

	str := hex.EncodeToString(bts)
	assert.True(t, checkSubstrings(str, fedInfo.ErpKeys...))
	assert.EqualValues(t, FlyoverErpScriptString, str)

	address, err := btcutil.NewAddressScriptHash(bts, &chaincfg.MainNetParams)
	expectedAddr := "3PS2FEphLJMbJURMdYYFNAZR6zLasX51RC"

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
		IrisActivationHeight: 0,
	}
}

func TestFederationHelper(t *testing.T) {
	t.Run("test derivation value hash", testGetDerivationValueHash)
	t.Run("test derivation complete", testDerivationComplete)
	t.Run("test get powpeg redeem script", testBuildPowPegRedeemScript)
	t.Run("test get erp redeem script", testBuildErpRedeemScript)
	t.Run("test get flyover redeem script", testBuildFlyoverRedeemScript)
	t.Run("test get flyover erp redeem script", testBuildFlyoverErpRedeemScript)
	t.Run("test get powpeg address hash", testBuildPowPegAddressHash)
	t.Run("test get flyover powpeg address hash", testBuildFlyoverPowPegAddressHash)
	t.Run("test get flyover erp address hash", testBuildFlyoverErpAddressHash)
}
