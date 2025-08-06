package integration_test

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin/mempool_space"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/datasets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"net/http"
	"testing"
	"time"
)

const (
	mainnetUrl = "https://mempool.space/api"
	testnetUrl = "https://mempool.space/testnet/api"
	regtestUrl = "http://localhost:1234/api"
)

const (
	testFilePath = "../mocks/"
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

// nolint:funlen
func TestMempoolSpaceApi_ValidateAddress(t *testing.T) {
	httpClient := http.DefaultClient
	mainnet := mempool_space.NewMempoolSpaceApi(httpClient, &chaincfg.MainNetParams, mainnetUrl)
	testnet := mempool_space.NewMempoolSpaceApi(httpClient, &chaincfg.TestNet3Params, testnetUrl)
	regtest := mempool_space.NewMempoolSpaceApi(httpClient, &chaincfg.RegressionNetParams, regtestUrl)
	invalid := mempool_space.NewMempoolSpaceApi(httpClient, &chaincfg.SigNetParams, "http://invalid-url")

	p2pkhMainnet := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
	p2pkhTestnet := "mipcBbFg9gMiCh81Kj8tqqdgoZub1ZJRfn"
	p2shMainnet := "3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy"
	p2shTestnet := "2MzQwSSnBHWHqSAqtTVQ6v47XtaisrJa1Vc"
	p2wpkhMainnet := "bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq"
	p2wpkhTestnet := "tb1qg9stx3w3xj0t5j8d5q2g4yvz5gj8v3tjcraugq"
	p2wshMainnet := "bc1quhruqrghgcca950rvhtrg7cpd7u8k6svpzgzmrjy8xyukacl5lkq0r8l2d"
	p2wshTestnet := "tb1qgpgtqj68zwsdz7xmvqxxxaan7dcfgu76jz0cfzynqgrtvdsxlyqsf7dfz8"
	p2trMainnet := "bc1pq2g3k9fleqcvu382g674psux05wwa08w6gw6022mr7sqla8009ws3p5054"
	p2trTestnet := "tb1p9jveg4j5mh2z3v6e6z93ln5jn4zfehd873ps2vv0g6k234tqw67sm08vk5"

	p2wpkhRegtest := "bcrt1q8gf8taa048wxjj90g7htpt5t5hzja4a2qz6use"
	p2wshRegtest := "bcrt1qtmm4qallkmnd2vl5y3w3an3uvq6w5v2ahqvfqm0mfxny8cnsdrashv8fsr"
	p2trRegtest := "bcrt1ptzfcz2r2uskkhq2yq3umhahf6y6algfrlwhu3eu8mjht44gu984q6ucjxd"

	require.NoError(t, mainnet.ValidateAddress(p2pkhMainnet))
	require.NoError(t, mainnet.ValidateAddress(p2shMainnet))
	require.NoError(t, mainnet.ValidateAddress(p2wpkhMainnet))
	require.NoError(t, mainnet.ValidateAddress(p2wshMainnet))
	require.NoError(t, mainnet.ValidateAddress(p2trMainnet))
	require.ErrorIs(t, mainnet.ValidateAddress(p2wpkhRegtest), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, mainnet.ValidateAddress(p2wshRegtest), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, mainnet.ValidateAddress(p2trRegtest), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, mainnet.ValidateAddress(p2pkhTestnet), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, mainnet.ValidateAddress(p2shTestnet), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, mainnet.ValidateAddress(p2wpkhTestnet), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, mainnet.ValidateAddress(p2wshTestnet), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, mainnet.ValidateAddress(p2trTestnet), blockchain.BtcAddressInvalidNetworkError)

	require.ErrorIs(t, testnet.ValidateAddress(p2pkhMainnet), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, testnet.ValidateAddress(p2shMainnet), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, testnet.ValidateAddress(p2wpkhMainnet), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, testnet.ValidateAddress(p2wshMainnet), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, testnet.ValidateAddress(p2trMainnet), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, testnet.ValidateAddress(p2wpkhRegtest), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, testnet.ValidateAddress(p2wshRegtest), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, testnet.ValidateAddress(p2trRegtest), blockchain.BtcAddressInvalidNetworkError)
	require.NoError(t, testnet.ValidateAddress(p2pkhTestnet))
	require.NoError(t, testnet.ValidateAddress(p2shTestnet))
	require.NoError(t, testnet.ValidateAddress(p2wpkhTestnet))
	require.NoError(t, testnet.ValidateAddress(p2wshTestnet))
	require.NoError(t, testnet.ValidateAddress(p2trTestnet))

	const unsupportedNetwork = "unsupported network"

	require.ErrorContains(t, regtest.ValidateAddress(p2pkhMainnet), unsupportedNetwork)
	require.ErrorContains(t, regtest.ValidateAddress(p2shMainnet), unsupportedNetwork)
	require.ErrorContains(t, regtest.ValidateAddress(p2wpkhMainnet), unsupportedNetwork)
	require.ErrorContains(t, regtest.ValidateAddress(p2wshMainnet), unsupportedNetwork)
	require.ErrorContains(t, regtest.ValidateAddress(p2trMainnet), unsupportedNetwork)
	require.ErrorContains(t, regtest.ValidateAddress(p2wpkhTestnet), unsupportedNetwork)
	require.ErrorContains(t, regtest.ValidateAddress(p2wshTestnet), unsupportedNetwork)
	require.ErrorContains(t, regtest.ValidateAddress(p2trTestnet), unsupportedNetwork)
	require.ErrorContains(t, regtest.ValidateAddress(p2pkhTestnet), unsupportedNetwork)
	require.ErrorContains(t, regtest.ValidateAddress(p2shTestnet), unsupportedNetwork)
	require.ErrorContains(t, regtest.ValidateAddress(p2wpkhRegtest), unsupportedNetwork)
	require.ErrorContains(t, regtest.ValidateAddress(p2wshRegtest), unsupportedNetwork)
	require.ErrorContains(t, regtest.ValidateAddress(p2trRegtest), unsupportedNetwork)
	require.ErrorContains(t, invalid.ValidateAddress(p2pkhMainnet), unsupportedNetwork)
	require.ErrorContains(t, invalid.ValidateAddress(p2shMainnet), unsupportedNetwork)
	require.ErrorContains(t, invalid.ValidateAddress(p2pkhTestnet), unsupportedNetwork)
	require.ErrorContains(t, invalid.ValidateAddress(p2shTestnet), unsupportedNetwork)
	require.ErrorContains(t, invalid.ValidateAddress(p2wpkhMainnet), unsupportedNetwork)
	require.ErrorContains(t, invalid.ValidateAddress(p2wpkhTestnet), unsupportedNetwork)

	const notSupported = "not-an-address"
	require.Error(t, testnet.ValidateAddress(notSupported))
	require.Error(t, mainnet.ValidateAddress(notSupported))
}

func TestMempoolSpaceApi_GetTransactionInfo(t *testing.T) {
	api := mempool_space.NewMempoolSpaceApi(http.DefaultClient, &chaincfg.TestNet3Params, testnetUrl)

	result, err := api.GetTransactionInfo(testnetTestTxHash)
	require.NoError(t, err)
	assert.Equal(t, testnetTestTxHash, result.Hash)
	assert.GreaterOrEqual(t, result.Confirmations, uint64(1965878))
	assert.False(t, result.HasWitness)
	assert.Equal(t, map[string][]*entities.Wei{
		"mr5FSft4PQvoWbf9Sf5iQXbw1u445iNksp": {entities.NewWei(0.004886 * 1e18)}, "mjaGtyj74LYn7gApr17prZxDPDnfuUnRa5": {entities.NewWei(0.005 * 1e18)},
	}, result.Outputs)

	result, err = api.GetTransactionInfo(flyoverTestnetPegoutTxHash)
	require.NoError(t, err)
	assert.Equal(t, flyoverTestnetPegoutTxHash, result.Hash)
	assert.GreaterOrEqual(t, result.Confirmations, uint64(1966871))
	assert.False(t, result.HasWitness)
	assert.Equal(t, map[string][]*entities.Wei{
		"mqbKtarYKnoEdPheFFDGRjksvEpb2vJGNh": {entities.NewWei(0.005 * 1e18)},
		"mowfvQDraTDvRgZowL4tx5EatL1u78w65v": {entities.NewWei(0.01956600 * 1e18)},
		"":                                   {entities.NewWei(0)}, // Null data script output
	}, result.Outputs)

	result, err = api.GetTransactionInfo(testnetWitnessTestTxHash)
	require.NoError(t, err)
	assert.Equal(t, testnetWitnessTestTxHash, result.Hash)
	assert.True(t, result.HasWitness)
	assert.GreaterOrEqual(t, result.Confirmations, uint64(1965878), result.Confirmations)
	assert.Equal(t, map[string][]*entities.Wei{
		"tb1p77g2j6wy827qanf5enq6yfzv0lhdvgahqkpexzpykkql4e5gtmns5cymxz": {entities.NewWei(0.00000546 * 1e18)},
	}, result.Outputs)
}

func TestMempoolSpaceApi_GetTransactionInfo_NotFound(t *testing.T) {
	api := mempool_space.NewMempoolSpaceApi(http.DefaultClient, &chaincfg.TestNet3Params, testnetUrl)
	tx, err := api.GetTransactionInfo(txNotPresentInTestnet)
	require.ErrorContains(t, err, "Transaction not found")
	assert.Empty(t, tx)
}

func TestMempoolSpaceApi_GetRawTransaction(t *testing.T) {
	txBytes, err := hex.DecodeString("0200000002ebf7c22a73f3baea460cad53a2788bd4f24020f6b374900a771d3422f128442e000000006a473044022062dae13ba281d0cf529b604bb59c1efcd7b83438af34d4a51acc6f31041be18c022044df281e688a52624f45f6c26662349d1f5efedd4d69530e65b7d7cec0d3792d0121038e509bc056004a5da7460b5acd5d4dcb2add41d53817180499e3814290ecc91efdffffffb5f09f38215b850f4ba644a7f7ab57efa8d10c5f4b5908e9aa980ff5ffa948f5000000006a47304402206538fc72b896e4c6e807a4daf56191f68dec307c3011d082e69eeb3d45d6d8c302205a329814ab87901ae56a82587e716fa2282ecc665ab203da14d93db71181ecd8012102498a833095175800f40b2c0ab23f108b47a319a94ccea826062bf66c827e91a9fdffffff0298740700000000001976a91473cce22e78ec61cd54a6438ca1210b88561ebcdd88ac20a10700000000001976a9142c81478132b5dda64ffc484a0d225096c4b22ad588acc3682700")
	require.NoError(t, err)
	api := mempool_space.NewMempoolSpaceApi(http.DefaultClient, &chaincfg.TestNet3Params, testnetUrl)
	result, err := api.GetRawTransaction(testnetTestTxHash)
	require.NoError(t, err)
	require.Equal(t, txBytes, result)
}

func TestMempoolSpaceApi_GetRawTransaction_FromBlock(t *testing.T) {
	cases := []datasets.RawTransaction{
		datasets.BtcCoinbaseTxNoWitness,
		datasets.BtcSegwitTxNoWitness,
		datasets.BtcLegacyTxNoWitness,
	}
	mainnetBlock := test.GetBitcoinTestBlock(t, testFilePath+mainnetBlockFile)

	for _, tx := range cases {
		parsedTx, err := mainnetBlock.Tx(tx.Index)
		require.NoError(t, err)
		api := mempool_space.NewMempoolSpaceApi(http.DefaultClient, &chaincfg.MainNetParams, mainnetUrl)
		result, err := api.GetRawTransaction(parsedTx.Hash().String())
		require.NoError(t, err)
		expectedBytes, err := hex.DecodeString(tx.Tx)
		require.NoError(t, err)
		require.Equal(t, expectedBytes, result)
	}
}

func TestMempoolSpaceApi_GetRawTransaction_NotFound(t *testing.T) {
	api := mempool_space.NewMempoolSpaceApi(http.DefaultClient, &chaincfg.TestNet3Params, testnetUrl)
	result, err := api.GetRawTransaction(txNotPresentInTestnet)
	require.ErrorContains(t, err, "Transaction not found")
	assert.Empty(t, result)
}

func TestMempoolSpaceApi_GetPartialMerkleTree(t *testing.T) {
	api := mempool_space.NewMempoolSpaceApi(http.DefaultClient, &chaincfg.TestNet3Params, testnetUrl)
	pmt, err := api.GetPartialMerkleTree(testnetTestTxHash)
	require.NoError(t, err)
	assert.Equal(t,
		[]byte{
			6, 4, 0, 0, 12, 29, 10, 48, 252, 165, 16, 116, 206, 186, 165, 191, 148, 88, 81, 235, 33, 6, 145, 112, 222,
			232, 231, 85, 28, 201, 91, 116, 169, 149, 80, 15, 95, 117, 217, 211, 63, 209, 219, 8, 23, 72, 47, 55, 198, 133,
			219, 135, 80, 89, 219, 165, 211, 247, 218, 41, 20, 73, 23, 209, 185, 188, 146, 133, 56, 177, 117, 67, 181, 113,
			150, 50, 210, 208, 19, 53, 20, 138, 242, 205, 142, 143, 247, 193, 151, 69, 86, 174, 87, 219, 201, 162, 199, 241,
			222, 116, 234, 120, 126, 84, 162, 64, 26, 114, 34, 125, 107, 21, 29, 235, 83, 183, 251, 182, 162, 236, 48, 177,
			234, 191, 154, 106, 245, 2, 248, 60, 140, 162, 59, 179, 189, 194, 141, 213, 255, 62, 144, 202, 243, 43, 60, 89,
			136, 233, 84, 49, 46, 13, 0, 168, 108, 85, 142, 115, 186, 30, 77, 32, 231, 170, 71, 41, 139, 79, 137, 234, 184,
			202, 172, 176, 219, 45, 144, 54, 141, 203, 222, 167, 51, 57, 34, 164, 37, 242, 240, 123, 199, 127, 113, 194, 6,
			7, 159, 8, 4, 222, 132, 219, 160, 32, 198, 48, 141, 82, 70, 31, 189, 207, 220, 35, 171, 87, 89, 222, 249, 62, 78,
			207, 254, 254, 224, 191, 207, 80, 155, 97, 41, 160, 63, 172, 96, 150, 64, 254, 21, 252, 74, 72, 104, 203, 221, 251,
			141, 161, 154, 163, 202, 141, 154, 158, 71, 120, 74, 22, 146, 7, 93, 192, 88, 92, 151, 235, 109, 80, 237, 85, 203,
			209, 135, 83, 0, 177, 136, 200, 102, 169, 57, 87, 114, 135, 127, 192, 182, 182, 221, 250, 8, 182, 141, 67, 116, 241,
			227, 210, 28, 199, 196, 232, 156, 246, 32, 136, 250, 35, 99, 201, 151, 13, 32, 104, 11, 75, 5, 151, 32, 13, 81, 196,
			34, 44, 23, 183, 83, 98, 136, 128, 229, 196, 252, 57, 214, 62, 144, 128, 41, 34, 30, 249, 147, 12, 201, 136, 227, 223,
			111, 145, 179, 232, 135, 219, 52, 139, 13, 138, 137, 254, 2, 114, 59, 80, 226, 147, 7, 3, 160, 70, 231, 81, 171, 188,
			100, 193, 106, 164, 143, 109, 225, 106, 24, 196, 49, 56, 158, 85, 22, 3, 215, 173, 1,
		},
		pmt)
	pmt, err = api.GetPartialMerkleTree(testnetWitnessTestTxHash)
	require.NoError(t, err)
	assert.Equal(t,
		[]byte{
			6, 4, 0, 0, 12, 57, 157, 211, 242, 2, 14, 119, 98, 227, 64, 109, 65, 224, 44, 139, 205, 254, 97, 203, 220, 209,
			254, 64, 12, 94, 182, 134, 176, 73, 161, 9, 78, 251, 191, 13, 139, 34, 244, 217, 55, 250, 112, 170, 25, 54, 178,
			227, 131, 30, 50, 50, 72, 134, 119, 126, 198, 157, 52, 74, 181, 121, 153, 117, 35, 236, 9, 149, 2, 45, 233, 151,
			96, 179, 240, 141, 236, 56, 147, 12, 242, 219, 111, 154, 154, 44, 242, 70, 35, 34, 27, 217, 204, 193, 203, 173,
			92, 248, 119, 38, 29, 41, 122, 65, 113, 78, 218, 253, 54, 48, 1, 26, 230, 168, 43, 76, 193, 77, 119, 239, 225, 227,
			193, 205, 229, 226, 208, 13, 80, 11, 174, 159, 83, 90, 23, 115, 148, 213, 115, 5, 193, 144, 93, 32, 12, 49, 105, 180,
			95, 107, 188, 18, 246, 66, 15, 135, 129, 234, 145, 222, 100, 60, 131, 94, 110, 160, 34, 109, 55, 122, 219, 68, 185,
			65, 207, 116, 40, 54, 179, 149, 245, 243, 114, 201, 69, 210, 237, 19, 23, 13, 193, 244, 83, 112, 43, 112, 250, 221,
			239, 182, 162, 137, 163, 199, 150, 140, 235, 125, 142, 254, 169, 194, 5, 202, 4, 224, 194, 252, 144, 246, 245, 217,
			5, 82, 130, 202, 249, 251, 225, 189, 208, 190, 246, 65, 123, 37, 52, 255, 10, 29, 131, 167, 127, 171, 54, 181, 151,
			127, 246, 242, 173, 199, 200, 221, 143, 232, 65, 65, 65, 209, 46, 168, 180, 14, 124, 58, 246, 71, 161, 71, 50, 4,
			216, 59, 87, 38, 188, 235, 138, 54, 169, 139, 181, 60, 235, 77, 68, 25, 182, 2, 125, 95, 34, 183, 83, 219, 67, 33, 76,
			127, 115, 205, 1, 101, 251, 121, 74, 9, 34, 152, 131, 240, 101, 1, 68, 216, 93, 194, 13, 111, 69, 183, 83, 98, 136, 128,
			229, 196, 252, 57, 214, 62, 144, 128, 41, 34, 30, 249, 147, 12, 201, 136, 227, 223, 111, 145, 179, 232, 135, 219, 52,
			139, 13, 138, 137, 254, 2, 114, 59, 80, 226, 147, 7, 3, 160, 70, 231, 81, 171, 188, 100, 193, 106, 164, 143, 109, 225,
			106, 24, 196, 49, 56, 158, 85, 22, 3, 255, 54, 0,
		}, pmt)
}

func TestMempoolSpaceApi_GetPartialMerkleTree_MainnetBlock(t *testing.T) {
	api := mempool_space.NewMempoolSpaceApi(http.DefaultClient, &chaincfg.MainNetParams, mainnetUrl)
	for _, c := range datasets.PartialMerkleTrees {
		serializedPMT, err := api.GetPartialMerkleTree(c.Tx)
		require.NoError(t, err)
		result := hex.EncodeToString(serializedPMT)
		assert.Equal(t, c.Pmt, result)
	}
}

func TestMempoolSpaceApi_GetPartialMerkleTree_NotFound(t *testing.T) {
	api := mempool_space.NewMempoolSpaceApi(http.DefaultClient, &chaincfg.TestNet3Params, testnetUrl)
	result, err := api.GetPartialMerkleTree(txNotPresentInTestnet)
	require.ErrorContains(t, err, "Transaction not found")
	assert.Empty(t, result)
}

func TestMempoolSpaceApi_GetHeight(t *testing.T) {
	httpClient := http.DefaultClient
	mainnet := mempool_space.NewMempoolSpaceApi(httpClient, &chaincfg.MainNetParams, mainnetUrl)
	testnet := mempool_space.NewMempoolSpaceApi(httpClient, &chaincfg.TestNet3Params, testnetUrl)

	mainnetHeight, mainnetErr := mainnet.GetHeight()
	testnetHeight, testnetErr := testnet.GetHeight()

	require.NoError(t, testnetErr)
	require.NoError(t, mainnetErr)
	assert.Less(t, big.NewInt(903475).Cmp(mainnetHeight), 1)
	assert.Less(t, big.NewInt(4548634).Cmp(testnetHeight), 1)
}

func TestMempoolSpaceApi_BuildMerkleBranch(t *testing.T) {
	httpClient := http.DefaultClient

	t.Run("Should build merkle branch for testnet transactions", func(t *testing.T) {
		testnet := mempool_space.NewMempoolSpaceApi(httpClient, &chaincfg.TestNet3Params, testnetUrl)
		legacyBranch, legacyErr := testnet.BuildMerkleBranch(datasets.TestnetLegacyMerkleBranch.Tx)
		require.NoError(t, legacyErr)
		assert.Equal(t, datasets.TestnetLegacyMerkleBranch.Branch, legacyBranch)
		witnessBranch, witnessErr := testnet.BuildMerkleBranch(datasets.TestnetWitnessMerkleBranch.Tx)
		require.NoError(t, witnessErr)
		assert.Equal(t, datasets.TestnetWitnessMerkleBranch.Branch, witnessBranch)
	})
	t.Run("Should build merkle branch for mainnet transactions", func(t *testing.T) {
		mainnet := mempool_space.NewMempoolSpaceApi(httpClient, &chaincfg.MainNetParams, mainnetUrl)
		legacyBranch, legacyErr := mainnet.BuildMerkleBranch(datasets.MainnetLegacyMerkleBranch.Tx)
		require.NoError(t, legacyErr)
		assert.Equal(t, datasets.MainnetLegacyMerkleBranch.Branch, legacyBranch)
		witnessBranch, witnessErr := mainnet.BuildMerkleBranch(datasets.MainnetWitnessMerkleBranch.Tx)
		require.NoError(t, witnessErr)
		assert.Equal(t, datasets.MainnetWitnessMerkleBranch.Branch, witnessBranch)
	})
}

func TestMempoolSpaceApi_BuildMerkleBranch_TxNotFound(t *testing.T) {
	api := mempool_space.NewMempoolSpaceApi(http.DefaultClient, &chaincfg.TestNet3Params, testnetUrl)
	branch, err := api.BuildMerkleBranch(txNotPresentInTestnet)
	require.ErrorContains(t, err, "Transaction not found")
	assert.Equal(t, blockchain.MerkleBranch{}, branch)
}

func TestMempoolSpaceApi_GetTransactionBlockInfo(t *testing.T) {
	api := mempool_space.NewMempoolSpaceApi(http.DefaultClient, &chaincfg.TestNet3Params, testnetUrl)
	result, err := api.GetTransactionBlockInfo(testnetTestTxHash)
	require.NoError(t, err)
	assert.Equal(t, [32]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x1e, 0x94, 0xd8, 0x5c, 0x3e, 0x73, 0x6a, 0xa4, 0x7, 0x1d, 0x36, 0xd2, 0x65, 0x47, 0x71, 0x38, 0x20, 0xa2, 0x7a, 0xf9, 0xed, 0xbe, 0x97, 0x48, 0x9c, 0x69, 0x6f}, result.Hash)
	assert.Equal(t, big.NewInt(2582756), result.Height)
	assert.Equal(t, time.Unix(1710931198, 0), result.Time)
}

func TestMempoolSpaceApi_GetTransactionBlockInfo_NotFound(t *testing.T) {
	api := mempool_space.NewMempoolSpaceApi(http.DefaultClient, &chaincfg.TestNet3Params, testnetUrl)
	result, err := api.GetTransactionBlockInfo(txNotPresentInTestnet)
	require.ErrorContains(t, err, "Transaction not found")
	assert.Empty(t, result)
}

func TestMempoolSpaceApi_GetCoinbaseInformation(t *testing.T) {
	var blockHash, witnessMerkleRoot [32]byte
	api := mempool_space.NewMempoolSpaceApi(http.DefaultClient, &chaincfg.TestNet3Params, testnetUrl)
	coinbaseInfo, err := api.GetCoinbaseInformation(testnetWitnessTestTxHash)
	require.NoError(t, err)
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
	assert.Equal(t, rootstock.BtcCoinbaseTransactionInformation{
		BtcTxSerialized:      tx,
		BlockHash:            blockHash,
		BlockHeight:          big.NewInt(2582756),
		SerializedPmt:        pmt,
		WitnessMerkleRoot:    witnessMerkleRoot,
		WitnessReservedValue: [32]byte{},
	}, coinbaseInfo)
}

func TestMempoolSpaceApi_GetCoinbaseInformation_NotFound(t *testing.T) {
	api := mempool_space.NewMempoolSpaceApi(http.DefaultClient, &chaincfg.TestNet3Params, testnetUrl)
	result, err := api.GetCoinbaseInformation(txNotPresentInTestnet)
	require.ErrorContains(t, err, "Transaction not found")
	assert.Empty(t, result)
}

func TestMempoolSpaceApi_GetBlockchainInfo(t *testing.T) {
	client := http.DefaultClient
	mainnet := mempool_space.NewMempoolSpaceApi(client, &chaincfg.MainNetParams, mainnetUrl)
	testnet := mempool_space.NewMempoolSpaceApi(client, &chaincfg.TestNet3Params, testnetUrl)

	mainnetInfo, err := mainnet.GetBlockchainInfo()
	require.NoError(t, err)
	testnetInfo, err := testnet.GetBlockchainInfo()
	require.NoError(t, err)

	assert.Equal(t, "mainnet", mainnetInfo.NetworkName)
	assert.Len(t, mainnetInfo.BestBlockHash, 64)
	assert.Less(t, big.NewInt(903566).Cmp(mainnetInfo.ValidatedHeaders), 1)
	assert.Less(t, big.NewInt(903566).Cmp(mainnetInfo.ValidatedBlocks), 1)

	assert.Equal(t, "testnet3", testnetInfo.NetworkName)
	assert.Len(t, testnetInfo.BestBlockHash, 64)
	assert.Less(t, big.NewInt(4548823).Cmp(testnetInfo.ValidatedHeaders), 1)
	assert.Less(t, big.NewInt(4548823).Cmp(testnetInfo.ValidatedBlocks), 1)
}
