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
// @Success 200 pkg.GetPegoutReportResponse
// @Route /reports/pegout [get]
func NewGetReportsPegoutHandler(useCase *pegout.GetPegoutReportUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error

		pegoutReport, err := useCase.Run(req.Context())
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
