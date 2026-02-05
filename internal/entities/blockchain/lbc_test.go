package blockchain_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

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

func TestRegisterPeginParams_String(t *testing.T) {
	params := blockchain.RegisterPeginParams{
		QuoteSignature:        []byte{0x01, 0x02, 0x03},
		BitcoinRawTransaction: []byte{0x04, 0x05, 0x06},
		PartialMerkleTree:     []byte{0x07, 0x08, 0x09},
		BlockHeight:           big.NewInt(1),
		Quote: quote.PeginQuote{
			FedBtcAddress:      test.AnyAddress,
			LbcAddress:         test.AnyAddress,
			LpRskAddress:       test.AnyAddress,
			BtcRefundAddress:   test.AnyAddress,
			RskRefundAddress:   test.AnyAddress,
			LpBtcAddress:       test.AnyAddress,
			CallFee:            entities.NewWei(3),
			PenaltyFee:         entities.NewWei(4),
			ContractAddress:    test.AnyAddress,
			Data:               "any data",
			GasLimit:           5,
			Nonce:              6,
			Value:              entities.NewWei(2),
			AgreementTimestamp: 7,
			TimeForDeposit:     8,
			LpCallTime:         9,
			Confirmations:      10,
			CallOnRegister:     true,
			GasFee:             entities.NewWei(1),
		},
	}
	assert.Equal(t,
		"RegisterPeginParams { QuoteSignature: 010203, BitcoinRawTransaction: 040506, PartialMerkleTree: 070809, "+
			"BlockHeight: 1, Quote: {FedBtcAddress:any address LbcAddress:any address LpRskAddress:any address "+
			"BtcRefundAddress:any address RskRefundAddress:any address LpBtcAddress:any address CallFee:3 PenaltyFee:4 "+
			"ContractAddress:any address Data:any data GasLimit:5 Nonce:6 Value:2 AgreementTimestamp:7 TimeForDeposit:8 "+
			"LpCallTime:9 Confirmations:10 CallOnRegister:true GasFee:1} }", params.String())
}
