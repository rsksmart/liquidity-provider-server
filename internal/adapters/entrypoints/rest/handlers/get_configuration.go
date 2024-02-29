package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"net/http"
)

// NewGetConfigurationHandler
// @Title Get configurations
// @Description Get all the configurations for the liquidity provider. Included in the management API.
// @Success 200 object
// @Route /configuration [get]
func NewGetConfigurationHandler(useCase *liquidity_provider.GetConfigUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		result := useCase.Run(req.Context())
		rest.JsonResponseWithBody(w, http.StatusOK, &result)
	}
}
