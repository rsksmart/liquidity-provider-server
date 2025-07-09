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

const (
	InvalidBtcLockingCapMessage  = "Invalid BTC locking cap value"
	InvalidRbtcLockingCapMessage = "Invalid RBTC locking cap value"
)

// NewAddTrustedAccountHandler
// @Title Add Trusted Account
// @Description Adds a new trusted account
// @Param TrustedAccountRequest body pkg.TrustedAccountRequest true "Details of the trusted account to add"
// @Success 204 object
// @Route /management/trusted-accounts [post]
func NewAddTrustedAccountHandler(useCase *lpuc.AddTrustedAccountUseCase) http.HandlerFunc {
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
		if errors.Is(err, lp.DuplicateTrustedAccountError) {
			jsonErr := rest.NewErrorResponse(lp.DuplicateTrustedAccountError.Error(), true)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		} else if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		rest.JsonResponse(w, http.StatusNoContent)
	}
}
