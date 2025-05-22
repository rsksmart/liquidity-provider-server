package dataproviders

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

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
	return strings.ToLower(lp.signer.Address().String())
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

	buf.WriteString(usecases.EthereumSignedMessagePrefix)
	buf.Write(hash)
	signatureBytes, err := lp.signer.SignBytes(crypto.Keccak256(buf.Bytes()))
	if err != nil {
		return "", err
	}
	signatureBytes[len(signatureBytes)-1] += 27 // v must be 27 or 28
	return hex.EncodeToString(signatureBytes), nil
}

func (lp *LocalLiquidityProvider) HasPegoutLiquidity(ctx context.Context, requiredLiquidity *entities.Wei) error {
	log.Debug("Verifying if has liquidity")
	availableLiquidity, err := lp.AvailablePegoutLiquidity(ctx)
	if err != nil {
		return err
	}
	if availableLiquidity.Cmp(requiredLiquidity) >= 0 {
		return nil
	} else {
		return fmt.Errorf(
			"not enough liquidity, missing %s satoshi",
			requiredLiquidity.Sub(requiredLiquidity, availableLiquidity).ToSatoshi().String(),
		)
	}
}

func (lp *LocalLiquidityProvider) AvailablePegoutLiquidity(ctx context.Context) (*entities.Wei, error) {
	lockedLiquidity := new(entities.Wei)
	liquidity, err := lp.btc.GetBalance()
	if err != nil {
		return nil, err
	}
	log.Debugf("Liquidity: %s satoshi", liquidity.ToSatoshi().String())
	quotes, err := lp.pegoutRepository.GetRetainedQuoteByState(ctx,
		quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations,
	)
	if err != nil {
		return nil, err
	}
	for _, retainedQuote := range quotes {
		lockedLiquidity.Add(lockedLiquidity, retainedQuote.RequiredLiquidity)
	}
	log.Debugf("Locked Liquidity: %s satoshi", lockedLiquidity.ToSatoshi().String())
	availableLiquidity := new(entities.Wei).Sub(liquidity, lockedLiquidity)
	return availableLiquidity, nil
}

func (lp *LocalLiquidityProvider) HasPeginLiquidity(ctx context.Context, requiredLiquidity *entities.Wei) error {
	log.Debug("Verifying if has liquidity")
	availableLiquidity, err := lp.AvailablePeginLiquidity(ctx)
	if err != nil {
		return err
	}
	if availableLiquidity.Cmp(requiredLiquidity) >= 0 {
		return nil
	} else {
		return fmt.Errorf(
			"%w missing %s wei",
			usecases.NoLiquidityError,
			requiredLiquidity.Sub(requiredLiquidity, availableLiquidity).String(),
		)
	}
}

func (lp *LocalLiquidityProvider) AvailablePeginLiquidity(ctx context.Context) (*entities.Wei, error) {
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
	log.Debugf("Liquidity: %s wei", liquidity.String())
	peginQuotes, err := lp.peginRepository.GetRetainedQuoteByState(ctx, quote.PeginStateWaitingForDeposit)
	if err != nil {
		return nil, err
	}
	for _, retainedQuote := range peginQuotes {
		lockedLiquidity.Add(lockedLiquidity, retainedQuote.RequiredLiquidity)
	}
	// we include this in the locked liquidity because the refund is done in RBTC, and it is converted to BTC once a threshold is reached
	pegoutQuotes, err := lp.pegoutRepository.GetRetainedQuoteByState(ctx, quote.PegoutStateRefundPegOutSucceeded)
	if err != nil {
		return nil, err
	}
	for _, retainedQuote := range pegoutQuotes {
		lockedLiquidity.Add(lockedLiquidity, retainedQuote.RequiredLiquidity)
	}
	log.Debugf("Locked Liquidity: %s wei", lockedLiquidity.String())
	return new(entities.Wei).Sub(liquidity, lockedLiquidity), nil
}

func (lp *LocalLiquidityProvider) GeneralConfiguration(ctx context.Context) liquidity_provider.GeneralConfiguration {
	configuration, err := liquidity_provider.ValidateConfiguration(lp.signer, func() (*entities.Signed[liquidity_provider.GeneralConfiguration], error) {
		return lp.lpRepository.GetGeneralConfiguration(ctx)
	})
	if err != nil {
		lp.logConfigError("general", err)
		return liquidity_provider.DefaultGeneralConfiguration()
	}
	return configuration.Value
}

func (lp *LocalLiquidityProvider) PegoutConfiguration(ctx context.Context) liquidity_provider.PegoutConfiguration {
	configuration, err := liquidity_provider.ValidateConfiguration(lp.signer, func() (*entities.Signed[liquidity_provider.PegoutConfiguration], error) {
		return lp.lpRepository.GetPegoutConfiguration(ctx)
	})
	if err != nil {
		lp.logConfigError("pegout", err)
		return liquidity_provider.DefaultPegoutConfiguration()
	}
	return configuration.Value
}

func (lp *LocalLiquidityProvider) PeginConfiguration(ctx context.Context) liquidity_provider.PeginConfiguration {
	configuration, err := liquidity_provider.ValidateConfiguration(lp.signer, func() (*entities.Signed[liquidity_provider.PeginConfiguration], error) {
		return lp.lpRepository.GetPeginConfiguration(ctx)
	})
	if err != nil {
		lp.logConfigError("pegin", err)
		return liquidity_provider.DefaultPeginConfiguration()
	}
	return configuration.Value
}

func (lp *LocalLiquidityProvider) GetSigner() entities.Signer {
	return lp.signer
}

func (lp *LocalLiquidityProvider) logConfigError(displayName string, err error) {
	if errors.Is(err, liquidity_provider.ConfigurationNotFoundError) {
		log.Warnf("Custom %s configuration not found. Using default configuration.", displayName)
	} else if errors.Is(err, liquidity_provider.InvalidSignatureError) {
		log.Errorf("Invalid %s configuration signature. Using default configuration.", displayName)
	} else if errors.Is(err, entities.IntegrityError) {
		log.Errorf("Tampered %s configuration. Using default configuration.", displayName)
	} else {
		log.Errorf("Error getting %s configuration, using default configuration. Error: %v", displayName, err)
	}
}
