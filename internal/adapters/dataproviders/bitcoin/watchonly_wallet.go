package bitcoin

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin/btcclient"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	log "github.com/sirupsen/logrus"
)

const (
	PeginWalletId = "pegin-watchonly-wallet"
)

type WatchOnlyWallet struct {
	conn *Connection
}

func NewWatchOnlyWallet(walletConnection *Connection) (blockchain.BitcoinWallet, error) {
	wallet := &WatchOnlyWallet{conn: walletConnection}

	var info *btcjson.GetWalletInfoResult
	var err error
	if info, err = wallet.conn.client.GetWalletInfo(); err != nil {
		if info, err = wallet.createWallet(); err != nil {
			return nil, fmt.Errorf("error creating watch-only wallet: %w", err)
		}
	}
	if info.PrivateKeysEnabled {
		return nil, errors.New("wallet is not watch-only")
	}
	return wallet, nil
}

func (wallet *WatchOnlyWallet) createWallet() (*btcjson.GetWalletInfoResult, error) {
	_, err := wallet.conn.client.LoadWallet(wallet.conn.WalletId)
	if err == nil {
		return wallet.conn.client.GetWalletInfo()
	}
	err = wallet.conn.client.CreateReadonlyWallet(btcclient.ReadonlyWalletRequest{
		WalletName:         wallet.conn.WalletId,
		DisablePrivateKeys: true,
		Blank:              true,
		AvoidReuse:         true,
		Descriptors:        false,
	})
	if err != nil {
		return nil, err
	}
	return wallet.conn.client.GetWalletInfo()
}

func (wallet *WatchOnlyWallet) EstimateTxFees(toAddress string, value *entities.Wei) (blockchain.BtcFeeEstimation, error) {
	return blockchain.BtcFeeEstimation{}, errors.New("cannot estimate from a watch-only wallet")
}

func (wallet *WatchOnlyWallet) GetBalance() (*entities.Wei, error) {
	return nil, errors.New("cannot get balance of a watch-only wallet since it may be tracking address from multiple wallets")
}

func (wallet *WatchOnlyWallet) SendWithOpReturn(address string, value *entities.Wei, opReturnContent []byte) (blockchain.BitcoinTransactionResult, error) {
	return blockchain.BitcoinTransactionResult{}, errors.New("cannot send from a watch-only wallet")
}

func (wallet *WatchOnlyWallet) ImportAddress(address string) error {
	_, err := btcutil.DecodeAddress(address, wallet.conn.NetworkParams)
	if err != nil {
		return err
	}
	if err = EnsureLoadedBtcWallet(wallet.conn); err != nil {
		return err
	}
	return wallet.conn.client.ImportAddressRescan(address, "", false)
}

func (wallet *WatchOnlyWallet) GetTransactions(address string) ([]blockchain.BitcoinTransactionInformation, error) {
	if err := EnsureLoadedBtcWallet(wallet.conn); err != nil {
		return nil, err
	}
	return getTransactionsToAddress(address, wallet.conn.NetworkParams, wallet.conn.client)
}

func (wallet *WatchOnlyWallet) Address() string {
	log.Warn("Trying to get main address from a watch-only wallet")
	return ""
}

func (wallet *WatchOnlyWallet) Unlock() error {
	return errors.New("watch-only wallet does not support unlocking as it only has monitoring purposes")
}

func (wallet *WatchOnlyWallet) Shutdown(closeChannel chan<- bool) {
	wallet.conn.Shutdown(closeChannel)
}
