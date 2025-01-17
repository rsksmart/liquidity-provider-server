package bitcoin

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/btcutil/bech32"
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
	MaxConfirmationsForUtxos = 9999999
	MinConfirmationsForUtxos = 1
)

func DecodeAddress(address string) ([]byte, error) {
	if blockchain.IsBtcP2PKHAddress(address) || blockchain.IsBtcP2SHAddress(address) {
		return DecodeAddressBase58(address, true)
	} else if blockchain.IsBtcP2WPKHAddress(address) || blockchain.IsBtcP2WSHAddress(address) || blockchain.IsBtcP2TRAddress(address) {
		_, data, err := bech32.Decode(address) // this function decodes both bech32 and bech32m
		return data, err
	}
	return nil, blockchain.BtcAddressNotSupportedError
}

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

func ToSwappedBytes32[T [32]byte | *chainhash.Hash | chainhash.Hash](hash T) [32]byte {
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
		hashes = append(hashes, ToSwappedBytes32(merkleTree[levelOffset+targetOffset]))

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
	var err error
	var tx blockchain.BitcoinTransactionInformation
	var rawTxInfo *btcjson.TxRawResult
	var outputs map[string][]*entities.Wei

	result := make([]blockchain.BitcoinTransactionInformation, 0)
	parsedAddress, err := btcutil.DecodeAddress(address, params)
	if err != nil {
		return nil, err
	}
	utxos, err := client.ListUnspentMinMaxAddresses(0, MaxConfirmationsForUtxos, []btcutil.Address{parsedAddress})
	if err != nil {
		return nil, err
	}

	txs := make(map[string]blockchain.BitcoinTransactionInformation)
	for _, utxo := range utxos {
		_, ok = txs[utxo.TxID]
		if !ok {
			txId, _ := chainhash.NewHashFromStr(utxo.TxID)
			if rawTxInfo, err = client.GetRawTransactionVerbose(txId); err != nil {
				return nil, err
			}
			if outputs, err = parseTransactionOutputs(rawTxInfo.Vout); err != nil {
				return nil, err
			}
			tx = blockchain.BitcoinTransactionInformation{
				Hash:          rawTxInfo.Txid,
				Confirmations: rawTxInfo.Confirmations,
				Outputs:       outputs,
				HasWitness:    rawTxInfo.Hash != rawTxInfo.Txid,
			}
			txs[utxo.TxID] = tx
		}
	}

	for key, value := range txs {
		result = append(result, value)
		delete(txs, key)
	}
	return result, nil
}

func parseTransactionOutputs(outputs []btcjson.Vout) (map[string][]*entities.Wei, error) {
	var err error
	var amount btcutil.Amount
	result := make(map[string][]*entities.Wei)
	for _, output := range outputs {
		if _, ok := result[output.ScriptPubKey.Address]; !ok {
			result[output.ScriptPubKey.Address] = make([]*entities.Wei, 0)
		}
		if amount, err = btcutil.NewAmount(output.Value); err != nil {
			return nil, err
		}
		result[output.ScriptPubKey.Address] = append(result[output.ScriptPubKey.Address], entities.SatoshiToWei(uint64(amount.ToUnit(btcutil.AmountSatoshi))))
	}
	return result, nil
}
