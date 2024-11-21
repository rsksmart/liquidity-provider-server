package liquidity_provider_test

import (
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServerInfoUseCase_Run(t *testing.T) {
	t.Run("Should return error if BuildRevision doesn't have any value", func(t *testing.T) {
		liquidity_provider.BuildVersion = "version"
		liquidity_provider.BuildRevision = ""
		useCase := liquidity_provider.NewServerInfoUseCase()
		result, err := useCase.Run()
		assert.Empty(t, result)
		require.Error(t, err)
	})
	t.Run("Should return error if BuildVersion doesn't have any value", func(t *testing.T) {
		liquidity_provider.BuildVersion = ""
		liquidity_provider.BuildRevision = "revision"
		useCase := liquidity_provider.NewServerInfoUseCase()
		result, err := useCase.Run()
		assert.Empty(t, result)
		require.Error(t, err)
	})
	t.Run("Should return ServerInfo with BuildVersion and BuildRevision", func(t *testing.T) {
		liquidity_provider.BuildVersion = "version"
		liquidity_provider.BuildRevision = "revision"
		useCase := liquidity_provider.NewServerInfoUseCase()
		result, err := useCase.Run()
		assert.Equal(t, lp.ServerInfo{Version: "version", Revision: "revision"}, result)
		require.NoError(t, err)
	})
}
