package handlers

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
)

// DateParamResult holds the result of date parameter validation
type DateParamResult struct {
	StartDate time.Time
	EndDate   time.Time
	Valid     bool
	Error     *rest.ErrorResponse
}

// validateDateParameters validates the startDate and endDate query parameters
func validateDateParameters(w http.ResponseWriter, req *http.Request) DateParamResult {
	result := DateParamResult{Valid: false}

	start := req.URL.Query().Get("startDate")
	end := req.URL.Query().Get("endDate")

	if start == "" || end == "" {
		missing := []string{}
		if start == "" {
			missing = append(missing, "startDate")
		}
		if end == "" {
			missing = append(missing, "endDate")
		}
		result.Error = rest.NewErrorResponseWithDetails("missing required parameters", map[string]any{
			"missing": missing,
		}, true)
		rest.JsonErrorResponse(w, http.StatusBadRequest, result.Error)
		return result
	}

	var err error
	result.StartDate, err = time.Parse(liquidity_provider.DateFormat, start)
	if err != nil {
		result.Error = rest.NewErrorResponseWithDetails("invalid date format", rest.DetailsFromError(err), true)
		rest.JsonErrorResponse(w, http.StatusBadRequest, result.Error)
		return result
	}

	result.EndDate, err = time.Parse(liquidity_provider.DateFormat, end)
	if err != nil {
		result.Error = rest.NewErrorResponseWithDetails("invalid date format", rest.DetailsFromError(err), true)
		rest.JsonErrorResponse(w, http.StatusBadRequest, result.Error)
		return result
	}

	// Set end date to end of day
	result.EndDate = time.Date(result.EndDate.Year(), result.EndDate.Month(), result.EndDate.Day(), 23, 59, 59, 0, time.UTC)

	if result.EndDate.Before(result.StartDate) {
		details := map[string]any{
			"startDate": result.StartDate.Format(liquidity_provider.DateFormat),
			"endDate":   result.EndDate.Format(liquidity_provider.DateFormat),
		}
		result.Error = rest.NewErrorResponseWithDetails("invalid date range", details, true)
		rest.JsonErrorResponse(w, http.StatusBadRequest, result.Error)
		return result
	}

	result.Valid = true
	return result
}

// NewGetReportSummariesHandler handles GET /report/summaries
// @Title Summaries
// @Description Returns financial data for a given period
// @Param startDate query string true "Start date in YYYY-MM-DD format" Format(date)
// @Param endDate query string true "End date in YYYY-MM-DD format" Format(date)
// @Success 200 {object} liquidity_provider.SummariesResponse "Financial data for the given period"
// @Router /report/summaries [get]
func NewGetReportSummariesHandler(useCase *liquidity_provider.SummariesUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		result := validateDateParameters(w, req)
		if !result.Valid {
			return
		}

		response, err := useCase.Run(req.Context(), result.StartDate, result.EndDate)
		if err != nil {
			log.Errorf("Error running summaries use case: %v", err)
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
