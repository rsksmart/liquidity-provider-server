package test

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/account"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strings"
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
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	for i := 0; i < value.NumField(); i++ {
		if !value.Field(i).IsZero() {
			count++
		}
	}
	return count
}

func AssertNonZeroValues(t *testing.T, aStruct any) {
	structType := reflect.TypeOf(aStruct)
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}
	require.Equal(t, structType.NumField(), CountNonZeroValues(aStruct))
}

func AssertMaxZeroValues(t *testing.T, aStruct any, max int) {
	structType := reflect.TypeOf(aStruct)
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}
	require.LessOrEqual(t, structType.NumField()-CountNonZeroValues(aStruct), max)
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
			t.Errorf("No log message found")
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

func AddDepositLogFromQuote(
	t *testing.T,
	receipt *blockchain.TransactionReceipt,
	pegoutQuote quote.PegoutQuote,
	retainedQuote quote.RetainedPegoutQuote,
) *blockchain.TransactionReceipt {
	const depositTopic = "b1bc7bfc0dab19777eb03aa0a5643378fc9f186c8fc5a36620d21136fbea570f"
	parsedDepositTopic, err := hex.DecodeString(depositTopic)
	require.NoError(t, err)
	quoteHashTopic, err := hex.DecodeString(retainedQuote.QuoteHash)
	require.NoError(t, err)
	senderTopic, err := hex.DecodeString(strings.TrimPrefix(receipt.From, "0x"))
	require.NoError(t, err)
	timestampHex := fmt.Sprintf("%064x", uint64(pegoutQuote.DepositDateLimit-500))
	timestampTopic, err := hex.DecodeString(timestampHex)
	require.NoError(t, err)

	amountHex := fmt.Sprintf("%064x", pegoutQuote.Total().AsBigInt())
	parsedData, err := hex.DecodeString(amountHex)
	require.NoError(t, err)

	log := blockchain.TransactionLog{
		Address: pegoutQuote.LbcAddress,
		Topics: [][32]byte{
			utils.To32Bytes(parsedDepositTopic),
			utils.To32Bytes(quoteHashTopic),
			utils.To32Bytes(senderTopic),
			utils.To32Bytes(timestampTopic),
		},
		Data:        parsedData,
		BlockNumber: receipt.BlockNumber,
		TxHash:      receipt.TransactionHash,
		TxIndex:     0,
		BlockHash:   receipt.BlockHash,
		Index:       0,
		Removed:     false,
	}
	receipt.Logs = slices.Insert(receipt.Logs, 0, log)
	return receipt
}

func GetBitcoinTestBlock(t *testing.T, path string) *btcutil.Block {
	absolutePath, err := filepath.Abs(path)
	require.NoError(t, err)
	blockFile, err := os.ReadFile(absolutePath)
	require.NoError(t, err)
	blockBytes, err := hex.DecodeString(string(blockFile))
	require.NoError(t, err)
	block, err := btcutil.NewBlockFromBytes(blockBytes)
	require.NoError(t, err)
	return block
}

func MustReadFileString(path string) string {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("MustRead: could not determine caller info")
	}
	baseDir := filepath.Dir(thisFile)
	fullPath := filepath.Join(baseDir, path)
	b, err := os.ReadFile(fullPath)
	if err != nil {
		panic(fmt.Errorf("MustRead: failed to read %q: %w", path, err))
	}
	return string(b)
}
