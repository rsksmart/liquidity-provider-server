package bitcoin_test

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/datasets"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDecodeAddressBase58(t *testing.T) {
	var decodedAddresses []datasets.DecodedAddress
	decodedAddresses = append(decodedAddresses, datasets.Base58Addresses...)
	cases := decodedAddresses
	for _, c := range cases {
		withVersion, err := bitcoin.DecodeAddressBase58(c.Address, true)
		require.NoError(t, err)
		withoutVersion, err := bitcoin.DecodeAddressBase58(c.Address, false)
		require.NoError(t, err)
		assert.Equal(t, c.Expected, withVersion)
		assert.Equal(t, c.Expected[1:], withoutVersion)
	}
}

func TestDecodeAddressBase58_ErrorHandling(t *testing.T) {
	var errorCases []string = []string{
		"5Hwgr3u458GLafKBgxtssHSPqJnYoGrSzgQsPwLFhLNYskDPyyA",
		"not in bas58",
		"A",
		"0x79568c2989232dCa1840087D73d403602364c0D4",
	}
	var result []byte
	var err error
	for _, c := range errorCases {
		result, err = bitcoin.DecodeAddressBase58(c, true)
		require.Error(t, err)
		assert.Nil(t, result)
		result, err = bitcoin.DecodeAddressBase58(c, false)
		require.Error(t, err)
		assert.Nil(t, result)
	}
}

func TestDecodeAddress(t *testing.T) {
	var decodedAddresses []datasets.DecodedAddress
	decodedAddresses = append(decodedAddresses, datasets.Base58Addresses...)
	decodedAddresses = append(decodedAddresses, datasets.Bech32Addresses...)
	decodedAddresses = append(decodedAddresses, datasets.Bech32mAddresses...)
	cases := decodedAddresses
	for _, c := range cases {
		decoded, err := bitcoin.DecodeAddress(c.Address)
		require.NoError(t, err)
		assert.Equal(t, c.Expected, decoded)
	}
}

func TestToSwappedBytes32(t *testing.T) {
	var bytes32 = [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20}
	var hash = chainhash.Hash([32]byte{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F, 0x40})
	var hashPointer, err = chainhash.NewHash([]byte{0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F, 0x60})

	require.NoError(t, err)
	swappedBytes := bitcoin.ToSwappedBytes32(bytes32)
	swappedHash := bitcoin.ToSwappedBytes32(hash)
	swappedHashPointer := bitcoin.ToSwappedBytes32(hashPointer)

	assert.Equal(t, [32]byte{0x20, 0x1F, 0x1E, 0x1D, 0x1C, 0x1B, 0x1A, 0x19, 0x18, 0x17, 0x16, 0x15, 0x14, 0x13, 0x12, 0x11, 0x10, 0x0F, 0x0E, 0x0D, 0x0C, 0x0B, 0x0A, 0x09, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01}, swappedBytes)
	assert.Equal(t, [32]byte{0x40, 0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38, 0x37, 0x36, 0x35, 0x34, 0x33, 0x32, 0x31, 0x30, 0x2F, 0x2E, 0x2D, 0x2C, 0x2B, 0x2A, 0x29, 0x28, 0x27, 0x26, 0x25, 0x24, 0x23, 0x22, 0x21}, swappedHash)
	assert.Equal(t, [32]byte{0x60, 0x5F, 0x5E, 0x5D, 0x5C, 0x5B, 0x5A, 0x59, 0x58, 0x57, 0x56, 0x55, 0x54, 0x53, 0x52, 0x51, 0x50, 0x4F, 0x4E, 0x4D, 0x4C, 0x4B, 0x4A, 0x49, 0x48, 0x47, 0x46, 0x45, 0x44, 0x43, 0x42, 0x41}, swappedHashPointer)
}

func TestEnsureLoadedBtcWallet(t *testing.T) {
	t.Run("Should return error if connection is not a wallet connection", func(t *testing.T) {
		conn := bitcoin.NewConnection(&chaincfg.TestNet3Params, new(mocks.ClientAdapterMock))
		err := bitcoin.EnsureLoadedBtcWallet(conn)
		require.ErrorContains(t, err, "connection is not a wallet connection")
	})
	t.Run("Shouldn't return error if wallet is loaded and responding", func(t *testing.T) {
		clientMock := new(mocks.ClientAdapterMock)
		clientMock.EXPECT().GetWalletInfo().Return(&btcjson.GetWalletInfoResult{WalletName: test.AnyString}, nil)
		conn := bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, clientMock, test.AnyString)
		err := bitcoin.EnsureLoadedBtcWallet(conn)
		require.NoError(t, err)
		clientMock.AssertExpectations(t)
	})
	t.Run("Should load wallet if wallet is not loaded", func(t *testing.T) {
		clientMock := new(mocks.ClientAdapterMock)
		clientMock.EXPECT().GetWalletInfo().Return(nil, assert.AnError)
		clientMock.EXPECT().LoadWallet(test.AnyString).Return(&btcjson.LoadWalletResult{Name: test.AnyString}, nil)
		conn := bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, clientMock, test.AnyString)
		err := bitcoin.EnsureLoadedBtcWallet(conn)
		require.NoError(t, err)
		clientMock.AssertExpectations(t)
	})
	t.Run("Should return error loading wallet if it fails", func(t *testing.T) {
		clientMock := new(mocks.ClientAdapterMock)
		clientMock.EXPECT().GetWalletInfo().Return(nil, assert.AnError)
		clientMock.EXPECT().LoadWallet(test.AnyString).Return(nil, assert.AnError)
		conn := bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, clientMock, test.AnyString)
		err := bitcoin.EnsureLoadedBtcWallet(conn)
		require.Error(t, err)
		clientMock.AssertExpectations(t)
	})
}
