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
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/account"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	log "github.com/sirupsen/logrus"
	"math/big"
	"time"
)

type RskWalletImpl struct {
	client        RpcClientBinding
	account       *account.RskAccount
	chainId       uint64
	miningTimeout time.Duration
}

func NewRskWalletImpl(client *RskClient, account *account.RskAccount, chainId uint64, miningTimeout time.Duration) *RskWalletImpl {
	return &RskWalletImpl{client: client.client, account: account, chainId: chainId, miningTimeout: miningTimeout}
}

func (wallet *RskWalletImpl) Address() common.Address {
	return wallet.account.Account.Address
}

func (wallet *RskWalletImpl) Sign(address common.Address, transaction *geth.Transaction) (*geth.Transaction, error) {
	var chainId big.Int
	if !bytes.Equal(address[:], wallet.account.Account.Address[:]) {
		return nil, fmt.Errorf("provider address %v is incorrect", address.Hex())
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

func (wallet *RskWalletImpl) SendRbtc(ctx context.Context, config blockchain.TransactionConfig, toAddress string) (blockchain.TransactionReceipt, error) {
	var to common.Address
	var signedTx *geth.Transaction
	var nonce uint64
	var err error

	receiptData := blockchain.TransactionReceipt{
		TransactionHash: "",
		GasUsed:         big.NewInt(0),
		GasPrice:        big.NewInt(0),
	}

	if err = ParseAddress(&to, toAddress); err != nil {
		return receiptData, err
	}

	if config.GasPrice == nil || config.Value == nil || config.GasLimit == nil {
		return receiptData, errors.New("incomplete transaction arguments")
	}

	if nonce, err = wallet.client.PendingNonceAt(ctx, wallet.Address()); err != nil {
		return receiptData, err
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
		return receiptData, err
	}

	sendError := wallet.client.SendTransaction(ctx, signedTx)
	receipt, err := AwaitTxWithCtx(wallet.client, wallet.miningTimeout, "SendRbtc", ctx, func() (*geth.Transaction, error) {
		return signedTx, sendError
	})

	receiptData.TransactionHash = signedTx.Hash().String()

	if err != nil {
		return receiptData, err
	} else if receipt == nil || receipt.Status == 0 {
		return receiptData, fmt.Errorf("%s transaction failed", receiptData.TransactionHash)
	}
	receiptData.GasUsed = big.NewInt(int64(receipt.GasUsed))
	return receiptData, nil
}

func (wallet *RskWalletImpl) GetBalance(ctx context.Context) (*entities.Wei, error) {
	balance, err := wallet.client.BalanceAt(ctx, wallet.Address(), nil)
	if err != nil {
		return nil, err
	}
	return entities.NewBigWei(balance), nil
}
