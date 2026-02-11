package btcclient_test

import (
	"fmt"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin/btcclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"unsafe"
)

var (
	rawTx               = "0200000001fef3d290059610a55329af5c7ae7074b4ccc3d19ece35cb04a6bc2c86d9d838f000000006a4730440220340159e3afac48b843b123b744e1b04faea64c274338dfc5f0e056ea4b09290e0220091e5568b6b0180aff47d384bf74b62d646159999a6ea770f0639a0acc07ca5801210232858a5faa413101831afe7a880da9a8ac4de6bd5e25b4358d762ba450b03c22fdffffff0280969800000000001976a914ba99c16de3f21befa3f43c1b0d257121f57a8f3188aceecaf505000000001976a914ddb677f36498f7a4901a74e882df68fd00cf473588ac00000000"
	signWithKeyResponse = fmt.Appendf(nil, "{\n  \"hex\": \"%s\",\n  \"complete\": true\n}", rawTx)
)

func TestFutureSignRawTransactionWithKeyResult_Receive(t *testing.T) {
	channel := make(btcclient.FutureSignRawTransactionWithKeyResult)
	responseBytes := signWithKeyResponse
	res := &rpcclient.Response{}

	// setting value using reflection since is a private field (from the library) and there is no other way to test this
	field := reflect.ValueOf(res).Elem().FieldByName("result")
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(responseBytes))

	go func() { channel <- res }()
	tx, complete, err := channel.Receive()

	require.NoError(t, err)
	assert.True(t, complete)
	assert.NotNil(t, tx)
	assert.Equal(t, "011b335a6c020543b42e801fac566d44e089cfa74ee2df44956be075ddd29691", tx.TxHash().String())
	assert.Len(t, tx.TxOut, 2)
	assert.Len(t, tx.TxIn, 1)
}
