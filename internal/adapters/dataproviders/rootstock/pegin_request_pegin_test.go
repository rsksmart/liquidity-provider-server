package rootstock_test

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	commitfirst "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_commit_first"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// pinnedRequestPegInABI is the IPegInCommitFirst requestPegIn signature from PR #517,
// parsed independently of the generated packer under test.
const pinnedRequestPegInABI = `[{"type":"function","name":"requestPegIn","stateMutability":"payable","inputs":[{"name":"rskAddr","type":"address"},{"name":"btcTxSerialized","type":"bytes"},{"name":"opReturn","type":"bytes"},{"name":"btcBlockHash","type":"bytes32"},{"name":"merkleBranchPath","type":"uint256"},{"name":"merkleBranchHashes","type":"bytes32[]"}],"outputs":[{"name":"pegInId","type":"bytes32"}]}]`

const requestPegInEstimatedGas = uint64(1000)

var (
	strippedRawTx    = []byte{0x01, 0x00, 0x00, 0x00, 0x01, 0xff, 0xaa, 0xbb}
	witnessRawTx     = []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x01, 0xaa, 0xbb}
	requestBlockHash = [32]byte{0x11}
	requestPath      = big.NewInt(1)
	requestHashes    = [][32]byte{{0x22}}
)

func packPinnedRequestPegIn(t *testing.T, rskAddr common.Address, rawTx []byte, blockHash [32]byte, path *big.Int, hashes [][32]byte) []byte {
	t.Helper()
	parsed, err := abi.JSON(strings.NewReader(pinnedRequestPegInABI))
	require.NoError(t, err)
	calldata, err := parsed.Pack("requestPegIn", rskAddr, rawTx, []byte{}, blockHash, path, hashes)
	require.NoError(t, err)
	inputs := parsed.Methods["requestPegIn"].Inputs
	require.Len(t, inputs, 6)
	for _, input := range inputs {
		assert.NotEqual(t, "amount", input.Name)
		assert.NotEqual(t, "btcTxHash", input.Name)
	}
	return calldata
}

func paddedRequestPegInGas() uint64 {
	return requestPegInEstimatedGas * 12 / 10
}

func stubEstimateRequestPegInGas(mockClient *mocks.RpcClientBindingMock) {
	mockClient.On("EstimateGas", mock.Anything, mock.Anything).Return(requestPegInEstimatedGas, nil).Once()
}

type requestPegInHarness struct {
	contractMock boundContractMock
	signerMock   *mocks.TransactionSignerMock
	mockClient   *mocks.RpcClientBindingMock
	pegin        blockchain.PeginContract
}

func newRequestPegInHarness(t *testing.T) requestPegInHarness {
	t.Helper()
	contractMock := createBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	return requestPegInHarness{
		contractMock: contractMock,
		signerMock:   signerMock,
		mockClient:   mockClient,
		pegin:        newRequestPegInContract(t, contractMock, mockClient, signerMock),
	}
}

func newRequestPegInContract(
	t *testing.T,
	contractMock boundContractMock,
	mockClient *mocks.RpcClientBindingMock,
	signerMock *mocks.TransactionSignerMock,
) blockchain.PeginContract {
	t.Helper()
	return rootstock.NewPeginContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		contractMock.contract,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
		nil,
		Abis,
	)
}

func sampleRequestPegInParams(amount, fee *entities.Wei) blockchain.RequestPegInParams {
	return blockchain.RequestPegInParams{
		RskAddress:         parsedAddress.Hex(),
		BitcoinRawTx:       strippedRawTx,
		BtcBlockHash:       requestBlockHash,
		MerkleBranchPath:   requestPath,
		MerkleBranchHashes: requestHashes,
		Amount:             amount,
		Fee:                fee,
	}
}

func revertHexFromErrorID(t *testing.T, id common.Hash, tail []byte) string {
	t.Helper()
	return "0x" + hex.EncodeToString(append(id.Bytes()[:4], tail...))
}

func mustPegInRequestedLog(t *testing.T, pegInId [32]byte, claimer, rskAddr common.Address, amount, net *big.Int) *geth.Log {
	t.Helper()
	parsed, err := commitfirst.PeginCommitFirstContractMetaData.ParseABI()
	require.NoError(t, err)
	event := parsed.Events["PegInRequested"]
	data, err := event.Inputs.NonIndexed().Pack(amount, net, true)
	require.NoError(t, err)
	return &geth.Log{
		Address: common.HexToAddress(test.AnyRskAddress),
		Topics: []common.Hash{
			event.ID,
			common.BytesToHash(pegInId[:]),
			common.BytesToHash(claimer.Bytes()),
			common.BytesToHash(rskAddr.Bytes()),
		},
		Data: data,
	}
}

func TestPeginContractImpl_RequestPegIn_PackingMatchesPinnedABI(t *testing.T) {
	h := newRequestPegInHarness(t)

	amount := entities.SatoshiToWei(1000)
	fee := entities.SatoshiToWei(100)
	expectedValue := new(entities.Wei).Sub(amount, fee)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	pegInId := [32]byte{0xab}
	eventLog := mustPegInRequestedLog(t, pegInId, parsedAddress, parsedAddress, amount.AsBigInt(), expectedValue.AsBigInt())
	gasLimit := paddedRequestPegInGas()

	h.contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, nil).Once()
	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, expectedValue.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&h.contractMock, h.mockClient, h.signerMock, true, eventLog)
	stubEstimateRequestPegInGas(h.mockClient)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(amount, fee))
	require.NoError(t, err)
	assert.Equal(t, pegInId, result.Event.PegInId)
	assert.Equal(t, expectedValue, result.Receipt.Value)
	assert.True(t, strings.HasPrefix(hex.EncodeToString(expectedData), "a355e935"))
	h.contractMock.transactor.AssertExpectations(t)
	h.mockClient.AssertCalled(t, "EstimateGas", mock.Anything, mock.Anything)
}

func TestPeginContractImpl_RequestPegIn_FirstOutputValueNotSum(t *testing.T) {
	address := "2N2Sg8C2uX1YtugYSxEQvRqf9V2EivxcWER"
	txInfo := blockchain.BitcoinTransactionInformation{
		Outputs: map[string][]*entities.Wei{
			address: {entities.SatoshiToWei(1000), entities.SatoshiToWei(2000)},
		},
	}
	first := txInfo.FirstOutputToAddress(address)
	sum := txInfo.AmountToAddress(address)
	require.Equal(t, entities.SatoshiToWei(1000), first)
	require.Equal(t, entities.SatoshiToWei(3000), sum)

	h := newRequestPegInHarness(t)
	fee := entities.SatoshiToWei(100)
	expectedValue := new(entities.Wei).Sub(first, fee)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	eventLog := mustPegInRequestedLog(t, [32]byte{0x01}, parsedAddress, parsedAddress, first.AsBigInt(), expectedValue.AsBigInt())
	gasLimit := paddedRequestPegInGas()

	h.contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, nil).Once()
	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, expectedValue.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&h.contractMock, h.mockClient, h.signerMock, true, eventLog)
	stubEstimateRequestPegInGas(h.mockClient)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(first, fee))
	require.NoError(t, err)
	assert.Equal(t, expectedValue, result.Receipt.Value)
	assert.NotEqual(t, new(entities.Wei).Sub(sum, fee), result.Receipt.Value)
}

func TestPeginContractImpl_RequestPegIn_SatToWeiBoundary(t *testing.T) {
	oneSat := entities.SatoshiToWei(1)
	require.Equal(t, entities.NewWei(10_000_000_000), oneSat)

	h := newRequestPegInHarness(t)
	fee := entities.NewWei(0)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	eventLog := mustPegInRequestedLog(t, [32]byte{0x02}, parsedAddress, parsedAddress, oneSat.AsBigInt(), oneSat.AsBigInt())
	gasLimit := paddedRequestPegInGas()

	h.contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, nil).Once()
	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, oneSat.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&h.contractMock, h.mockClient, h.signerMock, true, eventLog)
	stubEstimateRequestPegInGas(h.mockClient)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(oneSat, fee))
	require.NoError(t, err)
	assert.Equal(t, oneSat, result.Receipt.Value)
}

func TestPeginContractImpl_RequestPegIn_RejectsWitnessSerializedTx(t *testing.T) {
	h := newRequestPegInHarness(t)
	params := sampleRequestPegInParams(entities.SatoshiToWei(1000), entities.NewWei(0))
	params.BitcoinRawTx = witnessRawTx

	result, err := h.pegin.RequestPegIn(params)
	require.ErrorIs(t, err, blockchain.ErrWitnessSerializedTxNotAccepted)
	assert.Empty(t, result.Receipt.TransactionHash)
	h.contractMock.transactor.AssertNotCalled(t, "SendTransaction")
	h.mockClient.AssertNotCalled(t, "CallContract")
	h.mockClient.AssertNotCalled(t, "EstimateGas")
}

func TestPeginContractImpl_RequestPegIn_StatusZeroDoesNotClassifyRaceLoss(t *testing.T) {
	h := newRequestPegInHarness(t)
	amount := entities.SatoshiToWei(1000)
	fee := entities.NewWei(0)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	gasLimit := paddedRequestPegInGas()

	h.contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, nil).Once()
	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, amount.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&h.contractMock, h.mockClient, h.signerMock, false)
	stubEstimateRequestPegInGas(h.mockClient)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(amount, fee))
	require.Error(t, err)
	require.ErrorContains(t, err, "request pegin error: transaction reverted")
	require.NotErrorIs(t, err, blockchain.ErrPegInAlreadyProcessed)
	assert.NotEmpty(t, result.Receipt.TransactionHash)
	assert.Equal(t, blockchain.PegInRequestedEvent{}, result.Event)
}

func TestPeginContractImpl_RequestPegIn_DoesNotCheckPause(t *testing.T) {
	h := newRequestPegInHarness(t)
	amount := entities.SatoshiToWei(500)
	fee := entities.NewWei(0)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	eventLog := mustPegInRequestedLog(t, [32]byte{0x03}, parsedAddress, parsedAddress, amount.AsBigInt(), amount.AsBigInt())
	gasLimit := paddedRequestPegInGas()

	h.contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, nil).Once()
	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, amount.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&h.contractMock, h.mockClient, h.signerMock, true, eventLog)
	stubEstimateRequestPegInGas(h.mockClient)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(amount, fee))
	require.NoError(t, err)
	assert.NotEmpty(t, result.Receipt.TransactionHash)
	assert.Equal(t, [32]byte{0x03}, result.Event.PegInId)
	h.contractMock.transactor.AssertExpectations(t)
}

func TestPeginContractImpl_RequestPegIn_PreflightAlreadyProcessed(t *testing.T) {
	h := newRequestPegInHarness(t)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	pegInId := [32]byte{0x44}
	selector := commitfirst.PeginCommitFirstContractPegInAlreadyProcessedErrorID().Bytes()[:4]
	revertHex := "0x" + hex.EncodeToString(append(selector, mustPackBytes32(t, pegInId)...))

	h.contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, NewRskRpcError("execution reverted", revertHex)).Once()

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(entities.SatoshiToWei(1000), entities.NewWei(0)))
	require.ErrorIs(t, err, blockchain.ErrPegInAlreadyProcessed)
	assert.Empty(t, result.Receipt.TransactionHash)
	h.contractMock.transactor.AssertNotCalled(t, "SendTransaction")
	h.mockClient.AssertNotCalled(t, "EstimateGas")
}

func TestPeginContractImpl_RequestPegIn_AmountBelowFee(t *testing.T) {
	h := newRequestPegInHarness(t)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(entities.NewWei(1), entities.NewWei(2)))
	require.ErrorIs(t, err, blockchain.ErrIncorrectFronting)
	assert.Empty(t, result.Receipt.TransactionHash)
	h.contractMock.transactor.AssertNotCalled(t, "SendTransaction")
}

func TestPeginContractImpl_RequestPegIn_PreflightTypedErrors(t *testing.T) {
	pegInId := [32]byte{0x44}
	btcTxHash := [32]byte{0x55}
	cases := []struct {
		name    string
		errorID common.Hash
		tail    []byte
		want    error
	}{
		{"AlreadyProcessed", commitfirst.PeginCommitFirstContractPegInAlreadyProcessedErrorID(), mustPackBytes32(t, pegInId), blockchain.ErrPegInAlreadyProcessed},
		{"AddressNotRegistered", commitfirst.PeginCommitFirstContractAddressNotRegisteredErrorID(), mustPackAddress(t, parsedAddress), blockchain.ErrAddressNotRegistered},
		{"DepositOutputNotFound", commitfirst.PeginCommitFirstContractDepositOutputNotFoundErrorID(), append(mustPackAddress(t, parsedAddress), mustPackBytes32(t, btcTxHash)...), blockchain.ErrDepositOutputNotFound},
		{"InsufficientConfirmations", commitfirst.PeginCommitFirstContractInsufficientConfirmationsErrorID(), append(mustPackUint256(t, big.NewInt(1)), mustPackUint256(t, big.NewInt(6))...), blockchain.ErrInsufficientConfirmations},
		{"IncorrectFronting", commitfirst.PeginCommitFirstContractIncorrectFrontingErrorID(), append(mustPackUint256(t, big.NewInt(1000)), mustPackUint256(t, big.NewInt(500))...), blockchain.ErrIncorrectFronting},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertPreflightTypedError(t, tc.errorID, tc.tail, tc.want)
		})
	}
}

func assertPreflightTypedError(t *testing.T, errorID common.Hash, tail []byte, want error) {
	t.Helper()
	h := newRequestPegInHarness(t)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)

	h.contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, NewRskRpcError("execution reverted", revertHexFromErrorID(t, errorID, tail))).Once()

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(entities.SatoshiToWei(1000), entities.NewWei(0)))
	require.ErrorIs(t, err, want)
	assert.Empty(t, result.Receipt.TransactionHash)
	h.contractMock.transactor.AssertNotCalled(t, "SendTransaction")
}

func TestPeginContractImpl_RequestPegIn_PreflightDoesNotInventRaceLoss(t *testing.T) {
	errorTestHex := "0x08c379a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000047465737400000000000000000000000000000000000000000000000000000000"

	t.Run("generic Error(string) revert", func(t *testing.T) {
		assertPreflightJunkRevert(t, NewRskRpcError("execution reverted", errorTestHex), "requestPegIn reverted")
	})
	t.Run("non-DataError", func(t *testing.T) {
		assertPreflightJunkRevert(t, assert.AnError, "error parsing requestPegIn result")
	})
	t.Run("short revert data", func(t *testing.T) {
		for _, revertHex := range []string{"0x", "0xaabbcc"} {
			t.Run(revertHex, func(t *testing.T) {
				assertPreflightJunkRevert(t, NewRskRpcError("execution reverted", revertHex), "requestPegIn reverted")
			})
		}
	})
}

func assertPreflightJunkRevert(t *testing.T, revert error, wantContains string) {
	t.Helper()
	h := newRequestPegInHarness(t)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)

	h.contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, revert).Once()

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(entities.SatoshiToWei(1000), entities.NewWei(0)))
	require.ErrorContains(t, err, wantContains)
	require.NotErrorIs(t, err, blockchain.ErrPegInAlreadyProcessed)
	assert.Empty(t, result.Receipt.TransactionHash)
	h.contractMock.transactor.AssertNotCalled(t, "SendTransaction")
}

func TestPeginContractImpl_RequestPegIn_StatusOneWithoutPegInRequested(t *testing.T) {
	h := newRequestPegInHarness(t)
	amount := entities.SatoshiToWei(1000)
	fee := entities.NewWei(0)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	gasLimit := paddedRequestPegInGas()

	h.contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, nil).Once()
	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, amount.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&h.contractMock, h.mockClient, h.signerMock, true)
	stubEstimateRequestPegInGas(h.mockClient)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(amount, fee))
	require.ErrorContains(t, err, "PegInRequested event not found")
	assert.NotEmpty(t, result.Receipt.TransactionHash)
	assert.Equal(t, blockchain.PegInRequestedEvent{}, result.Event)
}

func TestPeginContractImpl_RequestPegIn_MatchingTopicUnpackError(t *testing.T) {
	h := newRequestPegInHarness(t)
	amount := entities.SatoshiToWei(1000)
	fee := entities.NewWei(0)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	gasLimit := paddedRequestPegInGas()
	parsed, err := commitfirst.PeginCommitFirstContractMetaData.ParseABI()
	require.NoError(t, err)
	badLog := &geth.Log{
		Address: common.HexToAddress(test.AnyRskAddress),
		Topics:  []common.Hash{parsed.Events["PegInRequested"].ID},
		Data:    []byte{0x01},
	}

	h.contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, nil).Once()
	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, amount.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&h.contractMock, h.mockClient, h.signerMock, true, badLog)
	stubEstimateRequestPegInGas(h.mockClient)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(amount, fee))
	require.Error(t, err)
	require.NotContains(t, err.Error(), "PegInRequested event not found")
	assert.Equal(t, blockchain.PegInRequestedEvent{}, result.Event)
}

func TestPeginContractImpl_RequestPegIn_RejectsShortRawTx(t *testing.T) {
	h := newRequestPegInHarness(t)
	params := sampleRequestPegInParams(entities.SatoshiToWei(1000), entities.NewWei(0))
	params.BitcoinRawTx = []byte{1, 0, 0, 0, 1}

	result, err := h.pegin.RequestPegIn(params)
	require.ErrorIs(t, err, blockchain.ErrWitnessSerializedTxNotAccepted)
	assert.Empty(t, result.Receipt.TransactionHash)
	h.contractMock.transactor.AssertNotCalled(t, "SendTransaction")
	h.mockClient.AssertNotCalled(t, "CallContract")
}

func TestPeginContractImpl_RequestPegIn_NilAmountOrFee(t *testing.T) {
	t.Run("nil amount", func(t *testing.T) {
		h := newRequestPegInHarness(t)

		result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(nil, entities.NewWei(0)))
		require.ErrorIs(t, err, blockchain.ErrIncorrectFronting)
		assert.Empty(t, result.Receipt.TransactionHash)
		h.contractMock.transactor.AssertNotCalled(t, "SendTransaction")
	})
	t.Run("nil fee", func(t *testing.T) {
		h := newRequestPegInHarness(t)

		result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(entities.SatoshiToWei(1000), nil))
		require.ErrorIs(t, err, blockchain.ErrIncorrectFronting)
		assert.Empty(t, result.Receipt.TransactionHash)
		h.contractMock.transactor.AssertNotCalled(t, "SendTransaction")
	})
}

func TestPeginContractImpl_RequestPegIn_InvalidAddress(t *testing.T) {
	h := newRequestPegInHarness(t)
	params := sampleRequestPegInParams(entities.SatoshiToWei(1000), entities.NewWei(0))
	params.RskAddress = "not-an-address"

	result, err := h.pegin.RequestPegIn(params)
	require.ErrorIs(t, err, blockchain.InvalidAddressError)
	assert.Empty(t, result.Receipt.TransactionHash)
	h.contractMock.transactor.AssertNotCalled(t, "SendTransaction")
	h.mockClient.AssertNotCalled(t, "CallContract")
}

func TestPeginContractImpl_RequestPegIn_SendError(t *testing.T) {
	h := newRequestPegInHarness(t)
	amount := entities.SatoshiToWei(1000)
	fee := entities.NewWei(0)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	gasLimit := paddedRequestPegInGas()

	h.contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, nil).Once()
	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, amount.AsBigInt(), expectedData),
	).Return(assert.AnError).Once()
	prepareTxMocks(&h.contractMock, h.mockClient, h.signerMock, true)
	stubEstimateRequestPegInGas(h.mockClient)

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(amount, fee))
	require.ErrorContains(t, err, "request pegin error")
	assert.Empty(t, result.Receipt.TransactionHash)
	assert.Equal(t, blockchain.PegInRequestedEvent{}, result.Event)
}

func TestPeginContractImpl_RequestPegIn_EstimateGasFailure(t *testing.T) {
	h := newRequestPegInHarness(t)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)

	h.contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, nil).Once()
	h.contractMock.caller.EXPECT().CodeAt(mock.Anything, mock.Anything, mock.Anything).Return([]byte{1}, nil).Maybe()
	h.signerMock.On("Address").Return(parsedAddress)
	h.mockClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(0), assert.AnError).Once()

	result, err := h.pegin.RequestPegIn(sampleRequestPegInParams(entities.SatoshiToWei(1000), entities.NewWei(0)))
	require.Error(t, err)
	assert.Empty(t, result.Receipt.TransactionHash)
	h.contractMock.transactor.AssertNotCalled(t, "SendTransaction")
}

func TestPeginContractImpl_EstimateRequestPegInGas(t *testing.T) {
	h := newRequestPegInHarness(t)
	h.signerMock.On("Address").Return(parsedAddress)
	h.mockClient.On("EstimateGas", mock.Anything, mock.Anything).Return(requestPegInEstimatedGas, nil).Once()

	gas, err := h.pegin.EstimateRequestPegInGas(sampleRequestPegInParams(entities.SatoshiToWei(1000), entities.NewWei(0)))
	require.NoError(t, err)
	assert.Equal(t, paddedRequestPegInGas(), gas)
	h.contractMock.transactor.AssertNotCalled(t, "SendTransaction")
}

func TestPeginContractImpl_IdentifyRequestPegIn_RejectsWithoutSending(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*blockchain.RequestPegInParams)
		want   error
	}{
		{
			name:   "invalid address",
			mutate: func(params *blockchain.RequestPegInParams) { params.RskAddress = "not-an-address" },
			want:   blockchain.InvalidAddressError,
		},
		{
			name:   "short raw tx",
			mutate: func(params *blockchain.RequestPegInParams) { params.BitcoinRawTx = []byte{1, 0, 0, 0, 1} },
			want:   blockchain.ErrWitnessSerializedTxNotAccepted,
		},
		{
			name:   "witness serialized tx",
			mutate: func(params *blockchain.RequestPegInParams) { params.BitcoinRawTx = witnessRawTx },
			want:   blockchain.ErrWitnessSerializedTxNotAccepted,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newRequestPegInHarness(t)
			params := sampleRequestPegInParams(entities.SatoshiToWei(1000), entities.NewWei(0))
			tc.mutate(&params)

			err := h.pegin.IdentifyRequestPegIn(params)
			require.ErrorIs(t, err, tc.want)
			h.contractMock.transactor.AssertNotCalled(t, "SendTransaction")
			h.mockClient.AssertNotCalled(t, "CallContract")
		})
	}
}

func TestPeginContractImpl_IdentifyRequestPegIn_DoesNotSend(t *testing.T) {
	contractMock := createBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginContract := newRequestPegInContract(t, contractMock, mockClient, signerMock)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)

	contractMock.caller.EXPECT().CodeAt(mock.Anything, mock.Anything, mock.Anything).Return([]byte{1}, nil).Maybe()
	contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, nil).Once()

	err := peginContract.IdentifyRequestPegIn(sampleRequestPegInParams(entities.SatoshiToWei(1000), entities.NewWei(0)))
	require.NoError(t, err)
	contractMock.transactor.AssertNotCalled(t, "SendTransaction")
	mockClient.AssertNotCalled(t, "TransactionReceipt")
}

func TestPeginContractImpl_UnpackPegInRequested(t *testing.T) {
	contractMock := createBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginContract := newRequestPegInContract(t, contractMock, mockClient, signerMock)
	pegInID := [32]byte{0x11, 0x22}
	amount := entities.SatoshiToWei(1000)
	eventLog := mustPegInRequestedLog(t, pegInID, parsedAddress, parsedAddress, amount.AsBigInt(), amount.AsBigInt())
	topics := make([][32]byte, len(eventLog.Topics))
	for i, topic := range eventLog.Topics {
		topics[i] = topic
	}

	event, err := peginContract.UnpackPegInRequested(blockchain.TransactionReceipt{
		Status: blockchain.SuccessfulTxStatus,
		Logs: []blockchain.TransactionLog{{
			Address: eventLog.Address.Hex(),
			Topics:  topics,
			Data:    eventLog.Data,
		}},
	})
	require.NoError(t, err)
	require.NotNil(t, event)
	assert.Equal(t, pegInID, event.PegInId)
	contractMock.transactor.AssertNotCalled(t, "SendTransaction")
}
