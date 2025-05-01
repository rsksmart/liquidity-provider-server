package handlers

import (
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	lpuc "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

// NewDeleteTrustedAccountHandler
// @Title Delete Trusted Account
// @Description Deletes a trusted account
// @Param TrustedAccountAddressRequest body pkg.TrustedAccountAddressRequest true "Address of the trusted account to delete"
// @Success 204 object
// @Route /management/trusted-accounts [delete]
func NewDeleteTrustedAccountHandler(useCase *lpuc.DeleteTrustedAccountUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		request := &pkg.TrustedAccountAddressRequest{}
		if err = rest.DecodeRequest(w, req, request); err != nil {
			return
		} else if err = rest.ValidateRequest(w, request); err != nil {
			return
		}
		err = useCase.Run(req.Context(), request.Address)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		rest.JsonResponse(w, http.StatusNoContent)
	}
}
