package connectors

import (
	"encoding/hex"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/rsksmart/liquidity-provider-server/connectors/testmocks"
	"github.com/stretchr/testify/mock"

	"github.com/btcsuite/btcutil"

	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"

	"sort"
	"strings"

	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
)

type testPmt struct {
	h   string
	pmt string
}

var (
	expectedPmts = [2]testPmt{
		{
			h: "07f8b22fa9a3b32e20b59bb90727de05fb634749519ebcb6a887aeaf2c7eb041",
			pmt: "f3080000" + "0d" +
				"41b07e2cafae87a8b6bc9e51494763fb05de2707b99bb5202eb3a3a92fb2f807" +
				"731c671fafb5d234834c726657f29c9af030ccf7068f1ef732af4efd8e146da0" +
				"a9d6075f4758821ceeef2c230cfd2497df2d1d1d02dd19e653d22b3dc271b393" +
				"94c2d2e51eae99af800c39575a11f3b4eb0fdf5d5deff0b9f5ff592566f4f173" +
				"2b3f6d41bded4344899a348d3053ccd68e922626e589da71d9a583ccfe9e3be6" +
				"393a24e14d18006b54a967963c56da58b18ce770cdb3b32e56d88c138c473a1a" +
				"37acc29a7788a88404fb9a05c416c2ad8f340a61c1ea331528a91ae6210db4f0" +
				"e22d21c55b3f6806387ca34aba54522a7fe15e593ab0d0ff89c6d826cf4e7455" +
				"99bdfed66a4ef1a366f056363158b2907388a5bd4013643d83a016469d392aab" +
				"2319910a0ac801d2c9793c661f1bb02863f672288ce3b4c5e365cff81932fbe4" +
				"3d1eba4b027f80b1240ab5b8c677167602ee63c5a6dad213777b6fbbd3dd9778" +
				"71f30b975d1a8f6cd62535985ea4d11f5c7f80549eab1b18a5cc011872b5403d" +
				"ee666302831a3c64d1604c5c0bec9c796d8dcace974ae97e5837ff0d446d060c" +
				"04ff1f0000",
		},
		{
			h: "ddf5061f9707f0c959bf24278d557b264716672c1b601ec50112d6dfe160d9d3",
			pmt: "f3080000" + "0d" +
				"c0746a357444e9948a18a612e02df5a99240e77f1ff75dd949d5b4038dcf3667" +
				"3a03c716cf722cff7d264c763088ceeb1665f26c6fdd5835d841eeee2f3ece4a" +
				"203a24db8b7a51e4ab0e35a6b4151f6d7f1eef96f32e4fceaac6127521911618" +
				"6efb7fdb763e821f99bd6af8d044cc6feadd7b4716e6938335a3e08548f5a077" +
				"5dd364971faab5cd089cd1fa713e8be658a67a704d39952218f6518e5045d269" +
				"d3d960e1dfd61201c51e601b2c671647267b558d2724bf59c9f007971f06f5dd" +
				"0eab2677f52c996a3f941bef3ec57ebdf22429c37dee5ae68892df30f8acfc22" +
				"5c6fed56bdff34686135e68fda4b716713e60258b6971c03091f25115c008eec" +
				"48a828c75ad7340fadbc368636b4014f6e8386c3990a35620cbddca933a72b02" +
				"d990fb8a602fcda9e1e41120c25f4981362a9dfc7f7ed1f5188482b8ee3f532f" +
				"0ee6234e44af99351ee430f4ac0fa7b71fe9c601c78480b9a97fea305d3abca2" +
				"35a6668846093803e07c48dc9a75be90ed6edd4debb0b7b49bc057e093ad395e" +
				"ee666302831a3c64d1604c5c0bec9c796d8dcace974ae97e5837ff0d446d060c" +
				"04dbd50300",
		},
	}
)

func testSerializeTx(t *testing.T) {
	expected := "01000000010000000000000000000000000000000000000000000000000000000000000000" +
		"ffffffff5f034aa00a1c2f5669614254432f4d696e656420627920797a33313936303538372f2cfabe" +
		"6d6dc0a3751203a336deb817199448996ebcb2a0e537b1ce9254fa3e9c3295ca196b10000000000000" +
		"0010c56e6700d262d24bd851bb829f9f0000ffffffff0401b3cc25000000001976a914536ffa992491" +
		"508dca0354e52f32a3a7a679a53a88ac00000000000000002b6a2952534b424c4f434b3a040c866ad2" +
		"fdb8b59b32dd17059edaeef11d295e279a74ab97125d2500371ce90000000000000000266a24b9e11b" +
		"6dab3e2ca50c1a6b01cf80eccb9d291aab8b095d653e348aa9d94a73964ff5cf1b0000000000000000" +
		"266a24aa21a9ed04f0bac0104f4fa47bec8058f2ebddd292dd85027ab0d6d95288d31f12c5a4b800000000"

	// this is block 0000000000000000000aca0460feaf0661f173b75d4cc824b57233aa7c6b7bc3
	f, err := os.Open("./testdata/test_block")
	if err != nil {
		t.Fatalf("error opening test block file: %v", err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("error reading test file: %v", err)
	}
	s := string(b)
	h, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("error decoding test file: %v", err)
	}

	block, err := btcutil.NewBlockFromBytes(h)
	if err != nil {
		t.Fatalf("error parsing test block: %v", err)
	}
	tx, err := block.Tx(0)
	if err != nil {
		t.Fatalf("error reading transaction from test block: %v", err)
	}
	result, err := serializeTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(result) != expected {
		t.Errorf("serialized tx does not match expected: %x \n----\n %v", result, expected)
	}
}

func testPMTSerialization(t *testing.T) {
	// this is block 0000000000000000000aca0460feaf0661f173b75d4cc824b57233aa7c6b7bc3
	f, err := os.Open("./testdata/test_block")
	if err != nil {
		t.Fatalf("error opening test block file: %v", err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("error reading test file: %v", err)
	}
	s := string(b)
	h, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("error decoding test file: %v", err)
	}

	block, err := btcutil.NewBlockFromBytes(h)
	if err != nil {
		t.Fatalf("error parsing test block: %v", err)
	}

	for _, p := range expectedPmts {
		serializedPMT, err := serializePMT(p.h, block)
		if err != nil {
			t.Errorf("error serializing PMT:\n %v", p.h)
			continue
		}
		result := hex.EncodeToString(serializedPMT)
		if result != p.pmt {
			t.Errorf("expected PMT:\n%v\n is different from serialized PMT:\n%v\n", p.pmt, result)
		}
	}
}

var testPegOutQuotes = []struct {
	BTCRefundAddr               string
	LBCAddr                     string
	LPBTCAddr                   string
	QuoteHash                   string
	UserPubKey                  string
	ExpectedDerivationValueHash string
	ExpectedAddressHash         string
	NetworkParams               string
	DerivationAddress           string
}{
	{
		LPBTCAddr:                   "mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk",
		BTCRefundAddr:               "mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk",
		LBCAddr:                     "2ff74F841b95E000625b3A77fed03714874C4fEa",
		QuoteHash:                   "4a3eca107f22707e5dbc79964f3e6c21ec5e354e0903391245d9fdbe6bd2b2f0",
		UserPubKey:                  "029448e68904c4ebd63a65a0e6b8a887fbd51bacd45092d3279902aba37fb7f994",
		ExpectedAddressHash:         "2Mx7jaPHtsgJTbqGnjU5UqBpkekHgfigXay",
		ExpectedDerivationValueHash: "ff883edd54f8cb22464a8181ed62652fcdb0028e0ada18f9828afd76e0df2c72",
		NetworkParams:               "testnet",
		DerivationAddress:           "2NE1twHkPzqvaJvuQWEUuAG2pUU337H7ots",
	},
	{
		LPBTCAddr:                   "2NDjJznHgtH1rzq63eeFG3SiDi5wxE25FSz",
		BTCRefundAddr:               "2NDjJznHgtH1rzq63eeFG3SiDi5wxE25FSz",
		LBCAddr:                     "2ff74F841b95E000625b3A77fed03714874C4fEa",
		QuoteHash:                   "4a3eca107f22707e5dbc79964f3e6c21ec5e354e0903391245d9fdbe6bd2b2f0",
		UserPubKey:                  "029448e68904c4ebd63a65a0e6b8a887fbd51bacd45092d3279902aba37fb7f994",
		ExpectedAddressHash:         "2N6LxcNDYkKzeyXh7xjZUNZnS9G4Sq3mysi",
		ExpectedDerivationValueHash: "4cd8a9037f5342217092a9ccc027ab0af1be60bf015e4228afc87214f86f2e51",
		NetworkParams:               "testnet",
		DerivationAddress:           "2NE1twHkPzqvaJvuQWEUuAG2pUU337H7ots",
	},
}

var testQuotes = []struct {
	BTCRefundAddr               string
	LBCAddr                     string
	LPBTCAddr                   string
	QuoteHash                   string
	ExpectedDerivationValueHash string
	ExpectedAddressHash         string
	NetworkParams               string
	FedInfo                     *FedInfo
}{
	{
		LPBTCAddr:                   "mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk",
		BTCRefundAddr:               "mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk",
		LBCAddr:                     "2ff74F841b95E000625b3A77fed03714874C4fEa",
		QuoteHash:                   "4a3eca107f22707e5dbc79964f3e6c21ec5e354e0903391245d9fdbe6bd2b2f0",
		ExpectedAddressHash:         "2Mx7jaPHtsgJTbqGnjU5UqBpkekHgfigXay",
		ExpectedDerivationValueHash: "ff883edd54f8cb22464a8181ed62652fcdb0028e0ada18f9828afd76e0df2c72",
		NetworkParams:               "testnet",
		FedInfo:                     getFakeFedInfo(),
	},
	{
		LPBTCAddr:                   "2NDjJznHgtH1rzq63eeFG3SiDi5wxE25FSz",
		BTCRefundAddr:               "2NDjJznHgtH1rzq63eeFG3SiDi5wxE25FSz",
		LBCAddr:                     "2ff74F841b95E000625b3A77fed03714874C4fEa",
		QuoteHash:                   "4a3eca107f22707e5dbc79964f3e6c21ec5e354e0903391245d9fdbe6bd2b2f0",
		ExpectedAddressHash:         "2N6LxcNDYkKzeyXh7xjZUNZnS9G4Sq3mysi",
		ExpectedDerivationValueHash: "4cd8a9037f5342217092a9ccc027ab0af1be60bf015e4228afc87214f86f2e51",
		NetworkParams:               "testnet",
		FedInfo:                     getFakeFedInfo(),
	},
	{
		LPBTCAddr:                   "17VZNX1SN5NtKa8UQFxwQbFeFc3iqRYhem",
		BTCRefundAddr:               "17VZNX1SN5NtKa8UQFxwQbFeFc3iqRYhem",
		LBCAddr:                     "2ff74F841b95E000625b3A77fed03714874C4fEa",
		QuoteHash:                   "4a3eca107f22707e5dbc79964f3e6c21ec5e354e0903391245d9fdbe6bd2b2f0",
		ExpectedAddressHash:         "38r8PQdgw5vdebE9h12Eum6saVnWEXxbve",
		ExpectedDerivationValueHash: "f07f644aa9123cd339f232be7f02ec536d40247f6f0c89a93d625ee57918c544",
		NetworkParams:               "mainnet",
		FedInfo:                     getFakeFedInfo(),
	},
	{
		LPBTCAddr:                   "3EktnHQD7RiAE6uzMj2ZifT9YgRrkSgzQX",
		BTCRefundAddr:               "3EktnHQD7RiAE6uzMj2ZifT9YgRrkSgzQX",
		LBCAddr:                     "2ff74F841b95E000625b3A77fed03714874C4fEa",
		QuoteHash:                   "4a3eca107f22707e5dbc79964f3e6c21ec5e354e0903391245d9fdbe6bd2b2f0",
		ExpectedAddressHash:         "33P85aACtqezxcGjhrferYkpg6djBtvstk",
		ExpectedDerivationValueHash: "edb9cfe28705fa1619fe1c1bc70e55d5eee4965aea0de631bcf56434a7c454cc",
		NetworkParams:               "mainnet",
		FedInfo:                     getFakeFedInfo(),
	},
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

const (
	PowPegScriptString     = "522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db53ae"
	ErpScriptString        = "64522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db536702cd50b27553210257c293086c4d4fe8943deda5f890a37d11bebd140e220faa76258a41d077b4d42103c2660a46aa73078ee6016dee953488566426cf55fc8011edd0085634d75395f92103cd3e383ec6e12719a6c69515e5559bcbe037d0aa24c187e1e26ce932e22ad7b32102370a9838e4d15708ad14a104ee5606b36caaaaf739d833e67770ce9fd9b3ec805468ae"
	FlyoverScriptString    = "20ffe4766f7b5f2fdf374f8ae02270d713c4dcb4b1c5d42bffda61b7f4c1c4c6c975522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db53ae"
	FlyoverErpScriptString = "20ffe4766f7b5f2fdf374f8ae02270d713c4dcb4b1c5d42bffda61b7f4c1c4c6c97564522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db536702cd50b27553210257c293086c4d4fe8943deda5f890a37d11bebd140e220faa76258a41d077b4d42103c2660a46aa73078ee6016dee953488566426cf55fc8011edd0085634d75395f92103cd3e383ec6e12719a6c69515e5559bcbe037d0aa24c187e1e26ce932e22ad7b32102370a9838e4d15708ad14a104ee5606b36caaaaf739d833e67770ce9fd9b3ec805468ae"
)

func testDerivationComplete(t *testing.T) {
	for _, tt := range testQuotes {
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
		scriptBuf, err := btc.getPowPegRedeemScriptBuf(fedInfo, true)
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
		err = btc.validateRedeemScript(fedInfo, scriptBuf.Bytes())
		if err != nil {
			t.Errorf("error in validateRedeemScript: %v", err)
		}
	}
}

func testBuildPowPegRedeemScript(t *testing.T) {
	btc, err := NewBTC("mainnet")
	if err != nil {
		t.Fatalf("error initializing BTC: %v", err)
	}

	fedInfo := getFakeFedInfo()
	buf, err := btc.getPowPegRedeemScriptBuf(fedInfo, true)
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

	buf2, err := btc.getPowPegRedeemScriptBuf(fedInfo, true)
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
	buf, err := btc.getErpRedeemScriptBuf(fedInfo)
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

	hash, err := getFlyoverDerivationHash()
	if err != nil {
		t.Fatalf("error in getFlyoverDerivationHash: %v", err)
	}

	bts, err := btc.getRedeemScript(fedInfo, hash)
	if err != nil {
		t.Fatalf("error in getRedeemScript: %v", err)
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

	hash, err := getFlyoverDerivationHash()
	if err != nil {
		t.Fatalf("error in getFlyoverDerivationHash: %v", err)
	}

	bts, err := btc.getRedeemScript(fedInfo, hash)
	if err != nil {
		t.Fatalf("error in getRedeemScript: %v", err)
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

	hash, err := getFlyoverDerivationHash()
	if err != nil {
		t.Fatalf("error in getFlyoverDerivationHash: %v", err)
	}

	bts, err := btc.getRedeemScript(fedInfo, hash)
	if err != nil {
		t.Fatalf("error in getRedeemScript: %v", err)
	}

	str := hex.EncodeToString(bts)
	assert.True(t, checkSubstrings(str, fedInfo.PubKeys...))
	assert.EqualValues(t, FlyoverScriptString, str)
}

func testBuildPowPegAddressHash(t *testing.T) {
	btc, err := NewBTC("mainnet")
	if err != nil {
		t.Fatalf("error initializing BTC: %v", err)
	}
	fedInfo := getFakeFedInfo()
	fedInfo.IrisActivationHeight = 1

	buf, err := btc.getPowPegRedeemScriptBuf(fedInfo, true)
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

	hash, err := getFlyoverDerivationHash()
	if err != nil {
		t.Fatalf("error in getFlyoverDerivationHash: %v", err)
	}
	bts, err := btc.getRedeemScript(fedInfo, hash)
	if err != nil {
		t.Fatalf("error in getRedeemScript: %v", err)
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

	hash, err := getFlyoverDerivationHash()
	if err != nil {
		t.Fatalf("error in getFlyoverDerivationHash: %v", err)
	}

	bts, err := btc.getRedeemScript(fedInfo, hash)
	if err != nil {
		t.Fatalf("error in getRedeemScript: %v", err)
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

	hash, err := getFlyoverDerivationHash()
	if err != nil {
		t.Fatalf("error in getFlyoverDerivationHash: %v", err)
	}

	bts, err := btc.getRedeemScript(fedInfo, hash)
	if err != nil {
		t.Fatalf("error in getRedeemScript: %v", err)
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

func testGetDerivedBitcoinAddress(t *testing.T) {
	for _, tt := range testQuotes {
		btc, err := NewBTC(tt.NetworkParams)
		if err != nil {
			t.Errorf("error initializing BTC: %v", err)
			continue
		}
		fedInfo := getFakeFedInfo()
		fedInfo.IrisActivationHeight = 1
		if btc.params.Name == chaincfg.TestNet3Params.Name {
			fedInfo.FedAddress = "2N5muMepJizJE1gR7FbHJU6CD18V3BpNF9p"
		} else {
			fedInfo.FedAddress = "3EDhHutH7XnsotnZaTfRr9CwnnGsNNrhCL"
		}
		lbcAddr, err := DecodeRSKAddress(tt.LBCAddr)
		if err != nil {
			t.Errorf("Unexpected error in DecodeRSKAddress. error: %v", err)
			continue
		}
		hashBytes, err := hex.DecodeString(tt.QuoteHash)
		if err != nil || len(hashBytes) == 0 {
			t.Errorf("Cannot parse QuoteHash correctly. value: %v, error: %v", tt.QuoteHash, err)
			continue
		}
		userBtcRefundAddr, err := DecodeBTCAddressWithVersion(tt.BTCRefundAddr)
		if err != nil {
			t.Errorf("Unexpected error in DecodeBTCAddressWithVersion. error: %v", err)
			continue
		}
		lpBtcAddress, err := DecodeBTCAddressWithVersion(tt.LPBTCAddr)
		if err != nil {
			t.Errorf("Unexpected error in DecodeBTCAddressWithVersion. error: %v", err)
			continue
		}
		addr, err := btc.GetDerivedBitcoinAddress(fedInfo, userBtcRefundAddr, lbcAddr, lpBtcAddress, hashBytes)
		if err != nil {
			t.Errorf("Unexpected error in GetDerivedBitcoinAddress. error: %v", err)
			continue
		}
		assert.EqualValues(t, tt.ExpectedAddressHash, addr)
	}
}

func testCheckBtcAddr(t *testing.T) {
	btcClientMock := new(testmocks.BTCClientMock)
	addrWatcherMock := new(testmocks.AddressWatcherMock)
	amountInBtc := float64(1)
	amount, err := btcutil.NewAmount(amountInBtc)
	assert.Nil(t, err)
	var confirmations int64

	btc, err := NewBTC("mainnet")
	if err != nil {
		t.Fatalf("error initializing BTC: %v", err)
	}
	btc.c = btcClientMock

	btcAddr, err := btcutil.DecodeAddress("38r8PQdgw5vdebE9h12Eum6saVnWEXxbve", &btc.params)
	if err != nil {
		t.Fatalf("error initializing BTC: %v", err)
	}

	// check error when retrieving unspent outputs for address
	btcClientMock.On("ListUnspentMinMaxAddresses", 0, 9999, mock.AnythingOfType("[]btcutil.Address")).Return([]btcjson.ListUnspentResult{}, errors.New("ListUnspentMinMaxAddresses failed")).Once()
	err = btc.checkBtcAddr(addrWatcherMock, btcAddr, btcutil.Amount(0), time.Unix(0, 0), &confirmations, func() time.Time { return time.Unix(0, 0) })
	assert.NotNil(t, err)
	assert.EqualValues(t, "error retrieving unspent outputs for address 38r8PQdgw5vdebE9h12Eum6saVnWEXxbve: ListUnspentMinMaxAddresses failed", err.Error())
	btcClientMock.AssertExpectations(t)
	addrWatcherMock.AssertExpectations(t)

	// check happy flow
	confirmations = 0
	btcClientMock.On("ListUnspentMinMaxAddresses", 0, 9999, mock.AnythingOfType("[]btcutil.Address")).Return([]btcjson.ListUnspentResult{{TxID: "0xabc", Confirmations: 1, Amount: amountInBtc}}, nil).Once()
	addrWatcherMock.On("OnNewConfirmation", "0xabc", int64(1), amount).Once()
	err = btc.checkBtcAddr(addrWatcherMock, btcAddr, amount, time.Unix(0, 0), &confirmations, func() time.Time { return time.Unix(0, 0) })
	assert.Nil(t, err)
	assert.EqualValues(t, 1, confirmations)
	btcClientMock.AssertExpectations(t)
	addrWatcherMock.AssertExpectations(t)

	// check happy flow #2: case when agreed amount has been deposited in the second tx (with two UTXOs)
	confirmations = 0
	btcClientMock.On("ListUnspentMinMaxAddresses", 0, 9999, mock.AnythingOfType("[]btcutil.Address")).Return([]btcjson.ListUnspentResult{
		{TxID: "0xabc", Confirmations: 1, Amount: float64(0.98)},
		{TxID: "0xdef", Confirmations: 1, Amount: float64(0.4)}, // \
		{TxID: "0xdef", Confirmations: 1, Amount: float64(0.6)}, // -- these two txs with hash 0xdef are going to be selected
		{TxID: "0xghi", Confirmations: 1, Amount: float64(0.99)},
		{TxID: "0xjkl", Confirmations: 1, Amount: float64(1.1)},
	}, nil).Once()
	addrWatcherMock.On("OnNewConfirmation", "0xdef", int64(1), amount).Once()
	err = btc.checkBtcAddr(addrWatcherMock, btcAddr, amount, time.Unix(0, 0), &confirmations, func() time.Time { return time.Unix(0, 0) })
	assert.Nil(t, err)
	assert.EqualValues(t, 1, confirmations)
	btcClientMock.AssertExpectations(t)
	addrWatcherMock.AssertExpectations(t)

	// check case when time for depositing has elapsed
	confirmations = 0
	btcClientMock.On("ListUnspentMinMaxAddresses", 0, 9999, mock.AnythingOfType("[]btcutil.Address")).Return([]btcjson.ListUnspentResult{{TxID: "0xabc", Confirmations: 1, Amount: float64(0.98)}}, nil).Once()
	addrWatcherMock.On("OnExpire").Once()
	err = btc.checkBtcAddr(addrWatcherMock, btcAddr, amount, time.Unix(0, 0), &confirmations, func() time.Time { return time.Unix(1, 0) })
	assert.NotNil(t, err)
	assert.EqualValues(t, "time for depositing 1 BTC has elapsed; addr: 38r8PQdgw5vdebE9h12Eum6saVnWEXxbve", err.Error())
	assert.EqualValues(t, 0, confirmations)
	btcClientMock.AssertExpectations(t)
	addrWatcherMock.AssertExpectations(t)

	// check case when time for depositing has elapsed, but agreed amount has been deposited
	confirmations = 0
	btcClientMock.On("ListUnspentMinMaxAddresses", 0, 9999, mock.AnythingOfType("[]btcutil.Address")).Return([]btcjson.ListUnspentResult{{TxID: "0xabc", Confirmations: 1, Amount: amountInBtc}}, nil).Times(1)
	addrWatcherMock.On("OnNewConfirmation", "0xabc", int64(1), amount).Times(0)
	err = btc.checkBtcAddr(addrWatcherMock, btcAddr, amount, time.Unix(0, 0), &confirmations, func() time.Time { return time.Unix(1, 0) })
	assert.Nil(t, err)
	assert.EqualValues(t, 1, confirmations)
	btcClientMock.AssertExpectations(t)
	addrWatcherMock.AssertExpectations(t)

	// check case when number of confirmations has not advanced after the previous check
	confirmations = 1
	btcClientMock.On("ListUnspentMinMaxAddresses", 0, 9999, mock.AnythingOfType("[]btcutil.Address")).Return([]btcjson.ListUnspentResult{{TxID: "0xabc", Confirmations: 1, Amount: amountInBtc}}, nil).Times(1)
	addrWatcherMock.On("OnNewConfirmation", "0xabc", int64(1), amount).Times(0)
	err = btc.checkBtcAddr(addrWatcherMock, btcAddr, amount, time.Unix(0, 0), &confirmations, func() time.Time { return time.Unix(0, 0) })
	assert.NotNil(t, err)
	assert.EqualValues(t, "num of confirmations has not advanced; conf: 1", err.Error())
	assert.EqualValues(t, 1, confirmations)
	btcClientMock.AssertExpectations(t)
	addrWatcherMock.AssertExpectations(t)
}

func testComputeDerivationAddress(t *testing.T) {
	for _, tt := range testPegOutQuotes {
		btc, err := NewBTC(tt.NetworkParams)
		assert.Nil(t, err)
		derivationAddress, err := btc.ComputeDerivationAddresss([]byte(tt.UserPubKey), []byte(tt.QuoteHash))
		assert.Nil(t, err)
		assert.Equal(t, tt.DerivationAddress, derivationAddress)
	}
}

func TestBitcoinConnector(t *testing.T) {
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
	t.Run("test pmt serialization", testPMTSerialization)
	t.Run("test tx serialization", testSerializeTx)
	t.Run("test get derived bitcoin address", testGetDerivedBitcoinAddress)
	t.Run("test check btc addr", testCheckBtcAddr)
	t.Run("test compute derivation address", testComputeDerivationAddress)
}
