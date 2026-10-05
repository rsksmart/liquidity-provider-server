package usecases

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

type RequestPegInInput struct {
	RskAddress  string
	DepositTxID string
	Amount      *entities.Wei
	Fee         *entities.Wei
	RawTx       []byte
}

func BuildRequestPegInParams(btc blockchain.BitcoinNetwork, input RequestPegInInput) (blockchain.RequestPegInParams, error) {
	block, err := btc.GetTransactionBlockInfo(input.DepositTxID)
	if err != nil {
		return blockchain.RequestPegInParams{}, err
	}
	merkle, err := btc.BuildMerkleBranch(input.DepositTxID)
	if err != nil {
		return blockchain.RequestPegInParams{}, err
	}
	return blockchain.RequestPegInParams{
		RskAddress:         input.RskAddress,
		BitcoinRawTx:       input.RawTx,
		BtcBlockHash:       block.Hash,
		MerkleBranchPath:   merkle.Path,
		MerkleBranchHashes: merkle.Hashes,
		Amount:             input.Amount,
		Fee:                input.Fee,
	}, nil
}
