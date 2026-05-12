package reports

import (
	"context"
	"errors"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

// TransactionItem represents a single transaction in the response
type TransactionItem struct {
	QuoteHash string        `json:"quoteHash"`
	Amount    *entities.Wei `json:"amount"`
	CallFee   *entities.Wei `json:"callFee"`
	GasFee    *entities.Wei `json:"gasFee"`
	Status    string        `json:"status"`
}

// PaginationMetadata provides pagination information for transaction responses
type PaginationMetadata struct {
	Total      int `json:"total"`
	PerPage    int `json:"perPage"`
	TotalPages int `json:"totalPages"`
	Page       int `json:"page"`
}

// GetTransactionsResult represents the complete paginated response for transactions
type GetTransactionsResult struct {
	Data       []TransactionItem  `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

type GetTransactionsUseCase struct {
	peginRepo  quote.PeginQuoteRepository
	pegoutRepo quote.PegoutQuoteRepository
}

func NewGetTransactionsUseCase(
	peginRepo quote.PeginQuoteRepository,
	pegoutRepo quote.PegoutQuoteRepository,
) *GetTransactionsUseCase {
	return &GetTransactionsUseCase{
		peginRepo:  peginRepo,
		pegoutRepo: pegoutRepo,
	}
}

func CalculatePaginationMetadata(page, perPage, totalCount int) PaginationMetadata {
	totalPages := (totalCount + perPage - 1) / perPage // Ceiling division
	if totalPages < 1 {
		totalPages = 1
	}

	return PaginationMetadata{
		Total:      totalCount,
		PerPage:    perPage,
		TotalPages: totalPages,
		Page:       page,
	}
}

func (useCase *GetTransactionsUseCase) Run(ctx context.Context, transactionType string, startTime, endTime time.Time, page, perPage int) (GetTransactionsResult, error) {
	var transactions []TransactionItem
	var totalCount int
	var err error

	switch transactionType {
	case "pegin":
		transactions, totalCount, err = useCase.getPeginTransactions(ctx, startTime, endTime, page, perPage)
	case "pegout":
		transactions, totalCount, err = useCase.getPegoutTransactions(ctx, startTime, endTime, page, perPage)
	default:
		return GetTransactionsResult{}, usecases.WrapUseCaseError(usecases.GetTransactionsReportId,
			errors.New("invalid transaction type: must be 'pegin' or 'pegout'"))
	}

	if err != nil {
		return GetTransactionsResult{}, usecases.WrapUseCaseError(usecases.GetTransactionsReportId, err)
	}

	paginationMetadata := CalculatePaginationMetadata(page, perPage, totalCount)

	response := GetTransactionsResult{
		Data:       transactions,
		Pagination: paginationMetadata,
	}

	return response, nil
}

func (useCase *GetTransactionsUseCase) getPeginTransactions(ctx context.Context, startTime, endTime time.Time, page, perPage int) ([]TransactionItem, int, error) {
	quotePairs, totalCount, err := useCase.peginRepo.ListQuotesByDateRange(ctx, startTime, endTime, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	transactions := make([]TransactionItem, 0, len(quotePairs))
	for _, pair := range quotePairs {
		// Skip quotes that have not been accepted
		if pair.RetainedQuote.QuoteHash == "" {
			totalCount--
			continue
		}

		transaction := TransactionItem{
			QuoteHash: pair.RetainedQuote.QuoteHash,
			Amount:    pair.Quote.Value,
			CallFee:   pair.Quote.CallFee,
			GasFee:    pair.Quote.GasFee,
			Status:    string(pair.RetainedQuote.State),
		}

		transactions = append(transactions, transaction)
	}

	return transactions, totalCount, nil
}

func (useCase *GetTransactionsUseCase) getPegoutTransactions(ctx context.Context, startTime, endTime time.Time, page, perPage int) ([]TransactionItem, int, error) {
	quotePairs, totalCount, err := useCase.pegoutRepo.ListQuotesByDateRange(ctx, startTime, endTime, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	transactions := make([]TransactionItem, 0, len(quotePairs))
	for _, pair := range quotePairs {
		// Skip quotes that have not been accepted
		if pair.RetainedQuote.QuoteHash == "" {
			totalCount--
			continue
		}

		transaction := TransactionItem{
			QuoteHash: pair.RetainedQuote.QuoteHash,
			Amount:    pair.Quote.Value,
			CallFee:   pair.Quote.CallFee,
			GasFee:    pair.Quote.GasFee,
			Status:    string(pair.RetainedQuote.State),
		}

		transactions = append(transactions, transaction)
	}

	return transactions, totalCount, nil
}
