package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
)

func NewGetUserQuotesHandler(useCase *pegout.GetUserDepositsUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		address := req.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "address parameter is required", http.StatusBadRequest)
			return
		}

		deposits, err := useCase.Run(req.Context(), address)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("unknown error", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}

		depositDtos := make([]pkg.DepositEventDTO, 0)

		for _, deposit := range deposits {
			depositDtos = append(depositDtos, pkg.DepositEventDTO{
				TxHash:      deposit.TxHash,
				QuoteHash:   deposit.QuoteHash,
				Amount:      deposit.Amount.AsBigInt(),
				Timestamp:   deposit.Timestamp,
				BlockNumber: deposit.BlockNumber,
				From:        deposit.From,
			})
		}

		rest.JsonResponseWithBody(w, http.StatusOK, &depositDtos)
	}
}
