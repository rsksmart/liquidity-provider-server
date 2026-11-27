package pegin_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	bindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin"
	commitFirstBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_commit_first"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/rsksmart/liquidity-provider-server/test/rootstocktest"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewPeginContractImpl(t *testing.T) {
	contract := rootstock.NewPeginContractImpl(
		rootstock.NewRskClient(&mocks.RpcClientBindingMock{}),
		test.AnyAddress,
		rootstocktest.NewBoundContractMock().Contract,
		&mocks.TransactionSignerMock{},
		rootstock.RetryParams{Retries: 1, Sleep: 1},
		time.Duration(1),
		bindings.NewPeginContract(),
		commitFirstBindings.NewPeginCommitFirstContract(),
		rootstocktest.Abis,
	)
	test.AssertNonZeroValues(t, contract)
}

// nolint:funlen
func TestPeginContractImpl_GetAddress(t *testing.T) {
	peginContract := rootstock.NewPeginContractImpl(rootstocktest.DummyClient, test.AnyAddress, nil, nil, rootstock.RetryParams{}, time.Duration(1), nil, nil, rootstocktest.Abis)
	assert.Equal(t, test.AnyAddress, peginContract.GetAddress())
}
