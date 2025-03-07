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
// @Param PeginReportsRequest  body pkg.PeginReportsRequest true "Date range for the report, case not provided will be last 30 days"
// @Success 200 pkg.GetPeginReportResponse
// @Route /reports/pegin [get]
func NewGetReportsPeginHandler(useCase *pegin.GetPeginReportUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error

		peginReport, err := useCase.Run(req.Context())
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
