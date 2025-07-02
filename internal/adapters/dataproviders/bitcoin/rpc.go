package bitcoin

import (
	"bytes"
	"fmt"
	merkle "github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcutil"
	"strings"
	"time"

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
	case wire.TestNet3:
		if !blockchain.IsTestnetBtcAddress(address) {
			return blockchain.BtcAddressInvalidNetworkError
		}
		return nil
	case wire.TestNet:
		if !blockchain.IsRegtestBtcAddress(address) {
			return blockchain.BtcAddressInvalidNetworkError
		}
		return nil
	default:
		return fmt.Errorf("unsupported network: %v", rpc.conn.NetworkParams.Net)
	}
}

func (rpc *bitcoindRpc) DecodeAddress(address string) ([]byte, error) {
	return DecodeAddress(address)
}

func (rpc *bitcoindRpc) GetTransactionInfo(hash string) (blockchain.BitcoinTransactionInformation, error) {
	// nolint:prealloc
	// false positive
	var amounts []*entities.Wei
	var btcAmount btcutil.Amount
	var ok bool

	parsedHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return blockchain.BitcoinTransactionInformation{}, err
	}

	tx, err := rpc.conn.client.GetRawTransactionVerbose(parsedHash)
	if err != nil {
		return blockchain.BitcoinTransactionInformation{}, err
	}

	outputs := make(map[string][]*entities.Wei)
	for _, output := range tx.Vout {
		amounts, ok = outputs[output.ScriptPubKey.Address]
		if !ok {
			amounts = make([]*entities.Wei, 0)
		}
		if btcAmount, err = btcutil.NewAmount(output.Value); err != nil {
			return blockchain.BitcoinTransactionInformation{}, err
		}
		amounts = append(amounts, entities.SatoshiToWei(uint64(btcAmount.ToUnit(btcutil.AmountSatoshi))))
		outputs[output.ScriptPubKey.Address] = amounts
	}
	return blockchain.BitcoinTransactionInformation{
		Hash:          tx.Txid,
		Confirmations: tx.Confirmations,
		Outputs:       outputs,
		HasWitness:    tx.Hash != tx.Txid,
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
	return SerializePartialMerkleTree(parsedTxHash, block)
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
	tx, err := rpc.conn.client.GetRawTransactionVerbose(parsedTxHash)
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

	blockHashBytes := ToSwappedBytes32(parsedBlockHash)
	return blockchain.BitcoinBlockInformation{
		Hash:   blockHashBytes,
		Height: big.NewInt(block.Height),
		Time:   time.Unix(block.Time, 0),
	}, nil
}

func (rpc *bitcoindRpc) GetCoinbaseInformation(txHash string) (blockchain.BtcCoinbaseTransactionInformation, error) {
	var coinbaseTxHash chainhash.Hash
	var witnessReservedValue [32]byte
	var err error

	block, _, err := rpc.getTxBlock(txHash)
	if err != nil {
		return blockchain.BtcCoinbaseTransactionInformation{}, err
	}
	txs := make([]*btcutil.Tx, 0)
	serializedCoinbase := bytes.NewBuffer([]byte{})

	for _, tx := range block.Transactions {
		if merkle.IsCoinBaseTx(tx) {
			if err = tx.SerializeNoWitness(serializedCoinbase); err != nil {
				return blockchain.BtcCoinbaseTransactionInformation{}, err
			}
			coinbaseTxHash = tx.TxHash()
			copy(witnessReservedValue[:], [][]byte(tx.TxIn[0].Witness)[0])
		}
		txs = append(txs, btcutil.NewTx(tx))
	}
	pmt, err := SerializePartialMerkleTree(&coinbaseTxHash, btcutil.NewBlock(block))
	if err != nil {
		return blockchain.BtcCoinbaseTransactionInformation{}, err
	}

	blockHash := block.BlockHash()
	blockVerboseInfo, err := rpc.conn.client.GetBlockVerbose(&blockHash)
	if err != nil {
		return blockchain.BtcCoinbaseTransactionInformation{}, err
	}

	return blockchain.BtcCoinbaseTransactionInformation{
		BtcTxSerialized:      serializedCoinbase.Bytes(),
		BlockHash:            ToSwappedBytes32(&blockHash),
		BlockHeight:          big.NewInt(blockVerboseInfo.Height),
		SerializedPmt:        pmt,
		WitnessMerkleRoot:    ToSwappedBytes32(merkle.CalcMerkleRoot(txs, true)),
		WitnessReservedValue: ToSwappedBytes32(witnessReservedValue),
	}, nil
}

func (rpc *bitcoindRpc) NetworkName() string {
	return strings.ToLower(rpc.conn.NetworkParams.Name)
}

func (rpc *bitcoindRpc) GetBlockchainInfo() (blockchain.BitcoinBlockchainInfo, error) {
	blockchainInfo, err := rpc.conn.client.GetBlockChainInfo()
	if err != nil {
		return blockchain.BitcoinBlockchainInfo{}, err
	}
	return blockchain.BitcoinBlockchainInfo{
		NetworkName:      blockchainInfo.Chain,
		ValidatedBlocks:  big.NewInt(int64(blockchainInfo.Blocks)),
		ValidatedHeaders: big.NewInt(int64(blockchainInfo.Headers)),
		BestBlockHash:    blockchainInfo.BestBlockHash,
	}, nil
}

func (rpc *bitcoindRpc) getTxBlock(txHash string) (*wire.MsgBlock, *chainhash.Hash, error) {
	parsedTxHash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return nil, nil, err
	}
	tx, err := rpc.conn.client.GetRawTransactionVerbose(parsedTxHash)
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
