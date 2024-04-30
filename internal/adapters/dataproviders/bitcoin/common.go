package bitcoin

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/btcutil/bloom"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin/btcclient"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"math/big"
)

const (
	BtcToSatoshi             = 100000000
	maxConfirmationsForUtxos = 9999999
	minConfirmationsForUtxos = 1
)

func DecodeAddressBase58(address string, keepVersion bool) ([]byte, error) {
	var buff bytes.Buffer
	addressBytes, version, err := base58.CheckDecode(address)
	if err != nil {
		return nil, err
	} else if len(addressBytes) != 20 {
		return nil, fmt.Errorf("decoded address exceeds 20 bytes (%s)", address)
	}
	if keepVersion {
		buff.WriteByte(version)
	}
	buff.Write(addressBytes)
	return buff.Bytes(), nil
}

func DecodeAddressBase58OnlyLegacy(address string, keepVersion bool) ([]byte, error) {
	if !blockchain.IsSupportedBtcAddress(address) {
		return nil, fmt.Errorf("only legacy address allowed (%s)", address)
	}
	return DecodeAddressBase58(address, keepVersion)
}

func toSwappedBytes32(hash *chainhash.Hash) [32]byte {
	var result [32]byte
	for i := 0; i < chainhash.HashSize/2; i++ {
		result[i], result[chainhash.HashSize-1-i] = hash[chainhash.HashSize-1-i], hash[i]
	}
	return result
}

func buildMerkleBranch(merkleTree []*chainhash.Hash, txCount uint32, txIndex uint32) blockchain.MerkleBranch {
	hashes := make([][32]byte, 0)
	path := 0
	pathIndex := 0
	var levelOffset uint32 = 0
	currentNodeOffset := txIndex

	for levelSize := txCount; levelSize > 1; levelSize = (levelSize + 1) / 2 {
		var targetOffset uint32
		if currentNodeOffset%2 == 0 {
			// Target is left hand side, use right hand side
			targetOffset = min(currentNodeOffset+1, levelSize-1)
		} else {
			// Target is right hand side, use left hand side
			targetOffset = currentNodeOffset - 1
			path = path + (1 << pathIndex)
		}
		hashes = append(hashes, toSwappedBytes32(merkleTree[levelOffset+targetOffset]))

		levelOffset += levelSize
		currentNodeOffset = currentNodeOffset / 2
		pathIndex++
	}

	return blockchain.MerkleBranch{
		Hashes: hashes,
		Path:   big.NewInt(int64(path)),
	}
}

func serializePartialMerkleTree(txHash *chainhash.Hash, block *btcutil.Block) ([]byte, error) {
	var err error
	filter := bloom.NewFilter(1, 0, 0, wire.BloomUpdateAll)
	filter.AddHash(txHash)

	msg, indices := bloom.NewMerkleBlock(block, filter)
	if len(indices) > 1 {
		return nil, fmt.Errorf("block matches more than one transaction (%v)", len(indices))
	}

	var buf bytes.Buffer
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(len(block.Transactions())))
	buf.Write(b)

	err = wire.WriteVarInt(&buf, wire.ProtocolVersion, uint64(len(msg.Hashes)))
	if err != nil {
		return nil, err
	}

	for _, h := range msg.Hashes {
		buf.Write(h[:])
	}
	err = wire.WriteVarInt(&buf, wire.ProtocolVersion, uint64(len(msg.Flags)))
	if err != nil {
		return nil, err
	}

	buf.Write(msg.Flags)
	return buf.Bytes(), nil
}

func getTransactionsToAddress(address string, params *chaincfg.Params, client btcclient.RpcClient) ([]blockchain.BitcoinTransactionInformation, error) {
	var ok bool
	var tx blockchain.BitcoinTransactionInformation
	var btcAmount btcutil.Amount
	result := make([]blockchain.BitcoinTransactionInformation, 0)
	parsedAddress, err := btcutil.DecodeAddress(address, params)
	if err != nil {
		return nil, err
	}
	utxos, err := client.ListUnspentMinMaxAddresses(0, maxConfirmationsForUtxos, []btcutil.Address{parsedAddress})
	if err != nil {
		return nil, err
	}

	txs := make(map[string]blockchain.BitcoinTransactionInformation)
	for _, utxo := range utxos {
		tx, ok = txs[utxo.TxID]
		if !ok {
			tx = blockchain.BitcoinTransactionInformation{
				Hash:          utxo.TxID,
				Confirmations: uint64(utxo.Confirmations),
				Outputs:       make(map[string][]*entities.Wei),
			}
			txs[utxo.TxID] = tx
		}
		if _, ok = tx.Outputs[utxo.Address]; !ok {
			tx.Outputs[utxo.Address] = make([]*entities.Wei, 0)
		}
		btcAmount, err = btcutil.NewAmount(utxo.Amount)
		if err != nil {
			return nil, err
		}
		tx.Outputs[utxo.Address] = append(tx.Outputs[utxo.Address], entities.SatoshiToWei(uint64(btcAmount.ToUnit(btcutil.AmountSatoshi))))
	}

	for key, value := range txs {
		result = append(result, value)
		delete(txs, key)
	}
	return result, nil
}
