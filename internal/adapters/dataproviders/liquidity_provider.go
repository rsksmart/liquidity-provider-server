package dataproviders

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
	"slices"
)

type LocalLiquidityProvider struct {
	env              *Configuration
	peginRepository  quote.PeginQuoteRepository
	pegoutRepository quote.PegoutQuoteRepository
	rsk              blockchain.RootstockRpcServer
	signer           rootstock.TransactionSigner
	btc              blockchain.BitcoinWallet
	lbc              blockchain.LiquidityBridgeContract
}

func NewLocalLiquidityProvider(
	env *Configuration,
	peginRepository quote.PeginQuoteRepository,
	pegoutRepository quote.PegoutQuoteRepository,
	rsk blockchain.RootstockRpcServer,
	signer rootstock.TransactionSigner,
	btc blockchain.BitcoinWallet,
	lbc blockchain.LiquidityBridgeContract,
) *LocalLiquidityProvider {
	return &LocalLiquidityProvider{
		env:              env,
		peginRepository:  peginRepository,
		pegoutRepository: pegoutRepository,
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

func (lp *LocalLiquidityProvider) ValidateAmountForPegout(amount *entities.Wei) error {
	if amount.Cmp(lp.MaxPegout()) <= 0 && amount.Cmp(lp.MinPegout()) >= 0 {
		return nil
	} else {
		return fmt.Errorf("%w [%v, %v]", usecases.AmountOutOfRangeError, lp.MinPegout(), lp.MaxPegout())
	}
}

func (lp *LocalLiquidityProvider) GetRootstockConfirmationsForValue(value *entities.Wei) uint16 {
	var values []int
	for key, _ := range lp.env.RskConfig.Confirmations {
		values = append(values, key)
	}
	slices.Sort(values)
	index := slices.IndexFunc(values, func(item int) bool {
		return int(value.AsBigInt().Int64()) < item
	})
	if index == -1 {
		return lp.env.RskConfig.Confirmations[values[len(values)-1]]
	} else {
		return lp.env.RskConfig.Confirmations[values[index]]
	}
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

func (lp *LocalLiquidityProvider) CallFeePegout() *entities.Wei {
	return lp.env.PegoutConfig.CallFee
}

func (lp *LocalLiquidityProvider) PenaltyFeePegout() *entities.Wei {
	return lp.env.PegoutConfig.PenaltyFee
}

func (lp *LocalLiquidityProvider) TimeForDepositPegout() uint32 {
	return lp.env.PegoutConfig.TimeForDeposit
}

func (lp *LocalLiquidityProvider) ExpireBlocksPegout() uint64 {
	return uint64(lp.env.PegoutConfig.ExpireBlocks)
}

func (lp *LocalLiquidityProvider) MaxPegout() *entities.Wei {
	return lp.env.PegoutConfig.MaxTransactionValue
}

func (lp *LocalLiquidityProvider) MinPegout() *entities.Wei {
	return lp.env.PegoutConfig.MinTransactionValue
}

func (lp *LocalLiquidityProvider) MaxPegoutConfirmations() uint16 {
	// TODO replace in go 1.22 with
	// return slices.Max(maps.Values(lp.env.RskConfig.Confirmations))
	var maxValue uint16
	for _, value := range lp.env.RskConfig.Confirmations {
		if maxValue < value {
			maxValue = value
		}
	}
	return maxValue
}

func (lp *LocalLiquidityProvider) ValidateAmountForPegin(amount *entities.Wei) error {
	if amount.Cmp(lp.MaxPegin()) <= 0 && amount.Cmp(lp.MinPegin()) >= 0 {
		return nil
	} else {
		return fmt.Errorf("%w [%v, %v]", usecases.AmountOutOfRangeError, lp.MinPegin(), lp.MaxPegin())
	}
}

func (lp *LocalLiquidityProvider) GetBitcoinConfirmationsForValue(value *entities.Wei) uint16 {
	var values []int
	for key, _ := range lp.env.BtcConfig.Confirmations {
		values = append(values, key)
	}
	slices.Sort(values)
	index := slices.IndexFunc(values, func(item int) bool {
		return int(value.AsBigInt().Int64()) < item
	})
	if index == -1 {
		return lp.env.BtcConfig.Confirmations[values[len(values)-1]]
	} else {
		return lp.env.BtcConfig.Confirmations[values[index]]
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

func (lp *LocalLiquidityProvider) CallTime() uint32 {
	return lp.env.PeginConfig.CallTime
}

func (lp *LocalLiquidityProvider) CallFeePegin() *entities.Wei {
	return lp.env.PeginConfig.CallFee
}

func (lp *LocalLiquidityProvider) PenaltyFeePegin() *entities.Wei {
	return lp.env.PeginConfig.PenaltyFee
}

func (lp *LocalLiquidityProvider) TimeForDepositPegin() uint32 {
	return lp.env.PeginConfig.TimeForDeposit
}

func (lp *LocalLiquidityProvider) MaxPegin() *entities.Wei {
	return lp.env.PeginConfig.MaxTransactionValue
}

func (lp *LocalLiquidityProvider) MinPegin() *entities.Wei {
	return lp.env.PeginConfig.MinTransactionValue
}

func (lp *LocalLiquidityProvider) MaxPeginConfirmations() uint16 {
	// TODO replace in go 1.22 with
	// return slices.Max(maps.Values(lp.env.BtcConfig.Confirmations))
	var maxValue uint16
	for _, value := range lp.env.BtcConfig.Confirmations {
		if maxValue < value {
			maxValue = value
		}
	}
	return maxValue
}
