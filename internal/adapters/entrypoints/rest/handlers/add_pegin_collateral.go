package handlers

import (
	"errors"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

type AddPeginCollateralUseCase interface {
	Run(amount *entities.Wei) (*entities.Wei, error)
}

// NewAddPeginCollateralHandler
// @Title Add PegIn Collateral
// @Description Adds PegIn Collateral
// @Param AddCollateralRequest  body pkg.AddCollateralRequest true "Add Collateral Request"
// @Success 200 object pkg.AddCollateralResponse
// @Route /pegin/addCollateral [post]
func NewAddPeginCollateralHandler(useCase AddPeginCollateralUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		request := pkg.AddCollateralRequest{}
		if err = rest.DecodeRequest(w, req, &request); err != nil {
			return
		} else if err = rest.ValidateRequest(w, &request); err != nil {
			return
		}

		result, err := useCase.Run(entities.NewBigWei(request.Amount))
		if errors.Is(err, usecases.InsufficientAmountError) {
			jsonErr := rest.NewErrorResponseWithDetails("not enough for minimum collateral", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusConflict, jsonErr)
			return
		} else if errors.Is(err, blockchain.ContractPausedError) {
			jsonErr := rest.NewErrorResponseWithDetails("protocol is paused", rest.DetailsFromError(err), true)
			rest.JsonErrorResponse(w, http.StatusServiceUnavailable, jsonErr)
			return
		} else if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		response := pkg.AddCollateralResponse{NewCollateralBalance: result.AsBigInt()}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
