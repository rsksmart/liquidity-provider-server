package dataproviders

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type LocalLiquidityProvider struct {
	peginRepository  quote.PeginQuoteRepository
	pegoutRepository quote.PegoutQuoteRepository
	lpRepository     liquidity_provider.LiquidityProviderRepository
	rpc              blockchain.Rpc
	signer           rootstock.TransactionSigner
	btc              blockchain.BitcoinWallet
	contracts        blockchain.RskContracts
}

func NewLocalLiquidityProvider(
	peginRepository quote.PeginQuoteRepository,
	pegoutRepository quote.PegoutQuoteRepository,
	lpRepository liquidity_provider.LiquidityProviderRepository,
	rpc blockchain.Rpc,
	signer rootstock.TransactionSigner,
	btc blockchain.BitcoinWallet,
	contracts blockchain.RskContracts,
) *LocalLiquidityProvider {
	return &LocalLiquidityProvider{
		peginRepository:  peginRepository,
		pegoutRepository: pegoutRepository,
		lpRepository:     lpRepository,
		rpc:              rpc,
		signer:           signer,
		btc:              btc,
		contracts:        contracts,
	}
}

func (lp *LocalLiquidityProvider) RskAddress() string {
	return lp.signer.Address().String()
}

func (lp *LocalLiquidityProvider) BtcAddress() string {
	return lp.btc.Address()
}

func (lp *LocalLiquidityProvider) SignQuote(quoteHash string) (string, error) {
	var buf bytes.Buffer

	hash, err := hex.DecodeString(quoteHash)
	if err != nil {
		return "", err
	}

	buf.WriteString("\x19Ethereum Signed Message:\n32")
	buf.Write(hash)
	signatureBytes, err := lp.signer.SignBytes(crypto.Keccak256(buf.Bytes()))
	if err != nil {
		return "", err
	}
	signatureBytes[len(signatureBytes)-1] += 27 // v must be 27 or 28
	return hex.EncodeToString(signatureBytes), nil
}

func (lp *LocalLiquidityProvider) CalculateLockedPegoutLiquidity(ctx context.Context) (*entities.Wei, error) {
	lockedLiquidity := new(entities.Wei)
	quotes, err := lp.pegoutRepository.GetRetainedQuoteByState(ctx,
		quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations, quote.PegoutStateSendPegoutFailed,
	)
	if err != nil {
		return nil, err
	}
	for _, retainedQuote := range quotes {
		lockedLiquidity.Add(lockedLiquidity, retainedQuote.RequiredLiquidity)
	}
	return lockedLiquidity, nil
}

func (lp *LocalLiquidityProvider) CalculateAvailablePegoutLiquidity(ctx context.Context) (*entities.Wei, error) {
	liquidity, err := lp.btc.GetBalance()
	if err != nil {
		return nil, err
	}
	lockedLiquidity, err := lp.CalculateLockedPegoutLiquidity(ctx)
	if err != nil {
		return nil, err
	}
	availableLiquidity := new(entities.Wei).Sub(liquidity, lockedLiquidity)
	return availableLiquidity, nil
}

func (lp *LocalLiquidityProvider) HasPegoutLiquidity(ctx context.Context, requiredLiquidity *entities.Wei) error {
	log.Debug("Verifying if has pegout liquidity")
	availableLiquidity, err := lp.CalculateAvailablePegoutLiquidity(ctx)
	if err != nil {
		return err
	}
	if availableLiquidity.Cmp(requiredLiquidity) >= 0 {
		return nil
	} else {
		return fmt.Errorf("not enough liquidity, missing %s satoshi\n", requiredLiquidity.Sub(requiredLiquidity, availableLiquidity).ToSatoshi().String())
	}
}

func (lp *LocalLiquidityProvider) CalculateAvailablePeginLiquidity(ctx context.Context) (*entities.Wei, error) {
	liquidity := new(entities.Wei)
	lockedLiquidity := new(entities.Wei)
	lpRskBalance, err := lp.rpc.Rsk.GetBalance(ctx, lp.RskAddress())
	if err != nil {
		return nil, err
	}
	lpLbcBalance, err := lp.contracts.Lbc.GetBalance(lp.RskAddress())
	if err != nil {
		return nil, err
	}
	liquidity.Add(lpRskBalance, lpLbcBalance)
	quotes, err := lp.peginRepository.GetRetainedQuoteByState(ctx, quote.PeginStateWaitingForDeposit)
	if err != nil {
		return nil, err
	}
	for _, retainedQuote := range quotes {
		lockedLiquidity.Add(lockedLiquidity, retainedQuote.RequiredLiquidity)
	}
	availableLiquidity := new(entities.Wei).Sub(liquidity, lockedLiquidity)
	return availableLiquidity, nil
}

func (lp *LocalLiquidityProvider) HasPeginLiquidity(ctx context.Context, requiredLiquidity *entities.Wei) error {
	log.Debug("Verifying if has liquidity")
	availableLiquidity, err := lp.CalculateAvailablePeginLiquidity(ctx)
	if err != nil {
		return err
	}
	if availableLiquidity.Cmp(requiredLiquidity) >= 0 {
		return nil
	} else {
		return fmt.Errorf(
			"%w missing %s wei\n",
			usecases.NoLiquidityError,
			requiredLiquidity.Sub(requiredLiquidity, availableLiquidity).String(),
		)
	}
}

func (lp *LocalLiquidityProvider) GeneralConfiguration(ctx context.Context) liquidity_provider.GeneralConfiguration {
	configuration, err := validateConfiguration("general", lp, func() (*entities.Signed[liquidity_provider.GeneralConfiguration], error) {
		return lp.lpRepository.GetGeneralConfiguration(ctx)
	})
	if err != nil {
		return liquidity_provider.DefaultGeneralConfiguration()
	}
	return configuration.Value
}

func (lp *LocalLiquidityProvider) PegoutConfiguration(ctx context.Context) liquidity_provider.PegoutConfiguration {
	configuration, err := validateConfiguration("pegout", lp, func() (*entities.Signed[liquidity_provider.PegoutConfiguration], error) {
		return lp.lpRepository.GetPegoutConfiguration(ctx)
	})
	if err != nil {
		return liquidity_provider.DefaultPegoutConfiguration()
	}
	return configuration.Value
}

func (lp *LocalLiquidityProvider) PeginConfiguration(ctx context.Context) liquidity_provider.PeginConfiguration {
	configuration, err := validateConfiguration("pegin", lp, func() (*entities.Signed[liquidity_provider.PeginConfiguration], error) {
		return lp.lpRepository.GetPeginConfiguration(ctx)
	})
	if err != nil {
		return liquidity_provider.DefaultPeginConfiguration()
	}
	return configuration.Value
}

func validateConfiguration[T liquidity_provider.ConfigurationType](
	displayName string,
	lp *LocalLiquidityProvider,
	readFunction func() (*entities.Signed[T], error),
) (*entities.Signed[T], error) {
	configuration, err := readFunction()
	if err != nil {
		log.Errorf("Error getting %s configuration, using default configuration. Error: %v", displayName, err)
		return nil, err
	}
	if configuration == nil {
		log.Warnf("Custom %s configuration not found. Using default configuration.", displayName)
		return nil, errors.New("configuration not found")
	}
	if err = configuration.CheckIntegrity(crypto.Keccak256); err != nil {
		log.Errorf("Tampered %s configuration. Using default configuration. Error: %v", displayName, err)
		return nil, err
	}
	if !lp.signer.Validate(configuration.Signature, configuration.Hash) {
		log.Errorf("Invalid %s configuration signature. Using default configuration.", displayName)
		return nil, errors.New("invalid signature")
	}
	return configuration, nil
}
