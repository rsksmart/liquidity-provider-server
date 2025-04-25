package bitcoin_test

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcutil/bloom"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"slices"
	"testing"
)

func TestNewMerkleBlock_Witness(t *testing.T) {
	t.Run("Should build the merkle block for testnet", func(t *testing.T) {
		txHash, err := chainhash.NewHashFromStr("06ad9b57a21d72c11d9aaf45fdbe61a2f8e0dae50712fe46ae8ebd93b5b7c91c")
		require.NoError(t, err)
		block := getTestBlock(t, testnetBlockFile)
		filter := bloom.NewFilter(1, 0, 0, wire.BloomUpdateAll)
		filter.AddHash(txHash)
		testnetMsg, testnetIndices := bitcoin.NewMerkleBlock(block, filter, true)
		assert.Equal(t, []uint32{3}, testnetIndices)
		assert.Equal(t, 1030, int(testnetMsg.Transactions))
		expectedHashes := []string{
			"0x25a831485c203f24db116e436acfdfa9d1f8945d34368e4ce234647646e8a298", "0x16a6aba7aba7d13d659842d19c6933750345a0187335736f77ff7cca3f1dc714", "0x4492ab35d1497292731f5a1df4198c517d347f75ca61ba68240840eab3fb1e1f",
			"0x8ee4b9bbd387c878f6bee9a2c37d9fe315e7bb5928349b79f88b088585318e55", "0xa8b7229a3da991db8845997d2c4b08ae3a4c4a2752eb8631e5889cd5b94bd7aa", "0xa7b662ef0c7ae8523e8839c6af00d3df52572b7aebe2189a510ff14e794c6052",
			"0xd375df3986f436ecbb3d7562027c59cdbd16be1dc5c2446a67607f267fe29b5b", "0xc8e6d4e1cb207cdb50cc9c3b7e230b221332415d3637033cb225014588dddcef", "0xbd0d5fa11ef62876ddbfd672c58e9e15c7715599858599901f1d45a7199b2651",
			"0x1ea28b8b9fd7c75bed9b0f2794e2bf75008db935e7ac389aee551de355fc0cff", "0x48018e83e9af4ae4f54559b3afe19428639c87e4a5107e93bf7ae0c7bf9ba640", "0x98a4f1165e6aeb9b73967ecdce2ee8c62cd25d553affcda476678bf5c8370ef9",
		}
		for i, v := range toChainHashSlice(expectedHashes) {
			assert.Equal(t, v, testnetMsg.Hashes[i])
		}
	})
	t.Run("Should build the merkle block for mainnet", func(t *testing.T) {
		txHash, err := chainhash.NewHashFromStr("7c39408eeda72542b182ddb4bc737f2f4a7cff9924a14d0426796e64df850b81")
		require.NoError(t, err)
		block := getTestBlock(t, mainnetBlockFile)
		filter := bloom.NewFilter(1, 0, 0, wire.BloomUpdateAll)
		filter.AddHash(txHash)
		mainnetMsg, mainnetIndices := bitcoin.NewMerkleBlock(block, filter, true)
		assert.Equal(t, []uint32{49}, mainnetIndices)
		assert.Equal(t, 2291, int(mainnetMsg.Transactions))
		expectedHashes := []string{
			"0x6dd06bb078072deef1872d3400e2cf46e7666ee7fe613560a1e06db84afcd44b", "0xa53dec627778de55c90ae467cec5c6c030d2aaa058f10ac542f59a8cd5352dd5", "0xe9d9c7aaf8b36f37a1b3e9a2174c1f88cc861e7427b2585f208ddd567fa04611",
			"0x09dd3d95c3caab691cc02ff6291a7fbc0c82e173e136cba88879892608ae8512", "0x2fbce604c42091494df59dfee70bbedc66e608cdf3c69d51b42ac5b0b17eaeca", "0x0b2ffb66037ffcc030250063e9ab23447b175ebe9df908cab5db2498b01f2cfb",
			"0x7bae18cb060c92965f060ee782ed9a272bd0d5d6f85f206f2998d1db7c7e3edc", "0xac65e93d925e1c2bccda7815ef22a99091ea873891a5a5db0754f973bf5f9cb7", "0xf92d171b8f3a83c30acd37abe19a7c58cacf323f8f73bd49edfbe7f1a7b88ac5",
			"0xcb0d16bb2dfce77bbd2089eea7088e1eea8242c9a0dd186d582ef2731abe7f08", "0x70f3ecd79b56311f70c6283c479af600ed23b85a94791f42f6c3c7bd8fee397a", "0xaf49942453445fd8bf64668b0401d3792a48984f5cf644b374d39e2a7daee82f",
			"0x5713a9308504cec9d4c0e06dee7cbe2a46bc558d42a569f2fcf68321ec47666b",
		}
		for i, v := range toChainHashSlice(expectedHashes) {
			assert.Equal(t, v, mainnetMsg.Hashes[i])
		}
	})
}

func TestNewMerkleBlock_NoWitness(t *testing.T) {
	t.Run("Should build the merkle block for testnet", func(t *testing.T) {
		txHash, err := chainhash.NewHashFromStr("06ad9b57a21d72c11d9aaf45fdbe61a2f8e0dae50712fe46ae8ebd93b5b7c91c")
		require.NoError(t, err)
		block := getTestBlock(t, testnetBlockFile)
		filter := bloom.NewFilter(1, 0, 0, wire.BloomUpdateAll)
		filter.AddHash(txHash)
		testnetMsg, testnetIndices := bitcoin.NewMerkleBlock(block, filter, false)
		assert.Equal(t, []uint32{3}, testnetIndices)
		assert.Equal(t, 1030, int(testnetMsg.Transactions))
		expectedHashes := []string{
			"0xca4d886dc09a86ef7a040e9e20b3bc885e4295f63951dffcc57a1a626ce24465", "0x1ac616516a8d94666293d5b36342c06be60d8df437e32278a22e02076a1b3eba", "0x06ad9b57a21d72c11d9aaf45fdbe61a2f8e0dae50712fe46ae8ebd93b5b7c91c",
			"0x80f69cc228447b204b31dc141cf2e5a6bd256184f384b10a265ca8011aa4ca48", "0xd109aa5364dbdbe72a53a10a4a6a30b69dc904c3a66fc7e80d02aa3797057779", "0x53f4c10d1713edd245c972f3f595b3362874cf41b944db7a376d22a06e5e833c",
			"0x825205d9f5f690fcc2e004ca05c2a9fe8e7deb8c96c7a389a2b6efddfa702b70", "0x41e88fddc8c7adf2f67f97b536ab7fa7831d0aff34257b41f6bed0bde1fbf9ca", "0xb619444deb3cb58ba9368aebbc26573bd8043247a147f63a7c0eb4a82ed14141",
			"0x456f0dc25dd8440165f0839822094a79fb6501cd737f4c2143db53b7225f7d02", "0x0d8b34db87e8b3916fdfe388c90c93f91e222980903ed639fcc4e580886253b7", "0x16559e3831c4186ae16d8fa46ac164bcab51e746a0030793e2503b7202fe898a",
		}
		for i, v := range toChainHashSlice(expectedHashes) {
			assert.Equal(t, v, testnetMsg.Hashes[i])
		}
	})
	t.Run("Should build the merkle block for mainnet", func(t *testing.T) {
		txHash, err := chainhash.NewHashFromStr("7c39408eeda72542b182ddb4bc737f2f4a7cff9924a14d0426796e64df850b81")
		require.NoError(t, err)
		block := getTestBlock(t, mainnetBlockFile)
		filter := bloom.NewFilter(1, 0, 0, wire.BloomUpdateAll)
		filter.AddHash(txHash)
		mainnetMsg, mainnetIndices := bitcoin.NewMerkleBlock(block, filter, false)
		assert.Equal(t, []uint32{49}, mainnetIndices)
		assert.Equal(t, 2291, int(mainnetMsg.Transactions))
		expectedHashes := []string{
			"0xcf432ddee041d48a985c3ca61de7509d8e05cb7bab29a60541343f5fb2927c4d", "0x74e3d16f82c17105f0f7ea1b3257fec1ac6ecfcb04cc0b05059dbafb6aebd721", "0xe9d9c7aaf8b36f37a1b3e9a2174c1f88cc861e7427b2585f208ddd567fa04611",
			"0x7c39408eeda72542b182ddb4bc737f2f4a7cff9924a14d0426796e64df850b81", "0x237f04ca70798d96e6102dae9c7afe870ce77704aaf2b4cc9759012828dde988", "0x202786f4974722189e29c57e307f208666939a0d715c780d8c1dea733dc3e901",
			"0xfd87109d9231c7ceb9abff4679313c226b7cf603a12e637a64b105941af6d97a", "0x55744ecf26d8c689ffd0b03a595ee17f2a5254ba4aa37c3806683f5bc5212de2", "0xab2a399d4616a0833d641340bda5887390b258313656f066a3f14e6ad6febd99",
			"0xe4fb3219f8cf65e3c5b4e38c2872f66328b01b1f663c79c9d201c80a0a911923", "0x7897ddd3bb6f7b7713d2daa6c563ee02761677c6b8b50a24b1807f024bba1e3d", "0x3d40b5721801cca5181bab9e54807f5c1fd1a45e983525d66c8f1a5d970bf371",
			"0x0c066d440dff37587ee94a97ceca8d6d799cec0b5c4c60d1643c1a83026366ee",
		}
		for i, v := range toChainHashSlice(expectedHashes) {
			assert.Equal(t, v, mainnetMsg.Hashes[i])
		}
	})
}

func toChainHashSlice(input []string) []*chainhash.Hash {
	result := make([]*chainhash.Hash, 0)
	for _, hash := range input {
		hashBytes, _ := hex.DecodeString(hash[2:]) // nolint:errcheck
		slices.Reverse(hashBytes)
		parsedBytes, _ := chainhash.NewHash(hashBytes) // nolint:errcheck
		result = append(result, parsedBytes)
	}
	return result
}
