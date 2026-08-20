package blockchain_test

import (
	"encoding/hex"
	"reflect"
	"testing"

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

func mustRoot(t *testing.T, value string) [32]byte {
	t.Helper()
	decoded, err := hex.DecodeString(value)
	require.NoError(t, err)
	require.Len(t, decoded, 32)
	return [32]byte(decoded)
}

func TestRootFoldMatchesPinnedContractVectors(t *testing.T) {
	addresses := []string{
		"0x00000000000000000000000000000000000000a1",
		"0x00000000000000000000000000000000000000a2",
		"0x00000000000000000000000000000000000000a3",
	}
	expectedRoots := [][32]byte{
		mustRoot(t, "5a16856e66cb2b1b463f7773c427085d55afdd19d778290b45fb959a6224877e"),
		mustRoot(t, "95bdcbc864b6248af173e5feb03f5830479e972a492e705c3c2a2bddcb8ca643"),
		mustRoot(t, "b80604f49f0685bb17dd0f5cc0f611383d724f10f53db8aebcec3e9541f552d8"),
	}

	t.Run("packs exactly 32 root bytes followed by 20 address bytes", func(t *testing.T) {
		addressBytes, err := hex.DecodeString(addresses[1][2:])
		require.NoError(t, err)
		packed := append(append([]byte{}, expectedRoots[0][:]...), addressBytes...)
		require.Len(t, expectedRoots[0], 32)
		require.Len(t, addressBytes, 20)
		require.Len(t, packed, 52)
		assert.Equal(t, expectedRoots[0][:], packed[:32])
		assert.Equal(t, addressBytes, packed[32:])
	})

	t.Run("matches every pinned root after each ordered fold", func(t *testing.T) {
		var root [32]byte
		for index, address := range addresses {
			var err error
			root, err = blockchain.FoldPegInAddressRegistryRoot(root, address)
			require.NoError(t, err)
			assert.Equal(t, expectedRoots[index], root)
		}
	})

	t.Run("depends on registration order", func(t *testing.T) {
		var reversedRoot [32]byte
		for index := len(addresses) - 1; index >= 0; index-- {
			var err error
			reversedRoot, err = blockchain.FoldPegInAddressRegistryRoot(reversedRoot, addresses[index])
			require.NoError(t, err)
		}
		assert.NotEqual(t, expectedRoots[len(expectedRoots)-1], reversedRoot)
	})

	t.Run("rejects malformed address formats", func(t *testing.T) {
		for _, address := range []string{
			"00000000000000000000000000000000000000a1",
			"0x000000000000000000000000000000000000a1",
			"0x0000000000000000000000000000000000000000a1",
			"0x00000000000000000000000000000000000000zz",
		} {
			t.Run(address, func(t *testing.T) {
				_, err := blockchain.FoldPegInAddressRegistryRoot([32]byte{}, address)
				require.ErrorIs(t, err, blockchain.InvalidAddressError)
			})
		}
	})
}

// PegInAddressRegistryContract must be a read-only port: registerAddress belongs to a
// separate on-chain watcher process, not the liquidity provider server, so it must never
// appear on this interface.
// nolint:funlen
func TestPegInAddressRegistryContract_MethodSet(t *testing.T) {
	contractType := reflect.TypeFor[blockchain.PegInAddressRegistryContract]()
	expectedMethods := []string{
		"GetAddress",
		"GetPegInAddress",
		"GetPegInAddresses",
		"IsRegistered",
		"GetRegistration",
		"GetRegistrationRoot",
		"GetAddressRegisteredEvents",
	}

	assert.Equal(t, len(expectedMethods), contractType.NumMethod(), "PegInAddressRegistryContract must expose exactly its intended read-only surface")

	actualMethods := make([]string, contractType.NumMethod())
	for i := 0; i < contractType.NumMethod(); i++ {
		actualMethods[i] = contractType.Method(i).Name
	}
	assert.ElementsMatch(t, expectedMethods, actualMethods)

	disallowedWriteMethods := []string{"RegisterAddress", "Register", "SetRegistration", "Write"}
	for _, disallowed := range disallowedWriteMethods {
		_, found := contractType.MethodByName(disallowed)
		assert.False(t, found, "PegInAddressRegistryContract must not expose write method %q", disallowed)
	}
}
