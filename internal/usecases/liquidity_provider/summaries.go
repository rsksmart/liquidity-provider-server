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
	peginData, err := u.aggregateData(ctx, startDate, endDate, true)
	if err != nil {
		log.Errorf("Error aggregating data: %v", err)
		return SummaryResult{}, err
	}
	pegoutData, err := u.aggregateData(ctx, startDate, endDate, false)
	if err != nil {
		log.Errorf("Error aggregating data: %v", err)
		return SummaryResult{}, err
	}
	return SummaryResult{PeginSummary: peginData, PegoutSummary: pegoutData}, nil
}

func buildQuotesHashMap(quotes []quote.Quote) map[string]quote.Quote {
	quotesByHash := make(map[string]quote.Quote)
	for i := range quotes {
		if rq, ok := quotes[i].(quote.RetainedQuote); ok {
			quotesByHash[rq.GetQuoteHash()] = quotes[i]
		} else if hasher, ok := quotes[i].(interface{ GetHash() string }); ok {
			quotesByHash[hasher.GetHash()] = quotes[i]
		}
	}
	return quotesByHash
}

func processRetainedQuote(
	ctx context.Context,
	retained quote.RetainedQuote,
	quotesByHash map[string]quote.Quote,
	getQuote func(context.Context, string) (quote.Quote, error),
	data *SummaryData,
	totalAmount, acceptedTotalAmount, totalFees, callFees, totalPenalty *entities.Wei,
) {
	hash := retained.GetQuoteHash()
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
	if q != nil {
		callFee, gasFee, productFee, penaltyFee := getFees(q)
		acceptedTotalAmount.Add(acceptedTotalAmount, q.Total())
		if isPaidQuote(retained) {
			data.PaidQuotesCount++
			callFees.Add(callFees, callFee)
			totalFees.Add(totalFees, callFee)
			totalFees.Add(totalFees, gasFee)
			totalFees.Add(totalFees, productFee)
		}
		if isRefundedQuote(retained) {
			data.RefundedQuotesCount++
			totalPenalty.Add(totalPenalty, penaltyFee)
		}
	}
}

func processQuoteData(ctx context.Context, quotes []quote.Quote, retainedQuotes []quote.RetainedQuote, getQuote func(context.Context, string) (quote.Quote, error)) SummaryData {
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
	quotesByHash := buildQuotesHashMap(quotes)
	acceptedTotalAmount := entities.NewWei(0)
	totalFees := entities.NewWei(0)
	callFees := entities.NewWei(0)
	totalPenalty := entities.NewWei(0)
	for _, retained := range retainedQuotes {
		processRetainedQuote(ctx, retained, quotesByHash, getQuote, &data,
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

func (u *SummariesUseCase) aggregateData(ctx context.Context, startDate, endDate time.Time, isPegin bool) (SummaryData, error) {
	var quotes []quote.Quote
	var retainedQuotes []quote.RetainedQuote
	var getQuote func(context.Context, string) (quote.Quote, error)
	if isPegin {
		peginResult, err := u.peginRepo.ListQuotesByDateRange(ctx, startDate, endDate)
		if err != nil {
			log.Errorf("Error listing quotes: %v", err)
			return NewSummaryData(), err
		}
		for i := range peginResult.Quotes {
			quotes = append(quotes, &peginResult.Quotes[i])
		}
		for _, rq := range peginResult.RetainedQuotes {
			retainedQuotes = append(retainedQuotes, rq)
		}
		getQuote = func(ctx context.Context, hash string) (quote.Quote, error) {
			return u.peginRepo.GetQuote(ctx, hash)
		}
	} else {
		pegoutResult, err := u.pegoutRepo.ListQuotesByDateRange(ctx, startDate, endDate)
		if err != nil {
			log.Errorf("Error listing quotes: %v", err)
			return NewSummaryData(), err
		}
		for i := range pegoutResult.Quotes {
			quotes = append(quotes, &pegoutResult.Quotes[i])
		}
		for _, rq := range pegoutResult.RetainedQuotes {
			retainedQuotes = append(retainedQuotes, rq)
		}
		getQuote = func(ctx context.Context, hash string) (quote.Quote, error) {
			return u.pegoutRepo.GetQuote(ctx, hash)
		}
	}
	return processQuoteData(ctx, quotes, retainedQuotes, getQuote), nil
}

func getFees(q quote.Quote) (callFee, gasFee, productFee, penaltyFee *entities.Wei) {
	callFee, gasFee, productFee, penaltyFee = entities.NewWei(0), entities.NewWei(0), entities.NewWei(0), entities.NewWei(0)
	switch v := q.(type) {
	case *quote.PeginQuote:
		if v != nil {
			callFee, gasFee = v.CallFee, v.GasFee
			productFee = entities.NewUWei(v.ProductFeeAmount)
			penaltyFee = v.PenaltyFee
		}
	case *quote.PegoutQuote:
		if v != nil {
			callFee, gasFee = v.CallFee, v.GasFee
			productFee = entities.NewUWei(v.ProductFeeAmount)
			penaltyFee = entities.NewUWei(v.PenaltyFee)
		}
	}
	return
}

func isPaidQuote(retained quote.RetainedQuote) bool {
	switch r := retained.(type) {
	case quote.RetainedPeginQuote:
		return r.State == quote.PeginStateCallForUserSucceeded
	case quote.RetainedPegoutQuote:
		return r.State == quote.PegoutStateSendPegoutSucceeded
	}
	return false
}

func isRefundedQuote(retained quote.RetainedQuote) bool {
	switch r := retained.(type) {
	case quote.RetainedPeginQuote:
		return r.State == quote.PeginStateRegisterPegInSucceeded
	case quote.RetainedPegoutQuote:
		return r.State == quote.PegoutStateRefundPegOutSucceeded
	}
	return false
}
