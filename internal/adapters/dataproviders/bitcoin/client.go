package bitcoin

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

type rpcWallet interface {
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
}

type rpcClient interface {
	rpcWallet
	Ping() error
	Disconnect()
	GetRawTransactionVerbose(txHash *chainhash.Hash) (*btcjson.TxRawResult, error)
	GetRawTransaction(txHash *chainhash.Hash) (*btcutil.Tx, error)
	GetBlockChainInfo() (*btcjson.GetBlockChainInfoResult, error)
	GetBlockVerbose(blockHash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error)
	GetBlock(blockHash *chainhash.Hash) (*wire.MsgBlock, error)
}
