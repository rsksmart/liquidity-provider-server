package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
)

// NewVersionInfoHandler
// @Title Get server version
// @Description Returns the server version and revision
// @Route /version [get]
// @Success 200  object pkg.ServerInfoDTO
func NewVersionInfoHandler(useCase *liquidity_provider.ServerInfoUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		version, err := useCase.Run()
		if err != nil {
			rest.JsonErrorResponse(w, http.StatusInternalServerError, rest.NewErrorResponse(err.Error(), false))
			return
		}

		dto := pkg.ToServerInfoDTO(version)
		rest.JsonResponseWithBody(w, http.StatusOK, &dto)
	}
}
