package connectors

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"regexp"
	"time"

	"github.com/btcsuite/btcd/btcjson"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/bloom"
	log "github.com/sirupsen/logrus"

	"encoding/hex"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"
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
	AddAddressPegOutWatcher(address string, minBtcAmount btcutil.Amount, interval time.Duration, exp time.Time, w AddressWatcher, cb AddressWatcherCompleteCallback) error
	GetParams() chaincfg.Params
	Close()
	SerializePMT(txHash string) ([]byte, error)
	SerializeTx(txHash string) ([]byte, error)
	GetBlockNumberByTx(txHash string) (int64, error)
	GetDerivedBitcoinAddress(fedInfo *FedInfo, userBtcRefundAddr []byte, lbcAddress []byte, lpBtcAddress []byte, derivationArgumentsHash []byte) (string, error)
	ComputeDerivationAddresss(userBtcRefundAddr []byte, quoteHash []byte) (string, error)
	BuildMerkleBranch(txHash string) (*MerkleBranch, error)
	BuildMerkleBranchByEndpoint(txHash string, btcAddress string) (*MerkleBranch, error)
	SendBTC(address string, amount uint) (string, error)
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
	SendToAddress(address btcutil.Address, amount btcutil.Amount) (*chainhash.Hash, error)
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
		return err
	}

	err = btc.c.ImportAddressRescan(address, "", false)
	if err != nil {
		return err
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

func (btc *BTC) AddAddressPegOutWatcher(address string, minBtcAmount btcutil.Amount, interval time.Duration, exp time.Time, w AddressWatcher, cb AddressWatcherCompleteCallback) error {
	btcAddr, err := btcutil.DecodeAddress(address, &btc.params)
	if err != nil {
		return err
	}

	err = btc.c.ImportAddressRescan(address, "", false)
	if err != nil {
		return err
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

func buildErrorImportAddress(address string, err error) error {
	log.Errorf("error importing address %v: %v", address, err)
	return fmt.Errorf("error importing address %v: %v", address, err)
}

func (btc *BTC) checkBtcAddr(w AddressWatcher, btcAddr btcutil.Address, minBtcAmount btcutil.Amount, expTime time.Time, confirmations *int64, now func() time.Time) error {
	log.Debugf("Derivation Address:: %v", btcAddr)
	log.Debugf("minBtcAmount:: %v", minBtcAmount)
	log.Debugf("confirmations:: %v", confirmations)
	conf, amount, txHash, err := btc.getConfirmations(btcAddr, minBtcAmount)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Debugf("amount:: %v", amount)
	log.Debugf("txHash:: %v", txHash)
	log.Debugf("expTime:: %v", expTime)
	log.Debugf("now:: %v", now())

	if amount < minBtcAmount && now().After(expTime) {
		w.OnExpire()
		return fmt.Errorf("time for depositing %v has elapsed; addr: %v", minBtcAmount, btcAddr)
	}

	if conf > *confirmations {
		*confirmations = conf
		w.OnNewConfirmation(txHash, conf, amount)
		branch, err := btc.BuildMerkleBranch(txHash)
		if err != nil {
			return err
		}

		log.Debugf("Merkle Branch info :::: path:%v hashes:%v ", branch.Path, branch.Hashes)

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
		return nil, buildErrorRetrievingBlock(blockHash, err)
	}
	block := btcutil.NewBlock(msgBlock)
	return serializePMT(txHash, block)
}

func buildErrorRetrievingBlock(blockHash *chainhash.Hash, err error) error {
	return fmt.Errorf("error retrieving block %v: %v", blockHash.String(), err)
}

func (btc *BTC) GetBlockNumberByTx(txHash string) (int64, error) {
	blockHash, err := btc.getBlockHash(txHash)
	if err != nil {
		return 0, err
	}
	msgBlock, err := btc.c.GetBlockVerbose(blockHash)
	if err != nil {
		return 0, buildErrorRetrievingBlock(blockHash, err)
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

func (btc *BTC) GetDerivedBitcoinAddress(fedInfo *FedInfo, userBtcRefundAddr []byte, lbcAddress []byte, lpBtcAddress []byte, derivationArgumentsHash []byte) (string, error) {
	derivationValue, err := getDerivationValueHash(userBtcRefundAddr, lbcAddress, lpBtcAddress, derivationArgumentsHash)
	if err != nil {
		return "", fmt.Errorf("error computing derivation value: %v", err)
	}
	flyoverScript, err := btc.getRedeemScript(fedInfo, derivationValue)
	if err != nil {
		return "", fmt.Errorf("error generating redeem script: %v", err)
	}
	addressScriptHash, err := btcutil.NewAddressScriptHash(flyoverScript, &btc.params)
	if err != nil {
		return "", err
	}
	return addressScriptHash.EncodeAddress(), nil
}

func (btc *BTC) ComputeDerivationAddresss(userPublicKey []byte, quoteHash []byte) (string, error) {

	rootScriptBuilder := txscript.NewScriptBuilder()

	rootScriptBuilder.AddData(quoteHash)

	rootScriptBuilder.AddOp(txscript.OP_DROP)

	rootScriptBuilder.AddOp(txscript.OP_1)

	rootScriptBuilder.AddData(userPublicKey)

	rootScriptBuilder.AddOp(txscript.OP_1)

	rootScriptBuilder.AddOp(txscript.OP_CHECKMULTISIG)

	rootScript, err := rootScriptBuilder.Script()

	if err != nil {
		return "", fmt.Errorf("error generating root script: %v", err)
	}

	redeemScript, err := btcutil.NewAddressScriptHash(rootScript[:], &btc.params)

	if err != nil {
		return "", err
	}

	return redeemScript.EncodeAddress(), nil
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

func (btc *BTC) BuildMerkleBranch(txHash string) (*MerkleBranch, error) {
	blockHash, err := btc.getBlockHash(txHash)
	if err != nil {
		return nil, err
	}
	msgBlock, err := btc.c.GetBlock(blockHash)
	if err != nil {
		return nil, buildErrorRetrievingBlock(blockHash, err)
	}
	block := btcutil.NewBlock(msgBlock)

	txs := make([]*btcutil.Tx, len(block.MsgBlock().Transactions))
	for i, t := range block.MsgBlock().Transactions {
		tx := btcutil.NewTx(t)
		txs[i] = tx
	}

	hash, err := chainhash.NewHashFromStr(txHash)

	if err != nil {
		return nil, fmt.Errorf("error parsing hash: %v", err)
	}

	store := blockchain.BuildMerkleTreeStore(txs, false)

	idx := FindInMerkleTreeStore(store, hash)
	if idx == -1 {
		return nil, fmt.Errorf("tx not found in merkle tree: %v", err)
	}

	branch := buildMerkleBranch(store, uint32(len(block.Transactions())), uint32(idx))

	return branch, nil
}

func (btc *BTC) BuildMerkleBranchByEndpoint(txHash string, btcAddress string) (*MerkleBranch, error) {

	btcAdd, err := btcutil.DecodeAddress(btcAddress, &btc.params)
	if err != nil {
		return nil, err
	}

	err = btc.c.ImportAddressRescan(btcAdd.String(), "", false)
	if err != nil {
		return nil, buildErrorImportAddress(btcAddress, err)
	}

	return btc.BuildMerkleBranch(txHash)
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

func FindInMerkleTreeStore(store []*chainhash.Hash, hash *chainhash.Hash) int {
	for i, h := range store {
		if h.IsEqual(hash) {
			return i
		}
	}
	return -1
}

func (btc *BTC) SendBTC(address string, amount uint) (string, error) {

	btcAdd, err := btcutil.DecodeAddress(address, &btc.params)
	if err != nil {
		return "", err
	}

	err = btc.c.ImportAddressRescan(btcAdd.String(), "", false)
	if err != nil {
		return "", err
	}

	hash, err := btc.c.SendToAddress(btcAdd, btcutil.Amount(btcutil.Amount(amount).ToBTC()))

	if err != nil {
		return "", err
	}

	return hash.String(), nil
}

type MerkleBranch struct {
	Hashes []*chainhash.Hash
	Path   int
}

func buildMerkleBranch(merkleTree []*chainhash.Hash, txCount uint32, txIndex uint32) *MerkleBranch {
	hashes := make([]*chainhash.Hash, 0)
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
		hashes = append(hashes, merkleTree[levelOffset+targetOffset])

		levelOffset += levelSize
		currentNodeOffset = currentNodeOffset / 2
		pathIndex++
	}

	return &MerkleBranch{hashes, path}
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
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

func (btc *BTC) validateRedeemScript(fedInfo *FedInfo, script []byte) error {
	addr, err := btcutil.NewAddressScriptHash(script, &btc.params)
	if err != nil {
		return err
	}

	fedAddress, err := btcutil.DecodeAddress(fedInfo.FedAddress, &btc.params)
	if err != nil {
		return err
	}

	if !bytes.Equal(addr.ScriptAddress(), fedAddress.ScriptAddress()) {
		return fmt.Errorf("the generated redeem script does not match with the federation redeem script")
	}
	return nil
}

func (btc *BTC) getRedeemScript(fedInfo *FedInfo, derivationValue []byte) ([]byte, error) {
	var hashBuf *bytes.Buffer

	buf, err := getFlyoverPrefix(derivationValue)
	if err != nil {
		return nil, err
	}

	// All federations activated AFTER Iris will be ERP, therefore we build erp redeem script.
	if fedInfo.ActiveFedBlockHeight < fedInfo.IrisActivationHeight {
		err := btc.buildRedeemScriptBuf(fedInfo, hashBuf, err)
		if err != nil {
			return nil, err
		}
	} else {
		hashBuf, err = btc.getErpRedeemScriptBuf(fedInfo)
		if err != nil {
			return nil, err
		}

		err = btc.validateRedeemScript(fedInfo, hashBuf.Bytes())
		if err != nil { // ok, it could be that ERP is not yet activated, falling back to Redeem Script
			err := btc.buildRedeemScriptBuf(fedInfo, hashBuf, err)
			if err != nil {
				return nil, err
			}
		}
	}

	buf.Write(hashBuf.Bytes())
	return buf.Bytes(), nil
}

func (btc *BTC) buildRedeemScriptBuf(fedInfo *FedInfo, hashBuf *bytes.Buffer, err error) error {
	hashBuf, err = btc.getRedeemScriptBuf(fedInfo, true)
	if err != nil {
		return err
	}

	err = btc.validateRedeemScript(fedInfo, hashBuf.Bytes())
	if err != nil {
		return err
	}

	return nil
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

func (btc *BTC) getRedeemScriptBuf(fedInfo *FedInfo, addMultiSig bool) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	sb := txscript.NewScriptBuilder()
	err := btc.addStdNToMScriptPart(fedInfo, sb)
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

func (btc *BTC) getErpRedeemScriptBuf(fedInfo *FedInfo) (*bytes.Buffer, error) {
	erpRedeemScriptBuf, err := btc.p2ms(fedInfo, false)
	if err != nil {
		return nil, err
	}
	redeemScriptBuf, err := btc.getRedeemScriptBuf(fedInfo, false)
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
	erpRedeemScriptBuffer.Write(redeemScriptBuf.Bytes())
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

func (btc *BTC) p2ms(fedInfo *FedInfo, addMultiSig bool) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	sb := txscript.NewScriptBuilder()
	err := btc.addErpNToMScriptPart(fedInfo, sb)
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

func (btc *BTC) addStdNToMScriptPart(fedInfo *FedInfo, builder *txscript.ScriptBuilder) error {
	builder.AddOp(getOpCodeFromInt(fedInfo.FedThreshold))

	for _, pubKey := range fedInfo.PubKeys {
		pkBuffer, err := hex.DecodeString(pubKey)
		if err != nil {
			return err
		}
		builder.AddData(pkBuffer)
	}

	builder.AddOp(getOpCodeFromInt(fedInfo.FedSize))

	return nil
}

func (btc *BTC) addErpNToMScriptPart(fedInfo *FedInfo, builder *txscript.ScriptBuilder) error {
	size := len(fedInfo.ErpKeys)
	min := size/2 + 1
	builder.AddOp(getOpCodeFromInt(min))

	for _, pubKey := range fedInfo.ErpKeys {
		pkBuffer, err := hex.DecodeString(pubKey)
		if err != nil {
			return err
		}
		builder.AddData(pkBuffer)
	}

	builder.AddOp(getOpCodeFromInt(len(fedInfo.ErpKeys)))

	return nil
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

func isP2PKH(address string) bool {
	pattern := regexp.MustCompile(`^[13][a-km-zA-HJ-NP-Z0-9]{25,34}$`)
	return pattern.MatchString(address)
}

func isP2SH(address string) bool {
	pattern := regexp.MustCompile(`^[32][a-km-zA-HJ-NP-Z0-9]{25,34}$`)
	return pattern.MatchString(address)
}

func isBech32(address string) bool {
	pattern := regexp.MustCompile(`^(bc1|tb1)[a-zA-HJ-NP-Z0-9]{8,87}$`)
	return pattern.MatchString(address)
}
