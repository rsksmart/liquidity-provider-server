package liquidity_provider_test

import (
	"context"
	"testing"

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
)

var peginConfigMock = entities.Signed[lp.PeginConfiguration]{
	Value: lp.PeginConfiguration{
		TimeForDeposit: 1,
		CallTime:       2,
		PenaltyFee:     entities.NewWei(3),
		FixedFee:       entities.NewWei(4),
		FeePercentage:  utils.NewBigFloat64(4.5),
		MaxValue:       entities.NewWei(5),
		MinValue:       entities.NewWei(1),
	},
	Signature: "010203",
	Hash:      "040506",
}

func TestSetPeginConfigUseCase_Run(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	lpRepository.On("UpsertPeginConfiguration", test.AnyCtx, peginConfigMock).Return(nil)
	walletMock := &mocks.RskWalletMock{}
	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
	hashMock := &mocks.HashMock{}
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})

	bridge := &mocks.BridgeMock{}
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1), nil)
	contracts := blockchain.RskContracts{Bridge: bridge}

	useCase := liquidity_provider.NewSetPeginConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)

	err := useCase.Run(context.Background(), peginConfigMock.Value)
	require.NoError(t, err)
	lpRepository.AssertExpectations(t)
	walletMock.AssertExpectations(t)
	hashMock.AssertExpectations(t)
	bridge.AssertExpectations(t)
}

func TestSetPeginConfigUseCase_Run_ErrorHandling(t *testing.T) {
	hashMock := &mocks.HashMock{}
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})

	errorSetups := []func(lpRepository *mocks.LiquidityProviderRepositoryMock, walletMock *mocks.RskWalletMock){
		func(lpRepository *mocks.LiquidityProviderRepositoryMock, walletMock *mocks.RskWalletMock) {
			walletMock.On("SignBytes", mock.Anything).Return(nil, assert.AnError)
		},
		func(lpRepository *mocks.LiquidityProviderRepositoryMock, walletMock *mocks.RskWalletMock) {
			walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
			lpRepository.On("UpsertPeginConfiguration", test.AnyCtx, peginConfigMock).Return(assert.AnError)
		},
	}

	for _, errorSetup := range errorSetups {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		walletMock := &mocks.RskWalletMock{}
		errorSetup(lpRepository, walletMock)
		bridge := &mocks.BridgeMock{}
		bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1), nil)
		contracts := blockchain.RskContracts{Bridge: bridge}
		useCase := liquidity_provider.NewSetPeginConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)
		err := useCase.Run(context.Background(), peginConfigMock.Value)
		require.Error(t, err)
		lpRepository.AssertExpectations(t)
		walletMock.AssertExpectations(t)
		bridge.AssertExpectations(t)
	}
}

func TestSetPeginConfigUseCase_Run_ValidateBridgeMin(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	bridge := &mocks.BridgeMock{}
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(10), nil)
	contracts := blockchain.RskContracts{Bridge: bridge}

	useCase := liquidity_provider.NewSetPeginConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)

	err := useCase.Run(context.Background(), peginConfigMock.Value)
	require.ErrorIs(t, err, usecases.TxBelowMinimumError)
	bridge.AssertExpectations(t)
}

func TestSetPeginConfigUseCase_Run_ValidatePositiveWei(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	invalidConfig := lp.PeginConfiguration{
		TimeForDeposit: 1,
		CallTime:       2,
		PenaltyFee:     entities.NewWei(3),
		FixedFee:       entities.NewWei(-1),
		FeePercentage:  utils.NewBigFloat64(4.5),
		MaxValue:       entities.NewWei(5),
		MinValue:       entities.NewWei(1),
	}

	bridge := &mocks.BridgeMock{}
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1), nil)
	contracts := blockchain.RskContracts{Bridge: bridge}

	useCase := liquidity_provider.NewSetPeginConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)

	err := useCase.Run(context.Background(), invalidConfig)
	require.ErrorIs(t, err, usecases.NonPositiveWeiError)
}

func TestSetPeginConfigUseCase_Run_ZeroFixedFee(t *testing.T) {
    lpRepository := &mocks.LiquidityProviderRepositoryMock{}
    walletMock := &mocks.RskWalletMock{}
    hashMock := &mocks.HashMock{}

    cfg := lp.PeginConfiguration{
        TimeForDeposit: 1,
        CallTime:       2,
        PenaltyFee:     entities.NewWei(3),
        FixedFee:       entities.NewWei(0),
        FeePercentage:  utils.NewBigFloat64(4.5),
        MaxValue:       entities.NewWei(5),
        MinValue:       entities.NewWei(1),
    }

    bridge := &mocks.BridgeMock{}
    bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1), nil)
    contracts := blockchain.RskContracts{Bridge: bridge}

    lpRepository.On("UpsertPeginConfiguration", test.AnyCtx, mock.Anything).Return(nil)
    walletMock.On("SignBytes", mock.Anything).Return([]byte{1,2,3}, nil)
    hashMock.On("Hash", mock.Anything).Return([]byte{4,5,6})

    useCase := liquidity_provider.NewSetPeginConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)

    err := useCase.Run(context.Background(), cfg)
    require.NoError(t, err)
    bridge.AssertExpectations(t)
    lpRepository.AssertExpectations(t)
}
