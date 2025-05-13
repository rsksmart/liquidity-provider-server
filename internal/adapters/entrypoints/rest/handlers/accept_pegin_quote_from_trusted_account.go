package handlers

import (
	"errors"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

// NewAcceptPeginQuoteFromTrustedAccountHandler
// @Title Accept Quote From Trusted Account
// @Description Accepts Quote with trusted account signature
// @Param Request body pkg.AcceptQuoteFromTrustedAccountRequest true "Quote Hash and Signature"
// @Success 200 object pkg.AcceptPeginRespose Interface that represents that the quote has been successfully accepted
// @Route /pegin/acceptQuoteFromTrustedAccount [post]
func NewAcceptPeginQuoteFromTrustedAccountHandler(useCase *pegin.AcceptQuoteUseCase) http.HandlerFunc {
	return NewAcceptPeginQuoteFromTrustedAccountHandlerWithInterface(useCase)
}

// NewAcceptPeginQuoteFromTrustedAccountHandlerWithInterface is like NewAcceptPeginQuoteFromTrustedAccountHandler but accepts an interface instead of a concrete type for testing
func NewAcceptPeginQuoteFromTrustedAccountHandlerWithInterface(useCase AcceptQuoteUseCaseInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		acceptRequest := pkg.AcceptQuoteFromTrustedAccountRequest{}
		if err = rest.DecodeRequest(w, req, &acceptRequest); err != nil {
			return
		} else if err = rest.ValidateRequest(w, &acceptRequest); err != nil {
			return
		}

		if err = quote.ValidateQuoteHash(acceptRequest.QuoteHash); err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("invalid quote hash", rest.DetailsFromError(err), true)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		acceptedQuote, err := useCase.Run(req.Context(), acceptRequest.QuoteHash, acceptRequest.Signature)
		if errors.Is(err, usecases.QuoteNotFoundError) {
			jsonErr := rest.NewErrorResponseWithDetails("quote not found", rest.DetailsFromError(err), true)
			rest.JsonErrorResponse(w, http.StatusNotFound, jsonErr)
			return
		} else if errors.Is(err, usecases.ExpiredQuoteError) {
			jsonErr := rest.NewErrorResponseWithDetails("expired quote", rest.DetailsFromError(err), true)
			rest.JsonErrorResponse(w, http.StatusGone, jsonErr)
			return
		} else if errors.Is(err, usecases.NoLiquidityError) {
			jsonErr := rest.NewErrorResponseWithDetails("not enough liquidity", rest.DetailsFromError(err), true)
			rest.JsonErrorResponse(w, http.StatusConflict, jsonErr)
			return
		} else if errors.Is(err, usecases.LockingCapExceededError) {
			jsonErr := rest.NewErrorResponseWithDetails("locking cap exceeded", rest.DetailsFromError(err), true)
			rest.JsonErrorResponse(w, http.StatusConflict, jsonErr)
			return
		} else if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}

		response := pkg.AcceptPeginRespose{
			Signature:                 acceptedQuote.Signature,
			BitcoinDepositAddressHash: acceptedQuote.DepositAddress,
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
