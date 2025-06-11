package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	log "github.com/sirupsen/logrus"
)

type PegoutStatusUseCase interface {
	Run(ctx context.Context, quoteHash string) (quote.WatchedPegoutQuote, error)
}

// NewGetPegoutQuoteStatusHandler
// @Title GetPegoutStatus
// @Description Returns the status of an accepted pegout quote
// @Param quoteHash query string true "Hash of the quote"
// @Success 200 {object} pkg.PegoutQuoteStatusDTO "Object containing the quote itself and its status"
// @Router /pegout/status [get]
func NewGetPegoutQuoteStatusHandler(useCase PegoutStatusUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		var result quote.WatchedPegoutQuote
		const paramName = "quoteHash"
		quoteHash := req.URL.Query().Get(paramName)
		if err = quote.ValidateQuoteHash(quoteHash); err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("invalid or missing parameter quoteHash", rest.DetailsFromError(err), true)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}
		result, err = useCase.Run(req.Context(), quoteHash)
		if errors.Is(err, usecases.QuoteNotFoundError) {
			rest.JsonErrorResponse(w, http.StatusNotFound, rest.NewErrorResponse("Quote not found", true))
			return
		} else if errors.Is(err, usecases.QuoteNotAcceptedError) {
			rest.JsonErrorResponse(w, http.StatusConflict, rest.NewErrorResponse(err.Error(), true))
			return
		} else if err != nil {
			log.Error("Unknown error: ", err)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, rest.NewErrorResponse("Internal server error", false))
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
