package liquidity_provider

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

var (
	BuildVersion  string
	BuildRevision string
)

type ServerInfoUseCase struct{}

func NewServerInfoUseCase() *ServerInfoUseCase {
	return &ServerInfoUseCase{}
}

func (useCase *ServerInfoUseCase) Run() (liquidity_provider.ServerInfo, error) {
	if BuildVersion == "" || BuildRevision == "" {
		return liquidity_provider.ServerInfo{}, usecases.WrapUseCaseError(usecases.ServerInfoId, errors.New("unable to read build info"))
	}
	return liquidity_provider.ServerInfo{
		Version:  BuildVersion,
		Revision: BuildRevision,
	}, nil
}
