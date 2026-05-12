package reports

import (
	"context"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

const (
	DateFormat = time.DateOnly
)

type SummaryResult struct {
	PeginSummary  SummaryData `json:"peginSummary"`
	PegoutSummary SummaryData `json:"pegoutSummary"`
}

type SummaryData struct {
	TotalQuotesCount          int64         `json:"totalQuotesCount"`
	AcceptedQuotesCount       int64         `json:"acceptedQuotesCount"` // value + gas fees
	TotalAcceptedQuotesAmount *entities.Wei `json:"totalAcceptedQuotesAmount"`
	PaidQuotesCount           int64         `json:"paidQuotesCount"` // value + gas fees
	PaidQuotesAmount          *entities.Wei `json:"paidQuotesAmount"`
	RefundedQuotesCount       int64         `json:"refundedQuotesCount"` // value + gas fees + call fee
	TotalRefundedQuotesAmount *entities.Wei `json:"totalRefundedQuotesAmount"`
	PenalizationsCount        int64         `json:"penalizationsCount"`
	TotalPenalizationsAmount  *entities.Wei `json:"totalPenalizationsAmount"`
}

type SummariesUseCase struct {
	peginRepo     quote.PeginQuoteRepository
	pegoutRepo    quote.PegoutQuoteRepository
	penalizedRepo penalization.PenalizedEventRepository
}

func NewSummariesUseCase(
	peginRepo quote.PeginQuoteRepository,
	pegoutRepo quote.PegoutQuoteRepository,
	penalizedRepo penalization.PenalizedEventRepository,
) *SummariesUseCase {
	return &SummariesUseCase{peginRepo: peginRepo, pegoutRepo: pegoutRepo, penalizedRepo: penalizedRepo}
}

func (u *SummariesUseCase) Run(ctx context.Context, startDate, endDate time.Time) (SummaryResult, error) {
	peginData, err := u.aggregatePeginData(ctx, startDate, endDate)
	if err != nil {
		return SummaryResult{}, usecases.WrapUseCaseError(usecases.SummariesUseCaseId, err)
	}
	pegoutData, err := u.aggregatePegoutData(ctx, startDate, endDate)
	if err != nil {
		return SummaryResult{}, usecases.WrapUseCaseError(usecases.SummariesUseCaseId, err)
	}
	return SummaryResult{PeginSummary: peginData, PegoutSummary: pegoutData}, nil
}

// ============================================================================
// Pegin aggregation
// ============================================================================

func (u *SummariesUseCase) aggregatePeginData(ctx context.Context, startDate, endDate time.Time) (SummaryData, error) {
	quotesWithRetained, err := u.peginRepo.GetQuotesWithRetainedByStateAndDate(ctx, getAllPeginStates(), startDate, endDate)
	if err != nil {
		return SummaryData{}, usecases.WrapUseCaseError(usecases.SummariesUseCaseId, err)
	}

	data, acceptedHashes := processPeginQuotes(quotesWithRetained, getPeginPaidStates(), getPeginRefundedStates())

	penalizations, err := u.getPenalizationsSummary(ctx, acceptedHashes)
	if err != nil {
		return SummaryData{}, usecases.WrapUseCaseError(usecases.SummariesUseCaseId, err)
	}

	data.PenalizationsCount = penalizations.Count
	data.TotalPenalizationsAmount = penalizations.Total

	return data, nil
}

func getAllPeginStates() []quote.PeginState {
	return []quote.PeginState{
		quote.PeginStateWaitingForDeposit,
		quote.PeginStateWaitingForDepositConfirmations,
		quote.PeginStateTimeForDepositElapsed,
		quote.PeginStateCallForUserSucceeded,
		quote.PeginStateCallForUserFailed,
		quote.PeginStateRegisterPegInSucceeded,
		quote.PeginStateRegisterPegInFailed,
	}
}

func getPeginPaidStates() map[quote.PeginState]bool {
	return map[quote.PeginState]bool{
		quote.PeginStateCallForUserSucceeded:   true,
		quote.PeginStateRegisterPegInSucceeded: true,
		quote.PeginStateRegisterPegInFailed:    true,
	}
}

func getPeginRefundedStates() map[quote.PeginState]bool {
	return map[quote.PeginState]bool{
		quote.PeginStateRegisterPegInSucceeded: true,
	}
}

func processPeginQuotes(quotesWithRetained []quote.PeginQuoteWithRetained, paidStates, refundedStates map[quote.PeginState]bool) (SummaryData, []string) {
	data := SummaryData{
		TotalAcceptedQuotesAmount: entities.NewWei(0),
		PaidQuotesAmount:          entities.NewWei(0),
		TotalRefundedQuotesAmount: entities.NewWei(0),
		TotalPenalizationsAmount:  entities.NewWei(0),
	}
	acceptedHashes := make([]string, 0)

	for _, quoteWithRetained := range quotesWithRetained {
		data.TotalQuotesCount++

		// From here on, we are only processing quotes with retained data (accepted quotes)
		if quoteWithRetained.RetainedQuote.QuoteHash != "" {
			data.AcceptedQuotesCount++
			data.TotalAcceptedQuotesAmount.Add(data.TotalAcceptedQuotesAmount, quoteWithRetained.Quote.Value)
			data.TotalAcceptedQuotesAmount.Add(data.TotalAcceptedQuotesAmount, quoteWithRetained.Quote.GasFee)
			acceptedHashes = append(acceptedHashes, quoteWithRetained.RetainedQuote.QuoteHash)

			if paidStates[quoteWithRetained.RetainedQuote.State] {
				data.PaidQuotesCount++
				data.PaidQuotesAmount.Add(data.PaidQuotesAmount, quoteWithRetained.Quote.Value)
				data.PaidQuotesAmount.Add(data.PaidQuotesAmount, quoteWithRetained.Quote.GasFee)
			}

			if refundedStates[quoteWithRetained.RetainedQuote.State] {
				data.RefundedQuotesCount++
				data.TotalRefundedQuotesAmount.Add(data.TotalRefundedQuotesAmount, quoteWithRetained.Quote.Value)
				data.TotalRefundedQuotesAmount.Add(data.TotalRefundedQuotesAmount, quoteWithRetained.Quote.GasFee)
				data.TotalRefundedQuotesAmount.Add(data.TotalRefundedQuotesAmount, quoteWithRetained.Quote.CallFee)
			}
		}
	}

	return data, acceptedHashes
}

// ============================================================================
// Pegout aggregation
// ============================================================================

func (u *SummariesUseCase) aggregatePegoutData(ctx context.Context, startDate, endDate time.Time) (SummaryData, error) {
	quotesWithRetained, err := u.pegoutRepo.GetQuotesWithRetainedByStateAndDate(ctx, getAllPegoutStates(), startDate, endDate)
	if err != nil {
		return SummaryData{}, usecases.WrapUseCaseError(usecases.SummariesUseCaseId, err)
	}

	data, acceptedHashes := processPegoutQuotes(quotesWithRetained, getPegoutPaidStates(), getPegoutRefundedStates())

	penalizations, err := u.getPenalizationsSummary(ctx, acceptedHashes)
	if err != nil {
		return SummaryData{}, usecases.WrapUseCaseError(usecases.SummariesUseCaseId, err)
	}

	data.PenalizationsCount = penalizations.Count
	data.TotalPenalizationsAmount = penalizations.Total

	return data, nil
}

func getAllPegoutStates() []quote.PegoutState {
	return []quote.PegoutState{
		quote.PegoutStateWaitingForDeposit,
		quote.PegoutStateWaitingForDepositConfirmations,
		quote.PegoutStateTimeForDepositElapsed,
		quote.PegoutStateSendPegoutSucceeded,
		quote.PegoutStateSendPegoutFailed,
		quote.PegoutStateRefundPegOutSucceeded,
		quote.PegoutStateRefundPegOutFailed,
		quote.PegoutStateBridgeTxSucceeded,
		quote.PegoutStateBridgeTxFailed,
		quote.PegoutStateBtcReleased,
	}
}

func getPegoutPaidStates() map[quote.PegoutState]bool {
	return map[quote.PegoutState]bool{
		quote.PegoutStateSendPegoutSucceeded:   true,
		quote.PegoutStateRefundPegOutSucceeded: true,
		quote.PegoutStateRefundPegOutFailed:    true,
		quote.PegoutStateBridgeTxSucceeded:     true,
		quote.PegoutStateBridgeTxFailed:        true,
		quote.PegoutStateBtcReleased:           true,
	}
}

func getPegoutRefundedStates() map[quote.PegoutState]bool {
	return map[quote.PegoutState]bool{
		quote.PegoutStateRefundPegOutSucceeded: true,
		quote.PegoutStateBridgeTxSucceeded:     true,
		quote.PegoutStateBridgeTxFailed:        true,
		quote.PegoutStateBtcReleased:           true,
	}
}

func processPegoutQuotes(quotesWithRetained []quote.PegoutQuoteWithRetained, paidStates, refundedStates map[quote.PegoutState]bool) (SummaryData, []string) {
	data := SummaryData{
		TotalAcceptedQuotesAmount: entities.NewWei(0),
		PaidQuotesAmount:          entities.NewWei(0),
		TotalRefundedQuotesAmount: entities.NewWei(0),
		TotalPenalizationsAmount:  entities.NewWei(0),
	}
	acceptedHashes := make([]string, 0)

	for _, quoteWithRetained := range quotesWithRetained {
		data.TotalQuotesCount++

		// From here on, we are only processing quotes with retained data (accepted quotes)
		if quoteWithRetained.RetainedQuote.QuoteHash != "" {
			data.AcceptedQuotesCount++
			data.TotalAcceptedQuotesAmount.Add(data.TotalAcceptedQuotesAmount, quoteWithRetained.Quote.Value)
			data.TotalAcceptedQuotesAmount.Add(data.TotalAcceptedQuotesAmount, quoteWithRetained.Quote.GasFee)
			acceptedHashes = append(acceptedHashes, quoteWithRetained.RetainedQuote.QuoteHash)

			if paidStates[quoteWithRetained.RetainedQuote.State] {
				data.PaidQuotesCount++
				data.PaidQuotesAmount.Add(data.PaidQuotesAmount, quoteWithRetained.Quote.Value)
				data.PaidQuotesAmount.Add(data.PaidQuotesAmount, quoteWithRetained.Quote.GasFee)
			}

			if refundedStates[quoteWithRetained.RetainedQuote.State] {
				data.RefundedQuotesCount++
				data.TotalRefundedQuotesAmount.Add(data.TotalRefundedQuotesAmount, quoteWithRetained.Quote.Value)
				data.TotalRefundedQuotesAmount.Add(data.TotalRefundedQuotesAmount, quoteWithRetained.Quote.GasFee)
				data.TotalRefundedQuotesAmount.Add(data.TotalRefundedQuotesAmount, quoteWithRetained.Quote.CallFee)
			}
		}
	}

	return data, acceptedHashes
}

// ============================================================================
// Shared utilities
// ============================================================================

type penalizationsSummary struct {
	Count int64
	Total *entities.Wei
}

func (u *SummariesUseCase) getPenalizationsSummary(ctx context.Context, quoteHashes []string) (penalizationsSummary, error) {
	result := penalizationsSummary{
		Count: 0,
		Total: entities.NewWei(0),
	}

	if len(quoteHashes) == 0 {
		return result, nil
	}

	penalizations, err := u.penalizedRepo.GetPenalizationsByQuoteHashes(ctx, quoteHashes)
	if err != nil {
		return penalizationsSummary{}, err
	}

	for _, p := range penalizations {
		result.Total.Add(result.Total, p.Penalty)
	}
	result.Count = int64(len(penalizations))

	return result, nil
}
