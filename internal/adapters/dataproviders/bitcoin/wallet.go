package bitcoin

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	WalletUnlockSeconds = 10
)

type bitcoindWallet struct {
	conn           *Connection
	fixedTxFeeRate float64
	isEncrypted    bool
	password       string
}

func NewBitcoindWallet(
	conn *Connection,
	fixedTxFeeRate float64,
	isEncrypted bool,
	password string,
) blockchain.BitcoinWallet {
	return &bitcoindWallet{
		conn:           conn,
		fixedTxFeeRate: fixedTxFeeRate,
		isEncrypted:    isEncrypted,
		password:       password,
	}
}

func (wallet *bitcoindWallet) EstimateTxFees(toAddress string, value *entities.Wei) (*entities.Wei, error) {
	_, err := btcutil.DecodeAddress(toAddress, wallet.conn.NetworkParams)
	if err != nil {
		return nil, err
	}

	amountInSatoshi, _ := value.ToSatoshi().Float64()
	output := []btcjson.PsbtOutput{
		{toAddress: float64(amountInSatoshi) / BtcToSatoshi},
		{"data": hex.EncodeToString(make([]byte, 32))}, // quote hash output
	}

	var changePosition int64 = 2
	feeRate := wallet.fixedTxFeeRate
	opts := btcjson.WalletCreateFundedPsbtOpts{
		ChangePosition: &changePosition,
		FeeRate:        &feeRate,
	}

	simulatedTx, err := wallet.conn.client.WalletCreateFundedPsbt(nil, output, nil, &opts, nil)
	if err != nil {
		return nil, err
	}
	return entities.SatoshiToWei(uint64(simulatedTx.Fee * BtcToSatoshi)), nil
}

func (wallet *bitcoindWallet) GetBalance() (*entities.Wei, error) {
	balance := new(entities.Wei)
	utxos, err := wallet.conn.client.ListUnspent()
	if err != nil {
		return nil, err
	}

	for _, utxo := range utxos {
		if utxo.Spendable {
			balance.Add(balance, entities.SatoshiToWei(uint64(utxo.Amount*BtcToSatoshi)))
		}
	}
	return balance, nil
}

func (wallet *bitcoindWallet) SendWithOpReturn(address string, value *entities.Wei, opReturnContent []byte) (string, error) {
	decodedAddress, err := btcutil.DecodeAddress(address, wallet.conn.NetworkParams)
	if err != nil {
		return "", err
	}

	satoshis, _ := value.ToSatoshi().Float64()
	output := map[btcutil.Address]btcutil.Amount{
		decodedAddress: btcutil.Amount(satoshis),
	}
	rawTx, err := wallet.conn.client.CreateRawTransaction(nil, output, nil)

	opReturnScript, err := txscript.NullDataScript(opReturnContent)
	if err != nil {
		return "", err
	}
	rawTx.AddTxOut(wire.NewTxOut(0, opReturnScript))

	changePosition := 2
	feeRate := wallet.fixedTxFeeRate
	opts := btcjson.FundRawTransactionOpts{
		ChangePosition: &changePosition,
		FeeRate:        &feeRate,
	}

	if wallet.isEncrypted {
		if err = wallet.unlock(); err != nil {
			return "", err
		}
	}

	log.Infof("Sending %v BTC to %s\n", value.ToRbtc(), address)
	fundedTx, err := wallet.conn.client.FundRawTransaction(rawTx, opts, nil)
	if err != nil {
		return "", err
	}
	signedTx, _, err := wallet.conn.client.SignRawTransactionWithWallet(fundedTx.Transaction)
	if err != nil {
		return "", err
	}
	txHash, err := wallet.conn.client.SendRawTransaction(signedTx, false)
	if err != nil {
		return "", err
	}
	return txHash.String(), nil
}

func (wallet *bitcoindWallet) unlock() error {
	info, err := wallet.conn.client.GetWalletInfo()
	if err != nil {
		return err
	}
	if info.UnlockedUntil != nil && time.Until(time.Unix(int64(*info.UnlockedUntil), 0)) > 0 {
		log.Debug("Wallet already unlocked")
		return nil
	}
	return wallet.conn.client.WalletPassphrase(wallet.password, WalletUnlockSeconds)
}

func (wallet *bitcoindWallet) ImportAddress(address string) error {
	_, err := btcutil.DecodeAddress(address, wallet.conn.NetworkParams)
	if err != nil {
		return err
	}
	return wallet.conn.client.ImportAddressRescan(address, "", false)
}

func (wallet *bitcoindWallet) GetTransactions(address string) ([]blockchain.BitcoinTransactionInformation, error) {
	var result []blockchain.BitcoinTransactionInformation
	var ok bool
	var tx blockchain.BitcoinTransactionInformation
	parsedAddress, err := btcutil.DecodeAddress(address, wallet.conn.NetworkParams)
	if err != nil {
		return nil, err
	}
	utxos, err := wallet.conn.client.ListUnspentMinMaxAddresses(0, 9999, []btcutil.Address{parsedAddress})

	txs := make(map[string]blockchain.BitcoinTransactionInformation)
	for _, utxo := range utxos {
		tx, ok = txs[utxo.TxID]
		if !ok {
			tx = blockchain.BitcoinTransactionInformation{
				Hash:          utxo.TxID,
				Confirmations: uint64(utxo.Confirmations),
				Outputs:       make(map[string][]*entities.Wei),
			}
			txs[utxo.TxID] = tx
		}
		if _, ok = tx.Outputs[address]; !ok {
			tx.Outputs[address] = make([]*entities.Wei, 0)
		}
		tx.Outputs[utxo.Address] = append(tx.Outputs[utxo.Address], entities.SatoshiToWei(uint64(utxo.Amount*BtcToSatoshi)))
	}

	for key, value := range txs {
		result = append(result, value)
		delete(txs, key)
	}
	return result, nil
}

func (wallet *bitcoindWallet) Unlock() error {
	return wallet.unlock()
}
