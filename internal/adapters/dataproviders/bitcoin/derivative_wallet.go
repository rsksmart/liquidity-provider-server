package bitcoin

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin/btcclient"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/account"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	log "github.com/sirupsen/logrus"
)

const (
	changePosition     = 2
	DerivativeWalletId = "rsk-wallet"
)

type DerivativeWallet struct {
	conn       *Connection
	rskAccount *account.RskAccount
}

func NewDerivativeWallet(
	conn *Connection,
	rskAccount *account.RskAccount,
) (blockchain.BitcoinWallet, error) {
	if conn.WalletId != DerivativeWalletId {
		return nil, errors.New("derivative wallet can only be created with wallet id " + DerivativeWalletId)
	}
	if _, err := rskAccount.BtcAddress(); err != nil {
		return nil, errors.New("derivative wallet can only be used if RSK account has derivation enabled")
	}
	wallet := &DerivativeWallet{conn: conn, rskAccount: rskAccount}
	err := wallet.initWallet()
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (wallet *DerivativeWallet) initWallet() error {
	const addressVerificationErrorTemplate = "error while verifying wallet has address: %w"
	var err error
	var info *btcjson.GetWalletInfoResult
	var btcAddress btcutil.Address
	var addressInfo *btcjson.GetAddressInfoResult

	if info, err = wallet.conn.client.GetWalletInfo(); err != nil || info.WalletName != wallet.conn.WalletId {
		if info, err = wallet.createWallet(); err != nil {
			return err
		}
	}
	_, ok := info.Scanning.Value.(btcjson.ScanProgress)
	if ok {
		return errors.New("wallet is still scanning, please wait for the scan to finish before initializing the server again")
	}

	if btcAddress, err = wallet.rskAccount.BtcAddress(); err != nil {
		return fmt.Errorf(addressVerificationErrorTemplate, err)
	}

	if addressInfo, err = wallet.conn.client.GetAddressInfo(btcAddress.EncodeAddress()); err != nil {
		return fmt.Errorf(addressVerificationErrorTemplate, err)
	} else if !addressInfo.Solvable || !addressInfo.IsWatchOnly {
		return wallet.importPublicKey()
	}

	return nil
}

func (wallet *DerivativeWallet) createWallet() (*btcjson.GetWalletInfoResult, error) {
	if _, err := wallet.conn.client.LoadWallet(wallet.conn.WalletId); err == nil {
		return wallet.conn.client.GetWalletInfo()
	}
	log.Infof("Wallet not found to be loaded, creating wallet %s...", wallet.conn.WalletId)
	err := wallet.conn.client.CreateReadonlyWallet(btcclient.ReadonlyWalletRequest{
		WalletName:         wallet.conn.WalletId,
		DisablePrivateKeys: true,
		Blank:              true,
		AvoidReuse:         false,
		Descriptors:        false,
	})

	if err != nil {
		return nil, fmt.Errorf("error while creating %s wallet: %w", wallet.conn.WalletId, err)
	}
	return wallet.conn.client.GetWalletInfo()
}

func (wallet *DerivativeWallet) importPublicKey() error {
	const errorTemplate = "error while importing public key: %w"
	pubKey, err := wallet.rskAccount.BtcPubKey()
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	err = wallet.conn.client.ImportPubKey(pubKey)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	err = wallet.conn.client.ImportAddressRescan(wallet.Address(), "", true)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	return errors.New("public key imported, rescan started, please wait for the rescan process to finish before initializing the server again")
}

func (wallet *DerivativeWallet) EstimateTxFees(toAddress string, value *entities.Wei) (blockchain.BtcFeeEstimation, error) {
	const quoteHashLength = 32

	if _, err := btcutil.DecodeAddress(toAddress, wallet.conn.NetworkParams); err != nil {
		return blockchain.BtcFeeEstimation{}, err
	}
	if err := EnsureLoadedBtcWallet(wallet.conn); err != nil {
		return blockchain.BtcFeeEstimation{}, err
	}

	amountInSatoshi, _ := value.ToSatoshi().Int64()
	output := []btcjson.PsbtOutput{
		{toAddress: btcutil.Amount(amountInSatoshi).ToUnit(btcutil.AmountBTC)},
		{"data": hex.EncodeToString(make([]byte, quoteHashLength))}, // quote hash output
	}

	feeRate, err := wallet.estimateFeeRate()
	if err != nil {
		return blockchain.BtcFeeEstimation{}, err
	}
	changeAddress, err := wallet.rskAccount.BtcAddress()
	if err != nil {
		return blockchain.BtcFeeEstimation{}, err
	}

	opts := btcjson.WalletCreateFundedPsbtOpts{
		ChangeAddress:   btcjson.String(changeAddress.EncodeAddress()),
		ChangePosition:  btcjson.Int64(changePosition),
		IncludeWatching: btcjson.Bool(true),
		FeeRate:         feeRate,
	}

	simulatedTx, err := wallet.conn.client.WalletCreateFundedPsbt(nil, output, nil, &opts, nil)
	if err != nil {
		return blockchain.BtcFeeEstimation{}, err
	}
	btcFee, err := btcutil.NewAmount(simulatedTx.Fee)
	if err != nil {
		return blockchain.BtcFeeEstimation{}, err
	}
	satoshiFee := btcFee.ToUnit(btcutil.AmountSatoshi)
	return blockchain.BtcFeeEstimation{
		Value:   entities.SatoshiToWei(uint64(satoshiFee)),
		FeeRate: utils.NewBigFloat64(*feeRate),
	}, nil
}

func (wallet *DerivativeWallet) GetBalance() (*entities.Wei, error) {
	var amount btcutil.Amount
	balance := new(entities.Wei)

	btcAddress, err := wallet.rskAccount.BtcAddress()
	if err != nil {
		return nil, err
	}
	if err = EnsureLoadedBtcWallet(wallet.conn); err != nil {
		return nil, err
	}

	utxos, err := wallet.conn.client.ListUnspentMinMaxAddresses(
		MinConfirmationsForUtxos,
		MaxConfirmationsForUtxos,
		[]btcutil.Address{btcAddress},
	)
	if err != nil {
		return nil, err
	}

	for _, utxo := range utxos {
		if amount, err = btcutil.NewAmount(utxo.Amount); err != nil {
			return nil, err
		}
		if utxo.Confirmations > 0 {
			balance.Add(balance, entities.SatoshiToWei(uint64(amount.ToUnit(btcutil.AmountSatoshi))))
		}
	}
	return balance, nil
}

func (wallet *DerivativeWallet) SendWithOpReturn(address string, value *entities.Wei, opReturnContent []byte) (transactionHash string, err error) {
	decodedAddress, err := btcutil.DecodeAddress(address, wallet.conn.NetworkParams)
	if err != nil {
		return "", err
	}
	if err = EnsureLoadedBtcWallet(wallet.conn); err != nil {
		return "", err
	}

	satoshis, _ := value.ToSatoshi().Float64()
	output := map[btcutil.Address]btcutil.Amount{decodedAddress: btcutil.Amount(satoshis)}
	rawTx, err := wallet.conn.client.CreateRawTransaction(nil, output, nil)
	if err != nil {
		return "", err
	}

	opReturnScript, err := txscript.NullDataScript(opReturnContent)
	if err != nil {
		return "", err
	}
	rawTx.AddTxOut(wire.NewTxOut(0, opReturnScript))

	opts, err := wallet.buildFundRawTransactionOpts()
	if err != nil {
		return "", err
	}
	fundedTx, err := wallet.conn.client.FundRawTransaction(rawTx, opts, nil)
	if err != nil {
		return "", err
	}

	signedTx, err := wallet.signFundedTransaction(fundedTx)
	if err != nil {
		return "", err
	}

	log.Infof("Sending %v BTC to %s\n", value.ToRbtc(), address)
	txHash, err := wallet.conn.client.SendRawTransaction(signedTx, false)
	if err != nil {
		return "", err
	}
	return txHash.String(), nil
}

func (wallet *DerivativeWallet) ImportAddress(address string) error {
	return errors.New("address importing is not supported in this type of wallet")
}

func (wallet *DerivativeWallet) GetTransactions(address string) ([]blockchain.BitcoinTransactionInformation, error) {
	if err := EnsureLoadedBtcWallet(wallet.conn); err != nil {
		return nil, err
	}
	return getTransactionsToAddress(address, wallet.conn.NetworkParams, wallet.conn.client)
}

func (wallet *DerivativeWallet) Address() string {
	address, err := wallet.rskAccount.BtcAddress()
	if err != nil {
		log.Errorf("error while getting address from rsk account %v", err)
		return ""
	}
	return address.EncodeAddress()
}

func (wallet *DerivativeWallet) Unlock() error {
	return errors.New("derivative wallet does not support unlocking as it is a watch-only wallet")
}

func (wallet *DerivativeWallet) estimateFeeRate() (*float64, error) {
	const (
		confirmationTargetForEstimation = 1
		minimumEstimatedConfirmations   = 2
		extraFeeMultiplier              = 0.1
		estimationMaxDecimals           = 8
	)
	estimationResult, err := wallet.conn.client.EstimateSmartFee(confirmationTargetForEstimation, &btcjson.EstimateModeEconomical)
	if err != nil {
		return nil, err
	} else if len(estimationResult.Errors) != 0 {
		return nil, errors.New(estimationResult.Errors[0])
	}
	// add 10% to the fee rate if result still over the target for the estimation
	if estimationResult.Blocks > confirmationTargetForEstimation && estimationResult.Blocks != minimumEstimatedConfirmations {
		return btcjson.Float64(utils.RoundToNDecimals(*estimationResult.FeeRate+*estimationResult.FeeRate*extraFeeMultiplier, estimationMaxDecimals)), nil
	}
	return btcjson.Float64(utils.RoundToNDecimals(*estimationResult.FeeRate, estimationMaxDecimals)), nil
}

func (wallet *DerivativeWallet) buildFundRawTransactionOpts() (btcjson.FundRawTransactionOpts, error) {
	feeRate, err := wallet.estimateFeeRate()
	if err != nil {
		return btcjson.FundRawTransactionOpts{}, err
	}
	changeAddress, err := wallet.rskAccount.BtcAddress()
	if err != nil {
		return btcjson.FundRawTransactionOpts{}, err
	}
	return btcjson.FundRawTransactionOpts{
		ChangeAddress:   btcjson.String(changeAddress.EncodeAddress()),
		ChangePosition:  btcjson.Int(changePosition),
		IncludeWatching: btcjson.Bool(true),
		LockUnspents:    btcjson.Bool(true),
		FeeRate:         feeRate,
		Replaceable:     btcjson.Bool(true),
	}, nil
}

func (wallet *DerivativeWallet) Shutdown(closeChannel chan<- bool) {
	wallet.conn.Shutdown(closeChannel)
}

func (wallet *DerivativeWallet) signFundedTransaction(fundedTx *btcjson.FundRawTransactionResult) (*wire.MsgTx, error) {
	var signedTx *wire.MsgTx
	var complete bool
	var err error
	signingErr := wallet.rskAccount.UsePrivateKeyWif(func(wif *btcutil.WIF) error {
		signedTx, complete, err = wallet.conn.client.SignRawTransactionWithKey(fundedTx.Transaction, []string{wif.String()})
		if err != nil {
			return err
		} else if !complete {
			return errors.New("trying to send a transaction without a complete set of signatures")
		}
		return nil
	})
	if signingErr != nil {
		return nil, signingErr
	}
	return signedTx, nil
}
