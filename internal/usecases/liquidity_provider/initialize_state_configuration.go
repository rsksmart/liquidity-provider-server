package liquidity_provider

import (
	"context"
	"encoding/hex"
	"errors"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/cold_wallet"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type InitializeStateConfigurationUseCase struct {
	provider     liquidity_provider.LiquidityProvider
	lpRepository liquidity_provider.LiquidityProviderRepository
	coldWallet   cold_wallet.ColdWallet
	signer       entities.Signer
	hashFunc     entities.HashFunction
}

func NewInitializeStateConfigurationUseCase(
	provider liquidity_provider.LiquidityProvider,
	lpRepository liquidity_provider.LiquidityProviderRepository,
	coldWallet cold_wallet.ColdWallet,
	signer entities.Signer,
	hashFunc entities.HashFunction,
) *InitializeStateConfigurationUseCase {
	return &InitializeStateConfigurationUseCase{
		provider:     provider,
		lpRepository: lpRepository,
		coldWallet:   coldWallet,
		signer:       signer,
		hashFunc:     hashFunc,
	}
}

type coldWalletAddresses struct {
	Btc string
	Rsk string
}

func (useCase *InitializeStateConfigurationUseCase) Run(ctx context.Context) error {
	addresses, err := useCase.validateColdWalletAddresses()
	if err != nil {
		return err
	}

	stateConfig, err := useCase.provider.StateConfiguration(ctx)
	if err != nil && !errors.Is(err, liquidity_provider.ConfigurationNotFoundError) {
		return usecases.WrapUseCaseError(usecases.InitializeStateConfigurationId, err)
	}

	modified := false
	now := time.Now().UTC().Unix()

	if stateConfig.LastBtcToColdWalletTransfer == 0 {
		log.Info("Initializing LastBtcToColdWalletTransfer with current timestamp")
		stateConfig.LastBtcToColdWalletTransfer = now
		modified = true
	}

	if stateConfig.LastRbtcToColdWalletTransfer == 0 {
		log.Info("Initializing LastRbtcToColdWalletTransfer with current timestamp")
		stateConfig.LastRbtcToColdWalletTransfer = now
		modified = true
	}

	if stateConfig.BtcColdWalletAddressHash == "" {
		log.Info("Initializing BtcColdWalletAddressHash")
		stateConfig.BtcColdWalletAddressHash = useCase.hashAddress(addresses.Btc)
		modified = true
	}

	if stateConfig.RskColdWalletAddressHash == "" {
		log.Info("Initializing RskColdWalletAddressHash")
		stateConfig.RskColdWalletAddressHash = useCase.hashAddress(addresses.Rsk)
		modified = true
	}

	if !modified {
		log.Debug("State configuration already fully initialized")
		return nil
	}

	return useCase.signAndPersist(ctx, stateConfig)
}

func (useCase *InitializeStateConfigurationUseCase) validateColdWalletAddresses() (coldWalletAddresses, error) {
	btcAddr := useCase.coldWallet.GetBtcAddress()
	rskAddr := useCase.coldWallet.GetRskAddress()
	if btcAddr == "" {
		return coldWalletAddresses{}, usecases.WrapUseCaseError(usecases.InitializeStateConfigurationId, errors.New("cold wallet BTC address not configured"))
	}
	if rskAddr == "" {
		return coldWalletAddresses{}, usecases.WrapUseCaseError(usecases.InitializeStateConfigurationId, errors.New("cold wallet RSK address not configured"))
	}
	return coldWalletAddresses{Btc: btcAddr, Rsk: rskAddr}, nil
}

func (useCase *InitializeStateConfigurationUseCase) hashAddress(address string) string {
	return hex.EncodeToString(useCase.hashFunc([]byte(address)))
}

func (useCase *InitializeStateConfigurationUseCase) signAndPersist(ctx context.Context, stateConfig liquidity_provider.StateConfiguration) error {
	signedConfig, err := usecases.SignConfiguration(usecases.InitializeStateConfigurationId, useCase.signer, useCase.hashFunc, stateConfig)
	if err != nil {
		return err
	}
	if err := useCase.lpRepository.UpsertStateConfiguration(ctx, signedConfig); err != nil {
		return usecases.WrapUseCaseError(usecases.InitializeStateConfigurationId, err)
	}
	log.Info("State configuration initialized successfully")
	return nil
}
