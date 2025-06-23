package watcher_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"sync"
	"testing"
	"time"
)

// nolint:funlen
func TestEclipseCheckUseCase_Run_Rootstock(t *testing.T) {
	const (
		recipient = "alertRecipient"
		blocks    = 500
		nonce     = 20
	)
	var nilBigInt *big.Int
	var timestamp = time.Now()
	config := w.EclipseCheckConfig{
		RskToleranceThreshold:    50,
		RskMaxMsWaitForBlock:     1000,
		RskWaitPollingMsInterval: 100,
		BtcToleranceThreshold:    0,
		BtcMaxMsWaitForBlock:     0,
		BtcWaitPollingMsInterval: 0,
	}
	t.Run("should not trigger the alert if no eclipse attack is detected", func(t *testing.T) {
		rskRpc := &mocks.RootstockRpcServerMock{}
		btcRpc := &mocks.BtcRpcMock{}
		eventBus := &mocks.EventBusMock{}
		alertSender := &mocks.AlertSenderMock{}
		rskExtra1 := &mocks.RootstockRpcServerMock{}
		rskExtra2 := &mocks.RootstockRpcServerMock{}
		rskExtras := []*mocks.RootstockRpcServerMock{rskExtra1, rskExtra2}
		rskRpc.EXPECT().GetBlockByNumber(mock.Anything, nilBigInt).Return(blockchain.BlockInfo{
			Hash:      test.AnyHash,
			Number:    blocks,
			Timestamp: timestamp,
			Nonce:     nonce,
		}, nil)
		// 50% should be enough to not trigger the alert
		rskExtras[0].EXPECT().GetBlockByNumber(mock.Anything, nilBigInt).Return(blockchain.BlockInfo{
			Hash:      test.AnyHash,
			Number:    blocks,
			Timestamp: timestamp,
			Nonce:     nonce,
		}, nil)
		rskExtras[1].EXPECT().GetBlockByNumber(mock.Anything, nilBigInt).Return(blockchain.BlockInfo{
			Hash:      "diffHash",
			Number:    300,
			Timestamp: timestamp,
			Nonce:     nonce,
		}, nil)
		mutex := &mocks.MutexMock{}
		mutex.On("Lock").Return()
		mutex.On("Unlock").Return()
		useCase := w.NewEclipseCheckUseCase(
			config,
			blockchain.Rpc{
				Btc: btcRpc,
				Rsk: rskRpc,
			},
			[]blockchain.BitcoinNetwork{},
			[]blockchain.RootstockRpcServer{rskExtra1, rskExtra2},
			eventBus,
			alertSender,
			recipient,
			mutex,
		)
		err := useCase.Run(context.Background(), entities.NodeTypeRootstock)
		require.NoError(t, err)
		rskRpc.AssertExpectations(t)
		for _, rskExtra := range rskExtras {
			rskExtra.AssertExpectations(t)
		}
		eventBus.AssertNotCalled(t, "Publish")
		alertSender.AssertNotCalled(t, "SendAlert")
		mutex.AssertExpectations(t)
	})
	t.Run("should trigger the alert if eclipse attack is detected", func(t *testing.T) {
		rskRpc := &mocks.RootstockRpcServerMock{}
		btcRpc := &mocks.BtcRpcMock{}
		eventBus := &mocks.EventBusMock{}
		alertSender := &mocks.AlertSenderMock{}
		rskExtra1 := &mocks.RootstockRpcServerMock{}
		rskExtra2 := &mocks.RootstockRpcServerMock{}
		rskExtras := []*mocks.RootstockRpcServerMock{rskExtra1, rskExtra2}
		rskRpc.EXPECT().GetBlockByNumber(mock.Anything, nilBigInt).Return(blockchain.BlockInfo{
			Hash:      test.AnyHash,
			Number:    blocks,
			Timestamp: timestamp,
			Nonce:     nonce,
		}, nil)
		for _, rskExtra := range rskExtras {
			rskExtra.EXPECT().GetBlockByNumber(mock.Anything, nilBigInt).Return(blockchain.BlockInfo{
				Hash:      "diffHash",
				Number:    700,
				Timestamp: timestamp,
				Nonce:     nonce,
			}, nil)
		}
		mutex := &mocks.MutexMock{}
		mutex.On("Lock").Return()
		mutex.On("Unlock").Return()
		alertSender.On(
			"SendAlert",
			mock.Anything,
			w.EclipseAlertSubject,
			"Your rootstock node is under eclipse attack. Please, check your node's connectivity and synchronization.",
			[]string{recipient},
		).Return(nil)
		eventBus.On("Publish", mock.MatchedBy(func(e blockchain.NodeEclipseEvent) bool {
			return assert.Equal(t, entities.NodeTypeRootstock, e.NodeType) &&
				assert.Equal(t, uint64(blocks), e.EclipsedBlockNumber) &&
				assert.Equal(t, test.AnyHash, e.EclipsedBlockHash)
		})).Return()
		useCase := w.NewEclipseCheckUseCase(
			config,
			blockchain.Rpc{
				Btc: btcRpc,
				Rsk: rskRpc,
			},
			[]blockchain.BitcoinNetwork{},
			[]blockchain.RootstockRpcServer{rskExtras[0], rskExtras[1]},
			eventBus,
			alertSender,
			recipient,
			mutex,
		)
		err := useCase.Run(context.Background(), entities.NodeTypeRootstock)
		require.ErrorIs(t, err, w.NodeEclipseDetectedError)
		rskRpc.AssertExpectations(t)
		for _, rskExtra := range rskExtras {
			rskExtra.AssertExpectations(t)
		}
		eventBus.AssertExpectations(t)
		alertSender.AssertExpectations(t)
		mutex.AssertExpectations(t)
	})
	t.Run("should not trigger alert if the node syncs during the tolerance threshold", func(t *testing.T) {
		rskRpc := &mocks.RootstockRpcServerMock{}
		btcRpc := &mocks.BtcRpcMock{}
		eventBus := &mocks.EventBusMock{}
		alertSender := &mocks.AlertSenderMock{}
		rskExtra1 := &mocks.RootstockRpcServerMock{}
		rskExtra2 := &mocks.RootstockRpcServerMock{}
		rskExtras := []*mocks.RootstockRpcServerMock{rskExtra1, rskExtra2}
		rskRpc.EXPECT().GetBlockByNumber(mock.Anything, nilBigInt).Return(blockchain.BlockInfo{
			Hash:      "otherHash",
			Number:    300,
			Timestamp: timestamp,
			Nonce:     nonce,
		}, nil).Once()
		rskRpc.EXPECT().GetBlockByNumber(mock.Anything, nilBigInt).Return(blockchain.BlockInfo{
			Hash:      test.AnyHash,
			Number:    blocks,
			Timestamp: timestamp,
			Nonce:     nonce,
		}, nil).Once()
		for _, rskExtra := range rskExtras {
			rskExtra.EXPECT().GetBlockByNumber(mock.Anything, nilBigInt).Return(blockchain.BlockInfo{
				Hash:      test.AnyHash,
				Number:    blocks,
				Timestamp: timestamp,
				Nonce:     nonce,
			}, nil)
			rskExtra.EXPECT().GetBlockByNumber(mock.Anything, nilBigInt).Return(blockchain.BlockInfo{
				Hash:      test.AnyHash,
				Number:    blocks,
				Timestamp: timestamp,
				Nonce:     nonce,
			}, nil)
		}
		mutex := &mocks.MutexMock{}
		mutex.On("Lock").Return()
		mutex.On("Unlock").Return()
		useCase := w.NewEclipseCheckUseCase(
			config,
			blockchain.Rpc{
				Btc: btcRpc,
				Rsk: rskRpc,
			},
			[]blockchain.BitcoinNetwork{},
			[]blockchain.RootstockRpcServer{rskExtra1, rskExtra2},
			eventBus,
			alertSender,
			recipient,
			mutex,
		)
		err := useCase.Run(context.Background(), entities.NodeTypeRootstock)
		require.NoError(t, err)
		rskRpc.AssertExpectations(t)
		for _, rskExtra := range rskExtras {
			rskExtra.AssertExpectations(t)
		}
		eventBus.AssertNotCalled(t, "Publish")
		alertSender.AssertNotCalled(t, "SendAlert")
		mutex.AssertExpectations(t)
	})
}

// nolint:funlen
func TestEclipseCheckUseCase_Run_Bitcoin(t *testing.T) {
	const (
		recipient   = "alertRecipient"
		networkName = "testnet"
		blocks      = 500
		headers     = 501
	)
	config := w.EclipseCheckConfig{
		RskToleranceThreshold:    0,
		RskMaxMsWaitForBlock:     0,
		RskWaitPollingMsInterval: 0,
		BtcToleranceThreshold:    50,
		BtcMaxMsWaitForBlock:     1000,
		BtcWaitPollingMsInterval: 100,
	}
	t.Run("should not trigger the alert if no eclipse attack is detected", func(t *testing.T) {
		rskRpc := &mocks.RootstockRpcServerMock{}
		btcRpc := &mocks.BtcRpcMock{}
		btcExtra1 := &mocks.BtcRpcMock{}
		btcExtra2 := &mocks.BtcRpcMock{}
		eventBus := &mocks.EventBusMock{}
		alertSender := &mocks.AlertSenderMock{}
		btcExtras := []*mocks.BtcRpcMock{btcExtra1, btcExtra2}
		btcRpc.On("GetBlockchainInfo").Return(blockchain.BitcoinBlockchainInfo{
			NetworkName:      networkName,
			ValidatedBlocks:  big.NewInt(blocks),
			ValidatedHeaders: big.NewInt(headers),
			BestBlockHash:    test.AnyHash,
		}, nil)
		// 50% should be enough to not trigger the alert
		btcExtras[0].On("GetBlockchainInfo").Return(blockchain.BitcoinBlockchainInfo{
			NetworkName:      networkName,
			ValidatedBlocks:  big.NewInt(blocks),
			ValidatedHeaders: big.NewInt(headers),
			BestBlockHash:    test.AnyHash,
		}, nil)
		btcExtras[1].On("GetBlockchainInfo").Return(blockchain.BitcoinBlockchainInfo{
			NetworkName:      networkName,
			ValidatedBlocks:  big.NewInt(100),
			ValidatedHeaders: big.NewInt(101),
			BestBlockHash:    "otherHash",
		}, nil)
		mutex := &mocks.MutexMock{}
		mutex.On("Lock").Return()
		mutex.On("Unlock").Return()
		useCase := w.NewEclipseCheckUseCase(
			config,
			blockchain.Rpc{
				Btc: btcRpc,
				Rsk: rskRpc,
			},
			[]blockchain.BitcoinNetwork{btcExtras[0], btcExtras[1]},
			[]blockchain.RootstockRpcServer{},
			eventBus,
			alertSender,
			recipient,
			mutex,
		)
		err := useCase.Run(context.Background(), entities.NodeTypeBitcoin)
		require.NoError(t, err)
		btcRpc.AssertExpectations(t)
		for _, btcExtra := range btcExtras {
			btcExtra.AssertExpectations(t)
		}
		eventBus.AssertNotCalled(t, "Publish")
		alertSender.AssertNotCalled(t, "SendAlert")
		mutex.AssertExpectations(t)
	})
	t.Run("should trigger the alert if eclipse attack is detected", func(t *testing.T) {
		rskRpc := &mocks.RootstockRpcServerMock{}
		btcRpc := &mocks.BtcRpcMock{}
		btcExtra1 := &mocks.BtcRpcMock{}
		btcExtra2 := &mocks.BtcRpcMock{}
		eventBus := &mocks.EventBusMock{}
		alertSender := &mocks.AlertSenderMock{}
		btcExtras := []*mocks.BtcRpcMock{btcExtra1, btcExtra2}
		btcRpc.On("GetBlockchainInfo").Return(blockchain.BitcoinBlockchainInfo{
			NetworkName:      networkName,
			ValidatedBlocks:  big.NewInt(blocks),
			ValidatedHeaders: big.NewInt(headers),
			BestBlockHash:    test.AnyHash,
		}, nil)
		for _, btcExtra := range btcExtras {
			btcExtra.On("GetBlockchainInfo").Return(blockchain.BitcoinBlockchainInfo{
				NetworkName:      networkName,
				ValidatedBlocks:  big.NewInt(600),
				ValidatedHeaders: big.NewInt(601),
				BestBlockHash:    "otherHash",
			}, nil)
		}
		mutex := &mocks.MutexMock{}
		mutex.On("Lock").Return()
		mutex.On("Unlock").Return()
		alertSender.On(
			"SendAlert",
			mock.Anything,
			w.EclipseAlertSubject,
			"Your bitcoin node is under eclipse attack. Please, check your node's connectivity and synchronization.",
			[]string{recipient},
		).Return(nil)
		eventBus.On("Publish", mock.MatchedBy(func(e blockchain.NodeEclipseEvent) bool {
			return assert.Equal(t, entities.NodeTypeBitcoin, e.NodeType) &&
				assert.Equal(t, uint64(blocks), e.EclipsedBlockNumber) &&
				assert.Equal(t, test.AnyHash, e.EclipsedBlockHash)
		})).Return()
		useCase := w.NewEclipseCheckUseCase(
			config,
			blockchain.Rpc{
				Btc: btcRpc,
				Rsk: rskRpc,
			},
			[]blockchain.BitcoinNetwork{btcExtras[0], btcExtras[1]},
			[]blockchain.RootstockRpcServer{},
			eventBus,
			alertSender,
			recipient,
			mutex,
		)
		err := useCase.Run(context.Background(), entities.NodeTypeBitcoin)
		require.ErrorIs(t, err, w.NodeEclipseDetectedError)
		btcRpc.AssertExpectations(t)
		for _, btcExtra := range btcExtras {
			btcExtra.AssertExpectations(t)
		}
		eventBus.AssertExpectations(t)
		alertSender.AssertExpectations(t)
		mutex.AssertExpectations(t)
	})
	t.Run("should not trigger alert if the node syncs during the tolerance threshold", func(t *testing.T) {
		rskRpc := &mocks.RootstockRpcServerMock{}
		btcRpc := &mocks.BtcRpcMock{}
		btcExtra1 := &mocks.BtcRpcMock{}
		btcExtra2 := &mocks.BtcRpcMock{}
		eventBus := &mocks.EventBusMock{}
		alertSender := &mocks.AlertSenderMock{}
		btcExtras := []*mocks.BtcRpcMock{btcExtra1, btcExtra2}
		btcRpc.On("GetBlockchainInfo").Return(blockchain.BitcoinBlockchainInfo{
			NetworkName:      networkName,
			ValidatedBlocks:  big.NewInt(300),
			ValidatedHeaders: big.NewInt(300),
			BestBlockHash:    "otherHash",
		}, nil).Once()
		btcRpc.On("GetBlockchainInfo").Return(blockchain.BitcoinBlockchainInfo{
			NetworkName:      networkName,
			ValidatedBlocks:  big.NewInt(blocks),
			ValidatedHeaders: big.NewInt(headers),
			BestBlockHash:    test.AnyHash,
		}, nil).Once()
		for _, btcExtra := range btcExtras {
			btcExtra.On("GetBlockchainInfo").Return(blockchain.BitcoinBlockchainInfo{
				NetworkName:      networkName,
				ValidatedBlocks:  big.NewInt(blocks),
				ValidatedHeaders: big.NewInt(headers),
				BestBlockHash:    test.AnyHash,
			}, nil)
		}
		mutex := &mocks.MutexMock{}
		mutex.On("Lock").Return()
		mutex.On("Unlock").Return()
		useCase := w.NewEclipseCheckUseCase(
			config,
			blockchain.Rpc{
				Btc: btcRpc,
				Rsk: rskRpc,
			},
			[]blockchain.BitcoinNetwork{btcExtras[0], btcExtras[1]},
			[]blockchain.RootstockRpcServer{},
			eventBus,
			alertSender,
			recipient,
			mutex,
		)
		err := useCase.Run(context.Background(), entities.NodeTypeBitcoin)
		require.NoError(t, err)
		btcRpc.AssertExpectations(t)
		for _, btcExtra := range btcExtras {
			btcExtra.AssertExpectations(t)
		}
		eventBus.AssertNotCalled(t, "Publish")
		alertSender.AssertNotCalled(t, "SendAlert")
		mutex.AssertExpectations(t)
	})
}

// nolint:funlen
func TestEclipseCheckUseCase_Run_ErrorCases(t *testing.T) {
	const (
		recipient = "alertRecipient"
	)
	t.Run("should handle unsupported node type", func(t *testing.T) {
		useCase := w.NewEclipseCheckUseCase(
			w.EclipseCheckConfig{},
			blockchain.Rpc{
				Btc: &mocks.BtcRpcMock{},
				Rsk: &mocks.RootstockRpcServerMock{},
			},
			[]blockchain.BitcoinNetwork{},
			[]blockchain.RootstockRpcServer{},
			&mocks.EventBusMock{},
			&mocks.AlertSenderMock{},
			recipient,
			&sync.Mutex{},
		)
		err := useCase.Run(context.Background(), "other type")
		require.ErrorContains(t, err, "unsupported node type: other type")
	})
	t.Run("should handle error getting block from main rsk source", func(t *testing.T) {
		rsk := &mocks.RootstockRpcServerMock{}
		rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.Anything).Return(blockchain.BlockInfo{}, assert.AnError)
		useCase := w.NewEclipseCheckUseCase(
			w.EclipseCheckConfig{},
			blockchain.Rpc{
				Btc: &mocks.BtcRpcMock{},
				Rsk: rsk,
			},
			[]blockchain.BitcoinNetwork{},
			[]blockchain.RootstockRpcServer{},
			&mocks.EventBusMock{},
			&mocks.AlertSenderMock{},
			recipient,
			&sync.Mutex{},
		)
		err := useCase.Run(context.Background(), entities.NodeTypeRootstock)
		require.Error(t, err)
		rsk.AssertExpectations(t)
	})
	t.Run("should handle error getting block from main btc source", func(t *testing.T) {
		btc := &mocks.BtcRpcMock{}
		btc.On("GetBlockchainInfo").Return(blockchain.BitcoinBlockchainInfo{}, assert.AnError)
		useCase := w.NewEclipseCheckUseCase(
			w.EclipseCheckConfig{},
			blockchain.Rpc{
				Btc: btc,
				Rsk: &mocks.RootstockRpcServerMock{},
			},
			[]blockchain.BitcoinNetwork{},
			[]blockchain.RootstockRpcServer{},
			&mocks.EventBusMock{},
			&mocks.AlertSenderMock{},
			recipient,
			&sync.Mutex{},
		)
		err := useCase.Run(context.Background(), entities.NodeTypeBitcoin)
		require.Error(t, err)
		btc.AssertExpectations(t)
	})
	t.Run("should not fail on error getting block from external rsk source", func(t *testing.T) {
		rsk := &mocks.RootstockRpcServerMock{}
		rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.Anything).Return(blockchain.BlockInfo{
			Hash:      test.AnyHash,
			Number:    200,
			Timestamp: time.Now(),
			Nonce:     123,
		}, nil)
		extra := &mocks.RootstockRpcServerMock{}
		extra.EXPECT().GetBlockByNumber(mock.Anything, mock.Anything).Return(blockchain.BlockInfo{}, assert.AnError)
		useCase := w.NewEclipseCheckUseCase(
			w.EclipseCheckConfig{},
			blockchain.Rpc{
				Btc: &mocks.BtcRpcMock{},
				Rsk: rsk,
			},
			[]blockchain.BitcoinNetwork{},
			[]blockchain.RootstockRpcServer{extra},
			&mocks.EventBusMock{},
			&mocks.AlertSenderMock{},
			recipient,
			&sync.Mutex{},
		)
		defer test.AssertLogContains(t, "Error getting block from external RSK source:")()
		err := useCase.Run(context.Background(), entities.NodeTypeRootstock)
		require.NoError(t, err)
		rsk.AssertExpectations(t)
		extra.AssertExpectations(t)
	})
	t.Run("should handle error getting block from external btc source", func(t *testing.T) {
		btc := &mocks.BtcRpcMock{}
		btc.On("GetBlockchainInfo").Return(blockchain.BitcoinBlockchainInfo{
			NetworkName:      "mainnet",
			ValidatedBlocks:  big.NewInt(123),
			ValidatedHeaders: big.NewInt(123),
			BestBlockHash:    test.AnyHash,
		}, nil)
		extra := &mocks.BtcRpcMock{}
		extra.On("GetBlockchainInfo").Return(blockchain.BitcoinBlockchainInfo{}, assert.AnError)
		useCase := w.NewEclipseCheckUseCase(
			w.EclipseCheckConfig{},
			blockchain.Rpc{
				Btc: btc,
				Rsk: &mocks.RootstockRpcServerMock{},
			},
			[]blockchain.BitcoinNetwork{extra},
			[]blockchain.RootstockRpcServer{},
			&mocks.EventBusMock{},
			&mocks.AlertSenderMock{},
			recipient,
			&sync.Mutex{},
		)
		defer test.AssertLogContains(t, "Error getting latest block from external Bitcoin source:")()
		err := useCase.Run(context.Background(), entities.NodeTypeBitcoin)
		require.NoError(t, err)
		btc.AssertExpectations(t)
		extra.AssertExpectations(t)
	})
	t.Run("should handle error sending alert", func(t *testing.T) {
		btc := &mocks.BtcRpcMock{}
		btc.On("GetBlockchainInfo").Return(blockchain.BitcoinBlockchainInfo{
			NetworkName:      "mainnet",
			ValidatedBlocks:  big.NewInt(123),
			ValidatedHeaders: big.NewInt(123),
			BestBlockHash:    test.AnyHash,
		}, nil)
		extra := &mocks.BtcRpcMock{}
		extra.On("GetBlockchainInfo").Return(blockchain.BitcoinBlockchainInfo{}, assert.AnError)
		alertSender := &mocks.AlertSenderMock{}
		alertSender.On("SendAlert", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(assert.AnError)
		eventBus := &mocks.EventBusMock{}
		eventBus.On("Publish", mock.AnythingOfType("blockchain.NodeEclipseEvent")).Return()
		useCase := w.NewEclipseCheckUseCase(
			w.EclipseCheckConfig{
				BtcToleranceThreshold:    10,
				BtcWaitPollingMsInterval: 100,
				BtcMaxMsWaitForBlock:     500,
			},
			blockchain.Rpc{
				Btc: btc,
				Rsk: &mocks.RootstockRpcServerMock{},
			},
			[]blockchain.BitcoinNetwork{extra},
			[]blockchain.RootstockRpcServer{},
			eventBus,
			alertSender,
			recipient,
			&sync.Mutex{},
		)
		defer test.AssertLogContains(t, "Error getting latest block from external Bitcoin source:")()
		err := useCase.Run(context.Background(), entities.NodeTypeBitcoin)
		require.ErrorIs(t, err, assert.AnError)
		btc.AssertExpectations(t)
		extra.AssertExpectations(t)
		alertSender.AssertExpectations(t)
		eventBus.AssertExpectations(t)
	})
	t.Run("should run when no extra data sources are provided", func(t *testing.T) {
		btc := &mocks.BtcRpcMock{}
		btc.On("GetBlockchainInfo").Return(blockchain.BitcoinBlockchainInfo{
			NetworkName:      "mainnet",
			ValidatedBlocks:  big.NewInt(123),
			ValidatedHeaders: big.NewInt(123),
			BestBlockHash:    test.AnyHash,
		}, nil)
		useCase := w.NewEclipseCheckUseCase(
			w.EclipseCheckConfig{},
			blockchain.Rpc{
				Btc: btc,
				Rsk: &mocks.RootstockRpcServerMock{},
			},
			[]blockchain.BitcoinNetwork{},
			[]blockchain.RootstockRpcServer{},
			&mocks.EventBusMock{},
			&mocks.AlertSenderMock{},
			recipient,
			&sync.Mutex{},
		)
		err := useCase.Run(context.Background(), entities.NodeTypeBitcoin)
		require.NoError(t, err)
		btc.AssertExpectations(t)
	})
}
