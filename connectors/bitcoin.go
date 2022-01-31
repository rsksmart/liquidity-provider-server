package connectors

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/bloom"
	log "github.com/sirupsen/logrus"

	"encoding/hex"

	"github.com/btcsuite/btcd/txscript"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"
)

type AddressWatcher interface {
	OnNewConfirmation(txHash string, confirmations int64, amount btcutil.Amount)
	OnExpire()
	Done() <-chan struct{}
}

type BTCConnector interface {
	Connect(endpoint string, username string, password string) error
	AddAddressWatcher(address string, minAmount btcutil.Amount, interval time.Duration, exp time.Time, w AddressWatcher) error
	GetParams() chaincfg.Params
	Close()
	SerializePMT(txHash string) ([]byte, error)
	SerializeTx(txHash string) ([]byte, error)
	GetBlockNumberByTx(txHash string) (int64, error)
	GetDerivedBitcoinAddress(userBtcRefundAddr []byte, lbcAddress []byte, lpBtcAddress []byte, derivationArgumentsHash []byte) (string, error)
}

type BTCClient interface {
	ImportAddressRescan(address string, account string, rescan bool) error
	GetTransaction(txHash *chainhash.Hash) (*btcjson.GetTransactionResult, error)
	GetBlockVerbose(blockHash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error)
	ListUnspentMinMaxAddresses(minConf, maxConf int, addrs []btcutil.Address) ([]btcjson.ListUnspentResult, error)
	GetBlock(blockHash *chainhash.Hash) (*wire.MsgBlock, error)
	GetRawTransaction(txHash *chainhash.Hash) (*btcutil.Tx, error)
	Disconnect()
}

type BTC struct {
	c       BTCClient
	params  chaincfg.Params
	fedInfo *FedInfo
}

type FedInfo struct {
	FedSize              int
	FedThreshold         int
	PubKeys              []string
	FedAddress           string
	ActiveFedBlockHeight int
	IrisActivationHeight int
	ErpKeys              []string
}

func NewBTC(network string, fedInfo FedInfo) (*BTC, error) {
	log.Debug("initializing BTC connector")
	btc := BTC{
		fedInfo: &fedInfo,
	}
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

	info, err := c.GetNetworkInfo()

	switch err := err.(type) {
	case nil:
		log.Debugf("detected btcd version: %v", info.Version)
	// Inspect the RPC error to ensure the method was not found, otherwise
	// we actually ran into an error.
	case *btcjson.RPCError:
		if err.Code != btcjson.ErrRPCMethodNotFound.Code {
			return fmt.Errorf("unable to detect btcd version: %v", err)
		}
	default:
		return fmt.Errorf("unable to detect btcd version: %v", err)
	}

	btc.c = c
	return nil
}

func (btc *BTC) AddAddressWatcher(address string, minBtcAmount btcutil.Amount, interval time.Duration, exp time.Time, w AddressWatcher) error {
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

func (btc *BTC) GetDerivedBitcoinAddress(userBtcRefundAddr []byte, lbcAddress []byte, lpBtcAddress []byte, derivationArgumentsHash []byte) (string, error) {
	derivationValue, err := getDerivationValueHash(userBtcRefundAddr, lbcAddress, lpBtcAddress, derivationArgumentsHash)
	if err != nil {
		return "", fmt.Errorf("error computing derivation value: %v", err)
	}
	flyoverScript, err := btc.getRedeemScript(derivationValue)
	if err != nil {
		return "", fmt.Errorf("error generating redeem script: %v", err)
	}
	addressScriptHash, err := btcutil.NewAddressScriptHash(flyoverScript, &btc.params)
	if err != nil {
		return "", err
	}
	return addressScriptHash.EncodeAddress(), nil
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

func getDerivationValueHash(userBtcRefundAddr []byte, lbcAddress []byte, lpBtcAddress []byte, derivationArgumentsHash []byte) ([]byte, error) {
	var buf bytes.Buffer
	buf.Write(derivationArgumentsHash)
	buf.Write(userBtcRefundAddr)
	buf.Write(lbcAddress)
	buf.Write(lpBtcAddress)

	derivationValueHash := crypto.Keccak256(buf.Bytes())

	return derivationValueHash, nil
}

func (btc *BTC) validateRedeemScript(script []byte) error {
	addr, err := btcutil.NewAddressScriptHash(script, &btc.params)
	if err != nil {
		return err
	}

	fedAddress, err := btcutil.DecodeAddress(btc.fedInfo.FedAddress, &btc.params)
	if err != nil {
		return err
	}

	if !bytes.Equal(addr.ScriptAddress(), fedAddress.ScriptAddress()) {
		return fmt.Errorf("the generated redeem script does not match with the federation redeem script")
	}
	return nil
}

func (btc *BTC) getRedeemScript(derivationValue []byte) ([]byte, error) {
	var hashBuf *bytes.Buffer

	buf, err := getFlyoverPrefix(derivationValue)
	if err != nil {
		return nil, err
	}

	// All federations activated AFTER Iris will be ERP, therefore we build erp redeem script.
	if btc.fedInfo.ActiveFedBlockHeight < btc.fedInfo.IrisActivationHeight {
		hashBuf, err = btc.getPowPegRedeemScriptBuf(true)
		if err != nil {
			return nil, err
		}

		err = btc.validateRedeemScript(hashBuf.Bytes())
		if err != nil {
			return nil, err
		}
	} else {
		hashBuf, err = btc.getErpRedeemScriptBuf()
		if err != nil {
			return nil, err
		}

		err = btc.validateRedeemScript(hashBuf.Bytes())
		if err != nil { // ok, it could be that ERP is not yet activated, falling back to PowPeg Redeem Script
			hashBuf, err = btc.getPowPegRedeemScriptBuf(true)
			if err != nil {
				return nil, err
			}

			err = btc.validateRedeemScript(hashBuf.Bytes())
			if err != nil {
				return nil, err
			}
		}
	}

	buf.Write(hashBuf.Bytes())
	return buf.Bytes(), nil
}

func getFlyoverPrefix(hash []byte) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	hashPrefix, err := hex.DecodeString("20")
	if err != nil {
		return nil, err
	}
	buf.Write(hashPrefix)
	buf.Write(hash)
	buf.WriteByte(txscript.OP_DROP)

	return &buf, nil
}

func (btc *BTC) getPowPegRedeemScriptBuf(addMultiSig bool) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	sb := txscript.NewScriptBuilder()
	err := btc.addStdNToMScriptPart(sb)
	if err != nil {
		return nil, err
	}
	if addMultiSig {
		sb.AddOp(txscript.OP_CHECKMULTISIG)
	}

	sbuf, err := sb.Script()
	if err != nil {
		return nil, err
	}
	buf.Write(sbuf)
	return &buf, nil
}

func (btc *BTC) getErpRedeemScriptBuf() (*bytes.Buffer, error) {
	erpRedeemScriptBuf, err := btc.p2ms(false)
	if err != nil {
		return nil, err
	}
	powPegRedeemScriptBuf, err := btc.getPowPegRedeemScriptBuf(false)
	if err != nil {
		return nil, err
	}
	scriptsA := txscript.NewScriptBuilder()
	scriptsA.AddOp(txscript.OP_NOTIF)
	var erpRedeemScriptBuffer bytes.Buffer
	scrA, err := scriptsA.Script()
	if err != nil {
		return nil, err
	}
	erpRedeemScriptBuffer.Write(scrA)
	erpRedeemScriptBuffer.Write(powPegRedeemScriptBuf.Bytes())
	erpRedeemScriptBuffer.WriteByte(txscript.OP_ELSE)
	byteArr, err := hex.DecodeString("02")
	if err != nil {
		return nil, err
	}
	erpRedeemScriptBuffer.Write(byteArr)

	csv, err := hex.DecodeString(btc.getCsvValueFromNetwork())
	if err != nil {
		return nil, err
	}
	erpRedeemScriptBuffer.Write(csv)
	erpRedeemScriptBuffer.WriteByte(txscript.OP_CHECKSEQUENCEVERIFY)
	erpRedeemScriptBuffer.WriteByte(txscript.OP_DROP)
	erpRedeemScriptBuffer.Write(erpRedeemScriptBuf.Bytes())
	erpRedeemScriptBuffer.WriteByte(txscript.OP_ENDIF)
	erpRedeemScriptBuffer.WriteByte(txscript.OP_CHECKMULTISIG)

	return &erpRedeemScriptBuffer, nil
}

func (btc *BTC) p2ms(addMultiSig bool) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	sb := txscript.NewScriptBuilder()
	err := btc.addErpNToMScriptPart(sb)
	if err != nil {
		return nil, err
	}
	if addMultiSig {
		sb.AddOp(txscript.OP_CHECKMULTISIG)
	}
	sbuf, err := sb.Script()
	if err != nil {
		return nil, err
	}
	buf.Write(sbuf)
	return &buf, nil
}

func (btc *BTC) getCsvValueFromNetwork() string {
	switch btc.params.Name {
	case chaincfg.MainNetParams.Name:
		return "CD50"
	case chaincfg.TestNet3Params.Name:
		return "CD50"
	default: // RegTest
		return "01F4"
	}
}

func (btc *BTC) addStdNToMScriptPart(builder *txscript.ScriptBuilder) error {
	builder.AddOp(getOpCodeFromInt(btc.fedInfo.FedThreshold))

	for _, pubKey := range btc.fedInfo.PubKeys {
		pkBuffer, err := hex.DecodeString(pubKey)
		if err != nil {
			return err
		}
		builder.AddData(pkBuffer)
	}

	builder.AddOp(getOpCodeFromInt(btc.fedInfo.FedSize))

	return nil
}

func (btc *BTC) addErpNToMScriptPart(builder *txscript.ScriptBuilder) error {
	size := len(btc.fedInfo.ErpKeys)
	min := size/2 + 1
	builder.AddOp(getOpCodeFromInt(min))

	for _, pubKey := range btc.fedInfo.ErpKeys {
		pkBuffer, err := hex.DecodeString(pubKey)
		if err != nil {
			return err
		}
		builder.AddData(pkBuffer)
	}

	builder.AddOp(getOpCodeFromInt(len(btc.fedInfo.ErpKeys)))

	return nil
}

func getOpCodeFromInt(val int) byte {
	switch val {
	case 1:
		return txscript.OP_1
	case 2:
		return txscript.OP_2
	case 3:
		return txscript.OP_3
	case 4:
		return txscript.OP_4
	case 5:
		return txscript.OP_5
	case 6:
		return txscript.OP_6
	case 7:
		return txscript.OP_7
	case 8:
		return txscript.OP_8
	case 9:
		return txscript.OP_9
	case 10:
		return txscript.OP_10
	case 11:
		return txscript.OP_11
	case 12:
		return txscript.OP_12
	case 13:
		return txscript.OP_13
	case 14:
		return txscript.OP_14
	case 15:
		return txscript.OP_15
	default:
		return txscript.OP_16
	}
}
