package pkg

import (
	"errors"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

// GetTransactionHistoryRequest extends GetReportsByPeriodRequest with pagination and type filtering
type GetTransactionHistoryRequest struct {
	GetReportsByPeriodRequest
	Type    string `json:"type" validate:"required,oneof=pegin pegout"`
	Page    int    `json:"page" validate:"min=1"`
	PerPage int    `json:"perPage" validate:"min=1,max=100"`
}

// TransactionHistoryItem represents a single transaction in the history response
type TransactionHistoryItem struct {
	QuoteHash string        `json:"quoteHash"`
	Amount    *entities.Wei `json:"amount"`
	CallFee   *entities.Wei `json:"callFee"`
	GasFee    *entities.Wei `json:"gasFee"`
	Status    string        `json:"status"`
}

// PaginationMetadata provides pagination information for paginated responses
type PaginationMetadata struct {
	Total      int `json:"total"`
	PerPage    int `json:"perPage"`
	TotalPages int `json:"totalPages"`
	Page       int `json:"page"`
}

// TransactionHistoryResponse represents the complete paginated response for transaction history
type TransactionHistoryResponse struct {
	Data       []TransactionHistoryItem `json:"data"`
	Pagination PaginationMetadata       `json:"pagination"`
}

// ValidateGetTransactionHistoryRequest validates the transaction history request including inherited validation
func (r *GetTransactionHistoryRequest) ValidateGetTransactionHistoryRequest() error {
	// First validate the inherited date range fields
	if err := r.ValidateGetReportsByPeriodRequest(); err != nil {
		return err
	}

	// Validate transaction type
	if r.Type == "" {
		return errors.New("type is required")
	}
	if r.Type != "pegin" && r.Type != "pegout" {
		return errors.New("type must be 'pegin' or 'pegout'")
	}

	// Validate pagination parameters before applying defaults
	// Check for explicit invalid values (negative values are invalid)
	if r.Page < 0 {
		return errors.New("page must be at least 1")
	}
	if r.PerPage < 0 {
		return errors.New("perPage must be at least 1")
	}
	if r.PerPage > 100 {
		return errors.New("perPage cannot exceed 100")
	}

	// Apply default values after validation
	r.applyDefaults()

	return nil
}

// applyDefaults sets default values for optional pagination parameters
func (r *GetTransactionHistoryRequest) applyDefaults() {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.PerPage == 0 {
		r.PerPage = 10
	}
}
