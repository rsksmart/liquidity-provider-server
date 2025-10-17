package handlers

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"math/big"
	"net/http"
	"strings"
)

type RecommendedPeginUseCase interface {
	Run(ctx context.Context, userBalance *entities.Wei, destinationAddress string, data []byte) (usecases.RecommendedOperationResult, error)
}

// NewRecommendedPeginHandler
// @Title Recommended pegin
// @Description Returns the recommended quote value to create a quote whose total payment is the input amount
// @Param amount query string true "Amount in wei expected to use as total payment for the quote"
// @Param destination_address query string false "Destination address for the pegin. Is optional, but if provided, it will increase the estimation accuracy."
// @Param data query string false "Hex-encoded data payload to include in the pegin transaction. Is optional, but if provided, it will increase the estimation accuracy."
// @Success 200 object pkg.RecommendedOperationDTO "Recommended operation object"
// @Route /pegin/recommended [get]
func NewRecommendedPeginHandler(useCase RecommendedPeginUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const (
			amountParam  = "amount"
			addressParam = "destination_address"
			dataParam    = "data"
		)

		amount := r.URL.Query().Get(amountParam)
		destinationAddress := r.URL.Query().Get(addressParam)
		data := r.URL.Query().Get(dataParam)

		parsedAmount, ok := new(big.Int).SetString(amount, 10)
		if !ok {
			jsonErr := rest.NewErrorResponse("invalid or missing parameter amount", true)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		parsedData, err := hex.DecodeString(strings.ToLower(strings.TrimPrefix(data, "0x")))
		if err != nil {
			jsonErr := rest.NewErrorResponse(err.Error(), true)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		result, err := useCase.Run(r.Context(), entities.NewBigWei(parsedAmount), destinationAddress, parsedData)

		if errors.Is(err, usecases.NoLiquidityError) ||
			errors.Is(err, usecases.TxBelowMinimumError) ||
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
