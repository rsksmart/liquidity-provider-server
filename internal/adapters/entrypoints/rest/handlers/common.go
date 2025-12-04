package handlers

import (
	"errors"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

const UnknownErrorMessage = "unknown error"

func HandleAcceptQuoteError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, usecases.QuoteNotFoundError):
		jsonErr := rest.NewErrorResponseWithDetails("quote not found", rest.DetailsFromError(err), true)
		rest.JsonErrorResponse(w, http.StatusNotFound, jsonErr)
	case errors.Is(err, usecases.ExpiredQuoteError):
		jsonErr := rest.NewErrorResponseWithDetails("expired quote", rest.DetailsFromError(err), true)
		rest.JsonErrorResponse(w, http.StatusGone, jsonErr)
	case errors.Is(err, usecases.NoLiquidityError):
		jsonErr := rest.NewErrorResponseWithDetails("not enough liquidity", rest.DetailsFromError(err), true)
		rest.JsonErrorResponse(w, http.StatusConflict, jsonErr)
	case errors.Is(err, usecases.LockingCapExceededError):
		jsonErr := rest.NewErrorResponseWithDetails("locking cap exceeded", rest.DetailsFromError(err), true)
		rest.JsonErrorResponse(w, http.StatusConflict, jsonErr)
	case errors.Is(err, liquidity_provider.TamperedTrustedAccountError):
		jsonErr := rest.NewErrorResponseWithDetails("error fetching trusted account", rest.DetailsFromError(err), true)
		rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
	default:
		jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
		rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
	}
}
