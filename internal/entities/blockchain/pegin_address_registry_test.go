package blockchain_test

import (
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/stretchr/testify/assert"
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
