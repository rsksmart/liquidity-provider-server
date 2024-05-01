package btcclient

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

func init() {
	btcjson.MustRegisterCmd("signrawtransactionwithkey", (*SignRawTransactionWithKeyCmd)(nil), btcjson.UsageFlag(0))
}

type RpcRequestParamsObject[T any] struct {
	Jsonrpc btcjson.RPCVersion `json:"jsonrpc"`
	Method  string             `json:"method"`
	Params  T                  `json:"params"`
	ID      interface{}        `json:"id"`
}

type RpcWallet interface {
	WalletCreateFundedPsbt(inputs []btcjson.PsbtInput, outputs []btcjson.PsbtOutput, locktime *uint32, options *btcjson.WalletCreateFundedPsbtOpts, bip32Derivs *bool) (*btcjson.WalletCreateFundedPsbtResult, error)
	ListUnspent() ([]btcjson.ListUnspentResult, error)
	CreateRawTransaction(inputs []btcjson.TransactionInput, amounts map[btcutil.Address]btcutil.Amount, lockTime *int64) (*wire.MsgTx, error)
	FundRawTransaction(tx *wire.MsgTx, opts btcjson.FundRawTransactionOpts, isWitness *bool) (*btcjson.FundRawTransactionResult, error)
	SignRawTransactionWithWallet(tx *wire.MsgTx) (*wire.MsgTx, bool, error)
	SendRawTransaction(tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error)
	GetWalletInfo() (*btcjson.GetWalletInfoResult, error)
	WalletPassphrase(passphrase string, timeoutSecs int64) error
	ImportAddressRescan(address string, account string, rescan bool) error
	ListUnspentMinMaxAddresses(minConf int, maxConf int, addrs []btcutil.Address) ([]btcjson.ListUnspentResult, error)
	GetTransaction(txHash *chainhash.Hash) (*btcjson.GetTransactionResult, error)
	GetAddressInfo(address string) (*btcjson.GetAddressInfoResult, error)
	ImportPubKeyRescan(pubKey string, rescan bool) error
	ImportPubKey(pubKey string) error
}

type RpcClient interface {
	RpcWallet
	SendCmd(cmd interface{}) chan *rpcclient.Response
	NextID() uint64
	Ping() error
	Disconnect()
	GetRawTransactionVerbose(txHash *chainhash.Hash) (*btcjson.TxRawResult, error)
	GetRawTransaction(txHash *chainhash.Hash) (*btcutil.Tx, error)
	GetBlockChainInfo() (*btcjson.GetBlockChainInfoResult, error)
	GetBlockVerbose(blockHash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error)
	GetBlock(blockHash *chainhash.Hash) (*wire.MsgBlock, error)
	CreateWallet(name string, opts ...rpcclient.CreateWalletOpt) (*btcjson.CreateWalletResult, error)
	LoadWallet(walletName string) (*btcjson.LoadWalletResult, error)
	EstimateSmartFee(confTarget int64, mode *btcjson.EstimateSmartFeeMode) (*btcjson.EstimateSmartFeeResult, error)
}

type ClientAdapter interface {
	RpcClient
	RpcWallet
	SignRawTransactionWithKey(tx *wire.MsgTx, privateKeysWIFs []string) (*wire.MsgTx, bool, error)
	CreateReadonlyWallet(bodyParams ReadonlyWalletRequest) error
}
