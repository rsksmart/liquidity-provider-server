package account_test

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/account"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"testing"
	"time"
)

const (
	keyAddress = "0x9d93929a9099be4355fc2389fbf253982f9df47c"
	keyPath    = "../../../../../docker-compose/localstack/local-key.json"
)

var derivationKeystore = fmt.Sprintf("test-account-derivation-%d", time.Now().UnixNano())

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
		testAccount, testError := account.GetRskAccount(account.CreationArgs{
			KeyDir:        testDir,
			AccountNum:    0,
			EncryptedJson: string(keyBytes),
			Password:      test.KeyPassword,
		})
		_, noExistError := os.Stat(testDir)
		assert.Falsef(t, os.IsNotExist(noExistError), "Key directory not created")
		require.NoError(t, testError)
		assert.Equal(t, common.HexToAddress(keyAddress), testAccount.Account.Address)
		assert.NotNil(t, 1, len(testAccount.Keystore.Accounts()))
	})
	t.Run("Retrieve created account new account", func(t *testing.T) {
		otherAccount, otherError := account.GetRskAccount(account.CreationArgs{
			KeyDir:        testDir,
			AccountNum:    0,
			EncryptedJson: string(keyBytes),
			Password:      test.KeyPassword,
		})
		require.NoError(t, otherError)
		assert.Equal(t, common.HexToAddress(keyAddress), otherAccount.Account.Address)
		assert.NotNil(t, 1, len(otherAccount.Keystore.Accounts()))
	})
	t.Run("Hasn't access to derivation methods", func(t *testing.T) {
		notDerivativeAccount, e := account.GetRskAccount(account.CreationArgs{
			KeyDir:        testDir,
			AccountNum:    0,
			EncryptedJson: string(keyBytes),
			Password:      test.KeyPassword,
		})
		require.NoError(t, e)
		btcAddress, e := notDerivativeAccount.BtcAddress()
		require.ErrorIs(t, e, account.NoDerivationError)
		assert.Empty(t, btcAddress)
		pubKey, e := notDerivativeAccount.BtcPubKey()
		require.ErrorIs(t, e, account.NoDerivationError)
		assert.Empty(t, pubKey)
		e = notDerivativeAccount.UsePrivateKeyWif(func(wif *btcutil.WIF) error {
			return nil
		})
		require.ErrorIs(t, e, account.NoDerivationError)
	})
}

func TestGetRskAccountWithDerivation(t *testing.T) {
	testDir := filepath.Join(t.TempDir(), derivationKeystore)
	keyFile, err := os.Open(keyPath)
	require.NoError(t, err)

	defer func(file *os.File) { closingErr := file.Close(); require.NoError(t, closingErr) }(keyFile)

	keyBytes, err := io.ReadAll(keyFile)
	require.NoError(t, err)
	t.Run("Create new account", func(t *testing.T) {
		testAccount, testError := account.GetRskAccountWithDerivation(account.CreationWithDerivationArgs{
			CreationArgs: account.CreationArgs{
				KeyDir:        testDir,
				AccountNum:    0,
				EncryptedJson: string(keyBytes),
				Password:      test.KeyPassword,
			},
			BtcParams: &chaincfg.TestNet3Params,
		})
		_, noExistError := os.Stat(testDir)
		assert.Falsef(t, os.IsNotExist(noExistError), "Key directory not created")
		require.NoError(t, testError)
		assert.Equal(t, common.HexToAddress(keyAddress), testAccount.Account.Address)
		assert.NotNil(t, 1, len(testAccount.Keystore.Accounts()))
	})
	t.Run("Retrieve created account", func(t *testing.T) {
		otherAccount, otherError := account.GetRskAccountWithDerivation(account.CreationWithDerivationArgs{
			CreationArgs: account.CreationArgs{
				KeyDir:        testDir,
				AccountNum:    0,
				EncryptedJson: string(keyBytes),
				Password:      test.KeyPassword,
			},
			BtcParams: &chaincfg.TestNet3Params,
		})
		require.NoError(t, otherError)
		assert.Equal(t, common.HexToAddress(keyAddress), otherAccount.Account.Address)
		assert.NotNil(t, 1, len(otherAccount.Keystore.Accounts()))
	})
	t.Run("Has access to derivation methods", func(t *testing.T) {
		derivativeAccount, e := account.GetRskAccountWithDerivation(account.CreationWithDerivationArgs{
			CreationArgs: account.CreationArgs{
				KeyDir:        testDir,
				AccountNum:    0,
				EncryptedJson: string(keyBytes),
				Password:      test.KeyPassword,
			},
			BtcParams: &chaincfg.TestNet3Params,
		})
		require.NoError(t, e)
		btcAddress, e := derivativeAccount.BtcAddress()
		require.NoError(t, e)
		assert.NotEmpty(t, btcAddress)
		pubKey, e := derivativeAccount.BtcPubKey()
		require.NoError(t, e)
		assert.NotEmpty(t, pubKey)
		e = derivativeAccount.UsePrivateKeyWif(func(wif *btcutil.WIF) error { return nil })
		require.NoError(t, e)
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
		testAccount, err := account.GetRskAccount(account.CreationArgs{
			KeyDir:        "/test",
			AccountNum:    0,
			EncryptedJson: string(keyBytes),
			Password:      test.KeyPassword,
		})
		assert.Nil(t, testAccount)
		require.Error(t, err)
	})
	t.Run("Invalid key", func(t *testing.T) {
		testAccount, err := account.GetRskAccount(account.CreationArgs{
			KeyDir:        testDir,
			AccountNum:    0,
			EncryptedJson: "any key",
			Password:      test.KeyPassword,
		})
		assert.Nil(t, testAccount)
		require.Error(t, err)
	})
	t.Run("Invalid password", func(t *testing.T) {
		testAccount, err := account.GetRskAccount(account.CreationArgs{
			KeyDir:        testDir,
			AccountNum:    0,
			EncryptedJson: string(keyBytes),
			Password:      "incorrect",
		})
		assert.Nil(t, testAccount)
		require.Error(t, err)
	})
	t.Run("Invalid account number", func(t *testing.T) {
		// we create a keystore first so in the second call we can try to get an account that doesn't exist
		_, err := account.GetRskAccount(account.CreationArgs{
			KeyDir:        testDir,
			AccountNum:    0,
			EncryptedJson: string(keyBytes),
			Password:      test.KeyPassword,
		})
		require.NoError(t, err)
		testAccount, err := account.GetRskAccount(account.CreationArgs{
			KeyDir:        testDir,
			AccountNum:    1,
			EncryptedJson: string(keyBytes),
			Password:      test.KeyPassword,
		})
		assert.Nil(t, testAccount)
		require.Error(t, err)
	})
}

func TestRskAccount(t *testing.T) {
	testDir := filepath.Join(t.TempDir(), fmt.Sprintf("test-derivation-methods-%d", time.Now().UnixNano()))
	keyFile, setupErr := os.Open(keyPath)
	require.NoError(t, setupErr)

	defer func(file *os.File) {
		closingErr := file.Close()
		require.NoError(t, closingErr)
	}(keyFile)

	keyBytes, setupErr := io.ReadAll(keyFile)
	require.NoError(t, setupErr)
	testnetAccount, err := account.GetRskAccountWithDerivation(account.CreationWithDerivationArgs{
		CreationArgs: account.CreationArgs{
			KeyDir:        testDir,
			AccountNum:    0,
			EncryptedJson: string(keyBytes),
			Password:      test.KeyPassword,
		},
		BtcParams: &chaincfg.TestNet3Params,
	})
	require.NoError(t, err)
	mainnetAccount, err := account.GetRskAccountWithDerivation(account.CreationWithDerivationArgs{
		CreationArgs: account.CreationArgs{
			KeyDir:        testDir,
			AccountNum:    0,
			EncryptedJson: string(keyBytes),
			Password:      test.KeyPassword,
		},
		BtcParams: &chaincfg.MainNetParams,
	})
	require.NoError(t, err)

	t.Run("Test BtcPubKey", func(t *testing.T) {
		pubTestnet, errTestnet := testnetAccount.BtcPubKey()
		assert.Equal(t, "0232858a5faa413101831afe7a880da9a8ac4de6bd5e25b4358d762ba450b03c22", pubTestnet)
		require.NoError(t, errTestnet)
		pubMainnet, errMainnet := mainnetAccount.BtcPubKey()
		assert.Equal(t, "0232858a5faa413101831afe7a880da9a8ac4de6bd5e25b4358d762ba450b03c22", pubMainnet)
		require.NoError(t, errMainnet)
	})

	t.Run("Test BtcAddress", func(t *testing.T) {
		testnetAddress, errTestnet := testnetAccount.BtcAddress()
		assert.Equal(t, "n1jGDaxCW6jemLZyd9wmDHddseZwEMV9C6", testnetAddress.EncodeAddress())
		require.NoError(t, errTestnet)
		mainnetAddress, errMainnet := mainnetAccount.BtcAddress()
		assert.Equal(t, "1MDJvXsDh5JPzE6MuayPPNRK1eyEMmmMCW", mainnetAddress.EncodeAddress())
		require.NoError(t, errMainnet)
	})

	t.Run("Test UsePrivateKeyWif", func(t *testing.T) {
		testUsePrivateKeyWif(t, testnetAccount, mainnetAccount)
	})
}

func testUsePrivateKeyWif(t *testing.T, testnetAccount, mainnetAccount *account.RskAccount) {
	t.Run("Should secure WIF pointer", func(t *testing.T) {
		var testnetPointer, mainnetPointer *btcutil.WIF

		debug.SetPanicOnFault(true)
		defer debug.SetPanicOnFault(false)

		err := testnetAccount.UsePrivateKeyWif(func(wif *btcutil.WIF) error {
			testnetPointer = wif
			testnetPointer.SerializePubKey()
			return nil
		})
		require.NoError(t, err)
		assert.Panics(t, func() { testnetPointer.SerializePubKey() })

		err = mainnetAccount.UsePrivateKeyWif(func(wif *btcutil.WIF) error {
			mainnetPointer = wif
			mainnetPointer.SerializePubKey()
			return nil
		})
		require.NoError(t, err)
		assert.Panics(t, func() { mainnetPointer.SerializePubKey() })
	})
	t.Run("Should return errors", func(t *testing.T) {
		testnetErr := testnetAccount.UsePrivateKeyWif(func(wif *btcutil.WIF) error { return assert.AnError })
		require.Error(t, testnetErr)
		mainnetErr := mainnetAccount.UsePrivateKeyWif(func(wif *btcutil.WIF) error { return assert.AnError })
		require.Error(t, mainnetErr)
	})
	t.Run("Should execute function", func(t *testing.T) {
		var testnet, mainnet string
		err := testnetAccount.UsePrivateKeyWif(func(wif *btcutil.WIF) error {
			testnet = test.AnyString
			return nil
		})
		require.NoError(t, err)
		err = mainnetAccount.UsePrivateKeyWif(func(wif *btcutil.WIF) error {
			mainnet = test.AnyAddress
			return nil
		})
		require.NoError(t, err)
		assert.Equal(t, test.AnyString, testnet)
		assert.Equal(t, test.AnyAddress, mainnet)
	})
}
