package pkg

import (
	"errors"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

// GetTransactionsRequest combines date range filtering with pagination and type filtering
type GetTransactionsRequest struct {
	DateRangeRequest
	Type    string `json:"type" validate:"required,oneof=pegin pegout"`
	Page    int    `json:"page" validate:"min=0"`
	PerPage int    `json:"perPage" validate:"min=0,max=100"`
}

// GetTransactionsItem represents a single transaction in the history response
type GetTransactionsItem struct {
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

// GetTransactionsResponse represents the complete paginated response for transaction history
type GetTransactionsResponse struct {
	Data       []GetTransactionsItem `json:"data"`
	Pagination PaginationMetadata    `json:"pagination"`
}

func (r *GetTransactionsRequest) ValidateGetTransactionsRequest() error {
	// Validate date range format and logical consistency
	if err := r.ValidateDateRange(); err != nil {
		return err
	}

	// Validate pagination parameters before applying defaults
	if r.Page < 0 {
		return errors.New("page must be at least 1")
	}
	if r.PerPage < 0 {
		return errors.New("perPage must be at least 1")
	}
	if r.PerPage > 100 {
		return errors.New("perPage cannot exceed 100")
	}

	r.applyDefaults()

	return nil
}

// applyDefaults sets default values for optional pagination parameters
func (r *GetTransactionsRequest) applyDefaults() {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.PerPage == 0 {
		r.PerPage = 10
	}
}
