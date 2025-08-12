package pkg

import (
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// nolint:funlen
func TestGetTransactionsRequest_ValidateGetTransactionsRequest(t *testing.T) {
	tests := []struct {
		name        string
		request     GetTransactionsRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid pegin request with default pagination",
			request: GetTransactionsRequest{
				DateRangeRequest: DateRangeRequest{
					StartDate: "2023-01-01",
					EndDate:   "2023-01-31",
				},
				Type: "pegin",
			},
			expectError: false,
		},
		{
			name: "Valid pegout request with explicit pagination",
			request: GetTransactionsRequest{
				DateRangeRequest: DateRangeRequest{
					StartDate: "2023-01-01",
					EndDate:   "2023-01-31",
				},
				Type:    "pegout",
				Page:    2,
				PerPage: 50,
			},
			expectError: false,
		},
		{
			name: "Valid request with ISO 8601 dates",
			request: GetTransactionsRequest{
				DateRangeRequest: DateRangeRequest{
					StartDate: "2023-01-01T00:00:00Z",
					EndDate:   "2023-01-31T23:59:59Z",
				},
				Type: "pegin",
			},
			expectError: false,
		},
		{
			name: "Valid request with zero page (should apply default)",
			request: GetTransactionsRequest{
				DateRangeRequest: DateRangeRequest{
					StartDate: "2023-01-01",
					EndDate:   "2023-01-31",
				},
				Type: "pegin",
				Page: 0,
			},
			expectError: false,
		},
		{
			name: "Valid request with zero perPage (should apply default)",
			request: GetTransactionsRequest{
				DateRangeRequest: DateRangeRequest{
					StartDate: "2023-01-01",
					EndDate:   "2023-01-31",
				},
				Type:    "pegin",
				PerPage: 0,
			},
			expectError: false,
		},
		{
			name: "Invalid page (negative)",
			request: GetTransactionsRequest{
				DateRangeRequest: DateRangeRequest{
					StartDate: "2023-01-01",
					EndDate:   "2023-01-31",
				},
				Type: "pegin",
				Page: -1,
			},
			expectError: true,
			errorMsg:    "page must be at least 1",
		},
		{
			name: "Invalid perPage (negative)",
			request: GetTransactionsRequest{
				DateRangeRequest: DateRangeRequest{
					StartDate: "2023-01-01",
					EndDate:   "2023-01-31",
				},
				Type:    "pegin",
				PerPage: -1,
			},
			expectError: true,
			errorMsg:    "perPage must be at least 1",
		},
		{
			name: "Invalid perPage (too large)",
			request: GetTransactionsRequest{
				DateRangeRequest: DateRangeRequest{
					StartDate: "2023-01-01",
					EndDate:   "2023-01-31",
				},
				Type:    "pegin",
				PerPage: 101,
			},
			expectError: true,
			errorMsg:    "perPage cannot exceed 100",
		},
		{
			name: "Invalid date format",
			request: GetTransactionsRequest{
				DateRangeRequest: DateRangeRequest{
					StartDate: "invalid-date",
					EndDate:   "2023-01-31",
				},
				Type: "pegin",
			},
			expectError: true,
			errorMsg:    "startDate invalid date format: must be YYYY-MM-DD or ISO 8601 UTC format (ending with Z)",
		},
		{
			name: "End date before start date",
			request: GetTransactionsRequest{
				DateRangeRequest: DateRangeRequest{
					StartDate: "2023-01-31",
					EndDate:   "2023-01-01",
				},
				Type: "pegin",
			},
			expectError: true,
			errorMsg:    "endDate must be on or after startDate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.ValidateGetTransactionsRequest()

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				// Check that defaults were applied
				assert.GreaterOrEqual(t, tt.request.Page, 1)
				assert.GreaterOrEqual(t, tt.request.PerPage, 1)
			}
		})
	}
}

func TestGetTransactionsRequest_applyDefaults(t *testing.T) {
	tests := []struct {
		name            string
		request         GetTransactionsRequest
		expectedPage    int
		expectedPerPage int
	}{
		{
			name: "Apply default page and perPage",
			request: GetTransactionsRequest{
				Type: "pegin",
			},
			expectedPage:    1,
			expectedPerPage: 10,
		},
		{
			name: "Apply default page only",
			request: GetTransactionsRequest{
				Type:    "pegin",
				PerPage: 25,
			},
			expectedPage:    1,
			expectedPerPage: 25,
		},
		{
			name: "Apply default perPage only",
			request: GetTransactionsRequest{
				Type: "pegin",
				Page: 3,
			},
			expectedPage:    3,
			expectedPerPage: 10,
		},
		{
			name: "No defaults needed",
			request: GetTransactionsRequest{
				Type:    "pegin",
				Page:    2,
				PerPage: 50,
			},
			expectedPage:    2,
			expectedPerPage: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.request.applyDefaults()

			assert.Equal(t, tt.expectedPage, tt.request.Page)
			assert.Equal(t, tt.expectedPerPage, tt.request.PerPage)
		})
	}
}

func TestTransactionHistoryResponse_Structure(t *testing.T) {
	// Test that the response structure matches the expected JSON format
	response := GetTransactionsResponse{
		Data: []GetTransactionsItem{
			{
				QuoteHash: "0x1234567890abcdef1234567890abcdef12345678",
				Amount:    entities.NewWei(1000000000000000000),
				CallFee:   entities.NewWei(50000000000000000),
				GasFee:    entities.NewWei(10000000000000000),
				Status:    "RegisterPegInSucceeded",
			},
		},
		Pagination: PaginationMetadata{
			Total:      500,
			PerPage:    10,
			TotalPages: 50,
			Page:       1,
		},
	}

	// Verify the structure is properly initialized
	assert.Len(t, response.Data, 1)
	assert.Equal(t, "0x1234567890abcdef1234567890abcdef12345678", response.Data[0].QuoteHash)
	assert.Equal(t, "RegisterPegInSucceeded", response.Data[0].Status)
	assert.Equal(t, 500, response.Pagination.Total)
	assert.Equal(t, 50, response.Pagination.TotalPages)
}
