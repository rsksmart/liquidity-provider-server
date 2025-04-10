package liquidity_provider

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
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
		log.Errorf("Error aggregating pegin data: %v", err)
		return SummaryResult{}, err
	}
	pegoutData, err := u.aggregatePegoutData(ctx, startDate, endDate)
	if err != nil {
		log.Errorf("Error aggregating pegout data: %v", err)
		return SummaryResult{}, err
	}
	return SummaryResult{PeginSummary: peginData, PegoutSummary: pegoutData}, nil
}

func processPeginRetainedQuote(
	ctx context.Context,
	retained quote.RetainedPeginQuote,
	quotesByHash map[string]*quote.PeginQuote,
	getQuote func(context.Context, string) (*quote.PeginQuote, error),
	data *SummaryData,
	totalAmount, acceptedTotalAmount, totalFees, callFees, totalPenalty *entities.Wei,
) {
	hash := retained.QuoteHash
	q, exists := quotesByHash[hash]
	if !exists {
		var err error
		q, err = getQuote(ctx, hash)
		if err != nil || q == nil {
			if err != nil {
				log.Errorf("Error getting quote %s: %v", hash, err)
			}
			return
		}
		quotesByHash[hash] = q
		totalAmount.Add(totalAmount, q.Total())
	}
	callFee, gasFee, productFee, penaltyFee := getPeginFees(q)
	acceptedTotalAmount.Add(acceptedTotalAmount, q.Total())
	if isPeginPaidQuote(retained) {
		data.PaidQuotesCount++
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
	retained quote.RetainedPegoutQuote,
	quotesByHash map[string]*quote.PegoutQuote,
	getQuote func(context.Context, string) (*quote.PegoutQuote, error),
	data *SummaryData,
	totalAmount, acceptedTotalAmount, totalFees, callFees, totalPenalty *entities.Wei,
) {
	hash := retained.QuoteHash
	q, exists := quotesByHash[hash]
	if !exists {
		var err error
		q, err = getQuote(ctx, hash)
		if err != nil || q == nil {
			if err != nil {
				log.Errorf("Error getting quote %s: %v", hash, err)
			}
			return
		}
		quotesByHash[hash] = q
		totalAmount.Add(totalAmount, q.Total())
	}
	callFee, gasFee, productFee, penaltyFee := getPegoutFees(q)
	acceptedTotalAmount.Add(acceptedTotalAmount, q.Total())
	if isPegoutPaidQuote(retained) {
		data.PaidQuotesCount++
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

func processPeginQuoteData(ctx context.Context, quotes []*quote.PeginQuote, retainedQuotes []quote.RetainedPeginQuote, getQuote func(context.Context, string) (*quote.PeginQuote, error)) SummaryData {
	data := NewSummaryData()
	data.TotalQuotesCount = int64(len(quotes))
	totalAmount := entities.NewWei(0)
	for i := range quotes {
		totalAmount.Add(totalAmount, quotes[i].Total())
	}
	data.PaidQuotesAmount = totalAmount
	if len(retainedQuotes) == 0 {
		return data
	}
	data.AcceptedQuotesCount = int64(len(retainedQuotes))
	quotesByHash := make(map[string]*quote.PeginQuote)
	acceptedTotalAmount := entities.NewWei(0)
	totalFees := entities.NewWei(0)
	callFees := entities.NewWei(0)
	totalPenalty := entities.NewWei(0)
	for _, retained := range retainedQuotes {
		processPeginRetainedQuote(ctx, retained, quotesByHash, getQuote, &data,
			totalAmount, acceptedTotalAmount, totalFees, callFees, totalPenalty)
	}
	lpEarnings := new(entities.Wei)
	lpEarnings.Sub(callFees, totalPenalty)
	data.TotalAcceptedQuotedAmount = acceptedTotalAmount
	data.TotalFeesCollected = totalFees
	data.TotalPenaltyAmount = totalPenalty
	data.LpEarnings = lpEarnings
	return data
}

func processPegoutQuoteData(ctx context.Context, quotes []*quote.PegoutQuote, retainedQuotes []quote.RetainedPegoutQuote, getQuote func(context.Context, string) (*quote.PegoutQuote, error)) SummaryData {
	data := NewSummaryData()
	data.TotalQuotesCount = int64(len(quotes))
	totalAmount := entities.NewWei(0)
	for i := range quotes {
		totalAmount.Add(totalAmount, quotes[i].Total())
	}
	data.PaidQuotesAmount = totalAmount
	if len(retainedQuotes) == 0 {
		return data
	}
	data.AcceptedQuotesCount = int64(len(retainedQuotes))
	quotesByHash := make(map[string]*quote.PegoutQuote)
	acceptedTotalAmount := entities.NewWei(0)
	totalFees := entities.NewWei(0)
	callFees := entities.NewWei(0)
	totalPenalty := entities.NewWei(0)
	for _, retained := range retainedQuotes {
		processPegoutRetainedQuote(ctx, retained, quotesByHash, getQuote, &data,
			totalAmount, acceptedTotalAmount, totalFees, callFees, totalPenalty)
	}
	lpEarnings := new(entities.Wei)
	lpEarnings.Sub(callFees, totalPenalty)
	data.TotalAcceptedQuotedAmount = acceptedTotalAmount
	data.TotalFeesCollected = totalFees
	data.TotalPenaltyAmount = totalPenalty
	data.LpEarnings = lpEarnings
	return data
}

func (u *SummariesUseCase) aggregatePeginData(ctx context.Context, startDate, endDate time.Time) (SummaryData, error) {
	var (
		peginResult quote.PeginQuoteResult
		err         error
	)
	peginResult, err = u.peginRepo.ListQuotesByDateRange(ctx, startDate, endDate)
	if err != nil {
		return NewSummaryData(), err
	}
	quotes := make([]*quote.PeginQuote, len(peginResult.Quotes))
	for i := range peginResult.Quotes {
		quotes[i] = &peginResult.Quotes[i]
	}
	retainedQuotes := peginResult.RetainedQuotes
	getQuote := func(ctx context.Context, hash string) (*quote.PeginQuote, error) {
		q, err := u.peginRepo.GetQuote(ctx, hash)
		return q, err
	}
	return processPeginQuoteData(ctx, quotes, retainedQuotes, getQuote), nil
}

func (u *SummariesUseCase) aggregatePegoutData(ctx context.Context, startDate, endDate time.Time) (SummaryData, error) {
	var (
		pegoutResult quote.PegoutQuoteResult
		err          error
	)
	pegoutResult, err = u.pegoutRepo.ListQuotesByDateRange(ctx, startDate, endDate)
	if err != nil {
		return NewSummaryData(), err
	}
	quotes := make([]*quote.PegoutQuote, len(pegoutResult.Quotes))
	for i := range pegoutResult.Quotes {
		quotes[i] = &pegoutResult.Quotes[i]
	}
	retainedQuotes := pegoutResult.RetainedQuotes
	getQuote := func(ctx context.Context, hash string) (*quote.PegoutQuote, error) {
		q, err := u.pegoutRepo.GetQuote(ctx, hash)
		return q, err
	}
	return processPegoutQuoteData(ctx, quotes, retainedQuotes, getQuote), nil
}

func getPeginFees(q *quote.PeginQuote) (callFee, gasFee, productFee, penaltyFee *entities.Wei) {
	callFee, gasFee, productFee, penaltyFee = entities.NewWei(0), entities.NewWei(0), entities.NewWei(0), entities.NewWei(0)
	if q != nil {
		callFee, gasFee = q.CallFee, q.GasFee
		productFee = entities.NewUWei(q.ProductFeeAmount)
		penaltyFee = q.PenaltyFee
	}
	return
}

func getPegoutFees(q *quote.PegoutQuote) (callFee, gasFee, productFee, penaltyFee *entities.Wei) {
	callFee, gasFee, productFee, penaltyFee = entities.NewWei(0), entities.NewWei(0), entities.NewWei(0), entities.NewWei(0)
	if q != nil {
		callFee, gasFee = q.CallFee, q.GasFee
		productFee = entities.NewUWei(q.ProductFeeAmount)
		penaltyFee = entities.NewUWei(q.PenaltyFee)
	}
	return
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
	return retained.State == quote.PegoutStateRefundPegOutSucceeded
}
