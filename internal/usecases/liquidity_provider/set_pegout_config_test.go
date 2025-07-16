package liquidity_provider_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

var pegoutConfigMock = entities.Signed[lp.PegoutConfiguration]{
	Value: lp.PegoutConfiguration{
		TimeForDeposit:       1,
		ExpireTime:           2,
		PenaltyFee:           entities.NewWei(3),
		FixedFee:             entities.NewWei(4),
		FeePercentage:        utils.NewBigFloat64(4.5),
		MaxValue:             entities.NewWei(5),
		MinValue:             entities.NewWei(1),
		ExpireBlocks:         10,
		BridgeTransactionMin: entities.NewWei(5),
	},
	Signature: "010203",
	Hash:      "040506",
}

func TestSetPegoutConfigUseCase_Run(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	lpRepository.On("UpsertPegoutConfiguration", test.AnyCtx, pegoutConfigMock).Return(nil)
	walletMock := &mocks.RskWalletMock{}
	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
	hashMock := &mocks.HashMock{}
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})
	bridge := &mocks.BridgeMock{}
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1), nil)
	contracts := blockchain.RskContracts{Bridge: bridge}

	useCase := liquidity_provider.NewSetPegoutConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)

	err := useCase.Run(context.Background(), pegoutConfigMock.Value)
	require.NoError(t, err)
	lpRepository.AssertExpectations(t)
	walletMock.AssertExpectations(t)
	hashMock.AssertExpectations(t)
	bridge.AssertExpectations(t)
}

func TestSetPegoutConfigUseCase_Run_ValidateBridgeMin(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}
	bridge := &mocks.BridgeMock{}
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(10), nil)
	contracts := blockchain.RskContracts{Bridge: bridge}

	useCase := liquidity_provider.NewSetPegoutConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)

	err := useCase.Run(context.Background(), pegoutConfigMock.Value)
	require.ErrorIs(t, err, usecases.TxBelowMinimumError)
	lpRepository.AssertExpectations(t)
	walletMock.AssertExpectations(t)
	hashMock.AssertExpectations(t)
	bridge.AssertExpectations(t)
}

func TestSetPegoutConfigUseCase_Run_ErrorHandling(t *testing.T) {
	hashMock := &mocks.HashMock{}
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})

	errorSetups := []func(lpRepository *mocks.LiquidityProviderRepositoryMock, walletMock *mocks.RskWalletMock){
		func(lpRepository *mocks.LiquidityProviderRepositoryMock, walletMock *mocks.RskWalletMock) {
			walletMock.On("SignBytes", mock.Anything).Return(nil, assert.AnError)
		},
		func(lpRepository *mocks.LiquidityProviderRepositoryMock, walletMock *mocks.RskWalletMock) {
			walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
			lpRepository.On("UpsertPegoutConfiguration", test.AnyCtx, pegoutConfigMock).Return(assert.AnError)
		},
	}

	for _, errorSetup := range errorSetups {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		walletMock := &mocks.RskWalletMock{}
		errorSetup(lpRepository, walletMock)
		bridge := &mocks.BridgeMock{}
		bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1), nil)
		contracts := blockchain.RskContracts{Bridge: bridge}
		useCase := liquidity_provider.NewSetPegoutConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)
		err := useCase.Run(context.Background(), pegoutConfigMock.Value)
		require.Error(t, err)
		lpRepository.AssertExpectations(t)
		walletMock.AssertExpectations(t)
	}
}
