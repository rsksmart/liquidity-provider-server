package test

import (
	"bytes"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/account"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

var AnyCtx = mock.AnythingOfType("context.backgroundCtx")

const (
	AnyAddress  = "any address"
	AnyString   = "any value"
	keyPath     = "../../docker-compose/localstack/local-key.json"
	KeyPassword = "test"
)

type Case[V, R any] struct {
	Value  V
	Result R
}

type Table[V, R any] []Case[V, R]

func RunTable[V, R any](t *testing.T, table Table[V, R], validationFunction func(V) R) {
	var result R
	for _, testCase := range table {
		result = validationFunction(testCase.Value)
		assert.Equal(t, testCase.Result, result)
	}
}

// CountNonZeroValues counts the number of non-zero values in a struct, it panics if value of aStruct parameter
// is not a struct
func CountNonZeroValues(aStruct any) int {
	value := reflect.ValueOf(aStruct)
	count := 0
	for i := 0; i < value.NumField(); i++ {
		if !value.Field(i).IsZero() {
			count++
		}
	}
	return count
}

func AssertNoLog(t *testing.T) (assertFunc func()) {
	buff := new(bytes.Buffer)
	log.SetOutput(buff)
	return func() {
		assert.Empty(t, buff.Bytes())
	}
}

func AssertLogContains(t *testing.T, expected string) (assertFunc func()) {
	message := make([]byte, 1024)
	buff := new(bytes.Buffer)
	log.SetOutput(buff)
	return func() {
		_, err := buff.Read(message)
		require.NoError(t, err, "Error reading log message")
		assert.Contains(t, string(message), expected, "Expected message not found")
	}
}

func OpenDerivativeWalletForTest(t *testing.T, testRef string) *account.RskAccount {
	_, currentPackageDir, _, _ := runtime.Caller(0)
	testDir := filepath.Join(t.TempDir(), fmt.Sprintf("test-derivative-%s-%d", testRef, time.Now().UnixNano()))
	keyFile, err := os.Open(filepath.Join(currentPackageDir, keyPath))
	require.NoError(t, err)

	defer func(file *os.File) { closingErr := file.Close(); require.NoError(t, closingErr) }(keyFile)

	keyBytes, err := io.ReadAll(keyFile)
	require.NoError(t, err)
	testAccount, err := account.GetRskAccountWithDerivation(account.CreationWithDerivationArgs{
		CreationArgs: account.CreationArgs{
			KeyDir:        testDir,
			AccountNum:    0,
			EncryptedJson: string(keyBytes),
			Password:      KeyPassword,
		},
		BtcParams: &chaincfg.TestNet3Params,
	})
	require.NoError(t, err)
	return testAccount
}

func OpenWalletForTest(t *testing.T, testRef string) *account.RskAccount {
	_, currentPackageDir, _, _ := runtime.Caller(0)
	testDir := filepath.Join(t.TempDir(), fmt.Sprintf("test-%s-%d", testRef, time.Now().UnixNano()))
	keyFile, err := os.Open(filepath.Join(currentPackageDir, keyPath))
	require.NoError(t, err)

	defer func(file *os.File) {
		closingErr := file.Close()
		require.NoError(t, closingErr)
	}(keyFile)

	keyBytes, err := io.ReadAll(keyFile)
	require.NoError(t, err)
	testAccount, err := account.GetRskAccount(account.CreationArgs{
		KeyDir:        testDir,
		AccountNum:    0,
		EncryptedJson: string(keyBytes),
		Password:      KeyPassword,
	})
	require.NoError(t, err)
	return testAccount
}
