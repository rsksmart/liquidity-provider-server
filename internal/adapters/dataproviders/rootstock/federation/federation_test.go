package federation_test

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/federation"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sort"
	"strings"
	"testing"
)

type testQuote struct {
	BtcRefundAddress            string
	LbcAddress                  string
	LpBtcAddress                string
	QuoteHash                   string
	ExpectedDerivationValueHash string
	ExpectedAddressHash         string
	ExpectedP2SHAddressHash     string
	NetworkParams               *chaincfg.Params
	FedInfo                     rootstock.FederationInfo
}

const (
	federationMainnetAddress       = "3MQ5CLsHWjzz24ZGnHDECjKxvPmdNdyEVX"
	federationTestnetAddress       = "2NCxHG5oK8CWLDrBpTQq6pgKE8jyoB2DpTe"
	federationLegacyMainnetAddress = "3PiQttcndFEfsE9jUjEdqAUeNESHgEa76T"
	federationLegacyTestnetAddress = "2NFGcxdYpEhk251nH9rrWT7TuaaeTVjRDQZ"
	invalidKey                     = "invalidKey"
)

const (
	powpegErpScriptString    = "64522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db53ae670350cd00b27553210257c293086c4d4fe8943deda5f890a37d11bebd140e220faa76258a41d077b4d42103c2660a46aa73078ee6016dee953488566426cf55fc8011edd0085634d75395f92103cd3e383ec6e12719a6c69515e5559bcbe037d0aa24c187e1e26ce932e22ad7b32102370a9838e4d15708ad14a104ee5606b36caaaaf739d833e67770ce9fd9b3ec8054ae68"
	powPegScriptString       = "522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db53ae"
	flyoverErpScriptString   = "20ffe4766f7b5f2fdf374f8ae02270d713c4dcb4b1c5d42bffda61b7f4c1c4c6c97564522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db53ae670350cd00b27553210257c293086c4d4fe8943deda5f890a37d11bebd140e220faa76258a41d077b4d42103c2660a46aa73078ee6016dee953488566426cf55fc8011edd0085634d75395f92103cd3e383ec6e12719a6c69515e5559bcbe037d0aa24c187e1e26ce932e22ad7b32102370a9838e4d15708ad14a104ee5606b36caaaaf739d833e67770ce9fd9b3ec8054ae68"
	flyoverDerivationHash    = "ffe4766f7b5f2fdf374f8ae02270d713c4dcb4b1c5d42bffda61b7f4c1c4c6c9"
	flyoverDerivationAddress = "3PMms7Q3FqUm3kf5YEpL2SGGNewuedm37i"
)

const invalidFailInfoTestName = "fail on invalid fed info"

var testQuotes = []testQuote{
	{
		LpBtcAddress:                "mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk",
		BtcRefundAddress:            "mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk",
		LbcAddress:                  "2ff74F841b95E000625b3A77fed03714874C4fEa",
		QuoteHash:                   "4a3eca107f22707e5dbc79964f3e6c21ec5e354e0903391245d9fdbe6bd2b2f0",
		ExpectedAddressHash:         "2MxeEHVx71taCeVsXFsfQ7TKK6v943PFVEu",
		ExpectedP2SHAddressHash:     "2NEpDCrxcENF6ySQVXvSEDEY7SdSnkNgd1U",
		ExpectedDerivationValueHash: "ff883edd54f8cb22464a8181ed62652fcdb0028e0ada18f9828afd76e0df2c72",
		NetworkParams:               &chaincfg.TestNet3Params,
		FedInfo:                     mocks.GetFakeFedInfo(),
	},
	{
		LpBtcAddress:                "2NDjJznHgtH1rzq63eeFG3SiDi5wxE25FSz",
		BtcRefundAddress:            "2NDjJznHgtH1rzq63eeFG3SiDi5wxE25FSz",
		LbcAddress:                  "2ff74F841b95E000625b3A77fed03714874C4fEa",
		QuoteHash:                   "4a3eca107f22707e5dbc79964f3e6c21ec5e354e0903391245d9fdbe6bd2b2f0",
		ExpectedAddressHash:         "2MwR3FgiXewcWM9xUEWDyZXjCxmT3p5ahrR",
		ExpectedP2SHAddressHash:     "2MwudEEwcB4Quo1NvGMnGkvvHBVTrkVuBTP",
		ExpectedDerivationValueHash: "4cd8a9037f5342217092a9ccc027ab0af1be60bf015e4228afc87214f86f2e51",
		NetworkParams:               &chaincfg.TestNet3Params,
		FedInfo:                     mocks.GetFakeFedInfo(),
	},
	{
		LpBtcAddress:                "17VZNX1SN5NtKa8UQFxwQbFeFc3iqRYhem",
		BtcRefundAddress:            "17VZNX1SN5NtKa8UQFxwQbFeFc3iqRYhem",
		LbcAddress:                  "2ff74F841b95E000625b3A77fed03714874C4fEa",
		QuoteHash:                   "4a3eca107f22707e5dbc79964f3e6c21ec5e354e0903391245d9fdbe6bd2b2f0",
		ExpectedAddressHash:         "3DPRV3zbqGkb3orYADRMrh3AhdZJiPteS4",
		ExpectedP2SHAddressHash:     "3CqwMXq1bLKjgb7XXVr12fsWX9Dmehs5xp",
		ExpectedDerivationValueHash: "f07f644aa9123cd339f232be7f02ec536d40247f6f0c89a93d625ee57918c544",
		NetworkParams:               &chaincfg.MainNetParams,
		FedInfo:                     mocks.GetFakeFedInfo(),
	},
	{
		LpBtcAddress:                "3EktnHQD7RiAE6uzMj2ZifT9YgRrkSgzQX",
		BtcRefundAddress:            "3EktnHQD7RiAE6uzMj2ZifT9YgRrkSgzQX",
		LbcAddress:                  "2ff74F841b95E000625b3A77fed03714874C4fEa",
		QuoteHash:                   "4a3eca107f22707e5dbc79964f3e6c21ec5e354e0903391245d9fdbe6bd2b2f0",
		ExpectedAddressHash:         "35PZtxos7QCwU3mYAiUQRgA5o1QFbkBycx",
		ExpectedP2SHAddressHash:     "33o1WZxf27visUUt3nN6GEVNQWQWw2Fv9S",
		ExpectedDerivationValueHash: "edb9cfe28705fa1619fe1c1bc70e55d5eee4965aea0de631bcf56434a7c454cc",
		NetworkParams:               &chaincfg.MainNetParams,
		FedInfo:                     mocks.GetFakeFedInfo(),
	},
}

func TestDerivationComplete(t *testing.T) {
	for _, test := range testQuotes {

		quoteHash, err := hex.DecodeString(test.QuoteHash)
		require.NoError(t, err)
		userBtcAddress, err := bitcoin.DecodeAddressBase58(test.BtcRefundAddress, true)
		require.NoError(t, err)
		lbcAddress, err := hex.DecodeString(test.LbcAddress)
		require.NoError(t, err)
		lpBtcAddress, err := bitcoin.DecodeAddressBase58(test.LpBtcAddress, true)
		require.NoError(t, err)

		args := rootstock.FlyoverDerivationArgs{
			QuoteHash:            quoteHash,
			UserBtcRefundAddress: userBtcAddress,
			LbcAdress:            lbcAddress,
			LpBtcAddress:         lpBtcAddress,
		}
		derivationHash := federation.GetDerivationValueHash(args)

		fedInfo := mocks.GetFakeFedInfo()
		if test.NetworkParams.Name == chaincfg.TestNet3Params.Name {
			fedInfo.FedAddress = federationTestnetAddress
		} else {
			fedInfo.FedAddress = federationMainnetAddress
		}

		fedRedeemScript, err := federation.GetFedRedeemScript(fedInfo, *test.NetworkParams)
		require.NoError(t, err)

		flyoverRedeemScript := federation.GetFlyoverRedeemScript(derivationHash, fedRedeemScript)
		address, err := bitcoin.ScriptToAddressP2shP2wsh(flyoverRedeemScript, test.NetworkParams)
		require.NoError(t, err)

		err = federation.ValidateRedeemScript(fedInfo, *test.NetworkParams, fedRedeemScript)
		require.NoError(t, err)
		require.EqualValues(t, test.ExpectedAddressHash, address.EncodeAddress())
	}
}

func TestLegacyDerivationComplete(t *testing.T) {
	for _, test := range testQuotes {

		quoteHash, err := hex.DecodeString(test.QuoteHash)
		require.NoError(t, err)
		userBtcAddress, err := bitcoin.DecodeAddressBase58(test.BtcRefundAddress, true)
		require.NoError(t, err)
		lbcAddress, err := hex.DecodeString(test.LbcAddress)
		require.NoError(t, err)
		lpBtcAddress, err := bitcoin.DecodeAddressBase58(test.LpBtcAddress, true)
		require.NoError(t, err)

		args := rootstock.FlyoverDerivationArgs{
			QuoteHash:            quoteHash,
			UserBtcRefundAddress: userBtcAddress,
			LbcAdress:            lbcAddress,
			LpBtcAddress:         lpBtcAddress,
		}
		derivationHash := federation.GetDerivationValueHash(args)

		fedInfo := mocks.GetFakeFedInfo()
		fedInfo.UseSegwit = false
		if test.NetworkParams.Name == chaincfg.TestNet3Params.Name {
			fedInfo.FedAddress = federationLegacyTestnetAddress
		} else {
			fedInfo.FedAddress = federationLegacyMainnetAddress
		}

		fedRedeemScript, err := federation.GetFedRedeemScript(fedInfo, *test.NetworkParams)
		require.NoError(t, err)

		flyoverRedeemScript := federation.GetFlyoverRedeemScript(derivationHash, fedRedeemScript)
		address, err := btcutil.NewAddressScriptHash(flyoverRedeemScript, test.NetworkParams)
		require.NoError(t, err)

		err = federation.ValidateRedeemScript(fedInfo, *test.NetworkParams, fedRedeemScript)
		require.NoError(t, err)
		require.EqualValues(t, test.ExpectedP2SHAddressHash, address.EncodeAddress())
	}
}

func TestGetDerivationValueHash(t *testing.T) {
	for _, test := range testQuotes {
		quoteHash, err := hex.DecodeString(test.QuoteHash)
		require.NoError(t, err)
		userBtcAddress, err := bitcoin.DecodeAddressBase58(test.BtcRefundAddress, true)
		require.NoError(t, err)
		lbcAddress, err := hex.DecodeString(test.LbcAddress)
		require.NoError(t, err)
		lpBtcAddress, err := bitcoin.DecodeAddressBase58(test.LpBtcAddress, true)
		require.NoError(t, err)
		args := rootstock.FlyoverDerivationArgs{
			QuoteHash:            quoteHash,
			UserBtcRefundAddress: userBtcAddress,
			LbcAdress:            lbcAddress,
			LpBtcAddress:         lpBtcAddress,
		}
		hash := hex.EncodeToString(federation.GetDerivationValueHash(args))
		assert.Equal(t, test.ExpectedDerivationValueHash, hash)
	}
}

func TestBuildPowPegRedeemScript(t *testing.T) {
	fedInfo := mocks.GetFakeFedInfo()
	fedRedeemScript, err := federation.GetRedeemScriptBuf(fedInfo, true)
	require.NoError(t, err)

	scriptString := hex.EncodeToString(fedRedeemScript.Bytes())
	assert.True(t, checkSubstrings(scriptString, fedInfo.PubKeys...))

	op2 := fmt.Sprintf("%02x", txscript.OP_2)
	assert.EqualValues(t, scriptString[0:2], op2)

	op3 := fmt.Sprintf("%02x", txscript.OP_3)
	assert.EqualValues(t, scriptString[len(scriptString)-4:len(scriptString)-2], op3)

	sort.Slice(fedInfo.PubKeys, func(i, j int) bool {
		return fedInfo.PubKeys[i] < fedInfo.PubKeys[j]
	})

	buf2, err := federation.GetRedeemScriptBuf(fedInfo, true)
	require.NoError(t, err)
	str2 := hex.EncodeToString(buf2.Bytes())

	assert.EqualValues(t, str2, scriptString)
}

func TestBuildErpRedeemScript(t *testing.T) {
	networkParams := chaincfg.MainNetParams
	fedInfo := mocks.GetFakeFedInfo()
	fedRedeemScript, err := federation.GetErpRedeemScriptBuf(fedInfo, networkParams)
	require.NoError(t, err)

	scriptString := hex.EncodeToString(fedRedeemScript.Bytes())
	assert.True(t, checkSubstrings(scriptString, fedInfo.ErpKeys...))
	assert.EqualValues(t, powpegErpScriptString, scriptString)
}

func TestBuildFlyoverErpRedeemScript(t *testing.T) {
	params := chaincfg.MainNetParams

	fedInfo := mocks.GetFakeFedInfo()
	fedInfo.FedAddress = federationMainnetAddress

	fedRedeemScript, err := federation.GetFedRedeemScript(fedInfo, params)
	require.NoError(t, err)
	derivationBytes, err := hex.DecodeString(flyoverDerivationHash)
	require.NoError(t, err)

	flyoverScript := federation.GetFlyoverRedeemScript(derivationBytes, fedRedeemScript)

	str := hex.EncodeToString(flyoverScript)
	fmt.Println(str)
	assert.True(t, checkSubstrings(str, fedInfo.ErpKeys...))
	assert.EqualValues(t, flyoverErpScriptString, str)
}

func TestBuildFlyoverErpRedeemScriptFallback(t *testing.T) {
	params := chaincfg.MainNetParams

	fedInfo := mocks.GetFakeFedInfo()
	fedInfo.FedAddress = federationMainnetAddress

	fedRedeemScript, err := federation.GetFedRedeemScript(fedInfo, params)
	require.NoError(t, err)
	derivationBytes, err := hex.DecodeString(flyoverDerivationHash)
	require.NoError(t, err)

	flyoverScript := federation.GetFlyoverRedeemScript(derivationBytes, fedRedeemScript)

	str := hex.EncodeToString(flyoverScript)
	assert.True(t, checkSubstrings(str, fedInfo.PubKeys...))
	assert.EqualValues(t, flyoverErpScriptString, str)
}

func TestBuildPowPegAddressHash(t *testing.T) {
	fedInfo := mocks.GetFakeFedInfo()

	buf, err := federation.GetRedeemScriptBuf(fedInfo, true)
	require.NoError(t, err)

	str := hex.EncodeToString(buf.Bytes())
	assert.True(t, checkSubstrings(str, fedInfo.PubKeys...))
	assert.EqualValues(t, powPegScriptString, str)

	address, err := bitcoin.ScriptToAddressP2shP2wsh(buf.Bytes(), &chaincfg.MainNetParams)
	require.NoError(t, err)

	// we expect the address of just the multisig, not the erp part
	assert.EqualValues(t, "3MxszLu3ZPvWY81ZXD5vAzd4yZRgSuQdpN", address.EncodeAddress())
}

func TestBuildFlyoverPowPegAddressHash(t *testing.T) {
	params := chaincfg.MainNetParams

	fedInfo := mocks.GetFakeFedInfo()
	fedInfo.FedAddress = federationMainnetAddress

	fedRedeemScript, err := federation.GetFedRedeemScript(fedInfo, params)
	require.NoError(t, err)
	derivationBytes, err := hex.DecodeString(flyoverDerivationHash)
	require.NoError(t, err)

	flyoverScript := federation.GetFlyoverRedeemScript(derivationBytes, fedRedeemScript)

	str := hex.EncodeToString(flyoverScript)
	assert.True(t, checkSubstrings(str, fedInfo.PubKeys...))
	assert.EqualValues(t, flyoverErpScriptString, str)

	address, err := bitcoin.ScriptToAddressP2shP2wsh(flyoverScript, &chaincfg.MainNetParams)
	require.NoError(t, err)
	expectedAddr := "3DuAtiHNNvbyb4D3pbW8z41VZT7X78TMa8"
	assert.EqualValues(t, expectedAddr, address.EncodeAddress())
}

func TestBuildFlyoverErpAddressHash(t *testing.T) {
	params := chaincfg.MainNetParams
	fedInfo := mocks.GetFakeFedInfo()
	fedInfo.FedAddress = federationMainnetAddress

	fedRedeemScript, err := federation.GetFedRedeemScript(fedInfo, params)
	require.NoError(t, err)
	derivationBytes, err := hex.DecodeString(flyoverDerivationHash)
	require.NoError(t, err)

	flyoverScript := federation.GetFlyoverRedeemScript(derivationBytes, fedRedeemScript)

	str := hex.EncodeToString(flyoverScript)
	assert.True(t, checkSubstrings(str, fedInfo.ErpKeys...))
	assert.EqualValues(t, flyoverErpScriptString, str)

	address, err := btcutil.NewAddressScriptHash(flyoverScript, &chaincfg.MainNetParams)
	require.NoError(t, err)

	assert.EqualValues(t, flyoverDerivationAddress, address.EncodeAddress())
}

func TestBuildFlyoverErpAddressHashFallback(t *testing.T) {
	params := chaincfg.MainNetParams
	fedInfo := mocks.GetFakeFedInfo()
	fedInfo.FedAddress = federationMainnetAddress

	fedRedeemScript, err := federation.GetFedRedeemScript(fedInfo, params)
	require.NoError(t, err)
	derivationBytes, err := hex.DecodeString(flyoverDerivationHash)
	require.NoError(t, err)

	flyoverScript := federation.GetFlyoverRedeemScript(derivationBytes, fedRedeemScript)

	str := hex.EncodeToString(flyoverScript)
	assert.True(t, checkSubstrings(str, fedInfo.PubKeys...))
	assert.EqualValues(t, flyoverErpScriptString, str)

	address, err := btcutil.NewAddressScriptHash(flyoverScript, &chaincfg.MainNetParams)
	require.NoError(t, err)

	assert.EqualValues(t, flyoverDerivationAddress, address.EncodeAddress())
}

func TestGetDerivedBitcoinAddress(t *testing.T) {
	for _, test := range testQuotes {
		params := test.NetworkParams
		fedInfo := mocks.GetFakeFedInfo()

		if params.Name == chaincfg.TestNet3Params.Name {
			fedInfo.FedAddress = federationTestnetAddress
		} else {
			fedInfo.FedAddress = federationMainnetAddress
		}

		quoteHash, err := hex.DecodeString(test.QuoteHash)
		require.NoError(t, err)
		userBtcAddress, err := bitcoin.DecodeAddressBase58(test.BtcRefundAddress, true)
		require.NoError(t, err)
		lbcAddress, err := hex.DecodeString(test.LbcAddress)
		require.NoError(t, err)
		lpBtcAddress, err := bitcoin.DecodeAddressBase58(test.LpBtcAddress, true)
		require.NoError(t, err)
		derivationArgs := rootstock.FlyoverDerivationArgs{
			QuoteHash:            quoteHash,
			UserBtcRefundAddress: userBtcAddress,
			LbcAdress:            lbcAddress,
			LpBtcAddress:         lpBtcAddress,
		}

		fedRedeemScript, err := federation.GetFedRedeemScript(fedInfo, *test.NetworkParams)
		require.NoError(t, err)
		derivationValue := federation.GetDerivationValueHash(derivationArgs)
		derivation, err := federation.CalculateFlyoverDerivationAddress(fedInfo, *params, fedRedeemScript, derivationValue)
		require.NoError(t, err)
		assert.EqualValues(t, test.ExpectedAddressHash, derivation.Address)
	}
}

func TestGetLegacyDerivedBitcoinAddress(t *testing.T) {
	for _, test := range testQuotes {
		params := test.NetworkParams
		fedInfo := mocks.GetFakeFedInfo()

		fedInfo.UseSegwit = false
		if params.Name == chaincfg.TestNet3Params.Name {
			fedInfo.FedAddress = federationLegacyTestnetAddress
		} else {
			fedInfo.FedAddress = federationLegacyMainnetAddress
		}

		quoteHash, err := hex.DecodeString(test.QuoteHash)
		require.NoError(t, err)
		userBtcAddress, err := bitcoin.DecodeAddressBase58(test.BtcRefundAddress, true)
		require.NoError(t, err)
		lbcAddress, err := hex.DecodeString(test.LbcAddress)
		require.NoError(t, err)
		lpBtcAddress, err := bitcoin.DecodeAddressBase58(test.LpBtcAddress, true)
		require.NoError(t, err)
		derivationArgs := rootstock.FlyoverDerivationArgs{
			QuoteHash:            quoteHash,
			UserBtcRefundAddress: userBtcAddress,
			LbcAdress:            lbcAddress,
			LpBtcAddress:         lpBtcAddress,
		}

		fedRedeemScript, err := federation.GetFedRedeemScript(fedInfo, *test.NetworkParams)
		require.NoError(t, err)
		derivationValue := federation.GetDerivationValueHash(derivationArgs)
		derivation, err := federation.CalculateFlyoverDerivationAddress(fedInfo, *params, fedRedeemScript, derivationValue)
		require.NoError(t, err)
		assert.EqualValues(t, test.ExpectedP2SHAddressHash, derivation.Address)
	}
}

func TestCalculateFlyoverDerivationAddress_ErrorHandling(t *testing.T) {
	t.Run("malformed redeem script", func(t *testing.T) {
		derivation, err := federation.CalculateFlyoverDerivationAddress(mocks.GetFakeFedInfo(), chaincfg.TestNet3Params, []byte{1}, []byte{1})
		assert.Equal(t, rootstock.FlyoverDerivation{}, derivation)
		require.Error(t, err)
	})
	t.Run("empty redeem script", func(t *testing.T) {
		derivationBytes, err := hex.DecodeString(flyoverDerivationHash)
		require.NoError(t, err)
		fedInfo := mocks.GetFakeFedInfo()
		fedInfo.FedAddress = federationMainnetAddress
		derivation, err := federation.CalculateFlyoverDerivationAddress(fedInfo, chaincfg.MainNetParams, []byte{}, derivationBytes)
		assert.EqualValues(t, flyoverErpScriptString, derivation.RedeemScript)
		require.NoError(t, err)
	})

	t.Run("invalid address", func(t *testing.T) {
		derivationBytes, err := hex.DecodeString(flyoverDerivationHash)
		require.NoError(t, err)
		fedInfo := mocks.GetFakeFedInfo()
		fedInfo.FedAddress = "invalid"
		derivation, err := federation.CalculateFlyoverDerivationAddress(fedInfo, chaincfg.MainNetParams, []byte{}, derivationBytes)
		assert.Equal(t, rootstock.FlyoverDerivation{}, derivation)
		require.ErrorContains(t, err, "error generating fed redeem script")
	})
}

func TestValidateRedeemScript_ErrorHandling(t *testing.T) {
	t.Run(invalidFailInfoTestName, func(t *testing.T) {
		require.Error(t, federation.ValidateRedeemScript(mocks.GetFakeFedInfo(), chaincfg.MainNetParams, []byte{1}))
	})
	t.Run("fail on invalid script", func(t *testing.T) {
		require.Error(t, federation.ValidateRedeemScript(mocks.GetFakeFedInfo(), chaincfg.MainNetParams, nil))
	})
}

func TestGetFedRedeemScript_ErrorHandling(t *testing.T) {
	t.Run(invalidFailInfoTestName, func(t *testing.T) {
		script, err := federation.GetFedRedeemScript(rootstock.FederationInfo{}, chaincfg.MainNetParams)
		assert.Nil(t, script)
		require.Error(t, err)
	})
}

func TestGetErpRedeemScriptBuf_ErrorHandling(t *testing.T) {
	t.Run(invalidFailInfoTestName, func(t *testing.T) {
		script, err := federation.GetErpRedeemScriptBuf(rootstock.FederationInfo{
			ErpKeys: []string{invalidKey},
		}, chaincfg.MainNetParams)
		assert.Nil(t, script)
		require.Error(t, err)
	})
}

func TestGetRedeemScriptBuf_ErrorHandling(t *testing.T) {
	t.Run(invalidFailInfoTestName, func(t *testing.T) {
		script, err := federation.GetRedeemScriptBuf(rootstock.FederationInfo{PubKeys: []string{invalidKey}}, true)
		assert.Nil(t, script)
		require.Error(t, err)
	})
}

func checkSubstrings(str string, subs ...string) bool {
	for _, sub := range subs {
		if !strings.Contains(str, sub) {
			return false
		}
	}
	return true
}
