package bitcoin

import (
	"bytes"
	"fmt"
	merkle "github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcutil"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"math/big"
	"slices"
)

type bitcoindRpc struct {
	conn *Connection
}

func NewBitcoindRpc(conn *Connection) blockchain.BitcoinNetwork {
	return &bitcoindRpc{conn: conn}
}

func (rpc *bitcoindRpc) ValidateAddress(address string) error {
	if err := rpc.validateNetwork(address); err != nil {
		return err
	}
	if !blockchain.IsSupportedBtcAddress(address) {
		return blockchain.BtcAddressNotSupportedError
	}
	return nil
}

func (rpc *bitcoindRpc) validateNetwork(address string) error {
	switch rpc.conn.NetworkParams.Net {
	case wire.MainNet:
		if !blockchain.IsMainnetBtcAddress(address) {
			return blockchain.BtcAddressInvalidNetworkError
		}
		return nil
	case wire.TestNet, wire.TestNet3:
		if !blockchain.IsTestnetBtcAddress(address) {
			return blockchain.BtcAddressInvalidNetworkError
		}
		return nil
	default:
		return fmt.Errorf("unsupported network: %v", rpc.conn.NetworkParams.Net)
	}
}

func (rpc *bitcoindRpc) DecodeAddress(address string, keepVersion bool) ([]byte, error) {
	return DecodeAddressBase58(address, keepVersion)
}

func (rpc *bitcoindRpc) GetTransactionInfo(hash string) (blockchain.BitcoinTransactionInformation, error) {
	var amounts []*entities.Wei
	var ok bool

	parsedHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return blockchain.BitcoinTransactionInformation{}, err
	}

	tx, err := rpc.conn.client.GetTransaction(parsedHash)
	if err != nil {
		return blockchain.BitcoinTransactionInformation{}, err
	}

	outputs := make(map[string][]*entities.Wei)
	for _, output := range tx.Details {
		amounts, ok = outputs[output.Address]
		if !ok {
			amounts = make([]*entities.Wei, 0)
		}
		if output.Category == "send" { //  send category has a negative value
			amounts = append(amounts, entities.SatoshiToWei(uint64(output.Amount*BtcToSatoshi*-1)))
		} else {
			amounts = append(amounts, entities.SatoshiToWei(uint64(output.Amount*BtcToSatoshi)))
		}
		outputs[output.Address] = amounts
	}
	return blockchain.BitcoinTransactionInformation{
		Hash:          tx.TxID,
		Confirmations: uint64(tx.Confirmations),
		Outputs:       outputs,
	}, nil
}

func (rpc *bitcoindRpc) GetRawTransaction(hash string) ([]byte, error) {
	parsedHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return nil, err
	}

	rawTx, err := rpc.conn.client.GetRawTransaction(parsedHash)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err = rawTx.MsgTx().SerializeNoWitness(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (rpc *bitcoindRpc) GetPartialMerkleTree(hash string) ([]byte, error) {
	rawBlock, parsedTxHash, err := rpc.getTxBlock(hash)
	if err != nil {
		return nil, err
	}

	block := btcutil.NewBlock(rawBlock)
	return serializePartialMerkleTree(parsedTxHash, block)
}

func (rpc *bitcoindRpc) GetHeight() (*big.Int, error) {
	info, err := rpc.conn.client.GetBlockChainInfo()
	if err != nil {
		return nil, err
	}
	return big.NewInt(int64(info.Blocks)), nil
}

func (rpc *bitcoindRpc) BuildMerkleBranch(txHash string) (blockchain.MerkleBranch, error) {
	rawBlock, parsedTxHash, err := rpc.getTxBlock(txHash)
	if err != nil {
		return blockchain.MerkleBranch{}, err
	}

	block := btcutil.NewBlock(rawBlock)
	txs := make([]*btcutil.Tx, 0)
	for _, t := range block.MsgBlock().Transactions {
		txs = append(txs, btcutil.NewTx(t))
	}

	var cleanStore []*chainhash.Hash
	store := merkle.BuildMerkleTreeStore(txs, false)
	for _, node := range store {
		if node != nil {
			cleanStore = append(cleanStore, node)
		}
	}

	index := slices.IndexFunc(cleanStore, func(h *chainhash.Hash) bool {
		return h != nil && h.IsEqual(parsedTxHash)
	})
	if index == -1 {
		return blockchain.MerkleBranch{}, fmt.Errorf("transaction %s not found in merkle tree", txHash)
	}

	branch := buildMerkleBranch(cleanStore, uint32(len(block.Transactions())), uint32(index))
	return branch, nil
}

func (rpc *bitcoindRpc) GetTransactionBlockInfo(transactionHash string) (blockchain.BitcoinBlockInformation, error) {
	parsedTxHash, err := chainhash.NewHashFromStr(transactionHash)
	if err != nil {
		return blockchain.BitcoinBlockInformation{}, err
	}
	tx, err := rpc.conn.client.GetTransaction(parsedTxHash)
	if err != nil {
		return blockchain.BitcoinBlockInformation{}, err
	}

	parsedBlockHash, err := chainhash.NewHashFromStr(tx.BlockHash)
	if err != nil {
		return blockchain.BitcoinBlockInformation{}, err
	}

	block, err := rpc.conn.client.GetBlockVerbose(parsedBlockHash)
	if err != nil {
		return blockchain.BitcoinBlockInformation{}, err
	}

	blockHashBytes := toSwappedBytes32(parsedBlockHash)
	return blockchain.BitcoinBlockInformation{
		Hash:   blockHashBytes,
		Height: big.NewInt(block.Height),
	}, nil
}

func (rpc *bitcoindRpc) getTxBlock(txHash string) (*wire.MsgBlock, *chainhash.Hash, error) {
	parsedTxHash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return nil, nil, err
	}
	tx, err := rpc.conn.client.GetTransaction(parsedTxHash)
	if err != nil {
		return nil, nil, err
	}
	parsedBlockHash, err := chainhash.NewHashFromStr(tx.BlockHash)
	if err != nil {
		return nil, nil, err
	}
	block, err := rpc.conn.client.GetBlock(parsedBlockHash)
	return block, parsedTxHash, err
}
