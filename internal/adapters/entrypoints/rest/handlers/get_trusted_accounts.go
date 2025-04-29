package handlers

import (
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

// NewGetTrustedAccountsHandler
// @Title Get Trusted Accounts
// @Description Returns all trusted accounts in the system
// @Success 200 object
// @Route /management/trusted-accounts [get]
func NewGetTrustedAccountsHandler(useCase *liquidity_provider.GetTrustedAccountsUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		accounts, err := useCase.Run(req.Context())
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		response := pkg.TrustedAccountsResponse{
			Accounts: accounts,
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
