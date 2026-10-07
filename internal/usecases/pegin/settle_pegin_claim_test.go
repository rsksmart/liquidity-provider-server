package pegin_test

import (
	"context"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type settleHarness struct {
	*claimHarness
	useCase *pegin.SettlePegInClaimUseCase
}

func newSettleHarness(t *testing.T, repo rootstock.PegInClaimRepository) *settleHarness {
	t.Helper()
	base := newClaimHarness(t, repo)
	return &settleHarness{
		claimHarness: base,
		useCase: pegin.NewSettlePegInClaimUseCase(
			base.repo,
			blockchain.RskContracts{
				PegIn:                 base.pegin,
				FlyoverConfigurations: base.configs,
				PauseRegistry:         base.pause,
			},
			blockchain.Rpc{Btc: base.btc, Rsk: base.rsk},
			base.eventBus,
			settleMaxReorgDepth,
		),
	}
}

// successReceipt is at block 100, so height 102 is final with depth 2.
const (
	settleMaxReorgDepth = 2
	settleFinalHeight   = 102
)

func (h *settleHarness) expectFinalHeight() {
	h.rsk.On("GetHeight", mock.Anything).Return(uint64(settleFinalHeight), nil).Once()
}

func failedReceipt() blockchain.TransactionReceipt {
	receipt := successReceipt()
	receipt.Status = 0
	return receipt
}

func newSettleUseCase(
	t *testing.T,
	claims rootstock.PegInClaimRepository,
	peginContract *mocks.PeginContractMock,
	rsk *mocks.RootstockRpcServerMock,
) *pegin.SettlePegInClaimUseCase {
	t.Helper()
	return pegin.NewSettlePegInClaimUseCase(
		claims,
		blockchain.RskContracts{
			PegIn:                 peginContract,
			FlyoverConfigurations: mocks.NewFlyoverConfigurationsContractMock(t),
		},
		blockchain.Rpc{Btc: mocks.NewBtcRpcMock(t), Rsk: rsk},
		claimEventBus(),
		settleMaxReorgDepth,
	)
}

func (h *settleHarness) expectReceipt() blockchain.TransactionReceipt {
	receipt := successReceipt()
	h.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(receipt, nil).Once()
	h.expectFinalHeight()
	return receipt
}

func TestSettlePegInClaimUseCase_EmptyHashDoesNotResubmit(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	peginContract := mocks.NewPeginContractMock(t)
	rsk := mocks.NewRootstockRpcServerMock(t)
	useCase := newSettleUseCase(t, claims, peginContract, rsk)

	claim := submittingClaim()
	claim.RequestTxHash = ""
	err := useCase.Run(context.Background(), claim)
	require.NoError(t, err)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	rsk.AssertNotCalled(t, "GetTransactionReceipt", mock.Anything, mock.Anything)
	peginContract.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestSettlePegInClaimUseCase_MissingReceiptDoesNotResubmit(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	peginContract := mocks.NewPeginContractMock(t)
	rsk := mocks.NewRootstockRpcServerMock(t)
	rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).
		Return(blockchain.TransactionReceipt{}, blockchain.ErrTransactionReceiptNotFound).Once()
	useCase := newSettleUseCase(t, claims, peginContract, rsk)

	err := useCase.Run(context.Background(), submittingClaim())
	require.NoError(t, err)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	peginContract.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	rsk.AssertExpectations(t)
}

func TestSettlePegInClaimUseCase_GenericRskRpcFailure(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	peginContract := mocks.NewPeginContractMock(t)
	rsk := mocks.NewRootstockRpcServerMock(t)
	rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).
		Return(blockchain.TransactionReceipt{}, assert.AnError).Once()
	useCase := newSettleUseCase(t, claims, peginContract, rsk)

	err := useCase.Run(context.Background(), submittingClaim())
	require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
	assert.Contains(t, err.Error(), string(usecases.SettlePegInClaimId))
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	peginContract.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestSettlePegInClaimUseCase_RemovedLogsStaySubmitting(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	peginContract := mocks.NewPeginContractMock(t)
	rsk := mocks.NewRootstockRpcServerMock(t)
	receipt := successReceipt()
	receipt.Logs = []blockchain.TransactionLog{{Removed: true}}
	rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(receipt, nil).Once()
	useCase := newSettleUseCase(t, claims, peginContract, rsk)

	err := useCase.Run(context.Background(), submittingClaim())
	require.NoError(t, err)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	peginContract.AssertNotCalled(t, "UnpackPegInRequested", mock.Anything)
	peginContract.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestSettlePegInClaimUseCase_EmptyBlockHashStaysSubmitting(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	peginContract := mocks.NewPeginContractMock(t)
	rsk := mocks.NewRootstockRpcServerMock(t)
	receipt := successReceipt()
	receipt.BlockHash = ""
	rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(receipt, nil).Once()
	useCase := newSettleUseCase(t, claims, peginContract, rsk)

	err := useCase.Run(context.Background(), submittingClaim())
	require.NoError(t, err)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	peginContract.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestSettlePegInClaimUseCase_StatusZeroClassifiesViaPreflight(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newSettleHarness(t, claims)
	lastSaved := recordClaimUpdates(claims)
	receipt := failedReceipt()
	harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).
		Return(receipt, nil).Once()
	harness.expectFinalHeight()
	harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(harness.payingTx(10), nil).Once()
	harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
	harness.btc.On("GetRawTransaction", claimDepositTxID).Return(harness.rawTx, nil).Once()
	harness.btc.On("GetTransactionBlockInfo", claimDepositTxID).Return(harness.block, nil).Once()
	harness.btc.On("BuildMerkleBranch", claimDepositTxID).Return(harness.merkle, nil).Once()
	harness.pegin.On("SimulateRequestPegIn", mock.Anything).Return(blockchain.ErrPegInAlreadyProcessed).Once()

	err := harness.useCase.Run(context.Background(), submittingClaim())
	require.NoError(t, err)

	stored := lastSaved()
	assert.Equal(t, rootstock.PegInClaimRaceLost, stored.State)
	harness.eventBus.AssertCalled(t, "Publish", matchClaimCompleted(rootstock.PegInClaimRaceLost))
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	harness.pegin.AssertNotCalled(t, "UnpackPegInRequested", mock.Anything)
}

func TestSettlePegInClaimUseCase_TxFailedIdentifyNilIsRetryable(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newSettleHarness(t, claims)
	lastSaved := recordClaimUpdates(claims)
	receipt := failedReceipt()
	harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).
		Return(receipt, nil).Once()
	harness.expectFinalHeight()
	harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(harness.payingTx(10), nil).Once()
	harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
	harness.expectBuildParams()
	harness.pegin.On("SimulateRequestPegIn", mock.Anything).Return(nil).Once()

	err := harness.useCase.Run(context.Background(), submittingClaim())
	require.ErrorIs(t, err, blockchain.TxFailedError)
	require.ErrorIs(t, err, usecases.NonRecoverableError)
	assert.Equal(t, rootstock.PegInClaimRetryableFailure, lastSaved().State)
	harness.eventBus.AssertNotCalled(t, "Publish", mock.Anything)
	assert.Empty(t, lastSaved().RequestTxHash)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestSettlePegInClaimUseCase_TxFailedIdentifyLookupErrors(t *testing.T) {
	t.Run("btc tx info", func(t *testing.T) {
		claims := mocks.NewPegInClaimRepositoryMock(t)
		harness := newSettleHarness(t, claims)
		receipt := failedReceipt()
		harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).
			Return(receipt, nil).Once()
		harness.expectFinalHeight()
		harness.btc.On("GetTransactionInfo", claimDepositTxID).
			Return(blockchain.BitcoinTransactionInformation{}, assert.AnError).Once()

		err := harness.useCase.Run(context.Background(), submittingClaim())
		require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
		claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
		harness.pegin.AssertNotCalled(t, "SimulateRequestPegIn", mock.Anything)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
	t.Run("fee oracle", func(t *testing.T) {
		claims := mocks.NewPegInClaimRepositoryMock(t)
		harness := newSettleHarness(t, claims)
		receipt := failedReceipt()
		harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).
			Return(receipt, nil).Once()
		harness.expectFinalHeight()
		harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(harness.payingTx(10), nil).Once()
		harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return((*entities.Wei)(nil), assert.AnError).Once()

		err := harness.useCase.Run(context.Background(), submittingClaim())
		require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
		claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
		harness.pegin.AssertNotCalled(t, "SimulateRequestPegIn", mock.Anything)
	})
}

func TestSettlePegInClaimUseCase_TxFailedIdentifyRequestParamsErrors(t *testing.T) {
	paramCases := []struct {
		name string
		stub func(*settleHarness)
	}{
		{
			name: "raw tx",
			stub: func(h *settleHarness) {
				h.btc.On("GetRawTransaction", claimDepositTxID).Return([]byte(nil), assert.AnError).Once()
			},
		},
		{
			name: "block info",
			stub: func(h *settleHarness) {
				h.btc.On("GetRawTransaction", claimDepositTxID).Return(h.rawTx, nil).Once()
				h.btc.On("GetTransactionBlockInfo", claimDepositTxID).
					Return(blockchain.BitcoinBlockInformation{}, assert.AnError).Once()
			},
		},
		{
			name: "merkle branch",
			stub: func(h *settleHarness) {
				h.btc.On("GetRawTransaction", claimDepositTxID).Return(h.rawTx, nil).Once()
				h.btc.On("GetTransactionBlockInfo", claimDepositTxID).Return(h.block, nil).Once()
				h.btc.On("BuildMerkleBranch", claimDepositTxID).Return(blockchain.MerkleBranch{}, assert.AnError).Once()
			},
		},
	}
	for _, tc := range paramCases {
		t.Run(tc.name, func(t *testing.T) {
			claims := mocks.NewPegInClaimRepositoryMock(t)
			harness := newSettleHarness(t, claims)
			harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(failedReceipt(), nil).Once()
			harness.expectFinalHeight()
			harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(harness.payingTx(10), nil).Once()
			harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
			tc.stub(harness)

			err := harness.useCase.Run(context.Background(), submittingClaim())
			require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
			require.ErrorIs(t, err, assert.AnError)
			claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			harness.pegin.AssertNotCalled(t, "SimulateRequestPegIn", mock.Anything)
		})
	}
}

func TestSettlePegInClaimUseCase_SuccessfulReceiptWithEventSetsPegInID(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newSettleHarness(t, claims)
	lastSaved := recordClaimUpdates(claims)
	pegInID := [32]byte{0xca, 0xfe}
	receipt := harness.expectReceipt()
	harness.pegin.On("UnpackPegInRequested", receipt).Return(blockchain.PegInRequestedEvent{
		PegInId:    pegInID,
		RskAddress: test.AnyRskAddress,
	}, nil).Once()

	err := harness.useCase.Run(context.Background(), submittingClaim())
	require.NoError(t, err)

	stored := lastSaved()
	assert.Equal(t, rootstock.PegInClaimClaimed, stored.State)
	harness.eventBus.AssertCalled(t, "Publish", matchClaimCompleted(rootstock.PegInClaimClaimed))
	assert.Equal(t, hex.EncodeToString(pegInID[:]), stored.PegInID)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestSettlePegInClaimUseCase_SuccessfulReceiptWithoutEventIsClaimed(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newSettleHarness(t, claims)
	lastSaved := recordClaimUpdates(claims)
	receipt := harness.expectReceipt()
	harness.pegin.On("UnpackPegInRequested", receipt).Return(blockchain.PegInRequestedEvent{}, errors.New("missing")).Once()

	err := harness.useCase.Run(context.Background(), submittingClaim())
	require.NoError(t, err)

	stored := lastSaved()
	assert.Equal(t, rootstock.PegInClaimClaimed, stored.State)
	assert.Empty(t, stored.PegInID)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestSettlePegInClaimUseCase_ZeroBlockNumberStaysSubmitting(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	peginContract := mocks.NewPeginContractMock(t)
	rsk := mocks.NewRootstockRpcServerMock(t)
	receipt := successReceipt()
	receipt.BlockNumber = 0
	rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(receipt, nil).Once()
	useCase := newSettleUseCase(t, claims, peginContract, rsk)

	err := useCase.Run(context.Background(), submittingClaim())
	require.NoError(t, err)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	peginContract.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestSettlePegInClaimUseCase_TxFailedIdentifyTypedContractErrorIsRetryable(t *testing.T) {
	cases := []error{
		blockchain.ErrAddressNotRegistered,
		blockchain.ErrDepositOutputNotFound,
		blockchain.ErrInsufficientConfirmations,
		blockchain.ErrIncorrectFronting,
		blockchain.ErrPegInBelowMinimum,
	}
	for _, simulateErr := range cases {
		t.Run(simulateErr.Error(), func(t *testing.T) {
			claims := mocks.NewPegInClaimRepositoryMock(t)
			harness := newSettleHarness(t, claims)
			lastSaved := recordClaimUpdates(claims)
			receipt := failedReceipt()
			harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).
				Return(receipt, nil).Once()
			harness.expectFinalHeight()
			harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(harness.payingTx(10), nil).Once()
			harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
			harness.expectBuildParams()
			harness.pegin.On("SimulateRequestPegIn", mock.Anything).Return(simulateErr).Once()

			err := harness.useCase.Run(context.Background(), submittingClaim())
			require.ErrorIs(t, err, simulateErr)
			require.ErrorIs(t, err, usecases.NonRecoverableError)
			require.NotErrorIs(t, err, usecases.InfrastructureUnavailableError)
			stored := lastSaved()
			assert.Equal(t, rootstock.PegInClaimRetryableFailure, stored.State)
			assert.Empty(t, stored.RequestTxHash)
			harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
		})
	}
}

func TestSettlePegInClaimUseCase_SimulateRpcErrorStaysSubmitting(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newSettleHarness(t, claims)
	receipt := failedReceipt()
	harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).
		Return(receipt, nil).Once()
	harness.expectFinalHeight()
	harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(harness.payingTx(10), nil).Once()
	harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
	harness.expectBuildParams()
	harness.pegin.On("SimulateRequestPegIn", mock.Anything).Return(assert.AnError).Once()

	err := harness.useCase.Run(context.Background(), submittingClaim())
	require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestSettlePegInClaimUseCase_ReceiptNotFinalStaysSubmitting(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	peginContract := mocks.NewPeginContractMock(t)
	rsk := mocks.NewRootstockRpcServerMock(t)
	rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(successReceipt(), nil).Once()
	rsk.On("GetHeight", mock.Anything).Return(uint64(settleFinalHeight-1), nil).Once()
	useCase := newSettleUseCase(t, claims, peginContract, rsk)

	err := useCase.Run(context.Background(), submittingClaim())
	require.NoError(t, err)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	peginContract.AssertNotCalled(t, "UnpackPegInRequested", mock.Anything)
}

func TestSettlePegInClaimUseCase_GetHeightErrorIsUnavailable(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	peginContract := mocks.NewPeginContractMock(t)
	rsk := mocks.NewRootstockRpcServerMock(t)
	rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(successReceipt(), nil).Once()
	rsk.On("GetHeight", mock.Anything).Return(uint64(0), assert.AnError).Once()
	useCase := newSettleUseCase(t, claims, peginContract, rsk)

	err := useCase.Run(context.Background(), submittingClaim())
	require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
	require.ErrorIs(t, err, assert.AnError)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestSettlePegInClaimUseCase_UnpackFailureKeepsStoredPegInID(t *testing.T) {
	stored := submittingClaim()
	stored.PegInID = "aabbcc"
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newSettleHarness(t, claims)
	lastSaved := recordClaimUpdates(claims)
	receipt := harness.expectReceipt()
	harness.pegin.On("UnpackPegInRequested", receipt).Return(blockchain.PegInRequestedEvent{}, errors.New("missing")).Once()

	err := harness.useCase.Run(context.Background(), stored)
	require.NoError(t, err)
	assert.Equal(t, rootstock.PegInClaimClaimed, lastSaved().State)
	assert.Equal(t, "aabbcc", lastSaved().PegInID)
}

func TestSettlePegInClaimUseCase_WitnessSerializedRawTxIsNonRecoverable(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newSettleHarness(t, claims)
	harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(failedReceipt(), nil).Once()
	harness.expectFinalHeight()
	harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(harness.payingTx(10), nil).Once()
	harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
	witnessRawTx := []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x01, 0xaa, 0xbb}
	harness.btc.On("GetRawTransaction", claimDepositTxID).Return(witnessRawTx, nil).Once()

	err := harness.useCase.Run(context.Background(), submittingClaim())
	require.ErrorIs(t, err, blockchain.ErrWitnessSerializedTxNotAccepted)
	require.ErrorIs(t, err, usecases.NonRecoverableError)
	require.NotErrorIs(t, err, usecases.InfrastructureUnavailableError)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	harness.pegin.AssertNotCalled(t, "SimulateRequestPegIn", mock.Anything)
}
