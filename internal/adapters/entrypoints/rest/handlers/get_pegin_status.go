package handlers

import (
	"errors"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

// NewGetPeginQuoteStatusHandler
// @Title GetPeginStatus
// @Description Returns the status of an accepted pegin quote
// @Param quoteHash query string true "Hash of the quote"
// @Success 200 {object} pkg.PeginQuoteStatusDTO "Object containing the quote itself and its status"
// @Router /pegin/status [get]
func NewGetPeginQuoteStatusHandler(useCase *pegin.StatusUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		const paramName = "quoteHash"
		quoteHash := req.URL.Query().Get(paramName)
		if quoteHash == "" {
			rest.ValidateRequestError(w, rest.RequiredQueryParam(paramName))
			return
		}
		result, err := useCase.Run(req.Context(), quoteHash)
		if errors.Is(err, usecases.QuoteNotFoundError) {
			rest.JsonErrorResponse(w, http.StatusNotFound, rest.NewErrorResponse("Quote not found", true))
			return
		} else if errors.Is(err, usecases.QuoteNotAcceptedError) {
			rest.JsonErrorResponse(w, http.StatusConflict, rest.NewErrorResponse(err.Error(), true))
			return
		}
		dto := pkg.PeginQuoteStatusDTO{
			Detail:       pkg.ToPeginQuoteDTO(result.PeginQuote),
			Status:       pkg.ToRetainedPeginQuoteDTO(result.RetainedQuote),
			CreationData: pkg.ToPeginCreationDataDTO(result.CreationData),
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &dto)
	}
}
