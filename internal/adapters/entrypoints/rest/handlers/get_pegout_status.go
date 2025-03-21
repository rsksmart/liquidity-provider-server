package handlers

import (
	"errors"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

// NewGetPegoutQuoteStatusHandler
// @Title GetPegoutStatus
// @Description Returns the status of an accepted pegout quote
// @Param quoteHash query string true "Hash of the quote"
// @Success 200 {object} pkg.PegoutQuoteStatusDTO "Object containing the quote itself and its status"
// @Router /pegout/status [get]
func NewGetPegoutQuoteStatusHandler(useCase *pegout.StatusUseCase) http.HandlerFunc {
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
		dto := pkg.PegoutQuoteStatusDTO{
			Detail:       pkg.ToPegoutQuoteDTO(result.PegoutQuote),
			Status:       pkg.ToRetainedPegoutQuoteDTO(result.RetainedQuote),
			CreationData: pkg.ToPegoutCreationDataDTO(result.CreationData),
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &dto)
	}
}
