package handlers

import (
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	entities_lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	uc_lp "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

// NewGetConfigurationHandler
// @Title Get configurations
// @Description Get all the configurations for the liquidity provider. Included in the management API.
// @Success 200 object
// @Route /configuration [get]
func NewGetConfigurationHandler(useCase *uc_lp.GetConfigUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		result := useCase.Run(req.Context())
		response := struct {
			General entities_lp.GeneralConfiguration `json:"general"`
			Pegin   pkg.PeginConfigurationDTO        `json:"pegin"`
			Pegout  pkg.PegoutConfigurationDTO       `json:"pegout"`
		}{
			General: result.General,
			Pegin:   pkg.ToPeginConfigurationDTO(result.Pegin),
			Pegout:  pkg.ToPegoutConfigurationDTO(result.Pegout),
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
