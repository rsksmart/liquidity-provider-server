package rootstocktest

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var Abis = rootstock.MustLoadFlyoverABIs()

var DummyClient = rootstock.NewRskClient(nil)

var ParsedAddress = common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")

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

type BoundContractMock struct {
	Contract   *bind.BoundContract
	Caller     *mocks.ContractCallerMock
	Transactor *mocks.ContractTransactorMock
	Filterer   *mocks.ContractFiltererMock
}

func NewBoundContractMock() BoundContractMock {
	callerMock := new(mocks.ContractCallerMock)
	transactorMock := new(mocks.ContractTransactorMock)
	filtererMock := new(mocks.ContractFiltererMock)
	contract := bind.NewBoundContract(common.HexToAddress(test.AnyRskAddress), abi.ABI{}, callerMock, transactorMock, filtererMock)
	return BoundContractMock{
		Contract:   contract,
		Caller:     callerMock,
		Transactor: transactorMock,
		Filterer:   filtererMock,
	}
}

func PrepareTxMocks(
	contractMock *BoundContractMock,
	mockClient *mocks.RpcClientBindingMock,
	signerMock *mocks.TransactionSignerMock,
	success bool,
	logs ...*geth.Log,
) {
	mockClient.Calls = []mock.Call{}
	mockClient.ExpectedCalls = []*mock.Call{}
	signerMock.Calls = []mock.Call{}
	signerMock.ExpectedCalls = []*mock.Call{}
	signerMock.EXPECT().Sign(mock.Anything, mock.Anything).RunAndReturn(func(addr common.Address, transaction *geth.Transaction) (*geth.Transaction, error) {
		return transaction, nil
	}).Once()
	receipt := &geth.Receipt{
		TxHash:            common.HexToHash(test.AnyHash),
		BlockNumber:       big.NewInt(123),
		BlockHash:         common.HexToHash("0x456"),
		GasUsed:           21000,
		CumulativeGasUsed: 50000,
		EffectiveGasPrice: big.NewInt(20000000000),
		Logs:              logs,
	}
	if success {
		receipt.Status = 1
	}
	mockClient.On("TransactionReceipt", mock.Anything, mock.Anything).Return(receipt, nil).Once()
	signerMock.On("Address").Return(ParsedAddress)
	contractMock.Transactor.EXPECT().PendingCodeAt(mock.Anything, mock.Anything).Return([]byte{1}, nil).Maybe()
	contractMock.Caller.EXPECT().CodeAt(mock.Anything, mock.Anything, mock.Anything).Return([]byte{1}, nil).Maybe()
	contractMock.Transactor.EXPECT().EstimateGas(mock.Anything, mock.Anything).Return(uint64(1), nil).Maybe()
}

func MatchCallData(expected []byte) any {
	return mock.MatchedBy(func(msg ethereum.CallMsg) bool {
		return bytes.Equal(msg.Data, expected)
	})
}

func MatchTransaction(transactor *mocks.ContractTransactorMock, expectedDestination common.Address, gasLimit uint64, expectedValue *big.Int, expectedData []byte) any {
	transactor.EXPECT().HeaderByNumber(mock.Anything, (*big.Int)(nil)).Return(&geth.Header{}, nil).Once()
	transactor.EXPECT().SuggestGasPrice(mock.Anything).Return(big.NewInt(1), nil).Once()
	transactor.EXPECT().PendingNonceAt(mock.Anything, mock.Anything).Return(uint64(1), nil).Once()
	gasMatch := true
	return mock.MatchedBy(func(tx *geth.Transaction) bool {
		if gasLimit != 0 {
			gasMatch = tx.Gas() == gasLimit
		}
		return bytes.Equal(tx.Data(), expectedData) && tx.Value().Cmp(expectedValue) == 0 && tx.To().Hex() == expectedDestination.Hex() && gasMatch
	})
}

func MustDecodeHex(t *testing.T, value string) []byte {
	t.Helper()
	decoded, err := hex.DecodeString(value)
	require.NoError(t, err)
	return decoded
}

func MustPackUint256(t *testing.T, v *big.Int) []byte {
	t.Helper()
	uint256, err := abi.NewType("uint256", "", nil)
	require.NoError(t, err)

	out, err := abi.Arguments{{Type: uint256}}.Pack(v)
	require.NoError(t, err)
	return out
}

type GeneralPauseStatus struct {
	IsPaused bool
	Reason   string
	Since    uint64
}

func MustPackPauseStatus(t *testing.T, status GeneralPauseStatus) []byte {
	t.Helper()
	uintType, uintErr := abi.NewType("uint64", "", nil)
	boolType, boolErr := abi.NewType("bool", "", nil)
	stringType, stringErr := abi.NewType("string", "", nil)
	if uintErr != nil || boolErr != nil || stringErr != nil {
		t.Fatal("error creating types for pause status packing")
	}
	args := abi.Arguments{
		{Type: boolType},
		{Type: stringType},
		{Type: uintType},
	}
	encoded, err := args.Pack(status.IsPaused, status.Reason, status.Since)
	if err != nil {
		t.Fatal(err)
	}
	return encoded
}
