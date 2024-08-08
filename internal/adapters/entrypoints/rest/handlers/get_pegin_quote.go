package handlers

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
)

// NewGetPeginQuoteHandler
// @Title Pegin GetQuote
// @Description Gets Pegin Quote
// @Param PeginQuoteRequest  body pkg.PeginQuoteRequest true "Interface with parameters for computing possible quotes for the service"
// @Success 200 array pkg.GetPeginQuoteResponse The quote structure defines the conditions of a service, and acts as a contract between users and LPs
// @Route /pegin/getQuote [post]
func NewGetPeginQuoteHandler(useCase *pegin.GetQuoteUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		var result pegin.GetPeginQuoteResult
		var callArgument []byte
		quoteRequest := pkg.PeginQuoteRequest{}
		if err = rest.DecodeRequest(w, req, &quoteRequest); err != nil {
			return
		} else if err = rest.ValidateRequest(w, &quoteRequest); err != nil {
			return
		}

		if callArgument, err = blockchain.DecodeStringTrimPrefix(quoteRequest.CallContractArguments); err != nil {
			jsonErr := rest.NewErrorResponse(err.Error(), true)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		peginRequest := pegin.NewQuoteRequest(
			quoteRequest.CallEoaOrContractAddress,
			callArgument,
			entities.NewUWei(quoteRequest.ValueToTransfer),
			quoteRequest.RskRefundAddress,
			quoteRequest.BitcoinRefundAddress,
		)

		result, err = useCase.Run(req.Context(), peginRequest)
		if errors.Is(err, blockchain.BtcAddressNotSupportedError) ||
			errors.Is(err, blockchain.BtcAddressInvalidNetworkError) ||
			errors.Is(err, usecases.RskAddressNotSupportedError) ||
			errors.Is(err, usecases.TxBelowMinimumError) ||
			errors.Is(err, liquidity_provider.AmountOutOfRangeError) {
			jsonErr := rest.NewErrorResponseWithDetails("invalid request", rest.DetailsFromError(err), true)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		} else if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		quoteDto := pkg.ToPeginQuoteDTO(result.PeginQuote)
		responseBody := []pkg.GetPeginQuoteResponse{{
			Quote:     quoteDto,
			QuoteHash: result.Hash,
		}} // to keep compatibility with legacy API
		rest.JsonResponseWithBody(w, http.StatusOK, &responseBody)
	}
}
