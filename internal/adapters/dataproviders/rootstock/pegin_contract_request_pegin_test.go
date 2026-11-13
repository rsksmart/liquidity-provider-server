package rootstock_test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	bindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin"
	commitFirstBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_commit_first"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	requestPegInKnownCallData = "fc73bbd3" +
		"00000000000000000000000079568c2989232dca1840087d73d403602364c0d4" +
		"00000000000000000000000000000000000000000000000000000000000000a0" +
		"1111111111111111111111111111111111111111111111111111111111111111" +
		"0000000000000000000000000000000000000000000000000000000000000005" +
		"00000000000000000000000000000000000000000000000000000000000000e0" +
		"0000000000000000000000000000000000000000000000000000000000000005" +
		"0102030405000000000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000000000000000000000000002" +
		"2222222222222222222222222222222222222222222222222222222222222222" +
		"3333333333333333333333333333333333333333333333333333333333333333"
	revertWordPadding         = "0000000000000000000000000000000000000000000000000000000000000000"
	requestPegInRevertMessage = "execution reverted"
)

type requestPegInHarness struct {
	contract     blockchain.PeginContract
	contractMock boundContractMock
	client       *mocks.RpcClientBindingMock
	signer       *mocks.TransactionSignerMock
	key          *ecdsa.PrivateKey
	from         common.Address
}

func newRequestPegInHarness(t *testing.T) requestPegInHarness {
	t.Helper()
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	contractMock := createBoundContractMock()
	client := &mocks.RpcClientBindingMock{}
	signer := &mocks.TransactionSignerMock{}
	contract := rootstock.NewPeginContractImpl(
		rootstock.NewRskClient(client),
		test.AnyRskAddress,
		contractMock.contract,
		signer,
		rootstock.RetryParams{Retries: 2},
		time.Duration(1),
		bindings.NewPeginContract(),
		commitFirstBindings.NewPeginCommitFirstContract(),
		Abis,
	)
	return requestPegInHarness{
		contract:     contract,
		contractMock: contractMock,
		client:       client,
		signer:       signer,
		key:          key,
		from:         crypto.PubkeyToAddress(key.PublicKey),
	}
}

type sentTransaction struct {
	tx *geth.Transaction
}

func (h requestPegInHarness) expectSignedSend(t *testing.T, gasLimit uint64, value *big.Int) *sentTransaction {
	t.Helper()
	sent := &sentTransaction{}
	h.signer.EXPECT().Address().Return(h.from)
	h.signer.EXPECT().Sign(mock.Anything, mock.Anything).
		RunAndReturn(func(_ common.Address, tx *geth.Transaction) (*geth.Transaction, error) {
			return geth.SignTx(tx, geth.HomesteadSigner{}, h.key)
		}).Once()
	h.contractMock.transactor.EXPECT().SendTransaction(
		mock.Anything,
		matchTransaction(h.contractMock.transactor, common.HexToAddress(test.AnyRskAddress), gasLimit, value, mustDecodeHex(t, requestPegInKnownCallData)),
	).RunAndReturn(func(_ context.Context, tx *geth.Transaction) error {
		sent.tx = tx
		return nil
	}).Once()
	return sent
}

func requestPegInTestParams() blockchain.RequestPegInParams {
	return blockchain.RequestPegInParams{
		RskAddress:         "0x79568c2989232dCa1840087D73d403602364c0D4",
		BtcRawTx:           []byte{1, 2, 3, 4, 5},
		BtcBlockHeaderHash: [32]byte(common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")),
		MerkleBranchPath:   big.NewInt(5),
		MerkleBranchHashes: [][32]byte{
			[32]byte(common.HexToHash("0x2222222222222222222222222222222222222222222222222222222222222222")),
			[32]byte(common.HexToHash("0x3333333333333333333333333333333333333333333333333333333333333333")),
		},
	}
}

func mustDecodeHex(t *testing.T, value string) []byte {
	t.Helper()
	decoded, err := hex.DecodeString(value)
	require.NoError(t, err)
	return decoded
}

func matchRequestPegInCallMsg(t *testing.T, from common.Address, value *big.Int) any {
	expectedData := mustDecodeHex(t, requestPegInKnownCallData)
	contract := common.HexToAddress(test.AnyRskAddress)
	return mock.MatchedBy(func(msg ethereum.CallMsg) bool {
		return msg.From == from && msg.To != nil && *msg.To == contract &&
			bytes.Equal(msg.Data, expectedData) && msg.Value != nil && msg.Value.Cmp(value) == 0
	})
}

func pegInRequestedLog(t *testing.T, emitter common.Address, event commitFirstBindings.PeginCommitFirstContractPegInRequested) *geth.Log {
	t.Helper()
	abiEvent := Abis.PegInCommitFirst.Events[commitFirstBindings.PeginCommitFirstContractPegInRequestedEventName]
	data, err := abiEvent.Inputs.NonIndexed().Pack(event.Amount, event.NetToUser, event.CallSuccess)
	require.NoError(t, err)
	return &geth.Log{
		Address: emitter,
		Topics: []common.Hash{
			abiEvent.ID,
			common.Hash(event.PegInId),
			common.BytesToHash(event.Claimer.Bytes()),
			common.BytesToHash(event.RskAddr.Bytes()),
		},
		Data: data,
	}
}

func requestPegInTestReceipt(status uint64, logs ...*geth.Log) *geth.Receipt {
	return &geth.Receipt{
		Status:            status,
		TxHash:            common.HexToHash("0x7a3c5e8f9b1d2a4c6e8f0a2b4c6d8e0f1a3b5c7d9e1f3a5b7c9d1e3f5a7b9c1d"),
		BlockNumber:       big.NewInt(123),
		BlockHash:         common.HexToHash("0x456"),
		GasUsed:           210000,
		CumulativeGasUsed: 300000,
		EffectiveGasPrice: big.NewInt(60000000),
		Logs:              logs,
	}
}

func TestPeginContractImpl_SimulateRequestPegIn(t *testing.T) {
	value := entities.NewWei(4880000000000000)
	txConfig := blockchain.TransactionConfig{Value: value}

	t.Run("accepts a call that does not revert", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		h.signer.EXPECT().Address().Return(h.from)
		h.client.EXPECT().CallContract(mock.Anything, matchRequestPegInCallMsg(t, h.from, value.AsBigInt()), (*big.Int)(nil)).
			Return(nil, nil).Once()

		require.NoError(t, h.contract.SimulateRequestPegIn(txConfig, requestPegInTestParams()))
		h.client.AssertExpectations(t)
		h.signer.AssertExpectations(t)
	})

	t.Run("names a contract revert the claim has no outcome for", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		h.signer.EXPECT().Address().Return(h.from)
		h.client.EXPECT().CallContract(mock.Anything, mock.Anything, (*big.Int)(nil)).
			Return(nil, NewRskRpcError(requestPegInRevertMessage, "0xf3a8c4ee"+revertWordPadding)).Once()

		require.EqualError(t, h.contract.SimulateRequestPegIn(txConfig, requestPegInTestParams()), "request pegin reverted with: AddressNotRegistered")
		h.client.AssertExpectations(t)
	})

	t.Run("returns a call error that carries no revert", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		h.signer.EXPECT().Address().Return(h.from)
		h.client.EXPECT().CallContract(mock.Anything, mock.Anything, (*big.Int)(nil)).Return(nil, assert.AnError).Once()

		require.ErrorIs(t, h.contract.SimulateRequestPegIn(txConfig, requestPegInTestParams()), assert.AnError)
		h.client.AssertExpectations(t)
	})

	t.Run("rejects an invalid rsk address without calling the node", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		params := requestPegInTestParams()
		params.RskAddress = test.AnyString

		require.ErrorIs(t, h.contract.SimulateRequestPegIn(txConfig, params), blockchain.InvalidAddressError)
		h.client.AssertNotCalled(t, "CallContract")
	})
}

func TestPeginContractImpl_SimulateRequestPegIn_Reverts(t *testing.T) {
	value := entities.NewWei(4880000000000000)
	txConfig := blockchain.TransactionConfig{Value: value}
	revertCases := []struct {
		name     string
		data     string
		expected error
	}{
		{name: "PegInAlreadyProcessed", data: "0x61ef2833" + revertWordPadding, expected: blockchain.PegInAlreadyProcessedError},
		{name: "PegInBelowMinimum", data: "0x2571dcd5" + revertWordPadding + revertWordPadding, expected: blockchain.PegInBelowMinimumError},
		{name: "DepositBelowBridgeMinimum", data: "0x4c93adb8" + revertWordPadding + revertWordPadding, expected: blockchain.DepositBelowBridgeMinimumError},
		{name: "InsufficientConfirmations", data: "0x22c53f1f" + revertWordPadding + revertWordPadding, expected: blockchain.InsufficientConfirmationsError},
		{name: "UnableToGetConfirmations", data: "0xd06e366a" + revertWordPadding, expected: blockchain.UnableToGetConfirmationsError},
		{name: "IncorrectFronting", data: "0x4c0c9bd9" + revertWordPadding + revertWordPadding, expected: blockchain.IncorrectFrontingError},
		{name: "EnforcedPause", data: "0xd93c0665", expected: blockchain.ContractPausedError},
	}
	for _, testCase := range revertCases {
		t.Run("maps the "+testCase.name+" revert", func(t *testing.T) {
			h := newRequestPegInHarness(t)
			h.signer.EXPECT().Address().Return(h.from)
			h.client.EXPECT().CallContract(mock.Anything, matchRequestPegInCallMsg(t, h.from, value.AsBigInt()), (*big.Int)(nil)).
				Return(nil, NewRskRpcError(requestPegInRevertMessage, testCase.data)).Once()

			require.ErrorIs(t, h.contract.SimulateRequestPegIn(txConfig, requestPegInTestParams()), testCase.expected)
			h.client.AssertExpectations(t)
		})
	}
}

func TestPeginContractImpl_EstimateRequestPegInGas(t *testing.T) {
	value := entities.NewWei(4880000000000000)
	txConfig := blockchain.TransactionConfig{Value: value}

	t.Run("returns the node estimate plus a 20% buffer", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		h.signer.EXPECT().Address().Return(h.from)
		h.client.EXPECT().EstimateGas(mock.Anything, matchRequestPegInCallMsg(t, h.from, value.AsBigInt())).
			Return(uint64(250000), nil).Once()

		gas, err := h.contract.EstimateRequestPegInGas(txConfig, requestPegInTestParams())
		require.NoError(t, err)
		assert.Equal(t, uint64(300000), gas)
		h.client.AssertExpectations(t)
	})

	t.Run("returns the estimate error", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		h.signer.EXPECT().Address().Return(h.from)
		h.client.EXPECT().EstimateGas(mock.Anything, mock.Anything).Return(uint64(0), assert.AnError).Once()

		gas, err := h.contract.EstimateRequestPegInGas(txConfig, requestPegInTestParams())
		require.ErrorIs(t, err, assert.AnError)
		assert.Zero(t, gas)
		h.client.AssertExpectations(t)
	})

	t.Run("maps a revert without asking the node again", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		h.signer.EXPECT().Address().Return(h.from)
		h.client.EXPECT().EstimateGas(mock.Anything, mock.Anything).
			Return(uint64(0), NewRskRpcError(requestPegInRevertMessage, "0x61ef2833"+revertWordPadding)).Once()

		gas, err := h.contract.EstimateRequestPegInGas(txConfig, requestPegInTestParams())
		require.ErrorIs(t, err, blockchain.PegInAlreadyProcessedError)
		assert.Zero(t, gas)
		h.client.AssertNumberOfCalls(t, "EstimateGas", 1)
	})

	t.Run("rejects an invalid rsk address without calling the node", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		params := requestPegInTestParams()
		params.RskAddress = test.AnyString

		gas, err := h.contract.EstimateRequestPegInGas(txConfig, params)
		require.ErrorIs(t, err, blockchain.InvalidAddressError)
		assert.Zero(t, gas)
		h.client.AssertNotCalled(t, "EstimateGas")
	})
}

// nolint:funlen
func TestPeginContractImpl_RequestPegIn(t *testing.T) {
	value := entities.NewWei(4880000000000000)
	gasLimit := uint64(300000)
	txConfig := blockchain.TransactionConfig{Value: value, GasLimit: &gasLimit}
	pegInContractAddress := common.HexToAddress(test.AnyRskAddress)
	event := commitFirstBindings.PeginCommitFirstContractPegInRequested{
		PegInId:     [32]byte(common.HexToHash("0x6c0c1c2ba8a9cd4f0a3f0e15cc1d0a1b2c3d4e5f60718293a4b5c6d7e8f90a1b")),
		Claimer:     common.HexToAddress("0x9D93929A9099be4355fC2389FbF253982F9dF47c"),
		RskAddr:     common.HexToAddress("0x79568c2989232dCa1840087D73d403602364c0D4"),
		Amount:      big.NewInt(5000000000000000),
		NetToUser:   value.AsBigInt(),
		CallSuccess: true,
	}

	t.Run("sends the value with the gas limit and returns the receipt and the PegInRequested event", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		sent := h.expectSignedSend(t, gasLimit, value.AsBigInt())
		receipt := requestPegInTestReceipt(1, pegInRequestedLog(t, pegInContractAddress, event))
		h.client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).Return(receipt, nil).Once()

		result, err := h.contract.RequestPegIn(txConfig, requestPegInTestParams())
		require.NoError(t, err)
		expectedReceipt, parseErr := rootstock.ParseReceipt(sent.tx, receipt)
		require.NoError(t, parseErr)
		assert.Equal(t, blockchain.RequestPegInResult{
			Receipt: expectedReceipt,
			Event: blockchain.PegInRequestedEvent{
				PegInID:     "6c0c1c2ba8a9cd4f0a3f0e15cc1d0a1b2c3d4e5f60718293a4b5c6d7e8f90a1b",
				Claimer:     "0x9D93929A9099be4355fC2389FbF253982F9dF47c",
				RskAddress:  "0x79568c2989232dCa1840087D73d403602364c0D4",
				Amount:      entities.NewWei(5000000000000000),
				NetToUser:   value,
				CallSuccess: true,
			},
		}, result)
		assert.Equal(t, h.from.String(), result.Receipt.From)
		h.contractMock.transactor.AssertExpectations(t)
		h.client.AssertExpectations(t)
	})

	t.Run("lets the node estimate the gas when no limit is given", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		h.contractMock.transactor.EXPECT().PendingCodeAt(mock.Anything, pegInContractAddress).Return([]byte{1}, nil).Once()
		h.contractMock.transactor.EXPECT().EstimateGas(mock.Anything, mock.Anything).Return(uint64(280000), nil).Once()
		sent := h.expectSignedSend(t, 0, value.AsBigInt())
		h.client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).
			Return(requestPegInTestReceipt(1, pegInRequestedLog(t, pegInContractAddress, event)), nil).Once()

		_, err := h.contract.RequestPegIn(blockchain.TransactionConfig{Value: value}, requestPegInTestParams())
		require.NoError(t, err)
		assert.Equal(t, uint64(280000), sent.tx.Gas())
		h.contractMock.transactor.AssertExpectations(t)
	})

	t.Run("signs with the given gas price without asking the node for one", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		gasPrice := entities.NewWei(65000000)
		h.signer.EXPECT().Address().Return(h.from)
		h.signer.EXPECT().Sign(mock.Anything, mock.Anything).
			RunAndReturn(func(_ common.Address, tx *geth.Transaction) (*geth.Transaction, error) {
				return geth.SignTx(tx, geth.HomesteadSigner{}, h.key)
			}).Once()
		h.contractMock.transactor.EXPECT().PendingNonceAt(mock.Anything, mock.Anything).Return(uint64(1), nil).Once()
		var sentTx *geth.Transaction
		h.contractMock.transactor.EXPECT().SendTransaction(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, tx *geth.Transaction) error {
				sentTx = tx
				return nil
			}).Once()
		h.client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).
			Return(requestPegInTestReceipt(1, pegInRequestedLog(t, pegInContractAddress, event)), nil).Once()

		_, err := h.contract.RequestPegIn(blockchain.TransactionConfig{Value: value, GasLimit: &gasLimit, GasPrice: gasPrice}, requestPegInTestParams())
		require.NoError(t, err)
		require.NotNil(t, sentTx)
		assert.Equal(t, gasPrice.AsBigInt(), sentTx.GasPrice())
		h.contractMock.transactor.AssertNotCalled(t, "SuggestGasPrice", mock.Anything)
		h.contractMock.transactor.AssertNotCalled(t, "HeaderByNumber", mock.Anything, mock.Anything)
	})

	t.Run("returns no transaction hash when nothing was broadcast", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		h.signer.EXPECT().Address().Return(h.from)
		h.signer.EXPECT().Sign(mock.Anything, mock.Anything).
			RunAndReturn(func(_ common.Address, tx *geth.Transaction) (*geth.Transaction, error) {
				return geth.SignTx(tx, geth.HomesteadSigner{}, h.key)
			}).Once()
		h.contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(h.contractMock.transactor, pegInContractAddress, gasLimit, value.AsBigInt(), mustDecodeHex(t, requestPegInKnownCallData)),
		).Return(assert.AnError).Once()

		result, err := h.contract.RequestPegIn(txConfig, requestPegInTestParams())
		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, blockchain.RequestPegInResult{}, result)
		h.client.AssertNotCalled(t, "TransactionReceipt", mock.Anything, mock.Anything)
	})

	t.Run("keeps the transaction hash when the receipt can't be read", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		sent := h.expectSignedSend(t, gasLimit, value.AsBigInt())
		h.client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		result, err := h.contract.RequestPegIn(txConfig, requestPegInTestParams())
		require.ErrorContains(t, err, "request pegin error")
		assert.Equal(t, blockchain.RequestPegInResult{
			Receipt: blockchain.TransactionReceipt{TransactionHash: sent.tx.Hash().String()},
		}, result)
		h.contractMock.transactor.AssertNumberOfCalls(t, "SendTransaction", 1)
	})

	t.Run("keeps the transaction hash when the node returns no receipt", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		sent := h.expectSignedSend(t, gasLimit, value.AsBigInt())
		h.client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).Return(nil, nil).Once()

		result, err := h.contract.RequestPegIn(txConfig, requestPegInTestParams())
		require.EqualError(t, err, "request pegin error: incomplete receipt")
		assert.Equal(t, blockchain.RequestPegInResult{
			Receipt: blockchain.TransactionReceipt{TransactionHash: sent.tx.Hash().String()},
		}, result)
		h.contractMock.transactor.AssertNumberOfCalls(t, "SendTransaction", 1)
	})

	t.Run("decodes the revert of a mined transaction at its block", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		sent := h.expectSignedSend(t, gasLimit, value.AsBigInt())
		receipt := requestPegInTestReceipt(0)
		h.client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).Return(receipt, nil).Once()
		h.client.EXPECT().CallContract(mock.Anything, matchRequestPegInCallMsg(t, h.from, value.AsBigInt()), big.NewInt(123)).
			Return(nil, NewRskRpcError(requestPegInRevertMessage, "0x61ef2833"+revertWordPadding)).Once()

		result, err := h.contract.RequestPegIn(txConfig, requestPegInTestParams())
		require.ErrorIs(t, err, blockchain.PegInAlreadyProcessedError)
		require.ErrorIs(t, err, blockchain.TxFailedError)
		expectedReceipt, parseErr := rootstock.ParseReceipt(sent.tx, receipt)
		require.NoError(t, parseErr)
		assert.Equal(t, blockchain.RequestPegInResult{Receipt: expectedReceipt}, result)
		h.client.AssertExpectations(t)
	})

	t.Run("reports a mined revert the replay can't explain as a failed transaction", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		h.expectSignedSend(t, gasLimit, value.AsBigInt())
		h.client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).Return(requestPegInTestReceipt(0), nil).Once()
		h.client.EXPECT().CallContract(mock.Anything, mock.Anything, big.NewInt(123)).Return(nil, nil).Once()

		_, err := h.contract.RequestPegIn(txConfig, requestPegInTestParams())
		require.ErrorIs(t, err, blockchain.TxFailedError)
		require.NotErrorIs(t, err, blockchain.PegInAlreadyProcessedError)
		h.client.AssertExpectations(t)
	})

	t.Run("fails when the receipt has no PegInRequested event", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		sent := h.expectSignedSend(t, gasLimit, value.AsBigInt())
		receipt := requestPegInTestReceipt(1)
		h.client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).Return(receipt, nil).Once()

		result, err := h.contract.RequestPegIn(txConfig, requestPegInTestParams())
		require.EqualError(t, err, "request pegin error: PegInRequested event not found")
		expectedReceipt, parseErr := rootstock.ParseReceipt(sent.tx, receipt)
		require.NoError(t, parseErr)
		assert.Equal(t, blockchain.RequestPegInResult{Receipt: expectedReceipt}, result)
	})

	t.Run("ignores a PegInRequested log from another contract", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		h.expectSignedSend(t, gasLimit, value.AsBigInt())
		otherContract := common.HexToAddress("0x5555555555555555555555555555555555555555")
		h.client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).
			Return(requestPegInTestReceipt(1, pegInRequestedLog(t, otherContract, event)), nil).Once()

		_, err := h.contract.RequestPegIn(txConfig, requestPegInTestParams())
		require.EqualError(t, err, "request pegin error: PegInRequested event not found")
	})

	t.Run("returns the unpack error of a malformed PegInRequested log", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		h.expectSignedSend(t, gasLimit, value.AsBigInt())
		malformed := pegInRequestedLog(t, pegInContractAddress, event)
		malformed.Data = []byte{1}
		h.client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).Return(requestPegInTestReceipt(1, malformed), nil).Once()

		result, err := h.contract.RequestPegIn(txConfig, requestPegInTestParams())
		require.ErrorContains(t, err, "request pegin error")
		assert.Equal(t, blockchain.PegInRequestedEvent{}, result.Event)
	})

	t.Run("rejects an invalid rsk address without sending", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		params := requestPegInTestParams()
		params.RskAddress = test.AnyString

		result, err := h.contract.RequestPegIn(txConfig, params)
		require.ErrorIs(t, err, blockchain.InvalidAddressError)
		assert.Equal(t, blockchain.RequestPegInResult{}, result)
		h.contractMock.transactor.AssertNotCalled(t, "SendTransaction", mock.Anything, mock.Anything)
	})
}
