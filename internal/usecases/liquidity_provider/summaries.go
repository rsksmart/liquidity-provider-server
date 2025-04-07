package liquidity_provider

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
)

type SummariesResponse struct {
	PeginSummary  SummaryData `json:"peginSummary"`
	PegoutSummary SummaryData `json:"pegoutSummary"`
}

type SummaryData struct {
	TotalQuotesCount          int64  `json:"totalQuotesCount"`
	AcceptedQuotesCount       int64  `json:"acceptedQuotesCount"`
	TotalQuotedAmount         string `json:"totalQuotedAmount"`
	TotalAcceptedQuotedAmount string `json:"totalAcceptedQuotedAmount"`
	TotalFeesCollected        string `json:"totalFeesCollected"`
	RefundedQuotesCount       int64  `json:"refundedQuotesCount"`
	TotalPenaltyAmount        string `json:"totalPenaltyAmount"`
	LpEarnings                string `json:"lpEarnings"`
}

type SummariesUseCase struct {
	peginRepo  quote.PeginQuoteRepository
	pegoutRepo quote.PegoutQuoteRepository
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

func (u *SummariesUseCase) Run(ctx context.Context, startDate, endDate time.Time) (SummariesResponse, error) {
	log.Printf("Running summaries from %s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	peginData, err := u.aggregatePeginData(ctx, startDate, endDate)
	if err != nil {
		log.Printf("Error aggregating pegin data: %v", err)
		return SummariesResponse{}, err
	}

	pegoutData, err := u.aggregatePegoutData(ctx, startDate, endDate)
	if err != nil {
		log.Printf("Error aggregating pegout data: %v", err)
		return SummariesResponse{}, err
	}

	return SummariesResponse{
		PeginSummary:  peginData,
		PegoutSummary: pegoutData,
	}, nil
}

func (u *SummariesUseCase) aggregatePeginData(ctx context.Context, startDate, endDate time.Time) (SummaryData, error) { //nolint:funlen,cyclop
	var data SummaryData
	quotes, retainedQuotes, err := u.peginRepo.ListQuotesByDateRange(ctx, startDate, endDate)
	if err != nil {
		log.Printf("Error listing pegin quotes: %v", err)
		return data, err
	}
	var totalAmount = entities.NewWei(0)
	var acceptedTotalAmount = entities.NewWei(0)
	var totalFees = entities.NewWei(0)
	var callFees = entities.NewWei(0)
	var totalPenalty = entities.NewWei(0)
	
	processedQuotes := make(map[string]*quote.PeginQuote)
	uniqueQuoteHashes := make(map[string]bool)
	
	for i := range quotes {
		quoteHash := strconv.FormatInt(quotes[i].Nonce, 10)
		uniqueQuoteHashes[quoteHash] = true
		totalAmount.Add(totalAmount, quotes[i].Total())		
		quoteCopy := quotes[i]
		processedQuotes[quoteHash] = &quoteCopy
	}
	
	if len(retainedQuotes) > 0 { //nolint:nestif
		data.TotalQuotesCount = int64(len(retainedQuotes))		
		var hashesToFetch []string
		for _, retainedQuote := range retainedQuotes {
			uniqueQuoteHashes[retainedQuote.QuoteHash] = true			
			if _, exists := processedQuotes[retainedQuote.QuoteHash]; !exists {
				hashesToFetch = append(hashesToFetch, retainedQuote.QuoteHash)
			}
		}
		for _, hash := range hashesToFetch {
			quoteObj, err := u.peginRepo.GetQuote(ctx, hash)
			if err != nil {
				log.Printf("Error getting pegin quote %s: %v", hash, err)
				continue
			}
			if quoteObj == nil {
				log.Printf("Pegin quote not found for hash %s", hash)
				continue
			}
			processedQuotes[hash] = quoteObj
			totalAmount.Add(totalAmount, quoteObj.Total())
		}
		
		for _, retainedQuote := range retainedQuotes {
			quoteObj, exists := processedQuotes[retainedQuote.QuoteHash]
			if !exists {
				continue
			}
			
			accepted := isAcceptedPegin(*quoteObj, retainedQuote)
			if accepted {
				data.AcceptedQuotesCount++
				acceptedTotalAmount.Add(acceptedTotalAmount, quoteObj.Total())
				callFees.Add(callFees, quoteObj.CallFee)				
				totalFees.Add(totalFees, quoteObj.CallFee)
				totalFees.Add(totalFees, quoteObj.GasFee)
				totalFees.Add(totalFees, entities.NewUWei(quoteObj.ProductFeeAmount))
			}
			
			refunded := isRefundedPegin(*quoteObj, retainedQuote)
			if refunded {
				data.RefundedQuotesCount++
				totalPenalty.Add(totalPenalty, quoteObj.PenaltyFee)
			}
		}
	} else {
		data.TotalQuotesCount = int64(len(quotes))
	}

	lpEarnings := new(entities.Wei)
	lpEarnings.Sub(callFees, totalPenalty)	
	data.TotalQuotedAmount = totalAmount.String()
	data.TotalAcceptedQuotedAmount = acceptedTotalAmount.String()
	data.TotalFeesCollected = totalFees.String()
	data.TotalPenaltyAmount = totalPenalty.String()
	data.LpEarnings = lpEarnings.String()
	return data, nil
}

func (u *SummariesUseCase) aggregatePegoutData(ctx context.Context, startDate, endDate time.Time) (SummaryData, error) { //nolint:funlen,cyclop
	var data SummaryData
	quotes, retainedQuotes, err := u.pegoutRepo.ListQuotesByDateRange(ctx, startDate, endDate)
	if err != nil {
		log.Printf("Error listing pegout quotes: %v", err)
		return data, err
	}
	
	var totalAmount = entities.NewWei(0)
	var acceptedTotalAmount = entities.NewWei(0)
	var totalFees = entities.NewWei(0)
	var callFees = entities.NewWei(0) 
	var totalPenalty = entities.NewWei(0)
	
	processedQuotes := make(map[string]*quote.PegoutQuote)
	uniqueQuoteHashes := make(map[string]bool)
	
	for i := range quotes {
		quoteHash := strconv.FormatInt(quotes[i].Nonce, 10)
		uniqueQuoteHashes[quoteHash] = true
		totalAmount.Add(totalAmount, quotes[i].Total())
		quoteCopy := quotes[i]
		processedQuotes[quoteHash] = &quoteCopy
	}
	
	if len(retainedQuotes) > 0 { //nolint:nestif
		data.TotalQuotesCount = int64(len(retainedQuotes))		
		var hashesToFetch []string
		for _, retainedQuote := range retainedQuotes {
			uniqueQuoteHashes[retainedQuote.QuoteHash] = true			
			if _, exists := processedQuotes[retainedQuote.QuoteHash]; !exists {
				hashesToFetch = append(hashesToFetch, retainedQuote.QuoteHash)
			}
		}
		for _, hash := range hashesToFetch {
			quoteObj, err := u.pegoutRepo.GetQuote(ctx, hash)
			if err != nil {
				log.Printf("Error getting pegout quote %s: %v", hash, err)
				continue
			}
			if quoteObj == nil {
				log.Printf("Pegout quote not found for hash %s", hash)
				continue
			}
			
			processedQuotes[hash] = quoteObj
			totalAmount.Add(totalAmount, quoteObj.Total())
		}
		
		for _, retainedQuote := range retainedQuotes {
			quoteObj, exists := processedQuotes[retainedQuote.QuoteHash]
			if !exists {
				continue
			}
			accepted := isAcceptedPegout(*quoteObj, retainedQuote)
			if accepted {
				data.AcceptedQuotesCount++
				acceptedTotalAmount.Add(acceptedTotalAmount, quoteObj.Total())
				callFees.Add(callFees, quoteObj.CallFee)				
				totalFees.Add(totalFees, quoteObj.CallFee)
				totalFees.Add(totalFees, quoteObj.GasFee)
				totalFees.Add(totalFees, entities.NewUWei(quoteObj.ProductFeeAmount))
			}			
			refunded := isRefundedPegout(*quoteObj, retainedQuote)
			if refunded {
				data.RefundedQuotesCount++
				totalPenalty.Add(totalPenalty, entities.NewUWei(quoteObj.PenaltyFee))
			}
		}
	} else {
		data.TotalQuotesCount = int64(len(quotes))
	}

	lpEarnings := new(entities.Wei)
	lpEarnings.Sub(callFees, totalPenalty)	
	data.TotalQuotedAmount = totalAmount.String()
	data.TotalAcceptedQuotedAmount = acceptedTotalAmount.String()
	data.TotalFeesCollected = totalFees.String()
	data.TotalPenaltyAmount = totalPenalty.String()
	data.LpEarnings = lpEarnings.String()
	return data, nil
}

func isAcceptedPegin(_ quote.PeginQuote, retained quote.RetainedPeginQuote) bool {
	return retained.State == quote.PeginStateCallForUserSucceeded
}

func isRefundedPegin(_ quote.PeginQuote, retained quote.RetainedPeginQuote) bool {
	return retained.State == quote.PeginStateCallForUserFailed
}

func isAcceptedPegout(_ quote.PegoutQuote, retained quote.RetainedPegoutQuote) bool {
	return retained.State == quote.PegoutStateBridgeTxSucceeded
}

func isRefundedPegout(_ quote.PegoutQuote, retained quote.RetainedPegoutQuote) bool {
	return retained.State == quote.PegoutStateBridgeTxFailed
}
