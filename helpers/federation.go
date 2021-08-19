package federation

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/crypto"
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

func GetDerivedFastBridgeFederationAddressHashString(fedPubKey string, derivationValue []byte, netParams *chaincfg.Params) *btcutil.AddressScriptHash {

	fastBridgeScript, err := buildRedeemScript(fedPubKey, derivationValue) // emulates RSKj BridgeSupport.java#L2607
	if err != nil {
		log.Fatal("There was an error while creating fast bridge redeem script", err)
	}

	addressScriptHash, err := getFastBridgeScriptHash(fastBridgeScript, netParams)
	if err != nil {
		log.Fatal("There was an error while creating fast bridge address script hash", err)
	}

	return addressScriptHash
}

func getFastBridgeScriptHash(script []byte, network *chaincfg.Params) (*btcutil.AddressScriptHash, error) {
	address, err := btcutil.NewAddressScriptHash(script, network)
	if err != nil {
		log.Fatal("There was an error creating an address from the script.", err)
	}
	return address, nil
}

func buildRedeemScript(fedPubKey string, derivationVal []byte) ([]byte, error) {
	// create redeem script
	builder := txscript.NewScriptBuilder()

	// add a random value as DROP
	builder.AddData(derivationVal)
	builder.AddOp(txscript.OP_DROP)

	// TODO: Confirm the script is not multisig and doesn't need all federation pub keys
	//// add the minimum number of needed signatures
	//builder.AddOp(txscript.OP_1)

	wif, err := btcutil.DecodeWIF(fedPubKey)
	if err != nil {
		return nil, nil
	}
	pk := wif.PrivKey.PubKey().SerializeCompressed()
	builder.AddData(pk)

	// TODO: Confirm the script is not multisig and doesn't need all federation pub keys
	//// add the total number of public keys in the multi-sig script
	//builder.AddOp(txscript.OP_1)

	// add the check-sig op-code
	builder.AddOp(txscript.OP_CHECKSIG)
	// redeem script is the script program in the format of []byte
	redeemScript, err := builder.Script()
	if err != nil {
		return nil, err
	}

	// disassemble the script program, so can see its structure
	redeemStr, err := txscript.DisasmString(redeemScript)
	if err != nil {
		return nil, err
	}

	log.Printf("Script: %v", redeemStr)

	return redeemScript, nil
}

// TODO: Confirm the script is not multisig and doesn't need all federation pub keys
//func getAmountOfSignatures(size int) byte {
//	switch size {
//	case 1:
//		return txscript.OP_1
//	case 2:
//		return txscript.OP_2
//	case 3:
//		return txscript.OP_3
//	case 4:
//		return txscript.OP_4
//	case 5:
//		return txscript.OP_5
//	case 6:
//		return txscript.OP_6
//	case 7:
//		return txscript.OP_7
//	case 8:
//		return txscript.OP_8
//	case 9:
//		return txscript.OP_9
//	case 10:
//		return txscript.OP_10
//	case 11:
//		return txscript.OP_11
//	case 12:
//		return txscript.OP_12
//	case 13:
//		return txscript.OP_13
//	case 14:
//		return txscript.OP_14
//	case 15:
//		return txscript.OP_15
//	case 16:
//		return txscript.OP_16
//	default:
//		log.Fatal("Amount of signatures required cannot be calculated. Please review the fedAddr config.")
//	}
//	return 0
//}
//
//func getPublicKeys(addresses []string, netParams *chaincfg.Params) ([]string, error) {
//	for _, address := range addresses {
//		decodedAddress, err := btcutil.DecodeAddress(address, &netParams)
//		if err != nil {
//			return nil, err
//		}
//		return decodedAddress., err
//	}
//	var result []string
//	for i := 0; i < size; i++ {
//		result = append(result, "bla")
//	}
//	return result
//}
