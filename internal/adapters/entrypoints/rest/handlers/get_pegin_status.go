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

type PeginStatusUseCase interface {
	Run(ctx context.Context, quoteHash string) (quote.WatchedPeginQuote, error)
}

// NewGetPeginQuoteStatusHandler
// @Title GetPeginStatus
// @Description Returns the status of an accepted pegin quote
// @Param quoteHash query string true "Hash of the quote"
// @Success 200 {object} pkg.PeginQuoteStatusDTO "Object containing the quote itself and its status"
// @Router /pegin/status [get]
func NewGetPeginQuoteStatusHandler(useCase PeginStatusUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		var result quote.WatchedPeginQuote

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
		dto := pkg.PeginQuoteStatusDTO{
			Detail:       pkg.ToPeginQuoteDTO(result.PeginQuote),
			Status:       pkg.ToRetainedPeginQuoteDTO(result.RetainedQuote),
			CreationData: pkg.ToPeginCreationDataDTO(result.CreationData),
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &dto)
	}
}
