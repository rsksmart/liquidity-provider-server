package rootstock_test

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	bindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/collateral_management"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
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

var penalizations = []types.Log{
	{
		Topics: []common.Hash{
			common.HexToHash("0x32d8dcdc3bd4d5d6dd9053c2e1d421c681715c97c6232e33a8658b7ae0bef13f"),
			common.HexToHash("0x00000000000000000000000079568c2989232dca1840087d73d403602364c0d4"),
			common.HexToHash("0x00000000000000000000000079568c2989232dca1840087d73d403602364c0d4"),
			common.HexToHash("0x0102030000000000000000000000000000000000000000000000000000000000"),
		},
		Data: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000022b0000000000000000000000000000000000000000000000000000000000000001"),
	},
	{
		Topics: []common.Hash{
			common.HexToHash("0x32d8dcdc3bd4d5d6dd9053c2e1d421c681715c97c6232e33a8658b7ae0bef13f"),
			common.HexToHash("0x00000000000000000000000079568c2989232dca1840087d73d403602364c0d4"),
			common.HexToHash("0x00000000000000000000000079568c2989232dca1840087d73d403602364c0d4"),
			common.HexToHash("0x0405060000000000000000000000000000000000000000000000000000000000"),
		},
		Data: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000029a0000000000000000000000000000000000000000000000000000000000000001"),
	},
	{
		Topics: []common.Hash{
			common.HexToHash("0x32d8dcdc3bd4d5d6dd9053c2e1d421c681715c97c6232e33a8658b7ae0bef13f"),
			common.HexToHash("0x00000000000000000000000079568c2989232dca1840087d73d403602364c0d4"),
			common.HexToHash("0x00000000000000000000000079568c2989232dca1840087d73d403602364c0d4"),
			common.HexToHash("0x0708090000000000000000000000000000000000000000000000000000000000"),
		},
		Data: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003090000000000000000000000000000000000000000000000000000000000000001"),
	},
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
	boundContract := bind.NewBoundContract(common.Address{}, abi.ABI{}, nil, nil, nil)
	contractBinding := bindings.NewCollateralManagementContract()
	contract := rootstock.NewCollateralManagementContractImpl(
		rootstock.NewRskClient(&mocks.RpcClientBindingMock{}),
		test.AnyAddress,
		test.AnyAddress,
		boundContract,
		&mocks.TransactionSignerMock{},
		contractBinding,
		rootstock.RetryParams{Retries: 1, Sleep: time.Duration(1)},
		time.Duration(1),
		Abis,
	)
	test.AssertNonZeroValues(t, contract)
}

func TestCollateralManagementContractImpl_GetAddress(t *testing.T) {
	collateral := rootstock.NewCollateralManagementContractImpl(dummyClient, test.AnyAddress, test.AnyAddress, nil, nil, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
	assert.Equal(t, test.AnyAddress, collateral.GetAddress())
}

func TestCollateralManagementContractImpl_ProviderResign(t *testing.T) {
	contractMock := createBoundContractMock()
	collateralBinding := bindings.NewCollateralManagementContract()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	collateral := rootstock.NewCollateralManagementContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		test.AnyAddress,
		contractMock.contract,
		signerMock,
		collateralBinding,
		rootstock.RetryParams{},
		time.Duration(1),
		Abis,
	)
	t.Run("Success", func(t *testing.T) {
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 0, big.NewInt(0), collateralBinding.PackResign()),
		).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, true)
		err := collateral.ProviderResign()
		require.NoError(t, err)
		contractMock.transactor.AssertExpectations(t)
	})
	t.Run("Error handling when sending resign tx", func(t *testing.T) {
		signerMock.EXPECT().Sign(mock.Anything, mock.Anything).RunAndReturn(func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return tx, nil
		})
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 0, big.NewInt(0), collateralBinding.PackResign()),
		).Return(assert.AnError).Once()
		err := collateral.ProviderResign()
		require.Error(t, err)
		contractMock.transactor.AssertExpectations(t)
	})
	t.Run("Error handling (resign tx reverted)", func(t *testing.T) {
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 0, big.NewInt(0), collateralBinding.PackResign()),
		).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, false)
		err := collateral.ProviderResign()
		require.ErrorContains(t, err, "resign transaction failed")
		contractMock.transactor.AssertExpectations(t)
	})
}

func TestCollateralManagementContractImpl_GetCollateral(t *testing.T) {
	collateralBinding := bindings.NewCollateralManagementContract()
	contractMock := createBoundContractMock()
	collateral := rootstock.NewCollateralManagementContractImpl(dummyClient, test.AnyAddress, test.AnyAddress, contractMock.contract, nil, collateralBinding, rootstock.RetryParams{}, time.Duration(1), Abis)
	t.Run("Success", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(collateralBinding.PackGetPegInCollateral(parsedAddress)),
			mock.Anything,
		).Return(mustPackUint256(t, big.NewInt(500)), nil).Once()
		result, err := collateral.GetCollateral(parsedAddress.String())
		require.NoError(t, err)
		assert.Equal(t, entities.NewWei(500), result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on GetCollateral call error", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(collateralBinding.PackGetPegInCollateral(parsedAddress)),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
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
	collateralBinding := bindings.NewCollateralManagementContract()
	contractMock := createBoundContractMock()
	collateral := rootstock.NewCollateralManagementContractImpl(dummyClient, test.AnyAddress, test.AnyAddress, contractMock.contract, nil, collateralBinding, rootstock.RetryParams{}, time.Duration(1), Abis)
	t.Run("Success", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(collateralBinding.PackGetPegOutCollateral(parsedAddress)),
			mock.Anything,
		).Return(mustPackUint256(t, big.NewInt(500)), nil).Once()
		result, err := collateral.GetPegoutCollateral(parsedAddress.String())
		require.NoError(t, err)
		assert.Equal(t, entities.NewWei(500), result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on GetPegoutCollateral call error", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(collateralBinding.PackGetPegOutCollateral(parsedAddress)),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
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
	collateralBinding := bindings.NewCollateralManagementContract()
	contractMock := createBoundContractMock()
	collateral := rootstock.NewCollateralManagementContractImpl(dummyClient, test.AnyAddress, test.AnyAddress, contractMock.contract, nil, collateralBinding, rootstock.RetryParams{}, time.Duration(1), Abis)
	t.Run("Success", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(collateralBinding.PackGetMinCollateral()),
			mock.Anything,
		).Return(mustPackUint256(t, big.NewInt(500)), nil).Once()
		result, err := collateral.GetMinimumCollateral()
		require.NoError(t, err)
		assert.Equal(t, entities.NewWei(500), result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on GetMinCollateral call fail", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(collateralBinding.PackGetMinCollateral()),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := collateral.GetMinimumCollateral()
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestCollateralManagementContractImpl_AddCollateral(t *testing.T) {
	collateralBinding := bindings.NewCollateralManagementContract()
	contractMock := createBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	collateral := rootstock.NewCollateralManagementContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		test.AnyAddress,
		contractMock.contract,
		signerMock,
		collateralBinding,
		rootstock.RetryParams{},
		time.Duration(1),
		Abis,
	)
	t.Run("Success", func(t *testing.T) {
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 0, big.NewInt(500), collateralBinding.PackAddPegInCollateral()),
		).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, true)
		err := collateral.AddCollateral(entities.NewWei(500))
		require.NoError(t, err)
		contractMock.transactor.AssertExpectations(t)
	})
	t.Run("Error handling when sending addCollateral tx", func(t *testing.T) {
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 0, big.NewInt(500), collateralBinding.PackAddPegInCollateral()),
		).Return(assert.AnError).Once()
		signerMock.EXPECT().Sign(mock.Anything, mock.Anything).RunAndReturn(func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return tx, nil
		})
		err := collateral.AddCollateral(entities.NewWei(500))
		require.Error(t, err)
		contractMock.transactor.AssertExpectations(t)
	})
	t.Run("Error handling (addCollateral tx reverted)", func(t *testing.T) {
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 0, big.NewInt(500), collateralBinding.PackAddPegInCollateral()),
		).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, false)
		err := collateral.AddCollateral(entities.NewWei(500))
		require.ErrorContains(t, err, "error adding pegin collateral")
		contractMock.transactor.AssertExpectations(t)
	})
}

func TestLiquidityBridgeContractImpl_AddPegoutCollateral(t *testing.T) {
	collateralBinding := bindings.NewCollateralManagementContract()
	contractMock := createBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	collateral := rootstock.NewCollateralManagementContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		test.AnyAddress,
		contractMock.contract,
		signerMock,
		collateralBinding,
		rootstock.RetryParams{},
		time.Duration(1),
		Abis,
	)
	t.Run("Success", func(t *testing.T) {
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 0, big.NewInt(777), collateralBinding.PackAddPegOutCollateral()),
		).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, true)
		err := collateral.AddPegoutCollateral(entities.NewWei(777))
		require.NoError(t, err)
		contractMock.transactor.AssertExpectations(t)
	})
	t.Run("Error handling when sending addPegoutCollateral tx", func(t *testing.T) {
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 0, big.NewInt(777), collateralBinding.PackAddPegOutCollateral()),
		).Return(assert.AnError).Once()
		signerMock.EXPECT().Sign(mock.Anything, mock.Anything).RunAndReturn(func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return tx, nil
		})
		err := collateral.AddPegoutCollateral(entities.NewWei(777))
		require.Error(t, err)
		contractMock.transactor.AssertExpectations(t)
	})
	t.Run("Error handling (addPegoutCollateral tx reverted)", func(t *testing.T) {
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 0, big.NewInt(777), collateralBinding.PackAddPegOutCollateral()),
		).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, false)
		err := collateral.AddPegoutCollateral(entities.NewWei(777))
		require.ErrorContains(t, err, "error adding pegout collateral")
		contractMock.transactor.AssertExpectations(t)
	})
}

// nolint:funlen
func TestLiquidityBridgeContractImpl_WithdrawCollateral(t *testing.T) {
	collateralBinding := bindings.NewCollateralManagementContract()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	t.Run("Success", func(t *testing.T) {
		contractMock := createBoundContractMock()
		collateral := rootstock.NewCollateralManagementContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, test.AnyAddress, contractMock.contract, signerMock, collateralBinding, rootstock.RetryParams{}, time.Duration(1), Abis)
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(collateralBinding.PackWithdrawCollateral()),
			mock.Anything,
		).Return(nil, nil).Once()
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 0, big.NewInt(0), collateralBinding.PackWithdrawCollateral()),
		).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, true)
		err := collateral.WithdrawCollateral()
		require.NoError(t, err)
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertExpectations(t)
	})
	t.Run("Error handling when sending withdrawCollateral tx", func(t *testing.T) {
		contractMock := createBoundContractMock()
		collateral := rootstock.NewCollateralManagementContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, test.AnyAddress, contractMock.contract, signerMock, collateralBinding, rootstock.RetryParams{}, time.Duration(1), Abis)
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(collateralBinding.PackWithdrawCollateral()),
			mock.Anything,
		).Return(nil, nil).Once()
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 0, big.NewInt(0), collateralBinding.PackWithdrawCollateral()),
		).Return(assert.AnError).Once()
		signerMock.EXPECT().Sign(mock.Anything, mock.Anything).RunAndReturn(func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return tx, nil
		})
		prepareTxMocks(&contractMock, mockClient, signerMock, true)
		err := collateral.WithdrawCollateral()
		require.Error(t, err)
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertExpectations(t)
	})
	t.Run("Error handling (withdrawCollateral tx reverted by panic)", func(t *testing.T) {
		contractMock := createBoundContractMock()
		collateral := rootstock.NewCollateralManagementContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, test.AnyAddress, contractMock.contract, signerMock, collateralBinding, rootstock.RetryParams{}, time.Duration(1), Abis)
		e := NewRskRpcError("division by zero", "0x4e487b710000000000000000000000000000000000000000000000000000000000000012")
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(collateralBinding.PackWithdrawCollateral()),
			mock.Anything,
		).Return(nil, e).Once()
		err := collateral.WithdrawCollateral()
		require.ErrorContains(t, err, "error parsing withdrawCollateral result")
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling (withdrawCollateral tx reverted by not resigned)", func(t *testing.T) {
		contractMock := createBoundContractMock()
		collateral := rootstock.NewCollateralManagementContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, test.AnyAddress, contractMock.contract, signerMock, collateralBinding, rootstock.RetryParams{}, time.Duration(1), Abis)
		e := NewRskRpcError("transaction reverted", "0x977254570000000000000000000000002279b7a0a67db372996a5fab50d91eaa73d2ebe6")
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(collateralBinding.PackWithdrawCollateral()),
			mock.Anything,
		).Return(nil, e).Once()
		err := collateral.WithdrawCollateral()
		require.ErrorContains(t, err, "provided hasn't completed resignation process")
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling (withdrawCollateral tx reverted by delay not passed)", func(t *testing.T) {
		contractMock := createBoundContractMock()
		collateral := rootstock.NewCollateralManagementContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, test.AnyAddress, contractMock.contract, signerMock, collateralBinding, rootstock.RetryParams{}, time.Duration(1), Abis)
		e := NewRskRpcError("transaction reverted", "0xf6cf33350000000000000000000000002279b7a0a67db372996a5fab50d91eaa73d2ebe600000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002")
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(collateralBinding.PackWithdrawCollateral()),
			mock.Anything,
		).Return(nil, e).Once()
		err := collateral.WithdrawCollateral()
		require.ErrorContains(t, err, "provided hasn't completed resignation process")
		contractMock.caller.AssertExpectations(t)
	})
}

func TestLiquidityBridgeContractImpl_GetPunishmentEvents(t *testing.T) {
	collateralBinding := bindings.NewCollateralManagementContract()
	contractMock := createBoundContractMock()
	collateral := rootstock.NewCollateralManagementContractImpl(dummyClient, test.AnyRskAddress, test.AnyAddress, contractMock.contract, nil, collateralBinding, rootstock.RetryParams{}, time.Duration(1), Abis)
	t.Run("Success", func(t *testing.T) {
		var from uint64 = 500
		var to uint64 = 1000
		contractMock.filterer.EXPECT().FilterLogs(mock.Anything, mock.MatchedBy(filterMatchFunc(from, to))).Return(penalizations, nil).Once()
		result, err := collateral.GetPenalizedEvents(context.Background(), from, &to)
		require.NoError(t, err)
		assert.Equal(t, parsedPenalizations, result)
		contractMock.filterer.AssertExpectations(t)
	})
	t.Run("Error handling when failed to get events", func(t *testing.T) {
		var from uint64 = 600
		var to uint64 = 1100
		contractMock.filterer.EXPECT().FilterLogs(mock.Anything, mock.MatchedBy(filterMatchFunc(from, to))).Return(nil, assert.AnError).Once()
		result, err := collateral.GetPenalizedEvents(context.Background(), from, &to)
		require.Error(t, err)
		assert.Nil(t, result)
		contractMock.filterer.AssertExpectations(t)
	})
}

func TestCollateralManagementContractImpl_PausedStatus(t *testing.T) {
	contractMock := createBoundContractMock()
	collateralBinding := bindings.NewCollateralManagementContract()
	contract := rootstock.NewCollateralManagementContractImpl(dummyClient, test.AnyRskAddress, test.AnyAddress, contractMock.contract, nil, collateralBinding, rootstock.RetryParams{}, time.Duration(1), Abis)
	t.Run("should return pause status result", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(collateralBinding.PackPauseStatus()),
			mock.Anything,
		).Return(mustPackPauseStatus(t, generalPauseStatus{IsPaused: true, Reason: "test", Since: 123}), nil).Once()
		result, err := contract.PausedStatus()
		require.NoError(t, err)
		assert.Equal(t, blockchain.PauseStatus{IsPaused: true, Reason: "test", Since: 123}, result)
	})
	t.Run("should handle error checking pause status", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(collateralBinding.PackPauseStatus()),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := contract.PausedStatus()
		require.Error(t, err)
		assert.Empty(t, result)
	})
	contractMock.caller.AssertExpectations(t)
}
