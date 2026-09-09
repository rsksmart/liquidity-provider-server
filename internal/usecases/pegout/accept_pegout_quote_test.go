package pegout_test

import (
	"context"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var acceptPegoutQuoteHash = "c8d4ad8d5d717371b92950cbe43a6a4e891cf27bcd7603c988595866944bd9cf"
var acceptPegoutQuoteHashSignature = "b062b09f5f3000f1092e606e90fa449e8527fb1bac20ff72897fd1d0a8aa3b18049d39d1956110992de0284e6d85223d4f69ed06e57184cad13abca7b421d6e41b"

func acceptTestQuote(now time.Time) quote.PegoutQuote {
	return quote.PegoutQuote{
		LbcAddress:            "0xabcd01",
		LpRskAddress:          "0xabcd02",
		BtcRefundAddress:      "hijk",
		RskRefundAddress:      "0xabcd04",
		LpBtcAddress:          "edfg",
		CallFee:               entities.NewWei(5),
		PenaltyFee:            entities.NewWei(1),
		Nonce:                 1,
		DepositAddress:        "address",
		Value:                 entities.NewWei(12),
		AgreementTimestamp:    uint32(now.Unix()),
		DepositDateLimit:      uint32(now.Unix() + 600),
		DepositConfirmations:  1,
		TransferConfirmations: 1,
		TransferTime:          600,
		ExpireDate:            uint32(now.Unix() + 600),
		ExpireBlock:           1,
		GasFee:                entities.NewWei(6),
		ChainId:               31,
	}
}

func TestAcceptQuoteUseCase_Run_Paused(t *testing.T) {
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	lp := new(mocks.ProviderMock)
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: true, Since: 5, Reason: "test"}, nil)
	pegoutContract.EXPECT().GetAddress().Return("test-contract")

	contracts := blockchain.RskContracts{PegOut: pegoutContract}
	useCase := pegout.NewAcceptQuoteUseCase(quoteRepository, contracts, lp)
	result, err := useCase.Run(context.Background(), acceptPegoutQuoteHash, "")
	assert.Empty(t, result)
	require.ErrorIs(t, err, blockchain.ContractPausedError)
}

func TestAcceptQuoteUseCase_Run(t *testing.T) {
	quoteHash := acceptPegoutQuoteHash
	quoteMock := acceptTestQuote(time.Now())
	quoteRepositoryMock := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepositoryMock.On("GetQuote", test.AnyCtx, quoteHash).Return(&quoteMock, nil).Once()
	quoteRepositoryMock.On("GetRetainedQuote", test.AnyCtx, quoteHash).Return(nil, nil).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.On("GetAddress").Return("0xabcd01").Once()
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	lp := new(mocks.ProviderMock)
	lp.On("SignPegoutQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(acceptPegoutQuoteHashSignature, nil)

	useCase := pegout.NewAcceptQuoteUseCase(quoteRepositoryMock, blockchain.RskContracts{PegOut: pegoutContract}, lp)
	result, err := useCase.Run(context.Background(), quoteHash, "")

	quoteRepositoryMock.AssertExpectations(t)
	pegoutContract.AssertExpectations(t)
	lp.AssertExpectations(t)
	quoteRepositoryMock.AssertNotCalled(t, "InsertRetainedQuote")
	lp.AssertNotCalled(t, "HasPegoutLiquidity")
	require.NoError(t, err)
	assert.Equal(t, "0xabcd01", result.DepositAddress)
	assert.Equal(t, acceptPegoutQuoteHashSignature, result.Signature)
}

func TestAcceptQuoteUseCase_Run_AlreadyAccepted(t *testing.T) {
	quoteMock := acceptTestQuote(time.Now())
	retainedQuote := quote.RetainedPegoutQuote{
		QuoteHash:         acceptPegoutQuoteHash,
		DepositAddress:    "0xexisting",
		Signature:         "existing-sig",
		RequiredLiquidity: entities.NewWei(18),
		State:             quote.PegoutStateClaimed,
	}
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(&quoteMock, nil).Once()
	quoteRepository.On("GetRetainedQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(&retainedQuote, nil).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	lp := new(mocks.ProviderMock)

	useCase := pegout.NewAcceptQuoteUseCase(quoteRepository, blockchain.RskContracts{PegOut: pegoutContract}, lp)
	result, err := useCase.Run(context.Background(), acceptPegoutQuoteHash, "")

	require.NoError(t, err)
	assert.Equal(t, retainedQuote.Signature, result.Signature)
	assert.Equal(t, retainedQuote.DepositAddress, result.DepositAddress)
	lp.AssertNotCalled(t, "SignPegoutQuote")
	quoteRepository.AssertNotCalled(t, "InsertRetainedQuote")
}

func TestAcceptQuoteUseCase_Run_QuoteNotFound(t *testing.T) {
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(nil, nil).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	lp := new(mocks.ProviderMock)

	useCase := pegout.NewAcceptQuoteUseCase(quoteRepository, blockchain.RskContracts{PegOut: pegoutContract}, lp)
	result, err := useCase.Run(context.Background(), acceptPegoutQuoteHash, "")

	assert.Empty(t, result)
	require.ErrorIs(t, err, usecases.QuoteNotFoundError)
}

func TestAcceptQuoteUseCase_Run_ExpiredQuote(t *testing.T) {
	now := time.Now()
	expired := acceptTestQuote(now)
	expired.ExpireDate = uint32(now.Unix() - 600)
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(&expired, nil).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	lp := new(mocks.ProviderMock)

	useCase := pegout.NewAcceptQuoteUseCase(quoteRepository, blockchain.RskContracts{PegOut: pegoutContract}, lp)
	result, err := useCase.Run(context.Background(), acceptPegoutQuoteHash, "")

	assert.Empty(t, result)
	require.ErrorIs(t, err, usecases.ExpiredQuoteError)
}

func TestAcceptQuoteUseCase_Run_SignError(t *testing.T) {
	quoteMock := acceptTestQuote(time.Now())
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(&quoteMock, nil).Once()
	quoteRepository.On("GetRetainedQuote", test.AnyCtx, acceptPegoutQuoteHash).Return(nil, nil).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	lp := new(mocks.ProviderMock)
	lp.On("SignPegoutQuote", test.AnyCtx, acceptPegoutQuoteHash).Return("", assert.AnError).Once()

	useCase := pegout.NewAcceptQuoteUseCase(quoteRepository, blockchain.RskContracts{PegOut: pegoutContract}, lp)
	result, err := useCase.Run(context.Background(), acceptPegoutQuoteHash, "")

	assert.Empty(t, result)
	require.Error(t, err)
	quoteRepository.AssertNotCalled(t, "InsertRetainedQuote")
}

func TestAcceptQuoteUseCase_S15_5_AcceptSpamLeavesLiquidityUnchanged(t *testing.T) {
	quoteHashes := []string{
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
		"dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
	}
	quoteMock := acceptTestQuote(time.Now())
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	for _, hash := range quoteHashes {
		q := quoteMock
		quoteRepository.On("GetQuote", test.AnyCtx, hash).Return(&q, nil).Once()
		quoteRepository.On("GetRetainedQuote", test.AnyCtx, hash).Return(nil, nil).Once()
	}
	quoteRepository.On("GetRetainedQuoteByState", test.AnyCtx, quote.PegoutStateWaitingForDepositConfirmations).
		Return([]quote.RetainedPegoutQuote{}, nil).Twice()

	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	pegoutContract.On("GetAddress").Return("0xabcd01")

	signLp := new(mocks.ProviderMock)
	for _, hash := range quoteHashes {
		signLp.On("SignPegoutQuote", test.AnyCtx, hash).Return(acceptPegoutQuoteHashSignature, nil).Once()
	}

	btcWallet := new(mocks.BitcoinWalletMock)
	btcWallet.On("GetBalance").Return(entities.NewWei(10_000), nil).Twice()
	liquidityLp := dataproviders.NewLocalLiquidityProvider(
		nil, quoteRepository, nil, blockchain.Rpc{}, nil, btcWallet, blockchain.RskContracts{},
	)

	before, err := liquidityLp.AvailablePegoutLiquidity(context.Background())
	require.NoError(t, err)

	useCase := pegout.NewAcceptQuoteUseCase(quoteRepository, blockchain.RskContracts{PegOut: pegoutContract}, signLp)
	for _, hash := range quoteHashes {
		result, runErr := useCase.Run(context.Background(), hash, "")
		require.NoError(t, runErr)
		assert.Equal(t, acceptPegoutQuoteHashSignature, result.Signature)
	}

	after, err := liquidityLp.AvailablePegoutLiquidity(context.Background())
	require.NoError(t, err)
	assert.Equal(t, before, after)
	assert.Equal(t, entities.NewWei(10_000), after)

	quoteRepository.AssertNotCalled(t, "InsertRetainedQuote")
	signLp.AssertNotCalled(t, "HasPegoutLiquidity")
	btcWallet.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
}
