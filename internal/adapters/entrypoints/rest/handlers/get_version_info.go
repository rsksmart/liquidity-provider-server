package handlers

import (
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

type ServerInfoUseCase interface {
	Run() (liquidity_provider.ServerInfo, error)
}

// NewVersionInfoHandler
// @Title Get server version
// @Description Returns the server version and revision
// @Route /version [get]
// @Success 200  object pkg.ServerInfoDTO
func NewVersionInfoHandler(useCase ServerInfoUseCase) http.HandlerFunc {
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
