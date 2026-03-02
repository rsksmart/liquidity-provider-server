package liquidity_provider_test

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/alerts"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var testHashFunction = crypto.Keccak256

func hashAddressForTest(address string) string {
	return hex.EncodeToString(testHashFunction([]byte(address)))
}

func stateConfigWithHashes(btcHash, rskHash string) lpEntity.StateConfiguration {
	now := time.Now().UTC().Unix()
	return lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  now,
		LastRbtcToColdWalletTransfer: now,
		BtcColdWalletAddressHash:     btcHash,
		RskColdWalletAddressHash:     rskHash,
	}
}

func TestCheckColdWalletAddressChangeUseCase_Run_FirstRun_NoStoredHashes(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	providerMock := new(mocks.ProviderMock)
	coldWallet := new(mocks.ColdWalletMock)
	alertSender := new(mocks.AlertSenderMock)
	walletMock := &mocks.RskWalletMock{}

	stateConfig := stateConfigWithHashes("", "")
	providerMock.On("StateConfiguration", test.AnyCtx).Return(stateConfig, nil).Once()

	coldWallet.On("GetBtcAddress").Return(test.AnyBtcAddress)
	coldWallet.On("GetRskAddress").Return(test.AnyRskAddress)

	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)

	var capturedConfig entities.Signed[lpEntity.StateConfiguration]
	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.MatchedBy(func(c entities.Signed[lpEntity.StateConfiguration]) bool {
		capturedConfig = c
		return c.Value.BtcColdWalletAddressHash == hashAddressForTest(test.AnyBtcAddress) &&
			c.Value.RskColdWalletAddressHash == hashAddressForTest(test.AnyRskAddress)
	})).Return(nil).Once()

	useCase := liquidity_provider.NewCheckColdWalletAddressChangeUseCase(
		lpRepository, providerMock, coldWallet, alertSender, "recipient@test.com", walletMock, testHashFunction,
	)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	providerMock.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	coldWallet.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	assert.Equal(t, hashAddressForTest(test.AnyBtcAddress), capturedConfig.Value.BtcColdWalletAddressHash)
	assert.Equal(t, hashAddressForTest(test.AnyRskAddress), capturedConfig.Value.RskColdWalletAddressHash)
}

func TestCheckColdWalletAddressChangeUseCase_Run_StoredHashesMatchCurrent_NoAlertNoWrite(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	providerMock := new(mocks.ProviderMock)
	coldWallet := new(mocks.ColdWalletMock)
	alertSender := new(mocks.AlertSenderMock)
	walletMock := &mocks.RskWalletMock{}

	btcHash := hashAddressForTest(test.AnyBtcAddress)
	rskHash := hashAddressForTest(test.AnyRskAddress)
	stateConfig := stateConfigWithHashes(btcHash, rskHash)
	providerMock.On("StateConfiguration", test.AnyCtx).Return(stateConfig, nil).Once()

	coldWallet.On("GetBtcAddress").Return(test.AnyBtcAddress)
	coldWallet.On("GetRskAddress").Return(test.AnyRskAddress)

	useCase := liquidity_provider.NewCheckColdWalletAddressChangeUseCase(
		lpRepository, providerMock, coldWallet, alertSender, "recipient@test.com", walletMock, testHashFunction,
	)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	providerMock.AssertExpectations(t)
	coldWallet.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
}

func TestCheckColdWalletAddressChangeUseCase_Run_BtcHashDifferent_AlertSent_StateConfigUpdated(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	providerMock := new(mocks.ProviderMock)
	coldWallet := new(mocks.ColdWalletMock)
	alertSender := new(mocks.AlertSenderMock)
	walletMock := &mocks.RskWalletMock{}

	newBtcAddr := "newBtcAddress"
	currentRskHash := hashAddressForTest(test.AnyRskAddress)
	stateConfig := stateConfigWithHashes("oldBtcHash", currentRskHash)
	providerMock.On("StateConfiguration", test.AnyCtx).Return(stateConfig, nil).Once()

	coldWallet.On("GetBtcAddress").Return(newBtcAddr)
	coldWallet.On("GetRskAddress").Return(test.AnyRskAddress)

	alertSender.On("SendAlert", test.AnyCtx, alerts.AlertSubjectColdWalletChange,
		"Cold wallet address change detected at startup | Network: BTC",
		[]string{"recipient@test.com"}).Return(nil).Once()

	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)

	var capturedConfig entities.Signed[lpEntity.StateConfiguration]
	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.MatchedBy(func(c entities.Signed[lpEntity.StateConfiguration]) bool {
		capturedConfig = c
		return c.Value.BtcColdWalletAddressHash == hashAddressForTest(newBtcAddr) &&
			c.Value.RskColdWalletAddressHash == currentRskHash
	})).Return(nil).Once()

	useCase := liquidity_provider.NewCheckColdWalletAddressChangeUseCase(
		lpRepository, providerMock, coldWallet, alertSender, "recipient@test.com", walletMock, testHashFunction,
	)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	providerMock.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	coldWallet.AssertExpectations(t)
	alertSender.AssertExpectations(t)
	assert.Equal(t, hashAddressForTest(newBtcAddr), capturedConfig.Value.BtcColdWalletAddressHash)
}

func TestCheckColdWalletAddressChangeUseCase_Run_RskHashDifferent_AlertSent_StateConfigUpdated(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	providerMock := new(mocks.ProviderMock)
	coldWallet := new(mocks.ColdWalletMock)
	alertSender := new(mocks.AlertSenderMock)
	walletMock := &mocks.RskWalletMock{}

	newRskAddr := "0xNewRskAddress123456789012345678901234567890"
	currentBtcHash := hashAddressForTest(test.AnyBtcAddress)
	stateConfig := stateConfigWithHashes(currentBtcHash, "oldRskHash")
	providerMock.On("StateConfiguration", test.AnyCtx).Return(stateConfig, nil).Once()

	coldWallet.On("GetBtcAddress").Return(test.AnyBtcAddress)
	coldWallet.On("GetRskAddress").Return(newRskAddr)

	alertSender.On("SendAlert", test.AnyCtx, alerts.AlertSubjectColdWalletChange,
		"Cold wallet address change detected at startup | Network: RSK",
		[]string{"recipient@test.com"}).Return(nil).Once()

	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)

	var capturedConfig entities.Signed[lpEntity.StateConfiguration]
	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.MatchedBy(func(c entities.Signed[lpEntity.StateConfiguration]) bool {
		capturedConfig = c
		return c.Value.RskColdWalletAddressHash == hashAddressForTest(newRskAddr) &&
			c.Value.BtcColdWalletAddressHash == currentBtcHash
	})).Return(nil).Once()

	useCase := liquidity_provider.NewCheckColdWalletAddressChangeUseCase(
		lpRepository, providerMock, coldWallet, alertSender, "recipient@test.com", walletMock, testHashFunction,
	)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	providerMock.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	coldWallet.AssertExpectations(t)
	alertSender.AssertExpectations(t)
	assert.Equal(t, hashAddressForTest(newRskAddr), capturedConfig.Value.RskColdWalletAddressHash)
}

func TestCheckColdWalletAddressChangeUseCase_Run_BothHashesDifferent_TwoAlerts_StateConfigUpdated(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	providerMock := new(mocks.ProviderMock)
	coldWallet := new(mocks.ColdWalletMock)
	alertSender := new(mocks.AlertSenderMock)
	walletMock := &mocks.RskWalletMock{}

	newBtcAddr := "newBtcAddr"
	newRskAddr := "0xNewRskAddr123456789012345678901234567890"
	stateConfig := stateConfigWithHashes("oldBtcHash", "oldRskHash")
	providerMock.On("StateConfiguration", test.AnyCtx).Return(stateConfig, nil).Once()

	coldWallet.On("GetBtcAddress").Return(newBtcAddr)
	coldWallet.On("GetRskAddress").Return(newRskAddr)

	alertSender.On("SendAlert", test.AnyCtx, alerts.AlertSubjectColdWalletChange,
		"Cold wallet address change detected at startup | Network: BTC",
		[]string{"recipient@test.com"}).Return(nil).Once()
	alertSender.On("SendAlert", test.AnyCtx, alerts.AlertSubjectColdWalletChange,
		"Cold wallet address change detected at startup | Network: RSK",
		[]string{"recipient@test.com"}).Return(nil).Once()

	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)

	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.MatchedBy(func(c entities.Signed[lpEntity.StateConfiguration]) bool {
		return c.Value.BtcColdWalletAddressHash == hashAddressForTest(newBtcAddr) &&
			c.Value.RskColdWalletAddressHash == hashAddressForTest(newRskAddr)
	})).Return(nil).Once()

	useCase := liquidity_provider.NewCheckColdWalletAddressChangeUseCase(
		lpRepository, providerMock, coldWallet, alertSender, "recipient@test.com", walletMock, testHashFunction,
	)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	providerMock.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	coldWallet.AssertExpectations(t)
	alertSender.AssertExpectations(t)
}

func TestCheckColdWalletAddressChangeUseCase_Run_StateConfigurationError_ReturnsError(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	providerMock := new(mocks.ProviderMock)
	coldWallet := new(mocks.ColdWalletMock)
	alertSender := new(mocks.AlertSenderMock)
	walletMock := &mocks.RskWalletMock{}

	providerMock.On("StateConfiguration", test.AnyCtx).Return(lpEntity.StateConfiguration{}, assert.AnError).Once()

	useCase := liquidity_provider.NewCheckColdWalletAddressChangeUseCase(
		lpRepository, providerMock, coldWallet, alertSender, "recipient@test.com", walletMock, testHashFunction,
	)
	err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "CheckColdWalletAddressChange")
	providerMock.AssertExpectations(t)
	coldWallet.AssertNotCalled(t, "GetBtcAddress")
	coldWallet.AssertNotCalled(t, "GetRskAddress")
	alertSender.AssertNotCalled(t, "SendAlert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
}

func TestCheckColdWalletAddressChangeUseCase_Run_SendAlertFails_ReturnsError_NoPersist(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	providerMock := new(mocks.ProviderMock)
	coldWallet := new(mocks.ColdWalletMock)
	alertSender := new(mocks.AlertSenderMock)
	walletMock := &mocks.RskWalletMock{}

	currentRskHash := hashAddressForTest(test.AnyRskAddress)
	stateConfig := stateConfigWithHashes("oldBtcHash", currentRskHash)
	providerMock.On("StateConfiguration", test.AnyCtx).Return(stateConfig, nil).Once()

	coldWallet.On("GetBtcAddress").Return("newBtcAddress")
	coldWallet.On("GetRskAddress").Return(test.AnyRskAddress)

	alertSender.On("SendAlert", test.AnyCtx, alerts.AlertSubjectColdWalletChange,
		"Cold wallet address change detected at startup | Network: BTC",
		[]string{"recipient@test.com"}).Return(assert.AnError).Once()

	useCase := liquidity_provider.NewCheckColdWalletAddressChangeUseCase(
		lpRepository, providerMock, coldWallet, alertSender, "recipient@test.com", walletMock, testHashFunction,
	)
	err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "send alert")
	providerMock.AssertExpectations(t)
	coldWallet.AssertExpectations(t)
	alertSender.AssertExpectations(t)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
}

func TestCheckColdWalletAddressChangeUseCase_Run_BtcAddressEmpty_ReturnsError_NoAlertNoPersist(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	providerMock := new(mocks.ProviderMock)
	coldWallet := new(mocks.ColdWalletMock)
	alertSender := new(mocks.AlertSenderMock)
	walletMock := &mocks.RskWalletMock{}

	stateConfig := stateConfigWithHashes("", "")
	providerMock.On("StateConfiguration", test.AnyCtx).Return(stateConfig, nil).Once()

	coldWallet.On("GetBtcAddress").Return("")
	coldWallet.On("GetRskAddress").Return(test.AnyRskAddress)

	useCase := liquidity_provider.NewCheckColdWalletAddressChangeUseCase(
		lpRepository, providerMock, coldWallet, alertSender, "recipient@test.com", walletMock, testHashFunction,
	)
	err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cold wallet BTC address not configured")
	providerMock.AssertExpectations(t)
	coldWallet.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
}

func TestCheckColdWalletAddressChangeUseCase_Run_RskAddressEmpty_ReturnsError_NoAlertNoPersist(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	providerMock := new(mocks.ProviderMock)
	coldWallet := new(mocks.ColdWalletMock)
	alertSender := new(mocks.AlertSenderMock)
	walletMock := &mocks.RskWalletMock{}

	stateConfig := stateConfigWithHashes("", "")
	providerMock.On("StateConfiguration", test.AnyCtx).Return(stateConfig, nil).Once()

	coldWallet.On("GetBtcAddress").Return(test.AnyBtcAddress)
	coldWallet.On("GetRskAddress").Return("")

	useCase := liquidity_provider.NewCheckColdWalletAddressChangeUseCase(
		lpRepository, providerMock, coldWallet, alertSender, "recipient@test.com", walletMock, testHashFunction,
	)
	err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cold wallet RSK address not configured")
	providerMock.AssertExpectations(t)
	coldWallet.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
}
