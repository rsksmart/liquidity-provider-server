package rootstock

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	log "github.com/sirupsen/logrus"
	"math/big"
)

type RskWalletImpl struct {
	client  *ethclient.Client
	account *RskAccount
	chainId uint64
}

func NewRskWalletImpl(client *RskClient, account *RskAccount, chainId uint64) *RskWalletImpl {
	return &RskWalletImpl{client: client.client, account: account, chainId: chainId}
}

func (wallet *RskWalletImpl) Address() common.Address {
	return wallet.account.Account.Address
}

func (wallet *RskWalletImpl) Sign(address common.Address, transaction *geth.Transaction) (*geth.Transaction, error) {
	var chainId big.Int
	if !bytes.Equal(address[:], wallet.account.Account.Address[:]) {
		return nil, fmt.Errorf("provider address %v is incorrect", address.Hash())
	}
	chainId.SetUint64(wallet.chainId)
	return wallet.account.Keystore.SignTx(*wallet.account.Account, transaction, &chainId)
}

func (wallet *RskWalletImpl) SignBytes(msg []byte) ([]byte, error) {
	return wallet.account.Keystore.SignHash(*wallet.account.Account, msg)
}

func (wallet *RskWalletImpl) Validate(signature, hash string) bool {
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		log.Error("Error decoding signature: ", err)
		return false
	}
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		log.Error("Error decoding hash: ", err)
		return false
	}
	pubKey, err := crypto.Ecrecover(hashBytes, signatureBytes)
	if err != nil {
		log.Error("Error recovering public key: ", err)
		return false
	}
	pubKeyHash := crypto.Keccak256Hash(pubKey[1:])
	return bytes.Equal(wallet.account.Account.Address.Bytes(), pubKeyHash[12:]) // the last 20 bytes of the hash
}

func (wallet *RskWalletImpl) SendRbtc(ctx context.Context, config blockchain.TransactionConfig, toAddress string) (string, error) {
	var to common.Address
	var signedTx *geth.Transaction
	var nonce uint64
	var err error

	if err = ParseAddress(&to, toAddress); err != nil {
		return "", err
	}

	newCtx, cancel := context.WithTimeout(ctx, txMiningWaitTimeout)
	defer cancel()

	if config.GasPrice == nil || config.Value == nil || config.GasLimit == nil {
		return "", errors.New("incomplete transaction arguments")
	}

	if nonce, err = wallet.client.PendingNonceAt(newCtx, wallet.Address()); err != nil {
		return "", err
	}

	tx := geth.NewTx(&geth.LegacyTx{
		To:       &to,
		Nonce:    nonce,
		GasPrice: config.GasPrice.AsBigInt(),
		Gas:      *config.GasLimit,
		Value:    config.Value.AsBigInt(),
	})
	log.Infof("Sending %v RBTC to %s\n", config.Value.ToRbtc(), toAddress)
	if signedTx, err = wallet.Sign(wallet.Address(), tx); err != nil {
		return "", err
	}
	err = wallet.client.SendTransaction(newCtx, signedTx)
	return tx.Hash().String(), err
}
