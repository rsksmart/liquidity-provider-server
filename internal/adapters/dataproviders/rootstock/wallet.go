package rootstock

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/account"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	log "github.com/sirupsen/logrus"
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
	to, nonce, err := wallet.validateAndPrepareSendRbtc(ctx, config, toAddress)
	if err != nil {
		return blockchain.TransactionReceipt{}, err
	}

	signedTx, err := wallet.createAndSignTransaction(to, nonce, config, toAddress)
	if err != nil {
		return blockchain.TransactionReceipt{}, err
	}

	receipt, err := wallet.sendAndAwaitTransaction(ctx, signedTx)
	if err != nil {
		return blockchain.TransactionReceipt{}, err
	}

	return wallet.buildTransactionReceipt(ctx, receipt)
}

func (wallet *RskWalletImpl) validateAndPrepareSendRbtc(ctx context.Context, config blockchain.TransactionConfig, toAddress string) (common.Address, uint64, error) {
	var to common.Address
	var nonce uint64
	var err error

	if err = ParseAddress(&to, toAddress); err != nil {
		return common.Address{}, 0, err
	}

	if config.GasPrice == nil || config.Value == nil || config.GasLimit == nil {
		return common.Address{}, 0, errors.New("incomplete transaction arguments")
	}

	if nonce, err = wallet.client.PendingNonceAt(ctx, wallet.Address()); err != nil {
		return common.Address{}, 0, err
	}

	return to, nonce, nil
}

func (wallet *RskWalletImpl) createAndSignTransaction(to common.Address, nonce uint64, config blockchain.TransactionConfig, toAddress string) (*geth.Transaction, error) {
	tx := geth.NewTx(&geth.LegacyTx{
		To:       &to,
		Nonce:    nonce,
		GasPrice: config.GasPrice.AsBigInt(),
		Gas:      *config.GasLimit,
		Value:    config.Value.AsBigInt(),
	})
	log.Infof("Sending %v RBTC to %s\n", config.Value.ToRbtc(), toAddress)

	return wallet.Sign(wallet.Address(), tx)
}

func (wallet *RskWalletImpl) sendAndAwaitTransaction(ctx context.Context, signedTx *geth.Transaction) (*geth.Receipt, error) {
	sendError := wallet.client.SendTransaction(ctx, signedTx)
	receipt, err := AwaitTxWithCtx(wallet.client, wallet.miningTimeout, "SendRbtc", ctx, func() (*geth.Transaction, error) {
		return signedTx, sendError
	})

	if err != nil {
		return nil, err
	}
	if receipt == nil {
		return nil, errors.New("send rbtc error: incomplete receipt")
	}

	return receipt, nil
}

func (wallet *RskWalletImpl) buildTransactionReceipt(ctx context.Context, receipt *geth.Receipt) (blockchain.TransactionReceipt, error) {
	// Fetch the transaction to get the "To" address and the Value
	toAddressStr := ""
	txValue := entities.NewWei(0)
	if tx, _, txErr := wallet.client.TransactionByHash(ctx, receipt.TxHash); txErr == nil {
		if tx.To() != nil {
			toAddressStr = tx.To().String()
		}
		txValue = entities.NewBigWei(tx.Value())
	}

	transactionReceipt := blockchain.TransactionReceipt{
		TransactionHash:   receipt.TxHash.String(),
		BlockHash:         receipt.BlockHash.String(),
		BlockNumber:       receipt.BlockNumber.Uint64(),
		From:              wallet.Address().String(),
		To:                toAddressStr,
		CumulativeGasUsed: new(big.Int).SetUint64(receipt.CumulativeGasUsed),
		GasUsed:           new(big.Int).SetUint64(receipt.GasUsed),
		Value:             txValue,
		GasPrice:          entities.NewWei(receipt.EffectiveGasPrice.Int64()),
	}

	if receipt.Status == 0 {
		return transactionReceipt, fmt.Errorf("send rbtc error: transaction reverted (%s)", receipt.TxHash.String())
	}

	return transactionReceipt, nil
}

func (wallet *RskWalletImpl) GetBalance(ctx context.Context) (*entities.Wei, error) {
	balance, err := wallet.client.BalanceAt(ctx, wallet.Address(), nil)
	if err != nil {
		return nil, err
	}
	return entities.NewBigWei(balance), nil
}
