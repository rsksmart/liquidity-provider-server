package liquidity_provider

import (
	"context"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
)

// Time format constants
const (
	DateFormat = time.DateOnly
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
	log.Debugf("Running summaries from %s to %s", startDate.Format(DateFormat), endDate.Format(DateFormat))

	peginData, err := u.aggregatePeginData(ctx, startDate, endDate)
	if err != nil {
		log.Errorf("Error aggregating pegin data: %v", err)
		return SummariesResponse{}, err
	}

	pegoutData, err := u.aggregatePegoutData(ctx, startDate, endDate)
	if err != nil {
		log.Errorf("Error aggregating pegout data: %v", err)
		return SummariesResponse{}, err
	}

	return SummariesResponse{
		PeginSummary:  peginData,
		PegoutSummary: pegoutData,
	}, nil
}

func (u *SummariesUseCase) aggregatePeginData(ctx context.Context, startDate, endDate time.Time) (SummaryData, error) {
	// Initialize data with default zero values for all fields
	data := SummaryData{
		TotalQuotesCount:          0,
		AcceptedQuotesCount:       0,
		TotalQuotedAmount:         "0",
		TotalAcceptedQuotedAmount: "0",
		TotalFeesCollected:        "0",
		RefundedQuotesCount:       0,
		TotalPenaltyAmount:        "0",
		LpEarnings:                "0",
	}

	result, err := u.peginRepo.ListQuotesByDateRange(ctx, startDate, endDate)
	if err != nil {
		log.Errorf("Error listing pegin quotes: %v", err)
		return data, err
	}

	return u.processPeginData(ctx, result.Quotes, result.RetainedQuotes), nil
}

// processPeginData processes the pegin quotes and retained quotes to calculate metrics
func (u *SummariesUseCase) processPeginData(ctx context.Context, quotes []quote.PeginQuote, retainedQuotes []quote.RetainedPeginQuote) SummaryData {
	// Create data structure with default values
	data := SummaryData{
		TotalQuotesCount:          0,
		AcceptedQuotesCount:       0,
		TotalQuotedAmount:         "0",
		TotalAcceptedQuotedAmount: "0",
		TotalFeesCollected:        "0",
		RefundedQuotesCount:       0,
		TotalPenaltyAmount:        "0",
		LpEarnings:                "0",
	}

	// Calculate base metrics from quotes
	totalAmount := calculateTotalAmount(quotes)

	// If no retained quotes, just return count of quotes and total amount
	if len(retainedQuotes) == 0 {
		data.TotalQuotesCount = int64(len(quotes))
		data.TotalQuotedAmount = totalAmount.String()
		return data
	}

	// Set total quotes count to retained quotes count
	data.TotalQuotesCount = int64(len(retainedQuotes))

	// Create a map for lookup of quotes by their hash
	quotesByHash := createQuoteHashMap(quotes)

	// Fetch any missing quotes
	u.fetchMissingPeginQuotes(ctx, quotesByHash, retainedQuotes, totalAmount)

	// Calculate metrics based on retained quotes and their quotes
	data.TotalQuotedAmount = totalAmount.String()

	// Calculate accepted and refunded metrics
	acceptedTotalAmount := entities.NewWei(0)
	totalFees := entities.NewWei(0)
	callFees := entities.NewWei(0)
	totalPenalty := entities.NewWei(0)

	for _, retainedQuote := range retainedQuotes {
		quoteObj, exists := quotesByHash[retainedQuote.QuoteHash]
		if !exists {
			continue
		}

		// Process accepted quotes
		if isAcceptedPegin(retainedQuote) {
			data.AcceptedQuotesCount++
			acceptedTotalAmount.Add(acceptedTotalAmount, quoteObj.Total())

			// Add fees
			callFees.Add(callFees, quoteObj.CallFee)
			totalFees.Add(totalFees, quoteObj.CallFee)
			totalFees.Add(totalFees, quoteObj.GasFee)
			totalFees.Add(totalFees, entities.NewUWei(quoteObj.ProductFeeAmount))
		}

		// Process refunded quotes
		if isRefundedPegin(retainedQuote) {
			data.RefundedQuotesCount++
			totalPenalty.Add(totalPenalty, quoteObj.PenaltyFee)
		}
	}

	// Calculate LP earnings: callFees - penalties
	lpEarnings := new(entities.Wei)
	lpEarnings.Sub(callFees, totalPenalty)

	// Set remaining data fields
	data.TotalAcceptedQuotedAmount = acceptedTotalAmount.String()
	data.TotalFeesCollected = totalFees.String()
	data.TotalPenaltyAmount = totalPenalty.String()
	data.LpEarnings = lpEarnings.String()

	return data
}

func calculateTotalAmount(quotes []quote.PeginQuote) *entities.Wei {
	totalAmount := entities.NewWei(0)
	for i := range quotes {
		totalAmount.Add(totalAmount, quotes[i].Total())
	}
	return totalAmount
}

// createQuoteHashMap creates a map of quotes indexed by quote hash
func createQuoteHashMap(quotes []quote.PeginQuote) map[string]*quote.PeginQuote {
	// Create a map of quotes by their hash or unique identifier
	quotesByHash := make(map[string]*quote.PeginQuote, len(quotes))

	// Store quotes by their nonce as a unique identifier
	for i := range quotes {
		quoteCopy := quotes[i]
		// Using Nonce as a unique identifier since Quote doesn't appear to have a Hash field
		nonceStr := strconv.FormatInt(quoteCopy.Nonce, 10)
		quotesByHash[nonceStr] = &quoteCopy
	}

	return quotesByHash
}

// fetchMissingPeginQuotes fetches quotes that are needed for calculations
func (u *SummariesUseCase) fetchMissingPeginQuotes(
	ctx context.Context,
	quotesByHash map[string]*quote.PeginQuote,
	retainedQuotes []quote.RetainedPeginQuote,
	totalAmount *entities.Wei,
) {
	// For each retained quote, fetch the quote if not present
	for _, retainedQuote := range retainedQuotes {
		quoteHash := retainedQuote.QuoteHash

		// Skip if already in map
		if _, exists := quotesByHash[quoteHash]; exists {
			continue
		}

		// Fetch the quote from repository
		quoteObj, err := u.peginRepo.GetQuote(ctx, quoteHash)
		if err != nil {
			log.Errorf("Error getting pegin quote %s: %v", quoteHash, err)
			continue
		}

		if quoteObj == nil {
			log.Debugf("Pegin quote not found for hash %s", quoteHash)
			continue
		}

		// Add to map and update total amount
		quotesByHash[quoteHash] = quoteObj
		totalAmount.Add(totalAmount, quoteObj.Total())
	}
}

// isAcceptedPegin checks if a retained pegin quote has been accepted
func isAcceptedPegin(retained quote.RetainedPeginQuote) bool {
	return retained.State == quote.PeginStateCallForUserSucceeded
}

// isRefundedPegin checks if a retained pegin quote has been refunded
func isRefundedPegin(retained quote.RetainedPeginQuote) bool {
	return retained.State == quote.PeginStateCallForUserFailed
}

// aggregatePegoutData aggregates the pegout data for the specified date range
func (u *SummariesUseCase) aggregatePegoutData(ctx context.Context, startDate, endDate time.Time) (SummaryData, error) {
	// Initialize data with default zero values for all fields
	data := SummaryData{
		TotalQuotesCount:          0,
		AcceptedQuotesCount:       0,
		TotalQuotedAmount:         "0",
		TotalAcceptedQuotedAmount: "0",
		TotalFeesCollected:        "0",
		RefundedQuotesCount:       0,
		TotalPenaltyAmount:        "0",
		LpEarnings:                "0",
	}

	result, err := u.pegoutRepo.ListQuotesByDateRange(ctx, startDate, endDate)
	if err != nil {
		log.Errorf("Error listing pegout quotes: %v", err)
		return data, err
	}

	return u.processPegoutData(ctx, result.Quotes, result.RetainedQuotes), nil
}

// processPegoutData processes the pegout quotes and retained quotes to calculate metrics
func (u *SummariesUseCase) processPegoutData(ctx context.Context, quotes []quote.PegoutQuote, retainedQuotes []quote.RetainedPegoutQuote) SummaryData {
	// Create data structure with default values
	data := SummaryData{
		TotalQuotesCount:          0,
		AcceptedQuotesCount:       0,
		TotalQuotedAmount:         "0",
		TotalAcceptedQuotedAmount: "0",
		TotalFeesCollected:        "0",
		RefundedQuotesCount:       0,
		TotalPenaltyAmount:        "0",
		LpEarnings:                "0",
	}

	// Calculate base metrics from quotes
	totalAmount := calculatePegoutTotalAmount(quotes)

	// If no retained quotes, just return count of quotes and total amount
	if len(retainedQuotes) == 0 {
		data.TotalQuotesCount = int64(len(quotes))
		data.TotalQuotedAmount = totalAmount.String()
		return data
	}

	// Set total quotes count to retained quotes count
	data.TotalQuotesCount = int64(len(retainedQuotes))

	// Create a map for lookup of quotes by their hash
	quotesByHash := createPegoutQuoteHashMap(quotes)

	// Fetch any missing quotes
	u.fetchMissingPegoutQuotes(ctx, quotesByHash, retainedQuotes, totalAmount)

	// Calculate metrics based on retained quotes and their quotes
	data.TotalQuotedAmount = totalAmount.String()

	// Calculate accepted and refunded metrics
	acceptedTotalAmount := entities.NewWei(0)
	totalFees := entities.NewWei(0)
	callFees := entities.NewWei(0)
	totalPenalty := entities.NewWei(0)

	for _, retainedQuote := range retainedQuotes {
		quoteObj, exists := quotesByHash[retainedQuote.QuoteHash]
		if !exists {
			continue
		}

		// Process accepted quotes
		if isAcceptedPegout(retainedQuote) {
			data.AcceptedQuotesCount++
			acceptedTotalAmount.Add(acceptedTotalAmount, quoteObj.Total())

			// Add fees
			callFees.Add(callFees, quoteObj.CallFee)
			totalFees.Add(totalFees, quoteObj.CallFee)
			totalFees.Add(totalFees, quoteObj.GasFee)
			totalFees.Add(totalFees, entities.NewUWei(quoteObj.ProductFeeAmount))
		}

		// Process refunded quotes
		if isRefundedPegout(retainedQuote) {
			data.RefundedQuotesCount++
			// Convert uint64 to Wei for PegoutQuote if needed
			totalPenalty.Add(totalPenalty, entities.NewUWei(quoteObj.PenaltyFee))
		}
	}

	// Calculate LP earnings: callFees - penalties
	lpEarnings := new(entities.Wei)
	lpEarnings.Sub(callFees, totalPenalty)

	// Set remaining data fields
	data.TotalAcceptedQuotedAmount = acceptedTotalAmount.String()
	data.TotalFeesCollected = totalFees.String()
	data.TotalPenaltyAmount = totalPenalty.String()
	data.LpEarnings = lpEarnings.String()

	return data
}

func calculatePegoutTotalAmount(quotes []quote.PegoutQuote) *entities.Wei {
	totalAmount := entities.NewWei(0)
	for i := range quotes {
		totalAmount.Add(totalAmount, quotes[i].Total())
	}
	return totalAmount
}

// createPegoutQuoteHashMap creates a map of pegout quotes indexed by quote hash
func createPegoutQuoteHashMap(quotes []quote.PegoutQuote) map[string]*quote.PegoutQuote {
	// Create a map of quotes by their hash or unique identifier
	quotesByHash := make(map[string]*quote.PegoutQuote, len(quotes))

	// Store quotes by their nonce as a unique identifier
	for i := range quotes {
		quoteCopy := quotes[i]
		// Using Nonce as a unique identifier
		nonceStr := strconv.FormatInt(quoteCopy.Nonce, 10)
		quotesByHash[nonceStr] = &quoteCopy
	}

	return quotesByHash
}

// fetchMissingPegoutQuotes fetches quotes that are needed for calculations
func (u *SummariesUseCase) fetchMissingPegoutQuotes(
	ctx context.Context,
	quotesByHash map[string]*quote.PegoutQuote,
	retainedQuotes []quote.RetainedPegoutQuote,
	totalAmount *entities.Wei,
) {
	// For each retained quote, fetch the quote if not present
	for _, retainedQuote := range retainedQuotes {
		quoteHash := retainedQuote.QuoteHash

		// Skip if already in map
		if _, exists := quotesByHash[quoteHash]; exists {
			continue
		}

		// Fetch the quote from repository
		quoteObj, err := u.pegoutRepo.GetQuote(ctx, quoteHash)
		if err != nil {
			log.Errorf("Error getting pegout quote %s: %v", quoteHash, err)
			continue
		}

		if quoteObj == nil {
			log.Debugf("Pegout quote not found for hash %s", quoteHash)
			continue
		}

		// Add to map and update total amount
		quotesByHash[quoteHash] = quoteObj
		totalAmount.Add(totalAmount, quoteObj.Total())
	}
}

// isAcceptedPegout checks if a retained pegout quote has been accepted
func isAcceptedPegout(retained quote.RetainedPegoutQuote) bool {
	return retained.State == quote.PegoutStateSendPegoutSucceeded
}

// isRefundedPegout checks if a retained pegout quote has been refunded
func isRefundedPegout(retained quote.RetainedPegoutQuote) bool {
	return retained.State == quote.PegoutStateSendPegoutFailed
}
