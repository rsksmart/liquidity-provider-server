package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
)

func NewHealthCheckHandler(useCase *usecases.HealthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		result := useCase.Run(req.Context())
		response := pkg.HealthResponse{
			Status: result.Status,
			Services: pkg.Services{
				Db:  result.Services.Db,
				Rsk: result.Services.Rsk,
				Btc: result.Services.Btc,
			},
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
