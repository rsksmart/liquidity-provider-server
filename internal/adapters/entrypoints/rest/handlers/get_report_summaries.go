package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
)

// NewGetReportSummariesHandler handles GET /report/summaries
// @Title Summaries
// @Description Returns financial data for a given period
// @Param startDate query string true "Start date in YYYY-MM-DD format" Format(date)
// @Param endDate query string true "End date in YYYY-MM-DD format" Format(date)
// @Success 200 {object} liquidity_provider.SummariesResponse "Financial data for the given period"
// @Router /report/summaries [get]
func NewGetReportSummariesHandler(useCase *liquidity_provider.SummariesUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		startDate, endDate, valid := rest.ValidateDateRange(w, req, liquidity_provider.DateFormat)
		if !valid {
			return
		}
		response, err := useCase.Run(req.Context(), startDate, endDate)
		if err != nil {
			log.Errorf("Error running summaries use case: %v", err)
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
