package pegout_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-playground/validator/v10"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
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

var acceptPegoutQuoteHash = "c8d4ad8d5d717371b92950cbe43a6a4e891cf27bcd7603c988595866944bd9cf"
var acceptPegoutQuoteHashSignature = "5f1a75f55f92c23be729adfb9eff21a00feb1ba99c5e7c2ea9c98a6430e3958f2db856b6260730b6aeeab83571bbafb77730ef1a9cb3a09ce3fa07065c8b200d1c"
var ownerAccountAddress = "0x233845a26a4dA08E16218e7B401501D048670674"

func TestAcceptQuoteUseCase_Run_Paused(t *testing.T) {
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	bridge := new(mocks.BridgeMock)
	lp := new(mocks.ProviderMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: true, Since: 5, Reason: "test"}, nil)
	pegoutContract.EXPECT().GetAddress().Return("test-contract")

	contracts := blockchain.RskContracts{Bridge: bridge, PegOut: pegoutContract}
	useCase := pegout.NewAcceptQuoteUseCase(quoteRepository, contracts, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
	result, err := useCase.Run(context.Background(), acceptPegoutQuoteHash, "")
	assert.Empty(t, result)
	require.ErrorIs(t, err, blockchain.ContractPausedError)
}

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
		ProductFeeAmount:      entities.NewWei(2),
	}
	retainedQuote := quote.RetainedPegoutQuote{
		QuoteHash: quoteHash, DepositAddress: quoteMock.LbcAddress, Signature: signature,
		RequiredLiquidity: entities.NewWei(18), State: quote.PegoutStateWaitingForDeposit,
	}
	creationData := quote.PegoutCreationData{FeeRate: utils.NewBigFloat64(1.5), FeePercentage: utils.NewBigFloat64(12.5), GasPrice: entities.NewWei(1), FixedFee: entities.NewWei(100)}
	quoteRepositoryMock := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepositoryMock.On("InsertRetainedQuote", test.AnyCtx, retainedQuote).Return(nil).Once()
	quoteRepositoryMock.On("GetQuote", test.AnyCtx, quoteHash).Return(&quoteMock, nil).Once()
	quoteRepositoryMock.On("GetRetainedQuote", test.AnyCtx, quoteHash).Return(nil, nil).Once()
	quoteRepositoryMock.EXPECT().GetPegoutCreationData(test.AnyCtx, quoteHash).Return(creationData).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.On("GetAddress").Return("0xabcd01").Once()
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
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
	contracts := blockchain.RskContracts{PegOut: pegoutContract}
	useCase := pegout.NewAcceptQuoteUseCase(quoteRepositoryMock, contracts, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
	result, err := useCase.Run(context.Background(), quoteHash, "")
	quoteRepositoryMock.AssertExpectations(t)
	pegoutContract.AssertExpectations(t)
	lp.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	mutex.AssertExpectations(t)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, quoteMock.LbcAddress, result.DepositAddress)
	assert.Equal(t, signature, result.Signature)
}

// nolint:funlen
func TestAcceptQuoteUseCase_Run_WithoutCaptcha(t *testing.T) {
	signerMock := &mocks.SignerMock{}
	signerMock.On("Validate", mock.Anything, mock.Anything).Return(true)

	lockingCap := entities.NewWei(100000)
	trustedAccountDetails := liquidity_provider.TrustedAccountDetails{
		Address:       ownerAccountAddress,
		BtcLockingCap: lockingCap,
	}
	trustedAccountBytes, err := json.Marshal(trustedAccountDetails)
	require.NoError(t, err)
	trustedAccountHash := hex.EncodeToString(crypto.Keccak256(trustedAccountBytes))

	accountSignature := "d1a9fe0de659875bc75252e6f5a73529ed6a5d88c9d97853ebf2ccc6e3080ecc423eee543470a80d373f1abb3a4f746264b47dda53252ddfc5d65989c1af34401c"
	trustedAccountRepository.On("GetTrustedAccount", mock.Anything, ownerAccountAddress).Return(&entities.Signed[liquidity_provider.TrustedAccountDetails]{
		Value:     trustedAccountDetails,
		Signature: accountSignature,
		Hash:      trustedAccountHash,
	}, nil)

	quoteRepositoryMock := new(mocks.PegoutQuoteRepositoryMock)
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	lp := new(mocks.ProviderMock)
	lp.On("GetSigner").Return(signerMock)
	contracts := blockchain.RskContracts{PegOut: pegoutContract}

	now := time.Now()
	quoteMock := quote.PegoutQuote{
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
		ProductFeeAmount:      entities.NewWei(2),
	}

	t.Run("happy path", func(t *testing.T) {
		quoteHash := acceptPegoutQuoteHash
		signature := "0x010203"

		requiredLiquidity := entities.NewWei(18)
		retainedQuote := quote.RetainedPegoutQuote{
			QuoteHash:           quoteHash,
			DepositAddress:      quoteMock.LbcAddress,
			Signature:           signature,
			RequiredLiquidity:   requiredLiquidity,
			State:               quote.PegoutStateWaitingForDeposit,
			OwnerAccountAddress: ownerAccountAddress,
		}

		creationData := quote.PegoutCreationData{FeeRate: utils.NewBigFloat64(1.5), FeePercentage: utils.NewBigFloat64(12.5), GasPrice: entities.NewWei(1), FixedFee: entities.NewWei(100)}

		quoteRepositoryMock.On("GetQuote", test.AnyCtx, quoteHash).Return(&quoteMock, nil)
		quoteRepositoryMock.On("GetRetainedQuote", test.AnyCtx, quoteHash).Return(nil, nil)
		quoteRepositoryMock.On("InsertRetainedQuote", test.AnyCtx, retainedQuote).Return(nil)
		quoteRepositoryMock.On("GetRetainedQuotesForAddress", test.AnyCtx, ownerAccountAddress, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations).Return([]quote.RetainedPegoutQuote{}, nil)
		quoteRepositoryMock.EXPECT().GetPegoutCreationData(test.AnyCtx, quoteHash).Return(creationData).Once()

		pegoutContract.On("GetAddress").Return("0xabcd01").Once()
		lp.On("HasPegoutLiquidity", test.AnyCtx, requiredLiquidity).Return(nil)
		lp.On("SignQuote", quoteHash).Return(signature, nil)

		eventBus.On("Publish", mock.MatchedBy(func(event quote.AcceptedPegoutQuoteEvent) bool {
			return assert.Equal(t, quoteMock, event.Quote) && assert.Equal(t, retainedQuote, event.RetainedQuote) && assert.Equal(t, quote.AcceptedPegoutQuoteEventId, event.Event.Id()) && assert.Equal(t, creationData, event.CreationData)
		})).Once()
		mutex.On("Lock").Return().On("Unlock").Return()

		useCase := pegout.NewAcceptQuoteUseCase(quoteRepositoryMock, contracts, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
		result, err := useCase.Run(context.Background(), quoteHash, acceptPegoutQuoteHashSignature)

		quoteRepositoryMock.AssertExpectations(t)
		trustedAccountRepository.AssertExpectations(t)
		pegoutContract.AssertExpectations(t)
		lp.AssertExpectations(t)
		eventBus.AssertExpectations(t)
		mutex.AssertExpectations(t)
		require.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Equal(t, quoteMock.LbcAddress, result.DepositAddress)
		assert.Equal(t, signature, result.Signature)
	})

	t.Run("invalid signature", func(t *testing.T) {
		quoteHash := acceptPegoutQuoteHash

		// Set up a well-formed but invalid signature
		invalidSignature := "5f1a75f55f92c23be729adfb9eff21a00feb1ba99c5e7c2ea9c98a6430e3958f2db856b6260730b6aeeab83571bbafb77730ef1a9cb3a09ce3fa07065c8b200d1d"

		quoteRepositoryMock.On("GetQuote", mock.Anything, mock.Anything).Return(&quoteMock, nil)

		useCase := pegout.NewAcceptQuoteUseCase(quoteRepositoryMock, contracts, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)

		result, err := useCase.Run(context.Background(), quoteHash, invalidSignature)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "recovery failed")
		assert.Empty(t, result)

		quoteRepositoryMock.AssertNotCalled(t, "InsertRetainedQuote")
		quoteRepositoryMock.AssertNotCalled(t, "GetRetainedQuotesForAddress")
		lp.AssertNotCalled(t, "HasPegoutLiquidity")
		lp.AssertNotCalled(t, "SignQuote")
		eventBus.AssertNotCalled(t, "Publish")
	})

	t.Run("locking cap exceeded", func(t *testing.T) {
		// Create two existing quotes that together with the new quote will exceed the locking cap
		existingQuote1 := quote.RetainedPegoutQuote{
			QuoteHash:           "existing-hash-1",
			DepositAddress:      "existing-address-1",
			Signature:           "existing-signature-1",
			RequiredLiquidity:   entities.NewWei(40000),
			State:               quote.PegoutStateWaitingForDeposit,
			OwnerAccountAddress: ownerAccountAddress,
		}

		existingQuote2 := quote.RetainedPegoutQuote{
			QuoteHash:           "existing-hash-2",
			DepositAddress:      "existing-address-2",
			Signature:           "existing-signature-2",
			RequiredLiquidity:   entities.NewWei(50000),
			State:               quote.PegoutStateWaitingForDeposit,
			OwnerAccountAddress: ownerAccountAddress,
		}

		// Set up the new pegout quote which would push the total over the locking cap
		lockingCapQuote := quote.PegoutQuote{
			LbcAddress:            "0xabcd01",
			LpRskAddress:          "0xabcd02",
			BtcRefundAddress:      "hijk",
			RskRefundAddress:      "0xabcd04",
			LpBtcAddress:          "edfg",
			CallFee:               entities.NewWei(5000),
			PenaltyFee:            entities.NewWei(1),
			Nonce:                 1,
			DepositAddress:        "address",
			Value:                 entities.NewWei(30000),
			AgreementTimestamp:    uint32(now.Unix()),
			DepositDateLimit:      uint32(now.Unix() + 600),
			DepositConfirmations:  1,
			TransferConfirmations: 1,
			TransferTime:          600,
			ExpireDate:            uint32(now.Unix() + 600),
			ExpireBlock:           1,
			GasFee:                entities.NewWei(20000),
			ProductFeeAmount:      entities.NewWei(2),
		}
		// Total required: 40000 + 50000 + 30000 + 20000 + 5000 = 145000 > locking cap of 100000

		quoteRepositoryMock := new(mocks.PegoutQuoteRepositoryMock)
		quoteRepositoryMock.On("GetQuote", mock.Anything, acceptPegoutQuoteHash).Return(&lockingCapQuote, nil)
		quoteRepositoryMock.On("GetRetainedQuote", mock.Anything, acceptPegoutQuoteHash).Return(nil, nil)
		quoteRepositoryMock.On("GetRetainedQuotesForAddress", mock.Anything, ownerAccountAddress, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations).Return([]quote.RetainedPegoutQuote{existingQuote1, existingQuote2}, nil)

		useCase := pegout.NewAcceptQuoteUseCase(quoteRepositoryMock, contracts, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
		result, err := useCase.Run(context.Background(), acceptPegoutQuoteHash, acceptPegoutQuoteHashSignature)

		require.Error(t, err)
		require.ErrorIs(t, err, usecases.LockingCapExceededError)
		assert.Empty(t, result)
	})
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
		PenaltyFee:            entities.NewWei(1),
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
		ProductFeeAmount:      entities.NewWei(1),
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
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	lp := new(mocks.ProviderMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return()
	mutex.On("Unlock").Return()
	contracts := blockchain.RskContracts{PegOut: pegoutContract}
	useCase := pegout.NewAcceptQuoteUseCase(quoteRepositoryMock, contracts, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
	result, err := useCase.Run(context.Background(), quoteHash, "")
	quoteRepositoryMock.AssertExpectations(t)
	pegoutContract.AssertNotCalled(t, "GetAddress")
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
		PenaltyFee:            entities.NewWei(1),
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
		ProductFeeAmount:      entities.NewWei(1),
	}
	quoteRepositoryMock := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepositoryMock.On("GetQuote", test.AnyCtx, quoteHash).Return(&quoteMock, nil).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	lp := new(mocks.ProviderMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	contracts := blockchain.RskContracts{PegOut: pegoutContract}
	useCase := pegout.NewAcceptQuoteUseCase(quoteRepositoryMock, contracts, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
	result, err := useCase.Run(context.Background(), quoteHash, "")
	quoteRepositoryMock.AssertExpectations(t)
	quoteRepositoryMock.AssertNotCalled(t, "GetRetainedQuote")
	pegoutContract.AssertNotCalled(t, "GetAddress")
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
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	lp := new(mocks.ProviderMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	contracts := blockchain.RskContracts{PegOut: pegoutContract}
	useCase := pegout.NewAcceptQuoteUseCase(quoteRepositoryMock, contracts, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
	result, err := useCase.Run(context.Background(), quoteHash, "")
	quoteRepositoryMock.AssertExpectations(t)
	quoteRepositoryMock.AssertNotCalled(t, "GetRetainedQuote")
	pegoutContract.AssertNotCalled(t, "GetAddress")
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
		PenaltyFee:            entities.NewWei(1),
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
		ProductFeeAmount:      entities.NewWei(8),
	}
	quoteRepositoryMock := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepositoryMock.On("GetQuote", test.AnyCtx, quoteHash).Return(&quoteMock, nil).Once()
	quoteRepositoryMock.On("GetRetainedQuote", test.AnyCtx, quoteHash).Return(nil, nil).Once()
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	lp := new(mocks.ProviderMock)
	lp.On("HasPegoutLiquidity", test.AnyCtx, entities.NewWei(65)).Return(usecases.NoLiquidityError).Once()
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Once()
	mutex.On("Unlock").Once()
	contracts := blockchain.RskContracts{PegOut: pegoutContract}
	useCase := pegout.NewAcceptQuoteUseCase(quoteRepositoryMock, contracts, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
	result, err := useCase.Run(context.Background(), quoteHash, "")
	quoteRepositoryMock.AssertExpectations(t)
	lp.AssertExpectations(t)
	mutex.AssertExpectations(t)
	pegoutContract.AssertNotCalled(t, "InsertRetainedQuote")
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
		ProductFeeAmount:      entities.NewWei(2),
	}
	retainedQuote := quote.RetainedPegoutQuote{
		QuoteHash:         quoteHash,
		DepositAddress:    quoteMock.LbcAddress,
		Signature:         signature,
		RequiredLiquidity: entities.NewWei(18),
		State:             quote.PegoutStateWaitingForDeposit,
	}

	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.On("GetAddress").Return("0xabcd01")
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
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
		contracts := blockchain.RskContracts{PegOut: pegoutContract}
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
		ProductFeeAmount:      entities.NewWei(2),
	}

	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.On("GetAddress").Return("")
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
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
	contracts := blockchain.RskContracts{PegOut: pegoutContract}
	useCase := pegout.NewAcceptQuoteUseCase(quoteRepositoryMock, contracts, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
	result, err := useCase.Run(context.Background(), quoteHash, "")
	quoteRepositoryMock.AssertExpectations(t)
	e := &validator.ValidationErrors{}
	require.ErrorAs(t, err, e)
	assert.Empty(t, result)
}
