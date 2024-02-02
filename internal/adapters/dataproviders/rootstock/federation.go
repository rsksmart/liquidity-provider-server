package rootstock

import (
	"bytes"
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

func getFedRedeemScript(fedInfo blockchain.FederationInfo, btcParams chaincfg.Params) ([]byte, error) {
	var buf *bytes.Buffer
	var err error

	// All Federations activated AFTER Iris will be ERP, therefore we build redeem script.
	if fedInfo.ActiveFedBlockHeight > fedInfo.IrisActivationHeight {
		buf, err = getFedRedeemScriptAfterIrisActivation(fedInfo, btcParams)
	} else {
		buf, err = getFedRedeemScriptBeforeIrisActivation(fedInfo, btcParams)
	}

	return buf.Bytes(), err
}

func getFedRedeemScriptAfterIrisActivation(fedInfo blockchain.FederationInfo, btcParams chaincfg.Params) (*bytes.Buffer, error) {
	buf, err := getRedeemScriptBuf(fedInfo, true)
	if err != nil {
		return nil, err
	}

	err = validateRedeemScript(fedInfo, btcParams, buf.Bytes())
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func getFedRedeemScriptBeforeIrisActivation(fedInfo blockchain.FederationInfo, btcParams chaincfg.Params) (*bytes.Buffer, error) {
	buf, err := getErpRedeemScriptBuf(fedInfo, btcParams)
	if err != nil {
		return nil, err
	}

	err = validateRedeemScript(fedInfo, btcParams, buf.Bytes())
	if err != nil { // ok, it could be that ERP is not yet activated, falling back to redeem Script
		buf, err = getRedeemScriptBuf(fedInfo, true)
		if err != nil {
			return nil, err
		}

		err = validateRedeemScript(fedInfo, btcParams, buf.Bytes())
		if err != nil {
			return nil, err
		}
	}
	return buf, nil
}

func getRedeemScriptBuf(fedInfo blockchain.FederationInfo, addMultiSig bool) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	sb := txscript.NewScriptBuilder()
	err := addStdNToMScriptPart(fedInfo, sb)
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

func addStdNToMScriptPart(fedInfo blockchain.FederationInfo, builder *txscript.ScriptBuilder) error {
	builder.AddOp(getOpCodeFromInt(int(fedInfo.FedThreshold)))

	for _, pubKey := range fedInfo.PubKeys {
		pkBuffer, err := hex.DecodeString(pubKey)
		if err != nil {
			return err
		}
		builder.AddData(pkBuffer)
	}

	builder.AddOp(getOpCodeFromInt(int(fedInfo.FedSize)))
	return nil
}

func validateRedeemScript(fedInfo blockchain.FederationInfo, btcParams chaincfg.Params, script []byte) error {
	addr, err := btcutil.NewAddressScriptHash(script, &btcParams)
	if err != nil {
		return err
	}

	fedAddress, err := btcutil.DecodeAddress(fedInfo.FedAddress, &btcParams)
	if err != nil {
		return err
	}
	if !bytes.Equal(addr.ScriptAddress(), fedAddress.ScriptAddress()) {
		return errors.New("the generated redeem script does not match with the federation redeem script")
	}

	return nil
}

func getErpRedeemScriptBuf(fedInfo blockchain.FederationInfo, btcParams chaincfg.Params) (*bytes.Buffer, error) {
	erpRedeemScriptBuf, err := p2ms(fedInfo, false)
	if err != nil {
		return nil, err
	}

	redeemScriptBuf, err := getRedeemScriptBuf(fedInfo, false)
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

func p2ms(fedInfo blockchain.FederationInfo, addMultiSig bool) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	sb := txscript.NewScriptBuilder()
	err := addErpNToMScriptPart(fedInfo, sb)
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

func addErpNToMScriptPart(fedInfo blockchain.FederationInfo, builder *txscript.ScriptBuilder) error {
	size := len(fedInfo.ErpKeys)
	minimum := size/2 + 1
	builder.AddOp(getOpCodeFromInt(minimum))

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

func getFlyoverRedeemScript(derivationValue []byte, fedRedeemScript []byte) []byte {
	var buf bytes.Buffer
	hashPrefix, _ := hex.DecodeString("20")
	buf.Write(hashPrefix)
	buf.Write(derivationValue)
	buf.WriteByte(txscript.OP_DROP)
	buf.Write(fedRedeemScript)
	return buf.Bytes()
}

// nolint:cyclop
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