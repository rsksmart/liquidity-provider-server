package pkg_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestToRecommendedOperationDTO(t *testing.T) {
	domain := usecases.RecommendedOperationResult{
		RecommendedQuoteValue: entities.NewWei(1),
		EstimatedCallFee:      entities.NewWei(2),
		EstimatedGasFee:       entities.NewWei(3),
	}
	dto := pkg.ToRecommendedOperationDTO(domain)
	assert.Equal(t, pkg.RecommendedOperationDTO{
		RecommendedQuoteValue: big.NewInt(1),
		EstimatedCallFee:      big.NewInt(2),
		EstimatedGasFee:       big.NewInt(3),
	}, dto)
}
