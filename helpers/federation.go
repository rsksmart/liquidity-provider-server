package federation

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
)

type FedInfo struct {
	FedSize              int
	FedThreshold         int
	PubKeys              []string
	FedAddress           []byte
	ActiveFedBlockHeight int
	IrisActivationHeight int
	ErpKeys              []string
}

func GetDerivationValueHash(userBtcRefundAddress []byte, lbcAddress []byte, lpBtcAddress []byte, derivationArgumentsHash []byte) ([]byte, error) {
	var resultData []byte
	resultData = append(resultData, derivationArgumentsHash...)
	resultData = append(resultData, userBtcRefundAddress...)
	resultData = append(resultData, lpBtcAddress...)
	resultData = append(resultData, lbcAddress...)

	derivationValueHash := crypto.Keccak256(resultData)

	return derivationValueHash, nil
}

func GetDerivedBitcoinAddressHash(derivationValue []byte, fedInfo *FedInfo, netParams *chaincfg.Params) (*btcutil.AddressScriptHash, error) {

	err := ensureRedeemScriptIsValid(fedInfo, netParams)
	if err != nil {
		return nil, err
	}

	flyoverScriptBuf, err := GetRedeemScriptBuffer(fedInfo, derivationValue, netParams)
	if err != nil {
		return nil, err
	}

	scriptString, err := txscript.DisasmString(flyoverScriptBuf.Bytes())
	if err != nil {
		return nil, err
	}
	log.Debug(scriptString)

	addressScriptHash, err := getAddressScriptHash(flyoverScriptBuf.Bytes(), netParams)
	if err != nil {
		return nil, err
	}

	return addressScriptHash, nil
}

func ensureRedeemScriptIsValid(info *FedInfo, params *chaincfg.Params) error {
	buf, err := GetRedeemScriptBufferWithoutPrefix(info)
	if err != nil {
		return err
	}

	script := buf.Bytes()
	scriptString, err := txscript.DisasmString(script)
	if err != nil {
		return err
	}
	log.Debug(scriptString)
	addr, err := btcutil.NewAddressScriptHash(script, params)
	if err != nil {
		return err
	}

	if bytes.Compare(addr.ScriptAddress(), info.FedAddress) != 0 {
		return fmt.Errorf("the generated redeem script does not match with the federation redeem script")
	}

	return nil
}

func GetRedeemScriptBuffer(info *FedInfo, derivationValue []byte, params *chaincfg.Params) (bytes.Buffer, error) {
	var buf bytes.Buffer
	// All federations activated AFTER Iris will be ERP, therefore we build erp redeem script.
	// TODO: Verify if bridge method that retrieves ActiveFedBlockHeight is giving correct results (0)
	if info.ActiveFedBlockHeight < info.IrisActivationHeight {
		sb, err := getFlyoverRedeemScriptBuf(info, getDerivationHashString(derivationValue))
		if err != nil {
			return bytes.Buffer{}, err
		}
		buf = *sb
	} else {
		sb, err := getFlyoverErpRedeemScriptBuf(info, getDerivationHashString(derivationValue), params)
		if err != nil {
			return bytes.Buffer{}, err
		}
		buf = *sb
	}
	return buf, nil
}

func GetRedeemScriptBufferWithoutPrefix(info *FedInfo) (bytes.Buffer, error) {
	var buf bytes.Buffer

	// TODO: verify whether we must check the erp fed activation to prevent comparing ERP vs. powPeg script.
	sb, err := getPowPegRedeemScriptBuf(info, true)
	if err != nil {
		return bytes.Buffer{}, err
	}
	buf = *sb

	return buf, nil
}

func getDerivationHashString(derivationValue []byte) string {
	return hex.EncodeToString(derivationValue)
}

func getFlyoverRedeemScriptBuf(info *FedInfo, hash string) (*bytes.Buffer, error) {
	buf, err := getFlyoverPrefix(hash)
	if err != nil {
		return nil, err
	}
	hashBuf, err := getPowPegRedeemScriptBuf(info, true)
	if err != nil {
		return nil, err
	}

	buf.Write(hashBuf.Bytes())

	return buf, nil
}

func getFlyoverErpRedeemScriptBuf(info *FedInfo, hash string, params *chaincfg.Params) (*bytes.Buffer, error) {
	buf, err := getFlyoverPrefix(hash)
	if err != nil {
		return nil, err
	}
	hashBuf, err := getErpRedeemScriptBuf(info, params)
	if err != nil {
		return nil, err
	}

	buf.Write(hashBuf.Bytes())

	return buf, nil
}

func getFlyoverPrefix(hash string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	hashLength, err := hex.DecodeString("20")
	if err != nil {
		return nil, err
	}
	encodedHash, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}
	buf.Write(hashLength)
	buf.Write(encodedHash)
	buf.WriteByte(txscript.OP_DROP)

	return &buf, nil
}

func getPowPegRedeemScriptBuf(info *FedInfo, addMultiSig bool) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	sb := txscript.NewScriptBuilder()
	err := addStdNToMScriptPart(sb, info)
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

func getErpRedeemScriptBuf(info *FedInfo, params *chaincfg.Params) (*bytes.Buffer, error) {

	erpRedeemScriptBuf, err := p2ms(info, false)
	if err != nil {
		return nil, err
	}
	powPegRedeemScriptBuf, err := getPowPegRedeemScriptBuf(info, false)

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

	csv, err := hex.DecodeString(getCsvValueFromNetwork(params))
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

func p2ms(info *FedInfo, addMultiSig bool) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	sb := txscript.NewScriptBuilder()
	err := addErpNToMScriptPart(sb, info)
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

func getCsvValueFromNetwork(params *chaincfg.Params) string {
	switch params.Name {
	case chaincfg.MainNetParams.Name:
		return "CD50"
	case chaincfg.TestNet3Params.Name:
		return "CD50"
	default: // RegTest
		return "01F4"
	}
}

func addStdNToMScriptPart(builder *txscript.ScriptBuilder, fedInfo *FedInfo) error {
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

func addErpNToMScriptPart(builder *txscript.ScriptBuilder, fedInfo *FedInfo) error {
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

func getAddressScriptHash(script []byte, network *chaincfg.Params) (*btcutil.AddressScriptHash, error) {
	// calculate the hash160 of the redeem script
	// TODO: Confirm that this is necessary.
	redeemHash := btcutil.Hash160(script)

	address, err := btcutil.NewAddressScriptHash(redeemHash, network)
	if err != nil {
		return nil, err
	}
	return address, nil
}
