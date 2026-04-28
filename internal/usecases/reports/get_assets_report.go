package reports

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
)

type BtcAssetLocation struct {
	BtcWallet  *entities.Wei `json:"btcWallet" validate:"required"`
	Federation *entities.Wei `json:"federation" validate:"required"`
	RskWallet  *entities.Wei `json:"rskWallet" validate:"required"`
	Lbc        *entities.Wei `json:"lbc" validate:"required"`
}

type BtcAssetAllocation struct {
	ReservedForUsers *entities.Wei `json:"reservedForUsers" validate:"required"`
	WaitingForRefund *entities.Wei `json:"waitingForRefund" validate:"required"`
	Available        *entities.Wei `json:"available" validate:"required"`
}

type BtcAssetReport struct {
	Total      *entities.Wei      `json:"total" validate:"required"`
	Location   BtcAssetLocation   `json:"location" validate:"required"`
	Allocation BtcAssetAllocation `json:"allocation" validate:"required"`
}

type RbtcAssetLocation struct {
	RskWallet  *entities.Wei `json:"rskWallet" validate:"required"`
	Lbc        *entities.Wei `json:"lbc" validate:"required"`
	Federation *entities.Wei `json:"federation" validate:"required"`
}

type RbtcAssetAllocation struct {
	ReservedForUsers *entities.Wei `json:"reservedForUsers" validate:"required"`
	WaitingForRefund *entities.Wei `json:"waitingForRefund" validate:"required"`
	Available        *entities.Wei `json:"available" validate:"required"`
}

type RbtcAssetReport struct {
	Total      *entities.Wei       `json:"total" validate:"required"`
	Location   RbtcAssetLocation   `json:"location" validate:"required"`
	Allocation RbtcAssetAllocation `json:"allocation" validate:"required"`
}
type GetAssetsReportResult struct {
	BtcAssetReport  BtcAssetReport
	RbtcAssetReport RbtcAssetReport
}

type GetAssetsReportUseCase struct {
	btcWallet        blockchain.BitcoinWallet
	rsk              blockchain.Rpc
	lp               liquidity_provider.LiquidityProvider
	peginProvider    liquidity_provider.PeginLiquidityProvider
	pegoutProvider   liquidity_provider.PegoutLiquidityProvider
	peginRepository  quote.PeginQuoteRepository
	pegoutRepository quote.PegoutQuoteRepository
	contracts        blockchain.RskContracts
}

func NewGetAssetsReportUseCase(
	wallet blockchain.BitcoinWallet,
	rsk blockchain.Rpc,
	lp liquidity_provider.LiquidityProvider,
	peginProvider liquidity_provider.PeginLiquidityProvider,
	pegoutProvider liquidity_provider.PegoutLiquidityProvider,
	peginRepository quote.PeginQuoteRepository,
	pegoutRepository quote.PegoutQuoteRepository,
	contracts blockchain.RskContracts,
) *GetAssetsReportUseCase {
	return &GetAssetsReportUseCase{
		btcWallet:        wallet,
		rsk:              rsk,
		lp:               lp,
		peginProvider:    peginProvider,
		pegoutProvider:   pegoutProvider,
		peginRepository:  peginRepository,
		pegoutRepository: pegoutRepository,
		contracts:        contracts,
	}
}

func (useCase *GetAssetsReportUseCase) Run(ctx context.Context) (GetAssetsReportResult, error) {
	btcReport, err := useCase.calculateBtcAssetReport(ctx)
	if err != nil {
		return GetAssetsReportResult{}, err
	}

	rbtcReport, err := useCase.calculateRbtcAssetReport(ctx, btcReport.Location.RskWallet)
	if err != nil {
		return GetAssetsReportResult{}, err
	}

	return GetAssetsReportResult{
		BtcAssetReport:  btcReport,
		RbtcAssetReport: rbtcReport,
	}, nil
}

func (useCase *GetAssetsReportUseCase) calculateBtcAssetReport(ctx context.Context) (BtcAssetReport, error) {
	btcWalletBalance, err := useCase.btcWallet.GetBalance()
	if err != nil {
		return BtcAssetReport{}, err
	}

	// A threshold of RBTC was reached and a bridge transaction initiated but not yet finished
	btcRebalancing, err := useCase.sumPegoutQuotesByState(ctx, quote.PegoutStateBridgeTxSucceeded)
	if err != nil {
		return BtcAssetReport{}, err
	}

	// A threshold of RBTC has not been reached yet and is sitting in the RBTC wallet
	btcWaitingForRebalancing, err := useCase.sumPegoutQuotesByState(ctx, quote.PegoutStateRefundPegOutSucceeded)
	if err != nil {
		return BtcAssetReport{}, err
	}

	// The LP already sent the BTC to the user, but the LBC has not yet sent the RBTC to the LP
	btcInLbc, err := useCase.sumPegoutQuotesByState(ctx, quote.PegoutStateSendPegoutSucceeded)
	if err != nil {
		return BtcAssetReport{}, err
	}

	// Already accepted pegout quotes
	btcReservedForUsers, err := useCase.sumPegoutQuotesByState(ctx, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations)
	if err != nil {
		return BtcAssetReport{}, err
	}

	btcWaitingForRefund, err := useCase.sumPegoutQuotesByState(ctx, quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateSendPegoutSucceeded, quote.PegoutStateBridgeTxSucceeded)
	if err != nil {
		return BtcAssetReport{}, err
	}

	// Calculate BTC total as sum of all location fields
	btcTotal := entities.NewWei(0)
	btcTotal.Add(btcTotal, btcWalletBalance)
	btcTotal.Add(btcTotal, btcRebalancing)
	btcTotal.Add(btcTotal, btcWaitingForRebalancing)
	btcTotal.Add(btcTotal, btcInLbc)

	return BtcAssetReport{
		Total: btcTotal,
		Location: BtcAssetLocation{
			BtcWallet:  btcWalletBalance,
			Federation: btcRebalancing,
			RskWallet:  btcWaitingForRebalancing,
			Lbc:        btcInLbc,
		},
		Allocation: BtcAssetAllocation{
			ReservedForUsers: btcReservedForUsers,
			WaitingForRefund: btcWaitingForRefund,
			Available:        entities.NewWei(0).Sub(btcWalletBalance, btcReservedForUsers),
		},
	}, nil
}

func (useCase *GetAssetsReportUseCase) calculateRbtcAssetReport(ctx context.Context, btcWaitingForRebalancing *entities.Wei) (RbtcAssetReport, error) {
	rbtcWalletBalance, err := useCase.rsk.Rsk.GetBalance(ctx, useCase.lp.RskAddress())
	if err != nil {
		return RbtcAssetReport{}, err
	}
	// A part of the RBTC in the RSK wallet is a representation of BTC waiting to be sent to the bridge for rebalancing
	rbtcInRskWallet := entities.NewWei(0).Sub(rbtcWalletBalance, btcWaitingForRebalancing)

	// Initial balance + the cases when registerPegin succeded and the LBC balance for the LP was increased
	rbtcLockedInLbc, err := useCase.contracts.PegIn.GetBalance(useCase.lp.RskAddress())
	if err != nil {
		return RbtcAssetReport{}, err
	}

	// LP spent RBTC by calling callForUser so the user now have the funds but the LP has not called registerPegin yet and is waiting the
	// number of blocks of the native pegin
	rbtcWaitingForRefund, err := useCase.sumPeginQuotesByState(ctx, quote.PeginStateCallForUserSucceeded)
	if err != nil {
		return RbtcAssetReport{}, err
	}

	// Already accepted pegin quotes
	rbtcReservedForUsers, err := useCase.sumPeginQuotesByState(ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations)
	if err != nil {
		return RbtcAssetReport{}, err
	}

	rbtcTotal := entities.NewWei(0)
	rbtcTotal.Add(rbtcTotal, rbtcInRskWallet)
	rbtcTotal.Add(rbtcTotal, rbtcLockedInLbc)
	rbtcTotal.Add(rbtcTotal, rbtcWaitingForRefund)

	return RbtcAssetReport{
		Total: rbtcTotal,
		Location: RbtcAssetLocation{
			RskWallet:  rbtcInRskWallet,
			Lbc:        rbtcLockedInLbc,
			Federation: rbtcWaitingForRefund,
		},
		Allocation: RbtcAssetAllocation{
			ReservedForUsers: rbtcReservedForUsers,
			WaitingForRefund: rbtcWaitingForRefund,
			Available: entities.NewWei(0).Add(
				entities.NewWei(0).Sub(rbtcInRskWallet, rbtcReservedForUsers),
				rbtcLockedInLbc,
			),
		},
	}, nil
}

func (useCase *GetAssetsReportUseCase) sumPegoutQuotesByState(ctx context.Context, states ...quote.PegoutState) (*entities.Wei, error) {
	total := entities.NewWei(0)

	pegoutQuotes, err := useCase.pegoutRepository.GetQuotesByState(ctx, states...)
	if err != nil {
		return nil, err
	}

	for _, pegoutQuote := range pegoutQuotes {
		total.Add(total, pegoutQuote.Total())
	}

	return total, nil
}

func (useCase *GetAssetsReportUseCase) sumPeginQuotesByState(ctx context.Context, states ...quote.PeginState) (*entities.Wei, error) {
	total := entities.NewWei(0)

	peginQuotes, err := useCase.peginRepository.GetQuotesByState(ctx, states...)
	if err != nil {
		return nil, err
	}

	for _, peginQuote := range peginQuotes {
		total.Add(total, peginQuote.Total())
	}

	return total, nil
}
