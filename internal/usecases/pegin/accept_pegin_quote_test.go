package pegin_test

import (
	"context"
	"encoding/hex"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

var acceptPeginSignature = "signature"
var acceptPeginDerivationAddress = "derivation address"
var acceptPeginQuoteHash = "1a1b1c"
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
	ProductFeeAmount:   10,
}

var federationInfo = blockchain.FederationInfo{
	FedSize:              1,
	FedThreshold:         2,
	PubKeys:              []string{"01", "02", "03"},
	FedAddress:           test.AnyAddress,
	ActiveFedBlockHeight: 500,
	IrisActivationHeight: 500,
	ErpKeys:              []string{"04", "05", "06"},
}

func TestAcceptQuoteUseCase_Run(t *testing.T) {
	requiredLiquidity := entities.NewWei(9280000)
	retainedQuote := quote.RetainedPeginQuote{
		QuoteHash:         acceptPeginQuoteHash,
		DepositAddress:    acceptPeginDerivationAddress,
		Signature:         acceptPeginSignature,
		RequiredLiquidity: requiredLiquidity,
		State:             quote.PeginStateWaitingForDeposit,
	}
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), acceptPeginQuoteHash).Return(&testPeginQuote, nil)
	quoteRepository.On("GetRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), acceptPeginQuoteHash).Return(nil, nil)
	quoteRepository.On("InsertRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuote).Return(nil)
	bridge := new(mocks.BridgeMock)
	bridge.On("FetchFederationInfo").Return(federationInfo, nil)
	lbcParsedAddress, _ := hex.DecodeString(strings.TrimPrefix(testPeginQuote.LbcAddress, "0x"))
	refundParsedAddress := []byte{4, 5, 6}
	lpParsedAddress := []byte{7, 8, 9}
	parsedHash, _ := hex.DecodeString(acceptPeginQuoteHash)
	bridge.On("GetFlyoverDerivationAddress", blockchain.FlyoverDerivationArgs{
		FedInfo:              federationInfo,
		LbcAdress:            lbcParsedAddress,
		UserBtcRefundAddress: refundParsedAddress,
		LpBtcAddress:         lpParsedAddress,
		QuoteHash:            parsedHash,
	}).Return(blockchain.FlyoverDerivation{Address: acceptPeginDerivationAddress, RedeemScript: "any script"}, nil)
	btc := new(mocks.BtcRpcMock)
	btc.On("DecodeAddress", testPeginQuote.BtcRefundAddress, true).Return(refundParsedAddress, nil)
	btc.On("DecodeAddress", testPeginQuote.LpBtcAddress, true).Return(lpParsedAddress, nil)
	lp := new(mocks.ProviderMock)
	lp.On("HasPeginLiquidity", mock.AnythingOfType("context.backgroundCtx"), requiredLiquidity).Return(nil)
	lp.On("SignQuote", acceptPeginQuoteHash).Return(acceptPeginSignature, nil)
	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.MatchedBy(func(event quote.AcceptedPeginQuoteEvent) bool {
		return assert.Equal(t, testPeginQuote, event.Quote) && assert.Equal(t, retainedQuote, event.RetainedQuote) && assert.Equal(t, quote.AcceptedPeginQuoteEventId, event.Event.Id())
	})).Once()
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return()
	mutex.On("Unlock").Return()
	rsk := new(mocks.RootstockRpcServerMock)
	rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(50), nil)

	contracts := blockchain.RskContracts{Bridge: bridge}
	useCase := pegin.NewAcceptQuoteUseCase(quoteRepository, contracts, blockchain.Rpc{Rsk: rsk, Btc: btc}, lp, lp, eventBus, mutex)
	result, err := useCase.Run(context.Background(), acceptPeginQuoteHash)

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

func TestAcceptQuoteUseCase_Run_AlreadyAccepted(t *testing.T) {
	retainedQuote := quote.RetainedPeginQuote{
		QuoteHash:         acceptPeginQuoteHash,
		DepositAddress:    acceptPeginDerivationAddress,
		Signature:         acceptPeginSignature,
		RequiredLiquidity: entities.NewWei(9280000),
		State:             quote.PeginStateWaitingForDeposit,
	}
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), acceptPeginQuoteHash).Return(&testPeginQuote, nil)
	quoteRepository.On("GetRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), acceptPeginQuoteHash).Return(&retainedQuote, nil)

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
	useCase := pegin.NewAcceptQuoteUseCase(quoteRepository, contracts, rpc, lp, lp, eventBus, mutex)
	result, err := useCase.Run(context.Background(), acceptPeginQuoteHash)

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
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), acceptPeginQuoteHash).Return(nil, nil)

	bridge := new(mocks.BridgeMock)
	btc := new(mocks.BtcRpcMock)
	lp := new(mocks.ProviderMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	rsk := new(mocks.RootstockRpcServerMock)

	contracts := blockchain.RskContracts{Bridge: bridge}
	rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
	useCase := pegin.NewAcceptQuoteUseCase(quoteRepository, contracts, rpc, lp, lp, eventBus, mutex)
	result, err := useCase.Run(context.Background(), acceptPeginQuoteHash)

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
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), acceptPeginQuoteHash).Return(&expiredQuote, nil)

	bridge := new(mocks.BridgeMock)
	btc := new(mocks.BtcRpcMock)
	lp := new(mocks.ProviderMock)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	rsk := new(mocks.RootstockRpcServerMock)

	contracts := blockchain.RskContracts{Bridge: bridge}
	rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
	useCase := pegin.NewAcceptQuoteUseCase(quoteRepository, contracts, rpc, lp, lp, eventBus, mutex)
	result, err := useCase.Run(context.Background(), acceptPeginQuoteHash)

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
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), acceptPeginQuoteHash).Return(&testPeginQuote, nil)
	quoteRepository.On("GetRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), acceptPeginQuoteHash).Return(nil, nil)
	bridge := new(mocks.BridgeMock)
	bridge.On("FetchFederationInfo").Return(federationInfo, nil)
	bridge.On("GetFlyoverDerivationAddress", mock.Anything).Return(blockchain.FlyoverDerivation{
		Address:      "derivation address",
		RedeemScript: "any script",
	}, nil)
	btc := new(mocks.BtcRpcMock)
	btc.On("DecodeAddress", testPeginQuote.BtcRefundAddress, true).Return([]byte{4, 5, 6}, nil)
	btc.On("DecodeAddress", testPeginQuote.LpBtcAddress, true).Return([]byte{7, 8, 9}, nil)
	lp := new(mocks.ProviderMock)
	lp.On("HasPeginLiquidity", mock.AnythingOfType("context.backgroundCtx"), requiredLiquidity).Return(assert.AnError)
	eventBus := new(mocks.EventBusMock)
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return()
	mutex.On("Unlock").Return()
	rsk := new(mocks.RootstockRpcServerMock)
	rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(50), nil)

	contracts := blockchain.RskContracts{Bridge: bridge}
	rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
	useCase := pegin.NewAcceptQuoteUseCase(quoteRepository, contracts, rpc, lp, lp, eventBus, mutex)
	result, err := useCase.Run(context.Background(), acceptPeginQuoteHash)

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
		useCase := pegin.NewAcceptQuoteUseCase(quoteRepository, contracts, rpc, lp, lp, eventBus, mutex)
		result, err := useCase.Run(context.Background(), caseHash)

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
	derivation := blockchain.FlyoverDerivation{Address: test.AnyAddress, RedeemScript: "any script"}
	return []func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
		btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock){
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil, nil).Once()
			*quoteHash = "malformed hash"
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil, nil).Once()
			btc.On("DecodeAddress", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil, nil).Once()
			btc.On("DecodeAddress", mock.Anything, mock.Anything).Return([]byte{1}, nil).Once()
			btc.On("DecodeAddress", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			caseQuote := testPeginQuote
			caseQuote.LbcAddress = "malformed address"
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(&caseQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil, nil).Once()
			btc.On("DecodeAddress", mock.Anything, mock.Anything).Return([]byte{1}, nil).Once()
			btc.On("DecodeAddress", mock.Anything, mock.Anything).Return([]byte{2}, nil).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil, nil).Once()
			btc.On("DecodeAddress", mock.Anything, mock.Anything).Return([]byte{1}, nil).Once()
			btc.On("DecodeAddress", mock.Anything, mock.Anything).Return([]byte{2}, nil).Once()
			bridge.On("FetchFederationInfo").Return(blockchain.FederationInfo{}, assert.AnError).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil, nil).Once()
			btc.On("DecodeAddress", mock.Anything, mock.Anything).Return([]byte{1}, nil).Once()
			btc.On("DecodeAddress", mock.Anything, mock.Anything).Return([]byte{2}, nil).Once()
			bridge.On("FetchFederationInfo").Return(federationInfo, nil).Once()
			bridge.On("GetFlyoverDerivationAddress", mock.Anything).Return(blockchain.FlyoverDerivation{}, assert.AnError).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil, nil).Once()
			btc.On("DecodeAddress", mock.Anything, mock.Anything).Return([]byte{1}, nil).Once()
			btc.On("DecodeAddress", mock.Anything, mock.Anything).Return([]byte{2}, nil).Once()
			bridge.On("FetchFederationInfo").Return(federationInfo, nil).Once()
			bridge.On("GetFlyoverDerivationAddress", mock.Anything).Return(derivation, nil).Once()
			rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(nil, assert.AnError).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil, nil).Once()
			btc.On("DecodeAddress", mock.Anything, mock.Anything).Return([]byte{1}, nil).Once()
			btc.On("DecodeAddress", mock.Anything, mock.Anything).Return([]byte{2}, nil).Once()
			bridge.On("FetchFederationInfo").Return(federationInfo, nil).Once()
			bridge.On("GetFlyoverDerivationAddress", mock.Anything).Return(derivation, nil).Once()
			rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(1), nil).Once()
			lp.On("HasPeginLiquidity", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil).Once()
			lp.On("SignQuote", mock.Anything).Return("", assert.AnError).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil, nil).Once()
			btc.On("DecodeAddress", mock.Anything, mock.Anything).Return([]byte{1}, nil).Once()
			btc.On("DecodeAddress", mock.Anything, mock.Anything).Return([]byte{2}, nil).Once()
			bridge.On("FetchFederationInfo").Return(federationInfo, nil).Once()
			rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(1), nil).Once()
			lp.On("HasPeginLiquidity", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil).Once()
			// set derivation and signature to empty to malform the retained quote
			bridge.On("GetFlyoverDerivationAddress", mock.Anything).Return(blockchain.FlyoverDerivation{}, nil).Once()
			lp.On("SignQuote", mock.Anything).Return("", nil).Once()
		},
		func(quoteHash *string, quoteRepository *mocks.PeginQuoteRepositoryMock, bridge *mocks.BridgeMock,
			btc *mocks.BtcRpcMock, lp *mocks.ProviderMock, rsk *mocks.RootstockRpcServerMock) {
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(&testPeginQuote, nil).Once()
			quoteRepository.On("GetRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil, nil).Once()
			btc.On("DecodeAddress", mock.Anything, mock.Anything).Return([]byte{1}, nil).Once()
			btc.On("DecodeAddress", mock.Anything, mock.Anything).Return([]byte{2}, nil).Once()
			bridge.On("FetchFederationInfo").Return(federationInfo, nil).Once()
			bridge.On("GetFlyoverDerivationAddress", mock.Anything).Return(derivation, nil).Once()
			rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(1), nil).Once()
			lp.On("HasPeginLiquidity", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(nil).Once()
			lp.On("SignQuote", mock.Anything).Return("signature", nil).Once()
			quoteRepository.On("InsertRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(assert.AnError).Once()
		},
	}
}
