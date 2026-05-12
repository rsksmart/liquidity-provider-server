package rootstock_test

import (
	"encoding/hex"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestBtcCoinbaseTransactionInformation_String(t *testing.T) {
	var blockHash, witnessMerkleRoot, witnessReservedValue [32]byte
	tx, err := hex.DecodeString("020000000001018e1993e43f182c6966ac011f12d82c18ee2b2e292b23f206e5c55d518cded7e80100000000fdffffff0300879303000000001976a914d60c3f1e0a8e76dd5ea1470c968b87b9b0339c4988ac0000000000000000226a2042be5ef1f59c24d2715f6f4b803a2acc66515554447f1a3e0abb99a3317aa6afc11fa11900000000160014ddb677f36498f7a4901a74e882df68fd00cf473502473044022077657caef5a7692e3ac1dffca4cfebea98029a21dbf5247a044ef4d2a8f2fdfd02206342029f868122a7c2321b67cea2440c98925728450f28b3e443e80c4f95765e01210232858a5faa413101831afe7a880da9a8ac4de6bd5e25b4358d762ba450b03c2200000000")
	require.NoError(t, err)
	blockHashBytes, err := hex.DecodeString("8e1993e43f182c6966ac011f12d82c18ee2b2e292b23f206e5c55d518cded7e8")
	require.NoError(t, err)
	copy(blockHash[:], blockHashBytes)
	witnessMerkleRootBytes, err := hex.DecodeString("42be5ef1f59c24d2715f6f4b803a2acc66515554447f1a3e0abb99a3317aa6af")
	require.NoError(t, err)
	copy(witnessMerkleRoot[:], witnessMerkleRootBytes)
	witnessReservedValueBytes, err := hex.DecodeString("ddb677f36498f7a4901a74e882df68fd00cf4735")
	require.NoError(t, err)
	copy(witnessReservedValue[:], witnessReservedValueBytes)
	transaction := rootstock.BtcCoinbaseTransactionInformation{
		BtcTxSerialized:      tx,
		BlockHash:            blockHash,
		BlockHeight:          big.NewInt(123456789),
		SerializedPmt:        []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
		WitnessMerkleRoot:    witnessMerkleRoot,
		WitnessReservedValue: witnessReservedValue,
	}
	expected := "RegisterPeginParams { BtcTxSerialized: 020000000001018e1993e43f182c6966ac011f12d82c18ee2b2e292b23f206e5c55d518cded7e80100000000fdffffff0300879303000000001976a914d60c3f1e0a8e76dd5ea1470c968b87b9b0339c4988ac0000000000000000226a2042be5ef1f59c24d2715f6f4b803a2acc66515554447f1a3e0abb99a3317aa6afc11fa11900000000160014ddb677f36498f7a4901a74e882df68fd00cf473502473044022077657caef5a7692e3ac1dffca4cfebea98029a21dbf5247a044ef4d2a8f2fdfd02206342029f868122a7c2321b67cea2440c98925728450f28b3e443e80c4f95765e01210232858a5faa413101831afe7a880da9a8ac4de6bd5e25b4358d762ba450b03c2200000000, BlockHash: 8e1993e43f182c6966ac011f12d82c18ee2b2e292b23f206e5c55d518cded7e8, BlockHeight: 123456789SerializedPmt: 010203040506070809, WitnessMerkleRoot: 42be5ef1f59c24d2715f6f4b803a2acc66515554447f1a3e0abb99a3317aa6af, WitnessReservedValue: ddb677f36498f7a4901a74e882df68fd00cf4735000000000000000000000000 }"
	assert.Equal(t, expected, transaction.String())
}
