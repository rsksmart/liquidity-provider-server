package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

// NewGetReportSummariesHandler handles GET /reports/summaries
// @Title Summaries
// @Description Returns financial data for a given period
// @Param startDate query string true "Start date in YYYY-MM-DD format" Format(date)
// @Param endDate query string true "End date in YYYY-MM-DD format" Format(date)
// @Success 200 {object} pkg.SummaryResultDTO "Financial data for the given period"
// @Router /reports/summaries [get]
func NewGetReportSummariesHandler(useCase *reports.SummariesUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		startDate, endDate, err := rest.ParseDateRange(req, reports.DateFormat)
		if err != nil {
			log.Errorf("Error parsing date range: %v", err)
			rest.JsonErrorResponse(w, http.StatusBadRequest,
				rest.NewErrorResponseWithDetails("Invalid date range", rest.DetailsFromError(err), true))
			return
		}
		validateErr := rest.ValidateDateRange(startDate, endDate, reports.DateFormat)
		if validateErr != nil {
			log.Errorf("Error validating date range: %v", validateErr)
			rest.JsonErrorResponse(w, http.StatusBadRequest,
				rest.NewErrorResponseWithDetails("Invalid date range", rest.DetailsFromError(validateErr), true))
			return
		}
		response, err := useCase.Run(req.Context(), startDate, endDate)
		if err != nil {
			log.Errorf("Error running summaries use case: %v", err)
			rest.JsonErrorResponse(w, http.StatusInternalServerError,
				rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false))
			return
		}
		dto := pkg.ToSummaryResultDTO(response)
		rest.JsonResponseWithBody(w, http.StatusOK, &dto)
	}
}
