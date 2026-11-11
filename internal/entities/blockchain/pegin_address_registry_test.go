package blockchain_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsSupportedPegInEncoding(t *testing.T) {
	assert.True(t, blockchain.IsSupportedPegInEncoding(blockchain.PegInAddressRegistryEncodingBase58))
	assert.False(t, blockchain.IsSupportedPegInEncoding(blockchain.PegInAddressRegistryEncodingBech32))
	assert.False(t, blockchain.IsSupportedPegInEncoding(blockchain.PegInAddressRegistryEncodingBech32M))
}

func TestNewAddressRegisteredFromWatchEntry(t *testing.T) {
	watch := rootstock.NewPegInWatch("0xhash", 4, 99, "0xrsk", "0xregistrant", [32]byte{7})

	event := blockchain.NewAddressRegisteredFromWatchEntry(watch)

	assert.Equal(t, blockchain.AddressRegistered{
		RskAddress:       "0xrsk",
		Registrant:       "0xregistrant",
		RegistrationRoot: [32]byte{7},
		TxHash:           "0xhash",
		BlockNumber:      99,
		LogIndex:         4,
	}, event)
}

func TestFoldPegInAddressRegistryRoot(t *testing.T) {
	t.Run("matches the contract fold vectors", func(t *testing.T) {
		expected := []struct {
			address string
			root    string
		}{
			{address: "0x00000000000000000000000000000000000000a1", root: "0x5a16856e66cb2b1b463f7773c427085d55afdd19d778290b45fb959a6224877e"},
			{address: "0x00000000000000000000000000000000000000a2", root: "0x95bdcbc864b6248af173e5feb03f5830479e972a492e705c3c2a2bddcb8ca643"},
			{address: "0x00000000000000000000000000000000000000a3", root: "0xb80604f49f0685bb17dd0f5cc0f611383d724f10f53db8aebcec3e9541f552d8"},
		}
		var root [32]byte
		for _, step := range expected {
			var err error
			root, err = blockchain.FoldPegInAddressRegistryRoot(crypto.Keccak256, root, step.address)
			require.NoError(t, err)
			assert.Equal(t, common.HexToHash(step.root), common.Hash(root))
		}
	})
	t.Run("depends on the registration order", func(t *testing.T) {
		first := "0x000000000000000000000000000000000000aaaa"
		second := "0x000000000000000000000000000000000000bbbb"
		firstThenSecond := foldAll(t, first, second)
		secondThenFirst := foldAll(t, second, first)
		assert.Equal(t, common.HexToHash("0x8cb7b8d9bfe92b5339d39bd1c1a03f82b52873855086b838b494eea818f66f2b"), common.Hash(firstThenSecond))
		assert.Equal(t, common.HexToHash("0x548d96cb839779d42379c1950321e7861349e8651f7e2f10ba06459c1c1ad4fa"), common.Hash(secondThenFirst))
	})
	t.Run("ignores address case and surrounding spaces", func(t *testing.T) {
		lower, err := blockchain.FoldPegInAddressRegistryRoot(crypto.Keccak256, [32]byte{}, "0x00000000000000000000000000000000000000ab")
		require.NoError(t, err)
		mixed, err := blockchain.FoldPegInAddressRegistryRoot(crypto.Keccak256, [32]byte{}, " 0x00000000000000000000000000000000000000AB ")
		require.NoError(t, err)
		assert.Equal(t, lower, mixed)
	})
	t.Run("rejects an invalid address", func(t *testing.T) {
		root, err := blockchain.FoldPegInAddressRegistryRoot(crypto.Keccak256, [32]byte{}, "0x1234")
		require.ErrorIs(t, err, blockchain.InvalidAddressError)
		assert.Equal(t, [32]byte{}, root)
	})
	t.Run("rejects a nil hash function", func(t *testing.T) {
		root, err := blockchain.FoldPegInAddressRegistryRoot(nil, [32]byte{}, "0x00000000000000000000000000000000000000a1")
		require.ErrorIs(t, err, blockchain.InvalidRootHashFunctionError)
		assert.Equal(t, [32]byte{}, root)
	})
	t.Run("rejects a hash function that does not return 32 bytes", func(t *testing.T) {
		shortHash := func(...[]byte) []byte { return []byte{1, 2, 3} }
		root, err := blockchain.FoldPegInAddressRegistryRoot(shortHash, [32]byte{}, "0x00000000000000000000000000000000000000a1")
		require.ErrorIs(t, err, blockchain.InvalidRootHashFunctionError)
		assert.Equal(t, [32]byte{}, root)
	})
}

func foldAll(t *testing.T, addresses ...string) [32]byte {
	t.Helper()
	var root [32]byte
	for _, address := range addresses {
		var err error
		root, err = blockchain.FoldPegInAddressRegistryRoot(crypto.Keccak256, root, address)
		require.NoError(t, err)
	}
	return root
}
