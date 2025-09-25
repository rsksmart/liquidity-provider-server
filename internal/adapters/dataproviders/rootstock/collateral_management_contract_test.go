package rootstock_test

import (
	"bytes"
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

var penalizations = []*bindings.ICollateralManagementPenalized{
	{QuoteHash: [32]byte{1, 2, 3}, LiquidityProvider: common.HexToAddress(test.AnyRskAddress), Penalty: big.NewInt(555)},
	{QuoteHash: [32]byte{4, 5, 6}, LiquidityProvider: common.HexToAddress(test.AnyRskAddress), Penalty: big.NewInt(666)},
	{QuoteHash: [32]byte{7, 8, 9}, LiquidityProvider: common.HexToAddress(test.AnyRskAddress), Penalty: big.NewInt(777)},
}

var parsedPenalizations = []penalization.PenalizedEvent{
	{
		QuoteHash:         "0102030000000000000000000000000000000000000000000000000000000000",
		Penalty:           entities.NewWei(555),
		LiquidityProvider: test.AnyRskAddress,
	},
	{
		QuoteHash:         "0405060000000000000000000000000000000000000000000000000000000000",
		Penalty:           entities.NewWei(666),
		LiquidityProvider: test.AnyRskAddress,
	},
	{
		QuoteHash:         "0708090000000000000000000000000000000000000000000000000000000000",
		Penalty:           entities.NewWei(777),
		LiquidityProvider: test.AnyRskAddress,
	},
}

func TestNewCollateralManagementContractImpl(t *testing.T) {
	contract := rootstock.NewCollateralManagementContractImpl(
		rootstock.NewRskClient(&mocks.RpcClientBindingMock{}),
		test.AnyAddress,
		test.AnyAddress,
		&mocks.CollateralManagementAdapterMock{},
		&mocks.TransactionSignerMock{},
		rootstock.RetryParams{Retries: 1, Sleep: time.Duration(1)},
		time.Duration(1),
		Abis,
	)
	test.AssertNonZeroValues(t, contract)
}

func TestCollateralManagementContractImpl_GetAddress(t *testing.T) {
	collateral := rootstock.NewCollateralManagementContractImpl(dummyClient, test.AnyAddress, test.AnyAddress, nil, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
	assert.Equal(t, test.AnyAddress, collateral.GetAddress())
}

func TestCollateralManagementContractImpl_ProviderResign(t *testing.T) {
	contractBinding := &mocks.CollateralManagementAdapterMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	collateral := rootstock.NewCollateralManagementContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		test.AnyAddress,
		contractBinding,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
		Abis,
	)
	t.Run("Success", func(t *testing.T) {
		tx := prepareTxMocks(mockClient, signerMock, true)
		contractBinding.On("Resign", mock.Anything).Return(tx, nil).Once()
		err := collateral.ProviderResign()
		require.NoError(t, err)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling when sending resign tx", func(t *testing.T) {
		_ = prepareTxMocks(mockClient, signerMock, true)
		contractBinding.On("Resign", mock.Anything).Return(nil, assert.AnError).Once()
		err := collateral.ProviderResign()
		require.Error(t, err)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling (resign tx reverted)", func(t *testing.T) {
		tx := prepareTxMocks(mockClient, signerMock, false)
		contractBinding.On("Resign", mock.Anything).Return(tx, nil).Once()
		err := collateral.ProviderResign()
		require.ErrorContains(t, err, "resign transaction failed")
		contractBinding.AssertExpectations(t)
	})
}

func TestCollateralManagementContractImpl_GetCollateral(t *testing.T) {
	contractBinding := &mocks.CollateralManagementAdapterMock{}
	collateral := rootstock.NewCollateralManagementContractImpl(dummyClient, test.AnyAddress, test.AnyAddress, contractBinding, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
	t.Run("Success", func(t *testing.T) {
		contractBinding.EXPECT().GetPegInCollateral(mock.Anything, parsedAddress).Return(big.NewInt(500), nil).Once()
		result, err := collateral.GetCollateral(parsedAddress.String())
		require.NoError(t, err)
		assert.Equal(t, entities.NewWei(500), result)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling on GetCollateral call error", func(t *testing.T) {
		contractBinding.EXPECT().GetPegInCollateral(mock.Anything, parsedAddress).Return(nil, assert.AnError).Once()
		result, err := collateral.GetCollateral(parsedAddress.String())
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Error handling on invalid address for getting collateral", func(t *testing.T) {
		result, err := collateral.GetCollateral(test.AnyString)
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestLiquidityBridgeContractImpl_GetPegoutCollateral(t *testing.T) {
	contractBinding := &mocks.CollateralManagementAdapterMock{}
	collateral := rootstock.NewCollateralManagementContractImpl(dummyClient, test.AnyAddress, test.AnyAddress, contractBinding, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
	t.Run("Success", func(t *testing.T) {
		contractBinding.EXPECT().GetPegOutCollateral(mock.Anything, parsedAddress).Return(big.NewInt(500), nil).Once()
		result, err := collateral.GetPegoutCollateral(parsedAddress.String())
		require.NoError(t, err)
		assert.Equal(t, entities.NewWei(500), result)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling on GetPegoutCollateral call error", func(t *testing.T) {
		contractBinding.EXPECT().GetPegOutCollateral(mock.Anything, parsedAddress).Return(nil, assert.AnError).Once()
		result, err := collateral.GetPegoutCollateral(parsedAddress.String())
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Error handling on invalid address for getting pegout collateral", func(t *testing.T) {
		result, err := collateral.GetPegoutCollateral(test.AnyString)
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestCollateralManagementContractImpl_GetMinimumCollateral(t *testing.T) {
	contractBinding := &mocks.CollateralManagementAdapterMock{}
	collateral := rootstock.NewCollateralManagementContractImpl(dummyClient, test.AnyAddress, test.AnyAddress, contractBinding, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
	t.Run("Success", func(t *testing.T) {
		contractBinding.On("GetMinCollateral", mock.Anything).Return(big.NewInt(500), nil).Once()
		result, err := collateral.GetMinimumCollateral()
		require.NoError(t, err)
		assert.Equal(t, entities.NewWei(500), result)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling on GetMinCollateral call fail", func(t *testing.T) {
		contractBinding.On("GetMinCollateral", mock.Anything).Return(nil, assert.AnError).Once()
		result, err := collateral.GetMinimumCollateral()
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestCollateralManagementContractImpl_AddCollateral(t *testing.T) {
	contractBinding := &mocks.CollateralManagementAdapterMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	collateral := rootstock.NewCollateralManagementContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		test.AnyAddress,
		contractBinding,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
		Abis,
	)
	txMatchFunction := mock.MatchedBy(func(opts *bind.TransactOpts) bool {
		return opts.Value.Cmp(big.NewInt(500)) == 0 && bytes.Equal(opts.From.Bytes(), parsedAddress.Bytes())
	})
	t.Run("Success", func(t *testing.T) {
		tx := prepareTxMocks(mockClient, signerMock, true, valueModifier(big.NewInt(500)))
		contractBinding.EXPECT().AddPegInCollateral(txMatchFunction).Return(tx, nil).Once()
		err := collateral.AddCollateral(entities.NewWei(500))
		require.NoError(t, err)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling when sending addCollateral tx", func(t *testing.T) {
		_ = prepareTxMocks(mockClient, signerMock, true)
		contractBinding.EXPECT().AddPegInCollateral(txMatchFunction).Return(nil, assert.AnError).Once()
		err := collateral.AddCollateral(entities.NewWei(500))
		require.Error(t, err)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling (addCollateral tx reverted)", func(t *testing.T) {
		tx := prepareTxMocks(mockClient, signerMock, false, valueModifier(big.NewInt(500)))
		contractBinding.EXPECT().AddPegInCollateral(txMatchFunction).Return(tx, nil).Once()
		err := collateral.AddCollateral(entities.NewWei(500))
		require.ErrorContains(t, err, "error adding pegin collateral")
		contractBinding.AssertExpectations(t)
	})
}

func TestLiquidityBridgeContractImpl_AddPegoutCollateral(t *testing.T) {
	contractBinding := &mocks.CollateralManagementAdapterMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	collateral := rootstock.NewCollateralManagementContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		test.AnyAddress,
		contractBinding,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
		Abis,
	)
	txMatchFunction := mock.MatchedBy(func(opts *bind.TransactOpts) bool {
		return opts.Value.Cmp(big.NewInt(777)) == 0 && bytes.Equal(opts.From.Bytes(), parsedAddress.Bytes())
	})
	t.Run("Success", func(t *testing.T) {
		tx := prepareTxMocks(mockClient, signerMock, true, valueModifier(big.NewInt(777)))
		contractBinding.EXPECT().AddPegOutCollateral(txMatchFunction).Return(tx, nil).Once()
		err := collateral.AddPegoutCollateral(entities.NewWei(777))
		require.NoError(t, err)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling when sending addPegoutCollateral tx", func(t *testing.T) {
		_ = prepareTxMocks(mockClient, signerMock, true)
		contractBinding.EXPECT().AddPegOutCollateral(txMatchFunction).Return(nil, assert.AnError).Once()
		err := collateral.AddPegoutCollateral(entities.NewWei(777))
		require.Error(t, err)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling (addPegoutCollateral tx reverted)", func(t *testing.T) {
		tx := prepareTxMocks(mockClient, signerMock, false, valueModifier(big.NewInt(777)))
		contractBinding.EXPECT().AddPegOutCollateral(txMatchFunction).Return(tx, nil).Once()
		err := collateral.AddPegoutCollateral(entities.NewWei(777))
		require.ErrorContains(t, err, "error adding pegout collateral")
		contractBinding.AssertExpectations(t)
	})
}

func TestLiquidityBridgeContractImpl_WithdrawCollateral(t *testing.T) {
	const functionName = "withdrawCollateral"
	contractBinding := &mocks.CollateralManagementAdapterMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	collateral := rootstock.NewCollateralManagementContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, test.AnyAddress, contractBinding, signerMock, rootstock.RetryParams{}, time.Duration(1), Abis)
	t.Run("Success", func(t *testing.T) {
		tx := prepareTxMocks(mockClient, signerMock, true)
		callerMock := &mocks.ContractCallerBindingMock{}
		callerMock.EXPECT().Call(mock.Anything, mock.Anything, functionName).Return(nil).Once()
		contractBinding.EXPECT().Caller().Return(callerMock)
		contractBinding.EXPECT().WithdrawCollateral(mock.Anything).Return(tx, nil).Once()
		err := collateral.WithdrawCollateral()
		require.NoError(t, err)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling when sending withdrawCollateral tx", func(t *testing.T) {
		_ = prepareTxMocks(mockClient, signerMock, true)
		contractBinding.Calls = []mock.Call{}
		contractBinding.ExpectedCalls = []*mock.Call{}
		callerMock := &mocks.ContractCallerBindingMock{}
		callerMock.EXPECT().Call(mock.Anything, mock.Anything, functionName).Return(nil).Once()
		contractBinding.EXPECT().Caller().Return(callerMock)
		contractBinding.EXPECT().WithdrawCollateral(mock.Anything).Return(nil, assert.AnError).Once()
		err := collateral.WithdrawCollateral()
		require.Error(t, err)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling (withdrawCollateral tx reverted by panic)", func(t *testing.T) {
		contractBinding.Calls = []mock.Call{}
		contractBinding.ExpectedCalls = []*mock.Call{}
		e := NewRskRpcError("division by zero", "0x4e487b710000000000000000000000000000000000000000000000000000000000000012")
		callerMock := &mocks.ContractCallerBindingMock{}
		callerMock.EXPECT().Call(mock.Anything, mock.Anything, functionName).Return(e).Once()
		contractBinding.EXPECT().Caller().Return(callerMock)
		err := collateral.WithdrawCollateral()
		require.ErrorContains(t, err, "error parsing withdrawCollateral result")
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling (withdrawCollateral tx reverted by not resigned)", func(t *testing.T) {
		contractBinding.Calls = []mock.Call{}
		contractBinding.ExpectedCalls = []*mock.Call{}
		e := NewRskRpcError("transaction reverted", "0x977254570000000000000000000000002279b7a0a67db372996a5fab50d91eaa73d2ebe6")
		callerMock := &mocks.ContractCallerBindingMock{}
		callerMock.EXPECT().Call(mock.Anything, mock.Anything, functionName).Return(e).Once()
		contractBinding.EXPECT().Caller().Return(callerMock)
		err := collateral.WithdrawCollateral()
		require.ErrorContains(t, err, "error parsing withdrawCollateral result")
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling (withdrawCollateral tx reverted by delay not passed)", func(t *testing.T) {
		contractBinding.Calls = []mock.Call{}
		contractBinding.ExpectedCalls = []*mock.Call{}
		e := NewRskRpcError("transaction reverted", "0xf6cf33350000000000000000000000002279b7a0a67db372996a5fab50d91eaa73d2ebe600000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002")
		callerMock := &mocks.ContractCallerBindingMock{}
		callerMock.EXPECT().Call(mock.Anything, mock.Anything, functionName).Return(e).Once()
		contractBinding.EXPECT().Caller().Return(callerMock)
		err := collateral.WithdrawCollateral()
		require.ErrorContains(t, err, "error parsing withdrawCollateral result")
		contractBinding.AssertExpectations(t)
	})
}

func TestLiquidityBridgeContractImpl_GetPunishmentEvents(t *testing.T) {
	contractBinding := &mocks.CollateralManagementAdapterMock{}
	iteratorMock := &mocks.EventIteratorAdapterMock[bindings.ICollateralManagementPenalized]{}
	filterMatchFunc := func(from uint64, to uint64) func(opts *bind.FilterOpts) bool {
		return func(opts *bind.FilterOpts) bool {
			return from == opts.Start && to == *opts.End && opts.Context != nil
		}
	}
	collateral := rootstock.NewCollateralManagementContractImpl(dummyClient, test.AnyRskAddress, test.AnyAddress, contractBinding, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
	t.Run("Success", func(t *testing.T) {
		var from uint64 = 500
		var to uint64 = 1000
		contractBinding.EXPECT().FilterPenalized(mock.MatchedBy(filterMatchFunc(from, to)), []common.Address{common.HexToAddress(test.AnyRskAddress)}, []common.Address(nil), [][32]uint8(nil)).
			Return(&bindings.ICollateralManagementPenalizedIterator{}, nil).Once()
		contractBinding.On("PenalizedEventIteratorAdapter", mock.AnythingOfType(penalizedIteratorString)).
			Return(iteratorMock)
		iteratorMock.On("Next").Return(true).Times(len(penalizations))
		iteratorMock.On("Next").Return(false).Once()
		for _, deposit := range penalizations {
			iteratorMock.On("Event").Return(deposit).Once()
		}
		iteratorMock.On("Error").Return(nil).Once()
		iteratorMock.On("Close").Return(nil).Once()
		result, err := collateral.GetPenalizedEvents(context.Background(), from, &to)
		require.NoError(t, err)
		assert.Equal(t, parsedPenalizations, result)
		contractBinding.AssertExpectations(t)
		iteratorMock.AssertExpectations(t)
	})
	t.Run("Error handling when failed to get iterator", func(t *testing.T) {
		var from uint64 = 600
		var to uint64 = 1100
		contractBinding.EXPECT().FilterPenalized(mock.MatchedBy(filterMatchFunc(from, to)), []common.Address{common.HexToAddress(test.AnyRskAddress)}, []common.Address(nil), [][32]uint8(nil)).
			Return(nil, assert.AnError).Once()
		contractBinding.On("PenalizedEventIteratorAdapter", mock.AnythingOfType(penalizedIteratorString)).
			Return(nil)
		result, err := collateral.GetPenalizedEvents(context.Background(), from, &to)
		require.Error(t, err)
		assert.Nil(t, result)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling on iterator error", func(t *testing.T) {
		var from uint64 = 700
		var to uint64 = 1200
		contractBinding.EXPECT().FilterPenalized(mock.MatchedBy(filterMatchFunc(from, to)), []common.Address{common.HexToAddress(test.AnyRskAddress)}, []common.Address(nil), [][32]uint8(nil)).
			Return(&bindings.ICollateralManagementPenalizedIterator{}, nil).Once()
		contractBinding.On("PenalizedEventIteratorAdapter", mock.AnythingOfType(penalizedIteratorString)).
			Return(iteratorMock)
		iteratorMock.On("Next").Return(false).Once()
		iteratorMock.On("Error").Return(assert.AnError).Once()
		iteratorMock.On("Close").Return(nil).Once()
		result, err := collateral.GetPenalizedEvents(context.Background(), from, &to)
		require.Error(t, err)
		assert.Nil(t, result)
		contractBinding.AssertExpectations(t)
		iteratorMock.AssertExpectations(t)
	})
}
