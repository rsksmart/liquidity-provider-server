package federation

import (
	"bytes"
	"encoding/binary"
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

	ensureRedeemScriptIsValid(fedInfo, derivationValue, netParams)

	flyoverScript, err := buildFlyOverRedeemScript(fedInfo, derivationValue, netParams, true)
	if err != nil {
		return nil, err
	}

	addressScriptHash, err := getAddressScriptHash(flyoverScript, netParams)
	if err != nil {
		return nil, err
	}

	return addressScriptHash, nil
}

func ensureRedeemScriptIsValid(info *FedInfo, derivationValue []byte, params *chaincfg.Params) error {
	newAddr, err := getStdRedeemScriptAddressWithoutPrefix(info, derivationValue, params)
	if err != nil {
		return err
	}

	if bytes.Compare(newAddr, info.FedAddress) != 0 {
		return fmt.Errorf("the generated redeem script does not match with the federation redeem script")
	}

	return nil
}

func getStdRedeemScriptAddressWithoutPrefix(fedInfo *FedInfo, derivationValue []byte, netParams *chaincfg.Params) ([]byte, error) {
	script, err := buildFlyOverRedeemScript(fedInfo, derivationValue, netParams, false)
	if err != nil {
		return nil, err
	}

	scriptString, err := txscript.DisasmString(script)
	if err != nil {
		return nil, err
	}
	log.Debug(scriptString)

	addr, err := btcutil.NewAddressScriptHash(script, netParams)
	if err != nil {
		return nil, err
	}
	return addr.ScriptAddress(), nil
}

func buildFlyOverRedeemScript(fedInfo *FedInfo, derivationValue []byte, netParams *chaincfg.Params, addFlyOverPrefix bool) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	// All federations activated AFTER Iris will be ERP, therefore we build erp redeem script.
	if fedInfo.ActiveFedBlockHeight < fedInfo.IrisActivationHeight {
		err := buildFlyOverStdRedeemScript(fedInfo, derivationValue, builder, addFlyOverPrefix)
		if err != nil {
			return nil, err
		}
	} else {
		var err = buildFlyOverErpRedeemScript(fedInfo, derivationValue, builder, netParams)
		if err != nil {
			return nil, err
		}
	}

	result, err := builder.Script()
	if err != nil {
		return nil, err
	}
	scriptString, err := txscript.DisasmString(result)
	if err != nil {
		return nil, err
	}

	log.Debug(scriptString)
	return result, nil
}

func buildFlyOverStdRedeemScript(fedInfo *FedInfo, derivationValue []byte, builder *txscript.ScriptBuilder, addFlyOverPrefix bool) error {

	if addFlyOverPrefix {
		addFlyOverPrefixHash(builder, derivationValue)
	}

	err := addStdNToMScriptPart(builder, fedInfo)
	if err != nil {
		return err
	}
	builder.AddOp(txscript.OP_CHECKMULTISIG)

	return nil
}

func buildFlyOverErpRedeemScript(fedInfo *FedInfo, derivationValue []byte, builder *txscript.ScriptBuilder, netParams *chaincfg.Params) error {
	addFlyOverPrefixHash(builder, derivationValue)
	builder.AddOp(txscript.OP_NOTIF)

	err := addStdNToMScriptPart(builder, fedInfo)
	if err != nil {
		return err
	}
	builder.AddOp(txscript.OP_ELSE)

	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, uint16(2))
	binary.LittleEndian.PutUint16(bs, uint16(0))
	builder.AddData(bs)

	//csv := getCsvValueFromNetwork(netParams)
	//if err != nil {
	//	return err
	//}
	//builder.AddData(csv)

	builder.AddOp(txscript.OP_CHECKSEQUENCEVERIFY)
	builder.AddOp(txscript.OP_DROP)

	err = addErpNToMScriptPart(builder, fedInfo)
	if err != nil {
		return err
	}

	builder.AddOp(txscript.OP_ENDIF)
	builder.AddOp(txscript.OP_CHECKMULTISIG)

	return nil
}

func getPowPegRedeemScriptBuf(info *FedInfo, addMultiSig bool) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	sb := txscript.NewScriptBuilder()
	addStdNToMScriptPart(sb, info)
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

func getErpRedeemScriptBuf(info *FedInfo) (*bytes.Buffer, error) {

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
	bytes, err := hex.DecodeString("02")
	if err != nil {
		return nil, err
	}
	erpRedeemScriptBuffer.Write(bytes)

	csv, err := hex.DecodeString(getCsvValueFromNetwork(&chaincfg.MainNetParams))
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
	addErpNToMScriptPart(sb, info)
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
	default: // regtest
		return "01F4"
	}
}

func addFlyOverPrefixHash(builder *txscript.ScriptBuilder, derivationValue []byte) {
	builder.AddData(derivationValue)
	builder.AddOp(txscript.OP_DROP)
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
