package liquidity_provider

import (
	"context"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
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

type SummariesUseCase struct {
	peginRepo  quote.PeginQuoteRepository
	pegoutRepo quote.PegoutQuoteRepository
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

func NewSummariesUseCase(peginRepo quote.PeginQuoteRepository, pegoutRepo quote.PegoutQuoteRepository) *SummariesUseCase {
	return &SummariesUseCase{peginRepo: peginRepo, pegoutRepo: pegoutRepo}
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

func (u *SummariesUseCase) aggregatePeginData(ctx context.Context, startDate, endDate time.Time) (SummaryData, error) {
	var (
		quotePairs []quote.PeginQuoteWithRetained
		err        error
	)
	quotePairs, err = u.peginRepo.ListQuotesByDateRange(ctx, startDate, endDate)
	if err != nil {
		return NewSummaryData(), usecases.WrapUseCaseError(usecases.SummariesUseCaseId, err)
	}
	data := NewSummaryData()
	acceptedQuotesCount := 0
	acceptedTotalAmount := entities.NewWei(0)
	totalFees := entities.NewWei(0)
	callFees := entities.NewWei(0)
	totalPenalty := entities.NewWei(0)
	for _, pair := range quotePairs {
		if pair.RetainedQuote.QuoteHash != "" {
			acceptedQuotesCount++
			processPeginPair(pair, &data, acceptedTotalAmount, totalFees, callFees, totalPenalty)
		}
	}
	data.TotalQuotesCount = int64(len(quotePairs))
	data.AcceptedQuotesCount = int64(acceptedQuotesCount)
	data.TotalAcceptedQuotedAmount = acceptedTotalAmount
	data.TotalFeesCollected = totalFees
	data.TotalPenaltyAmount = totalPenalty
	lpEarnings := new(entities.Wei)
	lpEarnings.Add(lpEarnings, callFees)
	lpEarnings.Sub(lpEarnings, totalPenalty)
	data.LpEarnings = lpEarnings
	return data, nil
}

func (u *SummariesUseCase) aggregatePegoutData(ctx context.Context, startDate, endDate time.Time) (SummaryData, error) {
	var (
		quotePairs []quote.PegoutQuoteWithRetained
		err        error
	)
	quotePairs, err = u.pegoutRepo.ListQuotesByDateRange(ctx, startDate, endDate)
	if err != nil {
		return NewSummaryData(), usecases.WrapUseCaseError(usecases.SummariesUseCaseId, err)
	}
	data := NewSummaryData()
	acceptedQuotesCount := 0
	acceptedTotalAmount := entities.NewWei(0)
	totalFees := entities.NewWei(0)
	callFees := entities.NewWei(0)
	totalPenalty := entities.NewWei(0)
	for _, pair := range quotePairs {
		if pair.RetainedQuote.QuoteHash != "" {
			acceptedQuotesCount++
			processPegoutPair(pair, &data, acceptedTotalAmount, totalFees, callFees, totalPenalty)
		}
	}
	data.TotalQuotesCount = int64(len(quotePairs))
	data.AcceptedQuotesCount = int64(acceptedQuotesCount)
	data.TotalAcceptedQuotedAmount = acceptedTotalAmount
	data.TotalFeesCollected = totalFees
	data.TotalPenaltyAmount = totalPenalty
	lpEarnings := new(entities.Wei)
	lpEarnings.Add(lpEarnings, callFees)
	lpEarnings.Sub(lpEarnings, totalPenalty)
	data.LpEarnings = lpEarnings
	return data, nil
}

func processPeginPair(
	pair quote.PeginQuoteWithRetained,
	data *SummaryData,
	acceptedTotalAmount, totalFees, callFees, totalPenalty *entities.Wei,
) {
	q := pair.Quote
	retained := pair.RetainedQuote
	acceptedTotalAmount.Add(acceptedTotalAmount, q.Total())
	callFee, gasFee := q.CallFee, q.GasFee
	productFee := entities.NewUWei(q.ProductFeeAmount)
	penaltyFee := q.PenaltyFee
	if isPeginPaidQuote(retained) {
		data.PaidQuotesCount++
		data.PaidQuotesAmount.Add(data.PaidQuotesAmount, q.Total())
		callFees.Add(callFees, callFee)
		totalFees.Add(totalFees, callFee)
		totalFees.Add(totalFees, gasFee)
		totalFees.Add(totalFees, productFee)
	}
	if isPeginRefundedQuote(retained) {
		data.RefundedQuotesCount++
		totalPenalty.Add(totalPenalty, penaltyFee)
	}
}

func processPegoutPair(
	pair quote.PegoutQuoteWithRetained,
	data *SummaryData,
	acceptedTotalAmount, totalFees, callFees, totalPenalty *entities.Wei,
) {
	q := pair.Quote
	retained := pair.RetainedQuote
	acceptedTotalAmount.Add(acceptedTotalAmount, q.Total())
	callFee, gasFee := q.CallFee, q.GasFee
	productFee := entities.NewUWei(q.ProductFeeAmount)
	penaltyFee := entities.NewUWei(q.PenaltyFee)
	if isPegoutPaidQuote(retained) {
		data.PaidQuotesCount++
		data.PaidQuotesAmount.Add(data.PaidQuotesAmount, q.Total())
		callFees.Add(callFees, callFee)
		totalFees.Add(totalFees, callFee)
		totalFees.Add(totalFees, gasFee)
		totalFees.Add(totalFees, productFee)
	}
	if isPegoutRefundedQuote(retained) {
		data.RefundedQuotesCount++
		totalPenalty.Add(totalPenalty, penaltyFee)
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
