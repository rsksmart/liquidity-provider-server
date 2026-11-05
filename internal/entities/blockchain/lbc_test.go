package blockchain_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestRegistrationState_AllowsRegistration(t *testing.T) {
	cases := map[blockchain.RegistrationState]bool{
		blockchain.RegistrationStateNone:      true,
		blockchain.RegistrationStateRejected:  true,
		blockchain.RegistrationStateWithdrawn: true,
		blockchain.RegistrationStatePending:   false,
		blockchain.RegistrationStateApproved:  false,
	}
	for state, expected := range cases {
		assert.Equal(t, expected, state.AllowsRegistration())
	}
}

func TestRefundPegoutParams_String(t *testing.T) {
	params := blockchain.RefundPegoutParams{
		QuoteHash:          [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
		BtcRawTx:           []byte{0x01, 0x02, 0x03},
		BtcBlockHeaderHash: [32]byte{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		MerkleBranchPath:   big.NewInt(1),
		MerkleBranchHashes: [][32]byte{
			{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
			{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		},
	}
	assert.Equal(t,
		"RefundPegoutParams { QuoteHash: 0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20, "+
			"BtcRawTx: 010203, BtcBlockHeaderHash: 201f1e1d1c1b1a191817161514131211100f0e0d0c0b0a090807060504030201, "+
			"MerkleBranchPath: 1, MerkleBranchHashes: ["+
			"[32 31 30 29 28 27 26 25 24 23 22 21 20 19 18 17 16 15 14 13 12 11 10 9 8 7 6 5 4 3 2 1] "+
			"[1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32] "+
			"[32 31 30 29 28 27 26 25 24 23 22 21 20 19 18 17 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16]] }",
		params.String())
}
