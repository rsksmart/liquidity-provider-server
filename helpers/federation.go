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
	FedSize      int
	FedThreshold int
	PubKeys      []string
	FedAddress   []byte
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

func GetDerivedFastBridgeFederationAddressHash(derivationValue []byte, fedInfo *FedInfo, netParams *chaincfg.Params) (*btcutil.AddressScriptHash, error) {

	testScript, err := buildRedeemScript(fedInfo, nil)
	if err != nil {
		return nil, err
	}

	newAddr, err := btcutil.NewAddressScriptHash(testScript, netParams)
	if err != nil {
		return nil, err
	}

	if bytes.Compare(newAddr.ScriptAddress(), fedInfo.FedAddress) != 0 {
		return nil, fmt.Errorf("the generated redeem script does not match with the federation redeem script")
	}
	flyoverScript, err := buildFlyOverScript(fedInfo, derivationValue)
	if err != nil {
		return nil, err
	}

	addressScriptHash, err := getAddressScriptHash(flyoverScript, netParams)
	if err != nil {
		return nil, err
	}

	return addressScriptHash, nil
}

func buildRedeemScript(fedInfo *FedInfo, scriptBuilder *txscript.ScriptBuilder) ([]byte, error) {
	var builder = txscript.NewScriptBuilder()
	if scriptBuilder != nil {
		builder = scriptBuilder
	}

	err := addStdLogicToScript(builder, fedInfo)
	if err != nil {
		return nil, err
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

	return builder.Script()
}

func addStdLogicToScript(builder *txscript.ScriptBuilder, fedInfo *FedInfo) error {
	builder.AddOp(getOpCodeFromInt(fedInfo.FedThreshold))

	for _, pubKey := range fedInfo.PubKeys {
		pkBuffer, err := hex.DecodeString(pubKey)
		if err != nil {
			return err
		}
		builder.AddData(pkBuffer)
	}

	builder.AddOp(getOpCodeFromInt(fedInfo.FedSize))
	builder.AddOp(txscript.OP_CHECKMULTISIG)

	return nil
}

func buildFlyOverScript(fedInfo *FedInfo, derivationValue []byte) ([]byte, error) {
	builder := txscript.NewScriptBuilder()

	// add
	builder.AddData(derivationValue)
	builder.AddOp(txscript.OP_DROP)

	// TODO: check if a simple concat of both script parts ([]byte) would work so we can remove this line.
	err := addStdLogicToScript(builder, fedInfo)
	if err != nil {
		return nil, err
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
