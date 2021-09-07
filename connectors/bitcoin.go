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

	"encoding/hex"

	"github.com/btcsuite/btcd/txscript"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"
)

type BTC struct {
	c       *rpcclient.Client
	chans   map[string]*chan bool
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

type AddressWatcher interface {
	OnNewConfirmation(txHash string, confirmations int64, amount float64)
}

func NewBTC(network string, fedInfo FedInfo) (*BTC, error) {
	btc := BTC{
		chans:   make(map[string]*chan bool),
		fedInfo: &fedInfo,
	}
	switch network {
	case "mainnet":
		btc.params = chaincfg.MainNetParams
	case "testnet":
		btc.params = chaincfg.TestNet3Params
	default:
		return nil, fmt.Errorf("invalid network name: %v", network)
	}
	return &btc, nil
}

func (btc *BTC) Connect(endpoint string, username string, password string) error {
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
	} else {
		hashBuf, err = btc.getErpRedeemScriptBuf()
		if err != nil {
			return nil, err
		}
	}

	err = btc.validateRedeemScript(hashBuf.Bytes())
	if err != nil {
		return nil, err
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
