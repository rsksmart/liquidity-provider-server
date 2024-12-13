package test

import (
	"bytes"
	"flag"
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
	"sync"
	"testing"
	"time"
)

var (
	AnyCtx = mock.AnythingOfType("context.backgroundCtx")
	AnyWei = mock.AnythingOfType("*entities.Wei")
)

const (
	AnyAddress    = "any address"
	AnyRskAddress = "0x79568c2989232dCa1840087D73d403602364c0D4"
	AnyBtcAddress = "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe"
	AnyString     = "any value"
	AnyHash       = "d8f5d705f146230553a8aec9a290a19bf4311187fa0489d41207d7215b0b65cb"
	AnyUrl        = "url.com"
	keyPath       = "../../docker-compose/localstack/local-key.json"
	KeyPassword   = "test"
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

type ThreadSafeBuffer struct {
	bytes.Buffer
	mutex sync.RWMutex
}

func (b *ThreadSafeBuffer) Write(p []byte) (n int, err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.Buffer.Write(p)
}

func (b *ThreadSafeBuffer) Read(p []byte) (n int, err error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.Buffer.Read(p)
}

func (b *ThreadSafeBuffer) Len() int {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.Buffer.Len()
}

func AssertNoLog(t *testing.T) (assertFunc func()) {
	buff := new(bytes.Buffer)
	log.SetOutput(buff)
	return func() {
		assert.Empty(t, buff.Bytes())
	}
}

func AssertLogContains(t *testing.T, expected string) (assertFunc func() bool) {
	message := make([]byte, 2048)
	buff := new(ThreadSafeBuffer)
	log.SetOutput(buff)
	return func() bool {
		if buff.Len() == 0 {
			return false
		}
		_, err := buff.Read(message)
		require.NoError(t, err, "Error reading log message")
		return assert.Contains(t, string(message), expected, "Expected message not found")
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

func ReadFile(t *testing.T, path string) []byte {
	_, currentPackageDir, _, _ := runtime.Caller(0)
	file, err := os.Open(filepath.Join(currentPackageDir, "../../", path))
	require.NoError(t, err)

	defer func(file *os.File) {
		closingErr := file.Close()
		require.NoError(t, closingErr)
	}(file)

	fileBytes, err := io.ReadAll(file)
	require.NoError(t, err)
	return fileBytes
}

// WriteTestFile writes a file to the temp directory with the given name and content, returning the path to the file.
// The file is written with 0644 permissions.
func WriteTestFile(t *testing.T, name string, content []byte) string {
	tempDir := os.TempDir()
	testFile := filepath.Join(tempDir, name)
	if _, err := os.Stat(testFile); err == nil {
		require.NoError(t, os.Remove(testFile))
	}
	err := os.WriteFile(testFile, content, os.FileMode(0644))
	require.NoError(t, err)
	return filepath.Join(tempDir, name)
}

func ResetFlagSet() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}
