package bitcoin

import (
	"crypto/sha256"
	"errors"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

func ScriptToP2shP2wsh(script []byte) ([]byte, error) {
	if len(script) == 0 {
		return nil, errors.New("script cannot be empty")
	}

	witnessScriptHash := sha256.Sum256(script)
	segwitScript, err := txscript.NewScriptBuilder().
		AddOp(txscript.OP_0).
		AddData(witnessScriptHash[:]).
		Script()
	if err != nil {
		return nil, err
	}

	return segwitScript, nil
}

func ScriptToAddressP2shP2wsh(script []byte, btcParams *chaincfg.Params) (*btcutil.AddressScriptHash, error) {
	if btcParams == nil {
		return nil, errors.New("bitcoin network parameters cannot be nil")
	}
	segwitScript, err := ScriptToP2shP2wsh(script)
	if err != nil {
		return nil, err
	}

	return btcutil.NewAddressScriptHash(segwitScript, btcParams)
}
