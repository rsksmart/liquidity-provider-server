package mempool_space_test

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin/mempool_space"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/datasets"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"testing"
	"time"
)

const (
	mainnetUrl = "https://mempool.space/api"
	testnetUrl = "https://mempool.space/testnet/api"
	regtestUrl = "http://localhost:1234/api"
)

const (
	testFilePath = "../../../../../test/mocks/"
)

const (
	testnetTestTxHash          = "9f0706c2717fc77bf0f225a4223933a7decb8d36902ddbb0accab8ea894f8b29"
	flyoverTestnetPegoutTxHash = "9b0c48b79fe40c67f7a2837e6e59a138a16671caf7685dcd831bd3c51b9f6d21"
	testnetWitnessTestTxHash   = "5cadcbc1ccd91b222346f22c9a9a6fdbf20c9338ec8df0b36097e92d029509ec"
	txNotPresentInTestnet      = "5badcbc1ccd91b222346f22c9a9a6fdbf20c9338ec8df0b36097e92d029509ec"
)

const (
	mainnetBlockFile = "block-696394-mainnet.txt"
)

func TestMempoolSpaceApi_ValidateAddress(t *testing.T) {
	t.Run("should validate correct address", func(t *testing.T) {
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			return assert.Equal(t, mainnetUrl+"/v1/validate-address/"+test.AnyBtcAddress, r.URL.String())
		})).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"isvalid":true,"address":"tb1q4kgratttzjvkxfmgd95z54qcq7y6hekdm3w56u","scriptPubKey":"0014ad903ead6b149963276869682a54180789abe6cd","isscript":false,"iswitness":true,"witness_version":0,"witness_program":"ad903ead6b149963276869682a54180789abe6cd"}`)),
		}, nil)
		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.MainNetParams, mainnetUrl)
		require.NoError(t, api.ValidateAddress(test.AnyBtcAddress))
		client.AssertExpectations(t)
	})
	t.Run("should return error in incorrect address", func(t *testing.T) {
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			parsedUrl, err := url.Parse(mainnetUrl + "/v1/validate-address/" + test.AnyString)
			require.NoError(t, err)
			return assert.Equal(t, parsedUrl.String(), r.URL.String())
		})).Return(&http.Response{
			StatusCode: http.StatusNotImplemented,
			Body:       io.NopCloser(bytes.NewBufferString(`{"error": "Invalid address"}`)),
		}, nil)
		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.MainNetParams, mainnetUrl)
		require.ErrorContains(t, api.ValidateAddress(test.AnyString), "Invalid address")
		client.AssertExpectations(t)
	})
	t.Run("should handle error in http request", func(t *testing.T) {
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.Anything).Return(nil, assert.AnError)
		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.MainNetParams, mainnetUrl)
		require.ErrorContains(t, api.ValidateAddress(test.AnyString), "unable to validate address")
		client.AssertExpectations(t)
	})
}

func TestMempoolSpaceApi_DecodeAddress(t *testing.T) {
	var decodedAddresses []datasets.DecodedAddress
	decodedAddresses = append(decodedAddresses, datasets.Base58Addresses...)
	decodedAddresses = append(decodedAddresses, datasets.Bech32Addresses...)
	decodedAddresses = append(decodedAddresses, datasets.Bech32mAddresses...)
	api := mempool_space.NewMempoolSpaceApi(&mocks.HttpClientMock{}, &chaincfg.MainNetParams, mainnetUrl)
	cases := decodedAddresses
	for _, c := range cases {
		decoded, err := api.DecodeAddress(c.Address)
		require.NoError(t, err)
		assert.Equal(t, c.Expected, decoded)
	}
}

// nolint:funlen
func TestMempoolSpaceApi_GetTransactionInfo(t *testing.T) {
	t.Run("should get transaction info if is confirmed", func(t *testing.T) {
		t.Run("any testnet transaction", func(t *testing.T) {
			client := &mocks.HttpClientMock{}
			client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
				return r != nil && r.URL.String() == testnetUrl+"/tx/"+testnetTestTxHash
			})).Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"9f0706c2717fc77bf0f225a4223933a7decb8d36902ddbb0accab8ea894f8b29","version":2,"locktime":2582723,"vin":[{"txid":"2e4428f122341d770a9074b3f62040f2d48b78a253ad0c46eabaf3732ac2f7eb","vout":0,"prevout":{"scriptpubkey":"76a91420d0cd81928e083da7e9ee690987df9f64bea9b088ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 20d0cd81928e083da7e9ee690987df9f64bea9b0 OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"miWU44anwP9mr3X4w6WTEAGuw6vkQwZQr1","value":489200},"scriptsig":"473044022062dae13ba281d0cf529b604bb59c1efcd7b83438af34d4a51acc6f31041be18c022044df281e688a52624f45f6c26662349d1f5efedd4d69530e65b7d7cec0d3792d0121038e509bc056004a5da7460b5acd5d4dcb2add41d53817180499e3814290ecc91e","scriptsig_asm":"OP_PUSHBYTES_71 3044022062dae13ba281d0cf529b604bb59c1efcd7b83438af34d4a51acc6f31041be18c022044df281e688a52624f45f6c26662349d1f5efedd4d69530e65b7d7cec0d3792d01 OP_PUSHBYTES_33 038e509bc056004a5da7460b5acd5d4dcb2add41d53817180499e3814290ecc91e","is_coinbase":false,"sequence":4294967293},{"txid":"f548a9fff50f98aae908594b5f0cd1a8ef57abf7a744a64b0f855b21389ff0b5","vout":0,"prevout":{"scriptpubkey":"76a91410c5469550e5b73b4c144e0b4bcd2e65fdeece3488ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 10c5469550e5b73b4c144e0b4bcd2e65fdeece34 OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"mh3dS5cYBfLGLvkGGc1LrFj7Qj7DAsccdT","value":500000},"scriptsig":"47304402206538fc72b896e4c6e807a4daf56191f68dec307c3011d082e69eeb3d45d6d8c302205a329814ab87901ae56a82587e716fa2282ecc665ab203da14d93db71181ecd8012102498a833095175800f40b2c0ab23f108b47a319a94ccea826062bf66c827e91a9","scriptsig_asm":"OP_PUSHBYTES_71 304402206538fc72b896e4c6e807a4daf56191f68dec307c3011d082e69eeb3d45d6d8c302205a329814ab87901ae56a82587e716fa2282ecc665ab203da14d93db71181ecd801 OP_PUSHBYTES_33 02498a833095175800f40b2c0ab23f108b47a319a94ccea826062bf66c827e91a9","is_coinbase":false,"sequence":4294967293}],"vout":[{"scriptpubkey":"76a91473cce22e78ec61cd54a6438ca1210b88561ebcdd88ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 73cce22e78ec61cd54a6438ca1210b88561ebcdd OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"mr5FSft4PQvoWbf9Sf5iQXbw1u445iNksp","value":488600},{"scriptpubkey":"76a9142c81478132b5dda64ffc484a0d225096c4b22ad588ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 2c81478132b5dda64ffc484a0d225096c4b22ad5 OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"mjaGtyj74LYn7gApr17prZxDPDnfuUnRa5","value":500000}],"size":372,"weight":1488,"sigops":8,"fee":600,"status":{"confirmed":true,"block_height":2582756,"block_hash":"00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f","block_time":1710931198}}`)),
			}, nil).Once()
			client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
				return r != nil && r.URL.String() == testnetUrl+"/tx/"+testnetTestTxHash+"/hex"
			})).Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`0200000002ebf7c22a73f3baea460cad53a2788bd4f24020f6b374900a771d3422f128442e000000006a473044022062dae13ba281d0cf529b604bb59c1efcd7b83438af34d4a51acc6f31041be18c022044df281e688a52624f45f6c26662349d1f5efedd4d69530e65b7d7cec0d3792d0121038e509bc056004a5da7460b5acd5d4dcb2add41d53817180499e3814290ecc91efdffffffb5f09f38215b850f4ba644a7f7ab57efa8d10c5f4b5908e9aa980ff5ffa948f5000000006a47304402206538fc72b896e4c6e807a4daf56191f68dec307c3011d082e69eeb3d45d6d8c302205a329814ab87901ae56a82587e716fa2282ecc665ab203da14d93db71181ecd8012102498a833095175800f40b2c0ab23f108b47a319a94ccea826062bf66c827e91a9fdffffff0298740700000000001976a91473cce22e78ec61cd54a6438ca1210b88561ebcdd88ac20a10700000000001976a9142c81478132b5dda64ffc484a0d225096c4b22ad588acc3682700`)),
			}, nil).Once()
			client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
				return r != nil && r.URL.String() == testnetUrl+"/blocks/tip/height"
			})).Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`4550094`)),
			}, nil).Once()
			api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)

			result, err := api.GetTransactionInfo(testnetTestTxHash)
			require.NoError(t, err)
			assert.Equal(t, testnetTestTxHash, result.Hash)
			assert.Equal(t, uint64(1967339), result.Confirmations)
			assert.False(t, result.HasWitness)
			assert.Equal(t, map[string][]*entities.Wei{
				"mr5FSft4PQvoWbf9Sf5iQXbw1u445iNksp": {entities.NewWei(0.004886 * 1e18)}, "mjaGtyj74LYn7gApr17prZxDPDnfuUnRa5": {entities.NewWei(0.005 * 1e18)},
			}, result.Outputs)
			client.AssertExpectations(t)
		})

		t.Run("flyover testnet transaction", func(t *testing.T) {
			client := &mocks.HttpClientMock{}
			client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
				return r != nil && r.URL.String() == testnetUrl+"/tx/"+flyoverTestnetPegoutTxHash
			})).Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"9b0c48b79fe40c67f7a2837e6e59a138a16671caf7685dcd831bd3c51b9f6d21","version":2,"locktime":0,"vin":[{"txid":"42a2276f05fc277a065151805e302e8f71952321ed76edc9a584c3e7cb0ebc87","vout":2,"prevout":{"scriptpubkey":"76a91417fd1347eda7590d18f1e2451cc4ea98fb0f2f2088ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 17fd1347eda7590d18f1e2451cc4ea98fb0f2f20 OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"mhho1SuBpFHiWVJB6J7v9FQXjcRLHw6pKD","value":2463300},"scriptsig":"47304402204992bb143d814f6b560198300ae0e4b0b1223c57a525cedd1a8d71238d8efd7802202f928d3a53f8a263528e7131c0cd5cb25715aad87769d5e47e99ee4876ee181a0121029e17dc44ad33f0f6f54d404d87441ef108873b451d4f236833ea9bb03ea2b04c","scriptsig_asm":"OP_PUSHBYTES_71 304402204992bb143d814f6b560198300ae0e4b0b1223c57a525cedd1a8d71238d8efd7802202f928d3a53f8a263528e7131c0cd5cb25715aad87769d5e47e99ee4876ee181a01 OP_PUSHBYTES_33 029e17dc44ad33f0f6f54d404d87441ef108873b451d4f236833ea9bb03ea2b04c","is_coinbase":false,"sequence":4294967293}],"vout":[{"scriptpubkey":"76a9146e84f2d601c6742a94bf9ba32bea7e3f3e377fbb88ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 6e84f2d601c6742a94bf9ba32bea7e3f3e377fbb OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"mqbKtarYKnoEdPheFFDGRjksvEpb2vJGNh","value":500000},{"scriptpubkey":"6a20ef4f66bcf4f55bc72c2c8e01894fe944a5fa3be3a8b2a2a474f0447838ddc1c2","scriptpubkey_asm":"OP_RETURN OP_PUSHBYTES_32 ef4f66bcf4f55bc72c2c8e01894fe944a5fa3be3a8b2a2a474f0447838ddc1c2","scriptpubkey_type":"op_return","value":0},{"scriptpubkey":"76a9145c6dbcd0321aadc2b51b564825aa444c886030b988ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 5c6dbcd0321aadc2b51b564825aa444c886030b9 OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"mowfvQDraTDvRgZowL4tx5EatL1u78w65v","value":1956600}],"size":268,"weight":1072,"sigops":8,"fee":6700,"status":{"confirmed":true,"block_height":2581763,"block_hash":"0000000000002d82f47f76d9b877af8f264504d6e0f89b82e89d2d84f64f269a","block_time":1710330877}}`)),
			}, nil).Once()
			client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
				return r != nil && r.URL.String() == testnetUrl+"/tx/"+flyoverTestnetPegoutTxHash+"/hex"
			})).Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`020000000187bc0ecbe7c384a5c9ed76ed212395718f2e305e805151067a27fc056f27a242020000006a47304402204992bb143d814f6b560198300ae0e4b0b1223c57a525cedd1a8d71238d8efd7802202f928d3a53f8a263528e7131c0cd5cb25715aad87769d5e47e99ee4876ee181a0121029e17dc44ad33f0f6f54d404d87441ef108873b451d4f236833ea9bb03ea2b04cfdffffff0320a10700000000001976a9146e84f2d601c6742a94bf9ba32bea7e3f3e377fbb88ac0000000000000000226a20ef4f66bcf4f55bc72c2c8e01894fe944a5fa3be3a8b2a2a474f0447838ddc1c2f8da1d00000000001976a9145c6dbcd0321aadc2b51b564825aa444c886030b988ac00000000`)),
			}, nil).Once()
			client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
				return r != nil && r.URL.String() == testnetUrl+"/blocks/tip/height"
			})).Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`4550094`)),
			}, nil).Once()

			api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
			result, err := api.GetTransactionInfo(flyoverTestnetPegoutTxHash)
			require.NoError(t, err)
			assert.Equal(t, flyoverTestnetPegoutTxHash, result.Hash)
			assert.Equal(t, uint64(1968332), result.Confirmations)
			assert.False(t, result.HasWitness)
			assert.Equal(t, map[string][]*entities.Wei{
				"mqbKtarYKnoEdPheFFDGRjksvEpb2vJGNh": {entities.NewWei(0.005 * 1e18)},
				"mowfvQDraTDvRgZowL4tx5EatL1u78w65v": {entities.NewWei(0.01956600 * 1e18)},
				"":                                   {entities.NewWei(0)}, // Null data script output
			}, result.Outputs)
			client.AssertExpectations(t)
		})

		t.Run("testnet witness transaction", func(t *testing.T) {
			client := &mocks.HttpClientMock{}
			client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
				return r != nil && r.URL.String() == testnetUrl+"/tx/"+testnetWitnessTestTxHash
			})).Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"5cadcbc1ccd91b222346f22c9a9a6fdbf20c9338ec8df0b36097e92d029509ec","version":2,"locktime":0,"vin":[{"txid":"1ac616516a8d94666293d5b36342c06be60d8df437e32278a22e02076a1b3eba","vout":7,"prevout":{"scriptpubkey":"5120588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1ptz8dytlnvmpx98q6mnj5x9dlc7pf06kykzqk4qgfya6mptsg3plsjkyqka","value":14246},"scriptsig":"","scriptsig_asm":"","witness":["881b79ef50067a758bd097077acd5df935cef7fb7bc670c24d537a2d4d72c0c0c65762e6dd4c240bf589f7e8928d413c3a3020f7f4d6f0080c54e012336b09f3","20d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555ac0063036f726401010a746578742f706c61696e000d656c7261756c69746f2e62746368","c0d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555"],"is_coinbase":false,"sequence":4294967295,"inner_witnessscript_asm":"OP_PUSHBYTES_32 d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555 OP_CHECKSIG OP_0 OP_IF OP_PUSHBYTES_3 6f7264 OP_PUSHBYTES_1 01 OP_PUSHBYTES_10 746578742f706c61696e OP_0 OP_PUSHBYTES_13 656c7261756c69746f2e627463 OP_ENDIF"}],"vout":[{"scriptpubkey":"5120f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1p77g2j6wy827qanf5enq6yfzv0lhdvgahqkpexzpykkql4e5gtmns5cymxz","value":546}],"size":266,"weight":548,"sigops":0,"fee":13700,"status":{"confirmed":true,"block_height":2582756,"block_hash":"00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f","block_time":1710931198}}`)),
			}, nil).Once()
			client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
				return r != nil && r.URL.String() == testnetUrl+"/tx/"+testnetWitnessTestTxHash+"/hex"
			})).Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`02000000000101ba3e1b6a07022ea27822e337f48d0de66bc04263b3d5936266948d6a5116c61a0700000000ffffffff012202000000000000225120f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee70340881b79ef50067a758bd097077acd5df935cef7fb7bc670c24d537a2d4d72c0c0c65762e6dd4c240bf589f7e8928d413c3a3020f7f4d6f0080c54e012336b09f34520d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555ac0063036f726401010a746578742f706c61696e000d656c7261756c69746f2e6274636821c0d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d55500000000`)),
			}, nil).Once()
			client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
				return r != nil && r.URL.String() == testnetUrl+"/blocks/tip/height"
			})).Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`4550094`)),
			}, nil).Once()

			api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
			result, err := api.GetTransactionInfo(testnetWitnessTestTxHash)
			require.NoError(t, err)
			assert.Equal(t, testnetWitnessTestTxHash, result.Hash)
			assert.True(t, result.HasWitness)
			assert.Equal(t, uint64(1967339), result.Confirmations)
			assert.Equal(t, map[string][]*entities.Wei{
				"tb1p77g2j6wy827qanf5enq6yfzv0lhdvgahqkpexzpykkql4e5gtmns5cymxz": {entities.NewWei(0.00000546 * 1e18)},
			}, result.Outputs)
			client.AssertExpectations(t)
		})
	})
}

func TestMempoolSpaceApi_GetTransactionInfo_ExpectedErrorHandling(t *testing.T) {
	t.Run("should get transaction info if is not confirmed", func(t *testing.T) {
		client := &mocks.HttpClientMock{}
		const unconfirmedTx = `fc76260e930dfc653b6d705d7c60a5566fde25061bcc3ab9d1c48138dd50ad0f`
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			return r != nil && r.URL.String() == testnetUrl+"/tx/"+unconfirmedTx
		})).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"fc76260e930dfc653b6d705d7c60a5566fde25061bcc3ab9d1c48138dd50ad0f","version":2,"locktime":0,"vin":[{"txid":"17fbf6f49e76eb196b4dd88d3e03a642aa0c37c36c729d3298d970288ff74398","vout":1,"prevout":{"scriptpubkey":"51203a9324a88db3c132354a4c6a3679d333a20d02bde5528aab67fbfb0202d59409","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 3a9324a88db3c132354a4c6a3679d333a20d02bde5528aab67fbfb0202d59409","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1p82fjf2ydk0qnyd22f34rv7wnxw3q6q4au4fg42m8l0asyqk4jsysvl9rlv","value":8035518},"scriptsig":"","scriptsig_asm":"","witness":["11351817526e1fd4be9b2701753a3f57ef41a76503a41d6923a9afdcf7822721523fd2db13163a6413b30c09369aa1d25628d259cf507cfb0bab0f8e43809254"],"is_coinbase":false,"sequence":4294967293}],"vout":[{"scriptpubkey":"5120c25276cda6ba8b1330ef55d445b89b44b5acc4c34747d2d8a6da8d32c0083a15","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 c25276cda6ba8b1330ef55d445b89b44b5acc4c34747d2d8a6da8d32c0083a15","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1pcff8dndxh293xv802h2ytwymgj66e3xrgara9k9xm2xn9sqg8g2shjzrft","value":69407},{"scriptpubkey":"51203a9324a88db3c132354a4c6a3679d333a20d02bde5528aab67fbfb0202d59409","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 3a9324a88db3c132354a4c6a3679d333a20d02bde5528aab67fbfb0202d59409","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1p82fjf2ydk0qnyd22f34rv7wnxw3q6q4au4fg42m8l0asyqk4jsysvl9rlv","value":7964667}],"size":205,"weight":616,"sigops":0,"fee":1444,"status":{"confirmed":false}}`)),
		}, nil).Once()
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			return r != nil && r.URL.String() == testnetUrl+"/tx/"+unconfirmedTx+"/hex"
		})).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`020000000001019843f78f2870d998329d726cc3370caa42a6033e8dd84d6b19eb769ef4f6fb170100000000fdffffff021f0f010000000000225120c25276cda6ba8b1330ef55d445b89b44b5acc4c34747d2d8a6da8d32c0083a15fb877900000000002251203a9324a88db3c132354a4c6a3679d333a20d02bde5528aab67fbfb0202d59409014011351817526e1fd4be9b2701753a3f57ef41a76503a41d6923a9afdcf7822721523fd2db13163a6413b30c09369aa1d25628d259cf507cfb0bab0f8e4380925400000000`)),
		}, nil).Once()

		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
		result, err := api.GetTransactionInfo(unconfirmedTx)
		require.NoError(t, err)
		assert.Equal(t, unconfirmedTx, result.Hash)
		assert.True(t, result.HasWitness)
		assert.Zero(t, result.Confirmations)
		assert.Equal(t, map[string][]*entities.Wei{
			"tb1pcff8dndxh293xv802h2ytwymgj66e3xrgara9k9xm2xn9sqg8g2shjzrft": {entities.NewWei(0.00069407 * 1e18)},
			"tb1p82fjf2ydk0qnyd22f34rv7wnxw3q6q4au4fg42m8l0asyqk4jsysvl9rlv": {entities.NewWei(0.07964667 * 1e18)},
		}, result.Outputs)
		client.AssertExpectations(t)
	})
	t.Run("should return error if transaction doesn't exists", func(t *testing.T) {
		client := &mocks.HttpClientMock{}
		const unconfirmedTx = `fc76260e930dfc653b6d705d7c60a5566fde25061bcc3ab9d1c48138dd50ad0f`
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			return r != nil && r.URL.String() == testnetUrl+"/tx/"+unconfirmedTx
		})).Return(&http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(bytes.NewBufferString(`Transaction not found`)),
		}, nil).Once()

		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
		result, err := api.GetTransactionInfo(unconfirmedTx)
		require.ErrorContains(t, err, "Transaction not found")
		assert.Empty(t, result)
		client.AssertExpectations(t)
	})
}

// nolint:funlen
func TestMempoolSpaceApi_GetTransactionInfo_UnexpectedErrorHandling(t *testing.T) {
	t.Run("should handle error if any request fails", func(t *testing.T) {
		errorSetUps := []func(client *mocks.HttpClientMock){
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(assert.AnError.Error())),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
					return r != nil && r.URL.String() == testnetUrl+"/tx/"+testnetWitnessTestTxHash
				})).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"txid"`)), // not a json
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
					return r != nil && r.URL.String() == testnetUrl+"/tx/"+testnetWitnessTestTxHash
				})).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"5cadcbc1ccd91b222346f22c9a9a6fdbf20c9338ec8df0b36097e92d029509ec","version":2,"locktime":0,"vin":[{"txid":"1ac616516a8d94666293d5b36342c06be60d8df437e32278a22e02076a1b3eba","vout":7,"prevout":{"scriptpubkey":"5120588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1ptz8dytlnvmpx98q6mnj5x9dlc7pf06kykzqk4qgfya6mptsg3plsjkyqka","value":14246},"scriptsig":"","scriptsig_asm":"","witness":["881b79ef50067a758bd097077acd5df935cef7fb7bc670c24d537a2d4d72c0c0c65762e6dd4c240bf589f7e8928d413c3a3020f7f4d6f0080c54e012336b09f3","20d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555ac0063036f726401010a746578742f706c61696e000d656c7261756c69746f2e62746368","c0d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555"],"is_coinbase":false,"sequence":4294967295,"inner_witnessscript_asm":"OP_PUSHBYTES_32 d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555 OP_CHECKSIG OP_0 OP_IF OP_PUSHBYTES_3 6f7264 OP_PUSHBYTES_1 01 OP_PUSHBYTES_10 746578742f706c61696e OP_0 OP_PUSHBYTES_13 656c7261756c69746f2e627463 OP_ENDIF"}],"vout":[{"scriptpubkey":"5120f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1p77g2j6wy827qanf5enq6yfzv0lhdvgahqkpexzpykkql4e5gtmns5cymxz","value":546}],"size":266,"weight":548,"sigops":0,"fee":13700,"status":{"confirmed":true,"block_height":2582756,"block_hash":"00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f","block_time":1710931198}}`)),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(assert.AnError.Error())),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
					return r != nil && r.URL.String() == testnetUrl+"/tx/"+testnetWitnessTestTxHash
				})).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"5cadcbc1ccd91b222346f22c9a9a6fdbf20c9338ec8df0b36097e92d029509ec","version":2,"locktime":0,"vin":[{"txid":"1ac616516a8d94666293d5b36342c06be60d8df437e32278a22e02076a1b3eba","vout":7,"prevout":{"scriptpubkey":"5120588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1ptz8dytlnvmpx98q6mnj5x9dlc7pf06kykzqk4qgfya6mptsg3plsjkyqka","value":14246},"scriptsig":"","scriptsig_asm":"","witness":["881b79ef50067a758bd097077acd5df935cef7fb7bc670c24d537a2d4d72c0c0c65762e6dd4c240bf589f7e8928d413c3a3020f7f4d6f0080c54e012336b09f3","20d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555ac0063036f726401010a746578742f706c61696e000d656c7261756c69746f2e62746368","c0d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555"],"is_coinbase":false,"sequence":4294967295,"inner_witnessscript_asm":"OP_PUSHBYTES_32 d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555 OP_CHECKSIG OP_0 OP_IF OP_PUSHBYTES_3 6f7264 OP_PUSHBYTES_1 01 OP_PUSHBYTES_10 746578742f706c61696e OP_0 OP_PUSHBYTES_13 656c7261756c69746f2e627463 OP_ENDIF"}],"vout":[{"scriptpubkey":"5120f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1p77g2j6wy827qanf5enq6yfzv0lhdvgahqkpexzpykkql4e5gtmns5cymxz","value":546}],"size":266,"weight":548,"sigops":0,"fee":13700,"status":{"confirmed":true,"block_height":2582756,"block_hash":"00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f","block_time":1710931198}}`)),
				}, nil).Once()
				client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
					return r != nil && r.URL.String() == testnetUrl+"/tx/"+testnetWitnessTestTxHash+"/hex"
				})).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`02000000000101ba3e1b6a07022ea27822e337f48d0de66bc04263b3d5936266948d6a5116c61a0700000000ffffffff012202000000000000225120f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee70340881b79ef50067a758bd097077acd5df935cef7fb7bc670c24d537a2d4d72c0c0c65762e6dd4c240bf589f7e8928d413c3a3020f7f4d6f0080c54e012336b09f34520d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555ac0063036f726401010a746578742f706c61696e000d656c7261756c69746f2e6274636821c0d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d55500000000`)),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(assert.AnError.Error())),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
					return r != nil && r.URL.String() == testnetUrl+"/tx/"+testnetWitnessTestTxHash
				})).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"5cadcbc1ccd91b222346f22c9a9a6fdbf20c9338ec8df0b36097e92d029509ec","version":2,"locktime":0,"vin":[{"txid":"1ac616516a8d94666293d5b36342c06be60d8df437e32278a22e02076a1b3eba","vout":7,"prevout":{"scriptpubkey":"5120588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1ptz8dytlnvmpx98q6mnj5x9dlc7pf06kykzqk4qgfya6mptsg3plsjkyqka","value":14246},"scriptsig":"","scriptsig_asm":"","witness":["881b79ef50067a758bd097077acd5df935cef7fb7bc670c24d537a2d4d72c0c0c65762e6dd4c240bf589f7e8928d413c3a3020f7f4d6f0080c54e012336b09f3","20d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555ac0063036f726401010a746578742f706c61696e000d656c7261756c69746f2e62746368","c0d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555"],"is_coinbase":false,"sequence":4294967295,"inner_witnessscript_asm":"OP_PUSHBYTES_32 d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555 OP_CHECKSIG OP_0 OP_IF OP_PUSHBYTES_3 6f7264 OP_PUSHBYTES_1 01 OP_PUSHBYTES_10 746578742f706c61696e OP_0 OP_PUSHBYTES_13 656c7261756c69746f2e627463 OP_ENDIF"}],"vout":[{"scriptpubkey":"5120f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1p77g2j6wy827qanf5enq6yfzv0lhdvgahqkpexzpykkql4e5gtmns5cymxz","value":546}],"size":266,"weight":548,"sigops":0,"fee":13700,"status":{"confirmed":true,"block_height":2582756,"block_hash":"00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f","block_time":1710931198}}`)),
				}, nil).Once()
				client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
					return r != nil && r.URL.String() == testnetUrl+"/tx/"+testnetWitnessTestTxHash+"/hex"
				})).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`not a hex`)),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
					return r != nil && r.URL.String() == testnetUrl+"/tx/"+testnetWitnessTestTxHash
				})).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"5cadcbc1ccd91b222346f22c9a9a6fdbf20c9338ec8df0b36097e92d029509ec","version":2,"locktime":0,"vin":[{"txid":"1ac616516a8d94666293d5b36342c06be60d8df437e32278a22e02076a1b3eba","vout":7,"prevout":{"scriptpubkey":"5120588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1ptz8dytlnvmpx98q6mnj5x9dlc7pf06kykzqk4qgfya6mptsg3plsjkyqka","value":14246},"scriptsig":"","scriptsig_asm":"","witness":["881b79ef50067a758bd097077acd5df935cef7fb7bc670c24d537a2d4d72c0c0c65762e6dd4c240bf589f7e8928d413c3a3020f7f4d6f0080c54e012336b09f3","20d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555ac0063036f726401010a746578742f706c61696e000d656c7261756c69746f2e62746368","c0d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555"],"is_coinbase":false,"sequence":4294967295,"inner_witnessscript_asm":"OP_PUSHBYTES_32 d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555 OP_CHECKSIG OP_0 OP_IF OP_PUSHBYTES_3 6f7264 OP_PUSHBYTES_1 01 OP_PUSHBYTES_10 746578742f706c61696e OP_0 OP_PUSHBYTES_13 656c7261756c69746f2e627463 OP_ENDIF"}],"vout":[{"scriptpubkey":"5120f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1p77g2j6wy827qanf5enq6yfzv0lhdvgahqkpexzpykkql4e5gtmns5cymxz","value":546}],"size":266,"weight":548,"sigops":0,"fee":13700,"status":{"confirmed":true,"block_height":2582756,"block_hash":"00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f","block_time":1710931198}}`)),
				}, nil).Once()
				client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
					return r != nil && r.URL.String() == testnetUrl+"/tx/"+testnetWitnessTestTxHash+"/hex"
				})).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`02000000000101ba3e1b6a07022ea27822e337f48d0de66bc04263b3d5936266948d6a5116c61a0700000000ffffffff012202000000000000225120f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee70340881b79ef50067a758bd097077acd5df935cef7fb7bc670c24d537a2d4d72c0c0c65762e6dd4c240bf589f7e8928d413c3a3020f7f4d6f0080c54e012336b09f34520d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555ac0063036f726401010a746578742f706c61696e000d656c7261756c69746f2e6274636821c0d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d55500000000`)),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("not a number")),
				}, nil).Once()
			},
		}

		for _, setUp := range errorSetUps {
			client := &mocks.HttpClientMock{}
			setUp(client)

			api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
			result, err := api.GetTransactionInfo(testnetWitnessTestTxHash)
			require.Error(t, err)
			assert.Empty(t, result)
			client.AssertExpectations(t)
		}
	})
}

func TestMempoolSpaceApi_GetRawTransaction(t *testing.T) {
	t.Run("should get raw transaction", func(t *testing.T) {
		cases := []datasets.RawTransaction{datasets.BtcCoinbaseTxNoWitness, datasets.BtcSegwitTxNoWitness, datasets.BtcLegacyTxNoWitness}
		mainnetBlock := test.GetBitcoinTestBlock(t, testFilePath+mainnetBlockFile)
		for _, tx := range cases {
			client := &mocks.HttpClientMock{}
			parsedTx, err := mainnetBlock.Tx(tx.Index)
			client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
				return r != nil && r.URL.String() == mainnetUrl+"/tx/"+parsedTx.Hash().String()+"/hex"
			})).Return(&http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString(tx.Tx))}, nil).Once()
			require.NoError(t, err)
			api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.MainNetParams, mainnetUrl)
			result, err := api.GetRawTransaction(parsedTx.Hash().String())
			require.NoError(t, err)
			expectedBytes, err := hex.DecodeString(tx.Tx)
			require.NoError(t, err)
			require.Equal(t, expectedBytes, result)
			client.AssertExpectations(t)
		}
	})
	t.Run("should return error if transaction is not found", func(t *testing.T) {
		client := &mocks.HttpClientMock{}
		const notFoundTx = `fc76260e930dfc653b6d705d7c60a5566fde25061bcc3ab9d1c48138dd50ad0f`
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			return r != nil && r.URL.String() == mainnetUrl+"/tx/"+notFoundTx+"/hex"
		})).Return(&http.Response{StatusCode: http.StatusNotFound, Body: io.NopCloser(bytes.NewBufferString("Transaction not found"))}, nil).Once()
		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.MainNetParams, mainnetUrl)
		result, err := api.GetRawTransaction(notFoundTx)
		require.ErrorContains(t, err, "Transaction not found")
		assert.Empty(t, result)
		client.AssertExpectations(t)
	})
	const errorTx = `fc76260e930dfc653b6d705d7c60a5566fde25061bcc3ab9d1c48138dd50ad0f`
	t.Run("should handle error if request fails", func(t *testing.T) {
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			return r != nil && r.URL.String() == mainnetUrl+"/tx/"+errorTx+"/hex"
		})).Return(&http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(bytes.NewBufferString(assert.AnError.Error()))}, nil).Once()
		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.MainNetParams, mainnetUrl)
		result, err := api.GetRawTransaction(errorTx)
		require.Error(t, err)
		assert.Empty(t, result)
		client.AssertExpectations(t)
	})
	t.Run("should handle wrong format", func(t *testing.T) {
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.Anything).Return(&http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString("not hex"))}, nil).Once()
		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.MainNetParams, mainnetUrl)
		result, err := api.GetRawTransaction(errorTx)
		require.Error(t, err)
		assert.Empty(t, result)
		client.AssertExpectations(t)
	})
}

// nolint:funlen
func TestMempoolSpaceApi_GetPartialMerkleTree(t *testing.T) {
	t.Run("should get partial merkle tree for confirmed transaction", func(t *testing.T) {
		for _, c := range datasets.PartialMerkleTrees {
			client := &mocks.HttpClientMock{}
			client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
				return r != nil && r.URL.String() == mainnetUrl+"/tx/"+c.Tx
			})).Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(c.TxInfo)),
			}, nil).Once()
			client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
				return r != nil && r.URL.String() == mainnetUrl+"/block/"+c.BlockHash
			})).Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(c.BlockInfo)),
			}, nil).Once()
			client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
				return r != nil && r.URL.String() == mainnetUrl+"/block/"+c.BlockHash+"/raw"
			})).RunAndReturn(func(request *http.Request) (*http.Response, error) {
				blockBytes, err := hex.DecodeString(c.BlockHex)
				require.NoError(t, err)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(blockBytes)),
				}, nil
			}).Once()
			api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.MainNetParams, mainnetUrl)
			serializedPMT, err := api.GetPartialMerkleTree(c.Tx)
			require.NoError(t, err)
			result := hex.EncodeToString(serializedPMT)
			assert.Equal(t, c.Pmt, result)
			client.AssertExpectations(t)
		}
	})
	t.Run("should return error when transaction is not confirmed", func(t *testing.T) {
		const unconfirmedTx = "fc76260e930dfc653b6d705d7c60a5566fde25061bcc3ab9d1c48138dd50ad0f"
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			return r != nil && r.URL.String() == testnetUrl+"/tx/"+unconfirmedTx
		})).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"ee4cf13e6ca315e692c015b0ad7f2198854aa3b2789ec5b7829c131ec92e140c","version":2,"locktime":0,"vin":[{"txid":"12db40361e73e9d27a9b5a9dff0c96c81cadae1cf024ee6cab06e8ab6e250b5c","vout":1,"prevout":{"scriptpubkey":"5120cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1pe7vkd2jppwepchr2tk34ehhr6ps4d2gkxgj04ujptpv408ftakzqgsvqqw","value":472098},"scriptsig":"","scriptsig_asm":"","witness":["d775cea49f6d0f4c9d40a5b3ab826eea776e0a1e970a3b3d7198bba5cfafcc969bb91e11f7d2d62f2b2f5d8025cea8f30a4cb3e2b5964f33dd1cc55eb16228fb"],"is_coinbase":false,"sequence":4294967293}],"vout":[{"scriptpubkey":"5120fe452cc5e20ed2f75bc267f7f1375dca2c95ca34dfc06d1c54e14c481feff7ca","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 fe452cc5e20ed2f75bc267f7f1375dca2c95ca34dfc06d1c54e14c481feff7ca","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1plezje30zpmf0wk7zvlmlzd6aegkftj35mlqx68z5u9xys8l07l9qptc4lv","value":3117},{"scriptpubkey":"5120cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1pe7vkd2jppwepchr2tk34ehhr6ps4d2gkxgj04ujptpv408ftakzqgsvqqw","value":468066}],"size":205,"weight":616,"sigops":0,"fee":915,"status":{"confirmed":false}}`)),
		}, nil).Once()
		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
		serializedPMT, err := api.GetPartialMerkleTree(unconfirmedTx)
		require.ErrorIs(t, err, mempool_space.BtcTxNotConfirmedError)
		assert.Empty(t, serializedPMT)
		client.AssertExpectations(t)
	})
	t.Run("should return error if transaction doesn't exists", func(t *testing.T) {
		const tx = "fc76260e930dfc653b6d705d7c60a5566fde25061bcc3ab9d1c48138dd50ad0e"
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			return r != nil && r.URL.String() == testnetUrl+"/tx/"+tx
		})).Return(&http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(bytes.NewBufferString(`Transaction not found`)),
		}, nil).Once()
		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
		serializedPMT, err := api.GetPartialMerkleTree(tx)
		require.ErrorContains(t, err, "Transaction not found")
		assert.Empty(t, serializedPMT)
		client.AssertExpectations(t)
	})
	t.Run("should handle error during any request", func(t *testing.T) {
		errorSetUps := []func(*mocks.HttpClientMock){
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(assert.AnError.Error())),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"ee4cf13e6ca315e692c015b0ad7f2198854aa3b2789ec5b7829c131ec92e140c","version":2,"locktime":0,"vin":[{"txid":"12db40361e73e9d27a9b5a9dff0c96c81cadae1cf024ee6cab06e8ab6e250b5c","vout":1,"prevout":{"scriptpubkey":"5120cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1pe7vkd2jppwepchr2tk34ehhr6ps4d2gkxgj04ujptpv408ftakzqgsvqqw","value":472098},"scriptsig":"","scriptsig_asm":"","witness":["d775cea49f6d0f4c9d40a5b3ab826eea776e0a1e970a3b3d7198bba5cfafcc969bb91e11f7d2d62f2b2f5d8025cea8f30a4cb3e2b5964f33dd1cc55eb16228fb"],"is_coinbase":false,"sequence":4294967293}],"vout":[{"scriptpubkey":"5120fe452cc5e20ed2f75bc267f7f1375dca2c95ca34dfc06d1c54e14c481feff7ca","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 fe452cc5e20ed2f75bc267f7f1375dca2c95ca34dfc06d1c54e14c481feff7ca","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1plezje30zpmf0wk7zvlmlzd6aegkftj35mlqx68z5u9xys8l07l9qptc4lv","value":3117},{"scriptpubkey":"5120cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1pe7vkd2jppwepchr2tk34ehhr6ps4d2gkxgj04ujptpv408ftakzqgsvqqw","value":468066}],"size":205,"weight":616,"sigops":0,"fee":915,"status":{"confirmed":true,"block_height":4550188,"block_hash":"000000000000004e34b46dcb0b48009ebf1522dbf3d24ff2f09361fea722f0f5","block_time":1751887409}}`)),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(assert.AnError.Error())),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"ee4cf13e6ca315`)),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"ee4cf13e6ca315e692c015b0ad7f2198854aa3b2789ec5b7829c131ec92e140c","version":2,"locktime":0,"vin":[{"txid":"12db40361e73e9d27a9b5a9dff0c96c81cadae1cf024ee6cab06e8ab6e250b5c","vout":1,"prevout":{"scriptpubkey":"5120cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1pe7vkd2jppwepchr2tk34ehhr6ps4d2gkxgj04ujptpv408ftakzqgsvqqw","value":472098},"scriptsig":"","scriptsig_asm":"","witness":["d775cea49f6d0f4c9d40a5b3ab826eea776e0a1e970a3b3d7198bba5cfafcc969bb91e11f7d2d62f2b2f5d8025cea8f30a4cb3e2b5964f33dd1cc55eb16228fb"],"is_coinbase":false,"sequence":4294967293}],"vout":[{"scriptpubkey":"5120fe452cc5e20ed2f75bc267f7f1375dca2c95ca34dfc06d1c54e14c481feff7ca","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 fe452cc5e20ed2f75bc267f7f1375dca2c95ca34dfc06d1c54e14c481feff7ca","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1plezje30zpmf0wk7zvlmlzd6aegkftj35mlqx68z5u9xys8l07l9qptc4lv","value":3117},{"scriptpubkey":"5120cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1pe7vkd2jppwepchr2tk34ehhr6ps4d2gkxgj04ujptpv408ftakzqgsvqqw","value":468066}],"size":205,"weight":616,"sigops":0,"fee":915,"status":{"confirmed":true,"block_height":4550188,"block_hash":"000000000000004e34b46dcb0b48009ebf1522dbf3d24ff2f09361fea722f0f5","block_time":1751887409}}`)),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"id":"0000`)),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"ee4cf13e6ca315e692c015b0ad7f2198854aa3b2789ec5b7829c131ec92e140c","version":2,"locktime":0,"vin":[{"txid":"12db40361e73e9d27a9b5a9dff0c96c81cadae1cf024ee6cab06e8ab6e250b5c","vout":1,"prevout":{"scriptpubkey":"5120cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1pe7vkd2jppwepchr2tk34ehhr6ps4d2gkxgj04ujptpv408ftakzqgsvqqw","value":472098},"scriptsig":"","scriptsig_asm":"","witness":["d775cea49f6d0f4c9d40a5b3ab826eea776e0a1e970a3b3d7198bba5cfafcc969bb91e11f7d2d62f2b2f5d8025cea8f30a4cb3e2b5964f33dd1cc55eb16228fb"],"is_coinbase":false,"sequence":4294967293}],"vout":[{"scriptpubkey":"5120fe452cc5e20ed2f75bc267f7f1375dca2c95ca34dfc06d1c54e14c481feff7ca","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 fe452cc5e20ed2f75bc267f7f1375dca2c95ca34dfc06d1c54e14c481feff7ca","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1plezje30zpmf0wk7zvlmlzd6aegkftj35mlqx68z5u9xys8l07l9qptc4lv","value":3117},{"scriptpubkey":"5120cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1pe7vkd2jppwepchr2tk34ehhr6ps4d2gkxgj04ujptpv408ftakzqgsvqqw","value":468066}],"size":205,"weight":616,"sigops":0,"fee":915,"status":{"confirmed":true,"block_height":4550188,"block_hash":"000000000000004e34b46dcb0b48009ebf1522dbf3d24ff2f09361fea722f0f5","block_time":1751887409}}`)),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"id":"000000000000004e34b46dcb0b48009ebf1522dbf3d24ff2f09361fea722f0f5","height":4550188,"version":727244800,"timestamp":1751887409,"tx_count":4972,"size":2113573,"weight":3992581,"merkle_root":"8403359b1c7c5bcf0df2cedf1d05dc6e3bbf3495df06655e3313d47c6f31e537","previousblockhash":"000000000000004fec1a9218d8c41a943d051be49fe342e396f721c83f4083e9","mediantime":1751886171,"nonce":431851372,"bits":424709359,"difficulty":53319353.63456318}`)),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(assert.AnError.Error())),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"ee4cf13e6ca315e692c015b0ad7f2198854aa3b2789ec5b7829c131ec92e140c","version":2,"locktime":0,"vin":[{"txid":"12db40361e73e9d27a9b5a9dff0c96c81cadae1cf024ee6cab06e8ab6e250b5c","vout":1,"prevout":{"scriptpubkey":"5120cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1pe7vkd2jppwepchr2tk34ehhr6ps4d2gkxgj04ujptpv408ftakzqgsvqqw","value":472098},"scriptsig":"","scriptsig_asm":"","witness":["d775cea49f6d0f4c9d40a5b3ab826eea776e0a1e970a3b3d7198bba5cfafcc969bb91e11f7d2d62f2b2f5d8025cea8f30a4cb3e2b5964f33dd1cc55eb16228fb"],"is_coinbase":false,"sequence":4294967293}],"vout":[{"scriptpubkey":"5120fe452cc5e20ed2f75bc267f7f1375dca2c95ca34dfc06d1c54e14c481feff7ca","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 fe452cc5e20ed2f75bc267f7f1375dca2c95ca34dfc06d1c54e14c481feff7ca","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1plezje30zpmf0wk7zvlmlzd6aegkftj35mlqx68z5u9xys8l07l9qptc4lv","value":3117},{"scriptpubkey":"5120cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 cf9966aa410bb21c5c6a5da35cdee3d06156a9163224faf2415859579d2bed84","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1pe7vkd2jppwepchr2tk34ehhr6ps4d2gkxgj04ujptpv408ftakzqgsvqqw","value":468066}],"size":205,"weight":616,"sigops":0,"fee":915,"status":{"confirmed":true,"block_height":4550188,"block_hash":"000000000000004e34b46dcb0b48009ebf1522dbf3d24ff2f09361fea722f0f5","block_time":1751887409}}`)),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"id":"000000000000004e34b46dcb0b48009ebf1522dbf3d24ff2f09361fea722f0f5","height":4550188,"version":727244800,"timestamp":1751887409,"tx_count":4972,"size":2113573,"weight":3992581,"merkle_root":"8403359b1c7c5bcf0df2cedf1d05dc6e3bbf3495df06655e3313d47c6f31e537","previousblockhash":"000000000000004fec1a9218d8c41a943d051be49fe342e396f721c83f4083e9","mediantime":1751886171,"nonce":431851372,"bits":424709359,"difficulty":53319353.63456318}`)),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("not hex")),
				}, nil).Once()
			},
		}

		for _, setUp := range errorSetUps {
			client := &mocks.HttpClientMock{}
			setUp(client)

			api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
			serializedPMT, err := api.GetPartialMerkleTree(testnetTestTxHash)
			require.Error(t, err)
			assert.Empty(t, serializedPMT)
			client.AssertExpectations(t)
		}
	})
}

func TestMempoolSpaceApi_GetHeight(t *testing.T) {
	t.Run("should return block height", func(t *testing.T) {
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			return r != nil && r.URL.String() == mainnetUrl+"/blocks/tip/height"
		})).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("4550189")),
		}, nil).Once()
		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.MainNetParams, mainnetUrl)
		height, err := api.GetHeight()
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(4550189), height)
		client.AssertExpectations(t)
	})
	t.Run("should handle request error", func(t *testing.T) {
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.Anything).Return(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewBufferString(assert.AnError.Error())),
		}, nil).Once()
		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.MainNetParams, mainnetUrl)
		height, err := api.GetHeight()
		require.Error(t, err)
		assert.Nil(t, height)
		client.AssertExpectations(t)
	})
	t.Run("should return error when parsing response", func(t *testing.T) {
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("not a number")),
		}, nil).Once()
		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.MainNetParams, mainnetUrl)
		height, err := api.GetHeight()
		require.Error(t, err)
		assert.Nil(t, height)
		client.AssertExpectations(t)
	})
}

// nolint:funlen
func TestMempoolSpaceApi_BuildMerkleBranch(t *testing.T) {
	t.Run("should build merkle branch for confirmed transaction", func(t *testing.T) {
		t.Run("Should build merkle branch for testnet transactions", func(t *testing.T) {
			t.Run("legacy transaction", func(t *testing.T) {
				client := &mocks.HttpClientMock{}
				client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
					return r != nil && r.URL.String() == testnetUrl+"/tx/"+datasets.TestnetLegacyMerkleBranch.Tx+"/merkle-proof"
				})).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"block_height":2582756,"merkle":["9b50cfbfe0fefecf4e3ef9de5957ab23dccfbd1f46528d30c620a0db84de0408","47aae7204d1eba738e556ca8000d2e3154e988593c2bf3ca903effd58dc2bdb3","3ba28c3cf802f56a9abfeab130eca2b6fbb753eb1d156b7d22721a40a2547e78","5d0792164a78479e9a8dcaa39aa18dfbddcb68484afc15fe409660ac3fa02961","ea74def1c7a2c9db57ae564597c1f78f8ecdf28a143513d0d2329671b54375b1","8db608faddb6b6c07f87725739a966c888b1005387d1cb55ed506deb975c58c0","172c22c4510d2097054b0b68200d97c96323fa8820f69ce8c4c71cd2e3f17443","388592bcb9d117491429daf7d3a5db595087db85c6372f481708dbd13fd3d975","5f0f5095a9745bc91c55e7e8de70910621eb515894bfa5bace7410a5fc300a1d","0d8b34db87e8b3916fdfe388c90c93f91e222980903ed639fcc4e580886253b7","16559e3831c4186ae16d8fa46ac164bcab51e746a0030793e2503b7202fe898a"],"pos":406}`)),
				}, nil)
				testnet := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
				legacyBranch, legacyErr := testnet.BuildMerkleBranch(datasets.TestnetLegacyMerkleBranch.Tx)
				require.NoError(t, legacyErr)
				assert.Equal(t, datasets.TestnetLegacyMerkleBranch.Branch, legacyBranch)
				client.AssertExpectations(t)
			})
			t.Run("witness transaction", func(t *testing.T) {
				client := &mocks.HttpClientMock{}
				client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
					return r != nil && r.URL.String() == testnetUrl+"/tx/"+datasets.TestnetWitnessMerkleBranch.Tx+"/merkle-proof"
				})).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"block_height":2582756,"merkle":["500dd0e2e5cdc1e3e1ef774dc14c2ba8e61a013036fdda4e71417a291d2677f8","23759979b54a349dc67e77864832321e83e3b23619aa70fa37d9f4228b0dbffb","64de91ea81870f42f612bc6b5fb469310c205d90c10573d59473175a539fae0b","4e09a149b086b65e0c40fed1dccb61fecd8b2ce0416d40e362770e02f2d39d39","53f4c10d1713edd245c972f3f595b3362874cf41b944db7a376d22a06e5e833c","825205d9f5f690fcc2e004ca05c2a9fe8e7deb8c96c7a389a2b6efddfa702b70","41e88fddc8c7adf2f67f97b536ab7fa7831d0aff34257b41f6bed0bde1fbf9ca","b619444deb3cb58ba9368aebbc26573bd8043247a147f63a7c0eb4a82ed14141","456f0dc25dd8440165f0839822094a79fb6501cd737f4c2143db53b7225f7d02","0d8b34db87e8b3916fdfe388c90c93f91e222980903ed639fcc4e580886253b7","16559e3831c4186ae16d8fa46ac164bcab51e746a0030793e2503b7202fe898a"],"pos":10}`)),
				}, nil)
				testnet := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
				witnessBranch, witnessErr := testnet.BuildMerkleBranch(datasets.TestnetWitnessMerkleBranch.Tx)
				require.NoError(t, witnessErr)
				assert.Equal(t, datasets.TestnetWitnessMerkleBranch.Branch, witnessBranch)
				client.AssertExpectations(t)
			})
		})
		t.Run("Should build merkle branch for mainnet transactions", func(t *testing.T) {
			t.Run("legacy transaction", func(t *testing.T) {
				client := &mocks.HttpClientMock{}
				client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
					return r != nil && r.URL.String() == mainnetUrl+"/tx/"+datasets.MainnetLegacyMerkleBranch.Tx+"/merkle-proof"
				})).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"block_height":696394,"merkle":["87ee06e07c184ac84a1b25232b2c7ff6e4147fee7d6f3154b780f76bcad6d3a6","6755308192f7a7c718243d06aa39e6d0646eb71c6ac5a0aeb5a8bb75fcba0de6","bcf67bc9e84c63f84a928c16100bb67bd52bc1e603ce410710470c70eaa116f1","e63b9efecc83a5d971da89e52626928ed6cc53308d349a894443edbd416d3f2b","1a3a478c138cd8562eb3b3cd70e78cb158da563c9667a9546b00184de1243a39","f0b40d21e61aa9281533eac1610a348fadc216c4059afb0484a888779ac2ac37","55744ecf26d8c689ffd0b03a595ee17f2a5254ba4aa37c3806683f5bc5212de2","ab2a399d4616a0833d641340bda5887390b258313656f066a3f14e6ad6febd99","e4fb3219f8cf65e3c5b4e38c2872f66328b01b1f663c79c9d201c80a0a911923","7897ddd3bb6f7b7713d2daa6c563ee02761677c6b8b50a24b1807f024bba1e3d","3d40b5721801cca5181bab9e54807f5c1fd1a45e983525d66c8f1a5d970bf371","0c066d440dff37587ee94a97ceca8d6d799cec0b5c4c60d1643c1a83026366ee"],"pos":6}`)),
				}, nil)
				mainnet := mempool_space.NewMempoolSpaceApi(client, &chaincfg.MainNetParams, mainnetUrl)
				legacyBranch, legacyErr := mainnet.BuildMerkleBranch(datasets.MainnetLegacyMerkleBranch.Tx)
				require.NoError(t, legacyErr)
				assert.Equal(t, datasets.MainnetLegacyMerkleBranch.Branch, legacyBranch)
			})
			t.Run("witness transaction", func(t *testing.T) {
				client := &mocks.HttpClientMock{}
				client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
					return r != nil && r.URL.String() == mainnetUrl+"/tx/"+datasets.MainnetWitnessMerkleBranch.Tx+"/merkle-proof"
				})).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"block_height":696394,"merkle":["6cecea4d6a6ca3cc045e431f56060286c77a4becfd5db08620dab17eb5d07ad2","35b223e6bcd4cd31a7789dc1f68c65e2317e5331e6845d750dad1a5f340497d2","f5bad3aa4464fac4869c5179907f305c6af73b815e56fba478a08ca8ec27636b","1116a57baf60c565a855bf772e2f36112dd3be6a98573acba02dd022d9726a7a","99ea21dff2a9d043ccbad843e14b47f36615ae11a6f6fd942b7de94b94d81b4f","f0b40d21e61aa9281533eac1610a348fadc216c4059afb0484a888779ac2ac37","55744ecf26d8c689ffd0b03a595ee17f2a5254ba4aa37c3806683f5bc5212de2","ab2a399d4616a0833d641340bda5887390b258313656f066a3f14e6ad6febd99","e4fb3219f8cf65e3c5b4e38c2872f66328b01b1f663c79c9d201c80a0a911923","7897ddd3bb6f7b7713d2daa6c563ee02761677c6b8b50a24b1807f024bba1e3d","3d40b5721801cca5181bab9e54807f5c1fd1a45e983525d66c8f1a5d970bf371","0c066d440dff37587ee94a97ceca8d6d799cec0b5c4c60d1643c1a83026366ee"],"pos":16}`)),
				}, nil)
				mainnet := mempool_space.NewMempoolSpaceApi(client, &chaincfg.MainNetParams, mainnetUrl)
				witnessBranch, witnessErr := mainnet.BuildMerkleBranch(datasets.MainnetWitnessMerkleBranch.Tx)
				require.NoError(t, witnessErr)
				assert.Equal(t, datasets.MainnetWitnessMerkleBranch.Branch, witnessBranch)
			})
		})
	})
	t.Run("should return error when transaction is not confirmed or doesn't exist", func(t *testing.T) {
		const notFoundError = "Transaction not found or is unconfirmed"
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.Anything).Return(&http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(bytes.NewBufferString(notFoundError)),
		}, nil)
		mainnet := mempool_space.NewMempoolSpaceApi(client, &chaincfg.MainNetParams, mainnetUrl)
		branch, err := mainnet.BuildMerkleBranch(datasets.MainnetWitnessMerkleBranch.Tx)
		require.ErrorContains(t, err, notFoundError)
		assert.Empty(t, branch)
	})
	t.Run("should handle error during any request", func(t *testing.T) {
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.Anything).Return(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewBufferString(assert.AnError.Error())),
		}, nil)
		mainnet := mempool_space.NewMempoolSpaceApi(client, &chaincfg.MainNetParams, mainnetUrl)
		branch, err := mainnet.BuildMerkleBranch(datasets.MainnetWitnessMerkleBranch.Tx)
		require.Error(t, err)
		assert.Empty(t, branch)
	})
}

// nolint:funlen
func TestMempoolSpaceApi_GetTransactionBlockInfo(t *testing.T) {
	t.Run("should return block information for confirmed transaction", func(t *testing.T) {
		const blockHash = "00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f"
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			return r != nil && r.URL.String() == testnetUrl+"/tx/"+testnetTestTxHash
		})).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"9f0706c2717fc77bf0f225a4223933a7decb8d36902ddbb0accab8ea894f8b29","version":2,"locktime":2582723,"vin":[{"txid":"2e4428f122341d770a9074b3f62040f2d48b78a253ad0c46eabaf3732ac2f7eb","vout":0,"prevout":{"scriptpubkey":"76a91420d0cd81928e083da7e9ee690987df9f64bea9b088ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 20d0cd81928e083da7e9ee690987df9f64bea9b0 OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"miWU44anwP9mr3X4w6WTEAGuw6vkQwZQr1","value":489200},"scriptsig":"473044022062dae13ba281d0cf529b604bb59c1efcd7b83438af34d4a51acc6f31041be18c022044df281e688a52624f45f6c26662349d1f5efedd4d69530e65b7d7cec0d3792d0121038e509bc056004a5da7460b5acd5d4dcb2add41d53817180499e3814290ecc91e","scriptsig_asm":"OP_PUSHBYTES_71 3044022062dae13ba281d0cf529b604bb59c1efcd7b83438af34d4a51acc6f31041be18c022044df281e688a52624f45f6c26662349d1f5efedd4d69530e65b7d7cec0d3792d01 OP_PUSHBYTES_33 038e509bc056004a5da7460b5acd5d4dcb2add41d53817180499e3814290ecc91e","is_coinbase":false,"sequence":4294967293},{"txid":"f548a9fff50f98aae908594b5f0cd1a8ef57abf7a744a64b0f855b21389ff0b5","vout":0,"prevout":{"scriptpubkey":"76a91410c5469550e5b73b4c144e0b4bcd2e65fdeece3488ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 10c5469550e5b73b4c144e0b4bcd2e65fdeece34 OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"mh3dS5cYBfLGLvkGGc1LrFj7Qj7DAsccdT","value":500000},"scriptsig":"47304402206538fc72b896e4c6e807a4daf56191f68dec307c3011d082e69eeb3d45d6d8c302205a329814ab87901ae56a82587e716fa2282ecc665ab203da14d93db71181ecd8012102498a833095175800f40b2c0ab23f108b47a319a94ccea826062bf66c827e91a9","scriptsig_asm":"OP_PUSHBYTES_71 304402206538fc72b896e4c6e807a4daf56191f68dec307c3011d082e69eeb3d45d6d8c302205a329814ab87901ae56a82587e716fa2282ecc665ab203da14d93db71181ecd801 OP_PUSHBYTES_33 02498a833095175800f40b2c0ab23f108b47a319a94ccea826062bf66c827e91a9","is_coinbase":false,"sequence":4294967293}],"vout":[{"scriptpubkey":"76a91473cce22e78ec61cd54a6438ca1210b88561ebcdd88ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 73cce22e78ec61cd54a6438ca1210b88561ebcdd OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"mr5FSft4PQvoWbf9Sf5iQXbw1u445iNksp","value":488600},{"scriptpubkey":"76a9142c81478132b5dda64ffc484a0d225096c4b22ad588ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 2c81478132b5dda64ffc484a0d225096c4b22ad5 OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"mjaGtyj74LYn7gApr17prZxDPDnfuUnRa5","value":500000}],"size":372,"weight":1488,"sigops":8,"fee":600,"status":{"confirmed":true,"block_height":2582756,"block_hash":"00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f","block_time":1710931198}}`)),
		}, nil).Once()
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			return r != nil && r.URL.String() == testnetUrl+"/block/"+blockHash
		})).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"id":"00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f","height":2582756,"version":536870912,"timestamp":1710931198,"tx_count":1030,"size":357157,"weight":848203,"merkle_root":"705c1425cc8ff977778acc9729fbd2dc4a8c21358f3e0b4fefa25ce5c343a3df","previousblockhash":"0000000000000014f9d3bd63b0257095d7d272134674285665f00d572fd2b0d8","mediantime":1710925891,"nonce":488178752,"bits":486604799,"difficulty":1.0}`)),
		}, nil).Once()
		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
		blockInfo, err := api.GetTransactionBlockInfo(testnetTestTxHash)
		require.NoError(t, err)
		blockHashBytes, err := hex.DecodeString(blockHash)
		require.NoError(t, err)
		assert.Equal(t, blockchain.BitcoinBlockInformation{
			Hash:   utils.To32Bytes(blockHashBytes),
			Height: big.NewInt(2582756),
			Time:   time.Unix(1710931198, 0),
		}, blockInfo)
		client.AssertExpectations(t)
	})
	t.Run("should return error for unconfirmed transaction", func(t *testing.T) {
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"7c0280790b24ecaad752a2be4982d2a32b6192b23a4f7fa8cef813dee6bdc873","version":2,"locktime":0,"vin":[{"txid":"76031ce026ac8cc3c1dee31fd402f2bfaf42add24a03a132cae7a6b318186c25","vout":0,"prevout":{"scriptpubkey":"512008fac5369147e53f6584ac9e437b4cc8e3a1dc55d7160f503915eee9d2700c0e","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 08fac5369147e53f6584ac9e437b4cc8e3a1dc55d7160f503915eee9d2700c0e","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1pprav2d53gljn7evy4j0yx76ver36rhz46utq75pezhhwn5nsps8q8zteyp","value":2289},"scriptsig":"","scriptsig_asm":"","witness":["7a0b58be893a250638cb2c95bf993ebe00b60779a4597b7c1ef0e76552c823ce","1102380d113cb61abf75d057e7da5dbcd9b3391aebdd8ddce97b24eba3165df4091f075c8cb9ac56a5bd80ae8b5c34799ac7cd54ad01c65b3768435989855f01","fcafa6b5806be712005ba3f5d00dcc65d8a6e4d16a685fce827ec97aa6db4823368044122f1e1b574f59fa3b60cec743d3994bc247df2835a75fcb016d8a3160","0c0300000000000000000000016b20000000000000002bad59995e2ae358ca0ee9e4ae433d4db67d88b8e86aa7d62d6b2039d771ed264db1b6ecc1268207b416c020c40c1f08c92c4da0780584f93efadc76aa201b6af0a8751b294c347ebb786029ba9fa3179dcfc1b6e0e042eec99e1b06681d88ad203807da2e6be6d76ee01c6db93f2e22d8ba66c17547917f1a06b5351dde265385ada9143ae86613ff51c5950a452bd7d07c99319801944e8874519c63026f704f3e1f8b08000000000002039bcfbd6325ff919f1b7c0f9ff93e73d5839d0b1c3639f59d60593a8b73da2cf7dcf5d363847e2dfb4f000000b193f10944000000675168","c139d771ed264db1b6ecc1268207b416c020c40c1f08c92c4da0780584f93efadc84654164130538fbb3f7b1539f3c5461c509d9f7a233c6fb59104ef8f0f8d7e5"],"is_coinbase":false,"sequence":4294967293,"inner_witnessscript_asm":"OP_PUSHBYTES_12 030000000000000000000001 OP_TOALTSTACK OP_PUSHBYTES_32 000000000000002bad59995e2ae358ca0ee9e4ae433d4db67d88b8e86aa7d62d OP_TOALTSTACK OP_PUSHBYTES_32 39d771ed264db1b6ecc1268207b416c020c40c1f08c92c4da0780584f93efadc OP_DUP OP_HASH256 OP_PUSHBYTES_32 1b6af0a8751b294c347ebb786029ba9fa3179dcfc1b6e0e042eec99e1b06681d OP_EQUALVERIFY OP_CHECKSIGVERIFY OP_PUSHBYTES_32 3807da2e6be6d76ee01c6db93f2e22d8ba66c17547917f1a06b5351dde265385 OP_CHECKSIGVERIFY OP_HASH160 OP_PUSHBYTES_20 3ae86613ff51c5950a452bd7d07c99319801944e OP_EQUALVERIFY OP_DEPTH OP_PUSHNUM_1 OP_NUMEQUAL OP_IF OP_PUSHBYTES_2 6f70 OP_PUSHNUM_NEG1 OP_PUSHBYTES_62 1f8b08000000000002039bcfbd6325ff919f1b7c0f9ff93e73d5839d0b1c3639f59d60593a8b73da2cf7dcf5d363847e2dfb4f000000b193f10944000000 OP_ELSE OP_PUSHNUM_1 OP_ENDIF"}],"vout":[{"scriptpubkey":"6015003ae86613ff51c5950a452bd7d07c99319801944e","scriptpubkey_asm":"OP_PUSHNUM_16 OP_PUSHBYTES_21 003ae86613ff51c5950a452bd7d07c99319801944e","scriptpubkey_type":"unknown","scriptpubkey_address":"tb1sqqawsesnlagut9g2g54a05runycesqv5fczw9mrr","value":405}],"size":565,"weight":814,"sigops":0,"fee":1884,"status":{"confirmed":false}}`)),
		}, nil).Once()
		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
		blockInfo, err := api.GetTransactionBlockInfo(testnetTestTxHash)
		require.ErrorIs(t, err, mempool_space.BtcTxNotConfirmedError)
		assert.Empty(t, blockInfo)
		client.AssertExpectations(t)
	})
	t.Run("should return error if transaction not found", func(t *testing.T) {
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.Anything).Return(&http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(bytes.NewBufferString(`"Transaction not found"`)),
		}, nil).Once()
		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
		blockInfo, err := api.GetTransactionBlockInfo(testnetTestTxHash)
		require.ErrorContains(t, err, "Transaction not found")
		assert.Empty(t, blockInfo)
		client.AssertExpectations(t)
	})
	t.Run("should handle request error", func(t *testing.T) {
		errorSetUps := []func(client *mocks.HttpClientMock){
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(assert.AnError.Error())),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"9f0706c2717fc77bf0f225a4223933a7decb8d36902ddbb0accab8ea894f8b29","version":2,"locktime":2582723,"vin":[{"txid":"2e4428f122341d770a9074b3f62040f2d48b78a253ad0c46eabaf3732ac2f7eb","vout":0,"prevout":{"scriptpubkey":"76a91420d0cd81928e083da7e9ee690987df9f64bea9b088ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 20d0cd81928e083da7e9ee690987df9f64bea9b0 OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"miWU44anwP9mr3X4w6WTEAGuw6vkQwZQr1","value":489200},"scriptsig":"473044022062dae13ba281d0cf529b604bb59c1efcd7b83438af34d4a51acc6f31041be18c022044df281e688a52624f45f6c26662349d1f5efedd4d69530e65b7d7cec0d3792d0121038e509bc056004a5da7460b5acd5d4dcb2add41d53817180499e3814290ecc91e","scriptsig_asm":"OP_PUSHBYTES_71 3044022062dae13ba281d0cf529b604bb59c1efcd7b83438af34d4a51acc6f31041be18c022044df281e688a52624f45f6c26662349d1f5efedd4d69530e65b7d7cec0d3792d01 OP_PUSHBYTES_33 038e509bc056004a5da7460b5acd5d4dcb2add41d53817180499e3814290ecc91e","is_coinbase":false,"sequence":4294967293},{"txid":"f548a9fff50f98aae908594b5f0cd1a8ef57abf7a744a64b0f855b21389ff0b5","vout":0,"prevout":{"scriptpubkey":"76a91410c5469550e5b73b4c144e0b4bcd2e65fdeece3488ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 10c5469550e5b73b4c144e0b4bcd2e65fdeece34 OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"mh3dS5cYBfLGLvkGGc1LrFj7Qj7DAsccdT","value":500000},"scriptsig":"47304402206538fc72b896e4c6e807a4daf56191f68dec307c3011d082e69eeb3d45d6d8c302205a329814ab87901ae56a82587e716fa2282ecc665ab203da14d93db71181ecd8012102498a833095175800f40b2c0ab23f108b47a319a94ccea826062bf66c827e91a9","scriptsig_asm":"OP_PUSHBYTES_71 304402206538fc72b896e4c6e807a4daf56191f68dec307c3011d082e69eeb3d45d6d8c302205a329814ab87901ae56a82587e716fa2282ecc665ab203da14d93db71181ecd801 OP_PUSHBYTES_33 02498a833095175800f40b2c0ab23f108b47a319a94ccea826062bf66c827e91a9","is_coinbase":false,"sequence":4294967293}],"vout":[{"scriptpubkey":"76a91473cce22e78ec61cd54a6438ca1210b88561ebcdd88ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 73cce22e78ec61cd54a6438ca1210b88561ebcdd OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"mr5FSft4PQvoWbf9Sf5iQXbw1u445iNksp","value":488600},{"scriptpubkey":"76a9142c81478132b5dda64ffc484a0d225096c4b22ad588ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 2c81478132b5dda64ffc484a0d225096c4b22ad5 OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"mjaGtyj74LYn7gApr17prZxDPDnfuUnRa5","value":500000}],"size":372,"weight":1488,"sigops":8,"fee":600,"status":{"confirmed":true,"block_height":2582756,"block_hash":"00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f","block_time":1710931198}}`)),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(assert.AnError.Error())),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"9f0706c2717fc77bf0f225a4223933a7decb8d36902ddbb0accab8ea894f8b29","version":2,"locktime":2582723,"vin":[{"txid":"2e4428f122341d770a9074b3f62040f2d48b78a253ad0c46eabaf3732ac2f7eb","vout":0,"prevout":{"scriptpubkey":"76a91420d0cd81928e083da7e9ee690987df9f64bea9b088ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 20d0cd81928e083da7e9ee690987df9f64bea9b0 OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"miWU44anwP9mr3X4w6WTEAGuw6vkQwZQr1","value":489200},"scriptsig":"473044022062dae13ba281d0cf529b604bb59c1efcd7b83438af34d4a51acc6f31041be18c022044df281e688a52624f45f6c26662349d1f5efedd4d69530e65b7d7cec0d3792d0121038e509bc056004a5da7460b5acd5d4dcb2add41d53817180499e3814290ecc91e","scriptsig_asm":"OP_PUSHBYTES_71 3044022062dae13ba281d0cf529b604bb59c1efcd7b83438af34d4a51acc6f31041be18c022044df281e688a52624f45f6c26662349d1f5efedd4d69530e65b7d7cec0d3792d01 OP_PUSHBYTES_33 038e509bc056004a5da7460b5acd5d4dcb2add41d53817180499e3814290ecc91e","is_coinbase":false,"sequence":4294967293},{"txid":"f548a9fff50f98aae908594b5f0cd1a8ef57abf7a744a64b0f855b21389ff0b5","vout":0,"prevout":{"scriptpubkey":"76a91410c5469550e5b73b4c144e0b4bcd2e65fdeece3488ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 10c5469550e5b73b4c144e0b4bcd2e65fdeece34 OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"mh3dS5cYBfLGLvkGGc1LrFj7Qj7DAsccdT","value":500000},"scriptsig":"47304402206538fc72b896e4c6e807a4daf56191f68dec307c3011d082e69eeb3d45d6d8c302205a329814ab87901ae56a82587e716fa2282ecc665ab203da14d93db71181ecd8012102498a833095175800f40b2c0ab23f108b47a319a94ccea826062bf66c827e91a9","scriptsig_asm":"OP_PUSHBYTES_71 304402206538fc72b896e4c6e807a4daf56191f68dec307c3011d082e69eeb3d45d6d8c302205a329814ab87901ae56a82587e716fa2282ecc665ab203da14d93db71181ecd801 OP_PUSHBYTES_33 02498a833095175800f40b2c0ab23f108b47a319a94ccea826062bf66c827e91a9","is_coinbase":false,"sequence":4294967293}],"vout":[{"scriptpubkey":"76a91473cce22e78ec61cd54a6438ca1210b88561ebcdd88ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 73cce22e78ec61cd54a6438ca1210b88561ebcdd OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"mr5FSft4PQvoWbf9Sf5iQXbw1u445iNksp","value":488600},{"scriptpubkey":"76a9142c81478132b5dda64ffc484a0d225096c4b22ad588ac","scriptpubkey_asm":"OP_DUP OP_HASH160 OP_PUSHBYTES_20 2c81478132b5dda64ffc484a0d225096c4b22ad5 OP_EQUALVERIFY OP_CHECKSIG","scriptpubkey_type":"p2pkh","scriptpubkey_address":"mjaGtyj74LYn7gApr17prZxDPDnfuUnRa5","value":500000}],"size":372,"weight":1488,"sigops":8,"fee":600,"status":{"confirmed":true,"block_height":2582756,"block_hash":"00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f","block_time":1710931198}}`)),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"id":"0000000000`)),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"9f0706c2717fc77bf0f225a422`)),
				}, nil).Once()
			},
		}
		for _, setup := range errorSetUps {
			client := &mocks.HttpClientMock{}
			setup(client)
			api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
			blockInfo, err := api.GetTransactionBlockInfo(testnetTestTxHash)
			require.Error(t, err)
			assert.Empty(t, blockInfo)
			client.AssertExpectations(t)
		}
	})
}

// nolint:funlen
func TestMempoolSpaceApi_GetCoinbaseInformation(t *testing.T) {
	t.Run("should return coinbase information for confirmed transaction", func(t *testing.T) {
		const blockHashString = "00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f"
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			return r != nil && r.URL.String() == testnetUrl+"/tx/"+testnetWitnessTestTxHash
		})).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"5cadcbc1ccd91b222346f22c9a9a6fdbf20c9338ec8df0b36097e92d029509ec","version":2,"locktime":0,"vin":[{"txid":"1ac616516a8d94666293d5b36342c06be60d8df437e32278a22e02076a1b3eba","vout":7,"prevout":{"scriptpubkey":"5120588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1ptz8dytlnvmpx98q6mnj5x9dlc7pf06kykzqk4qgfya6mptsg3plsjkyqka","value":14246},"scriptsig":"","scriptsig_asm":"","witness":["881b79ef50067a758bd097077acd5df935cef7fb7bc670c24d537a2d4d72c0c0c65762e6dd4c240bf589f7e8928d413c3a3020f7f4d6f0080c54e012336b09f3","20d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555ac0063036f726401010a746578742f706c61696e000d656c7261756c69746f2e62746368","c0d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555"],"is_coinbase":false,"sequence":4294967295,"inner_witnessscript_asm":"OP_PUSHBYTES_32 d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555 OP_CHECKSIG OP_0 OP_IF OP_PUSHBYTES_3 6f7264 OP_PUSHBYTES_1 01 OP_PUSHBYTES_10 746578742f706c61696e OP_0 OP_PUSHBYTES_13 656c7261756c69746f2e627463 OP_ENDIF"}],"vout":[{"scriptpubkey":"5120f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1p77g2j6wy827qanf5enq6yfzv0lhdvgahqkpexzpykkql4e5gtmns5cymxz","value":546}],"size":266,"weight":548,"sigops":0,"fee":13700,"status":{"confirmed":true,"block_height":2582756,"block_hash":"00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f","block_time":1710931198}}`)),
		}, nil).Once()
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			return r != nil && r.URL.String() == testnetUrl+"/block/"+blockHashString
		})).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"id":"00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f","height":2582756,"version":536870912,"timestamp":1710931198,"tx_count":1030,"size":357157,"weight":848203,"merkle_root":"705c1425cc8ff977778acc9729fbd2dc4a8c21358f3e0b4fefa25ce5c343a3df","previousblockhash":"0000000000000014f9d3bd63b0257095d7d272134674285665f00d572fd2b0d8","mediantime":1710925891,"nonce":488178752,"bits":486604799,"difficulty":1.0}`)),
		}, nil).Once()
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			return r != nil && r.URL.String() == testnetUrl+"/block/"+blockHashString+"/raw"
		})).RunAndReturn(func(request *http.Request) (*http.Response, error) {
			block := test.MustReadFileString("./datasets/block-00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f-testnet.txt")
			blockBytes, err := hex.DecodeString(block)
			require.NoError(t, err)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(blockBytes)),
			}, nil
		}).Once()

		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
		coinbaseInfo, err := api.GetCoinbaseInformation(testnetWitnessTestTxHash)
		require.NoError(t, err)
		var blockHash, witnessMerkleRoot [32]byte
		blockHashBytes, err := hex.DecodeString("00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f")
		require.NoError(t, err)
		copy(blockHash[:], blockHashBytes)
		witnessMerkleRootBytes, err := hex.DecodeString("ce489a9f5cd7351ed57606c373329a208d37a8654bf473916eeec4b8ffee55f4")
		require.NoError(t, err)
		copy(witnessMerkleRoot[:], witnessMerkleRootBytes)
		pmt, err := hex.DecodeString("060400000c89cc284039442363175d7f6d34492885618442a865381e6221a6223f5fa18af5ed018e20e419923fa29585d93161cae0a1a9ebb5095ca055e9888a8c0b956f4604987862683909d0749a969cbe1de8669b72c1174c0d4e8cd50238cc7198743148caa41a01a85c260ab184f3846125bda6e5f21c14dc314b207b4428c29cf6807977059737aa020de8c76fa6c304c99db6306a4a0aa1532ae7dbdb6453aa09d13c835e6ea0226d377adb44b941cf742836b395f5f372c945d2ed13170dc1f453702b70faddefb6a289a3c7968ceb7d8efea9c205ca04e0c2fc90f6f5d9055282caf9fbe1bdd0bef6417b2534ff0a1d83a77fab36b5977ff6f2adc7c8dd8fe8414141d12ea8b40e7c3af647a1473204d83b5726bceb8a36a98bb53ceb4d4419b6027d5f22b753db43214c7f73cd0165fb794a09229883f0650144d85dc20d6f45b753628880e5c4fc39d63e908029221ef9930cc988e3df6f91b3e887db348b0d8a89fe02723b50e2930703a046e751abbc64c16aa48f6de16a18c431389e551603ff0f00")
		require.NoError(t, err)
		tx, err := hex.DecodeString("04000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1203e4682704eebcfa650843f719b701000000000000000377781d00000000001600143758d475313d557dbe8b1d90406c5b3b4dbbc00df79900000000000017a914a775ee7e3118ac67f181faca330f1d5c7658d205870000000000000000266a24aa21a9ed3dccc6b158d03b681f5cd8c71653097d0e6a51ac3e19de0add0a2a43419622a500000000")
		require.NoError(t, err)
		assert.Equal(t, blockchain.BtcCoinbaseTransactionInformation{
			BtcTxSerialized:      tx,
			BlockHash:            blockHash,
			BlockHeight:          big.NewInt(2582756),
			SerializedPmt:        pmt,
			WitnessMerkleRoot:    witnessMerkleRoot,
			WitnessReservedValue: [32]byte{},
		}, coinbaseInfo)
		client.AssertExpectations(t)
	})
	t.Run("should return error for unconfirmed transaction", func(t *testing.T) {
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"e36742a7ead1a73a2b0961008fc7d26d97f801dfdd60bd66e8b2a53f63f3006b","version":2,"locktime":0,"vin":[{"txid":"ce59b1b999a95aacb504b9dd87e1ba405c30ec850f4532f5918c3f4728afca0c","vout":0,"prevout":{"scriptpubkey":"5120f6c9ea1ebf1046961ca082b0ecc1376ae5a155e09bb967c674c1ed33b4699b26","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 f6c9ea1ebf1046961ca082b0ecc1376ae5a155e09bb967c674c1ed33b4699b26","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1p7my7584lzprfv89qs2cwesfhdtj6z40qnwuk03n5c8kn8drfnvnqdz0fdl","value":10453},"scriptsig":"","scriptsig_asm":"","witness":["9e14fc4c4cfca73a89e25e1216ae3a22302a12a7c6e1e3a568e05e8cb824112b","b23ad59c094162e6ac022b18a3117acc543b8dd079a798b24afb348b4b3a5f792185f331b74d7a948262e14736c9f4ce2b50383a2ac772be7d014a2f96a875c0","5e45190fe5a7347cbea89f000d01ad6a79c95af34d1c9f078060a39c7b0ab0d4cdaa015128ce66589b35ab06640a3b0c021a307a4f4972950f463c16238f13c8","0c0200000000000000000000006b20000000000000004c30ab82074704408e2b69809854e3e34b4714fff85807216e6b20e1043abb4213cf0e6ff080d41d82e5a08264d0321743fdf8bf0bd38773b145df76aa200d78532557a79409654bdb2cd16328fb9100dbcdfdbdb5d73ec702682e36e05088ad20c872d5c3a25337722bcdc14b0aeeeabfa373e44f9426356a556ad32923cc2307ada9141813d6cfa01b528115c019382c1feaf427060a3a8874519c63026f704f4c971f8b08000000000002032bb462d461c001a40fe6a6fb9d60c00d4a2d5fffb7e47ecac054c51db1afd34a95cde2b4ced4fd33edf6316c63af5c12595d23f7e1796ad009e573b7a7f8559cba255b10f563d26ed684252b8fececcc4d79af7cff723bffdf00962f9f5ebcecda6513aa79cde6e1e9ba63f219cf9b5f4e335f96535acf7fd15b83b5d4c5076c9b6b9e1900d753e433ae000000675168","c0e1043abb4213cf0e6ff080d41d82e5a08264d0321743fdf8bf0bd38773b145df84654164130538fbb3f7b1539f3c5461c509d9f7a233c6fb59104ef8f0f8d7e5"],"is_coinbase":false,"sequence":4294967293,"inner_witnessscript_asm":"OP_PUSHBYTES_12 020000000000000000000000 OP_TOALTSTACK OP_PUSHBYTES_32 000000000000004c30ab82074704408e2b69809854e3e34b4714fff85807216e OP_TOALTSTACK OP_PUSHBYTES_32 e1043abb4213cf0e6ff080d41d82e5a08264d0321743fdf8bf0bd38773b145df OP_DUP OP_HASH256 OP_PUSHBYTES_32 0d78532557a79409654bdb2cd16328fb9100dbcdfdbdb5d73ec702682e36e050 OP_EQUALVERIFY OP_CHECKSIGVERIFY OP_PUSHBYTES_32 c872d5c3a25337722bcdc14b0aeeeabfa373e44f9426356a556ad32923cc2307 OP_CHECKSIGVERIFY OP_HASH160 OP_PUSHBYTES_20 1813d6cfa01b528115c019382c1feaf427060a3a OP_EQUALVERIFY OP_DEPTH OP_PUSHNUM_1 OP_NUMEQUAL OP_IF OP_PUSHBYTES_2 6f70 OP_PUSHNUM_NEG1 OP_PUSHDATA1 1f8b08000000000002032bb462d461c001a40fe6a6fb9d60c00d4a2d5fffb7e47ecac054c51db1afd34a95cde2b4ced4fd33edf6316c63af5c12595d23f7e1796ad009e573b7a7f8559cba255b10f563d26ed684252b8fececcc4d79af7cff723bffdf00962f9f5ebcecda6513aa79cde6e1e9ba63f219cf9b5f4e335f96535acf7fd15b83b5d4c5076c9b6b9e1900d753e433ae000000 OP_ELSE OP_PUSHNUM_1 OP_ENDIF"}],"vout":[{"scriptpubkey":"6015001813d6cfa01b528115c019382c1feaf427060a3a","scriptpubkey_asm":"OP_PUSHNUM_16 OP_PUSHBYTES_21 001813d6cfa01b528115c019382c1feaf427060a3a","scriptpubkey_type":"unknown","scriptpubkey_address":"tb1sqqvp84k05qd49qg4cqvnstqlat6zwps28gfcfhdy","value":297},{"scriptpubkey":"a914568f9c96289308a4c5f7f78b1880bed2d6ed10b587","scriptpubkey_asm":"OP_HASH160 OP_PUSHBYTES_20 568f9c96289308a4c5f7f78b1880bed2d6ed10b5 OP_EQUAL","scriptpubkey_type":"p2sh","scriptpubkey_address":"2N18v96TPxUUar2kYPXKBakoSGZ4y5DNV5V","value":8023}],"size":689,"weight":1034,"sigops":0,"fee":2133,"status":{"confirmed":false}}`)),
		}, nil).Once()

		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
		coinbaseInfo, err := api.GetCoinbaseInformation(testnetWitnessTestTxHash)
		require.ErrorIs(t, err, mempool_space.BtcTxNotConfirmedError)
		assert.Empty(t, coinbaseInfo)
		client.AssertExpectations(t)
	})
	t.Run("should return error if transaction not found", func(t *testing.T) {
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.Anything).Return(&http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(bytes.NewBufferString(`Transaction not found`)),
		}, nil).Once()

		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
		coinbaseInfo, err := api.GetCoinbaseInformation(testnetWitnessTestTxHash)
		require.ErrorContains(t, err, "Transaction not found")
		assert.Empty(t, coinbaseInfo)
		client.AssertExpectations(t)
	})
	t.Run("should handle request error", func(t *testing.T) {
		errorSetUps := []func(*mocks.HttpClientMock){
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"5cadcbc1ccd91b222346f22c9a9a6fdbf20c9338ec8df0b36097e92d029509ec","version":2,"locktime":0,"vin":[{"txid":"1ac616516a8d94666293d5b36342c06be60d8df437e32278a22e02076a1b3eba","vout":7,"prevout":{"scriptpubkey":"5120588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1ptz8dytlnvmpx98q6mnj5x9dlc7pf06kykzqk4qgfya6mptsg3plsjkyqka","value":14246},"scriptsig":"","scriptsig_asm":"","witness":["881b79ef50067a758bd097077acd5df935cef7fb7bc670c24d537a2d4d72c0c0c65762e6dd4c240bf589f7e8928d413c3a3020f7f4d6f0080c54e012336b09f3","20d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555ac0063036f726401010a746578742f706c61696e000d656c7261756c69746f2e62746368","c0d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555"],"is_coinbase":false,"sequence":4294967295,"inner_witnessscript_asm":"OP_PUSHBYTES_32 d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555 OP_CHECKSIG OP_0 OP_IF OP_PUSHBYTES_3 6f7264 OP_PUSHBYTES_1 01 OP_PUSHBYTES_10 746578742f706c61696e OP_0 OP_PUSHBYTES_13 656c7261756c69746f2e627463 OP_ENDIF"}],"vout":[{"scriptpubkey":"5120f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1p77g2j6wy827qanf5enq6yfzv0lhdvgahqkpexzpykkql4e5gtmns5cymxz","value":546}],"size":266,"weight":548,"sigops":0,"fee":13700,"status":{"confirmed":true,"block_height":2582756,"block_hash":"00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f","block_time":1710931198}}`)),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"id":"00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f","height":2582756,"version":536870912,"timestamp":1710931198,"tx_count":1030,"size":357157,"weight":848203,"merkle_root":"705c1425cc8ff977778acc9729fbd2dc4a8c21358f3e0b4fefa25ce5c343a3df","previousblockhash":"0000000000000014f9d3bd63b0257095d7d272134674285665f00d572fd2b0d8","mediantime":1710925891,"nonce":488178752,"bits":486604799,"difficulty":1.0}`)),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("not a hex")),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"5cadcbc1ccd91b222346f22c9a9a6fdbf20c9338ec8df0b36097e92d029509ec","version":2,"locktime":0,"vin":[{"txid":"1ac616516a8d94666293d5b36342c06be60d8df437e32278a22e02076a1b3eba","vout":7,"prevout":{"scriptpubkey":"5120588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1ptz8dytlnvmpx98q6mnj5x9dlc7pf06kykzqk4qgfya6mptsg3plsjkyqka","value":14246},"scriptsig":"","scriptsig_asm":"","witness":["881b79ef50067a758bd097077acd5df935cef7fb7bc670c24d537a2d4d72c0c0c65762e6dd4c240bf589f7e8928d413c3a3020f7f4d6f0080c54e012336b09f3","20d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555ac0063036f726401010a746578742f706c61696e000d656c7261756c69746f2e62746368","c0d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555"],"is_coinbase":false,"sequence":4294967295,"inner_witnessscript_asm":"OP_PUSHBYTES_32 d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555 OP_CHECKSIG OP_0 OP_IF OP_PUSHBYTES_3 6f7264 OP_PUSHBYTES_1 01 OP_PUSHBYTES_10 746578742f706c61696e OP_0 OP_PUSHBYTES_13 656c7261756c69746f2e627463 OP_ENDIF"}],"vout":[{"scriptpubkey":"5120f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1p77g2j6wy827qanf5enq6yfzv0lhdvgahqkpexzpykkql4e5gtmns5cymxz","value":546}],"size":266,"weight":548,"sigops":0,"fee":13700,"status":{"confirmed":true,"block_height":2582756,"block_hash":"00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f","block_time":1710931198}}`)),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"id":"00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f","height":2582756,"version":536870912,"timestamp":1710931198,"tx_count":1030,"size":357157,"weight":848203,"merkle_root":"705c1425cc8ff977778acc9729fbd2dc4a8c21358f3e0b4fefa25ce5c343a3df","previousblockhash":"0000000000000014f9d3bd63b0257095d7d272134674285665f00d572fd2b0d8","mediantime":1710925891,"nonce":488178752,"bits":486604799,"difficulty":1.0}`)),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(assert.AnError.Error())),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"5cadcbc1ccd91b222346f22c9a9a6fdbf20c9338ec8df0b36097e92d029509ec","version":2,"locktime":0,"vin":[{"txid":"1ac616516a8d94666293d5b36342c06be60d8df437e32278a22e02076a1b3eba","vout":7,"prevout":{"scriptpubkey":"5120588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1ptz8dytlnvmpx98q6mnj5x9dlc7pf06kykzqk4qgfya6mptsg3plsjkyqka","value":14246},"scriptsig":"","scriptsig_asm":"","witness":["881b79ef50067a758bd097077acd5df935cef7fb7bc670c24d537a2d4d72c0c0c65762e6dd4c240bf589f7e8928d413c3a3020f7f4d6f0080c54e012336b09f3","20d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555ac0063036f726401010a746578742f706c61696e000d656c7261756c69746f2e62746368","c0d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555"],"is_coinbase":false,"sequence":4294967295,"inner_witnessscript_asm":"OP_PUSHBYTES_32 d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555 OP_CHECKSIG OP_0 OP_IF OP_PUSHBYTES_3 6f7264 OP_PUSHBYTES_1 01 OP_PUSHBYTES_10 746578742f706c61696e OP_0 OP_PUSHBYTES_13 656c7261756c69746f2e627463 OP_ENDIF"}],"vout":[{"scriptpubkey":"5120f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1p77g2j6wy827qanf5enq6yfzv0lhdvgahqkpexzpykkql4e5gtmns5cymxz","value":546}],"size":266,"weight":548,"sigops":0,"fee":13700,"status":{"confirmed":true,"block_height":2582756,"block_hash":"00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f","block_time":1710931198}}`)),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"id":"00000000`)),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"5cadcbc1ccd91b222346f22c9a9a6fdbf20c9338ec8df0b36097e92d029509ec","version":2,"locktime":0,"vin":[{"txid":"1ac616516a8d94666293d5b36342c06be60d8df437e32278a22e02076a1b3eba","vout":7,"prevout":{"scriptpubkey":"5120588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 588ed22ff366c2629c1adce54315bfc78297eac4b0816a81092775b0ae08887f","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1ptz8dytlnvmpx98q6mnj5x9dlc7pf06kykzqk4qgfya6mptsg3plsjkyqka","value":14246},"scriptsig":"","scriptsig_asm":"","witness":["881b79ef50067a758bd097077acd5df935cef7fb7bc670c24d537a2d4d72c0c0c65762e6dd4c240bf589f7e8928d413c3a3020f7f4d6f0080c54e012336b09f3","20d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555ac0063036f726401010a746578742f706c61696e000d656c7261756c69746f2e62746368","c0d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555"],"is_coinbase":false,"sequence":4294967295,"inner_witnessscript_asm":"OP_PUSHBYTES_32 d90ef3f63c98d07f0458251f147d987b7b918ba54c92807513064271f028d555 OP_CHECKSIG OP_0 OP_IF OP_PUSHBYTES_3 6f7264 OP_PUSHBYTES_1 01 OP_PUSHBYTES_10 746578742f706c61696e OP_0 OP_PUSHBYTES_13 656c7261756c69746f2e627463 OP_ENDIF"}],"vout":[{"scriptpubkey":"5120f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_asm":"OP_PUSHNUM_1 OP_PUSHBYTES_32 f790a969c43abc0ecd34ccc1a2244c7feed623b70583930824b581fae6885ee7","scriptpubkey_type":"v1_p2tr","scriptpubkey_address":"tb1p77g2j6wy827qanf5enq6yfzv0lhdvgahqkpexzpykkql4e5gtmns5cymxz","value":546}],"size":266,"weight":548,"sigops":0,"fee":13700,"status":{"confirmed":true,"block_height":2582756,"block_hash":"00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f","block_time":1710931198}}`)),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(assert.AnError.Error())),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"txid":"5cadcbc1ccd`)),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(assert.AnError.Error())),
				}, nil).Once()
			},
		}
		for _, setup := range errorSetUps {
			client := &mocks.HttpClientMock{}
			setup(client)
			api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
			coinbaseInfo, err := api.GetCoinbaseInformation(testnetWitnessTestTxHash)
			require.Error(t, err)
			assert.Empty(t, coinbaseInfo)
		}
	})
}

func TestMempoolSpaceApi_NetworkName(t *testing.T) {
	type testValue struct {
		ChainParams *chaincfg.Params
		Url         string
	}
	const url = "http://localhost:8332"
	table := test.Table[testValue, string]{
		{Value: testValue{
			ChainParams: &chaincfg.MainNetParams,
			Url:         mainnetUrl,
		}, Result: "mainnet"},
		{Value: testValue{
			ChainParams: &chaincfg.TestNet3Params,
			Url:         testnetUrl,
		}, Result: "testnet3"},
		{Value: testValue{
			ChainParams: &chaincfg.RegressionNetParams,
			Url:         url,
		}, Result: "regtest"},
		{Value: testValue{
			ChainParams: &chaincfg.SigNetParams,
			Url:         url,
		}, Result: "signet"},
	}
	client := http.DefaultClient
	test.RunTable(t, table, func(v testValue) string {
		api := mempool_space.NewMempoolSpaceApi(client, v.ChainParams, v.Url)
		return api.NetworkName()
	})
}

// nolint:funlen
func TestMempoolSpaceApi_GetBlockchainInfo(t *testing.T) {
	t.Run("should return mainnet info", func(t *testing.T) {
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			return r != nil && r.URL.String() == mainnetUrl+"/blocks/tip/height"
		})).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("904453")),
		}, nil).Once()
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			return r != nil && r.URL.String() == mainnetUrl+"/blocks/tip/hash"
		})).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("0000000000000000000062a61c6753145aed01ee47420de08b1eb288c89ff6fa")),
		}, nil).Once()

		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.MainNetParams, mainnetUrl)
		info, err := api.GetBlockchainInfo()
		require.NoError(t, err)
		assert.Equal(t, blockchain.BitcoinBlockchainInfo{
			NetworkName:      "mainnet",
			ValidatedBlocks:  big.NewInt(904453),
			ValidatedHeaders: big.NewInt(904453),
			BestBlockHash:    "0000000000000000000062a61c6753145aed01ee47420de08b1eb288c89ff6fa",
		}, info)
		client.AssertExpectations(t)
	})
	t.Run("should return testnet info", func(t *testing.T) {
		client := &mocks.HttpClientMock{}
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			return r != nil && r.URL.String() == testnetUrl+"/blocks/tip/height"
		})).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("4550195")),
		}, nil).Once()
		client.EXPECT().Do(mock.MatchedBy(func(r *http.Request) bool {
			return r != nil && r.URL.String() == testnetUrl+"/blocks/tip/hash"
		})).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("0000000000000019c5b3123cbf53c67dd0458fb64e5be2c8421678e721c86a01")),
		}, nil).Once()

		api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)
		info, err := api.GetBlockchainInfo()
		require.NoError(t, err)
		assert.Equal(t, blockchain.BitcoinBlockchainInfo{
			NetworkName:      "testnet3",
			ValidatedBlocks:  big.NewInt(4550195),
			ValidatedHeaders: big.NewInt(4550195),
			BestBlockHash:    "0000000000000019c5b3123cbf53c67dd0458fb64e5be2c8421678e721c86a01",
		}, info)
		client.AssertExpectations(t)
	})
	t.Run("should return error if request fails", func(t *testing.T) {
		errorSetUps := []func(*mocks.HttpClientMock){
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("not a number")),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("123")),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("not a hex")),
				}, nil).Once()
			},
			func(client *mocks.HttpClientMock) {
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("123")),
				}, nil).Once()
				client.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(assert.AnError.Error())),
				}, nil).Once()
			},
		}
		for _, setup := range errorSetUps {
			client := &mocks.HttpClientMock{}
			setup(client)
			api := mempool_space.NewMempoolSpaceApi(client, &chaincfg.MainNetParams, mainnetUrl)
			info, err := api.GetBlockchainInfo()
			require.Error(t, err)
			assert.Empty(t, info)
		}
	})
}
