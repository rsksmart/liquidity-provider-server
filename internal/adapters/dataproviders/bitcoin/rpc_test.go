package bitcoin_test

import (
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	testnetTestBlockHash = "00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f"
	testnetTestTxHash    = "9f0706c2717fc77bf0f225a4223933a7decb8d36902ddbb0accab8ea894f8b29"
	witnessTestTxHash    = "5cadcbc1ccd91b222346f22c9a9a6fdbf20c9338ec8df0b36097e92d029509ec"
	testnetBlockFile     = "block-2582756-testnet.txt"
	mainnetBlockFile     = "block-696394-mainnet.txt"
)

// nolint:funlen
func TestBitcoindRpc_ValidateAddress(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	mainnet := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	testnet := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.TestNet3Params, client))
	regtest := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.RegressionNetParams, client))
	invalid := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.SigNetParams, client))

	p2pkhMainnet := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
	p2pkhTestnet := "mipcBbFg9gMiCh81Kj8tqqdgoZub1ZJRfn"
	p2shMainnet := "3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy"
	p2shTestnet := "2MzQwSSnBHWHqSAqtTVQ6v47XtaisrJa1Vc"
	p2wpkhMainnet := "bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq"
	p2wpkhTestnet := "tb1qg9stx3w3xj0t5j8d5q2g4yvz5gj8v3tjxjxw5v"
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

	require.ErrorIs(t, regtest.ValidateAddress(p2pkhMainnet), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, regtest.ValidateAddress(p2shMainnet), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, regtest.ValidateAddress(p2wpkhMainnet), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, regtest.ValidateAddress(p2wshMainnet), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, regtest.ValidateAddress(p2trMainnet), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, regtest.ValidateAddress(p2wpkhTestnet), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, regtest.ValidateAddress(p2wshTestnet), blockchain.BtcAddressInvalidNetworkError)
	require.ErrorIs(t, regtest.ValidateAddress(p2trTestnet), blockchain.BtcAddressInvalidNetworkError)
	require.NoError(t, regtest.ValidateAddress(p2pkhTestnet))
	require.NoError(t, regtest.ValidateAddress(p2shTestnet))
	require.NoError(t, regtest.ValidateAddress(p2wpkhRegtest))
	require.NoError(t, regtest.ValidateAddress(p2wshRegtest))
	require.NoError(t, regtest.ValidateAddress(p2trRegtest))

	const unsupportedNetwork = "unsupported network"
	require.Contains(t, invalid.ValidateAddress(p2pkhMainnet).Error(), unsupportedNetwork)
	require.Contains(t, invalid.ValidateAddress(p2shMainnet).Error(), unsupportedNetwork)
	require.Contains(t, invalid.ValidateAddress(p2pkhTestnet).Error(), unsupportedNetwork)
	require.Contains(t, invalid.ValidateAddress(p2shTestnet).Error(), unsupportedNetwork)
	require.Contains(t, invalid.ValidateAddress(p2wpkhMainnet).Error(), unsupportedNetwork)
	require.Contains(t, invalid.ValidateAddress(p2wpkhTestnet).Error(), unsupportedNetwork)
}

func TestBitcoindRpc_GetHeight(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))

	client.On("GetBlockChainInfo").Return(&btcjson.GetBlockChainInfoResult{Blocks: 123}, nil).Once()
	client.On("GetBlockChainInfo").Return(&btcjson.GetBlockChainInfoResult{}, assert.AnError).Once()

	height, err := rpc.GetHeight()
	require.NoError(t, err)
	require.Equal(t, big.NewInt(123), height)

	height, err = rpc.GetHeight()
	require.Error(t, err)
	require.Nil(t, height)
}

func TestBitcoindRpc_DecodeAddress(t *testing.T) {
	var decodedAddresses []decodedAddress
	decodedAddresses = append(decodedAddresses, base58Addresses...)
	decodedAddresses = append(decodedAddresses, bech32Addresses...)
	decodedAddresses = append(decodedAddresses, bech32mAddresses...)
	client := &mocks.ClientAdapterMock{}
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	cases := decodedAddresses
	for _, c := range cases {
		decoded, err := rpc.DecodeAddress(c.address)
		require.NoError(t, err)
		assert.Equal(t, c.expected, decoded)
	}
}

func TestBitcoindRpc_GetRawTransaction(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	txBytes, _ := hex.DecodeString("0200000002ebf7c22a73f3baea460cad53a2788bd4f24020f6b374900a771d3422f128442e000000006a473044022062dae13ba281d0cf529b604bb59c1efcd7b83438af34d4a51acc6f31041be18c022044df281e688a52624f45f6c26662349d1f5efedd4d69530e65b7d7cec0d3792d0121038e509bc056004a5da7460b5acd5d4dcb2add41d53817180499e3814290ecc91efdffffffb5f09f38215b850f4ba644a7f7ab57efa8d10c5f4b5908e9aa980ff5ffa948f5000000006a47304402206538fc72b896e4c6e807a4daf56191f68dec307c3011d082e69eeb3d45d6d8c302205a329814ab87901ae56a82587e716fa2282ecc665ab203da14d93db71181ecd8012102498a833095175800f40b2c0ab23f108b47a319a94ccea826062bf66c827e91a9fdffffff0298740700000000001976a91473cce22e78ec61cd54a6438ca1210b88561ebcdd88ac20a10700000000001976a9142c81478132b5dda64ffc484a0d225096c4b22ad588acc3682700")
	client.On("GetRawTransaction", mock.Anything).Return(btcutil.NewTxFromBytes(txBytes)).Once()
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	result, err := rpc.GetRawTransaction(testnetTestTxHash)
	require.NoError(t, err)
	require.Equal(t, txBytes, result)
}

func TestBitcoindRpc_GetRawTransaction_FromBlock(t *testing.T) {
	mainnetBlock := getTestBlock(t, mainnetBlockFile)
	tx, err := mainnetBlock.Tx(0)
	require.NoError(t, err)
	client := &mocks.ClientAdapterMock{}
	client.On("GetRawTransaction", mock.Anything).Return(tx, nil).Once()
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	result, err := rpc.GetRawTransaction(tx.Hash().String())
	require.NoError(t, err)
	expectedBytes, err := hex.DecodeString(
		"01000000010000000000000000000000000000000000000000000000000000000000000000" +
			"ffffffff5f034aa00a1c2f5669614254432f4d696e656420627920797a33313936303538372f2cfabe" +
			"6d6dc0a3751203a336deb817199448996ebcb2a0e537b1ce9254fa3e9c3295ca196b10000000000000" +
			"0010c56e6700d262d24bd851bb829f9f0000ffffffff0401b3cc25000000001976a914536ffa992491" +
			"508dca0354e52f32a3a7a679a53a88ac00000000000000002b6a2952534b424c4f434b3a040c866ad2" +
			"fdb8b59b32dd17059edaeef11d295e279a74ab97125d2500371ce90000000000000000266a24b9e11b" +
			"6dab3e2ca50c1a6b01cf80eccb9d291aab8b095d653e348aa9d94a73964ff5cf1b0000000000000000" +
			"266a24aa21a9ed04f0bac0104f4fa47bec8058f2ebddd292dd85027ab0d6d95288d31f12c5a4b800000000",
	)
	require.NoError(t, err)
	require.Equal(t, expectedBytes, result)
}

func TestBitcoindRpc_GetRawTransaction_ErrorHandling(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	txBytes, _ := hex.DecodeString("0200000002ebf7c22a73f3baea460cad53a2788bd4f24020f6b374900a771d3422f128442e000000006a473044022062dae13ba281d0cf529b604bb59c1efcd7b83438af34d4a51acc6f31041be18c022044df281e688a52624f45f6c26662349d1f5efedd4d69530e65b7d7cec0d3792d0121038e509bc056004a5da7460b5acd5d4dcb2add41d53817180499e3814290ecc91efdffffffb5f09f38215b850f4ba644a7f7ab57efa8d10c5f4b5908e9aa980ff5ffa948f5000000006a47304402206538fc72b896e4c6e807a4daf56191f68dec307c3011d082e69eeb3d45d6d8c302205a329814ab87901ae56a82587e716fa2282ecc665ab203da14d93db71181ecd8012102498a833095175800f40b2c0ab23f108b47a319a94ccea826062bf66c827e91a9fdffffff0298740700000000001976a91473cce22e78ec61cd54a6438ca1210b88561ebcdd88ac20a10700000000001976a9142c81478132b5dda64ffc484a0d225096c4b22ad588acc3682700")
	client.On("GetRawTransaction", mock.Anything).Return(btcutil.NewTxFromBytes(txBytes)).Once()
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))

	_, err := rpc.GetRawTransaction("invalidHash")
	require.Error(t, err)

	client = &mocks.ClientAdapterMock{}
	client.On("GetRawTransaction", mock.Anything).Return(nil, assert.AnError).Once()
	rpc = bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	_, err = rpc.GetRawTransaction(testnetTestTxHash)
	require.Error(t, err)

	client = &mocks.ClientAdapterMock{}
	client.On("GetRawTransaction", mock.Anything).Return(btcutil.NewTxFromBytes([]byte{01, 02, 03})).Once()
	rpc = bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	_, err = rpc.GetRawTransaction(testnetTestTxHash)
	require.Error(t, err)
}

func TestBitcoindRpc_GetTransactionBlockInfo(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	now := time.Now()
	client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: testnetTestBlockHash}, nil).Once()
	client.On("GetBlockVerbose", mock.Anything).Return(&btcjson.GetBlockVerboseResult{Height: 123, Time: now.Unix()}, nil).Once()
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	result, err := rpc.GetTransactionBlockInfo(testnetTestTxHash)
	require.NoError(t, err)
	assert.Equal(t, [32]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x1e, 0x94, 0xd8, 0x5c, 0x3e, 0x73, 0x6a, 0xa4, 0x7, 0x1d, 0x36, 0xd2, 0x65, 0x47, 0x71, 0x38, 0x20, 0xa2, 0x7a, 0xf9, 0xed, 0xbe, 0x97, 0x48, 0x9c, 0x69, 0x6f}, result.Hash)
	assert.Equal(t, big.NewInt(123), result.Height)
	assert.WithinDuration(t, now, result.Time, 1*time.Second)
}

func TestBitcoindRpc_GetTransactionBlockInfo_ErrorHandling(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	result, err := rpc.GetTransactionBlockInfo("txhash")
	assert.Equal(t, blockchain.BitcoinBlockInformation{}, result)
	require.Error(t, err)

	client = &mocks.ClientAdapterMock{}
	client.On("GetRawTransactionVerbose", mock.Anything).Return(nil, assert.AnError).Once()
	rpc = bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	result, err = rpc.GetTransactionBlockInfo(testnetTestTxHash)
	assert.Equal(t, blockchain.BitcoinBlockInformation{}, result)
	require.Error(t, err)

	client = &mocks.ClientAdapterMock{}
	client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: testnetTestBlockHash}, nil).Once()
	client.On("GetBlockVerbose", mock.Anything).Return(nil, assert.AnError).Once()
	rpc = bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	result, err = rpc.GetTransactionBlockInfo(testnetTestTxHash)
	assert.Equal(t, blockchain.BitcoinBlockInformation{}, result)
	require.Error(t, err)

	client = &mocks.ClientAdapterMock{}
	client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: "blk"}, nil).Once()
	rpc = bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	result, err = rpc.GetTransactionBlockInfo(testnetTestTxHash)
	assert.Equal(t, blockchain.BitcoinBlockInformation{}, result)
	require.Error(t, err)
}

func TestBitcoindRpc_BuildMerkleBranch(t *testing.T) {
	block := getTestBlock(t, testnetBlockFile)

	client := &mocks.ClientAdapterMock{}
	client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: testnetTestBlockHash}, nil).Once()
	client.On("GetBlock", mock.Anything).Return(block.MsgBlock(), nil).Once()
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	branch, err := rpc.BuildMerkleBranch(testnetTestTxHash)
	require.NoError(t, err)
	assert.Equal(t, blockchain.MerkleBranch{
		Hashes: [][32]byte{
			{155, 80, 207, 191, 224, 254, 254, 207, 78, 62, 249, 222, 89, 87, 171, 35, 220, 207, 189, 31, 70, 82, 141, 48, 198, 32, 160, 219, 132, 222, 4, 8},
			{71, 170, 231, 32, 77, 30, 186, 115, 142, 85, 108, 168, 0, 13, 46, 49, 84, 233, 136, 89, 60, 43, 243, 202, 144, 62, 255, 213, 141, 194, 189, 179},
			{59, 162, 140, 60, 248, 2, 245, 106, 154, 191, 234, 177, 48, 236, 162, 182, 251, 183, 83, 235, 29, 21, 107, 125, 34, 114, 26, 64, 162, 84, 126, 120},
			{93, 7, 146, 22, 74, 120, 71, 158, 154, 141, 202, 163, 154, 161, 141, 251, 221, 203, 104, 72, 74, 252, 21, 254, 64, 150, 96, 172, 63, 160, 41, 97},
			{234, 116, 222, 241, 199, 162, 201, 219, 87, 174, 86, 69, 151, 193, 247, 143, 142, 205, 242, 138, 20, 53, 19, 208, 210, 50, 150, 113, 181, 67, 117, 177},
			{141, 182, 8, 250, 221, 182, 182, 192, 127, 135, 114, 87, 57, 169, 102, 200, 136, 177, 0, 83, 135, 209, 203, 85, 237, 80, 109, 235, 151, 92, 88, 192},
			{23, 44, 34, 196, 81, 13, 32, 151, 5, 75, 11, 104, 32, 13, 151, 201, 99, 35, 250, 136, 32, 246, 156, 232, 196, 199, 28, 210, 227, 241, 116, 67},
			{56, 133, 146, 188, 185, 209, 23, 73, 20, 41, 218, 247, 211, 165, 219, 89, 80, 135, 219, 133, 198, 55, 47, 72, 23, 8, 219, 209, 63, 211, 217, 117},
			{95, 15, 80, 149, 169, 116, 91, 201, 28, 85, 231, 232, 222, 112, 145, 6, 33, 235, 81, 88, 148, 191, 165, 186, 206, 116, 16, 165, 252, 48, 10, 29},
			{13, 139, 52, 219, 135, 232, 179, 145, 111, 223, 227, 136, 201, 12, 147, 249, 30, 34, 41, 128, 144, 62, 214, 57, 252, 196, 229, 128, 136, 98, 83, 183},
			{22, 85, 158, 56, 49, 196, 24, 106, 225, 109, 143, 164, 106, 193, 100, 188, 171, 81, 231, 70, 160, 3, 7, 147, 226, 80, 59, 114, 2, 254, 137, 138},
		},
		Path: big.NewInt(406),
	}, branch)
}

func TestBitcoindRpc_BuildMerkleBranch_ErrorHandling(t *testing.T) {
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, &mocks.ClientAdapterMock{}))
	branch, err := rpc.BuildMerkleBranch("txhash")
	require.Error(t, err)
	require.Equal(t, blockchain.MerkleBranch{}, branch)

	client := &mocks.ClientAdapterMock{}
	client.On("GetRawTransactionVerbose", mock.Anything).Return(nil, assert.AnError).Once()
	rpc = bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	branch, err = rpc.BuildMerkleBranch(testnetTestTxHash)
	require.Error(t, err)
	require.Equal(t, blockchain.MerkleBranch{}, branch)

	client = &mocks.ClientAdapterMock{}
	client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: testnetTestBlockHash}, nil).Once()
	client.On("GetBlock", mock.Anything).Return(&wire.MsgBlock{}, assert.AnError).Once()
	rpc = bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	branch, err = rpc.BuildMerkleBranch(testnetTestTxHash)
	require.Error(t, err)
	require.Equal(t, blockchain.MerkleBranch{}, branch)

	client = &mocks.ClientAdapterMock{}
	client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: "blkhash"}, nil).Once()
	rpc = bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	branch, err = rpc.BuildMerkleBranch(testnetTestTxHash)
	require.Error(t, err)
	require.Equal(t, blockchain.MerkleBranch{}, branch)
}

func TestBitcoindRpc_BuildMerkleBranch_TxNotFound(t *testing.T) {
	block := getTestBlock(t, testnetBlockFile)

	client := &mocks.ClientAdapterMock{}
	client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: testnetTestBlockHash}, nil).Once()
	client.On("GetBlock", mock.Anything).Return(block.MsgBlock(), nil).Once()
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	branch, err := rpc.BuildMerkleBranch("9dd8911176857dff8244f75f7c95782b3495048ad75632f0a58c8e942cefb223")
	require.Error(t, err)
	assert.Equal(t, "transaction 9dd8911176857dff8244f75f7c95782b3495048ad75632f0a58c8e942cefb223 not found in merkle tree", err.Error())
	assert.Equal(t, blockchain.MerkleBranch{}, branch)
}

func TestBitcoindRpc_GetPartialMerkleTree(t *testing.T) {
	block := getTestBlock(t, testnetBlockFile)
	client := &mocks.ClientAdapterMock{}
	client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: testnetTestBlockHash}, nil).Twice()
	client.On("GetBlock", mock.Anything).Return(block.MsgBlock(), nil).Twice()
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.TestNet3Params, client))
	pmt, err := rpc.GetPartialMerkleTree(testnetTestTxHash)
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
	pmt, err = rpc.GetPartialMerkleTree(witnessTestTxHash)
	require.NoError(t, err)
	assert.Equal(t,
		[]byte{
			6, 4, 0, 0, 12, 58, 105, 110, 101, 15, 121, 221, 98, 48, 193, 196, 7, 210, 226, 25, 24, 11, 151, 82, 5, 90,
			125, 199, 5, 39, 243, 148, 209, 81, 45, 62, 11, 231, 132, 178, 158, 236, 169, 194, 214, 64, 17, 213, 126, 96,
			177, 12, 169, 128, 239, 100, 23, 187, 234, 136, 228, 130, 127, 66, 74, 182, 29, 24, 77, 101, 238, 142, 185, 83,
			221, 231, 68, 128, 143, 19, 65, 205, 234, 63, 12, 21, 150, 177, 4, 231, 65, 218, 211, 155, 236, 58, 187, 150,
			157, 88, 18, 101, 41, 41, 226, 101, 50, 151, 79, 20, 110, 100, 49, 230, 88, 121, 131, 72, 72, 20, 77, 21, 10,
			177, 47, 175, 163, 197, 197, 128, 6, 209, 202, 44, 28, 8, 153, 39, 77, 77, 10, 252, 90, 217, 130, 242, 218, 121,
			28, 149, 67, 204, 128, 48, 169, 116, 193, 31, 240, 7, 101, 152, 123, 108, 122, 82, 96, 76, 121, 78, 241, 15, 81,
			154, 24, 226, 235, 122, 43, 87, 82, 223, 211, 0, 175, 198, 57, 136, 62, 82, 232, 122, 12, 239, 98, 182, 167, 91,
			155, 226, 127, 38, 127, 96, 103, 106, 68, 194, 197, 29, 190, 22, 189, 205, 89, 124, 2, 98, 117, 61, 187, 236, 54,
			244, 134, 57, 223, 117, 211, 239, 220, 221, 136, 69, 1, 37, 178, 60, 3, 55, 54, 93, 65, 50, 19, 34, 11, 35, 126,
			59, 156, 204, 80, 219, 124, 32, 203, 225, 212, 230, 200, 81, 38, 155, 25, 167, 69, 29, 31, 144, 153, 133, 133, 153,
			85, 113, 199, 21, 158, 142, 197, 114, 214, 191, 221, 118, 40, 246, 30, 161, 95, 13, 189, 255, 12, 252, 85, 227, 29,
			85, 238, 154, 56, 172, 231, 53, 185, 141, 0, 117, 191, 226, 148, 39, 15, 155, 237, 91, 199, 215, 159, 139, 139, 162,
			30, 64, 166, 155, 191, 199, 224, 122, 191, 147, 126, 16, 165, 228, 135, 156, 99, 40, 148, 225, 175, 179, 89, 69, 245,
			228, 74, 175, 233, 131, 142, 1, 72, 249, 14, 55, 200, 245, 139, 103, 118, 164, 205, 255, 58, 85, 93, 210, 44, 198,
			232, 46, 206, 205, 126, 150, 115, 155, 235, 106, 94, 22, 241, 164, 152, 3, 255, 54, 0,
		}, pmt)
	client.AssertExpectations(t)
}

func TestBitcoindRpc_BuildMerkleBranch_MainnetBlock(t *testing.T) {
	cases := [2]struct {
		tx  string
		pmt string
	}{
		// first two are witness tx
		// TODO add non-witness scenario in the task to update the merkle proofs
		{
			tx: "07f8b22fa9a3b32e20b59bb90727de05fb634749519ebcb6a887aeaf2c7eb041",
			pmt: "f30800000d" +
				"000000000000000000000000000000000000000000000000000000000000000073" +
				"1c671fafb5d234834c726657f29c9af030ccf7068f1ef732af4efd8e146da0a9d6" +
				"075f4758821ceeef2c230cfd2497df2d1d1d02dd19e653d22b3dc271b3931b9aac" +
				"9d2faa0e1814b3ac46067e68ee5fa59b9e5e9f5eb60c6d00a746848af913ce9302" +
				"f52718a61abc3af42ed88a07341b19364cbac956a4519115837569da7bf6e51ea9" +
				"c98d57b3f8d3a781fb160ca04ebed21e80444036b70446f04a7838366fee70449b" +
				"9e89850eb3d54c2b81b63674aec107e16757ab404f6b8a35e775b79c5fbf73f954" +
				"07dba5a5913887ea9190a922ef1578dacc2b1c5e923de965acc58ab8a7f1e7fbed" +
				"49bd738f3f32cfca587c9ae1ab37cd0ac3833a8f1b172df9087fbe1a73f22e586d" +
				"18dda0c94282ea1e8e08a7ee8920bd7be7fc2dbb160dcb7a39ee8fbdc7c3f6421f" +
				"79945ab823ed00f69a473c28c6701f31569bd7ecf3702fe8ae7d2a9ed374b344f6" +
				"5c4f98482a79d301048b6664bfd85f4453249449af6b6647ec2183f6fcf269a542" +
				"8d55bc462abe7cee6de0c0d4c9ce048530a9135704ff1f0000",
		},
		{
			tx: "ddf5061f9707f0c959bf24278d557b264716672c1b601ec50112d6dfe160d9d3",
			pmt: "f30800000d" +
				"ae0f1f5cf24d3875e8f301fe178d8efd527c668d176709cf27730572392c4ad9" +
				"681d5442690792835378e68031a399bd8222e7b3d26b6ddadee237890087b516" +
				"820f677871c76a1b0cadd7d3315a413511c66feaa6f87ec9e820cd1479c5a054" +
				"f3ebc63fcc05543c5cf652d2aea6ab5f9be87cb91cd34bbadb9b04ba023f791a" +
				"3dfb1303b5a550d33774cd2dbb37f8d2522935f218e8b13e1dfdc1a0e824fb38" +
				"0d371f17a4c6328ec0860960bad8b6296483d55c155d216105ec0ade23da6c66" +
				"ad9c51f5f731c8f08214a3dd7f6eddae5a55076295190723c462ecfbfbeaeffb" +
				"ed89a4d28cdaec3e5039438d5134737dceeceff340d6920fab7518757d99e1c7" +
				"f305a80c97b49df6c0d5dc28c47c1bfc595c47bc820354a94b2d2d8bd6f75466" +
				"fbfb473898f2a5840d86338deae45a04f3912416b7e5526e37ce6842d43b05e7" +
				"a5e0ed8adc181d8e026ed6dc27a31b8b2729b76902b4e1a8758f2c70bbc3a442" +
				"0bf154f0aa8b60b415eaf70a9ada542f39b62baca6415c3611f8306e14a4d131" +
				"6b6647ec2183f6fcf269a5428d55bc462abe7cee6de0c0d4c9ce048530a91357" +
				"04dbd50300",
		},
	}
	mainnetBlock := getTestBlock(t, mainnetBlockFile)
	for _, c := range cases {
		client := &mocks.ClientAdapterMock{}
		client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: testnetTestBlockHash}, nil).Once()
		client.On("GetBlock", mock.Anything).Return(mainnetBlock.MsgBlock(), nil).Once()
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
		serializedPMT, err := rpc.GetPartialMerkleTree(c.tx)
		require.NoError(t, err)
		result := hex.EncodeToString(serializedPMT)
		assert.Equal(t, c.pmt, result)
	}
}

func TestBitcoindRpc_GetPartialMerkleTree_ErrorHandling(t *testing.T) {
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, &mocks.ClientAdapterMock{}))
	pmt, err := rpc.GetPartialMerkleTree("txhash")
	require.Error(t, err)
	require.Nil(t, pmt)

	client := &mocks.ClientAdapterMock{}
	block := getTestBlock(t, testnetBlockFile)
	msgBlock := block.MsgBlock()
	msgBlock.Transactions = append(msgBlock.Transactions, msgBlock.Transactions...)
	client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: testnetTestBlockHash}, nil).Once()
	client.On("GetBlock", mock.Anything).Return(msgBlock, nil).Once()
	rpc = bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	result, err := rpc.GetPartialMerkleTree(testnetTestTxHash)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "block matches more than one transaction")
	assert.Nil(t, result)
}

func TestBitcoindRpc_GetTransactionInfo(t *testing.T) {
	receivedTxPath, err := filepath.Abs("../../../../test/mocks/getRawTransactionVerboseReceived.json")
	require.NoError(t, err)
	receivedTxResponse, err := os.ReadFile(receivedTxPath)
	require.NoError(t, err)
	txReceiveDetails := btcjson.TxRawResult{}
	require.NoError(t, json.Unmarshal(receivedTxResponse, &txReceiveDetails))
	client := &mocks.ClientAdapterMock{}
	client.On("GetRawTransactionVerbose", mock.Anything).Return(&txReceiveDetails, nil).Once()
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.TestNet3Params, client))
	result, err := rpc.GetTransactionInfo(testnetTestTxHash)
	require.NoError(t, err)
	assert.Equal(t, blockchain.BitcoinTransactionInformation{
		Hash: testnetTestTxHash, Confirmations: uint64(105277),
		Outputs: map[string][]*entities.Wei{
			"mr5FSft4PQvoWbf9Sf5iQXbw1u445iNksp": {entities.NewWei(0.004886 * 1e18)}, "mjaGtyj74LYn7gApr17prZxDPDnfuUnRa5": {entities.NewWei(0.005 * 1e18)},
		},
	}, result)

	sentTxPath, err := filepath.Abs("../../../../test/mocks/getRawTransactionVerboseSent.json")
	require.NoError(t, err)
	sentTxResponse, err := os.ReadFile(sentTxPath)
	require.NoError(t, err)
	const sendTxHash = "9b0c48b79fe40c67f7a2837e6e59a138a16671caf7685dcd831bd3c51b9f6d21"
	txSendDetails := btcjson.TxRawResult{}
	err = json.Unmarshal(sentTxResponse, &txSendDetails)
	require.NoError(t, err)
	client = &mocks.ClientAdapterMock{}
	client.On("GetRawTransactionVerbose", mock.Anything).Return(&txSendDetails, nil).Once()
	rpc = bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.TestNet3Params, client))
	result, err = rpc.GetTransactionInfo(sendTxHash)
	require.NoError(t, err)
	assert.Equal(t, blockchain.BitcoinTransactionInformation{
		Hash: sendTxHash, Confirmations: uint64(106306),
		Outputs: map[string][]*entities.Wei{
			"mqbKtarYKnoEdPheFFDGRjksvEpb2vJGNh": {entities.NewWei(0.005 * 1e18)},
			"mowfvQDraTDvRgZowL4tx5EatL1u78w65v": {entities.NewWei(0.01956600 * 1e18)},
			"":                                   {entities.NewWei(0)}, // Null data script output
		},
	}, result)

	witnessTxPath, err := filepath.Abs("../../../../test/mocks/getRawTransactionVerboseWitness.json")
	require.NoError(t, err)
	witnessTxResponse, err := os.ReadFile(witnessTxPath)
	require.NoError(t, err)
	const witnessTxHash = "0b9b2c99aa47b7effdd1d945a9ebad5374666ed13883960e80f2e8ff92088687"
	witnessTxDetails := btcjson.TxRawResult{}
	err = json.Unmarshal(witnessTxResponse, &witnessTxDetails)
	require.NoError(t, err)
	client = &mocks.ClientAdapterMock{}
	client.On("GetRawTransactionVerbose", mock.Anything).Return(&witnessTxDetails, nil).Once()
	rpc = bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.TestNet3Params, client))
	result, err = rpc.GetTransactionInfo(witnessTxHash)
	require.NoError(t, err)
	assert.Equal(t, blockchain.BitcoinTransactionInformation{
		Hash: witnessTxHash, Confirmations: uint64(286342), HasWitness: true,
		Outputs: map[string][]*entities.Wei{
			"tb1q5tsjcyz7xmet07yxtumakt739y53hcttmntajq": {entities.NewWei(0.00049899 * 1e18)},
			"tb1q460pja0n0qk0a0mzl0amde5lmp9an5wc9tv9yz": {entities.NewWei(0.03220659 * 1e18)},
		},
	}, result)
}

func TestBitcoindRpc_GetTransactionInfo_ErrorHandling(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	result, err := rpc.GetTransactionInfo("txhash")
	assert.Equal(t, blockchain.BitcoinTransactionInformation{}, result)
	require.Error(t, err)

	client = &mocks.ClientAdapterMock{}
	client.On("GetRawTransactionVerbose", mock.Anything).Return(nil, assert.AnError).Once()
	rpc = bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	result, err = rpc.GetTransactionInfo(testnetTestTxHash)
	assert.Equal(t, blockchain.BitcoinTransactionInformation{}, result)
	require.Error(t, err)

	client = &mocks.ClientAdapterMock{}
	client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{Vout: []btcjson.Vout{{Value: math.NaN()}}}, nil).Once()
	rpc = bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	result, err = rpc.GetTransactionInfo(testnetTestTxHash)
	assert.Empty(t, blockchain.BitcoinTransactionInformation{}, result)
	require.Error(t, err)
}

func TestBitcoindRpc_GetCoinbaseInformation(t *testing.T) {
	t.Run("Should handle error getting transaction", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.On("GetRawTransactionVerbose", mock.Anything).Return(nil, assert.AnError).Once()
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
		coinbaseInfo, err := rpc.GetCoinbaseInformation(witnessTestTxHash)
		assert.Empty(t, coinbaseInfo)
		require.Error(t, err)
	})
	t.Run("Should handle error getting block", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: testnetTestBlockHash}, nil).Once()
		client.On("GetBlock", mock.Anything).Return(nil, assert.AnError).Once()
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
		coinbaseInfo, err := rpc.GetCoinbaseInformation(witnessTestTxHash)
		assert.Empty(t, coinbaseInfo)
		require.Error(t, err)
	})
	t.Run("Should handle error getting block verbose", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		testnetBlock := getTestBlock(t, testnetBlockFile)
		client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: testnetTestBlockHash}, nil).Once()
		client.On("GetBlock", mock.Anything).Return(testnetBlock.MsgBlock(), nil).Once()
		client.On("GetBlockVerbose", mock.Anything).Return(nil, assert.AnError).Once()
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
		coinbaseInfo, err := rpc.GetCoinbaseInformation(witnessTestTxHash)
		assert.Empty(t, coinbaseInfo)
		require.Error(t, err)
	})
	t.Run("Should build coinbase info", func(t *testing.T) {
		var blockHash, witnessMerkleRoot [32]byte
		testnetBlock := getTestBlock(t, testnetBlockFile)
		client := &mocks.ClientAdapterMock{}
		client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: testnetTestBlockHash}, nil).Once()
		client.On("GetBlock", mock.Anything).Return(testnetBlock.MsgBlock(), nil).Once()
		client.On("GetBlockVerbose", mock.Anything).Return(&btcjson.GetBlockVerboseResult{Height: 2582756}, nil).Once()
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
		coinbaseInfo, err := rpc.GetCoinbaseInformation(witnessTestTxHash)
		require.NoError(t, err)
		blockHashBytes, _ := hex.DecodeString("00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f")
		copy(blockHash[:], blockHashBytes)
		witnessMerkleRootBytes, _ := hex.DecodeString("ce489a9f5cd7351ed57606c373329a208d37a8654bf473916eeec4b8ffee55f4")
		copy(witnessMerkleRoot[:], witnessMerkleRootBytes)
		pmt, _ := hex.DecodeString("060400000c89cc284039442363175d7f6d34492885618442a865381e6221a6223f5fa18af5ed018e20e419923fa29585d93161cae0a1a9ebb5095ca055e9888a8c0b956f4604987862683909d0749a969cbe1de8669b72c1174c0d4e8cd50238cc7198743148caa41a01a85c260ab184f3846125bda6e5f21c14dc314b207b4428c29cf6807977059737aa020de8c76fa6c304c99db6306a4a0aa1532ae7dbdb6453aa09d13c835e6ea0226d377adb44b941cf742836b395f5f372c945d2ed13170dc1f453702b70faddefb6a289a3c7968ceb7d8efea9c205ca04e0c2fc90f6f5d9055282caf9fbe1bdd0bef6417b2534ff0a1d83a77fab36b5977ff6f2adc7c8dd8fe8414141d12ea8b40e7c3af647a1473204d83b5726bceb8a36a98bb53ceb4d4419b6027d5f22b753db43214c7f73cd0165fb794a09229883f0650144d85dc20d6f45b753628880e5c4fc39d63e908029221ef9930cc988e3df6f91b3e887db348b0d8a89fe02723b50e2930703a046e751abbc64c16aa48f6de16a18c431389e551603ff0f00")
		tx, _ := hex.DecodeString("04000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1203e4682704eebcfa650843f719b701000000000000000377781d00000000001600143758d475313d557dbe8b1d90406c5b3b4dbbc00df79900000000000017a914a775ee7e3118ac67f181faca330f1d5c7658d205870000000000000000266a24aa21a9ed3dccc6b158d03b681f5cd8c71653097d0e6a51ac3e19de0add0a2a43419622a500000000")
		assert.Equal(t, blockchain.BtcCoinbaseTransactionInformation{
			BtcTxSerialized:      tx,
			BlockHash:            blockHash,
			BlockHeight:          big.NewInt(2582756),
			SerializedPmt:        pmt,
			WitnessMerkleRoot:    witnessMerkleRoot,
			WitnessReservedValue: [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		}, coinbaseInfo)
	})
}

func getTestBlock(t *testing.T, filename string) *btcutil.Block {
	absolutePath, err := filepath.Abs("../../../../test/mocks/" + filename)
	require.NoError(t, err)
	blockFile, err := os.ReadFile(absolutePath)
	require.NoError(t, err)
	blockBytes, err := hex.DecodeString(string(blockFile))
	require.NoError(t, err)
	block, err := btcutil.NewBlockFromBytes(blockBytes)
	require.NoError(t, err)
	return block
}
