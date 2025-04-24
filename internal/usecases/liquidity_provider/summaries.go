package liquidity_provider

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

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

func processPeginRetainedQuote(
	ctx context.Context,
	repo quote.PeginQuoteRepository,
	retained quote.RetainedPeginQuote,
	data *SummaryData,
	acceptedTotalAmount, totalFees, callFees, totalPenalty *entities.Wei,
) {
	q, err := repo.GetQuote(ctx, retained.QuoteHash)
	if err != nil || q == nil {
		if err != nil {
			log.Errorf("Error getting quote %s: %v", retained.QuoteHash, err)
		}
		return
	}
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

func processPegoutRetainedQuote(
	ctx context.Context,
	repo quote.PegoutQuoteRepository,
	retained quote.RetainedPegoutQuote,
	data *SummaryData,
	acceptedTotalAmount, totalFees, callFees, totalPenalty *entities.Wei,
) {
	q, err := repo.GetQuote(ctx, retained.QuoteHash)
	if err != nil || q == nil {
		if err != nil {
			log.Errorf("Error getting quote %s: %v", retained.QuoteHash, err)
		}
		return
	}
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

func (u *SummariesUseCase) processPeginQuoteData(ctx context.Context, retainedQuotes []quote.RetainedPeginQuote) SummaryData {
	data := NewSummaryData()
	data.AcceptedQuotesCount = int64(len(retainedQuotes))
	acceptedTotalAmount := entities.NewWei(0)
	totalFees := entities.NewWei(0)
	callFees := entities.NewWei(0)
	totalPenalty := entities.NewWei(0)
	for _, retained := range retainedQuotes {
		processPeginRetainedQuote(ctx, u.peginRepo, retained, &data,
			acceptedTotalAmount, totalFees, callFees, totalPenalty)
	}
	data.TotalAcceptedQuotedAmount = acceptedTotalAmount
	data.TotalFeesCollected = totalFees
	data.TotalPenaltyAmount = totalPenalty
	lpEarnings := new(entities.Wei)
	lpEarnings.Add(lpEarnings, callFees)
	lpEarnings.Sub(lpEarnings, totalPenalty)
	data.LpEarnings = lpEarnings
	return data
}

func (u *SummariesUseCase) processPegoutQuoteData(ctx context.Context, retainedQuotes []quote.RetainedPegoutQuote) SummaryData {
	data := NewSummaryData()
	data.AcceptedQuotesCount = int64(len(retainedQuotes))
	acceptedTotalAmount := entities.NewWei(0)
	totalFees := entities.NewWei(0)
	callFees := entities.NewWei(0)
	totalPenalty := entities.NewWei(0)
	for _, retained := range retainedQuotes {
		processPegoutRetainedQuote(ctx, u.pegoutRepo, retained, &data,
			acceptedTotalAmount, totalFees, callFees, totalPenalty)
	}
	data.TotalAcceptedQuotedAmount = acceptedTotalAmount
	data.TotalFeesCollected = totalFees
	data.TotalPenaltyAmount = totalPenalty
	lpEarnings := new(entities.Wei)
	lpEarnings.Add(lpEarnings, callFees)
	lpEarnings.Sub(lpEarnings, totalPenalty)
	data.LpEarnings = lpEarnings
	return data
}

func (u *SummariesUseCase) aggregatePeginData(ctx context.Context, startDate, endDate time.Time) (SummaryData, error) {
	var (
		quotes         []quote.PeginQuote
		retainedQuotes []quote.RetainedPeginQuote
		err            error
	)
	quotes, retainedQuotes, err = u.peginRepo.ListQuotesByDateRange(ctx, startDate, endDate)
	if err != nil {
		return NewSummaryData(), usecases.WrapUseCaseError(usecases.SummariesUseCaseId, err)
	}
	data := u.processPeginQuoteData(ctx, retainedQuotes)
	data.TotalQuotesCount = int64(len(quotes))
	return data, nil
}

func (u *SummariesUseCase) aggregatePegoutData(ctx context.Context, startDate, endDate time.Time) (SummaryData, error) {
	var (
		quotes         []quote.PegoutQuote
		retainedQuotes []quote.RetainedPegoutQuote
		err            error
	)
	quotes, retainedQuotes, err = u.pegoutRepo.ListQuotesByDateRange(ctx, startDate, endDate)
	if err != nil {
		return NewSummaryData(), usecases.WrapUseCaseError(usecases.SummariesUseCaseId, err)
	}
	data := u.processPegoutQuoteData(ctx, retainedQuotes)
	data.TotalQuotesCount = int64(len(quotes))
	return data, nil
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
