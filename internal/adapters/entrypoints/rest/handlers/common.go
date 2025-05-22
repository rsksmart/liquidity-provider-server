package handlers

import (
	"errors"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

const UnknownErrorMessage = "unknown error"

func closeManagementSession(req *http.Request, w http.ResponseWriter, env environment.ManagementEnv) error {
	const errorMsg = "error closing session"
	cookieStore, err := cookies.GetSessionCookieStore(env)
	if err != nil {
		jsonErr := rest.NewErrorResponseWithDetails(errorMsg, rest.DetailsFromError(err), false)
		rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
		return err
	}

	err = cookies.CloseManagementSession(&cookies.CloseSessionArgs{
		Store:   cookieStore,
		Request: req,
		Writer:  w,
	})
	if err != nil {
		jsonErr := rest.NewErrorResponseWithDetails(errorMsg, rest.DetailsFromError(err), false)
		rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
		return err
	}
	return nil
}

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
	case errors.Is(err, liquidity_provider.ErrTamperedTrustedAccount):
		jsonErr := rest.NewErrorResponseWithDetails("error fetching trusted account", rest.DetailsFromError(err), true)
		rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
	default:
		jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
		rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
	}
}
