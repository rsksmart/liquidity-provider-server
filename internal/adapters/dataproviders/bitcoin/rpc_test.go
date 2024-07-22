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
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	testnetTestBlockHash = "00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f"
	testnetTestTxHash    = "9f0706c2717fc77bf0f225a4223933a7decb8d36902ddbb0accab8ea894f8b29"
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
	client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: testnetTestBlockHash}, nil).Once()
	client.On("GetBlock", mock.Anything).Return(block.MsgBlock(), nil).Once()
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
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
}

func TestBitcoindRpc_BuildMerkleBranch_MainnetBlock(t *testing.T) {
	cases := [2]struct {
		tx  string
		pmt string
	}{
		{
			tx: "07f8b22fa9a3b32e20b59bb90727de05fb634749519ebcb6a887aeaf2c7eb041",
			pmt: "f30800000d" +
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
			tx: "ddf5061f9707f0c959bf24278d557b264716672c1b601ec50112d6dfe160d9d3",
			pmt: "f30800000d" +
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
	err = json.Unmarshal(receivedTxResponse, &txReceiveDetails)
	require.NoError(t, err)
	client := &mocks.ClientAdapterMock{}
	client.On("GetRawTransactionVerbose", mock.Anything).Return(&txReceiveDetails, nil).Once()
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.TestNet3Params, client))
	result, err := rpc.GetTransactionInfo(testnetTestTxHash)
	require.NoError(t, err)
	assert.Equal(t, blockchain.BitcoinTransactionInformation{
		Hash:          testnetTestTxHash,
		Confirmations: uint64(105277),
		Outputs: map[string][]*entities.Wei{
			"mr5FSft4PQvoWbf9Sf5iQXbw1u445iNksp": {entities.NewWei(0.004886 * 1e18)},
			"mjaGtyj74LYn7gApr17prZxDPDnfuUnRa5": {entities.NewWei(0.005 * 1e18)},
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
		Hash:          sendTxHash,
		Confirmations: uint64(106306),
		Outputs: map[string][]*entities.Wei{
			"mqbKtarYKnoEdPheFFDGRjksvEpb2vJGNh": {entities.NewWei(0.005 * 1e18)},
			"mowfvQDraTDvRgZowL4tx5EatL1u78w65v": {entities.NewWei(0.01956600 * 1e18)},
			"":                                   {entities.NewWei(0)}, // Null data script output
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
