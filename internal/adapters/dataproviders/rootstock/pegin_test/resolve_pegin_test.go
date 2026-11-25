package pegin_test

import (
	"bytes"
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	commitFirstBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_commit_first"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/rootstocktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	resolvePegInKnownCallData = "1747d75e" +
		"00000000000000000000000079568c2989232dca1840087d73d403602364c0d4" +
		"0000000000000000000000000000000000000000000000000000000000000080" +
		"00000000000000000000000000000000000000000000000000000000000000c0" +
		"00000000000000000000000000000000000000000000000000000000000001f4" +
		"0000000000000000000000000000000000000000000000000000000000000005" +
		"0102030405000000000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000000000000000000000000003" +
		"0607080000000000000000000000000000000000000000000000000000000000"
	resolvePegInGasLimit         = uint64(2500000)
	bridgeValidationFailedRevert = "0x51a6e280"
)

func resolvePegInTestParams() blockchain.ResolvePegInParams {
	return blockchain.ResolvePegInParams{
		RskAddress:        "0x79568c2989232dCa1840087D73d403602364c0D4",
		BtcRawTx:          []byte{1, 2, 3, 4, 5},
		PartialMerkleTree: []byte{6, 7, 8},
		BlockHeight:       big.NewInt(500),
	}
}

func expectSignedResolve(t *testing.T, h requestPegInHarness) *sentTransaction {
	t.Helper()
	sent := &sentTransaction{}
	h.signer.EXPECT().Address().Return(h.from)
	h.signer.EXPECT().Sign(mock.Anything, mock.Anything).
		RunAndReturn(func(_ common.Address, tx *geth.Transaction) (*geth.Transaction, error) {
			return geth.SignTx(tx, geth.HomesteadSigner{}, h.key)
		}).Once()
	h.contractMock.Transactor.EXPECT().SendTransaction(
		mock.Anything,
		rootstocktest.MatchTransaction(h.contractMock.Transactor, common.HexToAddress(test.AnyRskAddress), resolvePegInGasLimit, big.NewInt(0), rootstocktest.MustDecodeHex(t, resolvePegInKnownCallData)),
	).RunAndReturn(func(_ context.Context, tx *geth.Transaction) error {
		sent.tx = tx
		return nil
	}).Once()
	return sent
}

func matchResolvePegInCallMsg(t *testing.T, from common.Address) any {
	expectedData := rootstocktest.MustDecodeHex(t, resolvePegInKnownCallData)
	contract := common.HexToAddress(test.AnyRskAddress)
	return mock.MatchedBy(func(msg ethereum.CallMsg) bool {
		return msg.From == from && msg.To != nil && *msg.To == contract &&
			bytes.Equal(msg.Data, expectedData) && msg.Value == nil
	})
}

func pegInResolvedTestEvent() commitFirstBindings.PeginCommitFirstContractPegInResolved {
	return commitFirstBindings.PeginCommitFirstContractPegInResolved{
		PegInId:       [32]byte(common.HexToHash("0x6c0c1c2ba8a9cd4f0a3f0e15cc1d0a1b2c3d4e5f60718293a4b5c6d7e8f90a1b")),
		Claimer:       common.HexToAddress("0x9D93929A9099be4355fC2389FbF253982F9dF47c"),
		Registrant:    common.HexToAddress("0x4a9D4b6E3C9f2D1B0aC5E7f8B9C0d1e2F3a4b5c6"),
		Released:      big.NewInt(5000000000000000),
		ClaimerPayout: big.NewInt(4970000000000000),
		RegistrantFee: big.NewInt(10000000000000),
		UserPayout:    big.NewInt(20000000000000),
	}
}

func pegInResolvedTransactionLog(t *testing.T, emitter common.Address, event commitFirstBindings.PeginCommitFirstContractPegInResolved) blockchain.TransactionLog {
	t.Helper()
	abiEvent := rootstocktest.Abis.PegInCommitFirst.Events[commitFirstBindings.PeginCommitFirstContractPegInResolvedEventName]
	data, err := abiEvent.Inputs.NonIndexed().Pack(event.Released, event.ClaimerPayout, event.RegistrantFee, event.UserPayout)
	require.NoError(t, err)
	return blockchain.TransactionLog{
		Address: emitter.String(),
		Topics: [][32]byte{
			abiEvent.ID,
			event.PegInId,
			common.BytesToHash(event.Claimer.Bytes()),
			common.BytesToHash(event.Registrant.Bytes()),
		},
		Data: data,
	}
}

func TestPeginContractImpl_SimulateResolvePegIn(t *testing.T) {
	t.Run("accepts a call that does not revert", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		h.signer.EXPECT().Address().Return(h.from)
		h.client.EXPECT().CallContract(mock.Anything, matchResolvePegInCallMsg(t, h.from), (*big.Int)(nil)).
			Return(nil, nil).Once()

		require.NoError(t, h.contract.SimulateResolvePegIn(resolvePegInTestParams()))
		h.client.AssertExpectations(t)
		h.signer.AssertExpectations(t)
	})

	t.Run("names a contract revert the resolve has no outcome for", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		h.signer.EXPECT().Address().Return(h.from)
		h.client.EXPECT().CallContract(mock.Anything, mock.Anything, (*big.Int)(nil)).
			Return(nil, rootstocktest.NewRskRpcError(requestPegInRevertMessage, "0x913cecc8"+revertWordPadding)).Once()

		require.EqualError(t, h.contract.SimulateResolvePegIn(resolvePegInTestParams()), "resolve pegin reverted with: PegInNotClaimed")
		h.client.AssertExpectations(t)
	})

	t.Run("returns a call error that carries no revert", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		h.signer.EXPECT().Address().Return(h.from)
		h.client.EXPECT().CallContract(mock.Anything, mock.Anything, (*big.Int)(nil)).Return(nil, assert.AnError).Once()

		err := h.contract.SimulateResolvePegIn(resolvePegInTestParams())
		require.ErrorIs(t, err, assert.AnError)
		require.ErrorContains(t, err, "resolve pegin call error: ")
		h.client.AssertExpectations(t)
	})

	t.Run("rejects an invalid rsk address without calling the node", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		params := resolvePegInTestParams()
		params.RskAddress = test.AnyString

		require.ErrorIs(t, h.contract.SimulateResolvePegIn(params), blockchain.InvalidAddressError)
	})

	t.Run("rejects a missing block height without calling the node", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		params := resolvePegInTestParams()
		params.BlockHeight = nil

		require.EqualError(t, h.contract.SimulateResolvePegIn(params), "resolve pegin error: missing block height")
	})
}

func TestPeginContractImpl_SimulateResolvePegIn_Reverts(t *testing.T) {
	revertCases := []struct {
		name     string
		data     string
		expected error
	}{
		{name: "PegInAlreadyProcessed", data: "0x61ef2833" + revertWordPadding, expected: blockchain.PegInAlreadyProcessedError},
		{name: "BridgeValidationFailed", data: bridgeValidationFailedRevert, expected: blockchain.WaitingForBridgeError},
		{name: "UnexpectedBridgeError", data: "0xab2e19b0" + revertWordPadding, expected: blockchain.BridgeRejectedPegInError},
		{name: "EnforcedPause", data: "0xd93c0665", expected: blockchain.ContractPausedError},
	}
	for _, testCase := range revertCases {
		t.Run("maps the "+testCase.name+" revert", func(t *testing.T) {
			h := newRequestPegInHarness(t)
			h.signer.EXPECT().Address().Return(h.from)
			h.client.EXPECT().CallContract(mock.Anything, matchResolvePegInCallMsg(t, h.from), (*big.Int)(nil)).
				Return(nil, rootstocktest.NewRskRpcError(requestPegInRevertMessage, testCase.data)).Once()

			require.ErrorIs(t, h.contract.SimulateResolvePegIn(resolvePegInTestParams()), testCase.expected)
			h.client.AssertExpectations(t)
		})
	}
}

func TestPeginContractImpl_ResolvePegIn(t *testing.T) {
	t.Run("sends the resolve with the fixed gas limit and no value, and returns the receipt", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		sent := expectSignedResolve(t, h)
		receipt := requestPegInTestReceipt(1)
		h.client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).Return(receipt, nil).Once()

		result, err := h.contract.ResolvePegIn(resolvePegInTestParams())
		require.NoError(t, err)
		expectedReceipt, parseErr := rootstock.ParseReceipt(sent.tx, receipt)
		require.NoError(t, parseErr)
		assert.Equal(t, expectedReceipt, result)
		assert.Equal(t, resolvePegInGasLimit, sent.tx.Gas())
		assert.Zero(t, sent.tx.Value().Sign())
		h.contractMock.Transactor.AssertExpectations(t)
		h.client.AssertExpectations(t)
	})

	t.Run("returns no transaction hash when nothing was broadcast", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		h.signer.EXPECT().Address().Return(h.from)
		h.signer.EXPECT().Sign(mock.Anything, mock.Anything).
			RunAndReturn(func(_ common.Address, tx *geth.Transaction) (*geth.Transaction, error) {
				return geth.SignTx(tx, geth.HomesteadSigner{}, h.key)
			}).Once()
		h.contractMock.Transactor.EXPECT().SendTransaction(
			mock.Anything,
			rootstocktest.MatchTransaction(h.contractMock.Transactor, common.HexToAddress(test.AnyRskAddress), resolvePegInGasLimit, big.NewInt(0), rootstocktest.MustDecodeHex(t, resolvePegInKnownCallData)),
		).Return(assert.AnError).Once()

		result, err := h.contract.ResolvePegIn(resolvePegInTestParams())
		require.ErrorIs(t, err, assert.AnError)
		require.ErrorContains(t, err, "resolve pegin error: ")
		assert.Equal(t, blockchain.TransactionReceipt{}, result)
		h.contractMock.Transactor.AssertExpectations(t)
	})
}

func TestPeginContractImpl_ResolvePegIn_ReceiptErrors(t *testing.T) {
	t.Run("keeps the transaction hash when the receipt can't be read", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		sent := expectSignedResolve(t, h)
		h.client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		result, err := h.contract.ResolvePegIn(resolvePegInTestParams())
		require.ErrorContains(t, err, "resolve pegin error")
		assert.Equal(t, blockchain.TransactionReceipt{TransactionHash: sent.tx.Hash().String()}, result)
		h.contractMock.Transactor.AssertExpectations(t)
	})

	t.Run("keeps the transaction hash when the node returns no receipt", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		sent := expectSignedResolve(t, h)
		h.client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).Return(nil, nil).Once()

		result, err := h.contract.ResolvePegIn(resolvePegInTestParams())
		require.EqualError(t, err, "resolve pegin error: incomplete receipt")
		assert.Equal(t, blockchain.TransactionReceipt{TransactionHash: sent.tx.Hash().String()}, result)
		h.contractMock.Transactor.AssertExpectations(t)
	})

	t.Run("keeps the transaction hash when the receipt is incomplete", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		sent := expectSignedResolve(t, h)
		receipt := requestPegInTestReceipt(1)
		receipt.BlockNumber = nil
		h.client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).Return(receipt, nil).Once()

		result, err := h.contract.ResolvePegIn(resolvePegInTestParams())
		require.EqualError(t, err, "resolve pegin error: receipt without block number")
		assert.Equal(t, blockchain.TransactionReceipt{TransactionHash: sent.tx.Hash().String()}, result)
	})
}

func TestPeginContractImpl_ResolvePegIn_MinedRevert(t *testing.T) {
	t.Run("decodes the revert of a mined transaction at its block", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		sent := expectSignedResolve(t, h)
		receipt := requestPegInTestReceipt(0)
		h.client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).Return(receipt, nil).Once()
		h.client.EXPECT().CallContract(mock.Anything, matchResolvePegInCallMsg(t, h.from), big.NewInt(123)).
			Return(nil, rootstocktest.NewRskRpcError(requestPegInRevertMessage, bridgeValidationFailedRevert)).Once()

		result, err := h.contract.ResolvePegIn(resolvePegInTestParams())
		require.ErrorIs(t, err, blockchain.WaitingForBridgeError)
		require.ErrorIs(t, err, blockchain.TxFailedError)
		require.ErrorContains(t, err, "resolve pegin transaction "+receipt.TxHash.String()+" reverted: ")
		expectedReceipt, parseErr := rootstock.ParseReceipt(sent.tx, receipt)
		require.NoError(t, parseErr)
		assert.Equal(t, expectedReceipt, result)
		h.client.AssertExpectations(t)
	})

	t.Run("reports a mined revert the replay can't explain as a failed transaction", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		expectSignedResolve(t, h)
		h.client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).Return(requestPegInTestReceipt(0), nil).Once()
		h.client.EXPECT().CallContract(mock.Anything, mock.Anything, big.NewInt(123)).Return(nil, nil).Once()

		_, err := h.contract.ResolvePegIn(resolvePegInTestParams())
		require.ErrorIs(t, err, blockchain.TxFailedError)
		require.NotErrorIs(t, err, blockchain.WaitingForBridgeError)
		h.client.AssertExpectations(t)
	})
}

func TestPeginContractImpl_ResolvePegIn_InvalidParams(t *testing.T) {
	t.Run("rejects an invalid rsk address without sending", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		params := resolvePegInTestParams()
		params.RskAddress = test.AnyString

		result, err := h.contract.ResolvePegIn(params)
		require.ErrorIs(t, err, blockchain.InvalidAddressError)
		assert.Equal(t, blockchain.TransactionReceipt{}, result)
	})

	t.Run("rejects a missing block height without sending", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		params := resolvePegInTestParams()
		params.BlockHeight = nil

		result, err := h.contract.ResolvePegIn(params)
		require.EqualError(t, err, "resolve pegin error: missing block height")
		assert.Equal(t, blockchain.TransactionReceipt{}, result)
	})
}

func TestPeginContractImpl_PegInResolvedEvent(t *testing.T) {
	pegInContractAddress := common.HexToAddress(test.AnyRskAddress)
	txHash := "0x7a3c5e8f9b1d2a4c6e8f0a2b4c6d8e0f1a3b5c7d9e1f3a5b7c9d1e3f5a7b9c1d"
	event := pegInResolvedTestEvent()
	receiptWith := func(logs ...blockchain.TransactionLog) blockchain.TransactionReceipt {
		return blockchain.TransactionReceipt{TransactionHash: txHash, Logs: logs}
	}

	t.Run("decodes the PegInResolved event of the pegin contract", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		unrelated := blockchain.TransactionLog{Address: pegInContractAddress.String()}

		result, err := h.contract.PegInResolvedEvent(receiptWith(unrelated, pegInResolvedTransactionLog(t, pegInContractAddress, event)))
		require.NoError(t, err)
		assert.Equal(t, blockchain.PegInResolvedEvent{
			PegInID:       "6c0c1c2ba8a9cd4f0a3f0e15cc1d0a1b2c3d4e5f60718293a4b5c6d7e8f90a1b",
			Claimer:       "0x9D93929A9099be4355fC2389FbF253982F9dF47c",
			Registrant:    "0x4a9D4b6E3C9f2D1B0aC5E7f8B9C0d1e2F3a4b5c6",
			Released:      entities.NewWei(5000000000000000),
			ClaimerPayout: entities.NewWei(4970000000000000),
			RegistrantFee: entities.NewWei(10000000000000),
			UserPayout:    entities.NewWei(20000000000000),
		}, result)
	})

	t.Run("fails when the receipt has no PegInResolved event", func(t *testing.T) {
		h := newRequestPegInHarness(t)

		result, err := h.contract.PegInResolvedEvent(receiptWith())
		require.EqualError(t, err, "PegInResolved event not found in transaction "+txHash)
		assert.Equal(t, blockchain.PegInResolvedEvent{}, result)
	})

	t.Run("ignores a PegInResolved log from another contract", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		otherContract := common.HexToAddress("0x5555555555555555555555555555555555555555")

		result, err := h.contract.PegInResolvedEvent(receiptWith(pegInResolvedTransactionLog(t, otherContract, event)))
		require.EqualError(t, err, "PegInResolved event not found in transaction "+txHash)
		assert.Equal(t, blockchain.PegInResolvedEvent{}, result)
	})

	t.Run("rejects a PegInResolved log with empty data instead of panicking", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		emptyData := pegInResolvedTransactionLog(t, pegInContractAddress, event)
		emptyData.Data = nil

		result, err := h.contract.PegInResolvedEvent(receiptWith(emptyData))
		require.EqualError(t, err, "PegInResolved event in transaction "+txHash+" has 0 data bytes, expected 128")
		assert.Equal(t, blockchain.PegInResolvedEvent{}, result)
	})

	t.Run("returns the unpack error of a PegInResolved log without its indexed topics", func(t *testing.T) {
		h := newRequestPegInHarness(t)
		malformed := pegInResolvedTransactionLog(t, pegInContractAddress, event)
		malformed.Topics = malformed.Topics[:1]

		result, err := h.contract.PegInResolvedEvent(receiptWith(malformed))
		require.ErrorContains(t, err, "PegInResolved event in transaction "+txHash+": ")
		assert.Equal(t, blockchain.PegInResolvedEvent{}, result)
	})
}
