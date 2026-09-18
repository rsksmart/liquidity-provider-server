package pegin_test

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
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
		),
	}
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
	)
}

func matchBlockNumber(number uint64) interface{} {
	return mock.MatchedBy(func(got *big.Int) bool {
		return got != nil && got.Cmp(big.NewInt(int64(number))) == 0
	})
}

func (h *settleHarness) expectCanonicalReceipt() blockchain.TransactionReceipt {
	receipt := successReceipt()
	h.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(receipt, nil).Once()
	h.rsk.On("GetBlockByNumber", mock.Anything, matchBlockNumber(receipt.BlockNumber)).
		Return(blockchain.BlockInfo{Hash: claimRskBlockHash}, nil).Once()
	return receipt
}

func TestSettlePegInClaimUseCase_EmptyHashDoesNotResubmit(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	peginContract := mocks.NewPeginContractMock(t)
	rsk := mocks.NewRootstockRpcServerMock(t)
	useCase := newSettleUseCase(t, claims, peginContract, rsk)

	claim := submittingClaim()
	claim.TxHash = ""
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

func TestSettlePegInClaimUseCase_GetBlockByNumberFailure(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	peginContract := mocks.NewPeginContractMock(t)
	rsk := mocks.NewRootstockRpcServerMock(t)
	rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(successReceipt(), nil).Once()
	rsk.On("GetBlockByNumber", mock.Anything, matchBlockNumber(successReceipt().BlockNumber)).
		Return(blockchain.BlockInfo{}, errors.New("not found")).Once()
	useCase := newSettleUseCase(t, claims, peginContract, rsk)

	err := useCase.Run(context.Background(), submittingClaim())
	require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	peginContract.AssertNotCalled(t, "UnpackPegInRequested", mock.Anything)
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
	rsk.AssertNotCalled(t, "GetBlockByNumber", mock.Anything, mock.Anything)
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
	rsk.AssertNotCalled(t, "GetBlockByNumber", mock.Anything, mock.Anything)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	peginContract.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestSettlePegInClaimUseCase_BlockHashMismatchStaysSubmitting(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	peginContract := mocks.NewPeginContractMock(t)
	rsk := mocks.NewRootstockRpcServerMock(t)
	rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(successReceipt(), nil).Once()
	rsk.On("GetBlockByNumber", mock.Anything, matchBlockNumber(successReceipt().BlockNumber)).
		Return(blockchain.BlockInfo{Hash: "0xotherblock"}, nil).Once()
	useCase := newSettleUseCase(t, claims, peginContract, rsk)

	err := useCase.Run(context.Background(), submittingClaim())
	require.NoError(t, err)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	peginContract.AssertNotCalled(t, "UnpackPegInRequested", mock.Anything)
	peginContract.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestSettlePegInClaimUseCase_TxFailedErrorClassifiesViaPreflight(t *testing.T) {
	repo := newMemoryClaimRepo(submittingClaim())
	harness := newSettleHarness(t, repo)
	receipt := successReceipt()
	harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).
		Return(receipt, blockchain.TxFailedError).Once()
	harness.rsk.On("GetBlockByNumber", mock.Anything, matchBlockNumber(receipt.BlockNumber)).
		Return(blockchain.BlockInfo{Hash: claimRskBlockHash}, nil).Once()
	harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(harness.payingTx(10), nil).Once()
	harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
	harness.btc.On("GetRawTransaction", claimDepositTxID).Return(harness.rawTx, nil).Once()
	harness.btc.On("GetTransactionBlockInfo", claimDepositTxID).Return(harness.block, nil).Once()
	harness.btc.On("BuildMerkleBranch", claimDepositTxID).Return(harness.merkle, nil).Once()
	harness.pegin.On("SimulateRequestPegIn", mock.Anything).Return(blockchain.ErrPegInAlreadyProcessed).Once()

	err := harness.useCase.Run(context.Background(), submittingClaim())
	require.NoError(t, err)

	stored := repo.stored()
	assert.Equal(t, rootstock.PegInClaimRaceLost, stored.State)
	assert.Equal(t, "0", stored.ReservedWei.String())
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	harness.pegin.AssertNotCalled(t, "UnpackPegInRequested", mock.Anything)
}

func TestSettlePegInClaimUseCase_TxFailedIdentifyNilIsRetryable(t *testing.T) {
	repo := newMemoryClaimRepo(submittingClaim())
	harness := newSettleHarness(t, repo)
	receipt := successReceipt()
	harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).
		Return(receipt, blockchain.TxFailedError).Once()
	harness.rsk.On("GetBlockByNumber", mock.Anything, matchBlockNumber(receipt.BlockNumber)).
		Return(blockchain.BlockInfo{Hash: claimRskBlockHash}, nil).Once()
	harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(harness.payingTx(10), nil).Once()
	harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
	harness.expectBuildParams()
	harness.pegin.On("SimulateRequestPegIn", mock.Anything).Return(nil).Once()

	err := harness.useCase.Run(context.Background(), submittingClaim())
	require.ErrorIs(t, err, pegin.ErrStatus0ReceiptStillCallable)
	assert.Equal(t, rootstock.PegInClaimRetryableFailure, repo.stored().State)
	assert.Empty(t, repo.stored().TxHash)
	assert.Equal(t, "0", repo.stored().ReservedWei.String())
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestSettlePegInClaimUseCase_TxFailedIdentifyLookupErrors(t *testing.T) {
	t.Run("btc tx info", func(t *testing.T) {
		claims := mocks.NewPegInClaimRepositoryMock(t)
		harness := newSettleHarness(t, claims)
		receipt := successReceipt()
		harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).
			Return(receipt, blockchain.TxFailedError).Once()
		harness.rsk.On("GetBlockByNumber", mock.Anything, matchBlockNumber(receipt.BlockNumber)).
			Return(blockchain.BlockInfo{Hash: claimRskBlockHash}, nil).Once()
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
		receipt := successReceipt()
		harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).
			Return(receipt, blockchain.TxFailedError).Once()
		harness.rsk.On("GetBlockByNumber", mock.Anything, matchBlockNumber(receipt.BlockNumber)).
			Return(blockchain.BlockInfo{Hash: claimRskBlockHash}, nil).Once()
		harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(harness.payingTx(10), nil).Once()
		harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return((*entities.Wei)(nil), assert.AnError).Once()

		err := harness.useCase.Run(context.Background(), submittingClaim())
		require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
		claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
		harness.pegin.AssertNotCalled(t, "SimulateRequestPegIn", mock.Anything)
	})
}

func TestSettlePegInClaimUseCase_SuccessfulReceiptWithEventSetsPegInID(t *testing.T) {
	repo := newMemoryClaimRepo(submittingClaim())
	harness := newSettleHarness(t, repo)
	pegInID := [32]byte{0xca, 0xfe}
	receipt := harness.expectCanonicalReceipt()
	harness.pegin.On("UnpackPegInRequested", receipt).Return(blockchain.PegInRequestedEvent{
		PegInId:    pegInID,
		RskAddress: test.AnyRskAddress,
	}, nil).Once()

	err := harness.useCase.Run(context.Background(), submittingClaim())
	require.NoError(t, err)

	stored := repo.stored()
	assert.Equal(t, rootstock.PegInClaimClaimed, stored.State)
	assert.Equal(t, hex.EncodeToString(pegInID[:]), stored.PegInID)
	assert.Equal(t, "0", stored.ReservedWei.String())
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestSettlePegInClaimUseCase_SuccessfulReceiptWithoutEventClearsReserve(t *testing.T) {
	repo := newMemoryClaimRepo(submittingClaim())
	harness := newSettleHarness(t, repo)
	receipt := harness.expectCanonicalReceipt()
	harness.pegin.On("UnpackPegInRequested", receipt).Return(blockchain.PegInRequestedEvent{}, errors.New("missing")).Once()

	err := harness.useCase.Run(context.Background(), submittingClaim())
	require.NoError(t, err)

	stored := repo.stored()
	assert.Equal(t, rootstock.PegInClaimClaimed, stored.State)
	assert.Empty(t, stored.PegInID)
	assert.Equal(t, "0", stored.ReservedWei.String())
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
	rsk.AssertNotCalled(t, "GetBlockByNumber", mock.Anything, mock.Anything)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	peginContract.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestSettlePegInClaimUseCase_SimulateRpcErrorStaysSubmitting(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newSettleHarness(t, claims)
	receipt := successReceipt()
	harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).
		Return(receipt, blockchain.TxFailedError).Once()
	harness.rsk.On("GetBlockByNumber", mock.Anything, matchBlockNumber(receipt.BlockNumber)).
		Return(blockchain.BlockInfo{Hash: claimRskBlockHash}, nil).Once()
	harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(harness.payingTx(10), nil).Once()
	harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
	harness.expectBuildParams()
	harness.pegin.On("SimulateRequestPegIn", mock.Anything).Return(assert.AnError).Once()

	err := harness.useCase.Run(context.Background(), submittingClaim())
	require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}
