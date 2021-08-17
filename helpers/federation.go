package federation

import (
	"github.com/ethereum/go-ethereum/crypto"
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

func GetDerivedFastBridgeFederationAddressHash(derivationValue []byte) []byte {
	return derivationValue // TODO: implement derivation
}