package rootstock_test

import (
	"bytes"
	"encoding/hex"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	discoveryBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/discovery"
	pegoutBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
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

type boundContractMock struct {
	contract   *bind.BoundContract
	caller     *mocks.ContractCallerMock
	transactor *mocks.ContractTransactorMock
	filterer   *mocks.ContractFiltererMock
}

func prepareTxMocks(
	contractMock *boundContractMock,
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
	if success == true {
		receipt.Status = 1
	}
	mockClient.On("TransactionReceipt", mock.Anything, mock.Anything).Return(receipt, nil).Once()
	signerMock.On("Address").Return(parsedAddress)
	contractMock.transactor.EXPECT().PendingCodeAt(mock.Anything, mock.Anything).Return([]byte{1}, nil).Maybe()
	contractMock.caller.EXPECT().CodeAt(mock.Anything, mock.Anything, mock.Anything).Return([]byte{1}, nil).Maybe()
	contractMock.transactor.EXPECT().EstimateGas(mock.Anything, mock.Anything).Return(uint64(1), nil).Maybe()
}

func matchCallData(expected []byte) any {
	return mock.MatchedBy(func(msg ethereum.CallMsg) bool {
		return bytes.Equal(msg.Data, expected)
	})
}

func filterMatchFunc(from uint64, to uint64) func(ethereum.FilterQuery) bool {
	return func(opts ethereum.FilterQuery) bool {
		return from == opts.FromBlock.Uint64() && to == opts.ToBlock.Uint64()
	}
}

func matchTransaction(transactor *mocks.ContractTransactorMock, expectedDestination common.Address, gasLimit uint64, expectedValue *big.Int, expectedData []byte) any {
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

func mustPackUint256(t *testing.T, v *big.Int) []byte {
	t.Helper()
	uint256, err := abi.NewType("uint256", "", nil)
	require.NoError(t, err)

	out, err := abi.Arguments{{Type: uint256}}.Pack(v)
	require.NoError(t, err)
	return out
}

func mustPackBool(t *testing.T, v bool) []byte {
	t.Helper()
	boolType, err := abi.NewType("bool", "", nil)
	require.NoError(t, err)

	out, err := abi.Arguments{{Type: boolType}}.Pack(v)
	require.NoError(t, err)
	return out
}

func mustPackString(t *testing.T, v string) []byte {
	t.Helper()
	strType, err := abi.NewType("string", "", nil)
	require.NoError(t, err)

	out, err := abi.Arguments{{Type: strType}}.Pack(v)
	require.NoError(t, err)
	return out
}

func mustPackBytes(t *testing.T, v string) []byte {
	t.Helper()
	strType, err := abi.NewType("bytes", "", nil)
	require.NoError(t, err)

	parsedValue, err := hex.DecodeString(v)
	require.NoError(t, err)

	out, err := abi.Arguments{{Type: strType}}.Pack(parsedValue)
	require.NoError(t, err)
	return out
}

func mustPackBytes32(t *testing.T, v [32]byte) []byte {
	t.Helper()
	byteType, err := abi.NewType("bytes32", "", nil)
	require.NoError(t, err)

	out, err := abi.Arguments{{Type: byteType}}.Pack(v)
	require.NoError(t, err)
	return out
}

func mustPackLiquidityProvider(t *testing.T, provider discoveryBindings.FlyoverLiquidityProvider) []byte {
	t.Helper()
	lpType, err := abi.NewType("tuple", "structFlyover.LiquidityProvider", []abi.ArgumentMarshaling{
		{Name: "id", Type: "uint256"},
		{Name: "providerAddress", Type: "address"},
		{Name: "status", Type: "bool"},
		{Name: "providerType", Type: "uint8"},
		{Name: "name", Type: "string"},
		{Name: "apiBaseUrl", Type: "string"},
	})
	if err != nil {
		t.Fatal(err)
	}
	args := abi.Arguments{{Type: lpType}}
	encoded, err := args.Pack(provider)
	if err != nil {
		t.Fatal(err)
	}
	return encoded
}

func mustPackLiquidityProviders(t *testing.T, provider []discoveryBindings.FlyoverLiquidityProvider) []byte {
	t.Helper()
	lpType, err := abi.NewType("tuple[]", "structFlyover.LiquidityProvider[]", []abi.ArgumentMarshaling{
		{Name: "id", Type: "uint256"},
		{Name: "providerAddress", Type: "address"},
		{Name: "status", Type: "bool"},
		{Name: "providerType", Type: "uint8"},
		{Name: "name", Type: "string"},
		{Name: "apiBaseUrl", Type: "string"},
	})
	if err != nil {
		t.Fatal(err)
	}
	args := abi.Arguments{{Type: lpType}}
	encoded, err := args.Pack(provider)
	if err != nil {
		t.Fatal(err)
	}
	return encoded
}

func mustPackPegoutQuote(t *testing.T, q pegoutBindings.QuotesPegOutQuote) []byte {
	t.Helper()
	quoteType, err := abi.NewType("tuple", "structQuotes.PegOutQuote", []abi.ArgumentMarshaling{
		{Name: "callFee", Type: "uint256"},
		{Name: "penaltyFee", Type: "uint256"},
		{Name: "value", Type: "uint256"},
		{Name: "productFeeAmount", Type: "uint256"},
		{Name: "gasFee", Type: "uint256"},
		{Name: "lbcAddress", Type: "address"},
		{Name: "lpRskAddress", Type: "address"},
		{Name: "rskRefundAddress", Type: "address"},
		{Name: "nonce", Type: "int64"},
		{Name: "agreementTimestamp", Type: "uint32"},
		{Name: "depositDateLimit", Type: "uint32"},
		{Name: "transferTime", Type: "uint32"},
		{Name: "expireDate", Type: "uint32"},
		{Name: "expireBlock", Type: "uint32"},
		{Name: "depositConfirmations", Type: "uint16"},
		{Name: "transferConfirmations", Type: "uint16"},
		{Name: "depositAddress", Type: "bytes"},
		{Name: "btcRefundAddress", Type: "bytes"},
		{Name: "lpBtcAddress", Type: "bytes"},
	})
	if err != nil {
		t.Fatal(err)
	}
	args := abi.Arguments{{Type: quoteType}}
	encoded, err := args.Pack(q)
	if err != nil {
		t.Fatal(err)
	}
	return encoded
}

type generalPauseStatus struct {
	IsPaused bool
	Reason   string
	Since    uint64
}

func mustPackPauseStatus(t *testing.T, status generalPauseStatus) []byte {
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

func createBoundContractMock() boundContractMock {
	callerMock := new(mocks.ContractCallerMock)
	transactorMock := new(mocks.ContractTransactorMock)
	filtererMock := new(mocks.ContractFiltererMock)
	contract := bind.NewBoundContract(common.HexToAddress(test.AnyRskAddress), abi.ABI{}, callerMock, transactorMock, filtererMock)
	return boundContractMock{
		contract:   contract,
		caller:     callerMock,
		transactor: transactorMock,
		filterer:   filtererMock,
	}
}
