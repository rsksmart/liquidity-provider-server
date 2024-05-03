package btcclient

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

type SignRawTransactionWithKeyCmd struct {
	RawTx   string
	WifKeys []string
}

type SignRawTransactionWithKeysRequest struct {
	Transaction  string   `json:"hexstring"`
	PrivKeysWIFs []string `json:"privkeys"`
}

type FutureSignRawTransactionWithKeyResult chan *rpcclient.Response

func (r FutureSignRawTransactionWithKeyResult) Receive() (*wire.MsgTx, bool, error) {
	res, err := rpcclient.ReceiveFuture(r)
	if err != nil {
		return nil, false, err
	}

	// Unmarshal as a signtransactionwithwallet since it has the same response as signrawtransactionwithkey
	// and this struct is already in the library
	var signRawTxWithWalletResult btcjson.SignRawTransactionWithWalletResult
	err = json.Unmarshal(res, &signRawTxWithWalletResult)
	if err != nil {
		return nil, false, err
	}

	serializedTx, err := hex.DecodeString(signRawTxWithWalletResult.Hex)
	if err != nil {
		return nil, false, err
	}

	var msgTx wire.MsgTx
	if witnessErr := msgTx.Deserialize(bytes.NewReader(serializedTx)); witnessErr != nil {
		if legacyErr := msgTx.DeserializeNoWitness(bytes.NewReader(serializedTx)); legacyErr != nil {
			return nil, false, legacyErr
		}
	}

	return &msgTx, signRawTxWithWalletResult.Complete, nil
}
