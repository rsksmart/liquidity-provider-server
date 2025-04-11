package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
)

// NewGetReportsPeginHandler
// @Title Get Pegin Reports
// @Description Get the last pegins on the API. Included in the management API.
// @Param startDate query string true "Start date for the report"
// @Param endDate query string true "End date for the report"
// @Success 200 pkg.GetPeginReportResponse
// @Route /reports/pegin [get]
func NewGetReportsPeginHandler(useCase *pegin.GetPeginReportUseCase) http.HandlerFunc {
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

		peginReport, err := useCase.Run(req.Context(), startTime, endTime)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		response := pkg.GetPeginReportResponse{
			NumberOfQuotes:     peginReport.NumberOfQuotes,
			MinimumQuoteValue:  peginReport.MinimumQuoteValue.AsBigInt(),
			MaximumQuoteValue:  peginReport.MaximumQuoteValue.AsBigInt(),
			AverageQuoteValue:  peginReport.AverageQuoteValue.AsBigInt(),
			TotalFeesCollected: peginReport.TotalFeesCollected.AsBigInt(),
			AverageFeePerQuote: peginReport.AverageFeePerQuote.AsBigInt(),
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
