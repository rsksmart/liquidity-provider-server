package rootstock_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestFeeCollectorImpl_DaoFeePercentage(t *testing.T) {
	lbcMock := &mocks.LbcBindingMock{}
	t.Run("Success", func(t *testing.T) {
		lbcMock.On("ProductFeePercentage", mock.Anything).Return(big.NewInt(1), nil).Once()
		feeCollector := rootstock.NewFeeCollectorImpl(lbcMock, rootstock.RetryParams{Retries: 0, Sleep: 0})
		percentage, err := feeCollector.DaoFeePercentage()
		require.NoError(t, err)
		require.Equal(t, uint64(1), percentage)
	})
	t.Run("Error handling on ProductFeePercentage call fail", func(t *testing.T) {
		lbcMock.On("ProductFeePercentage", mock.Anything).Return(nil, assert.AnError).Once()
		feeCollector := rootstock.NewFeeCollectorImpl(lbcMock, rootstock.RetryParams{Retries: 0, Sleep: 0})
		percentage, err := feeCollector.DaoFeePercentage()
		require.Error(t, err)
		require.Zero(t, percentage)
	})
}
