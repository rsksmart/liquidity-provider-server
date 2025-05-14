package handlers

import (
	"errors"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	lpuc "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

// NewUpdateTrustedAccountHandler
// @Title Update Trusted Account
// @Description Updates an existing trusted account
// @Param TrustedAccountRequest body pkg.TrustedAccountRequest true "Details of the trusted account to update"
// @Success 204 object
// @Route /management/trusted-accounts [put]
func NewUpdateTrustedAccountHandler(useCase *lpuc.UpdateTrustedAccountUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		request := &pkg.TrustedAccountRequest{}
		if err = rest.DecodeRequest(w, req, request); err != nil {
			return
		} else if err = rest.ValidateRequest(w, request); err != nil {
			return
		}
		accountDetails := lp.TrustedAccountDetails{
			Address:        request.Address,
			Name:           request.Name,
			BtcLockingCap:  entities.NewBigWei(request.BtcLockingCap),
			RbtcLockingCap: entities.NewBigWei(request.RbtcLockingCap),
		}
		err = useCase.Run(req.Context(), accountDetails)
		if errors.Is(err, lp.ErrTrustedAccountNotFound) {
			jsonErr := rest.NewErrorResponse(lp.ErrTrustedAccountNotFound.Error(), true)
			rest.JsonErrorResponse(w, http.StatusNotFound, jsonErr)
			return
		} else if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		rest.JsonResponse(w, http.StatusNoContent)
	}
}
