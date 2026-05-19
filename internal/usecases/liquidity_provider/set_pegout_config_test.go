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

func signedGeneralConfigMock(rskConfirmations, btcConfirmations uint16) *entities.Signed[lp.GeneralConfiguration] {
	return &entities.Signed[lp.GeneralConfiguration]{
		Value: lp.GeneralConfiguration{
			RskConfirmations: lp.ConfirmationsPerAmount{
				"1": rskConfirmations,
			},
			BtcConfirmations: lp.ConfirmationsPerAmount{
				"1": btcConfirmations,
			},
		},
	}
}

func TestSetPegoutConfigUseCase_Run(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	lpRepository.On("GetGeneralConfiguration", test.AnyCtx).Return(signedGeneralConfigMock(0, 0), nil)
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
		lpRepository.On("GetGeneralConfiguration", test.AnyCtx).Return(signedGeneralConfigMock(0, 0), nil)
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
		bridge.AssertExpectations(t)
	}
}

func TestSetPegoutConfigUseCase_Run_ValidatePositiveWei(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	invalidConfig := lp.PegoutConfiguration{
		TimeForDeposit:       1,
		ExpireTime:           2,
		PenaltyFee:           entities.NewWei(-3),
		FixedFee:             entities.NewWei(4),
		FeePercentage:        utils.NewBigFloat64(4.5),
		MaxValue:             entities.NewWei(5),
		MinValue:             entities.NewWei(1),
		ExpireBlocks:         10,
		BridgeTransactionMin: entities.NewWei(5),
	}

	bridge := &mocks.BridgeMock{}
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1), nil)
	contracts := blockchain.RskContracts{Bridge: bridge}

	useCase := liquidity_provider.NewSetPegoutConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)

	err := useCase.Run(context.Background(), invalidConfig)
	require.ErrorIs(t, err, usecases.NonPositiveWeiError)
}

func TestSetPegoutConfigUseCase_Run_ZeroFixedFee(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	cfg := lp.PegoutConfiguration{
		TimeForDeposit:       1,
		ExpireTime:           2,
		PenaltyFee:           entities.NewWei(3),
		FixedFee:             entities.NewWei(0),
		FeePercentage:        utils.NewBigFloat64(4.5),
		MaxValue:             entities.NewWei(5),
		MinValue:             entities.NewWei(1),
		ExpireBlocks:         10,
		BridgeTransactionMin: entities.NewWei(5),
	}

	bridge := &mocks.BridgeMock{}
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1), nil)
	contracts := blockchain.RskContracts{Bridge: bridge}

	lpRepository.On("GetGeneralConfiguration", test.AnyCtx).Return(signedGeneralConfigMock(0, 0), nil)
	lpRepository.On("UpsertPegoutConfiguration", test.AnyCtx, mock.Anything).Return(nil)
	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})

	useCase := liquidity_provider.NewSetPegoutConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)

	err := useCase.Run(context.Background(), cfg)
	require.NoError(t, err)
	bridge.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
}

func TestSetPegoutConfigUseCase_Run_ValidateExpiryAgainstConfirmations(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	cfg := lp.PegoutConfiguration{
		TimeForDeposit:       1,
		ExpireTime:           10,
		PenaltyFee:           entities.NewWei(3),
		FixedFee:             entities.NewWei(4),
		FeePercentage:        utils.NewBigFloat64(4.5),
		MaxValue:             entities.NewWei(5),
		MinValue:             entities.NewWei(1),
		ExpireBlocks:         10,
		BridgeTransactionMin: entities.NewWei(5),
	}

	lpRepository.On("GetGeneralConfiguration", test.AnyCtx).Return(signedGeneralConfigMock(1, 1), nil)

	bridge := &mocks.BridgeMock{}
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1), nil)
	contracts := blockchain.RskContracts{Bridge: bridge}

	useCase := liquidity_provider.NewSetPegoutConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)
	err := useCase.Run(context.Background(), cfg)

	require.ErrorIs(t, err, lp.PegoutExpiryTooShortForConfirmationsError)
	lpRepository.AssertNotCalled(t, "UpsertPegoutConfiguration", mock.Anything, mock.Anything)
	walletMock.AssertNotCalled(t, "SignBytes", mock.Anything)
}

func TestSetPegoutConfigUseCase_Run_DefaultGeneralConfigFallback(t *testing.T) {
	baseCfg := lp.PegoutConfiguration{
		TimeForDeposit:       1,
		ExpireTime:           10,
		PenaltyFee:           entities.NewWei(3),
		FixedFee:             entities.NewWei(4),
		FeePercentage:        utils.NewBigFloat64(4.5),
		MaxValue:             entities.NewWei(5),
		MinValue:             entities.NewWei(1),
		ExpireBlocks:         10,
		BridgeTransactionMin: entities.NewWei(5),
	}

	t.Run("should fallback when general config does not exist", func(t *testing.T) {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		lpRepository.On("GetGeneralConfiguration", test.AnyCtx).Return(nil, nil)
		walletMock := &mocks.RskWalletMock{}
		hashMock := &mocks.HashMock{}
		bridge := &mocks.BridgeMock{}
		bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1), nil)
		contracts := blockchain.RskContracts{Bridge: bridge}

		useCase := liquidity_provider.NewSetPegoutConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)
		err := useCase.Run(context.Background(), baseCfg)

		require.ErrorIs(t, err, lp.PegoutExpiryTooShortForConfirmationsError)
		lpRepository.AssertNotCalled(t, "UpsertPegoutConfiguration", mock.Anything, mock.Anything)
	})

	t.Run("should fallback when reading general config fails", func(t *testing.T) {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		lpRepository.On("GetGeneralConfiguration", test.AnyCtx).Return(nil, assert.AnError)
		walletMock := &mocks.RskWalletMock{}
		hashMock := &mocks.HashMock{}
		bridge := &mocks.BridgeMock{}
		bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1), nil)
		contracts := blockchain.RskContracts{Bridge: bridge}

		useCase := liquidity_provider.NewSetPegoutConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)
		err := useCase.Run(context.Background(), baseCfg)

		require.ErrorIs(t, err, lp.PegoutExpiryTooShortForConfirmationsError)
		lpRepository.AssertNotCalled(t, "UpsertPegoutConfiguration", mock.Anything, mock.Anything)
	})
}

func TestSetPeginConfigUseCase_Run_ValidatePositiveWei_EachField(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	baseCfg := lp.PeginConfiguration{
		TimeForDeposit: 1,
		CallTime:       2,
		PenaltyFee:     entities.NewWei(3),
		FixedFee:       entities.NewWei(4),
		FeePercentage:  utils.NewBigFloat64(4.5),
		MaxValue:       entities.NewWei(5),
		MinValue:       entities.NewWei(1),
	}

	makeCfg := func(modify func(*lp.PeginConfiguration)) lp.PeginConfiguration {
		cfg := baseCfg
		modify(&cfg)
		return cfg
	}

	cases := []lp.PeginConfiguration{
		makeCfg(func(c *lp.PeginConfiguration) { c.PenaltyFee = entities.NewWei(-1) }),
		makeCfg(func(c *lp.PeginConfiguration) { c.FixedFee = entities.NewWei(-1) }),
		makeCfg(func(c *lp.PeginConfiguration) { c.MaxValue = entities.NewWei(-1) }),
		makeCfg(func(c *lp.PeginConfiguration) { c.MinValue = entities.NewWei(-1) }),
	}

	for _, cfg := range cases {
		bridge := &mocks.BridgeMock{}
		bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1), nil)
		contracts := blockchain.RskContracts{Bridge: bridge}

		useCase := liquidity_provider.NewSetPeginConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)
		err := useCase.Run(context.Background(), cfg)
		require.ErrorIs(t, err, usecases.NonPositiveWeiError)
		bridge.AssertNotCalled(t, "GetMinimumLockTxValue")
	}
}
