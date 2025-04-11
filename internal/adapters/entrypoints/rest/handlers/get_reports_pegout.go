package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
)

// NewGetReportsPegoutHandler
// @Title Get Pegout Reports
// @Description Get the last pegouts on the API. Included in the management API.
// @Param startDate query string true "Start date for the report"
// @Param endDate query string true "End date for the report"
// @Success 200 pkg.GetPegoutReportResponse
// @Route /reports/pegout [get]
func NewGetReportsPegoutHandler(useCase *pegout.GetPegoutReportUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var requestParams pkg.GetReportsPeginPegoutRequest
		var err error
		requestParams.StartDate = req.URL.Query().Get("startDate")
		requestParams.EndDate = req.URL.Query().Get("endDate")

		if err = requestParams.ValidateGetReportsPeginPegoutRequest(); err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("Validation error", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		startTime, endTime, err := requestParams.GetTimestamps()
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("Date conversion error", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		pegoutReport, err := useCase.Run(req.Context(), startTime, endTime)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		response := pkg.GetPegoutReportResponse{
			NumberOfQuotes:     pegoutReport.NumberOfQuotes,
			MinimumQuoteValue:  pegoutReport.MinimumQuoteValue.AsBigInt(),
			MaximumQuoteValue:  pegoutReport.MaximumQuoteValue.AsBigInt(),
			AverageQuoteValue:  pegoutReport.AverageQuoteValue.AsBigInt(),
			TotalFeesCollected: pegoutReport.TotalFeesCollected.AsBigInt(),
			AverageFeePerQuote: pegoutReport.AverageFeePerQuote.AsBigInt(),
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
