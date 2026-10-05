package usecases_test

import (
	"math/big"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	u "github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const requestPegInDepositTxID = "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"

func requestPegInInput() u.RequestPegInInput {
	return u.RequestPegInInput{
		RskAddress:  "0x79568c2989232dCa1840087D73d403602364c0D4",
		DepositTxID: requestPegInDepositTxID,
		Amount:      entities.NewWei(1000),
		Fee:         entities.NewWei(10),
		RawTx:       []byte{0x01, 0x00, 0x00, 0x00, 0x01},
	}
}

func TestBuildRequestPegInParams(t *testing.T) {
	t.Run("maps the input and the deposit proofs", func(t *testing.T) {
		btc := mocks.NewBtcRpcMock(t)
		block := blockchain.BitcoinBlockInformation{Hash: [32]byte{0x09}, Height: big.NewInt(500)}
		merkle := blockchain.MerkleBranch{Path: big.NewInt(3), Hashes: [][32]byte{{0x0b}, {0x0c}}}
		btc.On("GetTransactionBlockInfo", requestPegInDepositTxID).Return(block, nil).Once()
		btc.On("BuildMerkleBranch", requestPegInDepositTxID).Return(merkle, nil).Once()
		input := requestPegInInput()

		params, err := u.BuildRequestPegInParams(btc, input)
		require.NoError(t, err)
		assert.Equal(t, blockchain.RequestPegInParams{
			RskAddress:         input.RskAddress,
			BitcoinRawTx:       input.RawTx,
			BtcBlockHash:       block.Hash,
			MerkleBranchPath:   merkle.Path,
			MerkleBranchHashes: merkle.Hashes,
			Amount:             input.Amount,
			Fee:                input.Fee,
		}, params)
	})
	t.Run("block info error skips the merkle branch", func(t *testing.T) {
		btc := mocks.NewBtcRpcMock(t)
		btc.On("GetTransactionBlockInfo", requestPegInDepositTxID).
			Return(blockchain.BitcoinBlockInformation{}, assert.AnError).Once()

		params, err := u.BuildRequestPegInParams(btc, requestPegInInput())
		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, params)
		btc.AssertNotCalled(t, "BuildMerkleBranch", mock.Anything)
	})
	t.Run("merkle branch error", func(t *testing.T) {
		btc := mocks.NewBtcRpcMock(t)
		btc.On("GetTransactionBlockInfo", requestPegInDepositTxID).
			Return(blockchain.BitcoinBlockInformation{Hash: [32]byte{0x09}}, nil).Once()
		btc.On("BuildMerkleBranch", requestPegInDepositTxID).Return(blockchain.MerkleBranch{}, assert.AnError).Once()

		params, err := u.BuildRequestPegInParams(btc, requestPegInInput())
		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, params)
	})
}
