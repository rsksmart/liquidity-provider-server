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
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/datasets"
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
	testFilePath = "../../../../test/mocks/"
)

const (
	testnetTestBlockHash     = "00000000001e94d85c3e736aa4071d36d26547713820a27af9edbe97489c696f"
	testnetTestTxHash        = "9f0706c2717fc77bf0f225a4223933a7decb8d36902ddbb0accab8ea894f8b29"
	testnetWitnessTestTxHash = "5cadcbc1ccd91b222346f22c9a9a6fdbf20c9338ec8df0b36097e92d029509ec"
	testnetBlockFile         = "block-2582756-testnet.txt"
)

const (
	mainnetTestBlockHash     = "0000000000000000000aca0460feaf0661f173b75d4cc824b57233aa7c6b7bc3"
	mainnetTestTxHash        = "e28bec3d29efce36405197d1255cfebde06ba9c193d8192d3825d6e9213b03ed"
	mainnetWitnessTestTxHash = "85c2fc50c70ceda8cb9f62aacc65a67b76411e442096d86649c95d7e9a28af8c"
	mainnetBlockFile         = "block-696394-mainnet.txt"
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
	p2wpkhTestnet := "tb1qg9stx3w3xj0t5j8d5q2g4yvz5gj8v3tjxjxw5v" // TODO add checksum validation to bitcoind implementation
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
	var decodedAddresses []datasets.DecodedAddress
	decodedAddresses = append(decodedAddresses, datasets.Base58Addresses...)
	decodedAddresses = append(decodedAddresses, datasets.Bech32Addresses...)
	decodedAddresses = append(decodedAddresses, datasets.Bech32mAddresses...)
	client := &mocks.ClientAdapterMock{}
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	cases := decodedAddresses
	for _, c := range cases {
		decoded, err := rpc.DecodeAddress(c.Address)
		require.NoError(t, err)
		assert.Equal(t, c.Expected, decoded)
	}
}

func TestBitcoindRpc_GetRawTransaction(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	txBytes, err := hex.DecodeString("0200000002ebf7c22a73f3baea460cad53a2788bd4f24020f6b374900a771d3422f128442e000000006a473044022062dae13ba281d0cf529b604bb59c1efcd7b83438af34d4a51acc6f31041be18c022044df281e688a52624f45f6c26662349d1f5efedd4d69530e65b7d7cec0d3792d0121038e509bc056004a5da7460b5acd5d4dcb2add41d53817180499e3814290ecc91efdffffffb5f09f38215b850f4ba644a7f7ab57efa8d10c5f4b5908e9aa980ff5ffa948f5000000006a47304402206538fc72b896e4c6e807a4daf56191f68dec307c3011d082e69eeb3d45d6d8c302205a329814ab87901ae56a82587e716fa2282ecc665ab203da14d93db71181ecd8012102498a833095175800f40b2c0ab23f108b47a319a94ccea826062bf66c827e91a9fdffffff0298740700000000001976a91473cce22e78ec61cd54a6438ca1210b88561ebcdd88ac20a10700000000001976a9142c81478132b5dda64ffc484a0d225096c4b22ad588acc3682700")
	require.NoError(t, err)
	client.On("GetRawTransaction", mock.Anything).Return(btcutil.NewTxFromBytes(txBytes)).Once()
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
	result, err := rpc.GetRawTransaction(testnetTestTxHash)
	require.NoError(t, err)
	require.Equal(t, txBytes, result)
}

func TestBitcoindRpc_GetRawTransaction_FromBlock(t *testing.T) {
	cases := []datasets.RawTransaction{
		datasets.BtcCoinbaseTxNoWitness,
		datasets.BtcSegwitTxNoWitness,
		datasets.BtcLegacyTxNoWitness,
	}
	mainnetBlock := test.GetBitcoinTestBlock(t, testFilePath+mainnetBlockFile)
	client := &mocks.ClientAdapterMock{}
	for _, tx := range cases {
		parsedTx, err := mainnetBlock.Tx(tx.Index)
		require.NoError(t, err)
		client.On("GetRawTransaction", mock.Anything).Return(parsedTx, nil).Once()
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
		result, err := rpc.GetRawTransaction(parsedTx.Hash().String())
		require.NoError(t, err)
		expectedBytes, err := hex.DecodeString(tx.Tx)
		require.NoError(t, err)
		require.Equal(t, expectedBytes, result)
	}
}

func TestBitcoindRpc_GetRawTransaction_ErrorHandling(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	txBytes, err := hex.DecodeString("0200000002ebf7c22a73f3baea460cad53a2788bd4f24020f6b374900a771d3422f128442e000000006a473044022062dae13ba281d0cf529b604bb59c1efcd7b83438af34d4a51acc6f31041be18c022044df281e688a52624f45f6c26662349d1f5efedd4d69530e65b7d7cec0d3792d0121038e509bc056004a5da7460b5acd5d4dcb2add41d53817180499e3814290ecc91efdffffffb5f09f38215b850f4ba644a7f7ab57efa8d10c5f4b5908e9aa980ff5ffa948f5000000006a47304402206538fc72b896e4c6e807a4daf56191f68dec307c3011d082e69eeb3d45d6d8c302205a329814ab87901ae56a82587e716fa2282ecc665ab203da14d93db71181ecd8012102498a833095175800f40b2c0ab23f108b47a319a94ccea826062bf66c827e91a9fdffffff0298740700000000001976a91473cce22e78ec61cd54a6438ca1210b88561ebcdd88ac20a10700000000001976a9142c81478132b5dda64ffc484a0d225096c4b22ad588acc3682700")
	require.NoError(t, err)
	client.On("GetRawTransaction", mock.Anything).Return(btcutil.NewTxFromBytes(txBytes)).Once()
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))

	_, err = rpc.GetRawTransaction("invalidHash")
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

// nolint:funlen
func TestBitcoindRpc_BuildMerkleBranch(t *testing.T) {
	testnetBlock := test.GetBitcoinTestBlock(t, testFilePath+testnetBlockFile)
	mainnetBlock := test.GetBitcoinTestBlock(t, testFilePath+mainnetBlockFile)
	t.Run("Should build merkle branch for testnet transactions", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: testnetTestBlockHash}, nil).Twice()
		client.On("GetBlock", mock.Anything).Return(testnetBlock.MsgBlock(), nil).Twice()
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.TestNet3Params, client))
		legacyBranch, legacyErr := rpc.BuildMerkleBranch(datasets.TestnetLegacyMerkleBranch.Tx)
		require.NoError(t, legacyErr)
		assert.Equal(t, datasets.TestnetLegacyMerkleBranch.Branch, legacyBranch)
		witnessBranch, witnessErr := rpc.BuildMerkleBranch(datasets.TestnetWitnessMerkleBranch.Tx)
		require.NoError(t, witnessErr)
		assert.Equal(t, datasets.TestnetWitnessMerkleBranch.Branch, witnessBranch)
		client.AssertExpectations(t)
	})
	t.Run("Should build merkle branch for mainnet transactions", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: mainnetTestBlockHash}, nil).Twice()
		client.On("GetBlock", mock.Anything).Return(mainnetBlock.MsgBlock(), nil).Twice()
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
		legacyBranch, legacyErr := rpc.BuildMerkleBranch(datasets.MainnetLegacyMerkleBranch.Tx)
		require.NoError(t, legacyErr)
		assert.Equal(t, datasets.MainnetLegacyMerkleBranch.Branch, legacyBranch)
		witnessBranch, witnessErr := rpc.BuildMerkleBranch(datasets.MainnetWitnessMerkleBranch.Tx)
		require.NoError(t, witnessErr)
		assert.Equal(t, datasets.MainnetWitnessMerkleBranch.Branch, witnessBranch)
		client.AssertExpectations(t)
	})
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
	block := test.GetBitcoinTestBlock(t, testFilePath+testnetBlockFile)

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
	block := test.GetBitcoinTestBlock(t, testFilePath+testnetBlockFile)
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
	pmt, err = rpc.GetPartialMerkleTree(testnetWitnessTestTxHash)
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
	client.AssertExpectations(t)
}

// nolint:funlen
func TestBitcoindRpc_GetPartialMerkleTree_MainnetBlock(t *testing.T) {
	cases := [3]struct {
		tx  string
		pmt string
	}{
		// first two are witness txs, last one is legacy
		{
			tx: "07f8b22fa9a3b32e20b59bb90727de05fb634749519ebcb6a887aeaf2c7eb041",
			pmt: "f30800000d" +
				"41b07e2cafae87a8b6bc9e51494763fb05de2707b99bb5202eb3a3a92fb2f80773" +
				"1c671fafb5d234834c726657f29c9af030ccf7068f1ef732af4efd8e146da0a9d6" +
				"075f4758821ceeef2c230cfd2497df2d1d1d02dd19e653d22b3dc271b39394c2d2" +
				"e51eae99af800c39575a11f3b4eb0fdf5d5deff0b9f5ff592566f4f1732b3f6d41" +
				"bded4344899a348d3053ccd68e922626e589da71d9a583ccfe9e3be6393a24e14d" +
				"18006b54a967963c56da58b18ce770cdb3b32e56d88c138c473a1a37acc29a7788" +
				"a88404fb9a05c416c2ad8f340a61c1ea331528a91ae6210db4f0e22d21c55b3f68" +
				"06387ca34aba54522a7fe15e593ab0d0ff89c6d826cf4e745599bdfed66a4ef1a3" +
				"66f056363158b2907388a5bd4013643d83a016469d392aab2319910a0ac801d2c9" +
				"793c661f1bb02863f672288ce3b4c5e365cff81932fbe43d1eba4b027f80b1240a" +
				"b5b8c677167602ee63c5a6dad213777b6fbbd3dd977871f30b975d1a8f6cd62535" +
				"985ea4d11f5c7f80549eab1b18a5cc011872b5403dee666302831a3c64d1604c5c" +
				"0bec9c796d8dcace974ae97e5837ff0d446d060c04ff1f0000",
		},
		{
			tx: "ddf5061f9707f0c959bf24278d557b264716672c1b601ec50112d6dfe160d9d3",
			pmt: "f30800000d" +
				"c0746a357444e9948a18a612e02df5a99240e77f1ff75dd949d5b4038dcf36673a" +
				"03c716cf722cff7d264c763088ceeb1665f26c6fdd5835d841eeee2f3ece4a203a" +
				"24db8b7a51e4ab0e35a6b4151f6d7f1eef96f32e4fceaac61275219116186efb7f" +
				"db763e821f99bd6af8d044cc6feadd7b4716e6938335a3e08548f5a0775dd36497" +
				"1faab5cd089cd1fa713e8be658a67a704d39952218f6518e5045d269d3d960e1df" +
				"d61201c51e601b2c671647267b558d2724bf59c9f007971f06f5dd0eab2677f52c" +
				"996a3f941bef3ec57ebdf22429c37dee5ae68892df30f8acfc225c6fed56bdff34" +
				"686135e68fda4b716713e60258b6971c03091f25115c008eec48a828c75ad7340f" +
				"adbc368636b4014f6e8386c3990a35620cbddca933a72b02d990fb8a602fcda9e1" +
				"e41120c25f4981362a9dfc7f7ed1f5188482b8ee3f532f0ee6234e44af99351ee4" +
				"30f4ac0fa7b71fe9c601c78480b9a97fea305d3abca235a6668846093803e07c48" +
				"dc9a75be90ed6edd4debb0b7b49bc057e093ad395eee666302831a3c64d1604c5c" +
				"0bec9c796d8dcace974ae97e5837ff0d446d060c04dbd50300",
		},
		{
			tx: "db0d1fe6384b5741ceb2e67f4b50372966e1bab2b50e91a597ca4170c5f281e9",
			pmt: "f30800000d" +
				"32f19bb610f0d51f754364ea4fda76ae2488b61fb4b2e0d966403c7c3544d20a" +
				"e981f2c57041ca97a5910eb5b2bae1662937504b7fe6b2ce41574b38e61f0ddb" +
				"48055668e4e3d31af8efdb4179740c9de740d57594efeee0535d572d9dbd5f95" +
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
				"04ff370000",
		},
	}
	mainnetBlock := test.GetBitcoinTestBlock(t, testFilePath+mainnetBlockFile)
	for _, c := range cases {
		client := &mocks.ClientAdapterMock{}
		client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: mainnetTestBlockHash}, nil).Once()
		client.On("GetBlock", mock.Anything).Return(mainnetBlock.MsgBlock(), nil).Once()
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
		serializedPMT, err := rpc.GetPartialMerkleTree(c.tx)
		require.NoError(t, err)
		result := hex.EncodeToString(serializedPMT)
		assert.Equal(t, c.pmt, result)
		client.AssertExpectations(t)
	}
}

func TestBitcoindRpc_GetPartialMerkleTree_ErrorHandling(t *testing.T) {
	rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, &mocks.ClientAdapterMock{}))
	pmt, err := rpc.GetPartialMerkleTree("txhash")
	require.Error(t, err)
	require.Nil(t, pmt)

	client := &mocks.ClientAdapterMock{}
	block := test.GetBitcoinTestBlock(t, testFilePath+testnetBlockFile)
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
		coinbaseInfo, err := rpc.GetCoinbaseInformation(testnetWitnessTestTxHash)
		assert.Empty(t, coinbaseInfo)
		require.Error(t, err)
	})
	t.Run("Should handle error getting block", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: testnetTestBlockHash}, nil).Once()
		client.On("GetBlock", mock.Anything).Return(nil, assert.AnError).Once()
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
		coinbaseInfo, err := rpc.GetCoinbaseInformation(testnetWitnessTestTxHash)
		assert.Empty(t, coinbaseInfo)
		require.Error(t, err)
	})
	t.Run("Should handle error getting block verbose", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		testnetBlock := test.GetBitcoinTestBlock(t, testFilePath+testnetBlockFile)
		client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: testnetTestBlockHash}, nil).Once()
		client.On("GetBlock", mock.Anything).Return(testnetBlock.MsgBlock(), nil).Once()
		client.On("GetBlockVerbose", mock.Anything).Return(nil, assert.AnError).Once()
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
		coinbaseInfo, err := rpc.GetCoinbaseInformation(testnetWitnessTestTxHash)
		assert.Empty(t, coinbaseInfo)
		require.Error(t, err)
	})
	t.Run("Should build coinbase info", func(t *testing.T) {
		var blockHash, witnessMerkleRoot [32]byte
		testnetBlock := test.GetBitcoinTestBlock(t, testFilePath+testnetBlockFile)
		client := &mocks.ClientAdapterMock{}
		client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{BlockHash: testnetTestBlockHash}, nil).Once()
		client.On("GetBlock", mock.Anything).Return(testnetBlock.MsgBlock(), nil).Once()
		client.On("GetBlockVerbose", mock.Anything).Return(&btcjson.GetBlockVerboseResult{Height: 2582756}, nil).Once()
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
		coinbaseInfo, err := rpc.GetCoinbaseInformation(testnetWitnessTestTxHash)
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
			WitnessReservedValue: [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		}, coinbaseInfo)
	})
}

func TestBitcoindRpc_NetworkName(t *testing.T) {
	table := test.Table[*chaincfg.Params, string]{
		{Value: &chaincfg.MainNetParams, Result: "mainnet"},
		{Value: &chaincfg.TestNet3Params, Result: "testnet3"},
		{Value: &chaincfg.RegressionNetParams, Result: "regtest"},
		{Value: &chaincfg.SimNetParams, Result: "simnet"},
	}
	test.RunTable(t, table, func(p *chaincfg.Params) string {
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(p, &mocks.ClientAdapterMock{}))
		return rpc.NetworkName()
	})
}

func TestBitcoindRpc_GetBlockchainInfo(t *testing.T) {
	t.Run("Should return blockchain info", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.EXPECT().GetBlockChainInfo().Return(&btcjson.GetBlockChainInfoResult{
			Chain:         "mainnet",
			Blocks:        300,
			Headers:       350,
			BestBlockHash: test.AnyHash,
		}, nil).Once()
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
		info, err := rpc.GetBlockchainInfo()
		require.NoError(t, err)
		assert.Equal(t, blockchain.BitcoinBlockchainInfo{
			NetworkName:      "mainnet",
			ValidatedBlocks:  big.NewInt(300),
			ValidatedHeaders: big.NewInt(350),
			BestBlockHash:    test.AnyHash,
		}, info)
		client.AssertExpectations(t)
	})
	t.Run("Should return error when client fails", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.EXPECT().GetBlockChainInfo().Return(nil, assert.AnError).Once()
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, client))
		info, err := rpc.GetBlockchainInfo()
		require.Error(t, err)
		assert.Empty(t, info)
		client.AssertExpectations(t)
	})
}

func TestBitcoindRpc_GetZeroAddress(t *testing.T) {
	t.Run("should return zero address in mainnet", func(t *testing.T) {
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.MainNetParams, &mocks.ClientAdapterMock{}))
		cases := test.Table[blockchain.BtcAddressType, string]{
			{Value: blockchain.BtcAddressTypeP2PKH, Result: blockchain.BitcoinMainnetP2PKHZeroAddress},
			{Value: blockchain.BtcAddressTypeP2SH, Result: blockchain.BitcoinMainnetP2SHZeroAddress},
			{Value: blockchain.BtcAddressTypeP2WPKH, Result: blockchain.BitcoinMainnetP2WPKHZeroAddress},
			{Value: blockchain.BtcAddressTypeP2WSH, Result: blockchain.BitcoinMainnetP2WSHZeroAddress},
			{Value: blockchain.BtcAddressTypeP2TR, Result: blockchain.BitcoinMainnetP2TRZeroAddress},
		}
		test.RunTable(t, cases, func(addressType blockchain.BtcAddressType) string {
			result, err := rpc.GetZeroAddress(addressType)
			require.NoError(t, err)
			return result
		})
	})
	t.Run("should return zero address in testnet", func(t *testing.T) {
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.TestNet3Params, &mocks.ClientAdapterMock{}))
		cases := test.Table[blockchain.BtcAddressType, string]{
			{Value: blockchain.BtcAddressTypeP2PKH, Result: blockchain.BitcoinTestnetP2PKHZeroAddress},
			{Value: blockchain.BtcAddressTypeP2SH, Result: blockchain.BitcoinTestnetP2SHZeroAddress},
			{Value: blockchain.BtcAddressTypeP2WPKH, Result: blockchain.BitcoinTestnetP2WPKHZeroAddress},
			{Value: blockchain.BtcAddressTypeP2WSH, Result: blockchain.BitcoinTestnetP2WSHZeroAddress},
			{Value: blockchain.BtcAddressTypeP2TR, Result: blockchain.BitcoinTestnetP2TRZeroAddress},
		}
		test.RunTable(t, cases, func(addressType blockchain.BtcAddressType) string {
			result, err := rpc.GetZeroAddress(addressType)
			require.NoError(t, err)
			return result
		})
	})
	t.Run("should return zero address in regtest", func(t *testing.T) {
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.RegressionNetParams, &mocks.ClientAdapterMock{}))
		cases := test.Table[blockchain.BtcAddressType, string]{
			{Value: blockchain.BtcAddressTypeP2PKH, Result: blockchain.BitcoinTestnetP2PKHZeroAddress},
			{Value: blockchain.BtcAddressTypeP2SH, Result: blockchain.BitcoinTestnetP2SHZeroAddress},
			{Value: blockchain.BtcAddressTypeP2WPKH, Result: blockchain.BitcoinRegtestP2WPKHZeroAddress},
			{Value: blockchain.BtcAddressTypeP2WSH, Result: blockchain.BitcoinRegtestP2WSHZeroAddress},
			{Value: blockchain.BtcAddressTypeP2TR, Result: blockchain.BitcoinRegtestP2TRZeroAddress},
		}
		test.RunTable(t, cases, func(addressType blockchain.BtcAddressType) string {
			result, err := rpc.GetZeroAddress(addressType)
			require.NoError(t, err)
			return result
		})
	})
	t.Run("should return error if network not supported", func(t *testing.T) {
		rpc := bitcoin.NewBitcoindRpc(bitcoin.NewConnection(&chaincfg.SigNetParams, &mocks.ClientAdapterMock{}))
		result, err := rpc.GetZeroAddress(blockchain.BtcAddressTypeP2PKH)
		require.Error(t, err)
		assert.Empty(t, result)
	})
}
