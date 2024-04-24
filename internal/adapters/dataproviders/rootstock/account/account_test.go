package account_test

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/account"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	keyAddress = "0x9d93929a9099be4355fc2389fbf253982f9df47c"
	keyPath    = "../../../../../docker-compose/localstack/local-key.json"
)

func TestGetAccount(t *testing.T) {
	testDir := filepath.Join(t.TempDir(), fmt.Sprintf("test-account-%d", time.Now().UnixNano()))
	keyFile, err := os.Open(keyPath)
	require.NoError(t, err)

	defer func(file *os.File) {
		closingErr := file.Close()
		require.NoError(t, closingErr)
	}(keyFile)

	keyBytes, err := io.ReadAll(keyFile)
	require.NoError(t, err)
	t.Run("Create new account", func(t *testing.T) {
		account, testError := account.GetRskAccount(testDir, 0, string(keyBytes), test.KeyPassword)
		_, noExistError := os.Stat(testDir)
		assert.Falsef(t, os.IsNotExist(noExistError), "Key directory not created")
		require.NoError(t, testError)
		assert.Equal(t, common.HexToAddress(keyAddress), account.Account.Address)
		assert.NotNil(t, 1, len(account.Keystore.Accounts()))
	})
	t.Run("Retrieve created account new account", func(t *testing.T) {
		otherAccount, otherError := account.GetRskAccount(testDir, 0, string(keyBytes), test.KeyPassword)
		require.NoError(t, otherError)
		assert.Equal(t, common.HexToAddress(keyAddress), otherAccount.Account.Address)
		assert.NotNil(t, 1, len(otherAccount.Keystore.Accounts()))
	})
}

func TestGetAccount_ErrorHandling(t *testing.T) {
	testDir := filepath.Join(t.TempDir(), fmt.Sprintf("test-%d", time.Now().UnixNano()))
	keyFile, setupErr := os.Open(keyPath)
	require.NoError(t, setupErr)

	defer func(file *os.File) {
		closingErr := file.Close()
		require.NoError(t, closingErr)
	}(keyFile)

	keyBytes, setupErr := io.ReadAll(keyFile)
	require.NoError(t, setupErr)
	t.Run("Invalid dir", func(t *testing.T) {
		account, err := account.GetRskAccount("/test", 0, string(keyBytes), test.KeyPassword)
		assert.Nil(t, account)
		require.Error(t, err)
	})
	t.Run("Invalid key", func(t *testing.T) {
		account, err := account.GetRskAccount(testDir, 0, "any key", test.KeyPassword)
		assert.Nil(t, account)
		require.Error(t, err)
	})
	t.Run("Invalid password", func(t *testing.T) {
		account, err := account.GetRskAccount(testDir, 0, string(keyBytes), "incorrect")
		assert.Nil(t, account)
		require.Error(t, err)
	})
	t.Run("Invalid account number", func(t *testing.T) {
		// we create a keystore first so in the second call we can try to get an account that doesn't exist
		_, err := account.GetRskAccount(testDir, 0, string(keyBytes), test.KeyPassword)
		require.NoError(t, err)
		account, err := account.GetRskAccount(testDir, 1, string(keyBytes), test.KeyPassword)
		assert.Nil(t, account)
		require.Error(t, err)
	})
}
