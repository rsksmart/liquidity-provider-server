package handlers

import (
	"context"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

type GetAssetsReportUseCase interface {
	Run(ctx context.Context) (reports.GetAssetsReportResult, error)
}

// NewGetReportsAssetsHandler
// @Title Get asset Reports
// @Description Get the asset information for the LPS including BTC and RBTC balances, locations, and allocations.
// @Success 200 {object} pkg.GetAssetsReportResponse "Detailed asset report with BTC and RBTC information"
// @Router /reports/assets [get]
func NewGetReportsAssetsHandler(useCase GetAssetsReportUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		response, err := useCase.Run(req.Context())
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}

		responseDto := pkg.ToGetAssetsReportResponse(response)

		rest.JsonResponseWithBody(w, http.StatusOK, &responseDto)
	}
}
