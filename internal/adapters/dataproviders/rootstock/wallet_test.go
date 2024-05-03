package rootstock_test

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/account"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

const (
	chainId            = 31
	timerContextString = "*context.timerCtx"
)

var (
	walletAddress = common.HexToAddress("0x9D93929A9099be4355fC2389FbF253982F9dF47c")
	signedBytes   = []byte{
		0x9f, 0x53, 0xc9, 0x15, 0x4, 0xf1, 0xbe, 0x5, 0x86, 0x9d, 0xd5, 0xaa, 0x8c, 0xd3,
		0x41, 0xde, 0x8e, 0x5e, 0x88, 0xaf, 0x8, 0xbe, 0xfe, 0x59, 0x37, 0x52, 0x8d, 0xba,
		0x95, 0xef, 0xd8, 0x8f, 0x60, 0x7c, 0xe1, 0x6c, 0xa3, 0x48, 0x65, 0x70, 0x38, 0x2a,
		0xc1, 0x8, 0xb8, 0x8d, 0x8d, 0xee, 0xff, 0x26, 0xb0, 0x14, 0xbd, 0x22, 0x27, 0xab,
		0x1e, 0x41, 0x75, 0xd7, 0x28, 0x5c, 0x86, 0x9d, 0x1,
	}
	bytesToSign = []byte{
		0xee, 0x08, 0x74, 0x17, 0x0b, 0x7f, 0x6f, 0x32, 0xb8, 0xc2, 0xac, 0x95, 0x73, 0xc4, 0x28, 0xd3,
		0x5b, 0x57, 0x52, 0x70, 0xa6, 0x6b, 0x75, 0x7c, 0x2c, 0x01, 0x85, 0xd2, 0xbd, 0x09, 0x71, 0x8d,
	}
)

// TestRskWalletImpl we use this function to run all the test related to the wallet to open the account only once
func TestRskWalletImpl(t *testing.T) {
	account := test.OpenWalletForTest(t, "wallet")
	t.Run("Address", createAddressTest(account))
	t.Run("Sign", creteSignTest(account))
	t.Run("SignBytes", createSignBytesTest(account))
	t.Run("Validate", createValidateTest(account))
	t.Run("SendRbtc", createSendRbtcTest(account))
	t.Run("SendRbtc error handling", createSendRbtcErrorHandlingTest(account))
}

func createSendRbtcTest(account *account.RskAccount) func(t *testing.T) {
	return func(t *testing.T) {
		const toAddress = "0x79568C2989232dcA1840087d73d403602364c0D4"
		var gasLimit uint64 = 21000
		t.Run("Success", func(t *testing.T) {
			clientMock := &mocks.RpcClientBindingMock{}
			clientMock.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *geth.Transaction) bool {
				v, r, s := tx.RawSignatureValues()
				return assert.NotNil(t, v) && assert.NotNil(t, r) && assert.NotNil(t, s)
			})).Return(nil)
			clientMock.On("PendingNonceAt", mock.AnythingOfType(timerContextString), walletAddress).Return(uint64(54), nil)
			wallet := rootstock.NewRskWalletImpl(rootstock.NewRskClient(clientMock), account, chainId)
			tx, err := wallet.SendRbtc(context.Background(), blockchain.TransactionConfig{
				Value:    entities.NewWei(89607151182921727),
				GasLimit: &gasLimit,
				GasPrice: entities.NewWei(65164000),
			}, toAddress)
			require.NoError(t, err)
			require.Equal(t, "0xa685c956bd47a5c6c9d66997a469f483447fb1366709f7374764ee597aeac266", tx)
		})
	}
}

func createSendRbtcErrorHandlingTest(account *account.RskAccount) func(t *testing.T) {
	return func(t *testing.T) {
		const toAddress = "0x79568C2989232dcA1840087d73d403602364c0D4"
		var gasLimit uint64 = 21000
		t.Run("Handle error on invalid address", func(t *testing.T) {
			clientMock := &mocks.RpcClientBindingMock{}
			wallet := rootstock.NewRskWalletImpl(rootstock.NewRskClient(clientMock), account, chainId)
			tx, err := wallet.SendRbtc(context.Background(), blockchain.TransactionConfig{}, test.AnyString)
			require.ErrorIs(t, err, blockchain.InvalidAddressError)
			require.Empty(t, tx)
		})
		t.Run("Handle error on incomplete config", func(t *testing.T) {
			const incompleteConfig = "incomplete transaction arguments"
			clientMock := &mocks.RpcClientBindingMock{}
			wallet := rootstock.NewRskWalletImpl(rootstock.NewRskClient(clientMock), account, chainId)
			t.Run("Missing gasPrice", func(t *testing.T) {
				tx, err := wallet.SendRbtc(context.Background(), blockchain.TransactionConfig{
					Value: entities.NewWei(1), GasLimit: &gasLimit}, toAddress)
				require.ErrorContains(t, err, incompleteConfig)
				require.Empty(t, tx)
			})
			t.Run("Missing value", func(t *testing.T) {
				tx, err := wallet.SendRbtc(context.Background(), blockchain.TransactionConfig{
					GasPrice: entities.NewWei(1), GasLimit: &gasLimit}, toAddress)
				require.ErrorContains(t, err, incompleteConfig)
				require.Empty(t, tx)
			})
			t.Run("Missing gasLimit", func(t *testing.T) {
				tx, err := wallet.SendRbtc(context.Background(), blockchain.TransactionConfig{
					Value: entities.NewWei(1), GasPrice: entities.NewWei(1)}, toAddress)
				require.ErrorContains(t, err, incompleteConfig)
				require.Empty(t, tx)
			})
		})
		t.Run("Handle error on failure when getting nonce", func(t *testing.T) {
			clientMock := &mocks.RpcClientBindingMock{}
			clientMock.On("PendingNonceAt", mock.AnythingOfType(timerContextString), walletAddress).Return(uint64(0), assert.AnError)
			wallet := rootstock.NewRskWalletImpl(rootstock.NewRskClient(clientMock), account, chainId)
			tx, err := wallet.SendRbtc(context.Background(), blockchain.TransactionConfig{
				Value:    entities.NewWei(1),
				GasLimit: &gasLimit,
				GasPrice: entities.NewWei(5),
			}, toAddress)
			require.Error(t, err)
			require.Empty(t, tx)
		})
		t.Run("Handle error on failure when broadcasting tx", func(t *testing.T) {
			clientMock := &mocks.RpcClientBindingMock{}
			clientMock.On("SendTransaction", mock.Anything, mock.Anything).Return(assert.AnError)
			clientMock.On("PendingNonceAt", mock.AnythingOfType(timerContextString), walletAddress).Return(uint64(54), nil)
			wallet := rootstock.NewRskWalletImpl(rootstock.NewRskClient(clientMock), account, chainId)
			tx, err := wallet.SendRbtc(context.Background(), blockchain.TransactionConfig{
				Value:    entities.NewWei(1),
				GasLimit: &gasLimit,
				GasPrice: entities.NewWei(5),
			}, toAddress)
			require.Error(t, err)
			require.Empty(t, tx)
		})
	}
}

func createAddressTest(account *account.RskAccount) func(t *testing.T) {
	return func(t *testing.T) {
		clientMock := &mocks.RpcClientBindingMock{}
		wallet := rootstock.NewRskWalletImpl(rootstock.NewRskClient(clientMock), account, chainId)
		address := wallet.Address()
		assert.Equal(t, walletAddress, address)
	}
}

func creteSignTest(account *account.RskAccount) func(t *testing.T) {
	return func(t *testing.T) {
		clientMock := &mocks.RpcClientBindingMock{}
		wallet := rootstock.NewRskWalletImpl(rootstock.NewRskClient(clientMock), account, chainId)
		toAddress := common.HexToAddress("0x8dccd82443B80DDdE3690af86746bfd9d766F8D2")
		tx := geth.NewTx(&geth.LegacyTx{
			To:       &toAddress,
			Nonce:    123,
			GasPrice: big.NewInt(700),
			Gas:      500,
			Value:    big.NewInt(800),
		})
		t.Run("Success", func(t *testing.T) {
			result, err := wallet.Sign(walletAddress, tx)
			require.NoError(t, err)
			v, r, s := result.RawSignatureValues()
			assert.NotNil(t, v)
			assert.NotNil(t, r)
			assert.NotNil(t, s)
		})
		t.Run("Error when signing with wrong account", func(t *testing.T) {
			result, err := wallet.Sign(toAddress, tx)
			require.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

func createSignBytesTest(account *account.RskAccount) func(t *testing.T) {
	return func(t *testing.T) {
		clientMock := &mocks.RpcClientBindingMock{}
		wallet := rootstock.NewRskWalletImpl(rootstock.NewRskClient(clientMock), account, chainId)
		signedHash, err := wallet.SignBytes(bytesToSign)
		require.NoError(t, err)
		assert.Equal(t, signedBytes, signedHash)
	}
}

func createValidateTest(account *account.RskAccount) func(t *testing.T) {
	return func(t *testing.T) {
		const noHex = "no hex"
		clientMock := &mocks.RpcClientBindingMock{}
		wallet := rootstock.NewRskWalletImpl(rootstock.NewRskClient(clientMock), account, chainId)
		t.Run("Success", func(t *testing.T) {
			isValid := wallet.Validate(hex.EncodeToString(signedBytes), hex.EncodeToString(bytesToSign))
			assert.True(t, isValid)
		})
		t.Run("Invalid signature", func(t *testing.T) {
			isValid := wallet.Validate(noHex, hex.EncodeToString(bytesToSign))
			assert.False(t, isValid)
		})
		t.Run("Invalid hash", func(t *testing.T) {
			isValid := wallet.Validate(hex.EncodeToString(signedBytes), noHex)
			assert.False(t, isValid)
		})
		t.Run("Signature mismatch", func(t *testing.T) {
			tamperedSignature := signedBytes
			tamperedSignature[0] = 0x12
			tamperedHash := bytesToSign
			tamperedHash[0] = 0x12
			assert.False(t, wallet.Validate(hex.EncodeToString(tamperedSignature), hex.EncodeToString(bytesToSign)))
			assert.False(t, wallet.Validate(hex.EncodeToString(signedBytes), hex.EncodeToString(tamperedHash)))
		})
	}
}
