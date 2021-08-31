package connectors

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/bloom"
	log "github.com/sirupsen/logrus"
)

type BTC struct {
	c      *rpcclient.Client
	chans  map[string]*chan bool
	params chaincfg.Params
}

type AddressWatcher interface {
	OnNewConfirmation(txHash string, confirmations int64, amount float64)
}

func NewBTC() *BTC {
	return &BTC{chans: make(map[string]*chan bool)}
}

func (btc *BTC) Connect(endpoint string, username string, password string, network string) error {
	switch network {
	case "mainnet":
		btc.params = chaincfg.MainNetParams
	case "testnet":
		btc.params = chaincfg.TestNet3Params
	default:
		return fmt.Errorf("invalid network name: %v", network)
	}
	config := rpcclient.ConnConfig{
		Host:         endpoint,
		User:         username,
		Pass:         password,
		Params:       btc.params.Name,
		DisableTLS:   true,
		HTTPPostMode: true,
	}
	c, err := rpcclient.New(&config, nil)
	if err != nil {
		return fmt.Errorf("RPC client error: %v", err)
	}
	btc.c = c
	return nil
}

func (btc *BTC) AddAddressWatcher(address string, interval time.Duration, w AddressWatcher) error {
	btcAddr, err := btcutil.DecodeAddress(address, &btc.params)
	if err != nil {
		return fmt.Errorf("error decoding address %v: %v", address, err)
	}

	err = btc.c.ImportAddressRescan(address, "", false)
	if err != nil {
		return fmt.Errorf("error importing address %v: %v", address, err)
	}
	ticker := time.NewTicker(interval)
	ch := make(chan bool)
	btc.chans[address] = &ch

	go func(w AddressWatcher) {
		var confirmations int64
		for {
			select {
			case <-ticker.C:
				conf, amount, txHash, err := btc.getConfirmations(btcAddr)
				if err != nil {
					log.Error(err)
				}
				if conf > int64(confirmations) {
					confirmations = conf
					w.OnNewConfirmation(txHash, confirmations, amount)
				}
			case <-ch:
				ticker.Stop()
				return
			}
		}
	}(w)
	return nil
}

func (btc *BTC) GetParams() chaincfg.Params {
	return btc.params
}

func (btc *BTC) RemoveAddressWatcher(address string) {
	*btc.chans[address] <- true
}

func (btc *BTC) Close() {
	btc.c.Disconnect()
}

// Computes the serialized partial merkle tree of a transaction in a block.
// The format of the serialized PMT is:
//
// - uint32     total_transactions (4 bytes)
// - varint     number of hashes   (1-3 bytes)
// - uint256[]  hashes in depth-first order (<= 32*N bytes)
// - varint     number of bytes of flag bits (1-3 bytes)
// - byte[]     flag bits, packed per 8 in a byte, least significant bit first (<= 2*N-1 bits)
//
func (btc *BTC) SerializePMT(txHash string) ([]byte, error) {
	blockHash, err := btc.getBlockHash(txHash)
	if err != nil {
		return nil, err
	}
	msgBlock, err := btc.c.GetBlock(blockHash)
	if err != nil {
		return nil, fmt.Errorf("error retrieving block %v: %v", blockHash.String(), err)
	}
	block := btcutil.NewBlock(msgBlock)
	return serializePMT(txHash, block)
}

func (btc *BTC) SerializeTx(txHash string) ([]byte, error) {
	h, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return nil, fmt.Errorf("invalid tx hash %v: %v", txHash, err)
	}
	rawTx, err := btc.c.GetRawTransaction(h)
	if err != nil {
		return nil, fmt.Errorf("error retrieving tx %v: %v", txHash, err)
	}
	return serializeTx(rawTx)
}

func serializeTx(tx *btcutil.Tx) ([]byte, error) {
	var buf bytes.Buffer
	err := tx.MsgTx().SerializeNoWitness(&buf)
	if err != nil {
		return nil, fmt.Errorf("error serializing tx: %v", err)
	}
	return buf.Bytes(), nil
}

func serializePMT(txHash string, block *btcutil.Block) ([]byte, error) {
	filter := bloom.NewFilter(1, 0, 0, wire.BloomUpdateAll)
	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return nil, fmt.Errorf("error parsing hash: %v", err)
	}
	filter.AddHash(hash)

	msg, indices := bloom.NewMerkleBlock(block, filter)
	if len(indices) > 1 {
		return nil, fmt.Errorf("block matches more than one transaction (%v)", len(indices))
	}

	var buf bytes.Buffer

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(len(block.Transactions())))
	buf.Write(b)

	wire.WriteVarInt(&buf, wire.ProtocolVersion, uint64(len(msg.Hashes)))

	for _, h := range msg.Hashes {
		buf.Write(h[:])
	}
	wire.WriteVarInt(&buf, wire.ProtocolVersion, uint64(len(msg.Flags)))

	buf.Write(msg.Flags)

	return buf.Bytes(), nil
}

func (btc *BTC) getConfirmations(address btcutil.Address) (int64, float64, string, error) {
	unspent, err := btc.c.ListUnspentMinMaxAddresses(0, 9999, []btcutil.Address{address})
	if err != nil {
		return 0, 0, "", fmt.Errorf("error retrieving unspent outputs for address %v: %v", address.EncodeAddress(), err)
	}
	if len(unspent) > 0 {
		return unspent[0].Confirmations, unspent[0].Amount, unspent[0].TxID, nil
	} else {
		return 0, 0, "", nil
	}
}

func (btc *BTC) getBlockHash(txHash string) (*chainhash.Hash, error) {
	h, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction hash %v: %v", txHash, err)
	}
	tx, err := btc.c.GetTransaction(h)
	if err != nil {
		return nil, fmt.Errorf("error getting transaction %v: %v", txHash, err)
	}
	blockHash, err := chainhash.NewHashFromStr(tx.BlockHash)
	if err != nil {
		return nil, fmt.Errorf("error parsing block hash %v: %v", tx.BlockHash, err)
	}
	return blockHash, nil
}
