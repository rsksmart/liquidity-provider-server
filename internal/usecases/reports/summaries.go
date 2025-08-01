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
	AcceptedQuotesCount       int64         `json:"acceptedQuotesCount"`
	PaidQuotesCount           int64         `json:"paidQuotesCount"`
	PaidQuotesAmount          *entities.Wei `json:"paidQuotesAmount"`
	TotalAcceptedQuotedAmount *entities.Wei `json:"totalAcceptedQuotedAmount"`
	TotalFeesCollected        *entities.Wei `json:"totalFeesCollected"`
	RefundedQuotesCount       int64         `json:"refundedQuotesCount"`
	TotalPenaltyAmount        *entities.Wei `json:"totalPenaltyAmount"`
	LpEarnings                *entities.Wei `json:"lpEarnings"`
}

type summaryTotals struct {
	AcceptedTotalAmount *entities.Wei
	TotalFees           *entities.Wei
	CallFees            *entities.Wei
	TotalPenalty        *entities.Wei
}

func newSummaryTotals() *summaryTotals {
	return &summaryTotals{
		AcceptedTotalAmount: entities.NewWei(0),
		TotalFees:           entities.NewWei(0),
		CallFees:            entities.NewWei(0),
		TotalPenalty:        entities.NewWei(0),
	}
}

type SummariesUseCase struct {
	peginRepo     quote.PeginQuoteRepository
	pegoutRepo    quote.PegoutQuoteRepository
	penalizedRepo penalization.PenalizedEventRepository
}

func NewSummaryData() SummaryData {
	return SummaryData{
		TotalQuotesCount:          0,
		AcceptedQuotesCount:       0,
		PaidQuotesCount:           0,
		PaidQuotesAmount:          entities.NewWei(0),
		TotalAcceptedQuotedAmount: entities.NewWei(0),
		TotalFeesCollected:        entities.NewWei(0),
		RefundedQuotesCount:       0,
		TotalPenaltyAmount:        entities.NewWei(0),
		LpEarnings:                entities.NewWei(0),
	}
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

func (u *SummariesUseCase) buildPenalizedMap(ctx context.Context, quoteHashes []string) (map[string]*entities.Wei, error) {
	penalizedMap := make(map[string]*entities.Wei)
	if u.penalizedRepo == nil || len(quoteHashes) == 0 {
		return penalizedMap, nil
	}
	events, err := u.penalizedRepo.GetPenalizationsByQuoteHashes(ctx, quoteHashes)
	if err != nil {
		return nil, err
	}
	for _, ev := range events {
		penalizedMap[ev.QuoteHash] = ev.Penalty
	}
	return penalizedMap, nil
}

func extractAcceptedPeginHashes(pairs []quote.PeginQuoteWithRetained) []string {
	hashes := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		if pair.RetainedQuote.QuoteHash != "" {
			hashes = append(hashes, pair.RetainedQuote.QuoteHash)
		}
	}
	return hashes
}

func extractAcceptedPegoutHashes(pairs []quote.PegoutQuoteWithRetained) []string {
	hashes := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		if pair.RetainedQuote.QuoteHash != "" {
			hashes = append(hashes, pair.RetainedQuote.QuoteHash)
		}
	}
	return hashes
}

func (u *SummariesUseCase) aggregatePeginData(ctx context.Context, startDate, endDate time.Time) (SummaryData, error) {
	var (
		quotePairs []quote.PeginQuoteWithRetained
		err        error
	)
	quotePairs, _, err = u.peginRepo.ListQuotesByDateRange(ctx, startDate, endDate, 0, 0)
	if err != nil {
		return NewSummaryData(), usecases.WrapUseCaseError(usecases.SummariesUseCaseId, err)
	}
	data := NewSummaryData()
	acceptedQuotesCount := 0
	totals := newSummaryTotals()
	penalizedMap, errPen := u.buildPenalizedMap(ctx, extractAcceptedPeginHashes(quotePairs))
	if errPen != nil {
		return NewSummaryData(), usecases.WrapUseCaseError(usecases.SummariesUseCaseId, errPen)
	}
	for _, pair := range quotePairs {
		if pair.RetainedQuote.QuoteHash != "" {
			acceptedQuotesCount++
			processPeginPair(pair, &data, totals)
			if penalty, ok := penalizedMap[pair.RetainedQuote.QuoteHash]; ok {
				totals.TotalPenalty.Add(totals.TotalPenalty, penalty)
			}
		}
	}
	data.TotalQuotesCount = int64(len(quotePairs))
	data.AcceptedQuotesCount = int64(acceptedQuotesCount)
	data.TotalAcceptedQuotedAmount = totals.AcceptedTotalAmount
	data.TotalFeesCollected = totals.TotalFees
	data.TotalPenaltyAmount = totals.TotalPenalty
	lpEarnings := new(entities.Wei)
	lpEarnings.Add(lpEarnings, totals.CallFees)
	data.LpEarnings = lpEarnings
	return data, nil
}

func (u *SummariesUseCase) aggregatePegoutData(ctx context.Context, startDate, endDate time.Time) (SummaryData, error) {
	var (
		quotePairs []quote.PegoutQuoteWithRetained
		err        error
	)
	quotePairs, _, err = u.pegoutRepo.ListQuotesByDateRange(ctx, startDate, endDate, 0, 0)
	if err != nil {
		return NewSummaryData(), usecases.WrapUseCaseError(usecases.SummariesUseCaseId, err)
	}
	data := NewSummaryData()
	acceptedQuotesCount := 0
	totals := newSummaryTotals()
	penalizedMap, errPen := u.buildPenalizedMap(ctx, extractAcceptedPegoutHashes(quotePairs))
	if errPen != nil {
		return NewSummaryData(), usecases.WrapUseCaseError(usecases.SummariesUseCaseId, errPen)
	}
	for _, pair := range quotePairs {
		if pair.RetainedQuote.QuoteHash != "" {
			acceptedQuotesCount++
			processPegoutPair(pair, &data, totals)
			if penalty, ok := penalizedMap[pair.RetainedQuote.QuoteHash]; ok {
				totals.TotalPenalty.Add(totals.TotalPenalty, penalty)
			}
		}
	}
	data.TotalQuotesCount = int64(len(quotePairs))
	data.AcceptedQuotesCount = int64(acceptedQuotesCount)
	data.TotalAcceptedQuotedAmount = totals.AcceptedTotalAmount
	data.TotalFeesCollected = totals.TotalFees
	data.TotalPenaltyAmount = totals.TotalPenalty
	lpEarnings := new(entities.Wei)
	lpEarnings.Add(lpEarnings, totals.CallFees)
	data.LpEarnings = lpEarnings
	return data, nil
}

func processPeginPair(
	pair quote.PeginQuoteWithRetained,
	data *SummaryData,
	totals *summaryTotals,
) {
	q := pair.Quote
	retained := pair.RetainedQuote
	totals.AcceptedTotalAmount.Add(totals.AcceptedTotalAmount, q.Value)
	callFee, gasFee := q.CallFee, q.GasFee
	if isPeginPaidQuote(retained) || isPeginRefundedQuote(retained) {
		data.PaidQuotesCount++
		quoteValue := q.Value
		data.PaidQuotesAmount.Add(data.PaidQuotesAmount, quoteValue)
		totals.CallFees.Add(totals.CallFees, callFee)
		totals.TotalFees.Add(totals.TotalFees, callFee)
		totals.TotalFees.Add(totals.TotalFees, gasFee)
	}
	if isPeginRefundedQuote(retained) {
		data.RefundedQuotesCount++
	}
}

func processPegoutPair(
	pair quote.PegoutQuoteWithRetained,
	data *SummaryData,
	totals *summaryTotals,
) {
	q := pair.Quote
	retained := pair.RetainedQuote
	totals.AcceptedTotalAmount.Add(totals.AcceptedTotalAmount, q.Value)
	callFee, gasFee := q.CallFee, q.GasFee
	if isPegoutPaidQuote(retained) || isPegoutRefundedQuote(retained) {
		data.PaidQuotesCount++
		quoteValue := q.Value
		data.PaidQuotesAmount.Add(data.PaidQuotesAmount, quoteValue)
		totals.CallFees.Add(totals.CallFees, callFee)
		totals.TotalFees.Add(totals.TotalFees, callFee)
		totals.TotalFees.Add(totals.TotalFees, gasFee)
	}
	if isPegoutRefundedQuote(retained) {
		data.RefundedQuotesCount++
	}
}

func isPeginPaidQuote(retained quote.RetainedPeginQuote) bool {
	return retained.State == quote.PeginStateCallForUserSucceeded
}

func isPegoutPaidQuote(retained quote.RetainedPegoutQuote) bool {
	return retained.State == quote.PegoutStateSendPegoutSucceeded
}

func isPeginRefundedQuote(retained quote.RetainedPeginQuote) bool {
	return retained.State == quote.PeginStateRegisterPegInSucceeded
}

func isPegoutRefundedQuote(retained quote.RetainedPegoutQuote) bool {
	return retained.State == quote.PegoutStateRefundPegOutSucceeded || retained.State == quote.PegoutStateBridgeTxSucceeded
}
