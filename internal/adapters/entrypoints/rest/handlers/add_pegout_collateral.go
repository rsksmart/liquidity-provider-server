package handlers

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
)

// NewAddPegoutCollateralHandler
// @Title Add PegOut Collateral
// @Description Adds PegOut Collateral
// @Param AddCollateralRequest body pkg.AddCollateralRequest true "Add Collateral Request"
// @Success 200 object pkg.AddCollateralResponse
// @Route /pegout/addCollateral [post]
func NewAddPegoutCollateralHandler(useCase *pegout.AddCollateralUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		request := pkg.AddCollateralRequest{}
		if err = rest.DecodeRequest(w, req, &request); err != nil {
			return
		} else if err = rest.ValidateRequest(w, &request); err != nil {
			return
		}

		result, err := useCase.Run(entities.NewUWei(request.Amount))
		if errors.Is(err, usecases.InsufficientAmountError) {
			jsonErr := rest.NewErrorResponseWithDetails("not enough for minimum collateral", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusConflict, jsonErr)
			return
		} else if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("unknown error", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		response := pkg.AddCollateralResponse{NewCollateralBalance: result.Uint64()}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
