package pegout_test

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-playground/validator/v10"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var trustedAccountRepository = new(mocks.TrustedAccountRepositoryMock)
var signingHashFunction = crypto.Keccak256

func TestAcceptQuoteUseCase_Run(t *testing.T) {
	quoteHash := "0x654321"
	now := time.Now()
	signature := "0x010203"
	quoteMock := quote.PegoutQuote{
		LbcAddress:            "0xabcd01",
		LpRskAddress:          "0xabcd02",
		BtcRefundAddress:      "hijk",
		RskRefundAddress:      "0xabcd04",
		LpBtcAddress:          "edfg",
		CallFee:               entities.NewWei(5),
		PenaltyFee:            1,
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
		ProductFeeAmount:      2,
	}
	retainedQuote := quote.RetainedPegoutQuote{
		QuoteHash: quoteHash, DepositAddress: quoteMock.LbcAddress,
		Signature: signature, RequiredLiquidity: entities.NewWei(18),
		State: quote.PegoutStateWaitingForDeposit,
	}
	creationData := quote.PegoutCreationData{FeeRate: utils.NewBigFloat64(1.5), FeePercentage: utils.NewBigFloat64(12.5), GasPrice: entities.NewWei(1), FixedFee: entities.NewWei(100)}
	quoteRepositoryMock := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepositoryMock.On("InsertRetainedQuote", test.AnyCtx, retainedQuote).Return(nil).Once()
	quoteRepositoryMock.On("GetQuote", test.AnyCtx, quoteHash).Return(&quoteMock, nil).Once()
	quoteRepositoryMock.On("GetRetainedQuote", test.AnyCtx, quoteHash).Return(nil, nil).Once()
	quoteRepositoryMock.EXPECT().GetPegoutCreationData(test.AnyCtx, quoteHash).Return(creationData).Once()
	lbc := new(mocks.LbcMock)
	lbc.On("GetAddress").Return("0xabcd01").Once()
	lp := new(mocks.ProviderMock)
	lp.On("HasPegoutLiquidity", test.AnyCtx, mock.Anything).Return(nil).Once()
	lp.On("SignQuote", mock.Anything).Return(signature, nil).Once()
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.AcceptedPegoutQuoteEvent) bool {
		return assert.Equal(t, quoteMock, event.Quote) && assert.Equal(t, retainedQuote, event.RetainedQuote) &&
			assert.Equal(t, quote.AcceptedPegoutQuoteEventId, event.Event.Id()) && assert.Equal(t, creationData, event.CreationData)
	})).Once()
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Once()
	mutex.On("Unlock").Once()
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := pegout.NewAcceptQuoteUseCase(quoteRepositoryMock, contracts, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
	result, err := useCase.Run(context.Background(), quoteHash, "")
	quoteRepositoryMock.AssertExpectations(t)
	lbc.AssertExpectations(t)
	lp.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	mutex.AssertExpectations(t)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, quoteMock.LbcAddress, result.DepositAddress)
	assert.Equal(t, signature, result.Signature)
}

func TestAcceptQuoteUseCase_Run_AlreadyAcceptedQuote(t *testing.T) {
	quoteHash := "0x654321"
	now := time.Now()
	quoteMock := quote.PegoutQuote{
		LbcAddress:            "0xabcd01",
		LpRskAddress:          "0xabcd02",
		BtcRefundAddress:      "hijk",
		RskRefundAddress:      "0xabcd04",
		LpBtcAddress:          "edfg",
		CallFee:               entities.NewWei(1),
		PenaltyFee:            1,
		Nonce:                 1,
		DepositAddress:        "address",
		Value:                 entities.NewWei(1),
		AgreementTimestamp:    uint32(now.Unix()),
		DepositDateLimit:      uint32(now.Unix() + 600),
		DepositConfirmations:  1,
		TransferConfirmations: 1,
		TransferTime:          600,
		ExpireDate:            uint32(now.Unix() + 600),
		ExpireBlock:           1,
		GasFee:                entities.NewWei(1),
		ProductFeeAmount:      1,
	}
	retainedQuote := quote.RetainedPegoutQuote{
		QuoteHash:         quoteHash,
		DepositAddress:    quoteMock.LbcAddress,
		Signature:         "signature",
		RequiredLiquidity: entities.NewWei(1),
		State:             quote.PegoutStateWaitingForDeposit,
	}
	quoteRepositoryMock := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepositoryMock.On("GetQuote", test.AnyCtx, quoteHash).Return(&quoteMock, nil).Once()
	quoteRepositoryMock.On("GetRetainedQuote", test.AnyCtx, quoteHash).Return(&retainedQuote, nil).Once()
	lbc := new(mocks.LbcMock)
	lp := new(mocks.ProviderMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return()
	mutex.On("Unlock").Return()
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := pegout.NewAcceptQuoteUseCase(quoteRepositoryMock, contracts, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
	result, err := useCase.Run(context.Background(), quoteHash, "")
	quoteRepositoryMock.AssertExpectations(t)
	lbc.AssertNotCalled(t, "GetAddress")
	lp.AssertNotCalled(t, "SignQuote")
	lp.AssertNotCalled(t, "HasPegoutLiquidity")
	eventBus.AssertNotCalled(t, "Publish")
	mutex.AssertExpectations(t)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, quoteMock.LbcAddress, result.DepositAddress)
	assert.Equal(t, "signature", result.Signature)
}

func TestAcceptQuoteUseCase_Run_ExpiredQuote(t *testing.T) {
	quoteHash := "0x654321"
	now := time.Now()
	quoteMock := quote.PegoutQuote{
		LbcAddress:            "0xabcd01",
		LpRskAddress:          "0xabcd02",
		BtcRefundAddress:      "hijk",
		RskRefundAddress:      "0xabcd04",
		LpBtcAddress:          "edfg",
		CallFee:               entities.NewWei(1),
		PenaltyFee:            1,
		Nonce:                 1,
		DepositAddress:        "address",
		Value:                 entities.NewWei(1),
		AgreementTimestamp:    uint32(now.Unix() - 600),
		DepositDateLimit:      uint32(now.Unix()),
		DepositConfirmations:  1,
		TransferConfirmations: 1,
		TransferTime:          600,
		ExpireDate:            uint32(now.Unix()),
		ExpireBlock:           1,
		GasFee:                entities.NewWei(1),
		ProductFeeAmount:      1,
	}
	quoteRepositoryMock := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepositoryMock.On("GetQuote", test.AnyCtx, quoteHash).Return(&quoteMock, nil).Once()
	lbc := new(mocks.LbcMock)
	lp := new(mocks.ProviderMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := pegout.NewAcceptQuoteUseCase(quoteRepositoryMock, contracts, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
	result, err := useCase.Run(context.Background(), quoteHash, "")
	quoteRepositoryMock.AssertExpectations(t)
	quoteRepositoryMock.AssertNotCalled(t, "GetRetainedQuote")
	lbc.AssertNotCalled(t, "GetAddress")
	lp.AssertNotCalled(t, "SignQuote")
	lp.AssertNotCalled(t, "HasPegoutLiquidity")
	eventBus.AssertNotCalled(t, "Publish")
	mutex.AssertNotCalled(t, "Lock")
	mutex.AssertNotCalled(t, "Unlock")
	assert.Empty(t, result)
	require.ErrorIs(t, err, usecases.ExpiredQuoteError)
}

func TestAcceptQuoteUseCase_Run_QuoteNotFound(t *testing.T) {
	quoteHash := "0x654321"
	quoteRepositoryMock := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepositoryMock.On("GetQuote", test.AnyCtx, quoteHash).Return(nil, nil).Once()
	lbc := new(mocks.LbcMock)
	lp := new(mocks.ProviderMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := pegout.NewAcceptQuoteUseCase(quoteRepositoryMock, contracts, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
	result, err := useCase.Run(context.Background(), quoteHash, "")
	quoteRepositoryMock.AssertExpectations(t)
	quoteRepositoryMock.AssertNotCalled(t, "GetRetainedQuote")
	lbc.AssertNotCalled(t, "GetAddress")
	lp.AssertNotCalled(t, "SignQuote")
	lp.AssertNotCalled(t, "HasPegoutLiquidity")
	eventBus.AssertNotCalled(t, "Publish")
	mutex.AssertNotCalled(t, "Lock")
	mutex.AssertNotCalled(t, "Unlock")
	assert.Empty(t, result)
	require.ErrorIs(t, err, usecases.QuoteNotFoundError)
}

func TestAcceptQuoteUseCase_Run_NoLiquidity(t *testing.T) {
	quoteHash := "0x654321"
	now := time.Now()
	quoteMock := quote.PegoutQuote{
		LbcAddress:            "0xabcd01",
		LpRskAddress:          "0xabcd02",
		BtcRefundAddress:      "hijk",
		RskRefundAddress:      "0xabcd04",
		LpBtcAddress:          "edfg",
		CallFee:               entities.NewWei(10),
		PenaltyFee:            1,
		Nonce:                 1,
		DepositAddress:        "address",
		Value:                 entities.NewWei(50),
		AgreementTimestamp:    uint32(now.Unix()),
		DepositDateLimit:      uint32(now.Unix() + 600),
		DepositConfirmations:  1,
		TransferConfirmations: 1,
		TransferTime:          600,
		ExpireDate:            uint32(now.Unix() + 600),
		ExpireBlock:           1,
		GasFee:                entities.NewWei(15),
		ProductFeeAmount:      8,
	}
	quoteRepositoryMock := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepositoryMock.On("GetQuote", test.AnyCtx, quoteHash).Return(&quoteMock, nil).Once()
	quoteRepositoryMock.On("GetRetainedQuote", test.AnyCtx, quoteHash).Return(nil, nil).Once()
	lbc := new(mocks.LbcMock)
	lp := new(mocks.ProviderMock)
	lp.On("HasPegoutLiquidity", test.AnyCtx, entities.NewWei(65)).Return(usecases.NoLiquidityError).Once()
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Once()
	mutex.On("Unlock").Once()
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := pegout.NewAcceptQuoteUseCase(quoteRepositoryMock, contracts, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
	result, err := useCase.Run(context.Background(), quoteHash, "")
	quoteRepositoryMock.AssertExpectations(t)
	lp.AssertExpectations(t)
	mutex.AssertExpectations(t)
	lbc.AssertNotCalled(t, "InsertRetainedQuote")
	lp.AssertNotCalled(t, "SignQuote")
	lp.AssertNotCalled(t, "GetAddress")
	eventBus.AssertNotCalled(t, "Publish")
	assert.Empty(t, result)
	require.ErrorIs(t, err, usecases.NoLiquidityError)
}

func TestAcceptQuoteUseCase_Run_ErrorHandling(t *testing.T) {
	quoteHash := "0x654321"
	now := time.Now()
	signature := "0x010203"
	quoteMock := quote.PegoutQuote{
		LbcAddress:            "0xabcd01",
		LpRskAddress:          "0xabcd02",
		BtcRefundAddress:      "hijk",
		RskRefundAddress:      "0xabcd04",
		LpBtcAddress:          "edfg",
		CallFee:               entities.NewWei(5),
		PenaltyFee:            1,
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
		ProductFeeAmount:      2,
	}
	retainedQuote := quote.RetainedPegoutQuote{
		QuoteHash:         quoteHash,
		DepositAddress:    quoteMock.LbcAddress,
		Signature:         signature,
		RequiredLiquidity: entities.NewWei(18),
		State:             quote.PegoutStateWaitingForDeposit,
	}

	lbc := new(mocks.LbcMock)
	lbc.On("GetAddress").Return("0xabcd01")
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.AcceptedPegoutQuoteEvent) bool {
		return assert.Equal(t, quoteMock, event.Quote) &&
			assert.Equal(t, retainedQuote, event.RetainedQuote) &&
			assert.Equal(t, quote.AcceptedPegoutQuoteEventId, event.Event.Id())
	}))
	mutex := new(mocks.MutexMock)
	mutex.On("Lock")
	mutex.On("Unlock")

	cases := acceptQuoteUseCaseUnexpectedErrorSetups(&quoteMock, quoteHash, signature)

	for _, c := range cases {
		quoteRepositoryMock := new(mocks.PegoutQuoteRepositoryMock)
		lp := new(mocks.ProviderMock)
		c.Value(quoteRepositoryMock, lp)
		contracts := blockchain.RskContracts{Lbc: lbc}
		useCase := pegout.NewAcceptQuoteUseCase(quoteRepositoryMock, contracts, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
		result, err := useCase.Run(context.Background(), quoteHash, "")
		quoteRepositoryMock.AssertExpectations(t)
		lp.AssertExpectations(t)
		require.Error(t, err)
		assert.Empty(t, result)
	}
}

func acceptQuoteUseCaseUnexpectedErrorSetups(quoteMock *quote.PegoutQuote, quoteHash, signature string) test.Table[func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lp *mocks.ProviderMock), error] {
	return test.Table[func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lp *mocks.ProviderMock), error]{
		{
			Value: func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lp *mocks.ProviderMock) {
				quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).
					Return(nil, assert.AnError).Once()
			},
		},
		{
			Value: func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lp *mocks.ProviderMock) {
				quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).Return(quoteMock, nil).Once()
				quoteRepository.On("GetRetainedQuote", test.AnyCtx, mock.Anything).
					Return(nil, assert.AnError).Once()
			},
		},
		{
			Value: func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lp *mocks.ProviderMock) {
				quoteRepository.On("GetQuote", test.AnyCtx, quoteHash).Return(quoteMock, nil).Once()
				quoteRepository.On("GetRetainedQuote", test.AnyCtx, quoteHash).Return(nil, nil).Once()
				lp.On("HasPegoutLiquidity", test.AnyCtx, mock.Anything).Return(assert.AnError).Once()
			},
		},
		{
			Value: func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lp *mocks.ProviderMock) {
				quoteRepository.On("GetQuote", test.AnyCtx, quoteHash).Return(quoteMock, nil).Once()
				quoteRepository.On("GetRetainedQuote", test.AnyCtx, quoteHash).Return(nil, nil).Once()
				lp.On("HasPegoutLiquidity", test.AnyCtx, mock.Anything).Return(nil).Once()
				lp.On("SignQuote", mock.Anything).Return("", assert.AnError).Once()
			},
		},
		{
			Value: func(quoteRepository *mocks.PegoutQuoteRepositoryMock, lp *mocks.ProviderMock) {
				quoteRepository.On("GetQuote", test.AnyCtx, quoteHash).Return(quoteMock, nil).Once()
				quoteRepository.On("GetRetainedQuote", test.AnyCtx, quoteHash).Return(nil, nil).Once()
				lp.On("HasPegoutLiquidity", test.AnyCtx, mock.Anything).Return(nil).Once()
				lp.On("SignQuote", mock.Anything).Return(signature, nil).Once()
				quoteRepository.EXPECT().GetPegoutCreationData(test.AnyCtx, quoteHash).Return(quote.PegoutCreationDataZeroValue()).Once()
				quoteRepository.On("InsertRetainedQuote", test.AnyCtx, mock.Anything).Return(assert.AnError).Once()
			},
		},
	}
}

func TestAcceptQuoteUseCase_Run_RetainedQuoteValidation(t *testing.T) {
	quoteHash := "0x654321"
	now := time.Now()
	signature := "0x010203"
	quoteMock := quote.PegoutQuote{
		LbcAddress:            "0xabcd01",
		LpRskAddress:          "0xabcd02",
		BtcRefundAddress:      "hijk",
		RskRefundAddress:      "0xabcd04",
		LpBtcAddress:          "edfg",
		CallFee:               entities.NewWei(5),
		PenaltyFee:            1,
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
		ProductFeeAmount:      2,
	}

	lbc := new(mocks.LbcMock)
	lbc.On("GetAddress").Return("")
	lp := new(mocks.ProviderMock)
	lp.On("HasPegoutLiquidity", test.AnyCtx, mock.Anything).Return(nil).Once()
	lp.On("SignQuote", mock.Anything).Return(signature, nil).Once()
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish").Once()
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Once()
	mutex.On("Unlock").Once()
	quoteRepositoryMock := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepositoryMock.On("GetQuote", test.AnyCtx, quoteHash).Return(&quoteMock, nil).Once()
	quoteRepositoryMock.On("GetRetainedQuote", test.AnyCtx, quoteHash).Return(nil, nil).Once()
	quoteRepositoryMock.EXPECT().GetPegoutCreationData(test.AnyCtx, quoteHash).Return(quote.PegoutCreationDataZeroValue()).Once()
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := pegout.NewAcceptQuoteUseCase(quoteRepositoryMock, contracts, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
	result, err := useCase.Run(context.Background(), quoteHash, "")
	quoteRepositoryMock.AssertExpectations(t)
	e := &validator.ValidationErrors{}
	require.ErrorAs(t, err, e)
	assert.Empty(t, result)
}
