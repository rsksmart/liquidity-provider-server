package handlers

import (
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

// NewAcceptPegoutAuthenticatedQuoteHandler
// @Title Accept authenticated quote
// @Description Accepts Quote with trusted account signature
// @Param Request body pkg.AcceptAuthenticatedQuoteRequest true "Quote Hash and Signature"
// @Success 200 object pkg.AcceptPegoutRespose Interface that represents that the quote has been successfully accepted
// @Route /pegout/acceptAuthenticatedQuote [post]
func NewAcceptPegoutAuthenticatedQuoteHandler(useCase AcceptQuoteUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		acceptRequest := pkg.AcceptAuthenticatedQuoteRequest{}
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
		if err != nil {
			HandleAcceptQuoteError(w, err)
			return
		}

		response := pkg.AcceptPeginRespose{
			Signature:                 acceptedQuote.Signature,
			BitcoinDepositAddressHash: acceptedQuote.DepositAddress,
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
