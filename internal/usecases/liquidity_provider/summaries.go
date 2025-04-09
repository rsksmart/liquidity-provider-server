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
	TotalAcceptedQuotesCount  int64         `json:"totalAcceptedQuotesCount"`
	ConfirmedQuotesCount      int64         `json:"confirmedQuotesCount"`
	TotalQuotedAmount         *entities.Wei `json:"totalQuotedAmount"`
	TotalAcceptedQuotedAmount *entities.Wei `json:"totalAcceptedQuotedAmount"`
	TotalFeesCollected        *entities.Wei `json:"totalFeesCollected"`
	RefundedQuotesCount       int64         `json:"refundedQuotesCount"`
	TotalPenaltyAmount        *entities.Wei `json:"totalPenaltyAmount"`
	LpEarnings                *entities.Wei `json:"lpEarnings"`
}

type feeAdapter struct {
	callFee    *entities.Wei
	gasFee     *entities.Wei
	productFee *entities.Wei
	penaltyFee *entities.Wei
}

type quoteResultAdapter[Q any, R quote.RetainedQuote] struct {
	quotes           []Q
	retainedQuotes   []R
	quoteHashToIndex map[string]int
}

type SummariesUseCase struct {
	peginRepo  quote.PeginQuoteRepository
	pegoutRepo quote.PegoutQuoteRepository
}

func NewSummaryData() SummaryData {
	return SummaryData{
		TotalAcceptedQuotesCount:  0,
		ConfirmedQuotesCount:      0,
		TotalQuotedAmount:         entities.NewWei(0),
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
) *SummariesUseCase {
	return &SummariesUseCase{
		peginRepo:  peginRepo,
		pegoutRepo: pegoutRepo,
	}
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
	return SummaryResult{
		PeginSummary:  peginData,
		PegoutSummary: pegoutData,
	}, nil
}

func mapPeginQuote(q *quote.PeginQuote) quote.FeeProvider {
	if q == nil {
		return &feeAdapter{}
	}
	return &feeAdapter{
		callFee:    q.CallFee,
		gasFee:     q.GasFee,
		productFee: entities.NewUWei(q.ProductFeeAmount),
		penaltyFee: q.PenaltyFee,
	}
}

func mapPegoutQuote(q *quote.PegoutQuote) quote.FeeProvider {
	if q == nil {
		return &feeAdapter{}
	}
	return &feeAdapter{
		callFee:    q.CallFee,
		gasFee:     q.GasFee,
		productFee: entities.NewUWei(q.ProductFeeAmount),
		penaltyFee: entities.NewUWei(q.PenaltyFee),
	}
}

func processQuoteData[Q any, R quote.RetainedQuote, F quote.FeeProvider](
	ctx context.Context,
	quotes []Q,
	retainedQuotes []R,
	quoteHashToIndex map[string]int,
	getQuote func(context.Context, string) (*Q, error),
	isPaid func(R) bool,
	isRefunded func(R) bool,
	feeProvider func(*Q) F,
) SummaryData {
	data := NewSummaryData()
	totalAmount := calculateTotalAmount(quotes)
	if len(retainedQuotes) == 0 {
		data.TotalAcceptedQuotesCount = int64(len(quotes))
		data.TotalQuotedAmount = totalAmount
		return data
	}
	data.TotalAcceptedQuotesCount = int64(len(retainedQuotes))
	quotesByHash := createQuoteHashMap(quotes, quoteHashToIndex)
	fetchMissingQuotes(ctx, quotesByHash, retainedQuotes, totalAmount, getQuote)
	data.TotalQuotedAmount = totalAmount
	acceptedTotalAmount := entities.NewWei(0)
	totalFees := entities.NewWei(0)
	callFees := entities.NewWei(0)
	totalPenalty := entities.NewWei(0)
	for _, retained := range retainedQuotes {
		quoteHash := retained.GetQuoteHash()
		quoteObj, exists := quotesByHash[quoteHash]
		if !exists {
			continue
		}
		fees := feeProvider(quoteObj)
		data.ConfirmedQuotesCount++
		if q, ok := any(quoteObj).(quote.Quote); ok {
			acceptedTotalAmount.Add(acceptedTotalAmount, q.Total())
		}
		if isPaid(retained) {
			callFee := fees.GetCallFee()
			callFees.Add(callFees, callFee)
			totalFees.Add(totalFees, callFee)
			totalFees.Add(totalFees, fees.GetGasFee())
			totalFees.Add(totalFees, fees.GetProductFee())
		}
		if isRefunded(retained) {
			data.RefundedQuotesCount++
			totalPenalty.Add(totalPenalty, fees.GetPenaltyFee())
		}
	}
	lpEarnings := new(entities.Wei)
	lpEarnings.Sub(callFees, totalPenalty)
	data.TotalAcceptedQuotedAmount = acceptedTotalAmount
	data.TotalFeesCollected = totalFees
	data.TotalPenaltyAmount = totalPenalty
	data.LpEarnings = lpEarnings
	return data
}

func calculateTotalAmount[T any](quotes []T) *entities.Wei {
	totalAmount := entities.NewWei(0)
	for i := range quotes {
		var total *entities.Wei
		if q, ok := any(&quotes[i]).(quote.Quote); ok {
			total = q.Total()
		}
		if total != nil {
			totalAmount.Add(totalAmount, total)
		}
	}
	return totalAmount
}

func createQuoteHashMap[T any](quotes []T, quoteHashToIndex map[string]int) map[string]*T {
	quotesByHash := make(map[string]*T, len(quoteHashToIndex))
	for hash, index := range quoteHashToIndex {
		if index >= 0 && index < len(quotes) {
			quoteCopy := quotes[index]
			quotesByHash[hash] = &quoteCopy
		}
	}
	return quotesByHash
}

func fetchMissingQuotes[Q any, R quote.RetainedQuote](
	ctx context.Context,
	quotesByHash map[string]*Q,
	retainedQuotes []R,
	totalAmount *entities.Wei,
	getQuote func(context.Context, string) (*Q, error),
) {
	for _, retained := range retainedQuotes {
		quoteHash := retained.GetQuoteHash()
		if _, exists := quotesByHash[quoteHash]; exists {
			continue
		}
		quoteObj, err := getQuote(ctx, quoteHash)
		if err != nil {
			log.Errorf("Error getting quote %s: %v", quoteHash, err)
			continue
		}
		if quoteObj == nil {
			log.Debugf("Quote not found for hash %s", quoteHash)
			continue
		}
		quotesByHash[quoteHash] = quoteObj
		if q, ok := any(quoteObj).(quote.Quote); ok {
			totalAmount.Add(totalAmount, q.Total())
		}
	}
}

func adaptPeginResult(result quote.PeginQuoteResult) quote.QuoteResult[quote.PeginQuote, quote.RetainedPeginQuote] {
	return quoteResultAdapter[quote.PeginQuote, quote.RetainedPeginQuote]{
		quotes:           result.Quotes,
		retainedQuotes:   result.RetainedQuotes,
		quoteHashToIndex: result.QuoteHashToIndex,
	}
}

func adaptPegoutResult(result quote.PegoutQuoteResult) quote.QuoteResult[quote.PegoutQuote, quote.RetainedPegoutQuote] {
	return quoteResultAdapter[quote.PegoutQuote, quote.RetainedPegoutQuote]{
		quotes:           result.Quotes,
		retainedQuotes:   result.RetainedQuotes,
		quoteHashToIndex: result.QuoteHashToIndex,
	}
}

func aggregateData[Q any, RQ quote.RetainedQuote](
	ctx context.Context,
	startDate, endDate time.Time,
	listQuotes func(context.Context, time.Time, time.Time) (quote.QuoteResult[Q, RQ], error),
	getQuote func(context.Context, string) (*Q, error),
	isPaid func(RQ) bool,
	isRefunded func(RQ) bool,
	toFeeProvider func(*Q) quote.FeeProvider,
) (SummaryData, error) {
	result, err := listQuotes(ctx, startDate, endDate)
	if err != nil {
		log.Errorf("Error listing quotes: %v", err)
		return NewSummaryData(), err
	}
	return processQuoteData(
		ctx,
		result.GetQuotes(),
		result.GetRetainedQuotes(),
		result.GetQuoteHashToIndex(),
		getQuote,
		isPaid,
		isRefunded,
		toFeeProvider,
	), nil
}

func (u *SummariesUseCase) aggregatePeginData(ctx context.Context, startDate, endDate time.Time) (SummaryData, error) {
	listPeginQuotes := func(ctx context.Context, start, end time.Time) (quote.QuoteResult[quote.PeginQuote, quote.RetainedPeginQuote], error) {
		result, err := u.peginRepo.ListQuotesByDateRange(ctx, start, end)
		if err != nil {
			return nil, err
		}
		return adaptPeginResult(result), nil
	}
	return aggregateData(
		ctx, startDate, endDate,
		listPeginQuotes,
		u.peginRepo.GetQuote,
		isPaidPegin,
		isRefundedPegin,
		mapPeginQuote,
	)
}

func (u *SummariesUseCase) aggregatePegoutData(ctx context.Context, startDate, endDate time.Time) (SummaryData, error) {
	listPegoutQuotes := func(ctx context.Context, start, end time.Time) (quote.QuoteResult[quote.PegoutQuote, quote.RetainedPegoutQuote], error) {
		result, err := u.pegoutRepo.ListQuotesByDateRange(ctx, start, end)
		if err != nil {
			return nil, err
		}
		return adaptPegoutResult(result), nil
	}
	return aggregateData(
		ctx, startDate, endDate,
		listPegoutQuotes,
		u.pegoutRepo.GetQuote,
		isPaidPegout,
		isRefundedPegout,
		mapPegoutQuote,
	)
}

func (a *feeAdapter) GetCallFee() *entities.Wei    { return a.callFee }
func (a *feeAdapter) GetGasFee() *entities.Wei     { return a.gasFee }
func (a *feeAdapter) GetProductFee() *entities.Wei { return a.productFee }
func (a *feeAdapter) GetPenaltyFee() *entities.Wei { return a.penaltyFee }

func isPaidPegout(retained quote.RetainedPegoutQuote) bool {
	return retained.State == quote.PegoutStateSendPegoutSucceeded
}

func isPaidPegin(retained quote.RetainedPeginQuote) bool {
	return retained.State == quote.PeginStateCallForUserSucceeded
}

func isRefundedPegout(retained quote.RetainedPegoutQuote) bool {
	return retained.State == quote.PegoutStateRefundPegOutSucceeded
}

func isRefundedPegin(retained quote.RetainedPeginQuote) bool {
	return retained.State == quote.PeginStateRegisterPegInSucceeded
}

func (a quoteResultAdapter[Q, R]) GetQuotes() []Q {
	return a.quotes
}

func (a quoteResultAdapter[Q, R]) GetRetainedQuotes() []R {
	return a.retainedQuotes
}

func (a quoteResultAdapter[Q, R]) GetQuoteHashToIndex() map[string]int {
	return a.quoteHashToIndex
}
