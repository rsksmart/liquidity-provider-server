package federation

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/connectors"
	"log"
)

func GetDerivationValueHash(userBtcRefundAddress []byte, lbcAddress []byte, lpBtcAddress []byte, derivationArgumentsHash []byte) ([]byte, error) {
	var resultData []byte
	resultData = append(resultData, derivationArgumentsHash...)
	resultData = append(resultData, userBtcRefundAddress...)
	resultData = append(resultData, lpBtcAddress...)
	resultData = append(resultData, lbcAddress...)

	derivationValueHash := crypto.Keccak256(resultData)

	return derivationValueHash, nil
}

func GetDerivedFastBridgeFederationAddressHashString(rsk *connectors.RSK, derivationValue []byte, netParams *chaincfg.Params) *btcutil.AddressScriptHash {

	fastBridgeScript, err := rsk.GetRedeemScript()
	if err != nil {
		log.Fatal("There was an error while creating fast bridge redeem script", err)
	}

	modifiedScript := getModifiedRedeemScript(fastBridgeScript, derivationValue)

	addressScriptHash, err := getAddressScriptHash(modifiedScript, netParams)
	if err != nil {
		log.Fatal("There was an error while creating fast bridge address script hash", err)
	}

	return addressScriptHash
}

func getModifiedRedeemScript(script []byte, derivationValue []byte) []byte {
	builder := txscript.NewScriptBuilder()
	builder.AddData(derivationValue)
	builder.AddOp(txscript.OP_DROP)
	builder.AddData(script)
	result, err := builder.Script()
	if err != nil {
		log.Fatal("There was an error while modifying fast bridge redeem script", err)
	}

	// TODO: Verify that the script is generated correctly.
	stringScript, err := txscript.DisasmString(result)
	if err != nil {
		log.Fatal("There was an error while parsing the modified fast bridge redeem script", err)
	}
	log.Printf("%v", stringScript)

	return script
}

func getAddressScriptHash(script []byte, network *chaincfg.Params) (*btcutil.AddressScriptHash, error) {
	// calculate the hash160 of the redeem script

	// TODO: Confirm that this is necessary.
	redeemHash := btcutil.Hash160(script)

	address, err := btcutil.NewAddressScriptHash(redeemHash, network)
	if err != nil {
		log.Fatal("There was an error creating an address from the script.", err)
	}
	return address, nil
}
