package handlers

import (
	"net/http"
	"strings"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

// NewAcceptPeginAuthenticatedQuoteHandler
// @Title Accept authenticated quote
// @Description Accepts Quote with trusted account signature
// @Param Request body pkg.AcceptAuthenticatedQuoteRequest true "Quote Hash and Signature"
// @Success 200 object pkg.AcceptPeginRespose Interface that represents that the quote has been successfully accepted
// @Route /pegin/acceptAuthenticatedQuote [post]
func NewAcceptPeginAuthenticatedQuoteHandler(useCase AcceptQuoteUseCase) http.HandlerFunc {
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

		acceptRequest.Signature = strings.TrimPrefix(acceptRequest.Signature, "0x")

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
