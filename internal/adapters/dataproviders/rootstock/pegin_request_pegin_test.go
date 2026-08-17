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

const (
	pauseRegistryABIJSON = `[{"inputs":[],"name":"pauseRegistry","outputs":[{"type":"address"}],"stateMutability":"view","type":"function"}]`
	pauseLevelABIJSON    = `[{"inputs":[],"name":"pauseLevel","outputs":[{"type":"uint8"}],"stateMutability":"view","type":"function"}]`
)

var (
	strippedRawTx     = []byte{0x01, 0x00, 0x00, 0x00, 0x01, 0xff, 0xaa, 0xbb}
	witnessRawTx      = []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x01, 0xaa, 0xbb}
	requestBlockHash  = [32]byte{0x11}
	requestPath       = big.NewInt(1)
	requestHashes     = [][32]byte{{0x22}}
	pauseRegistryAddr = common.HexToAddress("0x00000000000000000000000000000000000000b1")
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

func packPauseRegistry(t *testing.T) []byte {
	t.Helper()
	parsed, err := abi.JSON(strings.NewReader(pauseRegistryABIJSON))
	require.NoError(t, err)
	calldata, err := parsed.Pack("pauseRegistry")
	require.NoError(t, err)
	return calldata
}

func packPauseLevel(t *testing.T) []byte {
	t.Helper()
	parsed, err := abi.JSON(strings.NewReader(pauseLevelABIJSON))
	require.NoError(t, err)
	calldata, err := parsed.Pack("pauseLevel")
	require.NoError(t, err)
	return calldata
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

func stubPauseRegistry(t *testing.T, contractMock *boundContractMock) {
	t.Helper()
	contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(packPauseRegistry(t)),
		mock.Anything,
	).Return(mustPackAddress(t, pauseRegistryAddr), nil).Once()
}

func stubPauseLevelOnClient(t *testing.T, mockClient *mocks.RpcClientBindingMock, level uint8) {
	t.Helper()
	mockClient.On("CallContract", mock.Anything, matchCallData(packPauseLevel(t)), mock.Anything).
		Return(mustPackUint8(t, level), nil).Once()
}

func stubPauseLevel(
	t *testing.T,
	contractMock *boundContractMock,
	mockClient *mocks.RpcClientBindingMock,
	level uint8,
) {
	t.Helper()
	stubPauseRegistry(t, contractMock)
	stubPauseLevelOnClient(t, mockClient, level)
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
	contractMock := createBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginContract := newRequestPegInContract(t, contractMock, mockClient, signerMock)

	amount := entities.SatoshiToWei(1000)
	fee := entities.SatoshiToWei(100)
	expectedValue := new(entities.Wei).Sub(amount, fee)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	pegInId := [32]byte{0xab}
	eventLog := mustPegInRequestedLog(t, pegInId, parsedAddress, parsedAddress, amount.AsBigInt(), expectedValue.AsBigInt())

	stubPauseRegistry(t, &contractMock)
	contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, nil).Once()
	contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 2500000, expectedValue.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&contractMock, mockClient, signerMock, true, eventLog)
	stubPauseLevelOnClient(t, mockClient, blockchain.PauseLevelSoft)

	result, err := peginContract.RequestPegIn(sampleRequestPegInParams(amount, fee))
	require.NoError(t, err)
	assert.Equal(t, pegInId, result.Event.PegInId)
	assert.Equal(t, expectedValue, result.Receipt.Value)
	assert.True(t, strings.HasPrefix(hex.EncodeToString(expectedData), "a355e935"))
	contractMock.transactor.AssertExpectations(t)
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

	contractMock := createBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginContract := newRequestPegInContract(t, contractMock, mockClient, signerMock)
	fee := entities.SatoshiToWei(100)
	expectedValue := new(entities.Wei).Sub(first, fee)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	eventLog := mustPegInRequestedLog(t, [32]byte{0x01}, parsedAddress, parsedAddress, first.AsBigInt(), expectedValue.AsBigInt())

	stubPauseRegistry(t, &contractMock)
	contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, nil).Once()
	contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 2500000, expectedValue.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&contractMock, mockClient, signerMock, true, eventLog)
	stubPauseLevelOnClient(t, mockClient, blockchain.PauseLevelNone)

	result, err := peginContract.RequestPegIn(sampleRequestPegInParams(first, fee))
	require.NoError(t, err)
	assert.Equal(t, expectedValue, result.Receipt.Value)
	assert.NotEqual(t, new(entities.Wei).Sub(sum, fee), result.Receipt.Value)
}

func TestPeginContractImpl_RequestPegIn_SatToWeiBoundary(t *testing.T) {
	oneSat := entities.SatoshiToWei(1)
	require.Equal(t, entities.NewWei(10_000_000_000), oneSat)

	contractMock := createBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginContract := newRequestPegInContract(t, contractMock, mockClient, signerMock)
	fee := entities.NewWei(0)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	eventLog := mustPegInRequestedLog(t, [32]byte{0x02}, parsedAddress, parsedAddress, oneSat.AsBigInt(), oneSat.AsBigInt())

	stubPauseRegistry(t, &contractMock)
	contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, nil).Once()
	contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 2500000, oneSat.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&contractMock, mockClient, signerMock, true, eventLog)
	stubPauseLevelOnClient(t, mockClient, blockchain.PauseLevelNone)

	result, err := peginContract.RequestPegIn(sampleRequestPegInParams(oneSat, fee))
	require.NoError(t, err)
	assert.Equal(t, oneSat, result.Receipt.Value)
}

func TestPeginContractImpl_RequestPegIn_RejectsWitnessSerializedTx(t *testing.T) {
	contractMock := createBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginContract := newRequestPegInContract(t, contractMock, mockClient, signerMock)
	params := sampleRequestPegInParams(entities.SatoshiToWei(1000), entities.NewWei(0))
	params.BitcoinRawTx = witnessRawTx

	result, err := peginContract.RequestPegIn(params)
	require.ErrorIs(t, err, blockchain.ErrWitnessSerializedTxNotAccepted)
	assert.Empty(t, result.Receipt.TransactionHash)
	contractMock.transactor.AssertNotCalled(t, "SendTransaction")
	mockClient.AssertNotCalled(t, "CallContract")
}

func TestPeginContractImpl_RequestPegIn_StatusZeroDoesNotClassifyRaceLoss(t *testing.T) {
	contractMock := createBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginContract := newRequestPegInContract(t, contractMock, mockClient, signerMock)
	amount := entities.SatoshiToWei(1000)
	fee := entities.NewWei(0)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)

	stubPauseRegistry(t, &contractMock)
	contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, nil).Once()
	contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 2500000, amount.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&contractMock, mockClient, signerMock, false)
	stubPauseLevelOnClient(t, mockClient, blockchain.PauseLevelNone)

	result, err := peginContract.RequestPegIn(sampleRequestPegInParams(amount, fee))
	require.Error(t, err)
	require.ErrorContains(t, err, "request pegin error: transaction reverted")
	require.NotErrorIs(t, err, blockchain.ErrPegInAlreadyProcessed)
	assert.NotEmpty(t, result.Receipt.TransactionHash)
	assert.Nil(t, result.Event)
}

func TestPeginContractImpl_RequestPegIn_SoftPauseStillSends(t *testing.T) {
	contractMock := createBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginContract := newRequestPegInContract(t, contractMock, mockClient, signerMock)
	amount := entities.SatoshiToWei(500)
	fee := entities.NewWei(0)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	eventLog := mustPegInRequestedLog(t, [32]byte{0x03}, parsedAddress, parsedAddress, amount.AsBigInt(), amount.AsBigInt())

	stubPauseRegistry(t, &contractMock)
	contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, nil).Once()
	contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 2500000, amount.AsBigInt(), expectedData),
	).Return(nil).Once()
	prepareTxMocks(&contractMock, mockClient, signerMock, true, eventLog)
	stubPauseLevelOnClient(t, mockClient, blockchain.PauseLevelSoft)

	result, err := peginContract.RequestPegIn(sampleRequestPegInParams(amount, fee))
	require.NoError(t, err)
	assert.NotEmpty(t, result.Receipt.TransactionHash)
	contractMock.transactor.AssertExpectations(t)
}

func TestPeginContractImpl_RequestPegIn_HardPauseBlocksSend(t *testing.T) {
	contractMock := createBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginContract := newRequestPegInContract(t, contractMock, mockClient, signerMock)

	stubPauseLevel(t, &contractMock, mockClient, blockchain.PauseLevelHard)
	result, err := peginContract.RequestPegIn(sampleRequestPegInParams(entities.SatoshiToWei(500), entities.NewWei(0)))
	require.ErrorIs(t, err, blockchain.ErrHardPaused)
	assert.Empty(t, result.Receipt.TransactionHash)
	contractMock.transactor.AssertNotCalled(t, "SendTransaction")
}

func TestPeginContractImpl_IsHardPaused(t *testing.T) {
	t.Run("soft is not a send block", func(t *testing.T) {
		contractMock := createBoundContractMock()
		mockClient := &mocks.RpcClientBindingMock{}
		peginContract := newRequestPegInContract(t, contractMock, mockClient, &mocks.TransactionSignerMock{})
		stubPauseLevel(t, &contractMock, mockClient, blockchain.PauseLevelSoft)
		hard, err := peginContract.IsHardPaused()
		require.NoError(t, err)
		assert.False(t, hard)
	})
	t.Run("hard is a send block", func(t *testing.T) {
		contractMock := createBoundContractMock()
		mockClient := &mocks.RpcClientBindingMock{}
		peginContract := newRequestPegInContract(t, contractMock, mockClient, &mocks.TransactionSignerMock{})
		stubPauseLevel(t, &contractMock, mockClient, blockchain.PauseLevelHard)
		hard, err := peginContract.IsHardPaused()
		require.NoError(t, err)
		assert.True(t, hard)
	})
}

func TestPeginContractImpl_RequestPegIn_PreflightAlreadyProcessed(t *testing.T) {
	contractMock := createBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginContract := newRequestPegInContract(t, contractMock, mockClient, signerMock)
	expectedData := packPinnedRequestPegIn(t, parsedAddress, strippedRawTx, requestBlockHash, requestPath, requestHashes)
	pegInId := [32]byte{0x44}
	selector := commitfirst.PeginCommitFirstContractPegInAlreadyProcessedErrorID().Bytes()[:4]
	revertHex := "0x" + hex.EncodeToString(append(selector, mustPackBytes32(t, pegInId)...))

	stubPauseLevel(t, &contractMock, mockClient, blockchain.PauseLevelNone)
	contractMock.caller.EXPECT().CallContract(
		mock.Anything,
		matchCallData(expectedData),
		mock.Anything,
	).Return(nil, NewRskRpcError("execution reverted", revertHex)).Once()

	result, err := peginContract.RequestPegIn(sampleRequestPegInParams(entities.SatoshiToWei(1000), entities.NewWei(0)))
	require.ErrorIs(t, err, blockchain.ErrPegInAlreadyProcessed)
	assert.Empty(t, result.Receipt.TransactionHash)
	contractMock.transactor.AssertNotCalled(t, "SendTransaction")
}

func TestPeginContractImpl_RequestPegIn_AmountBelowFee(t *testing.T) {
	contractMock := createBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginContract := newRequestPegInContract(t, contractMock, mockClient, signerMock)
	stubPauseLevel(t, &contractMock, mockClient, blockchain.PauseLevelNone)

	result, err := peginContract.RequestPegIn(sampleRequestPegInParams(entities.NewWei(1), entities.NewWei(2)))
	require.ErrorIs(t, err, blockchain.ErrIncorrectFronting)
	assert.Empty(t, result.Receipt.TransactionHash)
	contractMock.transactor.AssertNotCalled(t, "SendTransaction")
}
