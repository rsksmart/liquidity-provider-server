package handlers

import (
	"errors"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	lpuc "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
)

// NewDeleteTrustedAccountHandler
// @Title Delete Trusted Account
// @Description Deletes a trusted account
// @Param address query string true "Address of the trusted account to delete"
// @Success 204 object
// @Route /management/trusted-accounts [delete]
func NewDeleteTrustedAccountHandler(useCase *lpuc.DeleteTrustedAccountUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		address := req.URL.Query().Get("address")
		if address == "" {
			rest.ValidateRequestError(w, rest.RequiredQueryParam(address))
			return
		}
		err = useCase.Run(req.Context(), address)
		if errors.Is(err, lp.ErrTrustedAccountNotFound) {
			jsonErr := rest.NewErrorResponse(err.Error(), true)
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
