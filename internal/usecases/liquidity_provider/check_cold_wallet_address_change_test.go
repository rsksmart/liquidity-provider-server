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

func hashAddressForTest(address string) string {
	h := crypto.Keccak256([]byte(address))
	return hex.EncodeToString(h)
}

func signedStateConfigWithHashes(btcHash, rskHash string) *entities.Signed[lpEntity.StateConfiguration] {
	now := time.Now().UTC().Unix()
	return &entities.Signed[lpEntity.StateConfiguration]{
		Value: lpEntity.StateConfiguration{
			LastBtcToColdWalletTransfer:       &now,
			LastRbtcToColdWalletTransfer:      &now,
			LastKnownBtcColdWalletAddressHash: btcHash,
			LastKnownRskColdWalletAddressHash: rskHash,
		},
		Hash:      "hash",
		Signature: "signature",
	}
}

func TestCheckColdWalletAddressChangeUseCase_Run_FirstRun_NoStoredHashes(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	coldWallet := new(mocks.ColdWalletMock)
	alertSender := new(mocks.AlertSenderMock)
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	stateConfig := signedStateConfigWithHashes("", "")
	lpRepository.On("GetStateConfiguration", test.AnyCtx).Return(stateConfig, nil).Once()

	coldWallet.On("GetBtcAddress").Return(test.AnyBtcAddress)
	coldWallet.On("GetRskAddress").Return(test.AnyRskAddress)

	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})

	var capturedConfig entities.Signed[lpEntity.StateConfiguration]
	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.MatchedBy(func(c entities.Signed[lpEntity.StateConfiguration]) bool {
		capturedConfig = c
		return c.Value.LastKnownBtcColdWalletAddressHash == hashAddressForTest(test.AnyBtcAddress) &&
			c.Value.LastKnownRskColdWalletAddressHash == hashAddressForTest(test.AnyRskAddress)
	})).Return(nil).Once()

	useCase := liquidity_provider.NewCheckColdWalletAddressChangeUseCase(
		lpRepository, coldWallet, alertSender, "recipient@test.com", walletMock, hashMock.Hash,
	)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	lpRepository.AssertExpectations(t)
	coldWallet.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	assert.Equal(t, hashAddressForTest(test.AnyBtcAddress), capturedConfig.Value.LastKnownBtcColdWalletAddressHash)
	assert.Equal(t, hashAddressForTest(test.AnyRskAddress), capturedConfig.Value.LastKnownRskColdWalletAddressHash)
}

func TestCheckColdWalletAddressChangeUseCase_Run_StoredHashesMatchCurrent_NoAlertNoWrite(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	coldWallet := new(mocks.ColdWalletMock)
	alertSender := new(mocks.AlertSenderMock)
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	btcHash := hashAddressForTest(test.AnyBtcAddress)
	rskHash := hashAddressForTest(test.AnyRskAddress)
	stateConfig := signedStateConfigWithHashes(btcHash, rskHash)
	lpRepository.On("GetStateConfiguration", test.AnyCtx).Return(stateConfig, nil).Once()

	coldWallet.On("GetBtcAddress").Return(test.AnyBtcAddress)
	coldWallet.On("GetRskAddress").Return(test.AnyRskAddress)

	useCase := liquidity_provider.NewCheckColdWalletAddressChangeUseCase(
		lpRepository, coldWallet, alertSender, "recipient@test.com", walletMock, hashMock.Hash,
	)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	lpRepository.AssertExpectations(t)
	coldWallet.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
}

func TestCheckColdWalletAddressChangeUseCase_Run_BtcHashDifferent_AlertSent_StateConfigUpdated(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	coldWallet := new(mocks.ColdWalletMock)
	alertSender := new(mocks.AlertSenderMock)
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	newBtcAddr := "newBtcAddress"
	currentRskHash := hashAddressForTest(test.AnyRskAddress)
	stateConfig := signedStateConfigWithHashes("oldBtcHash", currentRskHash)
	lpRepository.On("GetStateConfiguration", test.AnyCtx).Return(stateConfig, nil).Once()

	coldWallet.On("GetBtcAddress").Return(newBtcAddr)
	coldWallet.On("GetRskAddress").Return(test.AnyRskAddress)

	alertSender.On("SendAlert", test.AnyCtx, alerts.AlertSubjectColdWalletChange,
		"Cold wallet address change detected at startup | Network: BTC",
		[]string{"recipient@test.com"}).Return(nil).Once()

	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})

	var capturedConfig entities.Signed[lpEntity.StateConfiguration]
	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.MatchedBy(func(c entities.Signed[lpEntity.StateConfiguration]) bool {
		capturedConfig = c
		return c.Value.LastKnownBtcColdWalletAddressHash == hashAddressForTest(newBtcAddr) &&
			c.Value.LastKnownRskColdWalletAddressHash == currentRskHash
	})).Return(nil).Once()

	useCase := liquidity_provider.NewCheckColdWalletAddressChangeUseCase(
		lpRepository, coldWallet, alertSender, "recipient@test.com", walletMock, hashMock.Hash,
	)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	lpRepository.AssertExpectations(t)
	coldWallet.AssertExpectations(t)
	alertSender.AssertExpectations(t)
	assert.Equal(t, hashAddressForTest(newBtcAddr), capturedConfig.Value.LastKnownBtcColdWalletAddressHash)
}

func TestCheckColdWalletAddressChangeUseCase_Run_RskHashDifferent_AlertSent_StateConfigUpdated(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	coldWallet := new(mocks.ColdWalletMock)
	alertSender := new(mocks.AlertSenderMock)
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	newRskAddr := "0xNewRskAddress123456789012345678901234567890"
	currentBtcHash := hashAddressForTest(test.AnyBtcAddress)
	stateConfig := signedStateConfigWithHashes(currentBtcHash, "oldRskHash")
	lpRepository.On("GetStateConfiguration", test.AnyCtx).Return(stateConfig, nil).Once()

	coldWallet.On("GetBtcAddress").Return(test.AnyBtcAddress)
	coldWallet.On("GetRskAddress").Return(newRskAddr)

	alertSender.On("SendAlert", test.AnyCtx, alerts.AlertSubjectColdWalletChange,
		"Cold wallet address change detected at startup | Network: RSK",
		[]string{"recipient@test.com"}).Return(nil).Once()

	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})

	var capturedConfig entities.Signed[lpEntity.StateConfiguration]
	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.MatchedBy(func(c entities.Signed[lpEntity.StateConfiguration]) bool {
		capturedConfig = c
		return c.Value.LastKnownRskColdWalletAddressHash == hashAddressForTest(newRskAddr) &&
			c.Value.LastKnownBtcColdWalletAddressHash == currentBtcHash
	})).Return(nil).Once()

	useCase := liquidity_provider.NewCheckColdWalletAddressChangeUseCase(
		lpRepository, coldWallet, alertSender, "recipient@test.com", walletMock, hashMock.Hash,
	)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	lpRepository.AssertExpectations(t)
	coldWallet.AssertExpectations(t)
	alertSender.AssertExpectations(t)
	assert.Equal(t, hashAddressForTest(newRskAddr), capturedConfig.Value.LastKnownRskColdWalletAddressHash)
}

func TestCheckColdWalletAddressChangeUseCase_Run_BothHashesDifferent_TwoAlerts_StateConfigUpdated(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	coldWallet := new(mocks.ColdWalletMock)
	alertSender := new(mocks.AlertSenderMock)
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	newBtcAddr := "newBtcAddr"
	newRskAddr := "0xNewRskAddr123456789012345678901234567890"
	stateConfig := signedStateConfigWithHashes("oldBtcHash", "oldRskHash")
	lpRepository.On("GetStateConfiguration", test.AnyCtx).Return(stateConfig, nil).Once()

	coldWallet.On("GetBtcAddress").Return(newBtcAddr)
	coldWallet.On("GetRskAddress").Return(newRskAddr)

	alertSender.On("SendAlert", test.AnyCtx, alerts.AlertSubjectColdWalletChange,
		"Cold wallet address change detected at startup | Network: BTC",
		[]string{"recipient@test.com"}).Return(nil).Once()
	alertSender.On("SendAlert", test.AnyCtx, alerts.AlertSubjectColdWalletChange,
		"Cold wallet address change detected at startup | Network: RSK",
		[]string{"recipient@test.com"}).Return(nil).Once()

	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})

	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.MatchedBy(func(c entities.Signed[lpEntity.StateConfiguration]) bool {
		return c.Value.LastKnownBtcColdWalletAddressHash == hashAddressForTest(newBtcAddr) &&
			c.Value.LastKnownRskColdWalletAddressHash == hashAddressForTest(newRskAddr)
	})).Return(nil).Once()

	useCase := liquidity_provider.NewCheckColdWalletAddressChangeUseCase(
		lpRepository, coldWallet, alertSender, "recipient@test.com", walletMock, hashMock.Hash,
	)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	lpRepository.AssertExpectations(t)
	coldWallet.AssertExpectations(t)
	alertSender.AssertExpectations(t)
}

func TestCheckColdWalletAddressChangeUseCase_Run_BtcAddressEmpty_ReturnsError_NoAlertNoPersist(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	coldWallet := new(mocks.ColdWalletMock)
	alertSender := new(mocks.AlertSenderMock)
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	stateConfig := signedStateConfigWithHashes("", "")
	lpRepository.On("GetStateConfiguration", test.AnyCtx).Return(stateConfig, nil).Once()

	coldWallet.On("GetBtcAddress").Return("")
	coldWallet.On("GetRskAddress").Return(test.AnyRskAddress)

	useCase := liquidity_provider.NewCheckColdWalletAddressChangeUseCase(
		lpRepository, coldWallet, alertSender, "recipient@test.com", walletMock, hashMock.Hash,
	)
	err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cold wallet BTC address not configured")
	lpRepository.AssertExpectations(t)
	coldWallet.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
}

func TestCheckColdWalletAddressChangeUseCase_Run_RskAddressEmpty_ReturnsError_NoAlertNoPersist(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	coldWallet := new(mocks.ColdWalletMock)
	alertSender := new(mocks.AlertSenderMock)
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	stateConfig := signedStateConfigWithHashes("", "")
	lpRepository.On("GetStateConfiguration", test.AnyCtx).Return(stateConfig, nil).Once()

	coldWallet.On("GetBtcAddress").Return(test.AnyBtcAddress)
	coldWallet.On("GetRskAddress").Return("")

	useCase := liquidity_provider.NewCheckColdWalletAddressChangeUseCase(
		lpRepository, coldWallet, alertSender, "recipient@test.com", walletMock, hashMock.Hash,
	)
	err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cold wallet RSK address not configured")
	lpRepository.AssertExpectations(t)
	coldWallet.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
}
