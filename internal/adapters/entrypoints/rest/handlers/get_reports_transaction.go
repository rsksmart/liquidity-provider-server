package handlers

import (
	"errors"
	"strconv"

	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

// NewGetReportsTransactionHandler
// @Title Get Transaction Reports
// @Description Get a paginated list of individual transactions of a specific type processed by the liquidity provider within a specified time period
// @Param type query string true "Transaction type filter: 'pegin' or 'pegout'"
// @Param startDate query string false "Start date for the report. Supports YYYY-MM-DD (expands to full day) or ISO 8601 format (YYYY-MM-DDTHH:mm:ssZ)"
// @Param endDate query string false "End date for the report. Supports YYYY-MM-DD (expands to end of day) or ISO 8601 format (YYYY-MM-DDTHH:mm:ssZ)"
// @Param page query int false "Page number to retrieve (1-indexed, default: 1)"
// @Param perPage query int false "Number of transactions per page (max: 100, default: 10)"
// @Success 200 {object} pkg.TransactionHistoryResponse "Paginated list of transactions with metadata"
// @Router /reports/transactions [get]
func NewGetReportsTransactionHandler(useCase *reports.GetTransactionsUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Parse and validate request parameters
		requestParams, err := parseTransactionQueryParameters(req)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("Parameter parsing error", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		if err = requestParams.ValidateGetTransactionHistoryRequest(); err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("Validation error", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		// Convert dates to timestamps (leverages dual datetime format support)
		startTime, endTime, err := requestParams.GetTimestamps()
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("Date conversion error", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		// Call use case
		result, err := useCase.Run(req.Context(), requestParams.Type, startTime, endTime, requestParams.Page, requestParams.PerPage)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}

		// Convert use case result to API response DTO and return
		response := mapUseCaseResultToResponse(result)
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}

// parseTransactionQueryParameters extracts and parses all query parameters from the HTTP request
func parseTransactionQueryParameters(req *http.Request) (pkg.GetTransactionHistoryRequest, error) {
	var requestParams pkg.GetTransactionHistoryRequest
	var err error

	// Parse basic query parameters
	requestParams.StartDate = req.URL.Query().Get("startDate")
	requestParams.EndDate = req.URL.Query().Get("endDate")
	requestParams.Type = req.URL.Query().Get("type")

	// Set default values for optional parameters
	requestParams.Page = 0    // Will trigger default in validation
	requestParams.PerPage = 0 // Will trigger default in validation

	// Parse pagination parameters
	pageParam := req.URL.Query().Get("page")
	if pageParam != "" {
		requestParams.Page, err = strconv.Atoi(pageParam)
		if err != nil {
			return requestParams, errors.New("invalid page parameter: " + err.Error())
		}
	}

	perPageParam := req.URL.Query().Get("perPage")
	if perPageParam != "" {
		requestParams.PerPage, err = strconv.Atoi(perPageParam)
		if err != nil {
			return requestParams, errors.New("invalid perPage parameter: " + err.Error())
		}
	}

	return requestParams, nil
}

// mapUseCaseResultToResponse converts the use case result to the API response DTO
func mapUseCaseResultToResponse(result reports.GetTransactionsResult) pkg.TransactionHistoryResponse {
	response := pkg.TransactionHistoryResponse{
		Data: make([]pkg.TransactionHistoryItem, len(result.Data)),
		Pagination: pkg.PaginationMetadata{
			Total:      result.Pagination.Total,
			PerPage:    result.Pagination.PerPage,
			TotalPages: result.Pagination.TotalPages,
			Page:       result.Pagination.Page,
		},
	}

	// Map transaction items
	for i, item := range result.Data {
		response.Data[i] = pkg.TransactionHistoryItem{
			QuoteHash: item.QuoteHash,
			Amount:    item.Amount,
			CallFee:   item.CallFee,
			GasFee:    item.GasFee,
			Status:    item.Status,
		}
	}

	return response
}
