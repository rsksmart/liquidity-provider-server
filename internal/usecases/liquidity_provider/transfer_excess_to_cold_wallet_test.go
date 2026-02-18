package liquidity_provider_test

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/cold_wallet"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	// Wei conversion - 1 BTC/ETH = 10^18 wei
	oneEtherInWei = 1000000000000000000

	// Test configuration values
	testMaxLiquidityBtc              = 40     // Total liquidity across both networks (20 per network)
	testExcessTolerancePercent       = 20     // 20% tolerance
	testBtcMinTransferFeeMultiplier  = 5      // Minimum BTC transfer must be 5x the fee
	testRbtcMinTransferFeeMultiplier = 100    // Minimum RBTC transfer must be 100x the fee
	testForceTransferAfterSeconds    = 604800 // 1 week in seconds (7 * 24 * 60 * 60)

	// With 20% tolerance, threshold = 20 + (20 * 0.2) = 24 BTC per network
	testLiquidityAmountWithExcess     = "24500000000000000000" // 24.5 BTC/RBTC (above 24 threshold → has excess)
	testLiquidityAmountWithoutExcess  = 20                     // 20 BTC/RBTC (at target → no excess)
	testLiquidityAmountBelowThreshold = 21                     // 21 BTC/RBTC (above target but below threshold)
	testExcessAmount                  = "4500000000000000000"  // 4.5 BTC/RBTC (excess: 24.5 - 20)
	testExcessAmountBelowThreshold    = 1                      // 1 BTC/RBTC (excess: 21 - 20, forced by time)
	testBtcFeeAmount                  = 50000000000000         // BTC transaction fee in wei
)

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_HappyPathBtcExcess(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))

	btcLiquidityBig, _ := new(big.Int).SetString(testLiquidityAmountWithExcess, 10)
	btcLiquidity := entities.NewBigWei(btcLiquidityBig)
	rbtcLiquidity := new(entities.Wei).Mul(entities.NewWei(testLiquidityAmountWithoutExcess), entities.NewWei(oneEtherInWei))

	btcExcessBig, _ := new(big.Int).SetString(testExcessAmount, 10)
	btcExcess := entities.NewBigWei(btcExcessBig)

	btcFee := entities.NewWei(testBtcFeeAmount)

	btcTxHash := "btc_tx_hash_123"

	now := time.Now()
	nowUnix := now.Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &nowUnix,
		LastRbtcToColdWalletTransfer: &nowUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(btcLiquidity, nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(rbtcLiquidity, nil)

	btcWallet.On("EstimateTxFees", "cold_btc_address", btcExcess).Return(blockchain.BtcFeeEstimation{
		Value: btcFee,
	}, nil)
	btcWallet.On("Send", "cold_btc_address", btcExcess).Return(blockchain.BitcoinTransactionResult{
		Hash: btcTxHash,
		Fee:  btcFee,
	}, nil)

	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.Anything).Once()
	hashMock := &mocks.HashMock{}
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6}).Once()
	rskWallet.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil).Once()
	lpRepository.On("UpsertStateConfiguration", mock.Anything, mock.Anything).Return(nil).Once()

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, liquidity_provider.TransferStatusSuccess, result.BtcResult.Status)
	assert.Equal(t, btcTxHash, result.BtcResult.TxHash)
	assert.Equal(t, btcExcess.String(), result.BtcResult.Amount.String())
	assert.Equal(t, btcFee, result.BtcResult.Fee)
	require.NoError(t, result.BtcResult.Error)

	assert.Equal(t, liquidity_provider.TransferStatusSkippedNoExcess, result.RskResult.Status)
	assert.Empty(t, result.RskResult.TxHash)
	assert.Nil(t, result.RskResult.Amount)
	assert.Nil(t, result.RskResult.Fee)
	require.NoError(t, result.RskResult.Error)

	// Assert BtcTransferredDueToThresholdEvent was published with correct data
	eventBus.AssertCalled(t, "Publish", mock.MatchedBy(func(event cold_wallet.BtcTransferredDueToThresholdEvent) bool {
		return assert.Equal(t, cold_wallet.BtcTransferredDueToThresholdEventId, event.Event.Id()) &&
			assert.Equal(t, btcExcess.String(), event.Amount.String()) &&
			assert.Equal(t, btcTxHash, event.TxHash) &&
			assert.Equal(t, btcFee, event.Fee)
	}))

	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	btcWallet.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	hashMock.AssertExpectations(t)
	rskWallet.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_HappyPathRskExcess(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))

	btcLiquidity := new(entities.Wei).Mul(entities.NewWei(testLiquidityAmountWithoutExcess), entities.NewWei(oneEtherInWei))
	rbtcLiquidityBig, _ := new(big.Int).SetString(testLiquidityAmountWithExcess, 10)
	rbtcLiquidity := entities.NewBigWei(rbtcLiquidityBig)

	rbtcExcessBig, _ := new(big.Int).SetString(testExcessAmount, 10)
	rbtcExcess := entities.NewBigWei(rbtcExcessBig)

	rbtcGasPrice := entities.NewWei(1000000000)
	rbtcGasCost := new(entities.Wei).Mul(entities.NewWei(liquidity_provider.SimpleTransferGasLimit), rbtcGasPrice)
	rbtcAmountToTransfer := new(entities.Wei).Sub(rbtcExcess, rbtcGasCost)

	rskTxHash := "rsk_tx_hash_456"

	rskReceipt := blockchain.TransactionReceipt{
		TransactionHash: rskTxHash,
		GasUsed:         big.NewInt(liquidity_provider.SimpleTransferGasLimit),
		GasPrice:        rbtcGasPrice,
	}

	nowUnix := time.Now().Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &nowUnix,
		LastRbtcToColdWalletTransfer: &nowUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(btcLiquidity, nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(rbtcLiquidity, nil)

	rskRpcMock.On("GasPrice", ctx).Return(rbtcGasPrice, nil)
	rskWallet.On("SendRbtc", ctx,
		blockchain.NewTransactionConfig(
			rbtcAmountToTransfer,
			liquidity_provider.SimpleTransferGasLimit,
			rbtcGasPrice,
		),
		"cold_rsk_address",
	).Return(rskReceipt, nil)

	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.Anything).Once()
	hashMock := &mocks.HashMock{}
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6}).Once()
	rskWallet.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil).Once()
	lpRepository.On("UpsertStateConfiguration", mock.Anything, mock.Anything).Return(nil).Once()

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, liquidity_provider.TransferStatusSkippedNoExcess, result.BtcResult.Status)
	assert.Empty(t, result.BtcResult.TxHash)
	assert.Nil(t, result.BtcResult.Amount)
	assert.Nil(t, result.BtcResult.Fee)
	require.NoError(t, result.BtcResult.Error)

	rbtcActualFee := new(entities.Wei).Mul(rbtcGasPrice, entities.NewUWei(uint64(rskReceipt.GasUsed.Int64())))

	assert.Equal(t, liquidity_provider.TransferStatusSuccess, result.RskResult.Status)
	assert.Equal(t, rskTxHash, result.RskResult.TxHash)
	assert.Equal(t, rbtcAmountToTransfer.String(), result.RskResult.Amount.String())
	assert.Equal(t, rbtcActualFee.String(), result.RskResult.Fee.String())
	require.NoError(t, result.RskResult.Error)

	// Assert RbtcTransferredDueToThresholdEvent was published with correct data
	eventBus.AssertCalled(t, "Publish", mock.MatchedBy(func(event cold_wallet.RbtcTransferredDueToThresholdEvent) bool {
		return assert.Equal(t, cold_wallet.RbtcTransferredDueToThresholdEventId, event.Event.Id()) &&
			assert.Equal(t, rbtcAmountToTransfer.String(), event.Amount.String()) &&
			assert.Equal(t, rskTxHash, event.TxHash) &&
			assert.Equal(t, rbtcActualFee.String(), event.Fee.String())
	}))

	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	rskWallet.AssertExpectations(t)
	rskRpcMock.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	hashMock.AssertExpectations(t)
	rskWallet.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_HappyPathBothExcess(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))

	btcLiquidityBig, _ := new(big.Int).SetString(testLiquidityAmountWithExcess, 10)
	btcLiquidity := entities.NewBigWei(btcLiquidityBig)
	rbtcLiquidityBig, _ := new(big.Int).SetString(testLiquidityAmountWithExcess, 10)
	rbtcLiquidity := entities.NewBigWei(rbtcLiquidityBig)

	btcExcessBig, _ := new(big.Int).SetString(testExcessAmount, 10)
	btcExcess := entities.NewBigWei(btcExcessBig)
	rbtcExcessBig, _ := new(big.Int).SetString(testExcessAmount, 10)
	rbtcExcess := entities.NewBigWei(rbtcExcessBig)

	btcFee := entities.NewWei(testBtcFeeAmount)
	btcTxHash := "btc_tx_hash_123"

	rbtcGasPrice := entities.NewWei(1000000000)
	rbtcGasCost := new(entities.Wei).Mul(entities.NewWei(liquidity_provider.SimpleTransferGasLimit), rbtcGasPrice)
	rbtcAmountToTransfer := new(entities.Wei).Sub(rbtcExcess, rbtcGasCost)
	rskTxHash := "rsk_tx_hash_456"

	rskReceipt := blockchain.TransactionReceipt{
		TransactionHash: rskTxHash,
		GasUsed:         big.NewInt(liquidity_provider.SimpleTransferGasLimit),
		GasPrice:        rbtcGasPrice,
	}

	nowUnix := time.Now().Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &nowUnix,
		LastRbtcToColdWalletTransfer: &nowUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(btcLiquidity, nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(rbtcLiquidity, nil)

	btcWallet.On("EstimateTxFees", "cold_btc_address", btcExcess).Return(blockchain.BtcFeeEstimation{
		Value: btcFee,
	}, nil)
	btcWallet.On("Send", "cold_btc_address", btcExcess).Return(blockchain.BitcoinTransactionResult{
		Hash: btcTxHash,
		Fee:  btcFee,
	}, nil)

	rskRpcMock.On("GasPrice", ctx).Return(rbtcGasPrice, nil)
	rskWallet.On("SendRbtc", ctx,
		blockchain.NewTransactionConfig(
			rbtcAmountToTransfer,
			liquidity_provider.SimpleTransferGasLimit,
			rbtcGasPrice,
		),
		"cold_rsk_address",
	).Return(rskReceipt, nil)

	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.Anything).Times(2)
	hashMock := &mocks.HashMock{}
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6}).Once()
	rskWallet.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil).Once()
	lpRepository.On("UpsertStateConfiguration", mock.Anything, mock.Anything).Return(nil).Once()

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, liquidity_provider.TransferStatusSuccess, result.BtcResult.Status)
	assert.Equal(t, btcTxHash, result.BtcResult.TxHash)
	assert.Equal(t, btcExcess.String(), result.BtcResult.Amount.String())
	assert.Equal(t, btcFee, result.BtcResult.Fee)
	require.NoError(t, result.BtcResult.Error)

	rbtcActualFee := new(entities.Wei).Mul(rbtcGasPrice, entities.NewUWei(uint64(rskReceipt.GasUsed.Int64())))

	assert.Equal(t, liquidity_provider.TransferStatusSuccess, result.RskResult.Status)
	assert.Equal(t, rskTxHash, result.RskResult.TxHash)
	assert.Equal(t, rbtcAmountToTransfer.String(), result.RskResult.Amount.String())
	assert.Equal(t, rbtcActualFee.String(), result.RskResult.Fee.String())
	require.NoError(t, result.RskResult.Error)

	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	btcWallet.AssertExpectations(t)
	rskWallet.AssertExpectations(t)
	rskRpcMock.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	hashMock.AssertExpectations(t)
	rskWallet.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_NoExcess(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))

	btcLiquidity := new(entities.Wei).Mul(entities.NewWei(testLiquidityAmountWithoutExcess), entities.NewWei(oneEtherInWei))
	rbtcLiquidity := new(entities.Wei).Mul(entities.NewWei(testLiquidityAmountWithoutExcess), entities.NewWei(oneEtherInWei))

	nowUnix := time.Now().Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &nowUnix,
		LastRbtcToColdWalletTransfer: &nowUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(btcLiquidity, nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(rbtcLiquidity, nil)

	eventBus := new(mocks.EventBusMock)
	hashMock := &mocks.HashMock{}

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, liquidity_provider.TransferStatusSkippedNoExcess, result.BtcResult.Status)
	assert.Empty(t, result.BtcResult.TxHash)
	assert.Nil(t, result.BtcResult.Amount)
	assert.Nil(t, result.BtcResult.Fee)
	require.NoError(t, result.BtcResult.Error)

	assert.Equal(t, liquidity_provider.TransferStatusSkippedNoExcess, result.RskResult.Status)
	assert.Empty(t, result.RskResult.TxHash)
	assert.Nil(t, result.RskResult.Amount)
	assert.Nil(t, result.RskResult.Fee)
	require.NoError(t, result.RskResult.Error)

	eventBus.AssertNotCalled(t, "Publish")
	hashMock.AssertNotCalled(t, "Hash", mock.Anything)
	rskWallet.AssertNotCalled(t, "SignBytes", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)

	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_FixedToleranceInsteadOfPercentage(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))

	// Fixed tolerance: 2 BTC
	fixedTolerance := new(entities.Wei).Mul(entities.NewWei(2), entities.NewWei(oneEtherInWei))

	// Target: 20 BTC per network
	// Threshold with fixed tolerance: 20 + 2 = 22 BTC
	// BTC Liquidity: 23 BTC (above threshold) -> should transfer 3 BTC excess (23 - 20)
	// RBTC Liquidity: 21 BTC (below threshold) -> no transfer
	btcLiquidityBig, _ := new(big.Int).SetString("23000000000000000000", 10)
	btcLiquidity := entities.NewBigWei(btcLiquidityBig)

	rbtcLiquidityBig, _ := new(big.Int).SetString("21000000000000000000", 10)
	rbtcLiquidity := entities.NewBigWei(rbtcLiquidityBig)

	btcExcessBig, _ := new(big.Int).SetString("3000000000000000000", 10)
	btcExcess := entities.NewBigWei(btcExcessBig)
	btcFee := entities.NewWei(testBtcFeeAmount)
	btcTxHash := "btc_tx_hash_fixed_tolerance"

	nowUnix := time.Now().Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         true,
			PercentageValue: utils.NewBigFloat64(0),
			FixedValue:      fixedTolerance,
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &nowUnix,
		LastRbtcToColdWalletTransfer: &nowUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(btcLiquidity, nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(rbtcLiquidity, nil)

	btcWallet.On("EstimateTxFees", "cold_btc_address", btcExcess).Return(blockchain.BtcFeeEstimation{
		Value: btcFee,
	}, nil)
	btcWallet.On("Send", "cold_btc_address", btcExcess).Return(blockchain.BitcoinTransactionResult{
		Hash: btcTxHash,
		Fee:  btcFee,
	}, nil)

	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.Anything).Once()
	hashMock := &mocks.HashMock{}
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6}).Once()
	rskWallet.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil).Once()
	lpRepository.On("UpsertStateConfiguration", mock.Anything, mock.Anything).Return(nil).Once()

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)

	// BTC transfer should succeed with fixed tolerance
	assert.Equal(t, liquidity_provider.TransferStatusSuccess, result.BtcResult.Status)
	assert.Equal(t, btcTxHash, result.BtcResult.TxHash)
	assert.Equal(t, btcExcess, result.BtcResult.Amount)
	assert.Equal(t, btcFee, result.BtcResult.Fee)
	require.NoError(t, result.BtcResult.Error)

	// RBTC is below threshold (21 < 22), should be skipped
	assert.Equal(t, liquidity_provider.TransferStatusSkippedNoExcess, result.RskResult.Status)
	assert.Empty(t, result.RskResult.TxHash)
	assert.Nil(t, result.RskResult.Amount)
	assert.Nil(t, result.RskResult.Fee)
	require.NoError(t, result.RskResult.Error)

	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	btcWallet.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	hashMock.AssertExpectations(t)
	rskWallet.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_TimeForced(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))

	btcLiquidity := new(entities.Wei).Mul(entities.NewWei(testLiquidityAmountBelowThreshold), entities.NewWei(oneEtherInWei))
	rbtcLiquidity := new(entities.Wei).Mul(entities.NewWei(testLiquidityAmountBelowThreshold), entities.NewWei(oneEtherInWei))

	btcExcess := new(entities.Wei).Mul(entities.NewWei(testExcessAmountBelowThreshold), entities.NewWei(oneEtherInWei))
	rbtcExcess := new(entities.Wei).Mul(entities.NewWei(testExcessAmountBelowThreshold), entities.NewWei(oneEtherInWei))

	btcFee := entities.NewWei(testBtcFeeAmount)
	btcTxHash := "btc_tx_hash_123"

	rbtcGasPrice := entities.NewWei(1000000000)
	rbtcGasCost := new(entities.Wei).Mul(entities.NewWei(liquidity_provider.SimpleTransferGasLimit), rbtcGasPrice)
	rbtcAmountToTransfer := new(entities.Wei).Sub(rbtcExcess, rbtcGasCost)
	rskTxHash := "rsk_tx_hash_456"

	rskReceipt := blockchain.TransactionReceipt{
		TransactionHash: rskTxHash,
		GasUsed:         big.NewInt(liquidity_provider.SimpleTransferGasLimit),
		GasPrice:        rbtcGasPrice,
	}

	oldTransferTime := time.Now().Add(-time.Duration(testForceTransferAfterSeconds+3600) * time.Second)
	oldTransferUnix := oldTransferTime.Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &oldTransferUnix,
		LastRbtcToColdWalletTransfer: &oldTransferUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(btcLiquidity, nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(rbtcLiquidity, nil)

	btcWallet.On("EstimateTxFees", "cold_btc_address", btcExcess).Return(blockchain.BtcFeeEstimation{
		Value: btcFee,
	}, nil)
	btcWallet.On("Send", "cold_btc_address", btcExcess).Return(blockchain.BitcoinTransactionResult{
		Hash: btcTxHash,
		Fee:  btcFee,
	}, nil)

	rskRpcMock.On("GasPrice", ctx).Return(rbtcGasPrice, nil)
	rskWallet.On("SendRbtc", ctx,
		blockchain.NewTransactionConfig(
			rbtcAmountToTransfer,
			liquidity_provider.SimpleTransferGasLimit,
			rbtcGasPrice,
		),
		"cold_rsk_address",
	).Return(rskReceipt, nil)

	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.Anything).Times(2)
	hashMock := &mocks.HashMock{}
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6}).Once()
	rskWallet.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil).Once()
	lpRepository.On("UpsertStateConfiguration", mock.Anything, mock.Anything).Return(nil).Once()

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, liquidity_provider.TransferStatusSuccess, result.BtcResult.Status)
	assert.Equal(t, btcTxHash, result.BtcResult.TxHash)
	assert.Equal(t, btcExcess.String(), result.BtcResult.Amount.String())
	assert.Equal(t, btcFee, result.BtcResult.Fee)
	require.NoError(t, result.BtcResult.Error)

	rbtcActualFee := new(entities.Wei).Mul(rbtcGasPrice, entities.NewUWei(uint64(rskReceipt.GasUsed.Int64())))

	assert.Equal(t, liquidity_provider.TransferStatusSuccess, result.RskResult.Status)
	assert.Equal(t, rskTxHash, result.RskResult.TxHash)
	assert.Equal(t, rbtcAmountToTransfer.String(), result.RskResult.Amount.String())
	assert.Equal(t, rbtcActualFee.String(), result.RskResult.Fee.String())
	require.NoError(t, result.RskResult.Error)

	// Verify BtcTransferredDueToTimeForcingEvent was published
	eventBus.AssertCalled(t, "Publish", mock.MatchedBy(func(event cold_wallet.BtcTransferredDueToTimeForcingEvent) bool {
		return assert.Equal(t, cold_wallet.BtcTransferredDueToTimeForcingEventId, event.Event.Id()) &&
			assert.Equal(t, btcExcess.String(), event.Amount.String()) &&
			assert.Equal(t, btcTxHash, event.TxHash) &&
			assert.Equal(t, btcFee, event.Fee)
	}))

	// Verify RbtcTransferredDueToTimeForcingEvent was published
	eventBus.AssertCalled(t, "Publish", mock.MatchedBy(func(event cold_wallet.RbtcTransferredDueToTimeForcingEvent) bool {
		return assert.Equal(t, cold_wallet.RbtcTransferredDueToTimeForcingEventId, event.Event.Id()) &&
			assert.Equal(t, rbtcAmountToTransfer.String(), event.Amount.String()) &&
			assert.Equal(t, rskTxHash, event.TxHash) &&
			assert.Equal(t, rbtcActualFee.String(), event.Fee.String())
	}))

	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	btcWallet.AssertExpectations(t)
	rskWallet.AssertExpectations(t)
	rskRpcMock.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	hashMock.AssertExpectations(t)
	rskWallet.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_TimeForcedButNoExcess(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))

	btcLiquidity := new(entities.Wei).Mul(entities.NewWei(testLiquidityAmountWithoutExcess), entities.NewWei(oneEtherInWei))
	rbtcLiquidity := new(entities.Wei).Mul(entities.NewWei(testLiquidityAmountWithoutExcess), entities.NewWei(oneEtherInWei))

	oldTransferTime := time.Now().Add(-time.Duration(testForceTransferAfterSeconds+3600) * time.Second)
	oldTransferUnix := oldTransferTime.Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &oldTransferUnix,
		LastRbtcToColdWalletTransfer: &oldTransferUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(btcLiquidity, nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(rbtcLiquidity, nil)

	eventBus := new(mocks.EventBusMock)
	hashMock := &mocks.HashMock{}

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, liquidity_provider.TransferStatusSkippedNoExcess, result.BtcResult.Status)
	assert.Empty(t, result.BtcResult.TxHash)
	assert.Nil(t, result.BtcResult.Amount)
	assert.Nil(t, result.BtcResult.Fee)
	require.NoError(t, result.BtcResult.Error)

	assert.Equal(t, liquidity_provider.TransferStatusSkippedNoExcess, result.RskResult.Status)
	assert.Empty(t, result.RskResult.TxHash)
	assert.Nil(t, result.RskResult.Amount)
	assert.Nil(t, result.RskResult.Fee)
	require.NoError(t, result.RskResult.Error)

	eventBus.AssertNotCalled(t, "Publish")
	hashMock.AssertNotCalled(t, "Hash", mock.Anything)
	rskWallet.AssertNotCalled(t, "SignBytes", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)

	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_BtcTimeForcedRskThresholdExceeded(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))

	// BTC: 21 BTC (below 24 BTC threshold, but time forcing triggers transfer of 1 BTC)
	btcLiquidityBig, _ := new(big.Int).SetString("21000000000000000000", 10)
	btcLiquidity := entities.NewBigWei(btcLiquidityBig)
	btcExcessBig, _ := new(big.Int).SetString("1000000000000000000", 10)
	btcExcess := entities.NewBigWei(btcExcessBig)

	// RBTC: 24.5 BTC (exceeds 24 BTC threshold, transfers 4.5 BTC)
	rbtcLiquidityBig, _ := new(big.Int).SetString(testLiquidityAmountWithExcess, 10)
	rbtcLiquidity := entities.NewBigWei(rbtcLiquidityBig)
	rbtcExcessBig, _ := new(big.Int).SetString(testExcessAmount, 10)
	rbtcExcess := entities.NewBigWei(rbtcExcessBig)

	btcFee := entities.NewWei(testBtcFeeAmount)
	btcTxHash := "btc_tx_hash_time_forced"

	rbtcGasPrice := entities.NewWei(1000000000)
	rbtcGasCost := new(entities.Wei).Mul(entities.NewWei(liquidity_provider.SimpleTransferGasLimit), rbtcGasPrice)
	rbtcAmountToTransfer := new(entities.Wei).Sub(rbtcExcess, rbtcGasCost)
	rbtcGasUsed := uint64(21000)
	rbtcActualFee := new(entities.Wei).Mul(rbtcGasPrice, entities.NewUWei(rbtcGasUsed))
	rbtcTxHash := "rbtc_tx_hash_threshold"

	// BTC: old transfer time (time forcing enabled)
	// RBTC: recent transfer time (no time forcing)
	oldTransferTime := time.Now().Add(-time.Duration(testForceTransferAfterSeconds+3600) * time.Second)
	oldTransferUnix := oldTransferTime.Unix()
	nowUnix := time.Now().Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &oldTransferUnix,
		LastRbtcToColdWalletTransfer: &nowUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(btcLiquidity, nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(rbtcLiquidity, nil)

	btcWallet.On("EstimateTxFees", "cold_btc_address", btcExcess).Return(blockchain.BtcFeeEstimation{
		Value: btcFee,
	}, nil)
	btcWallet.On("Send", "cold_btc_address", btcExcess).Return(blockchain.BitcoinTransactionResult{
		Hash: btcTxHash,
		Fee:  btcFee,
	}, nil)

	rskRpcMock.On("GasPrice", ctx).Return(rbtcGasPrice, nil)
	rskWallet.On("SendRbtc", ctx, blockchain.NewTransactionConfig(rbtcAmountToTransfer, liquidity_provider.SimpleTransferGasLimit, rbtcGasPrice), "cold_rsk_address").Return(blockchain.TransactionReceipt{
		TransactionHash: rbtcTxHash,
		GasUsed:         big.NewInt(int64(rbtcGasUsed)),
		GasPrice:        rbtcGasPrice,
	}, nil)

	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.Anything).Times(2)
	hashMock := &mocks.HashMock{}
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6}).Once()
	rskWallet.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil).Once()
	lpRepository.On("UpsertStateConfiguration", mock.Anything, mock.Anything).Return(nil).Once()

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)

	// BTC transfer should succeed (time forced)
	assert.Equal(t, liquidity_provider.TransferStatusSuccess, result.BtcResult.Status)
	assert.Equal(t, btcTxHash, result.BtcResult.TxHash)
	assert.Equal(t, btcExcess, result.BtcResult.Amount)
	assert.Equal(t, btcFee, result.BtcResult.Fee)
	require.NoError(t, result.BtcResult.Error)

	// RBTC transfer should succeed (threshold exceeded)
	assert.Equal(t, liquidity_provider.TransferStatusSuccess, result.RskResult.Status)
	assert.Equal(t, rbtcTxHash, result.RskResult.TxHash)
	assert.Equal(t, rbtcAmountToTransfer, result.RskResult.Amount)
	assert.Equal(t, rbtcActualFee, result.RskResult.Fee)
	require.NoError(t, result.RskResult.Error)

	// Assert BtcTransferredDueToTimeForcingEvent was published with correct data
	eventBus.AssertCalled(t, "Publish", mock.MatchedBy(func(event cold_wallet.BtcTransferredDueToTimeForcingEvent) bool {
		return assert.Equal(t, cold_wallet.BtcTransferredDueToTimeForcingEventId, event.Event.Id()) &&
			assert.Equal(t, btcExcess.String(), event.Amount.String()) &&
			assert.Equal(t, btcTxHash, event.TxHash) &&
			assert.Equal(t, btcFee, event.Fee)
	}))

	// Assert RbtcTransferredDueToThresholdEvent was published with correct data
	eventBus.AssertCalled(t, "Publish", mock.MatchedBy(func(event cold_wallet.RbtcTransferredDueToThresholdEvent) bool {
		return assert.Equal(t, cold_wallet.RbtcTransferredDueToThresholdEventId, event.Event.Id()) &&
			assert.Equal(t, rbtcAmountToTransfer.String(), event.Amount.String()) &&
			assert.Equal(t, rbtcTxHash, event.TxHash) &&
			assert.Equal(t, rbtcActualFee.String(), event.Fee.String())
	}))

	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	btcWallet.AssertExpectations(t)
	rskRpcMock.AssertExpectations(t)
	rskWallet.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	hashMock.AssertExpectations(t)
	rskWallet.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_RskTimeForcedBtcThresholdExceeded(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))

	// BTC: 24.5 BTC (exceeds 24 BTC threshold, transfers 4.5 BTC)
	btcLiquidityBig, _ := new(big.Int).SetString(testLiquidityAmountWithExcess, 10)
	btcLiquidity := entities.NewBigWei(btcLiquidityBig)
	btcExcessBig, _ := new(big.Int).SetString(testExcessAmount, 10)
	btcExcess := entities.NewBigWei(btcExcessBig)

	// RBTC: 21 BTC (below 24 BTC threshold, but time forcing triggers transfer of 1 BTC)
	rbtcLiquidityBig, _ := new(big.Int).SetString("21000000000000000000", 10)
	rbtcLiquidity := entities.NewBigWei(rbtcLiquidityBig)
	rbtcExcessBig, _ := new(big.Int).SetString("1000000000000000000", 10)
	rbtcExcess := entities.NewBigWei(rbtcExcessBig)

	btcFee := entities.NewWei(testBtcFeeAmount)
	btcTxHash := "btc_tx_hash_threshold"

	rbtcGasPrice := entities.NewWei(1000000000)
	rbtcGasCost := new(entities.Wei).Mul(entities.NewWei(liquidity_provider.SimpleTransferGasLimit), rbtcGasPrice)
	rbtcAmountToTransfer := new(entities.Wei).Sub(rbtcExcess, rbtcGasCost)
	rbtcGasUsed := uint64(21000)
	rbtcActualFee := new(entities.Wei).Mul(rbtcGasPrice, entities.NewUWei(rbtcGasUsed))
	rbtcTxHash := "rbtc_tx_hash_time_forced"

	// BTC: recent transfer time (no time forcing)
	// RBTC: old transfer time (time forcing enabled)
	now := time.Now()
	nowUnix := now.Unix()
	oldTransferTime := time.Now().Add(-time.Duration(testForceTransferAfterSeconds+3600) * time.Second)
	oldTransferUnix := oldTransferTime.Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &nowUnix,
		LastRbtcToColdWalletTransfer: &oldTransferUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(btcLiquidity, nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(rbtcLiquidity, nil)

	btcWallet.On("EstimateTxFees", "cold_btc_address", btcExcess).Return(blockchain.BtcFeeEstimation{
		Value: btcFee,
	}, nil)
	btcWallet.On("Send", "cold_btc_address", btcExcess).Return(blockchain.BitcoinTransactionResult{
		Hash: btcTxHash,
		Fee:  btcFee,
	}, nil)

	rskRpcMock.On("GasPrice", ctx).Return(rbtcGasPrice, nil)
	rskWallet.On("SendRbtc", ctx, blockchain.NewTransactionConfig(rbtcAmountToTransfer, liquidity_provider.SimpleTransferGasLimit, rbtcGasPrice), "cold_rsk_address").Return(blockchain.TransactionReceipt{
		TransactionHash: rbtcTxHash,
		GasUsed:         big.NewInt(int64(rbtcGasUsed)),
		GasPrice:        rbtcGasPrice,
	}, nil)

	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.Anything).Times(2)
	hashMock := &mocks.HashMock{}
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6}).Once()
	rskWallet.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil).Once()
	lpRepository.On("UpsertStateConfiguration", mock.Anything, mock.Anything).Return(nil).Once()

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)

	// BTC transfer should succeed (threshold exceeded)
	assert.Equal(t, liquidity_provider.TransferStatusSuccess, result.BtcResult.Status)
	assert.Equal(t, btcTxHash, result.BtcResult.TxHash)
	assert.Equal(t, btcExcess, result.BtcResult.Amount)
	assert.Equal(t, btcFee, result.BtcResult.Fee)
	require.NoError(t, result.BtcResult.Error)

	// RBTC transfer should succeed (time forced)
	assert.Equal(t, liquidity_provider.TransferStatusSuccess, result.RskResult.Status)
	assert.Equal(t, rbtcTxHash, result.RskResult.TxHash)
	assert.Equal(t, rbtcAmountToTransfer, result.RskResult.Amount)
	assert.Equal(t, rbtcActualFee, result.RskResult.Fee)
	require.NoError(t, result.RskResult.Error)

	// Assert RbtcTransferredDueToTimeForcingEvent was published with correct data
	eventBus.AssertCalled(t, "Publish", mock.MatchedBy(func(event cold_wallet.RbtcTransferredDueToTimeForcingEvent) bool {
		return assert.Equal(t, cold_wallet.RbtcTransferredDueToTimeForcingEventId, event.Event.Id()) &&
			assert.Equal(t, rbtcAmountToTransfer.String(), event.Amount.String()) &&
			assert.Equal(t, rbtcTxHash, event.TxHash) &&
			assert.Equal(t, rbtcActualFee.String(), event.Fee.String())
	}))

	// Assert BtcTransferredDueToThresholdEvent was published with correct data
	eventBus.AssertCalled(t, "Publish", mock.MatchedBy(func(event cold_wallet.BtcTransferredDueToThresholdEvent) bool {
		return assert.Equal(t, cold_wallet.BtcTransferredDueToThresholdEventId, event.Event.Id()) &&
			assert.Equal(t, btcExcess.String(), event.Amount.String()) &&
			assert.Equal(t, btcTxHash, event.TxHash) &&
			assert.Equal(t, btcFee, event.Fee)
	}))

	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	btcWallet.AssertExpectations(t)
	rskRpcMock.AssertExpectations(t)
	rskWallet.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	hashMock.AssertExpectations(t)
	rskWallet.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
}

func TestTransferExcessToColdWalletUseCase_Run_BtcColdWalletAddressEmpty(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	coldWallet.On("GetBtcAddress").Return("")
	coldWallet.On("GetRskAddress").Return("0x1234567890abcdef")

	eventBus := new(mocks.EventBusMock)
	hashMock := &mocks.HashMock{}

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "cold wallet not configured")

	eventBus.AssertNotCalled(t, "Publish")
	hashMock.AssertNotCalled(t, "Hash", mock.Anything)
	rskWallet.AssertNotCalled(t, "SignBytes", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
	coldWallet.AssertExpectations(t)
}

func TestTransferExcessToColdWalletUseCase_Run_RskColdWalletAddressEmpty(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("")

	eventBus := new(mocks.EventBusMock)
	hashMock := &mocks.HashMock{}

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "cold wallet not configured")

	eventBus.AssertNotCalled(t, "Publish")
	hashMock.AssertNotCalled(t, "Hash", mock.Anything)
	rskWallet.AssertNotCalled(t, "SignBytes", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
	coldWallet.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_MaxLiquidityNotConfigured(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: nil,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	eventBus := new(mocks.EventBusMock)
	hashMock := &mocks.HashMock{}

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "max liquidity not configured")

	eventBus.AssertNotCalled(t, "Publish")
	hashMock.AssertNotCalled(t, "Hash", mock.Anything)
	rskWallet.AssertNotCalled(t, "SignBytes", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_BtcTransferHistoryNotConfigured(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))
	nowUnix := time.Now().Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  nil,
		LastRbtcToColdWalletTransfer: &nowUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	eventBus := new(mocks.EventBusMock)
	hashMock := &mocks.HashMock{}

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "no transfer history configured")

	eventBus.AssertNotCalled(t, "Publish")
	hashMock.AssertNotCalled(t, "Hash", mock.Anything)
	rskWallet.AssertNotCalled(t, "SignBytes", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_RskTransferHistoryNotConfigured(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))
	nowUnix := time.Now().Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &nowUnix,
		LastRbtcToColdWalletTransfer: nil,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	eventBus := new(mocks.EventBusMock)
	hashMock := &mocks.HashMock{}

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "no transfer history configured")

	eventBus.AssertNotCalled(t, "Publish")
	hashMock.AssertNotCalled(t, "Hash", mock.Anything)
	rskWallet.AssertNotCalled(t, "SignBytes", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_GetStateConfigurationFails(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	// Provider returns empty struct on validation failure (simulates validation error)
	emptyStateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  nil,
		LastRbtcToColdWalletTransfer: nil,
	}
	generalProvider.On("StateConfiguration", ctx).Return(emptyStateConfig)

	eventBus := new(mocks.EventBusMock)
	hashMock := &mocks.HashMock{}

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "no transfer history configured")

	eventBus.AssertNotCalled(t, "Publish")
	hashMock.AssertNotCalled(t, "Hash", mock.Anything)
	rskWallet.AssertNotCalled(t, "SignBytes", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
}

// nolint: errcheck, funlen
func TestTransferExcessToColdWalletUseCase_Run_BtcExcessNotEconomical(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))

	targetPerNetwork, _ := new(entities.Wei).Div(maxLiquidity, entities.NewWei(2))

	// Liquidity: 20.0006 BTC (below 24 BTC threshold but above 20 BTC target)
	// Time forcing will trigger the excess calculation even though below threshold
	smallExcess := new(entities.Wei).Mul(entities.NewWei(6), entities.NewWei(100000000000000))
	btcLiquidity := new(entities.Wei).Add(targetPerNetwork, smallExcess)
	rbtcLiquidity := new(entities.Wei).Mul(entities.NewWei(testLiquidityAmountWithoutExcess), entities.NewWei(oneEtherInWei))

	// Fee is calculated to fail economic check: excess (0.0006) < fee × multiplier (5)
	btcFee, _ := new(entities.Wei).Div(smallExcess, entities.NewWei(testBtcMinTransferFeeMultiplier-1))
	btcFee = new(entities.Wei).Add(btcFee, entities.NewWei(1))

	oldTransferTime := time.Now().Add(-time.Duration(testForceTransferAfterSeconds+3600) * time.Second)
	oldTransferUnix := oldTransferTime.Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &oldTransferUnix,
		LastRbtcToColdWalletTransfer: &oldTransferUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(btcLiquidity, nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(rbtcLiquidity, nil)

	btcWallet.On("EstimateTxFees", "cold_btc_address", smallExcess).Return(blockchain.BtcFeeEstimation{
		Value: btcFee,
	}, nil)

	eventBus := new(mocks.EventBusMock)
	hashMock := &mocks.HashMock{}

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)

	// BTC transfer should be skipped as not economical (time-forced, but excess is too small relative to fee)
	assert.Equal(t, liquidity_provider.TransferStatusSkippedNotEconomical, result.BtcResult.Status)
	assert.Empty(t, result.BtcResult.TxHash)
	assert.Nil(t, result.BtcResult.Amount)
	assert.Nil(t, result.BtcResult.Fee)
	require.NoError(t, result.BtcResult.Error)
	assert.Contains(t, result.BtcResult.Message, "not economical")

	// RSK has no excess, so it should be skipped normally
	assert.Equal(t, liquidity_provider.TransferStatusSkippedNoExcess, result.RskResult.Status)
	assert.Empty(t, result.RskResult.TxHash)
	assert.Nil(t, result.RskResult.Amount)
	assert.Nil(t, result.RskResult.Fee)
	require.NoError(t, result.RskResult.Error)

	eventBus.AssertNotCalled(t, "Publish")
	hashMock.AssertNotCalled(t, "Hash", mock.Anything)
	rskWallet.AssertNotCalled(t, "SignBytes", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)

	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	btcWallet.AssertExpectations(t)
}

// nolint:errcheck, funlen
func TestTransferExcessToColdWalletUseCase_Run_RbtcExcessNotEconomical(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))

	targetPerNetwork, _ := new(entities.Wei).Div(maxLiquidity, entities.NewWei(2))

	// Liquidity: 20.0006 RBTC (below 24 RBTC threshold but above 20 RBTC target)
	// Time forcing will trigger the excess calculation even though below threshold
	smallExcess := new(entities.Wei).Mul(entities.NewWei(6), entities.NewWei(100000000000000))
	rbtcLiquidity := new(entities.Wei).Add(targetPerNetwork, smallExcess)
	btcLiquidity := new(entities.Wei).Mul(entities.NewWei(testLiquidityAmountWithoutExcess), entities.NewWei(oneEtherInWei))

	// Gas price calculated to fail economic check: (excess - gasCost) < gasCost × multiplier (100)
	rbtcGasPrice, _ := new(entities.Wei).Div(smallExcess, entities.NewWei(liquidity_provider.SimpleTransferGasLimit*testRbtcMinTransferFeeMultiplier+liquidity_provider.SimpleTransferGasLimit))
	rbtcGasPrice = new(entities.Wei).Add(rbtcGasPrice, entities.NewWei(1))

	oldTransferTime := time.Now().Add(-time.Duration(testForceTransferAfterSeconds+3600) * time.Second)
	oldTransferUnix := oldTransferTime.Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &oldTransferUnix,
		LastRbtcToColdWalletTransfer: &oldTransferUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(btcLiquidity, nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(rbtcLiquidity, nil)

	rskRpcMock.On("GasPrice", ctx).Return(rbtcGasPrice, nil)

	eventBus := new(mocks.EventBusMock)
	hashMock := &mocks.HashMock{}

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)

	// BTC has no excess, so it should be skipped normally
	assert.Equal(t, liquidity_provider.TransferStatusSkippedNoExcess, result.BtcResult.Status)
	assert.Empty(t, result.BtcResult.TxHash)
	assert.Nil(t, result.BtcResult.Amount)
	assert.Nil(t, result.BtcResult.Fee)
	require.NoError(t, result.BtcResult.Error)

	assert.Equal(t, liquidity_provider.TransferStatusSkippedNotEconomical, result.RskResult.Status)
	assert.Empty(t, result.RskResult.TxHash)
	assert.Nil(t, result.RskResult.Amount)
	assert.Nil(t, result.RskResult.Fee)
	require.NoError(t, result.RskResult.Error)
	assert.Contains(t, result.RskResult.Message, "not economical")

	eventBus.AssertNotCalled(t, "Publish")
	hashMock.AssertNotCalled(t, "Hash", mock.Anything)
	rskWallet.AssertNotCalled(t, "SignBytes", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)

	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	rskRpcMock.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_GetBtcLiquidityFails(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))
	nowUnix := time.Now().Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &nowUnix,
		LastRbtcToColdWalletTransfer: &nowUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	// Pegout provider (BTC liquidity) returns an error
	expectedError := errors.New("btc wallet connection failed")
	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return((*entities.Wei)(nil), expectedError)

	eventBus := new(mocks.EventBusMock)
	hashMock := &mocks.HashMock{}

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), expectedError.Error())

	eventBus.AssertNotCalled(t, "Publish")
	hashMock.AssertNotCalled(t, "Hash", mock.Anything)
	rskWallet.AssertNotCalled(t, "SignBytes", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_GetRbtcLiquidityFails(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))
	btcLiquidity := new(entities.Wei).Mul(entities.NewWei(testLiquidityAmountWithoutExcess), entities.NewWei(oneEtherInWei))
	nowUnix := time.Now().Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &nowUnix,
		LastRbtcToColdWalletTransfer: &nowUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	// BTC liquidity succeeds
	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(btcLiquidity, nil)

	// Pegin provider (RBTC liquidity) returns an error
	expectedError := errors.New("rsk rpc connection failed")
	peginProvider.On("AvailablePeginLiquidity", ctx).Return((*entities.Wei)(nil), expectedError)

	eventBus := new(mocks.EventBusMock)
	hashMock := &mocks.HashMock{}

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), expectedError.Error())

	eventBus.AssertNotCalled(t, "Publish")
	hashMock.AssertNotCalled(t, "Hash", mock.Anything)
	rskWallet.AssertNotCalled(t, "SignBytes", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)

	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_BtcFeeEstimationFails(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))
	btcLiquidityBig, _ := new(big.Int).SetString(testLiquidityAmountWithExcess, 10)
	btcLiquidity := entities.NewBigWei(btcLiquidityBig)
	rbtcLiquidity := new(entities.Wei).Mul(entities.NewWei(testLiquidityAmountWithoutExcess), entities.NewWei(oneEtherInWei))

	btcExcessBig, _ := new(big.Int).SetString(testExcessAmount, 10)
	btcExcess := entities.NewBigWei(btcExcessBig)

	nowUnix := time.Now().Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &nowUnix,
		LastRbtcToColdWalletTransfer: &nowUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(btcLiquidity, nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(rbtcLiquidity, nil)

	// Fee estimation fails
	expectedError := errors.New("fee estimation service unavailable")
	btcWallet.On("EstimateTxFees", "cold_btc_address", btcExcess).Return(blockchain.BtcFeeEstimation{}, expectedError)

	eventBus := new(mocks.EventBusMock)
	hashMock := &mocks.HashMock{}

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.NotNil(t, result)
	assert.Contains(t, err.Error(), expectedError.Error())

	// BTC transfer should have failed
	assert.Equal(t, liquidity_provider.TransferStatusFailed, result.BtcResult.Status)
	assert.Empty(t, result.BtcResult.TxHash)
	assert.Nil(t, result.BtcResult.Amount)
	assert.Nil(t, result.BtcResult.Fee)
	require.Error(t, result.BtcResult.Error)
	assert.Contains(t, result.BtcResult.Error.Error(), expectedError.Error())

	eventBus.AssertNotCalled(t, "Publish")
	hashMock.AssertNotCalled(t, "Hash", mock.Anything)
	rskWallet.AssertNotCalled(t, "SignBytes", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)

	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	btcWallet.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_BtcTransferFails(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))
	btcLiquidityBig, _ := new(big.Int).SetString(testLiquidityAmountWithExcess, 10)
	btcLiquidity := entities.NewBigWei(btcLiquidityBig)
	rbtcLiquidity := new(entities.Wei).Mul(entities.NewWei(testLiquidityAmountWithoutExcess), entities.NewWei(oneEtherInWei))

	btcExcessBig, _ := new(big.Int).SetString(testExcessAmount, 10)
	btcExcess := entities.NewBigWei(btcExcessBig)
	btcFee := entities.NewWei(testBtcFeeAmount)

	nowUnix := time.Now().Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &nowUnix,
		LastRbtcToColdWalletTransfer: &nowUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(btcLiquidity, nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(rbtcLiquidity, nil)

	btcWallet.On("EstimateTxFees", "cold_btc_address", btcExcess).Return(blockchain.BtcFeeEstimation{
		Value: btcFee,
	}, nil)

	// BTC wallet send fails
	expectedError := errors.New("insufficient funds in wallet")
	btcWallet.On("Send", "cold_btc_address", btcExcess).Return(blockchain.BitcoinTransactionResult{}, expectedError)

	eventBus := new(mocks.EventBusMock)
	hashMock := &mocks.HashMock{}

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.NotNil(t, result)
	assert.Contains(t, err.Error(), expectedError.Error())

	// BTC transfer should have failed
	assert.Equal(t, liquidity_provider.TransferStatusFailed, result.BtcResult.Status)
	assert.Empty(t, result.BtcResult.TxHash)
	assert.Nil(t, result.BtcResult.Amount)
	assert.Nil(t, result.BtcResult.Fee)
	require.Error(t, result.BtcResult.Error)
	assert.Contains(t, result.BtcResult.Error.Error(), expectedError.Error())

	eventBus.AssertNotCalled(t, "Publish")
	hashMock.AssertNotCalled(t, "Hash", mock.Anything)
	rskWallet.AssertNotCalled(t, "SignBytes", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)

	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	btcWallet.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_RskGasPriceRetrievalFails(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))
	btcLiquidity := new(entities.Wei).Mul(entities.NewWei(testLiquidityAmountWithoutExcess), entities.NewWei(oneEtherInWei))
	rbtcLiquidityBig, _ := new(big.Int).SetString(testLiquidityAmountWithExcess, 10)
	rbtcLiquidity := entities.NewBigWei(rbtcLiquidityBig)

	nowUnix := time.Now().Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &nowUnix,
		LastRbtcToColdWalletTransfer: &nowUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(btcLiquidity, nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(rbtcLiquidity, nil)

	// Gas price retrieval fails
	expectedError := errors.New("rpc connection timeout")
	rskRpcMock.On("GasPrice", ctx).Return((*entities.Wei)(nil), expectedError)

	eventBus := new(mocks.EventBusMock)
	hashMock := &mocks.HashMock{}

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.NotNil(t, result)
	assert.Contains(t, err.Error(), expectedError.Error())

	// BTC has no excess, should be skipped
	assert.Equal(t, liquidity_provider.TransferStatusSkippedNoExcess, result.BtcResult.Status)

	// RSK transfer should have failed
	assert.Equal(t, liquidity_provider.TransferStatusFailed, result.RskResult.Status)
	assert.Empty(t, result.RskResult.TxHash)
	assert.Nil(t, result.RskResult.Amount)
	assert.Nil(t, result.RskResult.Fee)
	require.Error(t, result.RskResult.Error)
	assert.Contains(t, result.RskResult.Error.Error(), expectedError.Error())

	eventBus.AssertNotCalled(t, "Publish")
	hashMock.AssertNotCalled(t, "Hash", mock.Anything)
	rskWallet.AssertNotCalled(t, "SignBytes", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)

	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	rskRpcMock.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_RbtcTransferFails(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))
	btcLiquidity := new(entities.Wei).Mul(entities.NewWei(testLiquidityAmountWithoutExcess), entities.NewWei(oneEtherInWei))
	rbtcLiquidityBig, _ := new(big.Int).SetString(testLiquidityAmountWithExcess, 10)
	rbtcLiquidity := entities.NewBigWei(rbtcLiquidityBig)

	rbtcExcessBig, _ := new(big.Int).SetString(testExcessAmount, 10)
	rbtcExcess := entities.NewBigWei(rbtcExcessBig)
	rbtcGasPrice := entities.NewWei(1000000000)
	rbtcGasCost := new(entities.Wei).Mul(entities.NewWei(liquidity_provider.SimpleTransferGasLimit), rbtcGasPrice)
	rbtcAmountToTransfer := new(entities.Wei).Sub(rbtcExcess, rbtcGasCost)

	nowUnix := time.Now().Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &nowUnix,
		LastRbtcToColdWalletTransfer: &nowUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(btcLiquidity, nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(rbtcLiquidity, nil)

	rskRpcMock.On("GasPrice", ctx).Return(rbtcGasPrice, nil)

	// RSK wallet send fails
	expectedError := errors.New("nonce too low")
	rskWallet.On("SendRbtc", ctx, blockchain.NewTransactionConfig(rbtcAmountToTransfer, liquidity_provider.SimpleTransferGasLimit, rbtcGasPrice), "cold_rsk_address").Return(blockchain.TransactionReceipt{}, expectedError)

	eventBus := new(mocks.EventBusMock)
	hashMock := &mocks.HashMock{}

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.NotNil(t, result)
	assert.Contains(t, err.Error(), expectedError.Error())

	// BTC has no excess, should be skipped
	assert.Equal(t, liquidity_provider.TransferStatusSkippedNoExcess, result.BtcResult.Status)

	// RSK transfer should have failed
	assert.Equal(t, liquidity_provider.TransferStatusFailed, result.RskResult.Status)
	assert.Empty(t, result.RskResult.TxHash)
	assert.Nil(t, result.RskResult.Amount)
	assert.Nil(t, result.RskResult.Fee)
	require.Error(t, result.RskResult.Error)
	assert.Contains(t, result.RskResult.Error.Error(), expectedError.Error())

	eventBus.AssertNotCalled(t, "Publish")
	hashMock.AssertNotCalled(t, "Hash", mock.Anything)
	rskWallet.AssertNotCalled(t, "SignBytes", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)

	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	rskRpcMock.AssertExpectations(t)
	rskWallet.AssertExpectations(t)
}

// nolint:funlen
func TestTransferExcessToColdWalletUseCase_Run_BtcSucceedsRskFails(t *testing.T) {
	ctx := context.Background()

	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)
	generalProvider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	coldWallet := new(mocks.ColdWalletMock)
	btcWallet := new(mocks.BitcoinWalletMock)
	rskWallet := new(mocks.RskWalletMock)
	rskRpcMock := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{
		Rsk: rskRpcMock,
	}
	btcWalletMutex := &sync.Mutex{}
	rskWalletMutex := &sync.Mutex{}

	maxLiquidity := new(entities.Wei).Mul(entities.NewWei(testMaxLiquidityBtc), entities.NewWei(oneEtherInWei))
	btcLiquidityBig, _ := new(big.Int).SetString(testLiquidityAmountWithExcess, 10)
	btcLiquidity := entities.NewBigWei(btcLiquidityBig)
	rbtcLiquidityBig, _ := new(big.Int).SetString(testLiquidityAmountWithExcess, 10)
	rbtcLiquidity := entities.NewBigWei(rbtcLiquidityBig)

	btcExcessBig, _ := new(big.Int).SetString(testExcessAmount, 10)
	btcExcess := entities.NewBigWei(btcExcessBig)
	btcFee := entities.NewWei(testBtcFeeAmount)
	btcTxHash := "btc_tx_hash_123"

	rbtcExcessBig, _ := new(big.Int).SetString(testExcessAmount, 10)
	rbtcExcess := entities.NewBigWei(rbtcExcessBig)
	rbtcGasPrice := entities.NewWei(1000000000)
	rbtcGasCost := new(entities.Wei).Mul(entities.NewWei(liquidity_provider.SimpleTransferGasLimit), rbtcGasPrice)
	rbtcAmountToTransfer := new(entities.Wei).Sub(rbtcExcess, rbtcGasCost)

	nowUnix := time.Now().Unix()

	coldWallet.On("GetBtcAddress").Return("cold_btc_address")
	coldWallet.On("GetRskAddress").Return("cold_rsk_address")

	generalConfig := lpEntity.GeneralConfiguration{
		MaxLiquidity: maxLiquidity,
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(testExcessTolerancePercent),
			FixedValue:      entities.NewWei(0),
		},
	}
	generalProvider.On("GeneralConfiguration", ctx).Return(generalConfig)

	stateConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  &nowUnix,
		LastRbtcToColdWalletTransfer: &nowUnix,
	}
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig)

	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(btcLiquidity, nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(rbtcLiquidity, nil)

	// BTC transfer succeeds
	btcWallet.On("EstimateTxFees", "cold_btc_address", btcExcess).Return(blockchain.BtcFeeEstimation{
		Value: btcFee,
	}, nil)
	btcWallet.On("Send", "cold_btc_address", btcExcess).Return(blockchain.BitcoinTransactionResult{
		Hash: btcTxHash,
		Fee:  btcFee,
	}, nil)

	// RSK transfer fails
	rskRpcMock.On("GasPrice", ctx).Return(rbtcGasPrice, nil)
	expectedError := errors.New("transaction underpriced")
	rskWallet.On("SendRbtc", ctx, blockchain.NewTransactionConfig(rbtcAmountToTransfer, liquidity_provider.SimpleTransferGasLimit, rbtcGasPrice), "cold_rsk_address").Return(blockchain.TransactionReceipt{}, expectedError)

	eventBus := new(mocks.EventBusMock)
	eventBus.On("Publish", mock.Anything).Once()
	hashMock := &mocks.HashMock{}

	useCase := liquidity_provider.NewTransferExcessToColdWalletUseCase(
		peginProvider,
		pegoutProvider,
		generalProvider,
		lpRepository,
		coldWallet,
		btcWallet,
		rskWallet,
		rpc,
		btcWalletMutex,
		rskWalletMutex,
		testBtcMinTransferFeeMultiplier,
		testRbtcMinTransferFeeMultiplier,
		testForceTransferAfterSeconds,
		eventBus,
		rskWallet,
		hashMock.Hash,
	)

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.NotNil(t, result)
	assert.Contains(t, err.Error(), expectedError.Error())

	// BTC transfer should have succeeded
	assert.Equal(t, liquidity_provider.TransferStatusSuccess, result.BtcResult.Status)
	assert.Equal(t, btcTxHash, result.BtcResult.TxHash)
	assert.Equal(t, btcExcess, result.BtcResult.Amount)
	assert.Equal(t, btcFee, result.BtcResult.Fee)
	require.NoError(t, result.BtcResult.Error)

	// RSK transfer should have failed
	assert.Equal(t, liquidity_provider.TransferStatusFailed, result.RskResult.Status)
	assert.Empty(t, result.RskResult.TxHash)
	assert.Nil(t, result.RskResult.Amount)
	assert.Nil(t, result.RskResult.Fee)
	require.Error(t, result.RskResult.Error)
	assert.Contains(t, result.RskResult.Error.Error(), expectedError.Error())

	// Persist is not called when Run returns error (RSK failed)
	hashMock.AssertNotCalled(t, "Hash", mock.Anything)
	rskWallet.AssertNotCalled(t, "SignBytes", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)

	coldWallet.AssertExpectations(t)
	generalProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	btcWallet.AssertExpectations(t)
	rskRpcMock.AssertExpectations(t)
	rskWallet.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}
