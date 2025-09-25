package rootstock_test

import (
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/mock"
	"math/big"
)

var Abis = rootstock.MustLoadFlyoverABIs()

type RskRpcError struct {
	message string
	data    string
}

func NewRskRpcError(message string, data string) RskRpcError {
	return RskRpcError{
		message: message,
		data:    data,
	}
}

func (r RskRpcError) Error() string {
	return r.message
}

func (r RskRpcError) ErrorData() interface{} {
	return r.data
}

var parsedAddress = common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")

type txModifier func(tx *geth.LegacyTx)

func valueModifier(value *big.Int) txModifier {
	return func(tx *geth.LegacyTx) {
		tx.Value = value
	}
}
func gasLimitModifier(value uint64) txModifier {
	return func(tx *geth.LegacyTx) {
		tx.Gas = value
	}
}

func prepareTxMocks(
	mockClient *mocks.RpcClientBindingMock,
	signerMock *mocks.TransactionSignerMock,
	success bool,
	txModifiers ...txModifier,
) *geth.Transaction {
	legacyTx := &geth.LegacyTx{
		Nonce:    1,
		To:       &parsedAddress,
		Gas:      1,
		GasPrice: big.NewInt(1),
		Data:     nil,
	}

	mockClient.Calls = []mock.Call{}
	mockClient.ExpectedCalls = []*mock.Call{}
	signerMock.Calls = []mock.Call{}
	signerMock.ExpectedCalls = []*mock.Call{}

	for _, modifier := range txModifiers {
		modifier(legacyTx)
	}

	tx := geth.NewTx(legacyTx)

	receipt := &geth.Receipt{}
	receipt.TxHash = tx.Hash()
	if success == true {
		receipt.Status = 1
	}
	mockClient.On("TransactionReceipt", mock.Anything, mock.Anything).Return(receipt, nil).Once()
	signerMock.On("Sign", mock.Anything, mock.Anything).Return(tx, nil).Once()
	signerMock.On("Address").Return(parsedAddress)
	return tx
}
