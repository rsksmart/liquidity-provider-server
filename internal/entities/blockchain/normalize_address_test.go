package blockchain_test

import (
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeEthereumAddress(t *testing.T) {
	const wantLower = "0xd839c223634b224327430bb7062858109c850bf9"
	t.Run("lowercase input unchanged", func(t *testing.T) {
		got, err := blockchain.NormalizeEthereumAddress(wantLower)
		require.NoError(t, err)
		assert.Equal(t, wantLower, got)
	})
	t.Run("EIP-55 input becomes lowercase", func(t *testing.T) {
		got, err := blockchain.NormalizeEthereumAddress("0xD839C223634b224327430Bb7062858109C850bf9")
		require.NoError(t, err)
		assert.Equal(t, wantLower, got)
	})
	t.Run("uppercase hex becomes lowercase", func(t *testing.T) {
		got, err := blockchain.NormalizeEthereumAddress("0xD839C223634B224327430BB7062858109C850BF9")
		require.NoError(t, err)
		assert.Equal(t, wantLower, got)
	})
	t.Run("invalid address rejected", func(t *testing.T) {
		_, err := blockchain.NormalizeEthereumAddress("not-an-address")
		require.Error(t, err)
	})
	t.Run("too short hex rejected", func(t *testing.T) {
		_, err := blockchain.NormalizeEthereumAddress("0x1234")
		require.Error(t, err)
	})
}
