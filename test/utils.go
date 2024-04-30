package test

import (
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/account"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"path/filepath"
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
	account, err := account.GetRskAccount(account.CreationArgs{
		KeyDir:        testDir,
		AccountNum:    0,
		EncryptedJson: string(keyBytes),
		Password:      KeyPassword,
	})
	require.NoError(t, err)
	return account
}
