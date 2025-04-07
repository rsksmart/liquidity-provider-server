package liquidity_provider

import (
	"context"
	"log"
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

func (u *SummariesUseCase) aggregatePeginData(ctx context.Context, startDate, endDate time.Time) (SummaryData, error) {
	var data SummaryData
	quotes, retainedQuotes, err := u.peginRepo.ListQuotesByDateRange(ctx, startDate, endDate)
	if err != nil {
		log.Printf("Error listing pegin quotes: %v", err)
		return data, err
	}
	data.TotalQuotesCount = int64(len(quotes) + len(retainedQuotes))
	var totalAmount = entities.NewWei(0)
	var acceptedTotalAmount = entities.NewWei(0)
	var totalFees = entities.NewWei(0)
	var callFees = entities.NewWei(0)
	var totalPenalty = entities.NewWei(0)
	for _, retainedQuote := range retainedQuotes {
		quote, err := u.peginRepo.GetQuote(ctx, retainedQuote.QuoteHash)
		if err != nil {
			log.Printf("Error getting pegin quote %s: %v", retainedQuote.QuoteHash, err)
			continue
		}
		if quote == nil {
			log.Printf("Pegin quote not found for hash %s", retainedQuote.QuoteHash)
			continue
		}
		totalAmount.Add(totalAmount, quote.Total())
		accepted := isAcceptedPegin(*quote, retainedQuote)
		if accepted {
			data.AcceptedQuotesCount++
			acceptedTotalAmount.Add(acceptedTotalAmount, quote.Total())
			callFees.Add(callFees, quote.CallFee)
			totalFees.Add(totalFees, quote.CallFee)
			totalFees.Add(totalFees, quote.GasFee)
			totalFees.Add(totalFees, entities.NewUWei(quote.ProductFeeAmount))
		}
		refunded := isRefundedPegin(*quote, retainedQuote)
		if refunded {
			data.RefundedQuotesCount++
			totalPenalty.Add(totalPenalty, quote.PenaltyFee)
		}
	}
	for _, quote := range quotes {
		totalAmount.Add(totalAmount, quote.Total())
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

func (u *SummariesUseCase) aggregatePegoutData(ctx context.Context, startDate, endDate time.Time) (SummaryData, error) {
	var data SummaryData

	quotes, retainedQuotes, err := u.pegoutRepo.ListQuotesByDateRange(ctx, startDate, endDate)
	if err != nil {
		log.Printf("Error listing pegout quotes: %v", err)
		return data, err
	}
	data.TotalQuotesCount = int64(len(quotes) + len(retainedQuotes))
	var totalAmount = entities.NewWei(0)
	var acceptedTotalAmount = entities.NewWei(0)
	var totalFees = entities.NewWei(0)
	var callFees = entities.NewWei(0)
	var totalPenalty = entities.NewWei(0)

	for _, retainedQuote := range retainedQuotes {
		quote, err := u.pegoutRepo.GetQuote(ctx, retainedQuote.QuoteHash)
		if err != nil {
			log.Printf("Error getting pegout quote %s: %v", retainedQuote.QuoteHash, err)
			continue
		}
		if quote == nil {
			log.Printf("Pegout quote not found for hash %s", retainedQuote.QuoteHash)
			continue
		}
		totalAmount.Add(totalAmount, quote.Total())
		accepted := isAcceptedPegout(*quote, retainedQuote)

		if accepted {
			data.AcceptedQuotesCount++
			acceptedTotalAmount.Add(acceptedTotalAmount, quote.Total())
			callFees.Add(callFees, quote.CallFee)
			totalFees.Add(totalFees, quote.CallFee)
			totalFees.Add(totalFees, quote.GasFee)
			totalFees.Add(totalFees, entities.NewUWei(quote.ProductFeeAmount))
		}

		refunded := isRefundedPegout(*quote, retainedQuote)
		if refunded {
			data.RefundedQuotesCount++
			totalPenalty.Add(totalPenalty, entities.NewUWei(quote.PenaltyFee))
		}
	}
	for _, quote := range quotes {
		totalAmount.Add(totalAmount, quote.Total())
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

func isAcceptedPegin(q quote.PeginQuote, retained quote.RetainedPeginQuote) bool {
	return retained.Signature != "" &&
		retained.DepositAddress != "" &&
		(retained.State == quote.PeginStateCallForUserSucceeded)
}
func isRefundedPegin(q quote.PeginQuote, retained quote.RetainedPeginQuote) bool {
	return retained.State == quote.PeginStateCallForUserFailed ||
		retained.State == quote.PeginStateRegisterPegInFailed ||
		(retained.UserBtcTxHash != "" && retained.CallForUserTxHash == "")
}

func isAcceptedPegout(q quote.PegoutQuote, retained quote.RetainedPegoutQuote) bool {
	return retained.Signature != "" &&
		(retained.State == quote.PegoutStateBridgeTxSucceeded)
}

func isRefundedPegout(q quote.PegoutQuote, retained quote.RetainedPegoutQuote) bool {
	return retained.State == quote.PegoutStateBridgeTxFailed ||
		retained.State == quote.PegoutStateTimeForDepositElapsed
}
