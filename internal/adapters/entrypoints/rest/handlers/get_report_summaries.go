package handlers

import (
	"net/http"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
)

const dateFormat = "2006-01-02"

// NewGetReportSummariesHandler handles GET /report/summaries
// @Title Summaries
// @Description Returns financial data for a given period
// @Param startDate query string true "Start date in YYYY-MM-DD format"
// @Param endDate query string true "End date in YYYY-MM-DD format"
// @Success 200 {object} report.SummariesResponse "Financial data for the given period"
// @Router /report/summaries [get]
func NewGetReportSummariesHandler(useCase *liquidity_provider.SummariesUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
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
			return
		}

		startDate, err := time.Parse(dateFormat, start)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("invalid date format", rest.DetailsFromError(err), true)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		endDate, err := time.Parse(dateFormat, end)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("invalid date format", rest.DetailsFromError(err), true)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		if endDate.Before(startDate) {
			details := map[string]any{
				"startDate": startDate.Format(dateFormat),
				"endDate":   endDate.Format(dateFormat),
			}
			jsonErr := rest.NewErrorResponseWithDetails("invalid date range", details, true)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
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
