package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
)

// NewGetReportsAssetsHandler
// @Title Get asset Reports
// @Description Get the asset information for the LPS.
// @Success 200 pkg.GetAssetsReportDTO
// @Route /reports/assets [get]
func NewGetReportsAssetsHandler(useCase reports.AssetsReportUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		response, err := useCase.Run(req.Context())
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("Request error", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}
		responseDto := pkg.GetAssetsReportDTO{
			RbtcLocked:    response.RbtcLocked,
			RbtcBalance:   response.RbtcBalance,
			RbtcLiquidity: response.RbtcLiquidity,
			BtcLocked:     response.BtcLocked,
			BtcBalance:    response.BtcBalance,
			BtcLiquidity:  response.BtcLiquidity,
		}

		rest.JsonResponseWithBody(w, http.StatusOK, &responseDto)
	}
}
