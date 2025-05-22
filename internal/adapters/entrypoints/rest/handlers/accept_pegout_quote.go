package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

type AcceptPegoutQuoteUseCase interface {
	Run(ctx context.Context, quoteHash, signature string) (quote.AcceptedQuote, error)
}

// NewAcceptPegoutQuoteHandler
// @Title Accept Quote Pegout
// @Description Accepts Quote Pegout
// @Param QuoteHash body pkg.AcceptQuoteRequest true "Quote Hash"
// @Success 200 object pkg.AcceptPegoutResponse
// @Route /pegout/acceptQuote [post]
func NewAcceptPegoutQuoteHandler(useCase AcceptPegoutQuoteUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		acceptRequest := pkg.AcceptQuoteRequest{}
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

		acceptedQuote, err := useCase.Run(req.Context(), acceptRequest.QuoteHash, "")
		if errors.Is(err, usecases.QuoteNotFoundError) {
			jsonErr := rest.NewErrorResponseWithDetails("invalid quote hash", rest.DetailsFromError(err), true)
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
		} else if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}

		response := pkg.AcceptPegoutResponse{
			Signature:  acceptedQuote.Signature,
			LbcAddress: acceptedQuote.DepositAddress,
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
