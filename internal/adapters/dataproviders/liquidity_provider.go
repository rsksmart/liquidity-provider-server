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
	env              *Configuration
	peginRepository  quote.PeginQuoteRepository
	pegoutRepository quote.PegoutQuoteRepository
	lpRepository     liquidity_provider.LiquidityProviderRepository
	rsk              blockchain.RootstockRpcServer
	signer           rootstock.TransactionSigner
	btc              blockchain.BitcoinWallet
	lbc              blockchain.LiquidityBridgeContract
}

func NewLocalLiquidityProvider(
	env *Configuration,
	peginRepository quote.PeginQuoteRepository,
	pegoutRepository quote.PegoutQuoteRepository,
	lpRepository liquidity_provider.LiquidityProviderRepository,
	rsk blockchain.RootstockRpcServer,
	signer rootstock.TransactionSigner,
	btc blockchain.BitcoinWallet,
	lbc blockchain.LiquidityBridgeContract,
) *LocalLiquidityProvider {
	return &LocalLiquidityProvider{
		env:              env,
		peginRepository:  peginRepository,
		pegoutRepository: pegoutRepository,
		lpRepository:     lpRepository,
		rsk:              rsk,
		signer:           signer,
		btc:              btc,
		lbc:              lbc,
	}
}

func (lp *LocalLiquidityProvider) RskAddress() string {
	return lp.signer.Address().String()
}

func (lp *LocalLiquidityProvider) BtcAddress() string {
	return lp.env.BtcConfig.BtcAddress
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

func (lp *LocalLiquidityProvider) HasPegoutLiquidity(ctx context.Context, requiredLiquidity *entities.Wei) error {
	lockedLiquidity := new(entities.Wei)
	log.Debug("Verifying if has liquidity")
	liquidity, err := lp.btc.GetBalance()
	if err != nil {
		return err
	}
	log.Debugf("Liquidity: %s satoshi\n", liquidity.ToSatoshi().String())
	quotes, err := lp.pegoutRepository.GetRetainedQuoteByState(ctx,
		quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations, quote.PegoutStateSendPegoutFailed,
	)
	if err != nil {
		return err
	}
	for _, retainedQuote := range quotes {
		lockedLiquidity.Add(lockedLiquidity, retainedQuote.RequiredLiquidity)
	}
	log.Debugf("Locked Liquidity: %s satoshi\n", lockedLiquidity.ToSatoshi().String())
	availableLiquidity := new(entities.Wei).Sub(liquidity, lockedLiquidity)
	if availableLiquidity.Cmp(requiredLiquidity) >= 0 {
		return nil
	} else {
		return fmt.Errorf(
			"not enough liquidity, missing %s satoshi\n",
			requiredLiquidity.Sub(requiredLiquidity, availableLiquidity).ToSatoshi().String(),
		)
	}
}

func (lp *LocalLiquidityProvider) HasPeginLiquidity(ctx context.Context, requiredLiquidity *entities.Wei) error {
	liquidity := new(entities.Wei)
	lockedLiquidity := new(entities.Wei)
	log.Debug("Verifying if has liquidity")
	lpRskBalance, err := lp.rsk.GetBalance(ctx, lp.RskAddress())
	if err != nil {
		return err
	}
	lpLbcBalance, err := lp.lbc.GetBalance(lp.RskAddress())
	if err != nil {
		return err
	}
	liquidity.Add(lpRskBalance, lpLbcBalance)
	log.Debugf("Liquidity: %s wei\n", liquidity.String())
	quotes, err := lp.peginRepository.GetRetainedQuoteByState(ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateCallForUserFailed)
	if err != nil {
		return err
	}
	for _, retainedQuote := range quotes {
		lockedLiquidity.Add(lockedLiquidity, retainedQuote.RequiredLiquidity)
	}
	log.Debugf("Locked Liquidity: %s wei\n", lockedLiquidity.String())
	availableLiquidity := new(entities.Wei).Sub(liquidity, lockedLiquidity)
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
