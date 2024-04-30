package bitcoin_test

import (
	"context"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConnection_CheckConnection(t *testing.T) {
	networkParams := &chaincfg.Params{}
	client := &mocks.ClientAdapterMock{}
	client.On("Ping").Return(assert.AnError).Once()
	client.On("Ping").Return(nil).Once()
	conn := bitcoin.NewConnection(networkParams, client)
	conn.CheckConnection(context.Background())
	conn.CheckConnection(context.Background())
	client.AssertExpectations(t)
}

func TestConnection_Shutdown(t *testing.T) {
	endChannel := make(chan bool)
	client := &mocks.ClientAdapterMock{}
	client.On("Disconnect").Once()
	conn := bitcoin.NewConnection(&chaincfg.Params{}, client)
	go conn.Shutdown(endChannel)
	<-endChannel
	client.AssertExpectations(t)
}

var decodedAddresses = []struct {
	address  string
	expected []byte
}{
	{"n1BE7ioVukYS2GC88hT2K6cUvRiKwMwio7", []byte{111, 215, 167, 103, 99, 62, 208, 72, 131, 184, 122, 185, 112, 220, 93, 130, 94, 43, 74, 67, 67}},
	{"n2aSettzgmgwxMoaNQfNL58BqQtEAqCRp1", []byte{111, 231, 3, 144, 1, 201, 195, 77, 230, 3, 157, 153, 171, 189, 171, 215, 212, 180, 203, 182, 41}},
	{"mugStBVzw3ZRafzrYqzSrq5Qq9HH9cZJst", []byte{111, 155, 93, 94, 120, 235, 41, 106, 157, 161, 118, 38, 245, 188, 189, 62, 16, 242, 77, 50, 225}},
	{"mfng2KVQdxd1wmSEk5TRKVLwMarsX9sD7D", []byte{111, 2, 249, 18, 163, 98, 241, 57, 132, 12, 133, 190, 167, 51, 135, 39, 8, 238, 37, 238, 22}},
	{"miba3yyWiomi5HMXboKZ4YmEnj1T6m3RLx", []byte{111, 33, 199, 227, 213, 45, 69, 31, 199, 87, 161, 249, 199, 169, 25, 147, 105, 219, 96, 219, 46}},
	{"1Ab1Jfe6xQHzL8RHoHDukDQBEks35KFWHC", []byte{0, 105, 39, 135, 109, 109, 187, 191, 255, 241, 12, 238, 207, 1, 71, 119, 140, 83, 84, 176, 148}},
	{"19hiJTQpZyT3C7Hu29dJE2YYToCeKp6cGu", []byte{0, 95, 116, 30, 252, 161, 93, 178, 87, 245, 15, 180, 19, 235, 18, 247, 75, 185, 227, 86, 164}},
	{"1LQdpgVCY2nYzsoRNRHWhuCLxMpYzb6zzg", []byte{0, 212, 226, 177, 205, 89, 105, 145, 127, 135, 118, 204, 8, 84, 231, 158, 29, 254, 239, 62, 150}},
	{"17oKLsbZsd2BZdCDn1dbrbk2TT9HSzw2aM", []byte{0, 74, 147, 51, 90, 109, 233, 79, 47, 123, 200, 20, 172, 19, 242, 57, 107, 24, 145, 19, 190}},
	{"1NbJonAytRKfCFkvGcQNEUCXAFnf17bYQG", []byte{0, 236, 215, 165, 247, 66, 139, 143, 251, 27, 232, 167, 211, 106, 236, 65, 215, 144, 67, 50, 188}},
	{"1MwygkmvJHwwG934EbtkjhRUFyfMHLEPi9", []byte{0, 229, 200, 75, 213, 72, 142, 34, 190, 217, 45, 198, 157, 186, 72, 245, 47, 208, 252, 113, 119}},
	{"2MyKtiQcyAgQQdPDmBJcT4UMn6jyVdKVwxg", []byte{196, 66, 178, 190, 224, 218, 119, 248, 26, 82, 40, 48, 190, 49, 159, 125, 114, 243, 42, 175, 30}},
	{"2N9TujzoGXb4LkzCQscMuuRSwaTVPhFfYSy", []byte{196, 177, 232, 48, 84, 2, 170, 51, 48, 251, 57, 152, 125, 51, 230, 130, 144, 98, 27, 231, 156}},
	{"2NAvNk9t8fKwFRhMJGgY2wMLRboN8DHXeBB", []byte{196, 193, 225, 171, 152, 53, 142, 203, 39, 152, 199, 170, 199, 208, 233, 205, 65, 249, 239, 208, 91}},
	{"2MzqFJDg3QcbzeWX87XpxRpHZm6SSNoGdoF", []byte{196, 83, 56, 29, 183, 70, 201, 79, 21, 103, 129, 237, 73, 8, 214, 82, 67, 212, 124, 164, 213}},
	{"2MuBrh282woZkhycbpKAZ8zEptTAEtRSM62", []byte{196, 21, 77, 57, 173, 29, 183, 149, 178, 0, 34, 217, 55, 64, 184, 30, 90, 13, 91, 214, 241}},
	{"33BTakydSPnJfSfR13foniEsPCB2nuHiCb", []byte{5, 16, 89, 49, 17, 100, 144, 216, 162, 233, 218, 176, 118, 166, 174, 244, 82, 223, 195, 106, 191}},
	{"3Fi3ywSD7eBEtnoUuiW7zFSBmpATd3YDLs", []byte{5, 153, 195, 229, 51, 197, 179, 77, 218, 70, 180, 220, 57, 200, 231, 161, 90, 119, 51, 54, 31}},
	{"3DCT2YtzwZZYdr3pPEhqVjA1Amak89YrHf", []byte{5, 126, 58, 102, 228, 17, 153, 252, 64, 101, 122, 39, 140, 215, 39, 222, 50, 192, 152, 23, 248}},
	{"3EuQdJ651cN7Cv9jJk2EJPdRhKT9JJFpt8", []byte{5, 144, 241, 147, 242, 132, 35, 204, 174, 111, 172, 89, 117, 131, 75, 45, 4, 128, 176, 86, 42}},
	{"3BcLTSd24JRtJhLcKqkeF83rFmFxxY5qH9", []byte{5, 108, 206, 165, 252, 49, 214, 45, 206, 182, 122, 231, 101, 130, 27, 245, 51, 13, 14, 60, 174}},
}

func TestDecodeAddressBase58(t *testing.T) {
	cases := decodedAddresses
	for _, c := range cases {
		withVersion, err := bitcoin.DecodeAddressBase58(c.address, true)
		require.NoError(t, err)
		withoutVersion, err := bitcoin.DecodeAddressBase58(c.address, false)
		require.NoError(t, err)
		assert.Equal(t, c.expected, withVersion)
		assert.Equal(t, c.expected[1:], withoutVersion)
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

func TestDecodeAddressBase58OnlyLegacy(t *testing.T) {
	cases := decodedAddresses
	for _, c := range cases {
		withVersion, err := bitcoin.DecodeAddressBase58OnlyLegacy(c.address, true)
		require.NoError(t, err)
		withoutVersion, err := bitcoin.DecodeAddressBase58OnlyLegacy(c.address, false)
		require.NoError(t, err)
		assert.Equal(t, c.expected, withVersion)
		assert.Equal(t, c.expected[1:], withoutVersion)
	}
	bech32Addresses := []string{"tb1q2hxr4x5g4grwwrerf3y4tge776hmuw0wnh5vrd",
		"tb1qj9g0zjrj5r872hkkvxcedr3l504z50ayy5ercl",
		"tb1qpkv0lra0nz68ge5lzjjt6urdz2ejx8x4e9ell3",
		"tb1qqgzlw8yhyj6tmutat0u5n3dnxm3y6xnjp53wy9",
		"tb1qcc4j0tdu3lwfl05her3crlnvtqvltt90n5s5m0",
		"bc1qg5d579rlqmfekwx3m85a2sr8gy2s5dwfjj2lun",
		"bc1qtqxd29s9k3tj3rq9fzj7mnjknvlqzy8hsuzs5x",
		"bc1qv245zr29zw5urv5fy00c6km09l302fmlftf0aj",
		"bc1qw4z64jjvuxyddjdcm88yt0ln7fntkyw0w6wqhp",
		"bc1q8d7e3jrhsf8tj9q28x3msf8c644hdaetpqy7t4",
	}
	for _, c := range bech32Addresses {
		_, err := bitcoin.DecodeAddressBase58OnlyLegacy(c, true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "only legacy address allowed")
		_, err = bitcoin.DecodeAddressBase58OnlyLegacy(c, false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "only legacy address allowed")
	}
}
