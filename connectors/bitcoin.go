package connectors

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcjson"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/bloom"
	log "github.com/sirupsen/logrus"

	"github.com/btcsuite/btcutil/base58"
)

const unknownBtcdVersion = -1

type AddressWatcherCompleteCallback = func(w AddressWatcher)

type AddressWatcher interface {
	OnNewConfirmation(txHash string, confirmations int64, amount btcutil.Amount)
	OnExpire()
	Done() <-chan struct{}
}

type BTCConnector interface {
	Connect(endpoint string, username string, password string) error
	CheckConnection() error
	AddAddressWatcher(address string, minBtcAmount btcutil.Amount, interval time.Duration, exp time.Time, w AddressWatcher, cb AddressWatcherCompleteCallback) error
	GetParams() chaincfg.Params
	Close()
	SerializePMT(txHash string) ([]byte, error)
	SerializeTx(txHash string) ([]byte, error)
	GetBlockNumberByTx(txHash string) (int64, error)
}

type BTCClient interface {
	ImportAddressRescan(address string, account string, rescan bool) error
	GetTransaction(txHash *chainhash.Hash) (*btcjson.GetTransactionResult, error)
	GetBlockVerbose(blockHash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error)
	ListUnspentMinMaxAddresses(minConf, maxConf int, addrs []btcutil.Address) ([]btcjson.ListUnspentResult, error)
	GetBlock(blockHash *chainhash.Hash) (*wire.MsgBlock, error)
	GetRawTransaction(txHash *chainhash.Hash) (*btcutil.Tx, error)
	GetNetworkInfo() (*btcjson.GetNetworkInfoResult, error)
	Disconnect()
}

type BTC struct {
	c      BTCClient
	params chaincfg.Params
}

func NewBTC(network string) (*BTC, error) {
	log.Debug("initializing BTC connector")
	btc := BTC{}
	switch network {
	case "mainnet":
		btc.params = chaincfg.MainNetParams
	case "testnet":
		btc.params = chaincfg.TestNet3Params
	case "regtest":
		btc.params = chaincfg.RegressionNetParams
	default:
		return nil, fmt.Errorf("invalid network name: %v", network)
	}
	return &btc, nil
}

func (btc *BTC) Connect(endpoint string, username string, password string) error {
	log.Debug("connecting to BTC node")
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

	ver, err := checkBtcdVersion(c)
	if err != nil {
		return err
	}
	if ver == unknownBtcdVersion {
		log.Warn("unable to detect btcd version, but it is up and running")
	} else {
		log.Debugf("detected btcd version: %v", ver)
	}

	btc.c = c
	return nil
}

func (btc *BTC) CheckConnection() error {
	_, err := checkBtcdVersion(btc.c)
	return err
}

func (btc *BTC) AddAddressWatcher(address string, minBtcAmount btcutil.Amount, interval time.Duration, exp time.Time, w AddressWatcher, cb AddressWatcherCompleteCallback) error {
	btcAddr, err := btcutil.DecodeAddress(address, &btc.params)
	if err != nil {
		return fmt.Errorf("error decoding address %v: %v", address, err)
	}

	err = btc.c.ImportAddressRescan(address, "", false)
	if err != nil {
		return fmt.Errorf("error importing address %v: %v", address, err)
	}

	go func(w AddressWatcher) {
		ticker := time.NewTicker(interval)
		var confirmations int64
		for {
			select {
			case <-ticker.C:
				_ = btc.checkBtcAddr(w, btcAddr, minBtcAmount, exp, &confirmations, time.Now)
			case <-w.Done():
				ticker.Stop()
				cb(w)
				return
			}
		}
	}(w)
	return nil
}

func (btc *BTC) checkBtcAddr(w AddressWatcher, btcAddr btcutil.Address, minBtcAmount btcutil.Amount, expTime time.Time, confirmations *int64, now func() time.Time) error {
	conf, amount, txHash, err := btc.getConfirmations(btcAddr, minBtcAmount)
	if err != nil {
		log.Error(err)
		return err
	}
	if amount < minBtcAmount && now().After(expTime) {
		w.OnExpire()
		return fmt.Errorf("time for depositing %v has elapsed; addr: %v", minBtcAmount, btcAddr)
	}

	if conf > *confirmations {
		*confirmations = conf
		w.OnNewConfirmation(txHash, conf, amount)
		return nil
	}

	return fmt.Errorf("num of confirmations has not advanced; conf: %v", conf)
}

func (btc *BTC) GetParams() chaincfg.Params {
	return btc.params
}

func (btc *BTC) Close() {
	btc.c.Disconnect()
}

// SerializePMT computes the serialized partial merkle tree of a transaction in a block.
// The format of the serialized PMT is:
//
// - uint32     total_transactions (4 bytes)
// - varint     number of hashes   (1-3 bytes)
// - uint256[]  hashes in depth-first order (<= 32*N bytes)
// - varint     number of bytes of flag bits (1-3 bytes)
// - byte[]     flag bits, packed per 8 in a byte, least significant bit first (<= 2*N-1 bits)
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

func (btc *BTC) GetBlockNumberByTx(txHash string) (int64, error) {
	blockHash, err := btc.getBlockHash(txHash)
	if err != nil {
		return 0, err
	}
	msgBlock, err := btc.c.GetBlockVerbose(blockHash)
	if err != nil {
		return 0, fmt.Errorf("error retrieving block %v: %v", blockHash.String(), err)
	}
	return msgBlock.Height, nil
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

func DecodeBTCAddressWithVersion(address string) ([]byte, error) {
	addressBts, ver, err := base58.CheckDecode(address)
	if err != nil {
		return nil, fmt.Errorf("the provider address is not a valid base58 encoded address. address: %v", address)
	}
	var bts bytes.Buffer
	bts.WriteByte(ver)
	bts.Write(addressBts)
	return bts.Bytes(), nil
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

func (btc *BTC) getConfirmations(address btcutil.Address, minAmount btcutil.Amount) (int64, btcutil.Amount, string, error) {
	unspent, err := btc.c.ListUnspentMinMaxAddresses(0, 9999, []btcutil.Address{address})
	if err != nil {
		return 0, 0, "", fmt.Errorf("error retrieving unspent outputs for address %v: %v", address.EncodeAddress(), err)
	}
	if len(unspent) > 0 {
		type CumulativeResult struct {
			TxID          string
			Amount        btcutil.Amount
			Confirmations int64
		}
		var cumResults []*CumulativeResult

	outer:
		for _, u := range unspent {
			ua, err := btcutil.NewAmount(u.Amount)
			if err != nil {
				return 0, 0, "", fmt.Errorf("error retrieving unspent outputs for address %v: %v", address.EncodeAddress(), err)
			}

			for _, cr := range cumResults {
				if u.TxID == cr.TxID {
					cr.Amount += ua
					continue outer
				}
			}

			cumResults = append(cumResults, &CumulativeResult{
				TxID:          u.TxID,
				Amount:        ua,
				Confirmations: u.Confirmations,
			})
		}

		for _, cr := range cumResults {
			if cr.Amount >= minAmount {
				return cr.Confirmations, cr.Amount, cr.TxID, nil
			}
		}
	}

	return 0, 0, "", nil
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

func checkBtcdVersion(c BTCClient) (int32, error) {
	info, err := c.GetNetworkInfo()

	switch err := err.(type) {
	case nil:
		return info.Version, nil
	case *btcjson.RPCError:
		if err.Code != btcjson.ErrRPCMethodNotFound.Code {
			return 0, fmt.Errorf("unable to detect btcd version: %v", err)
		}
		return unknownBtcdVersion, nil
	default:
		return 0, fmt.Errorf("unable to detect btcd version: %v", err)
	}
}
