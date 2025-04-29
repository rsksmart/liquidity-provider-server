package handlers

import (
	"math/big"
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
// @Route /management/trusted-accounts/add [post]
func NewAddTrustedAccountHandler(useCase *lpuc.SetTrustedAccountUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		request := &pkg.TrustedAccountRequest{}
		if err = rest.DecodeRequest(w, req, request); err != nil {
			return
		} else if err = rest.ValidateRequest(w, request); err != nil {
			return
		}
		btcLockingCap, ok := new(big.Int).SetString(request.BtcLockingCap, 10)
		if !ok {
			jsonErr := rest.NewErrorResponse(InvalidBtcLockingCapMessage, true)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}
		rbtcLockingCap, ok := new(big.Int).SetString(request.RbtcLockingCap, 10)
		if !ok {
			jsonErr := rest.NewErrorResponse(InvalidRbtcLockingCapMessage, true)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}
		accountDetails := lp.TrustedAccountDetails{
			Address:          request.Address,
			Name:             request.Name,
			Btc_locking_cap:  entities.NewBigWei(btcLockingCap),
			Rbtc_locking_cap: entities.NewBigWei(rbtcLockingCap),
		}
		err = useCase.Add(req.Context(), accountDetails)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		rest.JsonResponse(w, http.StatusOK)
	}
}
