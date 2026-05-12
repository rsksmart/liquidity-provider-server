package handlers

import (
	"context"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

type GetTrustedAccountsUseCase interface {
	Run(ctx context.Context) ([]entities.Signed[liquidity_provider.TrustedAccountDetails], error)
}

// NewGetTrustedAccountsHandler
// @Title Get Trusted Accounts
// @Description Returns all trusted accounts
// @Success 200 object
// @Route /management/trusted-accounts [get]
func NewGetTrustedAccountsHandler(useCase GetTrustedAccountsUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		accounts, err := useCase.Run(req.Context())
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		response := pkg.TrustedAccountsResponse{
			Accounts: pkg.ToTrustedAccountsDTO(accounts),
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
