package handlers

import (
	"context"
	"errors"
	"math/big"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

type RecommendedPegoutUseCase interface {
	Run(ctx context.Context, userBalance *entities.Wei, destinationType blockchain.BtcAddressType) (usecases.RecommendedOperationResult, error)
}

// NewRecommendedPegoutHandler
// @Title Recommended pegout
// @Description Returns the recommended quote value to create a quote whose total payment is the input amount
// @Param amount query string true "Amount in wei expected to use as total payment for the quote"
// @Param destination_type query string false "Destination address type for the pegout. Is optional, but if provided, it will  increase the estimation accuracy. Must be one of: p2pkh, p2sh, p2wpkh, p2wsh, p2tr"
// @Success 200 object pkg.RecommendedOperationDTO "Recommended operation object"
// @Route /pegout/recommended [get]
func NewRecommendedPegoutHandler(useCase RecommendedPegoutUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const (
			amountParam = "amount"
			typeParam   = "destination_type"
		)
		var err error
		var parsedDestinationType blockchain.BtcAddressType

		amount := r.URL.Query().Get(amountParam)
		destinationType := r.URL.Query().Get(typeParam)
		parsedAmount, ok := new(big.Int).SetString(amount, 10)
		if !ok {
			jsonErr := rest.NewErrorResponse("invalid or missing parameter amount", true)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}
		if destinationType != "" {
			parsedDestinationType, err = blockchain.BtcAddressTypeFromString(destinationType)
		}
		if err != nil {
			jsonErr := rest.NewErrorResponse(err.Error(), true)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		result, err := useCase.Run(r.Context(), entities.NewBigWei(parsedAmount), parsedDestinationType)

		if errors.Is(err, usecases.NoLiquidityError) ||
			errors.Is(err, liquidity_provider.AmountOutOfRangeError) {
			jsonErr := rest.NewErrorResponse(err.Error(), true)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		} else if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		response := pkg.ToRecommendedOperationDTO(result)
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
