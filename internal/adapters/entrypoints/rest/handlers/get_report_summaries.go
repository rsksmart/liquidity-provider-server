package handlers

import (
	"net/http"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
)

const dateFormat = "2006-01-02"

func validateDateParameters(w http.ResponseWriter, req *http.Request) (startDate time.Time, endDate time.Time, valid bool) {
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
		jsonErr := rest.NewErrorResponseWithDetails("missing required parameters", map[string]any{
			"missing": missing,
		}, true)
		rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
		return time.Time{}, time.Time{}, false
	}

	var err error
	startDate, err = time.ParseInLocation(dateFormat, start, time.UTC)
	if err != nil {
		jsonErr := rest.NewErrorResponseWithDetails("invalid date format", rest.DetailsFromError(err), true)
		rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
		return time.Time{}, time.Time{}, false
	}

	endDate, err = time.ParseInLocation(dateFormat, end, time.UTC)
	if err != nil {
		jsonErr := rest.NewErrorResponseWithDetails("invalid date format", rest.DetailsFromError(err), true)
		rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
		return time.Time{}, time.Time{}, false
	}

	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())

	if endDate.Before(startDate) {
		details := map[string]any{
			"startDate": startDate.Format(dateFormat),
			"endDate":   endDate.Format(dateFormat),
		}
		jsonErr := rest.NewErrorResponseWithDetails("invalid date range", details, true)
		rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
		return time.Time{}, time.Time{}, false
	}

	return startDate, endDate, true
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
		startDate, endDate, valid := validateDateParameters(w, req)
		if !valid {
			return
		}

		response, err := useCase.Run(req.Context(), startDate, endDate)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
