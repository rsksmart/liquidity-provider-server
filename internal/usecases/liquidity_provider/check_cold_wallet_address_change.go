package liquidity_provider

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/alerts"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/cold_wallet"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

// stateAndCurrentHashes holds the loaded state config and current cold wallet address hashes.
type stateAndCurrentHashes struct {
	StateConfig liquidity_provider.StateConfiguration
	NewBtcHash  string
	NewRskHash  string
}

type CheckColdWalletAddressChangeUseCase struct {
	lpRepository    liquidity_provider.LiquidityProviderRepository
	generalProvider liquidity_provider.LiquidityProvider
	coldWallet      cold_wallet.ColdWallet
	alertSender     alerts.AlertSender
	alertRecipient  string
	signer          entities.Signer
	hashFunc        entities.HashFunction
}

func NewCheckColdWalletAddressChangeUseCase(
	lpRepository liquidity_provider.LiquidityProviderRepository,
	generalProvider liquidity_provider.LiquidityProvider,
	coldWallet cold_wallet.ColdWallet,
	alertSender alerts.AlertSender,
	alertRecipient string,
	signer entities.Signer,
	hashFunc entities.HashFunction,
) *CheckColdWalletAddressChangeUseCase {
	return &CheckColdWalletAddressChangeUseCase{
		lpRepository:    lpRepository,
		generalProvider: generalProvider,
		coldWallet:      coldWallet,
		alertSender:     alertSender,
		alertRecipient:  alertRecipient,
		signer:          signer,
		hashFunc:        hashFunc,
	}
}

func (useCase *CheckColdWalletAddressChangeUseCase) hashAddress(address string) string {
	return hex.EncodeToString(useCase.hashFunc([]byte(address)))
}

func (useCase *CheckColdWalletAddressChangeUseCase) Run(ctx context.Context) error {
	hashes, err := useCase.loadStateAndCurrentHashes(ctx)
	if err != nil {
		return err
	}

	storedBtcHash := hashes.StateConfig.BtcColdWalletAddressHash
	storedRskHash := hashes.StateConfig.RskColdWalletAddressHash

	firstRun := storedBtcHash == "" && storedRskHash == ""
	btcChanged := storedBtcHash != "" && storedBtcHash != hashes.NewBtcHash
	rskChanged := storedRskHash != "" && storedRskHash != hashes.NewRskHash

	if firstRun {
		return useCase.handleFirstRun(ctx, hashes)
	}
	if !btcChanged && !rskChanged {
		return nil
	}
	return useCase.handleAddressChange(ctx, hashes, btcChanged, rskChanged)
}

// loadStateAndCurrentHashes reads the state config through the LiquidityProvider interface
// (which validates signatures) and computes current cold wallet address hashes, or returns an error.
func (useCase *CheckColdWalletAddressChangeUseCase) loadStateAndCurrentHashes(ctx context.Context) (stateAndCurrentHashes, error) {
	stateConfig, err := useCase.generalProvider.StateConfiguration(ctx)
	if err != nil {
		return stateAndCurrentHashes{}, usecases.WrapUseCaseError(usecases.CheckColdWalletAddressChangeId, err)
	}

	btcAddr := useCase.coldWallet.GetBtcAddress()
	rskAddr := useCase.coldWallet.GetRskAddress()
	if btcAddr == "" {
		return stateAndCurrentHashes{}, usecases.WrapUseCaseError(usecases.CheckColdWalletAddressChangeId, errors.New("cold wallet BTC address not configured"))
	}
	if rskAddr == "" {
		return stateAndCurrentHashes{}, usecases.WrapUseCaseError(usecases.CheckColdWalletAddressChangeId, errors.New("cold wallet RSK address not configured"))
	}

	return stateAndCurrentHashes{
		StateConfig: stateConfig,
		NewBtcHash:  useCase.hashAddress(btcAddr),
		NewRskHash:  useCase.hashAddress(rskAddr),
	}, nil
}

func (useCase *CheckColdWalletAddressChangeUseCase) handleFirstRun(ctx context.Context, input stateAndCurrentHashes) error {
	log.Info("CheckColdWalletAddressChange: first run, persisting cold wallet address hashes (no alert)")
	input.StateConfig.BtcColdWalletAddressHash = input.NewBtcHash
	input.StateConfig.RskColdWalletAddressHash = input.NewRskHash
	return useCase.persistStateConfig(ctx, input.StateConfig)
}

func (useCase *CheckColdWalletAddressChangeUseCase) handleAddressChange(ctx context.Context, input stateAndCurrentHashes, btcChanged, rskChanged bool) error {
	if err := useCase.sendAddressChangeAlerts(ctx, btcChanged, rskChanged); err != nil {
		return err
	}
	input.StateConfig.BtcColdWalletAddressHash = input.NewBtcHash
	input.StateConfig.RskColdWalletAddressHash = input.NewRskHash
	return useCase.persistStateConfig(ctx, input.StateConfig)
}

func (useCase *CheckColdWalletAddressChangeUseCase) sendAddressChangeAlerts(ctx context.Context, btcChanged, rskChanged bool) error {
	const bodyPrefix = "Cold wallet address change detected at startup"
	if btcChanged {
		body := bodyPrefix + " | Network: BTC"
		log.Info("CheckColdWalletAddressChange: cold wallet address change detected at startup | Network: BTC")
		if err := useCase.alertSender.SendAlert(ctx, alerts.AlertSubjectColdWalletChange, body, []string{useCase.alertRecipient}); err != nil {
			log.Errorf("CheckColdWalletAddressChange: failed to send alert | Network: BTC | error: %v", err)
			return usecases.WrapUseCaseError(usecases.CheckColdWalletAddressChangeId, fmt.Errorf("send alert: %w", err))
		}
	}
	if rskChanged {
		body := bodyPrefix + " | Network: RSK"
		log.Info("CheckColdWalletAddressChange: cold wallet address change detected at startup | Network: RSK")
		if err := useCase.alertSender.SendAlert(ctx, alerts.AlertSubjectColdWalletChange, body, []string{useCase.alertRecipient}); err != nil {
			log.Errorf("CheckColdWalletAddressChange: failed to send alert | Network: RSK | error: %v", err)
			return usecases.WrapUseCaseError(usecases.CheckColdWalletAddressChangeId, fmt.Errorf("send alert: %w", err))
		}
	}
	return nil
}

func (useCase *CheckColdWalletAddressChangeUseCase) persistStateConfig(ctx context.Context, stateConfig liquidity_provider.StateConfiguration) error {
	signedConfig, err := usecases.SignConfiguration(usecases.CheckColdWalletAddressChangeId, useCase.signer, useCase.hashFunc, stateConfig)
	if err != nil {
		return err
	}
	if err := useCase.lpRepository.UpsertStateConfiguration(ctx, signedConfig); err != nil {
		return usecases.WrapUseCaseError(usecases.CheckColdWalletAddressChangeId, err)
	}
	return nil
}
