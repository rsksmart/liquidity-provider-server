package federation

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"
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

func GetDerivationValueHash(userBtcRefundAddr []byte, lbcAddress []byte, lpBtcAddress []byte, derivationArgumentsHash []byte) ([]byte, error) {
	var buf bytes.Buffer
	buf.Write(derivationArgumentsHash)
	buf.Write(userBtcRefundAddr)
	buf.Write(lbcAddress)
	buf.Write(lpBtcAddress)

	derivationValueHash := crypto.Keccak256(buf.Bytes())

	return derivationValueHash, nil
}

func GetBytesFromBtcAddress(encoded string) ([]byte, error) {
	addressBts, ver, err := base58.CheckDecode(encoded)
	if err != nil {
		return nil, err
	}
	var bts bytes.Buffer
	bts.WriteByte(ver)
	bts.Write(addressBts)

	return bts.Bytes(), nil
}

func GetDerivedBitcoinAddressHash(derivationValue []byte, fedInfo *FedInfo, netParams *chaincfg.Params) (*btcutil.AddressScriptHash, error) {
	flyoverScript, err := GetRedeemScript(fedInfo, derivationValue, netParams)
	if err != nil {
		return nil, err
	}
	addressScriptHash, err := btcutil.NewAddressScriptHash(flyoverScript, netParams)
	if err != nil {
		return nil, err
	}

	return addressScriptHash, nil
}

func validateRedeemScript(script []byte, expectedAddress []byte, params *chaincfg.Params) error {

	addr, err := btcutil.NewAddressScriptHash(script, params)
	if err != nil {
		return err
	}

	if !bytes.Equal(addr.ScriptAddress(), expectedAddress) {
		return fmt.Errorf("the generated redeem script does not match with the federation redeem script")
	}

	return nil
}

func GetRedeemScript(info *FedInfo, derivationValue []byte, params *chaincfg.Params) ([]byte, error) {
	var hashBuf *bytes.Buffer

	buf, err := getFlyoverPrefix(derivationValue)
	if err != nil {
		return nil, err
	}

	// All federations activated AFTER Iris will be ERP, therefore we build erp redeem script.
	if info.ActiveFedBlockHeight < info.IrisActivationHeight {
		hashBuf, err = getPowPegRedeemScriptBuf(info, true)
		if err != nil {
			return nil, err
		}
	} else {
		hashBuf, err = getErpRedeemScriptBuf(info, params)
		if err != nil {
			return nil, err
		}
	}

	err = validateRedeemScript(hashBuf.Bytes(), info.FedAddress, params)
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
