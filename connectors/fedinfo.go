package connectors

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type FedInfo struct {
	FedSize              int
	FedThreshold         int
	PubKeys              []string
	FedAddress           string
	ActiveFedBlockHeight int
	IrisActivationHeight int
	ErpKeys              []string
}

func (fedInfo *FedInfo) getFedRedeemScript(btcParams chaincfg.Params) ([]byte, error) {
	var buf *bytes.Buffer
	var err error

	// All Federations activated AFTER Iris will be ERP, therefore we build redeem script.
	if fedInfo.ActiveFedBlockHeight > fedInfo.IrisActivationHeight {
		buf, err = fedInfo.getPowPegRedeemScriptBuf(true)
		if err != nil {
			return nil, err
		}

		err = fedInfo.validateRedeemScript(btcParams, buf.Bytes())
		if err != nil {
			return nil, err
		}
	} else {
		buf, err = fedInfo.getErpRedeemScriptBuf(btcParams)
		if err != nil {
			return nil, err
		}

		err = fedInfo.validateRedeemScript(btcParams, buf.Bytes())
		if err != nil { // ok, it could be that ERP is not yet activated, falling back to PowPeg redeem Script
			buf, err = fedInfo.getPowPegRedeemScriptBuf(true)
			if err != nil {
				return nil, err
			}

			err = fedInfo.validateRedeemScript(btcParams, buf.Bytes())
			if err != nil {
				return nil, err
			}
		}
	}

	return buf.Bytes(), nil
}

func (fedInfo *FedInfo) getFlyoverRedeemScript(derivationValue []byte, fedRedeemScript []byte) ([]byte, error) {
	buf, err := getFlyoverPrefix(derivationValue)
	if err != nil {
		return nil, err
	}

	buf.Write(fedRedeemScript)
	return buf.Bytes(), nil
}

func (fedInfo *FedInfo) validateRedeemScript(btcParams chaincfg.Params, script []byte) error {
	addr, err := btcutil.NewAddressScriptHash(script, &btcParams)
	if err != nil {
		return err
	}

	fedAddress, err := btcutil.DecodeAddress(fedInfo.FedAddress, &btcParams)
	if err != nil {
		return err
	}

	if !bytes.Equal(addr.ScriptAddress(), fedAddress.ScriptAddress()) {
		return fmt.Errorf("The generated redeem script does not match with the federation redeem script")
	}

	return nil
}

func (fedInfo *FedInfo) getPowPegRedeemScriptBuf(addMultiSig bool) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	sb := txscript.NewScriptBuilder()
	err := fedInfo.addStdNToMScriptPart(sb)
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

func (fedInfo *FedInfo) getErpRedeemScriptBuf(btcParams chaincfg.Params) (*bytes.Buffer, error) {
	erpRedeemScriptBuf, err := fedInfo.p2ms(false)
	if err != nil {
		return nil, err
	}

	powpegRedeemScriptBuf, err := fedInfo.getPowPegRedeemScriptBuf(false)
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
	erpRedeemScriptBuffer.Write(powpegRedeemScriptBuf.Bytes())
	erpRedeemScriptBuffer.WriteByte(txscript.OP_ELSE)
	byteArr, err := hex.DecodeString("02")
	if err != nil {
		return nil, err
	}
	erpRedeemScriptBuffer.Write(byteArr)

	csv, err := hex.DecodeString(getCsvValueFromNetwork(btcParams))
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

func (fedInfo *FedInfo) p2ms(addMultiSig bool) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	sb := txscript.NewScriptBuilder()
	err := fedInfo.addErpNToMScriptPart(sb)
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

func (fedInfo *FedInfo) addStdNToMScriptPart(builder *txscript.ScriptBuilder) error {
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

func (fedInfo *FedInfo) addErpNToMScriptPart(builder *txscript.ScriptBuilder) error {
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

func getDerivationValueHash(userBtcRefundAddr []byte, lbcAddress []byte, lpBtcAddress []byte, derivationArgumentsHash []byte) ([]byte, error) {
	var buf bytes.Buffer
	buf.Write(derivationArgumentsHash)
	buf.Write(userBtcRefundAddr)
	buf.Write(lbcAddress)
	buf.Write(lpBtcAddress)

	derivationValueHash := crypto.Keccak256(buf.Bytes())

	return derivationValueHash, nil
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

func getCsvValueFromNetwork(btcParams chaincfg.Params) string {
	switch btcParams.Name {
	case chaincfg.MainNetParams.Name:
		return "CD50"
	case chaincfg.TestNet3Params.Name:
		return "CD50"
	default:
		return "01F4"
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
