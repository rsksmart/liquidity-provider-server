package handlers

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
)

type GetAssetsReportUseCase interface {
	Run(ctx context.Context) (reports.GetAssetsReportResult, error)
}

// NewGetReportsAssetsHandler
// @Title Get asset Reports
// @Description Get the asset information for the LPS.
// @Success 200 pkg.GetAssetsReportDTO
// @Route /reports/assets [get]
func NewGetReportsAssetsHandler(useCase GetAssetsReportUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		response, err := useCase.Run(req.Context())
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("Request error", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}
		responseDto := pkg.GetAssetsReportDTO{
			RbtcLockedLbc:      response.RbtcLockedLbc,
			RbtcLockedForUsers: response.RbtcLockedForUsers,
			RbtcWaitingRefund:  response.RbtcWaitingRefund,
			RbtcLiquidity:      response.RbtcLiquidity,
			RbtcWalletBalance:  response.RbtcWalletBalance,
			BtcLockedForUsers:  response.BtcLockedForUsers,
			BtcLiquidity:       response.BtcLiquidity,
			BtcWalletBalance:   response.BtcWalletBalance,
			BtcRebalancing:     response.BtcRebalancing,
		}

		rest.JsonResponseWithBody(w, http.StatusOK, &responseDto)
	}
}
