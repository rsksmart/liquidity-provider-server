package bindings_test

import (
	"testing"

	pauseregistry "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pause_registry"
	"github.com/stretchr/testify/assert"
)

func TestPauseRegistryBindingExposesPauseLevelAndStatus(t *testing.T) {
	contract := pauseregistry.NewPauseRegistryContract()

	assert.Contains(t, pauseregistry.PauseRegistryContractMetaData.ABI, `"name":"pauseLevel"`)
	assert.Contains(t, pauseregistry.PauseRegistryContractMetaData.ABI, `"name":"pauseStatus"`)

	var (
		_ = contract.PackPauseLevel
		_ = contract.UnpackPauseLevel
		_ = contract.PackPauseStatus
		_ = contract.UnpackPauseStatus
	)
}
