package pegin_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var anyScript = "any script"
var acceptPeginSignature = "signature"
var acceptPeginDerivationAddress = "derivation address"
var acceptPeginQuoteHash = "c8d4ad8d5d717371b92950cbe43a6a4e891cf27bcd7603c988595866944bd9cf"
var acceptPeginQuoteHashSignature = "5f1a75f55f92c23be729adfb9eff21a00feb1ba99c5e7c2ea9c98a6430e3958f2db856b6260730b6aeeab83571bbafb77730ef1a9cb3a09ce3fa07065c8b200d1c"
var ownerAccountAddress = "0x233845a26a4dA08E16218e7B401501D048670674"
var testPeginQuote = quote.PeginQuote{
	FedBtcAddress:      "2N4qmbZNDMyHDBEBKTCP218HV1LhxCMRMti",
	LbcAddress:         "0x79568c2989232dCa1840087D73d403602364c0D4",
	LpRskAddress:       "0x0D8Fb5d32704DB2931e05DB91F64BcA6f76Ce573",
	BtcRefundAddress:   "2N58BH8rEq9Ku7HuJbZvKX6WRywdNmoVrnA",
	RskRefundAddress:   "0xd5f00ABfbEA7A0B193836CAc6833c2Ad9D06cEa8",
	LpBtcAddress:       "mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6",
	CallFee:            entities.NewWei(1),
	PenaltyFee:         entities.NewWei(1),
	ContractAddress:    "0xd5f00ABfbEA7A0B193836CAc6833c2Ad9D06cEa8",
	Data:               "",
	GasLimit:           5000,
	Nonce:              654321,
	Value:              entities.NewWei(30000),
	AgreementTimestamp: uint32(time.Now().Unix()),
	TimeForDeposit:     600,
	LpCallTime:         600,
	Confirmations:      10,
	CallOnRegister:     false,
	GasFee:             entities.NewWei(1),
	ProductFeeAmount:   entities.NewWei(10),
}

var federationInfo = rootstock.FederationInfo{
	FedSize:              1,
	FedThreshold:         2,
	PubKeys:              []string{"01", "02", "03"},
	FedAddress:           test.AnyAddress,
	ActiveFedBlockHeight: 500,
	IrisActivationHeight: 500,
	ErpKeys:              []string{"04", "05", "06"},
}
var trustedAccountRepository = new(mocks.TrustedAccountRepositoryMock)

var signingHashFunction = crypto.Keccak256

// nolint:funlen
func TestAcceptQuoteUseCase_Run(t *testing.T) {
	requiredLiquidity := entities.NewWei(9280000)
	retainedQuote := quote.RetainedPeginQuote{
		QuoteHash:         acceptPeginQuoteHash,
		DepositAddress:    acceptPeginDerivationAddress,
		Signature:         acceptPeginSignature,
		RequiredLiquidity: requiredLiquidity,
		State:             quote.PeginStateWaitingForDeposit,
	}
	creationData := quote.PeginCreationData{GasPrice: entities.NewWei(5), FeePercentage: utils.NewBigFloat64(1.24), FixedFee: entities.NewWei(100)}
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, acceptPeginQuoteHash).Return(&testPeginQuote, nil)
	quoteRepository.On("GetRetainedQuote", test.AnyCtx, acceptPeginQuoteHash).Return(nil, nil)
	quoteRepository.On("InsertRetainedQuote", test.AnyCtx, retainedQuote).Return(nil)
	quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, acceptPeginQuoteHash).Return(creationData).Once()
	bridge := new(mocks.BridgeMock)
	bridge.On("FetchFederationInfo").Return(federationInfo, nil)
	lbcParsedAddress, err := hex.DecodeString(strings.TrimPrefix(testPeginQuote.LbcAddress, "0x"))
	require.NoError(t, err)
	refundParsedAddress, lpParsedAddress := []byte{4, 5, 6}, []byte{7, 8, 9}
	parsedHash, err := hex.DecodeString(acceptPeginQuoteHash)
	require.NoError(t, err)
	bridge.On("GetFlyoverDerivationAddress", rootstock.FlyoverDerivationArgs{
		FedInfo:              federationInfo,
		LbcAdress:            lbcParsedAddress,
		UserBtcRefundAddress: refundParsedAddress,
		LpBtcAddress:         lpParsedAddress,
		QuoteHash:            parsedHash,
	}).Return(rootstock.FlyoverDerivation{Address: acceptPeginDerivationAddress, RedeemScript: anyScript}, nil)
	btc := new(mocks.BtcRpcMock)
	btc.On("DecodeAddress", testPeginQuote.BtcRefundAddress).Return(refundParsedAddress, nil)
	btc.On("DecodeAddress", testPeginQuote.LpBtcAddress).Return(lpParsedAddress, nil)
	lp := new(mocks.ProviderMock)
	lp.On("HasPeginLiquidity", test.AnyCtx, requiredLiquidity).Return(nil)
	lp.On("SignQuote", acceptPeginQuoteHash).Return(acceptPeginSignature, nil)
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.AcceptedPeginQuoteEvent) bool {
		return assert.Equal(t, testPeginQuote, event.Quote) && assert.Equal(t, retainedQuote, event.RetainedQuote) && assert.Equal(t, quote.AcceptedPeginQuoteEventId, event.Event.Id())
	})).Once()
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return().On("Unlock").Return()
	rsk := new(mocks.RootstockRpcServerMock)
	rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(50), nil)

	contracts := blockchain.RskContracts{Bridge: bridge}
	useCase := pegin.NewAcceptQuoteUseCase(quoteRepository, contracts, blockchain.Rpc{Rsk: rsk, Btc: btc}, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
	result, err := useCase.Run(context.Background(), acceptPeginQuoteHash, "")

	rsk.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	bridge.AssertExpectations(t)
	btc.AssertExpectations(t)
	lp.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	mutex.AssertExpectations(t)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, acceptPeginDerivationAddress, result.DepositAddress)
	assert.Equal(t, acceptPeginSignature, result.Signature)
}

// nolint:funlen
func TestAcceptQuoteUseCase_Run_WithoutCaptcha(t *testing.T) {
	signerMock := &mocks.SignerMock{}
	signerMock.On("Validate", mock.Anything, mock.Anything).Return(true)

	lockingCap := entities.NewWei(100000)
	trustedAccountDetails := liquidity_provider.TrustedAccountDetails{
		Address:        ownerAccountAddress,
		RbtcLockingCap: lockingCap,
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

	quoteRepository := new(mocks.PeginQuoteRepositoryMock)

	btc := new(mocks.BtcRpcMock)
	bridge := new(mocks.BridgeMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	rsk := new(mocks.RootstockRpcServerMock)
	lp := new(mocks.ProviderMock)
	lp.On("GetSigner").Return(signerMock)
	contracts := blockchain.RskContracts{Bridge: bridge}

	t.Run("happy path", func(t *testing.T) {
		requiredLiquidity := entities.NewWei(9280000)
		retainedQuote := quote.RetainedPeginQuote{
			QuoteHash:           acceptPeginQuoteHash,
			DepositAddress:      acceptPeginDerivationAddress,
			Signature:           acceptPeginSignature,
			RequiredLiquidity:   requiredLiquidity,
			State:               quote.PeginStateWaitingForDeposit,
			OwnerAccountAddress: ownerAccountAddress,
		}

		creationData := quote.PeginCreationData{GasPrice: entities.NewWei(5), FeePercentage: utils.NewBigFloat64(1.24), FixedFee: entities.NewWei(100)}

		quoteRepository.On("GetQuote", test.AnyCtx, acceptPeginQuoteHash).Return(&testPeginQuote, nil)
		quoteRepository.On("GetRetainedQuote", test.AnyCtx, acceptPeginQuoteHash).Return(nil, nil)
		quoteRepository.On("InsertRetainedQuote", test.AnyCtx, retainedQuote).Return(nil)
		quoteRepository.On("GetRetainedQuotesForAddress", test.AnyCtx, ownerAccountAddress, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations).Return([]quote.RetainedPeginQuote{}, nil)
		quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, acceptPeginQuoteHash).Return(creationData).Once()

		bridge.On("FetchFederationInfo").Return(federationInfo, nil)
		lbcParsedAddress, err := hex.DecodeString(strings.TrimPrefix(testPeginQuote.LbcAddress, "0x"))
		require.NoError(t, err)
		refundParsedAddress, lpParsedAddress := []byte{4, 5, 6}, []byte{7, 8, 9}
		parsedHash, err := hex.DecodeString(acceptPeginQuoteHash)
		require.NoError(t, err)
		bridge.On("GetFlyoverDerivationAddress", rootstock.FlyoverDerivationArgs{
			FedInfo:              federationInfo,
			LbcAdress:            lbcParsedAddress,
			UserBtcRefundAddress: refundParsedAddress,
			LpBtcAddress:         lpParsedAddress,
			QuoteHash:            parsedHash,
		}).Return(rootstock.FlyoverDerivation{Address: acceptPeginDerivationAddress, RedeemScript: anyScript}, nil)

		btc.On("DecodeAddress", testPeginQuote.BtcRefundAddress).Return(refundParsedAddress, nil)
		btc.On("DecodeAddress", testPeginQuote.LpBtcAddress).Return(lpParsedAddress, nil)

		lp.On("HasPeginLiquidity", test.AnyCtx, requiredLiquidity).Return(nil)
		lp.On("SignQuote", acceptPeginQuoteHash).Return(acceptPeginSignature, nil)

		eventBus.On("Publish", mock.MatchedBy(func(event quote.AcceptedPeginQuoteEvent) bool {
			return assert.Equal(t, testPeginQuote, event.Quote) && assert.Equal(t, retainedQuote, event.RetainedQuote) && assert.Equal(t, quote.AcceptedPeginQuoteEventId, event.Event.Id())
		})).Once()
		mutex.On("Lock").Return().On("Unlock").Return()
		rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(50), nil)

		useCase := pegin.NewAcceptQuoteUseCase(quoteRepository, contracts, blockchain.Rpc{Rsk: rsk, Btc: btc}, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
		result, err := useCase.Run(context.Background(), acceptPeginQuoteHash, acceptPeginQuoteHashSignature)

		rsk.AssertExpectations(t)
		quoteRepository.AssertExpectations(t)
		trustedAccountRepository.AssertExpectations(t)
		bridge.AssertExpectations(t)
		btc.AssertExpectations(t)
		lp.AssertExpectations(t)
		eventBus.AssertExpectations(t)
		mutex.AssertExpectations(t)
		require.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Equal(t, acceptPeginDerivationAddress, result.DepositAddress)
		assert.Equal(t, acceptPeginSignature, result.Signature)
	})

	t.Run("invalid signature", func(t *testing.T) {
		// Set up the pegin quote
		newQuote := testPeginQuote

		// Set up a well-formed but invalid signature
		invalidSignature := "5f1a75f55f92c23be729adfb9eff21a00feb1ba99c5e7c2ea9c98a6430e3958f2db856b6260730b6aeeab83571bbafb77730ef1a9cb3a09ce3fa07065c8b200d1d"

		// Mock just what's necessary for the test path we're validating
		quoteRepository.On("GetQuote", mock.Anything, mock.Anything).Return(&newQuote, nil)

		// We don't expect these to be called because signature validation should fail first
		quoteRepository.On("GetRetainedQuote", mock.Anything, mock.Anything).Return(nil, nil).Maybe()

		useCase := pegin.NewAcceptQuoteUseCase(quoteRepository, contracts, blockchain.Rpc{Rsk: rsk, Btc: btc}, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)

		result, err := useCase.Run(context.Background(), acceptPeginQuoteHash, invalidSignature)

		require.Error(t, err)

		assert.Contains(t, err.Error(), "recovery failed")
		assert.Empty(t, result)

		quoteRepository.AssertNotCalled(t, "InsertRetainedQuote")
		quoteRepository.AssertNotCalled(t, "GetRetainedQuotesForAddress")
		lp.AssertNotCalled(t, "HasPeginLiquidity")
		lp.AssertNotCalled(t, "SignQuote")
		eventBus.AssertNotCalled(t, "Publish")
	})

	t.Run("locking cap exceeded", func(t *testing.T) {
		// Create two existing quotes that together with the new quote will exceed the locking cap
		existingQuote1 := quote.RetainedPeginQuote{
			QuoteHash:           "existing-hash-1",
			DepositAddress:      "existing-address-1",
			Signature:           "existing-signature-1",
			RequiredLiquidity:   entities.NewWei(40000),
			State:               quote.PeginStateWaitingForDeposit,
			OwnerAccountAddress: ownerAccountAddress,
		}

		existingQuote2 := quote.RetainedPeginQuote{
			QuoteHash:           "existing-hash-2",
			DepositAddress:      "existing-address-2",
			Signature:           "existing-signature-2",
			RequiredLiquidity:   entities.NewWei(50000),
			State:               quote.PeginStateWaitingForDeposit,
			OwnerAccountAddress: ownerAccountAddress,
		}

		// Set up the new pegin quote which would push the total over the locking cap
		newQuote := testPeginQuote
		newQuote.Value = entities.NewWei(30000)
		newQuote.GasFee = entities.NewWei(20000)
		// Total required: 40000 + 50000 + 30000 + 20000 = 140000 > locking cap of 100000

		quoteRepository := new(mocks.PeginQuoteRepositoryMock)
		quoteRepository.On("GetQuote", mock.Anything, acceptPeginQuoteHash).Return(&newQuote, nil)
		quoteRepository.On("GetRetainedQuote", mock.Anything, acceptPeginQuoteHash).Return(nil, nil)
		quoteRepository.On("GetRetainedQuotesForAddress", mock.Anything, ownerAccountAddress, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations).Return([]quote.RetainedPeginQuote{existingQuote1, existingQuote2}, nil)

		useCase := pegin.NewAcceptQuoteUseCase(quoteRepository, contracts, blockchain.Rpc{Rsk: rsk, Btc: btc}, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
		result, err := useCase.Run(context.Background(), acceptPeginQuoteHash, acceptPeginQuoteHashSignature)

		require.Error(t, err)
		require.ErrorIs(t, err, usecases.LockingCapExceededError)
		assert.Empty(t, result)
	})
}

func TestAcceptQuoteUseCase_Run_AlreadyAccepted(t *testing.T) {
	retainedQuote := quote.RetainedPeginQuote{
		QuoteHash:         acceptPeginQuoteHash,
		DepositAddress:    acceptPeginDerivationAddress,
		Signature:         acceptPeginSignature,
		RequiredLiquidity: entities.NewWei(9280000),
		State:             quote.PeginStateWaitingForDeposit,
	}
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, acceptPeginQuoteHash).Return(&testPeginQuote, nil)
	quoteRepository.On("GetRetainedQuote", test.AnyCtx, acceptPeginQuoteHash).Return(&retainedQuote, nil)

	bridge := new(mocks.BridgeMock)
	btc := new(mocks.BtcRpcMock)
	lp := new(mocks.ProviderMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return()
	mutex.On("Unlock").Return()
	rsk := new(mocks.RootstockRpcServerMock)

	contracts := blockchain.RskContracts{Bridge: bridge}
	rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
	useCase := pegin.NewAcceptQuoteUseCase(quoteRepository, contracts, rpc, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
	result, err := useCase.Run(context.Background(), acceptPeginQuoteHash, "")

	rsk.AssertNotCalled(t, "GasPrice")
	quoteRepository.AssertExpectations(t)
	quoteRepository.AssertNotCalled(t, "InsertRetainedQuote")
	btc.AssertNotCalled(t, "DecodeAddress")
	bridge.AssertNotCalled(t, "GetFlyoverDerivationAddress")
	bridge.AssertNotCalled(t, "FetchFederationInfo")
	lp.AssertNotCalled(t, "HasPeginLiquidity")
	lp.AssertNotCalled(t, "SignQuote")
	eventBus.AssertNotCalled(t, "Publish")
	mutex.AssertExpectations(t)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, acceptPeginDerivationAddress, result.DepositAddress)
	assert.Equal(t, acceptPeginSignature, result.Signature)
}

func TestAcceptQuoteUseCase_Run_QuoteNotFound(t *testing.T) {
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, acceptPeginQuoteHash).Return(nil, nil)

	bridge := new(mocks.BridgeMock)
	btc := new(mocks.BtcRpcMock)
	lp := new(mocks.ProviderMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	rsk := new(mocks.RootstockRpcServerMock)

	contracts := blockchain.RskContracts{Bridge: bridge}
	rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
	useCase := pegin.NewAcceptQuoteUseCase(quoteRepository, contracts, rpc, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
	result, err := useCase.Run(context.Background(), acceptPeginQuoteHash, "")

	rsk.AssertNotCalled(t, "GasPrice")
	quoteRepository.AssertExpectations(t)
	btc.AssertNotCalled(t, "DecodeAddress")
	bridge.AssertNotCalled(t, "GetFlyoverDerivationAddress")
	bridge.AssertNotCalled(t, "FetchFederationInfo")
	lp.AssertNotCalled(t, "HasPeginLiquidity")
	lp.AssertNotCalled(t, "SignQuote")
	eventBus.AssertNotCalled(t, "Publish")
	mutex.AssertNotCalled(t, "Unlock")
	mutex.AssertNotCalled(t, "Lock")
	require.ErrorIs(t, err, usecases.QuoteNotFoundError)
	assert.Empty(t, result)
}

func TestAcceptQuoteUseCase_Run_ExpiredQuote(t *testing.T) {
	expiredQuote := testPeginQuote
	expiredQuote.AgreementTimestamp = uint32(time.Now().Unix()) - 1000
	expiredQuote.TimeForDeposit = 500
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, acceptPeginQuoteHash).Return(&expiredQuote, nil)

	bridge := new(mocks.BridgeMock)
	btc := new(mocks.BtcRpcMock)
	lp := new(mocks.ProviderMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	rsk := new(mocks.RootstockRpcServerMock)

	contracts := blockchain.RskContracts{Bridge: bridge}
	rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
	useCase := pegin.NewAcceptQuoteUseCase(quoteRepository, contracts, rpc, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
	result, err := useCase.Run(context.Background(), acceptPeginQuoteHash, "")

	rsk.AssertNotCalled(t, "GasPrice")
	quoteRepository.AssertExpectations(t)
	btc.AssertNotCalled(t, "DecodeAddress")
	bridge.AssertNotCalled(t, "GetFlyoverDerivationAddress")
	bridge.AssertNotCalled(t, "FetchFederationInfo")
	lp.AssertNotCalled(t, "HasPeginLiquidity")
	lp.AssertNotCalled(t, "SignQuote")
	eventBus.AssertNotCalled(t, "Publish")
	mutex.AssertNotCalled(t, "Unlock")
	mutex.AssertNotCalled(t, "Lock")
	require.ErrorIs(t, err, usecases.ExpiredQuoteError)
	assert.Empty(t, result)
}

func TestAcceptQuoteUseCase_Run_NoLiquidity(t *testing.T) {
	requiredLiquidity := entities.NewWei(9280000)
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetQuote", test.AnyCtx, acceptPeginQuoteHash).Return(&testPeginQuote, nil)
	quoteRepository.On("GetRetainedQuote", test.AnyCtx, acceptPeginQuoteHash).Return(nil, nil)
	bridge := new(mocks.BridgeMock)
	bridge.On("FetchFederationInfo").Return(federationInfo, nil)
	bridge.On("GetFlyoverDerivationAddress", mock.Anything).Return(rootstock.FlyoverDerivation{
		Address:      "derivation address",
		RedeemScript: anyScript,
	}, nil)
	btc := new(mocks.BtcRpcMock)
	btc.On("DecodeAddress", testPeginQuote.BtcRefundAddress).Return([]byte{4, 5, 6}, nil)
	btc.On("DecodeAddress", testPeginQuote.LpBtcAddress).Return([]byte{7, 8, 9}, nil)
	lp := new(mocks.ProviderMock)
	lp.On("HasPeginLiquidity", test.AnyCtx, requiredLiquidity).Return(assert.AnError)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return()
	mutex.On("Unlock").Return()
	rsk := new(mocks.RootstockRpcServerMock)
	rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(50), nil)

	contracts := blockchain.RskContracts{Bridge: bridge}
	rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
	useCase := pegin.NewAcceptQuoteUseCase(quoteRepository, contracts, rpc, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
	result, err := useCase.Run(context.Background(), acceptPeginQuoteHash, "")

	rsk.AssertExpectations(t)
	quoteRepository.AssertExpectations(t)
	btc.AssertExpectations(t)
	bridge.AssertExpectations(t)
	lp.AssertExpectations(t)
	eventBus.AssertNotCalled(t, "Publish")
	mutex.AssertExpectations(t)
	require.ErrorIs(t, err, usecases.NoLiquidityError)
	assert.Empty(t, result)
}

func TestAcceptQuoteUseCase_Run_ErrorHandling(t *testing.T) {
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.Anything).Return()
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return()
	mutex.On("Unlock").Return()

	setups := acceptQuoteUseCaseUnexpectedErrorSetups()
	for _, setup := range setups {
		quoteRepository := new(mocks.PeginQuoteRepositoryMock)
		bridge := new(mocks.BridgeMock)
		btc := new(mocks.BtcRpcMock)
		lp := new(mocks.ProviderMock)
		rsk := new(mocks.RootstockRpcServerMock)
		caseHash := acceptPeginQuoteHash
		setup(&caseHash, quoteRepository, bridge, btc, lp, rsk)
		contracts := blockchain.RskContracts{Bridge: bridge}
		rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
		useCase := pegin.NewAcceptQuoteUseCase(quoteRepository, contracts, rpc, lp, lp, eventBus, mutex, trustedAccountRepository, signingHashFunction)
		result, err := useCase.Run(context.Background(), caseHash, "")

		rsk.AssertExpectations(t)
		quoteRepository.AssertExpectations(t)
		bridge.AssertExpectations(t)
		btc.AssertExpectations(t)
		lp.AssertExpectations(t)
		require.Error(t, err)
		assert.Empty(t, result)
	}
}

// nolint:funlen
func acceptQuoteUseCaseUnexpectedErrorSetups() []func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock,
	bridge *mocks.BridgeMock, btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
	derivation := rootstock.FlyoverDerivation{Address: test.AnyAddress, RedeemScript: anyScript}
	return []func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
		btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock){
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", test.AnyCtx, mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", test.AnyCtx, mock.Anything).Return(nil, nil).Once()
			*quoteHash = "malformed hash"
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", test.AnyCtx, mock.Anything).Return(nil, nil).Once()
			btc.On("DecodeAddress", mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", test.AnyCtx, mock.Anything).Return(nil, nil).Once()
			btc.On("DecodeAddress", mock.Anything).Return([]byte{1}, nil).Once()
			btc.On("DecodeAddress", mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			caseQuote := testPeginQuote
			caseQuote.LbcAddress = "malformed address"
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).Return(&caseQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", test.AnyCtx, mock.Anything).Return(nil, nil).Once()
			btc.On("DecodeAddress", mock.Anything).Return([]byte{1}, nil).Once()
			btc.On("DecodeAddress", mock.Anything).Return([]byte{2}, nil).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", test.AnyCtx, mock.Anything).Return(nil, nil).Once()
			btc.On("DecodeAddress", mock.Anything).Return([]byte{1}, nil).Once()
			btc.On("DecodeAddress", mock.Anything).Return([]byte{2}, nil).Once()
			bridge.On("FetchFederationInfo").Return(rootstock.FederationInfo{}, assert.AnError).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", test.AnyCtx, mock.Anything).Return(nil, nil).Once()
			btc.On("DecodeAddress", mock.Anything).Return([]byte{1}, nil).Once()
			btc.On("DecodeAddress", mock.Anything).Return([]byte{2}, nil).Once()
			bridge.On("FetchFederationInfo").Return(federationInfo, nil).Once()
			bridge.On("GetFlyoverDerivationAddress", mock.Anything).Return(rootstock.FlyoverDerivation{}, assert.AnError).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", test.AnyCtx, mock.Anything).Return(nil, nil).Once()
			btc.On("DecodeAddress", mock.Anything).Return([]byte{1}, nil).Once()
			btc.On("DecodeAddress", mock.Anything).Return([]byte{2}, nil).Once()
			bridge.On("FetchFederationInfo").Return(federationInfo, nil).Once()
			bridge.On("GetFlyoverDerivationAddress", mock.Anything).Return(derivation, nil).Once()
			rsk.On("GasPrice", test.AnyCtx).Return(nil, assert.AnError).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", test.AnyCtx, mock.Anything).Return(nil, nil).Once()
			btc.On("DecodeAddress", mock.Anything).Return([]byte{1}, nil).Once()
			btc.On("DecodeAddress", mock.Anything).Return([]byte{2}, nil).Once()
			bridge.On("FetchFederationInfo").Return(federationInfo, nil).Once()
			bridge.On("GetFlyoverDerivationAddress", mock.Anything).Return(derivation, nil).Once()
			rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(1), nil).Once()
			lp.On("HasPeginLiquidity", test.AnyCtx, mock.Anything).Return(nil).Once()
			lp.On("SignQuote", mock.Anything).Return("", assert.AnError).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", test.AnyCtx, mock.Anything).Return(nil, nil).Once()
			btc.On("DecodeAddress", mock.Anything).Return([]byte{1}, nil).Once()
			btc.On("DecodeAddress", mock.Anything).Return([]byte{2}, nil).Once()
			bridge.On("FetchFederationInfo").Return(federationInfo, nil).Once()
			rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(1), nil).Once()
			lp.On("HasPeginLiquidity", test.AnyCtx, mock.Anything).Return(nil).Once()
			// set derivation and signature to empty to malform the retained quote
			bridge.On("GetFlyoverDerivationAddress", mock.Anything).Return(rootstock.FlyoverDerivation{}, nil).Once()
			lp.On("SignQuote", mock.Anything).Return("", nil).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", test.AnyCtx, mock.Anything).Return(nil, nil).Once()
			btc.On("DecodeAddress", mock.Anything).Return([]byte{1}, nil).Once()
			btc.On("DecodeAddress", mock.Anything).Return([]byte{2}, nil).Once()
			bridge.On("FetchFederationInfo").Return(federationInfo, nil).Once()
			bridge.On("GetFlyoverDerivationAddress", mock.Anything).Return(derivation, nil).Once()
			rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(1), nil).Once()
			lp.On("HasPeginLiquidity", test.AnyCtx, mock.Anything).Return(nil).Once()
			lp.On("SignQuote", mock.Anything).Return("signature", nil).Once()
			quoteRepository.On("InsertRetainedQuote", test.AnyCtx, mock.Anything).Return(assert.AnError).Once()
		},
	}
}
