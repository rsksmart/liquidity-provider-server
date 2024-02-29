package handlers

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
)

// NewGetPegoutQuoteHandler
// @Title Pegout GetQuote
// @Description Gets Pegout Quote
// @Param PegoutQuoteRequest body pkg.PegoutQuoteRequest true "Interface with parameters for computing possible quotes for the service"
// @Success 200 array pkg.GetPegoutQuoteResponse The quote structure defines the conditions of a service, and acts as a contract between users and LPs
// @Route /pegout/getQuotes [post]
func NewGetPegoutQuoteHandler(useCase *pegout.GetQuoteUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		var result pegout.GetPegoutQuoteResult
		quoteRequest := pkg.PegoutQuoteRequest{}
		if err = rest.DecodeRequest(w, req, &quoteRequest); err != nil {
			return
		} else if err = rest.ValidateRequest(w, &quoteRequest); err != nil {
			return
		}

		pegoutRequest := pegout.NewQuoteRequest(
			quoteRequest.To,
			entities.NewUWei(quoteRequest.ValueToTransfer),
			quoteRequest.RskRefundAddress,
			quoteRequest.BitcoinRefundAddress,
		)

		result, err = useCase.Run(req.Context(), pegoutRequest)
		if errors.Is(err, blockchain.BtcAddressNotSupportedError) ||
			errors.Is(err, blockchain.BtcAddressInvalidNetworkError) ||
			errors.Is(err, usecases.RskAddressNotSupportedError) ||
			errors.Is(err, usecases.TxBelowMinimumError) ||
			errors.Is(err, liquidity_provider.AmountOutOfRangeError) {
			jsonErr := rest.NewErrorResponseWithDetails("invalid request", rest.DetailsFromError(err), true)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		} else if errors.Is(err, usecases.NoLiquidityError) {
			jsonErr := rest.NewErrorResponseWithDetails("no enough liquidity", rest.DetailsFromError(err), true)
			rest.JsonErrorResponse(w, http.StatusConflict, jsonErr)
			return
		} else if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("unknown error", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		quoteDto := pkg.ToPegoutQuoteDTO(result.PegoutQuote)
		responseBody := []pkg.GetPegoutQuoteResponse{{
			Quote:     quoteDto,
			QuoteHash: result.Hash,
		}} // to keep compatibility with legacy API
		rest.JsonResponseWithBody(w, http.StatusOK, &responseBody)
	}
}
